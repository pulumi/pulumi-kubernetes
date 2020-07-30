package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pulumi/pulumi-kubernetes/provider/cmd/crd2pulumi/nodejs"
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
		return CustomResourceGenerator{}, errors.New("invalid language: " + language)
	}

	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return CustomResourceGenerator{}, fmt.Errorf("read file %s: %v", yamlPath, err)
	}

	crd, err := UnmarshalYaml(yamlFile)
	if err != nil {
		return CustomResourceGenerator{}, fmt.Errorf("unmarshal %s: %v", yamlPath, err)
	}

	underscoreFields(crd.Object)

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
			}
			outputFileName := plural + "." + extension
			return path.Join(filepath.Dir(yamlPath), outputFileName)
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

func (gen *CustomResourceGenerator) GetName() string {
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
	objectTypeSpecs := GetObjectTypeSpecs(gen.Versions, gen.GetName(), gen.Kind)
	PrettyPrint(objectTypeSpecs)

	switch gen.Language {
	case NodeJS:
		types, err := nodejs.GenerateTypes(objectTypeSpecs)
		if err != nil {
			return fmt.Errorf("generate code: %v", err)
		}
		classes := gen.GenerateNodeJSClasses()

		file, err := os.Create(gen.OutputPath)
		if err != nil {
			return fmt.Errorf("creating file at %s: %v", gen.OutputPath, err)
		}
		file.WriteString(types)
		file.WriteString(classes)
		defer file.Close()
		return nil
	case DotNet:
	case Python:
	case Go:
	}

	return nil
}

func (gen *CustomResourceGenerator) GenerateNodeJSClasses() string {
	versionNames := gen.GetVersionNames()
	kind := gen.Kind
	group := gen.Group
	name := gen.GetName()

	// Generates a CustomResource sub-class for a single version
	generateResourceClass := func(sb *strings.Builder, version, kind, name, group string) {
		argsName := version + "." + kind
		apiVersion := fmt.Sprintf("%s/%s", group, version)
		sb.WriteString(fmt.Sprintf("export namespace %s {\n", version))
		sb.WriteString(fmt.Sprintf("\texport class %s extends k8s.apiextensions.CustomResource {\n", kind))
		sb.WriteString(fmt.Sprintf("\t\tconstructor(name: string, args?: %s, opts?: pulumi.CustomResourceOptions) {\n", argsName))
		sb.WriteString(fmt.Sprintf("\t\t\tsuper(name, { apiVersion: \"%s\", kind: \"%s\", ...args }, opts)\n", apiVersion, kind))
		sb.WriteString(fmt.Sprintf("\t\t}\n\t}\n}\n\n"))
	}

	// Generates a CustomResourceDefinition class for the entire CRD YAML
	generateDefinitionClass := func(sb *strings.Builder, kind, apiVersion string) {
		className := kind + "Definition"
		var superClassName string
		if apiVersion == v1 {
			superClassName = "k8s.apiextensions.v1.CustomResourceDefinition"
		} else {
			superClassName = "k8s.apiextensions.v1beta1.CustomResourceDefinition"
		}

		sb.WriteString(fmt.Sprintf("export class %s extends %s {\n", className, superClassName))
		sb.WriteString("\tconstructor(name: string, opts?: pulumi.CustomResourceOptions) {\n")
		sb.WriteString("\t\tsuper(name, ")
		definitionArgs, _ := json.MarshalIndent(gen.CustomResourceDefinition.Object, "\t\t", "\t")
		sb.Write(definitionArgs)
		sb.WriteString(", opts)\n\t}\n}\n")
	}

	var sb strings.Builder
	for _, version := range versionNames {
		generateResourceClass(&sb, version, kind, name, group)
	}
	generateDefinitionClass(&sb, kind, gen.APIVersion)

	return sb.String()
}
