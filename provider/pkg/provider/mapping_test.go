// Copyright 2025, Pulumi Corporation.
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

func newMappingTestProvider(t *testing.T) *kubeProvider {
	t.Helper()
	k, err := makeKubeProvider(
		newMockHost(nil),
		"kubernetes",
		testPluginVersion,
		[]byte(testPulumiSchema),
		[]byte(testTerraformMapping),
		[]byte(testHelmMapping),
	)
	require.NoError(t, err)
	return k
}

func TestGetMapping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		key              string
		provider         string
		expectedProvider string
		expectedData     []byte
	}{
		{
			name:             "terraform key, no provider (legacy) returns kubernetes mapping",
			key:              "terraform",
			provider:         "",
			expectedProvider: "kubernetes",
			expectedData:     []byte(testTerraformMapping),
		},
		{
			name:             "terraform key, kubernetes provider returns kubernetes mapping",
			key:              "terraform",
			provider:         "kubernetes",
			expectedProvider: "kubernetes",
			expectedData:     []byte(testTerraformMapping),
		},
		{
			name:             "terraform key, helm provider returns helm mapping",
			key:              "terraform",
			provider:         "helm",
			expectedProvider: "helm",
			expectedData:     []byte(testHelmMapping),
		},
		{
			name:             "terraform key, unknown provider returns empty",
			key:              "terraform",
			provider:         "unknown",
			expectedProvider: "",
			expectedData:     nil,
		},
		{
			name:             "non-terraform key returns empty",
			key:              "foo",
			provider:         "",
			expectedProvider: "",
			expectedData:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			k := newMappingTestProvider(t)
			resp, err := k.GetMapping(context.Background(), &pulumirpc.GetMappingRequest{
				Key:      tt.key,
				Provider: tt.provider,
			})
			require.NoError(t, err)
			assert.Equal(t, tt.expectedProvider, resp.Provider)
			assert.Equal(t, tt.expectedData, resp.Data)
		})
	}
}

func TestGetMappings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		key               string
		expectedProviders []string
	}{
		{
			name:              "terraform key advertises kubernetes and helm",
			key:               "terraform",
			expectedProviders: []string{"kubernetes", "helm"},
		},
		{
			name:              "non-terraform key returns empty",
			key:               "foo",
			expectedProviders: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			k := newMappingTestProvider(t)
			resp, err := k.GetMappings(context.Background(), &pulumirpc.GetMappingsRequest{
				Key: tt.key,
			})
			require.NoError(t, err)
			assert.ElementsMatch(t, tt.expectedProviders, resp.Providers)
		})
	}
}
