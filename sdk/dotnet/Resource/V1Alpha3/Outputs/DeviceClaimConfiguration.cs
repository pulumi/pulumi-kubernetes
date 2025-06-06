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
    /// DeviceClaimConfiguration is used for configuration parameters in DeviceClaim.
    /// </summary>
    [OutputType]
    public sealed class DeviceClaimConfiguration
    {
        /// <summary>
        /// Opaque provides driver-specific configuration parameters.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Resource.V1Alpha3.OpaqueDeviceConfiguration Opaque;
        /// <summary>
        /// Requests lists the names of requests where the configuration applies. If empty, it applies to all requests.
        /// 
        /// References to subrequests must include the name of the main request and may include the subrequest using the format &lt;main request&gt;[/&lt;subrequest&gt;]. If just the main request is given, the configuration applies to all subrequests.
        /// </summary>
        public readonly ImmutableArray<string> Requests;

        [OutputConstructor]
        private DeviceClaimConfiguration(
            Pulumi.Kubernetes.Types.Outputs.Resource.V1Alpha3.OpaqueDeviceConfiguration opaque,

            ImmutableArray<string> requests)
        {
            Opaque = opaque;
            Requests = requests;
        }
    }
}
