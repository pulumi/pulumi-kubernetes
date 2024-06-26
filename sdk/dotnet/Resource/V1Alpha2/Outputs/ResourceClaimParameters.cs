// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Resource.V1Alpha2
{

    /// <summary>
    /// ResourceClaimParameters defines resource requests for a ResourceClaim in an in-tree format understood by Kubernetes.
    /// </summary>
    [OutputType]
    public sealed class ResourceClaimParameters
    {
        /// <summary>
        /// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        /// </summary>
        public readonly string ApiVersion;
        /// <summary>
        /// DriverRequests describes all resources that are needed for the allocated claim. A single claim may use resources coming from different drivers. For each driver, this array has at most one entry which then may have one or more per-driver requests.
        /// 
        /// May be empty, in which case the claim can always be allocated.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Resource.V1Alpha2.DriverRequests> DriverRequests;
        /// <summary>
        /// If this object was created from some other resource, then this links back to that resource. This field is used to find the in-tree representation of the claim parameters when the parameter reference of the claim refers to some unknown type.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Resource.V1Alpha2.ResourceClaimParametersReference GeneratedFrom;
        /// <summary>
        /// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        /// </summary>
        public readonly string Kind;
        /// <summary>
        /// Standard object metadata
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Meta.V1.ObjectMeta Metadata;
        /// <summary>
        /// Shareable indicates whether the allocated claim is meant to be shareable by multiple consumers at the same time.
        /// </summary>
        public readonly bool Shareable;

        [OutputConstructor]
        private ResourceClaimParameters(
            string apiVersion,

            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Resource.V1Alpha2.DriverRequests> driverRequests,

            Pulumi.Kubernetes.Types.Outputs.Resource.V1Alpha2.ResourceClaimParametersReference generatedFrom,

            string kind,

            Pulumi.Kubernetes.Types.Outputs.Meta.V1.ObjectMeta metadata,

            bool shareable)
        {
            ApiVersion = apiVersion;
            DriverRequests = driverRequests;
            GeneratedFrom = generatedFrom;
            Kind = kind;
            Metadata = metadata;
            Shareable = shareable;
        }
    }
}
