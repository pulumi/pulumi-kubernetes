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

package v2

import (
	"context"
	"testing"

	structpb "github.com/golang/protobuf/ptypes/struct"
	gk "github.com/onsi/ginkgo/v2"
	gm "github.com/onsi/gomega"
	"google.golang.org/grpc"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	fakehost "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/host/fake"
)

func TestKustomizeV2(t *testing.T) {
	gm.RegisterFailHandler(gk.Fail)
	gk.RunSpecs(t, "provider/pkg/provider/kustomize/v2")
}

func unmarshalProperties(t testing.TB, props *structpb.Struct) resource.PropertyMap {
	pm, err := plugin.UnmarshalProperties(props, plugin.MarshalOptions{
		KeepUnknowns:     true,
		KeepSecrets:      true,
		KeepResources:    true,
		KeepOutputValues: true,
	})
	if err != nil {
		t.Fatalf("failed to unmarshal properties: %s", err)
	}
	return pm
}

type componentProviderTestContext struct {
	t           testing.TB
	engine      *fakehost.EngineServer
	engineAddr  string
	engineConn  *grpc.ClientConn
	monitor     *fakehost.ResourceMonitorServer
	monitorAddr string
}

func newTestContext(t testing.TB) *componentProviderTestContext {
	engine := fakehost.NewEngineServer(t)
	engineAddr := fakehost.StartEngineServer(t, engine)
	engineConn := fakehost.ConnectToEngine(t, engineAddr)

	monitor := fakehost.NewResourceMonitorServer(t, &fakehost.SimpleMonitor{})
	monitorAddr := fakehost.StartMonitorServer(t, monitor)

	return &componentProviderTestContext{
		t:           t,
		engine:      engine,
		engineAddr:  engineAddr,
		engineConn:  engineConn,
		monitor:     monitor,
		monitorAddr: monitorAddr,
	}
}

func (tc *componentProviderTestContext) NewConstructRequest() *pulumirpc.ConstructRequest {
	return &pulumirpc.ConstructRequest{
		Project:         "project",
		Stack:           "stack",
		MonitorEndpoint: tc.monitorAddr,
	}
}

func (tc *componentProviderTestContext) Engine() *fakehost.EngineServer {
	return tc.engine
}

func (tc *componentProviderTestContext) EngineConn() *grpc.ClientConn {
	return tc.engineConn
}

func (tc *componentProviderTestContext) Monitor() *fakehost.ResourceMonitorServer {
	return tc.monitor
}

func (tc *componentProviderTestContext) NewContext(ctx context.Context) *pulumi.Context {
	runInfo := pulumi.RunInfo{
		Project:     "project",
		Stack:       "stack",
		MonitorAddr: tc.monitorAddr,
		EngineAddr:  tc.engineAddr,
	}
	pulumiCtx, err := pulumi.NewContext(ctx, runInfo)
	if err != nil {
		tc.t.Fatalf("constructing run context: %s", err)
	}
	return pulumiCtx
}
