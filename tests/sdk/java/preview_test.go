package test

import (
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi-kubernetes/v4/tests"
	"github.com/stretchr/testify/require"
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

// TestCRDPreviews ensures that CRDs are correctly previewed, and are not created or updated on the cluster.
// https://github.com/pulumi/pulumi-kubernetes/issues/3094
func TestCRDPreviews(t *testing.T) {
	const (
		crdName    = "crontabs.previewtest.pulumi.com"
		testFolder = "testdata/crd-previews"
	)

	// 1. Create the CRD resource
	test := pulumitest.NewPulumiTest(t, testFolder, opttest.SkipInstall())
	t.Logf("into %s", test.Source())
	t.Cleanup(func() {
		test.Destroy()
	})
	test.Up()

	// 2. Preview should not actually update the CRD resource. Step 2 adds a new field ("testNewField") to the CRD.
	test.UpdateSource(testFolder, "step2")
	test.Preview()

	out, err := tests.Kubectl("get", "crd", crdName, "-o", "yaml")
	require.NoError(t, err, "unable to get CRD with kubectl")
	require.NotContains(t, string(out), "testNewField", "expected CRD to not have new field added in preview")

	// 3. Update should actually update the CRD resource.
	test.Up()
	out, err = tests.Kubectl("get", "crd", crdName, "-o", "yaml")
	require.NoError(t, err, "unable to get CRD with kubectl")
	require.Contains(t, string(out), "testNewField", "expected CRD to have new field added in update operation")
}
