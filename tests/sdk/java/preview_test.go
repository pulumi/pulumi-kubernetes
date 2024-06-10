package test

import (
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
)

// TestPreviewReplacements ensures that replacements for immutable fields are correctly previewed.
func TestPreviewReplacements(t *testing.T) {
	test := pulumitest.NewPulumiTest(t, "testdata/preview-replacements", opttest.SkipInstall())
	t.Logf("into %s", test.Source())
	t.Cleanup(func() {
		test.Destroy()
	})
	test.Preview()
	test.Up()

	// Preview should not fail when there is a replacement due to immutable fields.
	test.UpdateSource("testdata/preview-replacements", "step2")
	test.Preview()
}
