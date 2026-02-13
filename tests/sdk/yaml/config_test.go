package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pulumi/providertest/pulumitest"
)

func TestClusterIdentifier(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	test := pulumitest.NewPulumiTest(t, "config/cluster-identifier")
	t.Cleanup(func() {
		test.Destroy(t)
	})

	test.Up(t)

	test.UpdateSource(t, "config/cluster-identifier/step2")
	up, err := test.CurrentStack().Up(ctx)

	require.NoError(t, err)
	assert.Contains(t, up.StdOut, "updated")
	assert.NotContains(t, up.StdOut, "replaced")

	test.UpdateSource(t, "config/cluster-identifier/step3")
	up, err = test.CurrentStack().Up(ctx)

	require.NoError(t, err)
	assert.NotContains(t, up.StdOut, "updated")
	assert.Contains(t, up.StdOut, "replaced")
}
