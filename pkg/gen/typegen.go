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

package gen

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/go-wordwrap"
	"github.com/pulumi/pulumi-kubernetes/pkg/kinds"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"

	pycodegen "github.com/pulumi/pulumi/pkg/codegen/python"
)

const (
	tsObject  = "object"
	tsStringT = "string"
	pyDictT   = "dict"
	pyStringT = "str"
	pyIntT    = "int"
	pyListT   = "list"
	pyBoolT   = "bool"
	pyAnyT    = "Any"
)

const (
	apiRegistration = "apiregistration.k8s.io"
)

// --------------------------------------------------------------------------

// A collection of data structures and utility functions to transform an OpenAPI spec for the
// Kubernetes API into something that we can use for codegen'ing nodejs and Python clients.

// --------------------------------------------------------------------------

// GroupConfig represents a Kubernetes API group (e.g., core, apps, extensions, etc.)
type GroupConfig struct {
	group    string
	versions []*VersionConfig
}

// Group returns the name of the group (e.g., `core` for core, etc.)
func (gc *GroupConfig) Group() string { return gc.group }

// Versions returns the set of version for some Kubernetes API group. For example, the `apps` group
// has `v1beta1`, `v1beta2`, and `v1`.
func (gc *GroupConfig) Versions() []*VersionConfig { return gc.versions }

// VersionConfig represents a version of a Kubernetes API group (e.g., the `apps` group has
// `v1beta1`, `v1beta2`, and `v1`.)
type VersionConfig struct {
	version string
	kinds   []*KindConfig

	gv            *schema.GroupVersion // Used for sorting.
	apiVersion    string
	rawAPIVersion string
}

// Version returns the name of the version (e.g., `apps/v1beta1` would return `v1beta1`).
func (vc *VersionConfig) Version() string { return vc.version }

// Kinds returns the set of kinds in some Kubernetes API group/version combination (e.g.,
// `apps/v1beta1` has the `Deployment` kind, etc.).
func (vc *VersionConfig) Kinds() []*KindConfig { return vc.kinds }

// KindsAndAliases will produce a list of kinds, including aliases (e.g., both `apiregistration` and
// `apiregistration.k8s.io`).
func (vc *VersionConfig) KindsAndAliases() []*KindConfig {
	kindsAndAliases := []*KindConfig{}
	for _, kind := range vc.kinds {
		kindsAndAliases = append(kindsAndAliases, kind)
		if strings.HasPrefix(kind.APIVersion(), apiRegistration) {
			alias := KindConfig{}
			err := copier.Copy(&alias, kind)
			if err != nil {
				panic(err)
			}
			rawAPIVersion := "apiregistration" + strings.TrimPrefix(kind.APIVersion(), apiRegistration)
			alias.rawAPIVersion = rawAPIVersion
			kindsAndAliases = append(kindsAndAliases, &alias)
		}
	}
	return kindsAndAliases
}

// ListKindsAndAliases will return all known `Kind`s that are lists, or aliases of lists. These
// `Kind`s are not instantiated by the API server, and we must "flatten" them client-side to get an
// accurate view of what resource operations we need to perform.
func (vc *VersionConfig) ListKindsAndAliases() []*KindConfig {
	listKinds := []*KindConfig{}
	for _, kind := range vc.KindsAndAliases() {
		hasItems := false
		for _, prop := range kind.properties {
			if prop.name == "items" {
				hasItems = true
				break
			}
		}

		if strings.HasSuffix(kind.Kind(), "List") && hasItems {
			listKinds = append(listKinds, kind)
		}
	}

	return listKinds
}

// APIVersion returns the fully-qualified apiVersion (e.g., `storage.k8s.io/v1` for storage, etc.)
func (vc *VersionConfig) APIVersion() string { return vc.apiVersion }

// RawAPIVersion returns the "raw" apiVersion (e.g., `v1` rather than `core/v1`).
func (vc *VersionConfig) RawAPIVersion() string { return vc.rawAPIVersion }

// KindConfig represents a Kubernetes API kind (e.g., the `Deployment` type in
// `apps/v1beta1/Deployment`).
type KindConfig struct {
	kind                    string
	comment                 string
	pulumiComment           string
	properties              []*Property
	requiredInputProperties []*Property
	optionalInputProperties []*Property
	additionalSecretOutputs []string
	aliases                 []string

	gvk           *schema.GroupVersionKind // Used for sorting.
	apiVersion    string
	rawAPIVersion string
	typeGuard     string
}

// Kind returns the name of the Kubernetes API kind (e.g., `Deployment` for
// `apps/v1beta1/Deployment`).
func (kc *KindConfig) Kind() string { return kc.kind }

// Comment returns the comments associated with some Kubernetes API kind.
func (kc *KindConfig) Comment() string { return kc.comment }

// PulumiComment returns the await logic documentation associated with some Kubernetes API kind.
func (kc *KindConfig) PulumiComment() string { return kc.pulumiComment }

// Properties returns the list of properties that exist on some Kubernetes API kind (i.e., things
// that we will want to `.` into, like `thing.apiVersion`, `thing.kind`, `thing.metadata`, etc.).
func (kc *KindConfig) Properties() []*Property { return kc.properties }

// RequiredInputProperties returns the list of properties that are required input properties on some
// Kubernetes API kind (i.e., things that we will want to provide, like `thing.metadata`, etc.).
func (kc *KindConfig) RequiredInputProperties() []*Property { return kc.requiredInputProperties }

// OptionalInputProperties returns the list of properties that are optional input properties on some
// Kubernetes API kind (i.e., things that we will want to provide, like `thing.metadata`, etc.).
func (kc *KindConfig) OptionalInputProperties() []*Property { return kc.optionalInputProperties }

// AdditionalSecretOutputs returns the list of strings to set as additionalSecretOutputs on some
// Kubernetes API kind.
func (kc *KindConfig) AdditionalSecretOutputs() []string { return kc.additionalSecretOutputs }

// Aliases returns the list of aliases for a Kubernetes API kind.
func (kc *KindConfig) Aliases() []string { return kc.aliases }

// APIVersion returns the fully-qualified apiVersion (e.g., `storage.k8s.io/v1` for storage, etc.)
func (kc *KindConfig) APIVersion() string { return kc.apiVersion }

// RawAPIVersion returns the "raw" apiVersion (e.g., `v1` rather than `core/v1`).
func (kc *KindConfig) RawAPIVersion() string { return kc.rawAPIVersion }

// URNAPIVersion returns API version that can be used in a URN (e.g., using the backwards-compatible
// alias `apiextensions` instead of `apiextensions.k8s.io`).
func (kc *KindConfig) URNAPIVersion() string {
	if strings.HasPrefix(kc.apiVersion, apiRegistration) {
		return "apiregistration" + strings.TrimPrefix(kc.apiVersion, apiRegistration)
	}
	return kc.apiVersion
}

// TypeGuard returns the text of a TypeScript type guard for the given kind.
func (kc *KindConfig) TypeGuard() string { return kc.typeGuard }

// Property represents a property we want to expose on a Kubernetes API kind (i.e., things that we
// will want to `.` into, like `thing.apiVersion`, `thing.kind`, `thing.metadata`, etc.).
type Property struct {
	name                      string
	languageName              string
	comment                   string
	pythonConstructorComment  string
	propType                  string
	pythonConstructorPropType string
	defaultValue              string
	isLast                    bool
	dotnetVarName             string
	dotnetIsListOrMap         bool
}

// Name returns the name of the property.
func (p *Property) Name() string { return p.name }

// LanguageName returns the name of the property.
func (p *Property) LanguageName() string { return p.languageName }

// Comment returns the comments associated with some property.
func (p *Property) Comment() string { return p.comment }

// PythonConstructorComment returns the comments associated with some property, formatted for Python
// constructor documentation.
func (p *Property) PythonConstructorComment() string { return p.pythonConstructorComment }

// PropType returns the type of the property.
func (p *Property) PropType() string { return p.propType }

// PythonConstructorPropType returns the type of the property, typed for the Python constructor
// resource inputs.
func (p *Property) PythonConstructorPropType() string { return p.pythonConstructorPropType }

// DefaultValue returns the type of the property.
func (p *Property) DefaultValue() string { return p.defaultValue }

// IsLast returns whether the property is the last in the list of properties.
func (p *Property) IsLast() bool { return p.isLast }

// DotnetVarName returns a variable name safe to use in .NET (e.g. `@namespace` instead of `namespace`)
func (p *Property) DotnetVarName() string { return p.dotnetVarName }

// DotnetIsListOrMap returns whether the property type is a List or map
func (p *Property) DotnetIsListOrMap() bool { return p.dotnetIsListOrMap }

// --------------------------------------------------------------------------

// Utility functions.

// --------------------------------------------------------------------------

func gvkFromRef(ref string) schema.GroupVersionKind {
	// TODO(hausdorff): Surely there is an official k8s function somewhere for doing this.
	split := strings.Split(ref, ".")
	return schema.GroupVersionKind{
		Kind:    split[len(split)-1],
		Version: split[len(split)-2],
		Group:   split[len(split)-3],
	}
}

func stripPrefix(name string) string {
	const prefix = "#/definitions/"
	return strings.TrimPrefix(name, prefix)
}

func replaceDeprecationComment(comment string, gvk schema.GroupVersionKind, language language) string {
	// The deprecation warning doesn't always appear in the same place in the OpenAPI comments.
	// Standardize the message and where it appears in our docs.
	re1 := regexp.MustCompile(`^DEPRECATED - .* is deprecated by .* for more information\.\s*`)
	re2 := regexp.MustCompile(`DEPRECATED - .* is deprecated by .* for more information\.\s*`)

	var prefix, replacement string
	switch language {
	case typescript:
		prefix = "@deprecated "
		replacement = prefix + ApiVersionComment(gvk)
	case python, dotnet:
		prefix = "DEPRECATED - "
		replacement = prefix + ApiVersionComment(gvk)
	default:
		panic(fmt.Sprintf("Unsupported language '%s'", language))
	}

	if re1.MatchString(comment) {
		return re1.ReplaceAllString(comment, replacement)
	} else if re2.MatchString(comment) {
		return prefix + ApiVersionComment(gvk) + re2.ReplaceAllString(comment, "")
	} else {
		return comment
	}
}

func fmtComment(
	comment interface{}, prefix string, bareRender bool, opts groupOpts, gvk schema.GroupVersionKind,
) string {
	if comment == nil {
		return ""
	}

	var wrapParagraph func(line string) []string
	var renderComment func(lines []string) string
	switch opts.language {
	case python:
		wrapParagraph = func(paragraph string) []string {
			borderLen := len(prefix)
			wrapped := wordwrap.WrapString(paragraph, 100-uint(borderLen))
			return strings.Split(wrapped, "\n")
		}
		renderComment = func(lines []string) string {
			joined := strings.Join(lines, fmt.Sprintf("\n%s", prefix))
			if !bareRender {
				return fmt.Sprintf("\"\"\"\n%s%s\n%s\"\"\"", prefix, joined, prefix)
			}
			return joined
		}
	case typescript:
		wrapParagraph = func(paragraph string) []string {
			// Escape comment termination.
			escaped := strings.Replace(paragraph, "*/", "*&#8205;/", -1)
			borderLen := len(prefix + " * ")
			wrapped := wordwrap.WrapString(escaped, 100-uint(borderLen))
			return strings.Split(wrapped, "\n")
		}
		renderComment = func(lines []string) string {
			joined := strings.Join(lines, fmt.Sprintf("\n%s * ", prefix))
			if !bareRender {
				return fmt.Sprintf("/**\n%s * %s\n%s */", prefix, joined, prefix)
			}
			return joined
		}
	case dotnet:
		wrapParagraph = func(paragraph string) []string {
			escaped := strings.Replace(paragraph, "<", "&lt;", -1)
			escaped = strings.Replace(escaped, ">", "&gt;", -1)
			escaped = strings.Replace(escaped, "&", "&amp;", -1)
			borderLen := len(prefix + "/// ")
			wrapped := wordwrap.WrapString(escaped, 100-uint(borderLen))
			return strings.Split(wrapped, "\n")
		}
		renderComment = func(lines []string) string {
			joined := strings.Join(lines, fmt.Sprintf("\n%s/// ", prefix))
			if !bareRender {
				return fmt.Sprintf("/// <summary>\n%s/// %s\n%s/// </summary>", prefix, joined, prefix)
			}
			return joined
		}
	default:
		panic(fmt.Sprintf("Unsupported language '%s'", opts.language))
	}

	commentstr, _ := comment.(string)
	if len(commentstr) > 0 {

		// hack(levi): The OpenAPI docs currently include broken links to k8s docs. Until this is fixed
		// upstream, manually replace these with working links.
		// Upstream issue: https://github.com/kubernetes/kubernetes/issues/81526
		// Upstream PR: https://github.com/kubernetes/kubernetes/pull/74245
		commentstr = strings.Replace(
			commentstr,
			`https://git.k8s.io/community/contributors/devel/api-conventions.md`,
			`https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md`,
			-1)

		commentstr = replaceDeprecationComment(commentstr, gvk, opts.language)

		split := strings.Split(commentstr, "\n")
		var lines []string
		for _, paragraph := range split {
			lines = append(lines, wrapParagraph(paragraph)...)
		}
		return renderComment(lines)
	}
	return ""
}

const (
	apiextensionsV1beta1          = "io.k8s.apiextensions-apiserver.pkg.apis.apiextensions.v1beta1"
	apiextensionsV1               = "io.k8s.apiextensions-apiserver.pkg.apis.apiextensions.v1"
	quantity                      = "io.k8s.apimachinery.pkg.api.resource.Quantity"
	rawExtension                  = "io.k8s.apimachinery.pkg.runtime.RawExtension"
	intOrString                   = "io.k8s.apimachinery.pkg.util.intstr.IntOrString"
	v1Fields                      = "io.k8s.apimachinery.pkg.apis.meta.v1.Fields"
	v1FieldsV1                    = "io.k8s.apimachinery.pkg.apis.meta.v1.FieldsV1"
	v1Time                        = "io.k8s.apimachinery.pkg.apis.meta.v1.Time"
	v1MicroTime                   = "io.k8s.apimachinery.pkg.apis.meta.v1.MicroTime"
	v1beta1JSONSchemaPropsOrBool  = apiextensionsV1beta1 + ".JSONSchemaPropsOrBool"
	v1beta1JSONSchemaPropsOrArray = apiextensionsV1beta1 + ".JSONSchemaPropsOrArray"
	v1beta1JSON                   = apiextensionsV1beta1 + ".JSON"
	v1beta1CRSubresourceStatus    = apiextensionsV1beta1 + ".CustomResourceSubresourceStatus"
	v1JSONSchemaPropsOrBool       = apiextensionsV1 + ".JSONSchemaPropsOrBool"
	v1JSONSchemaPropsOrArray      = apiextensionsV1 + ".JSONSchemaPropsOrArray"
	v1JSON                        = apiextensionsV1 + ".JSON"
	v1CRSubresourceStatus         = apiextensionsV1 + ".CustomResourceSubresourceStatus"
)

func makeTypescriptType(resourceType, propName string, prop map[string]interface{}, opts groupOpts) string {
	wrapType := func(typ string) string {
		switch opts.generatorType {
		case provider:
			return fmt.Sprintf("pulumi.Output<%s>", typ)
		case outputsAPI:
			return typ
		case inputsAPI:
			return fmt.Sprintf("pulumi.Input<%s>", typ)
		default:
			panic(fmt.Sprintf("unrecognized generator type %d", opts.generatorType))
		}
	}

	refPrefix := ""
	if opts.generatorType == provider {
		refPrefix = "outputs"
	}

	if t, exists := prop["type"]; exists {
		tstr := t.(string)
		if tstr == "array" {
			elemType := makeTypescriptType(
				resourceType, propName, prop["items"].(map[string]interface{}), opts)
			switch opts.generatorType {
			case provider:
				return fmt.Sprintf("%s[]>", elemType[:len(elemType)-1])
			case outputsAPI:
				return fmt.Sprintf("%s[]", elemType)
			case inputsAPI:
				return fmt.Sprintf("pulumi.Input<%s[]>", elemType)
			}
		} else if tstr == "integer" {
			return wrapType("number")
		} else if tstr == tsObject {
			// `additionalProperties` with a single member, `type`, denotes a map whose keys and
			// values both have type `type`. This type is never a `$ref`.
			if additionalProperties, exists := prop["additionalProperties"]; exists {
				mapType := additionalProperties.(map[string]interface{})
				if ktype, exists := mapType["type"]; exists && len(mapType) == 1 {
					switch opts.generatorType {
					case inputsAPI:
						return fmt.Sprintf("pulumi.Input<{[key: %s]: pulumi.Input<%s>}>", ktype, ktype)
					case outputsAPI:
						return fmt.Sprintf("{[key: %s]: %s}", ktype, ktype)
					case provider:
						return fmt.Sprintf("pulumi.Output<{[key: %s]: pulumi.Output<%s>}>", ktype, ktype)
					}
				}
			}
		} else if tstr == "string" && resourceType == "io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta" && propName == "namespace" {
			// Special case: `.metadata.namespace` should either take a string or a namespace object
			// itself.

			switch opts.generatorType {
			case inputsAPI:
				// TODO: Enable metadata to take explicit namespaces, like:
				// return "pulumi.Input<string> | Namespace"
				return "pulumi.Input<string>"
			case outputsAPI:
				return "string"
			}
		}

		return wrapType(tstr)
	}

	ref := stripPrefix(prop["$ref"].(string))

	isSimpleRef := true
	switch ref {
	case quantity:
		ref = tsStringT
	case intOrString:
		ref = "number | string"
	case v1Fields, v1FieldsV1, rawExtension:
		ref = tsObject
	case v1Time, v1MicroTime:
		// TODO: Automatically deserialized with `DateConstructor`.
		ref = tsStringT
	case v1beta1JSONSchemaPropsOrBool:
		ref = "apiextensions.v1beta1.JSONSchemaProps | boolean"
	case v1JSONSchemaPropsOrBool:
		ref = "apiextensions.v1.JSONSchemaProps | boolean"
	case v1beta1JSONSchemaPropsOrArray:
		ref = "apiextensions.v1beta1.JSONSchemaProps | any[]"
	case v1JSONSchemaPropsOrArray:
		ref = "apiextensions.v1.JSONSchemaProps | any[]"
	case v1beta1JSON, v1beta1CRSubresourceStatus, v1JSON, v1CRSubresourceStatus:
		ref = "any"
	default:
		isSimpleRef = false
	}

	if isSimpleRef {
		return wrapType(ref)
	}

	gvk := gvkFromRef(ref)
	var gvkRefStr string
	if refPrefix == "" {
		gvkRefStr = fmt.Sprintf("%s.%s.%s", gvk.Group, gvk.Version, gvk.Kind)
	} else {
		gvkRefStr = fmt.Sprintf("%s.%s.%s.%s", refPrefix, gvk.Group, gvk.Version, gvk.Kind)
	}

	return wrapType(gvkRefStr)
}

func makePythonType(resourceType, propName string, prop map[string]interface{}, opts groupOpts) string {
	wrapType := func(typ string) string {
		switch opts.generatorType {
		case provider:
			return fmt.Sprintf("pulumi.Output[%s]", typ)
		case outputsAPI:
			return typ
		case inputsAPI:
			return fmt.Sprintf("pulumi.Input[%s]", typ)
		default:
			panic(fmt.Sprintf("unrecognized generator type %d", opts.generatorType))
		}
	}

	if t, exists := prop["type"]; exists {
		tstr := t.(string)
		if tstr == "array" {
			return wrapType(pyListT)
		} else if tstr == "integer" {
			return wrapType(pyIntT)
		} else if tstr == tsObject {
			return wrapType(pyDictT)
		} else if tstr == "string" {
			return wrapType(pyStringT)
		} else if tstr == "boolean" {
			return wrapType(pyBoolT)
		}

		return wrapType(tstr)
	}

	ref := stripPrefix(prop["$ref"].(string))

	switch ref {
	case quantity:
		ref = pyStringT
	case intOrString:
		ref = fmt.Sprintf("Union[%s, %s]", pyIntT, pyStringT)
	case v1Fields, v1FieldsV1, rawExtension:
		ref = pyDictT
	case v1Time, v1MicroTime:
		// TODO: Automatically deserialized with `DateConstructor`.
		ref = pyStringT
	case v1beta1JSONSchemaPropsOrBool, v1JSONSchemaPropsOrBool:
		ref = fmt.Sprintf("Union[%s, %s]", pyDictT, pyBoolT)
	case v1beta1JSONSchemaPropsOrArray, v1JSONSchemaPropsOrArray:
		ref = fmt.Sprintf("Union[%s, %s]", pyDictT, pyListT)
	case v1beta1JSON, v1beta1CRSubresourceStatus, v1JSON, v1CRSubresourceStatus:
		ref = pyAnyT
	default:
		ref = pyDictT
	}

	return wrapType(ref)
}

func makeDotnetType(resourceType, propName string, prop map[string]interface{}, opts groupOpts) string {
	wrapType := func(typ string) string {
		switch opts.generatorType {
		case provider:
			return fmt.Sprintf("Output<%s>", typ)
		case outputsAPI:
			return typ
		case inputsAPI:
			return fmt.Sprintf("Input<%s>", typ)
		default:
			panic(fmt.Sprintf("unrecognized generator type %d", opts.generatorType))
		}
	}

	refPrefix := ""
	if opts.generatorType == provider {
		refPrefix = "Types.Outputs"
	}

	if t, exists := prop["type"]; exists {
		tstr := t.(string)
		if tstr == "array" {
			elemType := makeDotnetType(
				resourceType, propName, prop["items"].(map[string]interface{}), opts)
			switch opts.generatorType {
			case provider:
				return fmt.Sprintf("%s[]>", elemType[:len(elemType)-1])
			case outputsAPI:
				return fmt.Sprintf("ImmutableArray<%s>", elemType)
			case inputsAPI:
				return fmt.Sprintf("InputList<%s>", elemType)
			}
		} else if tstr == "integer" {
			return wrapType("int")
		} else if tstr == "boolean" {
			return wrapType("bool")
		} else if tstr == "number" {
			return wrapType("double")
		} else if tstr == "object" {
			kvtype := "string"
			// `additionalProperties` with a single member, `type`, denotes a map whose keys and
			// values both have type `type`. This type is never a `$ref`.
			if additionalProperties, exists := prop["additionalProperties"]; exists {
				mapType := additionalProperties.(map[string]interface{})
				if ktype, exists := mapType["type"]; exists && len(mapType) == 1 {
					kvtype = ktype.(string)

				}
			}
			switch opts.generatorType {
			case inputsAPI:
				return fmt.Sprintf("InputMap<%s>", kvtype)
			case outputsAPI:
				return fmt.Sprintf("ImmutableDictionary<%s, %s>", kvtype, kvtype)
			case provider:
				return fmt.Sprintf("Output<ImmutableDictionary<%s, %s>>", kvtype, kvtype)
			}
		} else if tstr == "string" && resourceType == "io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta" && propName == "namespace" {
			// Special case: `.metadata.namespace` should either take a string or a namespace object
			// itself.

			switch opts.generatorType {
			case inputsAPI:
				// TODO: Enable metadata to take explicit namespaces, like:
				// return "pulumi.Input<string> | Namespace"
				return "Input<string>"
			case outputsAPI:
				return "string"
			}
		}

		return wrapType(tstr)
	}

	ref := stripPrefix(prop["$ref"].(string))

	isSimpleRef := true
	switch ref {
	case quantity:
		ref = "string"
	case intOrString:
		ref = "int /* TODO: or string */"
	case v1Fields, v1FieldsV1, rawExtension:
		ref = "string /* TODO: wrong!*/"
	case v1Time, v1MicroTime:
		// TODO: Automatically deserialized with `DateConstructor`.
		ref = "string"
	case v1beta1JSONSchemaPropsOrBool:
		ref = "ApiExtensions.V1Beta1.JSONSchemaProps /* TODO: or bool */"
	case v1JSONSchemaPropsOrBool:
		ref = "ApiExtensions.V1.JSONSchemaProps /* TODO: or bool */"
	case v1beta1JSONSchemaPropsOrArray:
		ref = "ApiExtensions.V1Beta1.JSONSchemaProps /* TODO: or array */"
	case v1JSONSchemaPropsOrArray:
		ref = "ApiExtensions.V1.JSONSchemaProps /* TODO: or array */"
	case v1beta1JSON, v1beta1CRSubresourceStatus, v1JSON, v1CRSubresourceStatus:
		ref = "string /* TODO: wrong!*/"
	default:
		isSimpleRef = false
	}

	if isSimpleRef {
		return wrapType(ref)
	}

	gvk := gvkFromRef(ref)
	group := pascalCase(gvk.Group)
	version := pascalCase(gvk.Version)

	var gvkRefStr string
	if refPrefix == "" {
		gvkRefStr = fmt.Sprintf("%s.%s.%s", group, version, gvk.Kind)
	} else {
		gvkRefStr = fmt.Sprintf("%s.%s.%s.%s", refPrefix, group, version, gvk.Kind)
	}

	return wrapType(gvkRefStr)
}

func makeType(resourceType, propName string, prop map[string]interface{}, opts groupOpts) string {
	switch opts.language {
	case typescript:
		return makeTypescriptType(resourceType, propName, prop, opts)
	case python:
		return makePythonType(resourceType, propName, prop, opts)
	case dotnet:
		return makeDotnetType(resourceType, propName, prop, opts)
	default:
		panic(fmt.Sprintf("Unsupported language '%s'", opts.language))
	}
}

func isTopLevel(d *definition) bool {
	gvks, gvkExists :=
		d.data["x-kubernetes-group-version-kind"].([]interface{})
	hasGVK := gvkExists && len(gvks) > 0
	if !hasGVK {
		return false
	}

	// Return `false` for the handful of top-level imperative resource types that can't be managed
	// by Pulumi.
	switch fmt.Sprintf("%s/%s", d.gvk.GroupVersion().String(), d.gvk.Kind) {
	case "policy/v1beta1/Eviction", "v1/Status", "apps/v1beta1/Scale", "apps/v1beta2/Scale",
		"autoscaling/v1/Scale", "extensions/v1beta1/Scale":
		return false
	}

	properties, hasProperties := d.data["properties"].(map[string]interface{})
	if !hasProperties {
		return false
	}

	meta, hasMetadata := properties["metadata"].(map[string]interface{})
	if !hasMetadata {
		return false
	}

	ref, hasRef := meta["$ref"]
	if !hasRef {
		return false
	}

	return ref == "#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta" ||
		ref == "#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.ListMeta"
}

// --------------------------------------------------------------------------

// Core grouping logic.

// --------------------------------------------------------------------------

type definition struct {
	gvk  schema.GroupVersionKind
	name string
	data map[string]interface{}
}

type gentype int

const (
	provider gentype = iota
	inputsAPI
	outputsAPI
)

type language string

const (
	python     = "python"
	typescript = "typescript"
	dotnet     = "dotnet"
)

type groupOpts struct {
	generatorType gentype
	language      language
}

func nodeJSInputs() groupOpts   { return groupOpts{generatorType: inputsAPI, language: typescript} }
func nodeJSOutputs() groupOpts  { return groupOpts{generatorType: outputsAPI, language: typescript} }
func nodeJSProvider() groupOpts { return groupOpts{generatorType: provider, language: typescript} }

func pythonProvider() groupOpts { return groupOpts{generatorType: provider, language: python} }

func dotnetInputs() groupOpts   { return groupOpts{generatorType: inputsAPI, language: dotnet} }
func dotnetOutputs() groupOpts  { return groupOpts{generatorType: outputsAPI, language: dotnet} }
func dotnetProvider() groupOpts { return groupOpts{generatorType: provider, language: dotnet} }

func allCamelCasePropertyNames(definitionsJSON map[string]interface{}, opts groupOpts) []string {
	// Map definition JSON object -> `definition` with metadata.
	definitions := make([]*definition, 0)
	linq.From(definitionsJSON).
		WhereT(func(kv linq.KeyValue) bool {
			// Skip these objects, special case. They're deprecated and empty.
			defName := kv.Key.(string)
			return !strings.HasPrefix(defName, "io.k8s.kubernetes.pkg")
		}).
		SelectT(func(kv linq.KeyValue) *definition {
			defName := kv.Key.(string)
			return &definition{
				gvk:  gvkFromRef(defName),
				name: defName,
				data: definitionsJSON[defName].(map[string]interface{}),
			}
		}).
		ToSlice(&definitions)

	properties := sets.String{}
	// Only select camel-cased property names
	re := regexp.MustCompile(`[a-z]+[A-Z]`)
	for _, d := range definitions {
		if pmap, exists := d.data["properties"]; exists {
			ps := pmap.(map[string]interface{})
			for p := range ps {
				if re.MatchString(p) {
					properties.Insert(p)
				}
			}
		}
	}

	return properties.List()
}

func createGroups(definitionsJSON map[string]interface{}, opts groupOpts) []*GroupConfig {
	// Map definition JSON object -> `definition` with metadata.
	definitions := []*definition{}
	linq.From(definitionsJSON).
		WhereT(func(kv linq.KeyValue) bool {
			// Skip these objects, special case. They're deprecated and empty.
			//
			// TODO(hausdorff): We can remove these now that we don't emit a `KindConfig` for an object
			// that has no properties.
			defName := kv.Key.(string)
			return !strings.HasPrefix(defName, "io.k8s.kubernetes.pkg")
		}).
		SelectT(func(kv linq.KeyValue) *definition {
			defName := kv.Key.(string)
			return &definition{
				gvk:  gvkFromRef(defName),
				name: defName,
				data: definitionsJSON[defName].(map[string]interface{}),
			}
		}).
		ToSlice(&definitions)

	//
	// Assemble a `KindConfig` for each Kubernetes kind.
	//

	kinds := []*KindConfig{}
	linq.From(definitions).
		OrderByT(func(d *definition) string { return d.gvk.String() }).
		SelectManyT(func(d *definition) linq.Query {
			// Skip if there are no properties on the type.
			if _, exists := d.data["properties"]; !exists {
				return linq.From([]KindConfig{})
			}

			// Make fully-qualified and "default" GroupVersion. The "default" GV is the `apiVersion` that
			// appears when writing Kubernetes YAML (e.g., `v1` instead of `core/v1`), while the
			// fully-qualified version is the "official" GV (e.g., `core/v1` instead of `v1` or
			// `admissionregistration.k8s.io/v1alpha1` instead of `admissionregistration/v1alpha1`).
			defaultGroupVersion := d.gvk.Group
			var fqGroupVersion string
			isTopLevel := isTopLevel(d)
			if gvks, gvkExists := d.data["x-kubernetes-group-version-kind"].([]interface{}); gvkExists && len(gvks) > 0 {
				gvk := gvks[0].(map[string]interface{})
				group := gvk["group"].(string)
				version := gvk["version"].(string)
				if group == "" {
					defaultGroupVersion = version
					fqGroupVersion = fmt.Sprintf(`core/%s`, version)
				} else {
					defaultGroupVersion = fmt.Sprintf(`%s/%s`, group, version)
					fqGroupVersion = fmt.Sprintf(`%s/%s`, group, version)
				}
			} else {
				gv := d.gvk.GroupVersion().String()
				if strings.HasPrefix(gv, "apiextensions/") && strings.HasPrefix(d.gvk.Kind, "CustomResource") {
					// Special case. Kubernetes OpenAPI spec should have an `x-kubernetes-group-version-kind`
					// CustomResource, but it doesn't. Hence, we hard-code it.
					gv = fmt.Sprintf("apiextensions.k8s.io/%s", d.gvk.Version)
				}
				defaultGroupVersion = gv
				fqGroupVersion = gv
			}

			ps := linq.From(d.data["properties"]).
				OrderByT(func(kv linq.KeyValue) string { return kv.Key.(string) }).
				WhereT(func(kv linq.KeyValue) bool {
					propName := kv.Key.(string)
					if (opts.language == python) && (propName == "apiVersion" || propName == "kind") {
						return false
					}
					return true
				}).
				SelectT(func(kv linq.KeyValue) *Property {
					propName := kv.Key.(string)
					prop := d.data["properties"].(map[string]interface{})[propName].(map[string]interface{})

					var prefix string
					var t, pyConstructorT string
					isListOrMap := false
					switch opts.language {
					case typescript:
						prefix = "      "
						t = makeType(d.name, propName, prop, opts)
					case python:
						prefix = "    "
						t = makeType(d.name, propName, prop, opts)
						pyConstructorT = makeType(d.name, propName, prop,
							groupOpts{language: python, generatorType: inputsAPI})
					case dotnet:
						prefix = "        "
						t = makeType(d.name, propName, prop, opts)
						if strings.HasPrefix(t, "InputList") || strings.HasPrefix(t, "InputMap") {
							isListOrMap = true
						}
					default:
						panic(fmt.Sprintf("Unsupported language '%s'", opts.language))
					}

					// `-` is invalid in TS variable names, so replace with `_`
					propName = strings.ReplaceAll(propName, "-", "_")

					// Create a default value for the field.
					defaultValue := fmt.Sprintf("args && args.%s || undefined", propName)
					switch propName {
					case "apiVersion":
						defaultValue = fmt.Sprintf(`"%s"`, defaultGroupVersion)
						if opts.language == typescript && isTopLevel {
							switch opts.generatorType {
							case provider:
								t = fmt.Sprintf(`pulumi.Output<"%s">`, defaultGroupVersion)
							case outputsAPI:
								t = fmt.Sprintf(`"%s"`, defaultGroupVersion)
							case inputsAPI:
								t = fmt.Sprintf(`pulumi.Input<"%s">`, defaultGroupVersion)
							}
						}
					case "kind":
						defaultValue = fmt.Sprintf(`"%s"`, d.gvk.Kind)
						if opts.language == typescript && isTopLevel {
							switch opts.generatorType {
							case provider:
								t = fmt.Sprintf(`pulumi.Output<"%s">`, d.gvk.Kind)
							case outputsAPI:
								t = fmt.Sprintf(`"%s"`, d.gvk.Kind)
							case inputsAPI:
								t = fmt.Sprintf(`pulumi.Input<"%s">`, d.gvk.Kind)
							}
						}
					}

					var languageName string
					var dotnetVarName string
					switch opts.language {
					case typescript:
						languageName = propName
					case python:
						languageName = pycodegen.PyName(propName)
					case dotnet:
						name := propName
						if name[0] == '$' {
							// $ref and $schema are property names that we want to special case into
							// just Ref and Schema.
							name = name[1:]
						}
						languageName = strings.ToUpper(name[:1]) + name[1:]
						if languageName == d.gvk.Kind {
							// .NET does not allow properties to be the same as the enclosing class - so special case these
							languageName = languageName + "Value"
						}
						dotnetVarName = "_" + name
					default:
						panic(fmt.Sprintf("Unsupported language '%s'", opts.language))
					}

					return &Property{
						comment:                   fmtComment(prop["description"], prefix, false, opts, d.gvk),
						pythonConstructorComment:  fmtComment(prop["description"], prefix+prefix+"       ", true, opts, d.gvk),
						propType:                  t,
						pythonConstructorPropType: pyConstructorT,
						name:                      propName,
						languageName:              languageName,
						dotnetVarName:             dotnetVarName,
						defaultValue:              defaultValue,
						isLast:                    false,
						dotnetIsListOrMap:         isListOrMap,
					}
				})

			// All properties.
			properties := []*Property{}
			ps.ToSlice(&properties)
			if len(properties) > 0 {
				properties[len(properties)-1].isLast = true
			}

			// Required properties.
			reqdProps := sets.NewString()
			if reqd, hasReqd := d.data["required"]; hasReqd {
				for _, propName := range reqd.([]interface{}) {
					reqdProps.Insert(propName.(string))
				}
			}

			requiredInputProperties := []*Property{}
			ps.
				WhereT(func(p *Property) bool {
					return reqdProps.Has(p.name)
				}).
				ToSlice(&requiredInputProperties)

			optionalInputProperties := []*Property{}
			ps.
				WhereT(func(p *Property) bool {
					return !reqdProps.Has(p.name) && p.name != "status"
				}).
				ToSlice(&optionalInputProperties)

			if len(properties) == 0 {
				return linq.From([]*KindConfig{})
			}

			if opts.generatorType == provider && (!isTopLevel) {
				return linq.From([]*KindConfig{})
			}

			var typeGuard string
			props := d.data["properties"].(map[string]interface{})
			_, apiVersionExists := props["apiVersion"]
			if apiVersionExists {
				typeGuard = fmt.Sprintf(`
    export function is%s(o: any): o is %s {
      return o.apiVersion == "%s" && o.kind == "%s";
    }`, d.gvk.Kind, d.gvk.Kind, defaultGroupVersion, d.gvk.Kind)
			}

			return linq.From([]*KindConfig{
				{
					kind: d.gvk.Kind,
					// NOTE: This transformation assumes git users on Windows to set
					// the "check in with UNIX line endings" setting.
					comment:                 fmtComment(d.data["description"], "    ", true, opts, d.gvk),
					pulumiComment:           fmtComment(PulumiComment(d.gvk.Kind), "    ", true, opts, d.gvk),
					properties:              properties,
					requiredInputProperties: requiredInputProperties,
					optionalInputProperties: optionalInputProperties,
					additionalSecretOutputs: additionalSecretOutputs(d.gvk),
					aliases:                 aliasesForGVK(d.gvk),
					gvk:                     &d.gvk,
					apiVersion:              fqGroupVersion,
					rawAPIVersion:           defaultGroupVersion,
					typeGuard:               typeGuard,
				},
			})
		}).
		ToSlice(&kinds)

	//
	// Assemble a `VersionConfig` for each group of kinds.
	//

	versions := []*VersionConfig{}
	linq.From(kinds).
		GroupByT(
			func(e *KindConfig) schema.GroupVersion { return e.gvk.GroupVersion() },
			func(e *KindConfig) *KindConfig { return e }).
		OrderByT(func(kinds linq.Group) string {
			return kinds.Key.(schema.GroupVersion).String()
		}).
		SelectManyT(func(kinds linq.Group) linq.Query {
			gv := kinds.Key.(schema.GroupVersion)
			kindsGroup := []*KindConfig{}
			linq.From(kinds.Group).ToSlice(&kindsGroup)
			if len(kindsGroup) == 0 {
				return linq.From([]*VersionConfig{})
			}

			version := gv.Version
			if opts.language == dotnet {
				version = pascalCase(version)
			}

			return linq.From([]*VersionConfig{
				{
					version:       version,
					kinds:         kindsGroup,
					gv:            &gv,
					apiVersion:    kindsGroup[0].apiVersion,    // NOTE: This is safe.
					rawAPIVersion: kindsGroup[0].rawAPIVersion, // NOTE: This is safe.
				},
			})
		}).
		ToSlice(&versions)

	//
	// Assemble a `GroupConfig` for each group of versions.
	//

	groups := []*GroupConfig{}
	linq.From(versions).
		GroupByT(
			func(e *VersionConfig) string { return e.gv.Group },
			func(e *VersionConfig) *VersionConfig { return e }).
		OrderByT(func(versions linq.Group) string { return versions.Key.(string) }).
		SelectManyT(func(versions linq.Group) linq.Query {
			versionsGroup := []*VersionConfig{}
			linq.From(versions.Group).ToSlice(&versionsGroup)
			if len(versionsGroup) == 0 {
				return linq.From([]*GroupConfig{})
			}

			group := versions.Key.(string)
			if opts.language == dotnet {
				group = pascalCase(group)
			}

			return linq.From([]*GroupConfig{
				{
					group:    group,
					versions: versionsGroup,
				},
			})
		}).
		WhereT(func(gc *GroupConfig) bool {
			return len(gc.Versions()) != 0
		}).
		ToSlice(&groups)

	return groups
}

func additionalSecretOutputs(gvk schema.GroupVersionKind) []string {
	kind := kinds.Kind(gvk.Kind)

	switch kind {
	case kinds.Secret:
		return []string{"data", "stringData"}
	default:
		return []string{}
	}
}

func aliasesForGVK(gvk schema.GroupVersionKind) []string {
	kind := kinds.Kind(gvk.Kind)

	switch kind {
	case kinds.DaemonSet:
		return []string{
			"kubernetes:apps/v1:DaemonSet",
			"kubernetes:apps/v1beta2:DaemonSet",
			"kubernetes:extensions/v1beta1:DaemonSet",
		}
	case kinds.Deployment:
		return []string{
			"kubernetes:apps/v1:Deployment",
			"kubernetes:apps/v1beta1:Deployment",
			"kubernetes:apps/v1beta2:Deployment",
			"kubernetes:extensions/v1beta1:Deployment",
		}
	case kinds.Ingress:
		return []string{
			"kubernetes:networking/v1beta1:Ingress",
			"kubernetes:extensions/v1beta1:Ingress",
		}
	case kinds.ReplicaSet:
		return []string{
			"kubernetes:apps/v1:ReplicaSet",
			"kubernetes:apps/v1beta2:ReplicaSet",
			"kubernetes:extensions/v1beta1:ReplicaSet",
		}
	case kinds.StatefulSet:
		return []string{
			"kubernetes:apps/v1:StatefulSet",
			"kubernetes:apps/v1beta1:StatefulSet",
			"kubernetes:apps/v1beta2:StatefulSet",
			"kubernetes:extensions/v1beta1:StatefulSet",
		}
	default:
		return []string{}
	}
}
