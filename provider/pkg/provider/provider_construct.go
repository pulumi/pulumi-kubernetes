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

	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	providerhelmv4 "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/helm/v4"
	providerkustomizev2 "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/kustomize/v2"
	providerresource "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/resource"
	provideryamlv2 "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/yaml/v2"
)

// resourceProviders contains factories for component resource providers.
var resourceProviders = map[string]providerresource.ResourceProviderFactory{
	"kubernetes:yaml/v2:ConfigFile":     provideryamlv2.NewConfigFileProvider,
	"kubernetes:yaml/v2:ConfigGroup":    provideryamlv2.NewConfigGroupProvider,
	"kubernetes:helm.sh/v4:Chart":       providerhelmv4.NewChartProvider,
	"kubernetes:kustomize/v2:Directory": providerkustomizev2.NewDirectoryProvider,
}

// getResourceProvider returns the resource provider for the given type, if a factory for one is registered.
func (k *kubeProvider) getResourceProvider(typ string) (providerresource.ResourceProvider, bool) {
	providerF, found := k.resourceProviders[typ]
	if !found {
		return nil, false
	}

	// In yamlRenderMode, defaultNamespace might not be set if the cluster is unreachable.
	// Provide a default value in this case.
	defaultNamespace := k.defaultNamespace
	if defaultNamespace == "" && k.yamlRenderMode {
		defaultNamespace = canonicalNamespace(defaultNamespace)
	}

	options := &providerresource.ResourceProviderOptions{
		ClientSet:        k.clientSet,
		DefaultNamespace: defaultNamespace,
		HelmOptions: &providerresource.HelmOptions{
			SuppressHelmHookWarnings: k.suppressHelmHookWarnings,
			HelmDriver:               k.helmDriver,
			EnvSettings:              k.helmSettings,
		},
	}
	return providerF(options), true
}

// Construct creates a new instance of the provided component resource and returns its state.
func (k *kubeProvider) Construct(
	ctx context.Context,
	req *pulumirpc.ConstructRequest,
) (*pulumirpc.ConstructResponse, error) {

	if k.clusterUnreachable && !k.yamlRenderMode {
		return nil, fmt.Errorf("configured Kubernetes cluster is unreachable: %s", k.clusterUnreachableReason)
	}
	// In yamlRenderMode we provide a default value for the default namespace.
	// In all other cases we need to assert a default namespace is set.
	if !k.yamlRenderMode {
		contract.Assertf(
			k.defaultNamespace != "" || k.yamlRenderMode,
			"expected defaultNamespace outside of render mode",
		)
	}
	contract.Assertf(k.helmDriver != "", "expected helmDriver")
	contract.Assertf(k.helmSettings != nil, "expected helmSettings")

	typ := req.GetType()
	provider, found := k.getResourceProvider(typ)
	if !found {
		return nil, fmt.Errorf("unknown resource type %q", typ)
	}
	return pulumiprovider.Construct(ctx, req, k.host.EngineConn(), provider.Construct)
}
