// Copyright 2016-2024, Pulumi Corporation.
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

package resourcex

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/stretchr/testify/require"
)

func Test_Traverse(t *testing.T) {

	// constants
	a := resource.NewStringProperty("a")
	b := resource.NewStringProperty("b")
	c := resource.NewStringProperty("c")

	// helper functions
	path := func(path string) resource.PropertyPath {
		result, err := resource.ParsePropertyPathStrict(path)
		require.NoError(t, err, "property path")
		return result
	}
	computed := func(v resource.PropertyValue) resource.PropertyValue {
		return resource.MakeComputed(v)
	}
	unknown := func() resource.PropertyValue {
		return resource.NewOutputProperty(resource.Output{
			Known: false,
		})
	}
	known := func(v resource.PropertyValue) resource.PropertyValue {
		return resource.NewOutputProperty(resource.Output{
			Element: v,
			Known:   true,
		})
	}
	secret := func(v resource.PropertyValue) resource.PropertyValue {
		return resource.MakeSecret(v)
	}
	object := func(props resource.PropertyMap) resource.PropertyValue {
		return resource.NewObjectProperty(props)
	}
	array := func(vs ...resource.PropertyValue) resource.PropertyValue {
		return resource.NewArrayProperty(vs)
	}

	// sub-test cases
	tests := []struct {
		name     string
		props    resource.PropertyMap
		path     resource.PropertyPath
		expected []resource.PropertyValue
	}{
		{
			name: "object",
			path: path("a"),
			props: resource.PropertyMap{
				"a": a,
				"b": b,
			},
			expected: []resource.PropertyValue{
				a,
			},
		},
		{
			name:     "object_undefined",
			path:     path("a"),
			props:    resource.PropertyMap{},
			expected: []resource.PropertyValue{},
		},
		{
			name: "object_null",
			path: path("a"),
			props: resource.PropertyMap{
				"a": resource.NewNullProperty(),
			},
			expected: []resource.PropertyValue{
				resource.NewNullProperty(),
			},
		},
		{
			name: "object_computed",
			path: path("a"),
			props: resource.PropertyMap{
				"a": computed(a),
			},
			expected: []resource.PropertyValue{
				computed(a),
			},
		},
		{
			name: "object_unknown",
			path: path("a"),
			props: resource.PropertyMap{
				"a": unknown(),
			},
			expected: []resource.PropertyValue{
				unknown(),
			},
		},
		{
			name: "object_known",
			path: path("a"),
			props: resource.PropertyMap{
				"a": known(a),
			},
			expected: []resource.PropertyValue{
				known(a),
				a,
			},
		},
		{
			name: "object_secret",
			path: path("a"),
			props: resource.PropertyMap{
				"a": secret(a),
			},
			expected: []resource.PropertyValue{
				secret(a),
				a,
			},
		},
		{
			name: "object_object",
			path: path("a.b"),
			props: resource.PropertyMap{
				"a": object(resource.PropertyMap{
					"b": b,
				}),
			},
			expected: []resource.PropertyValue{
				object(resource.PropertyMap{
					"b": b,
				}),
				b,
			},
		},
		{
			name: "object_secret_object",
			path: path("a.b"),
			props: resource.PropertyMap{
				"a": secret(object(resource.PropertyMap{
					"b": b,
				})),
			},
			expected: []resource.PropertyValue{
				secret(object(resource.PropertyMap{
					"b": b,
				})),
				object(resource.PropertyMap{
					"b": b,
				}),
				b,
			},
		},
		{
			name: "object_computed_object",
			path: path("a.b"),
			props: resource.PropertyMap{
				"a": computed(object(resource.PropertyMap{})),
			},
			expected: []resource.PropertyValue{
				computed(object(resource.PropertyMap{})),
			},
		},
		{
			name: "object_unknown_object",
			path: path("a.b"),
			props: resource.PropertyMap{
				"a": unknown(),
			},
			expected: []resource.PropertyValue{
				unknown(),
			},
		},
		{
			name: "object_known_object",
			path: path("a.b"),
			props: resource.PropertyMap{
				"a": known(object(resource.PropertyMap{
					"b": b,
				})),
			},
			expected: []resource.PropertyValue{
				known(object(resource.PropertyMap{
					"b": b,
				})),
				object(resource.PropertyMap{
					"b": b,
				}),
				b,
			},
		},
		{
			name: "array",
			path: path("a"),
			props: resource.PropertyMap{
				"a": array(),
			},
			expected: []resource.PropertyValue{
				array(),
			},
		},
		{
			name: "array_index",
			path: path("a[1]"),
			props: resource.PropertyMap{
				"a": array(
					b,
					c,
				),
			},
			expected: []resource.PropertyValue{
				array(
					b,
					c,
				),
				c,
			},
		},
		{
			name: "array_index_wildcard",
			path: path("a[*]"),
			props: resource.PropertyMap{
				"a": array(
					b,
					c,
				),
			},
			expected: []resource.PropertyValue{
				array(
					b,
					c,
				),
				b,
				c,
			},
		},
		{
			name: "array_element_null",
			path: path("a[0]"),
			props: resource.PropertyMap{
				"a": array(
					resource.NewNullProperty(),
				),
			},
			expected: []resource.PropertyValue{
				array(
					resource.NewNullProperty(),
				),
				resource.NewNullProperty(),
			},
		},
		{
			name: "array_element_computed",
			path: path("a[0]"),
			props: resource.PropertyMap{
				"a": array(
					computed(b),
				),
			},
			expected: []resource.PropertyValue{
				array(
					computed(b),
				),
				computed(b),
			},
		},
		{
			name: "array_element_unknown",
			path: path("a[0]"),
			props: resource.PropertyMap{
				"a": array(
					unknown(),
				),
			},
			expected: []resource.PropertyValue{
				array(
					unknown(),
				),
				unknown(),
			},
		},
		{
			name: "array_element_known",
			path: path("a[0]"),
			props: resource.PropertyMap{
				"a": array(
					known(b),
				),
			},
			expected: []resource.PropertyValue{
				array(
					known(b),
				),
				known(b),
				b,
			},
		},
		{
			name: "array_element_secret",
			path: path("a[0]"),
			props: resource.PropertyMap{
				"a": array(
					secret(b),
				),
			},
			expected: []resource.PropertyValue{
				array(
					secret(b),
				),
				secret(b),
				b,
			},
		},
		{
			name: "array_object",
			path: path("a[0].b"),
			props: resource.PropertyMap{
				"a": array(
					object(resource.PropertyMap{
						"b": b,
					}),
				),
			},
			expected: []resource.PropertyValue{
				array(
					object(resource.PropertyMap{
						"b": b,
					}),
				),
				object(resource.PropertyMap{
					"b": b,
				}),
				b,
			},
		},
		{
			name: "array_secret_object",
			path: path("a[0].b"),
			props: resource.PropertyMap{
				"a": array(
					secret(object(resource.PropertyMap{
						"b": b,
					})),
				),
			},
			expected: []resource.PropertyValue{
				array(
					secret(object(resource.PropertyMap{
						"b": b,
					})),
				),
				secret(object(resource.PropertyMap{
					"b": b,
				})),
				object(resource.PropertyMap{
					"b": b,
				}),
				b,
			},
		},
		{
			name: "array_computed_object",
			path: path("a[0].b"),
			props: resource.PropertyMap{
				"a": array(
					computed(object(resource.PropertyMap{})),
				),
			},
			expected: []resource.PropertyValue{
				array(
					computed(object(resource.PropertyMap{})),
				),
				computed(object(resource.PropertyMap{})),
			},
		},
		{
			name: "array_unknown_object",
			path: path("a[0].b"),
			props: resource.PropertyMap{
				"a": array(
					unknown(),
				),
			},
			expected: []resource.PropertyValue{
				array(
					unknown(),
				),
				unknown(),
			},
		},
		{
			name: "array_known_object",
			path: path("a[0].b"),
			props: resource.PropertyMap{
				"a": array(
					known(object(resource.PropertyMap{
						"b": b,
					})),
				),
			},
			expected: []resource.PropertyValue{
				array(
					known(object(resource.PropertyMap{
						"b": b,
					})),
				),
				known(object(resource.PropertyMap{
					"b": b,
				})),
				object(resource.PropertyMap{
					"b": b,
				}),
				b,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := resource.NewObjectProperty(tt.props)
			var actual []resource.PropertyValue
			Traverse(v, tt.path, func(v resource.PropertyValue) {
				actual = append(actual, v)
			})
			require.Len(t, actual, len(tt.expected)+1, "expected number of components")
			require.Equal(t, v, actual[0], "expected property map")
			actual = actual[1:]
			require.Equal(t, tt.expected, actual, "expected path components")
		})
	}
}
