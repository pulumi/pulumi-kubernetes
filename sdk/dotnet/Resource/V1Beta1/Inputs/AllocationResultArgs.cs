// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Inputs.Resource.V1Beta1
{

    /// <summary>
    /// AllocationResult contains attributes of an allocated resource.
    /// </summary>
    public class AllocationResultArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// Devices is the result of allocating devices.
        /// </summary>
        [Input("devices")]
        public Input<Pulumi.Kubernetes.Types.Inputs.Resource.V1Beta1.DeviceAllocationResultArgs>? Devices { get; set; }

        /// <summary>
        /// NodeSelector defines where the allocated resources are available. If unset, they are available everywhere.
        /// </summary>
        [Input("nodeSelector")]
        public Input<Pulumi.Kubernetes.Types.Inputs.Core.V1.NodeSelectorArgs>? NodeSelector { get; set; }

        public AllocationResultArgs()
        {
        }
        public static new AllocationResultArgs Empty => new AllocationResultArgs();
    }
}