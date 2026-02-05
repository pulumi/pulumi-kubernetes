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
	"encoding/json"
	"os"

	_ "embed"

	gk "github.com/onsi/ginkgo/v2"
	gm "github.com/onsi/gomega"
	gs "github.com/onsi/gomega/gstruct"
	kubeversion "k8s.io/apimachinery/pkg/version"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/utils/ptr"

	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

var _ = gk.Describe("RPC:Configure", func() {
	var opts []NewProviderOption
	var k *kubeProvider
	var req *pulumirpc.ConfigureRequest
	var ambient *clientcmdapi.Config

	gk.BeforeEach(func() {
		opts = []NewProviderOption{}

		// initialize the ConfigureRequest to be customized in nested BeforeEach blocks
		req = &pulumirpc.ConfigureRequest{
			AcceptSecrets: true,
			Variables:     map[string]string{},
		}

		// initialize a fake ambient kubeconfig
		ambient = pctx.NewConfig(WithContext("context1"))
	})

	gk.JustBeforeEach(func() {
		k = pctx.NewProvider(opts...)

		// set the KUBECONFIG environment variable
		path := WriteKubeconfigToFile(ambient)
		os.Setenv("KUBECONFIG", path)
		gk.DeferCleanup(func() {
			os.Unsetenv("KUBECONFIG")
		})
	})

	gk.It("should return a response detailing the provider's capabilities", func() {
		r, err := k.Configure(context.Background(), req)
		gm.Expect(err).ShouldNot(gm.HaveOccurred())
		gm.Expect(r.AcceptSecrets).Should(gm.BeTrue())
		gm.Expect(r.SupportsPreview).Should(gm.BeTrue())
		gm.Expect(r.AcceptResources).Should(gm.BeFalse())
		gm.Expect(r.AcceptOutputs).Should(gm.BeFalse())
	})

	gk.Describe("Secrets Support", func() {
		gk.Context("when configured to support secrets", func() {
			gk.BeforeEach(func() {
				req.AcceptSecrets = true
			})
			gk.It("should enable secrets support in subsequent RPC methods", func() {
				_, err := k.Configure(context.Background(), req)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(k.enableSecrets).Should(gm.BeTrue())
			})
		})

		gk.Context("when configured to NOT support secrets", func() {
			gk.BeforeEach(func() {
				req.AcceptSecrets = false
			})
			gk.It("should not enable secrets support in subsequent RPC methods", func() {
				_, err := k.Configure(context.Background(), req)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(k.enableSecrets).Should(gm.BeFalse())
			})
		})
	})

	gk.Describe("Namespacing", func() {
		gk.Context("when configured to use a particular namespace", func() {
			gk.JustBeforeEach(func() {
				req.Variables["kubernetes:config:namespace"] = "pulumi"
			})
			gk.It("should use the configured namespace as the default namespace", func() {
				_, err := k.Configure(context.Background(), req)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(k.defaultNamespace).To(gm.Equal("pulumi"))
				helmFlags := k.helmSettings.RESTClientGetter().(*genericclioptions.ConfigFlags)
				gm.Expect(helmFlags.Namespace).To(gs.PointTo(gm.Equal("pulumi")))
			})
		})
	})

	gk.Describe("Kubeconfig Parsing", func() {
		var other *clientcmdapi.Config

		gk.BeforeEach(func() {
			// make a (fake) kubeconfig to serve as the value of the 'kubeconfig' provider property
			other = pctx.NewConfig(WithContext("context2"), WithNamespace("other"))
		})

		// Define some "shared behaviors" that will be used to test various use cases.
		// pattern: https://onsi.github.io/ginkgo/#shared-behaviors

		commonChecks := func() {}

		connectedChecks := func(expectedNS string) {
			gk.It("should have an initialized client", func() {
				_, err := k.Configure(context.Background(), req)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())

				gk.By("creating strongly-typed clients")
				gm.Expect(k.clientSet).ToNot(gm.BeNil())
				gm.Expect(k.logClient).ToNot(gm.BeNil())
			})

			gk.It("should use the context namespace as the default namespace", func() {
				_, err := k.Configure(context.Background(), req)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(k.defaultNamespace).To(gm.Equal(expectedNS))
			})

			gk.It("should provide Helm settings", func() {
				_, err := k.Configure(context.Background(), req)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(k.helmSettings).ToNot(gm.BeNil())
			})
		}

		clusterUnreachableChecks := func() {
			gk.It("should be in clusterUnreachable mode", func() {
				_, err := k.Configure(context.Background(), req)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(k.clusterUnreachable).To(gm.BeTrue())
				gm.Expect(k.clientSet).ToNot(gm.BeNil())
				gm.Expect(k.logClient).ToNot(gm.BeNil())
			})
		}

		gk.Describe("use case: ambient kubeconfig", func() {
			commonChecks()
			connectedChecks("default")
		})

		gk.Describe("use case: kubeconfig string", func() {
			gk.Context("with an invalid value", func() {
				gk.JustBeforeEach(func() {
					req.Variables["kubernetes:config:kubeconfig"] = "invalid"
				})
				commonChecks()
				clusterUnreachableChecks()
			})

			gk.Context("with a valid kubeconfig as a string value", func() {
				gk.JustBeforeEach(func() {
					req.Variables["kubernetes:config:kubeconfig"] = WriteKubeconfigToString(other)
				})
				commonChecks()
				connectedChecks("other")

				gk.It("should set Helm's --kubeconfig", func() {
					_, err := k.Configure(context.Background(), req)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(k.helmSettings.KubeConfig).ToNot(gm.BeEmpty())
				})
			})
		})

		gk.Describe("use case: kubeconfig file", func() {
			gk.Context("with a non-existent config file", func() {
				gk.BeforeEach(func() {
					req.Variables["kubernetes:config:kubeconfig"] = "./nosuchfile"
				})
				commonChecks()
				clusterUnreachableChecks()
			})

			gk.Context("with an invalid config file", func() {
				gk.BeforeEach(func() {
					f, err := os.CreateTemp("", "kubeconfig-")
					gm.Expect(err).ToNot(gm.HaveOccurred())
					gk.DeferCleanup(func() {
						os.Remove(f.Name())
					})
					_, err = f.WriteString("invalid")
					gm.Expect(err).ToNot(gm.HaveOccurred())
					err = f.Close()
					gm.Expect(err).ToNot(gm.HaveOccurred())
					req.Variables["kubernetes:config:kubeconfig"] = f.Name()
				})
				commonChecks()
				clusterUnreachableChecks()
			})

			gk.Context("with a valid config file", func() {
				gk.BeforeEach(func() {
					req.Variables["kubernetes:config:kubeconfig"] = WriteKubeconfigToFile(other)
				})
				commonChecks()
				connectedChecks("other")

				gk.It("should set Helm's --kubeconfig", func() {
					_, err := k.Configure(context.Background(), req)
					gm.Expect(err).ShouldNot(gm.HaveOccurred())
					gm.Expect(k.helmSettings.KubeConfig).ToNot(gm.BeEmpty())
				})
			})
		})
	})

	gk.Describe("Kube Context", func() {
		gk.Context("when configured to use a particular context", func() {
			gk.JustBeforeEach(func() {
				req.Variables["kubernetes:config:context"] = "context2"
			})
			gk.It("should use the configured context", func() {
				_, err := k.Configure(context.Background(), req)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(k.helmSettings.KubeContext).To(gm.Equal("context2"))
			})
		})
	})

	gk.Describe("Kube Cluster", func() {
		gk.Context("when configured to use a particular cluster", func() {
			gk.JustBeforeEach(func() {
				req.Variables["kubernetes:config:cluster"] = "cluster2"
			})
			gk.It("should use the configured context", func() {
				_, err := k.Configure(context.Background(), req)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				helmFlags := k.helmSettings.RESTClientGetter().(*genericclioptions.ConfigFlags)
				gm.Expect(helmFlags.ClusterName).To(gs.PointTo(gm.Equal("cluster2")))
			})
		})
	})

	gk.Describe("Kube Client Settings", func() {
		gk.Context("when configured with Kube client settings", func() {
			var kubeClientSettings *KubeClientSettings
			gk.BeforeEach(func() {
				kubeClientSettings = &KubeClientSettings{
					Burst:   ptr.To(42),
					QPS:     ptr.To(42.),
					Timeout: ptr.To(42),
				}
			})
			gk.JustBeforeEach(func() {
				data, _ := json.Marshal(kubeClientSettings)
				req.Variables["kubernetes:config:kubeClientSettings"] = string(data)
			})
			gk.It("should use the configured settings", func() {
				_, err := k.Configure(context.Background(), req)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				helmFlags := k.helmSettings.RESTClientGetter().(*genericclioptions.ConfigFlags)
				gm.Expect(k.helmSettings.BurstLimit).To(gm.Equal(42))
				gm.Expect(k.helmSettings.QPS).To(gm.Equal(float32(42.)))
				gm.Expect(helmFlags.Timeout).To(gs.PointTo(gm.Equal("42")))
			})
		})
	})

	gk.Describe("Helm Release Settings", func() {
		gk.Context("given helmReleaseSettings", func() {
			var helmReleaseSettings *HelmReleaseSettings
			gk.BeforeEach(func() {
				helmReleaseSettings = &HelmReleaseSettings{
					Driver:               ptr.To("configmap"),
					PluginsPath:          ptr.To("plugins"),
					RegistryConfigPath:   ptr.To("registry"),
					RepositoryCache:      ptr.To("cache"),
					RepositoryConfigPath: ptr.To("config"),
				}
			})
			gk.JustBeforeEach(func() {
				data, _ := json.Marshal(helmReleaseSettings)
				req.Variables["kubernetes:config:helmReleaseSettings"] = string(data)
			})
			gk.It("should use the configured settings", func() {
				_, err := k.Configure(context.Background(), req)
				gm.Expect(err).ShouldNot(gm.HaveOccurred())
				gm.Expect(k.helmDriver).To(gm.Equal("configmap"))
				gm.Expect(k.helmSettings.PluginsDirectory).To(gm.Equal("plugins"))
				gm.Expect(k.helmSettings.RegistryConfig).To(gm.Equal("registry"))
				gm.Expect(k.helmSettings.RepositoryCache).To(gm.Equal("cache"))
				gm.Expect(k.helmSettings.RepositoryConfig).To(gm.Equal("config"))
			})
		})
	})

	gk.Describe("Discovery", func() {
		gk.It("should record the server version for use in subsequent RPC methods", func() {
			_, err := k.Configure(context.Background(), req)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())
			gm.Expect(k.k8sVersion).ToNot(gm.BeNil())
		})

		gk.It("should initialize a resource cache", func() {
			_, err := k.Configure(context.Background(), req)
			gm.Expect(err).ShouldNot(gm.HaveOccurred())

			gk.By("discovering the server resources")
			gm.Expect(k.resources).ToNot(gm.BeNil())

			gk.By("supporting invalidation")
			k.invalidateResources()
			gm.Expect(k.resources).To(gm.BeNil())
			gm.Expect(k.getResources()).ToNot(gm.BeNil())
		})

		gk.Context("when the server version is < 1.13", func() {
			gk.BeforeEach(func() {
				opts = append(opts, WithServerVersion(kubeversion.Info{Major: "1", Minor: "12"}))
			})

			gk.It("should fail to configure", func() {
				_, err := k.Configure(context.Background(), req)
				gm.Expect(err).Should(gm.HaveOccurred())
			})
		})
	})
})
