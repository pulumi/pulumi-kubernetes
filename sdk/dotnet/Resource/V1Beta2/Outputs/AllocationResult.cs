// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Resource.V1Beta2
{

    /// <summary>
    /// AllocationResult contains attributes of an allocated resource.
    /// </summary>
    [OutputType]
    public sealed class AllocationResult
    {
        /// <summary>
        /// Devices is the result of allocating devices.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Resource.V1Beta2.DeviceAllocationResult Devices;
        /// <summary>
        /// NodeSelector defines where the allocated resources are available. If unset, they are available everywhere.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeSelector NodeSelector;

        [OutputConstructor]
        private AllocationResult(
            Pulumi.Kubernetes.Types.Outputs.Resource.V1Beta2.DeviceAllocationResult devices,

            Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeSelector nodeSelector)
        {
            Devices = devices;
            NodeSelector = nodeSelector;
        }
    }
}
