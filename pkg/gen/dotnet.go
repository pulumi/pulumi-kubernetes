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
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cbroglie/mustache"
)

var pascalCaseMapping = map[string]string{
	"admissionregistration": "AdmissionRegistration",
	"apps":                  "Apps",
	"auditregistration":     "AuditRegistraion",
	"authentication":        "Authentication",
	"apiextensions":         "ApiExtensions",
	"authorization":         "Authorization",
	"autoscaling":           "Autoscaling",
	"apiregistration":       "ApiRegistration",
	"batch":                 "Batch",
	"certificates":          "Certificates",
	"coordination":          "Coordination",
	"core":                  "Core",
	"discovery":             "Discovery",
	"events":                "Events",
	"extensions":            "Extensions",
	"networking":            "Networking",
	"meta":                  "Meta",
	"node":                  "Node",
	"policy":                "Policy",
	"rbac":                  "Rbac",
	"scheduling":            "Scheduling",
	"settings":              "Settings",
	"storage":               "Storage",
	"v1":                    "V1",
	"v1alpha1":              "V1Alpha1",
	"v1beta1":               "V1Beta1",
	"v1beta2":               "V1Beta2",
	"v2":                    "V2",
	"v2alpha1":              "V2Alpha1",
	"v2beta1":               "V2Beta1",
	"v2beta2":               "V2Beta2",
}

func pascalCase(name string) string {
	return strings.ToUpper(name[:1]) + name[1:]
}

// --------------------------------------------------------------------------

// Main interface.

// --------------------------------------------------------------------------

// DotnetClient will generate a Pulumi Kubernetes provider client SDK for .NET.
func DotnetClient(
	swagger map[string]interface{},
	templateDir string,
) (inputsts, outputsts string, groups map[string]string, err error) {
	definitions := swagger["definitions"].(map[string]interface{})

	inputGroupsSlice := createGroups(definitions, dotnetInputs())
	inputsts, err = mustache.RenderFile(fmt.Sprintf("%s/typesInput.cs.mustache", templateDir),
		map[string]interface{}{
			"Groups": inputGroupsSlice,
		})
	if err != nil {
		return
	}

	outputGroupsSlice := createGroups(definitions, dotnetOutputs())
	outputsts, err = mustache.RenderFile(fmt.Sprintf("%s/typesOutput.cs.mustache", templateDir),
		map[string]interface{}{
			"Groups": outputGroupsSlice,
		})
	if err != nil {
		return
	}

	groupsSlice := createGroups(definitions, dotnetProvider())
	fmt.Printf("%v\n", groupsSlice)

	groups = make(map[string]string)

	for _, group := range groupsSlice {

		pascalGroup, ok := pascalCaseMapping[group.Group()]
		if !ok {
			return "", "", nil, fmt.Errorf("No pascal case mapping defined for %q", group.Group())
		}

		for _, version := range group.Versions() {

			pascalVersion, ok := pascalCaseMapping[version.Version()]
			if !ok {
				return "", "", nil, fmt.Errorf("No pascal case mapping defined for %q", version.Version())
			}

			for _, kind := range version.Kinds() {
				inputMap := map[string]interface{}{
					"RawAPIVersion":           kind.RawAPIVersion(),
					"Comment":                 kind.Comment(),
					"Group":                   pascalGroup,
					"Kind":                    kind.Kind(),
					"Properties":              kind.Properties(),
					"RequiredInputProperties": kind.RequiredInputProperties(),
					"OptionalInputProperties": kind.OptionalInputProperties(),
					"AdditionalSecretOutputs": kind.AdditionalSecretOutputs(),
					"Aliases":                 kind.Aliases(),
					"URNAPIVersion":           kind.URNAPIVersion(),
					"Version":                 pascalVersion,
					"PulumiComment":           kind.pulumiComment,
				}
				// Since mustache templates are logic-less, we have to add some extra variables
				// to selectively disable code generation for empty lists.
				additionalSecretOutputsPresent := len(kind.AdditionalSecretOutputs()) > 0
				aliasesPresent := len(kind.Aliases()) > 0
				inputMap["MergeOptsRequired"] = additionalSecretOutputsPresent || aliasesPresent
				inputMap["AdditionalSecretOutputsPresent"] = additionalSecretOutputsPresent
				inputMap["AliasesPresent"] = aliasesPresent

				kindCs, err := mustache.RenderFile(
					fmt.Sprintf("%s/kind.cs.mustache", templateDir), inputMap)
				if err != nil {
					return "", "", nil, err
				}

				groups[filepath.Join(group.Group(), version.Version(), kind.Kind()+".cs")] = kindCs
			}
		}
	}
	return
}
