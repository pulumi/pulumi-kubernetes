// Copyright 2016-2023, Pulumi Corporation.
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
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	. "github.com/onsi/gomega/gstruct"
	. "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/gomega"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/tests/v4"
	pulumirpctesting "github.com/pulumi/pulumi-kubernetes/tests/v4/pulumirpc"
	"github.com/pulumi/pulumi/pkg/v3/engine"
	"github.com/pulumi/pulumi/pkg/v3/resource/deploy/providers"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/fsutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/tools/clientcmd"
)

var baseOptions = &integration.ProgramTestOptions{
	Verbose: true,
	Dependencies: []string{
		"@pulumi/kubernetes",
	},
	PostPrepareProject: func(p *engine.Projinfo) error {
		return fsutil.CopyFile(filepath.Join(p.Root, "testdata"), filepath.Join("..", "..", "testdata"), nil)
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
			assert.Equal(t, "alias-test", deployment.URN.Name())
			assert.Equal(t, "kubernetes:apps/v1:Deployment", string(deployment.Type))
			containers, _ := openapi.Pluck(deployment.Outputs, "spec", "template", "spec", "containers")
			containerStatus := containers.([]any)[0].(map[string]any)
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
					assert.Equal(t, "alias-test", deployment.URN.Name())
					assert.Equal(t, "kubernetes:apps/v1:Deployment", string(deployment.Type))
					containers, _ := openapi.Pluck(deployment.Outputs, "spec", "template", "spec", "containers")
					containerStatus := containers.([]any)[0].(map[string]any)
					image := containerStatus["image"]
					assert.Equal(t, image.(string), "nginx:1.15")
				},
			},
		},
	})
	integration.ProgramTest(t, &test)
}

func TestAutonaming(t *testing.T) {
	var step1Name any
	var step2Name any
	var step3Name any
	var step4Name any

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
			assert.Equal(t, "autonaming-test", pod.URN.Name())
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
					assert.Equal(t, "autonaming-test", pod.URN.Name())
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
					assert.Equal(t, "autonaming-test", pod.URN.Name())
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
					// Assert Pod was NOT replaced, and has the same name, previously allocated by Pulumi.
					//

					pod := stackInfo.Deployment.Resources[1]
					assert.Equal(t, "autonaming-test", pod.URN.Name())
					step4Name, _ = openapi.Pluck(pod.Outputs, "metadata", "name")
					assert.True(t, strings.HasPrefix(step4Name.(string), "autonaming-test-"))

					autonamed, _ := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
					assert.Equal(t, "true", autonamed)

					assert.Equal(t, step3Name, step4Name)
				},
			},
			{
				Dir:      filepath.Join("autonaming", "step5"),
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
					assert.Equal(t, "autonaming-test", pod.URN.Name())
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

func TestGenerateName(t *testing.T) {
	var pt *integration.ProgramTester
	var step1Name any
	var step2Name any
	var step3Name any
	var step4Name any
	var step5Name any
	var step6Name any

	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("generatename", "step1"),
		Quick:                false,
		SkipRefresh:          false,
		ExpectRefreshChanges: false,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)

			//
			// Assert pod is successfully given a unique name by Kubernetes.
			//
			pod := tests.SearchResourcesByName(stackInfo, "", "kubernetes:core/v1:Pod", "generatename-test")
			step1Name, _ = openapi.Pluck(pod.Outputs, "metadata", "name")
			assert.True(t, strings.HasPrefix(step1Name.(string), "generatename-test-"))
			generateName, _ := openapi.Pluck(pod.Outputs, "metadata", "generateName")
			assert.Equal(t, "generatename-test-", generateName.(string))
			_, autonamed := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
			assert.False(t, autonamed)
		},
		Config: map[string]string{},

		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join("generatename", "step2"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)

					//
					// Assert pod was NOT replaced, and has the same name, previously allocated by Kubernetes.
					//
					pod := tests.SearchResourcesByName(stackInfo, "", "kubernetes:core/v1:Pod", "generatename-test")
					step2Name, _ = openapi.Pluck(pod.Outputs, "metadata", "name")
					assert.Equal(t, step1Name, step2Name)
					generateName, _ := openapi.Pluck(pod.Outputs, "metadata", "generateName")
					assert.Equal(t, "generatename-test-modified-", generateName.(string))
					_, autonamed := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
					assert.False(t, autonamed)
				},
			},
			{
				Dir:      filepath.Join("generatename", "step3"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)

					//
					// Assert pod was replaced, i.e., destroyed and re-created, with allocating a new name.
					//
					pod := tests.SearchResourcesByName(stackInfo, "", "kubernetes:core/v1:Pod", "generatename-test")
					step3Name, _ = openapi.Pluck(pod.Outputs, "metadata", "name")
					assert.NotEqual(t, step2Name, step3Name)
					assert.True(t, strings.HasPrefix(step3Name.(string), "generatename-test-modified-"))
					_, autonamed := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
					assert.False(t, autonamed)
				},
			},
			{
				Dir:      filepath.Join("generatename", "step4"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)

					//
					// Assert pod was NOT replaced, and has the same name, previously allocated by Kubernetes.
					//
					pod := tests.SearchResourcesByName(stackInfo, "", "kubernetes:core/v1:Pod", "generatename-test")
					step4Name, _ = openapi.Pluck(pod.Outputs, "metadata", "name")
					assert.Equal(t, step3Name, step4Name)
					_, autonamed := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
					assert.False(t, autonamed)

					// Update the configuration for subsequent steps.
					require.NoError(t,
						pt.RunPulumiCommand("config", "set", "podName", step4Name.(string)),
						"failed to set podName config")
				},
			},
			{
				Dir:             filepath.Join("generatename", "step5"),
				Additive:        true,
				ExpectNoChanges: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)

					//
					// User has explicitly set the name to the previously-generated name (maybe for clarity),
					// and Pulumi does NOT replace the pod.
					//
					pod := tests.SearchResourcesByName(stackInfo, "", "kubernetes:core/v1:Pod", "generatename-test")
					step5Name, _ = openapi.Pluck(pod.Outputs, "metadata", "name")
					assert.Equal(t, step4Name, step5Name)
					_, autonamed := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
					assert.False(t, autonamed)
				},
			},
			{
				Dir:      filepath.Join("generatename", "step6"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)

					//
					// User has specified their own name for the Pod, so we replace it, and Pulumi/Kubernetes does NOT
					// allocate a name on its own.
					//
					pod := tests.SearchResourcesByName(stackInfo, "", "kubernetes:core/v1:Pod", "generatename-test")
					step6Name, _ = openapi.Pluck(pod.Outputs, "metadata", "name")
					assert.NotEqual(t, step5Name, step6Name)
					assert.Equal(t, "generatename-test", step6Name.(string))
					_, autonamed := openapi.Pluck(pod.Outputs, "metadata", "annotations", "pulumi.com/autonamed")
					assert.False(t, autonamed)
				},
			},
		},
	})
	pt = integration.ProgramTestManualLifeCycle(t, &test)
	err := pt.TestLifeCycleInitAndDestroy()
	if !errors.Is(err, integration.ErrTestFailed) {
		assert.NoError(t, err)
	}
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
			outputs := stackInfo.Outputs

			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())

			//
			// Assert that CRD and CR exist
			//

			assert.Equal(t, "foobar", crd.URN.Name())
			assert.Equal(t, "my-new-foobar-object", ct1.URN.Name())

			// Assert that we can reference the x_kubernetes_preserve_unknown_fields field as we should correctly normalize the live object.
			assert.NotNil(t, outputs["preserveUnknownFields"], "preserveUnknownFields should be present")
			assert.True(t, outputs["preserveUnknownFields"].(bool), "preserveUnknownFields should be true")

			// Verify with kubectl that the CRD has `x-kubernetes-*` fields set correctly.
			output, err := tests.Kubectl("get crd foobars.stable.example.com -o json")
			require.NoError(t, err)
			assert.Contains(t, string(output), `"x-kubernetes-preserve-unknown-fields": true`)
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join("crds", "step2"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					outputs := stackInfo.Outputs

					//
					// Assert that the CR was gotten.
					//
					ct1ref := tests.SearchResourcesByName(stackInfo, "", "kubernetes:stable.example.com/v1:FooBar", "my-new-foobar-object-ref")
					assert.NotNil(t, ct1ref)

					// Assert that the x_kubernetes_preserve_unknown_fields field is now nil, as we remove this in step 2.
					assert.Nil(t, outputs["preserveUnknownFields"], "preserveUnknownFields should be present")

					// Verify with kubectl that the CRD does not have `x-kubernetes-*` fields.
					output, err := tests.Kubectl("get crd foobars.stable.example.com -o json")
					require.NoError(t, err)
					assert.NotContains(t, string(output), `"x-kubernetes-preserve-unknown-fields": true`)
				},
			},
			{
				Dir:      filepath.Join("crds", "step3"),
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
	getCondition := func(conditions []any, conditionType string) map[string]any {
		// Order of conditions is not guaranteed, so we need to search for the condition of the given type.
		for _, condition := range conditions {
			conditionMap, ok := condition.(map[string]any)
			require.True(t, ok, "condition items should be maps")
			if conditionMap["type"] == conditionType {
				return conditionMap
			}
		}
		t.Fatalf("condition of type %s not found", conditionType)
		return nil
	}

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
			ready := getCondition(conditions.([]any), "Ready")
			readyStatus := ready["status"]
			assert.Equal(t, "True", readyStatus)

			// Container is called "nginx" and uses image "docker.io/library/nginx:1.13-alpine".
			containerStatuses, _ := openapi.Pluck(pod.Outputs, "status", "containerStatuses")
			containerStatus := containerStatuses.([]any)[0].(map[string]any)
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
					ready := getCondition(conditions.([]any), "Ready")
					readyStatus := ready["status"]
					assert.Equal(t, "True", readyStatus)

					// Container is called "nginx" and uses image "docker.io/library/nginx:1.15-alpine".
					containerStatuses, _ := openapi.Pluck(pod.Outputs, "status", "containerStatuses")
					containerStatus := containerStatuses.([]any)[0].(map[string]any)
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
			containerStatus := containers.([]any)[0].(map[string]any)
			image := containerStatus["image"]
			assert.Equal(t, image.(string), "nginx")

			assert.True(t, strings.Contains(name.(string), "nginx"))
			containers, _ = openapi.Pluck(appsv1Deploy.Outputs, "spec", "template", "spec", "containers")
			containerStatus = containers.([]any)[0].(map[string]any)
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
					containerStatus := containers.([]any)[0].(map[string]any)
					image := containerStatus["image"]
					assert.Equal(t, image.(string), "nginx:stable")
				},
			},
			{
				// This is a deployment spec update that causes a no-op replica set update.
				Dir:      filepath.Join("deployment-rollout", "step3"),
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
					containerStatus := containers.([]any)[0].(map[string]any)
					image := containerStatus["image"]
					assert.Equal(t, image.(string), "nginx:stable")
				},
			},
			{
				// This is a deployment spec update that should cause a failure on await.
				// https://github.com/pulumi/pulumi-kubernetes/issues/2941
				Dir:           filepath.Join("deployment-rollout", "step4"),
				Additive:      true,
				ExpectFailure: true,
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
			assert.Equal(t, 7, len(stackInfo.Deployment.Resources))

			//
			// Assert we can use .get to retrieve the kube-api Service.
			//

			service := stackInfo.Outputs["svc"].(map[string]any)
			svcURN, _ := openapi.Pluck(service, "urn")
			assert.Containsf(t, svcURN, "kube-api", "urn missing expected name")
			svcName, _ := openapi.Pluck(service, "metadata", "name")
			assert.Equalf(t, "kubernetes", svcName, "unexpected service name")

			//
			// Assert that the uninitialized Service exists
			//

			awaitSvc := stackInfo.Outputs["awaitSvc"].(map[string]any)
			awaitSvcName, _ := openapi.Pluck(awaitSvc, "metadata", "name")
			assert.Equalf(t, "test", awaitSvcName, "unexpected service name")
			awaitSvcAnnotation, ok := openapi.Pluck(awaitSvc, "metadata", "annotations", "pulumi.com/skipAwait")
			assert.Truef(t, ok, "failed to find skipAwait annotation")
			assert.Equalf(t, "true", awaitSvcAnnotation, "expected annotation to be true")

			//
			// Assert that CRD and CR exist
			//

			crd := stackInfo.Outputs["ct"].(map[string]any)
			crdURN, _ := openapi.Pluck(crd, "urn")
			assert.Containsf(t, crdURN, "crontab", "urn missing expected name")

			cr := stackInfo.Outputs["cr"].(map[string]any)
			crURN, _ := openapi.Pluck(cr, "urn")
			assert.Containsf(t, crURN, "my-new-cron-object", "urn missing expected name")
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join("get", "step2"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					assert.NotNil(t, stackInfo.Deployment)
					assert.Equal(t, 10, len(stackInfo.Deployment.Resources))

					//
					// Assert we can use .get to retrieve the kube-api Service.
					//

					service := stackInfo.Outputs["svc"].(map[string]any)
					svcURN, _ := openapi.Pluck(service, "urn")
					assert.Containsf(t, svcURN, "kube-api", "urn missing expected name")
					svcName, _ := openapi.Pluck(service, "metadata", "name")
					assert.Equalf(t, "kubernetes", svcName, "unexpected service name")

					//
					// Assert that the uninitialized Service exists
					//

					awaitSvc := stackInfo.Outputs["awaitSvc"].(map[string]any)
					awaitSvcName, _ := openapi.Pluck(awaitSvc, "metadata", "name")
					assert.Equalf(t, "test", awaitSvcName, "unexpected service name")
					awaitSvcAnnotation, ok := openapi.Pluck(awaitSvc, "metadata", "annotations", "pulumi.com/skipAwait")
					assert.Truef(t, ok, "failed to find skipAwait annotation")
					assert.Equalf(t, "true", awaitSvcAnnotation, "expected annotation to be true")

					//
					// Assert we can use .get to retrieve a Service that would fail await logic.
					//

					awaitSvcGet := stackInfo.Outputs["awaitSvcGet"].(map[string]any)
					awaitSvcGetURN, _ := openapi.Pluck(awaitSvcGet, "urn")
					assert.Containsf(t, awaitSvcGetURN, "await", "urn missing expected name")

					//
					// Assert we can use an output from a Service that would fail await logic.
					//

					cm := stackInfo.Outputs["cm"].(map[string]any)
					cmURN, _ := openapi.Pluck(cm, "urn")
					assert.Containsf(t, cmURN, "svc-test", "urn missing expected name")
					clusterIP, _ := openapi.Pluck(cm, "data", "key")
					assert.NotEmptyf(t, clusterIP, "clusterIP should be set")

					//
					// Assert that CRD and CR exist
					//

					crd := stackInfo.Outputs["ct"].(map[string]any)
					crdURN, _ := openapi.Pluck(crd, "urn")
					assert.Containsf(t, crdURN, "crontab", "urn missing expected name")

					cr := stackInfo.Outputs["cr"].(map[string]any)
					crURN, _ := openapi.Pluck(cr, "urn")
					assert.Containsf(t, crURN, "my-new-cron-object", "urn missing expected name")

					//
					// Assert we can use .get to retrieve CRDs.
					//

					crGet := stackInfo.Outputs["crGet"].(map[string]any)
					crGetURN, _ := openapi.Pluck(crGet, "urn")
					assert.Containsf(t, crGetURN, "my-new-cron-object-get", "urn missing expected name")
				},
			},
		},
	})
	integration.ProgramTest(t, &test)
}

func TestIstio(t *testing.T) {
	tests.SkipIfShort(t, "test provisions a load balancer and requires a cloud provider cluster to run")
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

// TestKustomizeHelmChart verifies the helmChart plugin support for Kustomize. This requires the `helm` binary to be
// on the system PATH to succeed.
// Example based on https://github.com/kubernetes-sigs/kustomize/blob/v3.3.1/examples/chart.md
func TestKustomizeHelmChart(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir: filepath.Join("kustomizeHelmChart", "step1"),
		// Just test that the plugin integration is working with preview
		SkipUpdate:             true,
		SkipRefresh:            true,
		SkipEmptyPreviewUpdate: true,
		SkipExportImport:       true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		Env: []string{
			"PULUMI_K8S_KUSTOMIZE_HELM=true", // This experimental feature is currently gated behind a feature flag.
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
					assert.True(t, namespaceLabels.(map[string]any)["hello"] == "world")
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
			assert.Equal(t, "nginx1", pod1.URN.Name())
			podNamespace1, _ := openapi.Pluck(pod1.Outputs, "metadata", "namespace")
			assert.Equal(t, providerNsName.(string), podNamespace1.(string))

			// Assert the second Pod was created in the provider default namespace.
			pod2 := stackInfo.Deployment.Resources[5]
			assert.Equal(t, "nginx2", pod2.URN.Name())
			podNamespace2, _ := openapi.Pluck(pod2.Outputs, "metadata", "namespace")
			assert.Equal(t, providerNsName.(string), podNamespace2.(string))

			// Assert the Pod was created in the specified namespace rather than the provider default namespace.
			namespacedPod := stackInfo.Deployment.Resources[3]
			assert.Equal(t, "namespaced-nginx", namespacedPod.URN.Name())
			namespacedPodNamespace, _ := openapi.Pluck(namespacedPod.Outputs, "metadata", "namespace")
			assert.NotEqual(t, providerNsName.(string), namespacedPodNamespace.(string))
			assert.Equal(t, ns2Name.(string), namespacedPodNamespace.(string))
		},
	})
	integration.ProgramTest(t, &test)
}

// TestReadonlyMetadata tests the behavior of read-only metadata fields.
func TestReadonlyMetadata(t *testing.T) {

	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   filepath.Join("metadata", "step1"),
		Quick: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			// Verify that the resourceVersion (a read-only property) is available as an output.
			resourceVersion, ok := stackInfo.Outputs["resourceVersion"]
			assert.Truef(t, ok, "missing expected output \"resourceVersion\"")
			assert.NotEmptyf(t, resourceVersion, "resourceVersion is empty")
			assert.NotEqual(t, "invalid-step1", resourceVersion)
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join("metadata", "step2"),
				Additive: true,
				// Verify that changes to some read-only values are ignored
				ExpectNoChanges: true,
			},
		},
	})
	integration.ProgramTest(t, &test)
}

func TestRenderYAML(t *testing.T) {
	// Create a temporary directory to hold rendered YAML manifests.
	dir, err := os.MkdirTemp("", "render-yaml-test")
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
			files, err := os.ReadDir(dir)
			assert.NoError(t, err)
			assert.Equal(t, len(files), 2)

			// Verify that CRD manifest directory was created.
			files, err = os.ReadDir(filepath.Join(dir, "0-crd"))
			assert.NoError(t, err)
			assert.Equal(t, len(files), 0)

			// Verify that manifest directory was created.
			files, err = os.ReadDir(filepath.Join(dir, "1-manifest"))
			assert.NoError(t, err)
			assert.Equal(t, len(files), 2)
		},
		EditDirs: []integration.EditDir{
			{
				// Change some fields.
				Dir:      filepath.Join("render-yaml", "step2"),
				Additive: true,
			},
			{
				// Ensure updates do not cause a spurrious diff when re-running `pulumi up`.
				Dir:             filepath.Join("render-yaml", "step2"),
				Additive:        true,
				ExpectNoChanges: true,
			},
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
			assert.Equal(t, "nginx", pod.URN.Name())
			step1Name, _ := openapi.Pluck(pod.Outputs, "metadata", "name")
			assert.Equal(t, "nginx", step1Name.(string))
			step1PodNamespace, _ := openapi.Pluck(pod.Outputs, "metadata", "namespace")
			assert.Equal(t, namespace.ID.String(), step1PodNamespace.(string))
		},
	})
	integration.ProgramTest(t, &test)
}

func TestSecretFieldsMarkedAsSecrets(t *testing.T) {
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
		Quick: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			state, err := json.Marshal(stackInfo.Deployment)
			assert.NoError(t, err)

			ssStringDataData, ok := stackInfo.Outputs["ssStringDataData"]
			assert.Truef(t, ok, "missing expected output \"ssStringDataData\"")

			ssStringDataStringData, ok := stackInfo.Outputs["ssStringDataStringData"]
			assert.Truef(t, ok, "missing expected output \"ssStringDataStringData\"")

			cgSecretStringData, ok := stackInfo.Outputs["cgSecretStringData"]
			assert.Truef(t, ok, "missing expected output \"cgSecretStringData\"")

			assert.NotEmptyf(t, ssStringDataData, "data field is empty")
			assert.NotEmptyf(t, ssStringDataStringData, "stringData field is empty")
			assert.NotEmptyf(t, cgSecretStringData, "cgSecretStringData field is empty")

			assert.NotContains(t, string(state), secretMessage)

			// The program converts the secret message to base64, to make a ConfigMap from it, so the state
			// should also not contain the base64 encoding of secret message.
			assert.NotContains(t, string(state), b64.StdEncoding.EncodeToString([]byte(secretMessage)))
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join("secrets", "step2"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					secretMessage += "updated"

					assert.NotNil(t, stackInfo.Deployment)
					state, err := json.Marshal(stackInfo.Deployment)
					assert.NoError(t, err)

					ssStringDataData, ok := stackInfo.Outputs["ssStringDataData"]
					assert.Truef(t, ok, "missing expected output \"ssStringDataData\"")

					ssStringDataStringData, ok := stackInfo.Outputs["ssStringDataStringData"]
					assert.Truef(t, ok, "missing expected output \"ssStringDataStringData\"")

					cgSecretStringData, ok := stackInfo.Outputs["cgSecretStringData"]
					assert.Truef(t, ok, "missing expected output \"cgSecretStringData\"")

					assert.NotEmptyf(t, ssStringDataData, "data field is empty")
					assert.NotEmptyf(t, ssStringDataStringData, "stringData field is empty")
					assert.NotEmptyf(t, cgSecretStringData, "cgSecretStringData field is empty")

					assert.NotContains(t, string(state), secretMessage)

					// The program converts the secret message to base64, to make a ConfigMap from it, so the state
					// should also not contain the base64 encoding of secret message.
					assert.NotContains(t, string(state), b64.StdEncoding.EncodeToString([]byte(secretMessage)))
				},
			},
		},
	})
	integration.ProgramTest(t, &test)
}

// TestSecretDataNewLine tests that secrets with new lines in the base64 encoding are handled correctly.
// See: https://github.com/pulumi/pulumi-kubernetes/issues/2681
func TestSecretDataNewLine(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  "secrets-new-line",
		ExpectRefreshChanges: false,
		SkipRefresh:          false,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)

			data, ok := stackInfo.Outputs["data"]
			assert.Truef(t, ok, "missing expected output \"data\"")

			stringData, ok := stackInfo.Outputs["stringData"]
			assert.Falsef(t, ok, "unexpected non-empty output: \"stringData\"")

			assert.NotEmptyf(t, data, "data field is empty")
			assert.Emptyf(t, stringData, "stringData field is not empty")

		},
		EditDirs: []integration.EditDir{{
			Dir:             "secrets-new-line",
			ExpectNoChanges: true, // Re-running the same program should not cause any changes.
			Additive:        true,
		}},
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
			nsPatched := stackInfo.Outputs["nsPatched"].(map[string]any)
			fooV, ok, err := unstructured.NestedString(nsPatched, "metadata", "labels", "foo")
			assert.True(t, ok)
			assert.NoError(t, err)
			assert.Equal(t, "foo", fooV)

			// Validate patched CustomResource
			crPatched := stackInfo.Outputs["crPatched"].(map[string]any)
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
					deploymentPatched := stackInfo.Outputs["deploymentPatched"].(map[string]any)
					containersV, ok, err := unstructured.NestedSlice(
						deploymentPatched, "spec", "template", "spec", "containers")
					assert.True(t, ok)
					assert.NoError(t, err)
					assert.Len(t, containersV, 1)
					limitsV, ok, err := unstructured.NestedMap(
						containersV[0].(map[string]any), "resources", "limits")
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

// TestServerSideApplyEmptyMaps tests that we correctly handle merging structs containing empty maps when diffing live and wanted
// states. This is a regression test for issue #2332 to ensure Pulumi can handle updating a resource which has a
// map field that is empty in the live state but non-empty in the wanted state.
func TestServerSideApplyEmptyMaps(t *testing.T) {
	var ns, cmName string

	applyStep := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("server-side-apply-empty-maps", "configmap"),
		ExpectRefreshChanges: true,
		// Enable destroy-on-cleanup so we can shell out to kubectl to make external changes to the resource and reuse the same stack.
		DestroyOnCleanup: true,
		Quick:            true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			cm := stackInfo.Outputs["cm"].(map[string]any)
			// Save the name and namespace for later use with kubectl. We check that the vars are empty,
			// in case pulumi up creates a new ConfigMap/Namespace instead of updating the existing one on
			// subsequent runs.
			if ns == "" && cmName == "" {
				ns = cm["metadata"].(map[string]any)["namespace"].(string)
				cmName = cm["metadata"].(map[string]any)["name"].(string)
			}

			// Validate we applied ConfigMap with wanted labels.
			fooV, ok, err := unstructured.NestedString(cm, "metadata", "labels", "foo")
			assert.True(t, ok)
			assert.NoError(t, err)
			assert.Equal(t, "bar", fooV)
		},
	})

	// Use manual lifecycle management since we need to run external commands in between pulumi up steps, while referencing
	// the same stack.
	pt := integration.ProgramTestManualLifeCycle(t, &applyStep)
	err := pt.TestLifeCycleInitAndDestroy()
	assert.NoError(t, err)

	// Sanity check with kubectl to verify that the ConfigMap was created with the wanted label.
	out, err := exec.Command("kubectl", "get", "configmap", "-o", "yaml", "-n", ns, cmName).CombinedOutput()
	assert.NoError(t, err)
	assert.Contains(t, string(out), "bar") // ConfigMap should have been created with label foo=bar.

	// Update the ConfigMap and remove label using kubectl.
	out, err = exec.Command("kubectl", "label", "configmap", "-n", ns, cmName, "foo-").CombinedOutput()
	assert.NoError(t, err)
	assert.Contains(t, string(out), "configmap/"+cmName+" unlabeled") // Ensure CM was unlabeled.

	// Use kubectl to verify that the ConfigMap was updated and no longer has the label.
	out, err = exec.Command("kubectl", "get", "configmap", "-o", "yaml", "-n", ns, cmName).CombinedOutput()
	assert.NoError(t, err)
	assert.NotContains(t, string(out), "bar") // ConfigMap should no longer have label foo=bar.

	// Run `pulumi up` + `pulumi refresh` to refresh the state and detect the missing label.
	// (The program tester runs these as separate steps, so the `pulumi up` doesn't detect a change until after the
	// subsequent refresh is performed.)
	err = pt.TestPreviewUpdateAndEdits()
	assert.NoError(t, err)

	// Re-run `pulumi up` to update the ConfigMap and re-add the label.
	err = pt.TestPreviewUpdateAndEdits()
	assert.NoError(t, err)

	// Use kubectl to verify that the ConfigMap was updated and has the label again.
	out, err = exec.Command("kubectl", "get", "configmap", "-o", "yaml", "-n", ns, cmName).CombinedOutput()
	assert.NoError(t, err)
	assert.Contains(t, string(out), "bar") // ConfigMap should have been updated with label foo=bar.
}

func TestServerSideApplyUpgrade(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("server-side-apply-upgrade", "step1"),
		ExpectRefreshChanges: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			// Validate Provider config
			provider := stackInfo.Outputs["provider"].(map[string]any)
			enableSSA, ok, err := unstructured.NestedString(provider, "enableServerSideApply")
			assert.True(t, ok)
			assert.NoError(t, err)
			assert.Equalf(t, "false", enableSSA, "SSA should be disabled")
			for _, res := range stackInfo.Deployment.Resources {
				if res.Type == "kubernetes:apps/v1:Deployment" {
					live := &unstructured.Unstructured{Object: res.Outputs}
					foundExpectedManager := false
					for _, managedField := range live.GetManagedFields() {
						if managedField.Manager == "pulumi-kubernetes" && managedField.Operation == "Update" {
							foundExpectedManager = true
							break
						}
					}
					assert.Truef(t, foundExpectedManager, "missing expected pulumi-kubernetes field manager with operation type Update")
				}
			}
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join("server-side-apply-upgrade", "step2"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					// Validate Provider config
					provider := stackInfo.Outputs["provider"].(map[string]any)
					enableSSA, ok, err := unstructured.NestedString(provider, "enableServerSideApply")
					assert.True(t, ok)
					assert.NoError(t, err)
					assert.Equalf(t, "true", enableSSA, "SSA should be enabled")
					for _, res := range stackInfo.Deployment.Resources {
						if res.Type == "kubernetes:apps/v1:Deployment" {
							live := &unstructured.Unstructured{Object: res.Outputs}
							foundExpectedManager := false
							for _, managedField := range live.GetManagedFields() {
								if managedField.Manager == "pulumi-kubernetes" && managedField.Operation == "Update" {
									foundExpectedManager = true
									break
								}
							}
							assert.Truef(t, foundExpectedManager, "missing expected pulumi-kubernetes field manager with operation type Update")
						}
					}
				},
			},
			{
				Dir:      filepath.Join("server-side-apply-upgrade", "step3"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					for _, res := range stackInfo.Deployment.Resources {
						if res.Type == "kubernetes:apps/v1:Deployment" {
							live := &unstructured.Unstructured{Object: res.Outputs}
							foundExpectedManager := false
							for _, managedField := range live.GetManagedFields() {
								// Operation type should change to Apply on first update to the Deployment.
								if managedField.Manager == "pulumi-kubernetes" && managedField.Operation == "Apply" {
									foundExpectedManager = true
									break
								}
							}
							assert.Truef(t, foundExpectedManager, "missing expected pulumi-kubernetes field manager with operation type Apply")
							containers, ok := openapi.Pluck(res.Outputs, "spec", "template", "spec", "containers")
							assert.Truef(t, ok, "failed to get containers")
							containerStatus := containers.([]any)[0].(map[string]any)
							image := containerStatus["image"]
							assert.Equalf(t, image.(string), "nginx:1.17", "image should be updated")
							portsRaw := containerStatus["ports"]
							portsArray := portsRaw.([]any)
							assert.Equalf(t, len(portsArray), 1, "only one port should be set")
							portMap := portsArray[0].(map[string]any)
							assert.Equalf(t, portMap["containerPort"], float64(81), "port should be updated to 81")
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

			secretData := stackInfo.Outputs["data"].(map[string]any)

			assert.Contains(t, secretData, "ca.crt")
			assert.Contains(t, secretData, "token")
		},
	})
	integration.ProgramTest(t, &test)
}

func TestStrictMode(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:           filepath.Join("strict-mode", "step1"),
		Quick:         true,
		ExpectFailure: true,
		SkipRefresh:   true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "kubernetes:strictMode",
				Value: "true",
			},
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			// Check the event stream for a diagnostic event showing that a default provider is prohibited.
			foundMessage := false
			msg := "strict mode prohibits default provider"
			for _, e := range stackInfo.Events {
				if e.DiagnosticEvent != nil && strings.Contains(e.DiagnosticEvent.Message, msg) {
					foundMessage = true
					break
				}
			}
			assert.Truef(t, foundMessage, "did not find expected failure message: %q", msg)
		},
		EditDirs: []integration.EditDir{
			{
				Dir:           filepath.Join("strict-mode", "step2"),
				Additive:      true,
				ExpectFailure: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					// Check the event stream for a diagnostic event showing that a Provider requires a "context".
					foundMessage := false
					msg := `strict mode requires Provider "context" argument`
					for _, e := range stackInfo.Events {
						if e.DiagnosticEvent != nil && strings.Contains(e.DiagnosticEvent.Message, msg) {
							foundMessage = true
							break
						}
					}
					assert.Truef(t, foundMessage, "did not find expected failure message: %q", msg)
					_, ok := stackInfo.Outputs["cm"]
					assert.Falsef(t, ok, "ConfigMap should not be present since Provider is invalid")
				},
			},
			{
				Dir:           filepath.Join("strict-mode", "step3"),
				Additive:      true,
				ExpectFailure: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					// Check the event stream for a diagnostic event showing that a Provider requires a "kubeconfig".
					foundMessage := false
					msg := `strict mode requires Provider "kubeconfig" argument`
					for _, e := range stackInfo.Events {
						if e.DiagnosticEvent != nil && strings.Contains(e.DiagnosticEvent.Message, msg) {
							foundMessage = true
							break
						}
					}
					assert.Truef(t, foundMessage, "did not find expected failure message: %q", msg)
					_, ok := stackInfo.Outputs["cm"]
					assert.Falsef(t, ok, "ConfigMap should not be present since Provider is invalid")
				},
			},
		},
	})
	integration.ProgramTest(t, &test)
}

// TestClientSideDriftCorrectCSA tests that we can successfully reapply a resource that has been
// modified outside of Pulumi using the client-side apply strategy.
func TestClientSideDriftCorrectCSA(t *testing.T) {
	var ns, cmName string

	applyStep := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("drift-correct", "configmap-csa"),
		ExpectRefreshChanges: true,
		// Enable destroy-on-cleanup so we can shell out to kubectl to make external changes to the resource and reuse the same stack.
		DestroyOnCleanup: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			cm, ok := stackInfo.Outputs["cm"].(map[string]interface{})
			assert.True(t, ok)

			// Save the name and namespace for later use with kubectl. We check that the vars are empty,
			// in case pulumi up creates a new ConfigMap/Namespace instead of updating the existing one on
			// subsequent runs.
			if ns == "" && cmName == "" {
				ns = cm["metadata"].(map[string]interface{})["namespace"].(string)
				cmName = cm["metadata"].(map[string]interface{})["name"].(string)
			}

			// Validate we applied ConfigMap with data.
			fooV, ok, err := unstructured.NestedString(cm, "data", "foo")
			assert.True(t, ok)
			assert.NoError(t, err)
			assert.Equal(t, "bar", fooV)
		},
	})

	// Use manual lifecycle management since we need to run external commands in between pulumi up steps, while referencing
	// the same stack.
	pt := integration.ProgramTestManualLifeCycle(t, &applyStep)
	err := pt.TestLifeCycleInitAndDestroy()
	assert.NoError(t, err)

	// Sanity check with kubectl to verify that the ConfigMap was created with the wanted label.
	out, err := exec.Command("kubectl", "get", "configmap", "-o", "yaml", "-n", ns, cmName).CombinedOutput()
	assert.NoError(t, err)
	assert.Contains(t, string(out), "bar") // ConfigMap should have been created with data foo: bar.

	// Update the ConfigMap and change the data foo: bar to foo: baz.
	out, err = exec.Command("kubectl", "patch", "configmap", "-n", ns, cmName, "-p", `{"data":{"foo":"baz"}}`).CombinedOutput()
	assert.NoError(t, err)
	assert.Contains(t, string(out), "configmap/"+cmName+" patched") // Ensure CM was patched.

	// Use kubectl to verify that the ConfigMap was updated and now has data foo: baz.
	out, err = exec.Command("kubectl", "get", "configmap", "-o", "yaml", "-n", ns, cmName).CombinedOutput()
	assert.NoError(t, err)
	assert.NotContains(t, string(out), "foo: bar") // ConfigMap should no longer have data foo: bar.
	assert.Contains(t, string(out), "foo: baz")    // ConfigMap should have data foo: baz.

	// Run `pulumi up` + `pulumi refresh` to refresh the state and detect the missing label.
	// (The program tester runs these as separate steps, so the `pulumi up` doesn't detect a change until after the
	// subsequent refresh is performed.)
	err = pt.TestPreviewUpdateAndEdits()
	assert.NoError(t, err)

	// Re-run `pulumi up` to update the ConfigMap and re-add the label.
	err = pt.TestPreviewUpdateAndEdits()
	assert.NoError(t, err)

	// Use kubectl to verify that the ConfigMap was updated and has the label again.
	out, err = exec.Command("kubectl", "get", "configmap", "-o", "yaml", "-n", ns, cmName).CombinedOutput()
	assert.NoError(t, err)

	assert.Contains(t, string(out), "foo: bar")    // ConfigMap should have been updated with data foo: bar.
	assert.NotContains(t, string(out), "foo: baz") // ConfigMap should no longer have data foo: baz.
}

// TestClientSideDriftCorrectSSA tests that we can successfully reapply a resource that has been
// modified outside of Pulumi, with SSA enabled.
func TestClientSideDriftCorrectSSA(t *testing.T) {
	var ns, cmName string

	applyStep := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("drift-correct", "configmap-ssa"),
		ExpectRefreshChanges: true,
		// Enable destroy-on-cleanup so we can shell out to kubectl to make external changes to the resource and reuse the same stack.
		DestroyOnCleanup: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			cm, ok := stackInfo.Outputs["cm"].(map[string]interface{})
			assert.True(t, ok)

			// Save the name and namespace for later use with kubectl. We check that the vars are empty,
			// in case pulumi up creates a new ConfigMap/Namespace instead of updating the existing one on
			// subsequent runs.
			if ns == "" && cmName == "" {
				ns = cm["metadata"].(map[string]interface{})["namespace"].(string)
				cmName = cm["metadata"].(map[string]interface{})["name"].(string)
			}

			// Validate we applied ConfigMap with data.
			fooV, ok, err := unstructured.NestedString(cm, "data", "foo")
			assert.True(t, ok)
			assert.NoError(t, err)
			assert.Equal(t, "bar", fooV)
		},
	})

	// Use manual lifecycle management since we need to run external commands in between pulumi up steps, while referencing
	// the same stack.
	pt := integration.ProgramTestManualLifeCycle(t, &applyStep)
	err := pt.TestLifeCycleInitAndDestroy()
	assert.NoError(t, err)

	// Sanity check with kubectl to verify that the ConfigMap was created with the wanted label.
	out, err := exec.Command("kubectl", "get", "configmap", "-o", "yaml", "-n", ns, cmName).CombinedOutput()
	assert.NoError(t, err)
	assert.Contains(t, string(out), "bar") // ConfigMap should have been created with data foo: bar.

	// Update the ConfigMap and change the data foo: bar to foo: baz.
	out, err = exec.Command("kubectl", "patch", "configmap", "-n", ns, cmName, "-p", `{"data":{"foo":"baz"}}`).CombinedOutput()
	assert.NoError(t, err)
	assert.Contains(t, string(out), "configmap/"+cmName+" patched") // Ensure CM was patched.

	// Use kubectl to verify that the ConfigMap was updated and now has data foo: baz.
	out, err = exec.Command("kubectl", "get", "configmap", "-o", "yaml", "-n", ns, cmName).CombinedOutput()
	assert.NoError(t, err)
	assert.NotContains(t, string(out), "foo: bar") // ConfigMap should no longer have data foo: bar.
	assert.Contains(t, string(out), "foo: baz")    // ConfigMap should have data foo: baz.

	// Run `pulumi up` + `pulumi refresh` to refresh the state and detect the missing label.
	// (The program tester runs these as separate steps, so the `pulumi up` doesn't detect a change until after the
	// subsequent refresh is performed.)
	err = pt.TestPreviewUpdateAndEdits()
	assert.NoError(t, err)

	// Re-run `pulumi up` to update the ConfigMap and re-add the label.
	err = pt.TestPreviewUpdateAndEdits()
	assert.NoError(t, err)

	// Use kubectl to verify that the ConfigMap was updated and has the label again.
	out, err = exec.Command("kubectl", "get", "configmap", "-o", "yaml", "-n", ns, cmName).CombinedOutput()
	assert.NoError(t, err)

	assert.Contains(t, string(out), "foo: bar")    // ConfigMap should have been updated with data foo: bar.
	assert.NotContains(t, string(out), "foo: baz") // ConfigMap should no longer have data foo: baz.
}

// TestSkipUpdateUnreachableFlag tests that we can complete a pulumi up with the skipUpdateUnreachable flag set,
// while some clusters are unreachable. This is helpful when a Pulumi program is managing resources across multiple
// clusters, and some of the clusters are unreachable.
// See https://github.com/pulumi/pulumi-kubernetes/pull/2528 for more details.
// Steps:
//  1. Create a ConfigMap in 2 separate clusters. (We clone the same KUBECONFIG file to simulate 2 separate clusters.)
//  2. Disable access to one of the clusters.
//  3. Run `pulumi up` with the skip-update-unreachable flag set.
//  4. Validate that the resource in the reachable cluster was updated, and the resource in the unreachable cluster was not.
//  5. Re-enable access to the unreachable cluster and run `pulumi up` again.
//  6. Validate that the resource in the unreachable cluster was updated.
func TestSkipUpdateUnreachableFlag(t *testing.T) {
	var ns0, ns1, cm0, cm1 string

	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("skip-update-unreachable", "step1"),
		ExpectRefreshChanges: true,
		// Enable destroy-on-cleanup so we can shell out to kubectl to make external changes to the resource and reuse the same stack.
		DestroyOnCleanup: true,
		Env:              []string{"PULUMI_K8S_SKIP_UPDATE_UNREACHABLE=true"},
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			var ok bool
			cm0, ok = stackInfo.Outputs["cmName0"].(string)
			assert.True(t, ok)
			cm1, ok = stackInfo.Outputs["cmName1"].(string)
			assert.True(t, ok)
			ns0, ok = stackInfo.Outputs["nsName0"].(string)
			assert.True(t, ok)
			ns1, ok = stackInfo.Outputs["nsName1"].(string)
			assert.True(t, ok)
		},
	})

	t.Log("Creating kubeconfig files to simulate 2 clusters")
	kubeconfigs, err := simluateClusterKubeconfig(t, 2)
	assert.NoError(t, err)

	test = test.With(integration.ProgramTestOptions{
		Env: []string{"KUBECONFIG_CLUSTER_0=" + kubeconfigs[0], "KUBECONFIG_CLUSTER_1=" + kubeconfigs[0]},
	})

	// Use manual lifecycle management since we need to run external commands in between pulumi up steps, while referencing
	// the same stack.
	t.Log("Creating ConfigMaps in 2 separate clusters")
	pt := integration.ProgramTestManualLifeCycle(t, &test)
	err = pt.TestLifeCycleInitAndDestroy()
	assert.NoError(t, err)

	t.Log("Validating ConfigMaps were created in 2 separate clusters with data foo: step1")
	cm0Contents, err := tests.Kubectl("get configmap -o yaml -n", ns0, cm0)
	assert.NoError(t, err)
	assert.Contains(t, string(cm0Contents), "foo: step1") // ConfigMap should have been created with data foo: step1.
	cm1Contents, err := tests.Kubectl("get configmap -o yaml -n", ns1, cm1)
	assert.NoError(t, err)
	assert.Contains(t, string(cm1Contents), "foo: step1")

	t.Log("Disabling access to the second cluster by setting the KUBECONFIG to a fake filepath and re-running `pulumi up`")
	test.Env = append(test.Env, "KUBECONFIG_CLUSTER_1=/fake/filepath")

	t.Log("Updating the Pulumi program to use step2")
	test.EditDirs = []integration.EditDir{
		{
			Dir:      filepath.Join("skip-update-unreachable", "step2"),
			Additive: true,
		},
	}

	// Run `pulumi up`.
	err = pt.TestPreviewUpdateAndEdits()
	assert.NoError(t, err)

	t.Log("Validating ConfigMaps were updated in the reachable cluster, but not in the unreachable cluster")
	cm0Contents, err = tests.Kubectl("get configmap -o yaml -n", ns0, cm0)
	assert.NoError(t, err)
	assert.Contains(t, string(cm0Contents), "foo: step2") // ConfigMap should have been updated with data foo: step2.
	cm1Contents, err = tests.Kubectl("get configmap -o yaml -n", ns1, cm1)
	assert.NoError(t, err)
	assert.Contains(t, string(cm1Contents), "foo: step1") // ConfigMap should not have been updated.

	t.Log("Re-enabling access to the second cluster by setting the KUBECONFIG to the original kubeconfig and re-running `pulumi up`")
	test.Env = append(test.Env, "KUBECONFIG_CLUSTER_1="+kubeconfigs[1])

	err = pt.TestPreviewUpdateAndEdits()
	assert.NoError(t, err)

	t.Log("Validating ConfigMaps were updated in the unreachable cluster")
	cm1Contents, err = tests.Kubectl("get configmap -o yaml -n", ns1, cm1)
	assert.NoError(t, err)
	assert.Contains(t, string(cm1Contents), "foo: step2") // ConfigMap should have been updated with data foo: step2.
}

func simluateClusterKubeconfig(t *testing.T, numOfClusters int) ([]string, error) {
	t.Helper()

	// Load default kubeconfig as base.
	config, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return nil, err
	}

	// Create tmp dir to store kubeconfig.
	tmpDir, err := os.MkdirTemp("", "kubeconfig-preview")
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		log.Println("Deleting kubeconfig tmp dir")
		assert.NoError(t, os.RemoveAll(tmpDir))
	})

	// Write kubeconfig to tmp dir.
	kubeconfigs := make([]string, numOfClusters)
	for i := 0; i < numOfClusters; i++ {
		kubeconfigs[i] = filepath.Join(tmpDir, fmt.Sprintf("kubeconfig-%d.txt", i))
		err = clientcmd.WriteToFile(*config, kubeconfigs[i])
		if err != nil {
			return nil, err
		}
	}

	return kubeconfigs, err
}

// TestIgnoreChanges tests that we can successfully ignore changes to a resource without SSA conflicts,
// and that we use the right field value when ignoring changes obtained from the live cluster.
// SkipRefresh *must* be true to properly test that conflict is handled when the state is not refreshed
// and drift has occurred.
// https://github.com/pulumi/pulumi-kubernetes/issues/2542
func TestIgnoreChanges(t *testing.T) {
	testCases := []struct{ name, folderName string }{
		{name: "Server Side Apply Mode", folderName: "ignore-changes-ssa"},
		{name: "Client Side Apply Mode", folderName: "ignore-changes-csa"},
	}

	for _, tc := range testCases {
		// NB: the Pulumi Program test runs in parallel, so we need to shadow the tc var.
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ignoreChageTest(t, tc.folderName)
		})
	}
}

func ignoreChageTest(t *testing.T, testFolderName string) {
	var depName, depNS string

	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join(testFolderName, "step1"),
		ExpectRefreshChanges: true,
		// SkipRefresh MUST be true as the bug is not reproducible when the state is refreshed.
		SkipRefresh: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			var ok bool
			depName, ok = stackInfo.Outputs["deploymentName"].(string)
			assert.True(t, ok)
			depNS, ok = stackInfo.Outputs["deploymentNamespace"].(string)
			assert.True(t, ok)

			// Validate we applied the deployment with the correct name and namespace and spec.
			dep, err := tests.Kubectl("get deployment -o yaml -n", depNS, depName)
			assert.NoError(t, err)
			assert.NotContains(t, string(dep), "Error from server (NotFound)")

			// Validate image of deployment.
			depImage, err := tests.Kubectl("get deployment -o=jsonpath='{.spec.template.spec.containers[0].image}' -n", depNS, depName)
			assert.NoError(t, err)
			assert.Equal(t, "'nginx:1.25.2'", string(depImage))

			// Validate deployment replicas.
			depReplicas, err := tests.Kubectl("get deployment -o=jsonpath='{.spec.replicas}' -n", depNS, depName)
			assert.NoError(t, err)
			assert.Equal(t, "'2'", string(depReplicas))

			// Patch deployment replicas to 3 using patch file in preparation for ignore changes to be tested in step2.
			_, err = tests.Kubectl("patch --field-manager replica/manager deployment -n", depNS, depName, "--patch-file", filepath.Join(testFolderName, "deployment-patch.yaml"))
			assert.NoError(t, err)
			depReplicas, err = tests.Kubectl("get deployment -o=jsonpath='{.spec.replicas}' -n", depNS, depName)
			assert.NoError(t, err)
			assert.Equal(t, "'3'", string(depReplicas))
		},
		EditDirs: []integration.EditDir{
			{
				// Repeat step1 again where no changes are made to the deployment since we ignore changes to spec.replicas.
				Dir:      filepath.Join(testFolderName, "step1"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					// Validate replicas was not updated back to 1.
					depReplicas, err := tests.Kubectl("get deployment -o=jsonpath='{.spec.replicas}' -n", depNS, depName)
					assert.NoError(t, err)
					assert.Equal(t, "'3'", string(depReplicas))
				},
			},
			{
				Dir:      filepath.Join(testFolderName, "step2"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					// Validate image was updated, but spec.replicas was not.
					depImage, err := tests.Kubectl("get deployment -o=jsonpath='{.spec.template.spec.containers[0].image}' -n", depNS, depName)
					assert.NoError(t, err)
					assert.Equal(t, "'nginx:1.25.1'", string(depImage))

					depReplicas, err := tests.Kubectl("get deployment -o=jsonpath='{.spec.replicas}' -n", depNS, depName)
					assert.NoError(t, err)
					assert.Equal(t, "'3'", string(depReplicas))

					// Now use kubectl patch to update spec.replicas to 4 and see if we can correctly ignore changes to spec.replicas again when the field manager is
					// "kubectl-patch" since we have logic to override certain field managers with manager name prefixes. This is due to fluxssa.PatchReplaceFieldsManagers
					// doing a prefix match on the field manager name instead of an exact match on the given field manager name.
					_, err = tests.Kubectl("patch deployment -n", depNS, depName, "--patch-file", filepath.Join(testFolderName, "deployment-patch-2.yaml"))
					assert.NoError(t, err)
					depReplicas, err = tests.Kubectl("get deployment -o=jsonpath='{.spec.replicas}' -n", depNS, depName)
					assert.NoError(t, err)
					assert.Equal(t, "'4'", string(depReplicas))
				},
			},
			{
				Dir:      filepath.Join(testFolderName, "step3"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					// Validate image was updated, but spec.replicas was not.
					depImage, err := tests.Kubectl("get deployment -o=jsonpath='{.spec.template.spec.containers[0].image}' -n", depNS, depName)
					assert.NoError(t, err)
					assert.Equal(t, "'nginx:1.25'", string(depImage))

					depReplicas, err := tests.Kubectl("get deployment -o=jsonpath='{.spec.replicas}' -n", depNS, depName)
					assert.NoError(t, err)
					assert.Equal(t, "'4'", string(depReplicas))
				},
			},
		},
	})

	integration.ProgramTest(t, &test)
}

// TestEmptyItemNormalization tests that we correctly handle empty items when normalizing resources. We should
// not remove a list that contains an empty struct, as this is different from an empty list.
// See https://github.com/pulumi/pulumi-kubernetes/issues/2538 for more details.
// This test requires a cluster with NetworkPolicy support to ensure the NetworkPolicy resource is created
// and has a controller backing it. We create 2 pods to test egress between them, rather than hitting
// a live URL, to avoid flakiness.
func TestEmptyItemNormalization(t *testing.T) {
	tests.SkipIfShort(t, "test requires a cluster with NetworkPolicy support")
	validateProgram := func(networkingEnabled bool) func(*testing.T, integration.RuntimeValidationStackInfo) {
		return func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			ns, ok := stackInfo.Outputs["podANamespace"].(string)
			assert.True(t, ok)
			podA, ok := stackInfo.Outputs["podAName"].(string)
			assert.True(t, ok)
			nginxIP, ok := stackInfo.Outputs["nginxIP"].(string)
			assert.True(t, ok)
			np, ok := stackInfo.Outputs["networkPolicyName"].(string)
			assert.True(t, ok)

			// Sanity check with kubectl to verify that the NetworkPolicy was created with the wanted label.
			out, err := tests.Kubectl("get networkpolicies.networking.k8s.io", np, "-n", ns)
			assert.NoError(t, err)
			assert.NotContains(t, string(out), "Error from server (NotFound)")

			// Exec into pod and verify egress to/ ingress from podB based on step.
			// Step 1: Egress/Ingress should be blocked.
			// Step 2: Egress/Ingress should be allowed.
			out, err = tests.Kubectl("exec -i -n", ns, podA, "-- wget -qO- --timeout=5 ", nginxIP)
			if networkingEnabled {
				assert.NoError(t, err)
				assert.Contains(t, string(out), "Welcome to nginx!")
			} else {
				assert.Error(t, err)
				assert.Contains(t, string(out), "wget: download timed out")
			}
		}
	}

	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("network-policy", "step1"),
		ExpectRefreshChanges: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		EditDirs: []integration.EditDir{
			{
				Dir:                    filepath.Join("network-policy", "step2"),
				Additive:               true,
				ExtraRuntimeValidation: validateProgram(true),
			},
		},
		ExtraRuntimeValidation: validateProgram(false),
	})

	integration.ProgramTest(t, &test)
}

// TestFieldManagerPatchResources ensures field managers are not patched when dealing with a Patch resource
// (e.g., DaemonSetPatch). This precaution prevents the Patch resource from taking ownership of fields managed by kubectl field managers,
// which would otherwise result in the unintended unsetting of all unspecified fields of a resource via Server-Side Apply (SSA).
// For additional context, refer to: https://github.com/pulumi/pulumi-kubernetes/issues/2639.
// Test Steps:
//  1. Deploy an nginx:1.14.2 image with 2 replicas using kubectl which sets the field manager to kubectl-client-side-apply.
//  2. Apply a DeploymentPatch resource to modify the deployment by changing the image to nginx:1.14.1 (Step 1).
//  3. Update the DeploymentPatch resource to further patch the deployment, setting the image to nginx:1.14.0,
//     and verify that other fields managed by kubectl-client-side-apply remain unaffected (Step 2).
func TestFieldManagerPatchResources(t *testing.T) {
	testFolder := "field-manager-patch-resources"

	createDeployment := func() string {
		// Create a random namespace to deploy the nginx deployment to.
		ns := "test-filed-mgr-" + rand.String(5)
		_, err := tests.Kubectl("create namespace", ns)
		require.NoError(t, err)
		t.Cleanup(func() {
			_, _ = tests.Kubectl("delete namespace", ns)
		})

		// Create nginx deployment.
		_, err = tests.Kubectl("apply -f", filepath.Join(testFolder, "deployment.yaml"), "-n", ns)
		require.NoError(t, err)
		t.Cleanup(func() {
			_, _ = tests.Kubectl("delete namespace", ns)
		})

		return ns
	}

	namespace := createDeployment()

	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join(testFolder, "step1"),
		ExpectRefreshChanges: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		Config: map[string]string{
			"namespace": namespace,
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			// Ensure that the nginx deployment was patched with image nginx:1.14.1.
			depImage, err := tests.Kubectl("get deployment -o=jsonpath={.spec.template.spec.containers[0].image} -n", namespace, "test-mgr-nginx")
			assert.NoError(t, err)
			assert.Equal(t, "nginx:1.14.1", string(depImage))

			// Ensure that the nginx deployment replicas is still 2.
			depReplicas, err := tests.Kubectl("get deployment -o=jsonpath={.spec.replicas} -n", namespace, "test-mgr-nginx")
			assert.NoError(t, err)
			assert.Equal(t, "2", string(depReplicas))
		},
		EditDirs: []integration.EditDir{
			{
				Dir:      filepath.Join(testFolder, "step2"),
				Additive: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					// Ensure that the nginx deployment was patched with image nginx:1.14.1.
					depImage, err := tests.Kubectl("get deployment -o=jsonpath={.spec.template.spec.containers[0].image} -n", namespace, "test-mgr-nginx")
					assert.NoError(t, err)
					assert.Equal(t, "nginx:1.14.0", string(depImage))

					// Ensure that the nginx deployment replicas is still 2, and was not unset to the default 1 due to field manager being patched.
					depReplicas, err := tests.Kubectl("get deployment -o=jsonpath={.spec.replicas} -n", namespace, "test-mgr-nginx")
					assert.NoError(t, err)
					assert.Equal(t, "2", string(depReplicas))

					// Ensure that we don't inadvertently share ownership of nested fields that we specify in ignoreChanges.
					// See: https://github.com/pulumi/pulumi-kubernetes/issues/2714.
					liveObj, err := tests.Kubectl("get deployment -o yaml --show-managed-fields -n", namespace, "test-mgr-nginx")
					assert.NoError(t, err)
					wantString := ` - apiVersion: apps/v1
    fieldsType: FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          f:pulumi.com/patchForce: {}
      f:spec:
        f:template:
          f:spec:
            f:containers:
              k:{"name":"nginx"}:
                .: {}
                f:image: {}
                f:name: {}
    manager: pulumi-kubernetes-`
					assert.Contains(t, string(liveObj), wantString)
				},
			},
		},
	})

	integration.ProgramTest(t, &test)
}

// TestOptionPropagation tests the handling of resource options by the various compoonent resources.
// Component resources are responsible for implementing option propagation logic when creating
// child resources.
func TestOptionPropagation(t *testing.T) {
	g := NewWithT(t)
	format.MaxLength = 0
	format.MaxDepth = 5
	format.RegisterCustomFormatter(pulumirpctesting.FormatDebugInterceptorLog)

	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	grpcLog, err := pulumirpctesting.NewDebugInterceptorLog(t)
	require.NoError(t, err)

	options := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join(cwd, "options"),
		Env:                  []string{grpcLog.Env()},
		Quick:                true,
		ExpectRefreshChanges: false,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {

			// lookup some resources for later use
			providerA := tests.SearchResourcesByName(stackInfo, "", "pulumi:providers:kubernetes", "a")
			require.NotNil(t, providerA)
			providerB := tests.SearchResourcesByName(stackInfo, "", "pulumi:providers:kubernetes", "b")
			require.NotNil(t, providerB)
			providerNullOpts := tests.SearchResourcesByName(stackInfo, "", "pulumi:providers:kubernetes", "nullopts")
			require.NotNil(t, providerNullOpts)
			sleep := tests.SearchResourcesByName(stackInfo, "", "time:index/sleep:Sleep", "sleep")
			require.NotNil(t, sleep)

			// some helper functions for naming purposes
			providerUrn := func(prov *apitype.ResourceV3) resource.URN {
				return prov.URN + resource.URNNameDelimiter + resource.URN(prov.ID)
			}
			urn := func(parentType, baseType tokens.Type, name string) resource.URN {
				return resource.NewURN(stackInfo.StackName, "options-test", parentType, baseType, name)
			}

			// read the GRPC log file to inspect the RegisterResource calls, since they provide
			// the most detailed view of the resource's options as determined by the SDK.
			logEntries, err := grpcLog.ReadAll()
			require.NoError(t, err)
			rr := logEntries.ListRegisterResource()
			invokes := logEntries.Invokes()

			// Verify that the invokes for provider A contain version info across-the-board.
			// The Version and PluginDownloadURL options normally serve as hints when selecting
			// a default provider, and should be propagated. For testing purposes, we set the provider explicitly to avoid
			// any attempt to use the fake version/url.
			g.Expect(invokes.ByProvider(providerUrn(providerA))).To(HaveEach(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Version": Equal("1.2.3"),
						// bug: https://github.com/pulumi/pulumi/issues/14839
						// "PluginDownloadURL": Equal("https://a.pulumi.test"),
					}),
				}),
			))

			// --- ConfigGroup ---

			// ConfigGroup "cg-options" with most options
			g.Expect(rr.Named(stackInfo.RootResource.URN,
				"kubernetes:yaml:ConfigGroup", "cg-options")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           HaveExactElements(Alias("cg-options-old"), Alias("cg-options-aliased")),
						"Protect":           PointTo(BeTrue()),
						"Dependencies":      HaveExactElements(string(sleep.URN)),
						"Provider":          BeEmpty(),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerA)),
						}),
						"IgnoreChanges": HaveExactElements("ignored"),
					}),
				}),
			))
			g.Expect(rr.Named(urn("", "kubernetes:yaml:ConfigGroup", "cg-options"),
				"kubernetes:core/v1:ConfigMap", "cg-options-cg-options-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           HaveExactElements(Alias("cg-options-cm-1-k8s-aliased"), Alias("cg-options-cg-options-cm-1-aliased")),
						"Protect":           PointTo(BeTrue()),
						"Dependencies":      BeEmpty(),
						"Provider":          BeEquivalentTo(providerUrn(providerA)),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"Providers":         BeEmpty(),
						"IgnoreChanges":     BeEmpty(),
						"Object": PointTo(ProtobufStruct(MatchKeys(IgnoreExtras, Keys{
							"metadata": MatchKeys(IgnoreExtras, Keys{
								"name":        Equal("cg-options-cm-1"),
								"annotations": And(HaveKey("pulumi.com/skipAwait"), HaveKey("transformed")),
							}),
						}))),
					}),
				}),
			))
			g.Expect(rr.Named(urn("", "kubernetes:yaml:ConfigGroup", "cg-options"),
				"kubernetes:yaml:ConfigFile", "cg-options-./testdata/options/configgroup/manifest.yaml")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           HaveExactElements(Alias("cg-options-./testdata/options/configgroup/manifest.yaml-aliased")),
						"Protect":           PointTo(BeTrue()),
						"Dependencies":      BeEmpty(),
						"Provider":          BeEmpty(),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"IgnoreChanges":     BeEmpty(),
					}),
				}),
			))
			g.Expect(rr.Named(urn("kubernetes:yaml:ConfigGroup", "kubernetes:yaml:ConfigFile", "cg-options-./testdata/options/configgroup/manifest.yaml"),
				"kubernetes:core/v1:ConfigMap", "cg-options-configgroup-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           HaveExactElements(Alias("configgroup-cm-1-k8s-aliased"), Alias("cg-options-configgroup-cm-1-aliased")),
						"Protect":           PointTo(BeTrue()),
						"Dependencies":      BeEmpty(),
						"Provider":          BeEquivalentTo(providerUrn(providerA)),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"Providers":         BeEmpty(),
						"IgnoreChanges":     BeEmpty(),
						"Object": PointTo(ProtobufStruct(MatchKeys(IgnoreExtras, Keys{
							"metadata": MatchKeys(IgnoreExtras, Keys{
								"name":        Equal("configgroup-cm-1"),
								"annotations": And(HaveKey("pulumi.com/skipAwait"), HaveKey("transformed")),
							}),
						}))),
					}),
				}),
			))

			// ConfigGroup "cg-provider" with "provider" option that should propagate to children.
			g.Expect(rr.Named(stackInfo.RootResource.URN,
				"kubernetes:yaml:ConfigGroup", "cg-provider")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEmpty(),
						"Version":  BeEmpty(),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerB)),
						}),
					}),
				}),
			))

			// ConfigGroup "cg-nullopts" with a stack transform to apply a "provider" option.
			g.Expect(rr.Named(stackInfo.RootResource.URN,
				"kubernetes:yaml:ConfigGroup", "cg-nullopts")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEmpty(),
						"Version":  BeEmpty(),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerNullOpts)),
						}),
					}),
				}),
			))

			// --- ConfigFile ---

			// ConfigFile "cf-options" with most options
			g.Expect(rr.Named(stackInfo.RootResource.URN,
				"kubernetes:yaml:ConfigFile", "cf-options-cf-options")).To(HaveExactElements(
				// quirk: NodeJS SDK applies resource_prefix ("cf-options") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           HaveExactElements(Alias("cf-options-old"), Alias("cf-options-cf-options-aliased")),
						"Protect":           PointTo(BeTrue()),
						"Dependencies":      HaveExactElements(string(sleep.URN)),
						"Provider":          BeEmpty(),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerA)),
						}),
						"IgnoreChanges": HaveExactElements("ignored"),
					}),
				}),
			))
			g.Expect(rr.Named(urn("", "kubernetes:yaml:ConfigFile", "cf-options-cf-options"),
				"kubernetes:core/v1:ConfigMap", "cf-options-configfile-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           HaveExactElements(Alias("configfile-cm-1-k8s-aliased"), Alias("cf-options-configfile-cm-1-aliased")),
						"Protect":           PointTo(BeTrue()),
						"Dependencies":      BeEmpty(),
						"Provider":          BeEquivalentTo(providerUrn(providerA)),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"Providers":         BeEmpty(),
						"IgnoreChanges":     BeEmpty(),
						"Object": PointTo(ProtobufStruct(MatchKeys(IgnoreExtras, Keys{
							"metadata": MatchKeys(IgnoreExtras, Keys{
								"name":        Equal("configfile-cm-1"),
								"annotations": And(HaveKey("pulumi.com/skipAwait"), HaveKey("transformed")),
							}),
						}))),
					}),
				}),
			))

			// ConfigFile "cf-provider" with "provider" option that should propagate to children.
			g.Expect(rr.Named(stackInfo.RootResource.URN,
				"kubernetes:yaml:ConfigFile", "cf-provider-cf-provider")).To(HaveExactElements(
				// quirk: NodeJS SDK applies resource_prefix ("cf-provider") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEmpty(),
						"Version":  BeEmpty(),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerB)),
						}),
					}),
				}),
			))
			g.Expect(rr.Named(urn("", "kubernetes:yaml:ConfigFile", "cf-provider-cf-provider"),
				"kubernetes:core/v1:ConfigMap", "cf-provider-configfile-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider":  BeEquivalentTo(providerUrn(providerB)),
						"Version":   Not(BeEmpty()),
						"Providers": BeEmpty(),
						"Object": PointTo(ProtobufStruct(MatchKeys(IgnoreExtras, Keys{
							"metadata": MatchKeys(IgnoreExtras, Keys{
								"name": Equal("configfile-cm-1"),
							}),
						}))),
					}),
				}),
			))

			// ConfigFile "cf-nullopts" with a stack transform to apply a "provider" option.
			g.Expect(rr.Named(stackInfo.RootResource.URN,
				"kubernetes:yaml:ConfigFile", "cf-nullopts-cf-nullopts")).To(HaveExactElements(
				// quirk: NodeJS SDK applies resource_prefix ("cf-nullopts") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEmpty(),
						"Version":  BeEmpty(),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerNullOpts)),
						}),
					}),
				}),
			))

			// --- Directory ---

			// Directory "kustomize-options" with most options
			g.Expect(rr.Named(stackInfo.RootResource.URN,
				"kubernetes:kustomize:Directory", "kustomize-options-kustomize-options")).To(HaveExactElements(
				// quirk: NodeJS SDK applies resource_prefix ("kustomize-options") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           HaveExactElements(Alias("kustomize-options-old"), Alias("kustomize-options-kustomize-options-aliased")),
						"Protect":           PointTo(BeTrue()),
						"Dependencies":      HaveExactElements(string(sleep.URN)),
						"Provider":          BeEmpty(),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerA)),
						}),
						"IgnoreChanges": HaveExactElements("ignored"),
					}),
				}),
			))
			g.Expect(rr.Named(urn("", "kubernetes:kustomize:Directory", "kustomize-options-kustomize-options"),
				"kubernetes:core/v1:ConfigMap", "kustomize-options-kustomize-cm-1-2kkk4bthmg")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           HaveExactElements(Alias("kustomize-cm-1-2kkk4bthmg-k8s-aliased"), Alias("kustomize-options-kustomize-cm-1-2kkk4bthmg-aliased")),
						"Protect":           PointTo(BeTrue()),
						"Dependencies":      BeEmpty(),
						"Provider":          BeEquivalentTo(providerUrn(providerA)),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"Providers":         BeEmpty(),
						"IgnoreChanges":     BeEmpty(),
						"Object": PointTo(ProtobufStruct(MatchKeys(IgnoreExtras, Keys{
							"metadata": MatchKeys(IgnoreExtras, Keys{
								"name":        Equal("kustomize-cm-1-2kkk4bthmg"),
								"annotations": And(HaveKey("transformed")),
							}),
						}))),
					}),
				}),
			))

			// Directory "kustomize-provider" with "provider" option that should propagate to children.
			g.Expect(rr.Named(stackInfo.RootResource.URN,
				"kubernetes:kustomize:Directory", "kustomize-provider-kustomize-provider")).To(HaveExactElements(
				// quirk: NodeJS SDK applies resource_prefix ("kustomize-provider") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEmpty(),
						"Version":  BeEmpty(),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerB)),
						}),
					}),
				}),
			))
			g.Expect(rr.Named(urn("", "kubernetes:kustomize:Directory", "kustomize-provider-kustomize-provider"),
				"kubernetes:core/v1:ConfigMap", "kustomize-provider-kustomize-cm-1-2kkk4bthmg")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider":  BeEquivalentTo(providerUrn(providerB)),
						"Version":   Not(BeEmpty()),
						"Providers": BeEmpty(),
					}),
				}),
			))

			// Directory "kustomize-nullopts" with a stack transform to apply a "provider" option.
			g.Expect(rr.Named(stackInfo.RootResource.URN,
				"kubernetes:kustomize:Directory", "kustomize-nullopts-kustomize-nullopts")).To(HaveExactElements(
				// quirk: NodeJS SDK applies resource_prefix ("kustomize-nullopts") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEmpty(),
						"Version":  BeEmpty(),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerNullOpts)),
						}),
					}),
				}),
			))

			// --- Chart ---

			// Chart "chart-options"
			g.Expect(rr.Named(stackInfo.RootResource.URN,
				"kubernetes:helm.sh/v3:Chart", "chart-options-chart-options")).To(HaveExactElements(
				// quirk: NodeJS SDK applies resource_prefix ("chart-options") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases": HaveExactElements(
							Alias(tokens.Type("kubernetes:helm.sh/v2:Chart")),
							Alias("chart-options-old"),
							Alias("chart-options-chart-options-aliased")),
						"Protect":           PointTo(BeTrue()),
						"Dependencies":      HaveExactElements(string(sleep.URN)),
						"Provider":          BeEmpty(),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerA)),
						}),
						"IgnoreChanges": HaveExactElements("ignored"),
					}),
				}),
			))
			g.Expect(rr.Named(urn("", "kubernetes:helm.sh/v3:Chart", "chart-options-chart-options"),
				"kubernetes:core/v1:ConfigMap", "chart-options-chart-options-chart-options-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           HaveExactElements(Alias("chart-options-chart-options-cm-1-k8s-aliased"), Alias("chart-options-chart-options-chart-options-cm-1-aliased")),
						"Protect":           PointTo(BeTrue()),
						"Dependencies":      BeEmpty(),
						"Provider":          BeEquivalentTo(providerUrn(providerA)),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"Providers":         BeEmpty(),
						"IgnoreChanges":     BeEmpty(),
						"Object": PointTo(ProtobufStruct(MatchKeys(IgnoreExtras, Keys{
							"metadata": MatchKeys(IgnoreExtras, Keys{
								"name":        Equal("chart-options-chart-options-cm-1"), // note: based on the Helm Release name
								"annotations": And(HaveKey("pulumi.com/skipAwait")),
							}),
						}))),
					}),
				}),
			))

			// Chart "chart-provider" with "provider" option that should propagate to children.
			g.Expect(rr.Named(stackInfo.RootResource.URN,
				"kubernetes:helm.sh/v3:Chart", "chart-provider-chart-provider")).To(HaveExactElements(
				// quirk: NodeJS SDK applies resource_prefix ("chart-provider") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEmpty(),
						"Version":  BeEmpty(),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerB)),
						}),
					}),
				}),
			))
			g.Expect(rr.Named(urn("", "kubernetes:helm.sh/v3:Chart", "chart-provider-chart-provider"),
				"kubernetes:core/v1:ConfigMap", "chart-provider-chart-provider-chart-provider-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider":  BeEquivalentTo(providerUrn(providerB)),
						"Version":   Not(BeEmpty()),
						"Providers": BeEmpty(),
					}),
				}),
			))

			// Chart "chart-nullopts" with a stack transform to apply a "provider" option.
			g.Expect(rr.Named(stackInfo.RootResource.URN,
				"kubernetes:helm.sh/v3:Chart", "chart-nullopts-chart-nullopts")).To(HaveExactElements(
				// quirk: NodeJS SDK applies resource_prefix ("chart-options") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEmpty(),
						"Version":  BeEmpty(),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerNullOpts)),
						}),
					}),
				}),
			))
		},
	})

	pt := integration.ProgramTestManualLifeCycle(t, &options)

	err = pt.TestLifeCyclePrepare()
	require.NoError(t, err)
	t.Cleanup(pt.TestCleanUp)
	err = pt.TestLifeCycleInitialize()
	require.NoError(t, err)
	t.Cleanup(func() {
		// to ensure cleanup, we need to unprotected all resources
		err = pt.RunPulumiCommand("state", "unprotect", "--all", "-y")
		contract.IgnoreError(err)

		destroyErr := pt.TestLifeCycleDestroy()
		contract.IgnoreError(destroyErr)
	})

	err = pt.TestPreviewUpdateAndEdits()
	require.NoError(t, err)
}
