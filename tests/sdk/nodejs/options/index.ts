import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import * as time from "@pulumiverse/time";

const config = new pulumi.Config();

const bootstrapProvider = new k8s.Provider("bootstrap", {});

// create a set of providers across namespaces, simply to facilitate the reuse of manifests in the below tests.
const nulloptsNs = new k8s.core.v1.Namespace("nullopts", {}, { provider: bootstrapProvider });
const aNs = new k8s.core.v1.Namespace("a", {}, { provider: bootstrapProvider });
const bNs = new k8s.core.v1.Namespace("b", {}, { provider: bootstrapProvider });
const nulloptsProvider = new k8s.Provider("nullopts", { namespace: nulloptsNs.metadata["name"] })
const aProvider = new k8s.Provider("a", { namespace: aNs.metadata["name"] })
const bProvider = new k8s.Provider("b", { namespace: bNs.metadata["name"] })

// a sleep resource to exercise the "depends_on" component-level option
const sleep = new time.Sleep("sleep", { createDuration: "1s" }, { dependsOn: [aProvider, bProvider] });

// apply_default_opts is a stack transformation that applies default opts to any resource whose name ends with "-nullopts".
// this is intended to be applied to component resources only.
function applyDefaultOpts(args: pulumi.ResourceTransformationArgs): pulumi.ResourceTransformationResult | undefined {
    if (args.name.endsWith("-nullopts")) {
        return {
            props: args.props,
            opts: pulumi.mergeOptions(args.opts, {
                provider: nulloptsProvider,
            }),
        };
    }
    return undefined;
}
pulumi.runtime.registerStackTransformation(applyDefaultOpts);

// applyAlias is a Pulumi transformation that applies a unique alias to each resource.
function applyAlias(args: pulumi.ResourceTransformationArgs): pulumi.ResourceTransformationResult | undefined {
    return {
        props: args.props,
        opts: pulumi.mergeOptions(args.opts, {
            aliases: [{ name: `${args.name}-aliased` }],
            // aliases: [{ name: `${args.name}-aliased` }, ...(args.opts.aliases ?? [])],
        }),
    };
}

// transform_k8s is a Kubernetes transformation that applies a unique alias and annotation to each resource.
function transformK8s(obj: any, opts: pulumi.CustomResourceOptions) {
    opts.aliases = [{ name: `${obj.metadata.name}-k8s-aliased` }, ...(opts.aliases ?? [])]
    obj.metadata.annotations = { "transformed": "true" }
}

// --- ConfigGroup ---
// options: providers, aliases, depends_on, ignore_changes, protect, transformations
// args: skip_await, transformations
new k8s.yaml.ConfigGroup("cg-options", {
    resourcePrefix: "cg-options",
    skipAwait: true,
    transformations: [transformK8s],
    files: ["./testdata/options/configgroup/*.yaml"],
    yaml: [`
apiVersion: v1
kind: ConfigMap
metadata:
  name: cg-options-cm-1
`],
},
    {
        providers: [aProvider],
        aliases: [{ name: "cg-options-old" }],
        ignoreChanges: ["ignored"],
        protect: true,
        dependsOn: [sleep],
        transformations: [applyAlias],
    });

// "providers" option
new k8s.yaml.ConfigGroup("cg-providers", {
    resourcePrefix: "cg-providers",
    yaml: [`
apiVersion: v1
kind: ConfigMap
metadata:
    name: cg-providers-cm-1
`],
}, { providers: [bProvider] });

// "provider" option
new k8s.yaml.ConfigGroup("cg-provider", {
    resourcePrefix: "cg-provider",
    yaml: [`
apiVersion: v1
kind: ConfigMap
metadata:
    name: cg-provider-cm-1
`],
}, { provider: bProvider });

// null opts
new k8s.yaml.ConfigGroup("cg-nullopts", {
    resourcePrefix: "cg-nullopts",
    yaml: [`
apiVersion: v1
kind: ConfigMap
metadata:
    name: cg-nullopts-cm-1
`],
});

//  --- ConfigFile ---
new k8s.yaml.ConfigFile("cf-options", {
    resourcePrefix: "cf-options",
    skipAwait: true,
    transformations: [transformK8s],
    file: "./testdata/options/configfile/manifest.yaml",
}, {
    providers: [aProvider],
    aliases: [{ name: "cf-options-old" }],
    ignoreChanges: ["ignored"],
    protect: true,
    dependsOn: [sleep],
    transformations: [applyAlias],
});

// "provider" option
new k8s.yaml.ConfigFile("cf-provider", {
    resourcePrefix: "cf-provider",
    file: "./testdata/options/configfile/manifest.yaml",
}, { provider: bProvider });

// null opts
new k8s.yaml.ConfigFile("cf-nullopts", {
    resourcePrefix: "cf-nullopts",
    file: "./testdata/options/configfile/manifest.yaml",
});

// empty manifest
new k8s.yaml.ConfigFile("cf-empty", {
    resourcePrefix: "cf-empty",
    file: "./testdata/options/configfile/empty.yaml",
}, { providers: [bProvider] });

// --- Directory ---
new k8s.kustomize.Directory("kustomize-options", {
    directory: "./testdata/options/kustomize",
    resourcePrefix: "kustomize-options",
    transformations: [transformK8s],
}, {
    providers: [aProvider],
    aliases: [{ name: "kustomize-options-old" }],
    ignoreChanges: ["ignored"],
    protect: true,
    dependsOn: [sleep],
    transformations: [applyAlias],
});

// "provider" option
new k8s.kustomize.Directory("kustomize-provider", {
    directory: "./testdata/options/kustomize",
    resourcePrefix: "kustomize-provider",
}, { provider: bProvider });

// null opts
new k8s.kustomize.Directory("kustomize-nullopts", {
    directory: "./testdata/options/kustomize",
    resourcePrefix: "kustomize-nullopts",
});

// --- Chart ---
// options: providers, aliases, depends_on, ignore_changes, protect, transformations
// args: transformations
new k8s.helm.v3.Chart("chart-options", {
    path: "./testdata/options/chart",
    resourcePrefix: "chart-options",
    transformations: [transformK8s],
    skipAwait: true,
}, {
    providers: [aProvider],
    aliases: [{ name: "chart-options-old" }],
    ignoreChanges: ["ignored"],
    protect: true,
    dependsOn: [sleep],
    transformations: [applyAlias],
});

// "provider" option
new k8s.helm.v3.Chart("chart-provider", {
    path: "./testdata/options/chart",
    resourcePrefix: "chart-provider",
}, { provider: bProvider });

// null opts
new k8s.helm.v3.Chart("chart-nullopts", {
    path: "./testdata/options/chart",
    resourcePrefix: "chart-nullopts",
});


