package test

import (
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
)

// TestChartV4 deploys a complex stack using chart/v4 package.
func TestChartv4(t *testing.T) {
	test := pulumitest.NewPulumiTest(t, "testdata/chartv4", opttest.SkipInstall())
	t.Logf("into %s", test.Source())
	t.Cleanup(func() {
		test.Destroy()
	})
	test.Preview()
	test.Up()
	test.Up(optup.ExpectNoChanges())
}
