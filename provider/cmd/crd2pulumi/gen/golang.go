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
)

var unneededGoFiles = []string{
	"doc.go",
	"provider.go",
	"meta/v1/pulumiTypes.go",
}

func (pg *PackageGenerator) genGo(outputDir string) error {
	if files, err := pg.genGoFiles(); err != nil {
		return err
	} else if err := writeFiles(files, outputDir); err != nil {
		return err
	}
	return nil
}

func (pg *PackageGenerator) genGoFiles() (map[string]*bytes.Buffer, error) {
	pkg := pg.SchemaPackageWithObjectMetaType()

	moduleToPackage := pg.moduleToPackage()
	moduleToPackage["meta/v1"] = "meta/v1"
	pkg.Language["go"] = rawMessage(map[string]interface{}{
		"importBasePath":  "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes",
		"moduleToPackage": moduleToPackage,
		"packageImportAliases": map[string]interface{}{
			"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/meta/v1": "metav1",
		},
	})

	files, err := go_gen.GeneratePackage(tool, pkg)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate Go package")
	}

	delete(pkg.Language, Go)

	buffers := map[string]*bytes.Buffer{}

	// Now we remove the "crds/" file path prefix
	for path, code := range files {
		newPath, err := filepath.Rel(packageName, path)
		if err != nil {
			return nil, errors.Wrapf(err, "could not remove \"crds/\" prefix")
		}
		buffers[newPath] = bytes.NewBuffer(code)
	}

	for _, unneededFile := range unneededGoFiles {
		delete(buffers, unneededFile)
	}

	return buffers, nil
}
