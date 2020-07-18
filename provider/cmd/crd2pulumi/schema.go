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

package main

import (
	"fmt"
	"strconv"
	"strings"

	pschema "github.com/pulumi/pulumi/pkg/v2/codegen/schema"
	unstruct "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// GetTypes parses all the types within the given crd (an unmarshalled
// resourcedefinition.yaml file). Returns a mapping from each type's name in the
// format "{.spec.group}/{.spec.versions[*].name}:{.spec.names.kind}" to its
// corresponding pschema.ObjectTypeSpec.
func GetTypes(crd unstruct.Unstructured) map[string]pschema.ObjectTypeSpec {
	versions, _, _ := NestedMapSlice(crd.Object, "spec", "versions")
	metadataName, _, _ := unstruct.NestedString(crd.Object, "metadata", "name")
	metadataName = cleanMetadataName(metadataName)
	kind, _, _ := unstruct.NestedString(crd.Object, "spec", "names", "kind")

	types := map[string]pschema.ObjectTypeSpec{}
	for _, version := range versions {
		schema, _, _ := unstruct.NestedMap(version, "schema", "openAPIV3Schema")
		versionName, _, _ := unstruct.NestedString(version, "name")
		baseRef := fmt.Sprintf("crds:%s/%s:%s", metadataName, versionName, kind)
		addType(schema, baseRef, types)
	}

	return types
}

// addType converts the given OpenAPI v3 `schema` to a ObjectTypeSpec and adds
// it to the `types` map under the given `name`. Recursively converts and adds
// all nested schemas as well.
func addType(schema map[string]interface{}, name string, types map[string]pschema.ObjectTypeSpec) {
	description, _, _ := unstruct.NestedString(schema, "description")
	schemaType, _, _ := unstruct.NestedString(schema, "type")
	properties, _, _ := unstruct.NestedMap(schema, "properties")
	required, _, _ := unstruct.NestedStringSlice(schema, "required")

	propertySpecs := map[string]pschema.PropertySpec{}
	for propertyName := range properties {
		propertySchema, _, _ := unstruct.NestedMap(properties, propertyName)
		propertyDescription, _, _ := unstruct.NestedString(propertySchema, "description")
		defaultValue, _, _ := unstruct.NestedFieldNoCopy(propertySchema, "default")
		propertySpecs[propertyName] = pschema.PropertySpec{
			TypeSpec:    getTypeSpec(propertySchema, name+strings.Title(propertyName), types),
			Description: propertyDescription,
			Default:     defaultValue,
		}
	}

	types[name+"Args"] = pschema.ObjectTypeSpec{
		Type:        schemaType,
		Properties:  propertySpecs,
		Required:    required,
		Description: description,
	}
}

// getTypeSpec returns the corresponding pschema.TypeSpec for a OpenAPI v3
// schema. Handles nested pschema.TypeSpecs in case the schema type is an array,
// object, or "combined schema" (oneOf, allOf, anyOf). Also recursively converts
// and adds all schemas of type object to the types map.
func getTypeSpec(schema map[string]interface{}, name string, types map[string]pschema.ObjectTypeSpec) pschema.TypeSpec {
	if schema == nil {
		return pschema.TypeSpec{}
	}

	// If the schema is of the `oneOf` type: return a TypeSpec with the `OneOf`
	// field filled with the TypeSpec of all sub-schemas.
	oneOf, foundOneOf, _ := NestedMapSlice(schema, "oneOf")
	if foundOneOf {
		oneOfTypeSpecs := make([]pschema.TypeSpec, 0, len(oneOf))
		for i, oneOfSchema := range oneOf {
			oneOfTypeSpec := getTypeSpec(oneOfSchema, name+"OneOf"+strconv.Itoa(i), types)
			if isAnyType(oneOfTypeSpec) {
				return anyTypeSpec()
			}
			oneOfTypeSpecs = append(oneOfTypeSpecs, oneOfTypeSpec)
		}
		return pschema.TypeSpec{
			OneOf: oneOfTypeSpecs,
		}
	}

	// If the schema is of `allOf` type: combine properties, required
	// properties, and descriptions of sub-schemas into a single schema. Then
	// return the `TypeSpec` of that single, combined schema.
	allOf, foundAllOf, _ := NestedMapSlice(schema, "allOf")
	if foundAllOf {
		combinedSchema := combineSchemas(true, allOf...)
		return getTypeSpec(combinedSchema, name, types)
	}

	// If the schema is of `anyOf` type: combine properties and descriptions of
	// sub-schemas into a single schema, with all properties set to optional.
	// Then return the 'TypeSpec` of that single, combined schema.
	anyOf, foundAnyOf, _ := NestedMapSlice(schema, "anyOf")
	if foundAnyOf {
		combinedSchema := combineSchemas(false, anyOf...)
		return getTypeSpec(combinedSchema, name, types)
	}

	// If the the schema wasn't some combination of other types (`oneOf`,
	// `allOf`, `anyOf`), then it must have a "type" field, otherwise we
	// cannot represent it. If we cannot represent it, we simply set it to be
	// any type.
	schemaType, foundSchemaType, _ := unstruct.NestedString(schema, "type")
	if !foundSchemaType {
		return anyTypeSpec()
	}

	if schemaType == "array" {
		items, _, _ := unstruct.NestedMap(schema, "items")
		arrayTypeSpec := getTypeSpec(items, name, types)
		return pschema.TypeSpec{
			Type:  "array",
			Items: &arrayTypeSpec,
		}
	}

	if schemaType == "object" {
		addType(schema, name, types)
		additionalProperties, _, _ := unstruct.NestedMap(schema, "additionalProperties")
		additionalPropertiesTypeSpec := getTypeSpec(additionalProperties, name, types)
		return pschema.TypeSpec{
			Ref:                  "#/types/" + name + "Args",
			AdditionalProperties: &additionalPropertiesTypeSpec,
		}
	}

	// Then the schemaType must be a primitive (int, bool, string, number)
	return pschema.TypeSpec{
		Type: schemaType,
	}
}

// cleanMetadataName replaces all comma-seperated instances of "pulumi" in the
// metadataName with "pulumicorp." A metadata name can't contain the word
// "pulumi", since otherwise the namspace would conflict with the imported
// "pulumi" module.
func cleanMetadataName(metadataName string) string {
	words := strings.Split(metadataName, ".")
	for i, word := range words {
		if word == "pulumi" {
			words[i] = "pulumicorp"
		}
	}
	return strings.Join(words, ".")
}

// Given a list of schemas, returns a single schema representing their
// intersection type. Combines each of their properties, and descriptions into a single
// schema. If combineRequired is true, then required properties are also
// combined into the returned schema; otherwise all properties in the returned
// schema are optional.
func combineSchemas(combineRequired bool, schemas ...map[string]interface{}) map[string]interface{} {
	combinedProperties := map[string]interface{}{}
	combinedRequired := make([]string, 0)
	var combinedDescription strings.Builder
	combinedDescription.WriteString(fmt.Sprintf("Combines %d types: ", len(schemas)))

	for i, schema := range schemas {
		properties, _, _ := unstruct.NestedMap(schema, "properties")
		for propertyName := range properties {
			propertySchema, _, _ := unstruct.NestedMap(properties, propertyName)
			combinedProperties[propertyName] = propertySchema
		}
		if combineRequired {
			required, foundRequired, _ := unstruct.NestedStringSlice(schema, "required")
			if foundRequired {
				combinedRequired = append(combinedRequired, required...)
			}
		}
		description, foundDescription, _ := unstruct.NestedString(schema, "description")
		if foundDescription {
			combinedDescription.WriteString(fmt.Sprintf("(%d) %s", i, description))
		}
	}

	combinedSchema := map[string]interface{}{
		"type":        "object",
		"description": combinedDescription.String(),
		"properties":  combinedProperties,
	}
	if combineRequired {
		combinedSchema["required"] = GenericizeStringSlice(combinedRequired)
	}
	return combinedSchema
}

const anyTypeRef = "pulumi.json#/Any"

// Returns the designated "any" TypeSpec
func anyTypeSpec() pschema.TypeSpec {
	return pschema.TypeSpec{
		Ref: anyTypeRef,
	}
}

// Returns true if the given TypeSpec is of type any; returns false otherwise
func isAnyType(typeSpec pschema.TypeSpec) bool {
	return typeSpec.Ref == anyTypeRef
}
