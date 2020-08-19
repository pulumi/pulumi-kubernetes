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
	"io/ioutil"
	"os"
	"path/filepath"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
)

type HelmFetchOpts struct {
	CAFile      string `json:"ca_file"`
	CertFile    string `json:"cert_file"`
	Destination string `json:"destination"`
	Devel       bool   `json:"devel"`
	Home        string `json:"home"`
	KeyFile     string `json:"key_file"`
	Keyring     string `json:"keyring"`
	Password    string `json:"password"`
	Prov        bool   `json:"prov"`
	Repo        string `json:"repo"`
	UntarDir    string `json:"untar_dir"`
	Username    string `json:"username"`
	Verify      bool   `json:"verify"`
	Version     string `json:"version"`
}

type HelmChartOpts struct {
	HelmFetchOpts

	APIVersions []string               `json:"api_versions"`
	Chart       string                 `json:"chart"`
	Namespace   string                 `json:"namespace"`
	Repo        string                 `json:"repo"`
	Path        string                 `json:"path"`
	Values      map[string]interface{} `json:"values"`
	Version     string                 `json:"version"`
}

// HelmTemplate performs Helm fetch/pull + template operations and returns the resulting YAML manifest based on the
// provided chart options.
func HelmTemplate(opts HelmChartOpts) (string, error) {
	tempDir, err := ioutil.TempDir("", "helm")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir)

	chart := &chart{
		opts:     opts,
		chartDir: tempDir,
	}

	err = chart.fetch()
	if err != nil {
		return "", err
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
}

func (c *chart) fetch() error {
	p := action.NewPull()
	p.Settings = cli.New()
	p.RepoURL = c.opts.HelmFetchOpts.Repo
	p.Version = c.opts.Version
	p.Untar = true
	p.UntarDir = c.chartDir

	_, err := p.Run(c.opts.Chart)
	if err != nil {
		return err
	}
	return nil
}

func (c *chart) template() (string, error) {
	cfg := &action.Configuration{
		Capabilities: chartutil.DefaultCapabilities, // TODO: use caps if set in opts
		Releases:     storage.Init(driver.NewMemory()),
	}
	installAction := action.NewInstall(cfg)
	installAction.ClientOnly = true
	installAction.DryRun = true
	installAction.ReleaseName = c.opts.Chart // TODO: is this right?

	chart, err := loader.Load(filepath.Join(c.chartDir, c.opts.Chart))
	if err != nil {
		return "", err
	}

	rel, err := installAction.Run(chart, map[string]interface{}{})
	if err != nil {
		return "", err
	}

	return rel.Manifest, nil
}
