// Copyright 2016-2020, Pulumi Corporation.
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

package tests

import (
	"reflect"
	"testing"

	crdschema "github.com/pulumi/pulumi-kubernetes/provider/cmd/crd2pulumi/schema"
)

func TestCombineSchemas(t *testing.T) {
	// Test that supplying no schemas to CombineSchemas will return false
	nilSchemaFalse := crdschema.CombineSchemas(false)
	if nilSchemaFalse != nil {
		t.Errorf("CombineSchemas(false) = %v; want nil", nilSchemaFalse)
	}

	nilSchemaTrue := crdschema.CombineSchemas(true)
	if nilSchemaTrue != nil {
		t.Errorf("CombineSchemas(true) = %v; want nil", nilSchemaFalse)
	}

	// Create some test schemas
	person := map[string]interface{}{
		"type":        "object",
		"description": "Represents a person.",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type": "string",
			},
			"hometown": map[string]interface{}{
				"type": "string",
			},
			"age": map[string]interface{}{
				"type": "integer",
			},
		},
		"required": []interface{}{
			"name", "age",
		},
	}
	company := map[string]interface{}{
		"type":        "object",
		"description": "Represents a company.",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type": "string",
			},
			"address": map[string]interface{}{
				"type": "string",
			},
		},
	}
	employee := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"employeeID": map[string]interface{}{
				"type": "integer",
			},
			"company": company,
		},
		"required": []interface{}{
			"employeeID",
		},
	}
	schemas := []map[string]interface{}{
		person, company, employee,
	}

	// Test that combining just 1 schema will return that 1 schema
	for _, schema := range schemas {
		if combinedSchema := crdschema.CombineSchemas(true, schema); !reflect.DeepEqual(combinedSchema, schema) {
			t.Errorf("CombineSchemas(true, %v) = %v; want %v", schema, combinedSchema, schema)
		}
		if combinedSchema := crdschema.CombineSchemas(false, schema); !reflect.DeepEqual(combinedSchema, schema) {
			t.Errorf("CombineSchemas(false, %v) = %v; want %v", schema, combinedSchema, schema)
		}
	}

	// Test combining 2 schemas, set combineRequired = true
	personAndEmployeeActual := crdschema.CombineSchemas(true, person, employee)
	personAndEmployeeExpected := map[string]interface{}{
		"type":        "object",
		"description": "Combines 2 type(s): (1) Represents a person. (2) <no description found>",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type": "string",
			},
			"hometown": map[string]interface{}{
				"type": "string",
			},
			"age": map[string]interface{}{
				"type": "integer",
			},
			"employeeID": map[string]interface{}{
				"type": "integer",
			},
			"company": company,
		},
		"required": []interface{}{
			"name", "age", "employeeID",
		},
	}

	if !reflect.DeepEqual(personAndEmployeeActual, personAndEmployeeExpected) {
		t.Errorf("CombineSchemas(true, %v, %v) = %v; want %v", person, employee, personAndEmployeeActual, personAndEmployeeExpected)
	}

	// Test combining 2 schemas, set combineRequired = false
	personAndEmployeeNoRequiredActual := crdschema.CombineSchemas(false, person, employee)
	personAndEmployeeNoRequiredExpected := personAndEmployeeExpected
	delete(personAndEmployeeNoRequiredExpected, "required")
	if !reflect.DeepEqual(personAndEmployeeNoRequiredActual, personAndEmployeeNoRequiredExpected) {
		t.Errorf("CombineSchemas(false, %v, %v) = %v; want %v", person, employee, personAndEmployeeNoRequiredActual, personAndEmployeeNoRequiredExpected)
	}
}
