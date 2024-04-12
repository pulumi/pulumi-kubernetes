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
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/internals"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
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
	ResourcePrefix   pulumi.StringInput         `pulumi:"resourcePrefix,optional"`
	SkipAwait        pulumi.BoolInput           `pulumi:"skipAwait,optional"`
	CreateNamespace  pulumi.BoolInput           `pulumi:"createNamespace,optional"`
}

type chartArgs struct {
	Name             string
	Namespace        string
	Chart            string
	Version          string
	Devel            bool
	RepositoryOpts   helmv3.RepositoryOpts
	DependencyUpdate bool
	ResourcePrefix   *string
	SkipAwait        bool
	CreateNamespace  bool
}

func unwrapChartArgs(ctx context.Context, args *ChartArgs) (*chartArgs, internals.UnsafeAwaitOutputResult, error) {
	result, err := internals.UnsafeAwaitOutput(ctx, pulumi.All(
		args.Name, args.Namespace,
		args.Chart, args.Version, args.Devel, args.RepositoryOpts, args.DependencyUpdate,
		args.ResourcePrefix, args.SkipAwait, args.CreateNamespace))
	if err != nil {
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
	if v, ok := pop().(string); ok {
		r.ResourcePrefix = &v
	}
	r.SkipAwait, _ = pop().(bool)
	r.CreateNamespace, _ = pop().(bool)

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

	// Check if all the required args are known, and print a warning if not.
	result, err := internals.UnsafeAwaitOutput(ctx.Context(), pulumi.All(
		args.Name, args.Namespace,
		args.Chart, args.Version, args.Devel, args.RepositoryOpts,
		args.ResourcePrefix, args.SkipAwait, args.CreateNamespace))
	if err != nil {
		return nil, err
	}
	if !result.Known {
		msg := fmt.Sprintf("%s:%s -- Required input properties have unknown values. Preview is incomplete.\n", typ, name)
		_ = ctx.Log.Warn(msg, nil)
	}

	// Unpack the resolved inputs.
	chartArgs, result, err := unwrapChartArgs(ctx.Context(), args)
	if err != nil {
		return nil, err
	}
	if !result.Known {
		msg := fmt.Sprintf("%s:%s -- Required input properties have unknown values. Preview is incomplete.\n", typ, name)
		_ = ctx.Log.Warn(msg, nil)
		return pulumiprovider.NewConstructResult(comp)
	}
	if chartArgs.Name == "" {
		chartArgs.Name = name
	}
	if chartArgs.ResourcePrefix == nil {
		// use the name of the ConfigFile as the resource prefix to ensure uniqueness
		// across multiple instances of the component resource.
		chartArgs.ResourcePrefix = &name
	}

	// Prepare the `helm template` command
	tool := kubehelm.NewTool(r.opts.HelmOptions.EnvSettings)
	tool.HelmDriver = r.opts.HelmOptions.HelmDriver
	cmd := tool.Template()
	cmd.ReleaseName = chartArgs.Name
	cmd.Namespace = chartArgs.Namespace
	cmd.Chart = chartArgs.Chart
	cmd.Version = chartArgs.Version
	cmd.Devel = chartArgs.Devel
	cmd.DependencyUpdate = chartArgs.DependencyUpdate
	cmd.IncludeCRDs = true
	kubehelm.ApplyRepositoryOpts(&cmd.ChartPathOptions, chartArgs.RepositoryOpts)

	// Execute the Helm command
	release, err := cmd.Execute(ctx.Context())
	if err != nil {
		return nil, err
	}

	// Parse the YAML file into an array of Kubernetes objects.
	parseOpts := provideryamlv2.ParseOptions{
		YAML: release.Manifest,
	}
	objs, err := provideryamlv2.Parse(ctx.Context(), parseOpts)
	if err != nil {
		return nil, err
	}

	// https://github.com/helm/helm/issues/9813
	if chartArgs.CreateNamespace && chartArgs.Namespace != "" {
		ns := unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Namespace",
				"metadata": map[string]interface{}{
					"name": chartArgs.Namespace,
				},
			},
		}
		objs = append([]unstructured.Unstructured{ns}, objs...)
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
	}
	resources, err := provideryamlv2.Register(ctx, registerOpts)
	if err != nil {
		return nil, err
	}
	comp.Resources = resources

	return pulumiprovider.NewConstructResult(comp)
}
