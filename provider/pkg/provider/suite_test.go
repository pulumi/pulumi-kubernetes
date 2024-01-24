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
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/rpcutil"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"google.golang.org/grpc"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "provider/pkg/provider")
}

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
	// mock engine
	var buff bytes.Buffer
	engine := &mockEngine{
		t:      GinkgoT(),
		logger: log.New(&buff, "\t", 0),
	}

	// mock host
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
