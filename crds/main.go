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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pulumi/pulumi/pkg/v2/codegen/schema"
	pschema "github.com/pulumi/pulumi/pkg/v2/codegen/schema"
	unstruct "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const tool = "crdgen"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "error: expected <resourcedefinition.yaml> argument\n")
		os.Exit(-1)
	}

	yamlPath := os.Args[1]
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading the yaml file: %v\n", err)
		os.Exit(-1)
	}

	crd, err := UnmarshalYaml(yamlFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error unmarshalling yaml: %v\n", err)
		os.Exit(-1)
	}

	types := getTypes(crd)

	// Generate TypeScript for the types.
	code, err := generateTypeScriptTypes(types)
	if err != nil {
		log.Printf("error: %v", err)
	}

	// Write out the generated code to standard output.
	fmt.Println(code)
}

// getTypes parses all the types within the given crd (an unmarshalled
// resourcedefinition.yaml file). Returns a mapping from each type's name in the
// format "{.spec.group}/{.spec.versions[*].name}:{.spec.names.kind}" to its
// corresponding pschema.ObjectTypeSpec.
func getTypes(crd unstruct.Unstructured) map[string]pschema.ObjectTypeSpec {
	versions, _, _ := NestedMapSlice(crd.Object, "spec", "versions")
	metadataName, _, _ := unstruct.NestedString(crd.Object, "metadata", "name")
	metadataName = cleanMetadataName(metadataName)
	kind, _, _ := unstruct.NestedString(crd.Object, "spec", "names", "kind")

	types := make(map[string]pschema.ObjectTypeSpec)
	for _, version := range versions {
		schema, _, _ := unstruct.NestedMap(version, "schema", "openAPIV3Schema")
		versionName, _, _ := unstruct.NestedString(version, "name")
		baseRef := fmt.Sprintf("crds:%s/%s:%s", metadataName, versionName, kind)
		addType(schema, baseRef, types)
	}

	return types
}

/*
getType converts the given OpenAPI v3 `schema` to a ObjectTypeSpec and adds
it to the `types` map under the given `name`. Recursively converts and adds all
nested schemas as well.
*/
func addType(schema map[string]interface{}, name string, types map[string]pschema.ObjectTypeSpec) {
	description, _, _ := unstruct.NestedString(schema, "description")
	schemaType, _, _ := unstruct.NestedString(schema, "type")
	properties, _, _ := unstruct.NestedMap(schema, "properties")
	required, _, _ := unstruct.NestedStringSlice(schema, "required")

	propertySpecs := make(map[string]pschema.PropertySpec)
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
// schema. Handles nested pschema.TypeSpecs in case the schema type is of
// array, object, or a union type (anyOf). Also recursively converts and adds
// all schemas of type object to the types map.
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

	return pschema.TypeSpec{
		Type: schemaType,
	}
}

func genConstructor(argsName, apiVersion, kind string) []string {
	header := fmt.Sprintf("export class %s extends k8s.apiextensions.CustomResource {\n", kind)
	constructor := fmt.Sprintf("\tconstructor(name: string, args?: %s, opts?: pulumi.CustomResourceOptions) {\n", argsName)
	super := fmt.Sprintf("\t\tsuper(name, { apiVersion: \"%s\", kind: \"%s\", ...args }, opts)\n", apiVersion, kind)
	return []string{
		header,
		constructor,
		super,
		"\t}\n",
		"}\n",
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

func generateTypeScriptTypes(types map[string]schema.ObjectTypeSpec) (string, error) {
	// Create a fake `PackageSpec` that includes the specifed types.
	// This lets us call the public `nodejs.GeneratePackage` function, which generates code
	// for an entire package (which isn't really what we want to do, but is existing code
	// we have with a public entry point, that I'll leverage here for demonstration purposes).
	// We then return contents of a certain in-memory file that is generated.

	// Note: This is unlikely how we actually want to do the real code gen for the CRD types.
	// See my comments at the bottom of this function.

	// Create some a property map that references our types, that we can add
	// as properties to a fake resource, otherwise, our types won't be generated
	// by the Node.js codegen in `types/input.ts` or `types/output.ts` in-memory files.
	properties := map[string]schema.PropertySpec{}
	for name := range types {
		properties[name] = schema.PropertySpec{
			TypeSpec: schema.TypeSpec{
				Ref: fmt.Sprintf("#/types/%s", name),
			},
		}
	}

	// Create a fake package that includes the types passed-in to this function.
	var pkgSpec = schema.PackageSpec{
		// Include the passed-in types.
		Types: types,

		// Create a fake resource that has the properties.
		Resources: map[string]schema.ResourceSpec{
			"prov:module/resource:Resource": {
				ObjectTypeSpec: schema.ObjectTypeSpec{
					Properties: properties,
				},
				InputProperties: properties,
			},
		},

		// Apparently, the Node.js codegen is expected a non-nil map, so include it.
		Language: map[string]json.RawMessage{
			"nodejs": []byte("{}"),
		},
	}

	// Convert the PackageSpec into a Package.
	pkg, err := schema.ImportSpec(pkgSpec, nil)
	if err != nil {
		return "", err
	}

	// Generate all the code for the package.
	files, err := GeneratePackage(tool, pkg, nil)
	if err != nil {
		return "", err
	}

	// Extract the relevant generated code.
	// Here, I'm just returning the content of the in-memory generated `types/input.ts` file,
	// which has generated interfaces of the "input shape". Not sure if we want the "output shape"
	// instead (or in addition to) that are in `types/output.ts`, or if we want something entirely different.
	if file, ok := files["types/input.ts"]; ok {
		return string(file), nil
	}

	// Note: The generated code isn't quite right for what we need for the CRD code gen. For example,
	// There are some comments at the top (with "WARNING") about being generated code that we probably
	// don't want to include. There are some imports that don't make sense outside of package code
	// generation, and I'm not sure we want to use TypeScript namespaces.

	// But it does give a sense for what's possible calling into existing code.

	// I think you will need to refactor some of the existing code in
	// https://github.com/pulumi/pulumi/blob/master/pkg/codegen/nodejs/gen.go and expose a
	// public function we can call that will generate the right code. Perhaps leveraging/refactoring
	// some existing code we have and generalizing it so we can use it for this purpose.

	// Maybe something like generalizing, refactoring, and making public the existing `genPlainType` function?
	// https://github.com/pulumi/pulumi/blob/bb358c4d217334ad654f153cc155f84a4fb1933b/pkg/codegen/nodejs/gen.go#L243

	// And then we'll want to do something similar for each of the other languages (C#, Go, Python).

	return "", nil
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

// Converts a []string to []interface{}
func genericizeStringSlice(stringSlice []string) interface{} {
	genericSlice := make([]interface{}, len(stringSlice))
	for i, v := range stringSlice {
		genericSlice[i] = v
	}
	return genericSlice
}

// Given a list of schemas, returns a single schema representing their
// intersection type. Combines each of their properties, and descriptions into a single
// schema. If combineRequired is true, then required properties are also
// combined into the returned schema; otherwise all properties in the returned
// schema are optional.
func combineSchemas(combineRequired bool, schemas ...map[string]interface{}) map[string]interface{} {
	combinedProperties := make(map[string]interface{})
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
		combinedSchema["required"] = genericizeStringSlice(combinedRequired)
	}
	return combinedSchema
}

// PrettyPrint properly formats and indents an unstructured value, and prints it
// to stdout.
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
