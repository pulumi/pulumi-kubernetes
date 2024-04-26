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

package provider

import (
	"context"
	"fmt"

	providerhelmv4 "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/helm/v4"
	providerresource "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/resource"
	provideryamlv2 "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/yaml/v2"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

// resourceProviders contains factories for component resource providers.
var resourceProviders = map[string]providerresource.ResourceProviderFactory{
	"kubernetes:yaml/v2:ConfigFile":  provideryamlv2.NewConfigFileProvider,
	"kubernetes:yaml/v2:ConfigGroup": provideryamlv2.NewConfigGroupProvider,
	"kubernetes:helm.sh/v4:Chart":    providerhelmv4.NewChartProvider,
}

// getResourceProvider returns the resource provider for the given type, if a factory for one is registered.
func (k *kubeProvider) getResourceProvider(typ string) (providerresource.ResourceProvider, bool) {
	providerF, found := k.resourceProviders[typ]
	if !found {
		return nil, false
	}

	options := &providerresource.ResourceProviderOptions{
		ClientSet:        k.clientSet,
		DefaultNamespace: k.defaultNamespace,
		HelmOptions: &providerresource.HelmOptions{
			SuppressHelmHookWarnings: k.suppressHelmHookWarnings,
			HelmDriver:               k.helmDriver,
			EnvSettings:              k.helmSettings,
		},
	}
	return providerF(options), true
}

// Construct creates a new instance of the provided component resource and returns its state.
func (k *kubeProvider) Construct(ctx context.Context, req *pulumirpc.ConstructRequest) (*pulumirpc.ConstructResponse, error) {

	if k.clusterUnreachable {
		return nil, fmt.Errorf("configured Kubernetes cluster is unreachable: %s", k.clusterUnreachableReason)
	}

	typ := req.GetType()
	provider, found := k.getResourceProvider(typ)
	if !found {
		return nil, fmt.Errorf("unknown resource type %q", typ)
	}
	return pulumiprovider.Construct(ctx, req, k.host.EngineConn(), provider.Construct)
}
