package test

import (
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/stretchr/testify/assert"
)

func TestConfigMapAndSecretImmutability(t *testing.T) {
	t.Parallel()

	test := pulumitest.NewPulumiTest(t,
		"testdata/immutability",
		opttest.SkipInstall(),
	)
	t.Cleanup(func() {
		test.Destroy()
	})

	// Create the secrets/configmaps.
	up := test.Up()

	// We will detect update/replacement behavior by observing effects on our
	// downstream dependencies.
	secret := up.Outputs["secret"].Value.(string)
	configmap := up.Outputs["configmap"].Value.(string)
	autonamedSecret := up.Outputs["autonamedSecret"].Value.(string)
	autonamedConfigmap := up.Outputs["autonamedConfigmap"].Value.(string)
	mutableSecret := up.Outputs["mutableSecret"].Value.(string)
	mutableConfigmap := up.Outputs["mutableConfigmap"].Value.(string)

	// Update the data of all our secrets and configmaps.
	test.UpdateSource("testdata/immutability/step2")
	up = test.Up()

	// Only the mutable configmap and secret should have been updated -- so no
	// impact on those two downstreams.
	assert.Equal(t, mutableConfigmap, up.Outputs["mutableConfigmap"].Value.(string))
	assert.Equal(t, mutableSecret, up.Outputs["mutableSecret"].Value.(string))
	// All others should have been replaced, which should have regenerated our
	// random pets.
	assert.NotEqual(t, secret, up.Outputs["secret"].Value.(string))
	assert.NotEqual(t, configmap, up.Outputs["configmap"].Value.(string))
	assert.NotEqual(t, autonamedSecret, up.Outputs["autonamedSecret"].Value.(string))
	assert.NotEqual(t, autonamedConfigmap, up.Outputs["autonamedConfigmap"].Value.(string))
	// Record the new outputs.
	secret = up.Outputs["secret"].Value.(string)
	configmap = up.Outputs["configmap"].Value.(string)
	autonamedSecret = up.Outputs["autonamedSecret"].Value.(string)
	autonamedConfigmap = up.Outputs["autonamedConfigmap"].Value.(string)

	// The final step only touches annotations. All resources should have been
	// updated.
	test.UpdateSource("testdata/immutability/step3")
	up = test.Up()
	assert.Equal(t, secret, up.Outputs["secret"].Value.(string))
	assert.Equal(t, configmap, up.Outputs["configmap"].Value.(string))
	assert.Equal(t, autonamedSecret, up.Outputs["autonamedSecret"].Value.(string))
	assert.Equal(t, autonamedConfigmap, up.Outputs["autonamedConfigmap"].Value.(string))
	assert.Equal(t, mutableConfigmap, up.Outputs["mutableConfigmap"].Value.(string))
	assert.Equal(t, mutableSecret, up.Outputs["mutableSecret"].Value.(string))
}
