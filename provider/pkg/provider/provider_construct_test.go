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
	"github.com/pulumi/pulumi-kubernetes/v4/provider/pkg/clients"
	providerresource "github.com/pulumi/pulumi-kubernetes/v4/provider/pkg/provider/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	helmcli "helm.sh/helm/v3/pkg/cli"
)

var _ = Describe("RPC:Construct", func() {
	var opts []NewProviderOption
	var k *kubeProvider
	var req *pulumirpc.ConstructRequest

	BeforeEach(func() {
		opts = []NewProviderOption{}

		// initialize the ConstructRequest to be customized in nested BeforeEach blocks
		req = &pulumirpc.ConstructRequest{
			Project: "test",
			Stack:   "test",
		}
	})

	// factory is a helper function for using mock resource providers
	factory := func(rp *mockResourceProvider) providerresource.ResourceProviderFactory {
		return func(opts *providerresource.ResourceProviderOptions) providerresource.ResourceProvider {
			rp.opts = opts
			return rp
		}
	}

	JustBeforeEach(func() {
		k = pctx.NewProvider(opts...)
		k.clientSet = &clients.DynamicClientSet{}
		k.defaultNamespace = "default"
		k.helmDriver = "memory"
		k.helmSettings = helmcli.New()
	})

	Context("when the requested type is unknown", func() {
		BeforeEach(func() {
			req.Type = "kubernetes:test:UnknownComponent"
			req.Name = "testComponent"
		})
		It("should return an error", func() {
			_, err := k.Construct(context.Background(), req)
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("when the requested type is known", func() {
		var testComponent *mockResourceProvider
		BeforeEach(func() {
			testComponent = &mockResourceProvider{
				Result: &provider.ConstructResult{
					URN: pulumi.URN("urn:pulumi:test::test::test:TestComponent::testComponent"),
				},
			}
			opts = append(opts, WithResourceProvider("kubernetes:test:TestComponent", factory(testComponent)))

			req.Type = "kubernetes:test:TestComponent"
			req.Name = "testComponent"
		})

		It("should delegate to the provider", func() {
			result, err := k.Construct(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(testComponent.typ).Should(Equal("kubernetes:test:TestComponent"))
			Expect(testComponent.name).Should(Equal("testComponent"))
			Expect(result.Urn).Should(Equal("urn:pulumi:test::test::test:TestComponent::testComponent"))
		})

		It("should provide options", func() {
			_, err := k.Construct(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(testComponent.opts.ClientSet).ShouldNot(BeNil())
			Expect(testComponent.opts.DefaultNamespace).ShouldNot(BeEmpty())
			Expect(testComponent.opts.HelmOptions).ShouldNot(BeNil())
		})

		Context("when clusterUnreachable is true", func() {
			JustBeforeEach(func() {
				k.clusterUnreachable = true
				k.clusterUnreachableReason = "testing"
			})
			It("should return an error", func() {
				_, err := k.Construct(context.Background(), req)
				Expect(err).To(MatchError(ContainSubstring("configured Kubernetes cluster is unreachable")))
			})
		})
	})
})

type mockResourceProvider struct {
	Result *provider.ConstructResult
	Err    error

	opts *providerresource.ResourceProviderOptions

	ctx     *pulumi.Context
	typ     string
	name    string
	inputs  provider.ConstructInputs
	options pulumi.ResourceOption
}

var _ providerresource.ResourceProvider = &mockResourceProvider{}

func (t *mockResourceProvider) Construct(ctx *pulumi.Context, typ string, name string, inputs provider.ConstructInputs, options pulumi.ResourceOption) (*provider.ConstructResult, error) {
	t.ctx = ctx
	t.typ = typ
	t.name = name
	t.inputs = inputs
	t.options = options
	return t.Result, t.Err
}
