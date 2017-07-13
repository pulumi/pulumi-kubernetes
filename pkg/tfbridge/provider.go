// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"fmt"
	"strings"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	pbstruct "github.com/golang/protobuf/ptypes/struct"
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/flatmap"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	"github.com/pulumi/lumi/pkg/resource"
	"github.com/pulumi/lumi/pkg/resource/plugin"
	"github.com/pulumi/lumi/pkg/resource/provider"
	"github.com/pulumi/lumi/pkg/tokens"
	"github.com/pulumi/lumi/pkg/util/contract"
	"github.com/pulumi/lumi/sdk/go/pkg/lumirpc"
	"golang.org/x/net/context"
)

// Provider implements the Lumi resource provider operations for any Terraform plugin.
type Provider struct {
	host      *provider.HostClient                   // the RPC link back to the Lumi engine.
	tf        terraform.ResourceProvider             // the Terraform resource provider to use.
	module    string                                 // the Terraform module name.
	resources map[tokens.Type]terraform.ResourceType // a map of Lumi type tokens to Terraform type structs.
}

// NewProvider creates a new Lumi RPC server wired up to the given host and wrapping the given Terraform provider.
func NewProvider(host *provider.HostClient, tf terraform.ResourceProvider, module string) *Provider {
	// TODO: audit computed logic to ensure we flow from Lumi's notion of unknowns to TF computeds properly.
	p := &Provider{
		host:   host,
		tf:     tf,
		module: module,
	}
	p.initResourceMap()
	return p
}

var _ lumirpc.ResourceProviderServer = (*Provider)(nil)

func (p *Provider) pkg() tokens.Package      { return tokens.Package(p.module) }
func (p *Provider) indexMod() tokens.Module  { return tokens.Module(p.pkg() + ":index") }
func (p *Provider) configMod() tokens.Module { return tokens.Module(p.pkg() + ":config/vars") }

// NameProperty is the resource property used to assign names for URN assignment.
const NameProperty = "name"

// tfResource looks up the Terraform resource provider from its Lumi type token.
func (p *Provider) tfResource(t tokens.Type) (terraform.ResourceType, bool) {
	res, has := p.resources[t]
	return res, has
}

// initResourceMap creates a simple map from Lumi to Terraform resource type.
func (p *Provider) initResourceMap() {
	prefix := p.module + "_"        // all resources will have this prefix.
	provinfo := Providers[p.module] // fetch name/schema overrides, if any.

	// Fetch a list of all resource types handled by this provider and make a map.
	p.resources = make(map[tokens.Type]terraform.ResourceType)
	for _, res := range p.tf.Resources() {
		var tok tokens.Type

		// See if there is override information for this resource.  If yes, use that to decode the token.
		if provinfo.Resources != nil {
			if resinfo, has := provinfo.Resources[res.Name]; has {
				tok = resinfo.Tok
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

		p.resources[tok] = res
	}
}

// Some functions used below for name and value transformations.
var (
	// lumiKeyRepl swaps out Lumi names for Terraform names.
	lumiKeyRepl = func(k string) (string, bool) {
		return LumiToTerraformName(k), true
	}
	// terraformKeyRepl swaps out Terraform names for Lumi names.
	terraformKeyRepl = func(k string) (resource.PropertyKey, bool) {
		return resource.PropertyKey(TerraformToLumiName(k, false)), true
	}
	// lumiValueRepl swaps out Lumi-style float64 for Terraform-style int numbers.
	lumiValueRepl = func(v resource.PropertyValue) (interface{}, bool) {
		if v.IsNumber() {
			return int(v.NumberValue()), true
		}
		return nil, false
	}
	// terraformValueRepl does the reverse, and swaps out Terraform ints for Lumi float64s.
	terraformValueRepl = func(v interface{}) (resource.PropertyValue, bool) {
		if i, isint := v.(int); isint {
			return resource.NewNumberProperty(float64(i)), true
		}
		return resource.PropertyValue{}, false
	}
)

// terraformToLumiProps expands a Terraform-style flatmap into an expanded Lumi resource property map.
func terraformToLumiProps(props map[string]string) resource.PropertyMap {
	res := make(map[string]interface{})
	for _, key := range flatmap.Map(props).Keys() {
		res[key] = flatmap.Expand(props, key)
	}
	return resource.NewPropertyMapFromMapRepl(res, terraformKeyRepl, terraformValueRepl)
}

// makeTerraformConfig creates a Terraform config map, used in state and diff calculations, from a Lumi property map.
func makeTerraformConfig(m resource.PropertyMap) (*terraform.ResourceConfig, error) {
	// Convert the resource bag into an untyped map, and then create the resource config object.
	ma := m.MapRepl(lumiKeyRepl, lumiValueRepl)
	cfg, err := config.NewRawConfig(ma)
	if err != nil {
		return nil, err
	}
	return terraform.NewResourceConfig(cfg), nil
}

// makeTerraformConfigFromRPC creates a Terraform config map from a Lumi RPC property map.
func makeTerraformConfigFromRPC(m *pbstruct.Struct) (*terraform.ResourceConfig, error) {
	props := plugin.UnmarshalProperties(nil, m, plugin.MarshalOptions{SkipNulls: true})
	return makeTerraformConfig(props)
}

// makeTerraformPropertyMap converts a Lumi property bag into its Terraform equivalent.  This requires
// flattening everything and serializing individual properties as strings.  This is a little awkward, but it's how
// Terraform represents resource properties (schemas are simply sugar on top).
func makeTerraformPropertyMap(m resource.PropertyMap) map[string]string {
	// Turn the resource properties into a map.  For the most part, this is a straight Mappable, but we use MapReplace
	// because we use float64s and Terraform uses ints, to represent numbers.
	props := m.MapRepl(lumiKeyRepl, lumiValueRepl)
	// FIXME: marshal/unmarshal sets properly.
	return flatmap.Flatten(props)
}

// makeTerraformPropertyMapFromRPC unmarshals an RPC property map and calls through to makeTerraformPropertyMap.
func makeTerraformPropertyMapFromRPC(m *pbstruct.Struct) map[string]string {
	props := plugin.UnmarshalProperties(nil, m, plugin.MarshalOptions{SkipNulls: true})
	return makeTerraformPropertyMap(props)
}

// makeTerraformDiff takes a bag of old and new properties, and returns two things: the attribute state to use for the
// current resource alongside a Terraform diff for the old and new.  If there was no old state, the first return is nil.
func makeTerraformDiff(
	old resource.PropertyMap, new resource.PropertyMap) (map[string]string, *terraform.InstanceDiff) {
	var attrs map[string]string
	diff := make(map[string]*terraform.ResourceAttrDiff)
	// Add all new property values.
	if new != nil {
		for p, v := range makeTerraformPropertyMap(new) {
			if diff[p] == nil {
				diff[p] = &terraform.ResourceAttrDiff{}
			}
			diff[p].New = v
		}
	}
	// Now add all old property values, provided they exist in new.
	if old != nil {
		attrs = makeTerraformPropertyMap(old)
		for p, v := range attrs {
			if diff[p] != nil {
				diff[p].Old = v
			}
		}
	}
	return attrs, &terraform.InstanceDiff{Attributes: diff}
}

// makeTerraformDiffFromRPC takes RPC maps of old and new properties, unmarshals them, and calls into makeTerraformDiff.
func makeTerraformDiffFromRPC(
	old *pbstruct.Struct, new *pbstruct.Struct) (map[string]string, *terraform.InstanceDiff) {
	oldprops := plugin.UnmarshalProperties(nil, old, plugin.MarshalOptions{SkipNulls: true})
	newprops := plugin.UnmarshalProperties(nil, new, plugin.MarshalOptions{SkipNulls: true})
	return makeTerraformDiff(oldprops, newprops)
}

// Configure configures the underlying Terraform provider with the live Lumi variable state.
func (p *Provider) Configure() error {
	// Read all properties from the config module.
	props, err := p.host.ReadLocations(tokens.Token(p.configMod()), true)
	if err != nil {
		return err
	}

	// Now make a map of each of the config token values.
	config, err := makeTerraformConfig(props)
	if err != nil {
		return err
	}

	// Perform validation of the config state so we can offer nice errors.
	keys, errs := p.tf.Validate(config)
	if len(keys) > 0 {
		// TODO: unify this with check by adding a Configure RPC method to the gRPC interface.
		return errors.Errorf("One or more errors occurred while configuring: %v (%v)", keys[0], errs[0])
	}

	// Now actually attempt to do the configuring and return its resulting error (if any).
	return p.tf.Configure(config)
}

// Check validates that the given property bag is valid for a resource of the given type.
func (p *Provider) Check(ctx context.Context, req *lumirpc.CheckRequest) (*lumirpc.CheckResponse, error) {
	t := tokens.Type(req.GetType())
	res, has := p.tfResource(t)
	if !has {
		return nil, fmt.Errorf("Unrecognized resource type (Check): %v", t)
	}

	// Manufacture a resource config to check, check it, and return any failures that result.
	rescfg, err := makeTerraformConfigFromRPC(req.GetProperties())
	if err != nil {
		return nil, err
	}
	keys, errs := p.tf.ValidateResource(res.Name, rescfg)
	var failures []*lumirpc.CheckFailure
	for i, key := range keys {
		failures = append(failures, &lumirpc.CheckFailure{
			Property: key,
			Reason:   errs[i].Error(),
		})
	}
	return &lumirpc.CheckResponse{Failures: failures}, nil
}

// Name names a given resource.  Sometimes this will be assigned by a developer, and so the provider
// simply fetches it from the property bag; other times, the provider will assign this based on its own algorithm.
// In any case, resources with the same name must be safe to use interchangeably with one another.
func (p *Provider) Name(ctx context.Context, req *lumirpc.NameRequest) (*lumirpc.NameResponse, error) {
	t := tokens.Type(req.GetType())
	if _, has := p.tfResource(t); !has {
		return nil, fmt.Errorf("Unrecognized resource type (Name): %v", t)
	}

	// All Terraform bridge providers will have a name property that we use for URN naming purposes.
	props := plugin.UnmarshalProperties(nil, req.GetProperties(), plugin.MarshalOptions{})
	name, has := props[NameProperty]
	if !has {
		return nil, errors.Errorf("Missing a '%v' property", NameProperty)
	} else if !name.IsString() {
		return nil, errors.Errorf("Expected a string '%v' property; got %v", NameProperty, name)
	}
	namestr := name.StringValue()
	if namestr == "" {
		if req.GetUnknowns()[NameProperty] {
			return nil, errors.Errorf("The '%v' property cannot be a computed expression", NameProperty)
		}
		return nil, errors.Errorf("The '%v' property cannot be the empty string", NameProperty)
	}
	return &lumirpc.NameResponse{Name: namestr}, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.  (The input ID
// must be blank.)  If this call fails, the resource must not have been created (i.e., it is "transacational").
func (p *Provider) Create(ctx context.Context, req *lumirpc.CreateRequest) (*lumirpc.CreateResponse, error) {
	t := tokens.Type(req.GetType())
	res, has := p.tfResource(t)
	if !has {
		return nil, fmt.Errorf("Unrecognized resource type (Create): %v", t)
	}

	// Create a new state, with no diff, that is missing an ID.  Terraform will interpret this as a create operation.
	info := &terraform.InstanceInfo{Type: res.Name}
	_, diff := makeTerraformDiff(nil,
		plugin.UnmarshalProperties(nil, req.GetProperties(), plugin.MarshalOptions{SkipNulls: true}))
	newstate, err := p.tf.Apply(info, nil, diff)
	if err != nil {
		return nil, err
	}
	return &lumirpc.CreateResponse{Id: newstate.ID}, nil
}

// Get reads the instance state identified by ID, returning a populated resource object, or an error if not found.
func (p *Provider) Get(ctx context.Context, req *lumirpc.GetRequest) (*lumirpc.GetResponse, error) {
	t := tokens.Type(req.GetType())
	res, has := p.tfResource(t)
	if !has {
		return nil, fmt.Errorf("Unrecognized resource type (Get): %v", t)
	}

	// To read the instance state, create a blank bit of data and ask the resource provider to recompute it.
	info := &terraform.InstanceInfo{Type: res.Name}
	state := &terraform.InstanceState{ID: req.GetId()}
	getstate, err := p.tf.Refresh(info, state)
	if err != nil {
		return nil, err
	}
	props := terraformToLumiProps(getstate.Attributes)
	return &lumirpc.GetResponse{
		Properties: plugin.MarshalProperties(nil, props, plugin.MarshalOptions{SkipNulls: true}),
	}, nil
}

// InspectChange checks what impacts a hypothetical update will have on the resource's properties.
func (p *Provider) InspectChange(
	ctx context.Context, req *lumirpc.InspectChangeRequest) (*lumirpc.InspectChangeResponse, error) {
	t := tokens.Type(req.GetType())
	res, has := p.tfResource(t)
	if !has {
		return nil, fmt.Errorf("Unrecognized resource type (InspectChange): %v", t)
	}

	// To figure out if we have a replacement, perform the diff and then look for RequiresNew flags.
	info := &terraform.InstanceInfo{Type: res.Name}
	state := &terraform.InstanceState{
		ID:         req.GetId(),
		Attributes: makeTerraformPropertyMapFromRPC(req.GetOlds()),
	}
	config, err := makeTerraformConfigFromRPC(req.GetNews())
	if err != nil {
		return nil, err
	}
	diff, err := p.tf.Diff(info, state, config)
	if err != nil {
		return nil, err
	}

	// Each RequiresNew translates into a replacement.
	var replaces []string
	for k, attr := range diff.Attributes {
		if attr.RequiresNew {
			replaces = append(replaces, k)
		}
	}

	return &lumirpc.InspectChangeResponse{Replaces: replaces}, nil
}

// Update updates an existing resource with new values.  Only those values in the provided property bag are updated
// to new values.  The resource ID is returned and may be different if the resource had to be recreated.
func (p *Provider) Update(ctx context.Context, req *lumirpc.UpdateRequest) (*pbempty.Empty, error) {
	t := tokens.Type(req.GetType())
	res, has := p.tfResource(t)
	if !has {
		return nil, fmt.Errorf("Unrecognized resource type (Delete): %v", t)
	}

	// Create a state state with the ID to update, a diff with old and new states, and perform the apply.
	info := &terraform.InstanceInfo{Type: res.Name}
	attrs, diff := makeTerraformDiffFromRPC(req.GetOlds(), req.GetNews())
	state := &terraform.InstanceState{
		ID:         req.GetId(),
		Attributes: attrs,
	}
	if _, err := p.tf.Apply(info, state, diff); err != nil {
		return nil, err
	}
	return &pbempty.Empty{}, nil
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed to still exist.
func (p *Provider) Delete(ctx context.Context, req *lumirpc.DeleteRequest) (*pbempty.Empty, error) {
	t := tokens.Type(req.GetType())
	res, has := p.tfResource(t)
	if !has {
		return nil, fmt.Errorf("Unrecognized resource type (Delete): %v", t)
	}

	// Create a new state, with no diff, that is missing an ID.  Terraform will interpret this as a create operation.
	info := &terraform.InstanceInfo{Type: res.Name}
	state := &terraform.InstanceState{ID: req.GetId()}
	if _, err := p.tf.Apply(info, state, &terraform.InstanceDiff{Destroy: true}); err != nil {
		return nil, err
	}
	return &pbempty.Empty{}, nil
}
