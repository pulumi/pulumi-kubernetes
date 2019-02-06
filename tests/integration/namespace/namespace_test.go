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

	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"

	"github.com/pulumi/pulumi/pkg/tokens"

	"github.com/pulumi/pulumi-kubernetes/tests"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/deploy/providers"
	"github.com/pulumi/pulumi/pkg/testing/integration"
	"github.com/stretchr/testify/assert"
)

func TestNamespace(t *testing.T) {
	kubectx := os.Getenv("KUBERNETES_CONTEXT")

	if kubectx == "" {
		t.Skipf("Skipping test due to missing KUBERNETES_CONTEXT variable")
	}

	var nmPodName, defaultPodName string
	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir:          "step1",
		Dependencies: []string{"@pulumi/kubernetes"},
		Quick:        true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 5, len(stackInfo.Deployment.Resources))

			tests.SortResourcesByURN(stackInfo)

			stackRes := stackInfo.Deployment.Resources[4]
			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

			provRes := stackInfo.Deployment.Resources[3]
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			// Assert the Namespace was created
			namespace := stackInfo.Deployment.Resources[0]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())

			// Assert the "no metadata" Pod was created in the "default" namespace.
			nmPod := stackInfo.Deployment.Resources[2]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Pod"), nmPod.URN.Type())
			nmPodNamespace, _ := openapi.Pluck(nmPod.Outputs, "metadata", "namespace")
			assert.Equal(t, nmPodNamespace.(string), "default")
			nmPodNameRaw, _ := openapi.Pluck(nmPod.Outputs, "metadata", "name")
			nmPodName = nmPodNameRaw.(string)
			assert.True(t, strings.HasPrefix(nmPodName, "no-metadata-pod"))

			// Assert the "explicit default namespace" Pod was created in the "default" namespace.
			defaultPod := stackInfo.Deployment.Resources[1]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Pod"), defaultPod.URN.Type())
			defaultPodNamespace, _ := openapi.Pluck(defaultPod.Outputs, "metadata", "namespace")
			assert.Equal(t, defaultPodNamespace.(string), "default")
			defaultPodNameRaw, _ := openapi.Pluck(defaultPod.Outputs, "metadata", "name")
			defaultPodName = defaultPodNameRaw.(string)
			assert.True(t, strings.HasPrefix(defaultPodName, "default-ns-pod"))
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      "step2",
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 5, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					stackRes := stackInfo.Deployment.Resources[4]
					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

					provRes := stackInfo.Deployment.Resources[3]
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					// Assert that the Namespace was updated with the expected label.
					namespace := stackInfo.Deployment.Resources[0]
					assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())
					namespaceLabels, _ := openapi.Pluck(namespace.Outputs, "metadata", "labels")
					assert.Equal(t, map[string]interface{}{"hello": "world"},
						namespaceLabels.(map[string]interface{}))
				},
			},
			{
				Dir:      "step3",
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 5, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					stackRes := stackInfo.Deployment.Resources[4]
					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

					provRes := stackInfo.Deployment.Resources[3]
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					namespace := stackInfo.Deployment.Resources[0]
					assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())

					// Assert the "no metadata" Pod has not changed.
					nmPod := stackInfo.Deployment.Resources[2]
					assert.Equal(t, tokens.Type("kubernetes:core/v1:Pod"), nmPod.URN.Type())
					nmPodNamespace, _ := openapi.Pluck(nmPod.Outputs, "metadata", "namespace")
					assert.Equal(t, nmPodNamespace.(string), "default")
					nmPodNameRaw, _ := openapi.Pluck(nmPod.Outputs, "metadata", "name")
					assert.True(t, strings.HasPrefix(nmPodNameRaw.(string), "no-metadata-pod"))
					assert.Equal(t, nmPodNameRaw.(string), nmPodName)

					// Assert the "explicit default namespace" has not changed.
					defaultPod := stackInfo.Deployment.Resources[1]
					assert.Equal(t, tokens.Type("kubernetes:core/v1:Pod"), defaultPod.URN.Type())
					defaultPodNamespace, _ := openapi.Pluck(defaultPod.Outputs, "metadata", "namespace")
					assert.Equal(t, defaultPodNamespace.(string), "default")
					defaultPodNameRaw, _ := openapi.Pluck(defaultPod.Outputs, "metadata", "name")
					assert.True(t, strings.HasPrefix(defaultPodNameRaw.(string), "default-ns-pod"))
					assert.Equal(t, defaultPodNameRaw.(string), defaultPodName)
				},
			},
		},
	})
}
