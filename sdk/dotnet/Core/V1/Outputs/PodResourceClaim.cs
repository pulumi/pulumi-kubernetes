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
    /// PodResourceClaim references exactly one ResourceClaim through a ClaimSource. It adds a name to it that uniquely identifies the ResourceClaim inside the Pod. Containers that need access to the ResourceClaim reference it with this name.
    /// </summary>
    [OutputType]
    public sealed class PodResourceClaim
    {
        /// <summary>
        /// Name uniquely identifies this resource claim inside the pod. This must be a DNS_LABEL.
        /// </summary>
        public readonly string Name;
        /// <summary>
        /// Source describes where to find the ResourceClaim.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Core.V1.ClaimSource Source;

        [OutputConstructor]
        private PodResourceClaim(
            string name,

            Pulumi.Kubernetes.Types.Outputs.Core.V1.ClaimSource source)
        {
            Name = name;
            Source = source;
        }
    }
}
