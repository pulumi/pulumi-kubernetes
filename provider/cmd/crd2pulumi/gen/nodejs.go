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
	"regexp"
	"unicode"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/v2/codegen/nodejs"
	pschema "github.com/pulumi/pulumi/pkg/v2/codegen/schema"
	"github.com/pulumi/pulumi/sdk/v2/go/common/util/contract"
)

const metaPath = "meta/v1.ts"
const metaFile = `import * as k8s from "@pulumi/kubernetes";

export type ObjectMeta = k8s.types.input.meta.v1.ObjectMeta;
`

var alphanumericRegex = regexp.MustCompile("[^a-zA-Z0-9]+")

// removes all non-alphanumeric characters
func removeNonAlphanumeric(input string) string {
	return alphanumericRegex.ReplaceAllString(input, "")
}

// un-capitalizes the first character of a string
func toLowerFirst(input string) string {
	if input == "" {
		return ""
	}
	return string(unicode.ToLower(rune(input[0]))) + input[1:]
}

func (pg *PackageGenerator) genNodeJS(types map[string]pschema.ObjectTypeSpec, baseRefs []string) (map[string]*bytes.Buffer, error) {
	pkg, err := getPackage(types, baseRefs)
	if err != nil {
		return nil, errors.Wrap(err, "could not create package")
	}

	moduleToPackage := map[string]string{}
	for _, groupVersion := range pg.GroupVersions() {
		group, version := splitGroupVersion(groupVersion)
		moduleToPackage[groupVersion] = removeNonAlphanumeric(group) + "/" + version
	}
	pkg.Language["nodejs"] = rawMessage(map[string]interface{}{
		"moduleToPackage": moduleToPackage,
	})

	files, err := nodejs.GeneratePackage(tool, pkg, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate nodejs package")
	}

	// Search and replace ${VERSION} with 2.0.0 in package.json, so the
	// resources can be properly registered
	packageJSON, ok := files["package.json"]
	if !ok {
		return nil, errors.New("cannot find generated package.json")
	}
	files["package.json"] = bytes.ReplaceAll(packageJSON, []byte("${VERSION}"), []byte("2.0.0"))

	// Create a helper 'meta/v1.ts' script that just exports the ObjectMeta class
	files[metaPath] = []byte(metaFile)

	buffers := map[string]*bytes.Buffer{}
	for name, code := range files {
		buffers[name] = bytes.NewBuffer(code)
	}

	// Generates CustomResourceDefinition constructors. Soon this will be
	// replaced with `kube2pulumi`
	for _, crg := range pg.CustomResourceGenerators {
		path := filepath.Join(removeNonAlphanumeric(crg.Group), toLowerFirst(crg.Kind)+"Definition.ts")
		_, ok := buffers[path]
		contract.Assertf(!ok, "duplicate file at %s", path)
		buffer := &bytes.Buffer{}
		crg.genNodeJSDefinition(buffer)
		buffers[path] = buffer
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
