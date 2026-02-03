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
package v2

import (
	"context"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
	. "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/gomega"
	providerresource "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

var _ = Describe("Construct", func() {
	var tc *componentProviderTestContext
	var opts *providerresource.ResourceProviderOptions
	var req *pulumirpc.ConstructRequest
	var inputs resource.PropertyMap
	var k *ConfigGroupProvider

	BeforeEach(func() {
		tc = newTestContext(GinkgoTB())

		opts = &providerresource.ResourceProviderOptions{}
		opts.ClientSet, _, _, _ = fake.NewSimpleDynamicClient()

		// initialize the ConstructRequest to be customized in nested BeforeEach blocks
		req = tc.NewConstructRequest()
		req.Type = "kubernetes:yaml/v2:ConfigGroup"
		req.Name = "test"

		// initialize the input PropertyMap to be serialized into the request in JustBeforeEach
		inputs = make(resource.PropertyMap)
	})

	JustBeforeEach(func() {
		var err error
		k = NewConfigGroupProvider(opts).(*ConfigGroupProvider)
		req.Inputs, err = plugin.MarshalProperties(inputs, plugin.MarshalOptions{
			Label: "inputs", KeepSecrets: true, KeepResources: true, KeepUnknowns: true, KeepOutputValues: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("should register a component resource", func() {
		resp, err := pulumiprovider.Construct(context.Background(), req, tc.EngineConn(), k.Construct)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.Urn).Should(Equal("urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup::test"))

		Expect(tc.monitor.Resources()).To(MatchKeys(IgnoreExtras, Keys{
			"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup::test": MatchProps(IgnoreExtras, Props{
				"id": MatchValue("test"),
			}),
		}))
	})

	commonAssertions := func() {
		GinkgoHelper()

		It("should provide a 'resources' output property", func(ctx context.Context) {
			resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
			Expect(err).ShouldNot(HaveOccurred())
			outputs := unmarshalProperties(GinkgoTB(), resp.State)
			Expect(outputs).To(MatchProps(IgnoreExtras, Props{
				"resources": MatchArrayValue(ConsistOf(
					MatchResourceReferenceValue(
						"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup$kubernetes:core/v1:Namespace::test:my-namespace",
						"test:my-namespace",
					),
					MatchResourceReferenceValue(
						"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup$kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::test:crontabs.stable.example.com",
						"test:crontabs.stable.example.com",
					),
					MatchResourceReferenceValue(
						"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup$kubernetes:core/v1:ConfigMap::test:my-namespace/my-map",
						"test:my-namespace/my-map",
					),
					MatchResourceReferenceValue(
						"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup$kubernetes:stable.example.com/v1:CronTab::test:my-namespace/my-new-cron-object",
						"test:my-namespace/my-new-cron-object",
					),
				)),
			}))
		})

		Context("given a resource prefix", func() {
			BeforeEach(func() {
				inputs["resourcePrefix"] = resource.NewStringProperty("prefixed")
			})
			It("should use the prefix (instead of the component name)", func(ctx context.Context) {
				resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				outputs := unmarshalProperties(GinkgoTB(), resp.State)
				Expect(outputs).To(MatchProps(IgnoreExtras, Props{
					"resources": MatchArrayValue(ConsistOf(
						MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup$kubernetes:core/v1:Namespace::prefixed:my-namespace",
							"prefixed:my-namespace",
						),
						MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup$kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::prefixed:crontabs.stable.example.com",
							"prefixed:crontabs.stable.example.com",
						),
						MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup$kubernetes:core/v1:ConfigMap::prefixed:my-namespace/my-map",
							"prefixed:my-namespace/my-map",
						),
						MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup$kubernetes:stable.example.com/v1:CronTab::prefixed:my-namespace/my-new-cron-object",
							"prefixed:my-namespace/my-new-cron-object",
						),
					)),
				}))
			})
		})

		Context("given a blank resource prefix", func() {
			BeforeEach(func() {
				inputs["resourcePrefix"] = resource.NewStringProperty("")
			})
			It("should have no prefix", func(ctx context.Context) {
				resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				outputs := unmarshalProperties(GinkgoTB(), resp.State)
				Expect(outputs).To(MatchProps(IgnoreExtras, Props{
					"resources": MatchArrayValue(ConsistOf(
						MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup$kubernetes:core/v1:Namespace::my-namespace",
							"my-namespace",
						),
						MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup$kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::crontabs.stable.example.com",
							"crontabs.stable.example.com",
						),
						MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup$kubernetes:core/v1:ConfigMap::my-namespace/my-map",
							"my-namespace/my-map",
						),
						MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup$kubernetes:stable.example.com/v1:CronTab::my-namespace/my-new-cron-object",
							"my-namespace/my-new-cron-object",
						),
					)),
				}))
			})
		})
	}

	Describe("yamls", func() {
		Context("when the input is a valid YAML", func() {
			BeforeEach(func() {
				inputs["yaml"] = resource.NewStringProperty(manifest)
			})
			commonAssertions()
		})
	})

	Describe("objs", func() {
		decodeObjects := func(manifest string) []resource.PropertyValue {
			// decode the manifest to Unstructured objects, then convert to input properties
			resources, err := yamlDecode(manifest)
			Expect(err).ShouldNot(HaveOccurred())
			var objs []resource.PropertyValue
			for _, res := range resources {
				objs = append(objs, resource.NewPropertyValue(res.Object))
			}
			return objs
		}

		Context("when the input is a valid object literal", func() {
			BeforeEach(func() {
				inputs["objs"] = resource.NewArrayProperty(decodeObjects(manifest))
			})
			commonAssertions()
		})

		Context("when the object is a list", func() {
			BeforeEach(func() {
				inputs["objs"] = resource.NewArrayProperty(decodeObjects(list))
			})

			It("should expand the list", func(ctx context.Context) {
				resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				Expect(err).ShouldNot(HaveOccurred())
				outputs := unmarshalProperties(GinkgoTB(), resp.State)
				Expect(outputs).To(MatchProps(IgnoreExtras, Props{
					"resources": MatchArrayValue(HaveExactElements(
						MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup$kubernetes:core/v1:ConfigMap::test:map-1",
							"test:map-1",
						),
						MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup$kubernetes:core/v1:ConfigMap::test:map-2",
							"test:map-2",
						),
						MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigGroup$kubernetes:core/v1:ConfigMap::test:map-3",
							"test:map-3",
						),
					)),
				}))
			})
		})
	})

	Describe("files", func() {
		Context("when the input is a valid file", func() {
			BeforeEach(func() {
				tempDir := GinkgoTB().TempDir()
				err := os.WriteFile(filepath.Join(tempDir, "manifest.yaml"), []byte(manifest), 0o600)
				Expect(err).ShouldNot(HaveOccurred())
				inputs["files"] = resource.NewArrayProperty(
					[]resource.PropertyValue{resource.NewStringProperty(filepath.Join(tempDir, "manifest.yaml"))},
				)
			})
			commonAssertions()
		})
	})

	Describe("preview", func() {
		Context("when the input value(s) are unknown", func() {
			BeforeEach(func() {
				req.DryRun = true
				inputs["yaml"] = resource.MakeComputed(resource.NewStringProperty(""))
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
})
