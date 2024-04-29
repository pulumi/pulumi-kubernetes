package test

import (
	"testing"

	"github.com/pulumi/providertest/pulumitest"
)

// TestYamlV2 deploys a complex stack using yaml/v2 package.
// This test has the following features:
// - uses computed inputs
// - leverages DependsOn between components
// - installs and uses CRDs across components
// - uses implicit and explicit dependencies
func TestYamlV2(t *testing.T) {
	test := pulumitest.NewPulumiTest(t, "testdata/yamlv2")
	t.Logf("into %s", test.Source())
	t.Cleanup(func() {
		test.Destroy()
	})
	test.Preview()
	test.Up()
}
