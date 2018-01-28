// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfgen

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/golang/glog"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/tokens"
	"github.com/pulumi/pulumi/pkg/tools"
	"github.com/pulumi/pulumi/pkg/util/contract"

	"github.com/pulumi/pulumi-terraform/pkg/tfbridge"
)

// newNodeJSGenerator returns a language generator that understands how to produce Type/JavaScript packages.
func newNodeJSGenerator(pkg, version string, info tfbridge.ProviderInfo, overlaysDir, outDir string) langGenerator {
	return &nodeJSGenerator{
		pkg:         pkg,
		version:     version,
		info:        info,
		overlaysDir: overlaysDir,
		outDir:      outDir,
	}
}

type nodeJSGenerator struct {
	pkg         string
	version     string
	info        tfbridge.ProviderInfo
	overlaysDir string
	outDir      string
}

// commentChars returns the comment characters to use for single-line comments.
func (g *nodeJSGenerator) commentChars() string {
	return "//"
}

// moduleDir returns the directory for the given module.
func (g *nodeJSGenerator) moduleDir(mod *module) string {
	dir := g.outDir
	if mod.name != "" {
		dir = filepath.Join(dir, mod.name)
	}
	return dir
}

// openWriter opens a writer for the given module and file name, emitting the standard header automatically.
func (g *nodeJSGenerator) openWriter(mod *module, name string, needsSDK bool) (*tools.GenWriter, string, error) {
	dir := g.moduleDir(mod)
	file := filepath.Join(dir, name)
	w, err := tools.NewGenWriter(tfgen, file)
	if err != nil {
		return nil, "", err
	}

	// Emit a standard warning header ("do not edit", etc).
	w.EmitHeaderWarning(g.commentChars())

	// If needed, emit the standard Pulumi SDK import statement.
	if needsSDK {
		g.emitSDKImport(w)
	}

	return w, file, nil
}

func (g *nodeJSGenerator) emitSDKImport(w *tools.GenWriter) {
	w.Writefmtln("import * as pulumi from \"pulumi\";")
	w.Writefmtln("")
}

// emitPackage emits an entire package pack into the configured output directory with the configured settings.
func (g *nodeJSGenerator) emitPackage(pack *pkg) error {
	// First, generate the individual modules and their contents.
	files, submodules, err := g.emitModules(pack.modules)
	if err != nil {
		return err
	}

	// Generate a top-level index file that re-exports any child modules.
	index := pack.modules.ensureModule("")
	indexFiles, _, err := g.emitModule(index, submodules)
	if err != nil {
		return err
	}
	files = append(files, indexFiles...)

	// Finally emit the package metadata (NPM, TypeScript, and so on).
	sort.Strings(files)
	return g.emitPackageMetadata(pack, files)
}

// emitModules emits all modules in the given module map.  It returns a full list of files, a map of module to its
// associated index, and any error that occurred, if any.
func (g *nodeJSGenerator) emitModules(mmap moduleMap) ([]string, map[string]string, error) {
	var allFiles []string
	moduleMap := make(map[string]string)
	for _, mod := range mmap.values() {
		if mod.name == "" {
			continue // skip the root module, it is handled specially.
		}
		files, index, err := g.emitModule(mod, nil)
		if err != nil {
			return nil, nil, err
		}
		allFiles = append(allFiles, files...)
		moduleMap[mod.name] = index
	}
	return allFiles, moduleMap, nil
}

// emitModule emits a module.  This module ends up having many possible ES6 sub-modules which are then re-exported
// at the top level.  This is to make it convenient for overlays to import files within the same module without
// causing problematic cycles.  For example, imagine a module m with many members; the result is:
//
//     m/
//         index.ts
//         member1.ts
//         member<etc>.ts
//         memberN.ts
//
// The one special case is the configuration module, which yields a vars.ts file containing all exported variables.
////
// Note that the special module "" represents the top-most package module and won't be placed in a sub-directory.
//
// The return values are the full list of files generated, the index file, and any error that occurred, respectively.
func (g *nodeJSGenerator) emitModule(mod *module, submods map[string]string) ([]string, string, error) {
	glog.V(3).Infof("emitModule(%s)", mod.name)

	// Ensure that the target module directory exists.
	dir := g.moduleDir(mod)
	if err := tools.EnsureDir(dir); err != nil {
		return nil, "", errors.Wrapf(err, "creating module directory")
	}

	// Now, enumerate each module member, in the order presented to us, and do the right thing.
	var files []string
	for _, member := range mod.members {
		file, err := g.emitModuleMember(mod, member)
		if err != nil {
			return nil, "", errors.Wrapf(err, "emitting module %s member %s", mod.name, member.Name())
		} else if file != "" {
			files = append(files, file)
		}
	}

	// If this is a config module, we need to emit the configuration variables.
	if mod.config() {
		file, err := g.emitConfigVariables(mod)
		if err != nil {
			return nil, "", errors.Wrapf(err, "emitting config module variables")
		}
		files = append(files, file)
	}

	// Lastly, generate an index file for this module.
	index, err := g.emitIndex(mod, files, submods)
	if err != nil {
		return nil, "", errors.Wrapf(err, "emitting module %s index", mod.name)
	}
	files = append(files, index)

	return files, index, nil
}

// emitIndex emits an index module, optionally re-exporting other members or submodules.
func (g *nodeJSGenerator) emitIndex(mod *module, exports []string, submods map[string]string) (string, error) {
	// Open the index.ts file for this module, and ready it for writing.
	w, index, err := g.openWriter(mod, "index.ts", false)
	if err != nil {
		return "", err
	}
	defer contract.IgnoreClose(w)

	// Export anything flatly that is a direct export rather than sub-module.
	if len(exports) > 0 {
		w.Writefmtln("// Export members:")
		var exps []string
		exps = append(exps, exports...)
		sort.Strings(exps)
		for _, exp := range exps {
			rel, err := g.relModule(mod, exp)
			if err != nil {
				return "", err
			}
			w.Writefmtln("export * from \"%s\";", rel)
		}
	}

	// Finally, f there are submodules, export them.
	if len(submods) > 0 {
		if len(exports) > 0 {
			w.Writefmtln("")
		}
		w.Writefmtln("// Export sub-modules:")
		var subs []string
		for sub := range submods {
			subs = append(subs, sub)
		}
		sort.Strings(subs)
		for _, sub := range subs {
			rel, err := g.relModule(mod, submods[sub])
			if err != nil {
				return "", err
			}
			w.Writefmtln("import * as %s from \"%s\";", sub, rel)
		}
		w.Writefmt("export {")
		for i, sub := range subs {
			if i > 0 {
				w.Writefmt(", ")
			}
			w.Writefmt(sub)
		}
		w.Writefmtln("};")
	}

	return index, nil
}

// emitModuleMember emits the given member, and returns the module file that it was emitted into (if any).
func (g *nodeJSGenerator) emitModuleMember(mod *module, member moduleMember) (string, error) {
	glog.V(3).Infof("emitModuleMember(%s, %s)", mod, member.Name())

	switch t := member.(type) {
	case *resourceType:
		return g.emitResourceType(mod, t)
	case *resourceFunc:
		return g.emitResourceFunc(mod, t)
	case *variable:
		contract.Assertf(mod.config(),
			"only expected top-level variables in config module (%s is not one)", mod.name)
		// skip the variable, we will process it later.
		return "", nil
	case *overlayFile:
		return g.emitOverlay(mod, t)
	default:
		contract.Failf("unexpected member type: %v", reflect.TypeOf(member))
		return "", nil
	}
}

// emitConfigVariables emits all config vaiables in the given module, returning the resulting file.
func (g *nodeJSGenerator) emitConfigVariables(mod *module) (string, error) {
	// Create a vars.ts file into which all configuration variables will go.
	w, config, err := g.openWriter(mod, "vars.ts", true)
	if err != nil {
		return "", err
	}
	defer contract.IgnoreClose(w)

	// Ensure we import any custom schemas referenced by the variables.
	var infos []*tfbridge.SchemaInfo
	for _, member := range mod.members {
		if v, ok := member.(*variable); ok {
			infos = append(infos, v.info)
		}
	}
	if err = g.emitCustomImports(w, mod, infos); err != nil {
		return "", err
	}

	// Create a config bag for the variables to pull from.
	w.Writefmtln("let __config = new pulumi.Config(\"%v:config\");", g.pkg)
	w.Writefmtln("")

	// Emit an entry for all config variables.
	for _, member := range mod.members {
		if v, ok := member.(*variable); ok {
			g.emitConfigVariable(w, v)
		}
	}

	return config, nil
}

func (g *nodeJSGenerator) emitConfigVariable(w *tools.GenWriter, v *variable) {
	var getfunc string
	if v.optional() {
		getfunc = "get"
	} else {
		getfunc = "require"
	}
	if v.schema.Type != schema.TypeString {
		// Only try to parse a JSON object if the config isn't a straight string.
		getfunc = fmt.Sprintf("%sObject<%s>", getfunc, tsType(v, false /*noflags*/))
	}
	var anycast string
	if v.info != nil && v.info.Type != "" {
		// If there's a custom type, we need to inject a cast to silence the compiler.
		anycast = "<any>"
	}
	if v.doc != "" {
		g.emitDocComment(w, v.doc, "")
	} else if v.rawdoc != "" {
		g.emitRawDocComment(w, v.rawdoc, "")
	}
	w.Writefmtln("export let %[1]s: %[2]s = %[3]s__config.%[4]s(\"%[1]s\");",
		v.name, tsType(v, true /*noflags*/), anycast, getfunc)
}

// sanitizeForDocComment ensures that no `*/` sequence appears in the string, to avoid
// accidentally closing the comment block.
func sanitizeForDocComment(str string) string {
	return strings.Replace(str, "*/", "*&#47;", -1)
}

func (g *nodeJSGenerator) emitDocComment(w *tools.GenWriter, comment, prefix string) {
	if comment != "" {
		lines := strings.Split(comment, "\n")
		w.Writefmtln("%v/**", prefix)
		for i, docLine := range lines {
			docLine = sanitizeForDocComment(docLine)
			// Break if we get to the last line and it's empty
			if i == len(lines)-1 && strings.TrimSpace(docLine) == "" {
				break
			}
			// Print the line of documentation
			w.Writefmtln("%v * %s", prefix, docLine)
		}
		w.Writefmtln("%v */", prefix)
	}
}

func (g *nodeJSGenerator) emitRawDocComment(w *tools.GenWriter, comment, prefix string) {
	if comment != "" {
		curr := 0
		w.Writefmtln("%v/**", prefix)
		w.Writefmt("%v * ", prefix)
		for _, word := range strings.Fields(comment) {
			word = sanitizeForDocComment(word)
			if curr > 0 {
				if curr+len(word)+1 > (maxWidth - len(prefix)) {
					curr = 0
					w.Writefmt("\n%v * ", prefix)
				} else {
					w.Writefmt(" ")
					curr++
				}
			}
			w.Writefmt(word)
			curr += len(word)
		}
		w.Writefmtln("")
		w.Writefmtln("%v */", prefix)
	}
}

func (g *nodeJSGenerator) emitPlainOldType(w *tools.GenWriter, pot *plainOldType) {
	if pot.doc != "" {
		g.emitDocComment(w, pot.doc, "")
	}
	w.Writefmtln("export interface %s {", pot.name)
	for _, prop := range pot.props {
		if prop.doc != "" {
			g.emitDocComment(w, prop.doc, "    ")
		} else if prop.rawdoc != "" {
			g.emitRawDocComment(w, prop.rawdoc, "    ")
		}
		w.Writefmtln("    readonly %s%s: %s;", prop.name, tsFlags(prop), tsType(prop, false))
	}
	w.Writefmtln("}")
}

func (g *nodeJSGenerator) emitResourceType(mod *module, res *resourceType) (string, error) {
	// Create a vars.ts file into which all configuration variables will go.
	w, file, err := g.openWriter(mod, lowerFirst(res.name)+".ts", true)
	if err != nil {
		return "", err
	}
	defer contract.IgnoreClose(w)

	// Ensure that we've emitted any custom imports pertaining to any of the field types.
	var fldinfos []*tfbridge.SchemaInfo
	for _, fldinfo := range res.info.Fields {
		fldinfos = append(fldinfos, fldinfo)
	}
	if err = g.emitCustomImports(w, mod, fldinfos); err != nil {
		return "", err
	}

	// Write the TypeDoc/JSDoc for the resource class
	if res.doc != "" {
		g.emitDocComment(w, res.doc, "")
	}

	// Begin defining the class.
	w.Writefmtln("export class %s extends pulumi.CustomResource {", res.name)

	// Emit all properties (using their output types).
	// TODO[pulumi/pulumi#397]: represent sensitive types using a Secret<T> type.
	ins := make(map[string]bool)
	for _, prop := range res.inprops {
		ins[prop.name] = true
	}
	for _, prop := range res.outprops {
		if prop.doc != "" {
			g.emitDocComment(w, prop.doc, "    ")
		} else if prop.rawdoc != "" {
			g.emitRawDocComment(w, prop.rawdoc, "    ")
		}

		// Make a little comment in the code so it's easy to pick out output properties.
		var outcomment string
		if !ins[prop.name] {
			outcomment = "/*out*/ "
		}

		// Emit the property as a computed value; it has to carry undefined because of planning.
		w.Writefmtln("    public %sreadonly %s%s: pulumi.Computed<%s>;",
			outcomment, prop.name, tsFlags(prop), tsType(prop, false))
	}
	w.Writefmtln("")

	// Now create a constructor that chains supercalls and stores into properties.
	w.Writefmtln("    /**")
	w.Writefmtln("     * Create a %s resource with the given unique name, arguments, and options.", res.name)
	w.Writefmtln("     *")
	w.Writefmtln("     * @param name The _unique_ name of the resource.")
	if res.argst != nil {
		w.Writefmtln("     * @param args The arguments to use to populate this resource.")
	}
	w.Writefmtln("     * @param opts A bag of options that control this resource's behavior.")
	w.Writefmtln("     */")

	var argsarg string
	if res.argst != nil {
		var argsflags string
		if len(res.reqprops) == 0 {
			// If the number of input properties was zero, we make the args object optional.
			argsflags = "?"
		}
		argsarg = fmt.Sprintf("args%s: %s, ", argsflags, res.argst.name)
	}

	w.Writefmtln("    constructor(name: string, %sopts?: pulumi.ResourceOptions) {", argsarg)

	// If the property arg isn't required, zero-init it if it wasn't actually passed in.
	if res.argst != nil && len(res.reqprops) == 0 {
		w.Writefmtln("        args = args || {};")
	}

	// First, validate all required arguments.
	for _, prop := range res.inprops {
		if !prop.optional() {
			w.Writefmtln("        if (args.%s === undefined) {", prop.name)
			w.Writefmtln("            throw new Error(\"Missing required property '%s'\");", prop.name)
			w.Writefmtln("        }")
		}
	}

	// Now invoke the super constructor with the type, name, and a property map.
	w.Writefmtln("        super(\"%s\", name, {", res.info.Tok)
	for _, prop := range res.inprops {
		w.Writefmtln("            \"%[1]s\": args.%[1]s,", prop.name)
	}
	for _, prop := range res.outprops {
		if !ins[prop.name] {
			w.Writefmtln("            \"%s\": undefined,", prop.name)
		}
	}
	w.Writefmtln("        }, opts);")

	w.Writefmtln("    }")
	w.Writefmtln("}")

	// If there's an argument type, emit it.
	if res.argst != nil {
		w.Writefmtln("")
		g.emitPlainOldType(w, res.argst)
	}

	return file, nil
}

func (g *nodeJSGenerator) emitResourceFunc(mod *module, fun *resourceFunc) (string, error) {
	// Create a vars.ts file into which all configuration variables will go.
	w, file, err := g.openWriter(mod, fun.name+".ts", true)
	if err != nil {
		return "", err
	}
	defer contract.IgnoreClose(w)

	// Ensure that we've emitted any custom imports pertaining to any of the field types.
	var fldinfos []*tfbridge.SchemaInfo
	for _, fldinfo := range fun.info.Fields {
		fldinfos = append(fldinfos, fldinfo)
	}
	if err = g.emitCustomImports(w, mod, fldinfos); err != nil {
		return "", err
	}

	// Write the TypeDoc/JSDoc for the data source function.
	if fun.doc != "" {
		g.emitDocComment(w, fun.doc, "")
	}

	// Now, emit the function signature.
	var argsig string
	if fun.argst != nil {
		var optflag string
		if len(fun.reqargs) == 0 {
			optflag = "?"
		}
		argsig = fmt.Sprintf("args%s: %s", optflag, fun.argst.name)
	}
	var retty string
	if fun.retst == nil {
		retty = "void"
	} else {
		retty = fun.retst.name
	}
	w.Writefmtln("export function %s(%s): Promise<%s> {", fun.name, argsig, retty)

	// Zero initialize the args if empty and necessary.
	if len(fun.args) > 0 && len(fun.reqargs) == 0 {
		w.Writefmtln("    args = args || {};")
	}

	// Now simply invoke the runtime function with the arguments, returning the results.
	w.Writefmtln("    return pulumi.runtime.invoke(\"%s\", {", fun.info.Tok)
	for _, arg := range fun.args {
		// Pass the argument to the invocation.
		w.Writefmtln("        \"%[1]s\": args.%[1]s,", arg.name)
	}
	w.Writefmtln("    });")

	w.Writefmtln("}")

	// If there are argument and/or return types, emit them.
	if fun.argst != nil {
		w.Writefmtln("")
		g.emitPlainOldType(w, fun.argst)
	}
	if fun.retst != nil {
		w.Writefmtln("")
		g.emitPlainOldType(w, fun.retst)
	}

	return file, nil
}

// emitOverlay copies an overlay from its source to the target, and returns the resulting file to be exported.
func (g *nodeJSGenerator) emitOverlay(mod *module, overlay *overlayFile) (string, error) {
	// Copy the file from the overlays directory to the destination.
	dir := g.moduleDir(mod)
	dst := filepath.Join(dir, overlay.name)
	if err := copyFile(overlay.src, dst); err != nil {
		return "", err
	}

	// And then export the overlay's contents from the index.
	return dst, nil
}

// emitPackageMetadata generates all the non-code metadata required by a Pulumi package.
func (g *nodeJSGenerator) emitPackageMetadata(pack *pkg, files []string) error {
	// The generator already emitted Pulumi.yaml, so that leaves two more files to write out:
	//     1) package.json: minimal NPM package metadata
	//     2) tsconfig.json: instructions for TypeScript compilation
	if err := g.emitNPMPackageMetadata(pack); err != nil {
		return err
	}
	return g.emitTypeScriptProjectFile(pack, files)
}

type npmPackage struct {
	Name             string            `json:"name"`
	Version          string            `json:"version"`
	Description      string            `json:"description,omitempty"`
	Keywords         []string          `json:"keywords,omitempty"`
	Homepage         string            `json:"homepage,omitempty"`
	Repository       string            `json:"repository,omitempty"`
	License          string            `json:"license,omitempty"`
	Scripts          map[string]string `json:"scripts,omitempty"`
	Dependencies     map[string]string `json:"dependencies,omitempty"`
	DevDependencies  map[string]string `json:"devDependencies,omitempty"`
	PeerDependencies map[string]string `json:"peerDependencies,omitempty"`
}

func (g *nodeJSGenerator) emitNPMPackageMetadata(pack *pkg) error {
	w, err := tools.NewGenWriter(tfgen, filepath.Join(g.outDir, "package.json"))
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)

	// Create info that will get serialized into an NPM package.json.
	npminfo := npmPackage{
		Name:        fmt.Sprintf("@pulumi/%s", pack.name),
		Version:     pack.version,
		Description: g.info.Description,
		Keywords:    g.info.Keywords,
		Homepage:    g.info.Homepage,
		Repository:  g.info.Repository,
		License:     g.info.License,
		Scripts: map[string]string{
			"build": "tsc",
		},
		DevDependencies: map[string]string{
			"typescript": "^2.6.2",
		},
		PeerDependencies: map[string]string{
			"pulumi": "*",
		},
	}

	// Copy the overlay dependencies, if any.
	if overlay := g.info.Overlay; overlay != nil {
		for depk, depv := range overlay.Dependencies {
			npminfo.Dependencies[depk] = depv
		}
		for depk, depv := range overlay.DevDependencies {
			npminfo.DevDependencies[depk] = depv
		}
		for depk, depv := range overlay.PeerDependencies {
			npminfo.PeerDependencies[depk] = depv
		}
	}

	// Now write out the serialized form.
	npmjson, err := json.MarshalIndent(npminfo, "", "    ")
	if err != nil {
		return err
	}
	w.Writefmtln(string(npmjson))
	return nil
}

func (g *nodeJSGenerator) emitTypeScriptProjectFile(pack *pkg, files []string) error {
	w, err := tools.NewGenWriter(tfgen, filepath.Join(g.outDir, "tsconfig.json"))
	if err != nil {
		return err
	}
	defer contract.IgnoreClose(w)
	w.Writefmtln(`{
    "compilerOptions": {
        "outDir": "bin",
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
		relfile, err := filepath.Rel(g.outDir, file)
		if err != nil {
			return err
		}
		w.Writefmtln("        \"%s\"%s", relfile, suffix)
	}
	w.Writefmtln(`    ]
}
`)
	return nil
}

// relModule removes the path suffix from a module and makes it relative to the root path.
func (g *nodeJSGenerator) relModule(mod *module, path string) (string, error) {
	// Return the path as a relative path to the root, so that imports are relative.
	dir := g.moduleDir(mod)
	file, err := filepath.Rel(dir, path)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(file, ".") {
		file = "./" + file
	}
	return removeExtension(file, ".ts"), nil
}

// removeExtension removes the file extension, if any.
func removeExtension(file, ext string) string {
	if strings.HasSuffix(file, ext) {
		return file[:len(file)-len(ext)]
	}
	return file
}

// importMap is a map of module name to a map of members imported.
type importMap map[string]map[string]bool

// emitCustomImports traverses a custom schema map, deeply, to figure out the set of imported names and files that
// will be required to access those names.  WARNING: this routine doesn't (yet) attempt to eliminate naming collisions.
func (g *nodeJSGenerator) emitCustomImports(w *tools.GenWriter, mod *module, infos []*tfbridge.SchemaInfo) error {
	// First gather up all imports into a map of import module to a list of imported members.
	imports := make(importMap)
	for _, info := range infos {
		if err := g.gatherCustomImports(mod, info, imports); err != nil {
			return err
		}
	}

	// Next, if there were any imports, generate the import statement.  We must sort names to ensure determinism.
	if len(imports) > 0 {
		var files []string
		for file := range imports {
			files = append(files, file)
		}
		sort.Strings(files)

		for _, file := range files {
			var names []string
			for name := range imports[file] {
				names = append(names, name)
			}
			sort.Strings(names)

			w.Writefmt("import {")
			for i, name := range names {
				if i > 0 {
					w.Writefmt(", ")
				}
				w.Writefmt(name)
			}
			w.Writefmtln("} from \"%v\";", file)
		}
		w.Writefmtln("")
	}
	return nil
}

// gatherCustomImports gathers imports from an entire map of schema info, and places them into the target map.
func (g *nodeJSGenerator) gatherCustomImports(mod *module, info *tfbridge.SchemaInfo, imports importMap) error {
	if info != nil {
		// If this property has custom schema types that aren't "simple" (e.g., string, etc), then we need to
		// create a relative module import.  Note that we assume this is local to the current package!
		var custty []tokens.Type
		if info.Type != "" {
			custty = append(custty, info.Type)
			custty = append(custty, info.AltTypes...)
		}
		for _, ct := range custty {
			if !tokens.Token(ct).Simple() {
				// Make a relative module import, based on the module we are importing within.
				haspkg := string(ct.Module().Package().Name())
				if haspkg != g.pkg {
					return errors.Errorf("custom schema type %s was not in the current package %s", haspkg, g.pkg)
				}
				modname := ct.Module().Name()
				modfile := filepath.Join(g.outDir,
					strings.Replace(string(modname), tokens.TokenDelimiter, string(filepath.Separator), -1))
				relmod, err := g.relModule(mod, modfile)
				if err != nil {
					return err
				}

				// Now just mark the member in the resulting map.
				if imports[relmod] == nil {
					imports[relmod] = make(map[string]bool)
				}
				imports[relmod][string(ct.Name())] = true
			}
		}

		// If the property has an element type, recurse and propagate any results.
		if err := g.gatherCustomImports(mod, info.Elem, imports); err != nil {
			return err
		}

		// If the property has fields, then simply recurse and propagate any results, if any, to our map.
		for _, info := range info.Fields {
			if err := g.gatherCustomImports(mod, info, imports); err != nil {
				return err
			}
		}
	}

	return nil
}

// tsFlags returns the TypeScript flags for a given variable.
func tsFlags(v *variable) string {
	return tsFlagsComplex(v.schema, v.info, v.out)
}

// tsFlagsComplex is just like tsFlags, except that it permits recursing into component pieces individually.
func tsFlagsComplex(sch *schema.Schema, info *tfbridge.SchemaInfo, out bool) string {
	if optionalComplex(sch, info, out) {
		return "?"
	}
	return ""
}

// tsType returns the TypeScript type name for a given schema property.  noflags may be passed as true to create a
// type that represents the optional nature of a variable, even when flags will not be present; this is often needed
// when turning the type into a generic type argument, for example, since there will be no opportunity for "?" there.
func tsType(v *variable, noflags bool) string {
	return tsTypeComplex(v.schema, v.info, noflags, v.out)
}

// tsTypeComplex is just like tsType, but permits recursing using component pieces rather than a true variable.
func tsTypeComplex(sch *schema.Schema, info *tfbridge.SchemaInfo, noflags, out bool) string {
	// First, see if there is a custom override.  If yes, use it directly.
	var t string
	var elem *tfbridge.SchemaInfo
	if info != nil {
		if info.Type != "" {
			t = string(info.Type.Name())
			if len(info.AltTypes) > 0 {
				for _, at := range info.AltTypes {
					t = fmt.Sprintf("%s | %s", t, at.Name())
				}
			}
			if !out {
				t = fmt.Sprintf("pulumi.ComputedValue<%s>", t)
			}
		} else if info.Asset != nil {
			t = "pulumi.asset." + info.Asset.Type()
		}

		elem = info.Elem
	}

	// If nothing was found, generate the primitive type name for this.
	if t == "" {
		t = tsPrimitive(sch.Type, sch.Elem, elem, out)
	}

	// If we aren't using optional flags, we need to use TypeScript union types to permit undefined values.
	if noflags {
		if opt := optionalComplex(sch, info, out); opt {
			t += " | undefined"
		}
	}

	return t
}

// tsPrimitive returns the TypeScript type name for a given schema value type and element kind.
func tsPrimitive(vt schema.ValueType, elem interface{}, eleminfo *tfbridge.SchemaInfo, out bool) string {
	// First figure out the raw type.
	var t string
	var array bool
	switch vt {
	case schema.TypeBool:
		t = "boolean"
	case schema.TypeInt, schema.TypeFloat:
		t = "number"
	case schema.TypeString:
		t = "string"
	case schema.TypeList:
		t = tsElemType(elem, eleminfo, out)
		array = true
	case schema.TypeMap:
		t = fmt.Sprintf("{[key: string]: %v}", tsElemType(elem, eleminfo, out))
	case schema.TypeSet:
		// IDEA: we can't use ES6 sets here, because we're using values and not objects.  It would be possible to come
		//     up with a ValueSet of some sorts, but that depends on things like shallowEquals which is known to be
		//     brittle and implementation dependent.  For now, we will stick to arrays, and validate on the backend.
		t = tsElemType(elem, eleminfo, out)
		array = true
	default:
		contract.Failf("Unrecognized schema type: %v", vt)
	}

	// Now, if it is an input property value, it must be wrapped in a ComputedValue<T>.
	if !out {
		t = fmt.Sprintf("pulumi.ComputedValue<%s>", t)
	}

	// Finally make sure arrays are arrays; this must be done after the above, so we get a ComputedValue<T>[],
	// and not a ComputedValue<T[]>, which would constrain the ability to flexibly construct them.
	// BUGBUG[pulumi/pulumi-terraform#47]: this code needs to be removed -- it's just wrong.
	if array {
		t = fmt.Sprintf("%s[]", t)
	}

	return t
}

// tsElemType returns the TypeScript type for a given schema element.  This element may be either a simple schema
// property or a complex structure.  In the case of a complex structure, this will expand to its nominal type.
func tsElemType(elem interface{}, info *tfbridge.SchemaInfo, out bool) string {
	// If there is no element type specified, we will accept anything.
	if elem == nil {
		return "any"
	}

	switch e := elem.(type) {
	case schema.ValueType:
		return tsPrimitive(e, nil, info, out)
	case *schema.Schema:
		// A simple type, just return its type name.
		return tsTypeComplex(e, info, true /*noflags*/, out)
	case *schema.Resource:
		// A complex type, just expand to its nominal type name.
		// TODO: spill all complex structures in advance so that we don't have insane inline expansions.
		t := "{ "
		c := 0
		for _, s := range stableSchemas(e.Schema) {
			var fldinfo *tfbridge.SchemaInfo
			if info != nil {
				fldinfo = info.Fields[s]
			}
			sch := e.Schema[s]
			if name := propertyName(s, fldinfo); name != "" {
				if c > 0 {
					t += ", "
				}
				flg := tsFlagsComplex(sch, fldinfo, out)
				typ := tsTypeComplex(sch, fldinfo, false /*noflags*/, out)
				t += fmt.Sprintf("%s%s: %s", name, flg, typ)
				c++
			}
		}
		return t + " }"
	default:
		contract.Failf("Unrecognized schema element type: %v", e)
		return ""
	}
}
