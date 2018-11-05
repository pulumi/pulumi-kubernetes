// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package ints

import (
	"os"
	"testing"

	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/tests"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/deploy/providers"
	"github.com/pulumi/pulumi/pkg/testing/integration"
	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
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

			stackRes := stackInfo.Deployment.Resources[3]
			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

			provRes := stackInfo.Deployment.Resources[2]
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			// Assert the Pod was created
			pod := stackInfo.Deployment.Resources[1]
			assert.Equal(t, "nginx", string(pod.URN.Name()))
			step1Name, _ := openapi.Pluck(pod.Outputs, "metadata", "name")
			assert.Equal(t, "nginx", step1Name.(string))
			step1Namespace, _ := openapi.Pluck(pod.Outputs, "metadata", "namespace")
			assert.Equal(t, "test", step1Namespace.(string))
		},
	})
}
