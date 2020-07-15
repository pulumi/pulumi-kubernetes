package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/pulumi/pulumi/pkg/v2/codegen/nodejs"
	"github.com/pulumi/pulumi/pkg/v2/codegen/schema"
	pschema "github.com/pulumi/pulumi/pkg/v2/codegen/schema"
	unstruct "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const tool = "crdgen"

func main() {
	yamlPath := os.Args[1]
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		log.Printf("yamlFile.Get err	#%v ", err)
	}

	crd, err := UnmarshalYaml(yamlFile)
	if err != nil {
		log.Printf("%v", err)
	}

	versions, _, _ := NestedMapSlice(crd.Object, "spec", "versions")
	metadataName, _, _ := unstruct.NestedString(crd.Object, "metadata", "name")
	kind, _, _ := unstruct.NestedString(crd.Object, "spec", "names", "kind")

	// Populates the 'types' map with converted ObjectTypeSpecs for all objects
	// across all versions defined in the yamlFile.
	types := make(map[string]pschema.ObjectTypeSpec)
	for _, version := range versions {
		versionName, _, _ := unstruct.NestedString(version, "name")
		schema, _, _ := unstruct.NestedMap(version, "schema", "openAPIV3Schema")
		baseRef := fmt.Sprintf("crds:%s/%s:%s", metadataName, versionName, kind)
		generateType(schema, baseRef, types)
	}

	// Generate TypeScript for the types.
	code, err := generateTypeScriptTypes(types)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(-1)
	}

	// Write out the generated code to standard output.
	fmt.Println(code)
}

/*
generateType converts the given OpenAPI v3 `schema` to a ObjectTypeSpec and adds
it to the `types` map under the given `name`. Recursively converts and adds all
nested schemas as well.

For example, if we had the following inputs:
schema := map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"cronSpec:" map[string]interface{} {
			"type": "string",
		},
		"image:" map[string]interface{} {
			"type": "string",
		},
		"replicas:" map[string]interface{} {
			"type": "integer",
		},
	}
}
name := "crds:stable.example.com/v1:CronTab"
types := make(map[string]pschema.ObjectTypeSpec)

generateType adds the following ObjectTypeSpecs to the 'types' map:
types := map[string]schema.ObjectTypeSpec{
	"crds:stable.example.com/v1:CronTabArgs": {
		Type: "object",
		Properties: map[string]schema.PropertySpec{
			"spec": {
				TypeSpec: schema.TypeSpec{
					Ref: "#/types/crds:stable.example.com/v1:CronTabSpec",
				},
			},
		},
	},
	"crds:stable.example.com/v1:CronTabSpecArgs": {
		Type: "object",
		Properties: map[string]schema.PropertySpec{
			"cronSpec": {
				TypeSpec: schema.TypeSpec{
					Type: "string",
				},
			},
			"image": {
				TypeSpec: schema.TypeSpec{
					Type: "string",
				},
			},
			"replicas": {
				TypeSpec: schema.TypeSpec{
					Type: "string",
				},
			},
		},
	},
}
*/
func generateType(schema map[string]interface{}, name string, types map[string]pschema.ObjectTypeSpec) {
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
// array or object. Also recursively converts and adds all schemas of type
// object to the types map.
func getTypeSpec(schema map[string]interface{}, name string, types map[string]pschema.ObjectTypeSpec) pschema.TypeSpec {
	if schema == nil {
		return pschema.TypeSpec{}
	}

	schemaType, foundSchemaType, _ := unstruct.NestedString(schema, "type")
	if !foundSchemaType {
		log.Printf("could not find type in %v", schema)
	}

	if schemaType == "array" {
		items, _, _ := unstruct.NestedMap(schema, "items")
		arrayTypeSpec := getTypeSpec(items, name, types)
		return pschema.TypeSpec{
			Type:  "array",
			Items: &arrayTypeSpec,
		}
	} else if schemaType == "object" {
		generateType(schema, name, types)
		additionalProperties, _, _ := unstruct.NestedMap(schema, "additionalProperties")
		additionalPropertiesTypeSpec := getTypeSpec(additionalProperties, name, types)
		return pschema.TypeSpec{
			Ref:                  "#/types/" + name + "Args",
			AdditionalProperties: &additionalPropertiesTypeSpec,
		}
	} else {
		return pschema.TypeSpec{
			Type: schemaType,
		}
	}
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
	files, err := nodejs.GeneratePackage(tool, pkg, nil)
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

// PrettyPrint properly formats and indents an unstructured value, and prints it
// to stdout.
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
