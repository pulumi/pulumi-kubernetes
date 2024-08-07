// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Settings.V1Alpha1
{

    /// <summary>
    /// PodPresetSpec is a description of a pod preset.
    /// </summary>
    [OutputType]
    public sealed class PodPresetSpec
    {
        /// <summary>
        /// Env defines the collection of EnvVar to inject into containers.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.EnvVar> Env;
        /// <summary>
        /// EnvFrom defines the collection of EnvFromSource to inject into containers.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.EnvFromSource> EnvFrom;
        /// <summary>
        /// Selector is a label query over a set of resources, in this case pods. Required.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Meta.V1.LabelSelector Selector;
        /// <summary>
        /// VolumeMounts defines the collection of VolumeMount to inject into containers.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.VolumeMount> VolumeMounts;
        /// <summary>
        /// Volumes defines the collection of Volume to inject into the pod.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.Volume> Volumes;

        [OutputConstructor]
        private PodPresetSpec(
            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.EnvVar> env,

            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.EnvFromSource> envFrom,

            Pulumi.Kubernetes.Types.Outputs.Meta.V1.LabelSelector selector,

            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.VolumeMount> volumeMounts,

            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.Volume> volumes)
        {
            Env = env;
            EnvFrom = envFrom;
            Selector = selector;
            VolumeMounts = volumeMounts;
            Volumes = volumes;
        }
    }
}
