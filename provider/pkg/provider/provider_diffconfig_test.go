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

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

var _ = Describe("RPC:DiffConfig", func() {
	var k *kubeProvider
	var req *pulumirpc.DiffRequest
	var olds, news resource.PropertyMap
	var oldConfig, newConfig *clientcmdapi.Config

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
		oldConfig, err = clientcmd.LoadFromFile(filepath.Join(homeDir(), "/.kube/config"))
		Expect(err).ToNot(HaveOccurred())
		newConfig, err = clientcmd.LoadFromFile(filepath.Join(homeDir(), "/.kube/config"))
		Expect(err).ToNot(HaveOccurred())

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
		req.Olds, err = plugin.MarshalProperties(olds, plugin.MarshalOptions{
			Label: "olds", KeepUnknowns: true, SkipNulls: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
		req.News, err = plugin.MarshalProperties(news, plugin.MarshalOptions{
			Label: "news", KeepUnknowns: true, SkipNulls: true,
		})
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("Kubeconfig Parsing", func() {
		Context("when kubeconfig is a computed value", func() {
			BeforeEach(func() {
				olds["kubeconfig"] = resource.NewStringProperty(KubeconfigAsString(oldConfig))
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

		Context("when kubeconfig is a string value", func() {
			BeforeEach(func() {
				olds["kubeconfig"] = resource.NewStringProperty(KubeconfigAsString(oldConfig))
				news["kubeconfig"] = resource.NewStringProperty(KubeconfigAsString(newConfig))
			})

			It("should return an empty diff", func() {
				resp, err := k.DiffConfig(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_NONE))
				Expect(resp.Diffs).To(BeEmpty())
				Expect(resp.Replaces).To(BeEmpty())
			})
		})

		Context("when kubeconfig is a file", func() {
			BeforeEach(func() {
				olds["kubeconfig"] = resource.NewStringProperty(KubeconfigAsFile(oldConfig))
				news["kubeconfig"] = resource.NewStringProperty(KubeconfigAsFile(newConfig))
			})

			XIt("should return an empty diff", func() {
				resp, err := k.DiffConfig(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_NONE))
				Expect(resp.Diffs).To(BeEmpty())
				Expect(resp.Replaces).To(BeEmpty())
			})

			It("should fail due to pulumi/pulumi-kubernetes#2663", func() {
				_, err := k.DiffConfig(context.Background(), req)
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe("Cluster Change Detection", func() {

		Context("when the cluster info has changed", func() {
			BeforeEach(func() {
				clusterName := newConfig.Contexts[newConfig.CurrentContext].Cluster
				newConfig.Clusters[clusterName].Server = "https://updated.invalid"

				olds["kubeconfig"] = resource.NewStringProperty(KubeconfigAsString(oldConfig))
				news["kubeconfig"] = resource.NewStringProperty(KubeconfigAsString(newConfig))
			})

			It("should suggest replacement since the cluster itself may have been replaced", func() {
				resp, err := k.DiffConfig(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.Changes).To(Equal(pulumirpc.DiffResponse_DIFF_SOME))
				Expect(resp.Diffs).To(ContainElements("kubeconfig"))
				Expect(resp.Replaces).To(ContainElements("kubeconfig"))
			})
		})
	})
})
