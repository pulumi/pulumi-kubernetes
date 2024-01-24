// Copyright 2016-2023, Pulumi Corporation.
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
	"os/user"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	gomegatypes "github.com/onsi/gomega/types"
	"github.com/pulumi/pulumi/pkg/v3/resource/deploy/providers"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

func MatchCheckFailure(prop string) gomegatypes.GomegaMatcher {
	return WithTransform(func(failure *pulumirpc.CheckFailure) string {
		return failure.GetProperty()
	}, Equal(prop))
}

var _ = Describe("RPC:CheckConfig", func() {
	var k *kubeProvider
	var req *pulumirpc.CheckRequest
	var news resource.PropertyMap
	var config *clientcmdapi.Config

	BeforeEach(func() {
		var err error
		k, err = pctx.NewProvider()
		Expect(err).ShouldNot(HaveOccurred())

		// load the ambient kubeconfig for test purposes
		homeDir := func() string {
			// Ignore errors. The filepath will be checked later, so we can handle failures there.
			usr, _ := user.Current()
			return usr.HomeDir
		}
		config, err = clientcmd.LoadFromFile(filepath.Join(homeDir(), "/.kube/config"))
		Expect(err).ToNot(HaveOccurred())

		// initialize the CheckRequest to be customized in nested BeforeEach blocks
		req = &pulumirpc.CheckRequest{
			Urn: "urn:pulumi:test::test::pulumi:providers:kubernetes::k8s",
		}
		// initialize the 'new' PropertyMap to be serialized into the request in JustBeforeEach
		news = make(resource.PropertyMap)
	})

	JustBeforeEach(func() {
		var err error
		Expect(err).ShouldNot(HaveOccurred())
		req.News, err = plugin.MarshalProperties(news, plugin.MarshalOptions{
			Label: "news", KeepUnknowns: true, SkipNulls: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("Strict Mode", func() {
		BeforeEach(func() {
			news["strictMode"] = resource.NewStringProperty("true")
			news["kubeconfig"] = resource.NewStringProperty(KubeconfigAsString(config))
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
				Expect(resp.Failures).To(HaveExactElements(MatchCheckFailure("")))
			})
		})

		Context("when kubeconfig is NOT specified", func() {
			BeforeEach(func() {
				delete(news, "kubeconfig")
			})
			It("should fail because strict mode requires kubeconfig", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Failures).To(HaveExactElements(MatchCheckFailure("kubeconfig")))
			})
		})

		Context("when context is NOT specified", func() {
			BeforeEach(func() {
				delete(news, "context")
			})
			It("should fail because strict mode requires context", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Failures).To(HaveExactElements(MatchCheckFailure("context")))
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
				news["kubeconfig"] = resource.NewStringProperty(KubeconfigAsString(config))
			})
			It("should fail because yaml mode disallows kubeconfig", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Failures).To(HaveExactElements(MatchCheckFailure("kubeconfig")))
			})
		})

		Context("when context is specified", func() {
			BeforeEach(func() {
				news["context"] = resource.NewStringProperty(config.CurrentContext)
			})
			It("should fail because yaml mode disallows context", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Failures).To(HaveExactElements(MatchCheckFailure("context")))
			})
		})

		Context("when cluster is specified", func() {
			BeforeEach(func() {
				news["cluster"] = resource.NewStringProperty(config.Contexts[config.CurrentContext].Cluster)
			})
			It("should fail because yaml mode disallows cluster", func() {
				resp, err := k.CheckConfig(context.Background(), req)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Failures).To(HaveExactElements(MatchCheckFailure("cluster")))
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
