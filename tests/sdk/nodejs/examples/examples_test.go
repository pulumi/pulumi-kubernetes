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

package examples

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/tests/v4"
	"github.com/pulumi/pulumi/pkg/v3/resource/deploy/providers"
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/stretchr/testify/assert"
)

func TestAccMinimal(t *testing.T) {
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "minimal"),
		})

	integration.ProgramTest(t, &test)
}

func TestAccGuestbook(t *testing.T) {
	tests.SkipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:                  filepath.Join(getCwd(t), "guestbook"),
			ExpectRefreshChanges: true,
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

				var name any
				var status any

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
	tests.SkipIfShort(t)
	testNetworkingV1 := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:           filepath.Join(getCwd(t), "ingress"),
			Quick:         true,
			NoParallel:    true, // We want to run this and the next test serially so nginx ingress isn't clobbered.
			DebugLogLevel: 3,
			SkipRefresh:   true, // ingress may have changes during refresh.
			ExtraRuntimeValidation: func(
				t *testing.T, stackInfo integration.RuntimeValidationStackInfo,
			) {
				assert.NotNil(t, stackInfo.Deployment)
				assert.Equal(t, 15, len(stackInfo.Deployment.Resources))

				integration.AssertHTTPResultWithRetry(t,
					fmt.Sprintf("%s/index.html", stackInfo.Outputs["ingressIp"]),
					nil, 10*time.Minute, func(body string) bool {
						return assert.NotEmpty(t, body, "Body should not be empty")
					})

				integration.AssertHTTPResultWithRetry(t, fmt.Sprintf("%s/hello", stackInfo.Outputs["ingressNginxIp"]),
					map[string]string{"Host": "ingresshello.io"}, 10*time.Minute, func(body string) bool {
						return assert.NotEmpty(t, body, "Body should not be empty")
					})
			},
		})
	integration.ProgramTest(t, &testNetworkingV1)
}

func TestAccHelm(t *testing.T) {
	tests.SkipIfShort(t)
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
						specMap, is := spec.(map[string]any)
						assert.True(t, is)
						sigKey, has := specMap[resource.SigKey]
						assert.True(t, has)
						assert.Equal(t, resource.SecretSig, sigKey)
					}
				}
			},
		})

	integration.ProgramTest(t, &test)
}

func TestHelmNoDefaultProvider(t *testing.T) {
	tests.SkipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:         filepath.Join(getCwd(t), "helm-no-default-provider"),
			SkipRefresh: true,
			Verbose:     true,
			Quick:       true,
			Config:      map[string]string{"disable-default-providers": `["kubernetes"]`},
		})

	integration.ProgramTest(t, &test)
}

func TestAccHelmApiVersions(t *testing.T) {
	tests.SkipIfShort(t)
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
		})

	integration.ProgramTest(t, &test)
}

func TestAccHelmKubeVersion(t *testing.T) {
	tests.SkipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:         filepath.Join(getCwd(t), "helm-kube-version", "step1"),
			SkipRefresh: true,
			ExtraRuntimeValidation: func(
				t *testing.T, stackInfo integration.RuntimeValidationStackInfo,
			) {
				assert.NotNil(t, stackInfo.Deployment)
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
	tests.SkipIfShort(t)
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
		})

	integration.ProgramTest(t, &test)
}

func testAccPrometheusOperator(t *testing.T) {
	tests.SkipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:         filepath.Join(getCwd(t), "prometheus-operator"),
			SkipRefresh: true,
			NoParallel:  true,
			OrderedConfig: []integration.ConfigValue{
				{
					Key:   "kubernetes:enableServerSideApply",
					Value: "true",
				},
			},
			ExtraRuntimeValidation: func(
				t *testing.T, stackInfo integration.RuntimeValidationStackInfo,
			) {
				assert.NotNil(t, stackInfo.Deployment)
				assert.Equal(t, 18, len(stackInfo.Deployment.Resources))
			},
			EditDirs: []integration.EditDir{
				{
					Dir:      filepath.Join(getCwd(t), "prometheus-operator", "step1"),
					Additive: true,
					ExtraRuntimeValidation: func(
						t *testing.T, stackInfo integration.RuntimeValidationStackInfo,
					) {
						assert.NotNil(t, stackInfo.Deployment)
						assert.Equal(t, 18, len(stackInfo.Deployment.Resources))
					},
				},
			},
		})

	integration.ProgramTest(t, &test)
}

// TODO: Fix this example.
//func TestAccMariadb(t *testing.T) {
//	tests.SkipIfShort(t)
//	test := getBaseOptions(t).
//		With(integration.ProgramTestOptions{
//			Dir: filepath.Join(getCwd(t), "mariadb"),
//		})
//
//	integration.ProgramTest(t, &test)
//}

func TestAccProvider(t *testing.T) {
	tests.SkipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir: filepath.Join(getCwd(t), "provider"),
		})

	integration.ProgramTest(t, &test)
}

func TestHelmRelease(t *testing.T) {
	tests.SkipIfShort(t)
	validationFunc := func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
		assert.NotEmpty(t, stackInfo.Outputs["redisMasterClusterIP"].(string))
		assert.Equal(t, stackInfo.Outputs["status"], "deployed")
		for _, res := range stackInfo.Deployment.Resources {
			if res.Type == "kubernetes:helm.sh/v3:Release" {
				version, has := res.Inputs["version"]
				assert.True(t, has)
				stat, has := res.Outputs["status"]
				assert.True(t, has)
				specMap, is := stat.(map[string]any)
				assert.True(t, is)
				versionOut, has := specMap["version"]
				assert.True(t, has)
				assert.Equal(t, version, versionOut)
				values, has := res.Outputs["values"]
				assert.True(t, has)
				assert.Contains(t, values, "cluster")
				valMap := values.(map[string]any)
				assert.Equal(t, valMap["cluster"], map[string]any{
					"enabled":    true,
					"slaveCount": float64(2),
				})
				// not asserting contents since the secret is hard to assert equality on.
				assert.Contains(t, values, "global")
				assert.Contains(t, values, "metrics")
				assert.Equal(t, valMap["metrics"], map[string]any{
					"enabled": true,
					"service": map[string]any{
						"annotations": map[string]any{
							"prometheus.io/port": "9127",
						},
					},
				})
				assert.Contains(t, values, "rbac")
				assert.Equal(t, valMap["rbac"], map[string]any{
					"create": true,
				})
			}
		}
	}
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:                    filepath.Join(getCwd(t), "helm-release", "step1"),
			SkipRefresh:            false,
			Verbose:                true,
			ExtraRuntimeValidation: validationFunc,
			ExpectRefreshChanges:   true,
			EditDirs: []integration.EditDir{
				{
					Dir:                    filepath.Join(getCwd(t), "helm-release", "step2"),
					Additive:               true,
					ExtraRuntimeValidation: validationFunc,
				},
			},
		})

	integration.ProgramTest(t, &test)
}

func TestHelmReleaseCRD(t *testing.T) {
	// Validate that Helm charts with CRDs work across create/update/refresh/delete cycles.
	// https://github.com/pulumi/pulumi-kubernetes/issues/1712
	tests.SkipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:                  filepath.Join(getCwd(t), "helm-release-crd", "step1"),
			SkipRefresh:          false,
			ExpectRefreshChanges: true,
			Verbose:              true,
			EditDirs: []integration.EditDir{
				{
					Dir:      filepath.Join(getCwd(t), "helm-release-crd", "step2"),
					Additive: true,
				},
			},
		})

	integration.ProgramTest(t, &test)
}

func TestHelmReleaseNamespace(t *testing.T) {
	// Validate fix for https://github.com/pulumi/pulumi-kubernetes/issues/1710
	tests.SkipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:                  filepath.Join(getCwd(t), "helm-release-namespace", "step1"),
			SkipRefresh:          false,
			Verbose:              true,
			ExpectRefreshChanges: true,
			// Ensure that the rule was found in the release's namespace.
			ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
				assert.Equalf(t, stackInfo.Outputs["namespaceName"], stackInfo.Outputs["alertManagerNamespace"],
					"expected Helm resources to reside in the provided Namespace")
				assert.NotEmptyf(t, stackInfo.Outputs["alertManagerNamespace"].(string),
					"Helm resources should not be in the default Namespace")
				getResponse, err := tests.Kubectl(fmt.Sprintf("get statefulset -n %s alertmanager",
					stackInfo.Outputs["alertManagerNamespace"]))
				assert.NoErrorf(t, err, "kubectl command failed")
				assert.NotContainsf(t, getResponse, "No resources found",
					"kubectl did not find the expected resource")
			},
			EditDirs: []integration.EditDir{
				{
					Dir:      filepath.Join(getCwd(t), "helm-release-namespace", "step2"),
					Additive: true,
					// Ensure that the rule was found in the release's namespace.
					ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
						assert.Equalf(t, stackInfo.Outputs["namespaceName"], stackInfo.Outputs["alertManagerNamespace"],
							"expected Helm resources to reside in the provided Namespace")
						assert.NotEmptyf(t, stackInfo.Outputs["alertManagerNamespace"].(string),
							"Helm resources should not be in the default Namespace")
						getResponse, err := tests.Kubectl(fmt.Sprintf("get statefulset -n %s alertmanager",
							stackInfo.Outputs["alertManagerNamespace"]))
						assert.NoErrorf(t, err, "kubectl command failed")
						assert.NotContainsf(t, getResponse, "No resources found",
							"kubectl did not find the expected resource")
					},
				},
			},
		})

	integration.ProgramTest(t, &test)
}

// TestHelmReleaseProviderNamespace tests how Helm Release inherits provider namespace.
func TestHelmReleaseProviderNamespace(t *testing.T) {
	tests.SkipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:         filepath.Join(getCwd(t), "helm-release-provider-namespace"),
			SkipRefresh: true,
			Verbose:     true,
			Quick:       true,
			ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
				assert.NotNil(t, stackInfo.Outputs["providerNamespace"])
				assert.NotNil(t, stackInfo.Outputs["alertManagerNamespace"])
				assert.Equal(t, stackInfo.Outputs["providerNamespace"], stackInfo.Outputs["alertManagerNamespace"])
			},
		})

	integration.ProgramTest(t, &test)
}

func TestHelmReleaseRedis(t *testing.T) {
	expectKeyringInput := func(verifyVal bool, keyRingNonEmpty bool) func(t *testing.T,
		stackInfo integration.RuntimeValidationStackInfo) {
		// Validate that the keyring is omitted when verify is false/unspecified.
		// https://github.com/pulumi/pulumi-kubernetes/issues/1959
		return func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			seen := false
			for _, res := range stackInfo.Deployment.Resources {
				if res.Type == "kubernetes:helm.sh/v3:Release" {
					seen = true

					assert.Contains(t, res.Inputs, "verify")
					verify := res.Inputs["verify"].(bool)
					assert.Equal(t, verifyVal, verify)
					val := res.Inputs["keyring"]
					if keyRingNonEmpty {
						assert.NotEmpty(t, val)
					} else {
						assert.Empty(t, val)
					}

				}
			}
			assert.True(t, seen)
		}
	}

	// Validate fix for https://github.com/pulumi/pulumi-kubernetes/issues/1933
	tests.SkipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:                  filepath.Join(getCwd(t), "helm-release-redis", "step1"),
			SkipRefresh:          false,
			ExpectRefreshChanges: true,
			Verbose:              true,
			Quick:                true,
			EditDirs: []integration.EditDir{
				{
					Dir:      filepath.Join(getCwd(t), "helm-release-redis", "step2"),
					Additive: true,
					// The redis chart isn't signed so can't find provenance file for it.
					// TODO: Add a separate test for chart verification.
					ExtraRuntimeValidation: expectKeyringInput(false, false),
				},
			},
			ExtraRuntimeValidation: expectKeyringInput(false, false),
		})

	integration.ProgramTest(t, &test)
}

func testRancher(t *testing.T) {
	// Validate fix for https://github.com/pulumi/pulumi-kubernetes/issues/1848
	tests.SkipIfShort(t)
	test := getBaseOptions(t).
		With(integration.ProgramTestOptions{
			Dir:         filepath.Join(getCwd(t), "rancher", "step1"),
			SkipRefresh: true,
			Verbose:     true,
			NoParallel:  true,
			OrderedConfig: []integration.ConfigValue{
				{
					Key:   "kubernetes:enableServerSideApply",
					Value: "true",
				},
			},
			EditDirs: []integration.EditDir{
				{
					Dir:      filepath.Join(getCwd(t), "rancher", "step2"),
					Additive: true,
				},
			},
		})
	integration.ProgramTest(t, &test)
}

// TestCRDs runs 2 sub tests that cannot be parallelized as they touch
// the same cluster-scoped CRD. This is required until we can run tests
// in parallel with different clusters (tracked by: https://github.com/pulumi/pulumi-kubernetes/issues/2243).
func TestCRDs(t *testing.T) {
	t.Run("testAccPrometheusOperator", testAccPrometheusOperator)
	t.Run("testRancher", testRancher)
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
		Env: []string{
			"PULUMI_K8S_CLIENT_BURST=200",
			"PULUMI_K8S_CLIENT_QPS=100",
		},
	}
}
