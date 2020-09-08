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
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/v2/codegen/nodejs"
	"github.com/pulumi/pulumi/sdk/v2/go/common/util/contract"
)

const nodejsMetaPath = "meta/v1.ts"
const nodejsMetaFile = `import * as k8s from "@pulumi/kubernetes";

export type ObjectMeta = k8s.types.input.meta.v1.ObjectMeta;
`

func (pg *PackageGenerator) genNodeJS(outputDir string) error {
	if files, err := pg.genNodeJSFiles(); err != nil {
		return err
	} else if err := writeFiles(files, outputDir); err != nil {
		return err
	}
	return nil
}

func (pg *PackageGenerator) genNodeJSFiles() (map[string]*bytes.Buffer, error) {
	pkg := pg.SchemaPackage()

	pkg.Language["nodejs"] = rawMessage(map[string]interface{}{
		"moduleToPackage": pg.moduleToPackage(),
	})

	files, err := nodejs.GeneratePackage(tool, pkg, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate nodejs package")
	}

	delete(pkg.Language, NodeJS)

	// Search and replace ${VERSION} with the crd2pulumi version in package.json, so the resources can be properly
	// registered
	packageJSON, ok := files["package.json"]
	if !ok {
		return nil, errors.New("cannot find generated package.json")
	}
	files["package.json"] = bytes.ReplaceAll(packageJSON, []byte("${VERSION}"), []byte(Version))

	// Create a helper `meta/v1.ts` script that exports the ObjectMeta class from the SDK. If there happens to already
	// be a `meta/v1.ts` file, then just append the script.
	if code, ok := files[nodejsMetaPath]; !ok {
		files[nodejsMetaPath] = []byte(nodejsMetaFile)
	} else {
		files[nodejsMetaPath] = append(code, []byte("\n"+nodejsMetaFile)...)
	}

	buffers := map[string]*bytes.Buffer{}
	for name, code := range files {
		buffers[name] = bytes.NewBuffer(code)
	}

	// Generates CustomResourceDefinition constructors. Soon this will be
	// replaced with `kube2pulumi`
	for _, crg := range pg.CustomResourceGenerators {
		definitionFileName := toLowerFirst(crg.Kind) + "Definition"

		// Create the customResourceDefinition.ts class
		path := filepath.Join(groupPrefix(crg.Group), definitionFileName+".ts")
		_, ok := buffers[path]
		contract.Assertf(!ok, "duplicate file at %s", path)
		buffer := &bytes.Buffer{}
		crg.genNodeJSDefinition(buffer)
		buffers[path] = buffer

		// Export in the index.ts file
		indexPath := filepath.Join(filepath.Dir(path), "index.ts")
		indexBuffer := buffers[indexPath]
		definitionClassName := crg.Kind + "Definition"
		exportCode := fmt.Sprintf(
			"import {%s} from \"./%s\";\nexport {%s};\n",
			definitionClassName,
			definitionFileName,
			definitionClassName,
		)
		indexBuffer.WriteString(exportCode)
	}

	return buffers, nil
}

const definitionImports = `import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

`

// Outputs the code for a CustomResourceDefinition class to the given buffer.
// Mutates crg.CustomResourceDefinition.Object by underscoring all hyphenated
// fields.
func (crg *CustomResourceGenerator) genNodeJSDefinition(buffer *bytes.Buffer) {
	className := crg.Kind + "Definition"
	var superClassName string
	if crg.APIVersion == v1 {
		superClassName = "k8s.apiextensions.v1.CustomResourceDefinition"
	} else {
		superClassName = "k8s.apiextensions.v1beta1.CustomResourceDefinition"
	}

	fmt.Fprint(buffer, definitionImports)
	fmt.Fprintf(buffer, "export class %s extends %s {\n", className, superClassName)
	fmt.Fprint(buffer, "\tconstructor(name: string, opts?: pulumi.CustomResourceOptions) {\n")
	fmt.Fprint(buffer, "\t\tsuper(name, ")

	UnderscoreFields(crg.CustomResourceDefinition.Object)
	definitionArgs, _ := json.MarshalIndent(crg.CustomResourceDefinition.Object, "\t\t", "\t")
	buffer.Write(definitionArgs)

	fmt.Fprint(buffer, ", opts)\n\t}\n}\n")
}
