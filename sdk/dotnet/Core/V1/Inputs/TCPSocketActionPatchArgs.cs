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
    /// TCPSocketAction describes an action based on opening a socket
    /// </summary>
    public class TCPSocketActionPatchArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// Optional: Host name to connect to, defaults to the pod IP.
        /// </summary>
        [Input("host")]
        public Input<string>? Host { get; set; }

        /// <summary>
        /// Number or name of the port to access on the container. Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME.
        /// </summary>
        [Input("port")]
        public InputUnion<int, string>? Port { get; set; }

        public TCPSocketActionPatchArgs()
        {
        }
        public static new TCPSocketActionPatchArgs Empty => new TCPSocketActionPatchArgs();
    }
}
