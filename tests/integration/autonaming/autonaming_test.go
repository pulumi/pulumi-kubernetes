// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package ints

import (
	"os"
	"strings"
	"testing"

	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/tests"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/deploy/providers"
	"github.com/pulumi/pulumi/pkg/testing/integration"
	"github.com/stretchr/testify/assert"
)

var step1Name interface{}
var step2Name interface{}
var step3Name interface{}

func TestAutonaming(t *testing.T) {
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
			assert.Equal(t, 3, len(stackInfo.Deployment.Resources))

			tests.SortResourcesByURN(stackInfo)

			stackRes := stackInfo.Deployment.Resources[2]
			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

			provRes := stackInfo.Deployment.Resources[1]
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			//
			// Assert Pod is successfully given a unique name by Pulumi.
			//

			pod := stackInfo.Deployment.Resources[0]
			assert.Equal(t, "autonaming-test", string(pod.URN.Name()))
			step1Name, _ = openapi.Pluck(pod.Outputs, "metadata", "name")
			assert.True(t, strings.HasPrefix(step1Name.(string), "autonaming-test-"))

			autonamed, _ := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
			assert.Equal(t, "true", autonamed)
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      "step2",
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 3, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					stackRes := stackInfo.Deployment.Resources[2]
					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

					provRes := stackInfo.Deployment.Resources[1]
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					//
					// Assert Pod was replaced, i.e., destroyed and re-created, with allocating a new name.
					//

					pod := stackInfo.Deployment.Resources[0]
					assert.Equal(t, "autonaming-test", string(pod.URN.Name()))
					step2Name, _ = openapi.Pluck(pod.Outputs, "metadata", "name")
					assert.True(t, strings.HasPrefix(step2Name.(string), "autonaming-test-"))

					autonamed, _ := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
					assert.Equal(t, "true", autonamed)

					assert.NotEqual(t, step1Name, step2Name)

				},
			},
			{
				Dir:      "step3",
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 3, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					stackRes := stackInfo.Deployment.Resources[2]
					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

					provRes := stackInfo.Deployment.Resources[1]
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					//
					// Assert Pod was NOT replaced, and has the same name, previously allocated by Pulumi.
					//

					pod := stackInfo.Deployment.Resources[0]
					assert.Equal(t, "autonaming-test", string(pod.URN.Name()))
					step3Name, _ = openapi.Pluck(pod.Outputs, "metadata", "name")
					assert.True(t, strings.HasPrefix(step3Name.(string), "autonaming-test-"))

					autonamed, _ := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
					assert.Equal(t, "true", autonamed)

					assert.Equal(t, step2Name, step3Name)
				},
			},
			{
				Dir:      "step4",
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 3, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					stackRes := stackInfo.Deployment.Resources[2]
					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

					provRes := stackInfo.Deployment.Resources[1]
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					//
					// User has specified their own name for the Pod, so we replace it, and Pulumi does NOT
					// allocate a name on its own.
					//

					pod := stackInfo.Deployment.Resources[0]
					assert.Equal(t, "autonaming-test", string(pod.URN.Name()))
					name, _ := openapi.Pluck(pod.Outputs, "metadata", "name")
					assert.Equal(t, "autonaming-test", name.(string))

					_, autonamed := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
					assert.False(t, autonamed)
				},
			},
		},
	})
}
