// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Resource.V1Alpha3
{

    /// <summary>
    /// The device this taint is attached to has the "effect" on any claim which does not tolerate the taint and, through the claim, to pods using the claim.
    /// </summary>
    [OutputType]
    public sealed class DeviceTaintPatch
    {
        /// <summary>
        /// The effect of the taint on claims that do not tolerate the taint and through such claims on the pods using them. Valid effects are NoSchedule and NoExecute. PreferNoSchedule as used for nodes is not valid here.
        /// </summary>
        public readonly string Effect;
        /// <summary>
        /// The taint key to be applied to a device. Must be a label name.
        /// </summary>
        public readonly string Key;
        /// <summary>
        /// TimeAdded represents the time at which the taint was added. Added automatically during create or update if not set.
        /// </summary>
        public readonly string TimeAdded;
        /// <summary>
        /// The taint value corresponding to the taint key. Must be a label value.
        /// </summary>
        public readonly string Value;

        [OutputConstructor]
        private DeviceTaintPatch(
            string effect,

            string key,

            string timeAdded,

            string value)
        {
            Effect = effect;
            Key = key;
            TimeAdded = timeAdded;
            Value = value;
        }
    }
}
