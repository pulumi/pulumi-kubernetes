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
    /// Device represents one individual hardware instance that can be selected based on its attributes. Besides the name, exactly one field must be set.
    /// </summary>
    public class DevicePatchArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// Basic defines one device instance.
        /// </summary>
        [Input("basic")]
        public Input<Pulumi.Kubernetes.Types.Inputs.Resource.V1Beta1.BasicDevicePatchArgs>? Basic { get; set; }

        /// <summary>
        /// Name is unique identifier among all devices managed by the driver in the pool. It must be a DNS label.
        /// </summary>
        [Input("name")]
        public Input<string>? Name { get; set; }

        public DevicePatchArgs()
        {
        }
        public static new DevicePatchArgs Empty => new DevicePatchArgs();
    }
}