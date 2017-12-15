// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/tokens"
)

// ProviderInfo contains information about a Terraform provider plugin that we will use to generate the Pulumi
// metadata.  It primarily contains a pointer to the Terraform schema, but can also contain specific name translations.
type ProviderInfo struct {
	P           *schema.Provider           // the TF provider/schema.
	Name        string                     // the TF provider name (e.g. terraform-provider-XXXX).
	Config      map[string]*SchemaInfo     // a map of TF name to config schema overrides.
	Resources   map[string]*ResourceInfo   // a map of TF name to Pulumi name; standard mangling occurs if no entry.
	DataSources map[string]*DataSourceInfo // a map of TF name to Pulumi resource info.
	Overlay     *OverlayInfo               // optional overlay information for augmented code-generation.
}

// ResourceInfo is a top-level type exported by a provider.  This structure can override the type to generate.  It can
// also give custom metadata for fields, using the SchemaInfo structure below.  Finally, a set of composite keys can be
// given; this is used when Terraform needs more than just the ID to uniquely identify and query for a resource.
type ResourceInfo struct {
	Tok                 tokens.Type            // a type token to override the default; "" uses the default.
	Fields              map[string]*SchemaInfo // a map of custom field names; if a type is missing, uses the default.
	IDFields            []string               // an optional list of ID alias fields.
	Docs                *DocInfo               // overrides for finding and mapping TF docs.
	DeleteBeforeReplace bool                   // if true, Pulumi will delete before creating new replacement resources.
}

// DataSourceInfo can be used to override a data source's standard name mangling and argument/return information.
type DataSourceInfo struct {
	Tok    tokens.ModuleMember
	Fields map[string]*SchemaInfo
	Docs   *DocInfo // overrides for finding and mapping TF docs.
}

// SchemaInfo contains optional name transformations to apply.
type SchemaInfo struct {
	Name     string                 // a name to override the default; "" uses the default.
	Type     tokens.Type            // a type to override the default; "" uses the default.
	AltTypes []tokens.Type          // alternative types that can be used instead of the override.
	Elem     *SchemaInfo            // a schema override for elements for arrays, maps, and sets.
	Fields   map[string]*SchemaInfo // a map of custom field names; if a type is missing, the default is used.
	Asset    *AssetTranslation      // a map of asset translation information, if this is an asset.
	Default  *DefaultInfo           // an optional default directive to be applied if a value is missing.
	Stable   *bool                  // to override whether a property is stable or not.
}

// DocInfo contains optional overrids for finding and mapping TD docs.
type DocInfo struct {
	Source                         string // an optional override to locate TF docs; "" uses the default.
	IncludeAttributesFrom          string // optionally include attributes from another raw resource for docs.
	IncludeArgumentsFrom           string // optionally include arguments from another raw resource for docs.
	IncludeAttributesFromArguments string // optionally include attributes from another raw resource's arguments.
}

// HasDefault returns true if there is a default value for this property.
func (info SchemaInfo) HasDefault() bool {
	return info.Default != nil
}

// DefaultInfo lets fields get default values at runtime, before they are even passed to Terraform.
type DefaultInfo struct {
	From  func(res *PulumiResource) (interface{}, error) // a transformation from other resource properties.
	Value interface{}                                    // a raw value to inject.
}

// PulumiResource is just a little bundle that carries URN and properties around.
type PulumiResource struct {
	URN        resource.URN
	Properties resource.PropertyMap
}

// OverlayInfo contains optional overlay information.  Each info has a 1:1 correspondence with a module and permits
// extra files to be included from the overlays/ directory when building up packs/.  This allows augmented
// code-generation for convenient things like helper functions, modules, and gradual typing.
type OverlayInfo struct {
	Files           []string                // additional files to include in the index file.
	Modules         map[string]*OverlayInfo // extra modules to inject into the structure.
	Dependencies    map[string]string       // NPM dependencies to add to package.json.
	DevDependencies map[string]string       // NPM dev-dependencies to add to package.json.
}
