package test

import (
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/stretchr/testify/assert"
)

// TestAutonaming ensures that custom resource autonaming configuration works as expected.
func TestAutonaming(t *testing.T) {
	test := pulumitest.NewPulumiTest(
		t,
		"testdata/autonaming",
		opttest.SkipInstall(),
		opttest.Env("PULUMI_EXPERIMENTAL", "1"),
	)
	t.Logf("into %s", test.WorkingDir())
	t.Cleanup(func() {
		test.Destroy(t)
	})
	test.Preview(t)
	up := test.Up(t)
	nsname, ok := up.Outputs["nsname"].Value.(string)
	assert.True(t, ok)
	assert.Contains(t, nsname, "autonaming-ns-") // project + name + random suffix
}
