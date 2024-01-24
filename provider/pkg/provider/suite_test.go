package provider

import (
	"bytes"
	"context"
	"log"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
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
