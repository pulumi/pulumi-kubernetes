// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Extensions.V1Beta1
{

    /// <summary>
    /// HostPortRange defines a range of host ports that will be enabled by a policy for pods to use.  It requires both the start and end to be defined. Deprecated: use HostPortRange from policy API Group instead.
    /// </summary>
    [OutputType]
    public sealed class HostPortRange
    {
        /// <summary>
        /// max is the end of the range, inclusive.
        /// </summary>
        public readonly int Max;
        /// <summary>
        /// min is the start of the range, inclusive.
        /// </summary>
        public readonly int Min;

        [OutputConstructor]
        private HostPortRange(
            int max,

            int min)
        {
            Max = max;
            Min = min;
        }
    }
}
