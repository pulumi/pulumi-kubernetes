// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Resource.V1Alpha3
{

    /// <summary>
    /// PodSchedulingContextStatus describes where resources for the Pod can be allocated.
    /// </summary>
    [OutputType]
    public sealed class PodSchedulingContextStatus
    {
        /// <summary>
        /// ResourceClaims describes resource availability for each pod.spec.resourceClaim entry where the corresponding ResourceClaim uses "WaitForFirstConsumer" allocation mode.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Resource.V1Alpha3.ResourceClaimSchedulingStatus> ResourceClaims;

        [OutputConstructor]
        private PodSchedulingContextStatus(ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Resource.V1Alpha3.ResourceClaimSchedulingStatus> resourceClaims)
        {
            ResourceClaims = resourceClaims;
        }
    }
}
