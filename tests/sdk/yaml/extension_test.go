package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

	providerBin := filepath.Join(repoRoot(t), "bin", "pulumi-resource-kubernetes")

	// Declare the base kubernetes package explicitly so base-namespace types
	// (core/v1:Namespace) resolve without relying on the pulumi-yaml
	// package-resolution fix being released.
	// TODO #pulumi/pulumi-yaml/1162: drop the base add once the fix is released.
	base := exec.Command("pulumi", "package", "add", providerBin)
	base.Dir = test.WorkingDir()
	base.Env = os.Environ()
	if out, err := base.CombinedOutput(); err != nil {
		t.Fatalf("pulumi package add (base) failed: %v\n%s", err, out)
	}

	// Add the extension from the local provider binary: `package add --extension`
	// parameterizes the base provider and records the extension in Pulumi.yaml,
	// without downloading an (unpublished) dev provider.
	add := exec.Command("pulumi", "package", "add", providerBin,
		"--extension", "name=gateway-networking crd-manifest=gateway-api-crds.yaml")
	add.Dir = test.WorkingDir()
	add.Env = os.Environ()
	if out, err := add.CombinedOutput(); err != nil {
		t.Fatalf("pulumi package add --extension failed: %v\n%s", err, out)
	}

	test.Preview(t)

	up := test.Up(t)
	require.Equal(t, "example-class", up.Outputs["gatewayClassName"].Value,
		"extension-served GatewayClass should be created with its declared name")
	require.Equal(t, "gateway-system", up.Outputs["namespaceName"].Value,
		"base-provider Namespace should be created alongside the extension resource")
}

func repoRoot(t *testing.T) string {
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "could not determine test file path")
	// tests/sdk/yaml/extension_test.go -> repo root is three directories up.
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "..")
}
