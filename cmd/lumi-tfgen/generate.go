// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"github.com/pulumi/lumi/pkg/diag"
	"github.com/pulumi/lumi/pkg/tools"
	"github.com/pulumi/lumi/pkg/util/cmdutil"
	"github.com/pulumi/lumi/pkg/util/contract"

	"github.com/pulumi/terraform-bridge/pkg/tfbridge"
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
func (g *generator) Generate(provs map[string]tfbridge.ProviderInfo, path string) error {
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
func (g *generator) generateProvider(pkg string, provinfo tfbridge.ProviderInfo, root string) error {
	var files []string
	exports := make(map[string]string)               // a list of top-level exports.
	modules := make(map[string]string)               // a list of modules to export individually.
	submodules := make(map[string]map[string]string) // a map of sub-module name to exported members.

	// Ensure the root path exists.
	if err := tools.EnsureDir(root); err != nil {
		return err
	}

	// Place all configuration variables into a single config module.
	prov := provinfo.P
	if len(prov.Schema) > 0 {
		cfgfile, err := g.generateConfig(prov.Schema, root)
		if err != nil {
			return err
		}
		modules["config"] = cfgfile // ensure we reexport as config
		files = append(files, cfgfile)
	}

	// For each resource, create its own dedicated type and module export.
	resmap := prov.ResourcesMap
	reshits := make(map[string]bool)
	for _, r := range stableResources(prov.ResourcesMap) {
		var resinfo tfbridge.ResourceInfo
		if prov.Resources != nil {
			if ri, has := provinfo.Resources[r]; has {
				resinfo = ri
				reshits[r] = true
			} else {
				// if this has a map, but this resource wasn't found, issue a warning.
				cmdutil.Diag().Warningf(
					diag.Message("Resource %v not found in provider map; using default naming"), r)
			}
		}
		result, err := g.generateResource(pkg, r, resmap[r], resinfo, root)
		if err != nil {
			return err
		}
		if result.Submod == "" {
			// if no sub-module, export flatly in our own index.
			exports[result.Name] = result.File
		} else {
			// otherwise, make sure to track this in the submodule so we can create and export it correctly.
			submod := submodules[result.Submod]
			if submod == nil {
				submod = make(map[string]string)
				submodules[result.Submod] = submod
			}
			submod[result.Name] = result.File
		}
		files = append(files, result.File)
	}

	// Emit a warning if there is a map but some names didn't match.
	if provinfo.Resources != nil {
		var resnames []string
		for resname := range provinfo.Resources {
			resnames = append(resnames, resname)
		}
		sort.Strings(resnames)
		for _, resname := range resnames {
			if !reshits[resname] {
				cmdutil.Diag().Warningf(
					diag.Message("Resource %v (%v) wasn't found in the Terraform module; possible name mismatch?"),
					resname, provinfo.Resources[resname].Tok)
			}
		}
	}

	// Generate any submodules and add them to the export list.
	subs, err := g.generateSubmodules(submodules, root)
	if err != nil {
		return err
	}
	for sub, subf := range subs {
		if conflict, has := modules[sub]; has {
			cmdutil.Diag().Errorf(
				diag.Message("Conflicting submodule %v; exists for both %v and %v"), sub, conflict, subf)
		}
		modules[sub] = subf
	}

	// Generate the index.ts file that reexports everything at the entrypoint.
	ixfile, err := g.generateIndex(exports, modules, root)
	if err != nil {
		return err
	}
	files = append(files, ixfile)

	// Finally, generate all of the package metadata: Lumi.yaml, package.json, and tsconfig.json.
	return g.generatePackageMetadata(pkg, files, root)
}

// generateConfig takes a map of config variables and emits a config module to the given file.
func (g *generator) generateConfig(cfg map[string]*schema.Schema, root string) (string, error) {
	// Sort the config variables to ensure they are emitted in a deterministic order.
	var cfgkeys []string
	for key := range cfg {
		cfgkeys = append(cfgkeys, key)
	}
	sort.Strings(cfgkeys)

	// Open up the file and spew a standard "code-generated" warning header.
	file := filepath.Join(root, "config.ts")
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
		w.Writefmtln("export let %v: %v;", propertyName(key), g.tfToJSTypeFlags(sch))
	}
	w.Writefmtln("")
	return file, nil
}

type resourceResult struct {
	Name   string // the resource name.
	File   string // the resource filename.
	Submod string // the submodule name, if any.
}

// generateResource generates a single module for the given resource.
func (g *generator) generateResource(pkg string, rawname string,
	res *schema.Resource, resinfo tfbridge.ResourceInfo, root string) (resourceResult, error) {
	// Transform the name as necessary.
	resname, filename := resourceName(pkg, rawname, resinfo)

	// Make a fully qualified file path that we will write to.
	file := filepath.Join(root, filename+".ts")

	// If the filename contains slashes, it is a sub-module, and we must ensure it exists.
	var submod string
	if slix := strings.Index(filename, "/"); slix != -1 {
		// Extract the module and file parts.
		submod = filename[:slix]
		if strings.Index(filename[slix+1:], "/") != -1 {
			return resourceResult{},
				errors.Errorf("Modules nested more than one level deep not currently supported")
		}

		// Ensure the submodule directory exists.
		if err := tools.EnsureFileDir(file); err != nil {
			return resourceResult{}, err
		}
	}

	// Open up the file and spew a standard "code-generated" warning header.
	w, err := tools.NewGenWriter(tfgen, file)
	if err != nil {
		return resourceResult{}, err
	}
	defer contract.IgnoreClose(w)
	w.EmitHeaderWarning()

	// Now import the modules we need.
	w.Writefmtln("import * as lumi from \"@lumi/lumi\";")
	w.Writefmtln("")

	// Generate the resource class.
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

	return resourceResult{
		Name:   resname,
		File:   file,
		Submod: submod,
	}, nil
}

// generateSubmodules creates a set of index files, if necessary, for the given submodules.  It returns a map of
// submodule name to the generated index file, so that a caller can be sure to re-export it as necessary.
func (g *generator) generateSubmodules(submodules map[string]map[string]string,
	root string) (map[string]string, error) {
	results := make(map[string]string)
	var subs []string
	for submod := range submodules {
		subs = append(subs, submod)
	}
	sort.Strings(subs)
	for _, sub := range subs {
		index, err := g.generateIndex(submodules[sub], nil, filepath.Join(root, sub))
		if err != nil {
			return nil, err
		}
		results[sub] = index
	}
	return results, nil
}

// generateIndex creates a module index file for easy access to sub-modules and exports.
func (g *generator) generateIndex(exports, modules map[string]string, root string) (string, error) {
	// Open up the file and spew a standard "code-generated" warning header.
	file := filepath.Join(root, "index.ts")
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
			rel, err := relModule(root, modules[mod])
			if err != nil {
				return "", err
			}
			w.Writefmtln("import * as %v from \"%v\";", mod, rel)
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
			rel, err := relModule(root, exports[exp])
			if err != nil {
				return "", err
			}
			w.Writefmtln("export * from \"%v\";", rel)
		}
		w.Writefmtln("")
	}

	return file, nil
}

// relModule removes the path suffix from a module and makes it relative to the root path.
func relModule(root string, mod string) (string, error) {
	// Return the path as a relative path to the root, so that imports are relative.
	file, err := filepath.Rel(root, mod)
	if err != nil {
		return "", err
	}

	if strings.HasSuffix(file, ".ts") {
		file = file[:len(file)-3]
	}

	return "./" + file, nil
}

// generatePackageMetadata generates all the non-code metadata required by a Lumi package.
func (g *generator) generatePackageMetadata(pkg string, files []string, root string) error {
	// There are three files to write out:
	//     1) Lumi.yaml: Lumi package information
	//     2) package.json: minimal NPM package metadata
	//     3) tsconfig.json: instructions for TypeScript compilation
	if err := g.generateLumiPackageMetadata(pkg, root); err != nil {
		return err
	}
	if err := g.generateNPMPackageMetadata(pkg, root); err != nil {
		return err
	}
	return g.generateTypeScriptProjectFile(pkg, files, root)
}

func (g *generator) generateLumiPackageMetadata(pkg string, root string) error {
	w, err := tools.NewGenWriter(tfgen, filepath.Join(root, "Lumi.yaml"))
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

func (g *generator) generateNPMPackageMetadata(pkg string, root string) error {
	w, err := tools.NewGenWriter(tfgen, filepath.Join(root, "package.json"))
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

func (g *generator) generateTypeScriptProjectFile(pkg string, files []string, root string) error {
	w, err := tools.NewGenWriter(tfgen, filepath.Join(root, "tsconfig.json"))
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)
	w.Writefmtln(`{
    "compilerOptions": {
        "outDir": ".lumi/bin",
        "target": "es6",
        "module": "commonjs",
        "moduleResolution": "node",
        "declaration": true,
        "sourceMap": true
    },
    "files": [`)
	for i, file := range files {
		var suffix string
		if i != len(files)-1 {
			suffix = ","
		}
		relfile, err := filepath.Rel(root, file)
		if err != nil {
			return err
		}
		w.Writefmtln("        \"%v\"%v", relfile, suffix)
	}
	w.Writefmtln(`    ]
}
`)
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

// resourceName translates a Terraform name into its Lumi name equivalent, plus a suggested filename.
func resourceName(pkg string, rawname string, resinfo tfbridge.ResourceInfo) (string, string) {
	if resinfo.Tok == "" {
		// default transformations.
		contract.Assert(strings.HasPrefix(rawname, pkg+"_"))
		name := rawname[len(pkg)+1:]                     // strip off the pkg prefix.
		return tfbridge.TerraformToLumiName(name, true), // PascalCase the resource name.
			tfbridge.TerraformToLumiName(name, false) // camelCase the filename.
	}
	// otherwise, a custom transformation exists; use it.
	return string(resinfo.Tok.Name()), string(resinfo.Tok.Module().Name())
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
	} else {
		s = tfbridge.TerraformToLumiName(s, false /*no to PascalCase; we want camelCase*/)
	}

	return s
}

func stableProviders(provs map[string]tfbridge.ProviderInfo) []string {
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
