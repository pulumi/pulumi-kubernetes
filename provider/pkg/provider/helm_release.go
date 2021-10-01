// Copyright 2021, Pulumi Corporation.  All rights reserved.

package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"k8s.io/client-go/tools/clientcmd/api"

	jsonpatch "github.com/evanphx/json-patch"
	pbempty "github.com/golang/protobuf/ptypes/empty"
	pkgerrors "github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mitchellh/mapstructure"
	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/metadata"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/action"
	helmchart "helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/postrender"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// errReleaseNotFound is the error when a Helm release is not found
var errReleaseNotFound = errors.New("release not found")

// Release should explicitly track the shape of helm.sh/v3:Release resource
type Release struct {
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
	// Prevent CRD hooks from, running, but run other hooks.  See helm install --no-crd-hook
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
	RepositoryOpts RepositoryOpts `json:"repositoryOpts,omitempty"`
	// When upgrading, reset the values to the ones built into the chart
	ResetValues bool `json:"resetValues,omitempty"`
	// When upgrading, reuse the last release's values and merge in any overrides. If 'reset_values' is specified, this is ignored
	ReuseValues bool `json:"reuseValues,omitempty"`
	// Custom values to be merged with items loaded from values.
	Values map[string]interface{} `json:"values,omitempty"`
	// If set, no CRDs will be installed. By default, CRDs are installed if not already present
	SkipCrds bool `json:"skipCrds,omitempty"`
	// Time in seconds to wait for any individual kubernetes operation.
	Timeout int `json:"timeout,omitempty"`
	// ValueYamlFiles List of assets (raw yaml files) to pass to helm.
	//ValueYamlFiles []*resource.Asset `json:"valueYamlFiles,omitempty"`
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

type ReleaseSpec struct {
}

// Specification defining the Helm chart repository to use.
type RepositoryOpts struct {
	// Repository where to locate the requested chart. If is a URL the chart is installed without installing the repository.
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
	host             *provider.HostClient
	helmDriver       string
	apiConfig        *api.Config
	defaultOverrides *clientcmd.ConfigOverrides
	restConfig       *rest.Config
	defaultNamespace string
	enableSecrets    bool
	name             string
	settings         *cli.EnvSettings
}

func newHelmReleaseProvider(
	host *provider.HostClient,
	apiConfig *api.Config,
	defaultOverrides *clientcmd.ConfigOverrides,
	restConfig *rest.Config,
	helmDriver,
	namespace string,
	enableSecrets bool,
	pluginsDirectory,
	registryConfigPath,
	repositoryConfigPath,
	repositoryCache string,
) (customResourceProvider, error) {
	settings := cli.New()
	settings.PluginsDirectory = pluginsDirectory
	settings.RegistryConfig = registryConfigPath
	settings.RepositoryConfig = repositoryConfigPath
	settings.RepositoryCache = repositoryCache

	return &helmReleaseProvider{
		host:             host,
		apiConfig:        apiConfig,
		defaultOverrides: defaultOverrides,
		restConfig:       restConfig,
		helmDriver:       helmDriver,
		defaultNamespace: namespace,
		enableSecrets:    enableSecrets,
		name:             "kubernetes:helmrelease",
		settings:         settings,
	}, nil
}

func debug(format string, a ...interface{}) {
	logger.V(9).Infof("[DEBUG] %s", fmt.Sprintf(format, a...))
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
		// values would would reside in the $KUBECONFIG file, but can also be altered in several
		// places, including in env variables, client-go default values, and (if we allowed it) CLI
		// flags.
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
		clientConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &overrides)
	}
	kc := newKubeConfig(r.restConfig, clientConfig)
	if err := conf.Init(kc, namespace, r.helmDriver, debug); err != nil {
		return nil, err
	}
	return conf, nil
}

func decodeRelease(pm resource.PropertyMap) (*Release, error) {
	var release Release
	stripped := pm.MapRepl(nil, mapReplStripSecrets)
	logger.V(9).Infof("Decoding release: %#v", stripped)
	if err := mapstructure.Decode(stripped, &release); err != nil {
		return nil, fmt.Errorf("decoding failure: %w", err)
	}
	return &release, nil
}

func (r *helmReleaseProvider) Check(ctx context.Context, req *pulumirpc.CheckRequest, displayBetaWarning bool) (*pulumirpc.CheckResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("Provider[%s].Check(%s)", r.name, urn)

	if displayBetaWarning {
		_ = r.host.LogStatus(ctx, diag.Warning, urn, "Helm Release resource is currently in beta and may change. Use in production environments is discouraged.")
	}
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
		return nil, pkgerrors.Wrapf(err, "check failed because malformed resource inputs: %+v", err)
	}

	logger.V(9).Infof("Decoding new release.")
	new, err := decodeRelease(news)
	if err != nil {
		return nil, err
	}

	logger.V(9).Infof("Decoding old release.")
	old, err := decodeRelease(olds)
	if err != nil {
		return nil, err
	}

	if new.Namespace == "" {
		new.Namespace = r.defaultNamespace
	}

	if !new.SkipAwait && new.Timeout == 0 {
		new.Timeout = 300
	}

	if new.Keyring == "" {
		new.Keyring = os.ExpandEnv("$HOME/.gnupg/pubring.gpg")
	}

	var templateRelease bool
	if len(olds.Mappable()) > 0 {
		adoptOldNameIfUnnamed(new, old)

		templateRelease = true
	} else {
		assignNameIfAutonameable(new, news, "release")
		conf, err := r.getActionConfig(new.Namespace)
		if err != nil {
			return nil, err
		}
		exists, err := resourceReleaseExists(new.Name, conf)
		if err != nil {
			return nil, err
		}
		if !exists {
			templateRelease = true
		}
		// If resource exists, we are likely doing an import. We will just pass the inputs through.
	}

	if templateRelease {
		helmHome := os.Getenv("HELM_HOME")

		helmChartOpts := HelmChartOpts{
			HelmFetchOpts: HelmFetchOpts{
				CAFile:   new.RepositoryOpts.CAFile,
				CertFile: new.RepositoryOpts.CertFile,
				Devel:    new.Devel,
				Home:     helmHome,
				KeyFile:  new.RepositoryOpts.KeyFile,
				Keyring:  new.Keyring,
				Password: new.RepositoryOpts.Password,
				Repo:     new.RepositoryOpts.Repo,
				Username: new.RepositoryOpts.Password,
				Version:  new.Version,
			},
			APIVersions:              nil,
			Chart:                    new.Chart,
			IncludeTestHookResources: true,
			SkipCRDRendering:         new.SkipCrds,
			Namespace:                new.Namespace,
			Path:                     "",
			ReleaseName:              new.Name,
			Values:                   new.Values,
			Version:                  new.Version,
		}
		templ, err := helmTemplate(helmChartOpts)
		if err != nil {
			return nil, err
		}

		_, resources, err := convertYAMLManifestToJSON(templ)
		if err != nil {
			return nil, err
		}

		new.ResourceNames = resources
	}

	autonamed := resource.NewPropertyMap(new)
	annotateSecrets(autonamed, news)
	autonamedInputs, err := plugin.MarshalProperties(autonamed, plugin.MarshalOptions{
		Label:        fmt.Sprintf("%s.autonamedInputs", label),
		KeepUnknowns: true,
		SkipNulls:    true,
		KeepSecrets:  r.enableSecrets,
	})
	if err != nil {
		return nil, err
	}

	// Return new, possibly-autonamed inputs.
	return &pulumirpc.CheckResponse{Inputs: autonamedInputs}, nil
}

func (r *helmReleaseProvider) helmCreate(ctx context.Context, urn resource.URN, news resource.PropertyMap, newRelease *Release) error {
	conf, err := r.getActionConfig(newRelease.Namespace)
	if err != nil {
		return err
	}
	client := action.NewInstall(conf)
	logger.V(9).Infof("Looking up chart path options for release: %q", newRelease.Name)
	cpo, chartName, err := chartPathOptions(newRelease)
	if err != nil {
		return err
	}

	logger.V(9).Infof("getChart: %q settings: %#v, cpo: %+v", chartName, r.settings, cpo)
	c, path, err := getChart(chartName, r.settings, cpo)
	if err != nil {
		logger.V(9).Infof("getChart failed: %+v", err)
		return err
	}

	logger.V(9).Infof("Checking chart dependencies for chart: %q with path: %q", chartName, path)
	// check and update the chart's dependencies if needed
	updated, err := checkChartDependencies(
		c,
		path,
		newRelease.Keyring,
		r.settings,
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

	client.ChartPathOptions = *cpo
	client.ClientOnly = false
	client.DisableHooks = newRelease.DisableWebhooks
	client.Wait = !newRelease.SkipAwait
	client.WaitForJobs = !newRelease.SkipAwait && newRelease.WaitForJobs
	client.Devel = newRelease.Devel
	client.DependencyUpdate = newRelease.DependencyUpdate
	client.Timeout = time.Duration(newRelease.Timeout) * time.Second
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
	rel, err := client.Run(c, values)
	if err != nil && rel == nil {
		return err
	}

	if err != nil && rel != nil {
		exists, existsErr := resourceReleaseExists(newRelease.Name, conf)

		if existsErr != nil {
			return err
		}

		if !exists {
			return err
		}

		if err := setReleaseAttributes(newRelease, rel, false); err != nil {
			return err
		}

		_ = r.host.Log(ctx, diag.Warning, urn, fmt.Sprintf("Helm release %q was created but has a failed status. Use the `helm` command to investigate the error, correct it, then retry. Reason: %v", client.ReleaseName, err))
		return err

	}

	err = setReleaseAttributes(newRelease, rel, false)
	return err
}

func (r *helmReleaseProvider) helmUpdate(ctx context.Context, urn resource.URN, news resource.PropertyMap, newRelease, oldRelease *Release) error {
	cpo, chartName, err := chartPathOptions(newRelease)
	if err != nil {
		return err
	}

	logger.V(9).Infof("getChart: %q settings: %#v, cpo: %+v", chartName, r.settings, cpo)
	// Get Chart metadata, if we fail - we're done
	chart, path, err := getChart(chartName, r.settings, cpo)
	if err != nil {
		return err
	}

	// check and update the chart's dependencies if needed
	updated, err := checkChartDependencies(
		chart,
		path,
		newRelease.Keyring,
		r.settings,
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

	if newRelease.Lint {
		if err := resourceReleaseValidate(newRelease, news, r.settings, cpo); err != nil {
			return err
		}
	}

	actionConfig, err := r.getActionConfig(oldRelease.Namespace)
	if err != nil {
		return err
	}

	values, err := getValues(newRelease)
	if err != nil {
		return fmt.Errorf("error getting values for a diff: %w", err)
	}

	client := action.NewUpgrade(actionConfig)
	client.ChartPathOptions = *cpo
	client.Devel = newRelease.Devel
	client.Namespace = newRelease.Namespace
	client.Timeout = time.Duration(newRelease.Timeout) * time.Second
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

	rel, err := client.Run(newRelease.Name, chart, values)
	if err != nil && strings.Contains(err.Error(), "has no deployed releases") {
		logger.V(9).Infof("No existing release found.")
		return err
	} else if err != nil {
		return fmt.Errorf("error running update: %w", err)
	}

	err = setReleaseAttributes(newRelease, rel, false)
	return err
}

func adoptOldNameIfUnnamed(new, old *Release) {
	contract.Assert(old.Name != "")
	new.Name = old.Name
}

func assignNameIfAutonameable(release *Release, pm resource.PropertyMap, base tokens.QName) {
	if name, ok := pm["name"]; ok && name.IsComputed() {
		return
	}
	if name, ok := pm["name"]; !ok || name.StringValue() == "" {
		release.Name = fmt.Sprintf("%s-%s", base, metadata.RandString(8))
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
		return nil, pkgerrors.Wrapf(err, "diff failed because malformed resource inputs")
	}

	// Extract old inputs from the `__inputs` field of the old state.
	oldInputs, _ := parseCheckpointRelease(olds)
	diff := oldInputs.Diff(news)
	if diff == nil {
		logger.V(9).Infof("No diff found for %q", req.GetUrn())
		return &pulumirpc.DiffResponse{Changes: pulumirpc.DiffResponse_DIFF_NONE}, nil
	}

	oldRelease, err := decodeRelease(olds)
	if err != nil {
		return nil, err
	}
	newRelease, err := decodeRelease(news)
	if err != nil {
		return nil, err
	}

	logger.V(9).Infof("Diff: Old release: %#v", oldRelease)
	logger.V(9).Infof("Diff: New release: %#v", newRelease)

	// Always set desired state to DEPLOYED
	// TODO: This could be done in Check instead?
	if newRelease.Status == nil {
		newRelease.Status = &ReleaseStatus{}
	}
	newRelease.Status.Status = release.StatusDeployed.String()

	oldInputsJSON, err := json.Marshal(oldInputs.Mappable())
	if err != nil {
		return nil, err
	}
	newInputsJSON, err := json.Marshal(news.Mappable())
	if err != nil {
		return nil, err
	}

	logger.V(9).Infof("oldInputsJSON: %+v", string(oldInputsJSON))
	logger.V(9).Infof("newInputsJSON: %+v", string(newInputsJSON))
	patch, err := jsonpatch.CreateMergePatch(oldInputsJSON, newInputsJSON)
	if err != nil {
		return nil, err
	}
	patchObj := map[string]interface{}{}
	if err = json.Unmarshal(patch, &patchObj); err != nil {
		return nil, pkgerrors.Wrapf(
			err, "Failed to check for changes in Helm release %s because of an error serializing "+
				"the JSON patch describing resource changes",
			newRelease.Name)
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

		if detailedDiff, err = convertPatchToDiff(patchObj, olds.Mappable(), news.Mappable(), oldInputs.Mappable(), ".releaseSpec.name", ".releaseSpec.namespace"); err != nil {
			return nil, pkgerrors.Wrapf(
				err, "Failed to check for changes in helm release %s/%s because of an error "+
					"converting JSON patch describing resource changes to a diff",
				newRelease.Namespace, newRelease.Name)
		}
		for _, v := range detailedDiff {
			v.InputDiff = true
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
		HasDetailedDiff:     true,
	}, nil
}

func resourceReleaseValidate(release *Release, pm resource.PropertyMap, settings *cli.EnvSettings, cpo *action.ChartPathOptions) error {
	cpo, name, err := chartPathOptions(release)
	if err != nil {
		return fmt.Errorf("malformed values: \n\t%s", err)
	}

	values, err := getValues(release)
	if err != nil {
		return err
	}

	return lintChart(settings, name, cpo, values)
}

func lintChart(settings *cli.EnvSettings, name string, cpo *action.ChartPathOptions, values map[string]interface{}) (err error) {
	path, err := cpo.LocateChart(name, settings)
	if err != nil {
		return err
	}

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
		return nil, pkgerrors.Wrapf(err, "create failed because malformed resource inputs")
	}

	newRelease, err := decodeRelease(news)
	if err != nil {
		return nil, err
	}

	if !req.GetPreview() {
		if err = r.helmCreate(ctx, urn, news, newRelease); err != nil {
			return nil, err
		}
	}

	obj := checkpointRelease(news, newRelease)
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

	id := ""
	if !req.GetPreview() {
		id = fqName(newRelease.Namespace, newRelease.Name)
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
	oldInputs, err := plugin.UnmarshalProperties(req.GetInputs(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.oldInputs", label), KeepUnknowns: true, SkipNulls: true, KeepSecrets: true,
	})
	if err != nil {
		return nil, err
	}

	existingRelease, err := decodeRelease(oldState)
	if err != nil {
		return nil, err
	}
	logger.V(9).Infof("%s decoded release: %#v", label, existingRelease)

	var namespace, name string
	if len(oldState.Mappable()) == 0 {
		namespace, name = parseFqName(req.GetId())
	} else {
		name = existingRelease.Name
		namespace = existingRelease.Namespace
	}

	logger.V(9).Infof("%s Starting import for %s/%s", label, namespace, name)

	actionConfig, err := r.getActionConfig(namespace)
	if err != nil {
		return nil, err
	}
	exists, err := resourceReleaseExists(name, actionConfig)
	if err != nil {
		return nil, err
	}

	if !exists {
		// If not found, this resource was probably deleted.
		return deleteResponse, nil
	}

	liveObj, err := getRelease(actionConfig, name)
	if err != nil {
		return nil, err
	}

	err = setReleaseAttributes(existingRelease, liveObj, false)
	if err != nil {
		return nil, err
	}

	cpo, chartName, err := chartPathOptions(existingRelease)
	if err != nil {
		return nil, err
	}

	logger.V(9).Infof("Trying to get chart: %q, settings: %#v", chartName, r.settings)

	// Helm itself doesn't store any information about where the Chart was downloaded from.
	// We need the user to ensure the chart is downloadable by using `helm repo add` etc.
	_, _, err = getChart(chartName, r.settings, cpo)
	if err != nil {
		return nil, err
	}

	logger.V(9).Infof("%s Found release %s/%s", label, namespace, name)

	// Return a new "checkpoint object".
	state, err := plugin.MarshalProperties(
		checkpointRelease(oldInputs, existingRelease), plugin.MarshalOptions{
			Label:        fmt.Sprintf("%s.state", label),
			KeepUnknowns: true,
			SkipNulls:    true,
			KeepSecrets:  r.enableSecrets,
		})
	if err != nil {
		return nil, err
	}

	liveInputsPM := resource.NewPropertyMap(existingRelease)

	inputs, err := plugin.MarshalProperties(liveInputsPM, plugin.MarshalOptions{
		Label: label + ".inputs", KeepUnknowns: true, SkipNulls: true, KeepSecrets: r.enableSecrets,
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
		return nil, pkgerrors.Wrapf(err, "update failed because malformed resource inputs")
	}

	logger.V(9).Infof("%s executing", label)

	newRelease, err := decodeRelease(newResInputs)
	if err != nil {
		return nil, err
	}

	oldRelease, err := decodeRelease(oldState)
	if err != nil {
		return nil, err
	}

	if !req.GetPreview() {
		if err = r.helmUpdate(ctx, urn, newResInputs, newRelease, oldRelease); err != nil {
			return nil, err
		}
	}

	checkpointed := checkpointRelease(newResInputs, newRelease)
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
	return &pulumirpc.UpdateResponse{Properties: inputsAndComputed}, nil
}

func (r *helmReleaseProvider) Delete(ctx context.Context, req *pulumirpc.DeleteRequest) (*empty.Empty, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("Provider[%s].Delete(%s)", r.name, urn)

	// Obtain new properties, create a Kubernetes `unstructured.Unstructured`.
	olds, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true, KeepSecrets: true,
	})
	if err != nil {
		return nil, err
	}

	release, err := decodeRelease(olds)
	if err != nil {
		return nil, err
	}

	namespace := release.Namespace
	actionConfig, err := r.getActionConfig(namespace, true)
	if err != nil {
		return nil, err
	}

	name := release.Name

	res, err := action.NewUninstall(actionConfig).Run(name)
	if err != nil {
		return nil, err
	}

	if res.Info != "" {
		_ = r.host.Log(context.Background(), diag.Warning, "Helm uninstall returned information: %q", res.Info)
	}
	return &pbempty.Empty{}, nil
}

func checkpointRelease(inputs resource.PropertyMap, outputs *Release) resource.PropertyMap {
	object := resource.NewPropertyMap(outputs)
	object["__inputs"] = resource.MakeSecret(resource.NewObjectProperty(inputs))

	// Make sure parts of the inputs which are marked as secrets in the inputs are retained as
	// secrets in the outputs.
	annotateSecrets(object, inputs)
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
	if len(release.Values) == 0 {
		release.Values = r.Config
	}
	if release.Version == "" {
		release.Version = r.Chart.Metadata.Version
	}
	if release.Chart == "" {
		release.Chart = r.Chart.Metadata.Name
	}
	if release.Description == "" {
		release.Description = r.Info.Description
	}

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

func resourceReleaseExists(name string, actionConfig *action.Configuration) (bool, error) {
	logger.V(9).Infof("[resourceReleaseExists: %s]", name)
	_, err := getRelease(actionConfig, name)

	logger.V(9).Infof("[resourceReleaseExists: %s] Done", name)

	if err == nil {
		return true, nil
	}

	if err == errReleaseNotFound {
		return false, nil
	}

	return false, err
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

	logger.V(9).Infof("%s getRelease done", name)

	return res, nil
}

func isChartInstallable(ch *helmchart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return fmt.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func getValues(release *Release) (map[string]interface{}, error) {
	base := map[string]interface{}{}
	base = mergeMaps(base, release.Values)
	return base, logValues(base)
}

func logValues(values map[string]interface{}) error {
	// copy array to avoid change values by the cloak function.
	asJSON, _ := json.Marshal(values)
	var c map[string]interface{}
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

// Merges source and destination map, preferring values from the source map
// Taken from github.com/helm/pkg/cli/values/options.go
func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

func getChart(name string, settings *cli.EnvSettings, cpo *action.ChartPathOptions) (*helmchart.Chart, string, error) {
	path, err := cpo.LocateChart(name, settings)
	if err != nil {
		return nil, "", err
	}

	logger.V(9).Infof("Trying to load chart: %q from path: %q", name, path)
	c, err := loader.Load(path)
	if err != nil {
		return nil, "", err
	}

	return c, path, nil
}

func checkChartDependencies(c *helmchart.Chart, path, keyring string, settings *cli.EnvSettings, dependencyUpdate bool) (bool, error) {
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
					Debug:            settings.Debug,
				}
				log.Println("[DEBUG] Downloading chart dependencies...")
				return true, man.Update()
			}
			return false, err
		}
		return false, err
	}
	log.Println("[DEBUG] Chart dependencies are up to date.")
	return false, nil
}

func chartPathOptions(release *Release) (*action.ChartPathOptions, string, error) {
	chartName := release.Chart

	repository := release.RepositoryOpts.Repo
	repositoryURL, chartName, err := resolveChartName(repository, strings.TrimSpace(chartName))
	if err != nil {
		return nil, "", err
	}
	version := getVersion(release)

	return &action.ChartPathOptions{
		CaFile:   release.RepositoryOpts.CAFile,
		CertFile: release.RepositoryOpts.CertFile,
		KeyFile:  release.RepositoryOpts.KeyFile,
		Keyring:  release.Keyring,
		RepoURL:  repositoryURL,
		Verify:   release.Verify,
		Version:  version,
		Username: release.RepositoryOpts.Username,
		Password: release.RepositoryOpts.Password, // TODO: This should already be resolved.
	}, chartName, nil
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
	_, err := url.ParseRequestURI(repository)
	if err == nil {
		return repository, name, nil
	}

	if strings.Index(name, "/") == -1 && repository != "" {
		name = fmt.Sprintf("%s/%s", repository, name)
	}

	return "", name, nil
}

func isHelmRelease(urn resource.URN) bool {
	return urn.Type() == "kubernetes:helm.sh/v3:Release"
}
