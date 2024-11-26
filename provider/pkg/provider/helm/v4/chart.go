// Copyright 2016-2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v4

import (
	"context"
	"fmt"

	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	"github.com/pulumi/pulumi-kubernetes/v4/provider/pkg/clients"
	kubehelm "github.com/pulumi/pulumi-kubernetes/v4/provider/pkg/helm"
	providerresource "github.com/pulumi/pulumi-kubernetes/v4/provider/pkg/provider/resource"
	provideryamlv2 "github.com/pulumi/pulumi-kubernetes/v4/provider/pkg/provider/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/internals"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	helmkube "helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/postrender"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
)

type toolF func() *kubehelm.Tool

type ChartProvider struct {
	opts *providerresource.ResourceProviderOptions
	tool toolF
}

type ChartArgs struct {
	Name             pulumi.StringInput         `pulumi:"name,optional"`
	Namespace        pulumi.StringInput         `pulumi:"namespace,optional"`
	Chart            pulumi.StringInput         `pulumi:"chart"`
	Version          pulumi.StringInput         `pulumi:"version,optional"`
	Devel            pulumi.BoolInput           `pulumi:"devel,optional"`
	RepositoryOpts   helmv4.RepositoryOptsInput `pulumi:"repositoryOpts,optional"`
	DependencyUpdate pulumi.BoolInput           `pulumi:"dependencyUpdate,optional"`
	Verify           pulumi.BoolInput           `pulumi:"verify,optional"`
	Keyring          pulumi.AssetInput          `pulumi:"keyring,optional"`

	Values       pulumi.MapInput          `pulumi:"values,optional"`
	ValuesFiles  pulumi.AssetArrayInput   `pulumi:"valueYamlFiles,optional"`
	SkipCrds     pulumi.BoolInput         `pulumi:"skipCrds,optional"`
	PostRenderer helmv4.PostRendererInput `pulumi:"postRenderer,optional"`

	ResourcePrefix pulumi.StringInput `pulumi:"resourcePrefix,optional"`
	SkipAwait      pulumi.BoolInput   `pulumi:"skipAwait,optional"`
}

type chartArgs struct {
	Name             string
	Namespace        string
	Chart            string
	Version          string
	Devel            bool
	RepositoryOpts   helmv4.RepositoryOpts
	DependencyUpdate bool
	Verify           bool
	Keyring          pulumi.Asset

	Values       map[string]any
	ValuesFiles  []pulumi.Asset
	SkipCrds     bool
	PostRenderer *helmv4.PostRenderer

	ResourcePrefix *string
	SkipAwait      bool
}

func unwrapChartArgs(ctx context.Context, args *ChartArgs) (*chartArgs, internals.UnsafeAwaitOutputResult, error) {
	result, err := internals.UnsafeAwaitOutput(ctx, pulumi.All(
		args.Name, args.Namespace,
		args.Chart, args.Version, args.Devel, args.RepositoryOpts, args.DependencyUpdate, args.Verify, args.Keyring,
		args.Values, args.ValuesFiles, args.SkipCrds, args.PostRenderer,
		args.ResourcePrefix, args.SkipAwait))
	if err != nil || !result.Known {
		return nil, result, err
	}
	resolved := result.Value.([]any)
	pop := func() (r any) {
		r, resolved = resolved[0], resolved[1:]
		return
	}

	r := &chartArgs{}
	r.Name, _ = pop().(string)
	r.Namespace, _ = pop().(string)
	r.Chart, _ = pop().(string)
	r.Version, _ = pop().(string)
	r.Devel, _ = pop().(bool)
	r.RepositoryOpts, _ = pop().(helmv4.RepositoryOpts)
	r.DependencyUpdate, _ = pop().(bool)
	r.Verify, _ = pop().(bool)
	r.Keyring, _ = pop().(pulumi.Asset)

	r.Values, _ = pop().(map[string]any)
	r.ValuesFiles, _ = pop().([]pulumi.Asset)
	r.SkipCrds, _ = pop().(bool)
	if v, ok := pop().(helmv4.PostRenderer); ok {
		r.PostRenderer = &v
	}

	if v, ok := pop().(string); ok {
		r.ResourcePrefix = &v
	}
	r.SkipAwait, _ = pop().(bool)

	return r, result, nil
}

type ChartState struct {
	pulumi.ResourceState
	Resources pulumi.ArrayOutput `pulumi:"resources"`
}

var _ providerresource.ResourceProvider = &ChartProvider{}

func NewChartProvider(opts *providerresource.ResourceProviderOptions) providerresource.ResourceProvider {
	return &ChartProvider{
		opts: opts,
		tool: func() *kubehelm.Tool {
			return kubehelm.NewTool(opts.HelmOptions.EnvSettings)
		},
	}
}

func (r *ChartProvider) Construct(ctx *pulumi.Context, typ, name string, inputs pulumiprovider.ConstructInputs, options pulumi.ResourceOption) (*pulumiprovider.ConstructResult, error) {
	comp := &ChartState{}
	err := ctx.RegisterComponentResource(typ, name, comp, options)
	if err != nil {
		return nil, err
	}

	args := &ChartArgs{}
	if err := inputs.CopyTo(args); err != nil {
		return nil, fmt.Errorf("setting args: %w", err)
	}

	// Unpack the resolved inputs.
	chartArgs, result, err := unwrapChartArgs(ctx.Context(), args)
	if err != nil {
		return nil, fmt.Errorf("unwrapping args: %w", err)
	}
	if !result.Known {
		_ = ctx.Log.Warn("Input properties have unknown values. Preview is incomplete.", &pulumi.LogArgs{
			Resource: comp,
		})
		r, err := pulumiprovider.NewConstructResult(comp)
		return r, err
	}
	if chartArgs.Name == "" {
		chartArgs.Name = name
	}
	if chartArgs.ResourcePrefix == nil {
		// use the name of the Chart as the resource prefix to ensure uniqueness
		// across multiple instances of the component resource.
		chartArgs.ResourcePrefix = &name
	}

	// Prepare the `helm template` command
	tool := r.tool()
	tool.HelmDriver = r.opts.HelmOptions.HelmDriver
	p := tool.AllGetters()
	cmd := tool.Template()

	// connectivity: use server-side dry-run to enable the lookup function.
	// i.e. `helm template --dry-run=server --validate=false`.
	cmd.DisableOpenAPIValidation = true
	cmd.Validate = false
	cmd.ClientOnly = true
	cmd.DryRun = true
	cmd.DryRunOption = "server"
	if r.opts.ClientSet.DiscoveryClientCached != nil {
		if err := setKubeVersionAndAPIVersions(r.opts.ClientSet, cmd); err != nil {
			return nil, err
		}
	}

	// set chart resolution options
	cmd.Chart = chartArgs.Chart
	cmd.Version = chartArgs.Version
	cmd.Devel = chartArgs.Devel
	if err = kubehelm.ApplyRepositoryOpts(&cmd.ChartPathOptions, p, chartArgs.RepositoryOpts); err != nil {
		return nil, fmt.Errorf("repositoryOpts: %w", err)
	}
	cmd.DependencyUpdate = chartArgs.DependencyUpdate
	cmd.Verify = chartArgs.Verify

	if chartArgs.Keyring != nil {
		keyring, err := kubehelm.LocateKeyring(p, chartArgs.Keyring)
		if err != nil {
			return nil, fmt.Errorf("keyring: %w", err)
		}
		cmd.Keyring = keyring
	}

	// set templating options
	cmd.Values.Values = chartArgs.Values
	cmd.Values.ValuesFiles = chartArgs.ValuesFiles
	cmd.IncludeCRDs = !chartArgs.SkipCrds
	cmd.DisableHooks = true
	cmd.ReleaseName = chartArgs.Name
	cmd.Namespace = chartArgs.Namespace

	if chartArgs.PostRenderer != nil {
		postrenderer, err := postrender.NewExec(chartArgs.PostRenderer.Command, chartArgs.PostRenderer.Args...)
		if err != nil {
			return nil, err
		}
		cmd.PostRenderer = postrenderer
	}

	// Execute the Helm command
	release, err := cmd.Execute(ctx.Context())
	if err != nil {
		return nil, err
	}

	if release.Chart.Metadata.Deprecated {
		_ = ctx.Log.Warn(fmt.Sprintf("Using a deprecated Helm chart (%s)", release.Chart.Name()), &pulumi.LogArgs{
			Resource: comp,
		})
	}

	// Parse the YAML file into an array of Kubernetes objects.
	parseOpts := provideryamlv2.ParseOptions{
		YAML: release.Manifest,
	}
	objs, err := provideryamlv2.Parse(ctx.Context(), parseOpts)
	if err != nil {
		return nil, err
	}

	// Normalize the objects (apply a default namespace, etc.)
	ns := chartArgs.Namespace
	if ns == "" {
		ns = r.opts.DefaultNamespace
	}
	objs, err = provideryamlv2.Normalize(objs, ns, r.opts.ClientSet)
	if err != nil {
		return nil, err
	}

	// Register the objects as Pulumi resources.
	registerOpts := provideryamlv2.RegisterOptions{
		Objects:         objs,
		ResourcePrefix:  *chartArgs.ResourcePrefix,
		SkipAwait:       chartArgs.SkipAwait,
		ResourceOptions: []pulumi.ResourceOption{pulumi.Parent(comp)},
		PreRegisterF: func(ctx *pulumi.Context, apiVersion, kind, resourceName string, obj *unstructured.Unstructured,
			resourceOpts []pulumi.ResourceOption) (*unstructured.Unstructured, []pulumi.ResourceOption) {
			return preregister(ctx, comp, obj, resourceOpts)
		},
	}
	resources, err := provideryamlv2.Register(ctx, registerOpts)
	if err != nil {
		return nil, err
	}
	comp.Resources = resources

	return pulumiprovider.NewConstructResult(comp)
}

func preregister(ctx *pulumi.Context, comp *ChartState, obj *unstructured.Unstructured,
	resourceOpts []pulumi.ResourceOption) (*unstructured.Unstructured, []pulumi.ResourceOption) {

	// Implement support for Helm resource policies.
	// https://helm.sh/docs/howto/charts_tips_and_tricks/#tell-helm-not-to-uninstall-a-resource
	policy, hasPolicy, err := unstructured.NestedString(obj.Object, "metadata", "annotations", helmkube.ResourcePolicyAnno)
	if err == nil && hasPolicy {
		switch policy {
		case helmkube.KeepPolicy:
			resourceOpts = append(resourceOpts, pulumi.RetainOnDelete(true))
		default:
			_ = ctx.Log.Warn(fmt.Sprintf("Unsupported Helm resource policy %q", policy), &pulumi.LogArgs{
				Resource: comp,
			})
		}
	}

	return obj, resourceOpts
}

func setKubeVersionAndAPIVersions(clientSet *clients.DynamicClientSet, cmd *kubehelm.TemplateCommand) error {
	dc := clientSet.DiscoveryClientCached

	// https://github.com/helm/helm/blob/635b8cf33d25a86131635c32f35b2a76256e40cb/pkg/action/action.go#L246-L285

	// force a discovery cache invalidation to always fetch the latest server version/capabilities.
	dc.Invalidate()
	kubeVersion, err := dc.ServerVersion()
	if err != nil {
		return fmt.Errorf("could not get server version from Kubernetes: %w", err)
	}
	cmd.Install.KubeVersion = &chartutil.KubeVersion{
		Version: kubeVersion.GitVersion,
		Major:   kubeVersion.Major,
		Minor:   kubeVersion.Minor,
	}

	// Client-Go emits an error when an API service is registered but unimplemented.
	// Since the discovery client continues building the API object, it is correctly
	// populated with all valid APIs.
	// See https://github.com/kubernetes/kubernetes/issues/72051#issuecomment-521157642
	apiVersions, err := action.GetVersionSet(dc)
	if err != nil {
		if !discovery.IsGroupDiscoveryFailedError(err) {
			return fmt.Errorf("could not get apiVersions from Kubernetes: %w", err)
		}
	}
	cmd.Install.APIVersions = apiVersions

	return nil
}
