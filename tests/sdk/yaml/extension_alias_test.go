package test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

const (
	extensionGatewayClassType = "gateway-networking:gateway.networking.k8s.io/v1:GatewayClass"
	baseGatewayClassType      = "kubernetes:gateway.networking.k8s.io/v1:GatewayClass"
)

// Reproduces the state a crd2pulumi-generated SDK would have written, without generating one.
func rewriteToBaseNamespace(t *testing.T, exported apitype.UntypedDeployment) apitype.UntypedDeployment {
	t.Helper()

	var deployment apitype.DeploymentV3
	require.NoError(t, json.Unmarshal(exported.Deployment, &deployment))

	rewritten := 0
	for i := range deployment.Resources {
		r := &deployment.Resources[i]
		if string(r.Type) != extensionGatewayClassType {
			continue
		}
		r.Type = tokens.Type(baseGatewayClassType)
		r.URN = resource.URN(strings.Replace(string(r.URN),
			extensionGatewayClassType, baseGatewayClassType, 1))
		rewritten++
	}
	require.Equal(t, 1, rewritten, "expected exactly one extension-served GatewayClass in state")

	raw, err := json.Marshal(deployment)
	require.NoError(t, err)
	return apitype.UntypedDeployment{Version: exported.Version, Deployment: raw}
}

func TestExtensionAliasAdoptsBaseNamespacedState(t *testing.T) {
	test := pulumitest.NewPulumiTest(t, "testdata/extension-alias-migration", opttest.SkipInstall())
	t.Cleanup(func() {
		test.Destroy(t)
	})

	packageAdd := packageAddCmd(t, test)
	packageAdd()
	packageAdd("--extension", "name=gateway-networking crd-manifest=gateway-api-crds.yaml")

	created := test.Up(t)
	require.NotNil(t, created.Summary.ResourceChanges)
	require.NotZero(t, (*created.Summary.ResourceChanges)["create"],
		"the extension-served GatewayClass should be created first")

	test.ImportStack(t, rewriteToBaseNamespace(t, test.ExportStack(t)))

	migrated := test.Up(t)

	require.NotNil(t, migrated.Summary.ResourceChanges)
	changes := *migrated.Summary.ResourceChanges
	assert.Zero(t, changes["replace"],
		"the base-namespaced GatewayClass must be adopted via its alias, not replaced")
	assert.Zero(t, changes["delete"],
		"the base-namespaced GatewayClass must be adopted via its alias, not deleted")
	assert.Zero(t, changes["create"],
		"the base-namespaced GatewayClass must be adopted via its alias, not recreated")
}
