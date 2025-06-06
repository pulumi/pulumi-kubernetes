// Copyright 2016-2023, Pulumi Corporation.
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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/golang/protobuf/ptypes/empty"
	pbempty "github.com/golang/protobuf/ptypes/empty"
	structpb "github.com/golang/protobuf/ptypes/struct"
	checkjob "github.com/pulumi/cloud-ready-checks/pkg/kubernetes/job"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/cluster"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/gen"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/host"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/kinds"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/logging"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/metadata"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/openapi"
	providerresource "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/resource"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/ssa"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/version"
	pulumischema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/pkg/v3/resource/deploy/providers"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/rpcutil/rpcerror"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	helmcli "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/helmpath"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientapi "k8s.io/client-go/tools/clientcmd/api"
	k8sopenapi "k8s.io/kubectl/pkg/util/openapi"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"
)

// --------------------------------------------------------------------------

// Kubernetes resource provider.
//
// Implements functionality for the Pulumi Kubernetes Resource Provider. This code is responsible
// for producing sensible responses for the gRPC server to send back to a client when it requests
// something to do with the Kubernetes resources it's meant to manage.

// --------------------------------------------------------------------------

const (
	invokeDecodeYaml     = "kubernetes:yaml:decode"
	invokeHelmTemplate   = "kubernetes:helm:template"
	invokeKustomize      = "kubernetes:kustomize:directory"
	lastAppliedConfigKey = "kubectl.kubernetes.io/last-applied-configuration"
	initialAPIVersionKey = "__initialApiVersion"
	fieldManagerKey      = "__fieldManager"
	secretKind           = "Secret"
	clusterIdentifierKey = "clusterIdentifier"
)

type cancellationContext struct {
	context context.Context
	cancel  context.CancelFunc
}

func makeCancellationContext() *cancellationContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &cancellationContext{
		context: ctx,
		cancel:  cancel,
	}
}

type kubeOpts struct {
	rejectUnknownResources bool
}

type kubeProvider struct {
	pulumirpc.UnimplementedResourceProviderServer

	host             host.HostClient
	canceler         *cancellationContext
	name             string
	version          string
	pulumiSchema     []byte
	terraformMapping []byte
	providerPackage  string
	opts             kubeOpts
	defaultNamespace string

	deleteUnreachable           bool
	skipUpdateUnreachable       bool
	enableConfigMapMutable      bool
	enableSecretMutable         bool
	enableSecrets               bool
	suppressDeprecationWarnings bool
	suppressHelmHookWarnings    bool
	serverSideApplyMode         bool

	helmDriver               string
	helmPluginsPath          string
	helmRegistryConfigPath   string
	helmRepositoryConfigPath string
	helmRepositoryCache      string
	helmSettings             *helmcli.EnvSettings
	helmReleaseProvider      customResourceProvider

	yamlRenderMode bool
	yamlDirectory  string

	clusterUnreachable       bool   // Kubernetes cluster is unreachable.
	clusterUnreachableReason string // Detailed error message if cluster is unreachable.

	makeClient func(context.Context, *rest.Config) (*clients.DynamicClientSet, *clients.LogClient, error)
	clientSet  *clients.DynamicClientSet
	logClient  *clients.LogClient
	k8sVersion cluster.ServerVersion

	resources      k8sopenapi.Resources
	resourcesMutex sync.RWMutex

	resourceProviders map[string]providerresource.ResourceProviderFactory

	crdSchemas parameterizedPackageMap // In memory cache of CRD types from Parameterize calls.
}

var _ pulumirpc.ResourceProviderServer = (*kubeProvider)(nil)

func makeKubeProvider(
	host host.HostClient, name, version string, pulumiSchema, terraformMapping []byte,
) (*kubeProvider, error) {
	return &kubeProvider{
		host:                        host,
		canceler:                    makeCancellationContext(),
		name:                        name,
		version:                     version,
		pulumiSchema:                pulumiSchema,
		terraformMapping:            terraformMapping,
		providerPackage:             name,
		enableSecrets:               false,
		suppressDeprecationWarnings: false,
		deleteUnreachable:           false,
		skipUpdateUnreachable:       false,
		makeClient:                  makeClient,
		resourceProviders:           resourceProviders,
	}, nil
}

// makeClient makes a client to connect to a Kubernetes cluster using the given config.
// ctx is a cancellation context that may be used to cancel any subsequent requests made by the clients.
func makeClient(ctx context.Context, config *rest.Config) (*clients.DynamicClientSet, *clients.LogClient, error) {
	cs, err := clients.NewDynamicClientSet(config)
	if err != nil {
		return nil, nil, err
	}
	lc, err := clients.MakeLogClient(ctx, config)
	if err != nil {
		return nil, nil, err
	}
	return cs, lc, nil
}

func (k *kubeProvider) getResources() (k8sopenapi.Resources, error) {
	k.resourcesMutex.RLock()
	rs := k.resources
	k.resourcesMutex.RUnlock()

	if rs != nil {
		return rs, nil
	}

	k.resourcesMutex.Lock()
	defer k.resourcesMutex.Unlock()

	rs, err := openapi.GetResourceSchemasForClient(k.clientSet.DiscoveryClientCached)
	if err != nil {
		return nil, err
	}
	k.resources = rs
	return k.resources, nil
}

func (k *kubeProvider) invalidateResources() {
	k.resourcesMutex.Lock()
	defer k.resourcesMutex.Unlock()

	k.resources = nil
}

// Call dynamically executes a method in the provider associated with a component resource.
func (k *kubeProvider) Call(ctx context.Context, req *pulumirpc.CallRequest) (*pulumirpc.CallResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Call is not yet implemented")
}

// GetMapping fetches the mapping for this resource provider, if any. A provider should return an empty
// response (not an error) if it doesn't have a mapping for the given key.
func (k *kubeProvider) GetMapping(ctx context.Context, request *pulumirpc.GetMappingRequest) (*pulumirpc.GetMappingResponse, error) {
	// We only return a mapping for terraform
	if request.Key != "terraform" {
		// an empty response means no mapping, by design we don't return an error here
		return &pulumirpc.GetMappingResponse{}, nil
	}

	return &pulumirpc.GetMappingResponse{
		Provider: "kubernetes",
		Data:     k.terraformMapping,
	}, nil
}

// GetSchema returns the JSON-encoded schema for this provider's package.
func (k *kubeProvider) GetSchema(ctx context.Context, req *pulumirpc.GetSchemaRequest) (*pulumirpc.GetSchemaResponse, error) {
	if v := req.GetVersion(); v != 0 {
		return nil, fmt.Errorf("unsupported schema version %d", v)
	}

	if req.SubpackageName != "" && req.SubpackageVersion != "" {
		spec := k.crdSchemas.get(req.SubpackageName, req.SubpackageVersion)
		if spec == nil {
			return nil, fmt.Errorf("no schema found for %s", req.SubpackageName)
		}

		// Encode into JSON string.
		specJSON, err := json.Marshal(spec)
		if err != nil {
			return nil, err
		}

		return &pulumirpc.GetSchemaResponse{Schema: string(specJSON)}, nil
	}

	return &pulumirpc.GetSchemaResponse{Schema: string(k.pulumiSchema)}, nil
}

// CheckConfig validates the configuration for this provider.
func (k *kubeProvider) CheckConfig(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.CheckConfig(%s)", k.label(), urn)
	logger.V(9).Infof("%s executing", label)

	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.news", label),
		KeepUnknowns: true,
		SkipNulls:    true,
	})
	if err != nil {
		return nil, fmt.Errorf("CheckConfig failed because of malformed resource inputs: %w", err)
	}

	truthyValue := func(argName resource.PropertyKey, props resource.PropertyMap) bool {
		if arg := props[argName]; arg.HasValue() {
			switch {
			case arg.IsString() && len(arg.StringValue()) > 0:
				return true
			case arg.IsBool() && arg.BoolValue():
				return true
			default:
				return false
			}
		}
		return false
	}

	strictMode := false
	if pConfig, ok := k.loadPulumiConfig(); ok {
		if v, ok := pConfig["strictMode"]; ok {
			if v, ok := v.(string); ok {
				strictMode = v == "true"
			}
		}
	}
	if v := news["strictMode"]; v.HasValue() && v.IsString() {
		strictMode = v.StringValue() == "true"
	}

	if strictMode && providers.IsProviderType(urn.Type()) {
		var failures []*pulumirpc.CheckFailure

		if providers.IsDefaultProvider(urn) {
			failures = append(failures, &pulumirpc.CheckFailure{
				Reason: "strict mode prohibits default provider",
			})
		}
		if v := news["kubeconfig"]; !v.HasValue() || v.StringValue() == "" {
			failures = append(failures, &pulumirpc.CheckFailure{
				Property: "kubeconfig",
				Reason:   `strict mode requires Provider "kubeconfig" argument`,
			})
		}
		if v := news["context"]; !v.HasValue() || v.StringValue() == "" {
			failures = append(failures, &pulumirpc.CheckFailure{
				Property: "context",
				Reason:   `strict mode requires Provider "context" argument`,
			})
		}

		if len(failures) > 0 {
			return &pulumirpc.CheckResponse{Inputs: req.GetNews(), Failures: failures}, nil
		}
	}

	renderYamlEnabled := truthyValue("renderYamlToDirectory", news)

	errTemplate := `%q arg is not compatible with "renderYamlToDirectory" arg`
	if renderYamlEnabled {
		var failures []*pulumirpc.CheckFailure

		if truthyValue("cluster", news) {
			failures = append(failures, &pulumirpc.CheckFailure{
				Property: "cluster",
				Reason:   fmt.Sprintf(errTemplate, "cluster"),
			})
		}
		if truthyValue("context", news) {
			failures = append(failures, &pulumirpc.CheckFailure{
				Property: "context",
				Reason:   fmt.Sprintf(errTemplate, "context"),
			})
		}
		if truthyValue("kubeconfig", news) {
			failures = append(failures, &pulumirpc.CheckFailure{
				Property: "kubeconfig",
				Reason:   fmt.Sprintf(errTemplate, "kubeconfig"),
			})
		}

		if len(failures) > 0 {
			return &pulumirpc.CheckResponse{Inputs: req.GetNews(), Failures: failures}, nil
		}
	}

	return &pulumirpc.CheckResponse{Inputs: req.GetNews()}, nil
}

// DiffConfig diffs the configuration for this provider.
func (k *kubeProvider) DiffConfig(_ context.Context, req *pulumirpc.DiffRequest) (resp *pulumirpc.DiffResponse, err error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.DiffConfig(%s)", k.label(), urn)
	logger.V(9).Infof("%s executing", label)

	olds, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.olds", label),
		KeepUnknowns: true,
		SkipNulls:    true,
	})
	if err != nil {
		return nil, err
	}
	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.news", label),
		KeepUnknowns: true,
		SkipNulls:    true,
	})
	if err != nil {
		return nil, fmt.Errorf("DiffConfig failed because of malformed resource inputs: %w", err)
	}
	var diffs, replaces []string
	diff := olds.Diff(news)
	if diff == nil {
		return &pulumirpc.DiffResponse{
			Changes: pulumirpc.DiffResponse_DIFF_NONE,
		}, nil
	}

	// If clusterIdentifier is configured, we change our replacement behavior
	// to trigger when the identifier has changed. We don't do this when adding
	// or removing the identifier in order to make on- and off-boarding less
	// risky.
	if diff.Updated(clusterIdentifierKey) {
		replaces = append(replaces, clusterIdentifierKey)
	}
	// Similarly, we ignore replacements when the identifier is unchanged from
	// what it was previously set to.
	if diff.Same(clusterIdentifierKey) && news.HasValue(clusterIdentifierKey) {
		// Modify our response to no longer replace anything.
		defer func() {
			if resp != nil {
				resp.Replaces = nil
			}
		}()
	}

	// We can't tell for sure if a computed value has changed, so we make the conservative choice
	// and force a replacement. Note that getActiveClusterFromConfig relies on all three of the below properties.
	for _, key := range []resource.PropertyKey{"kubeconfig", "context", "cluster"} {
		if news[key].IsComputed() {
			replaces = append(replaces, string(key))
			resp = &pulumirpc.DiffResponse{
				Changes:  pulumirpc.DiffResponse_DIFF_SOME,
				Diffs:    []string{string(key)},
				Replaces: replaces,
			}
			return resp, nil
		}
	}

	oldConfig, err := parseKubeconfigPropertyValue(olds["kubeconfig"])
	if err != nil {
		return nil, err
	}
	newConfig, err := parseKubeconfigPropertyValue(news["kubeconfig"])
	if err != nil {
		return nil, err
	}

	// Check for differences in provider overrides.
	for _, k := range diff.ChangedKeys() {
		diffs = append(diffs, string(k))

		// Handle any special cases.
		switch k {
		case "renderYamlToDirectory":
			// If the render directory changes, all the manifests will be replaced.
			replaces = append(replaces, "renderYamlToDirectory")
		}
	}

	// In general, it's not possible to tell from a kubeconfig if the k8s cluster it points to has
	// changed. k8s clusters do not have a well defined identity, so the best we can do is check
	// if the settings for the active cluster have changed. This is not a foolproof method; a trivial
	// counterexample is changing the load balancer or DNS entry pointing to the same cluster.
	//
	// Given this limitation, we try to strike a reasonable balance by planning a replacement iff
	// the active cluster in the kubeconfig changes. This could still plan an erroneous replacement,
	// but should work for the majority of cases.
	//
	// The alternative of ignoring changes to the kubeconfig is untenable; if the k8s cluster has
	// changed, any dependent resources must be recreated, and ignoring changes prevents that from
	// happening.
	oldActiveCluster, oldFound := getActiveClusterFromConfig(oldConfig, olds)
	activeCluster, found := getActiveClusterFromConfig(newConfig, news)
	if oldFound && found && !reflect.DeepEqual(oldActiveCluster, activeCluster) {
		// one of these properties must have changed for the active cluster to change.
		for _, key := range []string{"kubeconfig", "context", "cluster"} {
			if slices.Contains(diffs, key) {
				replaces = append(replaces, key)
			}
		}
	}
	logger.V(7).Infof("%s: diffs %v / replaces %v", label, diffs, replaces)

	resp = &pulumirpc.DiffResponse{
		Changes:  pulumirpc.DiffResponse_DIFF_SOME,
		Diffs:    diffs,
		Replaces: replaces,
	}
	return resp, nil
}

// Configure configures the resource provider with "globals" that control its behavior.
func (k *kubeProvider) Configure(_ context.Context, req *pulumirpc.ConfigureRequest) (*pulumirpc.ConfigureResponse, error) {
	const trueStr = "true"

	// Configure Helm settings based on the ambient Helm environment,
	// using the provider configuration as overrides.
	helmSettings := helmcli.New()
	helmFlags := helmSettings.RESTClientGetter().(*genericclioptions.ConfigFlags)
	helmSettings.Debug = true // enable verbose logging (piped to glog at level 6)

	vars := req.GetVariables()

	//
	// Set simple configuration settings.
	//

	k.opts = kubeOpts{
		rejectUnknownResources: vars["kubernetes:config:rejectUnknownResources"] == trueStr,
	}
	k.enableSecrets = req.GetAcceptSecrets()

	//
	// Configure client-go using provided or ambient kubeconfig file.
	//
	if defaultNamespace := vars["kubernetes:config:namespace"]; defaultNamespace != "" {
		k.defaultNamespace = defaultNamespace
		helmSettings.SetNamespace(defaultNamespace)
		logger.V(9).Infof("namespace set to %v", defaultNamespace)
	}

	// Compute config overrides.
	overrides := &clientcmd.ConfigOverrides{
		Context: clientapi.Context{
			Cluster:   vars["kubernetes:config:cluster"],
			Namespace: k.defaultNamespace,
		},
		CurrentContext: vars["kubernetes:config:context"],
	}
	if overrides.Context.Cluster != "" {
		helmFlags.ClusterName = &overrides.Context.Cluster
	}
	if overrides.CurrentContext != "" {
		helmSettings.KubeContext = overrides.CurrentContext
	}

	deleteUnreachable := func() bool {
		// If the provider flag is set, use that value to determine behavior. This will override the ENV var.
		if enabled, exists := vars["kubernetes:config:deleteUnreachable"]; exists {
			return enabled == trueStr
		}
		// If the provider flag is not set, fall back to the ENV var.
		if enabled, exists := os.LookupEnv("PULUMI_K8S_DELETE_UNREACHABLE"); exists {
			return enabled == trueStr
		}
		// Default to false.
		return false
	}
	if deleteUnreachable() {
		k.deleteUnreachable = true
	}

	skipUpdateUnreachable := func() bool {
		// If the provider flag is set, use that value to determine behavior. This will override the ENV var.
		if enabled, exists := vars["kubernetes:config:skipUpdateUnreachable"]; exists {
			return enabled == trueStr
		}
		// If the provider flag is not set, fall back to the ENV var.
		if enabled, exists := os.LookupEnv("PULUMI_K8S_SKIP_UPDATE_UNREACHABLE"); exists {
			return enabled == trueStr
		}
		// Default to false.
		return false
	}
	if skipUpdateUnreachable() {
		k.skipUpdateUnreachable = true
	}

	enableServerSideApply := func() bool {
		// If the provider flag is set, use that value to determine behavior. This will override the ENV var.
		if enabled, exists := vars["kubernetes:config:enableServerSideApply"]; exists {
			return enabled == trueStr
		}
		// If the provider flag is not set, fall back to the ENV var.
		if enabled, exists := os.LookupEnv("PULUMI_K8S_ENABLE_SERVER_SIDE_APPLY"); exists {
			return enabled == trueStr
		}
		// Default to true.
		return true
	}
	if enableServerSideApply() {
		k.serverSideApplyMode = true
	}

	enableConfigMapMutable := func() bool {
		// If the provider flag is set, use that value to determine behavior. This will override the ENV var.
		if enabled, exists := vars["kubernetes:config:enableConfigMapMutable"]; exists {
			return enabled == trueStr
		}
		// If the provider flag is not set, fall back to the ENV var.
		if enabled, exists := os.LookupEnv("PULUMI_K8S_ENABLE_CONFIGMAP_MUTABLE"); exists {
			return enabled == trueStr
		}
		// Default to false.
		return false
	}
	if enableConfigMapMutable() {
		k.enableConfigMapMutable = true
	}

	enableSecretMutable := func() bool {
		// If the provider flag is set, use that value to determine behavior. This will override the ENV var.
		if enabled, exists := vars["kubernetes:config:enableSecretMutable"]; exists {
			return enabled == trueStr
		}
		// If the provider flag is not set, fall back to the ENV var.
		if enabled, exists := os.LookupEnv("PULUMI_K8S_ENABLE_SECRET_MUTABLE"); exists {
			return enabled == trueStr
		}
		// Default to false.
		return false
	}
	if enableSecretMutable() {
		k.enableSecretMutable = true
	}

	suppressDeprecationWarnings := func() bool {
		// If the provider flag is set, use that value to determine behavior. This will override the ENV var.
		if enabled, exists := vars["kubernetes:config:suppressDeprecationWarnings"]; exists {
			return enabled == trueStr
		}
		// If the provider flag is not set, fall back to the ENV var.
		if enabled, exists := os.LookupEnv("PULUMI_K8S_SUPPRESS_DEPRECATION_WARNINGS"); exists {
			return enabled == trueStr
		}
		// Default to false.
		return false
	}
	if suppressDeprecationWarnings() {
		k.suppressDeprecationWarnings = true
	}

	suppressHelmHookWarnings := func() bool {
		// If the provider flag is set, use that value to determine behavior. This will override the ENV var.
		if enabled, exists := vars["kubernetes:config:suppressHelmHookWarnings"]; exists {
			return enabled == trueStr
		}
		// If the provider flag is not set, fall back to the ENV var.
		if enabled, exists := os.LookupEnv("PULUMI_K8S_SUPPRESS_HELM_HOOK_WARNINGS"); exists {
			return enabled == trueStr
		}
		// Default to false.
		return false
	}
	if suppressHelmHookWarnings() {
		k.suppressHelmHookWarnings = true
	}

	renderYamlToDirectory := func() string {
		// Read the config from the Provider.
		if directory, exists := vars["kubernetes:config:renderYamlToDirectory"]; exists && directory != "" {
			return directory
		}
		return ""
	}
	k.yamlDirectory = renderYamlToDirectory()
	k.yamlRenderMode = len(k.yamlDirectory) > 0

	var helmReleaseSettings HelmReleaseSettings
	if obj, ok := vars["kubernetes:config:helmReleaseSettings"]; ok {
		err := json.Unmarshal([]byte(obj), &helmReleaseSettings)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal helmReleaseSettings option: %w", err)
		}
	}

	// TODO: Once https://github.com/pulumi/pulumi/issues/8132 is fixed, we can drop the env var handling logic.

	helmDriver := func() string {
		if helmReleaseSettings.Driver != nil {
			return *helmReleaseSettings.Driver
		}

		// If the provider flag is not set, fall back to the ENV var.
		if driver, exists := os.LookupEnv("PULUMI_K8S_HELM_DRIVER"); exists {
			return driver
		}
		return "secret"
	}
	k.helmDriver = helmDriver() // TODO: Make sure this is in provider state

	helmPluginsPath := func() string {
		if helmReleaseSettings.PluginsPath != nil {
			return *helmReleaseSettings.PluginsPath
		}

		// If the provider flag is not set, fall back to the ENV var.
		if pluginsPath, exists := os.LookupEnv("PULUMI_K8S_HELM_PLUGINS_PATH"); exists {
			return pluginsPath
		}
		return helmpath.DataPath("plugins")
	}
	k.helmPluginsPath = helmPluginsPath()
	if helmReleaseSettings.PluginsPath != nil {
		helmSettings.PluginsDirectory = *helmReleaseSettings.PluginsPath
	}

	helmRegistryConfigPath := func() string {
		if helmReleaseSettings.RegistryConfigPath != nil {
			return *helmReleaseSettings.RegistryConfigPath
		}

		// If the provider flag is not set, fall back to the ENV var.
		if registryPath, exists := os.LookupEnv("PULUMI_K8S_HELM_REGISTRY_CONFIG_PATH"); exists {
			return registryPath
		}
		legacyPath := helmpath.ConfigPath("registry.json")
		if _, err := os.Stat(legacyPath); err == nil {
			return legacyPath
		}
		return helmpath.ConfigPath("registry/config.json")
	}
	k.helmRegistryConfigPath = helmRegistryConfigPath()
	if helmReleaseSettings.RegistryConfigPath != nil {
		helmSettings.RegistryConfig = k.helmRegistryConfigPath
	}

	helmRepositoryConfigPath := func() string {
		if helmReleaseSettings.RepositoryConfigPath != nil {
			return *helmReleaseSettings.RepositoryConfigPath
		}

		if repositoryConfigPath, exists := os.LookupEnv("PULUMI_K8S_HELM_REPOSITORY_CONFIG_PATH"); exists {
			return repositoryConfigPath
		}
		return helmpath.ConfigPath("repositories.yaml")
	}
	k.helmRepositoryConfigPath = helmRepositoryConfigPath()
	if helmReleaseSettings.RepositoryConfigPath != nil {
		helmSettings.RepositoryConfig = k.helmRepositoryConfigPath
	}

	helmRepositoryCache := func() string {
		if helmReleaseSettings.RepositoryCache != nil {
			return *helmReleaseSettings.RepositoryCache
		}

		if repositoryCache, exists := os.LookupEnv("PULUMI_K8S_HELM_REPOSITORY_CACHE"); exists {
			return repositoryCache
		}
		return helmpath.CachePath("repository")
	}
	k.helmRepositoryCache = helmRepositoryCache()
	if helmReleaseSettings.RepositoryCache != nil {
		helmSettings.RepositoryCache = k.helmRepositoryCache
	}

	// Rather than erroring out on an invalid k8s config, mark the cluster as unreachable and conditionally bail out on
	// operations that require a valid cluster. This will allow us to perform invoke operations using the default
	// provider.
	unreachableCluster := func(err error) {
		k.clusterUnreachable = true
		k.clusterUnreachableReason = fmt.Sprintf(
			"failed to parse kubeconfig data in `kubernetes:config:kubeconfig`- %v", err)
	}

	var kubeconfig clientcmd.ClientConfig
	var apiConfig *clientapi.Config
	// Note: the Python SDK was setting the kubeconfig value to "" by default, so explicitly check for empty string.
	if pathOrContents, ok := vars["kubernetes:config:kubeconfig"]; ok && pathOrContents != "" {
		apiConfig, err := parseKubeconfigString(pathOrContents)
		if err != nil {
			unreachableCluster(err)
			// note: kubeconfig is not set when the cluster is unreachable
		} else {
			kubeconfig = clientcmd.NewDefaultClientConfig(*apiConfig, overrides)
			configurationNamespace, _, err := kubeconfig.Namespace()
			if err == nil {
				k.defaultNamespace = configurationNamespace
			}

			// initialize Helm settings to use the kubeconfig; use a generated file as necessary.
			configFile, err := writeKubeconfigToFile(apiConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to write kubeconfig file: %w", err)
			}
			helmSettings.KubeConfig = configFile
		}
	} else {
		// Use client-go to resolve the final configuration values for the client. Typically, these
		// values would reside in the $KUBECONFIG file, but can also be altered in several
		// places, including in env variables, client-go default values, and (if we allowed it) CLI
		// flags.
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
		kubeconfig = clientcmd.NewInteractiveDeferredLoadingClientConfig(loadingRules, overrides, os.Stdin)
		configurationNamespace, _, err := kubeconfig.Namespace()
		if err == nil {
			k.defaultNamespace = configurationNamespace
		}
	}

	var kubeClientSettings KubeClientSettings
	if obj, ok := vars["kubernetes:config:kubeClientSettings"]; ok {
		err := json.Unmarshal([]byte(obj), &kubeClientSettings)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal kubeClientSettings option: %w", err)
		}
	}

	// TODO: Once https://github.com/pulumi/pulumi/issues/8132 is fixed, we can drop the env var handling logic.
	if burst := os.Getenv("PULUMI_K8S_CLIENT_BURST"); burst != "" && kubeClientSettings.Burst == nil {
		asInt, err := strconv.Atoi(burst)
		if err != nil {
			return nil, fmt.Errorf("invalid value specified for PULUMI_K8S_CLIENT_BURST: %w", err)
		}
		kubeClientSettings.Burst = &asInt
	} else if kubeClientSettings.Burst == nil {
		v := 120 // Increased from default value of 10
		kubeClientSettings.Burst = &v
	}

	if qps := os.Getenv("PULUMI_K8S_CLIENT_QPS"); qps != "" && kubeClientSettings.QPS == nil {
		asFloat, err := strconv.ParseFloat(qps, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value specified for PULUMI_K8S_CLIENT_QPS: %w", err)
		}
		kubeClientSettings.QPS = &asFloat
	} else if kubeClientSettings.QPS == nil {
		v := 50.0 // Increased from default value of 5.0
		kubeClientSettings.QPS = &v
	}

	if timeout := os.Getenv("PULUMI_K8S_CLIENT_TIMEOUT"); timeout != "" && kubeClientSettings.Timeout == nil {
		asInt, err := strconv.Atoi(timeout)
		if err != nil {
			return nil, fmt.Errorf("invalid value specified for PULUMI_K8S_CLIENT_TIMEOUT: %w", err)
		}
		kubeClientSettings.Timeout = &asInt
	}

	// Attempt to load the configuration from the provided kubeconfig. If this fails, mark the cluster as unreachable.
	var config *rest.Config
	if !k.clusterUnreachable {
		contract.Assertf(kubeconfig != nil, "expected kubeconfig to be initialized")
		var err error
		config, err = kubeconfig.ClientConfig()
		if err != nil {
			k.clusterUnreachable = true
			k.clusterUnreachableReason = fmt.Sprintf("unable to load Kubernetes client configuration from kubeconfig file. Make sure you have: \n\n"+
				" \t • set up the provider as per https://www.pulumi.com/registry/packages/kubernetes/installation-configuration/ \n\n %v", err)
			config = nil
		} else {
			if kubeClientSettings.Burst != nil {
				config.Burst = *kubeClientSettings.Burst
				helmSettings.BurstLimit = *kubeClientSettings.Burst
				logger.V(9).Infof("kube client burst set to %v", config.Burst)
			}
			if kubeClientSettings.QPS != nil {
				config.QPS = float32(*kubeClientSettings.QPS)
				helmSettings.QPS = float32(*kubeClientSettings.QPS)
				logger.V(9).Infof("kube client QPS set to %v", config.QPS)
			}
			if kubeClientSettings.Timeout != nil {
				config.Timeout = time.Duration(*kubeClientSettings.Timeout) * time.Second
				helmFlags.Timeout = ptr.To(strconv.Itoa(*kubeClientSettings.Timeout))
				logger.V(9).Infof("kube client timeout set to %v", config.Timeout)
			}
			config.WarningHandler = rest.NoWarnings{}
			config.UserAgent = version.UserAgent
		}
	}

	var err error
	k.clientSet, k.logClient, err = k.makeClient(k.canceler.context, config)
	if err != nil {
		return nil, err
	}

	// These operations require a reachable cluster.
	if !k.clusterUnreachable {
		k.k8sVersion = cluster.TryGetServerVersion(k.clientSet.DiscoveryClientCached)

		if k.k8sVersion.Compare(cluster.ServerVersion{Major: 1, Minor: 13}) < 0 {
			return nil, fmt.Errorf("minimum supported cluster version is v1.13. found v%s", k.k8sVersion)
		}

		if _, err = k.getResources(); err != nil {
			k.clusterUnreachable = true
			k.clusterUnreachableReason = fmt.Sprintf(
				"unable to load schema information from the API server: %v", err)
		}
	}

	if !k.clusterUnreachable {
		k.helmSettings = helmSettings
	}

	k.helmReleaseProvider, err = newHelmReleaseProvider(
		k.host,
		k.canceler,
		apiConfig,
		overrides,
		config,
		k.clientSet,
		k.helmDriver,
		k.defaultNamespace,
		k.enableSecrets,
		k.helmPluginsPath,
		k.helmRegistryConfigPath,
		k.helmRepositoryConfigPath,
		k.helmRepositoryCache,
		k.clusterUnreachable,
		k.clusterUnreachableReason)
	if err != nil {
		return nil, err
	}

	return &pulumirpc.ConfigureResponse{
		AcceptSecrets:                   true,
		SupportsPreview:                 true,
		SupportsAutonamingConfiguration: true,
	}, nil
}

// Invoke dynamically executes a built-in function in the provider.
func (k *kubeProvider) Invoke(ctx context.Context,
	req *pulumirpc.InvokeRequest,
) (*pulumirpc.InvokeResponse, error) {
	// Important: Some invoke logic is intended to run during preview, and the Kubernetes provider
	// inputs may not have resolved yet. Any invoke logic that depends on an active cluster must check
	// k.clusterUnreachable and handle that condition appropriately.

	tok := req.GetTok()
	label := fmt.Sprintf("%s.Invoke(%s)", k.label(), tok)
	args, err := plugin.UnmarshalProperties(
		req.GetArgs(), plugin.MarshalOptions{Label: label, KeepUnknowns: true})
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal %v args during an Invoke call: %w", tok, err)
	}

	switch tok {
	case invokeDecodeYaml:
		var text, defaultNamespace string
		if textArg := args["text"]; textArg.HasValue() && textArg.IsString() {
			text = textArg.StringValue()
		} else {
			return nil, errors.New("missing required field 'text' of type string")
		}
		if defaultNsArg := args["defaultNamespace"]; defaultNsArg.HasValue() && defaultNsArg.IsString() {
			defaultNamespace = defaultNsArg.StringValue()
		}

		result, err := decodeYaml(text, defaultNamespace, k.clientSet)
		if err != nil {
			return nil, err
		}

		objProps, err := plugin.MarshalProperties(
			resource.NewPropertyMapFromMap(map[string]any{"result": result}),
			plugin.MarshalOptions{
				Label: label, KeepUnknowns: true, SkipNulls: true,
			})
		if err != nil {
			return nil, err
		}

		return &pulumirpc.InvokeResponse{Return: objProps}, nil
	case invokeHelmTemplate:
		var jsonOpts string
		if jsonOptsArgs := args["jsonOpts"]; jsonOptsArgs.HasValue() && jsonOptsArgs.IsString() {
			jsonOpts = jsonOptsArgs.StringValue()
		} else {
			return nil, errors.New("missing required field 'jsonOpts' of type string")
		}

		var opts HelmChartOpts
		err = json.Unmarshal([]byte(jsonOpts), &opts)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal 'jsonOpts': %w", err)
		}

		text, err := helmTemplate(k.host, opts, k.clientSet)
		if err != nil {
			return nil, fmt.Errorf("failed to generate YAML for specified Helm chart: %w", err)
		}

		// Decode the generated YAML here to avoid an extra invoke in the client.
		result, err := decodeYaml(text, opts.Namespace, k.clientSet)
		if err != nil {
			return nil, fmt.Errorf("failed to decode YAML for specified Helm chart: %w", err)
		}

		objProps, err := plugin.MarshalProperties(
			resource.NewPropertyMapFromMap(map[string]any{"result": result}),
			plugin.MarshalOptions{
				Label: label, KeepUnknowns: true, SkipNulls: true,
			})
		if err != nil {
			return nil, err
		}

		return &pulumirpc.InvokeResponse{Return: objProps}, nil

	case invokeKustomize:
		var directory string
		if directoryArg := args["directory"]; directoryArg.HasValue() && directoryArg.IsString() {
			directory = directoryArg.StringValue()
		} else {
			return nil, errors.New("missing required field 'directory' of type string")
		}

		result, err := kustomizeDirectory(ctx, directory, k.clientSet)
		if err != nil {
			return nil, err
		}

		objProps, err := plugin.MarshalProperties(
			resource.NewPropertyMapFromMap(map[string]any{"result": result}),
			plugin.MarshalOptions{
				Label: label, KeepUnknowns: true, SkipNulls: true,
			})
		if err != nil {
			return nil, err
		}

		return &pulumirpc.InvokeResponse{Return: objProps}, nil

	default:
		return nil, fmt.Errorf("unknown Invoke type %q", tok)
	}
}

// Attach sends the engine address to an already running plugin.
func (k *kubeProvider) Attach(_ context.Context, req *pulumirpc.PluginAttach) (*empty.Empty, error) {
	host, err := provider.NewHostClient(req.GetAddress())
	if err != nil {
		return nil, err
	}
	k.host = host
	return &empty.Empty{}, nil
}

// Check validates that the given property bag is valid for a resource of the given type and returns
// the inputs that should be passed to successive calls to Diff, Create, or Update for this
// resource. As a rule, the provider inputs returned by a call to Check should preserve the original
// representation of the properties as present in the program inputs. Though this rule is not
// required for correctness, violations thereof can negatively impact the end-user experience, as
// the provider inputs are using for detecting and rendering diffs.
func (k *kubeProvider) Check(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	//
	// Behavior as of v0.12.x: We take two inputs:
	//
	// 1. req.News, the new resource inputs, i.e., the property bag coming from a custom resource like
	//    k8s.core.v1.Service
	// 2. req.Olds, the last version submitted from a custom resource.
	//
	// `req.Olds` are ignored (and are sometimes nil). `req.News` are validated, and `.metadata.name`
	// is given to it if it's not already provided.
	//

	urn := resource.URN(req.GetUrn())

	if k.clusterUnreachable {
		if k.skipUpdateUnreachable {
			_ = k.host.Log(ctx, diag.Warning, urn, "Cluster is unreachable but skipUpdateUnreachable flag is set to true, skipping...")
			return &pulumirpc.CheckResponse{
				Inputs: req.GetOlds(),
			}, nil
		}
	}

	if isHelmRelease(urn) {
		return k.helmReleaseProvider.Check(ctx, req)
	}

	if kinds.IsListURN(urn) {
		// TODO: It might be possible to automatically expand List resources into a list of the underlying resources.
		//       Until then, return a descriptive error message. https://github.com/pulumi/pulumi-kubernetes/issues/2494
		return nil, fmt.Errorf("list resources exist for compatibility with YAML manifests and Helm charts, " +
			"and cannot be created directly. Use the underlying resource type instead")
	}

	label := fmt.Sprintf("%s.Check(%s)", k.label(), urn)
	logger.V(9).Infof("%s executing", label)

	// Obtain old resource inputs. This is the old version of the resource(s) supplied by the user as
	// an update.
	oldResInputs := req.GetOlds()
	olds, err := plugin.UnmarshalProperties(oldResInputs, plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true, KeepSecrets: true,
	})
	if err != nil {
		return nil, err
	}

	// Obtain new resource inputs. This is the new version of the resource(s) supplied by the user as
	// an update.
	newResInputs := req.GetNews()
	news, err := plugin.UnmarshalProperties(newResInputs, plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.news", label),
		KeepUnknowns: true,
		SkipNulls:    true,
		KeepSecrets:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("check failed because malformed resource inputs: %w", err)
	}

	oldInputs := propMapToUnstructured(olds)
	newInputs := propMapToUnstructured(news)

	if !k.serverSideApplyMode && kinds.IsPatchResource(urn, newInputs.GetKind()) {
		return nil, fmt.Errorf("patch resources require Server-Side Apply mode, which is enabled using the " +
			"`enableServerSideApply` Provider config")
	}

	newInputs, err = normalizeInputs(newInputs)
	if err != nil {
		return nil, err
	}

	if k.serverSideApplyMode && kinds.IsPatchResource(urn, newInputs.GetKind()) {
		if len(newInputs.GetName()) == 0 {
			return nil, fmt.Errorf("patch resources require the `.metadata.name` field to be set")
		}
	}

	var failures []*pulumirpc.CheckFailure

	k.helmHookWarning(ctx, newInputs, urn)

	// Adopt name from old object if appropriate.
	//
	// If the user HAS NOT assigned a name in the new inputs, we autoname it and mark the object as
	// autonamed in `.metadata.annotations`. This makes it easier for `Diff` to decide whether this
	// needs to be `DeleteBeforeReplace`'d. If the resource is marked `DeleteBeforeReplace`, then
	// `Create` will allocate it a new name later.
	if len(oldInputs.Object) > 0 {
		// NOTE: If old inputs exist, they MAY have a name, either provided by the user, or based on generateName,
		// or filled in with a previous run of `Check`.
		metadata.AdoptOldAutonameIfUnnamed(newInputs, oldInputs, news)

		// If the resource has existing state, we only set the "managed-by: pulumi" label if it is already present. This
		// avoids causing diffs for cases where the resource is being imported, or was created using SSA. The goal in
		// both cases is to leave the resource unchanged. The label is added if already present, or omitted if not.
		if metadata.HasManagedByLabel(oldInputs) {
			_, err = metadata.TrySetManagedByLabel(newInputs)
			if err != nil {
				return nil, fmt.Errorf("Failed to create object because of a problem setting managed-by labels: %w", err)
			}
		}
	} else {
		err = metadata.AssignNameIfAutonamable(req.RandomSeed, req.Autonaming, newInputs, news, urn)
		if err != nil {
			return nil, err
		}

		// Set a "managed-by: pulumi" label on resources created with Client-Side Apply. Do not set this label for SSA
		// resources since the fieldManagers field contains granular information about the managers.
		if !k.serverSideApplyMode {
			_, err = metadata.TrySetManagedByLabel(newInputs)
			if err != nil {
				return nil, fmt.Errorf("Failed to create object because of a problem setting managed-by labels: %w", err)
			}
		}
	}
	if metadata.IsGenerateName(newInputs, news) {
		if k.serverSideApplyMode {
			return nil, fmt.Errorf("the `.metadata.generateName` field is not supported in Server-Side Apply mode")
		}
		if k.yamlRenderMode {
			return nil, fmt.Errorf("the `.metadata.generateName` field is not supported in YAML rendering mode")
		}
	}

	gvk, err := k.gvkFromURN(urn)
	if err != nil {
		return nil, err
	}

	// Skip the API version check if the cluster is unreachable.
	if !k.clusterUnreachable {
		if removed, version := kinds.RemovedAPIVersion(gvk, k.k8sVersion); removed {
			_ = k.host.Log(ctx, diag.Warning, urn, (&kinds.RemovedAPIError{GVK: gvk, Version: version}).Error())
		} else if !k.suppressDeprecationWarnings && kinds.DeprecatedAPIVersion(gvk, &k.k8sVersion) {
			_ = k.host.Log(ctx, diag.Warning, urn, gen.APIVersionComment(gvk))
		}
	}

	// If a default namespace is set on the provider for this resource, check if the resource has Namespaced
	// or Global scope. For namespaced resources, set the namespace to the default value if unset.
	if k.defaultNamespace != "" && len(newInputs.GetNamespace()) == 0 {
		namespacedKind, err := clients.IsNamespacedKind(gvk, k.clientSet)
		if err != nil {
			if clients.IsNoNamespaceInfoErr(err) {
				// This is probably a CustomResource without a registered CustomResourceDefinition.
				// Since we can't tell for sure at this point, assume it is namespaced, and correct if
				// required during the Create step.
				namespacedKind = true
			} else {
				return nil, err
			}
		}

		if namespacedKind {
			newInputs.SetNamespace(k.defaultNamespace)
		}
	}

	// HACK: Do not validate against OpenAPI spec if there is a computed value. The OpenAPI spec
	// does not know how to deal with the placeholder values for computed values.
	if !hasComputedValue(newInputs) && !k.clusterUnreachable {
		resources, err := k.getResources()
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch OpenAPI schema from the API server: %w", err)
		}

		// Validate the object according to the OpenAPI schema for its GVK.
		err = openapi.ValidateAgainstSchema(resources, newInputs)
		if err != nil {
			resourceNotFound := apierrors.IsNotFound(err) ||
				strings.Contains(err.Error(), "is not supported by the server")
			k8sAPIUnreachable := strings.Contains(err.Error(), "connection refused")
			if resourceNotFound && k.gvkExists(newInputs) {
				failures = append(failures, &pulumirpc.CheckFailure{
					Reason: fmt.Sprintf(" Found API Group, but it did not contain a schema for %q", gvk),
				})
			} else if k8sAPIUnreachable {
				k8sURL := ""
				if err, ok := err.(*url.Error); ok {
					k8sURL = fmt.Sprintf("at %q", err.URL)
				}
				failures = append(failures, &pulumirpc.CheckFailure{
					Reason: fmt.Sprintf(" Kubernetes API server %s is unreachable. It's "+
						"possible that the URL or authentication information in your "+
						"kubeconfig is incorrect: %v", k8sURL, err),
				})
			} else if k.opts.rejectUnknownResources {
				// If the schema doesn't exist, it could still be a CRD (which may not have a
				// schema). Thus, if we are directed to check resources even if they have unknown
				// types, we fail here.
				return nil, fmt.Errorf("unable to fetch schema for resource type %s/%s: %w",
					newInputs.GetAPIVersion(), newInputs.GetKind(), err)
			}
		}
	}

	if clients.IsCRD(newInputs) {
		// add the CRD to the cache such that it contains all the CRDs that the program intends to create.
		// Do it now instead of later because update is called only if there's a non-empty diff,
		// and we want to ensure that the CRD is in the cache to support lookups by the component resources.
		if err := k.clientSet.CRDCache.AddCRD(newInputs); err != nil {
			return nil, err
		}
	}

	checkedInputs := resource.NewPropertyMapFromMap(newInputs.Object)
	annotateSecrets(checkedInputs, news)

	autonamedInputs, err := plugin.MarshalProperties(checkedInputs, plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.autonamedInputs", label),
		KeepUnknowns: true,
		SkipNulls:    true,
		KeepSecrets:  k.enableSecrets,
	})
	if err != nil {
		return nil, err
	}

	if k.yamlRenderMode {
		if checkedInputs.ContainsSecrets() {
			_ = k.host.Log(ctx, diag.Warning, urn, "rendered YAML will contain a secret value in plaintext")
		}
	}

	// Return new, possibly-autonamed inputs.
	return &pulumirpc.CheckResponse{Inputs: autonamedInputs, Failures: failures}, nil
}

var testHooks = regexp.MustCompile(`test|test-success|test-failure`)

// helmHookWarning logs a warning if a Chart contains unsupported hooks. The warning can be disabled by setting
// the suppressHelmHookWarnings provider flag or related ENV var.
func (k *kubeProvider) helmHookWarning(ctx context.Context, newInputs *unstructured.Unstructured, urn resource.URN) {
	hasHelmHook := false
	for key, value := range newInputs.GetAnnotations() {
		// If annotations with a reserved internal prefix exist, ignore them.
		if metadata.IsInternalAnnotation(key) {
			_ = k.host.Log(ctx, diag.Warning, urn,
				fmt.Sprintf("ignoring user-specified value for internal annotation %q", key))
		}

		// If the Helm hook annotation is found, set the hasHelmHook flag.
		if has := metadata.IsHelmHookAnnotation(key); has {
			// Test hooks are handled, so ignore this one.
			if testHooks.MatchString(value) {
				hasHelmHook = hasHelmHook || has
			}
		}
	}
	if hasHelmHook && !k.suppressHelmHookWarnings {
		_ = k.host.Log(ctx, diag.Warning, urn,
			"This resource contains Helm hooks that are not currently supported by Pulumi. The resource will "+
				"be created, but any hooks will not be executed. Hooks support is tracked at "+
				"https://github.com/pulumi/pulumi-kubernetes/issues/555 -- This warning can be disabled by setting "+
				"the PULUMI_K8S_SUPPRESS_HELM_HOOK_WARNINGS environment variable")
	}
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (k *kubeProvider) Diff(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	//
	// Behavior as of v4.0: We take 2 inputs:
	//
	// 1. req.News, the new resource inputs, i.e., the property bag coming from a custom resource like
	//    k8s.core.v1.Service
	// 2. req.Olds, the old _state_ returned by a `Create` or an `Update`.
	//
	// Kubernetes sets many additional fields that are not present in the resource inputs. We want to compare new inputs
	// to the most recent "live" values for the resource, but also don't want to show diffs for fields that are managed
	// by the cluster. We accomplish this by pruning the live state to match the shape of the old inputs, and then
	// comparing between the "pruned live" inputs and the new inputs.
	//
	// Note that comparing the old inputs to the new inputs will miss resource drift caused by other controllers
	// modifying the resources because the Pulumi inputs have not changed. Prior versions (pre-4.0) of the provider
	// used the "kubectl.kubernetes.io/last-applied-configuration" annotation to partially work around this problem.
	// This annotation was updated by the provider and some `kubectl` commands to capture the most recent set of inputs
	// used to produce the current resource state. When the annotation was present, this value was used instead of the
	// old inputs for the diff computation. This approach led to many problems, so the v4.0 release of the provider
	// removed the use of this annotation in favor of the "pruned live" input approach.
	//

	urn := resource.URN(req.GetUrn())
	if isHelmRelease(urn) {
		return k.helmReleaseProvider.Diff(ctx, req)
	}

	label := fmt.Sprintf("%s.Diff(%s)", k.label(), urn)
	logger.V(9).Infof("%s executing", label)

	// Get old state. This is an object that includes `inputs` previously supplied by the user, and the live state
	// we received back from the API server.
	oldState, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true, KeepSecrets: true,
	})
	if err != nil {
		return nil, err
	}

	// Get new resource inputs. The user is submitting these as an update.
	newResInputs, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.news", label),
		KeepUnknowns: true,
		SkipNulls:    true,
		KeepSecrets:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("diff failed because malformed resource inputs: %w", err)
	}

	newInputs := propMapToUnstructured(newResInputs)

	oldInputs, oldLive := parseCheckpointObject(oldState)
	contract.Assertf(oldLive.GetName() != "", "expected live object name to be nonempty: %v", oldLive)

	oldInputs, err = normalizeInputs(oldInputs)
	if err != nil {
		return nil, err
	}
	newInputs, err = normalizeInputs(newInputs)
	if err != nil {
		return nil, err
	}
	oldLivePruned := pruneLiveState(oldLive, oldInputs)

	gvk := k.gvkFromUnstructured(newInputs)

	namespacedKind, err := clients.IsNamespacedKind(gvk, k.clientSet)
	if err != nil {
		if clients.IsNoNamespaceInfoErr(err) {
			// This is probably a CustomResource without a registered CustomResourceDefinition.
			// Since we can't tell for sure at this point, assume it is namespaced, and correct if
			// required during the Create step.
			namespacedKind = true
		} else {
			return nil, fmt.Errorf(
				"API server returned error when asked if resource type %s is namespaced: %w", gvk, err)
		}
	}

	if namespacedKind {
		// Explicitly set the "default" namespace if unset so that the diff ignores it.
		oldLivePruned.SetNamespace(canonicalNamespace(oldLivePruned.GetNamespace()))
		newInputs.SetNamespace(canonicalNamespace(newInputs.GetNamespace()))
	} else {
		// Clear the namespace if it was set erroneously.
		oldLivePruned.SetNamespace("")
		newInputs.SetNamespace("")
	}
	if oldLivePruned.GroupVersionKind().Empty() {
		oldLivePruned.SetGroupVersionKind(gvk)
	}
	// If a resource was created without SSA enabled, and then the related provider was changed to enable SSA, a
	// resourceVersion may have been set on the old resource state. This produces erroneous diffs, so remove the
	// value from the oldLivePruned prior to computing the diff.
	if k.serverSideApplyMode && len(oldLivePruned.GetResourceVersion()) > 0 {
		oldLivePruned.SetResourceVersion("")
	}
	// If a name was specified in the new inputs, be sure that the old live object has the previous name.
	// This makes it possible to update the program to set `.metadata.name` to the name that was
	// made by `.metadata.generateName` without triggering replacement.
	if newInputs.GetName() != "" {
		oldLivePruned.SetName(oldLive.GetName())
	}

	var patch []byte
	patchBase := oldLivePruned.Object

	// Compute a diff between the pruned live state and the new inputs.
	patch, err = k.inputPatch(oldLivePruned, newInputs)
	if err != nil {
		return nil, fmt.Errorf("Failed to check for changes in resource %q: %w", urn, err)
	}

	patchObj := map[string]any{}
	if err = json.Unmarshal(patch, &patchObj); err != nil {
		return nil, fmt.Errorf(
			"Failed to check for changes in resource %q because of an error serializing "+
				"the JSON patch describing resource changes: %w", urn, err)
	}

	hasChanges := pulumirpc.DiffResponse_DIFF_NONE

	var replaces []string
	var detailedDiff map[string]*pulumirpc.PropertyDiff
	if len(patchObj) != 0 {
		// Changing the identity of the resource always causes a replacement.
		forceNewFields := []string{".metadata.name", ".metadata.namespace"}
		if !kinds.IsPatchResource(urn, newInputs.GetKind()) { // Patch resources can be updated in place for all other properties.
			forceNewFields = k.forceNewProperties(newInputs)
		}
		if detailedDiff, err = convertPatchToDiff(patchObj, patchBase, newInputs.Object, oldLivePruned.Object, forceNewFields...); err != nil {
			return nil, fmt.Errorf(
				"Failed to check for changes in resource %q because of an error "+
					"converting JSON patch describing resource changes to a diff: %w",
				urn, err)
		}

		// Remove any ignored changes from the computed diff.
		var ignorePaths []resource.PropertyPath
		for _, ignore := range req.IgnoreChanges {
			ignorePath, err := resource.ParsePropertyPath(ignore)
			if err != nil {
				return nil, fmt.Errorf("failed to parse ignore path: %w", err)
			}
			ignorePaths = append(ignorePaths, ignorePath)
		}
		if len(ignorePaths) > 0 {
			var diffPaths []resource.PropertyPath
			for p := range detailedDiff {
				diffPath, err := resource.ParsePropertyPath(p)
				if err != nil {
					return nil, fmt.Errorf("failed to parse diff path: %w", err)
				}
				diffPaths = append(diffPaths, diffPath)
			}
			for _, ignorePath := range ignorePaths {
				for _, diffPath := range diffPaths {
					if ignorePath.Contains(diffPath) {
						delete(detailedDiff, diffPath.String())
					}
				}
			}
		}

		if len(detailedDiff) > 0 {
			hasChanges = pulumirpc.DiffResponse_DIFF_SOME

			for k, v := range detailedDiff {
				switch v.Kind {
				case pulumirpc.PropertyDiff_ADD_REPLACE, pulumirpc.PropertyDiff_DELETE_REPLACE, pulumirpc.PropertyDiff_UPDATE_REPLACE:
					replaces = append(replaces, k)
				case pulumirpc.PropertyDiff_DELETE:
					if k == "metadata" {
						replaces = append(replaces, k)
					}
				}
			}
		}
	}

	if metadata.ReplaceUnready(newInputs) {
		switch {
		case k.clusterUnreachable:
			// Check if the cluster is unreachable. If it is, we can't check the status of the resource otherwise
			// a panic occurs due to the client being nil.
			_ = k.host.Log(ctx, diag.Warning, urn, "Cluster is unreachable, skipping replaceUnready check")
		case newInputs.GetKind() == "Job":
			// Fetch current Job status and check point-in-time readiness. Errors are ignored.
			if live, err := k.readLiveObject(oldLive); err == nil {
				jobChecker := checkjob.NewJobChecker()
				job, err := clients.FromUnstructured(live)
				if err == nil {
					if !jobChecker.Ready(job) {
						hasChanges = pulumirpc.DiffResponse_DIFF_SOME
						replaces = append(replaces, `.metadata.annotations["pulumi.com/replaceUnready"]`)
					}
				}
			}
		default:
			_ = k.host.Log(ctx, diag.Warning, urn, "replaceUnready annotation is not supported for this resource")
		}
	}

	// Delete before replacement if we are forced to replace the old object, and the new version of
	// that object MUST have the same name.
	deleteBeforeReplace := // 1. We know resource must be replaced.
		len(replaces) > 0 &&
			// 2. Object is named (i.e., not using metadata.generateName).
			metadata.IsNamed(newInputs, newResInputs) &&
			// 3. Object is NOT autonamed (i.e., user manually named it, and therefore we can't
			// auto-generate the name).
			!metadata.IsAutonamed(newInputs) &&
			// 4. The new, user-specified name is the same as the old name.
			newInputs.GetName() == oldLive.GetName() &&
			// 5. The resource is being deployed to the same namespace (i.e., we aren't creating the
			// object in a new namespace and then deleting the old one).
			newInputs.GetNamespace() == oldLive.GetNamespace()

	return &pulumirpc.DiffResponse{
		Changes:             hasChanges,
		Replaces:            replaces,
		Stables:             []string{},
		DeleteBeforeReplace: deleteBeforeReplace,
		DetailedDiff:        detailedDiff,
		HasDetailedDiff:     true,
	}, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID.
// (The input ID must be blank.)  If this call fails, the resource must not have been created (i.e.,
// it is "transactional").
func (k *kubeProvider) Create(
	ctx context.Context, req *pulumirpc.CreateRequest,
) (*pulumirpc.CreateResponse, error) {
	//
	// Behavior as of v4.0: We take 1 input:
	//
	// 1. `req.Properties`, the new resource inputs submitted by the user, after having been returned
	// by `Check`.
	//
	// This is used to create a new resource, and the computed values are returned. Importantly:
	//
	// * The return is formatted as a "checkpoint object", which includes both inputs and the live state of the
	//   resource. This is important both for `Diff` and for `Update`. See comments in those methods for details.
	//
	urn := resource.URN(req.GetUrn())
	if isHelmRelease(urn) {
		return k.helmReleaseProvider.Create(ctx, req)
	}

	label := fmt.Sprintf("%s.Create(%s)", k.label(), urn)
	logger.V(9).Infof("%s executing", label)

	// Except in the case of yamlRender mode, Create requires a connection to a k8s cluster, so bail out
	// immediately if it is unreachable.
	if !req.GetPreview() && k.clusterUnreachable && !k.yamlRenderMode {
		return nil, fmt.Errorf("configured Kubernetes cluster is unreachable: %s", k.clusterUnreachableReason)
	}

	// Parse inputs
	newResInputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.properties", label),
		KeepUnknowns: true,
		SkipNulls:    true,
		KeepSecrets:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("create failed because malformed resource inputs: %w", err)
	}

	newInputs := propMapToUnstructured(newResInputs)

	// Skip if:
	// 1: The input values contain unknowns
	// 2: The cluster is unreachable or the resource GVK does not exist
	// 3: The resource is a Patch resource
	// 4: We are in client-side-apply mode
	skipPreview := hasComputedValue(newInputs) || !k.gvkExists(newInputs) || kinds.IsPatchResource(urn, newInputs.GetKind()) || !k.serverSideApplyMode
	// If this is a preview and the input meets one of the skip criteria, then return them as-is. This is compatible
	// with prior behavior implemented by the Pulumi engine.
	if req.GetPreview() && skipPreview {
		logger.V(9).Infof("cannot preview Create(%v)", urn)
		return &pulumirpc.CreateResponse{Id: "", Properties: req.GetProperties()}, nil
	}

	initialAPIVersion := newInputs.GetAPIVersion()
	fieldManager := k.fieldManagerName(nil, newResInputs, newInputs)

	if k.yamlRenderMode {
		if newResInputs.ContainsSecrets() {
			_ = k.host.Log(ctx, diag.Warning, urn, fmt.Sprintf(
				"rendered file %s contains a secret value in plaintext",
				renderPathForResource(newInputs, k.yamlDirectory)))
		}
		err := renderYaml(newInputs, k.yamlDirectory)
		if err != nil {
			return nil, err
		}

		obj := checkpointObject(newInputs, newInputs, newResInputs, initialAPIVersion, fieldManager)
		inputsAndComputed, err := plugin.MarshalProperties(
			obj, plugin.MarshalOptions{
				Label:        fmt.Sprintf("%s.inputsAndComputed", label),
				KeepUnknowns: true,
				SkipNulls:    true,
				KeepSecrets:  k.enableSecrets,
			})
		if err != nil {
			return nil, err
		}

		_ = k.host.LogStatus(ctx, diag.Info, urn, fmt.Sprintf(
			"rendered %s", renderPathForResource(newInputs, k.yamlDirectory)))

		return &pulumirpc.CreateResponse{
			Id: fqObjName(newInputs), Properties: inputsAndComputed,
		}, nil
	}

	resources, err := k.getResources()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch OpenAPI schema from the API server: %w", err)
	}
	config := await.CreateConfig{
		ProviderConfig: await.ProviderConfig{
			Context:           k.canceler.context,
			Host:              k.host,
			URN:               urn,
			InitialAPIVersion: initialAPIVersion,
			FieldManager:      fieldManager,
			ClusterVersion:    &k.k8sVersion,
			ClientSet:         k.clientSet,
			DedupLogger:       logging.NewLogger(k.canceler.context, k.host, urn),
			Resources:         resources,
			ServerSideApply:   k.serverSideApplyMode,
		},
		Inputs:  newInputs,
		Timeout: req.Timeout,
		Preview: req.GetPreview(),
	}
	initialized, awaitErr := await.Creation(config)
	if awaitErr != nil {
		if req.GetPreview() {
			if apierrors.IsForbidden(awaitErr) {
				logger.V(1).Infof("unable to compute Server-side dry-run and defaulting to client-side: %s", awaitErr)
				return &pulumirpc.CreateResponse{Id: "", Properties: req.GetProperties()}, nil
			}

			failedPreview := false
			_, isPreviewErr := awaitErr.(await.PreviewError)
			if k.isDryRunDisabledError(awaitErr) || isPreviewErr {
				failedPreview = true
			}

			if k.serverSideApplyMode && apierrors.IsAlreadyExists(awaitErr) {
				failedPreview = true
			}

			if failedPreview {
				logger.V(9).Infof("could not preview Create(%v): %v", urn, err)
				return &pulumirpc.CreateResponse{Id: "", Properties: req.GetProperties()}, nil
			}
		}

		if meta.IsNoMatchError(awaitErr) {
			// If it's a "no match" error, this is probably a CustomResource with no corresponding
			// CustomResourceDefinition. This usually happens if the CRD was not created, and we
			// print a more useful error message in this case.
			gvk, err := k.gvkFromURN(urn)
			if err != nil {
				return nil, err
			}
			gvkStr := gvk.GroupVersion().String() + "/" + gvk.Kind
			return nil, fmt.Errorf(
				"creation of resource %q with kind %s failed because the Kubernetes API server "+
					"reported that the apiVersion for this resource does not exist. "+
					"Verify that any required CRDs have been created: %w", urn, gvkStr, awaitErr)

		}
		partialErr, isPartialErr := awaitErr.(await.PartialError)
		if !isPartialErr {
			// Object creation failed.
			return nil, fmt.Errorf(
				"resource %q was not successfully created by the Kubernetes API server: %w", urn, awaitErr)
		}

		// Resource was created, but failed to become fully initialized.
		initialized = partialErr.Object()
	}
	contract.Assertf(initialized.GetName() != "", "expected live object name to be nonempty: %v", initialized)

	// We need to delete the empty status field returned from the API server if we are in
	// preview mode. Having the status field set will cause a panic during preview if the Pulumi
	// program attempts to read the status field.
	if req.GetPreview() {
		unstructured.RemoveNestedField(initialized.Object, "status")
	}

	obj := checkpointObject(newInputs, initialized, newResInputs, initialAPIVersion, fieldManager)
	inputsAndComputed, err := plugin.MarshalProperties(
		obj, plugin.MarshalOptions{
			Label:        fmt.Sprintf("%s.inputsAndComputed", label),
			KeepUnknowns: true,
			SkipNulls:    true,
			KeepSecrets:  k.enableSecrets,
		})
	if err != nil {
		return nil, err
	}

	if awaitErr != nil {
		// Resource was created but failed to initialize. Return live version of object so it can be
		// checkpointed.
		return nil, partialError(
			fqObjName(initialized),
			fmt.Errorf(
				"resource %q was successfully created, but the Kubernetes API server "+
					"reported that it failed to fully initialize or become live: %w", urn, awaitErr),

			inputsAndComputed,
			nil)
	}

	// Invalidate the client cache if this was a CRD. This will require subsequent CR creations to
	// refresh the cache, at which point the CRD definition will be present, so that it doesn't fail
	// with an `apierrors.IsNotFound`.
	if clients.IsCRD(newInputs) {
		k.clientSet.RESTMapper.Reset()
		k.invalidateResources()
	}

	id := ""
	if !req.GetPreview() {
		id = fqObjName(initialized)
	}
	return &pulumirpc.CreateResponse{Id: id, Properties: inputsAndComputed}, nil
}

// Read the current live state associated with a resource.  Enough state must be included in the
// inputs to uniquely identify the resource; this is typically just the resource ID, but may also
// include some properties.
func (k *kubeProvider) Read(ctx context.Context, req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	//
	// Behavior as of v4.0: We take 2 inputs:
	//
	// 1. `req.Properties`, the previous state of the resource.
	// 2. `req.Inputs`, the old resource inputs submitted by the user, after having been persisted
	// (e.g., by `Create` or `Update`).
	//
	// We use this information to read the live version of a Kubernetes resource. This is sometimes
	// then checkpointed (e.g., in the case of `refresh`). Specifically:
	//
	// * The return is formatted as a "checkpoint object", which includes both inputs and the live state of the
	//   resource. This is important both for `Diff` and for `Update`. See comments in those methods for details.
	//

	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Read(%s)", k.label(), urn)
	logger.V(9).Infof("%s executing", label)

	// If the cluster is unreachable, return an error unless the user has opted in to mark the resources to be deleted
	if k.clusterUnreachable {
		_ = k.host.Log(ctx, diag.Warning, urn, fmt.Sprintf(
			"configured Kubernetes cluster is unreachable: %s", k.clusterUnreachableReason))
		if k.deleteUnreachable {
			_ = k.host.Log(ctx, diag.Info, urn, fmt.Sprintf(
				"configured Kubernetes cluster is unreachable and the `deleteUnreachable` option is enabled. "+
					"Deleting the unreachable resource from Pulumi state"))
			return deleteResponse, nil
		}
		if k.skipUpdateUnreachable {
			_ = k.host.Log(ctx, diag.Info, urn, fmt.Sprintf(
				"configured Kubernetes cluster is unreachable and the `skipUnreachable` option is enabled. "+
					"Returned data could not reflect the actual cluster configuration."))
			return &pulumirpc.ReadResponse{
				Id:         req.GetId(),
				Properties: req.GetProperties(),
				Inputs:     req.GetInputs(),
			}, nil
		}

		return nil, fmt.Errorf("failed to read resource state due to unreachable cluster. If the cluster was " +
			"deleted, you can remove this resource from Pulumi state by rerunning the operation with the " +
			"PULUMI_K8S_DELETE_UNREACHABLE environment variable set to \"true\"")
	}

	if isHelmRelease(urn) {
		contract.Assertf(k.helmReleaseProvider != nil, "helmReleaseProvider not initialized.")
		return k.helmReleaseProvider.Read(ctx, req)
	}

	// Obtain new properties, create a Kubernetes `unstructured.Unstructured` that we can pass to the
	// validation routines.
	oldState, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true, KeepSecrets: true,
	})
	if err != nil {
		return nil, err
	}
	oldInputsPM, err := plugin.UnmarshalProperties(req.GetInputs(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.oldInputs", label), KeepUnknowns: true, SkipNulls: true, KeepSecrets: true,
	})
	if err != nil {
		return nil, err
	}

	namespace, name := parseFqName(req.GetId())
	if name == "" {
		return nil, fmt.Errorf(
			"failed to read resource because of a failure to parse resource name from request ID: %s",
			req.GetId())
	}

	readFromCluster := false
	oldInputs := propMapToUnstructured(oldInputsPM)
	_, oldLive := parseCheckpointObject(oldState)

	oldInputs, err = normalizeInputs(oldInputs)
	if err != nil {
		return nil, err
	}

	if oldInputs.GroupVersionKind().Empty() {
		if oldLive.GroupVersionKind().Empty() {
			gvk, err := k.gvkFromURN(urn)
			if err != nil {
				return nil, err
			}
			oldInputs.SetGroupVersionKind(gvk)
			readFromCluster = true
		} else {
			oldInputs.SetGroupVersionKind(oldLive.GroupVersionKind())
		}

		if oldInputs.GetName() == "" {
			oldInputs.SetName(name)
		}

		if oldInputs.GetNamespace() == "" {
			oldInputs.SetNamespace(namespace)
		}
	}

	initialAPIVersion := initialAPIVersion(oldState, oldInputs)
	fieldManager := k.fieldManagerName(nil, oldState, oldInputs)

	if k.yamlRenderMode {
		// Return a new "checkpoint object".
		state, err := plugin.MarshalProperties(
			checkpointObject(oldInputs, oldLive, oldState, initialAPIVersion, fieldManager), plugin.MarshalOptions{
				Label:        fmt.Sprintf("%s.state", label),
				KeepUnknowns: true,
				SkipNulls:    true,
				KeepSecrets:  k.enableSecrets,
			})
		if err != nil {
			return nil, err
		}

		inputs, err := plugin.MarshalProperties(oldInputsPM, plugin.MarshalOptions{
			Label: label + ".inputs", KeepUnknowns: true, SkipNulls: true, KeepSecrets: k.enableSecrets,
		})
		if err != nil {
			return nil, err
		}

		return &pulumirpc.ReadResponse{Id: req.GetId(), Properties: state, Inputs: inputs}, nil
	}

	resources, err := k.getResources()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch OpenAPI schema from the API server: %w", err)
	}
	config := await.ReadConfig{
		ProviderConfig: await.ProviderConfig{
			Context:           k.canceler.context,
			Host:              k.host,
			URN:               urn,
			InitialAPIVersion: initialAPIVersion,
			FieldManager:      fieldManager,
			ClusterVersion:    &k.k8sVersion,
			ClientSet:         k.clientSet,
			DedupLogger:       logging.NewLogger(k.canceler.context, k.host, urn),
			Resources:         resources,
		},
		Inputs:          oldInputs,
		ReadFromCluster: readFromCluster,
		Name:            name,
	}
	liveObj, readErr := await.Read(config)
	if readErr != nil {
		logger.V(3).Infof("%v", readErr)

		if meta.IsNoMatchError(readErr) {
			// If it's a "no match" error, this is probably a CustomResource with no corresponding
			// CustomResourceDefinition. This usually happens if the CRD was deleted, and it's safe
			// to consider the CR to be deleted as well in this case.
			return deleteResponse, nil
		}

		statusErr, ok := readErr.(*apierrors.StatusError)
		if ok && statusErr.ErrStatus.Code == 404 {
			// If it's a 404 error, this resource was probably deleted.
			return deleteResponse, nil
		}

		if partialErr, ok := readErr.(await.PartialError); ok {
			liveObj = partialErr.Object()
		}

		// If `liveObj == nil` at this point, it means we've encountered an error that is neither a
		// 404, nor an `await.PartialError`. For example, the master could be unreachable. We
		// should fail in this case.
		if liveObj == nil {
			return nil, readErr
		}

		// If we get here, resource successfully registered with the API server, but failed to
		// initialize.
	}
	contract.Assertf(liveObj.GetName() != "", "expected live object name to be nonempty: %v", liveObj)

	// Prune the live inputs to remove properties that are not present in the program inputs.
	liveInputs := pruneLiveState(liveObj, oldInputs)

	if readFromCluster {
		// If no previous inputs were known, populate the inputs from the live cluster state for the resource.
		pkgSpec := pulumischema.PackageSpec{}
		if err := json.Unmarshal(k.pulumiSchema, &pkgSpec); err != nil {
			return nil, err
		}
		res := pkgSpec.Resources[urn.Type().String()]
		for k := range res.InputProperties {
			if liveVal, ok := liveObj.Object[k]; ok {
				if err = unstructured.SetNestedField(liveInputs.Object, liveVal, k); err != nil {
					return nil, fmt.Errorf("failure setting field %q for %q: %w", k, urn, err)
				}
			}
		}

		// Cleanup some obviously non-input-ty fields.
		unstructured.RemoveNestedField(liveInputs.Object, "metadata", "creationTimestamp")
		unstructured.RemoveNestedField(liveInputs.Object, "metadata", "generation")
		unstructured.RemoveNestedField(liveInputs.Object, "metadata", "managedFields")
		unstructured.RemoveNestedField(liveInputs.Object, "metadata", "resourceVersion")
		unstructured.RemoveNestedField(liveInputs.Object, "metadata", "uid")
		unstructured.RemoveNestedField(liveInputs.Object, "metadata", "annotations", "deployment.kubernetes.io/revision")
		unstructured.RemoveNestedField(liveInputs.Object, "metadata", "annotations", lastAppliedConfigKey)
	}

	// TODO(lblackstone): not sure why this is needed
	id := fqObjName(liveObj)
	if reqID := req.GetId(); len(reqID) > 0 {
		id = reqID
	}

	// Return a new "checkpoint object".
	state, err := plugin.MarshalProperties(
		checkpointObject(liveInputs, liveObj, oldInputsPM, initialAPIVersion, fieldManager), plugin.MarshalOptions{
			Label:        fmt.Sprintf("%s.state", label),
			KeepUnknowns: true,
			SkipNulls:    true,
			KeepSecrets:  k.enableSecrets,
		})
	if err != nil {
		return nil, err
	}

	liveInputsPM := resource.NewPropertyMapFromMap(liveInputs.Object)
	annotateSecrets(liveInputsPM, oldInputsPM)

	inputs, err := plugin.MarshalProperties(liveInputsPM, plugin.MarshalOptions{
		Label: label + ".inputs", KeepUnknowns: true, SkipNulls: true, KeepSecrets: k.enableSecrets,
	})
	if err != nil {
		return nil, err
	}

	if readErr != nil {
		// Resource was created but failed to initialize. Return live version of object so it can be
		// checkpointed.
		logger.V(3).Infof("%v", partialError(id, readErr, state, inputs))
		return nil, partialError(id, readErr, state, inputs)
	}

	return &pulumirpc.ReadResponse{Id: id, Properties: state, Inputs: inputs}, nil
}

// Update updates an existing resource with new values. This client uses a Server-side Apply (SSA) patch by default, but
// also supports the older three-way JSON patch and the strategic merge patch as fallback options.
// See references [1], [2], [3].
//
// nolint
// [1]:
// https://kubernetes.io/docs/tasks/run-application/update-api-object-kubectl-patch/#use-a-json-merge-patch-to-update-a-deployment
// [2]:
// https://kubernetes.io/docs/concepts/overview/object-management-kubectl/declarative-config/#how-apply-calculates-differences-and-merges-changes
// [3]:
// https://kubernetes.io/docs/reference/using-api/server-side-apply
func (k *kubeProvider) Update(
	ctx context.Context, req *pulumirpc.UpdateRequest,
) (*pulumirpc.UpdateResponse, error) {
	//
	// Behavior as of v4.0: We take 2 inputs:
	//
	// 1. req.News, the new resource inputs, i.e., the property bag coming from a custom resource like
	//    k8s.core.v1.Service
	// 2. req.Olds, the old _state_ returned by a `Create` or an `Update`. The old state is a struct that
	//    contains the old inputs as well as the last computed value obtained from the Kubernetes API server.
	//
	// The provider uses Server-side Apply (SSA) patches by default, which allows the provider to send the desired
	// state of the resource to Kubernetes, and let the server decide how to merge in the changes. Previews are
	// computed using the dry-run option for the API call, which computes the result on the server without persisting
	// the changes.
	//
	// Previous versions of the provider used Client-side Apply (CSA) instead, which required the provider to compute
	// the merge patch, and then send the patch to the server. This patch is computed as a three-way merge between:
	// (1) the new inputs, (2) the computed state returned by the API server, and (3) the old inputs. This is the
	// main reason why the old state is an object with both the old inputs and the live version of the
	// object. CSA is provided as a fallback option, but is generally less reliable than using SSA.
	//

	urn := resource.URN(req.GetUrn())
	if isHelmRelease(urn) {
		return k.helmReleaseProvider.Update(ctx, req)
	}

	label := fmt.Sprintf("%s.Update(%s)", k.label(), urn)
	logger.V(9).Infof("%s executing", label)

	// Except in the case of yamlRender mode, Update requires a connection to a k8s cluster, so bail out
	// immediately if it is unreachable.
	if !req.GetPreview() && k.clusterUnreachable && !k.yamlRenderMode {
		return nil, fmt.Errorf("configured Kubernetes cluster is unreachable: %s", k.clusterUnreachableReason)
	}

	// Obtain old properties, create a Kubernetes `unstructured.Unstructured`.
	oldState, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true, KeepSecrets: true,
	})
	if err != nil {
		return nil, err
	}

	// Obtain new properties, create a Kubernetes `unstructured.Unstructured`.
	newResInputs, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.news", label),
		KeepUnknowns: true,
		SkipNulls:    true,
		KeepSecrets:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("update failed because malformed resource inputs: %w", err)
	}
	newInputs := propMapToUnstructured(newResInputs)
	newInputs, err = normalizeInputs(newInputs)
	if err != nil {
		return nil, err
	}

	// If this is a preview and the input values contain unknowns, or an unregistered GVK, return them as-is. This is
	// compatible with prior behavior implemented by the Pulumi engine.
	if req.GetPreview() && (hasComputedValue(newInputs) || !k.gvkExists(newInputs)) {
		logger.V(9).Infof("cannot preview Update(%v)", urn)
		return &pulumirpc.UpdateResponse{Properties: req.News}, nil
	}

	// Ignore old state; we'll get it from Kubernetes later.
	oldInputs, oldLive := parseCheckpointObject(oldState)

	// Pre-4.0 versions of the provider used the last-applied-configuration annotation for client-side diffing. This
	// annotation was automatically added to all resources by Pulumi. This annotation is no longer used, and needs to
	// be removed from resources where it is present. To avoid causing a large update after upgrading to a 4.x release,
	// only perform this operation during Update for resources that include other changes. This removal will not show
	// up in the preview, which is symmetric with the previous creation behavior, which also does not show this
	// annotation during preview.
	removeLastAppliedConfigurationAnnotation(oldLive, oldInputs)

	oldInputs, err = normalizeInputs(oldInputs)
	if err != nil {
		return nil, err
	}
	oldLivePruned := pruneLiveState(oldLive, oldInputs)

	initialAPIVersion := initialAPIVersion(oldState, oldInputs)
	fieldManagerOld := k.fieldManagerName(nil, oldState, oldInputs)
	fieldManager := k.fieldManagerName(nil, oldState, newInputs)

	if k.yamlRenderMode {
		if newResInputs.ContainsSecrets() {
			_ = k.host.LogStatus(ctx, diag.Warning, urn, fmt.Sprintf(
				"rendered file %s contains a secret value in plaintext",
				renderPathForResource(newInputs, k.yamlDirectory)))
		}
		err := renderYaml(newInputs, k.yamlDirectory)
		if err != nil {
			return nil, err
		}

		obj := checkpointObject(newInputs, newInputs, newResInputs, initialAPIVersion, fieldManager)
		inputsAndComputed, err := plugin.MarshalProperties(
			obj, plugin.MarshalOptions{
				Label:        fmt.Sprintf("%s.inputsAndComputed", label),
				KeepUnknowns: true,
				SkipNulls:    true,
				KeepSecrets:  k.enableSecrets,
			})
		if err != nil {
			return nil, err
		}

		_ = k.host.LogStatus(ctx, diag.Info, urn, fmt.Sprintf(
			"rendered %s", renderPathForResource(newInputs, k.yamlDirectory)))

		return &pulumirpc.UpdateResponse{Properties: inputsAndComputed}, nil
	}

	resources, err := k.getResources()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch OpenAPI schema from the API server: %w", err)
	}
	config := await.UpdateConfig{
		ProviderConfig: await.ProviderConfig{
			Context:           k.canceler.context,
			Host:              k.host,
			URN:               urn,
			InitialAPIVersion: initialAPIVersion,
			FieldManager:      fieldManager,
			ClusterVersion:    &k.k8sVersion,
			ClientSet:         k.clientSet,
			DedupLogger:       logging.NewLogger(k.canceler.context, k.host, urn),
			Resources:         resources,
			ServerSideApply:   k.serverSideApplyMode,
		},
		OldInputs:     oldLivePruned,
		OldOutputs:    oldLive,
		Inputs:        newInputs,
		Timeout:       req.Timeout,
		Preview:       req.GetPreview(),
		IgnoreChanges: req.IgnoreChanges,
	}
	// Apply update.
	initialized, awaitErr := await.Update(config)
	if awaitErr != nil {
		if req.GetPreview() && apierrors.IsForbidden(awaitErr) {
			logger.V(9).Infof("could not preview Update(%v): %v", urn, err)
			return &pulumirpc.UpdateResponse{Properties: req.News}, nil
		}

		if meta.IsNoMatchError(awaitErr) {
			// If it's a "no match" error, this is probably a CustomResource with no corresponding
			// CustomResourceDefinition. This usually happens if the CRD was not created, and we
			// print a more useful error message in this case.
			return nil, fmt.Errorf(
				"update of resource %q failed because the Kubernetes API server "+
					"reported that the apiVersion for this resource does not exist. "+
					"Verify that any required CRDs have been created: %w", urn, awaitErr)
		}

		var getErr error
		initialized, getErr = k.readLiveObject(oldLive)
		if getErr != nil {
			// Object update/creation failed.
			return nil, errors.Join(
				fmt.Errorf("update of resource %q failed: %w", urn, awaitErr),
				fmt.Errorf("unable to get cluster state: %w", getErr),
			)
		}
		// If we get here, resource successfully registered with the API server, but failed to
		// initialize.
	}
	contract.Assertf(initialized.GetName() != "", "expected live object name to be nonempty: %v", initialized)

	// Return a new "checkpoint object".
	obj := checkpointObject(newInputs, initialized, newResInputs, initialAPIVersion, fieldManager)
	inputsAndComputed, err := plugin.MarshalProperties(
		obj, plugin.MarshalOptions{
			Label:        fmt.Sprintf("%s.inputsAndComputed", label),
			KeepUnknowns: true,
			SkipNulls:    true,
			KeepSecrets:  k.enableSecrets,
		})
	if err != nil {
		return nil, err
	}

	if awaitErr != nil {
		// Resource was updated/created but failed to initialize. Return live version of object so it
		// can be checkpointed.
		return nil, partialError(
			fqObjName(initialized),
			fmt.Errorf(
				"the Kubernetes API server reported that %q failed to fully initialize "+
					"or become live: %w", fqObjName(initialized), awaitErr),

			inputsAndComputed,
			nil)
	}

	if k.serverSideApplyMode {
		// For non-preview updates, drop the old fieldManager if the value changes.
		if !req.GetPreview() && fieldManagerOld != fieldManager {
			client, err := k.clientSet.ResourceClientForObject(initialized)
			if err != nil {
				return nil, err
			}

			err = ssa.Relinquish(k.canceler.context, client, initialized, fieldManagerOld)
			if err != nil {
				return nil, err
			}
		}
	}

	return &pulumirpc.UpdateResponse{Properties: inputsAndComputed}, nil
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed
// to still exist.
func (k *kubeProvider) Delete(ctx context.Context, req *pulumirpc.DeleteRequest) (*pbempty.Empty, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Delete(%s)", k.label(), urn)
	logger.V(9).Infof("%s executing", label)

	// TODO(hausdorff): Propagate other options, like grace period through flags.

	// Obtain new properties, create a Kubernetes `unstructured.Unstructured`.
	oldState, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true, KeepSecrets: true,
	})
	if err != nil {
		return nil, err
	}
	oldInputs, _ := parseCheckpointObject(oldState)

	if isHelmRelease(urn) {
		if k.clusterUnreachable {
			return nil, fmt.Errorf("can't delete Helm Release with unreachable cluster. Reason: %q", k.clusterUnreachableReason)
		}
		return k.helmReleaseProvider.Delete(ctx, req)
	}

	_, current := parseCheckpointObject(oldState)
	_, name := parseFqName(req.GetId())

	if k.yamlRenderMode {
		file := renderPathForResource(current, k.yamlDirectory)
		err := os.Remove(file)
		if err != nil {
			// Most of the time, errors will be because the file was already deleted. In this case,
			// the operation succeeds. It's also possible that deletion fails due to file permission if
			// the user changed the directory out-of-band, so log the error to help debug this scenario.
			logger.V(3).Infof("Failed to delete YAML file: %q - %v", file, err)
		}

		_ = k.host.LogStatus(ctx, diag.Info, urn, fmt.Sprintf("deleted %s", file))

		return &pbempty.Empty{}, nil
	}

	if k.clusterUnreachable {
		_ = k.host.Log(ctx, diag.Warning, urn, fmt.Sprintf(
			"configured Kubernetes cluster is unreachable: %s", k.clusterUnreachableReason))
		if k.deleteUnreachable {
			_ = k.host.Log(ctx, diag.Info, urn, fmt.Sprintf(
				"configured Kubernetes cluster is unreachable and the `deleteUnreachable` option is enabled. "+
					"Deleting the unreachable resource from Pulumi state"))
			return &pbempty.Empty{}, nil
		}

		return nil, fmt.Errorf("configured Kubernetes cluster is unreachable. If the cluster was deleted, " +
			"you can remove this resource from Pulumi state by rerunning the operation with the " +
			"PULUMI_K8S_DELETE_UNREACHABLE environment variable set to \"true\"")
	}

	initialAPIVersion := initialAPIVersion(oldState, &unstructured.Unstructured{})
	fieldManager := k.fieldManagerName(nil, oldState, oldInputs)
	resources, err := k.getResources()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch OpenAPI schema from the API server: %w", err)
	}

	config := await.DeleteConfig{
		ProviderConfig: await.ProviderConfig{
			Context:           k.canceler.context,
			Host:              k.host,
			URN:               urn,
			InitialAPIVersion: initialAPIVersion,
			FieldManager:      fieldManager,
			ClusterVersion:    &k.k8sVersion,
			ClientSet:         k.clientSet,
			DedupLogger:       logging.NewLogger(k.canceler.context, k.host, urn),
			Resources:         resources,
			ServerSideApply:   k.serverSideApplyMode,
		},
		Inputs:  oldInputs,
		Outputs: current,
		Name:    name,
		Timeout: req.Timeout,
	}

	awaitErr := await.Deletion(config)
	if awaitErr != nil {
		if meta.IsNoMatchError(awaitErr) {
			// If it's a "no match" error, this is probably a CustomResource with no corresponding
			// CustomResourceDefinition. This usually happens if the CRD was deleted, and it's safe
			// to consider the CR to be deleted as well in this case.
			return &pbempty.Empty{}, nil
		}
		if kinds.IsPatchResource(urn, current.GetKind()) && await.IsDeleteRequiredFieldErr(awaitErr) {
			if cause, ok := apierrors.StatusCause(awaitErr, metav1.CauseTypeFieldValueRequired); ok {
				awaitErr = fmt.Errorf(
					"this Patch resource is currently managing a required field, so it can't be deleted "+
						"directly. Either set the `retainOnDelete` resource option, or transfer ownership of the "+
						"field before deleting: %s", cause.Field)
			}
		}
		partialErr, isPartialErr := awaitErr.(await.PartialError)
		if !isPartialErr {
			// There was an error executing the delete operation. The resource is still present and tracked.
			return nil, awaitErr
		}

		lastKnownState := partialErr.Object()

		inputsAndComputed, err := plugin.MarshalProperties(
			checkpointObject(current, lastKnownState, oldState, initialAPIVersion, fieldManager), plugin.MarshalOptions{
				Label:        fmt.Sprintf("%s.inputsAndComputed", label),
				KeepUnknowns: true,
				SkipNulls:    true,
				KeepSecrets:  k.enableSecrets,
			})
		if err != nil {
			return nil, err
		}

		// Resource delete was issued, but failed to complete. Return live version of object so it can be
		// checkpointed.
		return nil, partialError(fqObjName(lastKnownState), awaitErr, inputsAndComputed, nil)
	}

	return &pbempty.Empty{}, nil
}

// GetPluginInfo returns generic information about this plugin, like its version.
func (k *kubeProvider) GetPluginInfo(context.Context, *pbempty.Empty) (*pulumirpc.PluginInfo, error) {
	return &pulumirpc.PluginInfo{
		Version: k.version,
	}, nil
}

// Cancel signals the provider to gracefully shut down and abort any ongoing resource operations.
// Operations aborted in this way will return an error (e.g., `Update` and `Create` will either a
// creation error or an initialization error). Since Cancel is advisory and non-blocking, it is up
// to the host to decide how long to wait after Cancel is called before (e.g.)
// hard-closing any gRPC connection.
func (k *kubeProvider) Cancel(context.Context, *pbempty.Empty) (*pbempty.Empty, error) {
	k.canceler.cancel()
	return &pbempty.Empty{}, nil
}

// --------------------------------------------------------------------------

// Private helpers.

// --------------------------------------------------------------------------

func (k *kubeProvider) label() string {
	return fmt.Sprintf("Provider[%s]", k.name)
}

func (k *kubeProvider) gvkFromUnstructured(input *unstructured.Unstructured) schema.GroupVersionKind {
	var group, version, kind string

	kind = input.GetKind()
	gv := strings.Split(input.GetAPIVersion(), "/")
	if len(gv) == 1 {
		version = input.GetAPIVersion()
	} else {
		group, version = gv[0], gv[1]
	}
	if group == "core" {
		group = ""
	}

	return schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	}
}

func (k *kubeProvider) gvkFromURN(urn resource.URN) (schema.GroupVersionKind, error) {
	if string(urn.Type().Package()) != k.providerPackage {
		return schema.GroupVersionKind{}, fmt.Errorf("unrecognized resource type: %q for this provider",
			urn.Type())
	}

	// Emit GVK.
	kind := string(urn.Type().Name())
	gv := strings.Split(string(urn.Type().Module().Name()), "/")
	if len(gv) != 2 {
		return schema.GroupVersionKind{},
			fmt.Errorf("apiVersion does not have both a group and a version: %q", urn.Type().Module().Name())
	}
	group, version := gv[0], gv[1]
	if group == "core" {
		group = ""
	}

	return schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	}, nil
}

func (k *kubeProvider) readLiveObject(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	contract.Assertf(obj.GetName() != "", "expected object name to be nonempty: %v", obj)
	contract.Assertf(k.clientSet != nil, "expected Kubernetes client-set to be non-nil")
	rc, err := k.clientSet.ResourceClientForObject(obj)
	if err != nil {
		return nil, err
	}

	// Get the "live" version of the last submitted object. This is necessary because the server may
	// have populated some fields automatically, updated status fields, and so on.
	return rc.Get(k.canceler.context, obj.GetName(), metav1.GetOptions{})
}

// inputPatch calculates a patch on the client-side by comparing old inputs to the current inputs.
func (k *kubeProvider) inputPatch(
	oldInputs, newInputs *unstructured.Unstructured,
) ([]byte, error) {
	oldInputsJSON, err := oldInputs.MarshalJSON()
	if err != nil {
		return nil, err
	}
	newInputsJSON, err := newInputs.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return jsonpatch.CreateMergePatch(oldInputsJSON, newInputsJSON)
}

func (k *kubeProvider) isDryRunDisabledError(err error) bool {
	se, isStatusError := err.(*apierrors.StatusError)
	if !isStatusError {
		return false
	}

	return se.Status().Code == http.StatusBadRequest &&
		(se.Status().Message == "the dryRun alpha feature is disabled" ||
			se.Status().Message == "the dryRun beta feature is disabled" ||
			strings.Contains(se.Status().Message, "does not support dry run"))
}

// fieldManagerName returns the name to use for the Server-Side Apply fieldManager. The values are looked up with the
// following precedence:
// 1. Resource annotation (this will likely change to a typed option field in the next major release)
// 2. Value from the Pulumi state
// 3. Randomly generated name
func (k *kubeProvider) fieldManagerName(
	randomSeed []byte, state resource.PropertyMap, inputs *unstructured.Unstructured,
) string {
	// Always use the same fieldManager name for Client-Side Apply mode to avoid conflicts based on the name of the
	// provider executable.
	if !k.serverSideApplyMode {
		return "pulumi-kubernetes"
	}

	if v := metadata.GetAnnotationValue(inputs, metadata.AnnotationPatchFieldManager); len(v) > 0 {
		return v
	}
	if v, ok := state[fieldManagerKey]; ok {
		return v.StringValue()
	}

	prefix := "pulumi-kubernetes-"
	// This function is called from other provider function apart from Check and so doesn't have a randomSeed
	// for those calls, but the field manager name should have already been filled in via Check so this case
	// shouldn't actually get hit.
	fieldManager, err := resource.NewUniqueName(randomSeed, prefix, 0, 0, nil)
	contract.AssertNoErrorf(err, "unexpected error while creating NewUniqueName")

	return fieldManager
}

// gvkExists attempts to load a REST mapping for the given resource and returns true on success. Since this operation
// will fail if the GVK has not been registered with the apiserver, it can be used to indirectly check if the resource
// may be an unregistered CustomResource.
func (k *kubeProvider) gvkExists(obj *unstructured.Unstructured) bool {
	gvk := obj.GroupVersionKind()
	if k.clusterUnreachable {
		logger.V(3).Infof("gvkExists check failed due to unreachable cluster")
		return false
	}
	if _, err := k.clientSet.RESTMapper.RESTMapping(gvk.GroupKind(), gvk.Version); err != nil {
		if !meta.IsNoMatchError(err) {
			logger.V(3).Infof("RESTMapping(%q) returned unexpected error %v", gvk, err)
		}
		return false
	}
	return true
}

// loadPulumiConfig loads the PULUMI_CONFIG environment variable set by the engine, unmarshals the JSON string into
// a map, and returns the map and a bool indicating if the operation succeeded.
func (k *kubeProvider) loadPulumiConfig() (map[string]any, bool) {
	configStr, ok := os.LookupEnv("PULUMI_CONFIG")
	// PULUMI_CONFIG is not set on older versions of the engine, so check if the lookup succeeds.
	if !ok || configStr == "" {
		return nil, false
	}

	// PULUMI_CONFIG should be a JSON string that looks something like this:
	// {"enableServerSideApply":"true","kubeClientSettings":"{\"burst\":120,\"qps\":50}","strictMode":"true"}
	// The keys correspond to any project/stack config with a "kubernetes" prefix.
	var pConfig map[string]any
	err := json.Unmarshal([]byte(configStr), &pConfig)
	if err != nil {
		logger.V(3).Infof("failed to load provider config from PULUMI_CONFIG: %v", err)
		return nil, false
	}

	return pConfig, true
}

// removeLastAppliedConfigurationAnnotation is used by the Update method to remove an existing
// last-applied-configuration annotation from a resource. This annotation was set automatically by the provider, so it
// does not show up in the resource inputs. If the value is present in the live state, copy that value into the old
// inputs so that a negative diff will be generated for it.
func removeLastAppliedConfigurationAnnotation(oldLive, oldInputs *unstructured.Unstructured) {
	oldLiveValue, existsInOldLive, _ := unstructured.NestedString(oldLive.Object,
		"metadata", "annotations", lastAppliedConfigKey)
	_, existsInOldInputs, _ := unstructured.NestedString(oldInputs.Object,
		"metadata", "annotations", lastAppliedConfigKey)

	if existsInOldLive && !existsInOldInputs {
		contract.IgnoreError(unstructured.SetNestedField(
			oldInputs.Object, oldLiveValue, "metadata", "annotations", lastAppliedConfigKey))
	}
}

// pruneLiveState prunes a live resource object to match the shape of the input object that created the resource.
func pruneLiveState(live, oldInputs *unstructured.Unstructured) *unstructured.Unstructured {
	oldLivePruned := &unstructured.Unstructured{
		Object: pruneMap(live.Object, oldInputs.Object),
	}

	return oldLivePruned
}

// shouldNormalize returns false for CustomResources, and true otherwise.
func shouldNormalize(uns *unstructured.Unstructured) bool {
	return kinds.KnownGroupVersions.Has(uns.GetAPIVersion())
}

// normalizeInputs converts an Unstructured resource into a normalized form so that semantically equivalent representations
// are set to the same output shape. This is important to avoid generating diffs for inputs that will produce the same
// result on the cluster.
func normalizeInputs(uns *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	if !hasComputedValue(uns) && shouldNormalize(uns) {
		normalized, err := clients.Normalize(uns)
		if err != nil {
			return nil, err
		}
		uns = pruneLiveState(normalized, uns)
	}

	// Remove read-only fields
	// see: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#objectmeta-v1-meta
	unstructured.RemoveNestedField(uns.Object, "metadata", "creationTimestamp")
	unstructured.RemoveNestedField(uns.Object, "metadata", "deletionGracePeriodSeconds")
	unstructured.RemoveNestedField(uns.Object, "metadata", "deletionTimestamp")
	unstructured.RemoveNestedField(uns.Object, "metadata", "generation")
	unstructured.RemoveNestedField(uns.Object, "metadata", "managedFields")
	unstructured.RemoveNestedField(uns.Object, "metadata", "resourceVersion")
	unstructured.RemoveNestedField(uns.Object, "metadata", "uid")

	return uns, nil
}

func combineMapReplv(replvs ...func(resource.PropertyValue) (any, bool)) func(resource.PropertyValue) (any, bool) {
	return func(v resource.PropertyValue) (any, bool) {
		for _, replv := range replvs {
			if r, ok := replv(v); ok {
				return r, true
			}
		}
		return "", false
	}
}

func mapReplStripSecrets(v resource.PropertyValue) (any, bool) {
	if v.IsSecret() {
		return v.SecretValue().Element.MapRepl(nil, mapReplStripSecrets), true
	}

	return nil, false
}

func mapReplStripComputed(v resource.PropertyValue) (any, bool) {
	if v.IsComputed() {
		return nil, true
	}

	return nil, false
}

// underscoreToDashMap holds the mappings between underscore and dash keys.
var underscoreToDashMap = map[string]string{
	"x_kubernetes_embedded_resource":       "x-kubernetes-embedded-resource",
	"x_kubernetes_int_or_string":           "x-kubernetes-int-or-string",
	"x_kubernetes_list_map_keys":           "x-kubernetes-list-map-keys",
	"x_kubernetes_list_type":               "x-kubernetes-list-type",
	"x_kubernetes_map_type":                "x-kubernetes-map-type",
	"x_kubernetes_preserve_unknown_fields": "x-kubernetes-preserve-unknown-fields",
	"x_kubernetes_validations":             "x-kubernetes-validations",
}

// dashedToUnderscoreMap holds the reverse mappings between dash and underscore keys. This
// is a precomputed map based on underscoreToDashMap at runtime to avoid duplicating
// code, or extra passes over the map.
var dashToUnderscoreMap map[string]string = func() map[string]string {
	dashToUnderscoreMap := make(map[string]string, len(underscoreToDashMap))
	for k, v := range underscoreToDashMap {
		dashToUnderscoreMap[v] = k
	}
	return dashToUnderscoreMap
}()

// mapReplUnderscoreToDash denormalizes keys by replacing underscores with dashes.
func mapReplUnderscoreToDash(v string) (string, bool) {
	val, ok := underscoreToDashMap[v]
	return val, ok
}

// mapReplDashToUnderscore normalizes keys by replacing dashes with underscores.
func mapReplDashToUnderscore(v string) (resource.PropertyKey, bool) {
	val, ok := dashToUnderscoreMap[v]
	return resource.PropertyKey(val), ok
}

// propMapToUnstructured converts a resource.PropertyMap to an *unstructured.Unstructured; and applies field name denormalization
// and secret stripping.
func propMapToUnstructured(pm resource.PropertyMap) *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: pm.MapRepl(mapReplUnderscoreToDash, mapReplStripSecrets)}
}

// initialAPIVersion retrieves the initialAPIVersion property from the checkpoint file and falls back to using
// the version from the resource metadata if that property is not present.
func initialAPIVersion(state resource.PropertyMap, oldInputs *unstructured.Unstructured) string {
	if v, ok := state[initialAPIVersionKey]; ok {
		return v.StringValue()
	}

	return oldInputs.GetAPIVersion()
}

// checkpointObject generates a checkpointed PropertyMap from the live and input Kubernetes objects.
// It normalizes `x-kubernetes-*` fields to their underscored equivalents, handles secret data annotations,
// processes `stringData` for secret kinds by marking corresponding `data` fields as secrets,
// and includes metadata such as the initial API version and field manager.
func checkpointObject(inputs, live *unstructured.Unstructured, fromInputs resource.PropertyMap,
	initialAPIVersion, fieldManager string,
) resource.PropertyMap {
	// When checkpointing the live object, we need to ensure we normalize any `x-kubernetes-*` fields to their
	// underscored versions so they can be correctly diffed, and deseriazlied to their typed SDK equivalents.
	object := resource.NewPropertyMapFromMapRepl(live.Object, mapReplDashToUnderscore, nil)
	inputsPM := resource.NewPropertyMapFromMapRepl(inputs.Object, mapReplDashToUnderscore, nil)

	annotateSecrets(object, fromInputs)
	annotateSecrets(inputsPM, fromInputs)

	isSecretKind := live.GetKind() == secretKind

	// For secrets, if `stringData` is present in the inputs, the API server will have filled in `data` based on it. By
	// base64 encoding the secrets. We should mark any of the values which were secrets in the `stringData` object
	// as secrets in the `data` field as well.
	if live.GetAPIVersion() == "v1" && isSecretKind {
		stringData, hasStringData := fromInputs["stringData"]
		data, hasData := object["data"]

		if hasStringData && hasData {
			if stringData.IsSecret() && !data.IsSecret() {
				object["data"] = resource.MakeSecret(data)
			}

			if stringData.IsObject() && data.IsObject() {
				annotateSecrets(data.ObjectValue(), stringData.ObjectValue())
			}
		}
	}

	object["__inputs"] = resource.NewObjectProperty(inputsPM)
	object[initialAPIVersionKey] = resource.NewStringProperty(initialAPIVersion)
	object[fieldManagerKey] = resource.NewStringProperty(fieldManager)

	return object
}

// parseCheckpointObject parses the given resource.PropertyMap, stripping sensitive information and normalizing field names.
// It returns two unstructured.Unstructured objects: oldInputs containing the input properties and live containing the live state.
func parseCheckpointObject(obj resource.PropertyMap) (oldInputs, live *unstructured.Unstructured) {
	// Since we are converting everything to unstructured's, we need to strip out any secretness that
	// may nested deep within the object.
	// Note: we also handle conversion of underscored `x_kubernetes_*` fields to their respective dashed
	// versions here.
	pm := obj.MapRepl(mapReplUnderscoreToDash, mapReplStripSecrets)

	//
	// NOTE: Inputs are now stored in `__inputs` to allow output properties to work. The inputs and
	// live properties used to be stored next to each other, in an object that looked like {live:
	// (...), inputs: (...)}, but this broke this resolution. See[1] for more information.
	//
	// [1]: https://github.com/pulumi/pulumi-kubernetes/issues/137
	//
	inputs, hasInputs := pm["inputs"]
	liveMap, hasLive := pm["live"]

	if !hasInputs || !hasLive {
		liveMap = pm

		inputs, hasInputs = pm["__inputs"]
		if hasInputs {
			delete(liveMap.(map[string]any), "__inputs")
		} else {
			inputs = map[string]any{}
		}
	}

	oldInputs = &unstructured.Unstructured{Object: inputs.(map[string]any)}
	live = &unstructured.Unstructured{Object: liveMap.(map[string]any)}
	return oldInputs, live
}

// partialError creates an error for resources that did not complete an operation in progress.
// The last known state of the object is included in the error so that it can be checkpointed.
func partialError(id string, err error, state *structpb.Struct, inputs *structpb.Struct) error {
	reasons := []string{err.Error()}
	var aggErr await.AggregatedError
	if errors.As(err, &aggErr) {
		reasons = append(reasons, aggErr.SubErrors()...)
	}
	detail := pulumirpc.ErrorResourceInitFailed{
		Id:         id,
		Properties: state,
		Reasons:    reasons,
		Inputs:     inputs,
	}
	return rpcerror.WithDetails(rpcerror.New(codes.Unknown, err.Error()), &detail)
}

// canonicalNamespace will provides the canonical name for a namespace. Specifically, if the
// namespace is "", the empty string, we report this as its canonical name, "default".
func canonicalNamespace(ns string) string {
	if ns == "" {
		return "default"
	}
	return ns
}

// deleteResponse causes the resource to be deleted from the state.
var deleteResponse = &pulumirpc.ReadResponse{Id: "", Properties: nil}

// convertPatchToDiff converts the given JSON merge patch to a Pulumi detailed diff.
func convertPatchToDiff(
	patch, oldLiveState, newInputs, oldInputs map[string]any, forceNewFields ...string,
) (map[string]*pulumirpc.PropertyDiff, error) {
	contract.Requiref(len(patch) != 0, "patch", "expected len() != 0")
	contract.Requiref(oldLiveState != nil, "oldLiveState", "expected != nil")

	pc := &patchConverter{
		forceNew: forceNewFields,
		diff:     map[string]*pulumirpc.PropertyDiff{},
	}
	err := pc.addPatchMapToDiff(nil, patch, oldLiveState, newInputs, oldInputs, false)
	return pc.diff, err
}

// makePatchSlice recursively processes the given path to create a slice of a POJO value that is appropriately shaped
// for querying using a JSON path. We use this in addPatchValueToDiff when deciding whether or not a particular
// property causes a replacement.
func makePatchSlice(path []any, v any) any {
	if len(path) == 0 {
		return v
	}
	switch p := path[0].(type) {
	case string:
		return map[string]any{
			p: makePatchSlice(path[1:], v),
		}
	case int:
		return []any{makePatchSlice(path[1:], v)}
	default:
		contract.Failf("unexpected element type in path: %T", p)
		return nil
	}
}

// equalNumbers returns true if both a and b are number values (int64 or float64). Note that if a this will fail if
// either value is not representable as a float64.
func equalNumbers(a, b any) bool {
	aKind, bKind := reflect.TypeOf(a).Kind(), reflect.TypeOf(b).Kind()
	if aKind == bKind {
		return reflect.DeepEqual(a, b)
	}

	toFloat := func(v any) (float64, bool) {
		switch field := v.(type) {
		case int64:
			return float64(field), true
		case float64:
			return field, true
		default:
			return 0, false
		}
	}

	aVal, aOk := toFloat(a)
	bVal, bOk := toFloat(b)
	return aOk && bOk && aVal == bVal
}

// patchConverter carries context for convertPatchToDiff.
type patchConverter struct {
	forceNew []string
	diff     map[string]*pulumirpc.PropertyDiff

	// missing is a placeholder used during diffing to distinguish between
	// `nil` values and absent ones.
	missing struct{}
}

// addPatchValueToDiff adds the given patched value to the detailed diff.
//
// Values for old, newInput, and oldInput should be `pc.missing` if they were
// originally absent from a map. This differenciates between the case where
// they were present in the map but had `nil` values.
//
// The diff that is recorded depends on the old and new values:
// - If the patched value is nil, the property is recorded as deleted.
// - If the old value is missing, the property is recorded as added.
// - If the old and patched values are both nil or missing, no diff is recorded.
// - If the types of the old and new values differ, the property is recorded as updated.
// - If both values are maps, the maps are recursively compared on a per-property basis and added to the diff.
// - If both values are arrays, the arrays are recursively compared on a per-element basis and added to the diff.
// - If both values are primitives and the values differ, the property is recorded as updated.
// - Otherwise, no diff is recorded.
//
// If a difference is present at the given path and the path matches one of the patterns in the database of
// force-new properties, the diff is amended to indicate that the resource needs to be replaced due to the change in
// this property.
func (pc *patchConverter) addPatchValueToDiff(
	path []any, v, old, newInput, oldInput any, inArray bool,
) error {
	contract.Assertf(v != nil || old != nil || oldInput != nil || newInput != nil,
		"path: %+v  |  v: %+v  | old: %+v  |  oldInput: %+v  |  newInput: %+v",
		path, v, old, oldInput, newInput)

	// If there is no new input, then the only possible diff here is a delete. All other diffs must be diffs between
	// old and new properties that are populated by the server. If there is also no old input, then there is no diff
	// whatsoever.
	if newInput == pc.missing && (v != nil || oldInput == pc.missing) {
		return nil
	}

	var diffKind pulumirpc.PropertyDiff_Kind
	inputDiff := false
	if v == nil {
		// computed values are rendered as null in the patch; handle this special case.
		if _, ok := newInput.(resource.Computed); ok {
			if old == pc.missing {
				diffKind = pulumirpc.PropertyDiff_ADD
			} else {
				diffKind = pulumirpc.PropertyDiff_UPDATE
			}
		} else {
			diffKind, inputDiff = pulumirpc.PropertyDiff_DELETE, true
		}
	} else if old == pc.missing {
		diffKind = pulumirpc.PropertyDiff_ADD
	} else {
		switch v := v.(type) {
		case map[string]any:
			if oldMap, ok := old.(map[string]any); ok {
				newInputMap, _ := newInput.(map[string]any)
				oldInputMap, _ := oldInput.(map[string]any)
				return pc.addPatchMapToDiff(path, v, oldMap, newInputMap, oldInputMap, inArray)
			}
			diffKind = pulumirpc.PropertyDiff_UPDATE
		case []any:
			if oldArray, ok := old.([]any); ok {
				newInputArray, _ := newInput.([]any)
				oldInputArray, _ := oldInput.([]any)
				return pc.addPatchArrayToDiff(path, v, oldArray, newInputArray, oldInputArray, inArray)
			}
			diffKind = pulumirpc.PropertyDiff_UPDATE
		default:
			if reflect.DeepEqual(v, old) || equalNumbers(v, old) {
				// From RFC 7386 (the JSON Merge Patch spec):
				//
				//   If the patch is anything other than an object, the result will always be to replace the entire
				//   target with the entire patch. Also, it is not possible to patch part of a target that is not an
				//   object, such as to replace just some of the values in an array.
				//
				// Because JSON merge patch does not allow array elements to be updated--instead, the array must be
				// replaced in full--the patch we have is an overestimate of the properties that changed. As such, we
				// only record updates for values that have in fact changed.
				return nil
			}
			diffKind = pulumirpc.PropertyDiff_UPDATE
		}
	}

	// Determine if this change causes a replace.
	matches, err := openapi.PatchPropertiesChanged(makePatchSlice(path, v).(map[string]any), pc.forceNew)
	if err != nil {
		return err
	}
	if len(matches) != 0 {
		switch diffKind {
		case pulumirpc.PropertyDiff_ADD:
			diffKind = pulumirpc.PropertyDiff_ADD_REPLACE
		case pulumirpc.PropertyDiff_DELETE:
			diffKind = pulumirpc.PropertyDiff_DELETE_REPLACE
		case pulumirpc.PropertyDiff_UPDATE:
			diffKind = pulumirpc.PropertyDiff_UPDATE_REPLACE
		}
	}

	pathStr := ""
	for _, v := range path {
		switch v := v.(type) {
		case string:
			if strings.ContainsAny(v, `."[]`) {
				pathStr = fmt.Sprintf(`%s["%s"]`, pathStr, strings.ReplaceAll(v, `"`, `\"`))
			} else if pathStr != "" {
				pathStr = fmt.Sprintf("%s.%s", pathStr, v)
			} else {
				pathStr = v
			}
		case int:
			pathStr = fmt.Sprintf("%s[%d]", pathStr, v)
		}
	}

	pc.diff[pathStr] = &pulumirpc.PropertyDiff{Kind: diffKind, InputDiff: inputDiff}
	return nil
}

// get will return the map's value for the given key, or the `missing`
// placeholder if no such key exists in the map.
func (pc *patchConverter) get(m map[string]any, k string) any {
	if v, ok := m[k]; ok {
		return v
	}
	return pc.missing
}

// addPatchMapToDiff adds the diffs in the given patched map to the detailed diff.
//
// If this map is contained within an array, we do a little bit more work to detect deletes, as they are not recorded
// in the patch in this case (see the note in addPatchValueToDiff for more details).
func (pc *patchConverter) addPatchMapToDiff(
	path []any, m, old, newInput, oldInput map[string]any, inArray bool,
) error {
	for k, v := range m {
		if err := pc.addPatchValueToDiff(append(path, k), v, pc.get(old, k), pc.get(newInput, k), pc.get(oldInput, k), inArray); err != nil {
			return err
		}
	}
	if inArray {
		for k, v := range old {
			if _, ok := m[k]; ok {
				continue
			}
			if err := pc.addPatchValueToDiff(append(path, k), nil, v, pc.get(newInput, k), pc.get(oldInput, k), inArray); err != nil {
				return err
			}
		}
	}
	return nil
}

// addPatchArrayToDiff adds the diffs in the given patched array to the detailed diff.
func (pc *patchConverter) addPatchArrayToDiff(
	path []any, a, old, newInput, oldInput []any, inArray bool,
) error {
	at := func(arr []any, i int) any {
		if i < len(arr) {
			return arr[i]
		}
		return nil
	}

	var i int
	for i = 0; i < len(a) && i < len(old); i++ {
		err := pc.addPatchValueToDiff(append(path, i), a[i], old[i], at(newInput, i), at(oldInput, i), true)
		if err != nil {
			return err
		}
	}

	if i < len(a) {
		for ; i < len(a); i++ {
			err := pc.addPatchValueToDiff(append(path, i), a[i], pc.missing, at(newInput, i), at(oldInput, i), true)
			if err != nil {
				return err
			}
		}
	} else {
		for ; i < len(old); i++ {
			err := pc.addPatchValueToDiff(append(path, i), nil, old[i], at(newInput, i), at(oldInput, i), true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// annotateSecrets copies the "secretness" from the ins to the outs. If there are values with the same keys for the
// outs and the ins, if they are both objects, they are transformed recursively. Otherwise, if the value in the ins
// contains a secret, the entire out value is marked as a secret.  This is very close to how we project secrets
// in the programming model, with one small difference, which is how we treat the case where both are objects. In the
// programming model, we would say the entire output object is a secret. Here, we actually recur in. We do this because
// we don't want a single secret value in a rich structure to taint the entire object. Doing so would mean things like
// the entire value in the deployment would be encrypted instead of a small chunk. It also means the entire property
// would be displayed as `[secret]` in the CLI instead of a small part.
//
// NOTE: This means that for an array, if any value in the input version is a secret, the entire output array is
// marked as a secret. This is actually a very nice result, because often arrays are treated like sets by providers
// and the order may not be preserved across an operation. This means we do end up encrypting the entire array
// but that's better than accidentally leaking a value which just moved to a different location.
func annotateSecrets(outs, ins resource.PropertyMap) {
	if outs == nil {
		return
	}

	// Note on IsString(): k8s "kind" is always a string, but we might be recursing into the structure and there
	// might be another, unrelated "kind" in there.
	if kind, ok := outs["kind"]; ok && kind.IsString() && kind.StringValue() == secretKind {
		if data, hasData := outs["data"]; hasData {
			outs["data"] = resource.MakeSecret(data)
		}
		if stringData, hasStringData := outs["stringData"]; hasStringData {
			outs["stringData"] = resource.MakeSecret(stringData)
		}
		return
	}

	if ins == nil {
		return
	}

	for key, inValue := range ins {
		outValue, has := outs[key]
		if !has {
			continue
		}
		if outValue.IsObject() && inValue.IsObject() {
			annotateSecrets(outValue.ObjectValue(), inValue.ObjectValue())
		} else if !outValue.IsSecret() && inValue.ContainsSecrets() {
			outs[key] = resource.MakeSecret(outValue)
		}
	}
}

// renderYaml marshals an Unstructured resource to YAML and writes it to the specified path on disk or returns an error.
func renderYaml(resource *unstructured.Unstructured, yamlDirectory string) error {
	jsonBytes, err := resource.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to render YAML file: %q: %w", yamlDirectory, err)
	}
	yamlBytes, err := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		return fmt.Errorf("failed to render YAML file: %q: %w", yamlDirectory, err)
	}

	crdDirectory := filepath.Join(yamlDirectory, "0-crd")
	manifestDirectory := filepath.Join(yamlDirectory, "1-manifest")

	if _, err := os.Stat(crdDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(crdDirectory, 0o700)
		if err != nil {
			return fmt.Errorf("failed to create directory for rendered YAML: %q: %w", crdDirectory, err)
		}
	}
	if _, err := os.Stat(manifestDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(manifestDirectory, 0o700)
		if err != nil {
			return fmt.Errorf("failed to create directory for rendered YAML: %q: %w", manifestDirectory, err)
		}
	}

	path := renderPathForResource(resource, yamlDirectory)
	err = os.WriteFile(path, yamlBytes, 0o600)
	if err != nil {
		return fmt.Errorf("failed to write YAML file: %q: %w", path, err)
	}

	return nil
}

// renderPathForResource determines the appropriate YAML render path depending on the resource kind.
func renderPathForResource(resource *unstructured.Unstructured, yamlDirectory string) string {
	contract.Assertf(resource.GetName() != "", "expected object name to be nonempty: %v", resource)
	crdDirectory := filepath.Join(yamlDirectory, "0-crd")
	manifestDirectory := filepath.Join(yamlDirectory, "1-manifest")

	namespace := "default"
	if "" != resource.GetNamespace() {
		namespace = resource.GetNamespace()
	}

	sanitise := func(name string) string {
		name = strings.NewReplacer("/", "_", ":", "_").Replace(name)
		return name
	}

	fileName := fmt.Sprintf("%s-%s-%s-%s.yaml", sanitise(resource.GetAPIVersion()), strings.ToLower(resource.GetKind()), namespace, resource.GetName())
	filepath.Join(yamlDirectory, fileName)

	var path string
	if kinds.KnownGroupVersions.Has(resource.GetAPIVersion()) && kinds.Kind(resource.GetKind()) == kinds.CustomResourceDefinition {
		path = filepath.Join(crdDirectory, fileName)
	} else {
		path = filepath.Join(manifestDirectory, fileName)
	}

	return path
}
