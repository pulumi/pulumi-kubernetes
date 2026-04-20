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

package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/asset"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
)


func TestDecodeRelease(t *testing.T) {
	bitnamiImage := `
image:
  repository: bitnami/nginx
  tag: latest
`

	tests := []struct {
		name  string
		given resource.PropertyMap
		want  *Release
	}{
		{
			name: "valueYamlFiles layering",
			given: resource.PropertyMap{
				"valueYamlFiles": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewAssetProperty(&asset.Asset{Text: `
image:
  repository: bitnami/nginx
`}),
					resource.NewAssetProperty(&asset.Asset{Text: `
image:
  tag: "1.25.0"
`}),
				}),
			},
			want: &Release{
				Values: map[string]any{
					"image": map[string]any{
						"tag":        "1.25.0",
						"repository": "bitnami/nginx",
					},
				},
			},
		},
		{
			name: "valueYamlFiles with literals",
			given: resource.PropertyMap{
				"values": resource.NewObjectProperty(resource.PropertyMap{
					"image": resource.NewObjectProperty(resource.PropertyMap{
						"tag": resource.NewStringProperty("patched"),
					}),
				},
				),
				"valueYamlFiles": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewAssetProperty(&asset.Asset{Text: bitnamiImage}),
				}),
			},
			want: &Release{
				Values: map[string]any{
					"image": map[string]any{
						"repository": "bitnami/nginx",
						"tag":        "patched",
					},
				},
			},
		},
		{
			name: "null in values overrides valueYamlFiles value",
			given: resource.PropertyMap{
				"values": resource.NewObjectProperty(resource.PropertyMap{
					"image": resource.NewObjectProperty(resource.PropertyMap{
						"tag": resource.NewNullProperty(),
					}),
				},
				),
				"valueYamlFiles": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewAssetProperty(&asset.Asset{Text: bitnamiImage}),
				}),
			},
			want: &Release{
				Values: map[string]any{
					"image": map[string]any{
						"repository": "bitnami/nginx",
						"tag":        nil,
					},
				},
			},
		},
		{
			name: "empty array in values overrides valueYamlFiles value",
			given: resource.PropertyMap{
				"values": resource.NewObjectProperty(resource.PropertyMap{
					"images": resource.NewArrayProperty([]resource.PropertyValue{}),
				},
				),
				"valueYamlFiles": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewAssetProperty(&asset.Asset{Text: `
images: ["bitnami/nginx"]
`}),
				}),
			},
			want: &Release{
				Values: map[string]any{
					"images": []any{},
				},
			},
		},
		{
			name: "explicit null in values overrides valueYamlFiles",
			given: resource.PropertyMap{
				"values": resource.NewObjectProperty(resource.PropertyMap{
					"image": resource.NewObjectProperty(resource.PropertyMap{
						"tag": resource.NewNullProperty(),
					}),
				},
				),
				"valueYamlFiles": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewAssetProperty(&asset.Asset{Text: bitnamiImage}),
				}),
			},
			want: &Release{
				Values: map[string]any{
					"image": map[string]any{
						"repository": "bitnami/nginx",
						"tag":        nil, // Cleared by the user's explicit null.
					},
				},
			},
		},
		{
			name: "valueYamlFiles with empty list overrides chart default",
			given: resource.PropertyMap{
				"valueYamlFiles": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewAssetProperty(&asset.Asset{Text: "items: []\n"}),
				}),
			},
			want: &Release{
				Values: map[string]any{
					"items": []any{},
				},
			},
		},
		{
			name: "valueYamlFiles provided, but is nil",
			given: resource.PropertyMap{
				"values": resource.NewObjectProperty(resource.PropertyMap{
					"image": resource.NewObjectProperty(resource.PropertyMap{
						"tag": resource.NewStringProperty("patched"),
					}),
				},
				),
				"valueYamlFiles": resource.NewPropertyValue(nil),
			},
			want: &Release{
				Values: map[string]any{
					"image": map[string]any{
						"tag": "patched",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := decodeRelease(tt.given, "")
			assert.NoError(t, err)
			assert.Equal(t, tt.want, actual)
		})
	}
}

// decodeReleaseThroughWire simulates the full path a real user request takes:
// PropertyMap -> marshal to structpb -> unmarshal with the options the
// provider's Check/Create/Update methods use -> decodeRelease.
//
// This lets tests exercise the actual behavior users see, including the
// null-stripping that happens at SkipNulls unmarshal time.
func decodeReleaseThroughWire(t *testing.T, userInputs resource.PropertyMap) *Release {
	t.Helper()

	wire, err := plugin.MarshalProperties(userInputs, plugin.MarshalOptions{
		Label:        "userInputs",
		KeepUnknowns: true,
		KeepSecrets:  true,
	})
	require.NoError(t, err)

	news, err := plugin.UnmarshalProperties(wire, plugin.MarshalOptions{
		Label:        "news",
		KeepUnknowns: true,
		SkipNulls:    false,
		KeepSecrets:  true,
	})
	require.NoError(t, err)

	release, err := decodeRelease(news, "")
	require.NoError(t, err)
	return release
}

// TestDecodeRelease_NullPreservedThroughWire verifies that a null value in
// Helm chart values survives the full wire round-trip (Marshal → Unmarshal →
// decodeRelease), matching the Helm CLI behavior of `--set key=null`.
func TestDecodeRelease_NullPreservedThroughWire(t *testing.T) {
	release := decodeReleaseThroughWire(t, resource.PropertyMap{
		// NOTE: allowNullValues is intentionally NOT set — we want the
		// default behavior to preserve nulls, same as Helm CLI.
		"chart": resource.NewStringProperty("nginx"),
		"values": resource.NewObjectProperty(resource.PropertyMap{
			"image": resource.NewObjectProperty(resource.PropertyMap{
				"tag": resource.NewNullProperty(),
			}),
		}),
	})

	require.Contains(t, release.Values, "image", "image key should be present in values")

	image, ok := release.Values["image"].(map[string]any)
	require.True(t, ok, "image should be a map")

	tag, hasTag := image["tag"]
	assert.True(t, hasTag, "tag key should be present in image values (got %v)", image)
	assert.Nil(t, tag, "tag value should be nil (the user's explicit null)")
}
