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
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/v2/codegen/nodejs"
)

// genNodeJS returns a buffer containing all the generated code
func (gen *CustomResourceGenerator) genNodeJS() (*bytes.Buffer, error) {
	objectTypeSpecs := gen.GetObjectTypeSpecs()
	baseRefs := gen.baseRefs()
	AddMetadataRefs(objectTypeSpecs, baseRefs)
	AddArgsSuffix(objectTypeSpecs, baseRefs)

	// Generate package
	pkg, err := genPackage(objectTypeSpecs, baseRefs, NodeJS)
	if err != nil {
		return nil, errors.Wrapf(err, "generating package")
	}

	// Generate all the code for the package
	buffer, err := nodejs.GenCRDTypes(tool, pkg)
	if err != nil {
		return nil, errors.Wrapf(err, "generating types")
	}
	gen.genNodeJSClasses(buffer)

	return buffer, nil
}

// Writes the namespaced NodeJS classes for each version to the given writer.
func (gen *CustomResourceGenerator) genNodeJSClasses(w io.Writer) {
	versions := gen.Versions()
	kind := gen.Kind
	group := gen.Group
	name := gen.Name()

	// Generates a CustomResource sub-class for a single version
	genResourceClass := func(w io.Writer, version, kind, name, group string) {
		argsName := fmt.Sprintf("%s.%sArgs", version, kind)
		apiVersion := fmt.Sprintf("%s/%s", group, version)
		fmt.Fprintf(w, "export namespace %s {\n", version)
		fmt.Fprintf(w, "\texport class %s extends k8s.apiextensions.CustomResource {\n", kind)
		fmt.Fprintf(w, "\t\tpublic static get%s(name: string, id: pulumi.Input<pulumi.ID>): %s {\n", kind, kind)
		fmt.Fprintf(w, "\t\t\treturn k8s.apiextensions.CustomResource.get(name, { apiVersion: \"%s\", kind: \"%s\", id: id })\n", apiVersion, kind)
		fmt.Fprint(w, "\t\t}\n\n")
		fmt.Fprintf(w, "\t\tconstructor(name: string, args?: %s, opts?: pulumi.CustomResourceOptions) {\n", argsName)
		fmt.Fprintf(w, "\t\t\tsuper(name, { apiVersion: \"%s\", kind: \"%s\", ...args }, opts)\n", apiVersion, kind)
		fmt.Fprint(w, "\t\t}\n\t}\n}\n\n")
	}

	// Generates a CustomResourceDefinition class for the entire CRD YAML
	genDefinitionClass := func(w io.Writer, kind, apiVersion string) {
		className := kind + "Definition"
		var superClassName string
		if apiVersion == v1 {
			superClassName = "k8s.apiextensions.v1.CustomResourceDefinition"
		} else {
			superClassName = "k8s.apiextensions.v1beta1.CustomResourceDefinition"
		}

		fmt.Fprintf(w, "export class %s extends %s {\n", className, superClassName)
		fmt.Fprint(w, "\tconstructor(name: string, opts?: pulumi.CustomResourceOptions) {\n")
		fmt.Fprint(w, "\t\tsuper(name, ")
		definitionArgs, _ := json.MarshalIndent(gen.CustomResourceDefinition.Object, "\t\t", "\t")
		w.Write(definitionArgs)
		fmt.Fprint(w, ", opts)\n\t}\n}\n")
	}

	for _, version := range versions {
		genResourceClass(w, version, kind, name, group)
	}
	genDefinitionClass(w, kind, gen.APIVersion)
}
