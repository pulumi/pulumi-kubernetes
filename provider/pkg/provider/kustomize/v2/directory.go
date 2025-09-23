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

package v2

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	providerresource "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/resource"
	provideryamlv2 "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/internals"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/krusty/localizer"
	kresmap "sigs.k8s.io/kustomize/api/resmap"
	ktypes "sigs.k8s.io/kustomize/api/types"
	kfilesys "sigs.k8s.io/kustomize/kyaml/filesys"
)

type kustomizer interface {
	Run(fSys kfilesys.FileSystem, path string) (kresmap.ResMap, error)
}

type DirectoryProvider struct {
	opts           *providerresource.ResourceProviderOptions
	makeKustomizer func(args *directoryArgs) kustomizer
}

type DirectoryArgs struct {
	Directory      pulumi.StringInput `pulumi:"directory"`
	Namespace      pulumi.StringInput `pulumi:"namespace,optional"`
	ResourcePrefix pulumi.StringInput `pulumi:"resourcePrefix,optional"`
	SkipAwait      pulumi.BoolInput   `pulumi:"skipAwait,optional"`
}

type directoryArgs struct {
	Directory      string
	Namespace      string
	ResourcePrefix *string
	SkipAwait      bool
}

func unwrapDirectoryArgs(ctx context.Context, args *DirectoryArgs) (*directoryArgs, internals.UnsafeAwaitOutputResult, error) {
	result, err := internals.UnsafeAwaitOutput(ctx, pulumi.All(
		args.Directory, args.Namespace, args.ResourcePrefix, args.SkipAwait))
	if err != nil || !result.Known {
		return nil, result, err
	}
	resolved := result.Value.([]any)
	pop := func() (r any) {
		r, resolved = resolved[0], resolved[1:]
		return
	}

	r := &directoryArgs{}
	r.Directory, _ = pop().(string)
	r.Namespace, _ = pop().(string)
	if v, ok := pop().(string); ok {
		r.ResourcePrefix = &v
	}
	r.SkipAwait, _ = pop().(bool)

	return r, result, nil
}

type DirectoryState struct {
	pulumi.ResourceState
	Resources pulumi.ArrayOutput `pulumi:"resources"`
}

var _ providerresource.ResourceProvider = &DirectoryProvider{}

func NewDirectoryProvider(opts *providerresource.ResourceProviderOptions) providerresource.ResourceProvider {
	return &DirectoryProvider{
		opts:           opts,
		makeKustomizer: makeKustomizer,
	}
}

func (r *DirectoryProvider) Construct(ctx *pulumi.Context, typ, name string, inputs pulumiprovider.ConstructInputs, options pulumi.ResourceOption) (*pulumiprovider.ConstructResult, error) {
	comp := &DirectoryState{}
	err := ctx.RegisterComponentResource(typ, name, comp, options)
	if err != nil {
		return nil, err
	}

	args := &DirectoryArgs{}
	if err := inputs.CopyTo(args); err != nil {
		return nil, fmt.Errorf("setting args: %w", err)
	}

	// Unpack the resolved inputs.
	directoryArgs, result, err := unwrapDirectoryArgs(ctx.Context(), args)
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
	if directoryArgs.ResourcePrefix == nil {
		// use the name of the Directory as the resource prefix to ensure uniqueness
		// across multiple instances of the component resource.
		directoryArgs.ResourcePrefix = &name
	}

	fs := kfilesys.MakeFsOnDisk()

	source := directoryArgs.Directory

	// If the directory is remote, localize it and then kustomize the local copy.
	if url, err := url.Parse(source); err == nil && url.Scheme != "" {
		target, err := os.MkdirTemp("", "pulumi-kustomize")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp directory: %w", err)
		}
		defer func() { _ = (os.RemoveAll(target)) }()

		source, err = localizer.Run(fs, source, "", filepath.Join(target, "local"))
		if err != nil {
			return nil, fmt.Errorf("kustomize localization failed: %w", err)
		}
	}

	if _, err := os.Stat(source); errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%q: %w", source, err)
	}

	// Execute the kustomize command to generate the Kubernetes manifest.
	k := r.makeKustomizer(directoryArgs)
	rm, err := k.Run(fs, source)
	if err != nil {
		return nil, fmt.Errorf("kustomize build error: %w", err)
	}
	manifest, err := rm.AsYaml()
	if err != nil {
		return nil, fmt.Errorf("failed to convert kustomize result to YAML: %w", err)
	}

	// Parse the YAML file into an array of Kubernetes objects.
	parseOpts := provideryamlv2.ParseOptions{
		YAML: string(manifest),
	}
	objs, err := provideryamlv2.Parse(ctx.Context(), parseOpts)
	if err != nil {
		return nil, err
	}

	// Normalize the objects (apply a default namespace, etc.)
	ns := r.opts.DefaultNamespace
	if directoryArgs.Namespace != "" {
		ns = directoryArgs.Namespace
	}
	objs, err = provideryamlv2.Normalize(objs, ns, r.opts.ClientSet)
	if err != nil {
		return nil, err
	}

	// Register the objects as Pulumi resources.
	registerOpts := provideryamlv2.RegisterOptions{
		Objects:         objs,
		ResourcePrefix:  *directoryArgs.ResourcePrefix,
		SkipAwait:       directoryArgs.SkipAwait,
		ResourceOptions: []pulumi.ResourceOption{pulumi.Parent(comp)},
	}
	resources, err := provideryamlv2.Register(ctx, registerOpts)
	if err != nil {
		return nil, err
	}
	comp.Resources = resources

	return pulumiprovider.NewConstructResult(comp)
}

// makeKustomizer prepares the kustomize tool with helm support, full permission to use plugins, and no load restrictions
func makeKustomizer(args *directoryArgs) kustomizer {
	opts := krusty.MakeDefaultOptions()
	opts.Reorder = krusty.ReorderOptionNone
	opts.AddManagedbyLabel = false
	opts.LoadRestrictions = ktypes.LoadRestrictionsNone
	opts.PluginConfig = ktypes.EnabledPluginConfig(ktypes.BploUseStaticallyLinked)
	opts.PluginConfig.HelmConfig.Enabled = true
	opts.PluginConfig.HelmConfig.Command = "helm"
	k := krusty.MakeKustomizer(opts)
	return k
}
