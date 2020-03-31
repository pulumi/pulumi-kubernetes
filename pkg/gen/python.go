// Copyright 2016-2018, Pulumi Corporation.
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
	"fmt"
	"io/ioutil"

	"github.com/cbroglie/mustache"

	pycodegen "github.com/pulumi/pulumi/pkg/codegen/python"
)

// --------------------------------------------------------------------------

// Main interface.

// --------------------------------------------------------------------------

// PythonClient will generate a Pulumi Kubernetes provider client SDK for Python.
func PythonClient(
	swagger map[string]interface{},
	templateDir string,
	rootInit func(initPy string) error,
	groupInit func(group, initPy string) error,
	customResource func(crPy string) error,
	versionInit func(group, version, initPy string) error,
	kindFile func(group, version, kind, kindPy string) error,
	casingFile func(casingPy string) error,
	yamlFile func(yamlPy string) error,
) error {
	definitions := swagger["definitions"].(map[string]interface{})

	// Generate casing tables from property names.
	// { properties: [ {name: fooBar, casedName: foo_bar}, ]}
	properties := allCamelCasePropertyNames(definitions, pythonOpts())
	cases := map[string][]map[string]string{"properties": make([]map[string]string, 0)}
	for _, name := range properties {
		cases["properties"] = append(cases["properties"],
			map[string]string{"name": name, "casedName": pycodegen.PyName(name)})
	}
	casingPy, err := mustache.RenderFile(
		fmt.Sprintf("%s/casing.py.mustache", templateDir), cases)
	if err != nil {
		return err
	}
	err = casingFile(casingPy)
	if err != nil {
		return err
	}

	groupsSlice := createGroups(definitions, pythonOpts())

	yamlPy, err := mustache.RenderFile(
		fmt.Sprintf("%s/yaml.py.mustache", templateDir),
		map[string]interface{}{
			"Groups": groupsSlice,
		})
	if err != nil {
		return err
	}
	err = yamlFile(yamlPy)
	if err != nil {
		return err
	}

	rootInitPy, err := mustache.RenderFile(
		fmt.Sprintf("%s/root__init__.py.mustache", templateDir),
		map[string]interface{}{
			"Groups": groupsSlice,
		})
	if err != nil {
		return err
	}
	err = rootInit(rootInitPy)
	if err != nil {
		return err
	}

	for _, group := range groupsSlice {
		if !group.HasTopLevelKinds() {
			continue
		}

		groupInitPy, err := mustache.RenderFile(
			fmt.Sprintf("%s/group__init__.py.mustache", templateDir), group)
		if err != nil {
			return err
		}
		if group.Group() == "apiextensions" {
			groupInitPy += fmt.Sprint(`
from .CustomResource import (CustomResource)
`)

			crBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/CustomResource.py", templateDir))
			if err != nil {
				return err
			}

			err = customResource(string(crBytes))
			if err != nil {
				return err
			}
		}

		err = groupInit(group.Group(), groupInitPy)
		if err != nil {
			return err
		}

		for _, version := range group.Versions() {
			if !version.HasTopLevelKinds() {
				continue
			}

			versionInitPy, err := mustache.RenderFile(
				fmt.Sprintf("%s/version__init__.py.mustache", templateDir), version)
			if err != nil {
				return err
			}

			err = versionInit(group.Group(), version.Version(), versionInitPy)
			if err != nil {
				return err
			}

			for _, kind := range version.TopLevelKinds() {
				inputMap := map[string]interface{}{
					"DefaultAPIVersion":       kind.DefaultAPIVersion(),
					"DeprecationComment":      kind.DeprecationComment(),
					"Comment":                 kind.Comment(),
					"Group":                   group.Group(),
					"Kind":                    kind.Kind(),
					"Properties":              kind.Properties(),
					"RequiredInputProperties": kind.RequiredInputProperties(),
					"OptionalInputProperties": kind.OptionalInputProperties(),
					"AdditionalSecretOutputs": kind.AdditionalSecretOutputs(),
					"Aliases":                 kind.Aliases(),
					"URNAPIVersion":           kind.URNAPIVersion(),
					"Version":                 version.Version(),
					"PulumiComment":           kind.PulumiComment(),
				}
				// Since mustache templates are logic-less, we have to add some extra variables
				// to selectively disable code generation for empty lists.
				additionalSecretOutputsPresent := len(kind.AdditionalSecretOutputs()) > 0
				aliasesPresent := len(kind.Aliases()) > 0
				inputMap["MergeOptsRequired"] = additionalSecretOutputsPresent || aliasesPresent
				inputMap["AdditionalSecretOutputsPresent"] = additionalSecretOutputsPresent
				inputMap["AliasesPresent"] = aliasesPresent

				kindPy, err := mustache.RenderFile(
					fmt.Sprintf("%s/kind.py.mustache", templateDir), inputMap)
				if err != nil {
					return err
				}

				err = kindFile(group.Group(), version.Version(), kind.Kind(), kindPy)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
