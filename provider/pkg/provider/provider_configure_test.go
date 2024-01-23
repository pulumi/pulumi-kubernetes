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
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/rpcutil"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// A mock engine for test purposes.
type mockEngine struct {
	pulumirpc.UnsafeEngineServer
	t            TestingTB
	logger       *log.Logger
	rootResource string
}

// Log logs a global message in the engine, including errors and warnings.
func (m *mockEngine) Log(ctx context.Context, in *pulumirpc.LogRequest) (*pbempty.Empty, error) {
	m.t.Logf("%s: %s", in.GetSeverity(), in.GetMessage())
	if m.logger != nil {
		m.logger.Printf("%s: %s", in.GetSeverity(), in.GetMessage())
	}
	return &pbempty.Empty{}, nil
}

// GetRootResource gets the URN of the root resource, the resource that should be the root of all
// otherwise-unparented resources.
func (m *mockEngine) GetRootResource(ctx context.Context, in *pulumirpc.GetRootResourceRequest) (*pulumirpc.GetRootResourceResponse, error) {
	return &pulumirpc.GetRootResourceResponse{
		Urn: m.rootResource,
	}, nil
}

// SetRootResource sets the URN of the root resource.
func (m *mockEngine) SetRootResource(ctx context.Context, in *pulumirpc.SetRootResourceRequest) (*pulumirpc.SetRootResourceResponse, error) {
	m.rootResource = in.GetUrn()
	return &pulumirpc.SetRootResourceResponse{}, nil
}

// newMockHost creates a mock host for test purposes and returns a client.
// Dispatches Engine RPC calls to the given engine.
func newMockHost(ctx context.Context, engine pulumirpc.EngineServer) *provider.HostClient {
	cancel := make(chan bool)
	go func() {
		<-ctx.Done()
		close(cancel)
	}()
	handle, err := rpcutil.ServeWithOptions(rpcutil.ServeOptions{
		Cancel: cancel,
		Init: func(srv *grpc.Server) error {
			pulumirpc.RegisterEngineServer(srv, engine)
			return nil
		},
		Options: rpcutil.OpenTracingServerInterceptorOptions(nil),
	})
	if err != nil {
		panic(fmt.Errorf("could not start host engine service: %v", err))
	}

	go func() {
		err := <-handle.Done
		if err != nil {
			panic(fmt.Errorf("host engine service failed: %v", err))
		}
	}()

	address := fmt.Sprintf("127.0.0.1:%v", handle.Port)
	hostClient, err := provider.NewHostClient(address)
	if err != nil {
		panic(fmt.Errorf("could not connect to host engine service: %v", err))
	}
	return hostClient
}

type providerTestContext struct {
	engine *mockEngine
	host   *provider.HostClient
}

func (c *providerTestContext) NewProvider() (*kubeProvider, error) {
	var pulumiSchema []byte
	var terraformMapping []byte
	return makeKubeProvider(c.host, "kubernetes", "v0.0.0", pulumiSchema, terraformMapping)
}

var pctx *providerTestContext

var _ = BeforeSuite(func() {
	var buff bytes.Buffer
	engine := &mockEngine{
		t:      GinkgoT(),
		logger: log.New(&buff, "\t", 0),
	}
	// t.Cleanup(func() {
	// 	log.Default().Printf("Engine Log:\n%s", buff.String())
	// })
	ctx, cancel := context.WithCancel(context.Background())
	host := newMockHost(ctx, engine)
	DeferCleanup(func() {
		cancel()
	})

	pctx = &providerTestContext{
		engine: engine,
		host:   host,
	}
})

var _ = Describe("RPC:Configure", func() {
	var k *kubeProvider
	var req *pulumirpc.ConfigureRequest

	BeforeEach(func() {
		var err error
		k, err = pctx.NewProvider()
		Expect(err).ShouldNot(HaveOccurred())

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

		clientChecks := func() {
			It("should have an initialized client", func() {
				_, err := k.Configure(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())

				By("creating a client-go config")
				Expect(k.config).ToNot(BeNil())
				Expect(k.kubeconfig).ToNot(BeNil())

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
				Expect(k.config).To(BeNil())
				Expect(k.kubeconfig).To(BeNil())
				Expect(k.clientSet).To(BeNil())
				Expect(k.logClient).To(BeNil())
			})
		}

		Describe("ambient kubeconfig", func() {
			commonChecks()
			clientChecks()
		})

		Describe("kubeconfig literal", func() {
			Context("with an invalid value", func() {
				JustBeforeEach(func() {
					req.Variables["kubernetes:config:kubeconfig"] = "invalid"
				})
				commonChecks()
				clusterUnreachableChecks()
			})

			Context("with a valid kubeconfig as a literal value", func() {
				JustBeforeEach(func() {
					contents, err := clientcmd.Write(*config)
					Expect(err).ToNot(HaveOccurred())
					req.Variables["kubernetes:config:kubeconfig"] = string(contents)
				})
				commonChecks()
				clientChecks()
			})
		})

		Describe("kubeconfig path", func() {
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
					f.WriteString("invalid")
					f.Close()
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
				clientChecks()
			})
		})
	})
})
