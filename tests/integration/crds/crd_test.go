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

	"github.com/pulumi/pulumi-kubernetes/tests/v2"
	"github.com/pulumi/pulumi/pkg/resource/deploy/providers"
	"github.com/pulumi/pulumi/pkg/testing/integration"
	"github.com/pulumi/pulumi/sdk/go/common/resource"
	"github.com/pulumi/pulumi/sdk/go/common/tokens"
	"github.com/stretchr/testify/assert"
)

func TestCRDs(t *testing.T) {
	kubectx := os.Getenv("KUBERNETES_CONTEXT")

	if kubectx == "" {
		t.Skipf("Skipping test due to missing KUBERNETES_CONTEXT variable")
	}

	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir:                  "step1",
		Dependencies:         []string{"@pulumi/kubernetes"},
		Quick:                false,
		ExpectRefreshChanges: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 5, len(stackInfo.Deployment.Resources))

			tests.SortResourcesByURN(stackInfo)

			crd := stackInfo.Deployment.Resources[0]
			namespace := stackInfo.Deployment.Resources[1]
			ct1 := stackInfo.Deployment.Resources[2]
			provRes := stackInfo.Deployment.Resources[3]
			stackRes := stackInfo.Deployment.Resources[4]

			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())

			//
			// Assert that CRD and CR exist
			//

			assert.Equal(t, "foobar", string(crd.URN.Name()))
			assert.Equal(t, "my-new-foobar-object", string(ct1.URN.Name()))
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      "step2",
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					namespace := stackInfo.Deployment.Resources[0]
					provRes := stackInfo.Deployment.Resources[2]
					stackRes := stackInfo.Deployment.Resources[3]

					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())
				},
			},
		},
	})
}
