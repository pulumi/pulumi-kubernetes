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
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/internals"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	providerresource "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/resource"
)

type ConfigGroupProvider struct {
	clientSet        *clients.DynamicClientSet
	defaultNamespace string
}

type ConfigGroupArgs struct {
	Files          pulumi.StringArrayInput `pulumi:"files"`
	YAML           pulumi.StringInput      `pulumi:"yaml,optional"`
	Objects        pulumi.MapArrayInput    `pulumi:"objs,optional"`
	ResourcePrefix pulumi.StringInput      `pulumi:"resourcePrefix,optional"`
	SkipAwait      pulumi.BoolInput        `pulumi:"skipAwait,optional"`
}

type ConfigGroupState struct {
	pulumi.ResourceState
	Resources pulumi.ArrayOutput `pulumi:"resources"`
}

var _ providerresource.ResourceProvider = &ConfigGroupProvider{}

func NewConfigGroupProvider(opts *providerresource.ResourceProviderOptions) providerresource.ResourceProvider {
	return &ConfigGroupProvider{
		clientSet:        opts.ClientSet,
		defaultNamespace: opts.DefaultNamespace,
	}
}

func (k *ConfigGroupProvider) Construct(
	ctx *pulumi.Context,
	typ, name string,
	inputs pulumiprovider.ConstructInputs,
	options pulumi.ResourceOption,
) (*pulumiprovider.ConstructResult, error) {
	comp := &ConfigGroupState{}
	err := ctx.RegisterComponentResource(typ, name, comp, options)
	if err != nil {
		return nil, err
	}

	args := &ConfigGroupArgs{}
	if err := inputs.CopyTo(args); err != nil {
		return nil, fmt.Errorf("setting args: %w", err)
	}

	// Check if all the required args are known, and print a warning if not.
	result, err := internals.UnsafeAwaitOutput(ctx.Context(), pulumi.All(
		args.Files, args.YAML, args.Objects, args.ResourcePrefix, args.SkipAwait))
	if err != nil {
		return nil, err
	}
	if !result.Known {
		msg := fmt.Sprintf(
			"%s:%s -- Required input properties have unknown values. Preview is incomplete.\n",
			typ,
			name,
		)
		_ = ctx.Log.Warn(msg, nil)
	}

	// Parse the manifest(s) and register the resources.

	comp.Resources = pulumi.All(args.Files, args.YAML, args.Objects, args.ResourcePrefix, args.SkipAwait).
		ApplyTWithContext(ctx.Context(), func(_ context.Context, args []any) (pulumi.ArrayOutput, error) {
			// make type assertions to get each value (or the zero value)
			// note: "objects" contains unwrapped values at this point
			files, _ := args[0].([]string)
			yaml, _ := args[1].(string)
			objects, _ := args[2].([]map[string]any)
			resourcePrefix, hasResourcePrefix := args[3].(string)
			skipAwait, _ := args[4].(bool)

			if !hasResourcePrefix {
				// use the name of the ConfigGroup as the resource prefix to ensure uniqueness
				// across multiple instances of the component resource.
				resourcePrefix = name
			}

			// Parse the YAML files and literals into an array of Kubernetes objects, plus the provided objects.
			parseOpts := ParseOptions{
				Files: files,
				Glob:  true,
				YAML:  yaml,
			}
			objs, err := Parse(ctx.Context(), parseOpts)
			if err != nil {
				return pulumi.ArrayOutput{}, err
			}
			for _, obj := range objects {
				objs = append(objs, unstructured.Unstructured{Object: obj})
			}

			// Normalize the objects (apply a default namespace, etc.)
			objs, err = Normalize(objs, k.defaultNamespace, k.clientSet)
			if err != nil {
				return pulumi.ArrayOutput{}, err
			}

			// Register the objects as Pulumi resources.
			registerOpts := RegisterOptions{
				Objects:         objs,
				ResourcePrefix:  resourcePrefix,
				SkipAwait:       skipAwait,
				ResourceOptions: []pulumi.ResourceOption{pulumi.Parent(comp)},
			}
			return Register(ctx, registerOpts)

		}).(pulumi.ArrayOutput)

	// issue: https://github.com/pulumi/pulumi/issues/15527
	_, _ = internals.UnsafeAwaitOutput(ctx.Context(), comp.Resources)

	return pulumiprovider.NewConstructResult(comp)
}
