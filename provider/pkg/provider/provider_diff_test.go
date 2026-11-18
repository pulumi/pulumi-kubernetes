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
	structpb "google.golang.org/protobuf/types/known/structpb"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	fakeclients "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
)

func TestDiff_UnchangedPatchResourceReportsNoChanges(t *testing.T) {
	daemonSet := func(containers ...any) map[string]any {
		return map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "DaemonSet",
			"metadata":   map[string]any{"name": "aws-node", "namespace": "kube-system"},
			"spec": map[string]any{"template": map[string]any{"spec": map[string]any{
				"containers": containers,
			}}},
		}
	}
	envVar := func(name, value string) any {
		return map[string]any{"name": name, "value": value}
	}

	// The managed container is second and its managed entry is last, which is the arrangement
	// positional matching gets wrong.
	declared := daemonSet(map[string]any{
		"name": "aws-node",
		"env":  []any{envVar("ENI_CONFIG_LABEL_DEF", "topology.kubernetes.io/zone")},
	})
	observed := daemonSet(
		map[string]any{"name": "aws-vpc-cni-init"},
		map[string]any{"name": "aws-node", "env": []any{
			envVar("ADDITIONAL_ENI_TAGS", "{}"),
			envVar("ENI_CONFIG_LABEL_DEF", "topology.kubernetes.io/zone"),
		}},
	)
	observed["__inputs"] = declared

	marshal := func(obj map[string]any) *structpb.Struct {
		pm, err := plugin.MarshalProperties(resource.NewPropertyMapFromMap(obj),
			plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: false, KeepSecrets: true})
		require.NoError(t, err)
		return pm
	}

	k, err := makeKubeProvider(nil, "kubernetes", "v0.0.0", []byte("{}"), []byte("{}"), []byte("{}"))
	require.NoError(t, err)
	cs, _, _, _ := fakeclients.NewSimpleDynamicClient()
	k.clientSet = cs

	resp, err := k.Diff(context.Background(), &pulumirpc.DiffRequest{
		Urn:  "urn:pulumi:test::test::kubernetes:apps/v1:DaemonSet::aws-node",
		Id:   "kube-system/aws-node",
		Olds: marshal(observed),
		News: marshal(declared),
	})
	require.NoError(t, err)
	assert.Equal(t, pulumirpc.DiffResponse_DIFF_NONE, resp.GetChanges(),
		"detailed diff: %v", resp.GetDetailedDiff())
}
