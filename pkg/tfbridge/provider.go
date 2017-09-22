// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/glog"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	pbstruct "github.com/golang/protobuf/ptypes/struct"
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/flatmap"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/diag"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/plugin"
	"github.com/pulumi/pulumi/pkg/resource/provider"
	"github.com/pulumi/pulumi/pkg/tokens"
	"github.com/pulumi/pulumi/pkg/util/contract"
	lumirpc "github.com/pulumi/pulumi/sdk/proto/go"
	"golang.org/x/net/context"
)

// Provider implements the Lumi resource provider operations for any Terraform plugin.
type Provider struct {
	host      *provider.HostClient       // the RPC link back to the Lumi engine.
	module    string                     // the Terraform module name.
	tf        terraform.ResourceProvider // the Terraform resource provider to use.
	info      ProviderInfo               // overlaid info about this provider.
	resources map[tokens.Type]Resource   // a map of Lumi type tokens to resource info.
}

// Resource wraps both the Terraform resource type info plus the overlay resource info.
type Resource struct {
	TF     terraform.ResourceType // Terraform resource info.
	Schema *ResourceInfo          // optional provider overrides.
}

// NewProvider creates a new Lumi RPC server wired up to the given host and wrapping the given Terraform provider.
func NewProvider(host *provider.HostClient, module string,
	tf terraform.ResourceProvider, info ProviderInfo) *Provider {
	p := &Provider{
		host:   host,
		module: module,
		tf:     tf,
		info:   info,
	}
	p.initResourceMap()
	return p
}

var _ lumirpc.ResourceProviderServer = (*Provider)(nil)

func (p *Provider) pkg() tokens.Package          { return tokens.Package(p.module) }
func (p *Provider) indexMod() tokens.Module      { return tokens.Module(p.pkg() + ":index") }
func (p *Provider) baseConfigMod() tokens.Module { return tokens.Module(p.pkg() + ":config") }
func (p *Provider) configMod() tokens.Module     { return p.baseConfigMod() + "/vars" }

// resource looks up the Terraform resource provider from its Lumi type token.
func (p *Provider) resource(t tokens.Type) (Resource, bool) {
	res, has := p.resources[t]
	return res, has
}

// initResourceMap creates a simple map from Lumi to Terraform resource type.
func (p *Provider) initResourceMap() {
	prefix := p.module + "_" // all resources will have this prefix.

	// Fetch a list of all resource types handled by this provider and make a map.
	p.resources = make(map[tokens.Type]Resource)
	for _, res := range p.tf.Resources() {
		var tok tokens.Type

		// See if there is override information for this resource.  If yes, use that to decode the token.
		var schema *ResourceInfo
		if p.info.Resources != nil {
			schema = p.info.Resources[res.Name]
			if schema != nil {
				tok = schema.Tok
			}
		}

		// Otherwise, we default to the standard naming scheme.
		if tok == "" {
			// Strip off the module prefix (e.g., "aws_").
			contract.Assertf(strings.HasPrefix(res.Name, prefix),
				"Expected all Terraform resources in this module to have a '%v' prefix", prefix)
			name := res.Name[len(prefix):]

			// Create a camel name for the module and pascal for the resource type.
			camelName := TerraformToLumiName(name, false)
			pascalName := TerraformToLumiName(name, true)

			// Now just manufacture a token with the package, module, and resource type name.
			tok = tokens.Type(string(p.pkg()) + ":" + camelName + ":" + pascalName)
		}

		p.resources[tok] = Resource{TF: res, Schema: schema}
	}
}

// getInfoFromTerraformName does a map lookup to find the Lumi name and schema info, if any.
func getInfoFromTerraformName(key string, schema map[string]*SchemaInfo) (resource.PropertyKey, *SchemaInfo) {
	info := schema[key]
	var name string
	if info != nil {
		name = info.Name
	}
	if name == "" {
		// If no name override exists, use the default name mangling scheme.
		name = TerraformToLumiName(key, false)
	}
	return resource.PropertyKey(name), info
}

// getInfoFromLumiName does a reverse map lookup to find the Terraform name and schema info for a Lumi name, if any.
func getInfoFromLumiName(key resource.PropertyKey, schema map[string]*SchemaInfo) (string, *SchemaInfo) {
	// To do this, we will first look to see if there's a known custom schema that uses this name.  If yes, we
	// prefer to use that.  To do this, we must use a reverse lookup.  (In the future we may want to make a
	// lookaside map to avoid the traversal of this map.)  Otherwise, use the standard name mangling scheme.
	ks := string(key)
	for tfname, schinfo := range schema {
		if schinfo != nil && schinfo.Name == ks {
			return tfname, schinfo
		}
	}
	return LumiToTerraformName(ks), schema[ks]
}

// makeTerraformInputs takes a property map plus custom schema info and does whatever is necessary to prepare it for
// use by Terraform.  Note that this function may have side effects, for instance if it is necessary to spill an asset
// to disk in order to create a name out of it.  Please take care not to call it superfluously!
func (p *Provider) makeTerraformInputs(res *LumiResource, m resource.PropertyMap,
	schema map[string]*SchemaInfo, defaults bool) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Enumerate the inputs provided and add them to the map using their Terraform names.
	for key, value := range m {
		// First translate the Lumi property name to a Terraform name.
		name, info := getInfoFromLumiName(key, schema)
		contract.Assert(name != "")

		// And then translate the property value.
		v, err := p.makeTerraformInput(res, name, value, info, defaults)
		if err != nil {
			return nil, err
		}
		result[name] = v
	}

	// Now enumerate and propagate defaults if the corresponding values are still missing.
	for key, info := range schema {
		if v, has := result[key]; has {
			glog.V(9).Infof("Created Terraform input: %v = %v", key, v)
		} else if defaults && info.HasDefault() {
			if info.Default.Value != nil {
				result[key] = info.Default.Value
				glog.V(9).Infof("Created Terraform input: %v = %v (default)", key, result[key])
			} else if from := info.Default.From; from != nil {
				result[key] = from(res)
				glog.V(9).Infof("Created Terraform input: %v = %v (default from fnc)", key, result[key])
			} else {
				contract.Failf("Default missing Value or From")
			}
		} else {
			glog.V(9).Infof("Skipped Terraform input: %v (skipped or no defaults)", key)
		}
	}

	if glog.V(5) {
		for k, v := range result {
			glog.V(5).Infof("Terraform input %v = %v", k, v)
		}
	}

	return result, nil
}

// makeTerraformInput takes a single property plus custom schema info and does whatever is necessary to prepare it for
// use by Terraform.  Note that this function may have side effects, for instance if it is necessary to spill an asset
// to disk in order to create a name out of it.  Please take care not to call it superfluously!
func (p *Provider) makeTerraformInput(res *LumiResource, name string,
	v resource.PropertyValue, schema *SchemaInfo, defaults bool) (interface{}, error) {
	if v.IsNull() {
		return nil, nil
	} else if v.IsBool() {
		return v.BoolValue(), nil
	} else if v.IsNumber() {
		return int(v.NumberValue()), nil // convert floats to ints.
	} else if v.IsString() {
		return v.StringValue(), nil
	} else if v.IsArray() {
		// FIXME: marshal/unmarshal sets properly.
		var arr []interface{}
		for i, elem := range v.ArrayValue() {
			var elemschema *SchemaInfo
			if schema != nil {
				elemschema = schema.Elem
			}
			e, err := p.makeTerraformInput(res, fmt.Sprintf("%v[%v]", name, i), elem, elemschema, defaults)
			if err != nil {
				return nil, err
			}
			arr = append(arr, e)
		}
		return arr, nil
	} else if v.IsAsset() {
		// We require that there be asset information, otherwise an error occurs.
		if schema == nil || schema.Asset == nil {
			return nil,
				errors.Errorf("Encountered an asset %v but asset translation instructions were missing", name)
		} else if !schema.Asset.IsAsset() {
			return nil,
				errors.Errorf("Invalid asset translation instructions for %v; expected an asset", name)
		}
		return schema.Asset.TranslateAsset(v.AssetValue())
	} else if v.IsArchive() {
		// We require that there be archive information, otherwise an error occurs.
		if schema == nil || schema.Asset == nil {
			return nil,
				errors.Errorf("Encountered an archive %v but asset translation instructions were missing", name)
		} else if !schema.Asset.IsArchive() {
			return nil,
				errors.Errorf("Invalid asset translation instructions for %v; expected an archive", name)
		}
		return schema.Asset.TranslateArchive(v.ArchiveValue())
	} else if v.IsObject() {
		var fldschemas map[string]*SchemaInfo
		if schema != nil {
			fldschemas = schema.Fields
		}
		return p.makeTerraformInputs(res, v.ObjectValue(), fldschemas, defaults)
	} else if v.IsComputed() || v.IsOutput() {
		// If any variables are unknown, we need to mark them in the inputs so the config map treats it right.  This
		// requires the use of the special UnknownVariableValue sentinel in Terraform, which is how it internally stores
		// interpolated variables whose inputs are currently unknown.
		return config.UnknownVariableValue, nil
	}

	contract.Failf("Unexpected value marshaled: %v", v)
	return nil, nil
}

// makeTerraformInputsFromRPC unmarshals an RPC payload of properties and turns the results into Terraform inputs.
func (p *Provider) makeTerraformInputsFromRPC(res *LumiResource, m *pbstruct.Struct,
	schema map[string]*SchemaInfo, allowUnknowns bool, defaults bool) (map[string]interface{}, error) {
	props, err := plugin.UnmarshalProperties(m,
		plugin.MarshalOptions{AllowUnknowns: allowUnknowns, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	return p.makeTerraformInputs(res, props, schema, defaults)
}

// makeTerraformResult expands a Terraform-style flatmap into an expanded Lumi resource property map.  This respects
// the property maps so that results end up with their correct Lumi names when shipping back to the engine.
func (p *Provider) makeTerraformResult(props map[string]string,
	schema map[string]*SchemaInfo) resource.PropertyMap {
	outs := make(map[string]interface{})
	for _, key := range flatmap.Map(props).Keys() {
		outs[key] = flatmap.Expand(props, key)
	}
	return p.makeTerraformOutputs(outs, schema)
}

// makeTerraformResultValue expands a single Terraform-style flatmap entry into a resource property value.
func (p *Provider) makeTerraformResultValue(props map[string]string,
	key string, schema *SchemaInfo) resource.PropertyValue {
	v := flatmap.Expand(props, key)
	return p.makeTerraformOutput(v, schema)
}

// makeTerraformOutputs takes an expanded Terraform property map and returns a Lumi equivalent.  This respects
// the property maps so that results end up with their correct Lumi names when shipping back to the engine.
func (p *Provider) makeTerraformOutputs(outs map[string]interface{},
	schema map[string]*SchemaInfo) resource.PropertyMap {
	result := make(resource.PropertyMap)
	for key, value := range outs {
		// First do a lookup of the name/info.
		name, info := getInfoFromTerraformName(key, schema)
		contract.Assert(name != "")

		// Next perform a translation of the value accordingly.
		result[name] = p.makeTerraformOutput(value, info)
	}

	if glog.V(5) {
		for k, v := range result {
			glog.V(5).Infof("Terraform output %v = %v", k, v)
		}
	}

	return result
}

// makeTerraformOutput takes a single Terraform property and returns the Lumi equivalent.
func (p *Provider) makeTerraformOutput(v interface{}, schema *SchemaInfo) resource.PropertyValue {
	if v == nil {
		return resource.NewNullProperty()
	}
	switch t := v.(type) {
	case bool:
		return resource.NewBoolProperty(t)
	case int:
		return resource.NewNumberProperty(float64(t))
	case string:
		// If the string is the special unknown property sentinel, reflect back an unknown computed property.  Note that
		// Terraform doesn't carry the types along with it, so the best we can do is give back a computed string.
		if t == config.UnknownVariableValue {
			elem := resource.Computed{Element: resource.NewStringProperty("")}
			return resource.NewComputedProperty(elem)
		}
		// Else it's just a string.
		return resource.NewStringProperty(t)
	case []interface{}:
		var elemschema *SchemaInfo
		if schema != nil {
			elemschema = schema.Elem
		}
		var arr []resource.PropertyValue
		for _, elem := range t {
			arr = append(arr, p.makeTerraformOutput(elem, elemschema))
		}
		return resource.NewArrayProperty(arr)
	case map[string]interface{}:
		var fldschemas map[string]*SchemaInfo
		if schema != nil {
			fldschemas = schema.Fields
		}
		obj := p.makeTerraformOutputs(t, fldschemas)
		return resource.NewObjectProperty(obj)
	default:
		contract.Failf("Unexpected TF output property value: %v", v)
		return resource.NewNullProperty()
	}
}

// makeTerraformConfig creates a Terraform config map, used in state and diff calculations, from a Lumi property map.
func (p *Provider) makeTerraformConfig(res *LumiResource, m resource.PropertyMap,
	schema map[string]*SchemaInfo, defaults bool) (*terraform.ResourceConfig, error) {
	// Convert the resource bag into an untyped map, and then create the resource config object.
	inputs, err := p.makeTerraformInputs(res, m, schema, defaults)
	if err != nil {
		return nil, err
	}
	return p.makeTerraformConfigFromInputs(inputs)
}

// makeTerraformConfigFromRPC creates a Terraform config map from a Lumi RPC property map.
func (p *Provider) makeTerraformConfigFromRPC(res *LumiResource, m *pbstruct.Struct,
	schema map[string]*SchemaInfo, allowUnknowns, defaults bool) (*terraform.ResourceConfig, error) {
	props, err := plugin.UnmarshalProperties(m,
		plugin.MarshalOptions{AllowUnknowns: allowUnknowns, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	return p.makeTerraformConfig(res, props, schema, defaults)
}

// makeTerraformConfigFromInputs creates a new Terraform configuration object from a set of Terraform inputs.
func (p *Provider) makeTerraformConfigFromInputs(inputs map[string]interface{}) (*terraform.ResourceConfig, error) {
	cfg, err := config.NewRawConfig(inputs)
	if err != nil {
		return nil, err
	}
	return terraform.NewResourceConfig(cfg), nil
}

// makeTerraformAttributes converts a Lumi property bag into its Terraform equivalent.  This requires
// flattening everything and serializing individual properties as strings.  This is a little awkward, but it's how
// Terraform represents resource properties (schemas are simply sugar on top).
func (p *Provider) makeTerraformAttributes(res *LumiResource, m resource.PropertyMap,
	schema map[string]*SchemaInfo, defaults bool) (map[string]string, error) {
	// Turn the resource properties into a map.  For the most part, this is a straight Mappable, but we use MapReplace
	// because we use float64s and Terraform uses ints, to represent numbers.
	inputs, err := p.makeTerraformInputs(res, m, schema, defaults)
	if err != nil {
		return nil, err
	}
	return p.makeTerraformAttributesFromInputs(inputs), nil
}

// makeTerraformAttributesFromRPC unmarshals an RPC property map and calls through to makeTerraformAttributes.
func (p *Provider) makeTerraformAttributesFromRPC(res *LumiResource, m *pbstruct.Struct,
	schema map[string]*SchemaInfo, allowUnknowns, defaults bool) (map[string]string, error) {
	props, err := plugin.UnmarshalProperties(m,
		plugin.MarshalOptions{AllowUnknowns: allowUnknowns, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	return p.makeTerraformAttributes(res, props, schema, defaults)
}

// makeTerraformAttributesFromInputs creates a flat Terraform map from a structured set of Terraform inputs.
func (p *Provider) makeTerraformAttributesFromInputs(inputs map[string]interface{}) map[string]string {
	return flatmap.Flatten(inputs)
}

// makeTerraformDiff takes a bag of old and new properties, and returns two things: the existing resource's state as
// an attribute map, alongside a Terraform diff for the old versus new state.  If there was no existing state, the
// returned attributes will be empty (because the resource doesn't yet exist).
func (p *Provider) makeTerraformDiff(old resource.PropertyMap, new resource.PropertyMap,
	schema map[string]*SchemaInfo) (*terraform.InstanceState, *terraform.InstanceDiff, error) {
	// BUGBUG[pulumi/pulumi-terraform#22]: avoid spilling except for during creation.
	diff := make(map[string]*terraform.ResourceAttrDiff)
	// Add all new property values.
	if new != nil {
		inputs, err := p.makeTerraformAttributes(nil, new, schema, false)
		if err != nil {
			return nil, nil, err
		}
		for p, v := range inputs {
			if diff[p] == nil {
				diff[p] = &terraform.ResourceAttrDiff{}
			}
			diff[p].New = v
		}
	}
	// Now add all old property values, provided they exist in new.
	existing := make(map[string]string)
	if old != nil {
		inputs, err := p.makeTerraformAttributes(nil, old, schema, false)
		if err != nil {
			return nil, nil, err
		}
		for p, v := range inputs {
			if d, has := diff[p]; has {
				d.Old = v
			}
			existing[p] = v
		}
	}
	return &terraform.InstanceState{Attributes: existing},
		&terraform.InstanceDiff{Attributes: diff}, nil
}

// makeTerraformDiffFromRPC takes RPC maps of old and new properties, unmarshals them, and calls into makeTerraformDiff.
func (p *Provider) makeTerraformDiffFromRPC(old *pbstruct.Struct, new *pbstruct.Struct,
	schema map[string]*SchemaInfo) (*terraform.InstanceState, *terraform.InstanceDiff, error) {
	var err error
	var oldprops resource.PropertyMap
	if old != nil {
		oldprops, err = plugin.UnmarshalProperties(old,
			plugin.MarshalOptions{SkipNulls: true})
		if err != nil {
			return nil, nil, err
		}
	}
	var newprops resource.PropertyMap
	if new != nil {
		newprops, err = plugin.UnmarshalProperties(new,
			plugin.MarshalOptions{AllowUnknowns: true, SkipNulls: true})
		if err != nil {
			return nil, nil, err
		}
	}
	return p.makeTerraformDiff(oldprops, newprops, schema)
}

// Configure configures the underlying Terraform provider with the live Lumi variable state.
func (p *Provider) Configure(ctx context.Context, req *lumirpc.ConfigureRequest) (*pbempty.Empty, error) {
	// Fetch the map of tokens to values.  It will be in the form of fully qualified tokens, so
	// we will need to translate into simply the configuration variable names.
	vars := make(resource.PropertyMap)
	for k, v := range req.GetVariables() {
		mm, err := tokens.ParseModuleMember(k)
		if err != nil {
			return nil, errors.Wrapf(err, "malformed configuration token '%v'", k)
		}
		if mm.Module() == p.baseConfigMod() || mm.Module() == p.configMod() {
			vars[resource.PropertyKey(mm.Name())] = resource.NewStringProperty(v)
		}
	}

	// Now make a Terraform config map out of the variables.
	config, err := p.makeTerraformConfig(nil, vars, p.info.Config, true)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling config state to Terraform")
	}

	// Perform validation of the config state so we can offer nice errors.
	keys, errs := p.tf.Validate(config)
	if len(keys) > 0 {
		return nil, errors.Errorf("one or more errors occurred while configuring key '%v' (%v)", keys[0], errs[0])
	}

	// Now actually attempt to do the configuring and return its resulting error (if any).
	if err = p.tf.Configure(config); err != nil {
		return nil, err
	}
	return &pbempty.Empty{}, nil
}

// Check validates that the given property bag is valid for a resource of the given type.
func (p *Provider) Check(ctx context.Context, req *lumirpc.CheckRequest) (*lumirpc.CheckResponse, error) {
	urn := resource.URN(req.GetUrn())
	t := urn.Type()
	res, has := p.resource(t)
	if !has {
		return nil, errors.Errorf("Unrecognized resource type (Check): %v", t)
	}
	props, err := plugin.UnmarshalProperties(req.GetProperties(),
		plugin.MarshalOptions{AllowUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	lumires := &LumiResource{URN: urn, Properties: props}

	// Step one is to populate any default values.  This is a two-stage process.  First we must create the
	// bridge-specific diffs, in cases where the overlays inject their own default values.
	inputs, err := p.makeTerraformInputs(lumires, props, res.Schema.Fields, true)
	if err != nil {
		return nil, errors.Errorf("Error preparing %v property state: %v", urn, err)
	}

	// Next, we let Terraform inject its own defaults, by way of Diff.  This may seem supremely bizarre, however, if
	// you carefully inspect how Terraform's pkg/helper/schema/ field readers work, default values are only injected
	// for the config variety.  The config variety is not chained in the multi-field reader structure during ordinary
	// CRUD operations, however; instead, it is chained only during ResourceConfig-related ones.  Diff is one such
	// operation that chains config in, which gives us back a Diff that is perfectly populated with the defaults.
	info := &terraform.InstanceInfo{Type: res.TF.Name}
	attrs := p.makeTerraformAttributesFromInputs(inputs)
	state := &terraform.InstanceState{Attributes: attrs}
	rescfg, err := p.makeTerraformConfigFromInputs(inputs)
	if err != nil {
		return nil, errors.Errorf("Error preparing %v's resource config state: %v", urn, err)
	}
	diff, err := p.tf.Diff(info, state, rescfg)
	if err != nil {
		return nil, errors.Errorf("Error preparing %v's resource diff (for defaults): %v", urn, err)
	}

	// After all is said and done, we need to go back and return only what got populated as a diff from the origin.
	defaults := make(resource.PropertyMap)
	outputs := p.makeTerraformResult(attrs, res.Schema.Fields)
	if outdiff := props.Diff(outputs); outdiff != nil {
		// Just recognized adds/changes, since these are defaults.
		for k := range outdiff.Adds {
			defaults[k] = outputs[k]
		}
		for k := range outdiff.Updates {
			defaults[k] = outputs[k]
		}
	}
	if diff != nil {
		// Expand the flatmap, so all arrays and sets are in their normal form, and then record any changes.
		flatolds := make(flatmap.Map)
		flatnews := make(flatmap.Map)
		for k, attr := range diff.Attributes {
			if attr.Old != "" {
				flatolds[k] = attr.Old
			}
			if attr.New != "" {
				flatnews[k] = attr.New
			}
		}
		for _, k := range flatnews.Keys() {
			var oldv interface{}
			if flatolds.Contains(k) {
				oldv = flatmap.Expand(flatolds, k)
			}
			newv := flatmap.Expand(flatnews, k)
			if !reflect.DeepEqual(oldv, newv) {
				name, info := getInfoFromTerraformName(k, res.Schema.Fields)
				defaults[name] = p.makeTerraformOutput(newv, info)
			}
		}
	}

	// Now check with the resource provider to see if the values pass muster.
	warns, errs := p.tf.ValidateResource(res.TF.Name, rescfg)

	// For each warning, emit a warning, but don't fail the check.
	for _, warn := range warns {
		if err = p.host.Log(diag.Warning, fmt.Sprintf("%v verification warning: %v", urn, warn)); err != nil {
			return nil, err
		}
	}

	// Now produce a return value of any properties that failed verification.
	var failures []*lumirpc.CheckFailure
	for _, err := range errs {
		failures = append(failures, &lumirpc.CheckFailure{
			Reason: err.Error(),
		})
	}

	defprops, err := plugin.MarshalProperties(defaults, plugin.MarshalOptions{AllowUnknowns: true})
	if err != nil {
		return nil, err
	}
	return &lumirpc.CheckResponse{Defaults: defprops, Failures: failures}, nil
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (p *Provider) Diff(ctx context.Context, req *lumirpc.DiffRequest) (*lumirpc.DiffResponse, error) {
	urn := resource.URN(req.GetUrn())
	t := urn.Type()
	res, has := p.resource(t)
	if !has {
		return nil, errors.Errorf("Unrecognized resource type (Diff): %v", urn)
	}
	glog.V(9).Infof("tfbridge/Provider.Diff: lumi='%v', tf=%v", urn, res.TF.Name)

	// To figure out if we have a replacement, perform the diff and then look for RequiresNew flags.
	inputs, err := p.makeTerraformAttributesFromRPC(
		nil, req.GetOlds(), res.Schema.Fields, false, false)
	if err != nil {
		return nil, errors.Errorf("Error preparing %v old property state: %v", urn, err)
	}
	info := &terraform.InstanceInfo{Type: res.TF.Name}
	state := &terraform.InstanceState{ID: req.GetId(), Attributes: inputs}
	config, err := p.makeTerraformConfigFromRPC(
		nil, req.GetNews(), res.Schema.Fields, true, false)
	if err != nil {
		return nil, errors.Errorf("Error preparing %v property state: %v", urn, err)
	}
	diff, err := p.tf.Diff(info, state, config)
	if err != nil {
		return nil, errors.Errorf("Error diffing %v old and new state: %v", urn, err)
	}

	// Each RequiresNew translates into a replacement.
	var replaces []string
	for k, attr := range diff.Attributes {
		if attr.RequiresNew {
			replaces = append(replaces, k)
		}
	}

	return &lumirpc.DiffResponse{Replaces: replaces}, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.  (The input ID
// must be blank.)  If this call fails, the resource must not have been created (i.e., it is "transacational").
func (p *Provider) Create(ctx context.Context, req *lumirpc.CreateRequest) (*lumirpc.CreateResponse, error) {
	urn := resource.URN(req.GetUrn())
	t := urn.Type()
	res, has := p.resource(t)
	if !has {
		return nil, errors.Errorf("Unrecognized resource type (Create): %v", t)
	}
	glog.V(9).Infof("tfbridge/Provider.Create: lumi='%v', tf=%v", urn, res.TF.Name)

	// To get Terraform to create a new resource, the ID msut be blank and existing state must be empty (since the
	// resource does not exist yet), and the diff object should have no old state and all of the new state.
	info := &terraform.InstanceInfo{Type: res.TF.Name}
	state, diff, err := p.makeTerraformDiffFromRPC(nil, req.GetProperties(), res.Schema.Fields)
	if err != nil {
		return nil, errors.Errorf("Error preparing %v's property state: %v", urn, err)
	}
	newstate, err := p.tf.Apply(info, state, diff)
	if err != nil {
		return nil, err
	}

	// Create the ID and property maps and return them.
	props := p.makeTerraformResult(newstate.Attributes, res.Schema.Fields)
	mprops, err := plugin.MarshalProperties(props, plugin.MarshalOptions{})
	if err != nil {
		return nil, err
	}
	return &lumirpc.CreateResponse{Id: newstate.ID, Properties: mprops}, nil
}

// Update updates an existing resource with new values.  Only those values in the provided property bag are updated
// to new values.  The resource ID is returned and may be different if the resource had to be recreated.
func (p *Provider) Update(ctx context.Context, req *lumirpc.UpdateRequest) (*lumirpc.UpdateResponse, error) {
	urn := resource.URN(req.GetUrn())
	t := urn.Type()
	res, has := p.resource(t)
	if !has {
		return nil, errors.Errorf("Unrecognized resource type (Update): %v", t)
	}
	glog.V(9).Infof("tfbridge/Provider.Update: lumi='%v', tf=%v", urn, res.TF.Name)

	// Create a state state with the ID to update, a diff with old and new states, and perform the apply.
	state, diff, err := p.makeTerraformDiffFromRPC(req.GetOlds(), req.GetNews(), res.Schema.Fields)
	if err != nil {
		return nil, errors.Errorf("Error preparing %v old and new state/diffs: %v", urn, err)
	}
	state.ID = req.GetId()

	info := &terraform.InstanceInfo{Type: res.TF.Name}
	newstate, err := p.tf.Apply(info, state, diff)
	if err != nil {
		return nil, errors.Errorf("Error applying %v update: %v", urn, err)
	}

	props := p.makeTerraformResult(newstate.Attributes, res.Schema.Fields)
	mprops, err := plugin.MarshalProperties(props, plugin.MarshalOptions{})
	if err != nil {
		return nil, err
	}
	return &lumirpc.UpdateResponse{Properties: mprops}, nil
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed to still exist.
func (p *Provider) Delete(ctx context.Context, req *lumirpc.DeleteRequest) (*pbempty.Empty, error) {
	urn := resource.URN(req.GetUrn())
	t := urn.Type()
	res, has := p.resource(t)
	if !has {
		return nil, errors.Errorf("Unrecognized resource type (Delete): %v", t)
	}
	glog.V(9).Infof("tfbridge/Provider.Delete: lumi='%v', tf=%v", urn, res.TF.Name)

	// Fetch the resource attributes since many providers need more than just the ID to perform the delete.
	attrs, err := p.makeTerraformAttributesFromRPC(
		nil, req.GetProperties(), res.Schema.Fields, false, false)
	if err != nil {
		return nil, err
	}

	// Create a new state, with no diff, that is missing an ID.  Terraform will interpret this as a create operation.
	info := &terraform.InstanceInfo{Type: res.TF.Name}
	state := &terraform.InstanceState{ID: req.GetId(), Attributes: attrs}
	if _, err := p.tf.Apply(info, state, &terraform.InstanceDiff{Destroy: true}); err != nil {
		return nil, errors.Errorf("Error apply %v deletion: %v", urn, err)
	}
	return &pbempty.Empty{}, nil
}
