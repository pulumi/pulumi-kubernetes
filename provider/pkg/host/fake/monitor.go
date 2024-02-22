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

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/rpcutil"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SimpleMonitor struct {
	CallF        func(args pulumi.MockCallArgs) (resource.PropertyMap, error)
	NewResourceF func(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error)
}

func (m *SimpleMonitor) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	if m.CallF == nil {
		return resource.PropertyMap{}, nil
	}
	return m.CallF(args)
}

func (m *SimpleMonitor) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	if m.NewResourceF == nil {
		return args.Name, args.Inputs, nil
	}
	return m.NewResourceF(args)
}

func StartMonitorServer(t testing.TB, monitor pulumirpc.ResourceMonitorServer) (addr string) {
	cancel := make(chan bool)
	t.Cleanup(func() {
		close(cancel)
	})
	handle, err := rpcutil.ServeWithOptions(rpcutil.ServeOptions{
		Cancel: cancel,
		Init: func(srv *grpc.Server) error {
			pulumirpc.RegisterResourceMonitorServer(srv, monitor)
			return nil
		},
		Options: rpcutil.OpenTracingServerInterceptorOptions(nil),
	})
	if err != nil {
		t.Fatalf("could not start resource monitor service: %s", err)
	}

	go func() {
		err := <-handle.Done
		if err != nil {
			t.Errorf("resource monitor service failed: %s", err)
		}
	}()

	return fmt.Sprintf("127.0.0.1:%v", handle.Port)
}

func NewResourceMonitorServer(t testing.TB, monitor pulumi.MockResourceMonitor) *ResourceMonitorServer {
	return &ResourceMonitorServer{
		t:         t,
		project:   "project",
		stack:     "stack",
		mocks:     monitor,
		resources: map[string]resource.PropertyMap{},
	}
}

type ResourceMonitorServer struct {
	pulumirpc.UnimplementedResourceMonitorServer
	t       testing.TB
	project string
	stack   string
	mocks   pulumi.MockResourceMonitor

	mu        sync.Mutex
	resources map[string]resource.PropertyMap
}

func (m *ResourceMonitorServer) Resources() map[string]resource.PropertyMap {
	m.mu.Lock()
	defer m.mu.Unlock()
	resources := map[string]resource.PropertyMap{}
	for k, v := range m.resources {
		resources[k] = v
	}
	return resources
}

var _ pulumirpc.ResourceMonitorServer = &ResourceMonitorServer{}

func (m *ResourceMonitorServer) newURN(parent, typ, name string) string {
	parentType := tokens.Type("")
	if parentURN := resource.URN(parent); parentURN != "" && parentURN.QualifiedType() != resource.RootStackType {
		parentType = parentURN.QualifiedType()
	}

	return string(resource.NewURN(tokens.QName(m.stack), tokens.PackageName(m.project), parentType, tokens.Type(typ),
		name))
}

func (m *ResourceMonitorServer) SupportsFeature(context.Context, *pulumirpc.SupportsFeatureRequest) (*pulumirpc.SupportsFeatureResponse, error) {
	return &pulumirpc.SupportsFeatureResponse{
		HasSupport: true,
	}, nil
}

func (m *ResourceMonitorServer) RegisterResourceOutputs(context.Context, *pulumirpc.RegisterResourceOutputsRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (m *ResourceMonitorServer) RegisterResource(ctx context.Context, in *pulumirpc.RegisterResourceRequest) (*pulumirpc.RegisterResourceResponse, error) {
	if in.GetType() == string(resource.RootStackType) && in.GetParent() == "" {
		return &pulumirpc.RegisterResourceResponse{
			Urn: m.newURN(in.GetParent(), in.GetType(), in.GetName()),
		}, nil
	}

	inputs, err := plugin.UnmarshalProperties(in.GetObject(), plugin.MarshalOptions{
		KeepSecrets:   true,
		KeepResources: true,
	})
	if err != nil {
		return nil, err
	}

	id, state, err := m.mocks.NewResource(pulumi.MockResourceArgs{
		TypeToken:   in.GetType(),
		Name:        in.GetName(),
		Inputs:      inputs,
		Provider:    in.GetProvider(),
		ID:          in.GetImportId(),
		Custom:      in.GetCustom(),
		RegisterRPC: in,
	})
	if err != nil {
		return nil, err
	}

	urn := m.newURN(in.GetParent(), in.GetType(), in.GetName())

	func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.resources[urn] = resource.PropertyMap{
			resource.PropertyKey("urn"):    resource.NewStringProperty(urn),
			resource.PropertyKey("id"):     resource.NewStringProperty(id),
			resource.PropertyKey("state"):  resource.NewObjectProperty(state),
			resource.PropertyKey("parent"): resource.NewStringProperty(in.GetParent()),
		}
	}()

	stateOut, err := plugin.MarshalProperties(state, plugin.MarshalOptions{
		KeepSecrets:   true,
		KeepResources: true,
	})
	if err != nil {
		return nil, err
	}

	return &pulumirpc.RegisterResourceResponse{
		Urn:    urn,
		Id:     id,
		Object: stateOut,
	}, nil
}
