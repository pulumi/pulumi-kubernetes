// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Core.V1
{

    /// <summary>
    /// PodResourceClaimStatus is stored in the PodStatus for each PodResourceClaim which references a ResourceClaimTemplate. It stores the generated name for the corresponding ResourceClaim.
    /// </summary>
    [OutputType]
    public sealed class PodResourceClaimStatus
    {
        /// <summary>
        /// Name uniquely identifies this resource claim inside the pod. This must match the name of an entry in pod.spec.resourceClaims, which implies that the string must be a DNS_LABEL.
        /// </summary>
        public readonly string Name;
        /// <summary>
        /// ResourceClaimName is the name of the ResourceClaim that was generated for the Pod in the namespace of the Pod. It this is unset, then generating a ResourceClaim was not necessary. The pod.spec.resourceClaims entry can be ignored in this case.
        /// </summary>
        public readonly string ResourceClaimName;

        [OutputConstructor]
        private PodResourceClaimStatus(
            string name,

            string resourceClaimName)
        {
            Name = name;
            ResourceClaimName = resourceClaimName;
        }
    }
}