package provider

import (
	"context"
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/stretchr/testify/assert"
)

func TestCreateWithoutUpsert(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	test := pulumitest.NewPulumiTest(t, "create/without-upsert")
	t.Cleanup(func() {
		test.Destroy()
	})

	test.Up()

	// Create a new pod resource referencing one that already exists -- should fail.
	test.UpdateSource("create/without-upsert/step2")
	_, err := test.CurrentStack().Up(ctx)

	assert.ErrorContains(t, err, "already exists")

	// Try again with upsert re-enabled -- should succeed.
	test.UpdateSource("create/without-upsert/step3")
	test.Up()
}
