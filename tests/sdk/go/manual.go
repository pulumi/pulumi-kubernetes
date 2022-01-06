// Copyright 2016-2021, Pulumi Corporation.
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

	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/provider"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func getRandomNamespace(prefix string) string {
	genRand := func(n int) string {
		const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
		b := make([]byte, n)
		for i := range b {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		}
		return string(b)
	}
	return prefix + "-" + genRand(7)
}

func createRelease(releaseName, releaseNamespace, baseDir string, createNamespace bool) error {
	chartPath := filepath.Join(baseDir, "./nginx")
	chart, err := loader.Load(chartPath)
	if err != nil {
		panic(err)
	}

	var overrides clientcmd.ConfigOverrides
	overrides.Context.Namespace = releaseNamespace

	actionConfig := new(action.Configuration)
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &overrides)
	restConfig, err := kubeconfig.ClientConfig()
	if err != nil {
		return err
	}
	provider.NewKubeConfig(restConfig, kubeconfig)
	if err := actionConfig.Init(provider.NewKubeConfig(restConfig, kubeconfig), releaseNamespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
	}); err != nil {
		panic(err)
	}

	action := action.NewInstall(actionConfig)
	action.Namespace = releaseNamespace
	action.ReleaseName = releaseName
	action.CreateNamespace = createNamespace
	rel, err := action.Run(chart, map[string]interface{}{"service": map[string]interface{}{"type": "ClusterIP"}})
	if err != nil {
		return err
	}
	fmt.Println("Successfully installed release: ", rel.Name)
	return nil
}

func deleteRelease(releaseName, releaseNamespace string) error {
	actionConfig := new(action.Configuration)
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	overrides := clientcmd.ConfigOverrides{Context: api.Context{Namespace: releaseNamespace}}
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &overrides)
	restConfig, err := kubeconfig.ClientConfig()
	if err != nil {
		return err
	}
	if err := actionConfig.Init(provider.NewKubeConfig(restConfig, kubeconfig), releaseNamespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Sprintf(format, v)
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
