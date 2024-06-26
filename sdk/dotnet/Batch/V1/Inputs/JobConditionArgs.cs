// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Inputs.Batch.V1
{

    /// <summary>
    /// JobCondition describes current state of a job.
    /// </summary>
    public class JobConditionArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// Last time the condition was checked.
        /// </summary>
        [Input("lastProbeTime")]
        public Input<string>? LastProbeTime { get; set; }

        /// <summary>
        /// Last time the condition transit from one status to another.
        /// </summary>
        [Input("lastTransitionTime")]
        public Input<string>? LastTransitionTime { get; set; }

        /// <summary>
        /// Human readable message indicating details about last transition.
        /// </summary>
        [Input("message")]
        public Input<string>? Message { get; set; }

        /// <summary>
        /// (brief) reason for the condition's last transition.
        /// </summary>
        [Input("reason")]
        public Input<string>? Reason { get; set; }

        /// <summary>
        /// Status of the condition, one of True, False, Unknown.
        /// </summary>
        [Input("status", required: true)]
        public Input<string> Status { get; set; } = null!;

        /// <summary>
        /// Type of job condition, Complete or Failed.
        /// </summary>
        [Input("type", required: true)]
        public Input<string> Type { get; set; } = null!;

        public JobConditionArgs()
        {
        }
        public static new JobConditionArgs Empty => new JobConditionArgs();
    }
}
