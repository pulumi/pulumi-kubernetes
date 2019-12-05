package kubernetes

import (
	"github.com/pulumi/pulumi/sdk/go/pulumi"
)

// Provider is the provider type for the kubernetes package.
type Provider struct {
	pulumi.ProviderResourceState
}

// ProviderArgs is the set of arguments for constructing a Provider.
type ProviderArgs struct {
	// If present, the name of the kubeconfig cluster to use.
	Cluster pulumi.StringInput `pulumi:"cluster"`

	// If present, the name of the kubeconfig context to use.
	Context pulumi.StringInput `pulumi:"context"`

	// The contents of a kubeconfig file. If this is set, this config will be used instead of $KUBECONFIG.
	Kubeconfig pulumi.StringInput `pulumi:"kubeconfig"`

	// If present, the default namespace to use. This flag is ignored for cluster-scoped resources.
	// Note: if .metadata.namespace is set on a resource, that value takes precedence over the provider default.
	Namespace pulumi.StringInput `pulumi:"namespace"`

	// BETA FEATURE - If present and set to true, enable server-side diff calculations.
	// This feature is in developer preview, and is disabled by default.
	EnableDryRun pulumi.BoolInput `pulumi:"enableDryRun"`

	// If present and set to true, suppress apiVersion deprecation warnings from the CLI.
	SuppressDeprecationWarnings pulumi.BoolInput `pulumi:"suppressDeprecationWarnings"`
}

// NewProvider registers a new resource with the given unique name, arguments, and options.
func NewProvider(ctx *pulumi.Context, name string, args *ProviderArgs, opts ...pulumi.ResourceOption) (*Provider, error) {
	inputs := map[string]pulumi.Input{}
	if args != nil {
		if i := args.Cluster; i != nil {
			inputs["cluster"] = i.ToStringOutput()
		}
		if i := args.Context; i != nil {
			inputs["context"] = i.ToStringOutput()
		}
		if i := args.Kubeconfig; i != nil {
			inputs["kubeconfig"] = i.ToStringOutput()
		}
		if i := args.Namespace; i != nil {
			inputs["namespace"] = i.ToStringOutput()
		}
		if i := args.EnableDryRun; i != nil {
			inputs["enableDryRun"] = i.ToStringOutput()
		}
		if i := args.SuppressDeprecationWarnings; i != nil {
			inputs["suppressDeprecationWarnings"] = i.ToStringOutput()
		}
	}
	var resource Provider
	err := ctx.RegisterResource("pulumi:providers:kubernetes", name, inputs, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}
