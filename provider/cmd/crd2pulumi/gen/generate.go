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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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

type CustomResourceGenerator struct {
	// CustomResourceDefinition represents unmarshalled CRD YAML
	CustomResourceDefinition unstruct.Unstructured
	// OutputDir represents the directory where generated code will output to
	OutputDir string
	// Language represents the target language to generate code
	Language string
	// ApiVersion represents the `apiVersion` field in the CRD YAML
	APIVersion string
	// Kind represents the `spec.names.kind` field in the CRD YAML
	Kind string
	// Plural represents the `spec.names.plural` field in the CRD YAML
	Plural string
	// Group represents the `spec.group` field in the CRD YAML
	Group string
	// Versions represents a mapping from each version name in the
	// `spec.versions` list to its corresponding `openAPIV3Schema` field in the
	// CRD YAML
	Versions map[string]map[string]interface{}
}

func NewCustomResourceGenerator(language, yamlPath, outputDir string) (CustomResourceGenerator, error) {
	if language != DotNet && language != Go && language != NodeJS && language != Python {
		return CustomResourceGenerator{}, errors.New("invalid language " + language)
	}

	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return CustomResourceGenerator{}, errors.Wrapf(err, "reading file %s", yamlPath)
	}

	crd, err := UnmarshalYaml(yamlFile)
	if err != nil {
		return CustomResourceGenerator{}, errors.Wrapf(err, "unmarshalling file %s", yamlPath)
	}

	UnderscoreFields(crd.Object)

	apiVersion := crd.GetAPIVersion()
	if apiVersion != v1beta1 && apiVersion != v1 {
		return CustomResourceGenerator{}, errors.New("invalid apiVersion " + apiVersion)
	}

	versions := map[string]map[string]interface{}{}

	validation, foundValidation, _ := unstruct.NestedMap(crd.Object, "spec", "validation", "openAPIV3Schema")
	if foundValidation { // If present, use the top-level schema to validate all versions
		versionMaps, _, _ := NestedMapSlice(crd.Object, "spec", "versions")
		for _, version := range versionMaps {
			name, _, _ := unstruct.NestedString(version, "name")
			versions[name] = validation
		}
	} else { // Otherwise use per-version schemas to validate each version
		versionMaps, _, _ := NestedMapSlice(crd.Object, "spec", "versions")
		for _, version := range versionMaps {
			name, _, _ := unstruct.NestedString(version, "name")
			schema, _, _ := unstruct.NestedMap(version, "schema", "openAPIV3Schema")
			versions[name] = schema
		}
	}

	kind, _, _ := unstruct.NestedString(crd.Object, "spec", "names", "kind")
	plural, _, _ := unstruct.NestedString(crd.Object, "spec", "names", "plural")
	group, _, _ := unstruct.NestedString(crd.Object, "spec", "group")

	if outputDir == "" {
		outputDir = filepath.Dir(yamlPath)
	} else {
		_, err := os.Stat(outputDir)
		if os.IsNotExist(err) {
			return CustomResourceGenerator{}, errors.Wrapf(err, "output directory does not exist")
		}
	}

	customResourceGenerator := CustomResourceGenerator{
		CustomResourceDefinition: crd,
		OutputDir:                outputDir,
		Language:                 language,
		APIVersion:               apiVersion,
		Kind:                     kind,
		Plural:                   plural,
		Group:                    group,
		Versions:                 versions,
	}
	return customResourceGenerator, nil
}

func (gen *CustomResourceGenerator) Name() string {
	return gen.Plural + "." + gen.Group
}

// VersionKeys returns a slice of the versions supported by this CRD.
func (gen *CustomResourceGenerator) VersionKeys() []string {
	versionNames := make([]string, 0, len(gen.Versions))
	for versionName := range gen.Versions {
		versionNames = append(versionNames, versionName)
	}
	return versionNames
}

// VersionNames returns a slice of the full names of each version, in the format
// <plural>.<group>/<version>.
func (gen *CustomResourceGenerator) VersionNames() []string {
	versions := gen.VersionKeys()
	name := gen.Name()
	for i, version := range versions {
		versions[i] = name + "/" + version
	}
	return versions
}

// getVersion returns the <version> field of a string in the format
// <plural>.<group>/<version>
func getVersion(versionName string) string {
	return strings.Split(versionName, "/")[1]
}

// Generate outputs strongly-typed args for the CustomResourceGenerator's
// CRD in the target language and output folder
func (gen *CustomResourceGenerator) Generate() error {
	switch gen.Language {
	case NodeJS:
		buffer, err := gen.genNodeJS()
		if err != nil {
			return errors.Wrapf(err, "generating nodeJS code")
		}

		outputFile := filepath.Join(gen.OutputDir, gen.Plural+".ts")
		file, err := os.Create(outputFile)
		if err != nil {
			return errors.Wrapf(err, "creating file %s", outputFile)
		}
		defer file.Close()

		_, err = buffer.WriteTo(file)
		if err != nil {
			return errors.Wrapf(err, "writing to %s", outputFile)
		}
	case Go:
		buffers, err := gen.genGo()
		if err != nil {
			return errors.Wrapf(err, "generating Go code")
		}

		for versionName, buffer := range buffers {
			packageDir := filepath.Join(gen.OutputDir, getVersion(versionName))
			if _, err := os.Stat(packageDir); os.IsNotExist(err) {
				err = os.Mkdir(packageDir, 0755)
				if err != nil {
					return errors.Wrapf(err, "creating directory %s", packageDir)
				}
			}
			outputFile := filepath.Join(packageDir, gen.Plural+".go")
			file, err := os.Create(outputFile)
			if err != nil {
				return errors.Wrapf(err, "creating file %s", outputFile)
			}
			defer file.Close()

			_, err = buffer.WriteTo(file)
			if err != nil {
				return errors.Wrapf(err, "writing to %s", outputFile)
			}
		}
	case DotNet:
		fallthrough
	case Python:
		return errors.Errorf("non-supported language %s", gen.Language)
	default:
		contract.Failf("unexpected language %s", gen.Language)
	}

	return nil
}
