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
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/provider/cmd/crd2pulumi/gen"
	"github.com/stretchr/testify/assert"
)

const TestUnderscoreFieldsYAML = "test-underscorefields.yaml"
const TestCombineSchemasYAML = "test-combineschemas.yaml"

func UnmarshalSchemas(yamlPath string) (map[string]interface{}, error) {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, errors.Wrapf(err, "reading file %s", yamlPath)
	}
	schemasUnstruct, err := gen.UnmarshalYaml(yamlFile)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal %s", yamlPath)
	}
	return schemasUnstruct.Object, nil
}

func TestUnderscoreFields(t *testing.T) {
	schemas, _ := UnmarshalSchemas("testUnderscoreFields.yaml")

	// Test that calling underscoreFields() on each initial schema changes it
	// to become the same as the expected schema
	for name := range schemas {
		schema := schemas[name].(map[string]interface{})
		expected := schema["expected"].(map[string]interface{})
		initial := schema["initial"].(map[string]interface{})
		gen.UnderscoreFields(initial)
		assert.EqualValues(t, expected, initial)
	}
}

func TestCombineSchemas(t *testing.T) {
	// Test that CombineSchemas on no schemas returns nil
	assert.Nil(t, gen.CombineSchemas(false))
	assert.Nil(t, gen.CombineSchemas(true))

	// Unmarshal some testing schemas
	schemas, _ := UnmarshalSchemas(TestCombineSchemasYAML)
	person := schemas["person"].(map[string]interface{})
	employee := schemas["employee"].(map[string]interface{})

	// Test that CombineSchemas with 1 schema returns the same schema
	assert.Equal(t, person, gen.CombineSchemas(true, person))
	assert.Equal(t, person, gen.CombineSchemas(false, person))

	// Test CombineSchemas with 2 schemas and combineSchemas = true
	personAndEmployeeWithRequiredExpected := schemas["personAndEmployeeWithRequired"].(map[string]interface{})
	personAndEmployeeWithRequiredActual := gen.CombineSchemas(true, person, employee)
	assert.EqualValues(t, personAndEmployeeWithRequiredExpected, personAndEmployeeWithRequiredActual)

	// Test CombineSchemas with 2 schemas and combineSchemas = false
	personAndEmployeeWithoutRequiredExpected := schemas["personAndEmployeeWithoutRequired"].(map[string]interface{})
	personAndEmployeeWithoutRequiredActual := gen.CombineSchemas(false, person, employee)
	assert.EqualValues(t, personAndEmployeeWithoutRequiredExpected, personAndEmployeeWithoutRequiredActual)
}
