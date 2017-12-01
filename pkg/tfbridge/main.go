// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"github.com/hashicorp/terraform/helper/logging"

	"github.com/pulumi/pulumi/pkg/util/cmdutil"
)

// Main launches the tfbridge plugin for a given package pkg and provider prov.
func Main(pkg string, version string, prov ProviderInfo) {
	// Initialize Terraform logging.
	logging.SetOutput()

	if err := Serve(pkg, version, prov); err != nil {
		cmdutil.ExitError(err.Error())
	}
}
