// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"io/ioutil"
	"log"

	"github.com/pulumi/lumi/pkg/util/cmdutil"
)

// Main launches the tfbridge plugin for a given package pkg and provider prov.
func Main(pkg string, prov ProviderInfo) {
	// Suppress logging, since Terraform plugins will echo to log and we want to intercept it ourselves.
	// IDEA: there is undoubtedly a better way to do this.  It's too bad we are smashing all logging that happens in
	//     this process (not that we have any).  We could instead fork the process and keep the current one pristine.
	log.SetOutput(ioutil.Discard)

	// Now serve it up!
	if err := Serve(pkg, prov); err != nil {
		cmdutil.ExitError(err.Error())
	}
}
