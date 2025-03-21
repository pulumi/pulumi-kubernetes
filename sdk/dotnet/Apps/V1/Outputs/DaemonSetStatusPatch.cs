// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Apps.V1
{

    /// <summary>
    /// DaemonSetStatus represents the current status of a daemon set.
    /// </summary>
    [OutputType]
    public sealed class DaemonSetStatusPatch
    {
        /// <summary>
        /// Count of hash collisions for the DaemonSet. The DaemonSet controller uses this field as a collision avoidance mechanism when it needs to create the name for the newest ControllerRevision.
        /// </summary>
        public readonly int CollisionCount;
        /// <summary>
        /// Represents the latest available observations of a DaemonSet's current state.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Apps.V1.DaemonSetConditionPatch> Conditions;
        /// <summary>
        /// The number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod. More info: https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/
        /// </summary>
        public readonly int CurrentNumberScheduled;
        /// <summary>
        /// The total number of nodes that should be running the daemon pod (including nodes correctly running the daemon pod). More info: https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/
        /// </summary>
        public readonly int DesiredNumberScheduled;
        /// <summary>
        /// The number of nodes that should be running the daemon pod and have one or more of the daemon pod running and available (ready for at least spec.minReadySeconds)
        /// </summary>
        public readonly int NumberAvailable;
        /// <summary>
        /// The number of nodes that are running the daemon pod, but are not supposed to run the daemon pod. More info: https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/
        /// </summary>
        public readonly int NumberMisscheduled;
        /// <summary>
        /// numberReady is the number of nodes that should be running the daemon pod and have one or more of the daemon pod running with a Ready Condition.
        /// </summary>
        public readonly int NumberReady;
        /// <summary>
        /// The number of nodes that should be running the daemon pod and have none of the daemon pod running and available (ready for at least spec.minReadySeconds)
        /// </summary>
        public readonly int NumberUnavailable;
        /// <summary>
        /// The most recent generation observed by the daemon set controller.
        /// </summary>
        public readonly int ObservedGeneration;
        /// <summary>
        /// The total number of nodes that are running updated daemon pod
        /// </summary>
        public readonly int UpdatedNumberScheduled;

        [OutputConstructor]
        private DaemonSetStatusPatch(
            int collisionCount,

            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Apps.V1.DaemonSetConditionPatch> conditions,

            int currentNumberScheduled,

            int desiredNumberScheduled,

            int numberAvailable,

            int numberMisscheduled,

            int numberReady,

            int numberUnavailable,

            int observedGeneration,

            int updatedNumberScheduled)
        {
            CollisionCount = collisionCount;
            Conditions = conditions;
            CurrentNumberScheduled = currentNumberScheduled;
            DesiredNumberScheduled = desiredNumberScheduled;
            NumberAvailable = numberAvailable;
            NumberMisscheduled = numberMisscheduled;
            NumberReady = numberReady;
            NumberUnavailable = numberUnavailable;
            ObservedGeneration = observedGeneration;
            UpdatedNumberScheduled = updatedNumberScheduled;
        }
    }
}
