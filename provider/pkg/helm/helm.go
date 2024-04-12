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

package helm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/util/homedir"
)

type Tool struct {
	EnvSettings *cli.EnvSettings
	HelmDriver  string
}

func NewTool(settings *cli.EnvSettings) *Tool {
	return &Tool{
		EnvSettings: settings,
	}
}

func (t *Tool) initialize(actionConfig *action.Configuration, namespaceOverride string) error {
	if namespaceOverride == "" {
		namespaceOverride = t.EnvSettings.Namespace()
	}
	return actionConfig.Init(t.EnvSettings.RESTClientGetter(), namespaceOverride, t.HelmDriver, debug)
}

type TemplateOrInstallCommand struct {
	*action.Install
	tool         *Tool
	actionConfig *action.Configuration
	valueOpts    *values.Options

	// Chart is the chart specification, which can be:
	// - a path to a local chart directory or archive file
	// - a qualified chart reference (e.g., "stable/mariadb") based on the local repository configuration
	// - a URL to a remote chart (https://, oci://, file://, etc.)
	Chart string

	// DependencyBuild performs "helm dependency build" before installing the chart.
	DependencyBuild bool
}

func (cmd *TemplateOrInstallCommand) addFlags() {
	cmd.DependencyBuild = false
	cmd.addInstallFlags()
}

func (cmd *TemplateOrInstallCommand) addInstallFlags() {
	client := cmd.Install
	client.CreateNamespace = false
	client.DryRunOption = "client"
	client.Force = false
	client.DisableHooks = false
	client.Replace = false
	client.Timeout = 300 * time.Second
	client.Wait = false
	client.WaitForJobs = false
	client.GenerateName = false
	client.NameTemplate = ""
	client.Description = ""
	client.Devel = false
	client.DependencyUpdate = false
	client.DisableOpenAPIValidation = false
	client.Atomic = false
	client.SkipCRDs = false
	client.SubNotes = false
	client.Labels = nil
	client.EnableDNS = false
	cmd.addValueOptionsFlags()
	cmd.addChartPathOptionsFlags()
}

func (cmd *TemplateOrInstallCommand) addValueOptionsFlags() {
	// https://github.com/helm/helm/blob/14d0c13e9eefff5b4a1b511cf50643529692ec94/cmd/helm/flags.go#L45-51
	v := cmd.valueOpts
	v.ValueFiles = []string{}
	v.Values = []string{}
	v.StringValues = []string{}
	v.FileValues = []string{}
	v.JSONValues = []string{}
	v.LiteralValues = []string{}
}

func (cmd *TemplateOrInstallCommand) addChartPathOptionsFlags() {
	// https://github.com/helm/helm/blob/14d0c13e9eefff5b4a1b511cf50643529692ec94/cmd/helm/flags.go#L54-66
	c := &cmd.Install.ChartPathOptions
	c.Version = ""
	c.Verify = false
	c.Keyring = defaultKeyring()
	c.RepoURL = ""
	c.Username = ""
	c.Password = ""
	c.CertFile = ""
	c.KeyFile = ""
	c.InsecureSkipTLSverify = false
	c.PlainHTTP = false
	c.CaFile = ""
	c.PassCredentialsAll = false
}

type InstallCommand struct {
	TemplateOrInstallCommand
}

type TemplateCommand struct {
	TemplateOrInstallCommand

	// https://github.com/helm/helm/blob/635b8cf33d25a86131635c32f35b2a76256e40cb/cmd/helm/template.go#L50-L57
	Validate    bool
	IncludeCRDs bool
	SkipTests   bool
	KubeVersion string
	ExtraAPIs   []string
}

func (t *Tool) Template() *TemplateCommand {
	actionConfig := new(action.Configuration)

	cmd := &TemplateCommand{
		TemplateOrInstallCommand: TemplateOrInstallCommand{
			tool:         t,
			actionConfig: actionConfig,
			valueOpts:    &values.Options{},
			Install:      action.NewInstall(actionConfig),
		},
	}

	cmd.addFlags()
	cmd.Install.OutputDir = ""
	cmd.Validate = false
	cmd.IncludeCRDs = false
	cmd.SkipTests = false
	cmd.Install.IsUpgrade = false
	cmd.KubeVersion = ""
	cmd.ExtraAPIs = []string{}
	cmd.Install.UseReleaseName = false
	return cmd
}

// Execute runs the `helm template` command.
func (cmd *TemplateCommand) Execute(ctx context.Context) (*release.Release, error) {
	client := cmd.Install

	err := cmd.tool.initialize(cmd.actionConfig, cmd.Namespace)
	if err != nil {
		return nil, err
	}

	// https://github.com/helm/helm/blob/635b8cf33d25a86131635c32f35b2a76256e40cb/cmd/helm/template.go#L68-L74
	if cmd.KubeVersion != "" {
		parsedKubeVersion, err := chartutil.ParseKubeVersion(cmd.KubeVersion)
		if err != nil {
			return nil, fmt.Errorf("invalid kube version '%s': %s", cmd.KubeVersion, err)
		}
		client.KubeVersion = parsedKubeVersion
	}

	// https://github.com/helm/helm/blob/635b8cf33d25a86131635c32f35b2a76256e40cb/cmd/helm/template.go#L76-L81
	registryClient, err := newRegistryClient(cmd.tool.EnvSettings, client.CertFile, client.KeyFile, client.CaFile,
		client.InsecureSkipTLSverify, client.PlainHTTP)
	if err != nil {
		return nil, fmt.Errorf("missing registry client: %w", err)
	}
	client.SetRegistryClient(registryClient)

	// https://github.com/helm/helm/blob/635b8cf33d25a86131635c32f35b2a76256e40cb/cmd/helm/template.go#L88-L94
	client.DryRunOption = "client"
	client.DryRun = true
	// client.ReleaseName = "release-name"
	client.Replace = true // Skip the name check
	client.ClientOnly = !cmd.Validate
	client.APIVersions = chartutil.VersionSet(cmd.ExtraAPIs)
	client.IncludeCRDs = cmd.IncludeCRDs

	return cmd.runInstall(ctx)
}

func (cmd *TemplateOrInstallCommand) runInstall(ctx context.Context) (*release.Release, error) {
	settings := cmd.tool.EnvSettings
	client := cmd.Install
	valueOpts := cmd.valueOpts

	if client.Version == "" && client.Devel {
		debug("setting version to >0.0.0-0")
		client.Version = ">0.0.0-0"
	}

	name, chart, err := client.NameAndChart([]string{cmd.Chart})
	if err != nil {
		return nil, err
	}
	client.ReleaseName = name

	cp, err := client.ChartPathOptions.LocateChart(chart, settings)
	if err != nil {
		return nil, err
	}

	debug("CHART PATH: %s\n", cp)

	p := getter.All(settings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, err
	}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	if err := checkIfInstallable(chartRequested); err != nil {
		return nil, err
	}

	if chartRequested.Metadata.Deprecated {
		warning("This chart is deprecated")
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			err = errors.Wrap(err, "An error occurred while checking for chart dependencies. You may need to run `helm dependency build` to fetch missing dependencies")
			if client.DependencyUpdate || cmd.DependencyBuild {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
					Debug:            settings.Debug,
					RegistryClient:   client.GetRegistryClient(),
				}
				if cmd.DependencyBuild {
					if err := man.Build(); err != nil {
						return nil, err
					}
				}
				if client.DependencyUpdate {
					if err := man.Update(); err != nil {
						return nil, err
					}
				}
				// Reload the chart with the updated Chart.lock file.
				if chartRequested, err = loader.Load(cp); err != nil {
					return nil, errors.Wrap(err, "failed reloading chart after repo update")
				}
			} else {
				return nil, err
			}
		}
	}

	if client.Namespace == "" {
		client.Namespace = settings.Namespace()
	}

	// Validate DryRunOption member is one of the allowed values
	if err := validateDryRunOptionFlag(client.DryRunOption); err != nil {
		return nil, err
	}

	return client.RunWithContext(ctx, chartRequested, vals)
}

func debug(format string, a ...any) {
	logger.V(6).Infof("[DEBUG] %s", fmt.Sprintf(format, a...))
}

func warning(format string, v ...interface{}) {
	format = fmt.Sprintf("WARNING: %s\n", format)
	fmt.Fprintf(os.Stderr, format, v...)
}

// defaultKeyring returns the expanded path to the default keyring.
func defaultKeyring() string {
	if v, ok := os.LookupEnv("GNUPGHOME"); ok {
		return filepath.Join(v, "pubring.gpg")
	}
	return filepath.Join(homedir.HomeDir(), ".gnupg", "pubring.gpg")
}

func newRegistryClient(settings *cli.EnvSettings, certFile, keyFile, caFile string, insecureSkipTLSverify, plainHTTP bool) (*registry.Client, error) {
	if certFile != "" && keyFile != "" || caFile != "" || insecureSkipTLSverify {
		registryClient, err := newRegistryClientWithTLS(settings, certFile, keyFile, caFile, insecureSkipTLSverify)
		if err != nil {
			return nil, err
		}
		return registryClient, nil
	}
	registryClient, err := newDefaultRegistryClient(settings, plainHTTP)
	if err != nil {
		return nil, err
	}
	return registryClient, nil
}

func newDefaultRegistryClient(settings *cli.EnvSettings, plainHTTP bool) (*registry.Client, error) {
	opts := []registry.ClientOption{
		registry.ClientOptDebug(settings.Debug),
		registry.ClientOptEnableCache(true),
		registry.ClientOptWriter(os.Stderr),
		registry.ClientOptCredentialsFile(settings.RegistryConfig),
	}
	if plainHTTP {
		opts = append(opts, registry.ClientOptPlainHTTP())
	}

	// Create a new registry client
	registryClient, err := registry.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return registryClient, nil
}

func newRegistryClientWithTLS(settings *cli.EnvSettings, certFile, keyFile, caFile string, insecureSkipTLSverify bool) (*registry.Client, error) {
	// Create a new registry client
	registryClient, err := registry.NewRegistryClientWithTLS(os.Stderr, certFile, keyFile, caFile, insecureSkipTLSverify,
		settings.RegistryConfig, settings.Debug,
	)
	if err != nil {
		return nil, err
	}
	return registryClient, nil
}

// checkIfInstallable validates if a chart can be installed
//
// Application chart type is only installable
func checkIfInstallable(ch *chart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func validateDryRunOptionFlag(dryRunOptionFlagValue string) error {
	// Validate dry-run flag value with a set of allowed value
	allowedDryRunValues := []string{"false", "true", "none", "client", "server"}
	isAllowed := false
	for _, v := range allowedDryRunValues {
		if dryRunOptionFlagValue == v {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return errors.New("Invalid dry-run flag. Flag must one of the following: false, true, none, client, server")
	}
	return nil
}

func ApplyRepositoryOpts(cpo *action.ChartPathOptions, repoOpts helmv3.RepositoryOpts) {
	if repoOpts.CaFile != nil {
		cpo.CaFile = *repoOpts.CaFile
	}
	if repoOpts.CertFile != nil {
		cpo.CertFile = *repoOpts.CertFile
	}
	if repoOpts.KeyFile != nil {
		cpo.KeyFile = *repoOpts.KeyFile
	}
	if repoOpts.Username != nil {
		cpo.Username = *repoOpts.Username
	}
	if repoOpts.Password != nil {
		cpo.Password = *repoOpts.Password
	}
	if repoOpts.Repo != nil {
		cpo.RepoURL = *repoOpts.Repo
	}
}
