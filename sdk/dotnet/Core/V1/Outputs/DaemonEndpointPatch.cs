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
    /// DaemonEndpoint contains information about a single Daemon endpoint.
    /// </summary>
    [OutputType]
    public sealed class DaemonEndpointPatch
    {
        /// <summary>
        /// Port number of the given endpoint.
        /// </summary>
        public readonly int Port;

        [OutputConstructor]
        private DaemonEndpointPatch(int Port)
        {
            this.Port = Port;
        }
    }
}
