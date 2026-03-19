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

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
)

// TestEnablePatchForceConfig tests that the enablePatchForce provider config option
// resolves SSA field conflicts. It does this by:
// 1. Creating a ConfigMap with SSA.
// 2. Adding a conflicting field manager via kubectl to simulate a conflict.
// 3. Updating the ConfigMap with enablePatchForce=true, which should succeed despite the conflict.
func TestEnablePatchForceConfig(t *testing.T) {
	ctx := context.Background()

	// Step 1: Create the ConfigMap using SSA.
	test := pulumitest.NewPulumiTest(t, "testdata/patch-force", opttest.SkipInstall())
	t.Logf("into %s", test.WorkingDir())
	t.Cleanup(func() {
		test.Destroy(t)
	})
	test.Preview(t)
	test.Up(t)

	outputs, err := test.CurrentStack().Outputs(ctx)
	require.NoError(t, err, "failed to get outputs")
	namespace, ok := outputs["namespace"].Value.(string)
	require.True(t, ok, "failed to get namespace output")
	require.NotEmpty(t, namespace)
	cmName, ok := outputs["configmap"].Value.(string)
	require.True(t, ok, "failed to get configmap name output")
	require.NotEmpty(t, cmName)

	// Step 2: Simulate a field conflict by patching the ConfigMap with a different field manager.
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	require.NoError(t, err, "failed to build kubeconfig")
	dynamicClient, err := dynamic.NewForConfig(config)
	require.NoError(t, err, "failed to create dynamic client")

	cmClient := dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "configmaps",
	}).Namespace(namespace)

	// Apply a conflicting field manager that owns the "key1" field.
	patch := []byte(`{"apiVersion":"v1","kind":"ConfigMap","data":{"key1":"conflicting-value"}}`)
	_, err = cmClient.Patch(ctx, cmName, types.ApplyPatchType, patch, metav1.PatchOptions{
		FieldManager: "conflict-manager",
		Force:        boolPtr(true),
	})
	require.NoError(t, err, "failed to apply conflicting field manager")

	// Step 3: Update with enablePatchForce=true to resolve the conflict.
	test.UpdateSource(t, "testdata/patch-force/step2")
	test.Preview(t)
	test.Up(t)

	// Verify the update succeeded and no further changes are needed.
	test.Up(t, optup.ExpectNoChanges())
}

func boolPtr(b bool) *bool {
	return &b
}
