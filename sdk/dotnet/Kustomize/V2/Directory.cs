// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Kustomize.V2
{
    /// <summary>
    /// Directory is a component representing a collection of resources described by a kustomize directory (kustomization).
    /// 
    /// ## Example Usage
    /// ### Local Kustomize Directory
    /// ```csharp
    /// using System.Threading.Tasks;
    /// using Pulumi;
    /// using Pulumi.Kubernetes.Kustomize.V2;
    /// 
    /// class KustomizeStack : Stack
    /// {
    ///     public KustomizeStack()
    ///     {
    ///         var helloWorld = new Directory("helloWorldLocal", new DirectoryArgs
    ///         {
    ///             Directory = "./helloWorld",
    ///         });
    ///     }
    /// }
    /// ```
    /// ### Kustomize Directory from a Git Repo
    /// ```csharp
    /// using System.Threading.Tasks;
    /// using Pulumi;
    /// using Pulumi.Kubernetes.Kustomize.V2;
    /// 
    /// class KustomizeStack : Stack
    /// {
    ///     public KustomizeStack()
    ///     {
    ///         var helloWorld = new Directory("helloWorldRemote", new DirectoryArgs
    ///         {
    ///             Directory = "https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld",
    ///         });
    ///     }
    /// }
    /// ```
    /// </summary>
    [KubernetesResourceType("kubernetes:kustomize/v2:Directory")]
    public partial class Directory : global::Pulumi.ComponentResource
    {
        /// <summary>
        /// Resources created by the Directory resource.
        /// </summary>
        [Output("resources")]
        public Output<string> Resources { get; private set; } = null!;


        /// <summary>
        /// Create a Directory resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public Directory(string name, Pulumi.Kubernetes.Types.Inputs.Kustomize.V2.DirectoryArgs? args = null, ComponentResourceOptions? options = null)
            : base("kubernetes:kustomize/v2:Directory", name, args ?? new Pulumi.Kubernetes.Types.Inputs.Kustomize.V2.DirectoryArgs(), MakeResourceOptions(options, ""), remote: true)
        {
        }

        private static ComponentResourceOptions MakeResourceOptions(ComponentResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new ComponentResourceOptions
            {
                Version = Utilities.Version,
            };
            var merged = ComponentResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
    }
}
namespace Pulumi.Kubernetes.Types.Inputs.Kustomize.V2
{

    public class DirectoryArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// The directory containing the kustomization to apply. The value can be a local directory or a folder in a
        /// git repository.
        /// Example: ./helloWorld
        /// Example: https://github.com/kubernetes-sigs/kustomize/tree/master/examples/helloWorld
        /// </summary>
        [Input("directory", required: true)]
        public Input<string> Directory { get; set; } = null!;

        /// <summary>
        /// The default namespace to apply to the resources. Defaults to the provider's namespace.
        /// </summary>
        [Input("namespace")]
        public Input<string>? Namespace { get; set; }

        /// <summary>
        /// A prefix for the auto-generated resource names. Defaults to the name of the Directory resource. Example: A resource created with resourcePrefix="foo" would produce a resource named "foo:resourceName".
        /// </summary>
        [Input("resourcePrefix")]
        public Input<string>? ResourcePrefix { get; set; }

        /// <summary>
        /// Indicates that child resources should skip the await logic.
        /// </summary>
        [Input("skipAwait")]
        public Input<bool>? SkipAwait { get; set; }

        public DirectoryArgs()
        {
        }
        public static new DirectoryArgs Empty => new DirectoryArgs();
    }
}
