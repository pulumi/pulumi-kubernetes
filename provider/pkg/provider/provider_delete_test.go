// Copyright 2016-2026, Pulumi Corporation.
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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	kubetesting "k8s.io/client-go/testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/rpcutil/rpcerror"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
	fakehost "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/host/fake"
)

func TestDelete_HelmReleaseWithUnreachableCluster(t *testing.T) {
	tests := []struct {
		name              string
		deleteUnreachable bool
		wantErr           string
	}{
		{
			name:              "deleteUnreachable disabled reports how to recover",
			deleteUnreachable: false,
			wantErr:           "PULUMI_K8S_DELETE_UNREACHABLE",
		},
		{
			name:              "deleteUnreachable enabled drops the release from state",
			deleteUnreachable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &kubeProvider{
				host:                     newMockHost(&mockEngine{t: t}),
				clusterUnreachable:       true,
				clusterUnreachableReason: "test-only: unable to load kubeconfig",
				deleteUnreachable:        tt.deleteUnreachable,
			}

			req := &pulumirpc.DeleteRequest{
				Urn: "urn:pulumi:test::test::kubernetes:helm.sh/v3:Release::test",
				Id:  "default/test",
			}

			resp, err := k.Delete(context.Background(), req)

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, resp)
		})
	}
}

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

	client, _, _, transport := fake.NewSimpleDynamicClient(fake.WithObjects(replacement))
	transport.PrependReactor("delete", "configmaps", func(action kubetesting.Action) (bool, runtime.Object, error) {
		options := action.(kubetesting.DeleteAction).GetDeleteOptions()
		require.NotNil(t, options.Preconditions)
		require.NotNil(t, options.Preconditions.UID)
		require.Equal(t, checkpoint.GetUID(), *options.Preconditions.UID)
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
	state := checkpointObject(inputs, checkpoint, nil, "v1", "test-manager")
	properties, err := plugin.MarshalProperties(state, plugin.MarshalOptions{})
	require.NoError(t, err)

	_, err = k.Delete(t.Context(), &pulumirpc.DeleteRequest{
		Urn: "urn:pulumi:test::test::kubernetes:core/v1:ConfigMap::foo",
		Id:  "default/foo", Properties: properties, Timeout: 1,
	})
	require.Error(t, err)
	rpcErr, ok := rpcerror.FromError(err)
	require.True(t, ok)
	require.Len(t, rpcErr.Details(), 1)
	detail, ok := rpcErr.Details()[0].(*pulumirpc.ErrorResourceInitFailed)
	require.True(t, ok)
	actual, err := plugin.UnmarshalProperties(detail.Properties, plugin.MarshalOptions{})
	require.NoError(t, err)
	require.Equal(t, state, actual)
	_, recorded := parseCheckpointObject(actual)
	require.Equal(t, "default/foo", detail.Id)
	require.Equal(t, checkpoint.GetUID(), recorded.GetUID())
}
