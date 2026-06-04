package test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
)

func TestHelmReleaseCRDScopeWritesToCache(t *testing.T) {
	test := pulumitest.NewPulumiTest(t, "testdata/helm-release-crd-scope", opttest.SkipInstall())
	t.Cleanup(func() {
		test.Destroy(t)
	})
	test.Preview(t)
	test.Up(t)

	var deployment apitype.DeploymentV3
	require.NoError(t, json.Unmarshal(test.ExportStack(t).Deployment, &deployment))

	var found bool
	for _, r := range deployment.Resources {
		if r.Type != "kubernetes:stable.example.com/v1:CronTab" {
			continue
		}
		found = true
		assert.Contains(t, string(r.URN), "default/my-crontab",
			"CronTab should resolve as namespaced via the cached CRD")
	}
	assert.True(t, found, "expected a CronTab custom resource in the stack")
}
