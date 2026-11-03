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

	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
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
