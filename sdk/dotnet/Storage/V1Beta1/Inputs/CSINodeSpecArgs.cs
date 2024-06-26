// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Inputs.Storage.V1Beta1
{

    /// <summary>
    /// CSINodeSpec holds information about the specification of all CSI drivers installed on a node
    /// </summary>
    public class CSINodeSpecArgs : global::Pulumi.ResourceArgs
    {
        [Input("drivers", required: true)]
        private InputList<Pulumi.Kubernetes.Types.Inputs.Storage.V1Beta1.CSINodeDriverArgs>? _drivers;

        /// <summary>
        /// drivers is a list of information of all CSI Drivers existing on a node. If all drivers in the list are uninstalled, this can become empty.
        /// </summary>
        public InputList<Pulumi.Kubernetes.Types.Inputs.Storage.V1Beta1.CSINodeDriverArgs> Drivers
        {
            get => _drivers ?? (_drivers = new InputList<Pulumi.Kubernetes.Types.Inputs.Storage.V1Beta1.CSINodeDriverArgs>());
            set => _drivers = value;
        }

        public CSINodeSpecArgs()
        {
        }
        public static new CSINodeSpecArgs Empty => new CSINodeSpecArgs();
    }
}
