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

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	providerresource "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/internals"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
)

type ConfigFileProvider struct {
	clientSet *clients.DynamicClientSet
}

type ConfigFileArgs struct {
	File           pulumi.StringInput `pulumi:"file"`
	ResourcePrefix pulumi.StringInput `pulumi:"resourcePrefix,optional"`
	SkipAwait      pulumi.BoolInput   `pulumi:"skipAwait,optional"`
}

type ConfigFileState struct {
	pulumi.ResourceState
	Resources pulumi.ArrayOutput `pulumi:"resources"`
}

var _ providerresource.ResourceProvider = &ConfigFileProvider{}

func NewConfigFileProvider(opts *providerresource.ResourceProviderOptions) providerresource.ResourceProvider {
	return &ConfigFileProvider{
		clientSet: opts.ClientSet,
	}
}

func (k *ConfigFileProvider) Construct(ctx *pulumi.Context, typ, name string, inputs pulumiprovider.ConstructInputs, options pulumi.ResourceOption) (*pulumiprovider.ConstructResult, error) {
	comp := &ConfigFileState{}
	err := ctx.RegisterComponentResource(typ, name, comp, options)
	if err != nil {
		return nil, err
	}

	args := &ConfigFileArgs{}
	if err := inputs.CopyTo(args); err != nil {
		return nil, fmt.Errorf("setting args: %w", err)
	}

	// Check if all the required args have resolved, and print a warning if not.
	result, err := internals.UnsafeAwaitOutput(ctx.Context(), pulumi.All(
		args.File, args.ResourcePrefix, args.SkipAwait))
	if err != nil {
		return nil, err
	}
	if !result.Known {
		msg := fmt.Sprintf("%s:%s -- Required input properties have unknown values. Preview is incomplete.\n", typ, name)
		_ = ctx.Log.Warn(msg, nil)
	}

	// Parse the manifest(s) and register the resources.

	comp.Resources = pulumi.All(args.File, args.ResourcePrefix, args.SkipAwait).ApplyTWithContext(ctx.Context(), func(_ context.Context, args []any) (pulumi.ArrayOutput, error) {
		// make type assertions to get each value (or the zero value)
		file, _ := args[0].(string)
		resourcePrefix, hasResourcePrefix := args[1].(string)
		skipAwait, _ := args[2].(bool)

		if !hasResourcePrefix {
			// use the name of the ConfigFile as the resource prefix to ensure uniqueness
			// across multiple instances of the component resource.
			resourcePrefix = name
		}

		return ParseDecodeYamlFiles(ctx, &ParseArgs{
			Files:          []string{file},
			ResourcePrefix: resourcePrefix,
			SkipAwait:      skipAwait,
		}, false, k.clientSet, pulumi.Parent(comp))
	}).(pulumi.ArrayOutput)

	// issue: https://github.com/pulumi/pulumi/issues/15527
	_, _ = internals.UnsafeAwaitOutput(ctx.Context(), comp.Resources)

	return pulumiprovider.NewConstructResult(comp)
}
