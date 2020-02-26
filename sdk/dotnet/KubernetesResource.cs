// Copyright 2016-2020, Pulumi Corporation

using System.Collections.Immutable;

namespace Pulumi.Kubernetes
{
    /// <summary>
    /// A base class for all Kubernetes resources.
    /// </summary>
    public abstract class KubernetesResource : CustomResource
    {
        /// <summary>
        /// Standard constructor passing arguments to <see cref="CustomResource"/>.
        /// </summary>
        internal KubernetesResource(string type, string name, ResourceArgs? args, CustomResourceOptions? options = null)
            : base(type, name, args, MakeResourceOptions(options))
        {
        }

        /// <summary>
        /// Additional constructor for dynamic arguments received from YAML-based sources.
        /// </summary>
        internal KubernetesResource(string type, string name, ImmutableDictionary<string, object?> dictionary,
            CustomResourceOptions? options = null)
            : base(type, name, new DictionaryResourceArgs(dictionary), MakeResourceOptions(options))
        {
        }

        private static CustomResourceOptions MakeResourceOptions(CustomResourceOptions? options)
        {
            var defaultOptions = new CustomResourceOptions
            {
                Version = Utilities.Version,
            };
            return CustomResourceOptions.Merge(defaultOptions, options);
        }
    }
}
