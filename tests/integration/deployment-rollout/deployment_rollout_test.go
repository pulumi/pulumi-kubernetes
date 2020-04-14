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
	"strings"
	"testing"

	"github.com/pulumi/pulumi-kubernetes/provider/v2/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/tests/v2"
	"github.com/pulumi/pulumi/pkg/v2/resource/deploy/providers"
	"github.com/pulumi/pulumi/pkg/v2/testing/integration"
	"github.com/pulumi/pulumi/sdk/v2/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v2/go/common/tokens"
	"github.com/stretchr/testify/assert"
)

func TestDeploymentRollout(t *testing.T) {
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
			assert.Equal(t, 5, len(stackInfo.Deployment.Resources))

			tests.SortResourcesByURN(stackInfo)

			appsv1Deploy := stackInfo.Deployment.Resources[0]
			namespace := stackInfo.Deployment.Resources[1]
			extensionsv1beta1Deploy := stackInfo.Deployment.Resources[2]
			provRes := stackInfo.Deployment.Resources[3]
			stackRes := stackInfo.Deployment.Resources[4]

			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())

			//
			// Assert deployments are successfully created.
			//

			name, _ := openapi.Pluck(appsv1Deploy.Outputs, "metadata", "name")
			assert.True(t, strings.Contains(name.(string), "nginx"))
			containers, _ := openapi.Pluck(appsv1Deploy.Outputs, "spec", "template", "spec", "containers")
			containerStatus := containers.([]interface{})[0].(map[string]interface{})
			image := containerStatus["image"]
			assert.Equal(t, image.(string), "nginx")

			name, _ = openapi.Pluck(extensionsv1beta1Deploy.Outputs, "metadata", "name")
			assert.True(t, strings.Contains(name.(string), "nginx"))
			containers, _ = openapi.Pluck(appsv1Deploy.Outputs, "spec", "template", "spec", "containers")
			containerStatus = containers.([]interface{})[0].(map[string]interface{})
			image = containerStatus["image"]
			assert.Equal(t, image.(string), "nginx")
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      "step2",
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 5, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					appsv1Deploy := stackInfo.Deployment.Resources[0]
					namespace := stackInfo.Deployment.Resources[1]
					extensionsv1beta1Deploy := stackInfo.Deployment.Resources[2]
					provRes := stackInfo.Deployment.Resources[3]
					stackRes := stackInfo.Deployment.Resources[4]

					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())

					//
					// Assert deployments are updated successfully.
					//

					name, _ := openapi.Pluck(appsv1Deploy.Outputs, "metadata", "name")
					assert.True(t, strings.Contains(name.(string), "nginx"))
					containers, _ := openapi.Pluck(appsv1Deploy.Outputs, "spec", "template", "spec", "containers")
					containerStatus := containers.([]interface{})[0].(map[string]interface{})
					image := containerStatus["image"]
					assert.Equal(t, image.(string), "nginx:stable")

					name, _ = openapi.Pluck(extensionsv1beta1Deploy.Outputs, "metadata", "name")
					assert.True(t, strings.Contains(name.(string), "nginx-ev1b1"))
					containers, _ = openapi.Pluck(appsv1Deploy.Outputs, "spec", "template", "spec", "containers")
					containerStatus = containers.([]interface{})[0].(map[string]interface{})
					image = containerStatus["image"]
					assert.Equal(t, image.(string), "nginx:stable")
				},
			},
		},
	})
}
