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
    /// ReplicaSetSpec is the specification of a ReplicaSet.
    /// </summary>
    [OutputType]
    public sealed class ReplicaSetSpec
    {
        /// <summary>
        /// Minimum number of seconds for which a newly created pod should be ready without any of its container crashing, for it to be considered available. Defaults to 0 (pod will be considered available as soon as it is ready)
        /// </summary>
        public readonly int MinReadySeconds;
        /// <summary>
        /// Replicas is the number of desired pods. This is a pointer to distinguish between explicit zero and unspecified. Defaults to 1. More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicaset
        /// </summary>
        public readonly int Replicas;
        /// <summary>
        /// Selector is a label query over pods that should match the replica count. Label keys and values that must match in order to be controlled by this replica set. It must match the pod template's labels. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Meta.V1.LabelSelector Selector;
        /// <summary>
        /// Template is the object that describes the pod that will be created if insufficient replicas are detected. More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/#pod-template
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Core.V1.PodTemplateSpec Template;

        [OutputConstructor]
        private ReplicaSetSpec(
            int minReadySeconds,

            int replicas,

            Pulumi.Kubernetes.Types.Outputs.Meta.V1.LabelSelector selector,

            Pulumi.Kubernetes.Types.Outputs.Core.V1.PodTemplateSpec template)
        {
            MinReadySeconds = minReadySeconds;
            Replicas = replicas;
            Selector = selector;
            Template = template;
        }
    }
}
