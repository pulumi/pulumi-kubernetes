// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package ints

import (
	"os"
	"testing"

	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/tests"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/testing/integration"
	"github.com/stretchr/testify/assert"
)

func TestPod(t *testing.T) {
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
			assert.Equal(t, 2, len(stackInfo.Deployment.Resources))

			tests.SortResourcesByURN(stackInfo)

			stackRes := stackInfo.Deployment.Resources[1]
			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

			//
			// Assert pod is successfully created.
			//

			pod := stackInfo.Deployment.Resources[0]
			name, _ := openapi.Pluck(pod.Outputs, "live", "metadata", "name")
			assert.Equal(t, name.(string), "pod-test")

			// Not autonamed.
			_, autonamed := openapi.Pluck(pod.Outputs, "live", "metadata", "annotations",
				"pulumi.com/autonamed")
			assert.False(t, autonamed)

			// Status is "Running"
			phase, _ := openapi.Pluck(pod.Outputs, "live", "status", "phase")
			assert.Equal(t, "Running", phase)

			// Status "Ready" is "True".
			conditions, _ := openapi.Pluck(pod.Outputs, "live", "status", "conditions")
			ready := conditions.([]interface{})[1].(map[string]interface{})
			readyType, _ := ready["type"]
			assert.Equal(t, "Ready", readyType)
			readyStatus, _ := ready["status"]
			assert.Equal(t, "True", readyStatus)

			// Container is called "nginx" and uses image "nginx:1.13-alpine".
			containerStatuses, _ := openapi.Pluck(pod.Outputs, "live", "status", "containerStatuses")
			containerStatus := containerStatuses.([]interface{})[0].(map[string]interface{})
			containerName, _ := containerStatus["name"]
			assert.Equal(t, "nginx", containerName)
			image, _ := containerStatus["image"]
			assert.Equal(t, "nginx:1.13-alpine", image)
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      "step2",
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 2, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					stackRes := stackInfo.Deployment.Resources[1]
					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

					//
					// Assert Pod is deleted before being replaced with the new Pod, running
					// nginx:1.15-alpine.
					//
					// Because the Pod name is supplied in the resource definition, we are forced to delete it
					// before replacing it, as otherwise Kubernetes would complain that it can't have two Pods
					// with the same name.
					//

					pod := stackInfo.Deployment.Resources[0]
					name, _ := openapi.Pluck(pod.Outputs, "live", "metadata", "name")
					assert.Equal(t, name.(string), "pod-test")

					// Not autonamed.
					_, autonamed := openapi.Pluck(pod.Outputs, "live", "metadata", "annotations",
						"pulumi.com/autonamed")
					assert.False(t, autonamed)

					// Status is "Running"
					phase, _ := openapi.Pluck(pod.Outputs, "live", "status", "phase")
					assert.Equal(t, "Running", phase)

					// Status "Ready" is "True".
					conditions, _ := openapi.Pluck(pod.Outputs, "live", "status", "conditions")
					ready := conditions.([]interface{})[1].(map[string]interface{})
					readyType, _ := ready["type"]
					assert.Equal(t, "Ready", readyType)
					readyStatus, _ := ready["status"]
					assert.Equal(t, "True", readyStatus)

					// Container is called "nginx" and uses image "nginx:1.13-alpine".
					containerStatuses, _ := openapi.Pluck(pod.Outputs, "live", "status", "containerStatuses")
					containerStatus := containerStatuses.([]interface{})[0].(map[string]interface{})
					containerName, _ := containerStatus["name"]
					assert.Equal(t, "nginx", containerName)
					image, _ := containerStatus["image"]
					assert.Equal(t, "nginx:1.15-alpine", image)
				},
			},
			{
				Dir:      "step3",
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 2, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					stackRes := stackInfo.Deployment.Resources[1]
					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

					//
					// Assert new Pod is deleted before being replaced with the new Pod, running
					// nginx:1.13-alpine, EVEN WHEN we change the namespace from "" -> "default". This
					// captures the case that we need to delete-before-replace if we're deploying to the same
					// namespace, as measured by canonical name, rather than literal string equality.
					//

					pod := stackInfo.Deployment.Resources[0]
					name, _ := openapi.Pluck(pod.Outputs, "live", "metadata", "name")
					assert.Equal(t, name.(string), "pod-test")

					// Not autonamed.
					_, autonamed := openapi.Pluck(pod.Outputs, "live", "metadata", "annotations",
						"pulumi.com/autonamed")
					assert.False(t, autonamed)

					// Status is "Running"
					phase, _ := openapi.Pluck(pod.Outputs, "live", "status", "phase")
					assert.Equal(t, "Running", phase)

					// Status "Ready" is "True".
					conditions, _ := openapi.Pluck(pod.Outputs, "live", "status", "conditions")
					ready := conditions.([]interface{})[1].(map[string]interface{})
					readyType, _ := ready["type"]
					assert.Equal(t, "Ready", readyType)
					readyStatus, _ := ready["status"]
					assert.Equal(t, "True", readyStatus)

					// Container is called "nginx" and uses image "nginx:1.13-alpine".
					containerStatuses, _ := openapi.Pluck(pod.Outputs, "live", "status", "containerStatuses")
					containerStatus := containerStatuses.([]interface{})[0].(map[string]interface{})
					containerName, _ := containerStatus["name"]
					assert.Equal(t, "nginx", containerName)
					image, _ := containerStatus["image"]
					assert.Equal(t, "nginx:1.13-alpine", image)
				},
			},
		},
	})
}
