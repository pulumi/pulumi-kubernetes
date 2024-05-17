// Copyright 2016-2022, Pulumi Corporation.
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

package helm

import (
	"bytes"
	"os"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/getter"
)

type mockGetter struct {
	data []byte
	err  error
}

var _ getter.Getter = (*mockGetter)(nil)

func (m *mockGetter) Get(url string, options ...getter.Option) (*bytes.Buffer, error) {
	return bytes.NewBuffer(m.data), m.err
}

func TestMergeValues(t *testing.T) {

	bitnamiImage := `
image:
  repository: bitnami/nginx
  tag: latest
`

	veleroConfiguration := `
configuration:
  backupStorageLocation:
    - name: default
`

	tests := []struct {
		name        string
		valuesFiles []pulumi.Asset
		values      map[string]any
		want        map[string]interface{}
	}{
		{
			name: "valueYamlFiles layering",
			valuesFiles: []pulumi.Asset{
				pulumi.NewStringAsset(`
image:
  repository: bitnami/nginx
`),
				pulumi.NewStringAsset(`
image:
  tag: "1.25.0"
`),
			},
			want: map[string]interface{}{
				"image": map[string]any{
					"tag":        "1.25.0",
					"repository": "bitnami/nginx",
				},
			},
		},
		{
			name: "valueYamlFiles with literals",
			valuesFiles: []pulumi.Asset{
				pulumi.NewStringAsset(bitnamiImage),
			},
			values: map[string]interface{}{
				"image": map[string]any{
					"tag": "patched",
				},
			},
			want: map[string]interface{}{
				"image": map[string]any{
					"repository": "bitnami/nginx",
					"tag":        "patched",
				},
			},
		},
		{
			name: "overrides: null values",
			valuesFiles: []pulumi.Asset{
				pulumi.NewStringAsset(bitnamiImage),
			},
			values: map[string]interface{}{
				"image": map[string]any{
					"tag": nil,
				},
			},
			want: map[string]interface{}{
				"image": map[string]any{
					"repository": "bitnami/nginx",
					"tag":        nil,
				},
			},
		},
		{
			name: "overrides: empty list (#2731)",
			valuesFiles: []pulumi.Asset{
				pulumi.NewStringAsset(veleroConfiguration),
			},
			values: map[string]interface{}{
				"configuration": map[string]any{
					"backupStorageLocation": []map[string]any{},
				},
			},
			want: map[string]interface{}{
				"configuration": map[string]any{
					"backupStorageLocation": []map[string]any{},
				},
			},
		},
		{
			name: "asset-type value",
			values: map[string]interface{}{
				"extra": map[string]any{
					"notes": pulumi.NewStringAsset("this is a note"),
				},
			},
			want: map[string]interface{}{
				"extra": map[string]any{
					"notes": "this is a note",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			merger := &ValueOpts{
				ValuesFiles: tt.valuesFiles,
				Values:      tt.values,
			}

			actual, err := merger.MergeValues(getter.Providers{})
			require.NoError(t, err)
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestReadAsset(t *testing.T) {

	bitnamiImage := `
image:
  repository: bitnami/nginx
  tag: latest
`
	bitnamiImageFile, err := os.CreateTemp("", "pulumi-TestReadAsset-*.yaml")
	require.NoError(t, err)
	_, _ = bitnamiImageFile.WriteString(bitnamiImage)
	_ = bitnamiImageFile.Close()
	defer os.Remove(bitnamiImageFile.Name())

	tests := []struct {
		name       string
		asset      pulumi.Asset
		mockGetter getter.Getter
		want       []byte
		wantErr    bool
	}{
		{
			name:  "string asset",
			asset: pulumi.NewStringAsset(bitnamiImage),
			want:  []byte(bitnamiImage),
		},
		{
			name:  "file asset",
			asset: pulumi.NewFileAsset(bitnamiImageFile.Name()),
			want:  []byte(bitnamiImage),
		},
		{
			name:    "file asset",
			asset:   pulumi.NewFileAsset("nosuchfile"),
			wantErr: true,
		},
		{
			name: "remote asset",
			mockGetter: &mockGetter{
				data: []byte(bitnamiImage),
			},
			asset: pulumi.NewRemoteAsset("mock://example.com/values.yaml"),
			want:  []byte(bitnamiImage),
		},
		{
			name: "remote asset (no protocol handler)",
			mockGetter: &mockGetter{
				data: []byte(bitnamiImage),
			},
			asset:   pulumi.NewRemoteAsset("invalid://example.com/values.yaml"),
			wantErr: true,
		},
		{
			name: "remote asset (remote error)",
			mockGetter: &mockGetter{
				err: assert.AnError,
			},
			asset:   pulumi.NewRemoteAsset("mock://example.com/values.yaml"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			p := getter.Provider{
				Schemes: []string{"mock"},
				New: func(options ...getter.Option) (getter.Getter, error) {
					return tt.mockGetter, nil
				},
			}

			actual, err := readAsset(getter.Providers{p}, tt.asset)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, actual)
		})
	}
}
