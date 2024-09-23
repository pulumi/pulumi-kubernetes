// Copyright 2016-2022, Pulumi Corporation.
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
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/helm"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/host"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"helm.sh/helm/v3/pkg/action"
	helmchart "helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/postrender"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/yaml"
)

// Default timeout for awaited install and uninstall operations
const defaultTimeoutSeconds = 300

// errReleaseNotFound is the error when a Helm release is not found
var errReleaseNotFound = errors.New("release not found")

// Release should explicitly track the shape of helm.sh/v3:Release resource
type Release struct {
	// When combinging Values with mergeMaps, allow Nulls
	AllowNullValues bool `json:"allowNullValues,omitempty"`
	// If set, installation process purges chart on fail. The wait flag will be set automatically if atomic is used
	Atomic bool `json:"atomic,omitempty"`
	// Chart name to be installed. A path may be used.
	Chart string `json:"chart,omitempty"`
	// Allow deletion of new resources created in this upgrade when upgrade fails
	CleanupOnFail bool `json:"cleanupOnFail,omitempty"`
	// Create the namespace if it does not exist
	CreateNamespace bool `json:"createNamespace,omitempty"`
	// Run helm dependency update before installing the chart
	DependencyUpdate bool `json:"dependencyUpdate,omitempty"`
	// Add a custom description
	Description string `json:"description,omitempty"`
	// Use chart development versions, too. Equivalent to version '>0.0.0-0'. If `version` is set, this is ignored
	Devel bool `json:"devel,omitempty"`
	// Prevent CRD hooks from running, but run other hooks.  See helm install --no-crd-hook
	DisableCRDHooks bool `json:"disableCRDHooks,omitempty"`
	// If set, the installation process will not validate rendered templates against the Kubernetes OpenAPI Schema
	DisableOpenapiValidation bool `json:"disableOpenapiValidation,omitempty"`
	// Prevent hooks from running.
	DisableWebhooks bool `json:"disableWebhooks,omitempty"`
	// Force resource update through delete/recreate if needed.
	ForceUpdate bool `json:"forceUpdate,omitempty"`
	// Location of public keys used for verification. Used only if `verify` is true
	Keyring string `json:"keyring,omitempty"`
	// Run helm lint when planning
	Lint bool `json:"lint,omitempty"`
	// Limit the maximum number of revisions saved per release. Use 0 for no limit
	MaxHistory *int `json:"maxHistory,omitempty"`
	// Release name.
	Name string `json:"name,omitempty"`
	// Namespace to install the release into.
	Namespace string `json:"namespace,omitempty"`
	// Postrender command to run.
	Postrender string `json:"postrender,omitempty"`
	// Perform pods restart during upgrade/rollback
	RecreatePods bool `json:"recreatePods,omitempty"`
	// If set, render subchart notes along with the parent
	RenderSubchartNotes bool `json:"renderSubchartNotes,omitempty"`
	// Re-use the given name, even if that name is already used. This is unsafe in production
	Replace bool `json:"replace,omitempty"`
	// Specification defining the Helm chart repository to use.
	RepositoryOpts *RepositoryOpts `json:"repositoryOpts,omitempty"`
	// When upgrading, reset the values to the ones built into the chart
	ResetValues bool `json:"resetValues,omitempty"`
	// When upgrading, reuse the last release's values and merge in any overrides. If 'reset_values' is specified, this is ignored
	ReuseValues bool `json:"reuseValues,omitempty"`
	// Custom values to be merged with items loaded from values.
	Values map[string]any `json:"values,omitempty"`
	// If set, no CRDs will be installed. By default, CRDs are installed if not already present
	SkipCrds bool `json:"skipCrds,omitempty"`
	// Time in seconds to wait for any individual kubernetes operation.
	Timeout int `json:"timeout,omitempty"`
	// Verify the package before installing it.
	Verify bool `json:"verify,omitempty"`
	// Specify the exact chart version to install. If this is not specified, the latest version is installed.
	Version string `json:"version,omitempty"`
	// By default, the provider waits until all resources are in a ready state before marking the release as successful. Setting this to true will skip such await logic.
	SkipAwait bool `json:"skipAwait,omitempty"`
	// Will wait until all Jobs have been completed before marking the release as successful. This is ignored if `skipAwait` is enabled.
	WaitForJobs bool `json:"waitForJobs,omitempty"`
	// The rendered manifests.
	// Manifest map[string]interface{} `json:"manifest,omitempty"`
	// Names of resources created by the release grouped by "kind/version".
	ResourceNames map[string][]string `json:"resourceNames,omitempty"`
	// Status of the deployed release.
	Status *ReleaseStatus `json:"status,omitempty"`
}

type ReleaseSpec struct{}

// Specification defining the Helm chart repository to use.
type RepositoryOpts struct {
	// Repository where to locate the requested chart. If it's a URL the chart is installed without installing the repository.
	Repo string `json:"repo,omitempty"`
	// The Repositories CA File
	CAFile string `json:"caFile,omitempty"`
	// The repositories cert file
	CertFile string `json:"certFile,omitempty"`
	// The repositories cert key file
	KeyFile string `json:"keyFile,omitempty"`
	// Password for HTTP basic authentication
	Password string `json:"password,omitempty"`
	// Username for HTTP basic authentication
	Username string `json:"username,omitempty"`
}

type ReleaseStatus struct {
	// The version number of the application being deployed.
	AppVersion string `json:"appVersion,omitempty"`
	// The name of the chart.
	Chart string `json:"chart,omitempty"`
	// Name is the name of the release.
	Name string `json:"name,omitempty"`
	// Namespace is the kubernetes namespace of the release.
	Namespace string `json:"namespace,omitempty"`
	// Version is an int32 which represents the version of the release.
	Revision *int `json:"revision,omitempty"`
	// Status of the release.
	Status string `json:"status,omitempty"`
	// A SemVer 2 conformant version string of the chart.
	Version string `json:"version,omitempty"`
}

type helmReleaseProvider struct {
	host                     host.HostClient
	canceler                 *cancellationContext
	helmDriver               string
	apiConfig                *api.Config
	defaultOverrides         *clientcmd.ConfigOverrides
	restConfig               *rest.Config
	clientSet                *clients.DynamicClientSet
	defaultNamespace         string
	enableSecrets            bool
	clusterUnreachable       bool
	clusterUnreachableReason string
	name                     string
	settings                 *cli.EnvSettings
}

func newHelmReleaseProvider(
	host host.HostClient,
	canceler *cancellationContext,
	apiConfig *api.Config,
	defaultOverrides *clientcmd.ConfigOverrides,
	restConfig *rest.Config,
	clientSet *clients.DynamicClientSet,
	helmDriver,
	namespace string,
	enableSecrets bool,
	pluginsDirectory,
	registryConfigPath,
	repositoryConfigPath,
	repositoryCache string,
	clusterUnreachable bool,
	clusterUnreachableReason string,
) (customResourceProvider, error) {
	settings := cli.New()
	settings.PluginsDirectory = pluginsDirectory
	settings.RegistryConfig = registryConfigPath
	settings.RepositoryConfig = repositoryConfigPath
	settings.RepositoryCache = repositoryCache
	settings.Debug = true

	return &helmReleaseProvider{
		host:                     host,
		canceler:                 canceler,
		apiConfig:                apiConfig,
		defaultOverrides:         defaultOverrides,
		restConfig:               restConfig,
		clientSet:                clientSet,
		helmDriver:               helmDriver,
		defaultNamespace:         namespace,
		enableSecrets:            enableSecrets,
		clusterUnreachable:       clusterUnreachable,
		clusterUnreachableReason: clusterUnreachableReason,
		name:                     "kubernetes:helmrelease",
		settings:                 settings,
	}, nil
}

func debug(format string, a ...any) {
	logger.V(6).Infof("[DEBUG] %s", fmt.Sprintf(format, a...))
}

func (r *helmReleaseProvider) getActionConfig(namespace string) (*action.Configuration, error) {
	conf := new(action.Configuration)
	var overrides clientcmd.ConfigOverrides
	if r.defaultOverrides != nil {
		overrides = *r.defaultOverrides
	}

	// This essentially points the client to use the specified namespace when a namespaced
	// object doesn't have the namespace specified. This allows us to interpolate the
	// release's namespace as the default namespace on charts with templates that don't
	// explicitly set the namespace (e.g. through namespace: {{ .Release.Namespace }}).
	overrides.Context.Namespace = namespace

	var clientConfig clientcmd.ClientConfig
	if r.apiConfig != nil {
		clientConfig = clientcmd.NewDefaultClientConfig(*r.apiConfig, &overrides)
	} else {
		// Use client-go to resolve the final configuration values for the client. Typically these
		// values would reside in the $KUBECONFIG file, but can also be altered in several
		// places, including in env variables, client-go default values, and (if we allowed it) CLI
		// flags.
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
		clientConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &overrides)
	}
	kc := NewKubeConfig(r.restConfig, clientConfig)

	if err := conf.Init(kc, namespace, r.helmDriver, debug); err != nil {
		return nil, err
	}
	logger.V(9).Infof("Setting registry client with config file: %q and debug: %v", r.settings.RegistryConfig,
		r.settings.Debug)
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(r.settings.Debug),
		registry.ClientOptCredentialsFile(r.settings.RegistryConfig),
	)
	if err != nil {
		return nil, err
	}
	conf.RegistryClient = registryClient
	return conf, nil
}

// mapReplExtractValues extracts pure values from the property map.
var mapReplExtractValues = combineMapReplv(mapReplStripSecrets, mapReplStripComputed)

func decodeRelease(pm resource.PropertyMap, label string) (*Release, error) {
	var release Release
	values := map[string]any{}
	stripped := pm.MapRepl(nil, mapReplExtractValues)
	logger.V(9).Infof("[%s] Decoding release: %#v", label, stripped)

	if v, ok := stripped["valueYamlFiles"]; ok {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice, reflect.Array:
			s := reflect.ValueOf(v)
			for i := 0; i < s.Len(); i++ {
				val := s.Index(i).Interface()
				switch t := val.(type) {
				case *resource.Asset:
					b, err := t.Bytes()
					if err != nil {
						return nil, err
					}
					valuesMap := map[string]any{}
					if err = yaml.Unmarshal(b, &valuesMap); err != nil {
						return nil, err
					}
					values = helm.MergeMaps(values, valuesMap)
				default:
					return nil, fmt.Errorf("unsupported type for 'valueYamlFiles' arg: %T", v)
				}
			}
		}
	}

	var err error
	if err = mapstructure.Decode(stripped, &release); err != nil {
		return nil, fmt.Errorf("decoding failure: %w", err)
	}
	release.Values, err = mergeMaps(values, release.Values, release.AllowNullValues)
	if err != nil {
		return nil, err
	}
	return &release, nil
}

func (r *helmReleaseProvider) Check(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("Provider[%s].Check(%s)", r.name, urn)

	var failures []*pulumirpc.CheckFailure
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

	if len(olds) > 0 {
		adoptOldNameIfUnnamed(news, olds)
	}
	assignNameIfAutonameable(news, urn)
	r.setDefaults(news)

	logger.V(9).Infof("Decoding new release.")
	new, err := decodeRelease(news, fmt.Sprintf("%s.news", label))
	if err != nil {
		return nil, err
	}

	if !news.ContainsUnknowns() {
		logger.V(9).Infof("Loading Helm chart.")
		chart, err := r.helmLoad(ctx, urn, new)
		if err != nil {
			failures = append(failures, &pulumirpc.CheckFailure{
				Property: "chart",
				Reason:   fmt.Sprintf("%v; check the chart name and repository configuration.", err),
			})
		} else {
			// determine the desired state of the resource, i.e the specific chart version
			// as opposed to the program input (which is a constraint such as ">= 1.2.3").
			// with this we may determine whether the Helm release needs to be upgraded.
			new.Version = chart.Metadata.Version
		}
	}

	logger.V(9).Infof("New: %+v", new)
	news = resource.NewPropertyMap(new)

	// remove deprecated inputs
	delete(news, "resourceNames")

	newInputs, err := plugin.UnmarshalProperties(newResInputs, plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.newInputs", label),
		KeepUnknowns: true,
		SkipNulls:    true,
		KeepSecrets:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("check failed because malformed resource inputs: %w", err)
	}
	// ensure we don't leak secrets into state, and preserve the computedness of inputs.
	annotateComputed(news, newInputs)
	annotateSecrets(news, newInputs)

	autonamedInputs, err := plugin.MarshalProperties(news, plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.autonamedInputs", label),
		KeepUnknowns: true,
		SkipNulls:    true,
		KeepSecrets:  r.enableSecrets,
	})
	if err != nil {
		return nil, err
	}

	// Return new, possibly-autonamed inputs.
	return &pulumirpc.CheckResponse{Inputs: autonamedInputs, Failures: failures}, nil
}

func (r *helmReleaseProvider) setDefaults(target resource.PropertyMap) {
	namespace, ok := target["namespace"]
	if !ok || (namespace.IsString() && namespace.StringValue() == "") {
		target["namespace"] = resource.NewStringProperty(r.defaultNamespace)
	}

	skipAwaitVal, ok := target["skipAwait"]
	if !ok || (skipAwaitVal.IsBool() && !skipAwaitVal.BoolValue()) {
		// If timeout is specified (even if zero value), use that. Otherwise use default.
		_, has := target["timeout"]
		if !has {
			target["timeout"] = resource.NewNumberProperty(defaultTimeoutSeconds)
		}
	}

	// Discover the keyring if chart verification is requested, and a keyring is not explicitly specified.
	verify, ok := target["verify"]
	if ok && verify.IsBool() && verify.BoolValue() {
		keyringVal, ok := target["keyring"]
		if !ok || (keyringVal.IsString() && keyringVal.StringValue() == "") {
			target["keyring"] = resource.NewStringProperty(os.ExpandEnv("$HOME/.gnupg/pubring.gpg"))
		}
	}
}

func (r *helmReleaseProvider) helmLoad(ctx context.Context, urn resource.URN, newRelease *Release) (*helmchart.Chart, error) {
	conf, err := r.getActionConfig(newRelease.Namespace)
	if err != nil {
		return nil, err
	}
	client := action.NewInstall(conf)
	c, path, err := getChart(&client.ChartPathOptions, conf.RegistryClient, r.settings, newRelease)
	if err != nil {
		logger.V(9).Infof("getChart failed: %v", err)
		logger.V(9).Infof("Settings: %#v", r.settings)
		return nil, err
	}

	logger.V(9).Infof("Checking chart dependencies for chart: %q with path: %q", newRelease.Chart, path)

	// check and update the chart's dependencies if needed
	updated, err := checkChartDependencies(
		c,
		path,
		newRelease.Keyring,
		r.settings,
		conf.RegistryClient,
		newRelease.DependencyUpdate)
	if err != nil {
		return nil, err
	} else if updated {
		// load the chart again if its dependencies have been updated
		c, err = loader.Load(path)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (r *helmReleaseProvider) helmCreate(ctx context.Context, urn resource.URN, newRelease *Release) error {
	conf, err := r.getActionConfig(newRelease.Namespace)
	if err != nil {
		return err
	}
	client := action.NewInstall(conf)
	c, path, err := getChart(&client.ChartPathOptions, conf.RegistryClient, r.settings, newRelease)
	if err != nil {
		logger.V(9).Infof("getChart failed: %+v", err)
		logger.V(9).Infof("Settings: %#v", r.settings)
		return err
	}

	logger.V(9).Infof("Checking chart dependencies for chart: %q with path: %q", newRelease.Chart, path)
	// check and update the chart's dependencies if needed
	updated, err := checkChartDependencies(
		c,
		path,
		newRelease.Keyring,
		r.settings,
		conf.RegistryClient,
		newRelease.DependencyUpdate)
	if err != nil {
		return err
	} else if updated {
		// load the chart again if its dependencies have been updated
		c, err = loader.Load(path)
		if err != nil {
			return err
		}
	}

	logger.V(9).Infof("Fetching values for release: %q", newRelease.Name)
	values, err := getValues(newRelease)
	if err != nil {
		return err
	}

	logger.V(9).Infof("Values: %+v", values)

	err = isChartInstallable(c)
	if err != nil {
		return err
	}

	client.ClientOnly = false
	client.DisableHooks = newRelease.DisableWebhooks
	client.Wait = !newRelease.SkipAwait
	client.WaitForJobs = !newRelease.SkipAwait && newRelease.WaitForJobs
	client.Devel = newRelease.Devel
	client.DependencyUpdate = newRelease.DependencyUpdate
	client.Timeout = getTimeoutOrDefault(newRelease.Timeout)
	client.Namespace = newRelease.Namespace
	client.ReleaseName = newRelease.Name
	client.GenerateName = false
	client.NameTemplate = ""
	client.OutputDir = ""
	client.Atomic = newRelease.Atomic
	client.SkipCRDs = newRelease.SkipCrds
	client.SubNotes = newRelease.RenderSubchartNotes
	client.DisableOpenAPIValidation = newRelease.DisableOpenapiValidation
	client.Replace = newRelease.Replace
	client.Description = newRelease.Description
	client.CreateNamespace = newRelease.CreateNamespace

	if cmd := newRelease.Postrender; cmd != "" {
		pr, err := postrender.NewExec(cmd)

		if err != nil {
			return err
		}

		client.PostRenderer = pr
	}

	logger.V(9).Infof("install helm chart")
	rel, err := client.RunWithContext(r.canceler.context, c, values)
	if err != nil && rel == nil {
		return err
	}

	if err != nil && rel != nil {
		_, exists, existsErr := resourceReleaseLookup(newRelease.Name, conf)
		if existsErr != nil {
			return err
		}
		if !exists {
			return err
		}

		// Don't expect this to fail
		if err := setReleaseAttributes(newRelease, rel, false); err != nil {
			return err
		}
		_ = r.host.Log(ctx, diag.Warning, urn, fmt.Sprintf("Helm release %q was created but has a failed status. Use the `helm` command to investigate the error, correct it, then retry. Reason: %v", client.ReleaseName, err))
		return &releaseFailedError{release: newRelease, err: err}
	}

	err = setReleaseAttributes(newRelease, rel, false)
	return err
}

type releaseFailedError struct {
	release *Release
	err     error
}

func (e *releaseFailedError) Error() string {
	var s strings.Builder
	s.WriteString("Helm Release ")
	if e.release != nil {
		s.WriteString(fmt.Sprintf("%s/%s: ", e.release.Namespace, e.release.Name))
	}
	s.WriteString(e.err.Error())
	return "failed to become available within allocated timeout. Error: " + s.String()
}

func (r *helmReleaseProvider) helmUpdate(newRelease, oldRelease *Release) error {
	logger.V(9).Infof("getChart: %q settings: %#v", newRelease.Chart, r.settings)

	actionConfig, err := r.getActionConfig(oldRelease.Namespace)
	if err != nil {
		return err
	}
	client := action.NewUpgrade(actionConfig)
	cpo := &client.ChartPathOptions
	// Get Chart metadata, if we fail - we're done
	chart, path, err := getChart(cpo, actionConfig.RegistryClient, r.settings, newRelease)
	if err != nil {
		return err
	}

	// check and update the chart's dependencies if needed
	updated, err := checkChartDependencies(
		chart,
		path,
		newRelease.Keyring,
		r.settings,
		actionConfig.RegistryClient,
		newRelease.DependencyUpdate)
	if err != nil {
		return err
	} else if updated {
		// load the chart again if its dependencies have been updated
		chart, err = loader.Load(path)
		if err != nil {
			return err
		}
	}

	values, err := getValues(newRelease)
	if err != nil {
		return fmt.Errorf("error getting values for a diff: %w", err)
	}

	if newRelease.Lint {
		if err := lintChart(path, values); err != nil {
			return err
		}
	}

	client.Devel = newRelease.Devel
	client.Namespace = newRelease.Namespace
	client.Timeout = getTimeoutOrDefault(newRelease.Timeout)
	client.Wait = !newRelease.SkipAwait
	client.DisableHooks = newRelease.DisableCRDHooks
	client.Atomic = newRelease.Atomic
	client.SubNotes = newRelease.RenderSubchartNotes
	client.WaitForJobs = !newRelease.SkipAwait && newRelease.WaitForJobs
	client.Force = newRelease.ForceUpdate
	client.ResetValues = newRelease.ResetValues
	client.ReuseValues = newRelease.ReuseValues
	client.Recreate = newRelease.RecreatePods
	client.MaxHistory = 0
	if newRelease.MaxHistory != nil {
		client.MaxHistory = *newRelease.MaxHistory
	}
	client.CleanupOnFail = newRelease.CleanupOnFail
	client.Description = newRelease.Description

	if cmd := newRelease.Postrender; cmd != "" {
		pr, err := postrender.NewExec(cmd)

		if err != nil {
			return err
		}
		client.PostRenderer = pr
	}

	rel, err := client.RunWithContext(r.canceler.context, newRelease.Name, chart, values)
	if err != nil && rel == nil {
		return err
	}
	if err != nil && errors.Is(err, driver.ErrNoDeployedReleases) {
		logger.V(9).Infof("No existing release found.")
		return err
	}
	if err != nil {
		if err := setReleaseAttributes(newRelease, rel, false); err != nil {
			return err
		}
		return fmt.Errorf("error running update: %w", &releaseFailedError{release: newRelease, err: err})
	}

	err = setReleaseAttributes(newRelease, rel, false)
	return err
}

func adoptOldNameIfUnnamed(new, old resource.PropertyMap) {
	if _, ok := new["name"]; ok {
		return
	}
	contract.Assertf(old["name"].StringValue() != "", "expected 'name' value to be nonempty: %v", old)
	new["name"] = old["name"]
}

func assignNameIfAutonameable(pm resource.PropertyMap, urn resource.URN) {
	name, ok := pm["name"]
	if !ok || (name.IsString() && name.StringValue() == "") {
		prefix := urn.Name() + "-"
		autoname, err := resource.NewUniqueHex(prefix, 0, 0)
		contract.AssertNoErrorf(err, "unexpected error while executing NewUniqueHex")
		pm["name"] = resource.NewStringProperty(autoname)
	}
}

func (r *helmReleaseProvider) Diff(ctx context.Context, req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("Provider[%s].Diff(%s)", r.name, urn)

	// Get old state. This is an object of the form {inputs: {...}, live: {...}} where `inputs` is the
	// previous resource inputs supplied by the user, and `live` is the computed state of that inputs
	// we received back from the API server.
	olds, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true, KeepSecrets: true,
	})
	if err != nil {
		return nil, err
	}

	// Get new resource inputs. The user is submitting these as an update.
	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.news", label),
		KeepUnknowns: true,
		SkipNulls:    true,
		KeepSecrets:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("diff failed because malformed resource inputs: %w", err)
	}

	// Extract old inputs from the `__inputs` field of the old state.
	oldInputs, _ := parseCheckpointRelease(olds)

	// apply ignoreChanges
	for _, ignore := range req.GetIgnoreChanges() {
		if ignore == "version" {
			news["version"] = olds["version"]
		}
	}

	// remove deprecated inputs from old inputs, to avoid producing a delete op w.r.t. the new state.
	delete(oldInputs, "checksum")
	delete(oldInputs, "resourceNames")

	oldRelease, err := decodeRelease(olds, fmt.Sprintf("%s.olds", label))
	if err != nil {
		return nil, err
	}
	newRelease, err := decodeRelease(news, fmt.Sprintf("%s.news", label))
	if err != nil {
		return nil, err
	}

	logger.V(9).Infof("Diff: Old release: %#v", oldRelease)
	logger.V(9).Infof("Diff: New release: %#v", newRelease)

	// Generate a patch to apply the new inputs to the old state, including deletions.
	// Computed values are mapped to null, and secrets are mapped to plain values.
	// Later, we'll use this patch to generate a diff response, with special handling for the computed values.
	oldInputsJSON, err := json.Marshal(oldInputs.MapRepl(nil, mapReplExtractValues))
	if err != nil {
		return nil, fmt.Errorf("internal error: json.Marshal(oldInputsJson): %w", err)
	}
	logger.V(9).Infof("oldInputsJSON: %s", string(oldInputsJSON))
	newInputsJSON, err := json.Marshal(news.MapRepl(nil, mapReplExtractValues))
	if err != nil {
		return nil, fmt.Errorf("internal error: json.Marshal(oldInputsJson): %w", err)
	}
	logger.V(9).Infof("newInputsJSON: %s", string(newInputsJSON))
	oldStateJSON, err := json.Marshal(olds.MapRepl(nil, mapReplExtractValues))
	if err != nil {
		return nil, fmt.Errorf("internal error: json.Marshal(oldStateJson): %w", err)
	}
	logger.V(9).Infof("oldStateJSON: %s", string(oldStateJSON))
	strategicPatchJSON, err := strategicpatch.CreateThreeWayMergePatch(oldInputsJSON, newInputsJSON, oldStateJSON, &noSchema{}, true)
	if err != nil {
		return nil, fmt.Errorf("internal error: CreateThreeWayMergePatch: %w", err)
	}
	logger.V(9).Infof("strategicPatchJSON: %s", string(strategicPatchJSON))
	patchObj := map[string]any{}
	if err = json.Unmarshal(strategicPatchJSON, &patchObj); err != nil {
		return nil, fmt.Errorf(
			"Failed to check for changes in Helm release %s/%s because of an error serializing "+
				"the JSON patch describing resource changes: %w",
			oldRelease.Namespace, oldRelease.Name, err)

	}

	// Pack up PB, ship response back.
	hasChanges := pulumirpc.DiffResponse_DIFF_NONE

	var changes, replaces []string
	var detailedDiff map[string]*pulumirpc.PropertyDiff
	if len(patchObj) != 0 {
		hasChanges = pulumirpc.DiffResponse_DIFF_SOME

		for k := range patchObj {
			changes = append(changes, k)
		}

		logger.V(9).Infof("patchObj: %+v", patchObj)
		logger.V(9).Infof("oldLiveState: %+v", olds.Mappable())
		logger.V(9).Infof("news: %+v", news.Mappable())
		logger.V(9).Infof("oldInputs: %+v", oldInputs.Mappable())

		strip := func(pm resource.PropertyMap) map[string]interface{} {
			// strip the secretness but retain computedness (as is understood by convertPatchToDiff)
			return pm.MapRepl(nil, mapReplStripSecrets)
		}
		forceNewFields := []string{".name", ".namespace"}
		if detailedDiff, err = convertPatchToDiff(patchObj, strip(olds), strip(news), strip(oldInputs), forceNewFields...); err != nil {
			return nil, fmt.Errorf(
				"Failed to check for changes in helm release %s/%s because of an error "+
					"converting JSON patch describing resource changes to a diff: %w",
				oldRelease.Namespace, oldRelease.Name, err)

		}

		for k, v := range detailedDiff {
			switch v.Kind {
			case pulumirpc.PropertyDiff_ADD_REPLACE, pulumirpc.PropertyDiff_DELETE_REPLACE, pulumirpc.PropertyDiff_UPDATE_REPLACE:
				replaces = append(replaces, k)
			}
		}
	}

	return &pulumirpc.DiffResponse{
		Changes:             hasChanges,
		Replaces:            replaces,
		Stables:             []string{},
		DeleteBeforeReplace: false, // TODO: revisit this.
		Diffs:               changes,
		DetailedDiff:        detailedDiff,
		HasDetailedDiff:     len(detailedDiff) > 0,
	}, nil
}

func lintChart(path string, values map[string]any) (err error) {
	l := action.NewLint()
	result := l.Run([]string{path}, values)

	return resultToError(result)
}

func resultToError(r *action.LintResult) error {
	if len(r.Errors) == 0 {
		return nil
	}

	messages := []string{}
	for _, msg := range r.Messages {
		for _, err := range r.Errors {
			if err == msg.Err {
				messages = append(messages, fmt.Sprintf("%s: %s", msg.Path, msg.Err))
				break
			}
		}
	}

	return fmt.Errorf("malformed chart or values: \n\t%s", strings.Join(messages, "\n\t"))
}

func (r *helmReleaseProvider) Create(ctx context.Context, req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("Provider[%s].Create(%s)", r.name, urn)

	news, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.properties", label),
		KeepUnknowns: true,
		SkipNulls:    true,
		KeepSecrets:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("create failed because malformed resource inputs: %w", err)
	}

	newRelease, err := decodeRelease(news, fmt.Sprintf("%s.news", label))
	if err != nil {
		return nil, err
	}

	id := ""

	var creationError error
	if !req.GetPreview() {
		if r.clusterUnreachable {
			return nil, fmt.Errorf("can't create Helm Release with unreachable cluster: %s", r.clusterUnreachableReason)
		}
		id = fqName(newRelease.Namespace, newRelease.Name)
		if err := r.helmCreate(ctx, urn, newRelease); err != nil {
			var failedErr *releaseFailedError
			if errors.As(err, &failedErr) {
				creationError = failedErr
			} else {
				return nil, err
			}
		}
	}

	obj := checkpointRelease(news, newRelease, fmt.Sprintf("%s.news", label), req.GetPreview())
	inputsAndComputed, err := plugin.MarshalProperties(
		obj, plugin.MarshalOptions{
			Label:        fmt.Sprintf("%s.inputsAndComputed", label),
			KeepUnknowns: true,
			SkipNulls:    true,
			KeepSecrets:  r.enableSecrets,
		})
	if err != nil {
		return nil, err
	}

	if creationError != nil {
		return nil, partialError(
			id,
			fmt.Errorf(
				"Helm release %q was created, but failed to initialize completely. "+
					"Use Helm CLI to investigate: %w", id, creationError),

			inputsAndComputed,
			nil)
	}

	logger.V(9).Infof("Create: [id: %q] properties: %+v", id, inputsAndComputed)
	return &pulumirpc.CreateResponse{Id: id, Properties: inputsAndComputed}, nil
}

func (r *helmReleaseProvider) Read(ctx context.Context, req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("Provider[%s].Read(%s)", r.name, urn)
	logger.V(9).Infof("%s Starting", label)

	oldState, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true, KeepSecrets: true,
	})
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	existingRelease, err := decodeRelease(oldState, fmt.Sprintf("%s.olds", label))
	if err != nil {
		return nil, err
	}
	logger.V(9).Infof("%s decoded release: %#v", label, existingRelease)

	var namespace, name string
	if len(oldState) == 0 {
		namespace, name = parseFqName(req.GetId())
		logger.V(9).Infof("%s Starting import for %s/%s", label, namespace, name)
	} else {
		name = existingRelease.Name
		namespace = existingRelease.Namespace
	}

	actionConfig, err := r.getActionConfig(namespace)
	if err != nil {
		return nil, err
	}
	liveObj, exists, err := resourceReleaseLookup(name, actionConfig)
	if !exists && err == nil {
		// If not found, this resource was probably deleted.
		return deleteResponse, nil
	}
	if err != nil {
		return nil, err
	}

	err = setReleaseAttributes(existingRelease, liveObj, false)
	if err != nil {
		return nil, err
	}

	logger.V(9).Infof("%s Found release %s/%s", label, namespace, name)

	oldInputs, _ := parseCheckpointRelease(oldState)
	if oldInputs == nil {
		// No old inputs suggests this is an import. Hydrate the state from the current live object.
		// A subsequent Check operation will apply the computed inputs.
		err = r.importRelease(ctx, urn, existingRelease, liveObj)
		if err != nil {
			return nil, err
		}
		logger.V(9).Infof("%s Imported release: %#v", label, existingRelease)

		oldInputs = r.serializeImportInputs(existingRelease)
		r.setDefaults(oldInputs)
	}

	// Return a new "checkpoint object".
	state, err := plugin.MarshalProperties(
		checkpointRelease(oldInputs, existingRelease, fmt.Sprintf("%s.olds", label), false), plugin.MarshalOptions{
			Label:        fmt.Sprintf("%s.state", label),
			KeepUnknowns: true,
			SkipNulls:    true,
			KeepSecrets:  r.enableSecrets,
		})
	if err != nil {
		return nil, err
	}

	inputs, err := plugin.MarshalProperties(oldInputs, plugin.MarshalOptions{
		Label: label + ".inputs", KeepUnknowns: true, SkipNulls: true, KeepSecrets: r.enableSecrets, //nolint:goconst
	})
	if err != nil {
		return nil, err
	}

	id := fqName(existingRelease.Namespace, existingRelease.Name)
	if reqID := req.GetId(); len(reqID) > 0 {
		id = reqID
	}

	return &pulumirpc.ReadResponse{Id: id, Properties: state, Inputs: inputs}, nil
}

func (r *helmReleaseProvider) serializeImportInputs(release *Release) resource.PropertyMap {
	inputs := resource.NewPropertyMap(release)
	delete(inputs, "resourceNames")
	delete(inputs, "status")
	return inputs
}

func (r *helmReleaseProvider) Update(ctx context.Context, req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("Provider[%s].Update(%s)", r.name, urn)
	oldState, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.olds", label),
		KeepUnknowns: true,
		SkipNulls:    true,
		KeepSecrets:  true,
	})
	if err != nil {
		return nil, err
	}
	newResInputs, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.news", label),
		KeepUnknowns: true,
		SkipNulls:    true,
		KeepSecrets:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("update failed because malformed resource inputs: %w", err)
	}

	logger.V(9).Infof("%s executing", label)

	newRelease, err := decodeRelease(newResInputs, fmt.Sprintf("%s.news", label))
	if err != nil {
		return nil, err
	}

	oldRelease, err := decodeRelease(oldState, fmt.Sprintf("%s.olds", label))
	if err != nil {
		return nil, err
	}

	var updateError error
	if !req.GetPreview() {
		if r.clusterUnreachable {
			return nil, fmt.Errorf("can't update Helm Release with unreachable cluster: %s", r.clusterUnreachableReason)
		}
		if err = r.helmUpdate(newRelease, oldRelease); err != nil {
			var failedErr *releaseFailedError
			if errors.As(err, &failedErr) {
				updateError = failedErr
			} else {
				return nil, err
			}
		}
	}

	checkpointed := checkpointRelease(newResInputs, newRelease, fmt.Sprintf("%s.news", label), req.GetPreview())
	inputsAndComputed, err := plugin.MarshalProperties(
		checkpointed, plugin.MarshalOptions{
			Label:        fmt.Sprintf("%s.inputsAndComputed", label),
			KeepUnknowns: true,
			SkipNulls:    true,
			KeepSecrets:  r.enableSecrets,
		})
	if err != nil {
		return nil, err
	}

	if updateError != nil {
		return nil, partialError(
			fqName(newRelease.Namespace, newRelease.Name),
			fmt.Errorf(
				"Helm release %q failed to initialize completely. "+
					"Use Helm CLI to investigate: %w", fqName(newRelease.Namespace, newRelease.Name), updateError),

			inputsAndComputed,
			nil)
	}
	return &pulumirpc.UpdateResponse{Properties: inputsAndComputed}, nil
}

func (r *helmReleaseProvider) Delete(ctx context.Context, req *pulumirpc.DeleteRequest) (*pbempty.Empty, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("Provider[%s].Delete(%s)", r.name, urn)

	// Obtain new properties, create a Kubernetes `unstructured.Unstructured`.
	olds, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true, KeepSecrets: true,
	})
	if err != nil {
		return nil, err
	}

	release, err := decodeRelease(olds, fmt.Sprintf("%s.olds", label))
	if err != nil {
		return nil, err
	}

	namespace := release.Namespace
	actionConfig, err := r.getActionConfig(namespace)
	if err != nil {
		return nil, err
	}

	name := release.Name

	uninstall := action.NewUninstall(actionConfig)
	if release.Atomic || !release.SkipAwait { // If the release was atomic or skipAwait was not set, block on deletion
		uninstall.Wait = true
		uninstall.Timeout = getTimeoutOrDefault(release.Timeout)
	}
	// TODO: once https://github.com/helm/helm/pull/12109 is merged, use uninstall.RunWithContext
	res, err := uninstall.Run(name)
	if err != nil {
		return nil, err
	}

	if res.Info != "" {
		_ = r.host.Log(context.Background(), diag.Warning, urn, fmt.Sprintf("Helm uninstall returned information: %q", res.Info))
	}
	return &pbempty.Empty{}, nil
}

func checkpointRelease(inputs resource.PropertyMap, outputs *Release, label string, isPreview bool) resource.PropertyMap {
	logger.V(9).Infof("[%s] Checkpointing outputs: %#v", label, outputs)
	logger.V(9).Infof("[%s] Checkpointing inputs: %#v", label, inputs)
	object := resource.NewPropertyMap(outputs)
	object["__inputs"] = resource.MakeSecret(resource.NewObjectProperty(inputs))

	// Make sure parts of the inputs which are marked as secrets in the inputs are retained as
	// secrets in the outputs. Likewise for computed values.
	annotateComputed(object, inputs)
	annotateSecrets(object, inputs)

	// If this is a preview, emit computed placeholders for the pure outputs.
	if isPreview {
		object["resourceNames"] = resource.MakeComputed(resource.NewStringProperty(""))
		object["status"] = resource.MakeComputed(resource.NewStringProperty(""))
	}

	return object
}

// parseCheckpointRelease returns inputs that are saved in the `__inputs` field of the state.
func parseCheckpointRelease(obj resource.PropertyMap) (resource.PropertyMap, resource.PropertyMap) {
	state := obj.Copy()
	if inputs, ok := obj["__inputs"]; ok {
		delete(state, "__inputs")
		return inputs.SecretValue().Element.ObjectValue(), state
	}

	return nil, state
}

func setReleaseAttributes(release *Release, r *release.Release, isPreview bool) error {
	logger.V(9).Infof("Will populate dest: %#v with data from release: %+v", release, r)

	// import
	if release.Name == "" {
		release.Name = r.Name
	}
	if release.Namespace == "" {
		release.Namespace = r.Namespace
	}
	if release.Chart == "" {
		release.Chart = r.Chart.Metadata.Name
	}
	var err error
	logger.V(9).Infof("Setting release values: %+v", r.Config)
	release.Values, err = mergeMaps(release.Values, r.Config, release.AllowNullValues)
	if err != nil {
		return err
	}
	release.Version = r.Chart.Metadata.Version

	_, resources, err := convertYAMLManifestToJSON(r.Manifest)
	if err != nil {
		return err
	}

	release.ResourceNames = resources

	// TODO: redact sensitive values and add manifest to releaseSpec

	if isPreview {
		return nil
	}

	if release.Status == nil {
		release.Status = &ReleaseStatus{}
	}

	release.Status.Version = r.Chart.Metadata.Version
	release.Status.Namespace = r.Namespace
	release.Status.Name = r.Name
	release.Status.Status = r.Info.Status.String()

	release.Status.Name = r.Name
	release.Status.Namespace = r.Namespace
	release.Status.Revision = &r.Version
	release.Status.Chart = r.Chart.Metadata.Name
	release.Status.Version = r.Chart.Metadata.Version
	release.Status.AppVersion = r.Chart.Metadata.AppVersion
	return nil
}

func resourceReleaseLookup(name string, actionConfig *action.Configuration) (*release.Release, bool, error) {
	logger.V(9).Infof("[resourceReleaseLookup: %s]", name)
	release, err := getRelease(actionConfig, name)
	logger.V(9).Infof("[resourceReleaseLookup: %s] Done", name)

	if err == nil {
		return release, true, nil
	}

	if err == errReleaseNotFound {
		return nil, false, nil
	}

	return nil, false, err
}

func getRelease(cfg *action.Configuration, name string) (*release.Release, error) {
	get := action.NewGet(cfg)
	logger.V(9).Infof("%s getRelease post action created", name)

	res, err := get.Run(name)
	logger.V(9).Infof("%s getRelease post run", name)

	if err != nil {
		logger.V(9).Infof("getRelease for %s errored", name)
		logger.V(9).Infof("%v", err)
		if strings.Contains(err.Error(), "release: not found") {
			return nil, errReleaseNotFound
		}

		logger.V(9).Infof("could not get release %s", err)

		return nil, err
	}

	logger.V(9).Infof("%s getRelease done: %+v", name, res)

	return res, nil
}

func (r *helmReleaseProvider) importRelease(ctx context.Context, urn resource.URN, release *Release, hr *release.Release) error {

	// note: setReleaseAttributes pre-populates some of the inputs

	// Attempt to resolve the chart's origin in either a local or remote repository.
	// Note that the local chart is not verified.
	if name, _, found := searchProgramDirectory(hr.Chart.Metadata.Name, hr.Chart.Metadata.Version); found {
		release.Chart = name
	} else if repo, chart, found := searchHelmRepositories(r.settings, hr.Chart.Metadata.Name, hr.Chart.Metadata.Version); found {
		// use a local repository reference, rather than reconstructing all the repository opts
		release.Chart = fmt.Sprintf("%s/%s", repo.Name, chart.Name)
	} else {
		// fallback to using a local chart reference
		release.Chart = hr.Chart.Metadata.Name
	}

	chart, err := r.helmLoad(ctx, urn, release)
	if err != nil {
		// Likely because the chart is not readily available (e.g. import of chart where no repo info is stored).
		// Eat the error to allow import to succeed, assuming that Check will report the failure later.
		contract.IgnoreError(err)
	} else {
		release.Version = chart.Metadata.Version
	}

	return nil
}

func isChartInstallable(ch *helmchart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return fmt.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func getValues(release *Release) (map[string]any, error) {
	var err error
	base := map[string]any{}
	base, err = mergeMaps(base, release.Values, release.AllowNullValues)
	if err != nil {
		return nil, err
	}
	return base, logValues(base)
}

func logValues(values map[string]any) error {
	// copy array to avoid change values by the cloak function.
	asJSON, _ := json.Marshal(values)
	var c map[string]any
	err := json.Unmarshal(asJSON, &c)
	if err != nil {
		return err
	}

	y, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	logger.V(9).Infof(
		"---[ values.yaml ]-----------------------------------\n%s\n",
		string(y),
	)

	return nil
}

// Merges a and b map, preferring values from b map
func mergeMaps(a, b map[string]any, allowNullValues bool) (map[string]any, error) {
	if allowNullValues {
		// Use upstream's behavior.
		return helm.MergeMaps(a, b), nil
	}

	a = excludeNulls(a).(map[string]any)
	b = excludeNulls(b).(map[string]any)
	if err := mergo.Merge(&a, b, mergo.WithOverride, mergo.WithTypeCheck); err != nil {
		return nil, err
	}
	return a, nil
}

func excludeNulls(in any) any {
	switch reflect.TypeOf(in).Kind() {
	case reflect.Map:
		out := map[string]any{}
		m := in.(map[string]any)
		for k, v := range m {
			val := reflect.ValueOf(v)
			if val.IsValid() {
				switch val.Kind() {
				case reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
					if val.IsNil() {
						continue
					}
				}
				out[k] = excludeNulls(v)
			}
		}
		return out
	case reflect.Slice, reflect.Array:
		var out []any
		s := in.([]any)
		for _, i := range s {
			if i != nil {
				out = append(out, excludeNulls(i))
			}
		}
		return out
	}
	return in
}

// searchProgramDirectory implements a best-effort search for a chart in the program directory.
// It searches for a chart archive and for an unpacked chart directory, with and without a version suffix.
// The search order is: "<name>-<version>/", "<name>-<version>.tgz", "<name>/", "<name>.tgz".
func searchProgramDirectory(name, version string) (string, *helmchart.Metadata, bool) {
	var file, dir string
	dir = fmt.Sprintf("%s-%s", name, version)
	if c, err := loader.LoadDir(dir); err == nil {
		return dir, c.Metadata, true
	}
	file = fmt.Sprintf("%s-%s.tgz", name, version)
	if c, err := loader.LoadFile(file); err == nil {
		return file, c.Metadata, true
	}
	dir = name
	if c, err := loader.LoadDir(dir); err == nil {
		return dir, c.Metadata, true
	}
	file = fmt.Sprintf("%s.tgz", name)
	if c, err := loader.LoadFile(file); err == nil {
		return file, c.Metadata, true
	}
	return "", nil, false
}

// searchHelmRepositories implements a best-effort search for a chart in the locally-configured repositories.
func searchHelmRepositories(settings *cli.EnvSettings, name, version string) (*repo.Entry, *repo.ChartVersion, bool) {
	repoFile := settings.RepositoryConfig
	repoCacheDir := settings.RepositoryCache

	// Load the repositories.yaml
	rf, err := repo.LoadFile(repoFile)
	if errors.Is(err, fs.ErrNotExist) || len(rf.Repositories) == 0 {
		logger.V(9).Infof("no repositories configured")
		return nil, nil, false
	}

	// Scan the repositories for a chart
	for _, re := range rf.Repositories {
		n := re.Name
		f := filepath.Join(repoCacheDir, helmpath.CacheIndexFile(n))
		ind, err := repo.LoadIndexFile(f)
		if err != nil {
			logger.V(9).Infof("Repo %q is corrupt or missing. Try 'helm repo update'.", n)
			continue
		}
		chartVersion, err := ind.Get(name, version)
		if err != nil {
			logger.V(9).Infof("No such chart: %v", err)
			continue
		}
		return re, chartVersion, true
	}
	return nil, nil, false
}

func getChart(cpo *action.ChartPathOptions, registryClient *registry.Client, settings *cli.EnvSettings,
	newRelease *Release) (*helmchart.Chart, string,
	error) {
	logger.V(9).Infof("Looking up chart path options for release: %q", newRelease.Name)

	chartName, err := chartPathOptionsFromRelease(cpo, newRelease)
	if err != nil {
		return nil, "", err
	}

	logger.V(9).Infof("Chart name: %q", chartName)
	path, err := locateChart(cpo, registryClient, chartName, settings)
	if err != nil {
		return nil, "", err
	}

	logger.V(9).Infof("Trying to load chart from path: %q", path)
	c, err := loader.Load(path)
	if err != nil {
		return nil, "", err
	}

	return c, path, nil
}

// localChart determines if the specified chart is available locally (either compressed or not),
// and if so, validates it and returns the path to the chart.
func localChart(name string, verify bool, keyring string) (string, bool, error) {
	fi, err := os.Stat(name)
	if err != nil {
		// Helm eats all errors at this point.
		return "", false, nil
	}

	// If a folder is of the same name as a chart, use the folder if it contains a Chart.yaml.
	if fi.IsDir() {
		if _, err := os.Stat(filepath.Join(name, "Chart.yaml")); err != nil {
			// This is not a chart directory, so do not error as Helm could still
			// resolve this as a locally added chart repository, eg. `helm repo add`.
			return "", false, nil
		}
	}

	// Get the absolute path to the local compressed chart archive if it's a file.
	absPath, err := filepath.Abs(name)
	if err != nil {
		return "", false, err
	}

	// Verify the chart with the specified keyring if enabled. The chart must be in a compressed
	// archive, with a valid adjancent provenance file.
	if verify {
		if _, err := downloader.VerifyChart(name, keyring); err != nil {
			return "", false, err
		}
	}

	return absPath, true, nil
}

// locateChart is a copy of cpo.LocateChart with specialized behavior around the resolution of
// local charts. Helm prefers a local chart over a remote chart with the same name, even if the local chart
// directory is not well-formed. Pulumi adds additional checks of the local chart directory.
func locateChart(cpo *action.ChartPathOptions, registryClient *registry.Client, name string,
	settings *cli.EnvSettings) (string, error) {
	name = strings.TrimSpace(name)
	version := strings.TrimSpace(cpo.Version)

	// Determine if chart is already available locally.
	if cpo.RepoURL == "" {
		abs, found, err := localChart(name, cpo.Verify, cpo.Keyring)
		if found || err != nil {
			return abs, err
		}

		// If not found, do more validations. This is from the original LocateChart.
		if filepath.IsAbs(name) || strings.HasPrefix(name, ".") {
			return name, fmt.Errorf("path %q not found", name)
		}
	}

	// If local chart not found, try to download it.
	dl := downloader.ChartDownloader{
		Out:     os.Stdout,
		Keyring: cpo.Keyring,
		Getters: getter.All(settings),
		Options: []getter.Option{
			getter.WithPassCredentialsAll(cpo.PassCredentialsAll),
			getter.WithTLSClientConfig(cpo.CertFile, cpo.KeyFile, cpo.CaFile),
			getter.WithInsecureSkipVerifyTLS(cpo.InsecureSkipTLSverify),
		},
		RepositoryConfig: settings.RepositoryConfig,
		RepositoryCache:  settings.RepositoryCache,
		RegistryClient:   registryClient,
	}

	if registry.IsOCI(name) {
		dl.Options = append(dl.Options, getter.WithRegistryClient(registryClient))
	}

	if cpo.Verify {
		dl.Verify = downloader.VerifyAlways
	}
	if cpo.RepoURL != "" {
		chartURL, err := repo.FindChartInAuthAndTLSAndPassRepoURL(
			cpo.RepoURL, cpo.Username, cpo.Password, name, version, cpo.CertFile,
			cpo.KeyFile, cpo.CaFile, cpo.InsecureSkipTLSverify, cpo.PassCredentialsAll,
			getter.All(settings))
		if err != nil {
			return "", err
		}
		name = chartURL

		// Only pass the user/pass on when the user has said to or when the
		// location of the chart repo and the chart are the same domain.
		u1, err := url.Parse(cpo.RepoURL)
		if err != nil {
			return "", err
		}
		u2, err := url.Parse(chartURL)
		if err != nil {
			return "", err
		}

		// Host on URL (returned from url.Parse) contains the port if present.
		// This check ensures credentials are not passed between different
		// services on different ports.
		if cpo.PassCredentialsAll || (u1.Scheme == u2.Scheme && u1.Host == u2.Host) {
			dl.Options = append(dl.Options, getter.WithBasicAuth(cpo.Username, cpo.Password))
		} else {
			dl.Options = append(dl.Options, getter.WithBasicAuth("", ""))
		}
	} else {
		dl.Options = append(dl.Options, getter.WithBasicAuth(cpo.Username, cpo.Password))
	}

	if err := os.MkdirAll(settings.RepositoryCache, 0755); err != nil {
		return "", err
	}

	filename, _, err := dl.DownloadTo(name, version, settings.RepositoryCache)
	if err == nil {
		lname, err := filepath.Abs(filename)
		if err != nil {
			return filename, err
		}
		return lname, nil
	} else if settings.Debug {
		return filename, err
	}

	atVersion := ""
	if version != "" {
		atVersion = fmt.Sprintf(" at version %q", version)
	}

	return filename, fmt.Errorf("failed to download %q%s", name, atVersion)
}

func checkChartDependencies(c *helmchart.Chart, path, keyring string, settings *cli.EnvSettings,
	registryClient *registry.Client, dependencyUpdate bool) (bool, error) {
	p := getter.All(settings)

	if req := c.Metadata.Dependencies; req != nil {
		err := action.CheckDependencies(c, req)
		if err != nil {
			if dependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        path,
					Keyring:          keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
					RegistryClient:   registryClient,
					Debug:            settings.Debug,
				}
				return true, man.Update()
			}
			return false, err
		}
		return false, err
	}
	return false, nil
}

func chartPathOptionsFromRelease(cpo *action.ChartPathOptions, release *Release) (string, error) {
	chartName := release.Chart

	version := getVersion(release)
	cpo.Keyring = release.Keyring
	cpo.Verify = release.Verify
	cpo.Version = version
	if release.RepositoryOpts != nil {
		var repositoryURL string
		var err error
		repository := release.RepositoryOpts.Repo
		repositoryURL, chartName, err = resolveChartName(repository, strings.TrimSpace(chartName))
		if err != nil {
			return "", err
		}
		cpo.CertFile = release.RepositoryOpts.CertFile
		cpo.CaFile = release.RepositoryOpts.CAFile
		cpo.KeyFile = release.RepositoryOpts.KeyFile
		cpo.Username = release.RepositoryOpts.Username
		cpo.Password = release.RepositoryOpts.Password
		cpo.RepoURL = repositoryURL
	}

	return chartName, nil
}

func getVersion(release *Release) (version string) {
	version = release.Version

	if version == "" && release.Devel {
		logger.V(9).Infof("setting version to >0.0.0-0")
		version = ">0.0.0-0"
	} else {
		version = strings.TrimSpace(version)
	}

	return
}

func resolveChartName(repository, name string) (string, string, error) {
	if registry.IsOCI(name) {
		return "", name, nil
	}
	_, err := url.ParseRequestURI(repository)
	if err == nil {
		return repository, name, nil
	}

	if !strings.Contains(name, "/") && repository != "" {
		name = fmt.Sprintf("%s/%s", repository, name)
	}

	return "", name, nil
}

func isHelmRelease(urn resource.URN) bool {
	return urn.Type() == "kubernetes:helm.sh/v3:Release"
}

func getTimeoutOrDefault(timeout int) time.Duration {
	if timeout == 0 {
		timeout = defaultTimeoutSeconds
	}
	return time.Duration(timeout) * time.Second
}

// noSchema implements a trivial lookup function for patch metadata (i.e. patch strategy and merge key).
// CreateThreeWayMergePatch supports various strategies for merging maps and slices, but we use the default strategy.
type noSchema struct{}

var _ strategicpatch.LookupPatchMeta = &noSchema{}

func (*noSchema) LookupPatchMetadataForSlice(key string) (strategicpatch.LookupPatchMeta, strategicpatch.PatchMeta, error) {
	return &noSchema{}, strategicpatch.PatchMeta{}, nil
}

func (*noSchema) LookupPatchMetadataForStruct(key string) (strategicpatch.LookupPatchMeta, strategicpatch.PatchMeta, error) {
	return &noSchema{}, strategicpatch.PatchMeta{}, nil
}

func (*noSchema) Name() string {
	return ""
}
