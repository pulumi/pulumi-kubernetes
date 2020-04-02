// Copyright 2016-2019, Pulumi Corporation.
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

package ints

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi/pkg/v2/testing/integration"
)

func TestRenderYAML(t *testing.T) {
	kubectx := os.Getenv("KUBERNETES_CONTEXT")

	if kubectx == "" {
		t.Skipf("Skipping test due to missing KUBERNETES_CONTEXT variable")
	}

	// Create a temporary directory to hold rendered YAML manifests.
	dir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Config: map[string]string{
			"renderDir": dir,
		},
		Dir:           "step1",
		Dependencies:  []string{"@pulumi/kubernetes"},
		Quick:         true,
		ExpectFailure: true, // step1 should fail because of an invalid Provider config.
		EditDirs: []integration.EditDir{
			{
				Dir:           "step2",
				Additive:      true,
				ExpectFailure: false,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

					// Verify that YAML directory was created.
					files, err := ioutil.ReadDir(dir)
					assert.NoError(t, err)
					assert.Equal(t, len(files), 2)

					// Verify that CRD manifest directory was created.
					files, err = ioutil.ReadDir(filepath.Join(dir, "0-crd"))
					assert.NoError(t, err)
					assert.Equal(t, len(files), 0)

					// Verify that manifest directory was created.
					files, err = ioutil.ReadDir(filepath.Join(dir, "1-manifest"))
					assert.NoError(t, err)
					assert.Equal(t, len(files), 2)
				},
			},
		},
	})
}
