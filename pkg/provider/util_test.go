// Copyright 2016-2018, Pulumi Corporation.
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

	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/stretchr/testify/assert"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestHasComputedValue(t *testing.T) {
	tests := []struct {
		name             string
		obj              *unstructured.Unstructured
		hasComputedValue bool
	}{
		{
			name:             "nil object does not have a computed value",
			obj:              nil,
			hasComputedValue: false,
		},
		{
			name:             "Empty object does not have a computed value",
			obj:              &unstructured.Unstructured{},
			hasComputedValue: false,
		},
		{
			name:             "Object with no computed values does not have a computed value",
			obj:              &unstructured.Unstructured{Object: map[string]interface{}{}},
			hasComputedValue: false,
		},
		{
			name: "Object with one concrete value does not have a computed value",
			obj: &unstructured.Unstructured{Object: map[string]interface{}{
				"field1": 1,
			}},
			hasComputedValue: false,
		},
		{
			name: "Object with one computed value does have a computed value",
			obj: &unstructured.Unstructured{Object: map[string]interface{}{
				"field1": 1,
				"field2": resource.Computed{},
			}},
			hasComputedValue: true,
		},
		{
			name: "Object with one nested computed value does have a computed value",
			obj: &unstructured.Unstructured{Object: map[string]interface{}{
				"field1": 1,
				"field2": map[string]interface{}{
					"field3": resource.Computed{},
				},
			}},
			hasComputedValue: true,
		},
		{
			name: "Object with nested maps and no computed values",
			obj: &unstructured.Unstructured{Object: map[string]interface{}{
				"field1": 1,
				"field2": map[string]interface{}{
					"field3": "3",
				},
			}},
			hasComputedValue: false,
		},
		{
			name: "Object with doubly nested maps and 1 computed value",
			obj: &unstructured.Unstructured{Object: map[string]interface{}{
				"field1": 1,
				"field2": map[string]interface{}{
					"field3": "3",
					"field4": map[string]interface{}{
						"field5": resource.Computed{},
					},
				},
			}},
			hasComputedValue: true,
		},
		{
			name: "Object with nested slice of map[string]interface{} has a computed value",
			obj: &unstructured.Unstructured{Object: map[string]interface{}{
				"field1": 1,
				"field2": []map[string]interface{}{
					{"field3": resource.Computed{}},
				},
			}},
			hasComputedValue: true,
		},
		{
			name: "Object with nested slice of interface{} has a computed value",
			obj: &unstructured.Unstructured{Object: map[string]interface{}{
				"field1": 1,
				"field2": []interface{}{
					resource.Computed{},
				},
			}},
			hasComputedValue: true,
		},
		{
			name: "Object with nested slice of map[string]interface{} with nested slice of interface{} has a computed value",
			obj: &unstructured.Unstructured{Object: map[string]interface{}{
				"field1": 1,
				"field2": []map[string]interface{}{
					{"field3": []interface{}{
						resource.Computed{},
					}},
				},
			}},
			hasComputedValue: true,
		},
		{
			name: "Complex nested object with computed value",
			obj: &unstructured.Unstructured{Object: map[string]interface{}{
				"field1": 1,
				"field2": []map[string]interface{}{
					{"field3": []interface{}{
						[]map[string]interface{}{
							{"field4": []interface{}{
								resource.Computed{},
							}},
						},
					}},
				},
			}},
			hasComputedValue: true,
		},
		{
			name: "Complex nested object with no computed value",
			obj: &unstructured.Unstructured{Object: map[string]interface{}{
				"field1": 1,
				"field2": []map[string]interface{}{
					{"field3": []interface{}{
						[]map[string]interface{}{
							{"field4": []interface{}{
								"field5",
							}},
						},
					}},
				},
			}},
			hasComputedValue: false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.hasComputedValue, hasComputedValue(test.obj), test.name)
	}
}

func TestFqName(t *testing.T) {
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "tests/v1alpha1",
			"kind":       "Test",
			"metadata": map[string]interface{}{
				"name": "myname",
			},
		},
	}

	if n := fqName(obj.GetNamespace(), obj.GetName()); n != "myname" {
		t.Errorf("Got %q for %v", n, obj)
	}

	obj.SetNamespace("mynamespace")
	if n := fqName(obj.GetNamespace(), obj.GetName()); n != "mynamespace/myname" {
		t.Errorf("Got %q for %v", n, obj)
	}
}
