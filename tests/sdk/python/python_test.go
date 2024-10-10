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

package python

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	. "github.com/onsi/gomega/gstruct"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/tests/v4"
	. "github.com/pulumi/pulumi-kubernetes/tests/v4/gomega"
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
)

var baseOptions = &integration.ProgramTestOptions{
	Verbose: true,
	Dependencies: []string{
		filepath.Join("..", "..", "..", "sdk", "python", "bin"),
	},
	PostPrepareProject: func(p *engine.Projinfo) error {
		return fsutil.CopyFile(filepath.Join(p.Root, "testdata"), filepath.Join("..", "..", "testdata"), nil)
	},
	Env: []string{
		"PULUMI_K8S_CLIENT_BURST=200",
		"PULUMI_K8S_CLIENT_QPS=100",
	},
}

func TestSmoke(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	for _, dir := range []string{"smoke-test", "smoke-test-old"} {
		t.Run(dir, func(t *testing.T) {
			options := baseOptions.With(integration.ProgramTestOptions{
				Dir:        filepath.Join(cwd, dir),
				NoParallel: true,
			})
			integration.ProgramTest(t, &options)
		})
	}
}

// Smoke test for .get support.
func TestGetPython(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	for _, dir := range []string{"get-old"} {
		t.Run(dir, func(t *testing.T) {
			options := baseOptions.With(integration.ProgramTestOptions{
				ExpectRefreshChanges: true, // CRD changes on refresh
				Dir:                  filepath.Join(cwd, dir, "step1"),
				NoParallel:           true,
				EditDirs: []integration.EditDir{
					{
						Dir:      filepath.Join(cwd, dir, "step2"),
						Additive: true,
					},
				},
			})
			integration.ProgramTest(t, &options)
		})
	}
}

func TestGetOneStep(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	for _, dir := range []string{"get-one-step"} {
		t.Run(dir, func(t *testing.T) {
			options := baseOptions.With(integration.ProgramTestOptions{
				ExpectRefreshChanges: true, // CRD changes on refresh
				Dir:                  filepath.Join(cwd, dir),
				NoParallel:           true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
					var success bool
					for _, res := range stackInfo.Deployment.Resources {
						if res.URN.Type() == "kubernetes:python.test/v1:GetTest" {
							nodeSelector, _ := openapi.Pluck(res.Outputs, "spec", "node_selector")
							assert.Equal(t, "kubernetes.io/hostname: \"docker-desktop\"\n", nodeSelector)
							foo, _ := openapi.Pluck(res.Outputs, "spec", "foo")
							assert.Equal(t, "bar", foo)
							apiVersion, _ := openapi.Pluck(res.Outputs, "apiVersion")
							assert.Equal(t, "python.test/v1", apiVersion)
							success = true
						}
					}
					assert.True(t, success)
				},
			})
			integration.ProgramTest(t, &options)
		})
	}
}

func TestYamlPython(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join(cwd, "yaml-test"),
		ExpectRefreshChanges: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 20, len(stackInfo.Deployment.Resources))

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
			var ns any
			var namespaceName, namespace2Name string

			// Verify CRD.
			crd := stackInfo.Deployment.Resources[0]
			assert.Equal(t, tokens.Type("kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition"),
				crd.URN.Type())
			name, _ = openapi.Pluck(crd.Outputs, "metadata", "name")
			assert.True(t, strings.HasPrefix(name.(string), "foos.bar.example.com"))

			// Verify CR.
			cr := stackInfo.Deployment.Resources[3]
			assert.Equal(t, tokens.Type("kubernetes:bar.example.com/v1:Foo"), cr.URN.Type())
			name, _ = openapi.Pluck(cr.Outputs, "metadata", "name")
			assert.True(t, strings.HasPrefix(name.(string), "foobar"))

			// Verify namespace1.
			namespace := stackInfo.Deployment.Resources[7]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())
			name, _ = openapi.Pluck(namespace.Outputs, "metadata", "name")
			namespaceName = name.(string)
			assert.True(t, strings.HasPrefix(namespaceName, "ns"), fmt.Sprintf("%s %s", name, "ns"))

			// Verify namespace2.
			namespace2 := stackInfo.Deployment.Resources[8]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace2.URN.Type())
			name, _ = openapi.Pluck(namespace2.Outputs, "metadata", "name")
			namespace2Name = name.(string)
			assert.True(t, strings.HasPrefix(namespace2Name, "ns2"), fmt.Sprintf("%s %s", name, "ns2"))

			// Verify Pod "bar".
			podBar := stackInfo.Deployment.Resources[9]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Pod"), podBar.URN.Type())
			name, _ = openapi.Pluck(podBar.Outputs, "metadata", "name")
			assert.True(t, strings.Contains(name.(string), "bar"), fmt.Sprintf("%s %s", name, "bar"))
			ns, _ = openapi.Pluck(podBar.Outputs, "metadata", "namespace")
			assert.Equal(t, ns, namespaceName)

			// Verify Pod "baz".
			podBaz := stackInfo.Deployment.Resources[10]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Pod"), podBaz.URN.Type())
			name, _ = openapi.Pluck(podBaz.Outputs, "metadata", "name")
			assert.True(t, strings.Contains(name.(string), "baz"), fmt.Sprintf("%s %s", name, "baz"))
			ns, _ = openapi.Pluck(podBaz.Outputs, "metadata", "namespace")
			assert.Equal(t, ns, namespaceName)

			// Verify Pod "foo".
			podFoo := stackInfo.Deployment.Resources[11]
			assert.Equal(t, tokens.Type("kubernetes:core/v1:Pod"), podFoo.URN.Type())
			name, _ = openapi.Pluck(podFoo.Outputs, "metadata", "name")
			assert.True(t, strings.Contains(name.(string), "foo"), fmt.Sprintf("%s %s", name, "foo"))
			ns, _ = openapi.Pluck(podFoo.Outputs, "metadata", "namespace")
			assert.Equal(t, ns, namespaceName)

			// Note: Skipping validation for the guestbook app in this test since it's no different from the
			// first ConfigFile.

			// Verify the provider resources.
			provRes := stackInfo.Deployment.Resources[18]
			assert.True(t, providers.IsProviderType(provRes.URN.Type()))

			// Verify root resource.
			stackRes := stackInfo.Deployment.Resources[19]
			assert.Equal(t, resource.RootStackType, stackRes.URN.Type())

			// Ensure that all `Pod`s have `apiVersion` marked as a `Secret`
			for _, res := range stackInfo.Deployment.Resources {
				if res.Type == "kubernetes:core/v1:Pod" {
					spec, has := res.Outputs["apiVersion"]
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
	integration.ProgramTest(t, &options)
}

// Regression Test for https://github.com/pulumi/pulumi-kubernetes/issues/2664.
// Ensure the program runs without an error being raised when an invoke is called
// using a provider that is not configured.
func TestYamlUnconfiguredProvider(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join(cwd, "yaml-test-unconfigured-provider"),
		ExpectRefreshChanges: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
	})
	integration.ProgramTest(t, &options)
}

func TestGuestbook(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	for _, dir := range []string{"guestbook", "guestbook-old"} {
		t.Run(dir, func(t *testing.T) {
			options := baseOptions.With(integration.ProgramTestOptions{
				Dir:                  filepath.Join(cwd, dir),
				NoParallel:           true,
				ExpectRefreshChanges: true,
				ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
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
					assert.True(t, strings.HasPrefix(name.(string), "frontend"))
					status, _ = openapi.Pluck(frontendDepl.Outputs, "status", "readyReplicas")
					assert.Equal(t, float64(3), status)

					// Verify redis-follower deployment.
					redisFollowerDepl := stackInfo.Deployment.Resources[1]
					assert.Equal(t, tokens.Type("kubernetes:apps/v1:Deployment"), redisFollowerDepl.URN.Type())
					name, _ = openapi.Pluck(redisFollowerDepl.Outputs, "metadata", "name")
					assert.True(t, strings.HasPrefix(name.(string), "redis-follower"), fmt.Sprintf("%s %s", name, "redis-slave"))
					status, _ = openapi.Pluck(redisFollowerDepl.Outputs, "status", "readyReplicas")
					assert.Equal(t, float64(1), status)

					// Verify redis-leader deployment.
					redisLeaderDepl := stackInfo.Deployment.Resources[2]
					assert.Equal(t, tokens.Type("kubernetes:apps/v1:Deployment"), redisLeaderDepl.URN.Type())
					name, _ = openapi.Pluck(redisLeaderDepl.Outputs, "metadata", "name")
					assert.True(t, strings.HasPrefix(name.(string), "redis-leader"), fmt.Sprintf("%s %s", name, "redis-master"))
					status, _ = openapi.Pluck(redisLeaderDepl.Outputs, "status", "readyReplicas")
					assert.Equal(t, float64(1), status)

					// Verify test namespace.
					namespace := stackInfo.Deployment.Resources[3]
					assert.Equal(t, tokens.Type("kubernetes:core/v1:Namespace"), namespace.URN.Type())
					name, _ = openapi.Pluck(namespace.Outputs, "metadata", "name")
					assert.True(t, strings.HasPrefix(name.(string), "test"), fmt.Sprintf("%s %s", name, "test"))

					// Verify frontend service.
					frontendService := stackInfo.Deployment.Resources[4]
					assert.Equal(t, tokens.Type("kubernetes:core/v1:Service"), frontendService.URN.Type())
					name, _ = openapi.Pluck(frontendService.Outputs, "metadata", "name")
					assert.True(t, strings.HasPrefix(name.(string), "frontend"), fmt.Sprintf("%s %s", name, "frontend"))
					status, _ = openapi.Pluck(frontendService.Outputs, "spec", "clusterIP")
					assert.True(t, len(status.(string)) > 1)

					// Verify redis-follower service.
					redisFollowerService := stackInfo.Deployment.Resources[5]
					assert.Equal(t, tokens.Type("kubernetes:core/v1:Service"), redisFollowerService.URN.Type())
					name, _ = openapi.Pluck(redisFollowerService.Outputs, "metadata", "name")
					assert.True(t, strings.HasPrefix(name.(string), "redis-follower"), fmt.Sprintf("%s %s", name, "redis-slave"))
					status, _ = openapi.Pluck(redisFollowerService.Outputs, "spec", "clusterIP")
					assert.True(t, len(status.(string)) > 1)

					// Verify redis-leader service.
					redisLeaderService := stackInfo.Deployment.Resources[6]
					assert.Equal(t, tokens.Type("kubernetes:core/v1:Service"), redisLeaderService.URN.Type())
					name, _ = openapi.Pluck(redisLeaderService.Outputs, "metadata", "name")
					assert.True(t, strings.HasPrefix(name.(string), "redis-leader"), fmt.Sprintf("%s %s", name, "redis-master"))
					status, _ = openapi.Pluck(redisLeaderService.Outputs, "spec", "clusterIP")
					assert.True(t, len(status.(string)) > 1)

					// Verify the provider resource.
					provRes := stackInfo.Deployment.Resources[7]
					assert.True(t, providers.IsProviderType(provRes.URN.Type()))

					// Verify root resource.
					stackRes := stackInfo.Deployment.Resources[8]
					assert.Equal(t, resource.RootStackType, stackRes.URN.Type())
				},
			})
			integration.ProgramTest(t, &options)
		})
	}

}

// Smoke test for first-class Kubernetes providers.
func TestProviderPython(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	for _, dir := range []string{"provider", "provider-old"} {
		t.Run(dir, func(t *testing.T) {
			options := baseOptions.With(integration.ProgramTestOptions{
				Dir:        filepath.Join(cwd, "provider"),
				NoParallel: true,
			})
			integration.ProgramTest(t, &options)
		})
	}
}

func TestHelm(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join(cwd, "helm", "step1"),
		ExpectRefreshChanges: true,
	})
	integration.ProgramTest(t, &options)
}

func TestHelmReleasePython(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join(cwd, "helm-release", "step1"),
		ExpectRefreshChanges: true,
		EditDirs: []integration.EditDir{
			{
				Dir:             filepath.Join(cwd, "helm-release", "step2"),
				Additive:        true,
				ExpectNoChanges: true,
			},
		},
	})
	integration.ProgramTest(t, &options)
}

func TestHelmLocalPython(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join(cwd, "helm-local", "step1"),
		ExpectRefreshChanges: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			// Verify resource creation order using the Event stream. The Chart resources must be created
			// first, followed by the dependent ConfigMap. (The ConfigMap doesn't actually need the Chart, but
			// it creates almost instantly, so it's a good choice to test creation ordering)
			cmRegex := regexp.MustCompile(`ConfigMap::nginx-server-block`)
			svcRegex := regexp.MustCompile(`Service::nginx`)
			deployRegex := regexp.MustCompile(`Deployment::nginx`)
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
	integration.ProgramTest(t, &options)
}

// Regression Test for https://github.com/pulumi/pulumi-kubernetes/issues/2664.
// Ensure the program runs without an error being raised when an invoke is called
// using a provider that is not configured.
func TestHelmLocalUnconfiguredProvider(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir: filepath.Join(cwd, "helm-local-unconfigured-provider"),
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
			{
				Key:   "path",
				Value: filepath.Join(cwd, "..", "..", "testdata", "helm", "nginx"),
			},
		},
		ExpectRefreshChanges: true,
	})
	integration.ProgramTest(t, &options)
}

func TestHelmApiVersions(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join(cwd, "helm-api-versions", "step1"),
		ExpectRefreshChanges: true,
	})
	integration.ProgramTest(t, &options)
}

func TestHelmKubeVersionPython(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join(cwd, "helm-kube-version", "step1"),
		ExpectRefreshChanges: true,
	})
	integration.ProgramTest(t, &options)
}

func TestHelmAllowCRDRendering(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
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

func TestKustomizePython(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir: filepath.Join(cwd, "kustomize"),
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
	})
	integration.ProgramTest(t, &options)
}

// Regression Test for https://github.com/pulumi/pulumi-kubernetes/issues/2664.
// Ensure the program runs without an error being raised when an invoke is called
// using a provider that is not configured.
func TestKustomizeUnconfiguredProvider(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir: filepath.Join(cwd, "kustomize-unconfigured-provider"),
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
	})
	integration.ProgramTest(t, &options)
}

func TestSecretsPython(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	secretMessage := "secret message for testing"

	options := baseOptions.With(integration.ProgramTestOptions{
		Dir: filepath.Join(cwd, "secrets"),
		Config: map[string]string{
			"message": secretMessage,
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
	integration.ProgramTest(t, &options)
}

func TestServerSideApplyPython(t *testing.T) {
	cwd, err := os.Getwd()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	options := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join(cwd, "server-side-apply"),
		ExpectRefreshChanges: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			// Validate patched CustomResource
			crPatched := stackInfo.Outputs["crPatched"].(map[string]any)
			fooV, ok, err := unstructured.NestedString(crPatched, "metadata", "labels", "foo")
			assert.True(t, ok)
			assert.NoError(t, err)
			assert.Equal(t, "foo", fooV)
		},
	})
	integration.ProgramTest(t, &options)
}

// TestOptionPropagationPython tests the handling of resource options by the various compoonent resources.
// Component resources are responsible for implementing option propagation logic when creating
// child resources.
func TestOptionPropagationPython(t *testing.T) {
	g := NewWithT(t)
	format.MaxLength = 0
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
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
					}),
				}),
			))

			// --- ConfigGroup ---

			// ConfigGroup "cg-options" with most options
			g.Expect(rr.Named(stackInfo.RootResource.URN,
				"kubernetes:yaml:ConfigGroup", "cg-options-cg-options")).To(HaveExactElements(
				// quirk: Python SDK applies resource_prefix ("cg-options") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           ConsistOf(Alias("cg-options-old"), Alias("cg-options-cg-options-aliased")),
						"Protect":           BeTrue(),
						"Dependencies":      ConsistOf(string(sleep.URN)),
						"Provider":          BeEmpty(),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerA)),
						}),
						"IgnoreChanges": ConsistOf("ignored"),
					}),
				}),
			))
			g.Expect(rr.Named(urn("", "kubernetes:yaml:ConfigGroup", "cg-options-cg-options"),
				"kubernetes:core/v1:ConfigMap", "cg-options-cg-options-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           ConsistOf(Alias("cg-options-cm-1-k8s-aliased"), Alias("cg-options-cg-options-cm-1-aliased")),
						"Protect":           BeTrue(),
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
			g.Expect(rr.Named(urn("", "kubernetes:yaml:ConfigGroup", "cg-options-cg-options"),
				"kubernetes:yaml:ConfigFile", "cg-options-./testdata/options/configgroup/manifest.yaml")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           ConsistOf(Alias("cg-options-./testdata/options/configgroup/manifest.yaml-aliased")),
						"Protect":           BeTrue(),
						"Dependencies":      BeEmpty(),
						"Provider":          BeEmpty(),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerA)),
						}),
					}),
				}),
			))
			g.Expect(rr.Named(urn("kubernetes:yaml:ConfigGroup", "kubernetes:yaml:ConfigFile", "cg-options-./testdata/options/configgroup/manifest.yaml"),
				"kubernetes:core/v1:ConfigMap", "cg-options-configgroup-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           ConsistOf(Alias("configgroup-cm-1-k8s-aliased"), Alias("cg-options-configgroup-cm-1-aliased")),
						"Protect":           BeTrue(),
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
				"kubernetes:yaml:ConfigGroup", "cg-provider-cg-provider")).To(HaveExactElements(
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
				"kubernetes:yaml:ConfigGroup", "cg-nullopts-cg-nullopts")).To(HaveExactElements(
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
				// quirk: Python SDK applies resource_prefix ("cf-options") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           ConsistOf(Alias("cf-options-old"), Alias("cf-options-cf-options-aliased")),
						"Protect":           BeTrue(),
						"Dependencies":      ConsistOf(string(sleep.URN)),
						"Provider":          BeEmpty(),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerA)),
						}),
					}),
				}),
			))
			g.Expect(rr.Named(urn("", "kubernetes:yaml:ConfigFile", "cf-options-cf-options"),
				"kubernetes:core/v1:ConfigMap", "cf-options-configfile-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           ConsistOf(Alias("configfile-cm-1-k8s-aliased"), Alias("cf-options-configfile-cm-1-aliased")),
						"Protect":           BeTrue(),
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
				// quirk: Python SDK applies resource_prefix ("cf-provider") to the component itself.
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
				// quirk: Python SDK applies resource_prefix ("cf-nullopts") to the component itself.
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
				// quirk: Python SDK applies resource_prefix ("kustomize-options") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           ConsistOf(Alias("kustomize-options-old"), Alias("kustomize-options-kustomize-options-aliased")),
						"Protect":           BeTrue(),
						"Dependencies":      ConsistOf(string(sleep.URN)),
						"Provider":          BeEmpty(),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerA)),
						}),
						"IgnoreChanges": ConsistOf("ignored"),
					}),
				}),
			))
			g.Expect(rr.Named(urn("", "kubernetes:kustomize:Directory", "kustomize-options-kustomize-options"),
				"kubernetes:core/v1:ConfigMap", "kustomize-options-kustomize-cm-1-2kkk4bthmg")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           ConsistOf(Alias("kustomize-cm-1-2kkk4bthmg-k8s-aliased"), Alias("kustomize-options-kustomize-cm-1-2kkk4bthmg-aliased")),
						"Protect":           BeTrue(),
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
				// quirk: Python SDK applies resource_prefix ("kustomize-provider") to the component itself.
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
				// quirk: Python SDK applies resource_prefix ("kustomize-nullopts") to the component itself.
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
				// quirk: Python SDK applies resource_prefix ("chart-options") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases": ConsistOf(
							Alias("chart-options-old"),
							Alias("chart-options-chart-options-aliased"),
							Alias(tokens.Type("kubernetes:helm.sh/v2:Chart"))),
						"Protect":           BeTrue(),
						"Dependencies":      ConsistOf(string(sleep.URN)),
						"Provider":          BeEmpty(),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"Providers": MatchAllKeys(Keys{
							"kubernetes": BeEquivalentTo(providerUrn(providerA)),
						}),
						"IgnoreChanges": ConsistOf("ignored"),
					}),
				}),
			))
			g.Expect(rr.Named(urn("", "kubernetes:helm.sh/v3:Chart", "chart-options-chart-options"),
				"kubernetes:core/v1:ConfigMap", "chart-options-chart-options-cm-1")).To(HaveExactElements(
				// quirk: Python SDK does NOT directly apply resource_prefix to the child resources (unlike other SDKs)
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           ConsistOf(Alias("chart-options-chart-options-cm-1-k8s-aliased"), Alias("chart-options-chart-options-cm-1-aliased")),
						"Protect":           BeTrue(),
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
				// quirk: Python SDK applies resource_prefix ("chart-provider") to the component itself.
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
				"kubernetes:core/v1:ConfigMap", "chart-provider-chart-provider-cm-1")).To(HaveExactElements(
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
				// quirk: Python SDK applies resource_prefix ("chart-options") to the component itself.
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
