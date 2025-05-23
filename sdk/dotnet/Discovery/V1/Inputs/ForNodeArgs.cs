// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Inputs.Discovery.V1
{

    /// <summary>
    /// ForNode provides information about which nodes should consume this endpoint.
    /// </summary>
    public class ForNodeArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// name represents the name of the node.
        /// </summary>
        [Input("name", required: true)]
        public Input<string> Name { get; set; } = null!;

        public ForNodeArgs()
        {
        }
        public static new ForNodeArgs Empty => new ForNodeArgs();
    }
}
