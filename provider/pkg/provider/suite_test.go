package provider

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/rpcutil"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubeversion "k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	discoveryfake "k8s.io/client-go/discovery/fake"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	kubetesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/kubectl/pkg/scheme"
)

const (
	testPluginVersion    = "v0.0.0"
	testPulumiSchema     = "{}"
	testTerraformMapping = "{}"
)

var (
	testServerVersion = kubeversion.Info{Major: "1", Minor: "29"}
)

// TestingTB is an interface that describes the implementation of the testing object.
// Using an interface that describes testing.TB instead of the actual implementation
// makes testutil usable in a wider variety of contexts (e.g. use with ginkgo : https://godoc.org/github.com/onsi/ginkgo#GinkgoT)
type TestingTB interface {
	Cleanup(func())
	Failed() bool
	Logf(format string, args ...interface{})
	Name() string
}

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

// KubeconfigAsString converts a Kubernetes configuration to a string.
func KubeconfigAsString(config *clientcmdapi.Config) string {
	contents, err := clientcmd.Write(*config)
	Expect(err).ToNot(HaveOccurred())
	return string(contents)
}

// KubeconfigAsFile converts a Kubernetes configuration to a string.
func KubeconfigAsFile(config *clientcmdapi.Config) string {
	f, _ := os.CreateTemp("", "kubeconfig-")
	DeferCleanup(func() {
		os.Remove(f.Name())
	})
	_ = f.Close()
	err := clientcmd.WriteToFile(*config, f.Name())
	Expect(err).ToNot(HaveOccurred())
	return f.Name()
}

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

type mockCachedDiscoveryClient struct {
	discovery.DiscoveryInterface
}

var _ discovery.CachedDiscoveryInterface = &mockCachedDiscoveryClient{}

func (*mockCachedDiscoveryClient) Fresh() bool {
	return true
}

func (*mockCachedDiscoveryClient) Invalidate() {}

type mockRestMapper struct {
	meta.RESTMapper
}

var _ meta.ResettableRESTMapper = &mockRestMapper{}

func (m *mockRestMapper) Reset() {}

type providerTestContext struct {
	engine *mockEngine
	host   *provider.HostClient
}

func (c *providerTestContext) NewProvider(objects ...runtime.Object) (*kubeProvider, error) {
	k, err := makeKubeProvider(c.host, "kubernetes", testPluginVersion, []byte(testPulumiSchema), []byte(testTerraformMapping))
	if err != nil {
		return nil, err
	}
	k.makeClient = func(ctx context.Context, config *rest.Config) (*clients.DynamicClientSet, *clients.LogClient, error) {
		// make a fake clientset for testing purposes, backed by an testing.ObjectTracker with pre-populated objects.
		// see also: https://github.com/kubernetes/client-go/blob/kubernetes-1.29.0/examples/fake-client/main_test.go
		disco := &discoveryfake.FakeDiscovery{
			Fake:               &kubetesting.Fake{},
			FakedServerVersion: &testServerVersion,
		}
		mapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{})
		client := dynamicfake.NewSimpleDynamicClient(scheme.Scheme, objects...)
		cs := &clients.DynamicClientSet{
			GenericClient:         client,
			DiscoveryClientCached: &mockCachedDiscoveryClient{DiscoveryInterface: disco},
			RESTMapper:            &mockRestMapper{RESTMapper: mapper},
		}

		// make a fake log client for testing purposes.
		clientset := fake.NewSimpleClientset(objects...)
		lc := clients.NewLogClient(ctx, clientset.CoreV1())

		return cs, lc, nil
	}
	return k, nil
}

// getDiscoveryClient returns the fake discovery client that the provider is using.
// use the Resources field to populate the server resources.
func getDiscoveryClient(k *kubeProvider) *discoveryfake.FakeDiscovery {
	return k.clientSet.DiscoveryClientCached.(*mockCachedDiscoveryClient).DiscoveryInterface.(*discoveryfake.FakeDiscovery)
}

// getDynamicClient returns the fake dynamic client that the provider is using.
// Use the Tracker() method to inject fake objects.
func getDynamicClient(k *kubeProvider) *dynamicfake.FakeDynamicClient {
	return k.clientSet.GenericClient.(*dynamicfake.FakeDynamicClient)
}

var pctx *providerTestContext

var _ = BeforeSuite(func() {
	// make a mock engine that simply buffers the log messages.
	var buff bytes.Buffer
	engine := &mockEngine{
		t:      GinkgoT(),
		logger: log.New(&buff, "\t", 0),
	}

	// make a mock host as an RPC server for the mock engine.
	ctx, cancel := context.WithCancel(context.Background())
	host := newMockHost(ctx, engine)
	DeferCleanup(func() {
		cancel()
	})

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
