// Copyright 2016-2023, Pulumi Corporation.
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

package crd

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/iancoleman/strcase"
	kversion "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/version"
	pschema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/sdk/v3/go/common/slice"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	Boolean string = "boolean"
	Integer string = "integer"
	Number  string = "number"
	String  string = "string"
	Array   string = "array"
	Object  string = "object"
)

const anyTypeRef = "pulumi.json#/Any"

var anyTypeSpec = pschema.TypeSpec{
	Ref: anyTypeRef,
}
var arbitraryJSONTypeSpec = pschema.TypeSpec{
	Type:                 Object,
	AdditionalProperties: &anyTypeSpec,
}

var emptySpec = pschema.ComplexTypeSpec{
	ObjectTypeSpec: pschema.ObjectTypeSpec{
		Type:       Object,
		Properties: map[string]pschema.PropertySpec{},
	},
}

const objectMetaRef = "#/types/kubernetes:meta/v1:ObjectMeta"
const objectMetaToken = "kubernetes:meta/v1:ObjectMeta"

// Union type of integer and string
var intOrStringTypeSpec = pschema.TypeSpec{
	OneOf: []pschema.TypeSpec{
		{
			Type: Integer,
		},
		{
			Type: String,
		},
	},
}

// Returns the Pulumi package given a types map and a slice of the token types
// of every CustomResource. If includeObjectMetaType is true, then a
// ObjectMetaType type is also generated.
func genPackage(name, version string, types map[string]pschema.ComplexTypeSpec, resourceTokens []string, resourceGenerators []CustomResourceGenerator) (*pschema.Package, error) {
	if name == "kubernetes" {
		types[objectMetaToken] = pschema.ComplexTypeSpec{
			ObjectTypeSpec: pschema.ObjectTypeSpec{
				Type: "object",
			},
		}
	}

	packages := map[string]bool{}
	resources := map[string]pschema.ResourceSpec{}
	for _, baseRef := range resourceTokens {
		complexTypeSpec := types[baseRef]

		tok := tokens.Type(baseRef)
		alias := fmt.Sprintf("kubernetes:%s:%s", tok.Module().Name().String(), tok.Name().String())
		resources[baseRef] = pschema.ResourceSpec{
			ObjectTypeSpec:  complexTypeSpec.ObjectTypeSpec,
			InputProperties: complexTypeSpec.Properties,
			// Accomodate for easy transition from old crd2pulumi generated libraries
			Aliases: []pschema.AliasSpec{
				{
					Type: &alias,
				},
			},
		}
		packages[string(tokens.ModuleMember(baseRef).Package())] = true
	}

	allowedPackages := make([]string, 0, len(packages))
	for pkg := range packages {
		allowedPackages = append(allowedPackages, pkg)
	}
	sort.Strings(allowedPackages)

	crds := slice.Map[CustomResourceGenerator, unstructured.Unstructured](resourceGenerators, func(resourceGenerator CustomResourceGenerator) unstructured.Unstructured {
		return resourceGenerator.CustomResourceDefinition
	})

	pkgSpec := pschema.PackageSpec{
		Name:                name,
		Version:             version,
		Types:               types,
		Resources:           resources,
		AllowedPackageNames: allowedPackages,
		Extension: &pschema.PackageExtensionSpec{
			Name:    "kubernetes",
			Version: kversion.Version,
		},
		Parameter: crds,
	}

	pkg, err := pschema.ImportSpec(pkgSpec, nil)
	if err != nil {
		return &pschema.Package{}, fmt.Errorf("could not import spec: %w", err)
	}

	delete(types, objectMetaToken)
	// CRD's don't need a provider, they rely on the Kubernetes provider
	pkg.Provider = nil

	return pkg, nil
}

// Returns true if the given TypeSpec is of type any; returns false otherwise
func isAnyType(typeSpec pschema.TypeSpec) bool {
	return typeSpec.Ref == anyTypeRef
}

// AddType converts the given OpenAPI `schema` to a ObjectTypeSpec and adds it
// to the `types` map under the given `name`. Recursively converts and adds all
// nested schemas as well.
func AddType(schema map[string]any, name string, types map[string]pschema.ComplexTypeSpec) {
	properties, foundProperties, _ := unstructured.NestedMap(schema, "properties")
	description, _, _ := unstructured.NestedString(schema, "description")
	schemaType, _, _ := unstructured.NestedString(schema, "type")
	required, _, _ := unstructured.NestedStringSlice(schema, "required")

	propertySpecs := map[string]pschema.PropertySpec{}
	for propertyName := range properties {
		propertySchema, _, _ := unstructured.NestedMap(properties, propertyName)
		propertyDescription, _, _ := unstructured.NestedString(propertySchema, "description")
		defaultValue, _, _ := unstructured.NestedFieldNoCopy(propertySchema, "default")
		propertySpecs[propertyName] = pschema.PropertySpec{
			TypeSpec:    GetTypeSpec(propertySchema, name+strcase.ToCamel(propertyName), types),
			Description: propertyDescription,
			Default:     defaultValue,
		}
	}

	// If the type wasn't specified but we found properties, then we can infer that the type is an object
	if foundProperties && schemaType == "" {
		schemaType = Object
	}

	types[name] = pschema.ComplexTypeSpec{
		ObjectTypeSpec: pschema.ObjectTypeSpec{
			Type:        schemaType,
			Properties:  propertySpecs,
			Required:    required,
			Description: description,
		}}
}

// GetTypeSpec returns the corresponding pschema.TypeSpec for a OpenAPI v3
// schema. Handles nested pschema.TypeSpecs in case the schema type is an array,
// object, or "combined schema" (oneOf, allOf, anyOf). Also recursively converts
// and adds all schemas of type object to the types map.
func GetTypeSpec(schema map[string]any, name string, types map[string]pschema.ComplexTypeSpec) pschema.TypeSpec {
	if schema == nil {
		return anyTypeSpec
	}

	intOrString, foundIntOrString, _ := unstructured.NestedBool(schema, "x-kubernetes-int-or-string")
	if foundIntOrString && intOrString {
		return intOrStringTypeSpec
	}

	// If the schema is of the `oneOf` type: return a TypeSpec with the `OneOf`
	// field filled with the TypeSpec of all sub-schemas.
	oneOf, foundOneOf, _ := NestedMapSlice(schema, "oneOf")
	if foundOneOf {
		oneOfTypeSpecs := make([]pschema.TypeSpec, 0, len(oneOf))
		for i, oneOfSchema := range oneOf {
			oneOfTypeSpec := GetTypeSpec(oneOfSchema, name+"OneOf"+strconv.Itoa(i), types)
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
		return GetTypeSpec(combinedSchema, name, types)
	}

	// If the schema is of `anyOf` type: combine only `properties` of
	// sub-schemas into a single schema, with all `properties` set to optional.
	// Then return the `TypeSpec` of that combined schema.
	anyOf, foundAnyOf, _ := NestedMapSlice(schema, "anyOf")
	if foundAnyOf {
		combinedSchema := CombineSchemas(false, anyOf...)
		return GetTypeSpec(combinedSchema, name, types)
	}

	preserveUnknownFields, foundPreserveUnknownFields, _ := unstructured.NestedBool(schema, "x-kubernetes-preserve-unknown-fields")
	if foundPreserveUnknownFields && preserveUnknownFields {
		return arbitraryJSONTypeSpec
	}

	// If the the schema wasn't some combination of other types (`oneOf`,
	// `allOf`, `anyOf`), then it must have a "type" field, otherwise we
	// cannot represent it. If we cannot represent it, we simply set it to be
	// any type.
	schemaType, foundSchemaType, _ := unstructured.NestedString(schema, "type")
	if !foundSchemaType {
		return anyTypeSpec
	}

	switch schemaType {
	case Array:
		items, _, _ := unstructured.NestedMap(schema, "items")
		arrayTypeSpec := GetTypeSpec(items, name, types)
		return pschema.TypeSpec{
			Type:  Array,
			Items: &arrayTypeSpec,
		}
	case Object:
		AddType(schema, name, types)
		// If `additionalProperties` has a sub-schema, then we generate a type for a map from string --> sub-schema type
		additionalProperties, foundAdditionalProperties, _ := unstructured.NestedMap(schema, "additionalProperties")
		if foundAdditionalProperties {
			additionalPropertiesTypeSpec := GetTypeSpec(additionalProperties, name, types)
			return pschema.TypeSpec{
				Type:                 Object,
				AdditionalProperties: &additionalPropertiesTypeSpec,
			}
		}
		// `additionalProperties: true` is equivalent to `additionalProperties: {}`, meaning a map from string -> any
		additionalPropertiesIsTrue, additionalPropertiesIsTrueFound, _ := unstructured.NestedBool(schema, "additionalProperties")
		if additionalPropertiesIsTrueFound && additionalPropertiesIsTrue {
			return pschema.TypeSpec{
				Type:                 Object,
				AdditionalProperties: &anyTypeSpec,
			}
		}
		// If no properties are found, then it can be arbitrary JSON
		_, foundProperties, _ := unstructured.NestedMap(schema, "properties")
		if !foundProperties {
			return arbitraryJSONTypeSpec
		}
		// If properties are found, then we must specify those in a seperate interface
		return pschema.TypeSpec{
			Type: Object,
			Ref:  "#/types/" + name,
		}
	case Integer:
		fallthrough
	case Boolean:
		fallthrough
	case String:
		fallthrough
	case Number:
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
// `required` field is of type []any, not []string.
func CombineSchemas(combineRequired bool, schemas ...map[string]any) map[string]any {
	if len(schemas) == 0 {
		return nil
	}
	if len(schemas) == 1 {
		return schemas[0]
	}

	combinedProperties := map[string]any{}
	combinedRequired := make([]string, 0)

	for _, schema := range schemas {
		properties, _, _ := unstructured.NestedMap(schema, "properties")
		for propertyName := range properties {
			propertySchema, _, _ := unstructured.NestedMap(properties, propertyName)
			combinedProperties[propertyName] = propertySchema
		}
		if combineRequired {
			required, foundRequired, _ := unstructured.NestedStringSlice(schema, "required")
			if foundRequired {
				combinedRequired = append(combinedRequired, required...)
			}
		}
	}

	combinedSchema := map[string]any{
		"type":       Object,
		"properties": combinedProperties,
	}
	if combineRequired {
		combinedSchema["required"] = ToAny(combinedRequired)
	}
	return combinedSchema
}

func getToken(packagename, group, version, kind string) string {
	return fmt.Sprintf("%s:%s/%s:%s", packagename, group, version, kind)
}
