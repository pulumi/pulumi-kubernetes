package main

// package main

import (
	"github.com/pulumi/pulumi-kubernetes/pkg/provider"
	"github.com/pulumi/pulumi-kubernetes/pkg/version"
)

var providerName = "kubernetes"

func main() {
	provider.Serve(providerName, version.Version)
}
