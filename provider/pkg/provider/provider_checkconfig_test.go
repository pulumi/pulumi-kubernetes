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

	"github.com/pulumi/pulumi/sdk/v3/go/common/providers"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

var _ = Describe("RPC:CheckConfig", func() {
	var opts []NewProviderOption
	var k *kubeProvider
	var req *pulumirpc.CheckRequest
	var news resource.PropertyMap
	var config *clientcmdapi.Config

	BeforeEach(func() {
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

	JustBeforeEach(func() {
		var err error
		k = pctx.NewProvider(opts...)

		req.News, err = plugin.MarshalProperties(news, plugin.MarshalOptions{
			Label: "news", KeepUnknowns: true, SkipNulls: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("Strict Mode", func() {
		BeforeEach(func() {
			news["strictMode"] = resource.NewStringProperty("true")
			news["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToString(config))
			news["context"] = resource.NewStringProperty(config.CurrentContext)
		})

		Context("when enabled on the default provider", func() {
			BeforeEach(func() {
				req.Urn = "urn:pulumi:test::test::pulumi:providers:kubernetes::default"
				Expect(providers.IsDefaultProvider(resource.URN(req.Urn))).To(BeTrue())
			})
			It("should fail because strict mode prohibits default provider", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Failures).To(HaveExactElements(
					CheckFailure("", Equal(`strict mode prohibits default provider`))))
			})
		})

		Context("when kubeconfig is NOT specified", func() {
			BeforeEach(func() {
				delete(news, "kubeconfig")
			})
			It("should fail because strict mode requires kubeconfig", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Failures).To(HaveExactElements(
					CheckFailure("kubeconfig", Equal(`strict mode requires Provider "kubeconfig" argument`))))
			})
		})

		Context("when context is NOT specified", func() {
			BeforeEach(func() {
				delete(news, "context")
			})
			It("should fail because strict mode requires context", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Failures).To(HaveExactElements(
					CheckFailure("context", Equal(`strict mode requires Provider "context" argument`))))
			})
		})

		Context("when properly configured", func() {
			It("should succeed", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Failures).To(BeEmpty())
			})
		})
	})

	Describe("Yaml Rendering Mode", func() {
		BeforeEach(func() {
			news["renderYamlToDirectory"] = resource.NewStringProperty("true")
		})

		Context("when kubeconfig is specified", func() {
			BeforeEach(func() {
				news["kubeconfig"] = resource.NewStringProperty(WriteKubeconfigToString(config))
			})
			It("should fail because yaml mode disallows kubeconfig", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Failures).To(HaveExactElements(
					CheckFailure(
						"kubeconfig",
						Equal(`"kubeconfig" arg is not compatible with "renderYamlToDirectory" arg`),
					),
				))
			})
		})

		Context("when context is specified", func() {
			BeforeEach(func() {
				news["context"] = resource.NewStringProperty(config.CurrentContext)
			})
			It("should fail because yaml mode disallows context", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Failures).To(HaveExactElements(
					CheckFailure("context", Equal(`"context" arg is not compatible with "renderYamlToDirectory" arg`))))
			})
		})

		Context("when cluster is specified", func() {
			BeforeEach(func() {
				news["cluster"] = resource.NewStringProperty(config.Contexts[config.CurrentContext].Cluster)
			})
			It("should fail because yaml mode disallows cluster", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Failures).To(HaveExactElements(
					CheckFailure("cluster", Equal(`"cluster" arg is not compatible with "renderYamlToDirectory" arg`))))
			})
		})

		Context("when properly configured", func() {
			It("should succeed", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Failures).To(BeEmpty())
			})
		})
	})
})
