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

// Providers returns a map of all known Terraform providers from which we will generate packages.  It would
// be nice to be able to do this dynamically, so that we could invoke the command line with whatever source providers we
// want.  Sadly, Go's dynamic plugin support is still iffy -- and non-existent on anything but Linux -- and so for now
// we will simply statically link in all of the source providers.  Hey, it works.
var Providers = map[string]ProviderInfo{
	"aws": awsProvider(),
}

// ProviderInfo contains information about a Terraform provider plugin that we will use to generate the Lumi
// metadata.  It primarily contains a pointer to the Terraform schema, but can also contain specific name translations.
type ProviderInfo struct {
	P         *schema.Provider        // the TF provider/schema.
	Git       GitInfo                 // the info about this provider's Git repo.
	Resources map[string]ResourceInfo // a map of TF name to Lumi name; if a type is missing, standard mangling occurs.
}

// ResourceInfo is a top-level type exported by a provider.
type ResourceInfo struct {
	Tok    tokens.Type           // a type token to override the default; "" uses the default.
	Fields map[string]SchemaInfo // a map of custom field names; if a type is missing, the default is used.
}

// SchemaInfo contains optional name transformations to apply.
type SchemaInfo struct {
	Name   string                // a name to override the default; "" uses the default.
	Fields map[string]SchemaInfo // a map of custom field names; if a type is missing, the default is used.
}

// GitInfo contains Git information about a provider.
type GitInfo struct {
	Repo      string // the Git repo for this provider.
	Taggish   string // the Git tag info for this provider.
	Commitish string // the Git commit info for this provider.
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
