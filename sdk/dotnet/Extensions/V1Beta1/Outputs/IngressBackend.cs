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
    /// IngressBackend describes all endpoints for a given service and port.
    /// </summary>
    [OutputType]
    public sealed class IngressBackend
    {
        /// <summary>
        /// Resource is an ObjectRef to another Kubernetes resource in the namespace of the Ingress object. If resource is specified, serviceName and servicePort must not be specified.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Core.V1.TypedLocalObjectReference Resource;
        /// <summary>
        /// Specifies the name of the referenced service.
        /// </summary>
        public readonly string ServiceName;
        /// <summary>
        /// Specifies the port of the referenced service.
        /// </summary>
        public readonly Union<int, string> ServicePort;

        [OutputConstructor]
        private IngressBackend(
            Pulumi.Kubernetes.Types.Outputs.Core.V1.TypedLocalObjectReference resource,

            string serviceName,

            Union<int, string> servicePort)
        {
            Resource = resource;
            ServiceName = serviceName;
            ServicePort = servicePort;
        }
    }
}
