// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Inputs.Core.V1
{

    /// <summary>
    /// The weights of all of the matched WeightedPodAffinityTerm fields are added per-node to find the most preferred node(s)
    /// </summary>
    public class WeightedPodAffinityTermPatchArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// Required. A pod affinity term, associated with the corresponding weight.
        /// </summary>
        [Input("podAffinityTerm")]
        public Input<Pulumi.Kubernetes.Types.Inputs.Core.V1.PodAffinityTermPatchArgs>? PodAffinityTerm { get; set; }

        /// <summary>
        /// weight associated with matching the corresponding podAffinityTerm, in the range 1-100.
        /// </summary>
        [Input("weight")]
        public Input<int>? Weight { get; set; }

        public WeightedPodAffinityTermPatchArgs()
        {
        }
        public static new WeightedPodAffinityTermPatchArgs Empty => new WeightedPodAffinityTermPatchArgs();
    }
}
