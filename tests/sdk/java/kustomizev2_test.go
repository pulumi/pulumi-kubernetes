package test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
)

// TestKustomizeV2 exercises the kustomize/v2 package.
// - local and remote targets
// - alpha plugin feature
// - helm chart support
func TestKustomizeV2(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("PrintPluginEnv is not supported on Windows")
	}
	pluginHome, _ := filepath.Abs("testdata/kustomizev2/pluginExample/plugin")
	test := pulumitest.NewPulumiTest(t, "testdata/kustomizev2",
		opttest.Env("KUSTOMIZE_PLUGIN_HOME", pluginHome),
	)
	t.Logf("into %s", test.Source())
	t.Cleanup(func() {
		test.Destroy()
	})
	test.Preview()
	test.Up()
}
