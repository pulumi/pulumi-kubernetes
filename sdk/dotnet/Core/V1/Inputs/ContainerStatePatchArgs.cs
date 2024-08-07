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
    /// ContainerState holds a possible state of container. Only one of its members may be specified. If none of them is specified, the default one is ContainerStateWaiting.
    /// </summary>
    public class ContainerStatePatchArgs : Pulumi.ResourceArgs
    {
        /// <summary>
        /// Details about a running container
        /// </summary>
        [Input("running")]
        public Input<Pulumi.Kubernetes.Types.Inputs.Core.V1.ContainerStateRunningPatchArgs>? Running { get; set; }

        /// <summary>
        /// Details about a terminated container
        /// </summary>
        [Input("terminated")]
        public Input<Pulumi.Kubernetes.Types.Inputs.Core.V1.ContainerStateTerminatedPatchArgs>? Terminated { get; set; }

        /// <summary>
        /// Details about a waiting container
        /// </summary>
        [Input("waiting")]
        public Input<Pulumi.Kubernetes.Types.Inputs.Core.V1.ContainerStateWaitingPatchArgs>? Waiting { get; set; }

        public ContainerStatePatchArgs()
        {
        }
    }
}
