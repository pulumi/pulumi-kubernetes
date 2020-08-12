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
	"path/filepath"

	"github.com/pkg/errors"
	go_gen "github.com/pulumi/pulumi/pkg/v2/codegen/go"
	pschema "github.com/pulumi/pulumi/pkg/v2/codegen/schema"
)

func (pg *PackageGenerator) genGo(types map[string]pschema.ObjectTypeSpec, baseRefs []string) (map[string]*bytes.Buffer, error) {
	AddPlaceholderMetadataSpec(types)

	// Generate the package
	pkg, err := getPackage(types, baseRefs)
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
	buffers, err := go_gen.CRDTypes(tool, pkg)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate Go code")
	}

	files := map[string]*bytes.Buffer{}

	for groupVersion, plural := range pg.GroupVersionsToPlural() {
		buffer, ok := buffers[groupVersion]
		if !ok {
			return nil, errors.New("could not find %s in generated Go code")
		}
		files[filepath.Join(plural, getVersion(groupVersion), plural+".go")] = buffer
	}

	return files, nil
}
