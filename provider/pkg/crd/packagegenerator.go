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
	"io"

	pschema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// PackageGenerator generates code for multiple CustomResources
type PackageGenerator struct {
	// CustomResourceGenerators contains a slice of all CustomResourceGenerators
	CustomResourceGenerators []CustomResourceGenerator
	// ResourceTokens is a slice of the token types of every CustomResource
	ResourceTokens []string
	// GroupVersions is a slice of the names of every CustomResource's versions,
	// in the format <group>/<version>
	GroupVersions []string
	// Types is a mapping from every type's token name to its ComplexTypeSpec
	Types map[string]pschema.ComplexTypeSpec
	// Version is the semver that will be stamped into the generated package
	Version string
	// schemaPackageWithObjectMetaType is the Pulumi schema package used to
	// generate code for languages that need an ObjectMeta type (Python, Go, and .NET)
	schemaPackageWithObjectMetaType *pschema.Package
}

// ReadPackagesFromSource reads one or more documents and returns a PackageGenerator that can be used to generate Pulumi code.
// Calling this function will fully read and close each document.
func ReadPackagesFromSource(version string, yamlSources []io.ReadCloser) (*PackageGenerator, error) {
	yamlData := make([][]byte, len(yamlSources))

	for i, yamlSource := range yamlSources {
		defer yamlSource.Close()
		var err error
		yamlData[i], err = io.ReadAll(yamlSource)
		if err != nil {
			return nil, fmt.Errorf("failed to read YAML: %w", err)
		}
	}

	crds, err := UnmarshalYamls(yamlData)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal yaml file(s): %w", err)
	}

	if len(crds) == 0 {
		return nil, fmt.Errorf("could not find any CRD YAML files")
	}

	resourceTokensSize := 0
	groupVersionsSize := 0

	crgs := make([]CustomResourceGenerator, 0, len(crds))
	for i, crd := range crds {
		crg, err := NewCustomResourceGenerator(crd)
		if err != nil {
			return nil, fmt.Errorf("could not parse crd %d: %w", i, err)
		}
		resourceTokensSize += len(crg.ResourceTokens)
		groupVersionsSize += len(crg.GroupVersions)
		crgs = append(crgs, crg)
	}

	baseRefs := make([]string, 0, resourceTokensSize)
	groupVersions := make([]string, 0, groupVersionsSize)
	for _, crg := range crgs {
		baseRefs = append(baseRefs, crg.ResourceTokens...)
		groupVersions = append(groupVersions, crg.GroupVersions...)
	}

	pg := &PackageGenerator{
		CustomResourceGenerators: crgs,
		ResourceTokens:           baseRefs,
		GroupVersions:            groupVersions,
		Version:                  version,
	}
	pg.Types = pg.GetTypes()
	return pg, nil
}

func (pg *PackageGenerator) SchemaPackage() *pschema.Package {
	if pg.schemaPackageWithObjectMetaType == nil {
		pkg, err := genPackage(pg.Version, pg.Types, pg.ResourceTokens)
		contract.AssertNoErrorf(err, "could not parse Pulumi package")
		pg.schemaPackageWithObjectMetaType = pkg
	}
	return pg.schemaPackageWithObjectMetaType
}

// Returns language-specific 'ModuleToPackage' map. Creates a mapping from
// every groupVersion string <group>/<version> to <groupPrefix>/<version>.
func (pg *PackageGenerator) ModuleToPackage() (map[string]string, error) {
	moduleToPackage := map[string]string{}
	for _, groupVersion := range pg.GroupVersions {
		group, version, err := SplitGroupVersion(groupVersion)
		if err != nil {
			return nil, fmt.Errorf("invalid version: %w", err)
		}
		prefix, err := GroupPrefix(group)
		if err != nil {
			return nil, fmt.Errorf("invalid version: %w", err)
		}
		moduleToPackage[groupVersion] = prefix + "/" + version
	}
	return moduleToPackage, nil
}

// HasSchemas returns true if there exists at least one CustomResource with a schema in this package.
func (pg *PackageGenerator) HasSchemas() bool {
	for _, crg := range pg.CustomResourceGenerators {
		if crg.HasSchemas() {
			return true
		}
	}
	return false
}

func (pg *PackageGenerator) GetTypes() map[string]pschema.ComplexTypeSpec {
	types := map[string]pschema.ComplexTypeSpec{}
	for _, crg := range pg.CustomResourceGenerators {
		for version, schema := range crg.Schemas {
			resourceToken := getToken(crg.Group, version, crg.Kind)
			_, foundProperties, _ := unstructured.NestedMap(schema, "properties")
			if foundProperties {
				AddType(schema, resourceToken, types)
			}
			preserveUnknownFields, _, _ := unstructured.NestedBool(schema, "x-kubernetes-preserve-unknown-fields")
			if preserveUnknownFields {
				types[resourceToken] = emptySpec
			}
			if foundProperties || preserveUnknownFields {
				types[resourceToken].Properties["apiVersion"] = pschema.PropertySpec{
					TypeSpec: pschema.TypeSpec{
						Type: String,
					},
					Const: crg.Group + "/" + version,
				}
				types[resourceToken].Properties["kind"] = pschema.PropertySpec{
					TypeSpec: pschema.TypeSpec{
						Type: String,
					},
					Const: crg.Kind,
				}
				types[resourceToken].Properties["metadata"] = pschema.PropertySpec{
					TypeSpec: pschema.TypeSpec{
						Ref: objectMetaRef,
					},
				}
			}
		}
	}
	return types
}
