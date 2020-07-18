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

package nodejs

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi/pkg/v2/codegen/schema"
)

const tool = "crdgen"

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

func GenerateTypeScriptTypes(types map[string]schema.ObjectTypeSpec) ([]byte, error) {
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
		return []byte{}, err
	}

	// Generate all the code for the package.
	files, err := GeneratePackage(tool, pkg, nil)
	if err != nil {
		return []byte{}, err
	}

	// Extract the relevant generated code.
	// Here, I'm just returning the content of the in-memory generated `types/input.ts` file,
	// which has generated interfaces of the "input shape". Not sure if we want the "output shape"
	// instead (or in addition to) that are in `types/output.ts`, or if we want something entirely different.
	if file, ok := files["types/input.ts"]; ok {
		return file, nil
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

	return []byte{}, nil
}
