/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// package helm contains code vendored from the upstream Helm project.
package helm

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/util/homedir"

	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/logging"
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
)

type InitActionConfigF func(actionConfig *action.Configuration, namespaceOverride string) error

type LocateChartF func(i *action.Install, name string, settings *cli.EnvSettings) (string, error)

type ExecuteF func(
	ctx context.Context, i *action.Install, chrt *chart.Chart, vals map[string]interface{},
) (*release.Release, error)

// Tool for executing Helm commands via the Helm library.
type Tool struct {
	EnvSettings *cli.EnvSettings
	HelmDriver  string

	initActionConfig InitActionConfigF
	locateChart      LocateChartF
	execute          ExecuteF
}

// NewTool creates a new Helm tool with the given environment settings.
func NewTool(settings *cli.EnvSettings) *Tool {
	helmDriver := os.Getenv("HELM_DRIVER")
	logger.V(6).Infof("initializing Helm tool: driver=%q, settings=%+v", helmDriver, *settings)

	return &Tool{
		EnvSettings: settings,
		HelmDriver:  helmDriver,
		initActionConfig: func(actionConfig *action.Configuration, namespaceOverride string) error {
			if namespaceOverride == "" {
				namespaceOverride = settings.Namespace()
			}
			// https://github.com/helm/helm/blob/635b8cf33d25a86131635c32f35b2a76256e40cb/cmd/helm/helm.go#L72-L81
			return actionConfig.Init(settings.RESTClientGetter(), namespaceOverride, helmDriver, debug)
		},
		locateChart: func(i *action.Install, name string, settings *cli.EnvSettings) (string, error) {
			// Perform a login against the registry if necessary.
			if i.Username != "" && i.Password != "" {
				u, err := url.Parse(name)
				if err != nil {
					return "", err
				}
				c := i.GetRegistryClient()
				// Login can fail for harmless reasons like already being
				// logged in. Optimistically ignore those errors.
				if err := c.Login(u.Host, registry.LoginOptBasicAuth(i.Username, i.Password)); err != nil {
					logger.V(6).Infof("[helm] %s", fmt.Sprintf("login error: %s", err))
				}
			}
			return i.LocateChart(name, settings)
		},
		execute: func(
			ctx context.Context, i *action.Install, chrt *chart.Chart, vals map[string]interface{},
		) (*release.Release, error) {
			return i.RunWithContext(ctx, chrt, vals)
		},
	}
}

func (t *Tool) AllGetters() getter.Providers {
	return getter.All(t.EnvSettings)
}

// TemplateOrInstallCommand for `helm template` or `helm install`.
type TemplateOrInstallCommand struct {
	// Install parameters.
	*action.Install

	// Chart is the chart specification, which can be:
	// - a path to a local chart directory or archive file
	// - a qualified chart reference (e.g., "stable/mariadb") based on the local repository configuration
	// - a URL to a remote chart (https://, oci://, file://, etc.)
	Chart string

	// Values to be applied to the chart.
	Values ValueOpts

	tool         *Tool
	actionConfig *action.Configuration
}

func (cmd *TemplateOrInstallCommand) addFlags() {
	cmd.addInstallFlags()
}

func (cmd *TemplateOrInstallCommand) addInstallFlags() {
	// https://github.com/helm/helm/blob/14d0c13e9eefff5b4a1b511cf50643529692ec94/cmd/helm/install.go#L176-L203
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
	// https://github.com/helm/helm/blob/14d0c13e9eefff5b4a1b511cf50643529692ec94/cmd/helm/flags.go#L45-L51
	v := cmd.Values
	v.Values = map[string]any{}
	v.ValuesFiles = []pulumi.Asset{}
}

func (cmd *TemplateOrInstallCommand) addChartPathOptionsFlags() {
	// https://github.com/helm/helm/blob/14d0c13e9eefff5b4a1b511cf50643529692ec94/cmd/helm/flags.go#L54-L66
	c := &cmd.ChartPathOptions
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

// TemplateCommand for `helm template`.
type TemplateCommand struct {
	TemplateOrInstallCommand

	// https://github.com/helm/helm/blob/635b8cf33d25a86131635c32f35b2a76256e40cb/cmd/helm/template.go#L50-L57

	Validate    bool
	IncludeCRDs bool
	SkipTests   bool
}

// Template returns a new `helm template` command.
func (t *Tool) Template() *TemplateCommand {
	actionConfig := new(action.Configuration)

	cmd := &TemplateCommand{
		TemplateOrInstallCommand: TemplateOrInstallCommand{
			tool:         t,
			actionConfig: actionConfig,
			Install:      action.NewInstall(actionConfig),
			Values:       ValueOpts{},
		},
	}

	// https://github.com/helm/helm/blob/635b8cf33d25a86131635c32f35b2a76256e40cb/cmd/helm/template.go#L192-L203
	cmd.addFlags()
	cmd.OutputDir = ""
	cmd.Validate = false
	cmd.IncludeCRDs = false
	cmd.SkipTests = false
	cmd.IsUpgrade = false
	cmd.UseReleaseName = false
	return cmd
}

// Execute runs the `helm template` command.
func (cmd *TemplateCommand) Execute(ctx context.Context) (*release.Release, error) {
	client := cmd.Install

	err := cmd.tool.initActionConfig(cmd.actionConfig, cmd.Namespace)
	if err != nil {
		return nil, err
	}

	// https://github.com/helm/helm/blob/635b8cf33d25a86131635c32f35b2a76256e40cb/cmd/helm/template.go#L76-L81
	registryClient, err := newRegistryClient(cmd.tool.EnvSettings, client.CertFile, client.KeyFile, client.CaFile,
		client.InsecureSkipTLSverify, client.PlainHTTP, client.Username, client.Password)
	if err != nil {
		return nil, fmt.Errorf("missing registry client: %w", err)
	}
	client.SetRegistryClient(registryClient)

	// https://github.com/helm/helm/blob/635b8cf33d25a86131635c32f35b2a76256e40cb/cmd/helm/template.go#L88-L94
	if client.DryRunOption == "" {
		client.DryRunOption = "true"
	}
	client.DryRun = true
	// client.ReleaseName = "release-name"
	client.Replace = true // Skip the name check
	client.ClientOnly = !cmd.Validate
	client.IncludeCRDs = cmd.IncludeCRDs
	client.PlainHTTP = cmd.PlainHTTP
	client.Username = cmd.Username
	client.Password = cmd.Password
	client.CaFile = cmd.CaFile
	client.KeyFile = cmd.KeyFile
	client.CertFile = cmd.CertFile

	return cmd.runInstall(ctx)
}

// runInstall runs the install action.
// https://github.com/helm/helm/blob/14d0c13e9eefff5b4a1b511cf50643529692ec94/cmd/helm/install.go#L221
func (cmd *TemplateOrInstallCommand) runInstall(ctx context.Context) (*release.Release, error) {
	settings := cmd.tool.EnvSettings
	client := cmd.Install
	valueOpts := cmd.Values

	if client.Version == "" && client.Devel {
		debug("setting version to >0.0.0-0")
		client.Version = ">0.0.0-0"
	}

	releaseName, chart, err := client.NameAndChart([]string{cmd.Chart})
	if err != nil {
		return nil, err
	}
	client.ReleaseName = releaseName

	debug("attempting to resolve the chart %q with version %q", chart, client.Version)
	cp, err := cmd.tool.locateChart(client, chart, settings)
	if err != nil {
		return nil, errors.Wrap(err, "unable to locate the chart")
	}
	debug("a chart was located at %s", cp)

	p := cmd.tool.AllGetters()
	// FUTURE: add a "file:" getter for parity with Pulumi resource package
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, errors.Wrap(err, "unable to process the chart values")
	}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load the chart")
	}

	if err := checkIfInstallable(chartRequested); err != nil {
		return nil, err
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			err = errors.Wrap(
				err,
				"An error occurred while checking for chart dependencies. "+
					"You may need to run `helm dependency build` to fetch missing dependencies",
			)
			if client.DependencyUpdate || chartRequested.Lock != nil {
				logStream := debugStream()
				defer logStream.Close()

				man := &downloader.Manager{
					Out:              logStream,
					ChartPath:        cp,
					Keyring:          client.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
					Debug:            settings.Debug,
					RegistryClient:   client.GetRegistryClient(),
				}
				if client.DependencyUpdate {
					if err2 := man.Update(); err2 != nil {
						debug("unable to update dependencies: %s", err2.Error())
						return nil, err
					}
				} else if chartRequested.Lock != nil {
					// Pulumi behavior: automatically build the dependencies if a lock file is present
					if err2 := man.Build(); err2 != nil {
						debug("unable to build dependencies: %s", err2.Error())
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

	return cmd.tool.execute(ctx, client, chartRequested, vals)
}

func debug(format string, a ...any) {
	logger.V(6).Infof("[helm] %s", fmt.Sprintf(format, a...))
}

func debugStream() *logging.LogWriter {
	// FUTURE: set log depth
	return logging.NewLogWriter(logger.V(6).Infof, logging.WithPrefix("[helm] "))
}

// defaultKeyring returns the expanded path to the default keyring.
// https://github.com/helm/helm/blob/635b8cf33d25a86131635c32f35b2a76256e40cb/cmd/helm/root.go#L276-L293
func defaultKeyring() string {
	if v, ok := os.LookupEnv("GNUPGHOME"); ok {
		return filepath.Join(v, "pubring.gpg")
	}
	return filepath.Join(homedir.HomeDir(), ".gnupg", "pubring.gpg")
}

// newRegistryClient returns a new registry client
// https://github.com/helm/helm/blob/01adbab466b6133936cac0c56a99274715f7c085/pkg/cmd/root.go#L345-L360
func newRegistryClient(settings *cli.EnvSettings,
	certFile, keyFile, caFile string, insecureSkipVerify, plainHTTP bool, username, password string,
) (*registry.Client, error) {
	logStream := debugStream()
	opts := []registry.ClientOption{
		registry.ClientOptDebug(settings.Debug),
		registry.ClientOptEnableCache(true),
		registry.ClientOptWriter(logStream),
		registry.ClientOptCredentialsFile(settings.RegistryConfig),
	}
	if plainHTTP {
		opts = append(opts, registry.ClientOptPlainHTTP())
	}

	if certFile != "" && keyFile != "" || caFile != "" || insecureSkipVerify {
		tlsConf, err := newTLSConfig(certFile, keyFile, caFile, insecureSkipVerify)
		if err != nil {
			return nil, err
		}
		tlsOpt := registry.ClientOptHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConf,
				Proxy:           http.ProxyFromEnvironment,
			},
		})
		opts = append(opts, tlsOpt)
	}

	opts = append(opts, registry.ClientOptBasicAuth(username, password))

	registryClient, err := registry.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return registryClient, nil
}

// checkIfInstallable validates if a chart can be installed
//
// Application chart type is only installable
// https://github.com/helm/helm/blob/635b8cf33d25a86131635c32f35b2a76256e40cb/cmd/helm/install.go#L317-L326
func checkIfInstallable(ch *chart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

// validateDryRunOptionFlag validates the dry-run flag value
// https://github.com/helm/helm/blob/635b8cf33d25a86131635c32f35b2a76256e40cb/cmd/helm/install.go#L340-L354
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

func ApplyRepositoryOpts(cpo *action.ChartPathOptions, p getter.Providers, repoOpts helmv4.RepositoryOpts) error {
	if repoOpts.CaFile != nil {
		file, _, err := downloadAsset(p, repoOpts.CaFile)
		if err != nil {
			return fmt.Errorf("cafile: %w", err)
		}
		cpo.CaFile = file
	}
	if repoOpts.CertFile != nil {
		file, _, err := downloadAsset(p, repoOpts.CertFile)
		if err != nil {
			return fmt.Errorf("certfile: %w", err)
		}
		cpo.CertFile = file
	}
	if repoOpts.KeyFile != nil {
		file, _, err := downloadAsset(p, repoOpts.KeyFile)
		if err != nil {
			return fmt.Errorf("keyfile: %w", err)
		}
		cpo.KeyFile = file
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
	return nil
}

type cleanupF func() error

// downloadAsset downloads an asset to the local filesystem.
func downloadAsset(p getter.Providers, asset pulumi.AssetOrArchive) (string, cleanupF, error) {
	a, isAsset := asset.(pulumi.Asset)
	if !isAsset {
		return "", nil, errors.New("expected an asset")
	}
	makeTemp := func(data []byte) (string, cleanupF, error) {
		file, err := os.CreateTemp("", "pulumi-")
		if err != nil {
			return "", nil, err
		}
		defer file.Close()
		if _, err := file.Write(data); err != nil {
			return "", nil, err
		}
		return file.Name(), func() error {
			return os.Remove(file.Name())
		}, err
	}

	switch {
	case a.Text() != "":
		return makeTemp([]byte(a.Text()))
	case a.Path() != "":
		return a.Path(), func() error { return nil }, nil
	case a.URI() != "":
		u, err := url.Parse(a.URI())
		if err != nil {
			return "", nil, err
		}
		g, err := p.ByScheme(u.Scheme)
		if err != nil {
			return "", nil, fmt.Errorf("no protocol handler for uri %q", a.URI())
		}
		data, err := g.Get(a.URI(), getter.WithURL(a.URI()))
		if err != nil {
			return "", nil, fmt.Errorf("failed to read uri %q: %w", a.URI(), err)
		}
		return makeTemp(data.Bytes())
	default:
		return "", nil, errors.New("unrecognized asset type")
	}
}
