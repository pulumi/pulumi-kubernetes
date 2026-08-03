package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
)

func TestExtensionGatewayAPI(t *testing.T) {
	test := pulumitest.NewPulumiTest(t, "testdata/extension-gateway-api", opttest.SkipInstall())
	t.Cleanup(func() {
		test.Destroy(t)
	})

	cwd, err := os.Getwd()
	require.NoError(t, err)
	providerBin := filepath.Join(cwd, "..", "..", "..", "bin", "pulumi-resource-kubernetes")

	pulumi := test.CurrentStack().Workspace().PulumiCommand()
	packageAdd := func(args ...string) {
		t.Helper()
		_, stderr, _, err := pulumi.Run(test.Context(), test.WorkingDir(),
			nil, nil, nil, nil, append([]string{"package", "add", providerBin}, args...)...)
		require.NoErrorf(t, err, "pulumi package add %v failed: %s", args, stderr)
	}

	// Declare the base kubernetes package explicitly so base-namespace types
	// (core/v1:Namespace) resolve without relying on the pulumi-yaml
	// package-resolution fix being released.
	// TODO #pulumi/pulumi-yaml/1162: drop the base add once the fix is released.
	packageAdd()
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
