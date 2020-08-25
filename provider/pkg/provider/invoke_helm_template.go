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
	"os"
	"path/filepath"
	"strings"

	pkgerrors "github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
)

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
	opts      HelmChartOpts
	chartDir  string
	chartName string
}

// fetch runs the `helm fetch` action to fetch a Chart from a remote URL.
func (c *chart) fetch() error {
	p := action.NewPull()
	p.Settings = cli.New()
	p.CaFile = c.opts.CAFile
	p.CertFile = c.opts.CertFile
	p.DestDir = c.chartDir
	//p.DestDir = c.opts.Destination // TODO: Not currently used, but might be useful for caching
	// c.opts.Home is unused
	p.KeyFile = c.opts.KeyFile
	p.Keyring = c.opts.Keyring
	p.Password = c.opts.Password
	// c.opts.Prov is unused
	p.RepoURL = c.opts.HelmFetchOpts.Repo
	p.Untar = true
	p.UntarDir = c.chartDir
	p.Username = c.opts.Username
	p.Verify = c.opts.Verify

	if c.opts.Version == "" && c.opts.Devel {
		p.Version = ">0.0.0-0"
	} else {
		p.Version = c.opts.HelmFetchOpts.Version
	}

	if c.opts.HelmFetchOpts.Repo == "" {
		splits := strings.Split(c.opts.Chart, "/")
		if len(splits) != 2 {
			return fmt.Errorf("chart repo not specified: %s", c.opts.Chart)
		}
		c.chartName = splits[1]
	} else {
		c.chartName = c.opts.Chart
	}

	_, err := p.Run(c.opts.Chart)
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

	if c.opts.APIVersions != nil {
		cfg.Capabilities.APIVersions = c.opts.APIVersions
		// TODO: add support for overriding kube version
	}

	installAction := action.NewInstall(cfg)
	installAction.ClientOnly = true
	installAction.DryRun = true
	installAction.IncludeCRDs = true // TODO: handle this conditionally?
	installAction.Namespace = c.opts.Namespace
	installAction.ReleaseName = c.opts.ReleaseName
	installAction.Version = c.opts.Version

	chart, err := loader.Load(filepath.Join(c.chartDir, c.chartName))
	if err != nil {
		return "", pkgerrors.Wrap(err, "failed to load chart from temp directory")
	}

	rel, err := installAction.Run(chart, c.opts.Values)
	if err != nil {
		return "", pkgerrors.Wrap(err, "failed to create chart from template")
	}

	return rel.Manifest, nil
}
