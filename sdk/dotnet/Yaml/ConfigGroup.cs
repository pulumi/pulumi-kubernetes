// Copyright 2016-2020, Pulumi Corporation

using System.Collections.Immutable;
using TransformationAction = System.Func<System.Collections.Immutable.ImmutableDictionary<string, object>, Pulumi.CustomResourceOptions, System.Collections.Immutable.ImmutableDictionary<string, object>>;

namespace Pulumi.Kubernetes.Yaml
{
    /// <summary>
    /// Creates a set of Kubernetes resources from Kubernetes YAML text. The YAML text
    /// may be supplied using any of the following <see cref="ConfigGroupArgs"/>:
    /// 1. Using a list of filenames: `Files = new[] { "foo.yaml", "bar.yaml" }`
    /// 2. Using a list of file patterns: `Files = new[] { "foo/*.yaml", "bar/*.yaml" }`
    /// 3. Using literal strings containing YAML: `Yaml = new[] { "(LITERAL YAML HERE)", "(MORE YAML)" }`
    /// 4. Any combination of files, patterns, or YAML strings.
    /// </summary>
    public sealed class ConfigGroup : CollectionComponentResource
    {
        public ConfigGroup(string name, ConfigGroupArgs config, ComponentResourceOptions? options = null)
            : base("kubernetes:yaml:ConfigGroup", name, options)
        {
            options ??= new ComponentResourceOptions();
            options.Parent ??= this;
            RegisterResources(Parser.Parse(config, options));
        }
    }
    
    /// <summary>
    /// Resource arguments for <see cref="ConfigGroup"/>.
    /// </summary>
    public class ConfigGroupArgs : ResourceArgs
    {
        /// <summary>
        /// Set of paths or a URLs that uniquely identify files.
        /// </summary>
        public string[]? Files { get; set; }
        
        private InputList<string>? _yaml;

        /// <summary>
        /// YAML text containing Kubernetes resource definitions.
        /// </summary>
        public InputList<string> Yaml
        {
            get => _yaml ??= new InputList<string>();
            set => _yaml = value;
        }

        private InputList<ImmutableDictionary<string, object>>? _objs;

        /// <summary>
        /// Objects representing Kubernetes resources.
        /// </summary>
        public InputList<ImmutableDictionary<string, object>> Objs
        {
            get => _objs ??= new InputList<ImmutableDictionary<string, object>>();
            set => _objs = value;
        }

        /// <summary>
        /// A set of transformations to apply to Kubernetes resource definitions before registering
        /// with engine.
        /// </summary>
        public TransformationAction[]? Transformations { get; set; }

        /// <summary>
        /// An optional prefix for the auto-generated resource names.
        /// Example: A resource created with ResourcePrefix="foo" would produce a resource named
        /// "foo-resourceName".
        /// </summary>
        public string? ResourcePrefix { get; set; }
    }
}
