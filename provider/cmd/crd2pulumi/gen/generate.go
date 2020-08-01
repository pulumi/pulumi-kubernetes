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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/provider/cmd/crd2pulumi/nodejs"
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
	// OutputPath represents the path of the file to generate code in
	OutputPath string
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

func NewCustomResourceGenerator(language, yamlPath, outputPath string) (CustomResourceGenerator, error) {
	if language != DotNet && language != Go && language != NodeJS && language != Python {
		return CustomResourceGenerator{}, errors.New("invalid language " + language)
	}

	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return CustomResourceGenerator{}, errors.Wrapf(err, "reading file: %s", yamlPath)
	}

	crd, err := UnmarshalYaml(yamlFile)
	if err != nil {
		return CustomResourceGenerator{}, fmt.Errorf("unmarshal %s: %v", yamlPath, err)
	}

	UnderscoreFields(crd.Object)

	apiVersion := crd.GetAPIVersion()
	if apiVersion != v1beta1 && apiVersion != v1 {
		return CustomResourceGenerator{}, fmt.Errorf("invalid apiVersion %s: %v", apiVersion, err)
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

	if outputPath == "" {
		defaultOutputPath := func(yamlPath, plural, language string) string {
			var extension string
			switch language {
			case NodeJS:
				extension = "ts"
			case DotNet:
				extension = "cs"
			case Python:
				extension = "py"
			case Go:
				extension = "go"
			default:
				contract.Failf("unexpected language %s", language)
			}
			outputFileName := plural + "." + extension
			return filepath.Join(filepath.Dir(yamlPath), outputFileName)
		}
		outputPath = defaultOutputPath(yamlPath, plural, language)
	}

	customResourceGenerator := CustomResourceGenerator{
		CustomResourceDefinition: crd,
		OutputPath:               outputPath,
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

// GetVersionNames returns a slice of the versions supported by this CRD.
func (gen *CustomResourceGenerator) GetVersionNames() []string {
	versionNames := make([]string, 0, len(gen.Versions))
	for versionName := range gen.Versions {
		versionNames = append(versionNames, versionName)
	}
	return versionNames
}

func (gen *CustomResourceGenerator) GenerateCode() error {
	objectTypeSpecs := GetObjectTypeSpecs(gen.Versions, gen.Name(), gen.Kind)

	switch gen.Language {
	case NodeJS:
		types, err := nodejs.GenerateTypes(objectTypeSpecs)
		if err != nil {
			return fmt.Errorf("generate code %v", err)
		}
		classes := gen.GenerateNodeJSClasses()

		file, err := os.Create(gen.OutputPath)
		if err != nil {
			return fmt.Errorf("creating file at %s: %v", gen.OutputPath, err)
		}
		defer file.Close()
		file.WriteString(types)
		file.WriteString(classes)
		return nil
	case DotNet:
		fallthrough
	case Python:
		fallthrough
	case Go:
		return errors.Errorf("non-supported language %s", gen.Language)
	default:
		contract.Failf("unexpected language %s", gen.Language)
	}

	return nil
}

func (gen *CustomResourceGenerator) GenerateNodeJSClasses() string {
	versionNames := gen.GetVersionNames()
	kind := gen.Kind
	group := gen.Group
	name := gen.Name()

	// Generates a CustomResource sub-class for a single version
	generateResourceClass := func(w io.Writer, version, kind, name, group string) {
		argsName := version + "." + kind + "Args"
		apiVersion := fmt.Sprintf("%s/%s", group, version)
		fmt.Fprintf(w, "export namespace %s {\n", version)
		fmt.Fprintf(w, "\texport class %s extends k8s.apiextensions.CustomResource {\n", kind)
		fmt.Fprintf(w, "\t\tconstructor(name: string, args?: %s, opts?: pulumi.CustomResourceOptions) {\n", argsName)
		fmt.Fprintf(w, "\t\t\tsuper(name, { apiVersion: \"%s\", kind: \"%s\", ...args }, opts)\n", apiVersion, kind)
		fmt.Fprintf(w, "\t\t}\n\t}\n}\n\n")
	}

	// Generates a CustomResourceDefinition class for the entire CRD YAML
	generateDefinitionClass := func(w io.Writer, kind, apiVersion string) {
		className := kind + "Definition"
		var superClassName string
		if apiVersion == v1 {
			superClassName = "k8s.apiextensions.v1.CustomResourceDefinition"
		} else {
			superClassName = "k8s.apiextensions.v1beta1.CustomResourceDefinition"
		}

		fmt.Fprintf(w, "export class %s extends %s {\n", className, superClassName)
		fmt.Fprintf(w, "\tconstructor(name: string, opts?: pulumi.CustomResourceOptions) {\n")
		fmt.Fprintf(w, "\t\tsuper(name, ")
		definitionArgs, _ := json.MarshalIndent(gen.CustomResourceDefinition.Object, "\t\t", "\t")
		fmt.Fprintf(w, "%s", definitionArgs)
		fmt.Fprintf(w, ", opts)\n\t}\n}\n")
	}

	buffer := &bytes.Buffer{}
	for _, version := range versionNames {
		generateResourceClass(buffer, version, kind, name, group)
	}
	generateDefinitionClass(buffer, kind, gen.APIVersion)

	return buffer.String()
}
