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
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"

	"github.com/pulumi/pulumi-kubernetes/tests/v4"
)

// setupPreexistingConfigMap creates a fresh namespace containing a plain,
// non-Helm-managed ConfigMap named "adopt-me". This is the resource that the
// take-ownership chart also renders, so it lets us exercise Helm's adoption of
// a resource it does not yet own. The namespace (and everything in it) is
// removed on test cleanup.
func setupPreexistingConfigMap(t *testing.T) string {
	t.Helper()

	ns := fmt.Sprintf("take-ownership-%d", time.Now().UnixNano())
	out, err := tests.Kubectl("create", "namespace", ns)
	require.NoError(t, err, "create namespace: %s", out)
	t.Cleanup(func() {
		_, _ = tests.Kubectl("delete", "namespace", ns, "--wait=false")
	})

	out, err = tests.Kubectl("create", "configmap", "adopt-me", "-n", ns, "--from-literal=key=preexisting")
	require.NoError(t, err, "create configmap: %s", out)

	return ns
}

// TestHelmReleaseTakeOwnershipAdopts verifies that a Release with
// takeOwnership=true adopts a pre-existing, unmanaged resource instead of
// failing with an ownership conflict, stamping it with Helm's ownership
// metadata.
func TestHelmReleaseTakeOwnershipAdopts(t *testing.T) {
	ns := setupPreexistingConfigMap(t)

	test := pulumitest.NewPulumiTest(t, "testdata/helm-release-take-ownership", opttest.SkipInstall())
	t.Logf("into %s", test.WorkingDir())
	test.SetConfig(t, "namespace", ns)
	test.SetConfig(t, "takeOwnership", "true")
	t.Cleanup(func() {
		test.Destroy(t)
	})

	// With takeOwnership enabled, the install adopts the existing ConfigMap
	// and succeeds.
	test.Up(t)

	// Helm should have stamped its ownership metadata onto the adopted object.
	managedBy, err := tests.Kubectl("get", "configmap", "adopt-me", "-n", ns,
		"-o", `jsonpath={.metadata.labels.app\.kubernetes\.io/managed-by}`)
	require.NoError(t, err)
	assert.Equal(t, "Helm", strings.TrimSpace(string(managedBy)),
		"adopted ConfigMap should carry the Helm managed-by label")

	releaseName, err := tests.Kubectl("get", "configmap", "adopt-me", "-n", ns,
		"-o", `jsonpath={.metadata.annotations.meta\.helm\.sh/release-name}`)
	require.NoError(t, err)
	assert.Equal(t, "adopt-test", strings.TrimSpace(string(releaseName)),
		"adopted ConfigMap should be annotated with the release name")

	// The pre-existing ConfigMap was created with key=preexisting; after
	// adoption Helm should have applied the chart's own value, proving it
	// actually manages the resource's contents and not just its metadata.
	value, err := tests.Kubectl("get", "configmap", "adopt-me", "-n", ns,
		"-o", `jsonpath={.data.key}`)
	require.NoError(t, err)
	assert.Equal(t, "from-helm", strings.TrimSpace(string(value)),
		"Helm should have written the chart's value into the adopted ConfigMap")
}

// TestHelmReleaseTakeOwnershipDisabledConflicts is the negative control: with
// takeOwnership unset the same install fails with Helm's resource-conflict
// error, proving the flag is what enables adoption.
func TestHelmReleaseTakeOwnershipDisabledConflicts(t *testing.T) {
	ctx := context.Background()
	ns := setupPreexistingConfigMap(t)

	test := pulumitest.NewPulumiTest(t, "testdata/helm-release-take-ownership", opttest.SkipInstall())
	t.Logf("into %s", test.WorkingDir())
	test.SetConfig(t, "namespace", ns)
	test.SetConfig(t, "takeOwnership", "false")
	t.Cleanup(func() {
		test.Destroy(t)
	})

	_, upErr := test.CurrentStack().Up(ctx)
	require.Error(t, upErr, "expected install to fail without takeOwnership")
	assert.Contains(t, upErr.Error(), "cannot be imported into the current release",
		"error should be Helm's ownership-conflict message")
}
