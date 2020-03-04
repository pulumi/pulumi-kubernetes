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
	versions []VersionConfig

	hasTopLevelKinds bool
}

// Group returns the name of the group (e.g., `core` for core, etc.)
func (gc GroupConfig) Group() string { return gc.group }

// Versions returns the set of version for some Kubernetes API group. For example, the `apps` group
// has `v1beta1`, `v1beta2`, and `v1`.
func (gc GroupConfig) Versions() []VersionConfig { return gc.versions }

// HasTopLevelKinds returns true if this group has top-level kinds.
func (gc GroupConfig) HasTopLevelKinds() bool { return gc.hasTopLevelKinds }

// VersionConfig represents a version of a Kubernetes API group (e.g., the `apps` group has
// `v1beta1`, `v1beta2`, and `v1`.)
type VersionConfig struct {
	version string
	kinds   []KindConfig

	gv            schema.GroupVersion // Used for sorting.
	apiVersion    string
	rawAPIVersion string

	hasTopLevelKinds bool
}

// Version returns the name of the version (e.g., `apps/v1beta1` would return `v1beta1`).
func (vc VersionConfig) Version() string { return vc.version }

// Kinds returns the set of kinds in some Kubernetes API group/version combination (e.g.,
// `apps/v1beta1` has the `Deployment` kind, etc.).
func (vc VersionConfig) Kinds() []KindConfig { return vc.kinds }

// HasTopLevelKinds returns true if this group has top-level kinds.
func (vc VersionConfig) HasTopLevelKinds() bool { return vc.hasTopLevelKinds }

// TopLevelKinds returns the set of kinds that are not nested.
func (vc VersionConfig) TopLevelKinds() []KindConfig {
	var kinds []KindConfig
	for _, k := range vc.kinds {
		if !k.IsNested() {
			kinds = append(kinds, k)
		}
	}
	return kinds
}

// TopLevelKindsAndAliases will produce a list of kinds, including aliases (e.g., both `apiregistration` and
// `apiregistration.k8s.io`).
func (vc VersionConfig) TopLevelKindsAndAliases() []KindConfig {
	var kindsAndAliases []KindConfig
	for _, kind := range vc.TopLevelKinds() {
		kindsAndAliases = append(kindsAndAliases, kind)
		if strings.HasPrefix(kind.APIVersion(), apiRegistration) {
			alias := KindConfig{}
			err := copier.Copy(&alias, kind)
			if err != nil {
				panic(err)
			}
			rawAPIVersion := "apiregistration" + strings.TrimPrefix(kind.APIVersion(), apiRegistration)
			alias.rawAPIVersion = rawAPIVersion
			kindsAndAliases = append(kindsAndAliases, alias)
		}
	}
	return kindsAndAliases
}

// ListTopLevelKindsAndAliases will return all known `Kind`s that are lists, or aliases of lists. These
// `Kind`s are not instantiated by the API server, and we must "flatten" them client-side to get an
// accurate view of what resource operations we need to perform.
func (vc VersionConfig) ListTopLevelKindsAndAliases() []KindConfig {
	var listKinds []KindConfig
	for _, kind := range vc.TopLevelKindsAndAliases() {
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
func (vc VersionConfig) APIVersion() string { return vc.apiVersion }

// RawAPIVersion returns the "raw" apiVersion (e.g., `v1` rather than `core/v1`).
func (vc VersionConfig) RawAPIVersion() string { return vc.rawAPIVersion }

// KindConfig represents a Kubernetes API kind (e.g., the `Deployment` type in
// `apps/v1beta1/Deployment`).
type KindConfig struct {
	kind                    string
	deprecationComment      string
	comment                 string
	pulumiComment           string
	properties              []Property
	requiredInputProperties []Property
	optionalInputProperties []Property
	additionalSecretOutputs []string
	aliases                 []string

	gvk           schema.GroupVersionKind // Used for sorting.
	apiVersion    string
	rawAPIVersion string
	typeGuard     string

	isNested bool
}

// Kind returns the name of the Kubernetes API kind (e.g., `Deployment` for
// `apps/v1beta1/Deployment`).
func (kc KindConfig) Kind() string { return kc.kind }

// DeprecationComment returns the deprecation comment for deprecated APIs, otherwise an empty string.
func (kc KindConfig) DeprecationComment() string { return kc.deprecationComment }

// Comment returns the comments associated with some Kubernetes API kind.
func (kc KindConfig) Comment() string { return kc.comment }

// PulumiComment returns the await logic documentation associated with some Kubernetes API kind.
func (kc KindConfig) PulumiComment() string { return kc.pulumiComment }

// Properties returns the list of properties that exist on some Kubernetes API kind (i.e., things
// that we will want to `.` into, like `thing.apiVersion`, `thing.kind`, `thing.metadata`, etc.).
func (kc KindConfig) Properties() []Property { return kc.properties }

// RequiredInputProperties returns the list of properties that are required input properties on some
// Kubernetes API kind (i.e., things that we will want to provide, like `thing.metadata`, etc.).
func (kc KindConfig) RequiredInputProperties() []Property { return kc.requiredInputProperties }

// OptionalInputProperties returns the list of properties that are optional input properties on some
// Kubernetes API kind (i.e., things that we will want to provide, like `thing.metadata`, etc.).
func (kc KindConfig) OptionalInputProperties() []Property { return kc.optionalInputProperties }

// AdditionalSecretOutputs returns the list of strings to set as additionalSecretOutputs on some
// Kubernetes API kind.
func (kc KindConfig) AdditionalSecretOutputs() []string { return kc.additionalSecretOutputs }

// Aliases returns the list of aliases for a Kubernetes API kind.
func (kc KindConfig) Aliases() []string { return kc.aliases }

// APIVersion returns the fully-qualified apiVersion (e.g., `storage.k8s.io/v1` for storage, etc.)
func (kc KindConfig) APIVersion() string { return kc.apiVersion }

// RawAPIVersion returns the "raw" apiVersion (e.g., `v1` rather than `core/v1`).
func (kc KindConfig) RawAPIVersion() string { return kc.rawAPIVersion }

// URNAPIVersion returns API version that can be used in a URN (e.g., using the backwards-compatible
// alias `apiextensions` instead of `apiextensions.k8s.io`).
func (kc KindConfig) URNAPIVersion() string {
	if strings.HasPrefix(kc.apiVersion, apiRegistration) {
		return "apiregistration" + strings.TrimPrefix(kc.apiVersion, apiRegistration)
	}
	return kc.apiVersion
}

// TypeGuard returns the text of a TypeScript type guard for the given kind.
func (kc KindConfig) TypeGuard() string { return kc.typeGuard }

// IsNested returns true if this is a nested kind.
func (kc KindConfig) IsNested() bool { return kc.isNested }

// Property represents a property we want to expose on a Kubernetes API kind (i.e., things that we
// will want to `.` into, like `thing.apiVersion`, `thing.kind`, `thing.metadata`, etc.).
type Property struct {
	name                     string
	languageName             string
	comment                  string
	pythonConstructorComment string
	inputsAPIType            string
	outputsAPIType           string
	providerType             string
	defaultValue             string
	isLast                   bool
	dotnetVarName            string
	dotnetIsListOrMap        bool
}

// Name returns the name of the property.
func (p Property) Name() string { return p.name }

// LanguageName returns the name of the property.
func (p Property) LanguageName() string { return p.languageName }

// Comment returns the comments associated with some property.
func (p Property) Comment() string { return p.comment }

// PythonConstructorComment returns the comments associated with some property, formatted for Python
// constructor documentation.
func (p Property) PythonConstructorComment() string { return p.pythonConstructorComment }

// InputsAPIType returns the type of the property for the inputs API.
func (p Property) InputsAPIType() string { return p.inputsAPIType }

// OutputsAPIType returns the type of the property for the outputs API.
func (p Property) OutputsAPIType() string { return p.outputsAPIType }

// ProviderType returns the type of the property for the provider API.
func (p Property) ProviderType() string { return p.providerType }

// DefaultValue returns the type of the property.
func (p Property) DefaultValue() string { return p.defaultValue }

// IsLast returns whether the property is the last in the list of properties.
func (p Property) IsLast() bool { return p.isLast }

// DotnetVarName returns a variable name safe to use in .NET (e.g. `@namespace` instead of `namespace`)
func (p Property) DotnetVarName() string { return p.dotnetVarName }

// DotnetIsListOrMap returns whether the property type is a List or map
func (p Property) DotnetIsListOrMap() bool { return p.dotnetIsListOrMap }

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

// extractDeprecationComment returns the comment with deprecation comment removed and the extracted deprecation
// comment, fixed-up for the specified language.
func extractDeprecationComment(comment interface{}, gvk schema.GroupVersionKind, language language) (string, string) {
	if comment == nil {
		return "", ""
	}

	commentstr, _ := comment.(string)
	if commentstr == "" {
		return "", ""
	}

	re := regexp.MustCompile(`DEPRECATED - .* is deprecated by .* for more information\.\s*`)

	var prefix, suffix string
	switch language {
	case typescript:
		prefix = "\n\n@deprecated "
		suffix = ""
	case python, dotnet:
		prefix = "DEPRECATED - "
		suffix = "\n\n"
	default:
		panic(fmt.Sprintf("Unsupported language '%s'", language))
	}

	if re.MatchString(commentstr) {
		deprecationMessage := prefix + ApiVersionComment(gvk) + suffix
		return re.ReplaceAllString(commentstr, ""), deprecationMessage
	}

	return commentstr, ""
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
	apiextensionsV1beta1                = "io.k8s.apiextensions-apiserver.pkg.apis.apiextensions.v1beta1"
	apiextensionsV1                     = "io.k8s.apiextensions-apiserver.pkg.apis.apiextensions.v1"
	quantity                            = "io.k8s.apimachinery.pkg.api.resource.Quantity"
	rawExtension                        = "io.k8s.apimachinery.pkg.runtime.RawExtension"
	intOrString                         = "io.k8s.apimachinery.pkg.util.intstr.IntOrString"
	v1Fields                            = "io.k8s.apimachinery.pkg.apis.meta.v1.Fields"
	v1FieldsV1                          = "io.k8s.apimachinery.pkg.apis.meta.v1.FieldsV1"
	v1Time                              = "io.k8s.apimachinery.pkg.apis.meta.v1.Time"
	v1MicroTime                         = "io.k8s.apimachinery.pkg.apis.meta.v1.MicroTime"
	v1beta1JSONSchemaPropsOrBool        = apiextensionsV1beta1 + ".JSONSchemaPropsOrBool"
	v1beta1JSONSchemaPropsOrArray       = apiextensionsV1beta1 + ".JSONSchemaPropsOrArray"
	v1beta1JSONSchemaPropsOrStringArray = apiextensionsV1beta1 + ".JSONSchemaPropsOrStringArray"
	v1beta1JSON                         = apiextensionsV1beta1 + ".JSON"
	v1beta1CRSubresourceStatus          = apiextensionsV1beta1 + ".CustomResourceSubresourceStatus"
	v1JSONSchemaPropsOrBool             = apiextensionsV1 + ".JSONSchemaPropsOrBool"
	v1JSONSchemaPropsOrArray            = apiextensionsV1 + ".JSONSchemaPropsOrArray"
	v1JSONSchemaPropsOrStringArray      = apiextensionsV1 + ".JSONSchemaPropsOrStringArray"
	v1JSON                              = apiextensionsV1 + ".JSON"
	v1CRSubresourceStatus               = apiextensionsV1 + ".CustomResourceSubresourceStatus"
)

func makeTypescriptType(resourceType, propName string, prop map[string]interface{}, gentype gentype) string {
	wrapType := func(typ string) string {
		switch gentype {
		case provider:
			return fmt.Sprintf("pulumi.Output<%s>", typ)
		case outputsAPI:
			return typ
		case inputsAPI:
			return fmt.Sprintf("pulumi.Input<%s>", typ)
		default:
			panic(fmt.Sprintf("unrecognized generator type %d", gentype))
		}
	}

	refPrefix := ""
	if gentype == provider {
		refPrefix = "outputs"
	}

	if t, exists := prop["type"]; exists {
		tstr := t.(string)
		if tstr == "array" {
			elemType := makeTypescriptType(
				resourceType, propName, prop["items"].(map[string]interface{}), gentype)
			switch gentype {
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
					switch gentype {
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

			switch gentype {
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

func makePythonType(resourceType, propName string, prop map[string]interface{}, gentype gentype) string {
	wrapType := func(typ string) string {
		switch gentype {
		case provider:
			return fmt.Sprintf("pulumi.Output[%s]", typ)
		case outputsAPI:
			return typ
		case inputsAPI:
			return fmt.Sprintf("pulumi.Input[%s]", typ)
		default:
			panic(fmt.Sprintf("unrecognized generator type %d", gentype))
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

func makeDotnetType(resourceType, propName string, prop map[string]interface{}, gentype gentype, forceNoWrap bool) string {

	refPrefix := ""
	if gentype == provider {
		refPrefix = "Types.Outputs"
	}

	wrapType := func(typ string) string {
		if forceNoWrap {
			return typ
		}
		switch gentype {
		case provider:
			return fmt.Sprintf("Output<%s>", typ)
		case outputsAPI:
			return typ
		case inputsAPI:
			return fmt.Sprintf("Input<%s>", typ)
		default:
			panic(fmt.Sprintf("unrecognized generator type %d", gentype))
		}
	}

	oneOf := func(typeA string, typeB string) string {
		if forceNoWrap {
			return fmt.Sprintf("Union<%s,%s>", typeA, typeB)
		}
		switch gentype {
		case provider:
			return fmt.Sprintf("Output<Union<%s,%s>>", typeA, typeB)
		case outputsAPI:
			return fmt.Sprintf("Union<%s,%s>", typeA, typeB)
		case inputsAPI:
			return fmt.Sprintf("InputUnion<%s,%s>", typeA, typeB)
		default:
			panic(fmt.Sprintf("unrecognized generator type %d", gentype))
		}
	}

	if t, exists := prop["type"]; exists {
		tstr := t.(string)
		if tstr == "array" {
			switch gentype {
			case provider:
				elemType := makeDotnetType(
					resourceType, propName, prop["items"].(map[string]interface{}), gentype, true)
				return fmt.Sprintf("Output<ImmutableArray<%s>>", elemType)
			case outputsAPI:
				elemType := makeDotnetType(
					resourceType, propName, prop["items"].(map[string]interface{}), gentype, forceNoWrap)
				return fmt.Sprintf("ImmutableArray<%s>", elemType)
			case inputsAPI:
				elemType := makeDotnetType(
					resourceType, propName, prop["items"].(map[string]interface{}), gentype, true)
				return fmt.Sprintf("InputList<%s>", elemType)
			}
		} else if tstr == "integer" {
			return wrapType("int")
		} else if tstr == "boolean" {
			return wrapType("bool")
		} else if tstr == "number" {
			return wrapType("double")
		} else if tstr == "object" {
			vtype := "string"
			if additionalProperties, exists := prop["additionalProperties"]; exists {
				switch gentype {
				case provider:
					vtype = makeDotnetType(
						resourceType, propName, additionalProperties.(map[string]interface{}), gentype, true)
				case outputsAPI:
					vtype = makeDotnetType(
						resourceType, propName, additionalProperties.(map[string]interface{}), gentype, forceNoWrap)
				case inputsAPI:
					vtype = makeDotnetType(
						resourceType, propName, additionalProperties.(map[string]interface{}), gentype, true)
				}
			}
			switch gentype {
			case inputsAPI:
				return fmt.Sprintf("InputMap<%s>", vtype)
			case outputsAPI:
				return fmt.Sprintf("ImmutableDictionary<string, %s>", vtype)
			case provider:
				return fmt.Sprintf("Output<ImmutableDictionary<string, %s>>", vtype)
			}
		} else if tstr == "string" && resourceType == "io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta" && propName == "namespace" {
			// Special case: `.metadata.namespace` should either take a string or a namespace object
			// itself.

			switch gentype {
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
	var argsSuffix string
	var stringArr string
	var jsonType string
	switch gentype {
	case inputsAPI:
		argsSuffix = "Args"
		stringArr = "InputList<string>"
		jsonType = "InputJson"
	case outputsAPI:
		argsSuffix = ""
		stringArr = "ImmutableArray<string>"
		jsonType = "System.Text.Json.JsonElement"
	case provider:
		argsSuffix = ""
		stringArr = "string[]"
		jsonType = "Output<System.Text.Json.JsonElement>"
	default:
		panic(fmt.Sprintf("unrecognized generator type %d", gentype))
	}

	isSimpleRef := true
	switch ref {
	case quantity:
		ref = "string"
	case intOrString:
		return oneOf("int", "string")
	case v1Fields, v1FieldsV1, rawExtension:
		return jsonType
	case v1Time, v1MicroTime:
		ref = "string"
	case v1beta1JSONSchemaPropsOrBool:
		return oneOf("ApiExtensions.V1Beta1.JSONSchemaProps"+argsSuffix, "bool")
	case v1JSONSchemaPropsOrBool:
		return oneOf("ApiExtensions.V1.JSONSchemaProps"+argsSuffix, "bool")
	case v1beta1JSONSchemaPropsOrArray, v1beta1JSONSchemaPropsOrStringArray:
		return oneOf("ApiExtensions.V1Beta1.JSONSchemaProps"+argsSuffix, stringArr)
	case v1JSONSchemaPropsOrArray, v1JSONSchemaPropsOrStringArray:
		return oneOf("ApiExtensions.V1.JSONSchemaProps"+argsSuffix, stringArr)
	case v1beta1JSON, v1beta1CRSubresourceStatus, v1JSON, v1CRSubresourceStatus:
		return jsonType
	default:
		isSimpleRef = false
	}

	if isSimpleRef {
		return wrapType(ref)
	}

	gvk := gvkFromRef(ref)
	group := pascalCase(gvk.Group)
	version := pascalCase(gvk.Version)
	kind := gvk.Kind
	if gentype == inputsAPI {
		kind = kind + "Args"
	}

	var gvkRefStr string
	if refPrefix == "" {
		gvkRefStr = fmt.Sprintf("%s.%s.%s", group, version, kind)
	} else {
		gvkRefStr = fmt.Sprintf("%s.%s.%s.%s", refPrefix, group, version, kind)
	}

	return wrapType(gvkRefStr)
}

func makeTypes(resourceType, propName string, prop map[string]interface{}, language language) (string, string, string) {
	inputsAPIType := makeType(resourceType, propName, prop, language, inputsAPI)
	outputsAPIType := makeType(resourceType, propName, prop, language, outputsAPI)
	providerType := makeType(resourceType, propName, prop, language, provider)
	return inputsAPIType, outputsAPIType, providerType
}

func makeType(resourceType, propName string, prop map[string]interface{}, language language, gentype gentype) string {
	switch language {
	case typescript:
		return makeTypescriptType(resourceType, propName, prop, gentype)
	case python:
		return makePythonType(resourceType, propName, prop, gentype)
	case dotnet:
		return makeDotnetType(resourceType, propName, prop, gentype, false)
	default:
		panic(fmt.Sprintf("Unsupported language '%s'", language))
	}
}

// --------------------------------------------------------------------------

// Core grouping logic.

// --------------------------------------------------------------------------

type definition struct {
	gvk  schema.GroupVersionKind
	name string
	data map[string]interface{}
}

// fqGroupVersion returns the fully-qualified GroupVersion, which is the "official" GV
// (e.g., `core/v1` instead of `v1` or `admissionregistration.k8s.io/v1alpha1` instead of
// `admissionregistration/v1alpha1`).
func (d definition) fqGroupVersion() string {
	if gvks, gvkExists := d.data["x-kubernetes-group-version-kind"].([]interface{}); gvkExists && len(gvks) > 0 {
		gvk := gvks[0].(map[string]interface{})
		group := gvk["group"].(string)
		version := gvk["version"].(string)
		if group == "" {
			return fmt.Sprintf(`core/%s`, version)
		} else {
			return fmt.Sprintf(`%s/%s`, group, version)
		}
	}

	return d.gvk.GroupVersion().String()
}

func (d definition) isTopLevel() bool {
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

type gentype int

const (
	provider gentype = iota
	inputsAPI
	outputsAPI
)

type language string

const (
	python     language = "python"
	typescript language = "typescript"
	dotnet     language = "dotnet"
)

type groupOpts struct {
	language language
}

func nodeJSOpts() groupOpts { return groupOpts{language: typescript} }
func pythonOpts() groupOpts { return groupOpts{language: python} }
func dotnetOpts() groupOpts { return groupOpts{language: dotnet} }

func allCamelCasePropertyNames(definitionsJSON map[string]interface{}, opts groupOpts) []string {
	// Map definition JSON object -> `definition` with metadata.
	var definitions []definition
	linq.From(definitionsJSON).
		WhereT(func(kv linq.KeyValue) bool {
			// Skip these objects, special case. They're deprecated and empty.
			defName := kv.Key.(string)
			return !strings.HasPrefix(defName, "io.k8s.kubernetes.pkg")
		}).
		SelectT(func(kv linq.KeyValue) definition {
			defName := kv.Key.(string)
			return definition{
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

func createGroups(definitionsJSON map[string]interface{}, opts groupOpts) []GroupConfig {
	// Map definition JSON object -> `definition` with metadata.
	var definitions []definition
	linq.From(definitionsJSON).
		SelectT(func(kv linq.KeyValue) definition {
			defName := kv.Key.(string)
			return definition{
				gvk:  gvkFromRef(defName),
				name: defName,
				data: definitionsJSON[defName].(map[string]interface{}),
			}
		}).
		ToSlice(&definitions)

	// Compute aliases for Kinds. Many k8s resources have multiple GVs, so create a map from Kind -> GV string.
	// For Kinds with more than one GV, create aliases in the SDKs.
	aliases := map[string][]interface{}{}
	linq.From(definitions).
		WhereT(func(d definition) bool { return d.isTopLevel() && !strings.HasSuffix(d.gvk.Kind, "List") }).
		OrderByT(func(d definition) string { return d.gvk.String() }).
		SelectManyT(func(d definition) linq.Query {
			return linq.From([]KindConfig{
				{
					kind:       d.gvk.Kind,
					apiVersion: d.fqGroupVersion(),
				},
			})
		}).
		GroupByT(
			func(kind KindConfig) string {
				return kind.kind
			},
			func(kind KindConfig) string {
				return fmt.Sprintf("kubernetes:%s:%s", kind.apiVersion, kind.kind)
			}).
		WhereT(func(group linq.Group) bool {
			return len(group.Group) > 1
		}).
		ToMapBy(&aliases,
			func(i interface{}) interface{} {
				return i.(linq.Group).Key
			},
			func(i interface{}) interface{} {
				return i.(linq.Group).Group
			})
	aliasesForKind := func(kind, fqGroupVersion string) []string {
		var results []string

		for _, alias := range aliases[kind] {
			aliasString := alias.(string)
			re := fmt.Sprintf(`:%s:`, fqGroupVersion)
			match, err := regexp.MatchString(re, aliasString)
			if err == nil && match {
				continue
			}
			results = append(results, aliasString)
		}
		return results
	}

	//
	// Assemble a `KindConfig` for each Kubernetes kind.
	//

	var kinds []KindConfig
	linq.From(definitions).
		OrderByT(func(d definition) string { return d.gvk.String() }).
		SelectManyT(func(d definition) linq.Query {
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
			isTopLevel := d.isTopLevel()
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
				SelectT(func(kv linq.KeyValue) Property {
					propName := kv.Key.(string)
					prop := d.data["properties"].(map[string]interface{})[propName].(map[string]interface{})

					var prefix string
					var inputsAPIType, outputsAPIType, providerType string
					isListOrMap := false
					switch opts.language {
					case typescript:
						prefix = "      "
						inputsAPIType, outputsAPIType, providerType = makeTypes(d.name, propName, prop, typescript)
					case python:
						prefix = "    "
						inputsAPIType, outputsAPIType, providerType = makeTypes(d.name, propName, prop, python)
					case dotnet:
						prefix = "        "
						inputsAPIType, outputsAPIType, providerType = makeTypes(d.name, propName, prop, dotnet)
						if strings.HasPrefix(inputsAPIType, "InputList") || strings.HasPrefix(inputsAPIType, "InputMap") {
							isListOrMap = true
						}
					default:
						panic(fmt.Sprintf("Unsupported language '%s'", opts.language))
					}

					// `-` is invalid in TS variable names, so replace with `_`
					propName = strings.ReplaceAll(propName, "-", "_")

					// Create a default value for the field.
					defaultValue := fmt.Sprintf("args?.%s", propName)
					switch propName {
					case "apiVersion":
						defaultValue = fmt.Sprintf(`"%s"`, defaultGroupVersion)
						if opts.language == typescript && isTopLevel {
							inputsAPIType = fmt.Sprintf(`pulumi.Input<"%s">`, defaultGroupVersion)
							outputsAPIType = fmt.Sprintf(`"%s"`, defaultGroupVersion)
							providerType = fmt.Sprintf(`pulumi.Output<"%s">`, defaultGroupVersion)
						}
					case "kind":
						defaultValue = fmt.Sprintf(`"%s"`, d.gvk.Kind)
						if opts.language == typescript && isTopLevel {
							inputsAPIType = fmt.Sprintf(`pulumi.Input<"%s">`, d.gvk.Kind)
							outputsAPIType = fmt.Sprintf(`"%s"`, d.gvk.Kind)
							providerType = fmt.Sprintf(`pulumi.Output<"%s">`, d.gvk.Kind)
						}
					}

					var languageName string
					var dotnetVarName string
					switch opts.language {
					case typescript:
						languageName = propName
					case python:
						languageName = pycodegen.PyName(propName)
						defaultValue = languageName
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
						defaultValue = ""
						dotnetVarName = "@" + name
					default:
						panic(fmt.Sprintf("Unsupported language '%s'", opts.language))
					}

					// Set Secret input fields to pulumi.secret
					if d.gvk.Kind == "Secret" && (propName == "stringData" || propName == "data") {
						switch opts.language {
						case typescript:
							defaultValue = fmt.Sprintf(
								"args?.%s === undefined ? undefined : pulumi.secret(args?.%s)", propName, propName)
						case python:
							defaultValue = fmt.Sprintf("pulumi.Output.secret(%s) if %s is not None else None",
								languageName, languageName)
						case dotnet:
							defaultValue = ".Apply(Output.CreateSecret)"
						}
					}

					return Property{
						comment:                  fmtComment(prop["description"], prefix, false, opts, d.gvk),
						pythonConstructorComment: fmtComment(prop["description"], prefix+prefix+"       ", true, opts, d.gvk),
						inputsAPIType:            inputsAPIType,
						outputsAPIType:           outputsAPIType,
						providerType:             providerType,
						name:                     propName,
						languageName:             languageName,
						dotnetVarName:            dotnetVarName,
						defaultValue:             defaultValue,
						isLast:                   false,
						dotnetIsListOrMap:        isListOrMap,
					}
				})

			// All properties.
			var properties []Property
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

			var requiredInputProperties []Property
			ps.
				WhereT(func(p Property) bool {
					return reqdProps.Has(p.name)
				}).
				ToSlice(&requiredInputProperties)

			var optionalInputProperties []Property
			ps.
				WhereT(func(p Property) bool {
					return !reqdProps.Has(p.name) && p.name != "status"
				}).
				ToSlice(&optionalInputProperties)

			if len(properties) == 0 {
				return linq.From([]KindConfig{})
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

			comment, deprecationComment := extractDeprecationComment(d.data["description"], d.gvk, opts.language)

			return linq.From([]KindConfig{
				{
					kind: d.gvk.Kind,
					// NOTE: This transformation assumes git users on Windows to set
					// the "check in with UNIX line endings" setting.
					deprecationComment:      fmtComment(deprecationComment, "    ", true, opts, d.gvk),
					comment:                 fmtComment(comment, "    ", true, opts, d.gvk),
					pulumiComment:           fmtComment(PulumiComment(d.gvk.Kind), "    ", true, opts, d.gvk),
					properties:              properties,
					requiredInputProperties: requiredInputProperties,
					optionalInputProperties: optionalInputProperties,
					additionalSecretOutputs: additionalSecretOutputs(d.gvk),
					aliases:                 aliasesForKind(d.gvk.Kind, fqGroupVersion),
					gvk:                     d.gvk,
					apiVersion:              fqGroupVersion,
					rawAPIVersion:           defaultGroupVersion,
					typeGuard:               typeGuard,
					isNested:                !isTopLevel,
				},
			})
		}).
		ToSlice(&kinds)

	//
	// Assemble a `VersionConfig` for each group of kinds.
	//

	var versions []VersionConfig
	linq.From(kinds).
		GroupByT(
			func(e KindConfig) schema.GroupVersion { return e.gvk.GroupVersion() },
			func(e KindConfig) KindConfig { return e }).
		OrderByT(func(kinds linq.Group) string {
			return kinds.Key.(schema.GroupVersion).String()
		}).
		SelectManyT(func(kinds linq.Group) linq.Query {
			gv := kinds.Key.(schema.GroupVersion)
			var kindsGroup []KindConfig
			linq.From(kinds.Group).ToSlice(&kindsGroup)
			if len(kindsGroup) == 0 {
				return linq.From([]VersionConfig{})
			}

			version := gv.Version
			if opts.language == dotnet {
				version = pascalCase(version)
			}

			hasTopLevelKinds := linq.From(kindsGroup).WhereT(func(k KindConfig) bool {
				return !k.IsNested()
			}).Any()

			return linq.From([]VersionConfig{
				{
					version:          version,
					kinds:            kindsGroup,
					gv:               gv,
					apiVersion:       kindsGroup[0].apiVersion,    // NOTE: This is safe.
					rawAPIVersion:    kindsGroup[0].rawAPIVersion, // NOTE: This is safe.
					hasTopLevelKinds: hasTopLevelKinds,
				},
			})
		}).
		ToSlice(&versions)

	//
	// Assemble a `GroupConfig` for each group of versions.
	//

	var groups []GroupConfig
	linq.From(versions).
		GroupByT(
			func(e VersionConfig) string { return e.gv.Group },
			func(e VersionConfig) VersionConfig { return e }).
		OrderByT(func(versions linq.Group) string { return versions.Key.(string) }).
		SelectManyT(func(versions linq.Group) linq.Query {
			var versionsGroup []VersionConfig
			linq.From(versions.Group).ToSlice(&versionsGroup)
			if len(versionsGroup) == 0 {
				return linq.From([]GroupConfig{})
			}

			group := versions.Key.(string)
			if opts.language == dotnet {
				group = pascalCase(group)
			}

			hasTopLevelKinds := linq.From(versionsGroup).WhereT(func(v VersionConfig) bool {
				return v.HasTopLevelKinds()
			}).Any()

			return linq.From([]GroupConfig{
				{
					group:            group,
					versions:         versionsGroup,
					hasTopLevelKinds: hasTopLevelKinds,
				},
			})
		}).
		WhereT(func(gc GroupConfig) bool {
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
