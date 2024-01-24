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
	"os"
	"os/user"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/tools/clientcmd"

	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var _ = Describe("RPC:Configure", func() {
	var k *kubeProvider
	var req *pulumirpc.ConfigureRequest

	BeforeEach(func() {
		var err error
		k, err = pctx.NewProvider()
		Expect(err).ShouldNot(HaveOccurred())

		// initialize the ConfigureRequest to be customized in nested BeforeEach blocks
		req = &pulumirpc.ConfigureRequest{
			AcceptSecrets: true,
			Variables:     map[string]string{},
		}
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

	Describe("Connectivity", func() {
		var config *clientcmdapi.Config

		BeforeEach(func() {
			var err error
			homeDir := func() string {
				// Ignore errors. The filepath will be checked later, so we can handle failures there.
				usr, _ := user.Current()
				return usr.HomeDir
			}
			// load the ambient kubeconfig for test purposes
			config, err = clientcmd.LoadFromFile(filepath.Join(homeDir(), "/.kube/config"))
			Expect(err).ToNot(HaveOccurred())
		})

		// Define some "shared behaviors" that will be used to test various use cases.
		// pattern: https://onsi.github.io/ginkgo/#shared-behaviors

		commonChecks := func() {
			Context("when configured to use a particular namespace", func() {
				JustBeforeEach(func() {
					req.Variables["kubernetes:config:namespace"] = "testns"
				})
				It("should use the configured namespace as the default namespace", func() {
					_, err := k.Configure(context.Background(), req)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(k.defaultNamespace).To(Equal("testns"))
				})
			})
		}

		connectedChecks := func() {
			It("should have an initialized client", func() {
				_, err := k.Configure(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())

				By("creating strongly-typed clients")
				Expect(k.clientSet).ToNot(BeNil())
				Expect(k.logClient).ToNot(BeNil())

				By("discovering the server version")
				Expect(k.k8sVersion).ToNot(BeNil())
			})

			It("should provide a resource cache", func() {
				_, err := k.Configure(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())

				By("discovering the server resources")
				Expect(k.resources).ToNot(BeNil())

				By("supporting invalidation")
				k.invalidateResources()
				Expect(k.resources).To(BeNil())
				Expect(k.getResources()).ToNot(BeNil())
			})

			It("should use the context namespace as the default namespace", func() {
				_, err := k.Configure(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				ns := config.Contexts[config.CurrentContext].Namespace
				if ns == "" {
					ns = "default"
				}
				Expect(k.defaultNamespace).To(Equal(ns))
			})
		}

		clusterUnreachableChecks := func() {
			It("should be in clusterUnreachable mode", func() {
				_, err := k.Configure(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(k.clusterUnreachable).To(BeTrue())
				Expect(k.clientSet).To(BeNil())
				Expect(k.logClient).To(BeNil())
			})
		}

		Describe("use case: ambient kubeconfig", func() {
			commonChecks()
			connectedChecks()
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
					req.Variables["kubernetes:config:kubeconfig"] = KubeconfigAsFile(config)
				})
				commonChecks()
				connectedChecks()
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
					f, _ := os.CreateTemp("", "kubeconfig-")
					DeferCleanup(func() {
						os.Remove(f.Name())
					})
					_, _ = f.WriteString("invalid")
					_ = f.Close()
					req.Variables["kubernetes:config:kubeconfig"] = f.Name()
				})
				commonChecks()
				clusterUnreachableChecks()
			})

			Context("with a valid config file", func() {
				BeforeEach(func() {
					req.Variables["kubernetes:config:kubeconfig"] = "~/.kube/config"
				})
				commonChecks()
				connectedChecks()
			})
		})
	})
})
