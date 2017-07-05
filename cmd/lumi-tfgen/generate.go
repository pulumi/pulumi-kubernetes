// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pulumi/lumi/pkg/tools"
	"github.com/pulumi/lumi/pkg/util/contract"
)

type generator struct {
}

func newGenerator() *generator {
	return &generator{}
}

const tfgen = "the Lumi Terraform Bridge (TFGEN) Tool"

// Generate creates Lumi packages out of one or more Terraform plugins.  It accepts a list of all of the input Terraform
// providers, already bound statically to the code (since we cannot obtain schema information dynamically), walks them
// and generates the Lumi code, and spews that code into the output directory, path.
func (g *generator) Generate(provs map[string]*schema.Provider, path string) error {
	// If the path is empty, default to the working directory.
	if path == "" {
		p, err := os.Getwd()
		if err != nil {
			return err
		}
		path = p
	}

	// Enumerate each provider and generate its code into a distinct directory.
	for _, p := range stableProviders(provs) {
		// IDEA: let command lines specify different locations for different provider outputs.
		if err := g.generateProvider(p, provs[p], filepath.Join(path, p)); err != nil {
			return err
		}
	}
	return nil
}

// generateProvider creates a single standalone Lumi package for the given provider.
func (g *generator) generateProvider(pkg string, prov *schema.Provider, path string) error {
	var files []string
	exports := make(map[string]string)
	modules := make(map[string]string)

	// Ensure the target exists.
	if err := tools.EnsureDir(path); err != nil {
		return err
	}

	// Place all configuration variables into a single config module.
	if len(prov.Schema) > 0 {
		cfgfile, err := g.generateConfig(prov.Schema, path)
		if err != nil {
			return err
		}
		modules["config"] = cfgfile // ensure we reexport as config
		files = append(files, cfgfile)
	}

	// For each resource, create its own dedicated file.
	// TODO: we need to split these into nested sub-modules (ec2, s3, etc).
	for _, r := range stableResources(prov.ResourcesMap) {
		resfile, err := g.generateResource(pkg, r, prov.ResourcesMap[r], path)
		if err != nil {
			return err
		}
		exports[r] = resfile // ensure we export flatly as part of index
		files = append(files, resfile)
	}

	// Generate the index.ts file that reexports everything at the entrypoint.
	ixfile, err := g.generateIndex(exports, modules, path)
	if err != nil {
		return err
	}
	files = append(files, ixfile)

	// Finally, generate all of the package metadata: Lumi.yaml, package.json, and tsconfig.json.
	return g.generatePackageMetadata(pkg, files, path)
}

// generateConfig takes a map of config variables and emits a config module to the given file.
func (g *generator) generateConfig(cfg map[string]*schema.Schema, path string) (string, error) {
	// Sort the config variables to ensure they are emitted in a deterministic order.
	var cfgkeys []string
	for key := range cfg {
		cfgkeys = append(cfgkeys, key)
	}
	sort.Strings(cfgkeys)

	// Open up the file and spew a standard "code-generated" warning header.
	file := filepath.Join(path, "config.ts")
	w, err := tools.NewGenWriter(tfgen, file)
	if err != nil {
		return "", err
	}
	defer contract.IgnoreClose(w)
	w.EmitHeaderWarning()

	// Now just emit a simple export for each variable.
	for _, key := range cfgkeys {
		sch := cfg[key]
		if sch.Description != "" {
			// TODO: If there's a description, print it in the comment.
		}
		w.Writefmtln("export let %v: %v;", key, g.tfToJSTypeFlags(sch))
	}
	w.Writefmtln("")
	return filepath.Rel(path, file)
}

// generateResource generates a single module for the given resource.
func (g *generator) generateResource(pkg string, name string, res *schema.Resource, path string) (string, error) {
	// Open up the file and spew a standard "code-generated" warning header.
	file := filepath.Join(path, name+".ts")
	w, err := tools.NewGenWriter(tfgen, file)
	if err != nil {
		return "", err
	}
	defer contract.IgnoreClose(w)
	w.EmitHeaderWarning()

	// Now import the modules we need.
	w.Writefmtln("import * as lumi from \"@lumi/lumi\";")
	w.Writefmtln("")

	// Generate the resource class.
	resname := resourceName(pkg, name)
	w.Writefmtln("export class %[1]v extends lumi.NamedResource implements %[1]vArgs {", resname)

	// First, generate all instance properties.
	var props []string
	var schemas []*schema.Schema
	for _, s := range stableSchemas(res.Schema) {
		if sch := res.Schema[s]; sch.Removed == "" {
			// TODO: print out the description.
			// TODO: should we skip deprecated fields?
			// TODO: figure out how to deal with sensitive fields.
			prop := propertyName(s)
			w.Writefmtln("    public readonly %v%v: %v;", prop, g.tfToJSFlags(sch), g.tfToJSType(sch))
			props = append(props, prop)
			schemas = append(schemas, sch)
		}
	}
	if len(res.Schema) > 0 {
		w.Writefmtln("")
	}

	// TODO: figure out how to add get/query methods.

	// Now create a constructor that chains supercalls and stores into properties.
	w.Writefmtln("    constructor(name: string, args: %vArgs) {", resname)
	w.Writefmtln("        super(name);")
	for _, prop := range props {
		w.Writefmtln("        this.%[1]v = args.%[1]v;", prop)
	}
	w.Writefmtln("    }")
	w.Writefmtln("}")

	w.Writefmtln("")

	// Next, generate the args interface for this class.
	w.Writefmtln("export interface %vArgs {", resname)
	for i, sch := range schemas {
		w.Writefmtln("    readonly %v%v: %v;", props[i], g.tfToJSFlags(sch), g.tfToJSType(sch))
	}
	w.Writefmtln("}")
	w.Writefmtln("")

	return filepath.Rel(path, file)
}

// generateIndex creates a module index file for easy access to sub-modules and exports.
func (g *generator) generateIndex(exports, modules map[string]string, path string) (string, error) {
	// Open up the file and spew a standard "code-generated" warning header.
	file := filepath.Join(path, "index.ts")
	w, err := tools.NewGenWriter(tfgen, file)
	if err != nil {
		return "", err
	}
	defer contract.IgnoreClose(w)
	w.EmitHeaderWarning()

	// Import anything we will export as a sub-module, and then re-export it.
	if len(modules) > 0 {
		w.Writefmtln("// Export sub-modules:")
		var mods []string
		for mod := range modules {
			mods = append(mods, mod)
		}
		sort.Strings(mods)
		for _, mod := range mods {
			w.Writefmtln("import * as %v from \"%v\";", mod, relModule(modules[mod]))
		}
		w.Writefmt("export {")
		for i, mod := range mods {
			if i > 0 {
				w.Writefmt(", ")
			}
			w.Writefmt(mod)
		}
		w.Writefmtln("};")
		w.Writefmtln("")
	}

	// Export anything flatly that is a direct export rather than sub-module.
	if len(exports) > 0 {
		w.Writefmtln("// Export members:")
		var exps []string
		for exp := range exports {
			exps = append(exps, exp)
		}
		sort.Strings(exps)
		for _, exp := range exps {
			w.Writefmtln("export * from \"%v\";", relModule(exports[exp]))
		}
		w.Writefmtln("")
	}

	return filepath.Rel(path, file)
}

// relModule removes the path suffix from a module path.
func relModule(mod string) string {
	if strings.HasSuffix(mod, ".ts") {
		mod = mod[:len(mod)-3]
	}
	return "./" + mod
}

// generatePackageMetadata generates all the non-code metadata required by a Lumi package.
func (g *generator) generatePackageMetadata(pkg string, files []string, path string) error {
	// There are three files to write out:
	//     1) Lumi.yaml: Lumi package information
	//     2) package.json: minimal NPM package metadata
	//     3) tsconfig.json: instructions for TypeScript compilation
	if err := g.generateLumiPackageMetadata(pkg, path); err != nil {
		return err
	}
	if err := g.generateNPMPackageMetadata(pkg, path); err != nil {
		return err
	}
	return g.generateTypeScriptProjectFile(pkg, files, path)
}

func (g *generator) generateLumiPackageMetadata(pkg string, path string) error {
	w, err := tools.NewGenWriter(tfgen, filepath.Join(path, "Lumi.yaml"))
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)
	w.Writefmtln("name: tf-%v", pkg) // TODO: remove "tf-" prefix after we shake out the kinks.
	w.Writefmtln("description: An auto-generated bridge to Terraform's %v provider.", pkg)
	w.Writefmtln("dependencies:")
	w.Writefmtln("    lumi: \"*\"")
	w.Writefmtln("")
	return nil
}

func (g *generator) generateNPMPackageMetadata(pkg string, path string) error {
	w, err := tools.NewGenWriter(tfgen, filepath.Join(path, "package.json"))
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)
	w.Writefmtln("{")
	w.Writefmtln("    \"name\": \"@lumi/tf-%v\"", pkg) // TODO: remove "tf-" prefix after we shake out the kinks.
	w.Writefmtln("}")
	w.Writefmtln("")
	return nil
}

func (g *generator) generateTypeScriptProjectFile(pkg string, files []string, path string) error {
	w, err := tools.NewGenWriter(tfgen, filepath.Join(path, "tsconfig.json"))
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)
	w.Writefmtln("{")
	w.Writefmtln("    \"compilerOptions\": {")
	w.Writefmtln("        \"outDir\": \".lumi/bin\",")
	w.Writefmtln("        \"target\": \"es6\",")
	w.Writefmtln("        \"module\": \"commonjs\",")
	w.Writefmtln("        \"moduleResolution\": \"node\",")
	w.Writefmtln("        \"declaration\": true,")
	w.Writefmtln("        \"sourceMap\": true")
	w.Writefmtln("    },")
	w.Writefmtln("    \"files\": [")
	for i, file := range files {
		var suffix string
		if i != len(files)-1 {
			suffix = ","
		}
		if err != nil {
			return err
		}
		w.Writefmtln("        \"%v\"%v", file, suffix)
	}
	w.Writefmtln("    ]")
	w.Writefmtln("}")
	w.Writefmtln("")
	return nil
}

// tsToJSFlags returns the JavaScript flags for a given schema property.
func (g *generator) tfToJSFlags(sch *schema.Schema) string {
	if sch.Optional || sch.Computed {
		return "?"
	}
	return ""
}

// tfToJSType returns the JavaScript type name for a given schema property.
func (g *generator) tfToJSType(sch *schema.Schema) string {
	return g.tfToJSValueType(sch.Type, sch.Elem)
}

// tfToJSValueType returns the JavaScript type name for a given schema value type and element kind.
func (g *generator) tfToJSValueType(vt schema.ValueType, elem interface{}) string {
	switch vt {
	case schema.TypeBool:
		return "boolean"
	case schema.TypeInt, schema.TypeFloat:
		return "number"
	case schema.TypeString:
		return "string"
	case schema.TypeList:
		return fmt.Sprintf("%v[]", g.tfElemToJSType(elem))
	case schema.TypeMap:
		return fmt.Sprintf("{[key: string]: %v}", g.tfElemToJSType(elem))
	case schema.TypeSet:
		// IDEA: we can't use ES6 sets here, because we're using values and not objects.  It would be possible to come
		//     up with a ValueSet of some sorts, but that depends on things like shallowEquals which is known to be
		//     brittle and implementation dependent.  For now, we will stick to arrays, and validate on the backend.
		return fmt.Sprintf("%v[]", g.tfElemToJSType(elem))
	default:
		contract.Failf("Unrecognized schema type: %v", vt)
		return ""
	}
}

// tfElemToJSType returns the JavaScript type for a given schema element.  This element may be either a simple schema
// property or a complex structure.  In the case of a complex structure, this will expand to its nominal type.
func (g *generator) tfElemToJSType(elem interface{}) string {
	// If there is no element type specified, we will accept anything.
	if elem == nil {
		return "any"
	}

	switch e := elem.(type) {
	case schema.ValueType:
		return g.tfToJSValueType(e, nil)
	case *schema.Schema:
		// A simple type, just return its type name.
		return g.tfToJSType(e)
	case *schema.Resource:
		// A complex type, just expand to its nominal type name.
		// TODO: spill all complex structures in advance so that we don't have insane inline expansions.
		t := "{ "
		for i, s := range stableSchemas(e.Schema) {
			if i > 0 {
				t += ", "
			}
			sch := e.Schema[s]
			s = propertyName(s)
			t += fmt.Sprintf("%v%v: %v", s, g.tfToJSFlags(sch), g.tfToJSType(sch))
		}
		return t + " }"
	default:
		contract.Failf("Unrecognized schema element type: %v", e)
		return ""
	}
}

// tfToJSTypeFlags returns the JavaScript type name for a given schema property, just like tfToJSType, except that if
// the schema is optional, we will emit an undefined union type (for non-field positions).
func (g *generator) tfToJSTypeFlags(sch *schema.Schema) string {
	ts := g.tfToJSType(sch)
	if sch.Optional {
		ts += " | undefined"
	}
	return ts
}

// resourceName translates a Terraform underscore_cased_resource_name into the JavaScript PascalCasedResourceName.
func resourceName(pkg string, res string) string {
	contract.Assert(strings.HasPrefix(res, pkg+"_"))
	return res[len(pkg)+1:]
}

// propertyName translates a Terraform underscore_cased_property_name into the JavaScript camelCasedPropertyName.
func propertyName(s string) string {
	// BUGBUG: work around issue in the Elastic Transcoder where a field has a trailing ":".
	if strings.HasSuffix(s, ":") {
		s = s[:len(s)-1]
	}

	// If this property conflicts with resource's name or ID properties, rename it.
	// BUGBUG: this isn't a great solution.  We need to figure out whether to keep the underlying ones or not.
	if s == "name" {
		s = "_name"
	} else if s == "id" {
		s = "_id"
	}

	return s
}

func stableProviders(provs map[string]*schema.Provider) []string {
	var ps []string
	for p := range provs {
		ps = append(ps, p)
	}
	sort.Strings(ps)
	return ps
}

func stableResources(resources map[string]*schema.Resource) []string {
	var rs []string
	for r := range resources {
		rs = append(rs, r)
	}
	sort.Strings(rs)
	return rs
}

func stableSchemas(schemas map[string]*schema.Schema) []string {
	var ss []string
	for s := range schemas {
		ss = append(ss, s)
	}
	sort.Strings(ss)
	return ss
}
