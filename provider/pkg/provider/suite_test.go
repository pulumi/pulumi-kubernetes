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
	"bytes"
	"context"
	"log"
	"os"
	"testing"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	fakeclients "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
	fakehost "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/host/fake"
	providerresource "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/provider/resource"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"k8s.io/apimachinery/pkg/runtime"
	kubeversion "k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	testPluginVersion    = "v0.0.0"
	testPulumiSchema     = "{}"
	testTerraformMapping = "{}"
)

// CheckFailure matches a CheckFailure by property and reason.
func CheckFailure(prop string, reason gomegatypes.GomegaMatcher) gomegatypes.GomegaMatcher {
	return And(
		WithTransform(func(failure *pulumirpc.CheckFailure) string {
			return failure.GetProperty()
		}, Equal(prop)),
		WithTransform(func(failure *pulumirpc.CheckFailure) string {
			return failure.GetReason()
		}, reason))
}

// WriteKubeconfigToString converts a clientcmdapi.Config to a string.
func WriteKubeconfigToString(config *clientcmdapi.Config) string {
	contents, err := clientcmd.Write(*config)
	Expect(err).ToNot(HaveOccurred())
	return string(contents)
}

// WriteKubeconfigToFile converts a clientcmdapi.Config to a temporary file.
func WriteKubeconfigToFile(config *clientcmdapi.Config) string {
	f, err := os.CreateTemp("", "kubeconfig-")
	Expect(err).ToNot(HaveOccurred())
	DeferCleanup(func() {
		os.Remove(f.Name())
	})
	err = f.Close()
	Expect(err).ToNot(HaveOccurred())
	err = clientcmd.WriteToFile(*config, f.Name())
	Expect(err).ToNot(HaveOccurred())
	return f.Name()
}

// A mock engine for test purposes.
type mockEngine struct {
	pulumirpc.UnimplementedEngineServer
	t            testing.TB
	logger       *log.Logger
	rootResource string
}

var _ pulumirpc.EngineServer = &mockEngine{}

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
func (m *mockEngine) GetRootResource(
	ctx context.Context,
	in *pulumirpc.GetRootResourceRequest,
) (*pulumirpc.GetRootResourceResponse, error) {
	return &pulumirpc.GetRootResourceResponse{
		Urn: m.rootResource,
	}, nil
}

// SetRootResource sets the URN of the root resource.
func (m *mockEngine) SetRootResource(
	ctx context.Context,
	in *pulumirpc.SetRootResourceRequest,
) (*pulumirpc.SetRootResourceResponse, error) {
	m.rootResource = in.GetUrn()
	return &pulumirpc.SetRootResourceResponse{}, nil
}

// newMockHost creates a mock host for test purposes and returns a client.
// Dispatches Engine RPC calls to the given engine.
func newMockHost(engine pulumirpc.EngineServer) *fakehost.HostClient {
	return &fakehost.HostClient{
		Engine: engine,
	}
}

type providerTestContext struct {
	engine *mockEngine
	host   *fakehost.HostClient
}

type NewConfigOption func(*newConfigOptions)

type newConfigOptions struct {
	CurrentContext   string
	CurrentNamespace string
}

func WithContext(context string) NewConfigOption {
	return func(options *newConfigOptions) {
		options.CurrentContext = context
	}
}

func WithNamespace(namespace string) NewConfigOption {
	return func(options *newConfigOptions) {
		options.CurrentNamespace = namespace
	}
}

func (c *providerTestContext) NewConfig(opts ...NewConfigOption) *clientcmdapi.Config {
	options := newConfigOptions{
		CurrentContext: "context1",
	}
	for _, opt := range opts {
		opt(&options)
	}

	config := &clientcmdapi.Config{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: map[string]*clientcmdapi.Cluster{
			"cluster1": {
				Server: "https://cluster1.test",
			},
			"cluster2": {
				Server: "https://cluster2.test",
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"user1": {
				Token: "secret",
			},
			"user2": {
				Token: "secret",
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"context1": {
				Cluster:  "cluster1",
				AuthInfo: "user1",
			},
			"context2": {
				Cluster:  "cluster2",
				AuthInfo: "user2",
			},
		},
		CurrentContext: options.CurrentContext,
	}

	if options.CurrentNamespace != "" {
		config.Contexts[options.CurrentContext].Namespace = options.CurrentNamespace
	}

	return config
}

type NewProviderOption func(*newProviderOptions)

type newProviderOptions struct {
	copts     []fakeclients.NewDynamicClientOption
	objects   []runtime.Object
	providers map[string]providerresource.ResourceProviderFactory
}

func WithObjects(objects ...runtime.Object) NewProviderOption {
	return func(options *newProviderOptions) {
		options.copts = append(options.copts, fakeclients.WithObjects(objects...))
		options.objects = append(options.objects, objects...)
	}
}

func WithServerVersion(version kubeversion.Info) NewProviderOption {
	return func(options *newProviderOptions) {
		options.copts = append(options.copts, fakeclients.WithServerVersion(version))
	}
}

func WithResourceProvider(typ string, provider providerresource.ResourceProviderFactory) NewProviderOption {
	return func(options *newProviderOptions) {
		if options.providers == nil {
			options.providers = make(map[string]providerresource.ResourceProviderFactory)
		}
		options.providers[typ] = provider
	}
}

func (c *providerTestContext) NewProvider(opts ...NewProviderOption) *kubeProvider {
	options := newProviderOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	k, err := makeKubeProvider(
		c.host,
		"kubernetes",
		testPluginVersion,
		[]byte(testPulumiSchema),
		[]byte(testTerraformMapping),
	)
	Expect(err).ShouldNot(HaveOccurred())

	k.makeClient = func(ctx context.Context, config *rest.Config) (*clients.DynamicClientSet, *clients.LogClient, error) {
		dc, _, _, _ := fakeclients.NewSimpleDynamicClient(options.copts...)
		lc, _ := fakeclients.NewSimpleLogClient(ctx, options.objects...)
		return dc, lc, nil
	}

	k.resourceProviders = options.providers

	return k
}

var pctx *providerTestContext

var _ = BeforeSuite(func() {
	// make a mock engine that simply buffers the log messages.
	var buff bytes.Buffer
	engine := &mockEngine{
		t:      GinkgoTB(),
		logger: log.New(&buff, "\t", 0),
	}

	// make a mock host as an RPC server for the mock engine.
	host := newMockHost(engine)

	// make a suite-level context for use in test specs, e.g. to make a provider instance.
	pctx = &providerTestContext{
		engine: engine,
		host:   host,
	}
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "provider/pkg/provider")
}
