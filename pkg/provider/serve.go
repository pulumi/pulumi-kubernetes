// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package provider

import (
	"github.com/pulumi/pulumi/pkg/resource/provider"
	"github.com/pulumi/pulumi/pkg/util/cmdutil"
	lumirpc "github.com/pulumi/pulumi/sdk/proto/go"

	// Load auth plugins. Removing this will likely cause compilation error.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// Serve launches the gRPC server for the Pulumi Kubernetes resource provider.
func Serve(providerName, version string) {
	// Start gRPC service.
	err := provider.Main(
		providerName, func(host *provider.HostClient) (lumirpc.ResourceProviderServer, error) {
			return makeKubeProvider(providerName, version)
		})

	if err != nil {
		cmdutil.ExitError(err.Error())
	}
}
