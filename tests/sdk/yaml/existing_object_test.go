// Copyright 2016-2026, Pulumi Corporation.
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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
)

// TestExistingObjectBlocksCreate tests that creating a resource via
// server-side apply fails when the object already exists in the cluster.
// The "default" namespace is a built-in that always exists, so no setup
// is needed.
func TestExistingObjectBlocksCreate(t *testing.T) {
	ctx := context.Background()

	test := pulumitest.NewPulumiTest(t, "testdata/existing-object", opttest.SkipInstall())
	t.Logf("into %s", test.WorkingDir())

	_, upErr := test.CurrentStack().Up(ctx)
	require.Error(t, upErr, "expected Up to fail when object already exists")
	assert.Contains(t, upErr.Error(), "already exists",
		"error should mention that the resource already exists")
}
