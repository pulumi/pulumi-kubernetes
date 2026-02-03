// Copyright 2016-2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fake

import (
	"context"
	"fmt"
	"sync"
	"testing"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/pulumi/pulumi/sdk/v3/go/common/util/rpcutil"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

func NewEngineServer(t testing.TB) *EngineServer {
	return &EngineServer{t: t}
}

func StartEngineServer(t testing.TB, engine pulumirpc.EngineServer) (addr string) {
	cancel := make(chan bool)
	t.Cleanup(func() {
		close(cancel)
	})
	handle, err := rpcutil.ServeWithOptions(rpcutil.ServeOptions{
		Cancel: cancel,
		Init: func(srv *grpc.Server) error {
			pulumirpc.RegisterEngineServer(srv, engine)
			return nil
		},
		Options: rpcutil.OpenTracingServerInterceptorOptions(nil),
	})
	if err != nil {
		t.Fatalf("could not start host engine service: %s", err)
	}

	go func() {
		err := <-handle.Done
		if err != nil {
			t.Errorf("host engine service failed: %s", err)
		}
	}()

	return fmt.Sprintf("127.0.0.1:%v", handle.Port)
}

func ConnectToEngine(t testing.TB, addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(rpcutil.OpenTracingClientInterceptor()),
		grpc.WithStreamInterceptor(rpcutil.OpenTracingStreamClientInterceptor()),
		rpcutil.GrpcChannelOptions(),
	)
	if err != nil {
		t.Fatalf("could not connect to host engine service: %s", err)
	}
	t.Cleanup(func() {
		conn.Close()
	})

	return conn
}

// A fake engine server for test purposes.
type EngineServer struct {
	pulumirpc.UnimplementedEngineServer
	t testing.TB

	mu           sync.Mutex
	rootResource string
	logs         []*pulumirpc.LogRequest
}

func (m *EngineServer) Logs() []*pulumirpc.LogRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	var logs []*pulumirpc.LogRequest
	logs = append(logs, m.logs...)
	return logs
}

var _ pulumirpc.EngineServer = &EngineServer{}

// Log logs a global message in the engine, including errors and warnings.
func (m *EngineServer) Log(_ context.Context, in *pulumirpc.LogRequest) (*pbempty.Empty, error) {
	m.t.Logf("%s: %s", in.GetSeverity(), in.GetMessage())
	m.mu.Lock()
	m.logs = append(m.logs, in)
	m.mu.Unlock()
	return &pbempty.Empty{}, nil
}

// GetRootResource gets the URN of the root resource, the resource that should be the root of all
// otherwise-unparented resources.
func (m *EngineServer) GetRootResource(
	_ context.Context,
	_ *pulumirpc.GetRootResourceRequest,
) (*pulumirpc.GetRootResourceResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return &pulumirpc.GetRootResourceResponse{
		Urn: m.rootResource,
	}, nil
}

// SetRootResource sets the URN of the root resource.
func (m *EngineServer) SetRootResource(
	_ context.Context,
	in *pulumirpc.SetRootResourceRequest,
) (*pulumirpc.SetRootResourceResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rootResource = in.GetUrn()
	return &pulumirpc.SetRootResourceResponse{}, nil
}
