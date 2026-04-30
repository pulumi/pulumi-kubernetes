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
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"

	"github.com/pulumi/pulumi-kubernetes/tests/v4"
)

// TestParameterizeE2E runs the full CRD parameterization round trip against a
// real cluster: apply a CRD, generate a parameterized SDK from it, run the
// program that uses that SDK, verify the resulting CR exists, destroy.
//
// This is the one test that exercises the provider's runtime gvkFromURN
// dispatch — the other parameterize tests only cover schema shape.
//
// Requires `make k8sprovider` and a reachable cluster; both will fail loudly
// via pulumitest / kubectl if absent.
func TestParameterizeE2E(t *testing.T) {
	tests.SkipIfShort(t, "runs against a real cluster")

	providerBin, err := filepath.Abs("../../../bin/pulumi-resource-kubernetes")
	require.NoError(t, err)

	fixtureDir, err := filepath.Abs("parameterize-e2e")
	require.NoError(t, err)
	crdPath := filepath.Join(fixtureDir, "gateway-crd.yaml")

	// Install the CRD into the cluster up front and schedule its removal.
	// Registering the CRD cleanup BEFORE the stack cleanup matters: t.Cleanup
	// runs LIFO, so the pulumi destroy registered later runs first, then the
	// CRD is removed.
	out, err := tests.Kubectl("apply -f " + crdPath)
	require.NoError(t, err, "kubectl apply CRD: %s", out)
	t.Cleanup(func() {
		_, _ = tests.Kubectl("delete -f " + crdPath + " --wait=false --ignore-not-found")
	})

	test := pulumitest.NewPulumiTest(
		t,
		"parameterize-e2e",
		opttest.LocalProviderPath("kubernetes", "../../../bin"),
		opttest.SkipInstall(),
	)
	workDir := test.WorkingDir()

	// Generate the parameterized SDK into the working dir. The replace in the
	// fixture's go.mod points at ./sdks/go, which is where this lands.
	cmd := exec.Command("pulumi", "package", "gen-sdk", providerBin,
		"--language", "go",
		"--out", filepath.Join(workDir, "sdks"),
		"--local",
		"--", "-v", "1.0.0", "-c", "gateway-crd.yaml")
	cmd.Dir = workDir
	genOut, err := cmd.CombinedOutput()
	require.NoError(t, err, "pulumi package gen-sdk: %s", genOut)

	test.Install(t)
	up := test.Up(t)
	t.Cleanup(func() {
		test.Destroy(t)
	})

	gwName, ok := up.Outputs["gatewayName"].Value.(string)
	require.True(t, ok, "gatewayName output should be a string, got %T", up.Outputs["gatewayName"].Value)
	assert.Equal(t, "e2e-gateway", gwName)

	// Verify via kubectl.
	out, err = tests.Kubectl("get gateway.gateway.pulumi.test e2e-gateway -n default -o name")
	require.NoError(t, err, "kubectl get gateway: %s", out)
	assert.Contains(t, string(out), "gateway.gateway.pulumi.test/e2e-gateway")
}
