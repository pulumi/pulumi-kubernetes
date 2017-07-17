// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

// Package providers contains all of the known Terraform bridge providers.
package providers

import (
	"github.com/pulumi/aws"
	"github.com/pulumi/azure"
	"github.com/pulumi/gcp"

	"github.com/pulumi/terraform-bridge/pkg/tfbridge"
)

// Providers returns a map of all known Terraform providers from which we will generate packages.  It would
// be nice to be able to do this dynamically, so that we could invoke the command line with whatever source providers we
// want.  Sadly, Go's dynamic plugin support is still iffy -- and non-existent on anything but Linux -- and so for now
// we will simply statically link in all of the source providers.  Hey, it works.
var Providers = map[string]tfbridge.ProviderInfo{
	"aws":   aws.Provider(),
	"azure": azure.Provider(),
	"gcp":   gcp.Provider(),
}
