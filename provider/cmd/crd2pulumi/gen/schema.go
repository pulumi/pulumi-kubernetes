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

package gen

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/v2/codegen"
	pschema "github.com/pulumi/pulumi/pkg/v2/codegen/schema"
	unstruct "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const tool = "crd2pulumi"

// hyphenedFields contains a set of the hyphened Kubernetes fields
var hyphenedFields = codegen.NewStringSet(
	"x-kubernetes-embedded-resource",
	"x-kubernetes-int-or-string",
	"x-kubernetes-list-map-keys",
	"x-kubernetes-list-type",
	"x-kubernetes-map-type",
	"x-kubernetes-preserve-unknown-fields",
)

// anyTypeRef specifies the designated reference string for an "any" type
const anyTypeRef = "pulumi.json#/Any"

var anyTypeSpec = pschema.TypeSpec{
	Ref: anyTypeRef,
}

var arbitraryJSONTypeSpec = pschema.TypeSpec{
	Type:                 "object",
	AdditionalProperties: &anyTypeSpec,
}

// ObjectMeta type
const objectMetaRef = "#/types/kubernetes:meta/v1:ObjectMeta"

// Union type of integer and string
var intOrStringTypeSpec = pschema.TypeSpec{
	OneOf: []pschema.TypeSpec{
		pschema.TypeSpec{
			Type: "integer",
		},
		pschema.TypeSpec{
			Type: "string",
		},
	},
}

var apiVersionPropertySpec = pschema.PropertySpec{
	TypeSpec: pschema.TypeSpec{
		Type: "string",
	},
	Description: `APIVersion defines the versioned schema of this representation
	of an object. Servers should convert recognized schemas to the latest
	internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources`,
}

var kindPropertySpec = pschema.PropertySpec{
	TypeSpec: pschema.TypeSpec{
		Type: "string",
	},
	Description: `Kind is a string value representing the REST resource this
	object represents. Servers may infer this from the endpoint the client
	submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds`,
}

// Returns true if the given TypeSpec is of type any; returns false otherwise
func isAnyType(typeSpec pschema.TypeSpec) bool {
	return typeSpec.Ref == anyTypeRef
}

// GetObjectTypeSpecs generates types for each versioned schema. Returns a
// mapping from each type's name in the format of
// "{.spec.group}/{.spec.versions[*].name}:{.spec.names.kind}" to its proper
// pschema.ObjectTypeSpec.
func (gen *CustomResourceGenerator) GetObjectTypeSpecs() map[string]pschema.ObjectTypeSpec {
	objectTypeSpecs := map[string]pschema.ObjectTypeSpec{}
	for version, schema := range gen.Schemas {
		baseRef := getBaseRef(gen.Group, version, gen.Kind)
		addType(schema, baseRef, objectTypeSpecs)
	}
	return objectTypeSpecs
}

func (gen *CustomResourceGenerator) baseRefs() []string {
	versions := gen.Versions()
	for i, version := range versions {
		versions[i] = getBaseRef(gen.Group, version, gen.Kind)
	}
	return versions
}

func getBaseRef(group, version, kind string) string {
	return fmt.Sprintf("kubernetes:%s/%s:%s", group, version, kind)
}

// AddArgsSuffix appends 'Args' to the name of versioned CustomResourceArgs
// interface. For example, this would generate "CronTabArgs" instead of
// "CronTab" for NodeJS.
func AddArgsSuffix(objectTypeSpecs map[string]pschema.ObjectTypeSpec, baseRefs []string) {
	for _, baseRef := range baseRefs {
		newBaseRef := baseRef + "Args"
		objectTypeSpecs[newBaseRef] = objectTypeSpecs[baseRef]
		delete(objectTypeSpecs, baseRef)
	}
}

// AddMetadataRefs adds a `metadata?: meta/v1/ObjectMeta` field to each
// versioned CustomResourceArgs interface.
func AddMetadataRefs(objectTypeSpecs map[string]pschema.ObjectTypeSpec, baseRefs []string) {
	for _, baseRef := range baseRefs {
		objectTypeSpecs[baseRef].Properties["metadata"] = pschema.PropertySpec{
			TypeSpec: pschema.TypeSpec{
				Ref: objectMetaRef,
			},
		}
	}
}

// AddAPIVersionAndKindProperties adds the `apiVersion` and `kind` properties to
// every version's schema, if it doesn't exist already.
func (gen *CustomResourceGenerator) AddAPIVersionAndKindProperties(objectTypeSpecs map[string]pschema.ObjectTypeSpec, baseRefs []string) {
	for _, version := range gen.Versions() {
		baseRef := getBaseRef(gen.Group, version, gen.Kind)
		objectTypeSpecs[baseRef].Properties["apiVersion"] = pschema.PropertySpec{
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Const: gen.Group + "/" + version,
		}
		objectTypeSpecs[baseRef].Properties["kind"] = pschema.PropertySpec{
			TypeSpec: pschema.TypeSpec{
				Type: "string",
			},
			Const: gen.Kind,
		}
	}
}

// AddPlaceholderMetadataSpec adds a placeholder `kubernetes:meta/v1:ObjectMeta`
// objectTypeSpec value. This is needed so that the Go codegen can properly
// import meta/v1:ObjectMeta.
func AddPlaceholderMetadataSpec(objectTypeSpecs map[string]pschema.ObjectTypeSpec) {
	objectTypeSpecs["kubernetes:meta/v1:ObjectMeta"] = pschema.ObjectTypeSpec{
		Type: "object",
	}
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

// UnderscoreFields replaces the hyphens in the x-kubernetes-... fields with
// underscores
func UnderscoreFields(schema map[string]interface{}) {
	for field, val := range schema {
		if hyphenedFields.Has(field) {
			delete(schema, field)
			underScoredField := strings.ReplaceAll(field, "-", "_")
			schema[underScoredField] = val
		}
		if subSchema, ok := val.(map[string]interface{}); ok {
			UnderscoreFields(subSchema)
		} else if subSchemaSlice, ok := val.([]interface{}); ok {
			for _, genericSubSchema := range subSchemaSlice {
				if subSchema, ok = genericSubSchema.(map[string]interface{}); ok {
					UnderscoreFields(subSchema)
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
		return anyTypeSpec
	}

	// If the schema is of the `oneOf` type: return a TypeSpec with the `OneOf`
	// field filled with the TypeSpec of all sub-schemas.
	oneOf, foundOneOf, _ := NestedMapSlice(schema, "oneOf")
	if foundOneOf {
		oneOfTypeSpecs := make([]pschema.TypeSpec, 0, len(oneOf))
		for i, oneOfSchema := range oneOf {
			oneOfTypeSpec := getTypeSpec(oneOfSchema, name+"OneOf"+strconv.Itoa(i), types)
			if isAnyType(oneOfTypeSpec) {
				return anyTypeSpec
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
		return intOrStringTypeSpec
	}

	preserveUnknownFields, foundPreserveUnknownFields, _ := unstruct.NestedBool(schema, "x_kubernetes_preserve_unknown_fields")
	if foundPreserveUnknownFields && preserveUnknownFields {
		return arbitraryJSONTypeSpec
	}

	// If the the schema wasn't some combination of other types (`oneOf`,
	// `allOf`, `anyOf`), then it must have a "type" field, otherwise we
	// cannot represent it. If we cannot represent it, we simply set it to be
	// any type.
	schemaType, foundSchemaType, _ := unstruct.NestedString(schema, "type")
	if !foundSchemaType {
		return anyTypeSpec
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
				AdditionalProperties: &anyTypeSpec,
			}
		}
		// If no properties are found, then it can be arbitrary JSON
		_, foundProperties, _ := unstruct.NestedMap(schema, "properties")
		if !foundProperties {
			return arbitraryJSONTypeSpec
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
		return anyTypeSpec
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

func genPackage(objectTypeSpecs map[string]pschema.ObjectTypeSpec, baseRefs []string, language string) (*pschema.Package, error) {
	resources := map[string]pschema.ResourceSpec{}
	for _, baseRef := range baseRefs {
		objectTypeSpec := objectTypeSpecs[baseRef]
		resources[baseRef] = pschema.ResourceSpec{
			ObjectTypeSpec:  objectTypeSpec,
			InputProperties: objectTypeSpec.Properties,
		}
	}

	pkgSpec := pschema.PackageSpec{
		Version:   "2.0.0",
		Types:     objectTypeSpecs,
		Resources: resources,
		Language: map[string]json.RawMessage{
			language: []byte("{}"),
		},
	}

	pkg, err := pschema.ImportSpec(pkgSpec, nil)
	if err != nil {
		return &pschema.Package{}, errors.Wrapf(err, "importing spec")
	}

	return pkg, nil
}
