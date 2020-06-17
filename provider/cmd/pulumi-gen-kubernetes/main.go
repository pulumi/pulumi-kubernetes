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
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/provider/v2/pkg/gen"
	"github.com/pulumi/pulumi/pkg/v2/codegen"
	dotnetgen "github.com/pulumi/pulumi/pkg/v2/codegen/dotnet"
	gogen "github.com/pulumi/pulumi/pkg/v2/codegen/go"
	nodejsgen "github.com/pulumi/pulumi/pkg/v2/codegen/nodejs"
	pythongen "github.com/pulumi/pulumi/pkg/v2/codegen/python"
	"github.com/pulumi/pulumi/pkg/v2/codegen/schema"
	"github.com/pulumi/pulumi/sdk/v2/go/common/util/contract"
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
	Python Language = "python"
	Schema Language = "schema"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Usage: gen <language> <swagger-file> <root-pulumi-kubernetes-dir>")
	}

	language := Language(os.Args[1])

	swagger, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}

	swaggerDir := filepath.Dir(os.Args[2])

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

	BaseDir = os.Args[3]
	TemplateDir = path.Join(BaseDir, "provider", "pkg", "gen")
	outdir := path.Join(BaseDir, "sdk", string(language))

	// Generate schema
	pkgSpec := gen.PulumiSchema(data)

	// Generate package from schema
	pkg := genPulumiSchemaPackage(pkgSpec)

	// Generate provider code
	genK8sResourceTypes(pkg)

	switch language {
	case NodeJS:
		templateDir := path.Join(TemplateDir, "nodejs-templates")
		writeNodeJSClient(pkg, outdir, templateDir)
	case Python:
		templateDir := path.Join(TemplateDir, "python-templates")
		writePythonClient(pkg, outdir, templateDir)
	case DotNet:
		templateDir := path.Join(TemplateDir, "dotnet-templates")
		writeDotnetClient(pkg, data, outdir, templateDir)
	case Go:
		templateDir := path.Join(TemplateDir, "go-templates")
		writeGoClient(pkg, outdir, templateDir)
	case Schema:
		mustWritePulumiSchema(pkgSpec, outdir)
	default:
		panic(fmt.Sprintf("Unrecognized language '%s'", language))
	}
}

func writeNodeJSClient(pkg *schema.Package, outdir, templateDir string) {
	resources, err := nodejsgen.LanguageResources(pkg)
	if err != nil {
		panic(err)
	}

	templateResources := gen.TemplateResources{}
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
	}
	sort.Slice(templateResources.Resources, func(i, j int) bool {
		return templateResources.Resources[i].Token < templateResources.Resources[j].Token
	})

	overlays := map[string][]byte{
		"apiextensions/customResource.ts": mustLoadFile(filepath.Join(templateDir, "apiextensions", "customResource.ts")),
		"helm/v2/helm.ts":                 mustLoadFile(filepath.Join(templateDir, "helm", "v2", "helm.ts")),
		"helm/v3/helm.ts":                 mustLoadFile(filepath.Join(templateDir, "helm", "v2", "helm.ts")), // v3 support is currently identical to v2
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
		"apiextensions/custom_resource.py": mustLoadFile(filepath.Join(templateDir, "apiextensions", "custom_resource.py")),
		"helm/v2/helm.py":                  mustLoadFile(filepath.Join(templateDir, "helm", "v2", "helm.py")),
		"helm/v3/helm.py":                  mustLoadFile(filepath.Join(templateDir, "helm", "v2", "helm.py")), // v3 support is currently identical to v2
		"yaml.py":                          mustRenderTemplate(filepath.Join(templateDir, "yaml", "yaml.tmpl"), templateResources),
	}

	files, err := pythongen.GeneratePackage("pulumigen", pkg, overlays)
	if err != nil {
		panic(err)
	}

	mustWriteFiles(outdir, files)
}

func writeDotnetClient(pkg *schema.Package, data map[string]interface{}, outdir, templateDir string) {
	files, err := dotnetgen.GeneratePackage("pulumigen", pkg, nil)
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

	yamlcs, err := gen.DotnetYaml(data, templateDir)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/Yaml/Yaml.cs", outdir), []byte(yamlcs), 0777)
	if err != nil {
		panic(err)
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

	files["kubernetes/types.go"] = mustRenderGoTemplate(filepath.Join(templateDir, "types.tmpl"), templateResources)
	files["kubernetes/apiextensions/customResource.go"] = mustRenderGoTemplate(filepath.Join(templateDir, "apiextensions", "customResource.tmpl"), templateResources)
	files["kubernetes/helm/v2/chart.go"] = mustRenderGoTemplate(filepath.Join(templateDir, "helm", "v2", "chart.tmpl"), templateResources)
	files["kubernetes/helm/v2/types.go"] = mustRenderGoTemplate(filepath.Join(templateDir, "helm", "v2", "types.tmpl"), templateResources)
	files["kubernetes/yaml/configFile.go"] = mustRenderGoTemplate(filepath.Join(templateDir, "yaml", "configFile.tmpl"), templateResources)
	files["kubernetes/yaml/configGroup.go"] = mustRenderGoTemplate(filepath.Join(templateDir, "yaml", "configGroup.tmpl"), templateResources)
	files["kubernetes/yaml/transformation.go"] = mustRenderGoTemplate(filepath.Join(templateDir, "yaml", "transformation.tmpl"), templateResources)
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

func genPulumiSchemaPackage(pkgSpec schema.PackageSpec) *schema.Package {

	pkg, err := schema.ImportSpec(pkgSpec, nil)
	if err != nil {
		panic(err)
	}
	return pkg
}

func genK8sResourceTypes(pkg *schema.Package) {
	groupVersions, kinds := codegen.NewStringSet(), codegen.NewStringSet()
	for _, resource := range pkg.Resources {
		parts := strings.Split(resource.Token, ":")
		contract.Assert(len(parts) == 3)

		groupVersion, kind := parts[1], parts[2]
		if strings.HasSuffix(kind, "List") {
			continue
		}
		groupVersions.Add(groupVersion)
		kinds.Add(kind)
	}

	gvk := gen.GVK{Kinds: kinds.SortedValues()}
	gvStrings := groupVersions.SortedValues()
	for _, gvString := range gvStrings {
		gvk.GroupVersions = append(gvk.GroupVersions, gen.GroupVersion(gvString))
	}

	files := map[string][]byte{}
	files["provider/pkg/kinds/kinds.go"] = mustRenderGoTemplate(path.Join(TemplateDir, "kinds", "kinds.tmpl"), gvk)
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

func mustWritePulumiSchema(pkgSpec schema.PackageSpec, outDir string) {
	schemaJSON, err := json.MarshalIndent(pkgSpec, "", "    ")
	if err != nil {
		panic(errors.Wrap(err, "marshaling Pulumi schema"))
	}

	mustWriteFile(outDir, "schema.json", schemaJSON)
}
