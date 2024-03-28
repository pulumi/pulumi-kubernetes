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
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	kubeversion "k8s.io/apimachinery/pkg/version"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var _ = Describe("RPC:Configure", func() {
	var opts []NewProviderOption
	var k *kubeProvider
	var req *pulumirpc.ConfigureRequest
	var ambient *clientcmdapi.Config

	BeforeEach(func() {
		opts = []NewProviderOption{}

		// initialize the ConfigureRequest to be customized in nested BeforeEach blocks
		req = &pulumirpc.ConfigureRequest{
			AcceptSecrets: true,
			Variables:     map[string]string{},
		}

		// initialize a fake ambient kubeconfig
		ambient = pctx.NewConfig(WithContext("context1"))
	})

	JustBeforeEach(func() {
		k = pctx.NewProvider(opts...)

		// set the KUBECONFIG environment variable
		path := WriteKubeconfigToFile(ambient)
		os.Setenv("KUBECONFIG", path)
		DeferCleanup(func() {
			os.Unsetenv("KUBECONFIG")
		})
	})

	It("should return a response detailing the provider's capabilities", func() {
		r, err := k.Configure(context.Background(), req)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(r.AcceptSecrets).Should(BeTrue())
		Expect(r.SupportsPreview).Should(BeTrue())
		Expect(r.AcceptResources).Should(BeFalse())
		Expect(r.AcceptOutputs).Should(BeFalse())
	})

	Describe("Secrets Support", func() {
		Context("when configured to support secrets", func() {
			BeforeEach(func() {
				req.AcceptSecrets = true
			})
			It("should enable secrets support in subsequent RPC methods", func() {
				_, err := k.Configure(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(k.enableSecrets).Should(BeTrue())
			})
		})

		Context("when configured to NOT support secrets", func() {
			BeforeEach(func() {
				req.AcceptSecrets = false
			})
			It("should not enable secrets support in subsequent RPC methods", func() {
				_, err := k.Configure(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(k.enableSecrets).Should(BeFalse())
			})
		})
	})

	Describe("Kubeconfig Parsing", func() {
		var other *clientcmdapi.Config

		BeforeEach(func() {
			// make a (fake) kubeconfig to serve as the value of the 'kubeconfig' provider property
			other = pctx.NewConfig(WithContext("context2"), WithNamespace("other"))
		})

		// Define some "shared behaviors" that will be used to test various use cases.
		// pattern: https://onsi.github.io/ginkgo/#shared-behaviors

		commonChecks := func() {
			Context("when configured to use a particular namespace", func() {
				JustBeforeEach(func() {
					req.Variables["kubernetes:config:namespace"] = "pulumi"
				})
				It("should use the configured namespace as the default namespace", func() {
					_, err := k.Configure(context.Background(), req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(k.defaultNamespace).To(Equal("pulumi"))
				})
			})
		}

		connectedChecks := func(expectedNS string) {
			It("should have an initialized client", func() {
				_, err := k.Configure(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())

				By("creating strongly-typed clients")
				Expect(k.clientSet).ToNot(BeNil())
				Expect(k.logClient).ToNot(BeNil())
			})

			It("should use the context namespace as the default namespace", func() {
				_, err := k.Configure(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(k.defaultNamespace).To(Equal(expectedNS))
			})
		}

		clusterUnreachableChecks := func() {
			It("should be in clusterUnreachable mode", func() {
				_, err := k.Configure(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(k.clusterUnreachable).To(BeTrue())
				Expect(k.clientSet).ToNot(BeNil())
				Expect(k.logClient).ToNot(BeNil())
			})
		}

		Describe("use case: ambient kubeconfig", func() {
			commonChecks()
			connectedChecks("default")
		})

		Describe("use case: kubeconfig string", func() {
			Context("with an invalid value", func() {
				JustBeforeEach(func() {
					req.Variables["kubernetes:config:kubeconfig"] = "invalid"
				})
				commonChecks()
				clusterUnreachableChecks()
			})

			Context("with a valid kubeconfig as a string value", func() {
				JustBeforeEach(func() {
					req.Variables["kubernetes:config:kubeconfig"] = WriteKubeconfigToString(other)
				})
				commonChecks()
				connectedChecks("other")
			})
		})

		Describe("use case: kubeconfig file", func() {
			Context("with a non-existent config file", func() {
				BeforeEach(func() {
					req.Variables["kubernetes:config:kubeconfig"] = "./nosuchfile"
				})
				commonChecks()
				clusterUnreachableChecks()
			})

			Context("with an invalid config file", func() {
				BeforeEach(func() {
					f, err := os.CreateTemp("", "kubeconfig-")
					Expect(err).ToNot(HaveOccurred())
					DeferCleanup(func() {
						os.Remove(f.Name())
					})
					_, err = f.WriteString("invalid")
					Expect(err).ToNot(HaveOccurred())
					err = f.Close()
					Expect(err).ToNot(HaveOccurred())
					req.Variables["kubernetes:config:kubeconfig"] = f.Name()
				})
				commonChecks()
				clusterUnreachableChecks()
			})

			Context("with a valid config file", func() {
				BeforeEach(func() {
					req.Variables["kubernetes:config:kubeconfig"] = WriteKubeconfigToFile(other)
				})
				commonChecks()
				connectedChecks("other")
			})
		})
	})

	Describe("Discovery", func() {
		It("should record the server version for use in subsequent RPC methods", func() {
			_, err := k.Configure(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(k.k8sVersion).ToNot(BeNil())
		})

		It("should initialize a resource cache", func() {
			_, err := k.Configure(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())

			By("discovering the server resources")
			Expect(k.resources).ToNot(BeNil())

			By("supporting invalidation")
			k.invalidateResources()
			Expect(k.resources).To(BeNil())
			Expect(k.getResources()).ToNot(BeNil())
		})

		Context("when the server version is < 1.13", func() {
			BeforeEach(func() {
				opts = append(opts, WithServerVersion(kubeversion.Info{Major: "1", Minor: "12"}))
			})

			It("should fail to configure", func() {
				_, err := k.Configure(context.Background(), req)
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
