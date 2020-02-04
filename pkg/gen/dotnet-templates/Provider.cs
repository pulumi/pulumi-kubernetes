using Pulumi.Serialization;

namespace Pulumi.Kubernetes
{
    /// <summary>
    /// The provider type for the kubernetes package.
    /// </summary>
    public class Provider : ProviderResource
    {
        /// <summary>
        /// Create a Provider resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource.</param>
        /// <param name="args">The arguments used to populate this resource's properties.</param>
        /// <param name="options">A bag of options that control this resource's behavior.</param>
        public Provider(string name, ProviderArgs? args = null, ResourceOptions? options = null)
            : base("kubernetes", name, args ?? ResourceArgs.Empty, MakeResourceOptions(options, ""))
        {
        }

        private static ResourceOptions MakeResourceOptions(ResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new ResourceOptions
            {
                Version = Utilities.Version,
            };
            var merged = ResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
    }

    /// <summary>
    /// The set of arguments for constructing a Provider.
    /// </summary>
    public sealed class ProviderArgs : ResourceArgs
    {
        /// <summary>
        /// If present, the name of the kubeconfig cluster to use.
        /// </summary>
        [Input("cluster")]
        public Input<string>? Cluster { get; set; }

        /// <summary>
        /// If present, the name of the kubeconfig context to use.
        /// </summary>
        [Input("context")]
        public Input<string>? Context { get; set; }

        /// <summary>
        /// The contents of a kubeconfig file. If this is set, this config will be used instead of $KUBECONFIG.
        /// </summary>
        [Input("kubeconfig")]
        public Input<string>? KubeConfig { get; set; }

        /// <summary>
        /// If present, the default namespace to use. This flag is ignored for cluster-scoped resources.
        ///
        /// A namespace can be specified in multiple places, and the precedence is as follows:
        /// 1. `.metadata.namespace` set on the resource.
        /// 2. This `namespace` parameter.
        /// 3. `namespace` set for the active context in the kubeconfig.
        /// </summary>
        [Input("namespace")]
        public Input<string>? Namespace { get; set; }

        /// <summary>
        /// BETA FEATURE - If present and set to true, enable server-side diff calculations.
        /// This feature is in developer preview, and is disabled by default.
        ///
        /// This config can be specified in the following ways, using this precedence:
        /// 1. This `enableDryRun` parameter.
        /// 2. The `PULUMI_K8S_ENABLE_DRY_RUN` environment variable.
        /// </summary>
        [Input("enableDryRun")]
        public Input<bool>? EnableDryRun { get; set; }

        /// <summary>
        /// If present and set to true, suppress apiVersion deprecation warnings from the CLI.
        ///
        /// This config can be specified in the following ways, using this precedence:
        /// 1. This `suppressDeprecationWarnings` parameter.
        /// 2. The `PULUMI_K8S_SUPPRESS_DEPRECATION_WARNINGS` environment variable.
        /// </summary>
        [Input("suppressDeprecationWarnings")]
        public Input<bool>? SuppressDeprecationWarnings { get; set; }

        /// <summary>
        /// If present, render resource manifests to this directory. In this mode, resources will not
        /// be created on a Kubernetes cluster, but the rendered manifests will be kept in sync with changes
        /// to the Pulumi program. Note that some computed Outputs such as status fields will not be populated
        /// since the resources are not created on a Kubernetes cluster. Attempting to reference these Outputs
        /// may result in an error, or the value may be empty/undefined.
        /// </summary>
        [Input("renderYamlToDirectory")]
        public Input<string>? RenderYamlToDirectory { get; set; }
    }
}
