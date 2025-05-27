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

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/asset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
