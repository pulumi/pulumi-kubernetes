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

// nolint: nakedret
package gen

import (
	"fmt"
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
	"flowcontrol":           "FlowControl",
	"networking":            "Networking",
	"meta":                  "Meta",
	"node":                  "Node",
	"policy":                "Policy",
	"rbac":                  "Rbac",
	"resource":              "Resource",
	"scheduling":            "Scheduling",
	"settings":              "Settings",
	"storage":               "Storage",
	"storagemigration":      "StorageMigration",
	"v1":                    "V1",
	"v1alpha1":              "V1Alpha1",
	"v1alpha2":              "V1Alpha2",
	"v1alpha3":              "V1Alpha3",
	"v1beta1":               "V1Beta1",
	"v1beta2":               "V1Beta2",
	"v1beta3":               "V1Beta3",
	"v2":                    "V2",
	"v2alpha1":              "V2Alpha1",
	"v2beta1":               "V2Beta1",
	"v2beta2":               "V2Beta2",

	// Not sure what these are - but they show up in input and output types.
	"version": "Version",
	"pkg":     "Pkg",
}

func pascalCase(name string) string {
	pascal, ok := pascalCaseMapping[name]
	if !ok {
		panic(fmt.Sprintf("no case mapping for %q", name))
	}
	return pascal
}
