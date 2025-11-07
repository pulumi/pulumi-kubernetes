// Copyright 2025, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"context"
	"testing"
	"time"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteDueToRename(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	test := pulumitest.NewPulumiTest(t, "testdata/delete/rename")
	t.Cleanup(func() {
		test.Destroy(t)
	})

	test.Up(t)

	// Change our Pod's resource name
	test.UpdateSource(t, "testdata/delete/rename/step2")
	test.Up(t)

	// Renaming the namespace should not have deleted it. Perform a refresh and
	// make sure our pod is still running -- if it's not, Pulumi will have
	// deleted it from our state.
	refresh, err := test.CurrentStack().Refresh(ctx)
	assert.NoError(t, err)
	assert.NotContains(t, refresh.StdOut, "deleted", refresh.StdOut)
}

func TestDeletePatchResource(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	test := pulumitest.NewPulumiTest(t, "testdata/delete/patch")
	t.Cleanup(func() {
		test.Destroy(t)
	})

	test.Up(t)

	time.Sleep(60 * time.Second)

	outputs, err := test.CurrentStack().Outputs(ctx)
	require.NoError(t, err)

	// The ConfigMap should have 2 managed fields.
	mf, ok := outputs["managedFields"]
	require.True(t, ok)
	assert.Len(t, mf.Value, 2)

	// Delete a patch.
	test.UpdateSource(t, "testdata/delete/patch/step2")
	test.Up(t)

	// One ConfigMapPatch should still be applied.
	outputs, err = test.CurrentStack().Outputs(ctx)
	require.NoError(t, err)
	mf, ok = outputs["managedFields"]
	require.True(t, ok)
	assert.Len(t, mf.Value, 1)
}
