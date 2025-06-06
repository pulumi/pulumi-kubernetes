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
    /// OpaqueDeviceConfiguration contains configuration parameters for a driver in a format defined by the driver vendor.
    /// </summary>
    [OutputType]
    public sealed class OpaqueDeviceConfiguration
    {
        /// <summary>
        /// Driver is used to determine which kubelet plugin needs to be passed these configuration parameters.
        /// 
        /// An admission policy provided by the driver developer could use this to decide whether it needs to validate them.
        /// 
        /// Must be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver.
        /// </summary>
        public readonly string Driver;
        /// <summary>
        /// Parameters can contain arbitrary data. It is the responsibility of the driver developer to handle validation and versioning. Typically this includes self-identification and a version ("kind" + "apiVersion" for Kubernetes types), with conversion between different versions.
        /// 
        /// The length of the raw data must be smaller or equal to 10 Ki.
        /// </summary>
        public readonly System.Text.Json.JsonElement Parameters;

        [OutputConstructor]
        private OpaqueDeviceConfiguration(
            string driver,

            System.Text.Json.JsonElement parameters)
        {
            Driver = driver;
            Parameters = parameters;
        }
    }
}
