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
	"github.com/pulumi/pulumi/pkg/v2/codegen/python"
	"github.com/pulumi/pulumi/sdk/v2/go/common/util/contract"
)

const pythonPackageDir = "pulumi_" + packageName
const pythonMetaFile = `from pulumi_kubernetes.meta.v1._inputs import *
import pulumi_kubernetes.meta.v1.outputs
`

var unneededPythonFiles = []string{
	filepath.Join(pythonPackageDir, "README.md"),
	"setup.py",
}

func (pg *PackageGenerator) genPython(outputDir string) error {
	if files, err := pg.genPythonFiles(); err != nil {
		return err
	} else if err := writeFiles(files, outputDir); err != nil {
		return err
	}
	return nil
}

func (pg *PackageGenerator) genPythonFiles() (map[string]*bytes.Buffer, error) {
	pkg := pg.SchemaPackageWithObjectMetaType()

	pkg.Language[Python] = rawMessage(map[string]interface{}{
		"compatibility":       "kubernetes20",
		"moduleNameOverrides": pg.moduleToPackage(),
		"requires": map[string]string{
			"pulumi":   "\u003e=2.0.0,\u003c3.0.0",
			"pyyaml":   "\u003e=5.1,\u003c5.2",
			"requests": "\u003e=2.21.0,\u003c2.22.0",
		},
		"ignorePyNamePanic": true,
	})

	files, err := python.GeneratePackage(tool, pkg, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate Go package")
	}

	delete(pkg.Language, Python)

	// Remove unneeded files
	for _, unneededFile := range unneededPythonFiles {
		delete(files, unneededFile)
	}

	// Replace _utilities.py with our own hard-coded version
	utilitiesPath := filepath.Join(pythonPackageDir, "_utilities.py")
	_, ok := files[utilitiesPath]
	contract.Assertf(ok, "missing _utilities.py file")
	files[utilitiesPath] = []byte(pythonUtilitiesFile)

	// Import the actual SDK ObjectMeta types in place of our placeholder ones
	metaPath := filepath.Join(pythonPackageDir, "meta_v1", "__init__.py")
	code, ok := files[metaPath]
	contract.Assertf(ok, "missing meta_v1/__init__.py file")
	files[metaPath] = append(code, []byte(pythonMetaFile)...)

	buffers := map[string]*bytes.Buffer{}
	for name, code := range files {
		buffers[name] = bytes.NewBuffer(code)
	}
	return buffers, nil
}

const pythonUtilitiesFile = `from pulumi_kubernetes import _utilities


def get_env(*args):
    return _utilities.get_env(*args)


def get_env_bool(*args):
    return _utilities.get_env_bool(*args)


def get_env_int(*args):
    return _utilities.get_env_int(*args)


def get_env_float(*args):
    return _utilities.get_env_float(*args)


def get_version():
    return _utilities.get_version()
`
