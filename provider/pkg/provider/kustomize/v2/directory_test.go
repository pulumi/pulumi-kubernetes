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
	"fmt"

	. "github.com/onsi/ginkgo/v2"      //nolint:golint // dot-imports
	. "github.com/onsi/gomega"         //nolint:golint // dot-imports
	. "github.com/onsi/gomega/gstruct" //nolint:golint // dot-imports
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
	providerresource "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/resource"
	. "github.com/pulumi/pulumi-kubernetes/tests/v4/gomega"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	kprovider "sigs.k8s.io/kustomize/api/provider"
	kresmap "sigs.k8s.io/kustomize/api/resmap"
	kresource "sigs.k8s.io/kustomize/api/resource"
	kfilesys "sigs.k8s.io/kustomize/kyaml/filesys"
)

var depProvider = kprovider.NewDefaultDepProvider()
var rf = depProvider.GetResourceFactory()

type fakeKustomizer struct {
	resmap kresmap.ResMap
	err    error
}

var _ kustomizer = &fakeKustomizer{}

func (k *fakeKustomizer) Run(fSys kfilesys.FileSystem, path string) (kresmap.ResMap, error) {
	return k.resmap, k.err
}

func makeCm(i int) *kresource.Resource {
	return rf.FromMap(
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name": fmt.Sprintf("cm%03d", i),
			},
		})
}

var _ = Describe("Construct", func() {
	var tc *componentProviderTestContext
	var opts *providerresource.ResourceProviderOptions
	var req *pulumirpc.ConstructRequest
	var inputs resource.PropertyMap
	var tool *fakeKustomizer
	var k *DirectoryProvider

	BeforeEach(func() {
		tc = newTestContext(GinkgoTB())

		opts = &providerresource.ResourceProviderOptions{}
		opts.ClientSet, _, _, _ = fake.NewSimpleDynamicClient()
		opts.DefaultNamespace = "default"

		// initialize the ConstructRequest to be customized in nested BeforeEach blocks
		req = tc.NewConstructRequest()
		req.Type = "kubernetes:kustomize/v2:Directory"
		req.Name = "test"

		// initialize the input PropertyMap to be serialized into the request in JustBeforeEach
		inputs = make(resource.PropertyMap)
		inputs["directory"] = resource.NewStringProperty("reference")

		// configure the fake Kustomize tool
		tool = &fakeKustomizer{
			resmap: kresmap.New(),
		}
		_ = tool.resmap.Append(makeCm(1))
	})

	JustBeforeEach(func() {
		var err error
		k = &DirectoryProvider{
			opts: opts,
			makeKustomizer: func(args *directoryArgs) kustomizer {
				return tool
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
		Expect(resp.Urn).Should(Equal("urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory::test"))

		Expect(tc.monitor.Resources()).To(MatchKeys(IgnoreExtras, Keys{
			"urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory::test": MatchProps(IgnoreExtras, Props{
				"id": MatchValue("test"),
			}),
		}))
	})

	Describe("Preview", func() {
		Context("when the input value(s) are unknown", func() {
			BeforeEach(func() {
				req.DryRun = true
				inputs["directory"] = resource.MakeComputed(resource.NewStringProperty(""))
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

	Describe("Templating", func() {
		Describe("Namespacing", func() {
			Context("by default", func() {
				It("should use the context namespace", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					Expect(err).ShouldNot(HaveOccurred())
					outputs := unmarshalProperties(GinkgoTB(), resp.State)
					Expect(outputs).To(MatchProps(IgnoreExtras, Props{
						"resources": MatchArrayValue(ContainElements(
							MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory$kubernetes:core/v1:ConfigMap::test:default/cm001", "test:default/cm001"),
						)),
					}))
				})
			})

			Context("given a provider namespace", func() {
				BeforeEach(func() {
					opts.DefaultNamespace = "provider"
				})
				It("should use the provider's namespace", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					Expect(err).ShouldNot(HaveOccurred())
					outputs := unmarshalProperties(GinkgoTB(), resp.State)
					Expect(outputs).To(MatchProps(IgnoreExtras, Props{
						"resources": MatchArrayValue(ContainElements(
							MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory$kubernetes:core/v1:ConfigMap::test:provider/cm001", "test:provider/cm001"),
						)),
					}))
				})
			})

			Context("given an configured namespace", func() {
				BeforeEach(func() {
					inputs["namespace"] = resource.NewStringProperty("override")
				})
				It("should use the configured namespace", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					Expect(err).ShouldNot(HaveOccurred())
					outputs := unmarshalProperties(GinkgoTB(), resp.State)
					Expect(outputs).To(MatchProps(IgnoreExtras, Props{
						"resources": MatchArrayValue(ContainElements(
							MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory$kubernetes:core/v1:ConfigMap::test:override/cm001", "test:override/cm001"),
						)),
					}))
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
					MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory$kubernetes:core/v1:ConfigMap::test:default/cm001", "test:default/cm001"),
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
						MatchResourceReferenceValue("urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory$kubernetes:core/v1:ConfigMap::prefixed:default/cm001", "prefixed:default/cm001"),
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
					"urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory$kubernetes:core/v1:ConfigMap::test:default/cm001": MatchFields(IgnoreExtras, Fields{
						"State": HaveSkipAwaitAnnotation(),
					}),
				}))
			})
		})
	})
})
