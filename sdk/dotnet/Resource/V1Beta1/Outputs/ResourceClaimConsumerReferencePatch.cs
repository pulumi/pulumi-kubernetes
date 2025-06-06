// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Resource.V1Beta1
{

    /// <summary>
    /// ResourceClaimConsumerReference contains enough information to let you locate the consumer of a ResourceClaim. The user must be a resource in the same namespace as the ResourceClaim.
    /// </summary>
    [OutputType]
    public sealed class ResourceClaimConsumerReferencePatch
    {
        /// <summary>
        /// APIGroup is the group for the resource being referenced. It is empty for the core API. This matches the group in the APIVersion that is used when creating the resources.
        /// </summary>
        public readonly string ApiGroup;
        /// <summary>
        /// Name is the name of resource being referenced.
        /// </summary>
        public readonly string Name;
        /// <summary>
        /// Resource is the type of resource being referenced, for example "pods".
        /// </summary>
        public readonly string Resource;
        /// <summary>
        /// UID identifies exactly one incarnation of the resource.
        /// </summary>
        public readonly string Uid;

        [OutputConstructor]
        private ResourceClaimConsumerReferencePatch(
            string apiGroup,

            string name,

            string resource,

            string uid)
        {
            ApiGroup = apiGroup;
            Name = name;
            Resource = resource;
            Uid = uid;
        }
    }
}
