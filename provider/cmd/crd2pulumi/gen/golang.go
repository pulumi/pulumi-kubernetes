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

	"github.com/pkg/errors"
	go_gen "github.com/pulumi/pulumi/pkg/v2/codegen/go"
)

// genGo returns a map from each version's name to a buffer containing
// its generated code.
func (gen *CustomResourceGenerator) genGo() (map[string]*bytes.Buffer, error) {
	// Set up objectTypeSpecs. Notice that since we can't properly reference
	// external types for the Go codegen, we add a fake Metadata spec
	objectTypeSpecs := gen.GetObjectTypeSpecs()
	AddPlaceholderMetadataSpec(objectTypeSpecs)
	baseRefs := gen.baseRefs()
	AddMetadataRefs(objectTypeSpecs, baseRefs)
	gen.AddAPIVersionAndKindProperties(objectTypeSpecs, baseRefs)

	// Generate the package
	pkg, err := genPackage(objectTypeSpecs, baseRefs, Go)
	if err != nil {
		return nil, errors.Wrapf(err, "generating package")
	}
	// We added a fake Metadata spec so we hard-code its import to point to its
	// actual package
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
	allTypes, err := go_gen.CRDTypes(tool, pkg)

	buffers := map[string]*bytes.Buffer{}
	for _, versionName := range gen.VersionNames() {
		types, ok := allTypes[versionName]
		if !ok {
			return nil, errors.Errorf("cannot find generated Go code for %s", versionName)
		}
		buffers[versionName] = types
	}

	return buffers, nil
}
