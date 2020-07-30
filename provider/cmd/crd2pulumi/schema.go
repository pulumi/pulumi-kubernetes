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

var HyphenedFields = [...]string{
	"x-kubernetes-embedded-resource",
	"x-kubernetes-int-or-string",
	"x-kubernetes-list-map-keys",
	"x-kubernetes-list-type",
	"x-kubernetes-map-type",
	"x-kubernetes-preserve-unknown-fields",
}

const AnyTypeRef = "pulumi.json#/Any"

var AnyTypeSpec = pschema.TypeSpec{
	Ref: AnyTypeRef,
}

var ArbitraryJSONTypeSpec = pschema.TypeSpec{
	Type:                 "object",
	AdditionalProperties: &AnyTypeSpec,
}

// ObjectMeta type
const ObjectMetaRef = "#/types/kubernetes:meta/v1:ObjectMeta"

// Union type of integer and string
var IntOrStringTypeSpec = pschema.TypeSpec{
	OneOf: []pschema.TypeSpec{
		pschema.TypeSpec{
			Type: "integer",
		},
		pschema.TypeSpec{
			Type: "string",
		},
	},
}

// Returns true if the given TypeSpec is of type any; returns false otherwise
func isAnyType(typeSpec pschema.TypeSpec) bool {
	return typeSpec.Ref == AnyTypeRef
}

// GetObjectTypeSpecs generates types for each versioned schema. Returns a
// mapping from each type's name in the format of
// "{.spec.group}/{.spec.versions[*].name}:{.spec.names.kind}" to its proper
// pschema.ObjectTypeSpec.
func GetObjectTypeSpecs(versions map[string]map[string]interface{}, name, kind string) map[string]pschema.ObjectTypeSpec {
	objectTypeSpecs := map[string]pschema.ObjectTypeSpec{}
	for version, schema := range versions {
		baseRef := fmt.Sprintf("crds:%s/%s:%s", name, version, kind)
		addType(schema, baseRef, objectTypeSpecs)
		// Adds "Args" to the baseRef name
		newBaseRef := baseRef + "Args"
		objectTypeSpecs[newBaseRef] = objectTypeSpecs[baseRef]
		objectTypeSpecs[newBaseRef].Properties["metadata"] = pschema.PropertySpec{
			TypeSpec: pschema.TypeSpec{
				Ref: ObjectMetaRef,
			},
		}
		delete(objectTypeSpecs, baseRef)
	}

	return objectTypeSpecs
}

// addType converts the given OpenAPI `schema` to a ObjectTypeSpec and adds it
// to the `types` map under the given `name`. Recursively converts and adds all
// nested schemas as well.
func addType(schema map[string]interface{}, name string, types map[string]pschema.ObjectTypeSpec) {
	properties, foundProperties, _ := unstruct.NestedMap(schema, "properties")
	if !foundProperties {
		return
	}

	description, _, _ := unstruct.NestedString(schema, "description")
	schemaType, _, _ := unstruct.NestedString(schema, "type")
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

	types[name] = pschema.ObjectTypeSpec{
		Type:        schemaType,
		Properties:  propertySpecs,
		Required:    required,
		Description: description,
	}
}

// Replaces the hyphens in the x-kubernetes-... fields with underscores
func underscoreFields(schema map[string]interface{}) {
	for field, val := range schema {
		for _, hyphenedField := range HyphenedFields {
			if field == hyphenedField {
				delete(schema, field)
				underScoredField := strings.ReplaceAll(field, "-", "_")
				schema[underScoredField] = val
			}
		}
		subSchema, ok := val.(map[string]interface{})
		if ok {
			underscoreFields(subSchema)
		} else {
			subSchemaSlice, ok := val.([]interface{})
			if ok {
				for _, genericSubSchema := range subSchemaSlice {
					subSchema, ok = genericSubSchema.(map[string]interface{})
					if ok {
						underscoreFields(subSchema)
					} else {
						break
					}
				}
			}
		}
	}
}

// getTypeSpec returns the corresponding pschema.TypeSpec for a OpenAPI v3
// schema. Handles nested pschema.TypeSpecs in case the schema type is an array,
// object, or "combined schema" (oneOf, allOf, anyOf). Also recursively converts
// and adds all schemas of type object to the types map.
func getTypeSpec(schema map[string]interface{}, name string, types map[string]pschema.ObjectTypeSpec) pschema.TypeSpec {
	if schema == nil {
		return AnyTypeSpec
	}

	// If the schema is of the `oneOf` type: return a TypeSpec with the `OneOf`
	// field filled with the TypeSpec of all sub-schemas.
	oneOf, foundOneOf, _ := NestedMapSlice(schema, "oneOf")
	if foundOneOf {
		oneOfTypeSpecs := make([]pschema.TypeSpec, 0, len(oneOf))
		for i, oneOfSchema := range oneOf {
			oneOfTypeSpec := getTypeSpec(oneOfSchema, name+"OneOf"+strconv.Itoa(i), types)
			if isAnyType(oneOfTypeSpec) {
				return AnyTypeSpec
			}
			oneOfTypeSpecs = append(oneOfTypeSpecs, oneOfTypeSpec)
		}
		return pschema.TypeSpec{
			OneOf: oneOfTypeSpecs,
		}
	}

	// If the schema is of `allOf` type: combine `properties` and `required`
	// fields of sub-schemas into a single schema. Then return the `TypeSpec`
	// of that combined schema.
	allOf, foundAllOf, _ := NestedMapSlice(schema, "allOf")
	if foundAllOf {
		combinedSchema := CombineSchemas(true, allOf...)
		return getTypeSpec(combinedSchema, name, types)
	}

	// If the schema is of `anyOf` type: combine only `properties` of
	// sub-schemas into a single schema, with all `properties` set to optional.
	// Then return the `TypeSpec` of that combined schema.
	anyOf, foundAnyOf, _ := NestedMapSlice(schema, "anyOf")
	if foundAnyOf {
		combinedSchema := CombineSchemas(false, anyOf...)
		return getTypeSpec(combinedSchema, name, types)
	}

	intOrString, foundIntOrString, _ := unstruct.NestedBool(schema, "x_kubernetes_int_or_string")
	if foundIntOrString && intOrString {
		return IntOrStringTypeSpec
	}

	preserveUnknownFields, foundPreserveUnknownFields, _ := unstruct.NestedBool(schema, "x_kubernetes_preserve_unknown_fields")
	if foundPreserveUnknownFields && preserveUnknownFields {
		return ArbitraryJSONTypeSpec
	}

	// If the the schema wasn't some combination of other types (`oneOf`,
	// `allOf`, `anyOf`), then it must have a "type" field, otherwise we
	// cannot represent it. If we cannot represent it, we simply set it to be
	// any type.
	schemaType, foundSchemaType, _ := unstruct.NestedString(schema, "type")
	if !foundSchemaType {
		return AnyTypeSpec
	}

	switch schemaType {
	case "array":
		items, _, _ := unstruct.NestedMap(schema, "items")
		arrayTypeSpec := getTypeSpec(items, name, types)
		return pschema.TypeSpec{
			Type:  "array",
			Items: &arrayTypeSpec,
		}
	case "object":
		addType(schema, name, types)
		// If `additionalProperties` has a sub-schema, then we generate a type for a map from string --> sub-schema type
		additionalProperties, foundAdditionalProperties, _ := unstruct.NestedMap(schema, "additionalProperties")
		if foundAdditionalProperties {
			additionalPropertiesTypeSpec := getTypeSpec(additionalProperties, name, types)
			return pschema.TypeSpec{
				Type:                 "object",
				AdditionalProperties: &additionalPropertiesTypeSpec,
			}
		}
		// `additionalProperties: true` is equivalent to `additionalProperties: {}`, meaning a map from string -> any
		additionalPropertiesIsTrue, additionalPropertiesIsTrueFound, _ := unstruct.NestedBool(schema, "additionalProperties")
		if additionalPropertiesIsTrueFound && additionalPropertiesIsTrue {
			return pschema.TypeSpec{
				Type:                 "object",
				AdditionalProperties: &AnyTypeSpec,
			}
		}
		// If no properties are found, then it can be arbitrary JSON
		_, foundProperties, _ := unstruct.NestedMap(schema, "properties")
		if !foundProperties {
			return ArbitraryJSONTypeSpec
		}
		// If properties are found, then we must specify those in a seperate interface
		return pschema.TypeSpec{
			Ref: "#/types/" + name,
		}
	case "integer":
		fallthrough
	case "boolean":
		fallthrough
	case "string":
		fallthrough
	case "number":
		return pschema.TypeSpec{
			Type: schemaType,
		}
	default:
		return AnyTypeSpec
	}
}

// CombineSchemas combines the `properties` fields of the given sub-schemas into
// a single schema. Returns nil if no schemas are given. Returns the schema if
// only 1 schema is given. If combineRequired == true, then each sub-schema's
// `required` fields are also combined. In this case the combined schema's
// `required` field is of type []interface{}, not []string.
func CombineSchemas(combineRequired bool, schemas ...map[string]interface{}) map[string]interface{} {
	if len(schemas) == 0 {
		return nil
	}
	if len(schemas) == 1 {
		return schemas[0]
	}

	combinedProperties := map[string]interface{}{}
	combinedRequired := make([]string, 0)

	for _, schema := range schemas {
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
	}

	combinedSchema := map[string]interface{}{
		"type":       "object",
		"properties": combinedProperties,
	}
	if combineRequired {
		combinedSchema["required"] = GenericizeStringSlice(combinedRequired)
	}
	return combinedSchema
}
