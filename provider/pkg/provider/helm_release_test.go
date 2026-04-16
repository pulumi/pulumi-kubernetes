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

func Test_MergeMaps(t *testing.T) {
	m := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"d": []any{
					"1", "2",
				},
			},
		},
	}

	override := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"d": []any{
					"3", "4",
				},
			},
		},
	}

	for _, test := range []struct {
		name     string
		allowNil bool
		dest     map[string]any
		src      map[string]any
		expected map[string]any
	}{
		{
			name:     "Precedence",
			allowNil: false,
			dest:     m,
			src:      override,
			expected: override, // Expect the override to take precedence
		},
		{
			name:     "Merge maps",
			allowNil: false,
			dest:     m,
			src: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": []any{
							"3", "4",
						},
						"f": true,
					},
				},
			},
			expected: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"d": []any{
							"1", "2",
						},
						"c": []any{
							"3", "4",
						},
						"f": true,
					},
				},
			},
		},
		{
			name:     "Dest Has Nil Values- disallow nil",
			allowNil: false,
			dest:     m,
			src: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": any(nil),
						"e": (*any)(nil),
					},
				},
			},
			expected: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"d": []any{
							"1", "2",
						},
					},
				},
			},
		},
		{
			name:     "Dest Has Nil Values- allow nil",
			allowNil: true,
			dest:     m,
			src: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": any(nil),
						"e": (*any)(nil),
					},
				},
			},
			expected: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": any(nil),
						"e": (*any)(nil),
						"d": []any{
							"1", "2",
						},
					},
				},
			},
		},
		{
			name:     "allow nil can clear values",
			allowNil: true,
			dest: map[string]any{
				"string": "foo",
			},
			src: map[string]any{
				"string": nil,
			},
			expected: map[string]any{
				"string": nil,
			},
		},
		{
			name:     "allow nil can clear lists",
			allowNil: true,
			dest: map[string]any{
				"list": []any{1, 2, 3},
			},
			src: map[string]any{
				"list": []any{},
			},
			expected: map[string]any{
				"list": []any{},
			},
		},
		{
			name:     "allow nil doesn't clear objects",
			allowNil: true,
			dest: map[string]any{
				"object": map[string]any{"foo": "bar"},
			},
			src: map[string]any{
				"object": map[string]any{},
			},
			expected: map[string]any{
				"object": map[string]any{"foo": "bar"},
			},
		},
		{
			name:     "allow nil can clear keys",
			allowNil: true,
			dest: map[string]any{
				"livenessProbe": map[string]any{
					"httpGet": map[string]any{
						"path": "/user/login",
						"port": "http",
					},
				},
			},
			src: map[string]any{
				"livenessProbe": map[string]any{
					"httpGet": nil,
					"exec": map[string]any{
						"command": []any{"foo"},
					},
				},
			},
			expected: map[string]any{
				"livenessProbe": map[string]any{
					"httpGet": nil,
					"exec": map[string]any{
						"command": []any{"foo"},
					},
				},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			merged, err := mergeMaps(test.dest, test.src, test.allowNil)
			require.NoError(t, err)
			assert.Equal(t, test.expected, merged)
		})
	}
}

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
			name: "valueYamlFiles with literals and allowNullValues=true",
			given: resource.PropertyMap{
				"allowNullValues": resource.NewBoolProperty(true),
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
				AllowNullValues: true,
				Values: map[string]any{
					"image": map[string]any{
						"repository": "bitnami/nginx",
						"tag":        nil,
					},
				},
			},
		},
		{
			name: "valueYamlFiles with array literal and allowNullValues=true",
			given: resource.PropertyMap{
				"allowNullValues": resource.NewBoolProperty(true),
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
				AllowNullValues: true,
				Values: map[string]any{
					"images": []any{},
				},
			},
		},
		{
			name: "valueYamlFiles with string literal and allowNullValues=false",
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
						"tag":        "latest", // Not removed.
					},
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

// TestDecodeRelease_AllowNullValuesThroughWire reproduces the bug reported
// in pulumi/pulumi-kubernetes#2034, #3178, and #3234: users set
// `allowNullValues: true` expecting the null in `values.image.tag` to be
// preserved through Helm's value merge. Instead, `SkipNulls: true` on
// UnmarshalProperties strips the null before AllowNullValues gets a chance
// to preserve it.
func TestDecodeRelease_AllowNullValuesThroughWire(t *testing.T) {
	release := decodeReleaseThroughWire(t, resource.PropertyMap{
		"allowNullValues": resource.NewBoolProperty(true),
		"chart":           resource.NewStringProperty("nginx"),
		"values": resource.NewObjectProperty(resource.PropertyMap{
			"image": resource.NewObjectProperty(resource.PropertyMap{
				"tag": resource.NewNullProperty(),
			}),
		}),
	})

	require.True(t, release.AllowNullValues, "allowNullValues flag should be decoded as true")
	require.Contains(t, release.Values, "image", "image key should be present in values")

	image, ok := release.Values["image"].(map[string]any)
	require.True(t, ok, "image should be a map")

	tag, hasTag := image["tag"]
	assert.True(t, hasTag, "tag key should be present in image values (got %v)", image)
	assert.Nil(t, tag, "tag value should be nil (the user's explicit null)")
}

// TestDecodeRelease_NullWithoutAllowNullValuesThroughWire covers the common
// case: a user sets a Helm chart value to null WITHOUT setting
// `allowNullValues: true`. They expect the standard Helm convention to work
// (null deletes the key), same as `helm install --set key=null`.
//
// This test fails because there are two places nulls get stripped:
// SkipNulls at UnmarshalProperties, and `excludeNulls` in `mergeMaps` (which
// only runs when AllowNullValues is false). For this test to pass, both
// strippers must go: SkipNulls must be false, and the mergeMaps default must
// preserve nulls without requiring an opt-in flag.
func TestDecodeRelease_NullWithoutAllowNullValuesThroughWire(t *testing.T) {
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
