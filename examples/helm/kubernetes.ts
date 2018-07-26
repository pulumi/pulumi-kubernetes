import * as pulumi from "@pulumi/pulumi";

export class KubernetesProvider extends pulumi.CustomResource {
    constructor(name: string, config: pulumi.Inputs, opts?: pulumi.ResourceOptions) {
        super("pulumi-providers:provider:kubernetes", name, config, opts)
    }
}