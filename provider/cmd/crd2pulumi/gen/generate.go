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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/pkg/errors"
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

// Generate parses the CRD(s) in the YAML file at the given path and outputs
// code in the given language to `outputDir/crds`. Only overwrites existing
// files if force is true.
func Generate(language, yamlPath, outputDir string, force bool) error {
	outputDir = filepath.Join(outputDir, "crds")
	if !force {
		if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
			return errors.Errorf("%s already exists; use --force to overwrite", outputDir)
		}
	}

	pg, err := NewPackageGenerator(language, yamlPath)
	if err != nil {
		return errors.Wrapf(err, "could not generate %s package for %s", language, yamlPath)
	}

	files, err := pg.GenerateFiles()
	if err != nil {
		return errors.Wrapf(err, "could not generate files for %s", yamlPath)
	}

	if err := writeFiles(files, outputDir); err != nil {
		return errors.Wrap(err, "could not create files/directories")
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

func NewPackageGenerator(language, yamlPath string) (PackageGenerator, error) {
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return PackageGenerator{}, errors.Wrapf(err, "could not read file %s", yamlPath)
	}

	crds, err := UnmarshalYamls(yamlFile)
	if err != nil {
		return PackageGenerator{}, errors.Wrapf(err, "could not unmarshal yaml file %s", yamlPath)
	}

	crgs := make([]CustomResourceGenerator, len(crds))
	for i, crd := range crds {
		crg, err := NewCustomResourceGenerator(crd)
		if err != nil {
			return PackageGenerator{}, errors.Wrapf(err, "could not parse crd %d", i)
		}
		crgs[i] = crg
	}

	return PackageGenerator{
		CustomResourceGenerators: crgs,
		Language:                 language,
	}, nil
}

func (pg *PackageGenerator) BaseRefs() []string {
	var baseRefs []string
	for _, crg := range pg.CustomResourceGenerators {
		baseRefs = append(baseRefs, crg.baseRefs()...)
	}
	return baseRefs
}

func (pg *PackageGenerator) GroupVersions() []string {
	var groupVersions []string
	for _, crg := range pg.CustomResourceGenerators {
		groupVersions = append(groupVersions, crg.GroupVersions()...)
	}
	return groupVersions
}

func (pg *PackageGenerator) GenerateFiles() (map[string]*bytes.Buffer, error) {
	types := pg.GetTypes()
	baseRefs := pg.BaseRefs()

	var files map[string]*bytes.Buffer
	var err error

	switch pg.Language {
	case NodeJS:
		files, err = pg.genNodeJS(types, baseRefs)
	case Go:
		files, err = pg.genGo(types, baseRefs)
	case Python:
		return nil, errors.Errorf("unsupported language %s", pg.Language)
	case DotNet:
		files, err = pg.genDotNet(types, baseRefs)
	default:
		contract.Failf("unexpected language %s", pg.Language)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "could not generate files for %s", pg.Language)
	}
	return files, nil
}

// Returns the 'moduleToPackage' language-specific map. Creates a mapping from
// every groupVersion string (<group>/<version>) to <groupPrefix>/<version>.
func (pg *PackageGenerator) moduleToPackage() map[string]string {
	moduleToPackage := map[string]string{}
	for _, groupVersion := range pg.GroupVersions() {
		group, version := splitGroupVersion(groupVersion)
		moduleToPackage[groupVersion] = groupPrefix(group) + "/" + version
	}
	return moduleToPackage
}

// PackageGenerator generates code for multiple CustomResources
type PackageGenerator struct {
	// CustomResourceGenerators contains a slice of all CustomResourceGenerators
	CustomResourceGenerators []CustomResourceGenerator
	// Language represents the target language to generate code
	Language string
}

// CustomResourceGenerator generates the Pulumi schema for a single CustomResource
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
}

func NewCustomResourceGenerator(crd unstruct.Unstructured) (CustomResourceGenerator, error) {
	apiVersion := crd.GetAPIVersion()
	if !IsValidAPIVersion(apiVersion) {
		return CustomResourceGenerator{},
			errors.Errorf("invalid apiVersion %s, only v1 and v1beta1 are supported", apiVersion)
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

	crg := CustomResourceGenerator{
		CustomResourceDefinition: crd,
		Schemas:                  schemas,
		APIVersion:               apiVersion,
		Kind:                     kind,
		Plural:                   plural,
		Group:                    group,
	}

	return crg, nil
}

// Versions returns a slice of the versions supported by this CRD.
func (gen *CustomResourceGenerator) Versions() []string {
	versions := make([]string, 0, len(gen.Schemas))
	for version := range gen.Schemas {
		versions = append(versions, version)
	}
	return versions
}

// GroupVersions returns a slice of the names of each version, in the format
// <group>/<version>.
func (gen *CustomResourceGenerator) GroupVersions() []string {
	versions := gen.Versions()
	for i, version := range versions {
		versions[i] = gen.Group + "/" + version
	}
	return versions
}

func IsValidAPIVersion(apiVersion string) bool {
	return apiVersion == v1 || apiVersion == v1beta1
}

func IsValidLanguage(language string) bool {
	return language == NodeJS || language == Go || language == Python || language == DotNet
}

// splitGroupVersion returns the <group> and <version> field of a string in the
// format <group>/<version>
func splitGroupVersion(groupVersion string) (string, string) {
	parts := strings.Split(groupVersion, "/")
	contract.Assert(len(parts) == 2)
	return parts[0], parts[1]
}

// groupPrefix returns the first word in the dot-seperated group string, with
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
		if unicode.IsLetter(r) && (i == 0 || !unicode.IsLetter(rune(version[i]))) {
			sb.WriteRune(unicode.ToUpper(r))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
