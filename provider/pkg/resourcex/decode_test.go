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

func Test_Decode(t *testing.T) {
	tests := []struct {
		name     string
		props    resource.PropertyMap
		expected interface{}
	}{
		{
			name: "null",
			props: resource.PropertyMap{
				"value": resource.NewNullProperty(),
			},
			expected: map[string]interface{}{
				"value": nil,
			},
		},
		{
			name: "bool",
			props: resource.PropertyMap{
				"value": resource.NewBoolProperty(true),
			},
			expected: map[string]interface{}{
				"value": true,
			},
		},
		{
			name: "number",
			props: resource.PropertyMap{
				"value": resource.NewNumberProperty(42),
			},
			expected: map[string]interface{}{
				"value": 42.,
			},
		},
		{
			name: "string",
			props: resource.PropertyMap{
				"value": resource.NewStringProperty("foo"),
			},
			expected: map[string]interface{}{
				"value": "foo",
			},
		},
		{
			name: "array_value",
			props: resource.PropertyMap{
				"value": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewStringProperty("foo"),
				}),
			},
			expected: map[string]interface{}{
				"value": []interface{}{"foo"},
			},
		},
		{
			name: "array_null",
			props: resource.PropertyMap{
				"value": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewNullProperty(),
				}),
			},
			expected: map[string]interface{}{
				"value": []interface{}{nil},
			},
		},
		{
			name: "array_secret",
			props: resource.PropertyMap{
				"value": resource.NewArrayProperty([]resource.PropertyValue{
					resource.MakeSecret(resource.NewStringProperty("foo")),
				}),
			},
			expected: map[string]interface{}{
				"value": []interface{}{"foo"},
			},
		},
		{
			name: "array_computed",
			props: resource.PropertyMap{
				"value": resource.NewArrayProperty([]resource.PropertyValue{
					resource.MakeComputed(resource.NewStringProperty("foo")),
				}),
			},
			expected: map[string]interface{}{
				"value": []interface{}{nil},
			},
		},
		{
			name: "computed",
			props: resource.PropertyMap{
				"value": resource.MakeComputed(resource.NewStringProperty("foo")),
			},
			expected: map[string]interface{}{
				"value": nil,
			},
		},
		{
			name: "output_unknown",
			props: resource.PropertyMap{
				"value": resource.NewOutputProperty(resource.Output{
					Element: resource.NewStringProperty("foo"),
					Known:   false,
				}),
			},
			expected: map[string]interface{}{
				"value": nil,
			},
		},
		{
			name: "output_known",
			props: resource.PropertyMap{
				"value": resource.NewOutputProperty(resource.Output{
					Element: resource.NewStringProperty("foo"),
					Known:   true,
				}),
			},
			expected: map[string]interface{}{
				"value": "foo",
			},
		},
		{
			name: "output_byzantine",
			props: resource.PropertyMap{
				"value": resource.NewOutputProperty(resource.Output{
					Element: resource.MakeSecret(resource.NewStringProperty("foo")),
					Known:   true,
				}),
			},
			expected: map[string]interface{}{
				"value": "foo",
			},
		},
		{
			name: "secret_value",
			props: resource.PropertyMap{
				"value": resource.MakeSecret(resource.NewStringProperty("foo")),
			},
			expected: map[string]interface{}{
				"value": "foo",
			},
		},
		{
			name: "secret_computed",
			props: resource.PropertyMap{
				"value": resource.MakeSecret(resource.MakeComputed(resource.NewStringProperty("foo"))),
			},
			expected: map[string]interface{}{
				"value": nil,
			},
		},
		{
			name: "object_value",
			props: resource.PropertyMap{
				"object": resource.NewObjectProperty(resource.PropertyMap{
					"value": resource.NewStringProperty("value"),
				}),
			},
			expected: map[string]interface{}{
				"object": map[string]interface{}{
					"value": "value",
				},
			},
		},
		{
			name: "object_computed",
			props: resource.PropertyMap{
				"object": resource.NewObjectProperty(resource.PropertyMap{
					"value": resource.MakeComputed(resource.NewStringProperty("value")),
				}),
			},
			expected: map[string]interface{}{
				"object": map[string]interface{}{
					"value": nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := DecodeValues(tt.props)
			require.Equal(t, tt.expected, actual, "expected result")
		})
	}
}
