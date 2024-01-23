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
	"os/user"
	"path/filepath"
	"testing"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/rpcutil"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

// A mock engine for test purposes.
type mockEngine struct {
	pulumirpc.UnsafeEngineServer
	t            testing.TB
	logger       *log.Logger
	rootResource string
}

// Log logs a global message in the engine, including errors and warnings.
func (m *mockEngine) Log(_ context.Context, in *pulumirpc.LogRequest) (*pbempty.Empty, error) {
	m.t.Logf("%s: %s", in.GetSeverity(), in.GetMessage())
	if m.logger != nil {
		m.logger.Printf("%s: %s", in.GetSeverity(), in.GetMessage())
	}
	return &pbempty.Empty{}, nil
}

// GetRootResource gets the URN of the root resource, the resource that should be the root of all
// otherwise-unparented resources.
func (m *mockEngine) GetRootResource(_ context.Context, _ *pulumirpc.GetRootResourceRequest) (*pulumirpc.GetRootResourceResponse, error) {
	return &pulumirpc.GetRootResourceResponse{
		Urn: m.rootResource,
	}, nil
}

// SetRootResource sets the URN of the root resource.
func (m *mockEngine) SetRootResource(_ context.Context, in *pulumirpc.SetRootResourceRequest) (*pulumirpc.SetRootResourceResponse, error) {
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

func newProvider(t *testing.T) *kubeProvider {
	var buff bytes.Buffer
	engine := &mockEngine{
		t:      t,
		logger: log.New(&buff, "\t", 0),
	}
	t.Cleanup(func() {
		log.Default().Printf("Engine Log:\n%s", buff.String())
	})
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	host := newMockHost(ctx, engine)

	var terraformMapping []byte
	var pulumiSchema []byte
	k, err := makeKubeProvider(host, "kubernetes", "v0.0.0", pulumiSchema, terraformMapping)
	require.NoError(t, err)
	return k
}

func TestConfigure(t *testing.T) {
	t.Run("should return a response detailing the provider's capabilities", func(t *testing.T) {
		t.Parallel()

		k := newProvider(t)
		r, err := k.Configure(context.Background(), &pulumirpc.ConfigureRequest{})
		assert.NoError(t, err)
		assert.True(t, r.AcceptSecrets)
		assert.True(t, r.SupportsPreview)
		assert.True(t, r.AcceptResources)
		assert.True(t, r.AcceptOutputs)
	})

	t.Run("secret support", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name string
			req  *pulumirpc.ConfigureRequest
			want bool
		}{
			{
				name: "should persist true",
				req:  &pulumirpc.ConfigureRequest{AcceptSecrets: true},
				want: true,
			},
			{
				name: "should persist false",
				req:  &pulumirpc.ConfigureRequest{AcceptSecrets: false},
				want: false,
			},
		}
		for _, tt := range tests {
			k := newProvider(t)
			_, err := k.Configure(context.Background(), tt.req)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, k.enableSecrets)
		}
	})

	t.Run("connectivity", func(t *testing.T) {
		t.Parallel()

		homeDir := func() string {
			// Ignore errors. The filepath will be checked later, so we can handle failures there.
			usr, _ := user.Current()
			return usr.HomeDir
		}
		// load the ambient kubeconfig for test purposes
		config, err := clientcmd.LoadFromFile(filepath.Join(homeDir(), "/.kube/config"))
		require.NoError(t, err)
		contents, err := clientcmd.Write(*config)
		require.NoError(t, err)

		tests := []struct {
			name            string
			req             *pulumirpc.ConfigureRequest
			wantUnreachable bool
		}{
			{
				name: "ambient kubeconfig",
				req:  &pulumirpc.ConfigureRequest{},
			},
			{
				name: "invalid kubeconfig literal",
				req: &pulumirpc.ConfigureRequest{
					Variables: map[string]string{"kubernetes:config:kubeconfig": "invalid"},
				},
				wantUnreachable: true,
			},
			{
				name: "valid kubeconfig literal",

				req: &pulumirpc.ConfigureRequest{
					Variables: map[string]string{"kubernetes:config:kubeconfig": string(contents)},
				},
			},
			{
				name: "non-existent config file",
				req: &pulumirpc.ConfigureRequest{
					Variables: map[string]string{"kubernetes:config:kubeconfig": "./nosuchfile"},
				},
				wantUnreachable: true,
			},
			{
				name: "invalid config file",
				req: &pulumirpc.ConfigureRequest{
					Variables: map[string]string{"kubernetes:config:kubeconfig": "./testdata/invalid.kubeconfig"},
				},
				wantUnreachable: true,
			},
			{
				name: "valid config file",
				req: &pulumirpc.ConfigureRequest{
					Variables: map[string]string{"kubernetes:config:kubeconfig": "~/.kube/config"},
				},
			},
		}

		for _, tt := range tests {
			k := newProvider(t)

			if tt.req.Variables == nil {
				tt.req.Variables = map[string]string{}
			}
			tt.req.Variables["kubernetes:config:namespace"] = "testns"
			_, err := k.Configure(context.Background(), tt.req)
			assert.NoError(t, err)

			// common checks
			assert.Equal(t, "testns", k.defaultNamespace, "should use the configured namespace as the default")

			if tt.wantUnreachable {
				// clusterUnreachableChecks
				assert.True(t, k.clusterUnreachable)
				assert.Nil(t, k.config)
				assert.Nil(t, k.kubeconfig)
				assert.Nil(t, k.clientSet)
				assert.Nil(t, k.logClient)
			} else {
				// clientChecks
				assert.NotNil(t, k.config)
				assert.NotNil(t, k.kubeconfig)
				assert.NotNil(t, k.clientSet)
				assert.NotNil(t, k.logClient)
				assert.NotNil(t, k.k8sVersion)
				assert.NotNil(t, k.resources)

				k.invalidateResources()
				assert.Nil(t, k.resources)
				resources, err := k.getResources()
				require.NoError(t, err)
				assert.NotNil(t, resources)

				ns := config.Contexts[config.CurrentContext].Namespace
				if ns == "" {
					ns = "default"
				}
				assert.Equal(t, ns, k.defaultNamespace, "should use default namespace")
			}
		}
	})
}
