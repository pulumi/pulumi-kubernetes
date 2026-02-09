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

	gk "github.com/onsi/ginkgo/v2"
	gm "github.com/onsi/gomega"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

var _ = gk.Describe("RPC:DiffConfig", func() {
	var opts []NewProviderOption
	var k *kubeProvider
	var req *pulumirpc.DiffRequest
	var olds, news resource.PropertyMap
	var oldConfig, newConfig *clientcmdapi.Config

	gk.BeforeEach(func() {
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

	gk.JustBeforeEach(func() {
		var err error
		k = pctx.NewProvider(opts...)

		req.Olds, err = plugin.MarshalProperties(olds, plugin.MarshalOptions{
			Label: "olds", KeepUnknowns: true, SkipNulls: true,
		})
		gm.Expect(err).ShouldNot(gm.HaveOccurred())
		req.News, err = plugin.MarshalProperties(news, plugin.MarshalOptions{
			Label: "news", KeepUnknowns: true, SkipNulls: true,
		})
		gm.Expect(err).ShouldNot(gm.HaveOccurred())
	})

	gk.Describe("Cluster Change Detection", func() {
		gk.Describe("kubeconfig", func() {
			gk.Context("when kubeconfig is a computed value", func() {
				gk.BeforeEach(func() {
					olds["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToString(oldConfig))
					news["kubeconfig"] = resource.MakeComputed(resource.NewStringProperty(""))
				})

				gk.It("should suggest replacement since a detailed diff cannot be performed", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
					gm.Expect(resp.Diffs).To(gm.ContainElements("kubeconfig"))
					gm.Expect(resp.Replaces).To(gm.ContainElements("kubeconfig"))
				})
			})

			gk.Context("when kubeconfig is ambient", func() {
				gk.BeforeEach(func() {
					delete(olds, "kubeconfig")
					delete(news, "kubeconfig")
				})

				gk.It("should report no diffs", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_NONE))
					gm.Expect(resp.Diffs).To(gm.BeEmpty())
					gm.Expect(resp.Replaces).To(gm.BeEmpty())
				})
			})

			gk.Context("when kubeconfig is changed from ambient to explicit", func() {
				gk.BeforeEach(func() {
					delete(olds, "kubeconfig")
					news["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToString(newConfig))
				})

				gk.It("should report a diff (no replace) on the kubeconfig property", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
					gm.Expect(resp.Diffs).To(gm.ContainElements("kubeconfig"))
					gm.Expect(resp.Replaces).To(gm.BeEmpty())
				})
			})

			gk.Context("when the file path is changed", func() {
				gk.BeforeEach(func() {
					olds["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToFile(oldConfig))
					news["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToFile(newConfig))
				})

				gk.It("should report a diff (no replace) on the kubeconfig property", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
					gm.Expect(resp.Diffs).To(gm.ContainElements("kubeconfig"))
					gm.Expect(resp.Replaces).To(gm.BeEmpty())
				})
			})

			gk.Context("when the cluster info has changed", func() {
				gk.BeforeEach(func() {
					clusterName := newConfig.Contexts[newConfig.CurrentContext].Cluster
					newConfig.Clusters[clusterName].Server = "https://updated.test"
					olds["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToFile(oldConfig))
					news["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToFile(newConfig))
				})

				gk.Context("without a clusterIdentifier", func() {
					gk.It("should suggest replacement since the underlying cluster may be different", func() {
						resp, err := k.DiffConfig(context.Background(), req)
						gm.Expect(err).ShouldNot(gm.HaveOccurred())
						gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
						gm.Expect(resp.Diffs).To(gm.ContainElements("kubeconfig"))
						gm.Expect(resp.Replaces).To(gm.ContainElements("kubeconfig"))
					})
				})

				gk.Context("with an unchanged cluster identifier", func() {
					gk.BeforeEach(func() {
						olds["clusterIdentifier"] = resource.NewStringProperty("foo")
						news["clusterIdentifier"] = resource.NewStringProperty("foo")
					})

					gk.It("shouldn't suggest replacement", func() {
						resp, err := k.DiffConfig(context.Background(), req)
						gm.Expect(err).ShouldNot(gm.HaveOccurred())
						gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
						gm.Expect(resp.Diffs).To(gm.ContainElements("kubeconfig"))
						gm.Expect(resp.Replaces).To(gm.BeEmpty())
					})
				})

				gk.Context("with a different cluster identifier", func() {
					gk.BeforeEach(func() {
						olds["clusterIdentifier"] = resource.NewStringProperty("foo")
						news["clusterIdentifier"] = resource.NewStringProperty("bar")
					})

					gk.It("should suggest replacement", func() {
						resp, err := k.DiffConfig(context.Background(), req)
						gm.Expect(err).ShouldNot(gm.HaveOccurred())
						gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
						gm.Expect(resp.Diffs).To(gm.ContainElements("kubeconfig"))
						gm.Expect(resp.Replaces).To(gm.ContainElements("kubeconfig"))
					})
				})
			})
		})

		gk.Describe("clusterIdentifier", func() {
			gk.Context("when added", func() {
				gk.BeforeEach(func() {
					news["clusterIdentifier"] = resource.NewStringProperty("foo")
				})

				gk.It("should report a diff (no replace) on the clusterIdentifier property", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
					gm.Expect(resp.Diffs).To(gm.ContainElements("clusterIdentifier"))
					gm.Expect(resp.Replaces).To(gm.BeEmpty())
				})
			})

			gk.Context("when removed", func() {
				gk.BeforeEach(func() {
					olds["clusterIdentifier"] = resource.NewStringProperty("foo")
				})

				gk.It("should report a diff (no replace) on the clusterIdentifier property", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
					gm.Expect(resp.Diffs).To(gm.ContainElements("clusterIdentifier"))
					gm.Expect(resp.Replaces).To(gm.BeEmpty())
				})
			})

			gk.Context("when changed", func() {
				gk.BeforeEach(func() {
					olds["clusterIdentifier"] = resource.NewStringProperty("foo")
					news["clusterIdentifier"] = resource.NewStringProperty("bar")
				})

				gk.It("should suggest replacement", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
					gm.Expect(resp.Diffs).To(gm.ContainElements("clusterIdentifier"))
					gm.Expect(resp.Replaces).To(gm.ContainElements("clusterIdentifier"))
				})
			})
		})

		gk.Describe("context", func() {
			gk.BeforeEach(func() {
				kubeconfig := WriteKubeconfigToFile(oldConfig)
				olds["kubeconfig"] = resource.NewStringProperty(kubeconfig)
				news["kubeconfig"] = resource.NewStringProperty(kubeconfig)
			})

			gk.Context("when context is a computed value", func() {
				gk.BeforeEach(func() {
					olds["context"] = resource.NewStringProperty("context1")
					news["context"] = resource.MakeComputed(resource.NewStringProperty(""))
				})

				gk.It("should suggest replacement since a detailed diff cannot be performed", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
					gm.Expect(resp.Diffs).To(gm.ContainElements("context"))
					gm.Expect(resp.Replaces).To(gm.ContainElements("context"))
				})
			})

			gk.Context("when the context is changed to refer to an invalid value", func() {
				gk.BeforeEach(func() {
					olds["context"] = resource.NewStringProperty("context1")
					news["context"] = resource.NewStringProperty("other")
				})

				gk.It("should report a diff (no replace) on the context property", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
					gm.Expect(resp.Diffs).To(gm.ContainElements("context"))
					gm.Expect(resp.Replaces).To(gm.BeEmpty())
				})
			})

			gk.Context("when the context is changed to refer to a different cluster", func() {
				gk.BeforeEach(func() {
					olds["context"] = resource.NewStringProperty("context1")
					news["context"] = resource.NewStringProperty("context2")
				})

				gk.Context("without a cluster identifier", func() {
					gk.It("should suggest replacement since the underlying cluster may be different", func() {
						resp, err := k.DiffConfig(context.Background(), req)
						gm.Expect(err).ShouldNot(gm.HaveOccurred())
						gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
						gm.Expect(resp.Diffs).To(gm.ContainElements("context"))
						gm.Expect(resp.Replaces).To(gm.ContainElements("context"))
					})
				})

				gk.Context("with an unchanged cluster identifier", func() {
					gk.BeforeEach(func() {
						olds["clusterIdentifier"] = resource.NewStringProperty("foo")
						news["clusterIdentifier"] = resource.NewStringProperty("foo")
					})

					gk.It("shouldn't suggest replacement", func() {
						resp, err := k.DiffConfig(context.Background(), req)
						gm.Expect(err).ShouldNot(gm.HaveOccurred())
						gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
						gm.Expect(resp.Diffs).To(gm.ContainElements("context"))
						gm.Expect(resp.Replaces).To(gm.BeEmpty())
					})
				})

				gk.Context("with a different cluster identifier", func() {
					gk.BeforeEach(func() {
						olds["clusterIdentifier"] = resource.NewStringProperty("foo")
						news["clusterIdentifier"] = resource.NewStringProperty("bar")
					})

					gk.It("should suggest replacement", func() {
						resp, err := k.DiffConfig(context.Background(), req)
						gm.Expect(err).ShouldNot(gm.HaveOccurred())
						gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
						gm.Expect(resp.Diffs).To(gm.ContainElements("context"))
						gm.Expect(resp.Replaces).To(gm.ContainElements("context"))
					})
				})
			})
		})

		gk.Describe("cluster", func() {
			gk.BeforeEach(func() {
				kubeconfig := WriteKubeconfigToFile(oldConfig)
				olds["kubeconfig"] = resource.NewStringProperty(kubeconfig)
				news["kubeconfig"] = resource.NewStringProperty(kubeconfig)
			})

			gk.Context("when cluster is a computed value", func() {
				gk.BeforeEach(func() {
					olds["cluster"] = resource.NewStringProperty("cluster1")
					news["cluster"] = resource.MakeComputed(resource.NewStringProperty(""))
				})

				gk.It("should suggest replacement since a detailed diff cannot be performed", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
					gm.Expect(resp.Diffs).To(gm.ContainElements("cluster"))
					gm.Expect(resp.Replaces).To(gm.ContainElements("cluster"))
				})
			})

			gk.Context("when the cluster is changed to refer to an invalid value", func() {
				gk.BeforeEach(func() {
					olds["cluster"] = resource.NewStringProperty("cluster1")
					news["cluster"] = resource.NewStringProperty("other")
				})

				gk.It("should report a diff (no replace) on the cluster property", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
					gm.Expect(resp.Diffs).To(gm.ContainElements("cluster"))
					gm.Expect(resp.Replaces).To(gm.BeEmpty())
				})
			})

			gk.Context("when the cluster is changed to refer to a different cluster", func() {
				gk.BeforeEach(func() {
					olds["cluster"] = resource.NewStringProperty("cluster1")
					news["cluster"] = resource.NewStringProperty("cluster2")
				})

				gk.It("should suggest replacement since the underlying cluster may be different", func() {
					resp, err := k.DiffConfig(context.Background(), req)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(resp.Changes).To(gm.Equal(pulumirpc.DiffResponse_DIFF_SOME))
					gm.Expect(resp.Diffs).To(gm.ContainElements("cluster"))
					gm.Expect(resp.Replaces).To(gm.ContainElements("cluster"))
				})
			})
		})
	})
})
