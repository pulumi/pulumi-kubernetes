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
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/pulumi-kubernetes/tests/v4"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJavaCanCreateCRD tests that we can create a CRD using the Java SDK, and that `x-kubernetes-*` fields are
// correctly serialized.
func TestJavaCanCreateCRD(t *testing.T) {
	// Step 1 creates a CRD with `x-kubernetes-preserve-unknown-fields` set to true.
	test := pulumitest.NewPulumiTest(t, "testdata/crd-java/step1")
	t.Logf("into %s", test.WorkingDir())
	t.Cleanup(func() {
		test.Destroy(t)
	})
	test.Preview(t)
	test.Up(t)

	// Step 2 adds a pulumi CRD get operation and ensures we can read its URN properly.
	test.UpdateSource(t, "testdata/crd-java/step2")
	test.Preview(t)
	test.Up(t)
	up := test.Up(t, optup.ExpectNoChanges())

	urn, ok := up.Outputs["urn"]
	require.True(t, ok)
	require.NotNil(t, urn)
	require.Equal(
		t,
		urn.Value,
		"urn:pulumi:test::crd_java::kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::getCRDUrn",
	)

	// Verify with kubectl that the CRD has `x-kubernetes-*` fields set correctly.
	output, err := tests.Kubectl("get crd javacrds.example.com -o json")
	require.NoError(t, err)
	assert.Contains(t, string(output), `"x-kubernetes-preserve-unknown-fields": true`)

	// Step 3 removes the `x-kubernetes-preserve-unknown-fields` field and ensures that the CRD is updated.
	test.UpdateSource(t, "testdata/crd-java/step3")
	test.Preview(t)
	test.Up(t)
	up = test.Up(t, optup.ExpectNoChanges())

	// Verify with kubectl that the CRD no longer has `x-kubernetes-*` fields set.
	output, err = tests.Kubectl("get crd javacrds.example.com -o json")
	require.NoError(t, err)
	assert.NotContains(t, string(output), `"x-kubernetes-preserve-unknown-fields": true`)
}
