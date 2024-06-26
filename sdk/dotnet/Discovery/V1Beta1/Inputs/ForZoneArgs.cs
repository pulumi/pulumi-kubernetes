// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Inputs.Discovery.V1Beta1
{

    /// <summary>
    /// ForZone provides information about which zones should consume this endpoint.
    /// </summary>
    public class ForZoneArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// name represents the name of the zone.
        /// </summary>
        [Input("name", required: true)]
        public Input<string> Name { get; set; } = null!;

        public ForZoneArgs()
        {
        }
        public static new ForZoneArgs Empty => new ForZoneArgs();
    }
}
