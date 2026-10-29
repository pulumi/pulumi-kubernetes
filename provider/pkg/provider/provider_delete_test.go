// Copyright 2026, Pulumi Corporation.
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
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	kubetesting "k8s.io/client-go/testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/rpcutil/rpcerror"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
	fakehost "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/host/fake"
)

func TestDeletePartialFailurePreservesCheckpoint(t *testing.T) {
	inputs := &unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata":   map[string]any{"name": "foo", "namespace": "default"},
		"data":       map[string]any{"key": "original"},
	}}
	checkpoint := inputs.DeepCopy()
	checkpoint.SetUID("original")
	replacement := inputs.DeepCopy()
	replacement.SetUID("replacement")
	replacement.Object["data"] = map[string]any{"key": "replacement"}

	client, _, _, transport := fake.NewSimpleDynamicClient(fake.WithObjects(checkpoint))
	transport.PrependReactor("delete", "configmaps", func(action kubetesting.Action) (bool, runtime.Object, error) {
		options := action.(kubetesting.DeleteAction).GetDeleteOptions()
		require.NotNil(t, options.Preconditions)
		require.NotNil(t, options.Preconditions.UID)
		require.Equal(t, checkpoint.GetUID(), *options.Preconditions.UID)
		require.NoError(t, transport.Tracker().Delete(action.GetResource(), action.GetNamespace(), checkpoint.GetName()))
		require.NoError(t, transport.Tracker().Create(action.GetResource(), replacement, action.GetNamespace()))
		return true, nil, nil
	})
	transport.PrependReactor("get", "configmaps", func(kubetesting.Action) (bool, runtime.Object, error) {
		return true, nil, apierrors.NewServiceUnavailable("live lookup unavailable")
	})

	k, err := makeKubeProvider(&fakehost.HostClient{}, "kubernetes", "4.0.0", nil, nil, nil)
	require.NoError(t, err)
	t.Cleanup(k.canceler.cancel)
	k.clientSet = client
	k.serverSideApplyMode = true
	k.enableSecrets = true
	state := checkpointObject(inputs, checkpoint, nil, "v1", "test-manager")
	state["__inputs"].ObjectValue()["bootstrapToken"] = resource.MakeSecret(
		resource.NewStringProperty("input-only-secret"),
	)
	properties, err := plugin.MarshalProperties(state, plugin.MarshalOptions{KeepSecrets: true})
	require.NoError(t, err)

	_, err = k.Delete(t.Context(), &pulumirpc.DeleteRequest{
		Urn:        "urn:pulumi:test::test::kubernetes:core/v1:ConfigMap::foo",
		Id:         "default/foo",
		Properties: properties,
		Timeout:    1,
	})
	require.Error(t, err)
	rpcErr, ok := rpcerror.FromError(err)
	require.True(t, ok)
	require.Len(t, rpcErr.Details(), 1)
	detail, ok := rpcErr.Details()[0].(*pulumirpc.ErrorResourceInitFailed)
	require.True(t, ok)
	require.True(t, proto.Equal(properties, detail.Properties))
	actual, err := plugin.UnmarshalProperties(detail.Properties, plugin.MarshalOptions{KeepSecrets: true})
	require.NoError(t, err)
	require.Equal(t, state, actual)
	require.True(t, actual["__inputs"].ObjectValue()["bootstrapToken"].IsSecret())
	_, recorded := parseCheckpointObject(actual)
	require.Equal(t, "default/foo", detail.Id)
	require.Equal(t, checkpoint.GetUID(), recorded.GetUID())
}
