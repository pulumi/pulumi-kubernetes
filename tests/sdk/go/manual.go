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
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/client-go/tools/clientcmd"
)

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

func createRelease(releaseName, releaseNamespace, baseDir string, createNamespace bool) error {
	chartPath := filepath.Join(baseDir, "./nginx")
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
