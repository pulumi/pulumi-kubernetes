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

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/v2/codegen/nodejs"
)

const metaPath = "meta/v1.ts"
const metaFile = `import * as k8s from "@pulumi/kubernetes";

export type ObjectMeta = k8s.types.input.meta.v1.ObjectMeta;
`

// genNodeJS returns a mapping from each file path to its generated code
func (gen *CustomResourceGenerator) genNodeJS() (map[string][]byte, error) {
	objectTypeSpecs := gen.GetObjectTypeSpecs()
	baseRefs := gen.baseRefs()
	AddMetadataRefs(objectTypeSpecs, baseRefs)
	gen.AddAPIVersionAndKindProperties(objectTypeSpecs, baseRefs)

	// Generate package
	pkg, err := genPackage(objectTypeSpecs, baseRefs, NodeJS)
	if err != nil {
		return nil, errors.Wrapf(err, "generating package")
	}

	moduleToPackage := map[string]string{}
	for _, versionName := range gen.VersionNames() {
		moduleToPackage[versionName] = getVersion(versionName)
	}

	pkg.Language["nodejs"] = rawMessage(map[string]interface{}{
		"moduleToPackage": moduleToPackage,
	})

	files, err := nodejs.GeneratePackage(tool, pkg, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "generating nodejs package")
	}

	packageJSON, ok := files["package.json"]
	if !ok {
		return nil, errors.New("cannot find generated package.json")
	}
	files["package.json"] = bytes.ReplaceAll(packageJSON, []byte("${VERSION}"), []byte("2.0.0"))

	files[metaPath] = []byte(metaFile)

	files[gen.Kind+"Definition.ts"] = gen.genNodeJSDefinition()

	return files, nil
}

const definitionImports = `import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

`

// Generates a CustomResourceDefinition class for the entire CRD YAML
func (gen *CustomResourceGenerator) genNodeJSDefinition() []byte {
	buffer := &bytes.Buffer{}

	className := gen.Kind + "Definition"
	var superClassName string
	if gen.APIVersion == v1 {
		superClassName = "k8s.apiextensions.v1.CustomResourceDefinition"
	} else {
		superClassName = "k8s.apiextensions.v1beta1.CustomResourceDefinition"
	}

	fmt.Fprint(buffer, definitionImports)
	fmt.Fprintf(buffer, "export class %s extends %s {\n", className, superClassName)
	fmt.Fprint(buffer, "\tconstructor(name: string, opts?: pulumi.CustomResourceOptions) {\n")
	fmt.Fprint(buffer, "\t\tsuper(name, ")
	definitionArgs, _ := json.MarshalIndent(gen.CustomResourceDefinition.Object, "\t\t", "\t")
	buffer.Write(definitionArgs)
	fmt.Fprint(buffer, ", opts)\n\t}\n}\n")

	return buffer.Bytes()
}
