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
	helmcli "helm.sh/helm/v3/pkg/cli"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	providerresource "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/resource"
)

const defaultNamespace = "default"

var _ = gk.Describe("RPC:Construct", func() {
	var opts []NewProviderOption
	var k *kubeProvider
	var req *pulumirpc.ConstructRequest

	gk.BeforeEach(func() {
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

	gk.JustBeforeEach(func() {
		k = pctx.NewProvider(opts...)
		k.clientSet = &clients.DynamicClientSet{}
		k.defaultNamespace = defaultNamespace
		k.helmDriver = "memory"
		k.helmSettings = helmcli.New()
	})

	gk.Context("when the requested type is unknown", func() {
		gk.BeforeEach(func() {
			req.Type = "kubernetes:test:UnknownComponent"
			req.Name = "testComponent"
		})
		gk.It("should return an error", func() {
			_, err := k.Construct(context.Background(), req)
			gm.Expect(err).Should(gm.HaveOccurred())
		})
	})

	gk.Context("when the requested type is known", func() {
		var testComponent *mockResourceProvider
		gk.BeforeEach(func() {
			testComponent = &mockResourceProvider{
				Result: &provider.ConstructResult{
					URN: pulumi.URN("urn:pulumi:test::test::test:TestComponent::testComponent"),
				},
			}
			opts = append(opts, WithResourceProvider("kubernetes:test:TestComponent", factory(testComponent)))

			req.Type = "kubernetes:test:TestComponent"
			req.Name = "testComponent"
		})

		gk.It("should delegate to the provider", func() {
			result, err := k.Construct(context.Background(), req)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())
			gm.Expect(testComponent.typ).Should(gm.Equal("kubernetes:test:TestComponent"))
			gm.Expect(testComponent.name).Should(gm.Equal("testComponent"))
			gm.Expect(result.Urn).Should(gm.Equal("urn:pulumi:test::test::test:TestComponent::testComponent"))
		})

		gk.It("should provide options", func() {
			_, err := k.Construct(context.Background(), req)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())
			gm.Expect(testComponent.opts.ClientSet).ShouldNot(gm.BeNil())
			gm.Expect(testComponent.opts.DefaultNamespace).ShouldNot(gm.BeEmpty())
			gm.Expect(testComponent.opts.HelmOptions).ShouldNot(gm.BeNil())
		})

		gk.Context("when clusterUnreachable is true", func() {
			gk.JustBeforeEach(func() {
				k.clusterUnreachable = true
				k.clusterUnreachableReason = "testing"
			})
			gk.It("should return an error", func() {
				_, err := k.Construct(context.Background(), req)
				gm.Expect(err).To(gm.MatchError(gm.ContainSubstring("configured Kubernetes cluster is unreachable")))
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

func (t *mockResourceProvider) Construct(
	ctx *pulumi.Context,
	typ string,
	name string,
	inputs provider.ConstructInputs,
	options pulumi.ResourceOption,
) (*provider.ConstructResult, error) {
	t.ctx = ctx
	t.typ = typ
	t.name = name
	t.inputs = inputs
	t.options = options
	return t.Result, t.Err
}
