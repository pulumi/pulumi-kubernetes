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
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/provider/v2/cmd/crd2pulumi/gen"
	pschema "github.com/pulumi/pulumi/pkg/v2/codegen/schema"
	"github.com/stretchr/testify/assert"
)

const TestUnderscoreFieldsYAML = "test-underscorefields.yaml"
const TestCombineSchemasYAML = "test-combineschemas.yaml"
const TestGetTypeSpecYAML = "test-gettypespec.yaml"
const TestGetTypeSpecJSON = "test-gettypespec.json"

func UnmarshalSchemas(yamlPath string) (map[string]interface{}, error) {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read file %s", yamlPath)
	}
	schema, err := gen.UnmarshalYaml(yamlFile)
	if err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal %s", yamlPath)
	}
	return schema, nil
}

func UnmarshalTypeSpecJSON(jsonPath string) (map[string]pschema.TypeSpec, error) {
	jsonFile, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read file %s", jsonPath)
	}
	var v map[string]pschema.TypeSpec
	err = json.Unmarshal(jsonFile, &v)
	if err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal %s", jsonPath)
	}
	return v, nil
}

func TestUnderscoreFields(t *testing.T) {
	schemas, err := UnmarshalSchemas(TestUnderscoreFieldsYAML)
	assert.NoError(t, err)

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
	schemas, err := UnmarshalSchemas(TestCombineSchemasYAML)
	assert.NoError(t, err)
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

func TestGetTypeSpec(t *testing.T) {
	// gen.GetTypeSpec wants us to pass in a types map
	// (map[string]pschema.ObjectTypeSpec{}) to add object refs when we see
	// them. However we only want the returned pschema.TypeSpec, so this
	// wrapper function creates a placeholder types map and just returns
	// the pschema.TypeSpec. Since our initial name arg is "", this causes all
	// objects to have the ref "#/types/"
	getOnlyTypeSpec := func(schema map[string]interface{}) pschema.TypeSpec {
		placeholderTypes := map[string]pschema.ObjectTypeSpec{}
		return gen.GetTypeSpec(schema, "", placeholderTypes)
	}

	// Load YAML schemas
	schemas, err := UnmarshalSchemas(TestGetTypeSpecYAML)
	assert.NoError(t, err)

	// Load expected TypeSpec outputs as JSON
	typeSpecs, err := UnmarshalTypeSpecJSON(TestGetTypeSpecJSON)
	assert.NoError(t, err)

	for name := range schemas {
		expected, ok := typeSpecs[name]
		assert.True(t, ok)

		schema := schemas[name].(map[string]interface{})
		actual := getOnlyTypeSpec(schema)

		assert.EqualValues(t, expected, actual)
	}
}
