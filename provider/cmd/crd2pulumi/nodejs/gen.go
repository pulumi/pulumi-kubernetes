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

	if file, ok := files["types/input.ts"]; ok {
		return file, nil
	}

	return []byte{}, nil
}
