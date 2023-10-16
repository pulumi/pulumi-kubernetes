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

package test

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/archive"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/client-go/tools/clientcmd"
)

// a Helm chart available in numerous formats for test purposes
type helmChart struct {
	// The chart name (e.g `nginx`)
	Name string
	// a path to a Helm chart directory, relative to `tests/sdk/go`
	TestPath string
	// a path to a helm chart archive, relative to `tests/sdk/go`
	TestArchive string
	// The Helm Repository to use for this chart
	HelmRepo repo.Entry
	// an HTTPS URL to the Helm chart archive
	ChartURL string
	// an OCI URL to the helm chart
	OciURL string
	// versions of this chart for test purposes
	Versions []helmChartVersion
}

func (hc helmChart) ChartReference() string {
	return fmt.Sprintf("%s/%s", hc.HelmRepo.Name, hc.Name)
}

type helmChartVersion struct {
	Version string
	Values  map[string]interface{}
}

var (
	// homepage: https://hub.docker.com/r/bitnamicharts/nginx
	bitnamiNginxChart = helmChart{
		Name:        "nginx",
		TestPath:    "../../testdata/helm/nginx",
		TestArchive: "../../testdata/helm/nginx-15.3.4.tgz",
		ChartURL:    "https://charts.bitnami.com/bitnami/nginx-15.3.4.tgz",
		OciURL:      "oci://registry-1.docker.io/bitnamicharts/nginx",
		HelmRepo: repo.Entry{
			Name: "bitnami",
			URL:  "https://charts.bitnami.com/bitnami",
		},
		Versions: []helmChartVersion{{
			Version: "15.3.4",
			Values: map[string]interface{}{
				"service": map[string]interface{}{
					"type": "ClusterIP",
				},
			},
		}},
	}
)

func (hc helmChart) ExtractTo(dir string) error {
	file, err := os.Open(hc.TestArchive)
	if err != nil {
		return err
	}
	defer file.Close()
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)

	err = archive.ExtractTGZ(tr, dir)
	if err != nil {
		return fmt.Errorf("failed to extract tarball: %w", err)
	}
	return nil
}

func getRandomNamespace(prefix string) string {
	genRand := func(n int) string {
		const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
		b := make([]byte, n)
		for i := range b {
			//nolint: gosec
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		}
		return string(b)
	}
	return prefix + "-" + genRand(7)
}

func namespacedKubeconfig(namespace string) (*provider.KubeConfig, error) {
	var overrides clientcmd.ConfigOverrides
	overrides.Context.Namespace = namespace

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &overrides)
	restConfig, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	return provider.NewKubeConfig(restConfig, kubeconfig), nil
}

func createRelease(releaseName, releaseNamespace, chartPath string, createNamespace bool) error {
	chart, err := loader.Load(chartPath)
	if err != nil {
		panic(err)
	}

	kubeconfig, err := namespacedKubeconfig(releaseNamespace)
	if err != nil {
		panic(err)
	}

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(kubeconfig, releaseNamespace, os.Getenv("HELM_DRIVER"), func(format string, v ...any) {
		_ = fmt.Sprintf(format, v)
	}); err != nil {
		panic(err)
	}

	install := action.NewInstall(actionConfig)
	install.Namespace = releaseNamespace
	install.ReleaseName = releaseName
	install.CreateNamespace = createNamespace
	// Block on helm install since otherwise if we import resources created by the release which are not ready,
	// we might end up with mysterious test failures (initErrors during updates might trigger an update).
	install.Wait = true
	install.Timeout = 5 * time.Minute
	rel, err := install.Run(chart, map[string]any{"service": map[string]any{"type": "ClusterIP"}})
	if err != nil {
		return err
	}
	fmt.Println("Successfully installed release: ", rel.Name)
	return nil
}

func listReleases(releaseNamespace string) ([]*release.Release, error) {
	kubeconfig, err := namespacedKubeconfig(releaseNamespace)
	if err != nil {
		return nil, err
	}
	actionConfig := new(action.Configuration)
	if err != nil {
		return nil, err
	}
	if err := actionConfig.Init(kubeconfig, releaseNamespace, os.Getenv("HELM_DRIVER"), func(format string,
		v ...any) {
		_ = fmt.Sprintf(format, v)
	}); err != nil {
		panic(err)
	}
	list := action.NewList(actionConfig)
	releases, err := list.Run()
	if err != nil {
		return nil, err
	}
	return releases, nil
}

func getRelease(releaseName, releaseNamespace string) (*release.Release, error) {
	kubeconfig, err := namespacedKubeconfig(releaseNamespace)
	if err != nil {
		return nil, err
	}
	actionConfig := new(action.Configuration)
	if err != nil {
		return nil, err
	}
	if err := actionConfig.Init(kubeconfig, releaseNamespace, os.Getenv("HELM_DRIVER"), func(format string,
		v ...any) {
		_ = fmt.Sprintf(format, v)
	}); err != nil {
		panic(err)
	}
	get := action.NewGet(actionConfig)
	release, err := get.Run(releaseName)
	if err != nil {
		return nil, err
	}
	return release, nil
}

func deleteRelease(releaseName, releaseNamespace string) error {
	kubeconfig, err := namespacedKubeconfig(releaseNamespace)
	if err != nil {
		panic(err)
	}

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(kubeconfig, releaseNamespace, os.Getenv("HELM_DRIVER"), func(format string, v ...any) {
		_ = fmt.Sprintf(format, v)
	}); err != nil {
		panic(err)
	}
	get := action.NewGet(actionConfig)
	_, err = get.Run(releaseName)
	if err != nil && strings.Contains(err.Error(), "release: not found") {
		return nil
	}
	if err != nil {
		return err
	}
	act := action.NewUninstall(actionConfig)
	act.Wait = true
	if _, err = act.Run(releaseName); err != nil {
		return err
	}
	fmt.Printf("Uninstall release: %s/%s\n", releaseNamespace, releaseName)
	return nil
}

type helmEnvironment struct {
	CacheHome  string
	ConfigHome string
	DataHome   string
}

func (he *helmEnvironment) EnvVars() []string {
	return []string{
		helmpath.CacheHomeEnvVar + "=" + he.CacheHome,
		helmpath.ConfigHomeEnvVar + "=" + he.ConfigHome,
		helmpath.DataHomeEnvVar + "=" + he.DataHome,
	}
}

var helmEnvLock sync.Mutex

// createHelmEnvironment creates an isolated Helm environment with some preconfigured repositories.
func createHelmEnvironment(t *testing.T, re ...repo.Entry) (he *helmEnvironment, cleanup func() error, err error) {
	tmpdir, err := os.MkdirTemp("", "helm-")
	if err != nil {
		return nil, nil, err
	}
	t.Logf("Initializing Helm environment (%s)", tmpdir)

	he = &helmEnvironment{
		CacheHome:  filepath.Join(tmpdir, "cache"),
		ConfigHome: filepath.Join(tmpdir, "config"),
		DataHome:   filepath.Join(tmpdir, "data"),
	}
	_ = os.MkdirAll(he.CacheHome, 0755)
	_ = os.MkdirAll(he.ConfigHome, 0755)
	_ = os.MkdirAll(he.DataHome, 0755)

	settings := func() *cli.EnvSettings {
		// magic up a Helm EnvSettings struct using the home directories we just created
		helmEnvLock.Lock()
		defer helmEnvLock.Unlock()
		setEnv := func(name, new string) func() {
			old, existed := os.LookupEnv(name)
			_ = os.Setenv(name, new)
			return func() {
				if existed {
					_ = os.Setenv(name, old)
				} else {
					_ = os.Unsetenv(name)
				}
			}
		}
		unsetCacheHome := setEnv(helmpath.CacheHomeEnvVar, he.CacheHome)
		defer unsetCacheHome()
		unsetConfigHome := setEnv(helmpath.ConfigHomeEnvVar, he.ConfigHome)
		defer unsetConfigHome()
		unsetDataHome := setEnv(helmpath.DataHomeEnvVar, he.DataHome)
		defer unsetDataHome()
		return cli.New()
	}()

	// Generate repositories.yaml as they do in `helm repo add`
	rf := repo.NewFile()
	for _, c := range re {
		c := c
		r, err := repo.NewChartRepository(&c, getter.All(settings))
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to add chart repository %q", c.Name)
		}
		r.CachePath = settings.RepositoryCache
		if _, err := r.DownloadIndexFile(); err != nil {
			return nil, nil, errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", c.URL)
		}
		rf.Add(&c)
	}
	if err = rf.WriteFile(settings.RepositoryConfig, 0600); err != nil {
		return nil, nil, err
	}

	cleanup = func() error {
		if t.Failed() {
			t.Logf("Skipping cleanup of Helm environment due to test failure (%s)", tmpdir)
			return nil
		}
		t.Logf("Cleaning up Helm environment (%s)", tmpdir)
		return os.RemoveAll(tmpdir)
	}
	return he, cleanup, nil
}
