// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Apps.V1Beta2
{

    /// <summary>
    /// ReplicaSetCondition describes the state of a replica set at a certain point.
    /// </summary>
    [OutputType]
    public sealed class ReplicaSetConditionPatch
    {
        /// <summary>
        /// The last time the condition transitioned from one status to another.
        /// </summary>
        public readonly string LastTransitionTime;
        /// <summary>
        /// A human readable message indicating details about the transition.
        /// </summary>
        public readonly string Message;
        /// <summary>
        /// The reason for the condition's last transition.
        /// </summary>
        public readonly string Reason;
        /// <summary>
        /// Status of the condition, one of True, False, Unknown.
        /// </summary>
        public readonly string Status;
        /// <summary>
        /// Type of replica set condition.
        /// </summary>
        public readonly string Type;

        [OutputConstructor]
        private ReplicaSetConditionPatch(
            string lastTransitionTime,

            string message,

            string reason,

            string status,

            string type)
        {
            LastTransitionTime = lastTransitionTime;
            Message = message;
            Reason = reason;
            Status = status;
            Type = type;
        }
    }
}
