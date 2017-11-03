// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/glog"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	pbstruct "github.com/golang/protobuf/ptypes/struct"
	"github.com/hashicorp/terraform/flatmap"
	"github.com/hashicorp/terraform/helper/schema"
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

// Provider implements the Pulumi resource provider operations for any Terraform plugin.
type Provider struct {
	host        *provider.HostClient               // the RPC link back to the Pulumi engine.
	module      string                             // the Terraform module name.
	tf          *schema.Provider                   // the Terraform resource provider to use.
	info        ProviderInfo                       // overlaid info about this provider.
	resources   map[tokens.Type]Resource           // a map of Pulumi type tokens to resource info.
	dataSources map[tokens.ModuleMember]DataSource // a map of Pulumi module tokens to data sources.
}

// Resource wraps both the Terraform resource type info plus the overlay resource info.
type Resource struct {
	Schema   *ResourceInfo             // optional provider overrides.
	TF       terraform.ResourceType    // Terraform resource info.
	TFSchema map[string]*schema.Schema // the Terraform resource schema.
}

// DataSource wraps both the Terraform data source (resource) type info plus the overlay resource info.
type DataSource struct {
	Schema   *DataSourceInfo           // optional provider overrides.
	TF       terraform.DataSource      // Terraform resource info.
	TFSchema map[string]*schema.Schema // the Terraform data source schema.
}

// NewProvider creates a new Pulumi RPC server wired up to the given host and wrapping the given Terraform provider.
func NewProvider(host *provider.HostClient, module string,
	tf *schema.Provider, info ProviderInfo) *Provider {
	p := &Provider{
		host:   host,
		module: module,
		tf:     tf,
		info:   info,
	}
	p.initResourceMaps()
	return p
}

var _ lumirpc.ResourceProviderServer = (*Provider)(nil)

func (p *Provider) pkg() tokens.Package          { return tokens.Package(p.module) }
func (p *Provider) indexMod() tokens.Module      { return tokens.Module(p.pkg() + ":index") }
func (p *Provider) baseConfigMod() tokens.Module { return tokens.Module(p.pkg() + ":config") }
func (p *Provider) baseDataMod() tokens.Module   { return tokens.Module(p.pkg() + ":data") }
func (p *Provider) configMod() tokens.Module     { return p.baseConfigMod() + "/vars" }

// initResourceMaps creates maps from Pulumi types and tokens to Terraform resource type.
func (p *Provider) initResourceMaps() {
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
			// Manufacture a token with the package, module, and resource type name.
			camelName, pascalName := p.camelPascalPulumiName(res.Name)
			tok = tokens.Type(string(p.pkg()) + ":" + camelName + ":" + pascalName)
		}

		p.resources[tok] = Resource{
			TF:       res,
			TFSchema: p.tf.ResourcesMap[res.Name].Schema,
			Schema:   schema,
		}
	}

	// Fetch a list of all data source types handled by this provider and make a similar map.
	p.dataSources = make(map[tokens.ModuleMember]DataSource)
	for _, ds := range p.tf.DataSources() {
		var tok tokens.ModuleMember

		// See if there is override information for this resource.  If yes, use that to decode the token.
		var schema *DataSourceInfo
		if p.info.DataSources != nil {
			schema = p.info.DataSources[ds.Name]
			if schema != nil {
				tok = schema.Tok
			}
		}

		// Otherwise, we default to the standard naming scheme.
		if tok == "" {
			// Manufacture a token with the data module and camel-cased name.
			camelName, _ := p.camelPascalPulumiName(ds.Name)
			tok = tokens.ModuleMember(string(p.baseDataMod()) + ":" + camelName)
		}

		p.dataSources[tok] = DataSource{
			TF:       ds,
			TFSchema: p.tf.DataSourcesMap[ds.Name].Schema,
			Schema:   schema,
		}
	}
}

// camelPascalPulumiName returns the camel and pascal cased name for a given terraform name.
func (p *Provider) camelPascalPulumiName(name string) (string, string) {
	// Strip off the module prefix (e.g., "aws_") and then return the camel- and Pascal-cased names.
	prefix := p.module + "_" // all resources will have this prefix.
	contract.Assertf(strings.HasPrefix(name, prefix),
		"Expected all Terraform resources in this module to have a '%v' prefix", prefix)
	name = name[len(prefix):]
	return TerraformToPulumiName(name, false), TerraformToPulumiName(name, true)
}

// Configure configures the underlying Terraform provider with the live Pulumi variable state.
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
	config, err := MakeTerraformConfig(nil, vars, p.tf.Schema, p.info.Config, true)
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

// prepareInputDefaults takes an input "resource" with state, does everything necessary to round-trip with the
// Terraform provider and populate the default variables, and returns the resulting config and state.
func (p *Provider) prepareInputsDefaults(tfname string, lumires *PulumiResource,
	tfs map[string]*schema.Schema, ps map[string]*SchemaInfo) (*terraform.InstanceInfo, *terraform.InstanceState,
	*terraform.InstanceDiff, *terraform.ResourceConfig, error) {
	// Step one is to populate any default values.  This is a two-stage process.  First we must create the
	// bridge-specific diffs, in cases where the overlays inject their own default values.
	inputs, err := MakeTerraformInputs(lumires, lumires.Properties, tfs, ps, true, false)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "couldn't prepare resource %v input state", tfname)
	}

	// Next, we let Terraform inject its own defaults, by way of Diff.  This may seem supremely bizarre, however, if
	// you carefully inspect how Terraform's pkg/helper/schema/ field readers work, default values are only injected
	// for the config variety.  The config variety is not chained in the multi-field reader structure during ordinary
	// CRUD operations, however; instead, it is chained only during ResourceConfig-related ones.  Diff is one such
	// operation that chains config in, which gives us back a Diff that is perfectly populated with the defaults.
	info := &terraform.InstanceInfo{Type: tfname}
	attrs := MakeTerraformAttributesFromInputs(inputs)
	state := &terraform.InstanceState{Attributes: attrs}
	rescfg, err := MakeTerraformConfigFromInputs(inputs)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "couldn't prepare resource %v config state", tfname)
	}
	diff, err := p.tf.Diff(info, state, rescfg)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrapf(err, "couldn't prepare resource %v diff (for defaults)", tfname)
	}

	return info, state, diff, rescfg, nil
}

// Check validates that the given property bag is valid for a resource of the given type.
func (p *Provider) Check(ctx context.Context, req *lumirpc.CheckRequest) (*lumirpc.CheckResponse, error) {
	urn := resource.URN(req.GetUrn())
	t := urn.Type()
	res, has := p.resources[t]
	if !has {
		return nil, errors.Errorf("Unrecognized resource type (Check): %v", t)
	}

	// Unmarshal the properties.
	props, err := plugin.UnmarshalProperties(req.GetProperties(),
		plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	// Now fetch the default values so that (a) we can return them to the caller and (b) so that validation
	// includes the default values.  Otherwise, the provider wouldn't be presented with its own defaults.
	tfname := res.TF.Name
	_, state, diff, rescfg, err := p.prepareInputsDefaults(
		tfname, &PulumiResource{URN: urn, Properties: props}, res.TFSchema, res.Schema.Fields)
	if err != nil {
		return nil, err
	}

	// Now check with the resource provider to see if the values pass muster.
	warns, errs := p.tf.ValidateResource(tfname, rescfg)
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

	// After all is said and done, we need to go back and return only what got populated as a diff from the origin.
	defaults := make(resource.PropertyMap)
	outputs := MakeTerraformResult(state.Attributes, res.TFSchema, res.Schema.Fields)
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
				name, tfi, psi := getInfoFromTerraformName(k, res.TFSchema, res.Schema.Fields, false)
				defaults[name] = MakeTerraformOutput(newv, tfi, psi, false)
			}
		}
	}
	defprops, err := plugin.MarshalProperties(defaults, plugin.MarshalOptions{KeepUnknowns: true})
	if err != nil {
		return nil, err
	}

	return &lumirpc.CheckResponse{Defaults: defprops, Failures: failures}, nil
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (p *Provider) Diff(ctx context.Context, req *lumirpc.DiffRequest) (*lumirpc.DiffResponse, error) {
	urn := resource.URN(req.GetUrn())
	t := urn.Type()
	res, has := p.resources[t]
	if !has {
		return nil, errors.Errorf("Unrecognized resource type (Diff): %v", urn)
	}
	glog.V(9).Infof("tfbridge/Provider.Diff: lumi='%v', tf=%v", urn, res.TF.Name)

	// To figure out if we have a replacement, perform the diff and then look for RequiresNew flags.
	inputs, err := MakeTerraformAttributesFromRPC(
		nil, req.GetOlds(), res.TFSchema, res.Schema.Fields, false, false)
	if err != nil {
		return nil, errors.Errorf("Error preparing %v old property state: %v", urn, err)
	}
	info := &terraform.InstanceInfo{Type: res.TF.Name}
	state := &terraform.InstanceState{ID: req.GetId(), Attributes: inputs}
	config, err := MakeTerraformConfigFromRPC(
		nil, req.GetNews(), res.TFSchema, res.Schema.Fields, true, false)
	if err != nil {
		return nil, errors.Errorf("Error preparing %v property state: %v", urn, err)
	}
	diff, err := p.tf.Diff(info, state, config)
	if err != nil {
		return nil, errors.Errorf("Error diffing %v old and new state: %v", urn, err)
	}

	// Each RequiresNew translates into a replacement.
	var replaces []string
	replaced := make(map[resource.PropertyKey]bool)
	for k, attr := range diff.Attributes {
		if attr.RequiresNew {
			name, _, _ := getInfoFromTerraformName(k, res.TFSchema, res.Schema.Fields, false)
			replaces = append(replaces, string(name))
			replaced[name] = true
		}
	}

	// For all properties that are ForceNew, but didn't change, assume they are stable.  Also recognize
	// overlays that have requested that we treat specific properties as stable.
	var stables []string
	for k, sch := range res.TFSchema {
		name, _, cust := getInfoFromTerraformName(k, res.TFSchema, res.Schema.Fields, false)
		if !replaced[name] &&
			(sch.ForceNew || (cust != nil && cust.Stable != nil && *cust.Stable)) {
			stables = append(stables, string(name))
		}
	}

	return &lumirpc.DiffResponse{
		Replaces: replaces,
		Stables:  stables,
	}, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.  (The input ID
// must be blank.)  If this call fails, the resource must not have been created (i.e., it is "transacational").
func (p *Provider) Create(ctx context.Context, req *lumirpc.CreateRequest) (*lumirpc.CreateResponse, error) {
	urn := resource.URN(req.GetUrn())
	t := urn.Type()
	res, has := p.resources[t]
	if !has {
		return nil, errors.Errorf("Unrecognized resource type (Create): %v", t)
	}
	glog.V(9).Infof("tfbridge/Provider.Create: lumi='%v', tf=%v", urn, res.TF.Name)

	// To get Terraform to create a new resource, the ID msut be blank and existing state must be empty (since the
	// resource does not exist yet), and the diff object should have no old state and all of the new state.
	info := &terraform.InstanceInfo{Type: res.TF.Name}
	state, diff, err := MakeTerraformDiffFromRPC(nil, req.GetProperties(), res.TFSchema, res.Schema.Fields)
	if err != nil {
		return nil, errors.Errorf("Error preparing %v's property state: %v", urn, err)
	}
	newstate, err := p.tf.Apply(info, state, diff)
	if err != nil {
		return nil, err
	}

	// Create the ID and property maps and return them.
	props := MakeTerraformResult(newstate.Attributes, res.TFSchema, res.Schema.Fields)
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
	res, has := p.resources[t]
	if !has {
		return nil, errors.Errorf("Unrecognized resource type (Update): %v", t)
	}
	glog.V(9).Infof("tfbridge/Provider.Update: lumi='%v', tf=%v", urn, res.TF.Name)

	// Create a state state with the ID to update, a diff with old and new states, and perform the apply.
	state, diff, err := MakeTerraformDiffFromRPC(
		req.GetOlds(), req.GetNews(), res.TFSchema, res.Schema.Fields)
	if err != nil {
		return nil, errors.Errorf("Error preparing %v old and new state/diffs: %v", urn, err)
	}
	state.ID = req.GetId()

	info := &terraform.InstanceInfo{Type: res.TF.Name}
	newstate, err := p.tf.Apply(info, state, diff)
	if err != nil {
		return nil, errors.Errorf("Error applying %v update: %v", urn, err)
	}

	props := MakeTerraformResult(newstate.Attributes, res.TFSchema, res.Schema.Fields)
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
	res, has := p.resources[t]
	if !has {
		return nil, errors.Errorf("Unrecognized resource type (Delete): %v", t)
	}
	glog.V(9).Infof("tfbridge/Provider.Delete: lumi='%v', tf=%v", urn, res.TF.Name)

	// Fetch the resource attributes since many providers need more than just the ID to perform the delete.
	attrs, err := MakeTerraformAttributesFromRPC(
		nil, req.GetProperties(), res.TFSchema, res.Schema.Fields, false, false)
	if err != nil {
		return nil, err
	}

	// Create a new state, with no diff, that is missing an ID.  Terraform will interpret this as a create operation.
	info := &terraform.InstanceInfo{Type: res.TF.Name}
	state := &terraform.InstanceState{ID: req.GetId(), Attributes: attrs}
	if _, err := p.tf.Apply(info, state, &terraform.InstanceDiff{Destroy: true}); err != nil {
		return nil, errors.Errorf("Error applying %v deletion: %v", urn, err)
	}
	return &pbempty.Empty{}, nil
}

// Invoke dynamically executes a built-in function in the provider.
func (p *Provider) Invoke(ctx context.Context, req *lumirpc.InvokeRequest) (*lumirpc.InvokeResponse, error) {
	tok := tokens.ModuleMember(req.GetTok())
	ds, has := p.dataSources[tok]
	if !has {
		return nil, errors.Errorf("Unrecognized data function (Invoke): %v", tok)
	}
	glog.V(9).Infof("tfbridge/Provider.Invoke: tok=%v", tok)

	// Unmarshal the arguments.
	args, err := plugin.UnmarshalProperties(req.GetArgs(),
		plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	// First, create the inputs.
	tfname := ds.TF.Name
	inputs, err := MakeTerraformInputs(
		&PulumiResource{Properties: args}, args, ds.TFSchema, ds.Schema.Fields, true, false)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't prepare resource %v input state", tfname)
	}

	// Next, ensure the inputs are valid before actually performing the invoaction.
	info := &terraform.InstanceInfo{Type: tfname}
	rescfg, err := MakeTerraformConfigFromInputs(inputs)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't make config for %v validation", tfname)
	}
	warns, errs := p.tf.ValidateDataSource(tfname, rescfg)
	for _, warn := range warns {
		if err = p.host.Log(diag.Warning, fmt.Sprintf("%v verification warning: %v", tok, warn)); err != nil {
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

	// If there are no failures in verification, go ahead and perform the invocation.
	var ret *pbstruct.Struct
	if len(failures) == 0 {
		diff, err := p.tf.ReadDataDiff(info, rescfg)
		if err != nil {
			return nil, errors.Wrapf(err, "error reading data source '%v' diff", tok)
		}

		invoke, err := p.tf.ReadDataApply(info, diff)
		if err != nil {
			return nil, errors.Wrapf(err, "error invoking %v", tok)
		}
		ret, err = plugin.MarshalProperties(
			MakeTerraformResult(invoke.Attributes, ds.TFSchema, ds.Schema.Fields),
			plugin.MarshalOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &lumirpc.InvokeResponse{
		Return:   ret,
		Failures: failures,
	}, nil
}
