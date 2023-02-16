// Copyright 2016-2022, Pulumi Corporation.
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

// nolint:goconst
package test

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/tests/v3"
	"github.com/pulumi/pulumi/pkg/v3/resource/deploy/providers"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var baseOptions = &integration.ProgramTestOptions{
	Verbose: true,
	Dependencies: []string{
		"@pulumi/kubernetes",
	},
	Env: []string{
		"PULUMI_K8S_CLIENT_BURST=200",
		"PULUMI_K8S_CLIENT_QPS=100",
	},
}

func TestAliases(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   filepath.Join("aliases", "step1"),
		Quick: true,
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
		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join("aliases", "step2"),
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
	integration.ProgramTest(t, &test)
}

func TestAutonaming(t *testing.T) {
	var step1Name interface{}
	var step2Name interface{}
	var step3Name interface{}

	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   filepath.Join("autonaming", "step1"),
		Quick: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

			tests.SortResourcesByURN(stackInfo)

			stackRes := stackInfo.Deployment.Resources[3]
			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

			provRes := stackInfo.Deployment.Resources[2]
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			//
			// Assert Pod is successfully given a unique name by Pulumi.
			//

			pod := stackInfo.Deployment.Resources[1]
			assert.Equal(t, "autonaming-test", string(pod.URN.Name()))
			step1Name, _ = openapi.Pluck(pod.Outputs, "metadata", "name")
			assert.True(t, strings.HasPrefix(step1Name.(string), "autonaming-test-"))

			autonamed, _ := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
			assert.Equal(t, "true", autonamed)
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join("autonaming", "step2"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					stackRes := stackInfo.Deployment.Resources[3]
					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

					provRes := stackInfo.Deployment.Resources[2]
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					//
					// Assert Pod was replaced, i.e., destroyed and re-created, with allocating a new name.
					//

					pod := stackInfo.Deployment.Resources[1]
					assert.Equal(t, "autonaming-test", string(pod.URN.Name()))
					step2Name, _ = openapi.Pluck(pod.Outputs, "metadata", "name")
					assert.True(t, strings.HasPrefix(step2Name.(string), "autonaming-test-"))

					autonamed, _ := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
					assert.Equal(t, "true", autonamed)

					assert.NotEqual(t, step1Name, step2Name)

				},
			},
			{
				Dir:      filepath.Join("autonaming", "step3"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					stackRes := stackInfo.Deployment.Resources[3]
					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

					provRes := stackInfo.Deployment.Resources[2]
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					//
					// Assert Pod was NOT replaced, and has the same name, previously allocated by Pulumi.
					//

					pod := stackInfo.Deployment.Resources[1]
					assert.Equal(t, "autonaming-test", string(pod.URN.Name()))
					step3Name, _ = openapi.Pluck(pod.Outputs, "metadata", "name")
					assert.True(t, strings.HasPrefix(step3Name.(string), "autonaming-test-"))

					autonamed, _ := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
					assert.Equal(t, "true", autonamed)

					assert.Equal(t, step2Name, step3Name)
				},
			},
			{
				Dir:      filepath.Join("autonaming", "step4"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					stackRes := stackInfo.Deployment.Resources[3]
					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

					provRes := stackInfo.Deployment.Resources[2]
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					//
					// User has specified their own name for the Pod, so we replace it, and Pulumi does NOT
					// allocate a name on its own.
					//

					pod := stackInfo.Deployment.Resources[1]
					assert.Equal(t, "autonaming-test", string(pod.URN.Name()))
					name, _ := openapi.Pluck(pod.Outputs, "metadata", "name")
					assert.Equal(t, "autonaming-test", name.(string))

					_, autonamed := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
					assert.False(t, autonamed)
				},
			},
		},
	})
	integration.ProgramTest(t, &test)
}

func TestCRDs(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("crds", "step1"),
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
				Dir:      filepath.Join("crds", "step2"),
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
	integration.ProgramTest(t, &test)
}

func TestPod(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   filepath.Join("delete-before-replace", "step1"),
		Quick: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

			tests.SortResourcesByURN(stackInfo)

			stackRes := stackInfo.Deployment.Resources[3]
			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

			provRes := stackInfo.Deployment.Resources[2]
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			//
			// Assert pod is successfully created.
			//

			pod := stackInfo.Deployment.Resources[1]
			name, _ := openapi.Pluck(pod.Outputs, "metadata", "name")
			assert.Equal(t, name.(string), "pod-test")

			// Not autonamed.
			_, autonamed := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
			assert.False(t, autonamed)

			// Status is "Running"
			phase, _ := openapi.Pluck(pod.Outputs, "status", "phase")
			assert.Equal(t, "Running", phase)

			// Status "Ready" is "True".
			conditions, _ := openapi.Pluck(pod.Outputs, "status", "conditions")
			ready := conditions.([]interface{})[1].(map[string]interface{})
			readyType := ready["type"]
			assert.Equal(t, "Ready", readyType)
			readyStatus := ready["status"]
			assert.Equal(t, "True", readyStatus)

			// Container is called "nginx" and uses image "docker.io/library/nginx:1.13-alpine".
			containerStatuses, _ := openapi.Pluck(pod.Outputs, "status", "containerStatuses")
			containerStatus := containerStatuses.([]interface{})[0].(map[string]interface{})
			containerName := containerStatus["name"]
			assert.Equal(t, "nginx", containerName)
			image := containerStatus["image"]
			assert.Equal(t, "docker.io/library/nginx:1.13-alpine", image)
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join("delete-before-replace", "step2"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					stackRes := stackInfo.Deployment.Resources[3]
					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

					provRes := stackInfo.Deployment.Resources[2]
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					//
					// Assert Pod is deleted before being replaced with the new Pod, running
					// nginx:1.15-alpine.
					//
					// Because the Pod name is supplied in the resource definition, we are forced to delete it
					// before replacing it, as otherwise Kubernetes would complain that it can't have two Pods
					// with the same name.
					//

					pod := stackInfo.Deployment.Resources[1]
					name, _ := openapi.Pluck(pod.Outputs, "metadata", "name")
					assert.Equal(t, name.(string), "pod-test")

					// Not autonamed.
					_, autonamed := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
					assert.False(t, autonamed)

					// Status is "Running"
					phase, _ := openapi.Pluck(pod.Outputs, "status", "phase")
					assert.Equal(t, "Running", phase)

					// Status "Ready" is "True".
					conditions, _ := openapi.Pluck(pod.Outputs, "status", "conditions")
					ready := conditions.([]interface{})[1].(map[string]interface{})
					readyType := ready["type"]
					assert.Equal(t, "Ready", readyType)
					readyStatus := ready["status"]
					assert.Equal(t, "True", readyStatus)

					// Container is called "nginx" and uses image "docker.io/library/nginx:1.15-alpine".
					containerStatuses, _ := openapi.Pluck(pod.Outputs, "status", "containerStatuses")
					containerStatus := containerStatuses.([]interface{})[0].(map[string]interface{})
					containerName := containerStatus["name"]
					assert.Equal(t, "nginx", containerName)
					image := containerStatus["image"]
					assert.Equal(t, "docker.io/library/nginx:1.15-alpine", image)
				},
			},
		},
	})
	integration.ProgramTest(t, &test)
}

func TestDeploymentRollout(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   filepath.Join("deployment-rollout", "step1"),
		Quick: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

			tests.SortResourcesByURN(stackInfo)

			appsv1Deploy := stackInfo.Deployment.Resources[0]
			namespace := stackInfo.Deployment.Resources[1]
			provRes := stackInfo.Deployment.Resources[2]
			stackRes := stackInfo.Deployment.Resources[3]

			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())

			//
			// Assert deployment is successfully created.
			//

			name, _ := openapi.Pluck(appsv1Deploy.Outputs, "metadata", "name")
			assert.True(t, strings.Contains(name.(string), "nginx"))
			containers, _ := openapi.Pluck(appsv1Deploy.Outputs, "spec", "template", "spec", "containers")
			containerStatus := containers.([]interface{})[0].(map[string]interface{})
			image := containerStatus["image"]
			assert.Equal(t, image.(string), "nginx")

			assert.True(t, strings.Contains(name.(string), "nginx"))
			containers, _ = openapi.Pluck(appsv1Deploy.Outputs, "spec", "template", "spec", "containers")
			containerStatus = containers.([]interface{})[0].(map[string]interface{})
			image = containerStatus["image"]
			assert.Equal(t, image.(string), "nginx")
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join("deployment-rollout", "step2"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					appsv1Deploy := stackInfo.Deployment.Resources[0]
					namespace := stackInfo.Deployment.Resources[1]
					provRes := stackInfo.Deployment.Resources[2]
					stackRes := stackInfo.Deployment.Resources[3]

					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())

					//
					// Assert deployment is updated successfully.
					//

					name, _ := openapi.Pluck(appsv1Deploy.Outputs, "metadata", "name")
					assert.True(t, strings.Contains(name.(string), "nginx"))
					containers, _ := openapi.Pluck(appsv1Deploy.Outputs, "spec", "template", "spec", "containers")
					containerStatus := containers.([]interface{})[0].(map[string]interface{})
					image := containerStatus["image"]
					assert.Equal(t, image.(string), "nginx:stable")
				},
			},
		},
	})
	integration.ProgramTest(t, &test)
}

func TestEmptyArray(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir: filepath.Join("empty-array", "step1"),
		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join("empty-array", "step2"),
				Additive: true,
			},
		},
	})
	integration.ProgramTest(t, &test)
}

func TestGet(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("get", "step1"),
		Quick:                true,
		ExpectRefreshChanges: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 6, len(stackInfo.Deployment.Resources))

			tests.SortResourcesByURN(stackInfo)

			stackRes := stackInfo.Deployment.Resources[5]
			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

			provRes := stackInfo.Deployment.Resources[4]
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			//
			// Assert we can use .get to retrieve the kube-api Service.
			//

			service := stackInfo.Deployment.Resources[2]
			assert.Equal(t, "kube-api", string(service.URN.Name()))
			step1Name, _ := openapi.Pluck(service.Outputs, "metadata", "name")
			assert.Equal(t, "kubernetes", step1Name.(string))

			//
			// Assert that CRD and CR exist
			//

			crd := stackInfo.Deployment.Resources[0]
			assert.Equal(t, "crontab", string(crd.URN.Name()))

			ct1 := stackInfo.Deployment.Resources[3]
			assert.Equal(t, "my-new-cron-object", string(ct1.URN.Name()))

		},
		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join("get", "step2"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 7, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					stackRes := stackInfo.Deployment.Resources[6]
					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

					provRes := stackInfo.Deployment.Resources[5]
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					//
					// Assert we can use .get to retrieve CRDs.
					//

					ct2 := stackInfo.Deployment.Resources[4]
					assert.Equal(t, "my-new-cron-object-get", string(ct2.URN.Name()))
					image, _ := openapi.Pluck(ct2.Outputs, "spec", "image")
					assert.Equal(t, "my-awesome-cron-image", image.(string))
				},
			},
		},
	})
	integration.ProgramTest(t, &test)
}

func TestIstio(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:         filepath.Join("istio", "step1"),
		Quick:       true,
		SkipRefresh: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			frontend := stackInfo.Outputs["frontendIp"].(string)

			// Retry the GET on the Istio gateway repeatedly. Istio doesn't publish `.status` on any
			// of its CRDs, so this is as reliable as we can be right now.
			for i := 1; i < 10; i++ {
				req, err := http.Get(fmt.Sprintf("http://%s", frontend))
				if err != nil {
					fmt.Printf("Request to Istio gateway failed: %v\n", err)
					time.Sleep(time.Second * 10)
				} else if req.StatusCode == 200 {
					return
				}
			}

			assert.Fail(t, "Maximum Istio gateway request retries exceeded")
		},
	})
	integration.ProgramTest(t, &test)
}

func TestKustomize(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   filepath.Join("kustomize", "step1"),
		Quick: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)

			// Assert that we've retrieved the kustomizations and provisioned them.
			assert.Equal(t, 12, len(stackInfo.Deployment.Resources))
		},
	})
	integration.ProgramTest(t, &test)
}

func TestNamespace(t *testing.T) {
	var nmPodName, defaultPodName string
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   filepath.Join("namespace", "step1"),
		Quick: true,
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
				Dir:      filepath.Join("namespace", "step2"),
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
					assert.True(t, namespaceLabels.(map[string]interface{})["hello"] == "world")
				},
			},
			{
				Dir:      filepath.Join("namespace", "step3"),
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
	integration.ProgramTest(t, &test)
}

// FIXME(levi): Figure out why this test is flaky, and re-enable it in CI.
//  https://github.com/pulumi/pulumi-kubernetes/issues/1016
//func TestPerformance(t *testing.T) {
//	test := baseOptions.With(integration.ProgramTestOptions{
//		Dir:                  filepath.Join("performance", "step1"),
//		ExpectRefreshChanges: true, // The Mutating and Validating webhooks update on refresh.
//	})
//	integration.ProgramTest(t, &test)
//}

func TestProvider(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   filepath.Join("provider", "step1"),
		Quick: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 10, len(stackInfo.Deployment.Resources))

			tests.SortResourcesByURN(stackInfo)

			stackRes := stackInfo.Deployment.Resources[9]
			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

			k8sPathProvider := stackInfo.Deployment.Resources[8]
			assert.True(t, providers.IsProviderType(k8sPathProvider.URN.Type()))

			k8sContentsProvider := stackInfo.Deployment.Resources[7]
			assert.True(t, providers.IsProviderType(k8sContentsProvider.URN.Type()))

			defaultProvider := stackInfo.Deployment.Resources[6]
			assert.True(t, providers.IsProviderType(defaultProvider.URN.Type()))

			// Assert the provider default Namespace (ns1) was created
			ns1 := stackInfo.Deployment.Resources[0]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), ns1.URN.Type())
			providerNsName, _ := openapi.Pluck(ns1.Outputs, "metadata", "name")

			// Assert the ns2 Namespace was created
			ns2 := stackInfo.Deployment.Resources[1]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), ns2.URN.Type())
			ns2Name, _ := openapi.Pluck(ns2.Outputs, "metadata", "name")

			// Assert the other Namespace was created and doesn't use the provider default.
			otherNamespace := stackInfo.Deployment.Resources[2]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), otherNamespace.URN.Type())
			nsName, _ := openapi.Pluck(otherNamespace.Outputs, "metadata", "name")
			assert.NotEqual(t, nsName.(string), providerNsName.(string))

			// Assert the first Pod was created in the provider default namespace.
			pod1 := stackInfo.Deployment.Resources[4]
			assert.Equal(t, "nginx1", string(pod1.URN.Name()))
			podNamespace1, _ := openapi.Pluck(pod1.Outputs, "metadata", "namespace")
			assert.Equal(t, providerNsName.(string), podNamespace1.(string))

			// Assert the second Pod was created in the provider default namespace.
			pod2 := stackInfo.Deployment.Resources[5]
			assert.Equal(t, "nginx2", string(pod2.URN.Name()))
			podNamespace2, _ := openapi.Pluck(pod2.Outputs, "metadata", "namespace")
			assert.Equal(t, providerNsName.(string), podNamespace2.(string))

			// Assert the Pod was created in the specified namespace rather than the provider default namespace.
			namespacedPod := stackInfo.Deployment.Resources[3]
			assert.Equal(t, "namespaced-nginx", string(namespacedPod.URN.Name()))
			namespacedPodNamespace, _ := openapi.Pluck(namespacedPod.Outputs, "metadata", "namespace")
			assert.NotEqual(t, providerNsName.(string), namespacedPodNamespace.(string))
			assert.Equal(t, ns2Name.(string), namespacedPodNamespace.(string))
		},
	})
	integration.ProgramTest(t, &test)
}

func TestQuery(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:       filepath.Join("query", "step1"),
		Quick:     true,
		StackName: "query-test-c186bcc3-1572-44d8-b7d5-1028853682c3", // Chosen by fair dice roll. Guaranteed to be random.
		CloudURL:  "file://~",                                        // Required; we hard-code the stack name
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

			tests.SortResourcesByURN(stackInfo)

			stackRes := stackInfo.Deployment.Resources[3]
			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

			provRes := stackInfo.Deployment.Resources[2]
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			//
			// Assert Pod is successfully given a unique name by Pulumi.
			//

			pod := stackInfo.Deployment.Resources[1]
			assert.Equal(t, "query-test", string(pod.URN.Name()))
		},
		EditDirs: []integration.EditDir{
			{
				Dir:       filepath.Join("query", "step2"),
				Additive:  true,
				QueryMode: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					//
					// Verify no resources were deleted.
					//
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					stackRes := stackInfo.Deployment.Resources[3]
					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

					provRes := stackInfo.Deployment.Resources[2]
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					//
					// If we pass this point, the query did NOT throw an error, and is therefore
					// successful.
					//
				},
			},
		},
	})
	integration.ProgramTest(t, &test)
}

func TestRenderYAML(t *testing.T) {
	// Create a temporary directory to hold rendered YAML manifests.
	dir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	test := baseOptions.With(integration.ProgramTestOptions{
		Config: map[string]string{
			"renderDir": dir,
		},
		Dir:   filepath.Join("render-yaml", "step1"),
		Quick: true,
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
	})
	integration.ProgramTest(t, &test)
}

func TestReplaceUnready(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("replace-unready", "step1"),
		Quick:                true,
		ExpectFailure:        true, // The Job is intended to fail.
		ExpectRefreshChanges: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 3, len(stackInfo.Deployment.Resources))

			tests.SortResourcesByURN(stackInfo)

			job := stackInfo.Deployment.Resources[0]
			provRes := stackInfo.Deployment.Resources[1]
			stackRes := stackInfo.Deployment.Resources[2]

			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			assert.Equal(t, tokens.Type("kubernetes:batch/v1:Job"), job.URN.Type())
		},
		EditDirs: []integration.EditDir{
			{
				Dir:           filepath.Join("replace-unready", "step2"),
				Additive:      true,
				ExpectFailure: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 3, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					job := stackInfo.Deployment.Resources[0]
					provRes := stackInfo.Deployment.Resources[1]
					stackRes := stackInfo.Deployment.Resources[2]

					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					assert.Equal(t, tokens.Type("kubernetes:batch/v1:Job"), job.URN.Type())

					// Check the event stream for a preview showing that the Job will be updated.
					for _, e := range stackInfo.Events {
						if e.ResourcePreEvent != nil && e.ResourcePreEvent.Metadata.Type == "kubernetes:batch/v1:Job" {
							assert.Equal(t, e.ResourcePreEvent.Metadata.Op, apitype.OpUpdate)
						}
					}
				},
			},
			{
				Dir:           filepath.Join("replace-unready", "step3"),
				Additive:      true,
				ExpectFailure: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					// The stack has an extra Job now from the failed update in step2.
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 4, len(stackInfo.Deployment.Resources))

					tests.SortResourcesByURN(stackInfo)

					job := stackInfo.Deployment.Resources[0]
					jobOld := stackInfo.Deployment.Resources[1]
					provRes := stackInfo.Deployment.Resources[2]
					stackRes := stackInfo.Deployment.Resources[3]

					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					assert.Equal(t, tokens.Type("kubernetes:batch/v1:Job"), job.URN.Type())
					assert.Equal(t, tokens.Type("kubernetes:batch/v1:Job"), jobOld.URN.Type())

					// Check the event stream for a preview showing that the Job will be replaced.
					for _, e := range stackInfo.Events {
						if e.ResourcePreEvent != nil && e.ResourcePreEvent.Metadata.Type == "kubernetes:batch/v1:Job" {
							assert.Equal(t, e.ResourcePreEvent.Metadata.Op, apitype.OpCreateReplacement)
						}
					}
				},
			},
		},
	})
	integration.ProgramTest(t, &test)
}

func TestRetry(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   filepath.Join("retry", "step1"),
		Quick: true,
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
	integration.ProgramTest(t, &test)
}

func TestSecrets(t *testing.T) {
	secretMessage := "secret message for testing"

	test := baseOptions.With(integration.ProgramTestOptions{
		Dir: filepath.Join("secrets", "step1"),
		Config: map[string]string{
			"message": secretMessage,
		},
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExpectRefreshChanges: true,
		Quick:                true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			state, err := json.Marshal(stackInfo.Deployment)
			assert.NoError(t, err)

			assert.NotContains(t, string(state), secretMessage)

			// The program converts the secret message to base64, to make a ConfigMap from it, so the state
			// should also not contain the base64 encoding of secret message.
			assert.NotContains(t, string(state), b64.StdEncoding.EncodeToString([]byte(secretMessage)))
		},
	})
	integration.ProgramTest(t, &test)
}

func TestServerSideApply(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("server-side-apply", "step1"),
		ExpectRefreshChanges: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			// Validate patched Namespace
			nsPatched := stackInfo.Outputs["nsPatched"].(map[string]interface{})
			fooV, ok, err := unstructured.NestedString(nsPatched, "metadata", "labels", "foo")
			assert.True(t, ok)
			assert.NoError(t, err)
			assert.Equal(t, "foo", fooV)

			// Validate patched CustomResource
			crPatched := stackInfo.Outputs["crPatched"].(map[string]interface{})
			fooV, ok, err = unstructured.NestedString(crPatched, "metadata", "labels", "foo")
			assert.True(t, ok)
			assert.NoError(t, err)
			assert.Equal(t, "foo", fooV)

			for _, res := range stackInfo.Deployment.Resources {
				// Validate that the last-applied-configuration annotation is not present on SSA resources.
				annotations, ok, err := unstructured.NestedStringMap(res.Outputs, "metadata", "labels")
				assert.NoError(t, err)
				if ok {
					assert.NotContains(t, annotations, "kubectl.kubernetes.io/last-applied-configuration")
				}

				// Validate that the managed-by label is not present on SSA resources.
				labels, ok, err := unstructured.NestedStringMap(res.Outputs, "metadata", "labels")
				assert.NoError(t, err)
				if ok {
					assert.NotContains(t, labels, "app.kubernetes.io/managed-by")
				}

				if res.Type == "kubernetes:core/v1:ConfigMap" {
					dataV, ok, err := unstructured.NestedString(res.Outputs, "data", "foo")
					assert.True(t, ok)
					assert.NoError(t, err)
					assert.Equal(t, "bar", dataV)
				}
			}
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join("server-side-apply", "step2"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					// Validate patched Deployment
					deploymentPatched := stackInfo.Outputs["deploymentPatched"].(map[string]interface{})
					containersV, ok, err := unstructured.NestedSlice(
						deploymentPatched, "spec", "template", "spec", "containers")
					assert.True(t, ok)
					assert.NoError(t, err)
					assert.Len(t, containersV, 1)
					limitsV, ok, err := unstructured.NestedMap(
						containersV[0].(map[string]interface{}), "resources", "limits")
					assert.True(t, ok)
					assert.NoError(t, err)
					assert.Contains(t, limitsV, "memory")

					for _, res := range stackInfo.Deployment.Resources {
						if res.Type == "kubernetes:core/v1:ConfigMap" {
							dataV, ok, err := unstructured.NestedString(res.Outputs, "data", "foo")
							assert.True(t, ok)
							assert.NoError(t, err)
							assert.Equal(t, "baz", dataV)
						}
					}
				},
			},
			{
				Dir:      filepath.Join("server-side-apply", "step3"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					for _, res := range stackInfo.Deployment.Resources {
						if res.Type == "kubernetes:core/v1:ConfigMap" {
							dataV, ok, err := unstructured.NestedString(res.Outputs, "data", "foo")
							assert.True(t, ok)
							assert.NoError(t, err)
							assert.Equal(t, "baz", dataV) // Data should be unchanged from step2.
						}
					}
				},
			},
		},
	})
	integration.ProgramTest(t, &test)
}

func TestYAMLURL(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   filepath.Join("yaml-url", "step1"),
		Quick: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)

			// Assert that we've retrieved the YAML from the URL and provisioned them.
			assert.Equal(t, 18, len(stackInfo.Deployment.Resources))
		},
	})
	integration.ProgramTest(t, &test)
}

func TestReplaceDaemonSet(t *testing.T) {
	daemonSetName := ""
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:           filepath.Join("replace-daemonset", "step1"),
		Quick:         true,
		ExpectFailure: false,
		SkipRefresh:   true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 3, len(stackInfo.Deployment.Resources))

			// Save the DaemonSet name to compare it in the step2
			daemonSetName = stackInfo.Outputs["name"].(string)

			// Assert that the DaemonSet was created
			assert.True(t, strings.HasPrefix(stackInfo.Outputs["name"].(string), "test-replacement-"))
		},
		EditDirs: []integration.EditDir{
			{
				Dir:           filepath.Join("replace-daemonset", "step2"),
				Additive:      true,
				ExpectFailure: false,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 3, len(stackInfo.Deployment.Resources))

					newDaemonSetName := stackInfo.Outputs["name"].(string)

					// Assert that the DaemonSet still exists
					assert.True(t, strings.HasPrefix(newDaemonSetName, "test-replacement-"))

					// DaemonSet should have a different name as it was replaced
					assert.True(t, daemonSetName != newDaemonSetName)
				},
			},
		},
	})
	integration.ProgramTest(t, &test)
}

func TestServiceAccountTokenSecret(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:           filepath.Join("service-account-token-secret", "step1"),
		Quick:         true,
		ExpectFailure: false,
		SkipRefresh:   true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			_, err := json.Marshal(stackInfo.Deployment)
			assert.NoError(t, err)

			secretData := stackInfo.Outputs["data"].(map[string]interface{})

			assert.Contains(t, secretData, "ca.crt")
			assert.Contains(t, secretData, "token")
		},
	})
	integration.ProgramTest(t, &test)
}
