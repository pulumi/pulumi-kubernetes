// Copyright 2016-2023, Pulumi Corporation.
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

// nolint:goconst
package test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi-kubernetes/tests/v4"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// TestPreview tests the `pulumi preview` CUJ with a serviceaccount that is not allowed to create a configmap. This test ensures that
// we do a client side diff if the server dry run fails and do not error on preview. We also ensure that we error on pulumi up, and that
// the configmap was not created using kubectl.
func TestPreview(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  "preview-auth",
		ExpectRefreshChanges: true,
		// Enable destroy-on-cleanup so we can shell out to kubectl to make external changes to the resource and reuse the same stack.
		DestroyOnCleanup: true,
		Quick:            true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
	})

	// Create service account and RBAC policies for the service account.
	out, err := tests.Kubectl("apply -f preview-auth/service-account.yaml")
	if err != nil {
		t.Fatalf("unable to create RBAC policies: %s, out: %s", err, string(out))
	}
	t.Cleanup(func() {
		log.Println("Deleting service-account and rbac")
		_, err = tests.Kubectl("delete -f preview-auth/service-account.yaml")
		assert.NoError(t, err)
	})

	// Create kubeconfig for service account.
	kubeconfigPath, err := createSAKubeconfig(t, "ci-robot")
	if err != nil {
		t.Fatalf("unable to create kubeconfig: %s", err)
	}
	// Set kubeconfig env var for the test to use.
	test = test.With(integration.ProgramTestOptions{
		Env: []string{
			"KUBECONFIG=" + kubeconfigPath,
		},
	})

	// Initialize and the test project.
	pt := integration.ProgramTestManualLifeCycle(t, &test)
	err = pt.TestLifeCyclePrepare()
	if err != nil {
		t.Fatalf("unable to create temp dir: %s", err)
	}
	t.Cleanup(pt.TestCleanUp)
	err = pt.TestLifeCycleInitialize()
	if err != nil {
		t.Fatalf("unable to init test project: %s", err)
	}
	t.Cleanup(func() {
		destroyErr := pt.TestLifeCycleDestroy()
		assert.NoError(t, destroyErr)
	})

	// Run a preview and assert no error, since we do a client side diff if server dry run fails.
	err = pt.RunPulumiCommand("preview", "--non-interactive", "--diff", "--refresh", "--show-config")
	assert.NoError(t, err)

	// Run pulumi up and assert error since our SA doesn't have permissions to create a configmap.
	err = pt.TestPreviewUpdateAndEdits()
	assert.Error(t, err)

	// Check that the configmap was not created using kubectl.
	out, err = tests.Kubectl("get configmap foo")
	assert.Error(t, err)
	assert.Contains(t, string(out), `Error from server (NotFound): configmaps "foo" not found`)
}

// createSAKubeconfig creates a modified kubeconfig for the service account. The kubeconfig is created in a
// tmp dir which is cleaned up after the test.
func createSAKubeconfig(t *testing.T, saName string) (string, error) {
	t.Helper()

	// Create token to use for the service account.
	token, err := tests.Kubectl(fmt.Sprintf("create token %s --duration=1h", saName))
	if err != nil {
		return "", err
	}

	// Load default kubeconfig as base.
	config, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return "", err
	}

	// Use current-context cluster as base for service account context/auth.
	config.AuthInfos[saName] = &clientcmdapi.AuthInfo{Token: string(token)}
	config.Contexts[saName] = &clientcmdapi.Context{
		Cluster:  config.Contexts[config.CurrentContext].Cluster,
		AuthInfo: saName,
	}
	config.CurrentContext = saName

	// Create tmp dir to store kubeconfig.
	tmpDir, err := os.MkdirTemp("", "kubeconfig-preview")
	if err != nil {
		return "", err
	}

	t.Cleanup(func() {
		log.Println("Deleting kubeconfig tmp dir")
		assert.NoError(t, os.RemoveAll(tmpDir))
	})

	// Write kubeconfig to tmp dir.
	kubeconfigPath := filepath.Join(tmpDir, "kubeconfig.txt")
	err = clientcmd.WriteToFile(*config, kubeconfigPath)

	return kubeconfigPath, err
}

// TestPreviewWithApply tests the `pulumi preview` CUJ where the user Pulumi program contains an Apply call on status subresoruces.
// This is to ensure we don't fail preview, since status fields are only populated after the resource is created on cluster.
func TestPreviewWithApply(t *testing.T) {
	var externalIP, nsName, svcName string
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  "preview-apply",
		ExpectRefreshChanges: false,
		// Enable destroy-on-cleanup so we can shell out to kubectl to make external changes to the resource and reuse the same stack.
		DestroyOnCleanup: true,
		Quick:            true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			var ok bool
			externalIP, ok = stackInfo.Outputs["ip"].(string)
			require.True(t, ok)
			nsName, ok = stackInfo.Outputs["nsName"].(string)
			require.True(t, ok)
			svcName, ok = stackInfo.Outputs["svcName"].(string)
			require.True(t, ok)
		},
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
	})

	// Initialize and the test project.
	pt := integration.ProgramTestManualLifeCycle(t, &test)
	err := pt.TestLifeCyclePrepare()
	if err != nil {
		t.Fatalf("unable to create temp dir: %s", err)
	}
	t.Cleanup(pt.TestCleanUp)

	err = pt.TestLifeCycleInitialize()
	if err != nil {
		t.Fatalf("unable to init test project: %s", err)
	}
	t.Cleanup(func() {
		destroyErr := pt.TestLifeCycleDestroy()
		assert.NoError(t, destroyErr)
	})

	// Run a preview and assert no error.
	err = pt.RunPulumiCommand("preview", "--non-interactive", "--diff", "--refresh", "--show-config")
	assert.NoError(t, err)
	assert.Equal(t, "", externalIP)

	// Run pulumi up and assert no error creating the resources.
	err = pt.TestPreviewUpdateAndEdits()
	require.NoError(t, err)

	// Ensure that the ip output is the same as the external ip of the service via kubectl.
	out, err := tests.Kubectl("get service", svcName, "-n", nsName, "-o jsonpath={.status.loadBalancer.ingress[0].ip}")
	require.NoError(t, err)
	assert.Equal(t, externalIP, string(out))
}
