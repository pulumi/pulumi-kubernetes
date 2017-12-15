// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"fmt"
	"log"
	"strings"

	"github.com/golang/glog"
	pbempty "github.com/golang/protobuf/ptypes/empty"
	pbstruct "github.com/golang/protobuf/ptypes/struct"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/diag"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/plugin"
	"github.com/pulumi/pulumi/pkg/resource/provider"
	"github.com/pulumi/pulumi/pkg/tokens"
	"github.com/pulumi/pulumi/pkg/util/contract"
	pulumirpc "github.com/pulumi/pulumi/sdk/proto/go"
	"golang.org/x/net/context"
)

// Provider implements the Pulumi resource provider operations for any Terraform plugin.
type Provider struct {
	host        *provider.HostClient               // the RPC link back to the Pulumi engine.
	module      string                             // the Terraform module name.
	version     string                             // the plugin version number.
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
func NewProvider(host *provider.HostClient, module string, version string,
	tf *schema.Provider, info ProviderInfo) *Provider {
	p := &Provider{
		host:    host,
		module:  module,
		version: version,
		tf:      tf,
		info:    info,
	}
	p.initResourceMaps()
	return p
}

var _ pulumirpc.ResourceProviderServer = (*Provider)(nil)

func (p *Provider) pkg() tokens.Package          { return tokens.Package(p.module) }
func (p *Provider) indexMod() tokens.Module      { return tokens.Module(p.pkg() + ":index") }
func (p *Provider) baseConfigMod() tokens.Module { return tokens.Module(p.pkg() + ":config") }
func (p *Provider) baseDataMod() tokens.Module   { return tokens.Module(p.pkg() + ":data") }
func (p *Provider) configMod() tokens.Module     { return p.baseConfigMod() + "/vars" }

func (p *Provider) setLoggingContext(ctx context.Context) {
	log.SetOutput(&LogRedirector{
		writers: map[string]func(string) error{
			tfTracePrefix: func(msg string) error { return p.host.Log(ctx, diag.Debug, msg) },
			tfDebugPrefix: func(msg string) error { return p.host.Log(ctx, diag.Debug, msg) },
			tfInfoPrefix:  func(msg string) error { return p.host.Log(ctx, diag.Info, msg) },
			tfWarnPrefix:  func(msg string) error { return p.host.Log(ctx, diag.Warning, msg) },
			tfErrorPrefix: func(msg string) error { return p.host.Log(ctx, diag.Error, msg) },
		},
	})
}

func (p *Provider) label() string {
	return fmt.Sprintf("tf.Provider[%s]", p.module)
}

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
func (p *Provider) Configure(ctx context.Context, req *pulumirpc.ConfigureRequest) (*pbempty.Empty, error) {
	p.setLoggingContext(ctx)
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
		return nil, errors.Wrap(err, "could not marshal config state")
	}

	// Perform validation of the config state so we can offer nice errors.
	keys, errs := p.tf.Validate(config)
	if len(keys) > 0 {
		var result error
		for i, key := range keys {
			result = multierror.Append(result,
				errors.Wrapf(errs[i], "could not configure key %s", key))
		}
		return nil, result
	}

	// Now actually attempt to do the configuring and return its resulting error (if any).
	if err = p.tf.Configure(config); err != nil {
		return nil, err
	}
	return &pbempty.Empty{}, nil
}

// Check validates that the given property bag is valid for a resource of the given type.
func (p *Provider) Check(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	p.setLoggingContext(ctx)
	urn := resource.URN(req.GetUrn())
	t := urn.Type()
	res, has := p.resources[t]
	if !has {
		return nil, errors.Errorf("unrecognized resource type (Check): %s", t)
	}

	label := fmt.Sprintf("%s.Check(%s/%s)", p.label(), urn, res.TF.Name)
	glog.V(9).Infof("%s executing", label)

	// Unmarshal the old and new properties.
	var olds resource.PropertyMap
	var err error
	if req.GetOlds() != nil {
		olds, err = plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
			Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true})
		if err != nil {
			return nil, err
		}
	}

	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.news", label), KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	// Now fetch the default values so that (a) we can return them to the caller and (b) so that validation
	// includes the default values.  Otherwise, the provider wouldn't be presented with its own defaults.
	tfname := res.TF.Name
	assets := make(AssetTable)
	inputs, err := MakeTerraformInputs(
		&PulumiResource{URN: urn, Properties: news},
		olds, news, res.TFSchema, res.Schema.Fields, assets, true, false)
	if err != nil {
		return nil, err
	}

	// Now check with the resource provider to see if the values pass muster.
	rescfg, err := MakeTerraformConfigFromInputs(inputs)
	if err != nil {
		return nil, err
	}
	warns, errs := p.tf.ValidateResource(tfname, rescfg)
	for _, warn := range warns {
		if err = p.host.Log(ctx, diag.Warning, fmt.Sprintf("%v verification warning: %v", urn, warn)); err != nil {
			return nil, err
		}
	}

	// Now produce a return value of any properties that failed verification.
	var failures []*pulumirpc.CheckFailure
	for _, err := range errs {
		failures = append(failures, &pulumirpc.CheckFailure{
			Reason: err.Error(),
		})
	}

	// After all is said and done, we need to go back and return only what got populated as a diff from the origin.
	pinputs := MakeTerraformOutputs(inputs, res.TFSchema, res.Schema.Fields, assets, false)
	minputs, err := plugin.MarshalProperties(pinputs, plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.inputs", label), KeepUnknowns: true})
	if err != nil {
		return nil, err
	}

	return &pulumirpc.CheckResponse{Inputs: minputs, Failures: failures}, nil
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (p *Provider) Diff(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	p.setLoggingContext(ctx)
	urn := resource.URN(req.GetUrn())
	t := urn.Type()
	res, has := p.resources[t]
	if !has {
		return nil, errors.Errorf("unrecognized resource type (Diff): %s", urn)
	}

	label := fmt.Sprintf("%s.Diff(%s/%s)", p.label(), urn, res.TF.Name)
	glog.V(9).Infof("%s executing", label)

	// To figure out if we have a replacement, perform the diff and then look for RequiresNew flags.
	inputs, err := MakeTerraformAttributesFromRPC(
		nil, req.GetOlds(), res.TFSchema, res.Schema.Fields, false, false, fmt.Sprintf("%s.olds", label))
	if err != nil {
		return nil, errors.Wrapf(err, "preparing %s's old property state", urn)
	}
	info := &terraform.InstanceInfo{Type: res.TF.Name}
	state := &terraform.InstanceState{ID: req.GetId(), Attributes: inputs}
	config, err := MakeTerraformConfigFromRPC(
		nil, req.GetNews(), res.TFSchema, res.Schema.Fields, true, false, fmt.Sprintf("%s.news", label))
	if err != nil {
		return nil, errors.Wrapf(err, "preparing %s's new property state", urn)
	}
	diff, err := p.tf.Diff(info, state, config)
	if err != nil {
		return nil, errors.Wrapf(err, "diffing %s", urn)
	}

	// Each RequiresNew translates into a replacement.
	var replaces []string
	replaced := make(map[resource.PropertyKey]bool)
	if diff != nil {
		for k, attr := range diff.Attributes {
			if attr.RequiresNew {
				name, _, _ := getInfoFromTerraformName(k, res.TFSchema, res.Schema.Fields, false)
				replaces = append(replaces, string(name))
				replaced[name] = true
			}
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

	return &pulumirpc.DiffResponse{
		Replaces:            replaces,
		Stables:             stables,
		DeleteBeforeReplace: len(replaces) > 0 && res.Schema.DeleteBeforeReplace,
	}, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.  (The input ID
// must be blank.)  If this call fails, the resource must not have been created (i.e., it is "transacational").
func (p *Provider) Create(ctx context.Context, req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	p.setLoggingContext(ctx)
	urn := resource.URN(req.GetUrn())
	t := urn.Type()
	res, has := p.resources[t]
	if !has {
		return nil, errors.Errorf("unrecognized resource type (Create): %s", t)
	}

	label := fmt.Sprintf("%s.Create(%s/%s)", p.label(), urn, res.TF.Name)
	glog.V(9).Infof("%s executing", label)

	// To get Terraform to create a new resource, the ID msut be blank and existing state must be empty (since the
	// resource does not exist yet), and the diff object should have no old state and all of the new state.
	info := &terraform.InstanceInfo{Type: res.TF.Name}
	state, diff, err := MakeTerraformDiffFromRPC(nil, req.GetProperties(), res.TFSchema, res.Schema.Fields)
	if err != nil {
		return nil, errors.Wrapf(err, "preparing %s's property state")
	}
	newstate, err := p.tf.Apply(info, state, diff)
	if err != nil {
		return nil, errors.Wrapf(err, "creating %s", urn)
	}
	contract.Assertf(newstate != nil, "expected non-nil TF state during Create; required to obtain ID")

	// Create the ID and property maps and return them.
	props := MakeTerraformResult(newstate, res.TFSchema, res.Schema.Fields)
	mprops, err := plugin.MarshalProperties(props, plugin.MarshalOptions{Label: fmt.Sprintf("%s.outs", label)})
	if err != nil {
		return nil, err
	}
	return &pulumirpc.CreateResponse{Id: newstate.ID, Properties: mprops}, nil
}

// Update updates an existing resource with new values.  Only those values in the provided property bag are updated
// to new values.  The resource ID is returned and may be different if the resource had to be recreated.
func (p *Provider) Update(ctx context.Context, req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	p.setLoggingContext(ctx)
	urn := resource.URN(req.GetUrn())
	t := urn.Type()
	res, has := p.resources[t]
	if !has {
		return nil, errors.Errorf("unrecognized resource type (Update): %s", t)
	}

	label := fmt.Sprintf("%s.Update(%s/%s)", p.label(), urn, res.TF.Name)
	glog.V(9).Infof("%s executing", label)

	// Create a state state with the ID to update, a diff with old and new states, and perform the apply.
	state, diff, err := MakeTerraformDiffFromRPC(
		req.GetOlds(), req.GetNews(), res.TFSchema, res.Schema.Fields)
	if err != nil {
		return nil, errors.Wrapf(err, "preparing %s's old and new state/diffs", urn)
	}
	state.ID = req.GetId()

	info := &terraform.InstanceInfo{Type: res.TF.Name}
	newstate, err := p.tf.Apply(info, state, diff)
	if err != nil {
		return nil, errors.Wrapf(err, "updating %s", urn)
	}

	props := MakeTerraformResult(newstate, res.TFSchema, res.Schema.Fields)
	mprops, err := plugin.MarshalProperties(props, plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.outs", label)})
	if err != nil {
		return nil, err
	}
	return &pulumirpc.UpdateResponse{Properties: mprops}, nil
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed to still exist.
func (p *Provider) Delete(ctx context.Context, req *pulumirpc.DeleteRequest) (*pbempty.Empty, error) {
	p.setLoggingContext(ctx)
	urn := resource.URN(req.GetUrn())
	t := urn.Type()
	res, has := p.resources[t]
	if !has {
		return nil, errors.Errorf("unrecognized resource type (Delete): %s", t)
	}

	label := fmt.Sprintf("%s.Delete(%s/%s)", p.label(), urn, res.TF.Name)
	glog.V(9).Infof("%s executing", label)

	// Fetch the resource attributes since many providers need more than just the ID to perform the delete.
	attrs, err := MakeTerraformAttributesFromRPC(
		nil, req.GetProperties(), res.TFSchema, res.Schema.Fields, false, false, label)
	if err != nil {
		return nil, err
	}

	// Create a new state, with no diff, that is missing an ID.  Terraform will interpret this as a create operation.
	info := &terraform.InstanceInfo{Type: res.TF.Name}
	state := &terraform.InstanceState{ID: req.GetId(), Attributes: attrs}
	if _, err := p.tf.Apply(info, state, &terraform.InstanceDiff{Destroy: true}); err != nil {
		return nil, errors.Wrapf(err, "deleting %s", urn)
	}
	return &pbempty.Empty{}, nil
}

// Invoke dynamically executes a built-in function in the provider.
func (p *Provider) Invoke(ctx context.Context, req *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	p.setLoggingContext(ctx)
	tok := tokens.ModuleMember(req.GetTok())
	ds, has := p.dataSources[tok]
	if !has {
		return nil, errors.Errorf("unrecognized data function (Invoke): %s", tok)
	}

	label := fmt.Sprintf("%s.Invoke(%s)", p.label(), tok)
	glog.V(9).Infof("%s executing", label)

	// Unmarshal the arguments.
	args, err := plugin.UnmarshalProperties(req.GetArgs(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.args", label), KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	// First, create the inputs.
	tfname := ds.TF.Name
	inputs, err := MakeTerraformInputs(
		&PulumiResource{Properties: args}, nil, args, ds.TFSchema, ds.Schema.Fields, nil, true, false)
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
		if err = p.host.Log(ctx, diag.Warning, fmt.Sprintf("%v verification warning: %v", tok, warn)); err != nil {
			return nil, err
		}
	}

	// Now produce a return value of any properties that failed verification.
	var failures []*pulumirpc.CheckFailure
	for _, err := range errs {
		failures = append(failures, &pulumirpc.CheckFailure{
			Reason: err.Error(),
		})
	}

	// If there are no failures in verification, go ahead and perform the invocation.
	var ret *pbstruct.Struct
	if len(failures) == 0 {
		diff, err := p.tf.ReadDataDiff(info, rescfg)
		if err != nil {
			return nil, errors.Wrapf(err, "reading data source diff for %s", tok)
		}

		invoke, err := p.tf.ReadDataApply(info, diff)
		if err != nil {
			return nil, errors.Wrapf(err, "invoking %s", tok)
		}
		ret, err = plugin.MarshalProperties(
			MakeTerraformResult(invoke, ds.TFSchema, ds.Schema.Fields),
			plugin.MarshalOptions{Label: fmt.Sprintf("%s.returns", label)})
		if err != nil {
			return nil, err
		}
	}

	return &pulumirpc.InvokeResponse{
		Return:   ret,
		Failures: failures,
	}, nil
}

// GetPluginInfo implements an RPC call that returns the version of this plugin.
func (p *Provider) GetPluginInfo(ctx context.Context, req *pbempty.Empty) (*pulumirpc.PluginInfo, error) {
	return &pulumirpc.PluginInfo{
		Version: p.version,
	}, nil
}
