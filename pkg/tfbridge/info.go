// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"github.com/pulumi/lumi/pkg/tokens"
)

// ProviderInfo contains information about a Terraform provider plugin that we will use to generate the Lumi
// metadata.  It primarily contains a pointer to the Terraform schema, but can also contain specific name translations.
type ProviderInfo struct {
	P         *schema.Provider        // the TF provider/schema.
	Git       GitInfo                 // the info about this provider's Git repo.
	Config    map[string]SchemaInfo   // a map of TF name to config schema overrides.
	Resources map[string]ResourceInfo // a map of TF name to Lumi name; if a type is missing, standard mangling occurs.
	Overlay   OverlayInfo             // optional overlay information for augmented code-generation.
}

// GitInfo contains Git information about a provider.
type GitInfo struct {
	Repo      string // the Git repo for this provider.
	Taggish   string // the Git tag info for this provider.
	Commitish string // the Git commit info for this provider.
}

// ResourceInfo is a top-level type exported by a provider.  This structure can override the type to generate.  It can
// also give custom metadata for fields, using the SchemaInfo structure below.  Finally, a set of composite keys can be
// given; this is used when Terraform needs more than just the ID to uniquely identify and query for a resource.
type ResourceInfo struct {
	Tok       tokens.Type           // a type token to override the default; "" uses the default.
	Fields    map[string]SchemaInfo // a map of custom field names; if a type is missing, the default is used.
	IDFields  []string              // an optional list of ID alias fields.
	KeyFields []string              // an optional list of composite key fields.
}

// SchemaInfo contains optional name transformations to apply.
type SchemaInfo struct {
	Name    string                // a name to override the default; "" uses the default.
	Type    tokens.Type           // a type to override the default; "" uses the default.
	Elem    *SchemaInfo           // a schema override for elements for arrays, maps, and sets.
	Fields  map[string]SchemaInfo // a map of custom field names; if a type is missing, the default is used.
	Asset   *AssetTranslation     // a map of asset translation information, if this is an asset.
	Default DefaultInfo           // an optional default directive to be applied if a value is missing.
}

// Optional returns true if the value for this entry isn't required.
func (info SchemaInfo) Optional() bool {
	return info.Default.From != "" || info.Default.Value != nil
}

// DefaultInfo lets fields get default values at runtime, before they are even passed to Terraform.
type DefaultInfo struct {
	From          string                        // to take a default from another field.
	FromTransform func(interface{}) interface{} // an optional transformation to apply to the from value.
	Value         interface{}                   // a raw value to inject.
}

// OverlayInfo contains optional overlay information.  Each info has a 1:1 correspondence with a module and permits
// extra files to be included from the overlays/ directory when building up packs/.  This allows augmented
// code-generation for convenient things like helper functions, modules, and gradual typing.
type OverlayInfo struct {
	Files   []string
	Modules map[string]OverlayInfo
}

const (
	tfGitHub         = "github.com"
	tfProvidersOrg   = "terraform-providers"
	tfProviderPrefix = "terraform-provider"
)

// GetGitInfo fetches the taggish and commitish info for a provider's repo using a standard GOPATH location.
func GetGitInfo(prov string) (GitInfo, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return GitInfo{}, errors.New("GOPATH is not set; canot read provider's Git info")
	}
	tfdir := filepath.Join(gopath, "src", tfGitHub, tfProvidersOrg, tfProviderPrefix+"-"+prov)
	descCmd := exec.Command("git", "describe", "--all", "--long")
	descCmd.Dir = tfdir
	descOut, err := descCmd.Output()
	if err != nil {
		return GitInfo{}, err
	}
	showRefCmd := exec.Command("git", "show-ref", "HEAD")
	showRefCmd.Dir = tfdir
	showRefOut, err := showRefCmd.Output()
	if err != nil {
		return GitInfo{}, err
	}
	return GitInfo{
		Repo:      tfGitHub + "/" + tfProvidersOrg + "/" + tfProviderPrefix + "-" + prov,
		Taggish:   string(descOut),
		Commitish: string(showRefOut),
	}, nil
}
