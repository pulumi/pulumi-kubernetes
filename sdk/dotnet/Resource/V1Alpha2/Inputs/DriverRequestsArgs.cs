// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Inputs.Resource.V1Alpha2
{

    /// <summary>
    /// DriverRequests describes all resources that are needed from one particular driver.
    /// </summary>
    public class DriverRequestsArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// DriverName is the name used by the DRA driver kubelet plugin.
        /// </summary>
        [Input("driverName")]
        public Input<string>? DriverName { get; set; }

        [Input("requests")]
        private InputList<Pulumi.Kubernetes.Types.Inputs.Resource.V1Alpha2.ResourceRequestArgs>? _requests;

        /// <summary>
        /// Requests describes all resources that are needed from the driver.
        /// </summary>
        public InputList<Pulumi.Kubernetes.Types.Inputs.Resource.V1Alpha2.ResourceRequestArgs> Requests
        {
            get => _requests ?? (_requests = new InputList<Pulumi.Kubernetes.Types.Inputs.Resource.V1Alpha2.ResourceRequestArgs>());
            set => _requests = value;
        }

        /// <summary>
        /// VendorParameters are arbitrary setup parameters for all requests of the claim. They are ignored while allocating the claim.
        /// </summary>
        [Input("vendorParameters")]
        public InputJson? VendorParameters { get; set; }

        public DriverRequestsArgs()
        {
        }
        public static new DriverRequestsArgs Empty => new DriverRequestsArgs();
    }
}
