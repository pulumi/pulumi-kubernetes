import * as pulumi from "@pulumi/pulumi";

/**
 * The provider type for the kubernetes package.
 */
export class Provider extends pulumi.ProviderResource {
    /**
     * Create a Provider resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args: ProviderArgs, opts?: pulumi.ResourceOptions) {
        const props: pulumi.Inputs = {
            "cluster": args ? args.cluster : undefined,
            "context": args ? args.context : undefined,
            "kubeconfig": args ? args.kubeconfig : undefined,
            "namespace": args ? args.namespace : undefined,
            "enableDryRun": args && args.enableDryRun ? "true" : undefined,
            "suppressDeprecationWarnings": args && args.suppressDeprecationWarnings ? "true" : undefined
        };
        super("kubernetes", name, props, opts);
    }
}

/**
 * The set of arguments for constructing a Provider.
 */
export interface ProviderArgs {
    /**
     * If present, the name of the kubeconfig cluster to use.
     */
    readonly cluster?: pulumi.Input<string>;
    /**
     * If present, the name of the kubeconfig context to use.
     */
    readonly context?: pulumi.Input<string>;
    /**
     * The contents of a kubeconfig file. If this is set, this config will be used instead of $KUBECONFIG.
     */
    readonly kubeconfig?: pulumi.Input<string>;
    /**
     * If present, the default namespace to use. This flag is ignored for cluster-scoped resources.
     *
     * A namespace can be specified in multiple places, and the precedence is as follows:
     * 1. `.metadata.namespace` set on the resource.
     * 2. This `namespace` parameter.
     * 3. `namespace` set for the active context in the kubeconfig.
     */
    readonly namespace?: pulumi.Input<string>;
    /**
     * BETA FEATURE - If present and set to true, enable server-side diff calculations.
     * This feature is in developer preview, and is disabled by default.
     */
    readonly enableDryRun?: pulumi.Input<boolean>;
    /**
     * If present and set to true, suppress apiVersion deprecation warnings from the CLI.
     */
    readonly suppressDeprecationWarnings?: pulumi.Input<boolean>;
}
