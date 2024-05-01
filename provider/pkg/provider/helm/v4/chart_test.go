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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
	kubehelm "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/helm"
	providerresource "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/resource"
	. "github.com/pulumi/pulumi-kubernetes/tests/v4/gomega"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
)

var _ = Describe("Construct", func() {
	var tc *componentProviderTestContext
	var opts *providerresource.ResourceProviderOptions
	var req *pulumirpc.ConstructRequest
	var inputs resource.PropertyMap
	var initActionConfig kubehelm.InitActionConfigF
	var locator *kubehelm.FakeLocator
	var executor *kubehelm.FakeExecutor
	var k *ChartProvider

	BeforeEach(func() {
		tc = newTestContext(GinkgoTB())

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
		locator = kubehelm.NewFakeLocator("../../../../../tests/testdata/helm/reference", nil)
		executor = kubehelm.NewFakeExecutor()
	})

	JustBeforeEach(func() {
		var err error
		k = &ChartProvider{
			opts: opts,
			tool: func() *kubehelm.Tool {
				// make a fake tool for testing purposes
				return kubehelm.NewFakeTool(opts.HelmOptions.EnvSettings, initActionConfig, locator.LocateChart, executor.Execute)
			},
		}
		req.Inputs, err = plugin.MarshalProperties(inputs, plugin.MarshalOptions{
			Label: "inputs", KeepSecrets: true, KeepResources: true, KeepUnknowns: true, KeepOutputValues: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("should register a component resource", func() {
		resp, err := pulumiprovider.Construct(context.Background(), req, tc.EngineConn(), k.Construct)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.Urn).Should(Equal("urn:pulumi:stack::project::kubernetes:helm/v4:Chart::test"))

		Expect(tc.monitor.Resources()).To(MatchKeys(IgnoreExtras, Keys{
			"urn:pulumi:stack::project::kubernetes:helm/v4:Chart::test": MatchProps(IgnoreExtras, Props{
				"id": MatchValue("test"),
			}),
		}))
	})

	Describe("Preview", func() {
		Context("when the input value(s) are unknown", func() {
			BeforeEach(func() {
				req.DryRun = true
				inputs["chart"] = resource.MakeComputed(resource.NewStringProperty(""))
			})

			It("should emit a warning", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(tc.engine.Logs()).To(HaveLen(1))
			})

			It("should provide a 'resources' output property", func(ctx context.Context) {
				resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				outputs := unmarshalProperties(GinkgoTB(), resp.State)
				Expect(outputs).To(MatchProps(IgnoreExtras, Props{
					"resources": BeComputed(),
				}))
			})
		})
	})

	Describe("Connectivity", func() {
		It("should use server-side dry-run mode", func(ctx context.Context) {
			_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(executor.Action().DryRun).To(BeTrue())
			Expect(executor.Action().DryRunOption).To(Equal("server"))
			Expect(executor.Action().ClientOnly).To(BeFalse())
		})
	})

	Describe("Chart Resolution", func() {
		Describe("Chart", func() {
			BeforeEach(func() {
				inputs["chart"] = resource.NewStringProperty("reference")
			})
			It("should configure the name", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(locator.Name()).To(Equal("reference"))
			})
		})

		Describe("Version", func() {
			BeforeEach(func() {
				inputs["version"] = resource.NewStringProperty("1.0.0")
			})
			It("should configure the version", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(locator.Action().Version).To(Equal("1.0.0"))
			})
		})

		Describe("Devel", func() {
			BeforeEach(func() {
				inputs["devel"] = resource.NewBoolProperty(true)
			})
			It("should enable the devel flag", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(locator.Action().Devel).To(BeTrue())
			})
		})

		Describe("Verify", func() {
			BeforeEach(func() {
				inputs["verify"] = resource.NewBoolProperty(true)
			})
			It("should enable the verify flag", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(locator.Action().Verify).To(BeTrue())
			})
		})

		Describe("Keyring", func() {
			var pub *resource.Asset
			BeforeEach(func() {
				var err error
				pub, err = resource.NewPathAsset("../../../../../tests/testdata/helm/pubring.gpg")
				Expect(err).ShouldNot(HaveOccurred())

				inputs["verify"] = resource.NewBoolProperty(true)
				inputs["keyring"] = resource.NewAssetProperty(pub)
			})
			It("should configure the keyring", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(locator.Action().Keyring).To(Equal(pub.Path))
			})
		})

		Describe("DependencyUpdate", func() {
			BeforeEach(func() {
				inputs["dependencyUpdate"] = resource.NewBoolProperty(true)
			})
			It("should enable the dependencyUpdate flag", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(locator.Action().DependencyUpdate).To(BeTrue())
			})
		})
	})

	Describe("Values", func() {
		Describe("Literals", func() {
			BeforeEach(func() {
				inputs["values"] = resource.NewObjectProperty(resource.NewPropertyMapFromMap(map[string]any{
					"fullnameOverride": "overridden",
				}))
			})
			It("should configure the values", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(executor.Values()).To(HaveKeyWithValue("fullnameOverride", "overridden"))
			})
		})

		Describe("Values Files", func() {
			var valuesFile *resource.Asset
			BeforeEach(func() {
				var err error
				valuesFile, err = resource.NewTextAsset("fullnameOverride: overridden")
				Expect(err).ShouldNot(HaveOccurred())

				inputs["valueYamlFiles"] = resource.NewArrayProperty([]resource.PropertyValue{resource.NewAssetProperty(valuesFile)})
			})
			It("should configure the values", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(executor.Values()).To(HaveKeyWithValue("fullnameOverride", "overridden"))
			})
		})
	})

	Describe("Templating", func() {
		Describe("Namespacing", func() {
			Context("by default", func() {
				It("should use the context namespace", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(executor.Action().Namespace).To(Equal("default"))
					outputs := unmarshalProperties(GinkgoTB(), resp.State)
					Expect(outputs).To(MatchProps(IgnoreExtras, Props{
						"resources": MatchArrayValue(ContainElements(
							MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::test:crontabs.stable.example.com", "test:crontabs.stable.example.com"),
							MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::test:default/test-reference", "test:default/test-reference"),
							MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:Service::test:default/test-reference", "test:default/test-reference"),
						)),
					}))
				})
			})

			Context("given a provider namespace", func() {
				BeforeEach(func() {
					opts.DefaultNamespace = "provider"
					opts.HelmOptions.EnvSettings.SetNamespace("provider")
				})
				It("should use the provider's namespace", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(executor.Action().Namespace).To(Equal("provider"))
					outputs := unmarshalProperties(GinkgoTB(), resp.State)
					Expect(outputs).To(MatchProps(IgnoreExtras, Props{
						"resources": MatchArrayValue(ContainElements(
							MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::test:crontabs.stable.example.com", "test:crontabs.stable.example.com"),
							MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::test:provider/test-reference", "test:provider/test-reference"),
							MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:Service::test:provider/test-reference", "test:provider/test-reference"),
						)),
					}))
				})
			})

			Context("given a release namespace", func() {
				BeforeEach(func() {
					inputs["namespace"] = resource.NewStringProperty("release")
				})
				It("should use the release namespace", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(executor.Action().Namespace).To(Equal("release"))
					outputs := unmarshalProperties(GinkgoTB(), resp.State)
					Expect(outputs).To(MatchProps(IgnoreExtras, Props{
						"resources": MatchArrayValue(ContainElements(
							MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::test:crontabs.stable.example.com", "test:crontabs.stable.example.com"),
							MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::test:release/test-reference", "test:release/test-reference"),
							MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:Service::test:release/test-reference", "test:release/test-reference"),
						)),
					}))
				})
			})
		})

		Describe("Release Name", func() {
			It("should use the component name by default", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(executor.Action().ReleaseName).To(Equal("test"))
			})

			Context("given a release name", func() {
				BeforeEach(func() {
					inputs["name"] = resource.NewStringProperty("release")
				})
				It("should use the release name (instead of the component name)", func(ctx context.Context) {
					_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(executor.Action().ReleaseName).To(Equal("release"))
				})
			})
		})

		Describe("Skip CRDs", func() {
			Context("given skipCrds", func() {
				BeforeEach(func() {
					inputs["skipCrds"] = resource.NewBoolProperty(true)
				})
				It("should not produce CRDs from the 'crds/' directory of the chart", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					Expect(err).ShouldNot(HaveOccurred())
					outputs := unmarshalProperties(GinkgoTB(), resp.State)
					Expect(outputs).To(MatchProps(IgnoreExtras, Props{
						"resources": MatchArrayValue(ContainElements(
							MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::test:default/test-reference", "test:default/test-reference"),
							MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:Service::test:default/test-reference", "test:default/test-reference"),
						)),
					}))
					Expect(outputs).To(MatchProps(IgnoreExtras, Props{
						"resources": MatchArrayValue(Not(ContainElement(
							MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::test:crontabs.stable.example.com", "test:crontabs.stable.example.com"),
						))),
					}))
				})
			})
		})

		Describe("Post Renderer", func() {
			Context("given a postRenderer", func() {
				var tempdir string
				BeforeEach(func() {
					_, err := exec.LookPath("touch")
					if err != nil {
						Skip("touch command is not available")
					}
					tempdir, err = os.MkdirTemp("", "test")
					Expect(err).ShouldNot(HaveOccurred())
					DeferCleanup(func() {
						os.RemoveAll(tempdir)
					})
					inputs["postRenderer"] = resource.NewObjectProperty(resource.PropertyMap{
						"command": resource.NewStringProperty("touch"),
						"args":    resource.NewArrayProperty([]resource.PropertyValue{resource.NewStringProperty(filepath.Join(tempdir, "touched.txt"))}),
					})
				})
				It("should run the postrender command", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					Expect(err).ShouldNot(HaveOccurred())
					outputs := unmarshalProperties(GinkgoTB(), resp.State)
					Expect(outputs).To(MatchProps(IgnoreExtras, Props{
						"resources": MatchArrayValue(BeEmpty()),
					}))
					_, err = os.Stat(filepath.Join(tempdir, "touched.txt"))
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})
	})

	Describe("Resource Registration", func() {
		It("should provide a 'resources' output property", func(ctx context.Context) {
			resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
			Expect(err).ShouldNot(HaveOccurred())
			outputs := unmarshalProperties(GinkgoTB(), resp.State)
			Expect(outputs).To(MatchProps(IgnoreExtras, Props{
				"resources": MatchArrayValue(ConsistOf(
					MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::test:crontabs.stable.example.com", "test:crontabs.stable.example.com"),
					MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::test:default/test-reference", "test:default/test-reference"),
					MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:Service::test:default/test-reference", "test:default/test-reference"),
				)),
			}))
		})

		Describe("Resource Prefix", func() {
			BeforeEach(func() {
				inputs["resourcePrefix"] = resource.NewStringProperty("prefixed")
			})
			It("should use the prefix (instead of the component name)", func(ctx context.Context) {
				resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				outputs := unmarshalProperties(GinkgoTB(), resp.State)
				Expect(outputs).To(MatchProps(IgnoreExtras, Props{
					"resources": MatchArrayValue(ConsistOf(
						MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::prefixed:crontabs.stable.example.com", "prefixed:crontabs.stable.example.com"),
						MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::prefixed:default/test-reference", "prefixed:default/test-reference"),
						MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:Service::prefixed:default/test-reference", "prefixed:default/test-reference"),
					)),
				}))
			})
		})

		Describe("Skip Await", func() {
			BeforeEach(func() {
				inputs["skipAwait"] = resource.NewBoolProperty(true)
			})
			It("should not await the resources", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(tc.monitor.Registrations()).To(MatchKeys(IgnoreExtras, Keys{
					"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::test:crontabs.stable.example.com": MatchFields(IgnoreExtras, Fields{
						"State": HaveSkipAwaitAnnotation(),
					}),
					"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::test:default/test-reference": MatchFields(IgnoreExtras, Fields{
						"State": HaveSkipAwaitAnnotation(),
					}),
					"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:Service::test:default/test-reference": MatchFields(IgnoreExtras, Fields{
						"State": HaveSkipAwaitAnnotation(),
					}),
				}))
			})
		})

		Describe("helm.sh/resource-policy: keep", func() {
			BeforeEach(func() {
				inputs["values"] = resource.NewObjectProperty(resource.NewPropertyMapFromMap(map[string]any{
					"serviceAccount": map[string]any{
						"annotations": map[string]any{
							"helm.sh/resource-policy": "keep",
						},
					},
				}))
			})
			It("should enable the RetainOnDelete option", func(ctx context.Context) {
				_, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(tc.monitor.Registrations()).To(MatchKeys(IgnoreExtras, Keys{
					"urn:pulumi:stack::project::kubernetes:helm/v4:Chart$kubernetes:core/v1:ServiceAccount::test:default/test-reference": MatchFields(IgnoreExtras, Fields{
						"Request": MatchFields(IgnoreExtras, Fields{
							"RetainOnDelete": BeTrue(),
						}),
					}),
				}))
			})
		})
	})
})
