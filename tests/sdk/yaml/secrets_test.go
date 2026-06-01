package test

import (
	"encoding/base64"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

func TestSecrets(t *testing.T) {
	t.Parallel()

	const secretMessage = "secret message for testing"

	test := pulumitest.NewPulumiTest(t,
		"testdata/secrets",
		opttest.SkipInstall(),
	)
	t.Cleanup(func() {
		test.Destroy(t)
	})

	err := test.CurrentStack().SetConfig(test.Context(), "message", auto.ConfigValue{
		Value:  secretMessage,
		Secret: true,
	})
	require.NoError(t, err)

	test.Up(t)

	// auto.Stack.Export always passes --show-secrets, which decrypts secrets
	// to plaintext in the export. Shell out directly so we get the encrypted
	// state, matching the assertion in the legacy integration.ProgramTest
	// version of this scenario.
	state := exportState(t, test)

	require.NotContains(t, state, secretMessage,
		"secret value must not appear in plaintext state")
	require.NotContains(t, state, base64.StdEncoding.EncodeToString([]byte(secretMessage)),
		"base64 of secret value must not appear in state")
}

func exportState(t *testing.T, test *pulumitest.PulumiTest) string {
	t.Helper()
	ws := test.CurrentStack().Workspace()
	// #nosec G204 -- stack name comes from pulumitest internals, not user input.
	cmd := exec.CommandContext(test.Context(), "pulumi",
		"stack", "export", "--stack", test.CurrentStack().Name())
	cmd.Dir = ws.WorkDir()
	cmd.Env = os.Environ()
	for k, v := range ws.GetEnvVars() {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	out, err := cmd.Output()
	require.NoError(t, err, "pulumi stack export failed")
	return string(out)
}
