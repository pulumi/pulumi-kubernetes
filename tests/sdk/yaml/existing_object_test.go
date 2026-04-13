// Copyright 2016-2026, Pulumi Corporation.
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"

	"github.com/pulumi/pulumi-kubernetes/tests/v4"
)

// TestExistingObjectBlocksCreate tests that creating a resource that
// already exists in the cluster fails with AlreadyExists. The "default"
// namespace is a built-in that always exists, so no setup is needed.
func TestExistingObjectBlocksCreate(t *testing.T) {
	ctx := context.Background()

	test := pulumitest.NewPulumiTest(t, "testdata/existing-object", opttest.SkipInstall())
	t.Logf("into %s", test.WorkingDir())

	_, upErr := test.CurrentStack().Up(ctx)
	require.Error(t, upErr, "expected Up to fail when object already exists")
	assert.Contains(t, upErr.Error(), "already exists",
		"error should mention that the resource already exists")
}

// TestCSACreateToSSAUpdate tests that resources created with client-side
// apply (the new default for creates) can be successfully updated with
// server-side apply on subsequent runs. The field manager migration
// (fixCSAFieldManagers) handles the transition transparently.
func TestCSACreateToSSAUpdate(t *testing.T) {
	// Step 1: Create a ConfigMap with an explicit name. With server-side
	// apply enabled, the create internally uses client-side apply.
	test := pulumitest.NewPulumiTest(t, "testdata/csa-create-ssa-update", opttest.SkipInstall())
	t.Logf("into %s", test.WorkingDir())
	t.Cleanup(func() {
		test.Destroy(t)
	})
	result := test.Up(t)
	namespace := result.Outputs["namespace"].Value.(string)

	// Verify step 1 data.
	outB, err := tests.Kubectl("get", "configmap", "csa-ssa-transition-test",
		"-n", namespace, "-o", "jsonpath={.data.key}")
	require.NoError(t, err)
	assert.Equal(t, "step1", string(outB))

	// Verify that re-running with no changes produces no diff. This
	// confirms the client-side create to server-side apply update
	// transition doesn't cause spurious diffs from field manager
	// differences.
	test.Up(t, optup.ExpectNoChanges())

	// Step 2: Update the ConfigMap data. This goes through the SSA
	// update path, which triggers fixCSAFieldManagers to transfer
	// field ownership from the CSA field manager to SSA.
	test.UpdateSource(t, "testdata/csa-create-ssa-update/step2")
	test.Up(t)

	// Verify step 2 data.
	outB, err = tests.Kubectl("get", "configmap", "csa-ssa-transition-test",
		"-n", namespace, "-o", "jsonpath={.data.key}")
	require.NoError(t, err)
	assert.Equal(t, "step2", string(outB))
}
