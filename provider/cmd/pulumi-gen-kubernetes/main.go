// Copyright 2016-2018, Pulumi Corporation.
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

// nolint:gosec
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/gen"
	providerVersion "github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/version"
	"github.com/pulumi/pulumi/pkg/v3/codegen"
	dotnetgen "github.com/pulumi/pulumi/pkg/v3/codegen/dotnet"
	gogen "github.com/pulumi/pulumi/pkg/v3/codegen/go"
	nodejsgen "github.com/pulumi/pulumi/pkg/v3/codegen/nodejs"
	pythongen "github.com/pulumi/pulumi/pkg/v3/codegen/python"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
)

// This is the URL for the v1.17.0 swagger spec. This is the last version of the spec containing the following
// deprecated resources:
// - extensions/v1beta1/*
// - apps/v1beta1/*
// - apps/v1beta2/*
// Since these resources will continue to be important to users for the foreseeable future, we will merge in
// newer specs on top of this spec so that these resources continue to be available in our SDKs.
const Swagger117Url = "https://raw.githubusercontent.com/kubernetes/kubernetes/v1.17.0/api/openapi-spec/swagger.json"
const Swagger117FileName = "swagger-v1.17.0.json"

// TemplateDir is the path to the base directory for code generator templates.
var TemplateDir string

// BaseDir is the path to the base pulumi-kubernetes directory.
var BaseDir string

// Language is the SDK language.
type Language string

const (
	DotNet Language = "dotnet"
	Go     Language = "go"
	NodeJS Language = "nodejs"
	Kinds  Language = "kinds"
	Python Language = "python"
	Schema Language = "schema"
)

func main() {
	flag.Usage = func() {
		const usageFormat = "Usage: %s <language> <swagger-or-schema-file> <root-pulumi-kubernetes-dir>"
		_, err := fmt.Fprintf(flag.CommandLine.Output(), usageFormat, os.Args[0])
		contract.IgnoreError(err)
		flag.PrintDefaults()
	}

	var version string
	flag.StringVar(&version, "version", providerVersion.Version, "the provider version to record in the generated code")

	flag.Parse()
	args := flag.Args()
	if len(args) < 3 {
		flag.Usage()
		return
	}

	language, inputFile := Language(args[0]), args[1]

	BaseDir = args[2]
	TemplateDir = filepath.Join(BaseDir, "provider", "pkg", "gen")
	outdir := filepath.Join(BaseDir, "sdk", string(language))

	switch language {
	case NodeJS:
		templateDir := filepath.Join(TemplateDir, "nodejs-templates")
		writeNodeJSClient(readSchema(inputFile, version), outdir, templateDir)
	case Python:
		templateDir := filepath.Join(TemplateDir, "python-templates")
		writePythonClient(readSchema(inputFile, version), outdir, templateDir)
	case DotNet:
		templateDir := filepath.Join(TemplateDir, "dotnet-templates")
		writeDotnetClient(readSchema(inputFile, version), outdir, templateDir)
	case Go:
		templateDir := filepath.Join(TemplateDir, "_go-templates")
		writeGoClient(readSchema(inputFile, version), outdir, templateDir)
	case Kinds:
		pkg := readSchema(inputFile, version)
		genK8sResourceTypes(pkg)
	case Schema:
		pkgSpec := generateSchema(inputFile)
		mustWritePulumiSchema(pkgSpec, version)
	default:
		panic(fmt.Sprintf("Unrecognized language '%s'", language))
	}
}

func readSchema(schemaPath string, version string) *schema.Package {
	// Read in, decode, and import the schema.
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		panic(err)
	}

	var pkgSpec schema.PackageSpec
	if err = json.Unmarshal(schemaBytes, &pkgSpec); err != nil {
		panic(err)
	}
	pkgSpec.Version = version

	pkg, err := schema.ImportSpec(pkgSpec, nil)
	if err != nil {
		panic(err)
	}
	return pkg
}

func generateSchema(swaggerPath string) schema.PackageSpec {
	swagger, err := ioutil.ReadFile(swaggerPath)
	if err != nil {
		panic(err)
	}

	swaggerDir := filepath.Dir(swaggerPath)

	legacySwaggerPath := filepath.Join(swaggerDir, Swagger117FileName)
	err = DownloadFile(legacySwaggerPath, Swagger117Url)
	if err != nil {
		panic(err)
	}
	legacySwagger, err := ioutil.ReadFile(legacySwaggerPath)
	if err != nil {
		panic(err)
	}
	mergedSwagger := mergeSwaggerSpecs(legacySwagger, swagger)
	data := mergedSwagger.(map[string]interface{})

	// Generate schema
	return gen.PulumiSchema(data)
}

func writeNodeJSClient(pkg *schema.Package, outdir, templateDir string) {
	resources, err := nodejsgen.LanguageResources(pkg)
	if err != nil {
		panic(err)
	}

	templateResources := gen.TemplateResources{}
	packages := codegen.StringSet{}
	for _, resource := range resources {
		if resource.Package == "" {
			continue
		}
		tr := gen.TemplateResource{
			Name:    resource.Name,
			Package: resource.Package,
			Token:   resource.Token,
		}
		for _, property := range resource.Properties {
			tp := gen.TemplateProperty{
				ConstValue: property.ConstValue,
				Name:       property.Name,
				Package:    property.Package,
			}
			tr.Properties = append(tr.Properties, tp)
		}
		templateResources.Resources = append(templateResources.Resources, tr)
		groupPackage := strings.Split(resource.Package, ".")[0]
		packages.Add(groupPackage)
	}
	sort.Slice(templateResources.Resources, func(i, j int) bool {
		return templateResources.Resources[i].Token < templateResources.Resources[j].Token
	})
	templateResources.Packages = packages.SortedValues()

	overlays := map[string][]byte{
		"apiextensions/customResource.ts": mustLoadFile(filepath.Join(templateDir, "apiextensions", "customResource.ts")),
		"helm/v2/helm.ts":                 mustLoadFile(filepath.Join(templateDir, "helm", "v2", "helm.ts")),
		"helm/v3/helm.ts":                 mustLoadFile(filepath.Join(templateDir, "helm", "v3", "helm.ts")),
		"kustomize/kustomize.ts":          mustLoadFile(filepath.Join(templateDir, "kustomize", "kustomize.ts")),
		"yaml/yaml.ts":                    mustRenderTemplate(filepath.Join(templateDir, "yaml", "yaml.tmpl"), templateResources),
	}
	files, err := nodejsgen.GeneratePackage("pulumigen", pkg, overlays)
	if err != nil {
		panic(err)
	}

	// Internal files that don't need to be exported
	files["path.ts"] = mustLoadFile(filepath.Join(templateDir, "path.ts"))
	files["tests/path.ts"] = mustLoadFile(filepath.Join(templateDir, "tests", "path.ts"))

	mustWriteFiles(outdir, files)
}

func writePythonClient(pkg *schema.Package, outdir string, templateDir string) {
	resources, err := pythongen.LanguageResources("pulumigen", pkg)
	if err != nil {
		panic(err)
	}

	templateResources := gen.TemplateResources{}
	for _, resource := range resources {
		r := gen.TemplateResource{
			Name:    resource.Name,
			Package: resource.Package,
			Token:   resource.Token,
		}
		templateResources.Resources = append(templateResources.Resources, r)
	}
	sort.Slice(templateResources.Resources, func(i, j int) bool {
		return templateResources.Resources[i].Token < templateResources.Resources[j].Token
	})

	overlays := map[string][]byte{
		"apiextensions/CustomResource.py": mustLoadFile(filepath.Join(templateDir, "apiextensions", "CustomResource.py")),
		"helm/v2/helm.py":                 mustLoadFile(filepath.Join(templateDir, "helm", "v2", "helm.py")),
		"helm/v3/helm.py":                 mustLoadFile(filepath.Join(templateDir, "helm", "v3", "helm.py")),
		"kustomize.py":                    mustLoadFile(filepath.Join(templateDir, "kustomize", "kustomize.py")),
		"yaml.py":                         mustRenderTemplate(filepath.Join(templateDir, "yaml", "yaml.tmpl"), templateResources),
	}

	files, err := pythongen.GeneratePackage("pulumigen", pkg, overlays)
	if err != nil {
		panic(err)
	}

	mustWriteFiles(outdir, files)
}

func writeDotnetClient(pkg *schema.Package, outdir, templateDir string) {
	resources, err := dotnetgen.LanguageResources("pulumigen", pkg)
	if err != nil {
		panic(err)
	}

	templateResources := gen.TemplateResources{}
	for _, resource := range resources {
		r := gen.TemplateResource{
			Name:    resource.Name,
			Package: resource.Package,
			Token:   resource.Token,
		}
		templateResources.Resources = append(templateResources.Resources, r)
	}
	sort.Slice(templateResources.Resources, func(i, j int) bool {
		return templateResources.Resources[i].Token < templateResources.Resources[j].Token
	})
	overlays := map[string][]byte{
		"ApiExtensions/CustomResource.cs": mustLoadFile(filepath.Join(templateDir, "apiextensions", "CustomResource.cs")),
		"Helm/ChartBase.cs":               mustLoadFile(filepath.Join(templateDir, "helm", "ChartBase.cs")),
		"Helm/Unwraps.cs":                 mustLoadFile(filepath.Join(templateDir, "helm", "Unwraps.cs")),
		"Helm/V2/Chart.cs":                mustLoadFile(filepath.Join(templateDir, "helm", "v2", "Chart.cs")),
		"Helm/V3/Chart.cs":                mustLoadFile(filepath.Join(templateDir, "helm", "v3", "Chart.cs")),
		"Helm/V3/Invokes.cs":              mustLoadFile(filepath.Join(templateDir, "helm", "v3", "Invokes.cs")),
		"Kustomize/Directory.cs":          mustLoadFile(filepath.Join(templateDir, "kustomize", "Directory.cs")),
		"Kustomize/Invokes.cs":            mustLoadFile(filepath.Join(templateDir, "kustomize", "Invokes.cs")),
		"Yaml/ConfigFile.cs":              mustLoadFile(filepath.Join(templateDir, "yaml", "ConfigFile.cs")),
		"Yaml/ConfigGroup.cs":             mustLoadFile(filepath.Join(templateDir, "yaml", "ConfigGroup.cs")),
		"Yaml/Invokes.cs":                 mustLoadFile(filepath.Join(templateDir, "yaml", "Invokes.cs")),
		"Yaml/TransformationAction.cs":    mustLoadFile(filepath.Join(templateDir, "yaml", "TransformationAction.cs")),
		"Yaml/Yaml.cs":                    mustRenderTemplate(filepath.Join(templateDir, "yaml", "yaml.tmpl"), templateResources),
	}

	files, err := dotnetgen.GeneratePackage("pulumigen", pkg, overlays)
	if err != nil {
		panic(err)
	}

	for filename, contents := range files {
		path := filepath.Join(outdir, filename)

		if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			panic(err)
		}
		err := ioutil.WriteFile(path, contents, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func writeGoClient(pkg *schema.Package, outdir string, templateDir string) {
	files, err := gogen.GeneratePackage("pulumigen", pkg)
	if err != nil {
		panic(err)
	}

	resources, err := gogen.LanguageResources("pulumigen", pkg)
	if err != nil {
		panic(err)
	}

	templateResources := gen.GoTemplateResources{}
	for _, resource := range resources {
		r := gen.TemplateResource{
			Alias:   resource.Alias,
			Name:    resource.Name,
			Package: resource.Package,
			Token:   resource.Token,
		}
		templateResources.Resources = append(templateResources.Resources, r)
	}
	sort.Slice(templateResources.Resources, func(i, j int) bool {
		return templateResources.Resources[i].Token < templateResources.Resources[j].Token
	})

	files["kubernetes/pulumiTypes.go"] = mustLoadGoFile(filepath.Join(templateDir, "pulumiTypes.go"))
	files["kubernetes/apiextensions/customResource.go"] = mustLoadGoFile(filepath.Join(templateDir, "apiextensions", "customResource.go"))
	files["kubernetes/helm/v2/chart.go"] = mustLoadGoFile(filepath.Join(templateDir, "helm", "v2", "chart.go"))
	files["kubernetes/helm/v2/pulumiTypes.go"] = mustLoadGoFile(filepath.Join(templateDir, "helm", "v2", "pulumiTypes.go"))
	files["kubernetes/helm/v3/chart.go"] = mustLoadGoFile(filepath.Join(templateDir, "helm", "v3", "chart.go"))
	files["kubernetes/helm/v3/pulumiTypes.go"] = mustLoadGoFile(filepath.Join(templateDir, "helm", "v3", "pulumiTypes.go"))
	files["kubernetes/kustomize/directory.go"] = mustLoadGoFile(filepath.Join(templateDir, "kustomize", "directory.go"))
	files["kubernetes/kustomize/pulumiTypes.go"] = mustLoadGoFile(filepath.Join(templateDir, "kustomize", "pulumiTypes.go"))
	files["kubernetes/yaml/configFile.go"] = mustLoadGoFile(filepath.Join(templateDir, "yaml", "configFile.go"))
	files["kubernetes/yaml/configGroup.go"] = mustLoadGoFile(filepath.Join(templateDir, "yaml", "configGroup.go"))
	files["kubernetes/yaml/transformation.go"] = mustLoadGoFile(filepath.Join(templateDir, "yaml", "transformation.go"))
	files["kubernetes/yaml/yaml.go"] = mustRenderGoTemplate(filepath.Join(templateDir, "yaml", "yaml.tmpl"), templateResources)

	mustWriteFiles(outdir, files)
}

func mustLoadFile(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return b
}

func mustLoadGoFile(path string) []byte {
	b := mustLoadFile(path)

	formattedSource, err := format.Source(b)
	if err != nil {
		panic(err)
	}
	return formattedSource
}

func mustRenderTemplate(path string, resources interface{}) []byte {
	b := mustLoadFile(path)
	t := template.Must(template.New("resources").Parse(string(b)))

	var buf bytes.Buffer
	err := t.Execute(&buf, resources)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func mustRenderGoTemplate(path string, resources interface{}) []byte {
	bytes := mustRenderTemplate(path, resources)

	formattedSource, err := format.Source(bytes)
	if err != nil {
		panic(err)
	}
	return formattedSource
}

func genK8sResourceTypes(pkg *schema.Package) {
	groupVersions, kinds := codegen.NewStringSet(), codegen.NewStringSet()
	for _, resource := range pkg.Resources {
		parts := strings.Split(resource.Token, ":")
		contract.Assert(len(parts) == 3)

		groupVersion, kind := parts[1], parts[2]
		groupVersions.Add(groupVersion)
		kinds.Add(kind)
	}

	gvk := gen.GVK{Kinds: kinds.SortedValues()}
	gvStrings := groupVersions.SortedValues()
	for _, gvString := range gvStrings {
		gvk.GroupVersions = append(gvk.GroupVersions, gen.GroupVersion(gvString))
	}

	files := map[string][]byte{}
	files["provider/pkg/kinds/kinds.go"] = mustRenderGoTemplate(filepath.Join(TemplateDir, "kinds", "kinds.tmpl"), gvk)
	mustWriteFiles(BaseDir, files)
}

func mustWriteFiles(rootDir string, files map[string][]byte) {
	for filename, contents := range files {
		mustWriteFile(rootDir, filename, contents)
	}
}

func mustWriteFile(rootDir, filename string, contents []byte) {
	outPath := filepath.Join(rootDir, filename)

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		panic(err)
	}
	err := ioutil.WriteFile(outPath, contents, 0644)
	if err != nil {
		panic(err)
	}
}

func mustWritePulumiSchema(pkgSpec schema.PackageSpec, version string) {
	schemaJSON, err := json.MarshalIndent(pkgSpec, "", "    ")
	if err != nil {
		panic(errors.Wrap(err, "marshaling Pulumi schema"))
	}

	mustWriteFile(BaseDir, filepath.Join("provider", "cmd", "pulumi-resource-kubernetes", "schema.json"), schemaJSON)

	versionedPkgSpec := pkgSpec
	versionedPkgSpec.Version = version
	versionedSchemaJSON, err := json.MarshalIndent(versionedPkgSpec, "", "    ")
	if err != nil {
		panic(errors.Wrap(err, "marshaling Pulumi schema"))
	}
	mustWriteFile(BaseDir, filepath.Join("sdk", "schema", "schema.json"), versionedSchemaJSON)
}
