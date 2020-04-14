// Copyright 2016-2020, Pulumi Corporation

using Pulumi.Kubernetes.Types.Inputs.Meta.V1;
using Pulumi.Kubernetes.Types.Outputs.Meta.V1;

namespace Pulumi.Kubernetes.ApiExtensions
{
    /// <summary>
    /// Represents an instance of a <see cref="V1.CustomResourceDefinition"/> (CRD). For example, the
    /// CoreOS Prometheus operator exposes a CRD `monitoring.coreos.com/ServiceMonitor`; to
    /// instantiate this as a Pulumi resource, one could call `new CustomResource`, passing the
    /// `ServiceMonitor` resource definition as an argument.
    /// </summary>
    public class CustomResource : KubernetesResource
    {
        /// <summary>
        /// APIVersion defines the versioned schema of this representation of an object. Servers
        /// should convert recognized schemas to the latest internal value, and may reject
        /// unrecognized values. More info:
        /// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        /// </summary>
        [Output("apiVersion")]
        public Output<string> ApiVersion { get; private set; } = null!;

        /// <summary>
        /// Kind is a string value representing the REST resource this object represents. Servers
        /// may infer this from the endpoint the client submits requests to. Cannot be updated. In
        /// CamelCase. More info:
        /// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        /// </summary>
        [Output("kind")]
        public Output<string> Kind { get; private set; } = null!;
        
        /// <summary>
        /// Standard object metadata; More info:
        /// https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
        /// </summary>
        [Output("metadata")]
        public Output<ObjectMeta> Metadata { get; private set; } = null!;

        public CustomResource(string name, CustomResourceArgs args, CustomResourceOptions? options = null)
            : base(args.Type, name, args, options)
        {
        }
    }
    
    /// <summary>
    /// Represents a resource definition we'd use to create an instance of a Kubernetes
    /// <see cref="V1.CustomResourceDefinition"/> (CRD). For example, the CoreOS Prometheus operator
    /// exposes a CRD `monitoring.coreos.com/ServiceMonitor`; to create a `ServiceMonitor`, we'd
    /// pass a <see cref="CustomResourceArgs"/> containing the `ServiceMonitor` definition to
    /// <see cref="CustomResource"/>.
    ///
    /// NOTE: This type is abstract. You need to inherit from it and define all the properties of
    /// a specific custom resource in the derived class.
    /// </summary>
    public abstract class CustomResourceArgs : ResourceArgs
    {
        /// <summary>
        /// APIVersion defines the versioned schema of this representation of an object. Servers should
        /// convert recognized schemas to the latest internal value, and may reject unrecognized
        /// values. More info:
        /// https://git.k8s.io/community/contributors/devel/api-conventions.md#resources
        /// </summary>
        [Input("apiVersion")]
        public Input<string> ApiVersion { get; }

        /// <summary>
        /// Kind is a string value representing the REST resource this object represents. Servers may
        /// infer this from the endpoint the client submits requests to. Cannot be updated. In
        /// CamelCase. More info:
        /// https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
        /// </summary>
        [Input("kind")]
        public Input<string> Kind { get; }

        /// <summary>
        /// Standard object's metadata. More info:
        /// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        /// </summary>
        [Input("metadata")]
        public Input<ObjectMetaArgs>? Metadata { get; set; }

        /// <summary>
        /// Resource type name, e.g. `kubernetes:stable.example.com:CronTab`.
        /// We can't extract it from ApiVersion and Kind because they are Inputs, while we need
        /// a plain string in the CustomResource constructor.
        /// </summary>
        internal string Type { get; }

        protected CustomResourceArgs(string apiVersion, string kind)
        {
            this.ApiVersion = apiVersion;
            this.Kind = kind;
            this.Type = $"kubernetes:{apiVersion}:{kind}";
        }
    }
}
