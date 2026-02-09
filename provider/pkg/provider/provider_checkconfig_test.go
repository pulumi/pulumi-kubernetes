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

	"github.com/pulumi/pulumi/sdk/v3/go/common/providers"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

var _ = gk.Describe("RPC:CheckConfig", func() {
	var opts []NewProviderOption
	var k *kubeProvider
	var req *pulumirpc.CheckRequest
	var news resource.PropertyMap
	var config *clientcmdapi.Config

	gk.BeforeEach(func() {
		opts = []NewProviderOption{}

		// make a (fake) kubeconfig.
		config = pctx.NewConfig()

		// initialize the CheckRequest to be customized in nested BeforeEach blocks
		req = &pulumirpc.CheckRequest{
			Urn: "urn:pulumi:test::test::pulumi:providers:kubernetes::k8s",
		}
		// initialize the 'new' PropertyMap to be serialized into the request in JustBeforeEach
		news = make(resource.PropertyMap)
	})

	gk.JustBeforeEach(func() {
		var err error
		k = pctx.NewProvider(opts...)

		req.News, err = plugin.MarshalProperties(news, plugin.MarshalOptions{
			Label: "news", KeepUnknowns: true, SkipNulls: true,
		})
		gm.Expect(err).ShouldNot(gm.HaveOccurred())
	})

	gk.Describe("Strict Mode", func() {
		gk.BeforeEach(func() {
			news["strictMode"] = resource.NewStringProperty("true")
			news["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToString(config))
			news["context"] = resource.NewStringProperty(config.CurrentContext)
		})

		gk.Context("when enabled on the default provider", func() {
			gk.BeforeEach(func() {
				req.Urn = "urn:pulumi:test::test::pulumi:providers:kubernetes::default"
				gm.Expect(providers.IsDefaultProvider(resource.URN(req.Urn))).To(gm.BeTrue())
			})
			gk.It("should fail because strict mode prohibits default provider", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				gm.Expect(err).ToNot(gm.HaveOccurred())
				gm.Expect(resp.Failures).To(gm.HaveExactElements(
					CheckFailure("", gm.Equal(`strict mode prohibits default provider`))))
			})
		})

		gk.Context("when kubeconfig is NOT specified", func() {
			gk.BeforeEach(func() {
				delete(news, "kubeconfig")
			})
			gk.It("should fail because strict mode requires kubeconfig", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				gm.Expect(err).ToNot(gm.HaveOccurred())
				gm.Expect(resp.Failures).To(gm.HaveExactElements(
					CheckFailure("kubeconfig", gm.Equal(`strict mode requires Provider "kubeconfig" argument`))))
			})
		})

		gk.Context("when context is NOT specified", func() {
			gk.BeforeEach(func() {
				delete(news, "context")
			})
			gk.It("should fail because strict mode requires context", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				gm.Expect(err).ToNot(gm.HaveOccurred())
				gm.Expect(resp.Failures).To(gm.HaveExactElements(
					CheckFailure("context", gm.Equal(`strict mode requires Provider "context" argument`))))
			})
		})

		gk.Context("when properly configured", func() {
			gk.It("should succeed", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				gm.Expect(err).ToNot(gm.HaveOccurred())
				gm.Expect(resp.Failures).To(gm.BeEmpty())
			})
		})
	})

	gk.Describe("Yaml Rendering Mode", func() {
		gk.BeforeEach(func() {
			news["renderYamlToDirectory"] = resource.NewStringProperty("true")
		})

		gk.Context("when kubeconfig is specified", func() {
			gk.BeforeEach(func() {
				news["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToString(config))
			})
			gk.It("should fail because yaml mode disallows kubeconfig", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				gm.Expect(err).ToNot(gm.HaveOccurred())
				gm.Expect(resp.Failures).To(gm.HaveExactElements(
					CheckFailure(
						"kubeconfig",
						gm.Equal(`"kubeconfig" arg is not compatible with "renderYamlToDirectory" arg`),
					),
				))
			})
		})

		gk.Context("when context is specified", func() {
			gk.BeforeEach(func() {
				news["context"] = resource.NewStringProperty(config.CurrentContext)
			})
			gk.It("should fail because yaml mode disallows context", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				gm.Expect(err).ToNot(gm.HaveOccurred())
				gm.Expect(resp.Failures).To(gm.HaveExactElements(
					CheckFailure("context", gm.Equal(`"context" arg is not compatible with "renderYamlToDirectory" arg`))))
			})
		})

		gk.Context("when cluster is specified", func() {
			gk.BeforeEach(func() {
				news["cluster"] = resource.NewStringProperty(config.Contexts[config.CurrentContext].Cluster)
			})
			gk.It("should fail because yaml mode disallows cluster", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				gm.Expect(err).ToNot(gm.HaveOccurred())
				gm.Expect(resp.Failures).To(gm.HaveExactElements(
					CheckFailure("cluster", gm.Equal(`"cluster" arg is not compatible with "renderYamlToDirectory" arg`))))
			})
		})

		gk.Context("when properly configured", func() {
			gk.It("should succeed", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				gm.Expect(err).ToNot(gm.HaveOccurred())
				gm.Expect(resp.Failures).To(gm.BeEmpty())
			})
		})
	})
})
