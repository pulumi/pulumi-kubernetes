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

// nolint:gosec
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/gen"
	providerVersion "github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/version"
	"github.com/pulumi/pulumi/pkg/v3/codegen"
	dotnetgen "github.com/pulumi/pulumi/pkg/v3/codegen/dotnet"
	gogen "github.com/pulumi/pulumi/pkg/v3/codegen/go"
	nodejsgen "github.com/pulumi/pulumi/pkg/v3/codegen/nodejs"
	pythongen "github.com/pulumi/pulumi/pkg/v3/codegen/python"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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
		mustWriteTerraformMapping(pkgSpec)
	default:
		panic(fmt.Sprintf("Unrecognized language '%s'", language))
	}
}

func readSchema(schemaPath string, version string) *schema.Package {
	// Read in, decode, and import the schema.
	schemaBytes, err := os.ReadFile(schemaPath)
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
	swagger, err := os.ReadFile(swaggerPath)
	if err != nil {
		panic(err)
	}

	swaggerDir := filepath.Dir(swaggerPath)

	// The following APIs have been deprecated and removed in the more recent versions of k8s:
	// - extensions/v1beta1/*
	// - apps/v1beta1/*
	// - apps/v1beta2/*
	// - networking/v1beta1/IngressClass
	// - flowcontrol/v1beta2/* (removed in v1.29, and has schema changes in v1.28)
	// Removed in v1.31:
	// - resource.k8s.io/v1alpha2
	// - networking.k8s.io/v1alpha1
	// Removed in v1.32:
	// - coordination.k8s.io/v1alpha1
	// - resource.k8s.io/v1alpha3
	// Since these resources will continue to be important to users for the foreseeable future, we will merge in
	// newer specs on top of this spec so that these resources continue to be available in our SDKs.
	urlFmt := "https://raw.githubusercontent.com/kubernetes/kubernetes/v1.%s.0/api/openapi-spec/swagger.json"
	filenameFmt := "swagger-v1.%s.0.json"
	for _, v := range []string{"17", "18", "19", "20", "26", "28", "30", "31"} {
		legacySwaggerPath := filepath.Join(swaggerDir, fmt.Sprintf(filenameFmt, v))
		err = DownloadFile(legacySwaggerPath, fmt.Sprintf(urlFmt, v))
		if err != nil {
			panic(err)
		}
		legacySwagger, err := os.ReadFile(legacySwaggerPath)
		if err != nil {
			panic(err)
		}
		swagger = mergeSwaggerSpecs(legacySwagger, swagger)
	}

	var schemaMap map[string]any
	err = json.Unmarshal(swagger, &schemaMap)
	if err != nil {
		panic(err)
	}

	// Generate schema
	return gen.PulumiSchema(schemaMap, gen.WithResourceOverlays(gen.ResourceOverlays), gen.WithTypeOverlays(gen.TypeOverlays))
}

// This is to mostly filter resources from the spec.
var resourcesToFilterFromTemplate = codegen.NewStringSet(
	"kubernetes:helm.sh/v3:Release",
	"kubernetes:helm.sh/v4:Chart",
	"kubernetes:kustomize/v2:Directory",
	"kubernetes:yaml/v2:ConfigFile",
	"kubernetes:yaml/v2:ConfigGroup",
)

func writeNodeJSClient(pkg *schema.Package, outdir, templateDir string) {
	resources, err := nodejsgen.LanguageResources(pkg)
	if err != nil {
		panic(err)
	}

	templateResources := gen.TemplateResources{}
	packages := codegen.StringSet{}
	for tok, resource := range resources {
		if resourcesToFilterFromTemplate.Has(tok) {
			continue
		}
		if resource.Package == "" {
			continue
		}
		if strings.HasSuffix(resource.Name, "Patch") {
			continue
		}
		tr := gen.TemplateResource{
			Name:    resource.Name,
			Package: resource.Package,
			Token:   resource.Token,
		}
		for _, property := range resource.Properties {
			// hack(levi): manually remove `| undefined` from the ConstValue and Package until https://github.com/pulumi/pulumi-kubernetes/issues/1650 is resolved.
			cv := strings.TrimSuffix(property.ConstValue, " | undefined")
			pkg := strings.TrimSuffix(property.Package, " | undefined")

			tp := gen.TemplateProperty{
				ConstValue: cv,
				Name:       property.Name,
				Package:    pkg,
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
		"apiextensions/customResource.ts":      mustLoadFile(filepath.Join(templateDir, "apiextensions", "customResource.ts")),
		"apiextensions/customResourcePatch.ts": mustLoadFile(filepath.Join(templateDir, "apiextensions", "customResourcePatch.ts")),
		"helm/v3/helm.ts":                      mustLoadFile(filepath.Join(templateDir, "helm", "v3", "helm.ts")),
		"kustomize/kustomize.ts":               mustLoadFile(filepath.Join(templateDir, "kustomize", "kustomize.ts")),
		"yaml/yaml.ts":                         mustRenderTemplate(filepath.Join(templateDir, "yaml", "yaml.tmpl"), templateResources),
	}
	files, err := nodejsgen.GeneratePackage("pulumigen", pkg, overlays, nil, false, nil)
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
	for tok, resource := range resources {
		if resourcesToFilterFromTemplate.Has(tok) {
			continue
		}
		if resource.Name == "CustomResourceDefinition" { // Use manual overlay in yaml.tmpl
			continue
		}
		if strings.HasSuffix(resource.Name, "Patch") {
			continue
		}
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
		"apiextensions/CustomResource.py":      mustLoadFile(filepath.Join(templateDir, "apiextensions", "CustomResource.py")),
		"apiextensions/CustomResourcePatch.py": mustLoadFile(filepath.Join(templateDir, "apiextensions", "CustomResourcePatch.py")),
		"helm/v3/helm.py":                      mustLoadFile(filepath.Join(templateDir, "helm", "v3", "helm.py")),
		"kustomize/kustomize.py":               mustLoadFile(filepath.Join(templateDir, "kustomize", "kustomize.py")),
		"yaml/yaml.py":                         mustRenderTemplate(filepath.Join(templateDir, "yaml", "yaml.tmpl"), templateResources),
	}

	files, err := pythongen.GeneratePackage("pulumigen", pkg, overlays, nil)
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
	for tok, resource := range resources {
		if resourcesToFilterFromTemplate.Has(tok) {
			continue
		}
		if strings.HasSuffix(resource.Name, "Patch") {
			continue
		}
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
		"ApiExtensions/CustomResource.cs":      mustLoadFile(filepath.Join(templateDir, "apiextensions", "CustomResource.cs")),
		"ApiExtensions/CustomResourcePatch.cs": mustLoadFile(filepath.Join(templateDir, "apiextensions", "CustomResourcePatch.cs")),
		"Helm/ChartBase.cs":                    mustLoadFile(filepath.Join(templateDir, "helm", "ChartBase.cs")),
		"Helm/Unwraps.cs":                      mustLoadFile(filepath.Join(templateDir, "helm", "Unwraps.cs")),
		"Helm/V3/Chart.cs":                     mustLoadFile(filepath.Join(templateDir, "helm", "v3", "Chart.cs")),
		"Helm/V3/Invokes.cs":                   mustLoadFile(filepath.Join(templateDir, "helm", "v3", "Invokes.cs")),
		"Kustomize/Directory.cs":               mustLoadFile(filepath.Join(templateDir, "kustomize", "Directory.cs")),
		"Kustomize/Invokes.cs":                 mustLoadFile(filepath.Join(templateDir, "kustomize", "Invokes.cs")),
		"Yaml/ConfigFile.cs":                   mustLoadFile(filepath.Join(templateDir, "yaml", "ConfigFile.cs")),
		"Yaml/ConfigGroup.cs":                  mustLoadFile(filepath.Join(templateDir, "yaml", "ConfigGroup.cs")),
		"Yaml/Invokes.cs":                      mustLoadFile(filepath.Join(templateDir, "yaml", "Invokes.cs")),
		"Yaml/TransformationAction.cs":         mustLoadFile(filepath.Join(templateDir, "yaml", "TransformationAction.cs")),
		"Yaml/Yaml.cs":                         mustRenderTemplate(filepath.Join(templateDir, "yaml", "yaml.tmpl"), templateResources),
	}

	files, err := dotnetgen.GeneratePackage("pulumigen", pkg, overlays, nil)
	if err != nil {
		panic(err)
	}

	for filename, contents := range files {
		path := filepath.Join(outdir, filename)

		if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			panic(err)
		}
		err := os.WriteFile(path, contents, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func writeGoClient(pkg *schema.Package, outdir string, templateDir string) {
	files, err := gogen.GeneratePackage("pulumigen", pkg, nil)
	if err != nil {
		panic(err)
	}
	renamePackage := func(fileNames []string, sourcePackage, renameTo string) {
		re := regexp.MustCompile(fmt.Sprintf(`(%s)`, sourcePackage))

		for _, f := range fileNames {
			content, ok := files[f]
			if !ok {
				contract.Failf("Expected file: %q but not found.", f)
			}
			files[f] = re.ReplaceAll(content, []byte(renameTo))
		}
	}

	// Go codegen maps package to "v3" for Helm Release. Manually rename to
	// helm to avoid conflict with existing templates.
	renamePackage([]string{
		"kubernetes/helm/v3/pulumiTypes.go",
		"kubernetes/helm/v3/init.go",
		"kubernetes/helm/v3/release.go",
	},
		"package v3",
		"package helm")

	resources, err := gogen.LanguageResources("pulumigen", pkg)
	if err != nil {
		panic(err)
	}

	templateResources := gen.GoTemplateResources{}
	for tok, resource := range resources {
		if resourcesToFilterFromTemplate.Has(tok) {
			continue
		}
		if strings.HasSuffix(resource.Name, "Patch") {
			continue
		}
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

	files["kubernetes/customPulumiTypes.go"] = mustLoadGoFile(filepath.Join(templateDir, "customPulumiTypes.go"))
	files["kubernetes/apiextensions/customResource.go"] = mustLoadGoFile(filepath.Join(templateDir, "apiextensions", "customResource.go"))
	files["kubernetes/apiextensions/customResourcePatch.go"] = mustLoadGoFile(filepath.Join(templateDir, "apiextensions", "customResourcePatch.go"))
	files["kubernetes/helm/v3/chart.go"] = mustLoadGoFile(filepath.Join(templateDir, "helm", "v3", "chart.go"))
	// Rename pulumiTypes.go to avoid conflict with schema generated Helm Release types.
	files["kubernetes/helm/v3/chartPulumiTypes.go"] = mustLoadGoFile(filepath.Join(templateDir, "helm", "v3", "pulumiTypes.go"))
	files["kubernetes/kustomize/directory.go"] = mustLoadGoFile(filepath.Join(templateDir, "kustomize", "directory.go"))
	files["kubernetes/kustomize/pulumiTypes.go"] = mustLoadGoFile(filepath.Join(templateDir, "kustomize", "pulumiTypes.go"))
	files["kubernetes/yaml/configFile.go"] = mustLoadGoFile(filepath.Join(templateDir, "yaml", "configFile.go"))
	files["kubernetes/yaml/configGroup.go"] = mustLoadGoFile(filepath.Join(templateDir, "yaml", "configGroup.go"))
	files["kubernetes/yaml/transformation.go"] = mustLoadGoFile(filepath.Join(templateDir, "yaml", "transformation.go"))
	files["kubernetes/yaml/yaml.go"] = mustRenderGoTemplate(filepath.Join(templateDir, "yaml", "yaml.tmpl"), templateResources)
	files["kubernetes/yaml/v2/kinds.go"] = mustRenderGoTemplate(filepath.Join(templateDir, "yaml", "v2", "kinds.tmpl"), templateResources)

	mustWriteFiles(outdir, files)
}

func mustLoadFile(path string) []byte {
	b, err := os.ReadFile(path)
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

func mustRenderTemplate(path string, resources any) []byte {
	b := mustLoadFile(path)
	t := template.Must(template.New("resources").Parse(string(b)))

	var buf bytes.Buffer
	err := t.Execute(&buf, resources)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func mustRenderGoTemplate(path string, resources any) []byte {
	bytes := mustRenderTemplate(path, resources)

	formattedSource, err := format.Source(bytes)
	contract.AssertNoErrorf(err, "err: %+v path: %q source:\n%s", err, path, string(bytes))
	return formattedSource
}

func genK8sResourceTypes(pkg *schema.Package) {
	groupVersions, kinds, patchKinds, listKinds := codegen.NewStringSet(), codegen.NewStringSet(), codegen.NewStringSet(), codegen.NewStringSet()
	for _, resource := range pkg.Resources {
		if resourcesToFilterFromTemplate.Has(resource.Token) {
			continue
		}
		parts := strings.Split(resource.Token, ":")
		contract.Assertf(len(parts) == 3, "expected resource token to have three elements: %s", resource.Token)

		groupVersion, kind := parts[1], parts[2]

		if resource.IsOverlay {
			continue
		}
		if strings.HasSuffix(kind, "Patch") {
			patchKinds.Add(resource.Token)
			continue
		}
		if strings.HasSuffix(kind, "List") {
			listKinds.Add(resource.Token)
		}

		groupVersions.Add(groupVersion)
		kinds.Add(kind)
	}

	gvk := gen.GVK{Kinds: kinds.SortedValues(), PatchKinds: patchKinds.SortedValues(), ListKinds: listKinds.SortedValues()}
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
	err := os.WriteFile(outPath, contents, 0644)
	if err != nil {
		panic(err)
	}
}

func makeJSONString(v any) ([]byte, error) {
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func mustWritePulumiSchema(pkgSpec schema.PackageSpec, version string) {
	schemaJSON, err := makeJSONString(pkgSpec)
	if err != nil {
		panic(fmt.Errorf("marshaling Pulumi schema: %w", err))
	}

	mustWriteFile(BaseDir, filepath.Join("provider", "cmd", "pulumi-resource-kubernetes", "schema.json"), schemaJSON)

	versionedPkgSpec := pkgSpec
	versionedPkgSpec.Version = version
	versionedSchemaJSON, err := makeJSONString(versionedPkgSpec)
	if err != nil {
		panic(fmt.Errorf("marshaling Pulumi schema: %w", err))
	}
	mustWriteFile(BaseDir, filepath.Join("sdk", "schema", "schema.json"), versionedSchemaJSON)
}

// Minimal types for reading the terraform schema file
type TerraformAttributeSchema struct{}

type TerraformBlockTypeSchema struct {
	Block    *TerraformBlockSchema `json:"block"`
	MaxItems int                   `json:"max_items"`
}

type TerraformBlockSchema struct {
	Attributes map[string]*TerraformAttributeSchema `json:"attributes"`
	BlockTypes map[string]*TerraformBlockTypeSchema `json:"block_types"`
}

type TerraformResourceSchema struct {
	Block *TerraformBlockSchema `json:"block"`
}

type TerraformSchema struct {
	Provider          *TerraformResourceSchema            `json:"provider"`
	ResourceSchemas   map[string]*TerraformResourceSchema `json:"resource_schemas"`
	DataSourceSchemas map[string]*TerraformResourceSchema `json:"data_source_schemas"`
}

func findPulumiTokenFromTerraformToken(pkgSpec schema.PackageSpec, token string) tokens.Token {
	// strip off the leading "kubernetes_"
	token = strings.TrimPrefix(token, "kubernetes_")
	// split of any _v1, _v2, etc. suffix
	tokenParts := strings.Split(token, "_")
	maybeVersion := tokenParts[len(tokenParts)-1]
	versions := []string{
		"v1alpha1",
		"v1beta1",
		"v1beta2",
		"v1",
		"v2alpha1",
		"v2beta1",
		"v2beta2",
		"v2",
	}
	versionIndex := -1
	if maybeVersion[0] == 'v' && unicode.IsDigit(rune(maybeVersion[1])) {
		for i, v := range versions {
			if maybeVersion == v {
				versionIndex = i
				break
			}
		}
		if versionIndex == -1 {
			panic(fmt.Sprintf("unexpected version suffix %q in token %q", maybeVersion, token))
		}
	}

	// Get the full token as a pulumi camel case name
	length := len(tokenParts) - 1
	if versionIndex == -1 {
		// If the last item wasn't a version add it back to the token and clear maybeVersion
		length++
		maybeVersion = ""
	}
	caser := cases.Title(language.English)
	for i := 0; i < length; i++ {
		tokenParts[i] = caser.String(tokenParts[i])
	}
	searchToken := strings.Join(tokenParts[:length], "")

	foundTokens := make([]tokens.Token, 0)
	for t := range pkgSpec.Resources {
		pulumiToken := tokens.Token(t)
		member := pulumiToken.ModuleMember()
		memberName := member.Name()
		if string(memberName) == searchToken {
			foundTokens = append(foundTokens, pulumiToken)
		}
	}

	// If we didn't find any tokens, return an empty string
	if len(foundTokens) == 0 {
		return ""
	}

	// Else try and workout the right version to use. If we have a version suffix, try to use that...
	var foundToken tokens.Token
	if versionIndex != -1 {
		contract.Assertf(maybeVersion != "", "expected maybeVersion to be set")

		for _, t := range foundTokens {
			module := t.Module()
			if strings.HasSuffix(string(module), maybeVersion) {
				// Exact match, but we might see this from multiple modules. Prefer the the non-extension one.
				if foundToken == "" {
					foundToken = t
				}

				if !strings.Contains(string(t), ":extensions/") {
					foundToken = t
				}
			}
		}
		if foundToken != "" {
			return foundToken
		}
	}

	// ...otherwise use the _newest_ version. e.g. kubernetes_thing should map to kubernetes:core/v2:Thing,
	// not kubernetes:core/v1:Thing. `versions` is sorted for this purpose.
	highestIndex := -1
	for _, t := range foundTokens {
		module := t.Module()
		modulesVersionIndex := -1
		for i, v := range versions {
			if strings.HasSuffix(string(module), v) {
				modulesVersionIndex = i
				break
			}
		}

		if modulesVersionIndex == -1 {
			panic(fmt.Sprintf("unexpected module version %q in token %q", module, t))
		}

		if highestIndex == modulesVersionIndex {
			// If we've seen this version before prefer the non-extension module
			if !strings.Contains(string(t), ":extensions/") {
				foundToken = t
			}
		} else if highestIndex < modulesVersionIndex {
			highestIndex = modulesVersionIndex
			foundToken = t
		}
	}

	return foundToken
}

func buildPulumiFieldsFromTerraform(path string, block *TerraformBlockSchema) map[string]any {
	// Recursively build up the fields for this resource
	fields := make(map[string]any)

	// Attributes _might_ need to be renamed
	for attrName := range block.Attributes {
		field := map[string]any{}

		// Manual fixups for the schema

		// Only add this field if it says something meaningful
		if len(field) > 0 {
			fields[attrName] = field
		}
	}

	for blockName, blockType := range block.BlockTypes {
		field := map[string]any{}

		// If the block has a max_items of 1, then we need to tell the converter that
		if blockType.MaxItems == 1 {
			field["maxItemsOne"] = true
		}

		// Recurse to see if the block needs to return any fields
		elem := buildPulumiFieldsFromTerraform(path+"."+blockName, blockType.Block)
		if len(elem) > 0 {
			// Based on if we're treating this as a list of not elem should either be added to the "fields"
			// field or nested under the "element" field
			if field["maxItemsOne"] == true {
				field["fields"] = elem
			} else {
				field["element"] = map[string]any{
					"fields": elem,
				}
			}
		}

		// Manual fixups for the schema, most of these look like pluralization issues, but not sure if there's
		// a safe way to do this automatically.

		//1. kubernetes_deployment has a field "container" which is a list, but we call it "containers"
		if path == "kubernetes_deployment.spec.template.spec" && blockName == "container" {
			field["name"] = "containers"
		}
		// 2. kubernetes_deployment has a field "port" which is a list, but we call it "ports"
		if path == "kubernetes_deployment.spec.template.spec.container" && blockName == "port" {
			field["name"] = "ports"
		}
		// 3. kubernetes_service has a field "port" which is a list, but we call it "ports"
		if path == "kubernetes_service.spec" && blockName == "port" {
			field["name"] = "ports"
		}

		// Only add this field if it says something meaningful
		if len(field) > 0 {
			fields[blockName] = field
		}
	}
	return fields
}

func mustWriteTerraformMapping(pkgSpec schema.PackageSpec) {
	// The terraform converter expects the mapping to be the JSON serialization of it's ProviderInfo
	// structure. We can get away with returning a _very_ limited subset of the information here, since the
	// converter only cares about a few fields and is tolerant of missing fields. We get the terraform
	// kubernetes schema by running `terraform providers schema -json` in a minimal terraform project that
	// defines a kubernetes provider.
	rawTerraformSchema := mustLoadFile(filepath.Join(BaseDir, "provider", "cmd", "pulumi-gen-kubernetes", "terraform.json"))

	var terraformSchema TerraformSchema
	err := json.Unmarshal(rawTerraformSchema, &terraformSchema)
	if err != nil {
		panic(err)
	}

	resources := make(map[string]any)
	for tftok, resource := range terraformSchema.ResourceSchemas {
		putok := findPulumiTokenFromTerraformToken(pkgSpec, tftok)
		// Skip if the token is empty.
		if putok == "" {
			continue
		}

		// Need to fill in just enough fields so that MaxItemsOne is set correctly for things.
		resources[tftok] = map[string]any{
			"tok":    putok,
			"fields": buildPulumiFieldsFromTerraform(tftok, resource.Block),
		}
	}

	info := map[string]any{
		"name":        "kubernetes",
		"provider":    map[string]any{},
		"resources":   resources,
		"dataSources": map[string]any{},
	}

	data, err := makeJSONString(info)
	if err != nil {
		panic(err)
	}

	mustWriteFile(BaseDir, filepath.Join("provider", "cmd", "pulumi-resource-kubernetes", "terraform-mapping.json"), data)
}
