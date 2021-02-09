// Copyright 2016-2020, Pulumi Corporation.
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
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	pkgerrors "github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
)

// testAnnotation matches test-related Helm hook annotations (test, test-success, test-failure)
var testAnnotation = regexp.MustCompile(`"?helm.sh\/hook"?:.*test`)

type HelmFetchOpts struct {
	CAFile      string `json:"ca_file,omitempty"`
	CertFile    string `json:"cert_file,omitempty"`
	Destination string `json:"destination,omitempty"`
	Devel       bool   `json:"devel,omitempty"`
	Home        string `json:"home,omitempty"`
	KeyFile     string `json:"key_file,omitempty"`
	Keyring     string `json:"keyring,omitempty"`
	Password    string `json:"password,omitempty"`
	Prov        bool   `json:"prov,omitempty"`
	Repo        string `json:"repo,omitempty"`
	UntarDir    string `json:"untar_dir,omitempty"`
	Username    string `json:"username,omitempty"`
	Verify      bool   `json:"verify,omitempty"`
	Version     string `json:"version,omitempty"`
}

type HelmChartOpts struct {
	HelmFetchOpts `json:"fetch_opts,omitempty"`

	APIVersions []string               `json:"api_versions,omitempty"`
	Chart       string                 `json:"chart,omitempty"`
	Namespace   string                 `json:"namespace,omitempty"`
	Path        string                 `json:"path,omitempty"`
	ReleaseName string                 `json:"release_name,omitempty"`
	Repo        string                 `json:"repo,omitempty"`
	Values      map[string]interface{} `json:"values,omitempty"`
	Version     string                 `json:"version,omitempty"`
}

// helmTemplate performs Helm fetch/pull + template operations and returns the resulting YAML manifest based on the
// provided chart options.
func helmTemplate(opts HelmChartOpts) (string, error) {
	tempDir, err := ioutil.TempDir("", "helm")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir)

	chart := &chart{
		opts:     opts,
		chartDir: tempDir,
	}

	// If the 'home' option is specified, set the HELM_HOME env var for the duration of the invoke and then reset it
	// to its previous state.
	if chart.opts.Home != "" {
		if helmHome, ok := os.LookupEnv("HELM_HOME"); ok {
			chart.helmHome = &helmHome
		}
		err := os.Setenv("HELM_HOME", chart.opts.Home)
		if err != nil {
			return "", pkgerrors.Wrap(err, "failed to set HELM_HOME")
		}
		defer func() {
			if chart.helmHome != nil {
				_ = os.Setenv("HELM_HOME", *chart.helmHome)
			} else {
				_ = os.Unsetenv("HELM_HOME")
			}
		}()
	}

	// If Path is set, use a local Chart, otherwise fetch from a remote.
	if len(chart.opts.Path) > 0 {
		chart.chartDir = chart.opts.Path
	} else {
		err = chart.fetch()
		if err != nil {
			return "", err
		}
	}

	result, err := chart.template()
	if err != nil {
		return "", err
	}

	return result, nil
}

type chart struct {
	opts     HelmChartOpts
	chartDir string
	helmHome *string // Previous setting of HELM_HOME env var (if any)
}

// fetch runs the `helm fetch` action to fetch a Chart from a remote URL.
func (c *chart) fetch() error {
	p := action.NewPull()
	p.Settings = cli.New()
	p.CaFile = c.opts.CAFile
	p.CertFile = c.opts.CertFile
	p.DestDir = c.chartDir
	//p.DestDir = c.opts.Destination // TODO: Not currently used, but might be useful for caching
	p.KeyFile = c.opts.KeyFile
	p.Keyring = c.opts.Keyring
	p.Password = c.opts.Password
	// c.opts.Prov is unused
	p.RepoURL = c.opts.HelmFetchOpts.Repo
	p.Untar = true
	p.UntarDir = c.chartDir
	p.Username = c.opts.Username
	p.Verify = c.opts.Verify

	if len(c.opts.Repo) > 0 && strings.HasPrefix(c.opts.Repo, "http") {
		return pkgerrors.New("'repo' option specifies the name of the Helm Chart repo, not the URL." +
			"Use 'fetchOpts.repo' to specify a URL for a remote Chart")
	}

	// TODO: We have two different version parameters, but it doesn't make sense
	// 		 to specify both. We should deprecate the FetchOpts one.

	if len(c.opts.Version) == 0 && len(c.opts.HelmFetchOpts.Version) == 0 {
		if c.opts.Devel {
			p.Version = ">0.0.0-0"
		}
	} else if len(c.opts.Version) > 0 {
		p.Version = c.opts.Version
	} else if len(c.opts.HelmFetchOpts.Version) > 0 {
		p.Version = c.opts.HelmFetchOpts.Version
	} // If both are set, prefer the top-level version over the FetchOpts version.

	chartRef := func() string {
		if len(c.opts.Repo) > 0 {
			return fmt.Sprintf("%s/%s", strings.TrimSuffix(c.opts.Repo, "/"), c.opts.Chart)
		}

		return c.opts.Chart
	}

	_, err := p.Run(chartRef())
	if err != nil {
		return pkgerrors.Wrap(err, "failed to pull chart")
	}
	return nil
}

// template runs the `helm template` action to produce YAML from the Chart configuration.
func (c *chart) template() (string, error) {
	cfg := &action.Configuration{
		Capabilities: chartutil.DefaultCapabilities,
		Releases:     storage.Init(driver.NewMemory()),
	}
	if len(c.opts.APIVersions) > 0 {
		cfg.Capabilities.APIVersions = append(cfg.Capabilities.APIVersions, c.opts.APIVersions...)
	}

	// If the namespace isn't set, explicitly set it to "default".
	if len(c.opts.Namespace) == 0 {
		c.opts.Namespace = "default" // nolint: goconst
	}

	installAction := action.NewInstall(cfg)
	installAction.APIVersions = c.opts.APIVersions
	installAction.ClientOnly = true
	installAction.DryRun = true
	installAction.IncludeCRDs = true // TODO: handle this conditionally?
	installAction.Namespace = c.opts.Namespace
	installAction.NameTemplate = c.opts.ReleaseName
	installAction.ReleaseName = c.opts.ReleaseName
	installAction.Version = c.opts.Version

	chartName := func() string {
		// Check if the chart value is a URL with a defined scheme.
		if _url, err := url.Parse(c.opts.Chart); err == nil && len(_url.Scheme) > 0 {
			// Chart path will be of the form `/name-version.tgz`
			re := regexp.MustCompile(`/(\w+)-(\S+)\.tgz$`)
			matches := re.FindStringSubmatch(_url.Path)
			if len(matches) > 1 {
				return matches[1]
			}
		}

		splits := strings.Split(c.opts.Chart, "/")
		if len(splits) == 2 {
			return splits[1]
		}
		return c.opts.Chart
	}

	chart, err := loader.Load(filepath.Join(c.chartDir, chartName()))
	if err != nil {
		return "", pkgerrors.Wrap(err, "failed to load chart from temp directory")
	}

	rel, err := installAction.Run(chart, c.opts.Values)
	if err != nil {
		return "", pkgerrors.Wrap(err, "failed to create chart from template")
	}
	manifests := strings.Builder{}
	manifests.WriteString(rel.Manifest)
	for _, hook := range rel.Hooks {
		switch {
		case testAnnotation.MatchString(hook.Manifest):
			// TODO: allow user to opt out via flag
			// Skip test hook.
		default:
			manifests.WriteString("\n---\n")
			manifests.WriteString(hook.Manifest)
		}
	}

	return manifests.String(), nil
}
