package gen

import (
	"fmt"
	"strings"

	linq "github.com/ahmetb/go-linq"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
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

	gv *schema.GroupVersion // Used for sorting.
}

// Version returns the name of the version (e.g., `apps/v1beta1` would return `v1beta1`).
func (vc *VersionConfig) Version() string { return vc.version }

// Kinds returns the set of kinds in some Kubernetes API group/version combination (e.g.,
// `apps/v1beta1` has the `Deployment` kind, etc.).
func (vc *VersionConfig) Kinds() []*KindConfig { return vc.kinds }

// KindConfig represents a Kubernetes API kind (e.g., the `Deployment` type in
// `apps/v1beta1/Deployment`).
type KindConfig struct {
	kind               string
	comment            string
	properties         []*Property
	requiredProperties []*Property
	optionalProperties []*Property

	gvk *schema.GroupVersionKind // Used for sorting.
}

// Kind returns the name of the Kubernetes API kind (e.g., `Deployment` for
// `apps/v1beta1/Deployment`).
func (kc *KindConfig) Kind() string { return kc.kind }

// Comment returns the comments associated with some Kubernetes API kind.
func (kc *KindConfig) Comment() string { return kc.comment }

// Properties returns the list of properties that exist on some Kubernetes API kind (i.e., things
// that we will want to `.` into, like `thing.apiVersion`, `thing.kind`, `thing.metadata`, etc.).
func (kc *KindConfig) Properties() []*Property { return kc.properties }

// RequiredProperties returns the list of properties that are required to exist on some Kubernetes
// API kind (i.e., things that we will want to `.` into, like `thing.apiVersion`, `thing.kind`,
// `thing.metadata`, etc.).
func (kc *KindConfig) RequiredProperties() []*Property { return kc.requiredProperties }

// OptionalProperties returns the list of properties that are optional on some Kubernetes API kind
// (i.e., things that we will want to `.` into, like `thing.apiVersion`, `thing.kind`,
// `thing.metadata`, etc.).
func (kc *KindConfig) OptionalProperties() []*Property { return kc.optionalProperties }

// Property represents a property we want to expose on a Kubernetes API kind (i.e., things that we
// will want to `.` into, like `thing.apiVersion`, `thing.kind`, `thing.metadata`, etc.).
type Property struct {
	comment  string
	propType string
	name     string
}

// Name returns the name of the property
func (p *Property) Name() string { return p.name }

// Comment returns the comments associated with some property.
func (p *Property) Comment() string { return p.comment }

// PropType returns the type of the property.
func (p *Property) PropType() string { return p.propType }

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

func fmtComment(comment interface{}, prefix string) string {
	if comment == nil {
		return ""
	}
	commentstr, _ := comment.(string)
	if len(commentstr) > 0 {
		split := strings.Split(commentstr, "\n")
		joined := strings.Join(split, fmt.Sprintf("\n%s// ", prefix))
		return fmt.Sprintf(`// %s`, joined)
	}
	return ""
}

func makeType(prop map[string]interface{}, refPrefix string) string {
	if t, exists := prop["type"]; exists {
		tstr := t.(string)
		if tstr == "array" {
			return fmt.Sprintf("%s[]", makeType(prop["items"].(map[string]interface{}), refPrefix))
		} else if tstr == "integer" {
			return "number"
		}
		return tstr
	}

	ref := stripPrefix(prop["$ref"].(string))
	if ref == "io.k8s.apimachinery.pkg.api.resource.Quantity" {
		return "string"
	} else if ref == "io.k8s.apimachinery.pkg.util.intstr.IntOrString" {
		return "number | string"
	} else if ref == "io.k8s.apimachinery.pkg.apis.meta.v1.Time" ||
		ref == "io.k8s.apimachinery.pkg.apis.meta.v1.MicroTime" {
		// TODO: Automatically deserialized with `DateConstructor`.
		return "string"
	}

	gvk := gvkFromRef(ref)
	if refPrefix == "" {
		return fmt.Sprintf("%s.%s.%s", gvk.Group, gvk.Version, gvk.Kind)
	}
	return fmt.Sprintf("%s.%s.%s.%s", refPrefix, gvk.Group, gvk.Version, gvk.Kind)
}

func makeTypeLiteral(prop map[string]interface{}) string {
	return makeType(prop, "")
}

func makeAPITypeRef(prop map[string]interface{}) string {
	return makeType(prop, "api")
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
	api
)

func createGroups(definitionsJSON map[string]interface{}, generatorType gentype) []*GroupConfig {
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

			ps := linq.From(d.data["properties"]).
				OrderByT(func(kv linq.KeyValue) string { return kv.Key.(string) }).
				SelectT(func(kv linq.KeyValue) *Property {
					propName := kv.Key.(string)
					prop := d.data["properties"].(map[string]interface{})[propName].(map[string]interface{})
					var typeLiteral string
					switch generatorType {
					case api:
						typeLiteral = makeTypeLiteral(prop)
					case provider:
						typeLiteral = makeAPITypeRef(prop)
					default:
						panic("Unrecognized generator type")
					}

					return &Property{
						comment:  fmtComment(prop["description"], "      "),
						propType: typeLiteral,
						name:     propName,
					}
				})

			// All properties.
			properties := []*Property{}
			ps.ToSlice(&properties)

			// Required properties.
			reqdProps := sets.NewString()
			if reqd, hasReqd := d.data["required"]; hasReqd {
				for _, propName := range reqd.([]interface{}) {
					reqdProps.Insert(propName.(string))
				}
			}

			requiredProperties := []*Property{}
			ps.
				WhereT(func(p *Property) bool {
					return reqdProps.Has(p.name)
				}).
				ToSlice(&requiredProperties)

			optionalProperties := []*Property{}
			ps.
				WhereT(func(p *Property) bool {
					return !reqdProps.Has(p.name)
				}).
				ToSlice(&optionalProperties)

			if len(properties) == 0 {
				return linq.From([]*KindConfig{})
			}

			props := d.data["properties"].(map[string]interface{})
			_, kindExists := props["kind"]
			_, apiVersionExists := props["apiVersion"]
			if generatorType == provider && (!kindExists || !apiVersionExists) {
				return linq.From([]*KindConfig{})
			}

			return linq.From([]*KindConfig{
				{
					kind: d.gvk.Kind,
					// NOTE: This transformation assumes git users on Windows to set
					// the "check in with UNIX line endings" setting.
					comment:            fmtComment(d.data["description"], "    "),
					properties:         properties,
					requiredProperties: requiredProperties,
					optionalProperties: optionalProperties,
					gvk:                &d.gvk,
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
			return linq.From([]*VersionConfig{
				{
					version: gv.Version,
					kinds:   kindsGroup,
					gv:      &gv,
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
		SelectT(func(versions linq.Group) *GroupConfig {
			versionsGroup := []*VersionConfig{}
			linq.From(versions.Group).ToSlice(&versionsGroup)
			return &GroupConfig{
				group:    versions.Key.(string),
				versions: versionsGroup,
			}
		}).
		WhereT(func(gc *GroupConfig) bool {
			return len(gc.Versions()) != 0
		}).
		ToSlice(&groups)

	return groups
}
