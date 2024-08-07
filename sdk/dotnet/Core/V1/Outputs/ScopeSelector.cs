// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Core.V1
{

    /// <summary>
    /// A scope selector represents the AND of the selectors represented by the scoped-resource selector requirements.
    /// </summary>
    [OutputType]
    public sealed class ScopeSelector
    {
        /// <summary>
        /// A list of scope selector requirements by scope of the resources.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.ScopedResourceSelectorRequirement> MatchExpressions;

        [OutputConstructor]
        private ScopeSelector(ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.ScopedResourceSelectorRequirement> matchExpressions)
        {
            MatchExpressions = matchExpressions;
        }
    }
}
