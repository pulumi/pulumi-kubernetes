// Copyright 2016-2018, Pulumi Corporation.
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
	"os"
	"testing"

	"fmt"
	"log"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	gomegatypes "github.com/onsi/gomega/types"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/kinds"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/rpcutil"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	objInputs = map[string]any{
		"foo": "bar",
		"baz": float64(1234),
		"qux": map[string]any{
			"xuq": "oof",
		},
	}
	objLive = map[string]any{
		initialAPIVersionKey: "",
		fieldManagerKey:      "",
		"oof":                "bar",
		"zab":                float64(4321),
		"xuq": map[string]any{
			"qux": "foo",
		},
	}
)

func TestParseOldCheckpointObject(t *testing.T) {
	old := resource.NewPropertyMapFromMap(map[string]any{
		"inputs": objInputs,
		"live":   objLive,
	})

	oldInputs, live := parseCheckpointObject(old)
	assert.Equal(t, objInputs, oldInputs.Object)
	assert.Equal(t, objLive, live.Object)
}

func TestParseNewCheckpointObject(t *testing.T) {
	old := resource.NewPropertyMapFromMap(objLive)
	old["__inputs"] = resource.NewObjectProperty(resource.NewPropertyMapFromMap(objInputs))

	oldInputs, live := parseCheckpointObject(old)
	assert.Equal(t, objInputs, oldInputs.Object)
	assert.Equal(t, objLive, live.Object)
}

func TestCheckpointObject(t *testing.T) {
	inputs := &unstructured.Unstructured{Object: objInputs}
	live := &unstructured.Unstructured{Object: objLive}

	obj := checkpointObject(inputs, live, nil, "", "")
	assert.NotNil(t, obj)

	oldInputs := obj["__inputs"]
	assert.Equal(t, objInputs, oldInputs.Mappable())

	delete(obj, "__inputs")
	assert.Equal(t, objLive, obj.Mappable())
}

// #2300 - Read() on top-level k8s objects of kind "secret" led to unencrypted __input
func TestCheckpointSecretObject(t *testing.T) {
	objInputSecret := map[string]any{
		"kind": "Secret",
		"data": map[string]any{
			"password": "verysecret",
		},
	}
	objSecretLive := map[string]any{
		initialAPIVersionKey: "",
		fieldManagerKey:      "",
		"kind":               "Secret",
		"data": map[string]any{
			"password": "verysecret",
		},
	}

	// Questionable but correct pinning test as of the time of writing
	assert.False(t, resource.NewPropertyMapFromMap(objInputSecret).ContainsSecrets())
	assert.False(t, resource.NewPropertyMapFromMap(objSecretLive).ContainsSecrets())

	inputs := &unstructured.Unstructured{Object: objInputSecret}
	live := &unstructured.Unstructured{Object: objSecretLive}

	obj := checkpointObject(inputs, live, nil, "", "")
	assert.NotNil(t, obj)

	oldInputs := obj["__inputs"]
	assert.True(t, oldInputs.IsObject())
	oldInputsVal := oldInputs.ObjectValue()
	assert.True(t, oldInputsVal["data"].ContainsSecrets())
}

func TestRoundtripCheckpointObject(t *testing.T) {
	old := resource.NewPropertyMapFromMap(objLive)
	old["__inputs"] = resource.NewObjectProperty(resource.NewPropertyMapFromMap(objInputs))

	oldInputs, oldLive := parseCheckpointObject(old)
	assert.Equal(t, objInputs, oldInputs.Object)
	assert.Equal(t, objLive, oldLive.Object)

	obj := checkpointObject(oldInputs, oldLive, nil, "", "")
	assert.Equal(t, old, obj)

	newInputs, newLive := parseCheckpointObject(obj)
	assert.Equal(t, oldInputs, newInputs)
	assert.Equal(t, oldLive, newLive)
}

func Test_equalNumbers(t *testing.T) {
	type args struct {
		a any
		b any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"a = b, int64", args{a: int64(1), b: int64(1)}, true},
		{"a = b, float64", args{a: float64(1), b: float64(1)}, true},
		{"a = b, int64, float64", args{a: int64(1), b: float64(1)}, true},
		{"a = b, float64, int64", args{a: float64(1), b: int64(1)}, true},
		{"a != b, int64", args{a: int64(1), b: int64(2)}, false},
		{"a != b, float64", args{a: float64(1), b: float64(2)}, false},
		{"a != b, int64, float64", args{a: int64(1), b: float64(2)}, false},
		{"a != b, float64, int64", args{a: float64(1), b: int64(2)}, false},
		{"unsupported a", args{a: "", b: int64(1)}, false},
		{"unsupported b", args{a: int64(1), b: ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := equalNumbers(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("equalNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isPatchURN(t *testing.T) {
	type args struct {
		urn resource.URN
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "patch URN",
			args: args{
				urn: resource.NewURN("test", "test", "", "kubernetes:apps/v1:DeploymentPatch", "test"),
			},
			want: true,
		},
		{
			name: "regular URN",
			args: args{
				urn: resource.NewURN("test", "test", "", "kubernetes:apps/v1:Deployment", "test"),
			},
			want: false,
		},
		{
			name: "CustomResource with Patch suffix",
			args: args{
				urn: resource.NewURN("test", "test", "", "kubernetes:kuma.io/v1alpha1:MeshProxyPatch", "test"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, kinds.IsPatchURN(tt.args.urn), "isPatchURN(%v)", tt.args.urn)
		})
	}
}

func Test_isListURN(t *testing.T) {
	type args struct {
		urn resource.URN
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "list URN",
			args: args{
				urn: resource.NewURN("test", "test", "", "kubernetes:apps/v1:DeploymentList", "test"),
			},
			want: true,
		},
		{
			name: "regular URN",
			args: args{
				urn: resource.NewURN("test", "test", "", "kubernetes:apps/v1:Deployment", "test"),
			},
			want: false,
		},
		{
			name: "CustomResource with List suffix",
			args: args{
				urn: resource.NewURN("test", "test", "", "kubernetes:example/v1alpha1:ExampleList", "test"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, kinds.IsListURN(tt.args.urn), "isListURN(%v)", tt.args.urn)
		})
	}
}

// CheckFailure matches a CheckFailure by property and reason.
func CheckFailure(prop string, reason types.GomegaMatcher) gomegatypes.GomegaMatcher {
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
	f.Close()
	DeferCleanup(func() {
		os.Remove(f.Name())
	})
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
