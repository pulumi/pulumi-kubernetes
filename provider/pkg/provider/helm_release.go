package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
	pbempty "github.com/golang/protobuf/ptypes/empty"
	pkgerrors "github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mitchellh/mapstructure"
	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/metadata"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
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
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// errReleaseNotFound is the error when a Helm release is not found
var errReleaseNotFound = errors.New("release not found")

type Release struct {
	ResourceType string       `json:"resourceType,omitempty"`
	ReleaseSpec  *ReleaseSpec `json:"releaseSpec,omitempty"`
	// Status of the deployed release.
	Status *ReleaseStatus `json:"status,omitempty"`
}

type ReleaseSpec struct {
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
	RepositorySpec RepositorySpec `json:"repositorySpec,omitempty"`
	// When upgrading, reset the values to the ones built into the chart
	ResetValues bool `json:"resetValues,omitempty"`
	// When upgrading, reuse the last release's values and merge in any overrides. If 'reset_values' is specified, this is ignored
	ReuseValues bool `json:"reuseValues,omitempty"`
	// Custom values to be merged with the values.
	Set []*SetValue `json:"set,omitempty"`
	// If set, no CRDs will be installed. By default, CRDs are installed if not already present
	SkipCrds bool `json:"skipCrds,omitempty"`
	// Time in seconds to wait for any individual kubernetes operation.
	Timeout int `json:"timeout,omitempty"`
	// List of values in raw yaml format to pass to helm.
	Values []string `json:"values,omitempty"`
	// Verify the package before installing it.
	Verify bool `json:"verify,omitempty"`
	// Specify the exact chart version to install. If this is not specified, the latest version is installed.
	Version string `json:"version,omitempty"`
	// Will wait until all resources are in a ready state before marking the release as successful.
	Wait bool `json:"wait,omitempty"`
	// If wait is enabled, will wait until all Jobs have been completed before marking the release as successful.
	WaitForJobs bool `json:"waitForJobs,omitempty"`
}

// Specification defining the Helm chart repository to use.
type RepositorySpec struct {
	// Repository where to locate the requested chart. If is a URL the chart is installed without installing the repository.
	Repository string `json:"repository,omitempty"`
	// The Repositories CA File
	RepositoryCAFile string `json:"repositoryCAFile,omitempty"`
	// The repositories cert file
	RepositoryCertFile string `json:"repositoryCertFile,omitempty"`
	// The repositories cert key file
	RepositoryKeyFile string `json:"repositoryKeyFile,omitempty"`
	// Password for HTTP basic authentication
	RepositoryPassword string `json:"repositoryPassword,omitempty"`
	// Username for HTTP basic authentication
	RepositoryUsername string `json:"repositoryUsername,omitempty"`
}

type SetValue struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
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
	// Set of extra values, added to the chart. The sensitive data is cloaked. JSON encoded.
	Values string `json:"values,omitempty"`
	// A SemVer 2 conformant version string of the chart.
	Version string `json:"version,omitempty"`
	// The rendered manifest as JSON.
	Manifest string `json:"manifest,omitempty"`
}

type helmReleaseProvider struct {
	host             *provider.HostClient
	helmDriver       string
	kubeConfig       *KubeConfig
	defaultNamespace string
	enableSecrets    bool
	name             string
	settings         *cli.EnvSettings
}

func newHelmReleaseProvider(
	host *provider.HostClient,
	config *rest.Config,
	clientConfig clientcmd.ClientConfig,
	helmDriver,
	namespace string,
	enableSecrets bool,
	pluginsDirectory,
	registryConfigPath,
	repositoryConfigPath,
	repositoryCache string,
) (customResourceProvider, error) {
	kc := newKubeConfig(config, clientConfig)
	settings := cli.New()
	settings.PluginsDirectory = pluginsDirectory
	settings.RegistryConfig = registryConfigPath
	settings.RepositoryConfig = repositoryConfigPath
	settings.RepositoryCache = repositoryCache

	return &helmReleaseProvider{
		host:             host,
		kubeConfig:       kc,
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
	if err := conf.Init(r.kubeConfig, namespace, r.helmDriver, debug); err != nil {
		return nil, err
	}
	return conf, nil
}

func decodeRelease(pm resource.PropertyMap) (*Release, error) {
	var release Release
	stripped := pm.MapRepl(nil, mapReplStripSecrets)
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &release,
		TagName: "json",
	})
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(stripped); err != nil {
		return nil, err
	}
	return &release, nil
}

func (r *helmReleaseProvider) Check(ctx context.Context, req *pulumirpc.CheckRequest, olds, news resource.PropertyMap) (*pulumirpc.CheckResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("Provider[%s].Check(%s)", r.name, urn)

	new, err := decodeRelease(news)
	if err != nil {
		return nil, err
	}

	if len(olds.Mappable()) > 0 {
		old, err := decodeRelease(olds)
		if err != nil {
			return nil, err
		}
		adoptOldNameIfUnnamed(new, old)

		if new.ReleaseSpec.Namespace == "" {
			new.ReleaseSpec.Namespace = old.ReleaseSpec.Namespace
		}
	} else {
		assignNameIfAutonammable(new, news, "release")
	}

	if new.ReleaseSpec.Namespace == "" {
		new.ReleaseSpec.Namespace = r.defaultNamespace
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

func adoptOldNameIfUnnamed(new, old *Release) {
	contract.Assert(old.ReleaseSpec.Name != "")
	new.ReleaseSpec.Name = old.ReleaseSpec.Name
}

func assignNameIfAutonammable(release *Release, pm resource.PropertyMap, base tokens.QName) {
	if rs, ok := pm["resourceSpec"].V.(resource.PropertyMap); ok {
		if name, ok := rs["name"]; ok && name.IsComputed() {
			return
		}
		if name, ok := rs["name"]; !ok || name.StringValue() == "" {
			release.ReleaseSpec.Name = fmt.Sprintf("%s-%s", base, metadata.RandString(8))
		}
	}
}

func (r *helmReleaseProvider) Diff(
	ctx context.Context,
	request *pulumirpc.DiffRequest,
	olds,
	news resource.PropertyMap,
) (*pulumirpc.DiffResponse, error) {
	// Extract old inputs from the `__inputs` field of the old state.
	oldInputs := parseCheckpointRelease(olds)
	diff := oldInputs.Diff(news)
	if diff == nil {
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

	// Always set desired state to DEPLOYED
	// TODO: This could be done in Check instead?
	newRelease.Status.Status = release.StatusDeployed.String()

	cpo, chartName, err := chartPathOptions(newRelease.ReleaseSpec)
	if err != nil {
		return nil, err
	}

	// Get Chart metadata, if we fail - we're done
	chart, _, err := getChart(chartName, r.settings, cpo)
	if err != nil {
		return nil, err
	}

	if newRelease.ReleaseSpec.Lint {
		if err := resourceReleaseValidate(newRelease.ReleaseSpec, r.settings, cpo); err != nil {
			return nil, err
		}
	}

	actionConfig, err := r.getActionConfig(oldRelease.ReleaseSpec.Namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewUpgrade(actionConfig)
	client.ChartPathOptions = *cpo
	client.Devel = newRelease.ReleaseSpec.Devel
	client.Namespace = newRelease.ReleaseSpec.Namespace
	client.Timeout = time.Duration(newRelease.ReleaseSpec.Timeout) * time.Second
	client.Wait = newRelease.ReleaseSpec.Wait
	client.DryRun = true // do not apply changes
	client.DisableHooks = newRelease.ReleaseSpec.DisableCRDHooks
	client.Atomic = newRelease.ReleaseSpec.Atomic
	client.SubNotes = newRelease.ReleaseSpec.RenderSubchartNotes
	client.WaitForJobs = newRelease.ReleaseSpec.WaitForJobs
	client.Force = newRelease.ReleaseSpec.ForceUpdate
	client.ResetValues = newRelease.ReleaseSpec.ResetValues
	client.ReuseValues = newRelease.ReleaseSpec.ReuseValues
	client.Recreate = newRelease.ReleaseSpec.RecreatePods
	client.MaxHistory = 0
	if newRelease.ReleaseSpec.MaxHistory != nil {
		client.MaxHistory = *newRelease.ReleaseSpec.MaxHistory
	}
	client.CleanupOnFail = newRelease.ReleaseSpec.CleanupOnFail
	client.Description = newRelease.ReleaseSpec.Description

	if cmd := newRelease.ReleaseSpec.Postrender; cmd != "" {
		pr, err := postrender.NewExec(cmd)
		if err != nil {
			return nil, err
		}
		client.PostRenderer = pr
	}

	values, err := getValues(newRelease.ReleaseSpec)
	if err != nil {
		return nil, fmt.Errorf("error getting values for a diff: %v", err)
	}

	dry, err := client.Run(newRelease.ReleaseSpec.Name, chart, values)
	if err != nil && strings.Contains(err.Error(), "has no deployed releases") {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("error running dry run for a diff: %v", err)
	}

	jsonManifest, err := convertYAMLManifestToJSON(dry.Manifest)
	if err != nil {
		return nil, err
	}
	newRelease.Status.Manifest = jsonManifest

	oldInputsJSON, err := json.Marshal(oldInputs)
	if err != nil {
		return nil, err
	}
	newInputsJSON, err := json.Marshal(news.Mappable())
	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.CreateMergePatch(oldInputsJSON, newInputsJSON)
	if err != nil {
		return nil, err
	}
	patchObj := map[string]interface{}{}
	if err = json.Unmarshal(patch, &patchObj); err != nil {
		return nil, pkgerrors.Wrapf(
			err, "Failed to check for changes in Helm release %s because of an error serializing "+
				"the JSON patch describing resource changes",
			newRelease.ReleaseSpec.Name)
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

		if detailedDiff, err = convertPatchToDiff(patchObj, oldInputs.Mappable(), news.Mappable(), oldInputs.Mappable(), ".releaseSpec.name", ".releaseSpec.namespace"); err != nil {
			return nil, pkgerrors.Wrapf(
				err, "Failed to check for changes in helm release %s/%s because of an error "+
					"converting JSON patch describing resource changes to a diff",
				newRelease.ReleaseSpec.Namespace, newRelease.ReleaseSpec.Name)
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

func resourceReleaseValidate(releaseSpec *ReleaseSpec, settings *cli.EnvSettings, cpo *action.ChartPathOptions) error {
	cpo, name, err := chartPathOptions(releaseSpec)
	if err != nil {
		return fmt.Errorf("malformed values: \n\t%s", err)
	}

	values, err := getValues(releaseSpec)
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

func (r *helmReleaseProvider) Create(ctx context.Context, req *pulumirpc.CreateRequest, news resource.PropertyMap) (*pulumirpc.CreateResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("Provider[%s].Create(%s)", r.name, urn)

	// If this is a preview and the input values contain unknowns, return them as-is. This is compatible with
	// prior behavior implemented by the Pulumi engine. Similarly, if the server does not support server-side
	// dry run, return the inputs as-is.
	if req.GetPreview() && news.ContainsUnknowns() {
		logger.V(9).Infof("cannot preview Create(%v)", urn)
		return &pulumirpc.CreateResponse{Id: "", Properties: req.GetProperties()}, nil
	}

	newRelease, err := decodeRelease(news)
	if err != nil {
		return nil, err
	}

	conf, err := r.getActionConfig(newRelease.ReleaseSpec.Namespace)
	if err != nil {
		return nil, err
	}
	client := action.NewInstall(conf)
	cpo, chartName, err := chartPathOptions(newRelease.ReleaseSpec)

	c, path, err := getChart(chartName, r.settings, cpo)
	if err != nil {
		return nil, err
	}

	// check and update the chart's dependencies if needed
	updated, err := checkChartDependencies(
		c,
		path,
		newRelease.ReleaseSpec.Keyring,
		r.settings,
		newRelease.ReleaseSpec.DependencyUpdate)
	if err != nil {
		return nil, err
	} else if updated {
		// load the chart again if its dependencies have been updated
		c, err = loader.Load(path)
		if err != nil {
			return nil, err
		}
	}

	values, err := getValues(newRelease.ReleaseSpec)
	if err != nil {
		return nil, err
	}

	err = isChartInstallable(c)
	if err != nil {
		return nil, err
	}

	client.ChartPathOptions = *cpo
	client.ClientOnly = false
	client.DryRun = false
	client.DisableHooks = newRelease.ReleaseSpec.DisableWebhooks
	client.Wait = newRelease.ReleaseSpec.Wait
	client.WaitForJobs = newRelease.ReleaseSpec.WaitForJobs
	client.Devel = newRelease.ReleaseSpec.Devel
	client.DependencyUpdate = newRelease.ReleaseSpec.DependencyUpdate
	client.Timeout = time.Duration(newRelease.ReleaseSpec.Timeout) * time.Second
	client.Namespace = newRelease.ReleaseSpec.Namespace
	client.ReleaseName = newRelease.ReleaseSpec.Name
	client.GenerateName = false
	client.NameTemplate = ""
	client.OutputDir = ""
	client.Atomic = newRelease.ReleaseSpec.Atomic
	client.SkipCRDs = newRelease.ReleaseSpec.SkipCrds
	client.SubNotes = newRelease.ReleaseSpec.RenderSubchartNotes
	client.DisableOpenAPIValidation = newRelease.ReleaseSpec.DisableOpenapiValidation
	client.Replace = newRelease.ReleaseSpec.Replace
	client.Description = newRelease.ReleaseSpec.Description
	client.CreateNamespace = newRelease.ReleaseSpec.CreateNamespace

	if cmd := newRelease.ReleaseSpec.Postrender; cmd != "" {
		pr, err := postrender.NewExec(cmd)

		if err != nil {
			return nil, err
		}

		client.PostRenderer = pr
	}

	rel, err := client.Run(c, values)
	if err != nil && rel == nil {
		return nil, err
	}

	if err != nil && rel != nil {
		actionConfig, err := r.getActionConfig(newRelease.ReleaseSpec.Namespace)
		if err != nil {
			return nil, err
		}
		exists, existsErr := resourceReleaseExists(newRelease.ReleaseSpec, actionConfig)

		if existsErr != nil {
			return nil, err
		}

		if !exists {
			return nil, err
		}

		//debug("%s Release was created but returned an error", logID)

		if err := setReleaseAttributes(newRelease, rel); err != nil {
			return nil, err
		}

		//return diag.Diagnostics{
		//	{
		//		Severity: diag.Warning,
		//		Summary:  fmt.Sprintf("Helm release %q was created but has a failed status. Use the `helm` command to investigate the error, correct it, then run Terraform again.", client.ReleaseName),
		//	},
		//	{
		//		Severity: diag.Error,
		//		Summary:  err.Error(),
		//	},
		//}
		// TODO: k.host.LogStatus instead.
		return nil, err

	}

	err = setReleaseAttributes(newRelease, rel)
	if err != nil {
		return nil, err
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
		id = fqName(newRelease.ReleaseSpec.Namespace, newRelease.ReleaseSpec.Name)
	}
	return &pulumirpc.CreateResponse{Id: id, Properties: inputsAndComputed}, nil
}

func (r *helmReleaseProvider) Read(ctx context.Context, req *pulumirpc.ReadRequest, oldState, oldInputs resource.PropertyMap) (*pulumirpc.ReadResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("Provider[%s].Read(%s)", r.name, urn)

	existingRelease, err := decodeRelease(oldState)
	if err != nil {
		return nil, err
	}
	actionConfig, err := r.getActionConfig(existingRelease.ReleaseSpec.Namespace)
	if err != nil {
		return nil, err
	}
	exists, err := resourceReleaseExists(existingRelease.ReleaseSpec, actionConfig)
	if err != nil {
		return nil, err
	}

	if !exists {
		// If not found, this resource was probably deleted.
		return deleteResponse, nil
	}

	name := existingRelease.ReleaseSpec.Name
	liveObj, err := getRelease(actionConfig, name)
	if err != nil {
		return nil, err
	}

	err = setReleaseAttributes(existingRelease, liveObj)
	if err != nil {
		return nil, err
	}

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

	id := fqName(existingRelease.ReleaseSpec.Namespace, existingRelease.ReleaseSpec.Name)
	if reqID := req.GetId(); len(reqID) > 0 {
		id = reqID
	}

	return &pulumirpc.ReadResponse{Id: id, Properties: state, Inputs: inputs}, nil
}

func (r *helmReleaseProvider) Update(ctx context.Context, req *pulumirpc.UpdateRequest, oldState, newResInputs resource.PropertyMap) (*pulumirpc.UpdateResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("Provider[%s].Update(%s)", r.name, urn)
	logger.V(9).Infof("%s executing", label)

	newRelease, err := decodeRelease(newResInputs)
	if err != nil {
		return nil, err
	}
	actionConfig, err := r.getActionConfig(newRelease.ReleaseSpec.Namespace)
	if err != nil {
		return nil, err
	}

	cpo, chartName, err := chartPathOptions(newRelease.ReleaseSpec)
	if err != nil {
		return nil, err
	}

	c, path, err := getChart(chartName, r.settings, cpo)
	if err != nil {
		return nil, err
	}

	// check and update the chart's dependencies if needed
	updated, err := checkChartDependencies(
		c,
		path,
		newRelease.ReleaseSpec.Keyring,
		r.settings,
		newRelease.ReleaseSpec.DependencyUpdate)
	if err != nil {
		return nil, err
	} else if updated {
		// load the chart again if its dependencies have been updated
		c, err = loader.Load(path)
		if err != nil {
			return nil, err
		}
	}

	client := action.NewUpgrade(actionConfig)
	client.ChartPathOptions = *cpo
	client.Devel = newRelease.ReleaseSpec.Devel
	client.Namespace = newRelease.ReleaseSpec.Namespace
	client.Timeout = time.Duration(newRelease.ReleaseSpec.Timeout) * time.Second
	client.Wait = newRelease.ReleaseSpec.Wait
	client.WaitForJobs = newRelease.ReleaseSpec.WaitForJobs
	client.DryRun = false
	client.DisableHooks = newRelease.ReleaseSpec.DisableWebhooks
	client.Atomic = newRelease.ReleaseSpec.Atomic
	client.SkipCRDs = newRelease.ReleaseSpec.SkipCrds
	client.SubNotes = newRelease.ReleaseSpec.RenderSubchartNotes
	client.Force = newRelease.ReleaseSpec.ForceUpdate
	client.ResetValues = newRelease.ReleaseSpec.ResetValues
	client.ReuseValues = newRelease.ReleaseSpec.ReuseValues
	client.Recreate = newRelease.ReleaseSpec.RecreatePods
	client.MaxHistory = 0
	if newRelease.ReleaseSpec.MaxHistory != nil {
		client.MaxHistory = *newRelease.ReleaseSpec.MaxHistory
	}
	client.CleanupOnFail = newRelease.ReleaseSpec.CleanupOnFail
	client.Description = newRelease.ReleaseSpec.Description

	if cmd := newRelease.ReleaseSpec.Postrender; cmd != "" {
		pr, err := postrender.NewExec(cmd)

		if err != nil {
			return nil, err
		}

		client.PostRenderer = pr
	}

	values, err := getValues(newRelease.ReleaseSpec)
	if err != nil {
		return nil, err
	}

	name := newRelease.ReleaseSpec.Name
	updatedRelease, err := client.Run(name, c, values)
	if err != nil {
		return nil, err
	}

	err = setReleaseAttributes(newRelease, updatedRelease)
	if err != nil {
		return nil, err
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

func (r *helmReleaseProvider) Delete(ctx context.Context, request *pulumirpc.DeleteRequest, olds resource.PropertyMap) (*empty.Empty, error) {
	release, err := decodeRelease(olds)
	if err != nil {
		return nil, err
	}

	namespace := release.ReleaseSpec.Namespace
	actionConfig, err := r.getActionConfig(namespace)
	if err != nil {
		return nil, err
	}

	name := release.ReleaseSpec.Name

	res, err := action.NewUninstall(actionConfig).Run(name)
	if err != nil {
		return nil, err
	}

	if res.Info != "" {
		r.host.Log(context.Background(), diag.Warning, "Helm uninstall returned information: %q", res.Info)
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
func parseCheckpointRelease(obj resource.PropertyMap) resource.PropertyMap {
	if inputs, ok := obj["__inputs"]; ok {
		return inputs.SecretValue().Element.ObjectValue()
	}

	return nil
}

func setReleaseAttributes(release *Release, r *release.Release) error {
	release.Status.Version = r.Chart.Metadata.Version
	release.Status.Namespace = r.Namespace
	release.Status.Name = r.Name
	release.Status.Status = r.Info.Status.String()

	//cloakSetValues(r.Config, news)
	values, err := json.Marshal(r.Config)
	if err != nil {
		return err
	}

	jsonManifest, err := convertYAMLManifestToJSON(r.Manifest)
	if err != nil {
		return err
	}
	//manifest := redactSensitiveValues(jsonManifest, release.ReleaseSpec.Values)
	release.Status.Manifest = jsonManifest

	release.Status.Name = r.Name
	release.Status.Namespace = r.Namespace
	release.Status.Revision = &r.Version
	release.Status.Chart = r.Chart.Metadata.Name
	release.Status.Version = r.Chart.Metadata.Version
	release.Status.AppVersion = r.Chart.Metadata.AppVersion
	release.Status.Values = string(values)
	return nil
}

func resourceReleaseExists(releaseSpec *ReleaseSpec, actionConfig *action.Configuration) (bool, error) {
	logger.V(9).Infof("[resourceReleaseExists: %s]", releaseSpec.Name)
	name := releaseSpec.Name
	_, err := getRelease(actionConfig, name)

	logger.V(9).Infof("[resourceReleaseExists: %s] Done", releaseSpec.Name)

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
	debug("%s getRelease post action created", name)

	res, err := get.Run(name)
	debug("%s getRelease post run", name)

	if err != nil {
		debug("getRelease for %s errored", name)
		debug("%v", err)
		if strings.Contains(err.Error(), "release: not found") {
			return nil, errReleaseNotFound
		}

		debug("could not get release %s", err)

		return nil, err
	}

	debug("%s getRelease done", name)

	return res, nil
}

func isChartInstallable(ch *helmchart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return fmt.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func getValues(spec *ReleaseSpec) (map[string]interface{}, error) {
	base := map[string]interface{}{}

	for _, value := range spec.Values {
		if value == "" {
			continue
		}

		currentMap := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(value), &currentMap); err != nil {
			return nil, fmt.Errorf("---> %v %s", err, value)
		}

		base = mergeMaps(base, currentMap)
	}

	for _, set := range spec.Set {
		if err := getValue(base, set); err != nil {
			return nil, err
		}
	}

	//for _, set := range spec.SetSensitive {
	//	if err := getValue(base, set); err != nil {
	//		return nil, err
	//	}
	//}

	return base, logValues(base, spec)
}

func getValue(base map[string]interface{}, set *SetValue) error {
	name := set.Name
	value := set.Value
	valueType := set.Type

	switch valueType {
	case "auto", "":
		if err := strvals.ParseInto(fmt.Sprintf("%s=%s", name, value), base); err != nil {
			return fmt.Errorf("failed parsing key %q with value %s, %s", name, value, err)
		}
	case "string":
		if err := strvals.ParseIntoString(fmt.Sprintf("%s=%s", name, value), base); err != nil {
			return fmt.Errorf("failed parsing key %q with value %s, %s", name, value, err)
		}
	default:
		return fmt.Errorf("unexpected type: %s", valueType)
	}

	return nil
}

func logValues(values map[string]interface{}, spec *ReleaseSpec) error {
	// copy array to avoid change values by the cloak function.
	//asJSON, _ := json.Marshal(values)
	//var c map[string]interface{}
	//err := json.Unmarshal(asJSON, &c)
	//if err != nil {
	//	return err
	//}
	//
	//cloakSetValues(c, spec)
	//
	//y, err := yaml.Marshal(c)
	//if err != nil {
	//	return err
	//}
	//
	//log.Printf(
	//	"---[ values.yaml ]-----------------------------------\n%s\n",
	//	string(y),
	//)

	return nil
}

func cloakSetValues(config map[string]interface{}, pm resource.PropertyMap) {
	//if rs, ok := pm["resourceSpec"].V.(resource.PropertyMap); ok {
	//	if set, ok := rs["set"]; ok && set.ContainsSecrets() {
	//		set.SecretValue().Element
	//	}
	//}
	//
	//for _, raw := range d.Get("set_sensitive").(*schema.Set).List() {
	//	set := raw.(map[string]interface{})
	//	cloakSetValue(config, set["name"].(string))
	//}
}

const sensitiveContentValue = "(sensitive value)"

func cloakSetValue(values map[string]interface{}, valuePath string) {
	pathKeys := strings.Split(valuePath, ".")
	sensitiveKey := pathKeys[len(pathKeys)-1]
	parentPathKeys := pathKeys[:len(pathKeys)-1]

	m := values
	for _, key := range parentPathKeys {
		v, ok := m[key].(map[string]interface{})
		if !ok {
			return
		}
		m = v
	}

	m[sensitiveKey] = sensitiveContentValue
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
	//Load function blows up if accessed concurrently
	path, err := cpo.LocateChart(name, settings)
	if err != nil {
		return nil, "", err
	}

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

func chartPathOptions(releaseSpec *ReleaseSpec) (*action.ChartPathOptions, string, error) {
	chartName := releaseSpec.Chart

	repository := releaseSpec.RepositorySpec.Repository
	repositoryURL, chartName, err := resolveChartName(repository, strings.TrimSpace(chartName))

	if err != nil {
		return nil, "", err
	}
	version := getVersion(releaseSpec)

	return &action.ChartPathOptions{
		CaFile:   releaseSpec.RepositorySpec.RepositoryCAFile,
		CertFile: releaseSpec.RepositorySpec.RepositoryCertFile,
		KeyFile:  releaseSpec.RepositorySpec.RepositoryKeyFile,
		//Keyring:  d.Get("keyring").(string),
		RepoURL:  repositoryURL,
		Verify:   releaseSpec.Verify,
		Version:  version,
		Username: releaseSpec.RepositorySpec.RepositoryUsername,
		Password: releaseSpec.RepositorySpec.RepositoryPassword, // TODO: This should already be resolved.
	}, chartName, nil
}

func getVersion(releaseSpec *ReleaseSpec) (version string) {
	version = releaseSpec.Version

	if version == "" && releaseSpec.Devel {
		debug("setting version to >0.0.0-0")
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

func asHelmRelease(pm resource.PropertyMap) (*Release, error) {
	obj := pm.MapRepl(nil, mapReplStripSecrets)
	var release Release
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &release,
		TagName: "json",
	})
	contract.AssertNoError(err)
	if err := decoder.Decode(obj); err != nil {
		return nil, err
	}

	return &release, nil
}
