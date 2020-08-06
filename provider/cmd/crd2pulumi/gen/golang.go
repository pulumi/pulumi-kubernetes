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
	"fmt"
	"io"

	"github.com/pkg/errors"
	go_gen "github.com/pulumi/pulumi/pkg/v2/codegen/go"
)

// genGo returns a map from each version's name to a buffer containing
// its generated code.
func (gen *CustomResourceGenerator) genGo() (map[string]*bytes.Buffer, error) {
	objectTypeSpecs := GetObjectTypeSpecs(gen.Versions, gen.Name(), gen.Kind)
	AddPlaceholderMetadataSpec(objectTypeSpecs)
	AddMetadataRefs(gen.Versions, gen.Name(), gen.Kind, objectTypeSpecs)

	pkg, err := genPackage(objectTypeSpecs, Go)
	if err != nil {
		return nil, errors.Wrapf(err, "generating package")
	}
	pkg.Language["go"] = rawMessage(map[string]interface{}{
		"importBasePath": "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes",
		"moduleToPackage": map[string]interface{}{
			"meta/v1": "meta/v1",
		},
		"packageImportAliases": map[string]interface{}{
			"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/meta/v1": "metav1",
		},
	})

	// Generate all the code for the package.
	allTypes, err := go_gen.GenCRDTypes(tool, pkg)

	buffers := map[string]*bytes.Buffer{}

	versionNames := gen.VersionNames()
	for _, versionName := range versionNames {
		types, ok := allTypes[versionName]
		if !ok {
			return nil, errors.Errorf("cannot find generated Go code for %s", versionName)
		}
		gen.genGoConstructor(types, getVersion(versionName))
		buffers[versionName] = types
	}

	return buffers, nil
}

func (gen *CustomResourceGenerator) genGoConstructor(w io.Writer, version string) {
	fmt.Fprintf(w, "\nfunc New%s(ctx *pulumi.Context, name string, args *%sArgs, opts ...pulumi.ResourceOption) (*apiextensions.CustomResource, error) {", gen.Kind, gen.Kind)
	fmt.Fprint(w, unmarshalCode)
	fmt.Fprintf(w, "\tcustomResourceArgs := apiextensions.CustomResourceArgs{\n")
	apiVersion := fmt.Sprintf("%s/%s", gen.Group, version)
	fmt.Fprintf(w, "\t\tApiVersion: pulumi.String(\"%s\"),\n", apiVersion)
	fmt.Fprintf(w, "\t\tKind: pulumi.String(\"%s\"),\n", gen.Kind)
	fmt.Fprint(w, "\t\tMetadata: args.Metadata,\n")
	fmt.Fprint(w, "\t\tOtherFields: otherFields,\n")
	fmt.Fprint(w, "\t}\n\n")
	fmt.Fprint(w, "\treturn apiextensions.NewCustomResource(ctx, name, &customResourceArgs, opts...)\n")
	fmt.Fprint(w, "}\n")
	fmt.Fprint(w, lowerCode)
}

const unmarshalCode = `
	m := structs.Map(args.Spec)
	otherFields, ok := lower(m).(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("could not parse other fields %v", m)
	}
`

const lowerCode = `
func lower(object interface{}) interface{} {
	switch object := object.(type) {
	case map[string]interface{}:
		lowerObject := make(map[string]interface{}, len(object))
		for upperKey, value := range object {
			lowerKey := ""
			if upperKey != "" {
				lowerKey = strings.ToLower(string(upperKey[0])) + upperKey[1:]
			}
			lowerObject[lowerKey] = lower(value)
		}
		return lowerObject
	case []interface{}:
		for i := range object {
			object[i] = lower(object[i])
		}
		return object
	default:
		return object
	}
}
`
