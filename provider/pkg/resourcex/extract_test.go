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

func Test_Extract(t *testing.T) {

	res1 := resource.URN("urn:pulumi:test::test::kubernetes:core/v1:Namespace::some-namespace")

	pointer := func(i int) *int {
		return &i
	}

	type Nested struct {
		String string `json:"string"`
	}

	type Required struct {
		Number  int      `json:"number"`
		Numbers []int    `json:"numbers"`
		Struct  Nested   `json:"struct"`
		Structs []Nested `json:"structs"`
	}

	type Optional struct {
		Number  *int      `json:"number"`
		Numbers []*int    `json:"numbers"`
		Struct  *Nested   `json:"struct"`
		Structs []*Nested `json:"structs"`
	}

	tests := []struct {
		name     string
		opts     ExtractOptions
		props    resource.PropertyMap
		actual   interface{}
		expected interface{}
		result   ExtractResult
		err      error
	}{
		{
			name: "Options_RejectUnknowns",
			opts: ExtractOptions{
				RejectUnknowns: true,
			},
			props: resource.PropertyMap{
				"number": resource.NewOutputProperty(resource.Output{
					Element:      resource.NewNumberProperty(42),
					Known:        false,
					Secret:       false,
					Dependencies: []resource.URN{res1},
				}),
			},
			err:    &ContainsUnknownsError{[]resource.URN{res1}},
			actual: Required{},
		},
		{
			name: "Null_Required",
			props: resource.PropertyMap{
				"number": resource.NewNullProperty(),
			},
			expected: Required{
				Number: 0,
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false},
			actual: Required{},
		},
		{
			name: "Null_Optional",
			props: resource.PropertyMap{
				"number": resource.NewNullProperty(),
			},
			expected: Optional{
				Number: nil,
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false},
			actual: Optional{},
		},
		{
			name: "Value",
			props: resource.PropertyMap{
				"number": resource.NewNumberProperty(42),
			},
			expected: Required{
				Number: 42,
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false},
			actual: Required{},
		},
		{
			name: "Secret_Value",
			props: resource.PropertyMap{
				"number": resource.MakeSecret(resource.NewNumberProperty(42)),
			},
			expected: Required{
				Number: 42,
			},
			result: ExtractResult{ContainsSecrets: true, ContainsUnknowns: false},
			actual: Required{},
		},
		{
			name: "Secret_Byzantine",
			props: resource.PropertyMap{
				"number": resource.MakeSecret(resource.MakeComputed(resource.NewNumberProperty(42))),
			},
			expected: Required{
				Number: 0,
			},
			result: ExtractResult{ContainsSecrets: true, ContainsUnknowns: true},
			actual: Required{},
		},
		{
			name: "Computed_Required",
			props: resource.PropertyMap{
				"number": resource.MakeComputed(resource.NewNumberProperty(42)),
			},
			expected: Required{
				Number: 0,
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: true},
			actual: Required{},
		},
		{
			name: "Computed_Optional",
			props: resource.PropertyMap{
				"number": resource.MakeComputed(resource.NewNumberProperty(42)),
			},
			expected: Optional{
				Number: nil,
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: true},
			actual: Optional{},
		},
		{
			name: "Output_Unknown",
			props: resource.PropertyMap{
				"number": resource.NewOutputProperty(resource.Output{
					Element:      resource.NewNumberProperty(42),
					Known:        false,
					Secret:       false,
					Dependencies: []resource.URN{res1},
				}),
			},
			expected: Required{
				Number: 0,
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: true, Dependencies: []resource.URN{res1}},
			actual: Required{},
		},
		{
			name: "Output_Unknown_Secret",
			props: resource.PropertyMap{
				"number": resource.NewOutputProperty(resource.Output{
					Element:      resource.NewNumberProperty(42),
					Known:        false,
					Secret:       true,
					Dependencies: []resource.URN{res1},
				}),
			},
			expected: Required{
				Number: 0,
			},
			result: ExtractResult{ContainsSecrets: true, ContainsUnknowns: true, Dependencies: []resource.URN{res1}},
			actual: Required{},
		},
		{
			name: "Output_Known",
			props: resource.PropertyMap{
				"number": resource.NewOutputProperty(resource.Output{
					Element:      resource.NewNumberProperty(42),
					Known:        true,
					Secret:       false,
					Dependencies: []resource.URN{res1},
				}),
			},
			expected: Required{
				Number: 42,
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false, Dependencies: []resource.URN{res1}},
			actual: Required{},
		},
		{
			name: "Output_Known_Secret",
			props: resource.PropertyMap{
				"number": resource.NewOutputProperty(resource.Output{
					Element:      resource.NewNumberProperty(42),
					Known:        true,
					Secret:       true,
					Dependencies: []resource.URN{res1},
				}),
			},
			expected: Required{
				Number: 42,
			},
			result: ExtractResult{ContainsSecrets: true, ContainsUnknowns: false, Dependencies: []resource.URN{res1}},
			actual: Required{},
		},
		{
			name: "Output_Known_Byzantine",
			props: resource.PropertyMap{
				"number": resource.NewOutputProperty(resource.Output{
					Element:      resource.MakeSecret(resource.NewNumberProperty(42)),
					Known:        true,
					Secret:       false,
					Dependencies: []resource.URN{res1},
				}),
			},
			expected: Required{
				Number: 42,
			},
			result: ExtractResult{ContainsSecrets: true, ContainsUnknowns: false, Dependencies: []resource.URN{res1}},
			actual: Required{},
		},
		{
			name: "Array_Null",
			props: resource.PropertyMap{
				"numbers": resource.NewNullProperty(),
			},
			expected: Required{
				Numbers: nil,
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false},
			actual: Required{},
		},
		{
			name: "Array_Computed",
			props: resource.PropertyMap{
				"numbers": resource.MakeComputed(resource.NewArrayProperty([]resource.PropertyValue{})),
			},
			expected: Required{
				Numbers: nil,
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: true},
			actual: Required{},
		},
		{
			name: "Array_Secret",
			props: resource.PropertyMap{
				"numbers": resource.MakeSecret(resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewNumberProperty(42),
				})),
			},
			expected: Required{
				Numbers: []int{42},
			},
			result: ExtractResult{ContainsSecrets: true, ContainsUnknowns: false},
			actual: Required{},
		},
		{
			name: "Array_Element_Null",
			props: resource.PropertyMap{
				"numbers": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewNullProperty(),
				}),
			},
			expected: Required{
				Numbers: []int{0},
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false},
			actual: Required{},
		},
		{
			name: "Array_Element_Required",
			props: resource.PropertyMap{
				"numbers": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewNumberProperty(42),
				}),
			},
			expected: Required{
				Numbers: []int{42},
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false},
			actual: Required{},
		},
		{
			name: "Array_Element_Optional",
			props: resource.PropertyMap{
				"numbers": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewNumberProperty(42),
				}),
			},
			expected: Optional{
				Numbers: []*int{pointer(42)},
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false},
			actual: Optional{},
		},
		{
			name: "Array_Element_Computed",
			props: resource.PropertyMap{
				"numbers": resource.NewArrayProperty([]resource.PropertyValue{
					resource.MakeComputed(resource.NewNumberProperty(42)),
				}),
			},
			expected: Required{
				Numbers: []int{0},
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: true},
			actual: Required{},
		},
		{
			name: "Array_Element_Struct_Secret",
			props: resource.PropertyMap{
				"structs": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewObjectProperty(resource.PropertyMap{
						"string": resource.MakeSecret(resource.NewStringProperty("foo")),
					}),
				}),
			},
			expected: Required{
				Structs: []Nested{{String: "foo"}},
			},
			result: ExtractResult{ContainsSecrets: true, ContainsUnknowns: false},
			actual: Required{},
		},
		{
			name: "Array_Element_Struct_Computed",
			props: resource.PropertyMap{
				"structs": resource.NewArrayProperty([]resource.PropertyValue{
					resource.NewObjectProperty(resource.PropertyMap{
						"string": resource.MakeComputed(resource.NewStringProperty("foo")),
					}),
				}),
			},
			expected: Required{
				Structs: []Nested{{String: ""}},
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: true},
			actual: Required{},
		},
		{
			name: "Object_Null_Required",
			props: resource.PropertyMap{
				"struct": resource.NewNullProperty(),
			},
			expected: Required{
				Struct: Nested{},
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false},
			actual: Required{},
		},
		{
			name: "Object_Null_Optional",
			props: resource.PropertyMap{
				"struct": resource.NewNullProperty(),
			},
			expected: Optional{
				Struct: nil,
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false},
			actual: Optional{},
		},
		{
			name: "Object_Null_Required",
			props: resource.PropertyMap{
				"struct": resource.NewObjectProperty(resource.PropertyMap{}),
			},
			expected: Required{
				Struct: Nested{},
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false},
			actual: Required{},
		},
		{
			name: "Object_Computed",
			props: resource.PropertyMap{
				"struct": resource.MakeComputed(resource.NewObjectProperty(resource.PropertyMap{})),
			},
			expected: Required{
				Struct: Nested{},
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: true},
			actual: Required{},
		},
		{
			name: "Object_Secret",
			props: resource.PropertyMap{
				"struct": resource.MakeSecret(resource.NewObjectProperty(resource.PropertyMap{})),
			},
			expected: Required{
				Struct: Nested{},
			},
			result: ExtractResult{ContainsSecrets: true, ContainsUnknowns: false},
			actual: Required{},
		},
		{
			name: "Object_Element",
			props: resource.PropertyMap{
				"struct": resource.NewObjectProperty(resource.PropertyMap{
					"string": resource.NewStringProperty("foo"),
				}),
			},
			expected: Required{
				Struct: Nested{String: "foo"},
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false},
			actual: Required{},
		},
		{
			name: "Object_Element_Ignored",
			props: resource.PropertyMap{
				"struct": resource.NewObjectProperty(resource.PropertyMap{
					"string":  resource.NewStringProperty("foo"),
					"ignored": resource.MakeComputed(resource.NewStringProperty("")),
				}),
			},
			expected: Required{
				Struct: Nested{String: "foo"},
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false},
			actual: Required{},
		},
		{
			name: "Ignored_Computed",
			props: resource.PropertyMap{
				"number":  resource.NewNumberProperty(42),
				"ignored": resource.MakeComputed(resource.NewStringProperty("")),
			},
			expected: Required{
				Number: 42,
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false},
			actual: Required{},
		},
		{
			name: "Ignored_Output",
			props: resource.PropertyMap{
				"number": resource.NewNumberProperty(42),
				"ignored": resource.NewOutputProperty(resource.Output{
					Element:      resource.NewStringProperty("ignored"),
					Known:        false,
					Secret:       false,
					Dependencies: []resource.URN{res1},
				}),
			},
			expected: Required{
				Number: 42,
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false},
			actual: Required{},
		},
		{
			name: "Ignored_Secret",
			props: resource.PropertyMap{
				"number":  resource.NewNumberProperty(42),
				"ignored": resource.MakeSecret(resource.NewStringProperty("foo")),
			},
			expected: Required{
				Number: 42,
			},
			result: ExtractResult{ContainsSecrets: false, ContainsUnknowns: false},
			actual: Required{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result, err := Extract(&tt.actual, tt.props, tt.opts)
			if tt.err != nil {
				require.Equal(t, tt.err, err, "expected error")
				return
			}
			require.NoError(t, err, "expected no error")
			require.Equal(t, tt.result, result, "expected result")
			require.Equal(t, tt.expected, tt.actual, "expected target")
		})
	}
}
