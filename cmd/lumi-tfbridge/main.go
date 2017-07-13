// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pulumi/lumi/pkg/resource/plugin"
	"github.com/pulumi/lumi/pkg/util/cmdutil"

	"github.com/pulumi/terraform-bridge/pkg/tfbridge"
)

func main() {
	// To figure out the package that this process handles, we will look at the binary name.  It will be of the form:
	//
	//     lumi-resource-xyz
	//
	// where xyz is the name of the package.  For example, lumi-resource-aws.  This avoids needing any sort of
	// command line configuration and takes advantage of the way in which resource provider plugins are installed.
	bin := filepath.Base(os.Args[0])
	prefix := plugin.ProviderPluginPrefix
	if !strings.HasPrefix(bin, prefix) {
		cmdutil.ExitError("fatal: missing expected plugin prefix '%v': %v", prefix, bin)
	}
	module := bin[len(prefix):]
	if module == "" {
		cmdutil.ExitError("fatal: malformed plugin name; missing a module part: %v", bin)
	}

	// Suppress logging, since Terraform plugins will echo to log and we want to intercept it ourselves.
	// IDEA: there is undoubtedly a better way to do this.  It's too bad we are smashing all logging that happens in
	//     this process (not that we have any).  We could instead fork the process and keep the current one pristine.
	log.SetOutput(ioutil.Discard)

	// Now serve it up!
	if err := tfbridge.Serve(module); err != nil {
		cmdutil.ExitError(err.Error())
	}
}
