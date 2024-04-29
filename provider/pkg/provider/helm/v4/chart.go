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

	kubehelm "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/helm"
	providerresource "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/resource"
	provideryamlv2 "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/yaml/v2"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/internals"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
	helmkube "helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/postrender"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ChartProvider struct {
	opts *providerresource.ResourceProviderOptions
}

type ChartArgs struct {
	Name             pulumi.StringInput         `pulumi:"name"`
	Namespace        pulumi.StringInput         `pulumi:"namespace"`
	Chart            pulumi.StringInput         `pulumi:"chart"`
	Version          pulumi.StringInput         `pulumi:"version"`
	Devel            pulumi.BoolInput           `pulumi:"devel,optional"`
	RepositoryOpts   helmv3.RepositoryOptsInput `pulumi:"repositoryOpts,optional"`
	DependencyUpdate pulumi.BoolInput           `pulumi:"dependencyUpdate,optional"`

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
	RepositoryOpts   helmv3.RepositoryOpts
	DependencyUpdate bool

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
		args.Chart, args.Version, args.Devel, args.RepositoryOpts, args.DependencyUpdate,
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
	r.RepositoryOpts, _ = pop().(helmv3.RepositoryOpts)
	r.DependencyUpdate, _ = pop().(bool)

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
	tool := kubehelm.NewTool(r.opts.HelmOptions.EnvSettings)
	tool.HelmDriver = r.opts.HelmOptions.HelmDriver
	cmd := tool.Template()

	cmd.Validate = true
	cmd.DryRun = true
	cmd.DryRunOption = "server"

	cmd.Chart = chartArgs.Chart
	cmd.Version = chartArgs.Version
	cmd.Devel = chartArgs.Devel
	cmd.DependencyUpdate = chartArgs.DependencyUpdate
	kubehelm.ApplyRepositoryOpts(&cmd.ChartPathOptions, chartArgs.RepositoryOpts)

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
		PreRegisterF:    preregister,
	}
	resources, err := provideryamlv2.Register(ctx, registerOpts)
	if err != nil {
		return nil, err
	}
	comp.Resources = resources

	return pulumiprovider.NewConstructResult(comp)
}

func preregister(ctx *pulumi.Context, apiVersion, kind, resourceName string, obj *unstructured.Unstructured,
	resourceOpts []pulumi.ResourceOption) (*unstructured.Unstructured, []pulumi.ResourceOption) {

	// Implement support for Helm resource policies.
	// https://helm.sh/docs/howto/charts_tips_and_tricks/#tell-helm-not-to-uninstall-a-resource
	policy, hasPolicy, err := unstructured.NestedString(obj.Object, "metadata", "annotations", helmkube.ResourcePolicyAnno)
	if err == nil && hasPolicy {
		switch policy {
		case helmkube.KeepPolicy:
			resourceOpts = append(resourceOpts, pulumi.RetainOnDelete(true))
		}
	}

	return obj, resourceOpts
}
