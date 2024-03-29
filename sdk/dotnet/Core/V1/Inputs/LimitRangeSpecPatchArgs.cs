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
    /// LimitRangeSpec defines a min/max usage limit for resources that match on kind.
    /// </summary>
    public class LimitRangeSpecPatchArgs : global::Pulumi.ResourceArgs
    {
        [Input("limits")]
        private InputList<Pulumi.Kubernetes.Types.Inputs.Core.V1.LimitRangeItemPatchArgs>? _limits;

        /// <summary>
        /// Limits is the list of LimitRangeItem objects that are enforced.
        /// </summary>
        public InputList<Pulumi.Kubernetes.Types.Inputs.Core.V1.LimitRangeItemPatchArgs> Limits
        {
            get => _limits ?? (_limits = new InputList<Pulumi.Kubernetes.Types.Inputs.Core.V1.LimitRangeItemPatchArgs>());
            set => _limits = value;
        }

        public LimitRangeSpecPatchArgs()
        {
        }
        public static new LimitRangeSpecPatchArgs Empty => new LimitRangeSpecPatchArgs();
    }
}
