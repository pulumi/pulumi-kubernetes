// Copyright 2016-2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package v4

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	gk "github.com/onsi/ginkgo/v2"
	gm "github.com/onsi/gomega"
	gs "github.com/onsi/gomega/gstruct"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
	pgm "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/gomega"
	kubehelm "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/helm"
	providerresource "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/resource"
)

var _ = gk.Describe("Construct", func() {
	var tc *componentProviderTestContext
	var opts *providerresource.ResourceProviderOptions
	var req *pulumirpc.ConstructRequest
	var inputs resource.PropertyMap
	var initActionConfig kubehelm.InitActionConfigF
	var locator *kubehelm.FakeLocator
	var executor *kubehelm.FakeExecutor
	var k *ChartProvider

	gk.BeforeEach(func() {
		tc = newTestContext(gk.GinkgoTB())

		opts = &providerresource.ResourceProviderOptions{}
		opts.ClientSet, _, _, _ = fake.NewSimpleDynamicClient()
		opts.DefaultNamespace = "default"
		opts.HelmOptions = &providerresource.HelmOptions{
			SuppressHelmHookWarnings: false,
			EnvSettings:              cli.New(),
		}
		opts.HelmOptions.EnvSettings.SetNamespace("default")

		// initialize the ConstructRequest to be customized in nested BeforeEach blocks
		req = tc.NewConstructRequest()
		req.Type = "kubernetes:helm/v4:Chart"
		req.Name = "test"

		// initialize the input PropertyMap to be serialized into the request in JustBeforeEach
		inputs = make(resource.PropertyMap)
		inputs["chart"] = resource.NewStringProperty("reference")

		// configure the fake Helm tool
		initActionConfig = kubehelm.FakeInitActionConfig("default", chartutil.DefaultCapabilities)
		locator = kubehelm.NewFakeLocator("./testdata/reference", nil)
		executor = kubehelm.NewFakeExecutor()
	})

	gk.JustBeforeEach(func() {
		var err error
		k = &ChartProvider{
			opts: opts,
			tool: func() *kubehelm.Tool {
				// make a fake tool for testing purposes
				return kubehelm.NewFakeTool(
					opts.HelmOptions.EnvSettings,
					initActionConfig,
					locator.LocateChart,
					executor.Execute,
				)
			},
		}
		req.Inputs, err = plugin.MarshalProperties(inputs, plugin.MarshalOptions{
			Label: "inputs", KeepSecrets: true, KeepResources: true, KeepUnknowns: true, KeepOutputValues: true,
		})
		gm.Expect(err).ShouldNot(gm.HaveOccurred())
	})

	gk.It("should register a component resource", func() {
		resp, err := pulumiprovider.Construct(context.Background(), req, tc.EngineConn(), k.Construct)
		gm.Expect(err).ShouldNot(gm.HaveOccurred())
		gm.Expect(resp.Urn).Should(gm.Equal("urn:pulumi:stack::project::kubernetes:helm/v4:Chart::test"))

		gm.Expect(tc.monitor.Resources()).To(gs.MatchKeys(gs.IgnoreExtras, gs.Keys{
			"urn:pulumi:stack::project::kubernetes:helm/v4:Chart::test": pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
				"id": pgm.MatchValue("test"),
			}),
		}))
	})

	gk.Describe("Preview", func() {
		gk.Context("when the input value(s) are unknown", func() {
			gk.BeforeEach(func() {
				req.DryRun = true
				inputs["chart"] = resource.MakeComputed(resource.NewStringProperty(""))
			})

			gk.It("should emit a warning", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(tc.engine.Logs()).To(gm.HaveLen(1))
			})

			gk.It("should provide a 'resources' output property", func(ctx context.Context) {
				resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				outputs := unmarshalProperties(gk.GinkgoTB(), resp.State)
				gm.Expect(outputs).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
					"resources": pgm.BeComputed(),
				}))
			})
		})
	})

	gk.Describe("Connectivity", func() {
		gk.It("should use server-side dry-run mode", func(ctx context.Context) {
			_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())
			gm.Expect(executor.Action().DryRun).To(gm.BeTrue())
			gm.Expect(executor.Action().DryRunOption).To(gm.Equal("server"))
			gm.Expect(executor.Action().ClientOnly).To(gm.BeTrue())
		})
		gk.Describe("Capabilities", func() {
			gk.BeforeEach(func() {
				inputs["values"] = resource.NewObjectProperty(resource.NewPropertyMapFromMap(map[string]any{
					"versionCheck": ">=1.21-0",
				}))
			})
			gk.It("should have the correct kubeversion", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				v := fake.DefaultServerVersion
				gm.Expect(executor.Action().KubeVersion).To(gs.PointTo(gm.Equal(
					chartutil.KubeVersion{Version: v.GitVersion, Major: v.Major, Minor: v.Minor})))
			})
			gk.It("should have the correct apiversions", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(executor.Action().APIVersions).To(gm.Not(gm.BeEmpty()))
			})
		})
	})

	gk.Describe("Chart Resolution", func() {
		gk.Describe("Chart", func() {
			gk.BeforeEach(func() {
				inputs["chart"] = resource.NewStringProperty("reference")
			})
			gk.It("should configure the name", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(locator.Name()).To(gm.Equal("reference"))
			})
		})

		gk.Describe("Version", func() {
			gk.BeforeEach(func() {
				inputs["version"] = resource.NewStringProperty("1.0.0")
			})
			gk.It("should configure the version", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(locator.Action().Version).To(gm.Equal("1.0.0"))
			})
		})

		gk.Describe("Devel", func() {
			gk.BeforeEach(func() {
				inputs["devel"] = resource.NewBoolProperty(true)
			})
			gk.It("should enable the devel flag", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(locator.Action().Devel).To(gm.BeTrue())
			})
		})

		gk.Describe("Verify", func() {
			gk.BeforeEach(func() {
				inputs["verify"] = resource.NewBoolProperty(true)
			})
			gk.It("should enable the verify flag", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(locator.Action().Verify).To(gm.BeTrue())
			})
		})

		gk.Describe("Keyring", func() {
			var pub *resource.Asset
			gk.BeforeEach(func() {
				var err error
				pub, err = resource.NewPathAsset("./testdata/pubring.gpg")
				gm.Expect(err).ShouldNot(gm.HaveOccurred())

				inputs["verify"] = resource.NewBoolProperty(true)
				inputs["keyring"] = resource.NewAssetProperty(pub)
			})
			gk.It("should configure the keyring", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(locator.Action().Keyring).To(gm.Equal(pub.Path))
			})
		})

		gk.Describe("DependencyUpdate", func() {
			gk.BeforeEach(func() {
				inputs["dependencyUpdate"] = resource.NewBoolProperty(true)
			})
			gk.It("should enable the dependencyUpdate flag", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(locator.Action().DependencyUpdate).To(gm.BeTrue())
			})
		})
	})

	gk.Describe("Values", func() {
		gk.Describe("Literals", func() {
			gk.BeforeEach(func() {
				inputs["values"] = resource.NewObjectProperty(resource.NewPropertyMapFromMap(map[string]any{
					"fullnameOverride": "overridden",
				}))
			})
			gk.It("should configure the values", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(executor.Values()).To(gm.HaveKeyWithValue("fullnameOverride", "overridden"))
			})
		})

		gk.Describe("Values Files", func() {
			var valuesFile *resource.Asset
			gk.BeforeEach(func() {
				var err error
				valuesFile, err = resource.NewTextAsset("fullnameOverride: overridden")
				gm.Expect(err).ShouldNot(gm.HaveOccurred())

				inputs["valueYamlFiles"] = resource.NewArrayProperty(
					[]resource.PropertyValue{resource.NewAssetProperty(valuesFile)},
				)
			})
			gk.It("should configure the values", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(executor.Values()).To(gm.HaveKeyWithValue("fullnameOverride", "overridden"))
			})
		})
	})

	gk.Describe("Templating", func() {
		gk.Describe("Namespacing", func() {
			gk.Context("by default", func() {
				gk.It("should use the context namespace", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(executor.Action().Namespace).To(gm.Equal("default"))
					outputs := unmarshalProperties(gk.GinkgoTB(), resp.State)
					gm.Expect(outputs).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
						"resources": pgm.MatchArrayValue(gm.ContainElements(
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$"+
									"kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::test:crontabs.stable.example.com",
								"test:crontabs.stable.example.com",
							),
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::"+
									"test:default/test-reference",
								"test:default/test-reference",
							),
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:Service::test:default/test-reference",
								"test:default/test-reference",
							),
						)),
					}))
				})
			})

			gk.Context("given a provider namespace", func() {
				gk.BeforeEach(func() {
					opts.DefaultNamespace = "provider"
					opts.HelmOptions.EnvSettings.SetNamespace("provider")
				})
				gk.It("should use the provider's namespace", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(executor.Action().Namespace).To(gm.Equal("provider"))
					outputs := unmarshalProperties(gk.GinkgoTB(), resp.State)
					gm.Expect(outputs).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
						"resources": pgm.MatchArrayValue(gm.ContainElements(
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$"+
									"kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::test:crontabs.stable.example.com",
								"test:crontabs.stable.example.com",
							),
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::"+
									"test:provider/test-reference",
								"test:provider/test-reference",
							),
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:Service::test:provider/test-reference",
								"test:provider/test-reference",
							),
						)),
					}))
				})
			})

			gk.Context("given a release namespace", func() {
				gk.BeforeEach(func() {
					inputs["namespace"] = resource.NewStringProperty("release")
				})
				gk.It("should use the release namespace", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(executor.Action().Namespace).To(gm.Equal("release"))
					outputs := unmarshalProperties(gk.GinkgoTB(), resp.State)
					gm.Expect(outputs).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
						"resources": pgm.MatchArrayValue(gm.ContainElements(
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$"+
									"kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::test:crontabs.stable.example.com",
								"test:crontabs.stable.example.com",
							),
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::"+
									"test:release/test-reference",
								"test:release/test-reference",
							),
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:Service::test:release/test-reference",
								"test:release/test-reference",
							),
						)),
					}))
				})
			})
		})

		gk.Describe("Release Name", func() {
			gk.It("should use the component name by default", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(executor.Action().ReleaseName).To(gm.Equal("test"))
			})

			gk.Context("given a release name", func() {
				gk.BeforeEach(func() {
					inputs["name"] = resource.NewStringProperty("release")
				})
				gk.It("should use the release name (instead of the component name)", func(ctx context.Context) {
					_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(executor.Action().ReleaseName).To(gm.Equal("release"))
				})
			})
		})

		gk.Describe("Skip CRDs", func() {
			gk.Context("given skipCrds", func() {
				gk.BeforeEach(func() {
					inputs["skipCrds"] = resource.NewBoolProperty(true)
				})
				gk.It("should not produce CRDs from the 'crds/' directory of the chart", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					outputs := unmarshalProperties(gk.GinkgoTB(), resp.State)
					gm.Expect(outputs).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
						"resources": pgm.MatchArrayValue(gm.ContainElements(
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::"+
									"test:default/test-reference",
								"test:default/test-reference",
							),
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:Service::test:default/test-reference",
								"test:default/test-reference",
							),
						)),
					}))
					gm.Expect(outputs).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
						"resources": pgm.MatchArrayValue(gm.Not(gm.ContainElement(
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$"+
									"kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::test:crontabs.stable.example.com",
								"test:crontabs.stable.example.com",
							),
						))),
					}))
				})
			})
		})

		gk.Describe("Post Renderer", func() {
			gk.Context("given a postRenderer", func() {
				var tempdir string
				gk.BeforeEach(func() {
					_, err := exec.LookPath("touch")
					if err != nil {
						gk.Skip("touch command is not available")
					}
					tempdir, err = os.MkdirTemp("", "test")
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gk.DeferCleanup(func() {
						os.RemoveAll(tempdir)
					})
					inputs["postRenderer"] = resource.NewObjectProperty(resource.PropertyMap{
						"command": resource.NewStringProperty("touch"),
						"args": resource.NewArrayProperty(
							[]resource.PropertyValue{resource.NewStringProperty(filepath.Join(tempdir, "touched.txt"))},
						),
					})
				})
				gk.It("should run the postrender command", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					outputs := unmarshalProperties(gk.GinkgoTB(), resp.State)
					gm.Expect(outputs).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
						"resources": pgm.MatchArrayValue(gm.BeEmpty()),
					}))
					_, err = os.Stat(filepath.Join(tempdir, "touched.txt"))
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
				})
			})
		})
	})

	gk.Describe("Resource Registration", func() {
		gk.It("should provide a 'resources' output property", func(ctx context.Context) {
			resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())
			outputs := unmarshalProperties(gk.GinkgoTB(), resp.State)
			gm.Expect(outputs).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
				"resources": pgm.MatchArrayValue(gm.ConsistOf(
					pgm.MatchResourceReferenceValue(
						"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$"+
							"kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::test:crontabs.stable.example.com",
						"test:crontabs.stable.example.com",
					),
					pgm.MatchResourceReferenceValue(
						"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::"+
							"test:default/test-reference",
						"test:default/test-reference",
					),
					pgm.MatchResourceReferenceValue(
						"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:Service::test:default/test-reference",
						"test:default/test-reference",
					),
				)),
			}))
		})

		gk.Describe("Resource Prefix", func() {
			gk.BeforeEach(func() {
				inputs["resourcePrefix"] = resource.NewStringProperty("prefixed")
			})
			gk.It("should use the prefix (instead of the component name)", func(ctx context.Context) {
				resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				outputs := unmarshalProperties(gk.GinkgoTB(), resp.State)
				gm.Expect(outputs).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
					"resources": pgm.MatchArrayValue(gm.ConsistOf(
						pgm.MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$"+
								"kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::prefixed:crontabs.stable.example.com",
							"prefixed:crontabs.stable.example.com",
						),
						pgm.MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::"+
								"prefixed:default/test-reference",
							"prefixed:default/test-reference",
						),
						pgm.MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:Service::"+
								"prefixed:default/test-reference",
							"prefixed:default/test-reference",
						),
					)),
				}))
			})
		})

		gk.Describe("Skip Await", func() {
			gk.BeforeEach(func() {
				inputs["skipAwait"] = resource.NewBoolProperty(true)
			})
			gk.It("should not await the resources", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(tc.monitor.Registrations()).To(gs.MatchKeys(gs.IgnoreExtras, gs.Keys{
					"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$" +
						"kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::test:crontabs.stable.example.com": gs.MatchFields(
						gs.IgnoreExtras,
						gs.Fields{
							"State": pgm.HaveSkipAwaitAnnotation(),
						},
					),
					"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::" +
						"test:default/test-reference": gs.MatchFields(
						gs.IgnoreExtras,
						gs.Fields{
							"State": pgm.HaveSkipAwaitAnnotation(),
						},
					),
					"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:Service::" +
						"test:default/test-reference": gs.MatchFields(
						gs.IgnoreExtras,
						gs.Fields{
							"State": pgm.HaveSkipAwaitAnnotation(),
						},
					),
				}))
			})
		})

		gk.Describe("Plain HTTP set", func() {
			gk.BeforeEach(func() {
				inputs["plainHttp"] = resource.NewBoolProperty(true)
			})
			gk.It("should use plain HTTP", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(executor.Action().PlainHTTP).To(gm.Equal(true))
			})
		})

		gk.Describe("Plain HTTP unset", func() {
			gk.It("should not use plain HTTP by default", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(executor.Action().PlainHTTP).To(gm.Equal(false))
			})
		})

		gk.Describe("helm.sh/resource-policy: keep", func() {
			gk.BeforeEach(func() {
				inputs["values"] = resource.NewObjectProperty(resource.NewPropertyMapFromMap(map[string]any{
					"serviceAccount": map[string]any{
						"annotations": map[string]any{
							"helm.sh/resource-policy": "keep",
						},
					},
				}))
			})
			gk.It("should enable the RetainOnDelete option", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(tc.monitor.Registrations()).To(gs.MatchKeys(gs.IgnoreExtras, gs.Keys{
					"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::" +
						"test:default/test-reference": gs.MatchFields(
						gs.IgnoreExtras,
						gs.Fields{
							"Request": gs.MatchFields(gs.IgnoreExtras, gs.Fields{
								"RetainOnDelete": gs.PointTo(gm.BeTrue()),
							}),
						},
					),
				}))
			})
		})
	})
})
