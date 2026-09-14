package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
)

func TestExtensionTwoPackagesOnOneProvider(t *testing.T) {
	test := pulumitest.NewPulumiTest(t, "testdata/extension-multi", opttest.SkipInstall())
	t.Cleanup(func() {
		test.Destroy(t)
	})

	packageAdd := packageAddCmd(t, test)

	packageAdd("--extension", "name=dragon-ext crd-manifest=dragon-crds.yaml")
	packageAdd("--extension", "name=unicorn-ext crd-manifest=unicorn-crds.yaml")

	// Preview is where one extension used to displace the other, failing with
	// "unrecognized resource type" for whichever package lost.
	test.Preview(t)

	up := test.Up(t)
	require.Equal(t, "dragon", up.Outputs["dragonName"].Value,
		"the dragon extension should be served after the unicorn extension parameterizes the provider")
	require.Equal(t, "unicorn", up.Outputs["unicornName"].Value,
		"the unicorn extension should be served after the dragon extension parameterizes the provider")
	require.Equal(t, "multiext-system", up.Outputs["namespaceName"].Value,
		"base-provider resources should be served alongside both extensions")

	require.Equal(t, "green", up.Outputs["dragonScaleColor"].Value,
		"the dragon extension's own schema should survive the unicorn parameterization")
	require.EqualValues(t, 3, up.Outputs["unicornHornLength"].Value,
		"the unicorn extension's own schema should survive the dragon parameterization")
}
