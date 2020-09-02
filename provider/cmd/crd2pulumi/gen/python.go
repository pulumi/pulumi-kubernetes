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
	pschema "github.com/pulumi/pulumi/pkg/v2/codegen/schema"
	"github.com/pulumi/pulumi/sdk/v2/go/common/util/contract"
)

const pythonPackageDir = "pulumi_" + packageName

var unneededPythonFiles = []string{
	filepath.Join(pythonPackageDir, "README.md"),
	filepath.Join(pythonPackageDir, "provider.py"),
	"setup.py",
}

func (pg *PackageGenerator) genPython(types map[string]pschema.ObjectTypeSpec, baseRefs []string) (map[string]*bytes.Buffer, error) {
	pkg, err := getPackage(types, baseRefs)
	if err != nil {
		return nil, errors.Wrapf(err, "generating package")
	}
	pkg.Language["python"] = rawMessage(map[string]interface{}{
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

	// Remove unneeded files
	for _, unneededFile := range unneededPythonFiles {
		delete(files, unneededFile)
	}

	// Replace _utilities.py with our own hard-coded version
	utilitiesPath := filepath.Join("pulumi_"+packageName, "_utilities.py")
	_, ok := files[utilitiesPath]
	contract.Assertf(ok, "missing _utilities.py file")
	files[utilitiesPath] = []byte(pythonUtilitiesFile)

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
