// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"github.com/pulumi/lumi/pkg/resource"
	"github.com/pulumi/lumi/pkg/tokens"
	"github.com/pulumi/lumi/pkg/util/contract"
)

// Providers returns a map of all known Terraform providers from which we will generate packages.  It would
// be nice to be able to do this dynamically, so that we could invoke the command line with whatever source providers we
// want.  Sadly, Go's dynamic plugin support is still iffy -- and non-existent on anything but Linux -- and so for now
// we will simply statically link in all of the source providers.  Hey, it works.
var Providers = map[string]ProviderInfo{
	"aws":   awsProvider(),
	"azure": azureProvider(),
	"gcp":   gcpProvider(),
}

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

// ResourceInfo is a top-level type exported by a provider.
type ResourceInfo struct {
	Tok    tokens.Type           // a type token to override the default; "" uses the default.
	Fields map[string]SchemaInfo // a map of custom field names; if a type is missing, the default is used.
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

// getGitInfo fetches the taggish and commitish info for a provider's repo using a standard GOPATH location.
func getGitInfo(prov string) (GitInfo, error) {
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

// autoName creates custom schema for a Terraform name property.  It uses the new name (which must not be "name"), and
// figures out given the schema information how to populate the property.  maxlen specifies the maximum length.
func autoName(new string, schema *schema.Schema, maxlen int) SchemaInfo {
	contract.Assert(new != NameProperty)
	info := SchemaInfo{Name: new}
	if !schema.Optional {
		info.Default = DefaultInfo{
			From: NameProperty,
			FromTransform: func(v interface{}) interface{} {
				return resource.NewUniqueHex(v.(string)+"-", maxlen, -1)
			},
		}
	}
	return info
}
