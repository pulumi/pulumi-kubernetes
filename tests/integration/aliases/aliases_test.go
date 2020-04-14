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
	"os"
	"testing"

	"github.com/pulumi/pulumi-kubernetes/provider/v2/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/tests/v2"
	"github.com/pulumi/pulumi/pkg/v2/testing/integration"
	"github.com/stretchr/testify/assert"
)

func TestAliases(t *testing.T) {
	kubectx := os.Getenv("KUBERNETES_CONTEXT")

	if kubectx == "" {
		t.Skipf("Skipping test due to missing KUBERNETES_CONTEXT variable")
	}

	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir:          "step1",
		Dependencies: []string{"@pulumi/kubernetes"},
		Quick:        true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

			tests.SortResourcesByURN(stackInfo)

			deployment := stackInfo.Deployment.Resources[0]
			assert.Equal(t, "alias-test", string(deployment.URN.Name()))
			assert.Equal(t, "kubernetes:extensions/v1beta1:Deployment", string(deployment.Type))
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      "step2",
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					deployment := stackInfo.Deployment.Resources[0]
					assert.Equal(t, "alias-test", string(deployment.URN.Name()))
					assert.Equal(t, "kubernetes:apps/v1:Deployment", string(deployment.Type))
					containers, _ := openapi.Pluck(deployment.Outputs, "spec", "template", "spec", "containers")
					containerStatus := containers.([]interface{})[0].(map[string]interface{})
					image := containerStatus["image"]
					assert.Equal(t, image.(string), "nginx:1.14")
				},
			},
			{
				Dir:      "step3",
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					deployment := stackInfo.Deployment.Resources[0]
					assert.Equal(t, "alias-test", string(deployment.URN.Name()))
					assert.Equal(t, "kubernetes:apps/v1:Deployment", string(deployment.Type))
					containers, _ := openapi.Pluck(deployment.Outputs, "spec", "template", "spec", "containers")
					containerStatus := containers.([]interface{})[0].(map[string]interface{})
					image := containerStatus["image"]
					assert.Equal(t, image.(string), "nginx:1.15")
				},
			},
		},
	})
}
