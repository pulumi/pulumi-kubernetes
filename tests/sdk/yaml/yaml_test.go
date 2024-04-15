// Copyright 2024, Pulumi Corporation.
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
	"os"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAwaitDaemonSet(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	dir := filepath.Join(cwd, "await-daemonset")

	assertReady := func(t *testing.T, stack integration.RuntimeValidationStackInfo) {
		cns, ok := stack.Outputs["currentNumberScheduled"]
		require.True(t, ok)

		dns, ok := stack.Outputs["desiredNumberScheduled"]
		require.True(t, ok)

		nms, ok := stack.Outputs["numberMisscheduled"]
		require.True(t, ok)

		nr, ok := stack.Outputs["numberReady"]
		require.True(t, ok)

		assert.Greater(t, cns, float64(0))
		assert.Equal(t, float64(0), nms)
		assert.Equal(t, dns, cns)
		assert.Equal(t, cns, nr)
	}

	test := integration.ProgramTestOptions{
		// Await creation.
		Dir:                    dir,
		ExtraRuntimeValidation: assertReady,
		EditDirs: []integration.EditDir{
			{
				// Await successful update.
				Dir:                    filepath.Join(dir, "step2"),
				Additive:               true,
				ExtraRuntimeValidation: assertReady,
			},
			{
				// Await unsuccessful update -- should fail.
				Dir:           filepath.Join(dir, "step3"),
				Additive:      true,
				ExpectFailure: true,
			},
		},
	}

	integration.ProgramTest(t, &test)
}
