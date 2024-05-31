
# Provider Package

This package implements the Kubernetes provider. It provides Pulumi custom resources representing Kubernetes kinds
and Pulumi component resources for working with Kubernetes manifests and with tools such as Helm and Kustomize. 
It also provides some invokes to support various "overlay" resources that are implemented at the SDK level.

## Custom Resources

## Component Resources

Each component is implemented via a sub-provider as defined by the `resource.ResourceProvider` interface.

Steps to add a new component provider:
1. Create a new package for the implementation, e.g. `provider/pkg/provider/helm/v4`.
2. Add an alias to the package for each SDK in `provider/pkg/gen/schema.go`. You want your package name to be idiomatic to the SDK.
3. Add an exclusion to `resourcesToFilterFromTemplate` in `provider/cmd/provider-gen-kubernetes/main.go`.
4. Define a resource schema in `provider/pkg/gen/overlays.go`. The term *overlay* here means "not a Kubernetes kind".
5. Create a documentation file for the resource in `provider/pkg/gen/examples/overlays/kustomizeDirectory.md` and link from `overlays.go`.
6. Create an implementation of `resource.ResourceProvider` into your new implementation package.
7. Register the implementation into `resourceProviders` in `provider/pkg/provider/provider_construct.go`.
