// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pulumi/lumi/pkg/tokens"

	"github.com/terraform-providers/terraform-provider-aws/aws"
)

// Providers returns a map of all known Terraform providers from which we will generate packages.  It would
// be nice to be able to do this dynamically, so that we could invoke the command line with whatever source providers we
// want.  Sadly, Go's dynamic plugin support is still iffy -- and non-existent on anything but Linux -- and so for now
// we will simply statically link in all of the source providers.  Hey, it works.
func Providers() map[string]ProviderInfo {
	return map[string]ProviderInfo{
		"aws": {
			P: aws.Provider().(*schema.Provider),
		},
	}
}

// ProviderInfo contains information about a Terraform provider plugin that we will use to generate the Lumi
// metadata.  It primarily contains a pointer to the Terraform schema, but can also contain specific name translations.
type ProviderInfo struct {
	P     *schema.Provider    // the TF provider/schema.
	Types map[string]TypeInfo // a map of TF name to Lumi name; if a type is missing, standard mangling occurs.
}

// TypeInfo is a top-level type exported by a provider.
type TypeInfo struct {
	Name   tokens.Type           // a type token to override the default; "" uses the default.
	Fields map[string]SchemaInfo // a map of custom field names; if a type is missing, the default is used.
}

// SchemaInfo contains optional name transformations to apply.
type SchemaInfo struct {
	Name   tokens.Type           // a name to override the default; "" uses the default.
	Fields map[string]SchemaInfo // a map of custom field names; if a type is missing, the default is used.
}
