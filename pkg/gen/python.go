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

	"github.com/cbroglie/mustache"
)

// --------------------------------------------------------------------------

// Main interface.

// --------------------------------------------------------------------------

// PythonClient will generate a Pulumi Kubernetes provider client SDK for nodejs.
func PythonClient(
	swagger map[string]interface{}, templateDir string,
	rootInit func(initPy string) error,
	groupInit func(group, initPy string) error,
	versionInit func(group, version, initPy string) error,
	kindFile func(group, version, kind, kindPy string) error,
) error {
	definitions := swagger["definitions"].(map[string]interface{})

	groupsSlice := createGroups(definitions, pythonProvider())

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

		groupInitPy, err := mustache.RenderFile(
			fmt.Sprintf("%s/group__init__.py.mustache", templateDir), group)
		if err != nil {
			return err
		}

		err = groupInit(group.Group(), groupInitPy)
		if err != nil {
			return err
		}

		for _, version := range group.Versions() {
			versionInitPy, err := mustache.RenderFile(
				fmt.Sprintf("%s/version__init__.py.mustache", templateDir), version)
			if err != nil {
				return err
			}

			err = versionInit(group.Group(), version.Version(), versionInitPy)
			if err != nil {
				return err
			}

			for _, kind := range version.Kinds() {
				kindPy, err := mustache.RenderFile(
					fmt.Sprintf("%s/kind.py.mustache", templateDir), kind)
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
