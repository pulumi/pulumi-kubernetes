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

package examples

import (
	"fmt"
	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/v3/resource/deploy/providers"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func TestAccMinimal(t *testing.T) {
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "minimal"),
		})

	integration.ProgramTest(t, &test)
}

func TestAccGuestbook(t *testing.T) {
	skipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "guestbook"),
			ExtraRuntimeValidation: func(
				t *testing.T, stackInfo integration.RuntimeValidationStackInfo,
			) {
				assert.NotNil(t, stackInfo.Deployment)
				assert.Equal(t, 9, len(stackInfo.Deployment.Resources))

				sort.Slice(stackInfo.Deployment.Resources, func(i, j int) bool {
					ri := stackInfo.Deployment.Resources[i]
					rj := stackInfo.Deployment.Resources[j]
					riname, _ := openapi.Pluck(ri.Outputs, "metadata", "name")
					rinamespace, _ := openapi.Pluck(ri.Outputs, "metadata", "namespace")
					rjname, _ := openapi.Pluck(rj.Outputs, "metadata", "name")
					rjnamespace, _ := openapi.Pluck(rj.Outputs, "metadata", "namespace")
					return fmt.Sprintf("%s/%s/%s", ri.URN.Type(), rinamespace, riname) <
						fmt.Sprintf("%s/%s/%s", rj.URN.Type(), rjnamespace, rjname)
				})

				var name interface{}
				var status interface{}

				// Verify frontend deployment.
				frontendDepl := stackInfo.Deployment.Resources[0]
				assert.Equal(t, tokens.Type("kubernetes:apps/v1:Deployment"), frontendDepl.URN.Type())
				name, _ = openapi.Pluck(frontendDepl.Outputs, "metadata", "name")
				assert.Equal(t, "frontend", name)
				status, _ = openapi.Pluck(frontendDepl.Outputs, "status", "readyReplicas")
				assert.Equal(t, float64(3), status)

				// Verify redis-master deployment.
				redisMasterDepl := stackInfo.Deployment.Resources[1]
				assert.Equal(t, tokens.Type("kubernetes:apps/v1:Deployment"), redisMasterDepl.URN.Type())
				name, _ = openapi.Pluck(redisMasterDepl.Outputs, "metadata", "name")
				assert.Equal(t, "redis-master", name)
				status, _ = openapi.Pluck(redisMasterDepl.Outputs, "status", "readyReplicas")
				assert.Equal(t, float64(1), status)

				// Verify redis-slave deployment.
				redisSlaveDepl := stackInfo.Deployment.Resources[2]
				assert.Equal(t, tokens.Type("kubernetes:apps/v1:Deployment"), redisSlaveDepl.URN.Type())
				name, _ = openapi.Pluck(redisSlaveDepl.Outputs, "metadata", "name")
				assert.Equal(t, "redis-slave", name)
				status, _ = openapi.Pluck(redisSlaveDepl.Outputs, "status", "readyReplicas")
				assert.Equal(t, float64(1), status)

				// Verify test namespace.
				namespace := stackInfo.Deployment.Resources[3]
				assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())

				// Verify frontend service.
				frontentService := stackInfo.Deployment.Resources[4]
				assert.Equal(t, tokens.Type("kubernetes:core/v1:Service"), frontentService.URN.Type())
				name, _ = openapi.Pluck(frontentService.Outputs, "metadata", "name")
				assert.Equal(t, "frontend", name)
				status, _ = openapi.Pluck(frontentService.Outputs, "spec", "clusterIP")
				assert.True(t, len(status.(string)) > 1)

				// Verify redis-master service.
				redisMasterService := stackInfo.Deployment.Resources[5]
				assert.Equal(t, tokens.Type("kubernetes:core/v1:Service"), redisMasterService.URN.Type())
				name, _ = openapi.Pluck(redisMasterService.Outputs, "metadata", "name")
				assert.Equal(t, "redis-master", name)
				status, _ = openapi.Pluck(redisMasterService.Outputs, "spec", "clusterIP")
				assert.True(t, len(status.(string)) > 1)

				// Verify redis-slave service.
				redisSlaveService := stackInfo.Deployment.Resources[6]
				assert.Equal(t, tokens.Type("kubernetes:core/v1:Service"), redisSlaveService.URN.Type())
				name, _ = openapi.Pluck(redisSlaveService.Outputs, "metadata", "name")
				assert.Equal(t, "redis-slave", name)
				status, _ = openapi.Pluck(redisSlaveService.Outputs, "spec", "clusterIP")
				assert.True(t, len(status.(string)) > 1)

				// Verify the provider resource.
				provRes := stackInfo.Deployment.Resources[7]
				assert.True(t, providers.IsProviderType(provRes.URN.Type()))

				// Verify root resource.
				stackRes := stackInfo.Deployment.Resources[8]
				assert.Equal(t, resource.RootStackType, stackRes.URN.Type())
			},
		})

	integration.ProgramTest(t, &test)
}

func TestAccIngress(t *testing.T) {
	skipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:         filepath.Join(getCwd(t), "ingress"),
			Quick:       true,
			SkipRefresh: true, // ingress may have changes during refresh.
			ExtraRuntimeValidation: func(
				t *testing.T, stackInfo integration.RuntimeValidationStackInfo,
			) {
				assert.NotNil(t, stackInfo.Deployment)
				assert.Equal(t, 10, len(stackInfo.Deployment.Resources))

				sort.Slice(stackInfo.Deployment.Resources, func(i, j int) bool {
					ri := stackInfo.Deployment.Resources[i]
					rj := stackInfo.Deployment.Resources[j]
					riname, _ := openapi.Pluck(ri.Outputs, "metadata", "name")
					rinamespace, _ := openapi.Pluck(ri.Outputs, "metadata", "namespace")
					rjname, _ := openapi.Pluck(rj.Outputs, "metadata", "name")
					rjnamespace, _ := openapi.Pluck(rj.Outputs, "metadata", "namespace")
					return fmt.Sprintf("%s/%s/%s", ri.URN.Type(), rinamespace, riname) <
						fmt.Sprintf("%s/%s/%s", rj.URN.Type(), rjnamespace, rjname)
				})

				var name interface{}
				var status interface{}

				// Verify frontend deployment.
				frontendDepl := stackInfo.Deployment.Resources[0]
				assert.Equal(t, tokens.Type("kubernetes:apps/v1:Deployment"), frontendDepl.URN.Type())
				name, _ = openapi.Pluck(frontendDepl.Outputs, "metadata", "name")
				assert.Equal(t, "frontend", name)
				status, _ = openapi.Pluck(frontendDepl.Outputs, "status", "readyReplicas")
				assert.Equal(t, float64(3), status)

				// Verify redis-master deployment.
				redisMasterDepl := stackInfo.Deployment.Resources[1]
				assert.Equal(t, tokens.Type("kubernetes:apps/v1:Deployment"), redisMasterDepl.URN.Type())
				name, _ = openapi.Pluck(redisMasterDepl.Outputs, "metadata", "name")
				assert.Equal(t, "redis-master", name)
				status, _ = openapi.Pluck(redisMasterDepl.Outputs, "status", "readyReplicas")
				assert.Equal(t, float64(1), status)

				// Verify redis-slave deployment.
				redisSlaveDepl := stackInfo.Deployment.Resources[2]
				assert.Equal(t, tokens.Type("kubernetes:apps/v1:Deployment"), redisSlaveDepl.URN.Type())
				name, _ = openapi.Pluck(redisSlaveDepl.Outputs, "metadata", "name")
				assert.Equal(t, "redis-slave", name)
				status, _ = openapi.Pluck(redisSlaveDepl.Outputs, "status", "readyReplicas")
				assert.Equal(t, float64(1), status)

				// Verify test namespace.
				namespace := stackInfo.Deployment.Resources[3]
				assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())

				// Verify frontend service.
				frontentService := stackInfo.Deployment.Resources[4]
				assert.Equal(t, tokens.Type("kubernetes:core/v1:Service"), frontentService.URN.Type())
				name, _ = openapi.Pluck(frontentService.Outputs, "metadata", "name")
				assert.Equal(t, "frontend", name)

				// Verify redis-master service.
				redisMasterService := stackInfo.Deployment.Resources[5]
				assert.Equal(t, tokens.Type("kubernetes:core/v1:Service"), redisMasterService.URN.Type())
				name, _ = openapi.Pluck(redisMasterService.Outputs, "metadata", "name")
				assert.Equal(t, "redis-master", name)
				status, _ = openapi.Pluck(redisMasterService.Outputs, "spec", "clusterIP")
				assert.True(t, len(status.(string)) > 1)

				// Verify redis-slave service.
				redisSlaveService := stackInfo.Deployment.Resources[6]
				assert.Equal(t, tokens.Type("kubernetes:core/v1:Service"), redisSlaveService.URN.Type())
				name, _ = openapi.Pluck(redisSlaveService.Outputs, "metadata", "name")
				assert.Equal(t, "redis-slave", name)
				status, _ = openapi.Pluck(redisSlaveService.Outputs, "spec", "clusterIP")
				assert.True(t, len(status.(string)) > 1)

				// TODO: revisit this as part of https://github.com/pulumi/pulumi-kubernetes/issues/1649.
				// Verify the ingress endpoint is accessible
				// integration.AssertHTTPResultWithRetry(t, fmt.Sprintf("%s/index.html", stackInfo.Outputs["ingressIP"]), nil, 5*time.Minute, func(body string) bool {
				//	return assert.NotEmpty(t, body, "Body should not be empty")
				//})
			},
		})

	integration.ProgramTest(t, &test)
}

func TestAccHelm(t *testing.T) {
	skipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:         filepath.Join(getCwd(t), "helm", "step1"),
			SkipRefresh: true,
			Verbose:     true,
			ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
				// Ensure that all `Services` have `status` marked as a `Secret`
				for _, res := range stackInfo.Deployment.Resources {
					if res.Type == "kubernetes:core/v1:Service" {
						spec, has := res.Outputs["status"]
						assert.True(t, has)
						specMap, is := spec.(map[string]interface{})
						assert.True(t, is)
						sigKey, has := specMap[resource.SigKey]
						assert.True(t, has)
						assert.Equal(t, resource.SecretSig, sigKey)
					}
				}
			},
			EditDirs: []integration.EditDir{
				{
					Dir:             filepath.Join(getCwd(t), "helm", "step2"),
					Additive:        true,
					ExpectNoChanges: true,
				},
			},
		})

	integration.ProgramTest(t, &test)
}

func TestAccHelmApiVersions(t *testing.T) {
	skipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:         filepath.Join(getCwd(t), "helm-api-versions", "step1"),
			SkipRefresh: true,
			ExtraRuntimeValidation: func(
				t *testing.T, stackInfo integration.RuntimeValidationStackInfo,
			) {
				assert.NotNil(t, stackInfo.Deployment)
				assert.Equal(t, 7, len(stackInfo.Deployment.Resources))
			},
			EditDirs: []integration.EditDir{
				{
					Dir:             filepath.Join(getCwd(t), "helm-api-versions", "step2"),
					Additive:        true,
					ExpectNoChanges: true,
				},
			},
		})

	integration.ProgramTest(t, &test)
}

func TestAccHelmAllowCRDRendering(t *testing.T) {
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:         filepath.Join("helm-skip-crd-rendering", "step1"),
			Quick:       true,
			SkipRefresh: true,
			ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
				assert.NotNil(t, stackInfo.Deployment)
				assert.Equal(t, 8, len(stackInfo.Deployment.Resources))

				for _, res := range stackInfo.Deployment.Resources {
					if res.Type == "kubernetes:core/v1:Pod" {
						annotations, ok := openapi.Pluck(res.Inputs, "metadata", "annotations")
						if strings.Contains(res.ID.String(), "skip-crd") {
							assert.False(t, ok)
						} else {
							assert.True(t, ok)
							assert.Contains(t, annotations, "pulumi.com/skipAwait")
						}
					}
				}
			},
		})
	integration.ProgramTest(t, &test)
}

func TestAccHelmLocal(t *testing.T) {
	skipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:                  filepath.Join(getCwd(t), "helm-local", "step1"),
			ExpectRefreshChanges: true,
			ExtraRuntimeValidation: func(
				t *testing.T, stackInfo integration.RuntimeValidationStackInfo,
			) {
				assert.NotNil(t, stackInfo.Deployment)
				assert.Equal(t, 12, len(stackInfo.Deployment.Resources))

				// Verify resource creation order using the Event stream. The Chart resources must be created
				// first, followed by the dependent ConfigMap. (The ConfigMap doesn't actually need the Chart, but
				// it creates almost instantly, so it's a good choice to test creation ordering)
				cmRegex := regexp.MustCompile(`ConfigMap::test-.*/nginx-server-block`)
				svcRegex := regexp.MustCompile(`Service::test-.*/nginx`)
				deployRegex := regexp.MustCompile(`Deployment::test-.*/nginx`)
				dependentRegex := regexp.MustCompile(`ConfigMap::foo`)

				var configmapFound, serviceFound, deploymentFound, dependentFound bool
				for _, e := range stackInfo.Events {
					if e.ResOutputsEvent != nil {
						switch {
						case cmRegex.MatchString(e.ResOutputsEvent.Metadata.URN):
							configmapFound = true
						case svcRegex.MatchString(e.ResOutputsEvent.Metadata.URN):
							serviceFound = true
						case deployRegex.MatchString(e.ResOutputsEvent.Metadata.URN):
							deploymentFound = true
						case dependentRegex.MatchString(e.ResOutputsEvent.Metadata.URN):
							dependentFound = true
						}
						assert.Falsef(t, dependentFound && !(configmapFound && serviceFound && deploymentFound),
							"dependent ConfigMap created before all chart resources were ready")
						fmt.Println(e.ResOutputsEvent.Metadata.URN)
					}
				}
			},
			EditDirs: []integration.EditDir{
				{
					Dir:             filepath.Join(getCwd(t), "helm-local", "step2"),
					Additive:        true,
					ExpectNoChanges: true,
				},
			},
		})

	integration.ProgramTest(t, &test)
}

func TestAccPrometheusOperator(t *testing.T) {
	skipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:         filepath.Join(getCwd(t), "prometheus-operator"),
			SkipRefresh: true,
			ExtraRuntimeValidation: func(
				t *testing.T, stackInfo integration.RuntimeValidationStackInfo,
			) {
				assert.NotNil(t, stackInfo.Deployment)
				assert.Equal(t, 10, len(stackInfo.Deployment.Resources))
			},
			EditDirs: []integration.EditDir{
				{
					Dir:      filepath.Join(getCwd(t), "prometheus-operator", "step1"),
					Additive: true,
					ExtraRuntimeValidation: func(
						t *testing.T, stackInfo integration.RuntimeValidationStackInfo,
					) {
						assert.NotNil(t, stackInfo.Deployment)
						assert.Equal(t, 10, len(stackInfo.Deployment.Resources))
					},
				},
			},
		})

	integration.ProgramTest(t, &test)
}

func TestAccMariadb(t *testing.T) {
	skipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "mariadb"),
		})

	integration.ProgramTest(t, &test)
}

func TestAccProvider(t *testing.T) {
	skipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "provider"),
		})

	integration.ProgramTest(t, &test)
}

func skipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}
}

func getCwd(t *testing.T) string {
	cwd, err := os.Getwd()
	if err != nil {
		t.FailNow()
	}

	return cwd
}

func getBaseOptions(t *testing.T) integration.ProgramTestOptions {
	return integration.ProgramTestOptions{
		Dependencies: []string{
			"@pulumi/kubernetes",
		},
	}
}
