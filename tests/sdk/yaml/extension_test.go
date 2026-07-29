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
	t.Skip("Skipping until YAML resolver bug is fixed")
	t.Cleanup(func() {
		test.Destroy(t)
	})

	// Make the extension package from the local provider binary before running the
	// program: `package add` parameterizes the base provider and records the
	// extension in Pulumi.yaml, without downloading an (unpublished) dev provider.
	providerBin := filepath.Join(repoRoot(t), "bin", "pulumi-resource-kubernetes")
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
