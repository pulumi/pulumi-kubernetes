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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	pschema "github.com/pulumi/pulumi/pkg/v2/codegen/schema"
	"github.com/pulumi/pulumi/sdk/v2/go/common/util/contract"
	unstruct "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	DotNet string = "dotnet"
	Go     string = "go"
	NodeJS string = "nodejs"
	Python string = "python"
)

const (
	v1beta1 string = "apiextensions.k8s.io/v1beta1"
	v1      string = "apiextensions.k8s.io/v1"
)

// Version specifies the crd2pulumi version. It should be set by the linker via LDFLAGS.
var Version string

// LanguageSettings contains the output paths for each language. If a path is
// null, then that language will not be generated at all.
type LanguageSettings struct {
	NodeJSPath *string
	PythonPath *string
	DotNetPath *string
	GoPath     *string
}

// Returns true if at least one of the language-specific output paths already exists. If true, then a slice of the
// paths that already exist are also returned.
func (ls LanguageSettings) hasExistingPaths() (bool, []string) {
	pathExists := func(path string) bool {
		_, err := os.Stat(path)
		return !os.IsNotExist(err)
	}
	var existingPaths []string
	if ls.NodeJSPath != nil && pathExists(*ls.NodeJSPath) {
		existingPaths = append(existingPaths, *ls.NodeJSPath)
	}
	if ls.PythonPath != nil && pathExists(*ls.PythonPath) {
		existingPaths = append(existingPaths, *ls.PythonPath)
	}
	if ls.DotNetPath != nil && pathExists(*ls.DotNetPath) {
		existingPaths = append(existingPaths, *ls.DotNetPath)
	}
	if ls.GoPath != nil && pathExists(*ls.GoPath) {
		existingPaths = append(existingPaths, *ls.GoPath)
	}
	return len(existingPaths) > 0, existingPaths
}

// Generate parses the CRDs at the given yamlPaths and outputs the generated
// code according to the language settings. Only overwrites existing files if
// force is true.
func Generate(ls LanguageSettings, yamlPaths []string, force bool) error {
	if !force {
		if exists, paths := ls.hasExistingPaths(); exists {
			return errors.Errorf("path(s) %s already exists; use --force to overwrite", paths)
		}
	}

	pg, err := NewPackageGenerator(yamlPaths)
	if err != nil {
		return err
	}

	if ls.NodeJSPath != nil {
		if err := pg.genNodeJS(*ls.NodeJSPath); err != nil {
			return err
		}
	}
	if ls.PythonPath != nil {
		if err := pg.genPython(*ls.PythonPath); err != nil {
			return err
		}
	}
	if ls.GoPath != nil {
		if err := pg.genGo(*ls.GoPath); err != nil {
			return err
		}
	}
	if ls.DotNetPath != nil {
		if err := pg.genDotNet(*ls.DotNetPath); err != nil {
			return err
		}
	}

	return nil
}

// Writes the contents of each buffer to its file path, relative to `outputDir`.
// `files` should be a mapping from file path strings to buffers.
func writeFiles(files map[string]*bytes.Buffer, outputDir string) error {
	for path, code := range files {
		outputFilePath := filepath.Join(outputDir, path)
		err := os.MkdirAll(filepath.Dir(outputFilePath), 0755)
		if err != nil {
			return errors.Wrapf(err, "could not create directory to %s", outputFilePath)
		}
		file, err := os.Create(outputFilePath)
		if err != nil {
			return errors.Wrapf(err, "could not create file %s", outputFilePath)
		}
		defer file.Close()
		if _, err := code.WriteTo(file); err != nil {
			return errors.Wrapf(err, "could not write to file %s", outputFilePath)
		}
	}
	return nil
}

// PackageGenerator generates code for multiple CustomResources
type PackageGenerator struct {
	// CustomResourceGenerators contains a slice of all CustomResourceGenerators
	CustomResourceGenerators []CustomResourceGenerator
	// ResourceTokens is a slice of the token types of every CustomResource
	ResourceTokens []string
	// GroupVersions is a slice of the names of every CustomResource's versions,
	// in the format <group>/<version>
	GroupVersions []string
	// Types is a mapping from every type's token name to its ObjectTypeSpec
	Types map[string]pschema.ObjectTypeSpec
	// schemaPackage is the Pulumi schema package used to generate code for
	// languages that do not need an ObjectMeta type (NodeJS)
	schemaPackage *pschema.Package
	// schemaPackageWithObjectMetaType is the Pulumi schema package used to
	// generate code for languages that need an ObjectMeta type (Python, Go, and .NET)
	schemaPackageWithObjectMetaType *pschema.Package
}

func NewPackageGenerator(yamlPaths []string) (PackageGenerator, error) {
	yamlFiles := make([][]byte, 0, len(yamlPaths))

	for _, yamlPath := range yamlPaths {
		yamlFile, err := ioutil.ReadFile(yamlPath)
		if err != nil {
			return PackageGenerator{}, errors.Wrapf(err, "could not read file %s", yamlPath)
		}
		yamlFiles = append(yamlFiles, yamlFile)
	}

	crds, err := UnmarshalYamls(yamlFiles)
	if err != nil {
		return PackageGenerator{}, errors.Wrapf(err, "could not unmarshal yaml file(s)")
	}

	resourceTokensSize := 0
	groupVersionsSize := 0

	crgs := make([]CustomResourceGenerator, 0, len(crds))
	for i, crd := range crds {
		crg, err := NewCustomResourceGenerator(crd)
		if err != nil {
			return PackageGenerator{}, errors.Wrapf(err, "could not parse crd %d", i)
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

	pg := PackageGenerator{
		CustomResourceGenerators: crgs,
		ResourceTokens:           baseRefs,
		GroupVersions:            groupVersions,
	}
	pg.Types = pg.GetTypes()
	return pg, nil
}

// SchemaPackage returns the Pulumi schema package with no ObjectMeta type.
// This is only necessary for NodeJS and Python.
func (pg *PackageGenerator) SchemaPackage() *pschema.Package {
	if pg.schemaPackage == nil {
		pkg, err := genPackage(pg.Types, pg.ResourceTokens, false)
		contract.AssertNoErrorf(err, "could not parse Pulumi package")
		pg.schemaPackage = pkg
	}
	return pg.schemaPackage
}

// SchemaPackageWithObjectMetaType returns the Pulumi schema package with
// an ObjectMeta type. This is only necessary for Go and .NET.
func (pg *PackageGenerator) SchemaPackageWithObjectMetaType() *pschema.Package {
	if pg.schemaPackageWithObjectMetaType == nil {
		pkg, err := genPackage(pg.Types, pg.ResourceTokens, true)
		contract.AssertNoErrorf(err, "could not parse Pulumi package")
		pg.schemaPackageWithObjectMetaType = pkg
	}
	return pg.schemaPackageWithObjectMetaType
}

// Returns language-specific 'moduleToPackage' map. Creates a mapping from
// every groupVersion string <group>/<version> to <groupPrefix>/<version>.
func (pg *PackageGenerator) moduleToPackage() map[string]string {
	moduleToPackage := map[string]string{}
	for _, groupVersion := range pg.GroupVersions {
		group, version := splitGroupVersion(groupVersion)
		moduleToPackage[groupVersion] = groupPrefix(group) + "/" + version
	}
	return moduleToPackage
}

// CustomResourceGenerator generates a Pulumi schema for a single CustomResource
type CustomResourceGenerator struct {
	// CustomResourceDefinition contains the unmarshalled CRD YAML
	CustomResourceDefinition unstruct.Unstructured
	// Schemas represents a mapping from each version in the `spec.versions`
	// list to its corresponding `openAPIV3Schema` field in the CRD YAML
	Schemas map[string]map[string]interface{}
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

func NewCustomResourceGenerator(crd unstruct.Unstructured) (CustomResourceGenerator, error) {
	apiVersion := crd.GetAPIVersion()
	if !IsValidAPIVersion(apiVersion) {
		return CustomResourceGenerator{},
			errors.Errorf("invalid apiVersion %s; only v1 and v1beta1 are supported", apiVersion)
	}

	schemas := map[string]map[string]interface{}{}
	validation, foundValidation, _ := unstruct.NestedMap(crd.Object, "spec", "validation", "openAPIV3Schema")
	if foundValidation { // If present, use the top-level schema to validate all versions
		versionMaps, _, _ := NestedMapSlice(crd.Object, "spec", "versions")
		for _, version := range versionMaps {
			name, _, _ := unstruct.NestedString(version, "name")
			schemas[name] = validation
		}
	} else { // Otherwise use per-version schemas to validate each version
		versionMaps, _, _ := NestedMapSlice(crd.Object, "spec", "versions")
		for _, version := range versionMaps {
			name, _, _ := unstruct.NestedString(version, "name")
			schema, _, _ := unstruct.NestedMap(version, "schema", "openAPIV3Schema")
			schemas[name] = schema
		}
	}

	kind, foundKind, _ := unstruct.NestedString(crd.Object, "spec", "names", "kind")
	if !foundKind {
		return CustomResourceGenerator{}, errors.New("could not find `spec.names.kind` field in the CRD")
	}
	plural, foundPlural, _ := unstruct.NestedString(crd.Object, "spec", "names", "plural")
	if !foundPlural {
		return CustomResourceGenerator{}, errors.New("could not find `spec.names.plural` field in the CRD")
	}
	group, foundGroup, _ := unstruct.NestedString(crd.Object, "spec", "group")
	if !foundGroup {
		return CustomResourceGenerator{}, errors.New("could not find `spec.group` field in the CRD")
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

// Returns the type token for a Kubernetes CustomResource with the given group,
// version, and kind.
func getToken(group, version, kind string) string {
	return fmt.Sprintf("kubernetes:%s/%s:%s", group, version, kind)
}

// IsValidAPIVersion returns true if and only if the given apiVersion is
// supported (apiextensions.k8s.io/v1beta1 or apiextensions.k8s.io/v1).
func IsValidAPIVersion(apiVersion string) bool {
	return apiVersion == v1 || apiVersion == v1beta1
}

// splitGroupVersion returns the <group> and <version> field of a string in the
// format <group>/<version>
func splitGroupVersion(groupVersion string) (string, string) {
	parts := strings.Split(groupVersion, "/")
	contract.Assert(len(parts) == 2)
	return parts[0], parts[1]
}

// groupPrefix returns the first word in the dot-separated group string, with
// all non-alphanumeric characters removed.
func groupPrefix(group string) string {
	contract.Assert(group != "")
	return removeNonAlphanumeric(strings.Split(group, ".")[0])
}

// Capitalizes and returns the given version. For example,
// versionToUpper("v2beta1") returns "V2Beta1".
func versionToUpper(version string) string {
	var sb strings.Builder
	for i, r := range version {
		if unicode.IsLetter(r) && (i == 0 || !unicode.IsLetter(rune(version[i-1]))) {
			sb.WriteRune(unicode.ToUpper(r))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
