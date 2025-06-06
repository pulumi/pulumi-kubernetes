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
    /// DeviceClassSpec is used in a [DeviceClass] to define what can be allocated and how to configure it.
    /// </summary>
    [OutputType]
    public sealed class DeviceClassSpecPatch
    {
        /// <summary>
        /// Config defines configuration parameters that apply to each device that is claimed via this class. Some classses may potentially be satisfied by multiple drivers, so each instance of a vendor configuration applies to exactly one driver.
        /// 
        /// They are passed to the driver, but are not considered while allocating the claim.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Resource.V1Beta2.DeviceClassConfigurationPatch> Config;
        /// <summary>
        /// Each selector must be satisfied by a device which is claimed via this class.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Resource.V1Beta2.DeviceSelectorPatch> Selectors;

        [OutputConstructor]
        private DeviceClassSpecPatch(
            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Resource.V1Beta2.DeviceClassConfigurationPatch> config,

            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Resource.V1Beta2.DeviceSelectorPatch> selectors)
        {
            Config = config;
            Selectors = selectors;
        }
    }
}
