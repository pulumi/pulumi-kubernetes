// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	"log"

	"github.com/pulumi/pulumi-fabric/pkg/diag"
	"github.com/pulumi/pulumi-fabric/pkg/resource/provider"
	lumirpc "github.com/pulumi/pulumi-fabric/sdk/proto/go"
)

// Serve fires up a Lumi resource provider listening to inbound gRPC traffic,
// and translates calls from Lumi into actions against the provided Terraform Provider.
func Serve(module string, info ProviderInfo) error {
	// Create a new resource provider server and listen for and serve incoming connections.
	return provider.Main(func(host *provider.HostClient) (lumirpc.ResourceProviderServer, error) {
		// Set up a log redirector to capture Terraform provider logging and only pass through those that we need.
		log.SetOutput(&LogRedirector{
			writers: map[string]func(string) error{
				tfTracePrefix: func(msg string) error { return host.Log(diag.Debug, msg) },
				tfDebugPrefix: func(msg string) error { return host.Log(diag.Debug, msg) },
				tfInfoPrefix:  func(msg string) error { return host.Log(diag.Info, msg) },
				tfWarnPrefix:  func(msg string) error { return host.Log(diag.Warning, msg) },
				tfErrorPrefix: func(msg string) error { return host.Log(diag.Error, msg) },
			},
		})

		// Create a new bridge provider.
		return NewProvider(host, module, info.P, info), nil
	})
}
