package provider

import (
	"context"
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRename(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	test := pulumitest.NewPulumiTest(t, "rename")
	t.Cleanup(func() {
		test.Destroy()
	})

	test.Up()

	outputs, err := test.CurrentStack().Outputs(ctx)
	require.NoError(t, err)

	// The ConfigMap should have 3 managed fields -- one for the URN annotation
	// and 2 for the patches.
	mf, ok := outputs["managedFields"]
	require.True(t, ok)
	assert.Len(t, mf.Value, 3)

	// Change our Namespace's resource name and delete a patch.
	test.UpdateSource("rename/step2")
	test.Up()

	// Renaming the namespace should not have deleted it. Perform a refresh and
	// make sure our pod is still running -- if it's not, Pulumi will have
	// deleted it from our state.
	refresh, err := test.CurrentStack().Refresh(ctx)
	assert.NoError(t, err)
	assert.NotContains(t, refresh.StdOut, "deleted", refresh.StdOut)

	// One ConfigMapPatch should still be applied, plus a manager for our URN
	// annotation.
	outputs, err = test.CurrentStack().Outputs(ctx)
	require.NoError(t, err)
	mf, ok = outputs["managedFields"]
	require.True(t, ok)
	assert.Len(t, mf.Value, 2)
}
