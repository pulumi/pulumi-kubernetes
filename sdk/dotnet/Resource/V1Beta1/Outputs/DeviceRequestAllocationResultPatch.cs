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
    /// DeviceRequestAllocationResult contains the allocation result for one request.
    /// </summary>
    [OutputType]
    public sealed class DeviceRequestAllocationResultPatch
    {
        /// <summary>
        /// AdminAccess indicates that this device was allocated for administrative access. See the corresponding request field for a definition of mode.
        /// 
        /// This is an alpha field and requires enabling the DRAAdminAccess feature gate. Admin access is disabled if this field is unset or set to false, otherwise it is enabled.
        /// </summary>
        public readonly bool AdminAccess;
        /// <summary>
        /// Device references one device instance via its name in the driver's resource pool. It must be a DNS label.
        /// </summary>
        public readonly string Device;
        /// <summary>
        /// Driver specifies the name of the DRA driver whose kubelet plugin should be invoked to process the allocation once the claim is needed on a node.
        /// 
        /// Must be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver.
        /// </summary>
        public readonly string Driver;
        /// <summary>
        /// This name together with the driver name and the device name field identify which device was allocated (`&lt;driver name&gt;/&lt;pool name&gt;/&lt;device name&gt;`).
        /// 
        /// Must not be longer than 253 characters and may contain one or more DNS sub-domains separated by slashes.
        /// </summary>
        public readonly string Pool;
        /// <summary>
        /// Request is the name of the request in the claim which caused this device to be allocated. Multiple devices may have been allocated per request.
        /// </summary>
        public readonly string Request;

        [OutputConstructor]
        private DeviceRequestAllocationResultPatch(
            bool adminAccess,

            string device,

            string driver,

            string pool,

            string request)
        {
            AdminAccess = adminAccess;
            Device = device;
            Driver = driver;
            Pool = pool;
            Request = request;
        }
    }
}