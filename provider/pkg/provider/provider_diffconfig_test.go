// Copyright 2016-2024, Pulumi Corporation.
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

package provider

import (
	"context"

	_ "embed"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

var _ = Describe("RPC:DiffConfig", func() {
	var opts []NewProviderOption
	var k *kubeProvider
	var req *pulumirpc.DiffRequest
	var olds, news resource.PropertyMap
	var oldConfig, newConfig *clientcmdapi.Config

	BeforeEach(func() {
		opts = []NewProviderOption{}

		// make a pair of (fake) kubeconfigs.
		oldConfig = pctx.NewConfig()
		newConfig = pctx.NewConfig()

		// initialize the DiffRequest to be customized in nested BeforeEach blocks
		req = &pulumirpc.DiffRequest{
			Urn: "urn:pulumi:test::test::pulumi:providers:kubernetes::k8s",
		}
		// initialize the 'old'/'new' PropertyMaps to be serialized into the request in JustBeforeEach
		olds = make(resource.PropertyMap)
		news = make(resource.PropertyMap)
	})

	JustBeforeEach(func() {
		var err error
		k = pctx.NewProvider(opts...)

		req.Olds, err = plugin.MarshalProperties(olds, plugin.MarshalOptions{
			Label: "olds", KeepUnknowns: true, SkipNulls: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		req.News, err = plugin.MarshalProperties(news, plugin.MarshalOptions{
			Label: "news", KeepUnknowns: true, SkipNulls: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("Cluster Change Detection", func() {
		Describe("kubeconfig", func() {
			Context("when kubeconfig is a computed value", func() {
				BeforeEach(func() {
					olds["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToString(oldConfig))
					news["kubeconfig"] = resource.MakeComputed(resource.NewStringProperty(""))
				})

				It("should suggest replacement since a detailed diff cannot be performed", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
					Expect(resp.Diffs).To(ContainElements("kubeconfig"))
					Expect(resp.Replaces).To(ContainElements("kubeconfig"))
				})
			})

			Context("when kubeconfig is ambient", func() {
				BeforeEach(func() {
					delete(olds, "kubeconfig")
					delete(news, "kubeconfig")
				})

				It("should report no diffs", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_NONE))
					Expect(resp.Diffs).To(BeEmpty())
					Expect(resp.Replaces).To(BeEmpty())
				})
			})

			Context("when kubeconfig is changed from ambient to explicit", func() {
				BeforeEach(func() {
					delete(olds, "kubeconfig")
					news["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToString(newConfig))
				})

				It("should report a diff (no replace) on the kubeconfig property", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
					Expect(resp.Diffs).To(ContainElements("kubeconfig"))
					Expect(resp.Replaces).To(BeEmpty())
				})
			})

			Context("when the file path is changed", func() {
				BeforeEach(func() {
					olds["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToFile(oldConfig))
					news["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToFile(newConfig))
				})

				It("should report a diff (no replace) on the kubeconfig property", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
					Expect(resp.Diffs).To(ContainElements("kubeconfig"))
					Expect(resp.Replaces).To(BeEmpty())
				})
			})

			Context("when the cluster info has changed", func() {
				BeforeEach(func() {
					clusterName := newConfig.Contexts[newConfig.CurrentContext].Cluster
					newConfig.Clusters[clusterName].Server = "https://updated.test"
					olds["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToFile(oldConfig))
					news["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToFile(newConfig))
				})

				Context("without a clusterIdentifier", func() {
					It("should suggest replacement since the underlying cluster may be different", func() {
						resp, err := k.DiffConfig(context.Background(), req)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
						Expect(resp.Diffs).To(ContainElements("kubeconfig"))
						Expect(resp.Replaces).To(ContainElements("kubeconfig"))
					})
				})

				Context("with an unchanged cluster identifier", func() {
					BeforeEach(func() {
						olds["clusterIdentifier"] = resource.NewStringProperty("foo")
						news["clusterIdentifier"] = resource.NewStringProperty("foo")
					})

					It("shouldn't suggest replacement", func() {
						resp, err := k.DiffConfig(context.Background(), req)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
						Expect(resp.Diffs).To(ContainElements("kubeconfig"))
						Expect(resp.Replaces).To(BeEmpty())
					})
				})

				Context("with a different cluster identifier", func() {
					BeforeEach(func() {
						olds["clusterIdentifier"] = resource.NewStringProperty("foo")
						news["clusterIdentifier"] = resource.NewStringProperty("bar")
					})

					It("should suggest replacement", func() {
						resp, err := k.DiffConfig(context.Background(), req)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
						Expect(resp.Diffs).To(ContainElements("kubeconfig"))
						Expect(resp.Replaces).To(ContainElements("kubeconfig"))
					})
				})
			})
		})

		Describe("clusterIdentifier", func() {
			Context("when added", func() {
				BeforeEach(func() {
					news["clusterIdentifier"] = resource.NewStringProperty("foo")
				})

				It("should report a diff (no replace) on the clusterIdentifier property", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
					Expect(resp.Diffs).To(ContainElements("clusterIdentifier"))
					Expect(resp.Replaces).To(BeEmpty())
				})
			})

			Context("when removed", func() {
				BeforeEach(func() {
					olds["clusterIdentifier"] = resource.NewStringProperty("foo")
				})

				It("should report a diff (no replace) on the clusterIdentifier property", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
					Expect(resp.Diffs).To(ContainElements("clusterIdentifier"))
					Expect(resp.Replaces).To(BeEmpty())
				})
			})

			Context("when changed", func() {
				BeforeEach(func() {
					olds["clusterIdentifier"] = resource.NewStringProperty("foo")
					news["clusterIdentifier"] = resource.NewStringProperty("bar")
				})

				It("should suggest replacement", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
					Expect(resp.Diffs).To(ContainElements("clusterIdentifier"))
					Expect(resp.Replaces).To(ContainElements("clusterIdentifier"))
				})
			})
		})

		Describe("context", func() {
			BeforeEach(func() {
				kubeconfig := WriteKubeconfigToFile(oldConfig)
				olds["kubeconfig"] = resource.NewStringProperty(kubeconfig)
				news["kubeconfig"] = resource.NewStringProperty(kubeconfig)
			})

			Context("when context is a computed value", func() {
				BeforeEach(func() {
					olds["context"] = resource.NewStringProperty("context1")
					news["context"] = resource.MakeComputed(resource.NewStringProperty(""))
				})

				It("should suggest replacement since a detailed diff cannot be performed", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
					Expect(resp.Diffs).To(ContainElements("context"))
					Expect(resp.Replaces).To(ContainElements("context"))
				})
			})

			Context("when the context is changed to refer to an invalid value", func() {
				BeforeEach(func() {
					olds["context"] = resource.NewStringProperty("context1")
					news["context"] = resource.NewStringProperty("other")
				})

				It("should report a diff (no replace) on the context property", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
					Expect(resp.Diffs).To(ContainElements("context"))
					Expect(resp.Replaces).To(BeEmpty())
				})
			})

			Context("when the context is changed to refer to a different cluster", func() {
				BeforeEach(func() {
					olds["context"] = resource.NewStringProperty("context1")
					news["context"] = resource.NewStringProperty("context2")
				})

				Context("without a cluster identifier", func() {
					It("should suggest replacement since the underlying cluster may be different", func() {
						resp, err := k.DiffConfig(context.Background(), req)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
						Expect(resp.Diffs).To(ContainElements("context"))
						Expect(resp.Replaces).To(ContainElements("context"))
					})
				})

				Context("with an unchanged cluster identifier", func() {
					BeforeEach(func() {
						olds["clusterIdentifier"] = resource.NewStringProperty("foo")
						news["clusterIdentifier"] = resource.NewStringProperty("foo")
					})

					It("shouldn't suggest replacement", func() {
						resp, err := k.DiffConfig(context.Background(), req)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
						Expect(resp.Diffs).To(ContainElements("context"))
						Expect(resp.Replaces).To(BeEmpty())
					})
				})

				Context("with a different cluster identifier", func() {
					BeforeEach(func() {
						olds["clusterIdentifier"] = resource.NewStringProperty("foo")
						news["clusterIdentifier"] = resource.NewStringProperty("bar")
					})

					It("should suggest replacement", func() {
						resp, err := k.DiffConfig(context.Background(), req)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
						Expect(resp.Diffs).To(ContainElements("context"))
						Expect(resp.Replaces).To(ContainElements("context"))
					})
				})
			})
		})

		Describe("cluster", func() {
			BeforeEach(func() {
				kubeconfig := WriteKubeconfigToFile(oldConfig)
				olds["kubeconfig"] = resource.NewStringProperty(kubeconfig)
				news["kubeconfig"] = resource.NewStringProperty(kubeconfig)
			})

			Context("when cluster is a computed value", func() {
				BeforeEach(func() {
					olds["cluster"] = resource.NewStringProperty("cluster1")
					news["cluster"] = resource.MakeComputed(resource.NewStringProperty(""))
				})

				It("should suggest replacement since a detailed diff cannot be performed", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
					Expect(resp.Diffs).To(ContainElements("cluster"))
					Expect(resp.Replaces).To(ContainElements("cluster"))
				})
			})

			Context("when the cluster is changed to refer to an invalid value", func() {
				BeforeEach(func() {
					olds["cluster"] = resource.NewStringProperty("cluster1")
					news["cluster"] = resource.NewStringProperty("other")
				})

				It("should report a diff (no replace) on the cluster property", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
					Expect(resp.Diffs).To(ContainElements("cluster"))
					Expect(resp.Replaces).To(BeEmpty())
				})
			})

			Context("when the cluster is changed to refer to a different cluster", func() {
				BeforeEach(func() {
					olds["cluster"] = resource.NewStringProperty("cluster1")
					news["cluster"] = resource.NewStringProperty("cluster2")
				})

				It("should suggest replacement since the underlying cluster may be different", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
					Expect(resp.Diffs).To(ContainElements("cluster"))
					Expect(resp.Replaces).To(ContainElements("cluster"))
				})
			})
		})
	})
})
