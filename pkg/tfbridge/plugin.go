// Copyright 2016-2017, Pulumi Corporation.  All rights reserved.

package tfbridge

import (
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/plugin/discovery"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pkg/errors"
)

// Plugin resolves the path to a Terraform plugin, loads it, and returns two connections to it: one is a standard
// plugin client that can be used to manage its lifetime and the other is a typed provider interface.
func Plugin(provBin string) (*goplugin.Client, terraform.ResourceProvider, error) {
	// Resolve the path to a plugin.
	plugins := discovery.ResolvePluginPaths([]string{provBin})
	if len(plugins) == 0 {
		return nil, nil, errors.Errorf("No Terraform plugin found at path '%v'", provBin)
	}
	// If multiple were returned (e.g., the path wasn't specific enough), we will choose the newest one.
	plug := plugins.Newest()

	// Now fire up the plugin process and connect to it with a client.
	client := plugin.Client(plug)
	rpcClient, err := client.Client()
	if err != nil {
		return nil, nil, err
	}
	raw, err := rpcClient.Dispense(plugin.ProviderPluginName)
	if err != nil {
		return nil, nil, err
	}
	return client, raw.(terraform.ResourceProvider), nil
}
