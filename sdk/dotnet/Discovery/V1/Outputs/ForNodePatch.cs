// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Discovery.V1
{

    /// <summary>
    /// ForNode provides information about which nodes should consume this endpoint.
    /// </summary>
    [OutputType]
    public sealed class ForNodePatch
    {
        /// <summary>
        /// name represents the name of the node.
        /// </summary>
        public readonly string Name;

        [OutputConstructor]
        private ForNodePatch(string name)
        {
            Name = name;
        }
    }
}
