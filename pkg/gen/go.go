// Copyright 2016-2019, Pulumi Corporation.
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
	"path/filepath"

	"github.com/cbroglie/mustache"
	"github.com/pkg/errors"
)

// --------------------------------------------------------------------------

// Main interface.

// --------------------------------------------------------------------------

// GoClient will generate a Pulumi Kubernetes provider client SDK for Go.
func GoClient(
	swagger map[string]interface{},
	templateDir string,
) (files map[string]string, err error) {
	definitions := swagger["definitions"].(map[string]interface{})

	files = make(map[string]string)

	groups := createGroups(definitions, goOpts())
	for _, group := range groups {
		for _, version := range group.Versions() {
			for _, kind := range version.Kinds() {
				inputMap := map[string]interface{}{
					"RawAPIVersion":           kind.RawAPIVersion(),
					"Comment":                 kind.Comment(),
					"Group":                   group.Group(),
					"Kind":                    kind.Kind(),
					"PropertyKind":            kind.Kind(),
					"Properties":              kind.Properties(),
					"RequiredInputProperties": kind.RequiredInputProperties(),
					"OptionalInputProperties": kind.OptionalInputProperties(),
					"AdditionalSecretOutputs": kind.AdditionalSecretOutputs(),
					"Aliases":                 kind.Aliases(),
					"URNAPIVersion":           kind.URNAPIVersion(),
					"Version":                 version.Version(),
					"PulumiComment":           kind.pulumiComment,
					"GoImports":               kind.GoImports(),
					"IsPropertyType":          kind.IsPropertyType(),
					"IsArrayElement":          kind.IsArrayElement(),
					"IsMapElement":            kind.IsMapElement(),
					"NeedsErrors":             len(kind.RequiredInputProperties()) > 0,
				}
				// Since mustache templates are logic-less, we have to add some extra variables
				// to selectively disable code generation for empty lists.
				additionalSecretOutputsPresent := len(kind.AdditionalSecretOutputs()) > 0
				aliasesPresent := len(kind.Aliases()) > 0
				inputMap["MergeOptsRequired"] = additionalSecretOutputsPresent || aliasesPresent
				inputMap["AdditionalSecretOutputsPresent"] = additionalSecretOutputsPresent
				inputMap["AliasesPresent"] = aliasesPresent

				if !kind.IsNested() && kind.IsPropertyType() {
					inputMap["PropertyKind"] = kind.Kind() + "Property"
				}

				templateFile := "kind.go.mustache"
				if kind.IsNested() {
					templateFile = "nestedKind.go.mustache"
				}

				text, err := mustache.RenderFile(filepath.Join(templateDir, templateFile), inputMap)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to generate code for %s/%s/%s (%s)",
						group.Group(), version.Version(), kind.Kind(), templateFile)
				}
				files[filepath.Join(group.Group(), version.Version(), kind.Kind()+".go")] = text
			}
		}
	}
	return
}
