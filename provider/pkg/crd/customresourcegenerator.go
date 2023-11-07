// Copyright 2016-2023, Pulumi Corporation.
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

package crd

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// CustomResourceGenerator generates a Pulumi schema for a single CustomResource
type CustomResourceGenerator struct {
	// CustomResourceDefinition contains the unmarshalled CRD YAML
	CustomResourceDefinition unstructured.Unstructured
	// Schemas represents a mapping from each version in the `spec.versions`
	// list to its corresponding `openAPIV3Schema` field in the CRD YAML
	Schemas map[string]map[string]any
	// ApiVersion represents the `apiVersion` field in the CRD YAML
	APIVersion string
	// Kind represents the `spec.names.kind` field in the CRD YAML
	Kind string
	// Plural represents the `spec.names.plural` field in the CRD YAML
	Plural string
	// Group represents the `spec.group` field in the CRD YAML
	Group string
	// Versions is a slice of names of each version supported by this CRD
	Versions []string
	// GroupVersions is a slice of names of each version, in the format
	// <group>/<version>.
	GroupVersions []string
	// ResourceTokens is a slice of the token types of every versioned
	// CustomResource
	ResourceTokens []string
}

func NewCustomResourceGenerator(crd unstructured.Unstructured) (CustomResourceGenerator, error) {
	apiVersion := crd.GetAPIVersion()
	schemas := map[string]map[string]any{}

	validation, foundValidation, _ := unstructured.NestedMap(crd.Object, "spec", "validation", "openAPIV3Schema")
	if foundValidation { // If present, use the top-level schema to validate all versions
		versionName, foundVersionName, _ := unstructured.NestedString(crd.Object, "spec", "version")
		if foundVersionName {
			schemas[versionName] = validation
		} else if versionInfos, foundVersionInfos, _ := NestedMapSlice(crd.Object, "spec", "versions"); foundVersionInfos {
			for _, versionInfo := range versionInfos {
				versionName, _, _ := unstructured.NestedString(versionInfo, "name")
				schemas[versionName] = validation
			}
		}
	} else { // Otherwise use per-version schemas to validate each version
		versionInfos, foundVersionInfos, _ := NestedMapSlice(crd.Object, "spec", "versions")
		if foundVersionInfos {
			for _, version := range versionInfos {
				name, _, _ := unstructured.NestedString(version, "name")
				if schema, foundSchema, _ := unstructured.NestedMap(version, "schema", "openAPIV3Schema"); foundSchema {
					schemas[name] = schema
				}
			}
		}
	}

	kind, foundKind, _ := unstructured.NestedString(crd.Object, "spec", "names", "kind")
	if !foundKind {
		return CustomResourceGenerator{}, fmt.Errorf("could not find `spec.names.kind` field in the CRD")
	}
	plural, foundPlural, _ := unstructured.NestedString(crd.Object, "spec", "names", "plural")
	if !foundPlural {
		return CustomResourceGenerator{}, fmt.Errorf("could not find `spec.names.plural` field in the CRD")
	}
	group, foundGroup, _ := unstructured.NestedString(crd.Object, "spec", "group")
	if !foundGroup {
		return CustomResourceGenerator{}, fmt.Errorf("could not find `spec.group` field in the CRD")
	}

	versions := make([]string, 0, len(schemas))
	groupVersions := make([]string, 0, len(schemas))
	resourceTokens := make([]string, 0, len(schemas))
	for version := range schemas {
		versions = append(versions, version)
		groupVersions = append(groupVersions, group+"/"+version)
		resourceTokens = append(resourceTokens, getToken(group, version, kind))
	}

	crg := CustomResourceGenerator{
		CustomResourceDefinition: crd,
		Schemas:                  schemas,
		APIVersion:               apiVersion,
		Kind:                     kind,
		Plural:                   plural,
		Group:                    group,
		Versions:                 versions,
		GroupVersions:            groupVersions,
		ResourceTokens:           resourceTokens,
	}

	return crg, nil
}

// HasSchemas returns true if the CustomResource specifies at least some schema, and false otherwise.
func (crg *CustomResourceGenerator) HasSchemas() bool {
	return len(crg.Schemas) > 0
}
