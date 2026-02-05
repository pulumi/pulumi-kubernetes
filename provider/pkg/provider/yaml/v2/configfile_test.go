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

	gk "github.com/onsi/ginkgo/v2"
	gm "github.com/onsi/gomega"
	gs "github.com/onsi/gomega/gstruct"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumiprovider "github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
	pgm "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/gomega"
	providerresource "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/resource"
)

var _ = gk.Describe("ConfigFile.Construct", func() {
	var tc *componentProviderTestContext
	var opts *providerresource.ResourceProviderOptions
	var req *pulumirpc.ConstructRequest
	var inputs resource.PropertyMap
	var k *ConfigFileProvider

	gk.BeforeEach(func() {
		tc = newTestContext(gk.GinkgoTB())

		opts = &providerresource.ResourceProviderOptions{}
		opts.ClientSet, _, _, _ = fake.NewSimpleDynamicClient()

		// initialize the ConstructRequest to be customized in nested BeforeEach blocks
		req = tc.NewConstructRequest()
		req.Type = "kubernetes:yaml/v2:ConfigFile"
		req.Name = "test"

		// initialize the input PropertyMap to be serialized into the request in JustBeforeEach
		inputs = make(resource.PropertyMap)
	})

	gk.JustBeforeEach(func() {
		var err error
		k = NewConfigFileProvider(opts).(*ConfigFileProvider)
		req.Inputs, err = plugin.MarshalProperties(inputs, plugin.MarshalOptions{
			Label: "inputs", KeepSecrets: true, KeepResources: true, KeepUnknowns: true, KeepOutputValues: true,
		})
		gm.Expect(err).ShouldNot(gm.HaveOccurred())
	})

	componentAssertions := func() {
		gk.It("should register a component resource", func() {
			resp, err := pulumiprovider.Construct(context.Background(), req, tc.EngineConn(), k.Construct)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())
			gm.Expect(resp.Urn).Should(gm.Equal("urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigFile::test"))

			gm.Expect(tc.monitor.Resources()).To(gs.MatchKeys(gs.IgnoreExtras, gs.Keys{
				"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigFile::test": pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
					"id": pgm.MatchValue("test"),
				}),
			}))
		})
	}

	gk.Describe("file", func() {
		gk.Context("when the input is a valid file", func() {
			gk.BeforeEach(func() {
				tempDir := gk.GinkgoTB().TempDir()
				err := os.WriteFile(filepath.Join(tempDir, "manifest.yaml"), []byte(manifest), 0o600)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				inputs["file"] = resource.NewStringProperty(filepath.Join(tempDir, "manifest.yaml"))
			})

			componentAssertions()

			gk.It("should provide a 'resources' output property", func(ctx context.Context) {
				resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				outputs := unmarshalProperties(gk.GinkgoTB(), resp.State)
				gm.Expect(outputs).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
					"resources": pgm.MatchArrayValue(gm.ConsistOf(
						pgm.MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigFile$kubernetes:core/v1:Namespace::test:my-namespace",
							"test:my-namespace",
						),
						pgm.MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigFile"+
								"$kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::test:crontabs.stable.example.com",
							"test:crontabs.stable.example.com",
						),
						pgm.MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigFile"+
								"$kubernetes:core/v1:ConfigMap::test:my-namespace/my-map",
							"test:my-namespace/my-map",
						),
						pgm.MatchResourceReferenceValue(
							"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigFile"+
								"$kubernetes:stable.example.com/v1:CronTab::test:my-namespace/my-new-cron-object",
							"test:my-namespace/my-new-cron-object",
						),
					)),
				}))
			})

			gk.Context("given a resource prefix", func() {
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
								"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigFile$kubernetes:core/v1:Namespace::prefixed:my-namespace",
								"prefixed:my-namespace",
							),
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigFile"+
									"$kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::prefixed:crontabs.stable.example.com",
								"prefixed:crontabs.stable.example.com",
							),
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigFile"+
									"$kubernetes:core/v1:ConfigMap::prefixed:my-namespace/my-map",
								"prefixed:my-namespace/my-map",
							),
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigFile"+
									"$kubernetes:stable.example.com/v1:CronTab::prefixed:my-namespace/my-new-cron-object",
								"prefixed:my-namespace/my-new-cron-object",
							),
						)),
					}))
				})
			})

			gk.Context("given a blank resource prefix", func() {
				gk.BeforeEach(func() {
					inputs["resourcePrefix"] = resource.NewStringProperty("")
				})
				gk.It("should have no prefix", func(ctx context.Context) {
					resp, err := pulumiprovider.Construct(ctx, req, tc.EngineConn(), k.Construct)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					outputs := unmarshalProperties(gk.GinkgoTB(), resp.State)
					gm.Expect(outputs).To(pgm.MatchProps(gs.IgnoreExtras, pgm.Props{
						"resources": pgm.MatchArrayValue(gm.ConsistOf(
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigFile$kubernetes:core/v1:Namespace::my-namespace",
								"my-namespace",
							),
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigFile"+
									"$kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition::crontabs.stable.example.com",
								"crontabs.stable.example.com",
							),
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigFile"+
									"$kubernetes:core/v1:ConfigMap::my-namespace/my-map",
								"my-namespace/my-map",
							),
							pgm.MatchResourceReferenceValue(
								"urn:pulumi:stack::project::kubernetes:yaml/v2:ConfigFile"+
									"$kubernetes:stable.example.com/v1:CronTab::my-namespace/my-new-cron-object",
								"my-namespace/my-new-cron-object",
							),
						)),
					}))
				})
			})
		})
	})

	gk.Describe("preview", func() {
		gk.Context("when the input value(s) are unknown", func() {
			gk.BeforeEach(func() {
				req.DryRun = true
				inputs["file"] = resource.MakeComputed(resource.NewStringProperty(""))
			})

			componentAssertions()

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
})
