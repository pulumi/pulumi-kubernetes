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

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

// TestCSAToSSANoManagedFields tests that the conversion from CSA to SSA works as expected when the object
// being targeted does not have any managedFields entries. This occurs when the object was created before
// beta.2 of SSA (Kubernetes < 1.18). In this case, a conflict with the default `before-first-apply` occurs
// on the first apply with SSA.
// Note: interestingly this issue is not observed with ConfigMaps.
func TestCSAToSSANoManagedFields(t *testing.T) {
	ctx := context.Background()

	// 1. Create the deployment using pulumi-kubernetes in CSA mode.
	test := pulumitest.NewPulumiTest(t, "testdata/csa-to-ssa", opttest.SkipInstall())
	t.Logf("into %s", test.WorkingDir())
	t.Cleanup(func() {
		test.Destroy(t)
	})
	test.Preview(t)
	test.Up(t)

	outputs, err := test.CurrentStack().Outputs(ctx)
	require.NoError(t, err, "Failed to get outputs from stack")
	namespace, ok := outputs["namespace"].Value.(string)
	require.True(t, ok, "Failed to get namespace output as string")
	require.NotEmpty(t, namespace, "Namespace output is empty")
	depName, ok := outputs["deployment"].Value.(string)
	require.True(t, ok, "Failed to get deployment name output as string")
	require.NotEmpty(t, depName, "Deployment name output is empty")

	// 2. We need to nuke the .metadata.managedFields to simulate SSA takeover from an old CSA object. This has
	// to be done after the first apply, as the object's lifecycle should be managed by Pulumi, and newer Kubernetes
	// versions automatically populate this field.
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		t.Fatalf("Failed to build kubeconfig: %v", err)
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		t.Fatalf("Failed to create dynamic client: %v", err)
	}

	depClientNamespaced := dynamicClient.Resource(schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}).Namespace(namespace)

	dep, err := depClientNamespaced.Get(ctx, depName, metav1.GetOptions{})
	require.NoError(t, err, "Failed to get deployment to purge managedFields")

	// Remove the managedFields from the object by setting it to an empty array. Deleting the field entirely will not
	// remove it from the object.
	dep.SetManagedFields([]metav1.ManagedFieldsEntry{})

	// The update needs to be a PUT request, otherwise the server will just repopulate the managedFields.
	dep, err = depClientNamespaced.Update(ctx, dep, metav1.UpdateOptions{})
	require.NoError(t, err, "Failed to update deployment to purge managedFields")
	require.Empty(t, dep.GetManagedFields(), "Failed to remove managedFields from deployment object")

	// 3. Apply step 2 of testdata where SSA mode is enabled, with a change in the image spec field.
	test.UpdateSource(t, "testdata/csa-to-ssa/step2")
	test.Preview(t)
	test.Up(t)
	test.Up(t, optup.ExpectNoChanges())
}
