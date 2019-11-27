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
        /// Note: if .metadata.namespace is set on a resource, that value takes precedence over the provider default.
        /// </summary>
        [Input("namespace")]
        public Input<string>? Namespace { get; set; }

        /// <summary>
        /// BETA FEATURE - If present and set to true, enable server-side diff calculations.
        /// This feature is in developer preview, and is disabled by default.
        /// </summary>
        [Input("enableDryRun")]
        public Input<bool>? EnableDryRun { get; set; }

        /// <summary>
        /// If present and set to true, suppress apiVersion deprecation warnings from the CLI.
        /// </summary>
        [Input("suppressDeprecationWarnings")]
        public Input<bool>? SuppressDeprecationWarnings { get; set; }
    }
}
