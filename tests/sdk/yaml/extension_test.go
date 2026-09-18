package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"

	"github.com/pulumi/pulumi-kubernetes/tests/v4"
)

func packageAddCmd(t *testing.T, test *pulumitest.PulumiTest) func(args ...string) {
	t.Helper()

	cwd, err := os.Getwd()
	require.NoError(t, err)
	providerBin := filepath.Join(cwd, "..", "..", "..", "bin", "pulumi-resource-kubernetes")

	pulumi := test.CurrentStack().Workspace().PulumiCommand()
	return func(args ...string) {
		t.Helper()
		_, stderr, _, err := pulumi.Run(test.Context(), test.WorkingDir(),
			nil, nil, nil, nil, append([]string{"package", "add", providerBin}, args...)...)
		require.NoErrorf(t, err, "pulumi package add %v failed: %s", args, stderr)
	}
}

func TestExtensionGatewayAPI(t *testing.T) {
	test := pulumitest.NewPulumiTest(t, "testdata/extension-gateway-api", opttest.SkipInstall())
	t.Cleanup(func() {
		test.Destroy(t)
	})

	packageAdd := packageAddCmd(t, test)

	packageAdd("--extension", "name=gateway-networking crd-manifest=gateway-api-crds.yaml")

	preview := test.Preview(t)
	require.Contains(t, preview.StdOut, "GatewayClass",
		"preview should plan the extension-served GatewayClass")

	up := test.Up(t)
	require.Equal(t, "example-class", up.Outputs["gatewayClassName"].Value,
		"extension-served GatewayClass should be created with its declared name")
	require.Equal(t, "gateway-system", up.Outputs["namespaceName"].Value,
		"base-provider Namespace should be created alongside the extension resource")
}

func TestExtensionHyphenatedPropertyNames(t *testing.T) {
	testFolder := "testdata/extension-hyphenated-properties"
	crdManifest := filepath.Join(testFolder, "hyphenated-crd.yaml")

	_, err := tests.Kubectl("apply", "-f", crdManifest)
	require.NoError(t, err, "failed to apply the CRD")
	t.Cleanup(func() {
		_, err := tests.Kubectl("delete", "-f", crdManifest)
		contract.AssertNoErrorf(err, "failed to delete the CRD during cleanup")
	})

	test := pulumitest.NewPulumiTest(t, testFolder, opttest.SkipInstall())
	t.Cleanup(func() {
		test.Destroy(t)
	})

	packageAdd := packageAddCmd(t, test)
	packageAdd("--extension", "name=hyphenprops crd-manifest=hyphenated-crd.yaml")

	up := test.Up(t)

	spec, ok := up.Outputs["widgetSpec"].Value.(map[string]any)
	require.Truef(t, ok, "widgetSpec should be a map, got %T", up.Outputs["widgetSpec"].Value)

	// The API server prunes fields its structural schema does not know, so a value
	// read back under the hyphenated key proves both directions of the round trip.
	require.Equal(t, "reserved:world", spec["security-labels"])
	require.NotContains(t, spec, "security_labels")
	require.EqualValues(t, 3, spec["replicas"])
}
