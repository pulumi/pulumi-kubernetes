// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package main

import (
	"github.com/pulumi/lumi/pkg/util/cmdutil"
	"github.com/pulumi/terraform-bridge/pkg/tfbridge"
)

func main() {
	// TODO: don't hard-code AWS; we need to arrange for this to be passed somehow from the provider package.
	//     I am currently thinking of adding the ability for resource provider package manifests (Lumi.yamls) to
	//     specify additional arguments passed to their providers at startup time.
	if err := tfbridge.Serve("aws"); err != nil {
		cmdutil.ExitError(err.Error())
	}
}
