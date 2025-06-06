// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as inputs from "../../types/input";
import * as outputs from "../../types/output";
import * as enums from "../../types/enums";
import * as utilities from "../../utilities";

/**
 * _Looking for the Release resource? Please use the [v3 package](/registry/packages/kubernetes/api-docs/helm/v3/release/)
 * for production use cases, and stay tuned for an updated Release resource, coming soon._
 *
 * _See also: [New: Helm Chart v4 resource with new features and languages](/blog/kubernetes-chart-v4/)_
 *
 * Chart is a component representing a collection of resources described by a Helm Chart.
 * Helm charts are a popular packaging format for Kubernetes applications, and published
 * to registries such as [Artifact Hub](https://artifacthub.io/packages/search?kind=0&sort=relevance&page=1).
 *
 * Chart does not use Tiller or create a Helm Release; the semantics are equivalent to
 * running `helm template --dry-run=server` and then using Pulumi to deploy the resulting YAML manifests.
 * This allows you to apply [Pulumi Transformations](https://www.pulumi.com/docs/concepts/options/transformations/) and
 * [Pulumi Policies](https://www.pulumi.com/docs/using-pulumi/crossguard/) to the Kubernetes resources.
 *
 * You may also want to consider the `Release` resource as an alternative method for managing helm charts. For more
 * information about the trade-offs between these options, see: [Choosing the right Helm resource for your use case](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/choosing-the-right-helm-resource-for-your-use-case).
 *
 * ### Chart Resolution
 *
 * The Helm Chart can be fetched from any source that is accessible to the `helm` command line.
 * The following variations are supported:
 *
 * 1. By chart reference with repo prefix: `chart: "example/mariadb"`
 * 2. By path to a packaged chart: `chart: "./nginx-1.2.3.tgz"`
 * 3. By path to an unpacked chart directory: `chart: "./nginx"`
 * 4. By absolute URL: `chart: "https://example.com/charts/nginx-1.2.3.tgz"`
 * 5. By chart reference with repo URL: `chart: "nginx", repositoryOpts: { repo: "https://example.com/charts/" }`
 * 6. By OCI registry: `chart: "oci://example.com/charts/nginx", version: "1.2.3"`
 *
 * A chart reference is a convenient way of referencing a chart in a chart repository.
 *
 * When you use a chart reference with a repo prefix (`example/mariadb`), Pulumi will look in Helm's local configuration
 * for a chart repository named `example`, and will then look for a chart in that repository whose name is `mariadb`.
 * It will install the latest stable version of that chart, unless you specify `devel` to also include
 * development versions (alpha, beta, and release candidate releases), or supply a version number with `version`.
 *
 * Use the `verify` and optional `keyring` inputs to enable Chart verification.
 * By default, Pulumi uses the keyring at `$HOME/.gnupg/pubring.gpg`. See: [Helm Provenance and Integrity](https://helm.sh/docs/topics/provenance/).
 *
 * ### Chart Values
 *
 * [Values files](https://helm.sh/docs/chart_template_guide/values_files/#helm) (`values.yaml`) may be supplied
 * with the `valueYamlFiles` input, accepting [Pulumi Assets](https://www.pulumi.com/docs/concepts/assets-archives/#assets).
 *
 * A map of chart values may also be supplied with the `values` input, with highest precedence. You're able to use literals,
 * nested maps, [Pulumi outputs](https://www.pulumi.com/docs/concepts/inputs-outputs/), and Pulumi assets as values.
 * Assets are automatically opened and converted to a string.
 *
 * Note that the use of expressions (e.g. `--set service.type`) is not supported.
 *
 * ### Chart Dependency Resolution
 *
 * For unpacked chart directories, Pulumi automatically rebuilds the dependencies if dependencies are missing
 * and a `Chart.lock` file is present (see: [Helm Dependency Build](https://helm.sh/docs/helm/helm_dependency_build/)).
 * Use the `dependencyUpdate` input to have Pulumi update the dependencies (see: [Helm Dependency Update](https://helm.sh/docs/helm/helm_dependency_update/)).
 *
 * ### Templating
 *
 * The `Chart` resource renders the templates from your chart and then manages the resources directly with the
 * Pulumi Kubernetes provider. A default namespace is applied based on the `namespace` input, the provider's
 * configured namespace, and the active Kubernetes context. Use the `skipCrds` option to skip installing the
 * Custom Resource Definition (CRD) objects located in the chart's `crds/` special directory.
 *
 * Use the `postRenderer` input to pipe the rendered manifest through a [post-rendering command](https://helm.sh/docs/topics/advanced/#post-rendering).
 *
 * ### Resource Ordering
 *
 * Sometimes resources must be applied in a specific order. For example, a namespace resource must be
 * created before any namespaced resources, or a Custom Resource Definition (CRD) must be pre-installed.
 *
 * Pulumi uses heuristics to determine which order to apply and delete objects within the Chart.  Pulumi also
 * waits for each object to be fully reconciled, unless `skipAwait` is enabled.
 *
 * Pulumi supports the `config.kubernetes.io/depends-on` annotation to declare an explicit dependency on a given resource.
 * The annotation accepts a list of resource references, delimited by commas.
 *
 * Note that references to resources outside the Chart aren't supported.
 *
 * **Resource reference**
 *
 * A resource reference is a string that uniquely identifies a resource.
 *
 * It consists of the group, kind, name, and optionally the namespace, delimited by forward slashes.
 *
 * | Resource Scope   | Format                                         |
 * | :--------------- | :--------------------------------------------- |
 * | namespace-scoped | `<group>/namespaces/<namespace>/<kind>/<name>` |
 * | cluster-scoped   | `<group>/<kind>/<name>`                        |
 *
 * For resources in the “core” group, the empty string is used instead (for example: `/namespaces/test/Pod/pod-a`).
 *
 * ## Example Usage
 * ### Local Chart Directory
 *
 * ```typescript
 * import * as k8s from "@pulumi/kubernetes";
 *
 * const nginx = new k8s.helm.v4.Chart("nginx", {
 *     chart: "./nginx",
 * });
 * ```
 * ### Repository Chart
 *
 * ```typescript
 * import * as k8s from "@pulumi/kubernetes";
 *
 * const nginx = new k8s.helm.v4.Chart("nginx", {
 *     chart: "nginx",
 *     repositoryOpts: {
 *         repo: "https://charts.bitnami.com/bitnami",
 *     },
 * });
 * ```
 * ### OCI Chart
 *
 * ```typescript
 * import * as k8s from "@pulumi/kubernetes";
 *
 * const nginx = new k8s.helm.v4.Chart("nginx", {
 *     chart: "oci://registry-1.docker.io/bitnamicharts/nginx",
 *     version: "16.0.7",
 * });
 * ```
 * ### Chart Values
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as k8s from "@pulumi/kubernetes";
 *
 * const nginx = new k8s.helm.v4.Chart("nginx", {
 *     chart: "nginx",
 *     repositoryOpts: {
 *         repo: "https://charts.bitnami.com/bitnami",
 *     },
 *     valueYamlFiles: [
 *         new pulumi.asset.FileAsset("./values.yaml")
 *     ],
 *     values: {
 *         service: {
 *             type: "ClusterIP",
 *         },
 *         notes: new pulumi.asset.FileAsset("./notes.txt"),
 *     },
 * });
 * ```
 * ### Chart Namespace
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as k8s from "@pulumi/kubernetes";
 *
 * const ns = new k8s.core.v1.Namespace("nginx", {
 *     metadata: { name: "nginx" },
 * });
 * const nginx = new k8s.helm.v4.Chart("nginx", {
 *     namespace: ns.metadata.name,
 *     chart: "nginx",
 *     repositoryOpts: {
 *         repo: "https://charts.bitnami.com/bitnami",
 *     }
 * });
 * ```
 */
export class Chart extends pulumi.ComponentResource {
    /** @internal */
    public static readonly __pulumiType = 'kubernetes:helm.sh/v4:Chart';

    /**
     * Returns true if the given object is an instance of Chart.  This is designed to work even
     * when multiple copies of the Pulumi SDK have been loaded into the same process.
     */
    public static isInstance(obj: any): obj is Chart {
        if (obj === undefined || obj === null) {
            return false;
        }
        return obj['__pulumiType'] === Chart.__pulumiType;
    }

    /**
     * Resources created by the Chart.
     */
    public /*out*/ readonly resources!: pulumi.Output<any[]>;

    /**
     * Create a Chart resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args?: ChartArgs, opts?: pulumi.ComponentResourceOptions) {
        let resourceInputs: pulumi.Inputs = {};
        opts = opts || {};
        if (!opts.id) {
            if ((!args || args.chart === undefined) && !opts.urn) {
                throw new Error("Missing required property 'chart'");
            }
            resourceInputs["chart"] = args ? args.chart : undefined;
            resourceInputs["dependencyUpdate"] = args ? args.dependencyUpdate : undefined;
            resourceInputs["devel"] = args ? args.devel : undefined;
            resourceInputs["keyring"] = args ? args.keyring : undefined;
            resourceInputs["name"] = args ? args.name : undefined;
            resourceInputs["namespace"] = args ? args.namespace : undefined;
            resourceInputs["plainHttp"] = args ? args.plainHttp : undefined;
            resourceInputs["postRenderer"] = args ? args.postRenderer : undefined;
            resourceInputs["repositoryOpts"] = args ? args.repositoryOpts : undefined;
            resourceInputs["resourcePrefix"] = args ? args.resourcePrefix : undefined;
            resourceInputs["skipAwait"] = args ? args.skipAwait : undefined;
            resourceInputs["skipCrds"] = args ? args.skipCrds : undefined;
            resourceInputs["valueYamlFiles"] = args ? args.valueYamlFiles : undefined;
            resourceInputs["values"] = args ? args.values : undefined;
            resourceInputs["verify"] = args ? args.verify : undefined;
            resourceInputs["version"] = args ? args.version : undefined;
            resourceInputs["resources"] = undefined /*out*/;
        } else {
            resourceInputs["resources"] = undefined /*out*/;
        }
        opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts);
        super(Chart.__pulumiType, name, resourceInputs, opts, true /*remote*/);
    }
}

/**
 * The set of arguments for constructing a Chart resource.
 */
export interface ChartArgs {
    /**
     * Chart name to be installed. A path may be used.
     */
    chart: pulumi.Input<string>;
    /**
     * Run helm dependency update before installing the chart.
     */
    dependencyUpdate?: pulumi.Input<boolean>;
    /**
     * Use chart development versions, too. Equivalent to version '>0.0.0-0'. If `version` is set, this is ignored.
     */
    devel?: pulumi.Input<boolean>;
    /**
     * Location of public keys used for verification. Used only if `verify` is true
     */
    keyring?: pulumi.Input<pulumi.asset.Asset | pulumi.asset.Archive>;
    /**
     * Release name.
     */
    name?: pulumi.Input<string>;
    /**
     * Namespace for the release.
     */
    namespace?: pulumi.Input<string>;
    /**
     * Use insecure HTTP for the chart download instead of HTTPS.
     */
    plainHttp?: pulumi.Input<boolean>;
    /**
     * Specification defining the post-renderer to use.
     */
    postRenderer?: pulumi.Input<inputs.helm.v4.PostRenderer>;
    /**
     * Specification defining the Helm chart repository to use.
     */
    repositoryOpts?: pulumi.Input<inputs.helm.v4.RepositoryOpts>;
    /**
     * An optional prefix for the auto-generated resource names. Example: A resource created with resourcePrefix="foo" would produce a resource named "foo:resourceName".
     */
    resourcePrefix?: pulumi.Input<string>;
    /**
     * By default, the provider waits until all resources are in a ready state before marking the release as successful. Setting this to true will skip such await logic.
     */
    skipAwait?: pulumi.Input<boolean>;
    /**
     * If set, no CRDs will be installed. By default, CRDs are installed if not already present.
     */
    skipCrds?: pulumi.Input<boolean>;
    /**
     * List of assets (raw yaml files). Content is read and merged with values.
     */
    valueYamlFiles?: pulumi.Input<pulumi.Input<pulumi.asset.Asset | pulumi.asset.Archive>[]>;
    /**
     * Custom values set for the release.
     */
    values?: pulumi.Input<{[key: string]: any}>;
    /**
     * Verify the chart's integrity.
     */
    verify?: pulumi.Input<boolean>;
    /**
     * Specify the chart version to install. If this is not specified, the latest version is installed.
     */
    version?: pulumi.Input<string>;
}
