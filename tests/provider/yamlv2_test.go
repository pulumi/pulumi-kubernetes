package provider

import (
	"testing"

	"github.com/pulumi/providertest/pulumitest"
)

// TestHelmUnknowns tests the handling of unknowns in the Helm provider.
// Test steps:
// 1. Preview a program that has computed inputs; expected computed outputs.
// 2. Deploy a program; expected real outputs.
// 3. Preview an update involving a change to the release name; expect replacement.
func TestYamlV2(t *testing.T) {
	// Copy test_dir to temp directory, install deps and create "my-stack"
	test := pulumitest.NewPulumiTest(t, "yamlv2")
	t.Logf("into %s", test.Source())
	t.Cleanup(func() {
		test.Destroy()
	})

	test.Preview()
	test.Up()
}
