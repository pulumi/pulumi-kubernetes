package test

import (
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/stretchr/testify/assert"
)

// TestYamlV2 deploys a complex stack using yaml/v2 package.
// This test has the following features:
// - uses computed inputs
// - leverages DependsOn between components
// - installs and uses CRDs across components
// - uses implicit and explicit dependencies
func TestYamlV2(t *testing.T) {
	test := pulumitest.NewPulumiTest(t, "testdata/yamlv2", opttest.SkipInstall())
	t.Logf("into %s", test.Source())
	t.Cleanup(func() {
		test.Destroy(t)
	})
	test.Preview(t)
	test.Up(t)
}

// TestJobUnreachable ensures that a panic does not occur when diffing Job resources against an unreachable API server.
// https://github.com/pulumi/pulumi-kubernetes/issues/3022
func TestJobUnreachable(t *testing.T) {
	test := pulumitest.NewPulumiTest(t, "testdata/job-unreachable", opttest.SkipInstall())
	t.Logf("into %s", test.Source())
	t.Cleanup(func() {
		test.Destroy(t)
	})
	test.Preview(t)

	// Create the job, but expect it to fail as the job is meant to fail.
	_, err := test.CurrentStack().Up(test.Context())
	assert.ErrorContains(t, err, `but the Kubernetes API server reported that it failed to fully initialize or become live`)

	// Re-run the Pulumi program with a malformed kubeconfig to simulate an unreachable API server.
	// This should not panic annd preview should succeed.
	test.UpdateSource(t, "testdata/job-unreachable/step2")
	test.Preview(t)
}
