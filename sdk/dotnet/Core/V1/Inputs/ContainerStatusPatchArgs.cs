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
    /// ContainerStatus contains details for the current status of this container.
    /// </summary>
    public class ContainerStatusPatchArgs : Pulumi.ResourceArgs
    {
        /// <summary>
        /// Container's ID in the format '&lt;type&gt;://&lt;container_id&gt;'.
        /// </summary>
        [Input("containerID")]
        public Input<string>? ContainerID { get; set; }

        /// <summary>
        /// The image the container is running. More info: https://kubernetes.io/docs/concepts/containers/images.
        /// </summary>
        [Input("image")]
        public Input<string>? Image { get; set; }

        /// <summary>
        /// ImageID of the container's image.
        /// </summary>
        [Input("imageID")]
        public Input<string>? ImageID { get; set; }

        /// <summary>
        /// Details about the container's last termination condition.
        /// </summary>
        [Input("lastState")]
        public Input<Pulumi.Kubernetes.Types.Inputs.Core.V1.ContainerStatePatchArgs>? LastState { get; set; }

        /// <summary>
        /// This must be a DNS_LABEL. Each container in a pod must have a unique name. Cannot be updated.
        /// </summary>
        [Input("name")]
        public Input<string>? Name { get; set; }

        /// <summary>
        /// Specifies whether the container has passed its readiness probe.
        /// </summary>
        [Input("ready")]
        public Input<bool>? Ready { get; set; }

        /// <summary>
        /// The number of times the container has been restarted.
        /// </summary>
        [Input("restartCount")]
        public Input<int>? RestartCount { get; set; }

        /// <summary>
        /// Specifies whether the container has passed its startup probe. Initialized as false, becomes true after startupProbe is considered successful. Resets to false when the container is restarted, or if kubelet loses state temporarily. Is always true when no startupProbe is defined.
        /// </summary>
        [Input("started")]
        public Input<bool>? Started { get; set; }

        /// <summary>
        /// Details about the container's current condition.
        /// </summary>
        [Input("state")]
        public Input<Pulumi.Kubernetes.Types.Inputs.Core.V1.ContainerStatePatchArgs>? State { get; set; }

        public ContainerStatusPatchArgs()
        {
        }
    }
}
