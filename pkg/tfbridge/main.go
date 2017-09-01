// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"github.com/pulumi/pulumi-fabric/pkg/util/cmdutil"
)

// Main launches the tfbridge plugin for a given package pkg and provider prov.
func Main(pkg string, prov ProviderInfo) {
	if err := Serve(pkg, prov); err != nil {
		cmdutil.ExitError(err.Error())
	}
}
