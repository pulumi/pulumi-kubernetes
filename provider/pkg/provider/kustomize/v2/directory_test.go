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

	gk "github.com/onsi/ginkgo/v2"
	gm "github.com/onsi/gomega"
	gs "github.com/onsi/gomega/gstruct"
	kprovider "sigs.k8s.io/kustomize/api/provider"
	kresmap "sigs.k8s.io/kustomize/api/resmap"
	kresource "sigs.k8s.io/kustomize/api/resource"
	kfilesys "sigs.k8s.io/kustomize/kyaml/filesys"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
	pgm "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/gomega"
	providerresource "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/resource"
)

var (
	depProvider = kprovider.NewDefaultDepProvider()
	rf          = depProvider.GetResourceFactory()
)

type fakeKustomizer struct {
	resmap kresmap.ResMap
	err    error
}

var _ kustomizer = &fakeKustomizer{}

func (k *fakeKustomizer) Run(_ kfilesys.FileSystem, _ /* path */ string) (kresmap.ResMap, error) {
	return k.resmap, k.err
}

func makeCm(i int) *kresource.Resource {
	resource, err := rf.FromMap(
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name": fmt.Sprintf("cm%03d", i),
			},
		})
	if err != nil {
		panic(err)
	}
	return resource
}

var _ = gk.Describe("Construct", func() {
	var tc *componentProviderTestContext
	var opts *providerresource.ResourceProviderOptions
	var req *pulumirpc.ConstructRequest
	var inputs resource.PropertyMap
	var tool *fakeKustomizer
	var k *DirectoryProvider

	gk.BeforeEach(func() {
		tc = newTestContext(gk.GinkgoTB())

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

	gk.JustBeforeEach(func() {
		var err error
		k = &DirectoryProvider{
			opts: opts,
			makeKustomizer: func(_ *directoryArgs) kustomizer {
				return tool
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
		gm.Expect(resp.Urn).Should(gm.Equal("urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory::test"))

		gm.Expect(tc.monitor.Resources()).To(gs.MatchKeys(gs.IgnoreExtras, gs.Keys{
			"urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory::test": pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
				"id": pgm.MatchValue("test"),
			}),
		}))
	})

	gk.Describe("Preview", func() {
		gk.Context("when the input value(s) are unknown", func() {
			gk.BeforeEach(func() {
				req.DryRun = true
				inputs["directory"] = resource.MakeComputed(resource.NewStringProperty(""))
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

	gk.Describe("Templating", func() {
		gk.Describe("Namespacing", func() {
			gk.Context("by default", func() {
				gk.It("should use the context namespace", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					outputs := unmarshalProperties(gk.GinkgoTB(), resp.State)
					gm.Expect(outputs).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
						"resources": pgm.MatchArrayValue(gm.ContainElements(
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory$kubernetes:core/v1:ConfigMap::test:default/cm001",
								"test:default/cm001",
							),
						)),
					}))
				})
			})

			gk.Context("given a provider namespace", func() {
				gk.BeforeEach(func() {
					opts.DefaultNamespace = "provider"
				})
				gk.It("should use the provider's namespace", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					outputs := unmarshalProperties(gk.GinkgoTB(), resp.State)
					gm.Expect(outputs).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
						"resources": pgm.MatchArrayValue(gm.ContainElements(
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory$kubernetes:core/v1:ConfigMap::"+
									"test:provider/cm001",
								"test:provider/cm001",
							),
						)),
					}))
				})
			})

			gk.Context("given an configured namespace", func() {
				gk.BeforeEach(func() {
					inputs["namespace"] = resource.NewStringProperty("override")
				})
				gk.It("should use the configured namespace", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					outputs := unmarshalProperties(gk.GinkgoTB(), resp.State)
					gm.Expect(outputs).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
						"resources": pgm.MatchArrayValue(gm.ContainElements(
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory$kubernetes:core/v1:ConfigMap::"+
									"test:override/cm001",
								"test:override/cm001",
							),
						)),
					}))
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
						"urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory$kubernetes:core/v1:ConfigMap::test:default/cm001",
						"test:default/cm001",
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
							"urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory$kubernetes:core/v1:ConfigMap::"+
								"prefixed:default/cm001",
							"prefixed:default/cm001",
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
					"urn:pulumi:stack::project::kubernetes:kustomize/v2:Directory$kubernetes:core/v1:ConfigMap::" +
						"test:default/cm001": gs.MatchFields(
						gs.IgnoreExtras,
						gs.Fields{
							"State": pgm.HaveSkipAwaitAnnotation(),
						},
					),
				}))
			})
		})
	})
})
