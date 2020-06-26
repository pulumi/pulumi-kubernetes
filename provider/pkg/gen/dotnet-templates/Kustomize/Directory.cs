// Copyright 2016-2020, Pulumi Corporation

using System.Collections.Generic;
using Pulumi.Kubernetes.Yaml;

namespace Pulumi.Kubernetes.Kustomize
{
    /// <summary>
    /// Directory is a component representing a collection of resources described by a kustomize directory (kustomization).
    /// </summary>
    public sealed class Directory : CollectionComponentResource
    {
        /// <summary>
        /// Directory is a component representing a collection of resources described by a kustomize directory (kustomization).
        /// </summary>
        /// <param name="name">Name of the kustomization (e.g., nginx-ingress).</param>
        /// <param name="args">Configuration options for the kustomization.</param>
        /// <param name="options">Resource options.</param>
        public Directory(string name, DirectoryArgs args, ComponentResourceOptions? options = null)
            : base("kubernetes:kustomize:Directory", MakeName(args, name), options)
        {
            name = GetName(args, name);
            var objs = Invokes.KustomizeDirectory(new KustomizeDirectoryArgs { Directory = args.Directory });
            var configGroupArgs = new ConfigGroupArgs
            {
                ResourcePrefix = args.ResourcePrefix,
                Objs = objs,
                Transformations = args.Transformations
            };
            var opts = ComponentResourceOptions.Merge(options, new ComponentResourceOptions { Parent = this });
            var resources = Parser.Parse(configGroupArgs, opts);
            RegisterResources(resources);
        }
        private static string MakeName(DirectoryArgs? args, string name)
            => args?.ResourcePrefix != null ? $"{args.ResourcePrefix}-{name}" : name;

        private static string GetName(DirectoryArgs config, string releaseName)
        {
            var prefix = config.ResourcePrefix;
            return string.IsNullOrEmpty(prefix) ? releaseName : $"{prefix}-{releaseName}";
        }

    }

    /// <summary>
    /// Resource arguments for <see cref="Directory"/>.
    /// </summary>
    public class DirectoryArgs : ResourceArgs
    {
        /// <summary>
        /// The directory containing the kustomization to apply. The value can be a local directory or a folder in a
        /// git repository.
        /// Example: ./helloWorld
        /// Example: https://github.com/kubernetes-sigs/kustomize/tree/master/examples/helloWorld
        /// </summary>
        public string? Directory { get; set; }

        private List<TransformationAction>? _transformations;

        /// <summary>
        /// Optional array of transformations to apply to resources that will be created by this chart prior to
        /// creation. Allows customization of the chart behaviour without directly modifying the chart itself.
        /// </summary>
        public List<TransformationAction> Transformations
        {
            get => _transformations ??= new List<TransformationAction>();
            set => _transformations = value;
        }

        /// <summary>
        /// An optional prefix for the auto-generated resource names.
        /// Example: A resource created with resourcePrefix="foo" would produce a resource named "foo-resourceName".
        /// </summary>
        public string? ResourcePrefix { get; set; }
    }
}
