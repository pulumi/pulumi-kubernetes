// Copyright 2016-2020, Pulumi Corporation

using System;
using System.IO;
using TransformationAction = System.Func<System.Collections.Immutable.ImmutableDictionary<string, object>, Pulumi.CustomResourceOptions, System.Collections.Immutable.ImmutableDictionary<string, object>>;

namespace Pulumi.Kubernetes.Yaml
{
    /// <summary>
    /// Defines a set of Kubernetes resources from a Kubernetes YAML file. 
    /// </summary>
    public sealed class ConfigFile : CollectionComponentResource
    {
        /// <summary>
        /// Defines a set of Kubernetes resources from a Kubernetes YAML file.
        /// </summary>
        /// <param name="name">Component name. If `args` is not specified, also treated as the file name.
        /// </param>
        /// <param name="args">Resource arguments, including the YAML file name.</param>
        /// <param name="options">Resource options.</param>
        public ConfigFile(string name, ConfigFileArgs? args = null, ComponentResourceOptions? options = null)
            : base("kubernetes:yaml:ConfigFile", MakeName(args, name), options)
        {
            name = MakeName(args, name);
            options ??= new ComponentResourceOptions();
            options.Parent ??= this;
            
            var fileOutput = args?.File.ToOutput() ?? Output.Create(name);
            var resources = fileOutput.Apply(fileId =>
            {
                try
                {
                    if (Parser.IsUrl(fileId))
                    {
                        using var wc = new System.Net.WebClient();
                        return wc.DownloadString(fileId);
                    }

                    return File.ReadAllText(fileId);
                }
                catch (Exception e)
                {
                    throw new ResourceException($"Error fetching YAML file '{fileId}': {e.Message}", this);
                }
            }).Apply(text =>
                Parser.ParseYamlDocument(new ParseArgs
                {
                    Objs = Invokes.YamlDecode(new YamlDecodeArgs {Text = text}),
                    Transformations = args?.Transformations,
                    ResourcePrefix = args?.ResourcePrefix
                }, options));

            RegisterResources(resources);
        }

        private static string MakeName(ConfigFileArgs? args, string name)
            => args?.ResourcePrefix != null ? $"{args.ResourcePrefix}-{name}" : name;
    }
    
    /// <summary>
    /// Resource arguments for <see cref="ConfigFile"/>.
    /// </summary>
    public class ConfigFileArgs : ResourceArgs
    {
        /// <summary>
        /// Path or a URL that uniquely identifies a file.
        /// </summary>
        public Input<string>? File { get; set; }

        /// <summary>
        /// A set of transformations to apply to Kubernetes resource definitions before registering
        /// with engine.
        /// </summary>
        public TransformationAction[]? Transformations { get; set; }

        /// <summary>
        /// An optional prefix for the auto-generated resource names.
        /// Example: A resource created with resourcePrefix="foo" would produce a resource named "foo-resourceName".
        /// </summary>
        public string? ResourcePrefix { get; set; }
    }
}
