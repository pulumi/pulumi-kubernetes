package provider

import (
	"context"
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/stretchr/testify/assert"
)

func TestRename(t *testing.T) {
	test := pulumitest.NewPulumiTest(t, "rename")
	t.Cleanup(func() {
		test.Destroy()
	})
	test.Up()

	test.UpdateSource("rename/step2")
	test.Up()

	r, err := test.CurrentStack().Refresh(context.Background())
	assert.NoError(t, err)

	assert.NotContains(t, r.StdOut, "deleted", r.StdOut)
}
