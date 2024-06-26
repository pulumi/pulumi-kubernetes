// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Inputs.Events.V1Beta1
{

    /// <summary>
    /// EventSeries contain information on series of events, i.e. thing that was/is happening continuously for some time.
    /// </summary>
    public class EventSeriesPatchArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// Number of occurrences in this series up to the last heartbeat time
        /// </summary>
        [Input("count")]
        public Input<int>? Count { get; set; }

        /// <summary>
        /// Time when last Event from the series was seen before last heartbeat.
        /// </summary>
        [Input("lastObservedTime")]
        public Input<string>? LastObservedTime { get; set; }

        /// <summary>
        /// Information whether this series is ongoing or finished. Deprecated. Planned removal for 1.18
        /// </summary>
        [Input("state")]
        public Input<string>? State { get; set; }

        public EventSeriesPatchArgs()
        {
        }
        public static new EventSeriesPatchArgs Empty => new EventSeriesPatchArgs();
    }
}
