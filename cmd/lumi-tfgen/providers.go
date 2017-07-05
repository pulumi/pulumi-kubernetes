// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package main

import (
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-aws/aws"
)

// terraformProviders returns a map of all known Terraform providers from which we will generate packages.  It would
// be nice to be able to do this dynamically, so that we could invoke the command line with whatever source providers we
// want.  Sadly, Go's dynamic plugin support is still iffy -- and non-existent on anything but Linux -- and so for now
// we will simply statically link in all of the source providers.  Hey, it works.
func terraformProviders() map[string]*schema.Provider {
	return map[string]*schema.Provider{
		"aws": aws.Provider().(*schema.Provider),
	}
}
