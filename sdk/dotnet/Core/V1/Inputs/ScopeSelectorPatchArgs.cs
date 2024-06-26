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
    /// A scope selector represents the AND of the selectors represented by the scoped-resource selector requirements.
    /// </summary>
    public class ScopeSelectorPatchArgs : global::Pulumi.ResourceArgs
    {
        [Input("matchExpressions")]
        private InputList<Pulumi.Kubernetes.Types.Inputs.Core.V1.ScopedResourceSelectorRequirementPatchArgs>? _matchExpressions;

        /// <summary>
        /// A list of scope selector requirements by scope of the resources.
        /// </summary>
        public InputList<Pulumi.Kubernetes.Types.Inputs.Core.V1.ScopedResourceSelectorRequirementPatchArgs> MatchExpressions
        {
            get => _matchExpressions ?? (_matchExpressions = new InputList<Pulumi.Kubernetes.Types.Inputs.Core.V1.ScopedResourceSelectorRequirementPatchArgs>());
            set => _matchExpressions = value;
        }

        public ScopeSelectorPatchArgs()
        {
        }
        public static new ScopeSelectorPatchArgs Empty => new ScopeSelectorPatchArgs();
    }
}
