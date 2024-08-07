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
    /// ClientIPConfig represents the configurations of Client IP based session affinity.
    /// </summary>
    [OutputType]
    public sealed class ClientIPConfig
    {
        /// <summary>
        /// timeoutSeconds specifies the seconds of ClientIP type session sticky time. The value must be &gt;0 &amp;&amp; &lt;=86400(for 1 day) if ServiceAffinity == "ClientIP". Default value is 10800(for 3 hours).
        /// </summary>
        public readonly int TimeoutSeconds;

        [OutputConstructor]
        private ClientIPConfig(int timeoutSeconds)
        {
            TimeoutSeconds = timeoutSeconds;
        }
    }
}
