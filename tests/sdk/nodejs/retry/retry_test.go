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
	"testing"

	"github.com/pulumi/pulumi/sdk/v2/go/common/tokens"

	"github.com/pulumi/pulumi-kubernetes/provider/v2/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/tests/v2"
	"github.com/pulumi/pulumi/pkg/v2/resource/deploy/providers"
	"github.com/pulumi/pulumi/pkg/v2/testing/integration"
	"github.com/pulumi/pulumi/sdk/v2/go/common/resource"
	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	integration.ProgramTest(t, &integration.ProgramTestOptions{
		Dir:          "step1",
		Dependencies: []string{"@pulumi/kubernetes"},
		Quick:        true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 6, len(stackInfo.Deployment.Resources))

			tests.SortResourcesByURN(stackInfo)

			stackRes := stackInfo.Deployment.Resources[4]
			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

			provRes := stackInfo.Deployment.Resources[2]
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			// Assert the Namespace was created
			namespace := stackInfo.Deployment.Resources[0]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())

			// Assert the Pod was created
			pod := stackInfo.Deployment.Resources[1]
			assert.Equal(t, "nginx", string(pod.URN.Name()))
			step1Name, _ := openapi.Pluck(pod.Outputs, "metadata", "name")
			assert.Equal(t, "nginx", step1Name.(string))
			step1PodNamespace, _ := openapi.Pluck(pod.Outputs, "metadata", "namespace")
			assert.Equal(t, namespace.ID.String(), step1PodNamespace.(string))
		},
	})
}
