package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

func TestSecrets(t *testing.T) {
	t.Parallel()

	test := pulumitest.NewPulumiTest(t,
		"testdata/secrets",
		opttest.SkipInstall(),
	)
	t.Cleanup(func() {
		test.Destroy(t)
	})

	err := test.CurrentStack().SetConfig(test.Context(), "message", auto.ConfigValue{
		Value:  "secret message for testing",
		Secret: true,
	})
	require.NoError(t, err)

	up := test.Up(t)

	for _, k := range []string{"cmDataPassword", "cmBinaryDataPassword", "sStringDataPassword", "sDataPassword"} {
		require.Truef(t, up.Outputs[k].Secret, "output %q should be marked secret", k)
	}
}
