// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-fabric/pkg/tokens"
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
	Repo   string // the Git repo for this provider.
	Tag    string // the Git tag info for this provider.
	Commit string // the Git commit info for this provider.
}

// ResourceInfo is a top-level type exported by a provider.  This structure can override the type to generate.  It can
// also give custom metadata for fields, using the SchemaInfo structure below.  Finally, a set of composite keys can be
// given; this is used when Terraform needs more than just the ID to uniquely identify and query for a resource.
type ResourceInfo struct {
	Tok                 tokens.Type           // a type token to override the default; "" uses the default.
	Fields              map[string]SchemaInfo // a map of custom field names; if a type is missing, the default is used.
	NameFields          []string              // an optional list of fields to use as name (if not the default).
	NameFieldsDelimiter string                // an optional delimiter for name fields (if multiple).
	IDFields            []string              // an optional list of ID alias fields.
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

// HasDefault returns true if there is a default value for this property.
func (info SchemaInfo) HasDefault() bool {
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

// GetGitInfo fetches the taggish and commitish info for a provider's repo.  It prefers to use a Gopkg.lock file, in
// case dep is being used to vendor, and falls back to looking at the raw Git repo using a standard GOPATH location
// otherwise.  If neither is found, an error is returned.
func GetGitInfo(prov string) (GitInfo, error) {
	repo := tfGitHub + "/" + tfProvidersOrg + "/" + tfProviderPrefix + "-" + prov

	// First look for a Gopkg.lock file.
	pkglock, err := toml.LoadFile("Gopkg.lock")
	if err == nil {
		// If no error, attempt to use the file.  Otherwise, keep looking for a Git repo.
		if projs, isprojs := pkglock.Get("projects").([]*toml.Tree); isprojs {
			for _, proj := range projs {
				if name, isname := proj.Get("name").(string); isname && name == repo {
					var tag string
					if vers, isvers := proj.Get("version").(string); isvers {
						tag = vers
					}
					var commit string
					if revs, isrevs := proj.Get("revision").(string); isrevs {
						commit = revs
					}
					if tag != "" || commit != "" {
						return GitInfo{
							Repo:   repo,
							Tag:    tag,
							Commit: commit,
						}, nil
					}
				}
			}
		}
	}

	// If that didn't work, try the GOPATH for a Git repo.
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return GitInfo{}, errors.New("GOPATH is not set; canot read provider's Git info")
	}
	repodir := filepath.Join(gopath, "src", tfGitHub, tfProvidersOrg, tfProviderPrefix+"-"+prov)

	// Make sure the target is actually a Git repository so we can fail with a pretty error if not.
	if _, staterr := os.Stat(filepath.Join(repodir, ".git")); staterr != nil {
		return GitInfo{}, errors.Errorf("%v is not a Git repo, and no vendored copy was found", repodir)
	}

	// Now launch the Git commands.
	descCmd := exec.Command("git", "describe", "--all", "--long")
	descCmd.Dir = repodir
	descOut, err := descCmd.Output()
	if err != nil {
		return GitInfo{}, err
	} else if strings.HasSuffix(string(descOut), "\n") {
		descOut = descOut[:len(descOut)-1]
	}
	showRefCmd := exec.Command("git", "show-ref", "HEAD")
	showRefCmd.Dir = repodir
	showRefOut, err := showRefCmd.Output()
	if err != nil {
		return GitInfo{}, err
	} else if strings.HasSuffix(string(showRefOut), "\n") {
		showRefOut = showRefOut[:len(showRefOut)-1]
	}
	return GitInfo{
		Repo:   repo,
		Tag:    string(descOut),
		Commit: string(showRefOut),
	}, nil
}
