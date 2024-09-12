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

// nolint: lll
package gen

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/ahmetb/go-linq"
	pschema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
)

var validCharRegex = regexp.MustCompile(`[^a-zA-Z0-9]`)

// --------------------------------------------------------------------------

// A collection of data structures and utility functions to transform an OpenAPI spec for the
// Kubernetes API into something that we can use for codegen'ing nodejs and Python clients.

// --------------------------------------------------------------------------

// GroupConfig represents a Kubernetes API group (e.g., core, apps, extensions, etc.)
type GroupConfig struct {
	group    string
	versions []VersionConfig
}

// Group returns the name of the group (e.g., `core` for core, etc.)
func (gc GroupConfig) Group() string { return gc.group }

// Versions returns the set of version for some Kubernetes API group. For example, the `apps` group
// has `v1beta1`, `v1beta2`, and `v1`.
func (gc GroupConfig) Versions() []VersionConfig { return gc.versions }

// VersionConfig represents a version of a Kubernetes API group (e.g., the `apps` group has
// `v1beta1`, `v1beta2`, and `v1`.)
type VersionConfig struct {
	version string
	kinds   []KindConfig

	gv                schema.GroupVersion // Used for sorting.
	apiVersion        string
	defaultAPIVersion string
}

// Version returns the name of the version (e.g., `apps/v1beta1` would return `v1beta1`).
func (vc VersionConfig) Version() string { return vc.version }

// Kinds returns the set of kinds in some Kubernetes API group/version combination (e.g.,
// `apps/v1beta1` has the `Deployment` kind, etc.).
func (vc VersionConfig) Kinds() []KindConfig { return vc.kinds }

// APIVersion returns the fully-qualified apiVersion (e.g., `storage.k8s.io/v1` for storage, etc.)
func (vc VersionConfig) APIVersion() string { return vc.apiVersion }

// DefaultAPIVersion returns the default apiVersion (e.g., `v1` rather than `core/v1`).
func (vc VersionConfig) DefaultAPIVersion() string { return vc.defaultAPIVersion }

// KindConfig represents a Kubernetes API kind (e.g., the `Deployment` type in `apps/v1beta1/Deployment`).
type KindConfig struct {
	kind                    string
	deprecationComment      string
	comment                 string
	pulumiComment           string
	properties              []Property
	requiredInputProperties []Property
	optionalInputProperties []Property
	aliases                 []string

	gvk               schema.GroupVersionKind // Used for sorting.
	apiVersion        string
	defaultAPIVersion string

	isNested bool

	schemaPkgName string
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

// Aliases returns the list of aliases for a Kubernetes API kind.
func (kc KindConfig) Aliases() []string { return kc.aliases }

// IsNested returns true if this is a nested kind.
func (kc KindConfig) IsNested() bool { return kc.isNested }

// Property represents a property we want to expose on a Kubernetes API kind (i.e., things that we
// will want to `.` into, like `thing.apiVersion`, `thing.kind`, `thing.metadata`, etc.).
type Property struct {
	name         string
	comment      string
	schemaType   string
	constValue   string
	defaultValue string
}

// Name returns the name of the property.
func (p Property) Name() string { return p.name }

// Comment returns the comments associated with some property.
func (p Property) Comment() string { return p.comment }

// SchemaType returns the type of the property for the schema.
func (p Property) SchemaType() string { return p.schemaType }

// DefaultValue returns the constant value of the property.
func (p Property) ConstValue() string { return p.constValue }

// DefaultValue returns the default value of the property.
func (p Property) DefaultValue() string { return p.defaultValue }

type definition struct {
	gvk            schema.GroupVersionKind
	name           string
	data           map[string]any
	canonicalGroup string
}

// apiVersion creates a GV string in the canonical format.
func (d definition) apiVersion(canonicalGroups map[string]string) string {
	gvFmt := `%s/%s`

	// If the canonical group is set for this definition (i.e., it is a top-level resource), use that.
	if d.canonicalGroup != "" {
		return fmt.Sprintf(gvFmt, d.canonicalGroup, d.gvk.Version)
	}

	// Otherwise, look up the canonical group and use it.
	canonicalGroup := canonicalGroups[d.gvk.Group]
	return fmt.Sprintf(gvFmt, canonicalGroup, d.gvk.Version)
}

// defaultAPIVersion returns the "default" apiVersion that appears when writing Kubernetes
// YAML (e.g., `v1` instead of `core/v1`).
func (d definition) defaultAPIVersion() string {
	// Pull the canonical GVK from the OpenAPI `x-kubernetes-group-version-kind` field if it exists.
	if gvks, gvkExists := d.data["x-kubernetes-group-version-kind"].([]any); gvkExists && len(gvks) > 0 {
		gvk := gvks[0].(map[string]any)
		group := gvk["group"].(string)
		version := gvk["version"].(string)

		// Special case for the "core" group, which was historically called "".
		if group == "" {
			return version
		}

		return fmt.Sprintf(`%s/%s`, group, version)
	}

	// Fall back to using a GVK derived from the definition name.
	return d.gvk.GroupVersion().String()
}

func (d definition) isTopLevel() bool {
	gvks, gvkExists := d.data["x-kubernetes-group-version-kind"].([]any)
	hasGVK := gvkExists && len(gvks) > 0
	if !hasGVK {
		return false
	}

	// Return `false` for the handful of top-level imperative resource types that can't be managed
	// by Pulumi.
	switch fmt.Sprintf("%s/%s", d.gvk.GroupVersion().String(), d.gvk.Kind) {
	case
		"v1/Status",
		"apps/v1beta1/Scale",
		"apps/v1beta2/Scale",
		"authentication/v1/TokenRequest",
		"authentication/v1/TokenReview",
		"authentication/v1alpha1/SelfSubjectReview",
		"authentication/v1beta1/SelfSubjectReview",
		"authentication/v1/SelfSubjectReview",
		"authentication/v1beta1/TokenReview",
		"authorization/v1/LocalSubjectAccessReview",
		"authorization/v1/SelfSubjectAccessReview",
		"authorization/v1/SelfSubjectRulesReview",
		"authorization/v1/SubjectAccessReview",
		"authorization/v1beta1/LocalSubjectAccessReview",
		"authorization/v1beta1/SelfSubjectAccessReview",
		"authorization/v1beta1/SelfSubjectRulesReview",
		"authorization/v1beta1/SubjectAccessReview",
		"autoscaling/v1/Scale",
		"core/v1/ComponentStatus",
		"core/v1/ComponentStatusList",
		"extensions/v1beta1/Scale",
		"policy/v1beta1/Eviction",
		"policy/v1/Eviction":
		return false
	}

	properties, hasProperties := d.data["properties"].(map[string]any)
	if !hasProperties {
		return false
	}

	meta, hasMetadata := properties["metadata"].(map[string]any)
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

// Utility functions.

// --------------------------------------------------------------------------

func GVKFromRef(ref string) schema.GroupVersionKind {
	split := strings.Split(ref, ".")
	gvk := schema.GroupVersionKind{
		Kind:    split[len(split)-1],
		Version: split[len(split)-2],
		Group:   split[len(split)-3],
	}
	return gvk
}

func stripPrefix(name string) string {
	const prefix = "#/definitions/"
	return strings.TrimPrefix(name, prefix)
}

// extractDeprecationComment returns the comment with deprecation comment removed and the extracted deprecation comment.
func extractDeprecationComment(comment any, gvk schema.GroupVersionKind) (string, string) {
	if comment == nil {
		return "", ""
	}

	commentstr, _ := comment.(string)
	if commentstr == "" {
		return "", ""
	}

	re := regexp.MustCompile(`DEPRECATED - .* is deprecated by .* for more information\.\s*`)

	if re.MatchString(commentstr) {
		deprecationMessage := APIVersionComment(gvk)
		return re.ReplaceAllString(commentstr, ""), deprecationMessage
	}

	return commentstr, ""
}

func fmtComment(comment any) string {
	if comment == nil {
		return ""
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

		return commentstr
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

func makeSchemaTypeSpec(prop map[string]any, canonicalGroups map[string]string) pschema.TypeSpec {
	if t, exists := prop["type"]; exists {
		switch t := t.(string); t {
		case "array":
			elemSpec := makeSchemaTypeSpec(prop["items"].(map[string]any), canonicalGroups)
			return pschema.TypeSpec{
				Type:  "array",
				Items: &elemSpec,
			}
		case "object":
			additionalProperties, ok := prop["additionalProperties"]
			if !ok {
				return pschema.TypeSpec{Type: "object"}
			}

			elemSpec := makeSchemaTypeSpec(additionalProperties.(map[string]any), canonicalGroups)
			return pschema.TypeSpec{
				Type:                 "object",
				AdditionalProperties: &elemSpec,
			}
		default:
			return pschema.TypeSpec{Type: t}
		}
	}

	// Handle objects with `x-preserve-unknown-fields` set to true.
	if preserveUnknownFields, ok := prop["x-kubernetes-preserve-unknown-fields"]; ok {
		if preserveUnknownFields.(bool) {
			return pschema.TypeSpec{
				Type:                 "object",
				AdditionalProperties: &pschema.TypeSpec{Ref: "pulumi.json#/Any"},
			}
		}
	}

	// Handle objects with `x-kubernetes-int-or-string`.
	if _, ok := prop["x-kubernetes-int-or-string"]; ok {
		return pschema.TypeSpec{OneOf: []pschema.TypeSpec{
			{Type: "integer"},
			{Type: "string"},
		}}
	}

	ref := stripPrefix(prop["$ref"].(string))
	switch ref {
	case quantity:
		return pschema.TypeSpec{Type: "string"}
	case intOrString:
		return pschema.TypeSpec{OneOf: []pschema.TypeSpec{
			{Type: "integer"},
			{Type: "string"},
		}}
	case v1Fields, v1FieldsV1, rawExtension:
		return pschema.TypeSpec{
			Type: "object",
			Ref:  "pulumi.json#/Json",
		}
	case v1Time, v1MicroTime:
		return pschema.TypeSpec{Type: "string"}
	case v1beta1JSONSchemaPropsOrBool:
		return pschema.TypeSpec{OneOf: []pschema.TypeSpec{
			{Ref: "#/types/kubernetes:apiextensions.k8s.io/v1beta1:JSONSchemaProps"},
			{Type: "boolean"},
		}}
	case v1JSONSchemaPropsOrBool:
		return pschema.TypeSpec{OneOf: []pschema.TypeSpec{
			{Ref: "#/types/kubernetes:apiextensions.k8s.io/v1:JSONSchemaProps"},
			{Type: "boolean"},
		}}
	case v1beta1JSONSchemaPropsOrArray:
		return pschema.TypeSpec{OneOf: []pschema.TypeSpec{
			{Ref: "#/types/kubernetes:apiextensions.k8s.io/v1beta1:JSONSchemaProps"},
			{
				Type:  "array",
				Items: &pschema.TypeSpec{Ref: "pulumi.json#/Json"},
			},
		}}
	case v1JSONSchemaPropsOrArray:
		return pschema.TypeSpec{OneOf: []pschema.TypeSpec{
			{Ref: "#/types/kubernetes:apiextensions.k8s.io/v1:JSONSchemaProps"},
			{
				Type:  "array",
				Items: &pschema.TypeSpec{Ref: "pulumi.json#/Json"},
			},
		}}
	case v1beta1JSONSchemaPropsOrStringArray:
		return pschema.TypeSpec{OneOf: []pschema.TypeSpec{
			{Ref: "#/types/kubernetes:apiextensions.k8s.io/v1beta1:JSONSchemaProps"},
			{
				Type:  "array",
				Items: &pschema.TypeSpec{Type: "string"},
			},
		}}
	case v1JSONSchemaPropsOrStringArray:
		return pschema.TypeSpec{OneOf: []pschema.TypeSpec{
			{Ref: "#/types/kubernetes:apiextensions.k8s.io/v1:JSONSchemaProps"},
			{
				Type:  "array",
				Items: &pschema.TypeSpec{Type: "string"},
			},
		}}
	case v1beta1JSON, v1beta1CRSubresourceStatus, v1JSON, v1CRSubresourceStatus:
		return pschema.TypeSpec{Ref: "pulumi.json#/Json"}
	}

	gvk := GVKFromRef(ref)
	if canonicalGroup, ok := canonicalGroups[gvk.Group]; ok {
		return pschema.TypeSpec{Ref: fmt.Sprintf("#/types/kubernetes:%s/%s:%s",
			canonicalGroup, gvk.Version, gvk.Kind)}
	}
	panic("Canonical group not set for ref: " + ref)
}

func makeSchemaType(prop map[string]any, canonicalGroups map[string]string) string {
	spec := makeSchemaTypeSpec(prop, canonicalGroups)
	b, err := json.Marshal(spec)
	contract.AssertNoErrorf(err, "unexpected error while marshaling JSON")
	return string(b)
}

// --------------------------------------------------------------------------

// Core grouping logic.

// --------------------------------------------------------------------------

func createGroups(definitionsJSON map[string]any, allowHyphens bool) []GroupConfig {
	// Map Group -> canonical Group
	// e.g., flowcontrol -> flowcontrol.apiserver.k8s.io
	canonicalGroups := map[string]string{
		"meta": "meta", // "meta" Group doesn't include the `x-kubernetes-group-version-kind` field.
	}
	linq.From(definitionsJSON).
		SelectT(func(kv linq.KeyValue) definition {
			defName := kv.Key.(string)
			gvk := GVKFromRef(defName)
			def := definition{
				gvk:  gvk,
				name: defName,
				data: definitionsJSON[defName].(map[string]any),
			}
			// Top-level kinds include a canonical GVK.
			if gvks, gvkExists := def.data["x-kubernetes-group-version-kind"].([]any); gvkExists && len(gvks) > 0 {
				gvk := gvks[0].(map[string]any)
				group := gvk["group"].(string)
				// The "core" group shows up as "" in the OpenAPI spec.
				if group == "" && def.gvk.Group == "core" {
					group = "core"
				}
				def.canonicalGroup = group
			}
			return def
		}).
		WhereT(func(d definition) bool { return d.canonicalGroup != "" }).
		ToMapByT(&canonicalGroups,
			func(d definition) string { return d.gvk.Group },
			func(d definition) string { return d.canonicalGroup })

	// Map definition JSON object -> `definition` with metadata.
	var definitions []definition
	linq.From(definitionsJSON).
		SelectT(func(kv linq.KeyValue) definition {
			defName := kv.Key.(string)
			gvk := GVKFromRef(defName)
			def := definition{
				gvk:  gvk,
				name: defName,
				data: definitionsJSON[defName].(map[string]any),
			}
			if canonicalGroup, ok := canonicalGroups[gvk.Group]; ok {
				def.canonicalGroup = canonicalGroup
			} else {
				def.canonicalGroup = gvk.Group
			}
			return def
		}).
		ToSlice(&definitions)

	// Compute aliases for Kinds. Many k8s resources have multiple GVs, so create a map from Kind -> GV string.
	// For Kinds with more than one GV, create aliases in the SDKs.
	aliases := map[string][]any{}
	linq.From(definitions).
		WhereT(func(d definition) bool { return d.isTopLevel() && !strings.HasSuffix(d.gvk.Kind, "List") }).
		OrderByT(func(d definition) string { return d.gvk.String() }).
		SelectManyT(func(d definition) linq.Query {
			return linq.From([]KindConfig{
				{
					kind:       d.gvk.Kind,
					apiVersion: d.apiVersion(canonicalGroups),
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
			func(i any) any {
				return i.(linq.Group).Key
			},
			func(i any) any {
				return i.(linq.Group).Group
			})
	aliasesForKind := func(kind, apiVersion string) []string {
		var results []string

		for _, alias := range aliases[kind] {
			aliasString := alias.(string)
			re := fmt.Sprintf(`:%s:`, apiVersion)
			match, err := regexp.MatchString(re, aliasString)
			if err == nil && match {
				continue
			}
			results = append(results, aliasString)

			switch kind {
			case "CSIStorageCapacity":
				results = append(results, "kubernetes:storage.k8s.io/v1alpha1:CSIStorageCapacity")
			}

			// "apiregistration.k8s.io" was previously called "apiregistration", so create aliases for backward compat
			if strings.Contains(apiVersion, "apiregistration.k8s.io") {
				parts := strings.Split(aliasString, ":")
				parts[1] = "apiregistration" + strings.TrimPrefix(parts[1], "apiregistration.k8s.io")
				results = append(results, strings.Join(parts, ":"))
			}
		}

		// "apiregistration.k8s.io" was previously called "apiregistration", so create aliases for backward compat
		if strings.Contains(apiVersion, "apiregistration.k8s.io") {
			results = append(results, fmt.Sprintf("kubernetes:%s:%s",
				"apiregistration"+strings.TrimPrefix(apiVersion, "apiregistration.k8s.io"), kind))
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

			defaultAPIVersion := d.defaultAPIVersion()
			isTopLevel := d.isTopLevel()

			ps := linq.From(d.data["properties"]).
				OrderByT(func(kv linq.KeyValue) string { return kv.Key.(string) }).
				SelectT(func(kv linq.KeyValue) Property {
					propName := kv.Key.(string)
					prop := d.data["properties"].(map[string]any)[propName].(map[string]any)

					schemaType := makeSchemaType(prop, canonicalGroups)

					// `-` is invalid in variable names, so replace with `_`
					switch propName {
					case "x-kubernetes-embedded-resource":
						propName = "x_kubernetes_embedded_resource"
					case "x-kubernetes-int-or-string":
						propName = "x_kubernetes_int_or_string"
					case "x-kubernetes-list-map-keys":
						propName = "x_kubernetes_list_map_keys"
					case "x-kubernetes-list-type":
						propName = "x_kubernetes_list_type"
					case "x-kubernetes-map-type":
						propName = "x_kubernetes_map_type"
					case "x-kubernetes-preserve-unknown-fields":
						propName = "x_kubernetes_preserve_unknown_fields"
					case "x-kubernetes-validations":
						propName = "x_kubernetes_validations" //nolint:gosec
					}

					if !allowHyphens {
						contract.Assertf(!strings.Contains(propName, "-"), "property names may not contain `-`")
					}

					// Create a const value for the field.
					var constValue string

					// Create a default value for the field.
					switch propName {
					case "apiVersion":
						if d.isTopLevel() {
							constValue = defaultAPIVersion
						}
					case "kind":
						if d.isTopLevel() {
							constValue = d.gvk.Kind
						}
					}

					return Property{
						comment:    fmtComment(prop["description"]),
						schemaType: schemaType,
						name:       propName,
						constValue: constValue,
					}
				})

			// All properties.
			var properties []Property
			ps.ToSlice(&properties)

			// Required properties.
			reqdProps := sets.NewString()
			if reqd, hasReqd := d.data["required"]; hasReqd {
				for _, propName := range reqd.([]any) {
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

			comment, deprecationComment := extractDeprecationComment(d.data["description"], d.gvk)

			apiVersion := d.apiVersion(canonicalGroups)
			schemaPkgName := func(gv string) string {
				pkgName := strings.Replace(gv, ".k8s.io", "", -1)
				parts := strings.Split(pkgName, "/")
				contract.Assertf(len(parts) == 2, "expected package name to have two parts: %s", pkgName)
				g, v := parts[0], parts[1]
				gParts := strings.Split(g, ".")

				// We need to sanitize versions to be valid package names.
				v = validCharRegex.ReplaceAllString(v, "_")
				gStripped := validCharRegex.ReplaceAllString(gParts[0], "_")

				return fmt.Sprintf("%s/%s", gStripped, v)
			}
			return linq.From([]KindConfig{
				{
					kind:                    d.gvk.Kind,
					deprecationComment:      fmtComment(deprecationComment),
					comment:                 fmtComment(comment),
					pulumiComment:           fmtComment(PulumiComment(d.gvk.Kind)),
					properties:              properties,
					requiredInputProperties: requiredInputProperties,
					optionalInputProperties: optionalInputProperties,
					aliases:                 aliasesForKind(d.gvk.Kind, apiVersion),
					gvk:                     d.gvk,
					apiVersion:              apiVersion,
					defaultAPIVersion:       defaultAPIVersion,
					isNested:                !isTopLevel,
					schemaPkgName:           schemaPkgName(apiVersion),
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

			return linq.From([]VersionConfig{
				{
					version:           gv.Version,
					kinds:             kindsGroup,
					gv:                gv,
					apiVersion:        kindsGroup[0].apiVersion,        // NOTE: This is safe.
					defaultAPIVersion: kindsGroup[0].defaultAPIVersion, // NOTE: This is safe.
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

			return linq.From([]GroupConfig{
				{
					group:    group,
					versions: versionsGroup,
				},
			})
		}).
		WhereT(func(gc GroupConfig) bool {
			return len(gc.Versions()) != 0
		}).
		ToSlice(&groups)

	return groups
}
