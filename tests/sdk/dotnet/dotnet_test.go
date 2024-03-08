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

package test

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	"github.com/pulumi/pulumi/pkg/v3/testing/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/fsutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var baseOptions = &integration.ProgramTestOptions{
	Verbose: true,
	Dependencies: []string{
		"Pulumi.Kubernetes",
	},
	PostPrepareProject: func(p *engine.Projinfo) error {
		return fsutil.CopyFile(filepath.Join(p.Root, "testdata"), filepath.Join("..", "..", "testdata"), nil)
	},
	Env: []string{
		"PULUMI_K8S_CLIENT_BURST=200",
		"PULUMI_K8S_CLIENT_QPS=100",
	},
}

func TestDotnet_Basic(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  "basic",
		Quick:                true,
		ExpectRefreshChanges: true, // The CRD sometimes, but not always, has changes during refresh.
	})
	test, _ = testClusters.WrapProviderTestOptions(test)
	integration.ProgramTest(t, &test)
}

func TestDotnet_Guestbook(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   "guestbook",
		Quick: true,
	})
	test, _ = testClusters.WrapProviderTestOptions(test)
	integration.ProgramTest(t, &test)
}

func TestDotnet_YamlUrl(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   "yaml-url",
		Quick: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		ExtraRuntimeValidation: func(
			t *testing.T, stackInfo integration.RuntimeValidationStackInfo,
		) {
			assert.NotNil(t, stackInfo.Deployment)
			assert.Equal(t, 18, len(stackInfo.Deployment.Resources))
		},
	})
	test, _ = testClusters.WrapProviderTestOptions(test)
	integration.ProgramTest(t, &test)
}

func TestDotnet_YamlLocal(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   "yaml-local",
		Quick: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
	})
	test, _ = testClusters.WrapProviderTestOptions(test)
	integration.ProgramTest(t, &test)
}

func TestDotnet_Helm(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("helm", "step1"),
		Quick:                true,
		ExpectRefreshChanges: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
			// Ensure that all `Services` have `status` marked as a `Secret`
			for _, res := range stackInfo.Deployment.Resources {
				if res.Type == tokens.Type("kubernetes:core/v1:Service") {
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
	test, _ = testClusters.WrapProviderTestOptions(test)
	integration.ProgramTest(t, &test)
}

func TestDotnet_HelmLocal(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("helm-local", "step1"),
		Quick:                true,
		ExpectRefreshChanges: true,
		ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
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
	test, _ = testClusters.WrapProviderTestOptions(test)
	integration.ProgramTest(t, &test)
}

func TestDotnet_HelmApiVersions(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("helm-api-versions", "step1"),
		Quick:                true,
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
			assert.Equal(t, 7, len(stackInfo.Deployment.Resources))
		},
	})
	test, _ = testClusters.WrapProviderTestOptions(test)
	integration.ProgramTest(t, &test)
}

func TestDotnet_HelmKubeVersion(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("helm-kube-version", "step1"),
		Quick:                true,
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
		},
	})
	test, _ = testClusters.WrapProviderTestOptions(test)
	integration.ProgramTest(t, &test)
}

func TestDotnet_HelmAllowCRDRendering(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  filepath.Join("helm-skip-crd-rendering", "step1"),
		Quick:                true,
		SkipRefresh:          true,
		ExpectRefreshChanges: true,
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
	test, _ = testClusters.WrapProviderTestOptions(test)
	integration.ProgramTest(t, &test)
}

func TestDotnet_CustomResource(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  "custom-resource",
		Quick:                true,
		ExpectRefreshChanges: true, // The CRD sometimes, but not always, has changes during refresh.
	})
	test, _ = testClusters.WrapProviderTestOptions(test)
	integration.ProgramTest(t, &test)
}

func TestDotnet_Kustomize(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   "kustomize",
		Quick: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
	})
	test, _ = testClusters.WrapProviderTestOptions(test)
	integration.ProgramTest(t, &test)
}

func TestDotnet_Secrets(t *testing.T) {
	secretMessage := "secret message for testing"

	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:   "secrets",
		Quick: true,
		Config: map[string]string{
			"message": secretMessage,
		},
		ExpectRefreshChanges: true,
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
	test, _ = testClusters.WrapProviderTestOptions(test)
	integration.ProgramTest(t, &test)
}

func TestDotnet_ServerSideApply(t *testing.T) {
	test := baseOptions.With(integration.ProgramTestOptions{
		Dir:                  "server-side-apply",
		ExpectRefreshChanges: true,
		OrderedConfig: []integration.ConfigValue{
			{
				Key:   "pulumi:disable-default-providers[0]",
				Value: "kubernetes",
				Path:  true,
			},
		},
		// TODO: Need to support CustomResource.Get to get the required info here.
		//ExtraRuntimeValidation: func(t *testing.T, stackInfo integration.RuntimeValidationStackInfo) {
		//	// Validate patched CustomResource
		//	crPatched := stackInfo.Outputs["crPatched"].(map[string]interface{})
		//	fooV, ok, err := unstructured.NestedString(crPatched, "metadata", "labels", "foo")
		//	assert.True(t, ok)
		//	assert.NoError(t, err)
		//	assert.Equal(t, "foo", fooV)
		//},
	})
	test, _ = testClusters.WrapProviderTestOptions(test)
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
			parentA := tests.SearchResourcesByName(stackInfo, "", "pkg:index:MyComponent", "a")

			// some helper functions for naming purposes
			providerUrn := func(prov *apitype.ResourceV3) resource.URN {
				return prov.URN + resource.URNNameDelimiter + resource.URN(prov.ID)
			}
			urn := func(parts ...string) resource.URN {
				parentType := tokens.Type(strings.Join(parts[0:len(parts)-2], resource.URNTypeDelimiter))
				baseType := tokens.Type(parts[len(parts)-2])
				name := parts[len(parts)-1]
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
			g.Expect(rr.Named(urn("pkg:index:MyComponent", "a"),
				"kubernetes:yaml:ConfigGroup", "cg-options")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           HaveExactElements(Alias("cg-options-old"), Alias("cg-options-aliased")),
						"Protect":           BeTrue(),
						"Dependencies":      HaveExactElements(string(sleep.URN)),
						"Provider":          BeEmpty(),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						// quirk: RegisterResource for component resources doesn't include provider info.
						"Providers":     BeEmpty(),
						"IgnoreChanges": HaveExactElements("ignored"),
					}),
				}),
			))
			g.Expect(rr.Named(urn("pkg:index:MyComponent", "kubernetes:yaml:ConfigGroup", "cg-options"),
				"kubernetes:core/v1:ConfigMap", "cg-options-cg-options-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases": HaveExactElements(
							Alias("cg-options-cm-1-k8s-aliased"),
							Alias(parentA.URN, tokens.Type("kubernetes:core/v1:ConfigMap"), "cg-options-cg-options-cm-1"),
							Alias("cg-options-cg-options-cm-1-aliased")),
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
			g.Expect(rr.Named(urn("pkg:index:MyComponent", "kubernetes:yaml:ConfigGroup", "cg-options"),
				"kubernetes:yaml:ConfigFile", "cg-options-testdata/options/configgroup/manifest.yaml")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases": HaveExactElements(
							Alias(parentA.URN, tokens.Type("kubernetes:yaml:ConfigFile"), "cg-options-testdata/options/configgroup/manifest.yaml"),
							Alias("cg-options-testdata/options/configgroup/manifest.yaml-aliased")),
						"Protect":           BeTrue(),
						"Dependencies":      BeEmpty(),
						"Provider":          BeEmpty(),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						"IgnoreChanges":     BeEmpty(),
					}),
				}),
			))
			g.Expect(rr.Named(urn("pkg:index:MyComponent", "kubernetes:yaml:ConfigGroup", "kubernetes:yaml:ConfigFile", "cg-options-testdata/options/configgroup/manifest.yaml"),
				"kubernetes:core/v1:ConfigMap", "cg-options-configgroup-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases": HaveExactElements(
							Alias("configgroup-cm-1-k8s-aliased"),
							Alias(parentA.URN, tokens.Type("kubernetes:core/v1:ConfigMap"), "cg-options-configgroup-cm-1"),
							Alias("cg-options-configgroup-cm-1-aliased")),
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
				"kubernetes:yaml:ConfigGroup", "cg-provider")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEmpty(),
						"Version":  BeEmpty(),
						// quirk: RegisterResource for component resources doesn't include provider info.
						"Providers": BeEmpty(),
					}),
				}),
			))
			g.Expect(rr.Named(urn("kubernetes:yaml:ConfigGroup", "cg-provider"),
				"kubernetes:core/v1:ConfigMap", "cg-provider-cg-provider-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEquivalentTo(providerUrn(providerB)),
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
						// quirk: RegisterResource for component resources doesn't include provider info.
						"Providers": BeEmpty(),
					}),
				}),
			))
			g.Expect(rr.Named(urn("kubernetes:yaml:ConfigGroup", "cg-nullopts"),
				"kubernetes:core/v1:ConfigMap", "cg-nullopts-cg-nullopts-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEquivalentTo(providerUrn(providerNullOpts)),
					}),
				}),
			))

			// --- ConfigFile ---

			// ConfigFile "cf-options" with most options
			g.Expect(rr.Named(urn("pkg:index:MyComponent", "a"),
				"kubernetes:yaml:ConfigFile", "cf-options-cf-options")).To(HaveExactElements(
				// quirk: dotnet SDK applies resource_prefix ("cf-options") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           HaveExactElements(Alias("cf-options-old"), Alias("cf-options-cf-options-aliased")),
						"Protect":           BeTrue(),
						"Dependencies":      HaveExactElements(string(sleep.URN)),
						"Provider":          BeEmpty(),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						// quirk: RegisterResource for component resources doesn't include provider info.
						"Providers":     BeEmpty(),
						"IgnoreChanges": HaveExactElements("ignored"),
					}),
				}),
			))
			g.Expect(rr.Named(urn("pkg:index:MyComponent", "kubernetes:yaml:ConfigFile", "cf-options-cf-options"),
				"kubernetes:core/v1:ConfigMap", "cf-options-configfile-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases": HaveExactElements(
							Alias("configfile-cm-1-k8s-aliased"),
							Alias(parentA.URN, tokens.Type("kubernetes:core/v1:ConfigMap"), "cf-options-configfile-cm-1"),
							Alias("cf-options-configfile-cm-1-aliased")),
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
				// quirk: dotnet SDK applies resource_prefix ("cf-provider") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEmpty(),
						"Version":  BeEmpty(),
						// quirk: RegisterResource for component resources doesn't include provider info.
						"Providers": BeEmpty(),
					}),
				}),
			))
			g.Expect(rr.Named(urn("kubernetes:yaml:ConfigFile", "cf-provider-cf-provider"),
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
				// quirk: dotnet SDK applies resource_prefix ("cf-nullopts") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEmpty(),
						"Version":  BeEmpty(),
						// quirk: RegisterResource for component resources doesn't include provider info.
						"Providers": BeEmpty(),
					}),
				}),
			))
			g.Expect(rr.Named(urn("kubernetes:yaml:ConfigFile", "cf-nullopts-cf-nullopts"),
				"kubernetes:core/v1:ConfigMap", "cf-nullopts-configfile-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEquivalentTo(providerUrn(providerNullOpts)),
					}),
				}),
			))

			// --- Directory ---

			// Directory "kustomize-options" with most options
			g.Expect(rr.Named(urn("pkg:index:MyComponent", "a"),
				"kubernetes:kustomize:Directory", "kustomize-options-kustomize-options")).To(HaveExactElements(
				// quirk: dotnet SDK applies resource_prefix ("kustomize-options") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           HaveExactElements(Alias("kustomize-options-old"), Alias("kustomize-options-kustomize-options-aliased")),
						"Protect":           BeTrue(),
						"Dependencies":      HaveExactElements(string(sleep.URN)),
						"Provider":          BeEmpty(),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						// quirk: RegisterResource for component resources doesn't include provider info.
						"Providers":     BeEmpty(),
						"IgnoreChanges": HaveExactElements("ignored"),
					}),
				}),
			))
			// urn:pulumi:p-it-pulumitron-options-a5535ee6::options-test::pkg:index:MyComponent$kubernetes:kustomize:Directory::kustomize-options-kustomize-options",
			g.Expect(rr.Named(urn("pkg:index:MyComponent", "kubernetes:kustomize:Directory", "kustomize-options-kustomize-options"),
				"kubernetes:core/v1:ConfigMap", "kustomize-options-kustomize-cm-1-2kkk4bthmg")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           HaveExactElements(Alias("kustomize-cm-1-2kkk4bthmg-k8s-aliased"), Alias("kustomize-options-kustomize-cm-1-2kkk4bthmg-aliased")),
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
				// quirk: dotnet SDK applies resource_prefix ("kustomize-provider") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEmpty(),
						"Version":  BeEmpty(),
						// quirk: RegisterResource for component resources doesn't include provider info.
						"Providers": BeEmpty(),
					}),
				}),
			))
			g.Expect(rr.Named(urn("kubernetes:kustomize:Directory", "kustomize-provider-kustomize-provider"),
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
				// quirk: dotnet SDK applies resource_prefix ("kustomize-nullopts") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEmpty(),
						"Version":  BeEmpty(),
						// quirk: RegisterResource for component resources doesn't include provider info.
						"Providers": BeEmpty(),
					}),
				}),
			))
			g.Expect(rr.Named(urn("kubernetes:kustomize:Directory", "kustomize-nullopts-kustomize-nullopts"),
				"kubernetes:core/v1:ConfigMap", "kustomize-nullopts-kustomize-cm-1-2kkk4bthmg")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider":  BeEquivalentTo(providerUrn(providerNullOpts)),
						"Version":   Not(BeEmpty()),
						"Providers": BeEmpty(),
					}),
				}),
			))

			// --- Chart ---

			// Chart "chart-options"
			g.Expect(rr.Named(urn("pkg:index:MyComponent", "a"),
				"kubernetes:helm.sh/v3:Chart", "chart-options-chart-options")).To(HaveExactElements(
				// quirk: dotnet SDK applies resource_prefix ("chart-options") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           HaveExactElements(Alias(tokens.Type("kubernetes:helm.sh/v2:Chart")), Alias("chart-options-old"), Alias("chart-options-chart-options-aliased")),
						"Protect":           BeTrue(),
						"Dependencies":      HaveExactElements(string(sleep.URN)),
						"Provider":          BeEmpty(),
						"Version":           Equal("1.2.3"),
						"PluginDownloadURL": Equal("https://a.pulumi.test"),
						// quirk: RegisterResource for component resources doesn't include provider info.
						"Providers":     BeEmpty(),
						"IgnoreChanges": HaveExactElements("ignored"),
					}),
				}),
			))
			g.Expect(rr.Named(urn("pkg:index:MyComponent", "kubernetes:helm.sh/v3:Chart", "chart-options-chart-options"),
				"kubernetes:core/v1:ConfigMap", "chart-options-chart-options-chart-options-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Aliases":           HaveExactElements(Alias("chart-options-chart-options-cm-1-k8s-aliased"), Alias("chart-options-chart-options-chart-options-cm-1-aliased")),
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
				// quirk: dotnet SDK applies resource_prefix ("chart-provider") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEmpty(),
						// quirk: dotnet SDK applies a default version
						"Version": Not(BeEmpty()),
						// quirk: RegisterResource for component resources doesn't include provider info.
						"Providers": BeEmpty(),
					}),
				}),
			))
			g.Expect(rr.Named(urn("kubernetes:helm.sh/v3:Chart", "chart-provider-chart-provider"),
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
				// quirk: dotnet SDK applies resource_prefix ("chart-options") to the component itself.
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider": BeEmpty(),
						// quirk: dotnet SDK applies a default version
						"Version": Not(BeEmpty()),
						// quirk: RegisterResource for component resources doesn't include provider info.
						"Providers": BeEmpty(),
					}),
				}),
			))
			g.Expect(rr.Named(urn("kubernetes:helm.sh/v3:Chart", "chart-nullopts-chart-nullopts"),
				"kubernetes:core/v1:ConfigMap", "chart-nullopts-chart-nullopts-chart-nullopts-cm-1")).To(HaveExactElements(
				MatchFields(IgnoreExtras, Fields{
					"Request": MatchFields(IgnoreExtras, Fields{
						"Provider":  BeEquivalentTo(providerUrn(providerNullOpts)),
						"Version":   Not(BeEmpty()),
						"Providers": BeEmpty(),
					}),
				}),
			))
		},
	})

	options, _ = testClusters.WrapProviderTestOptions(options)
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
