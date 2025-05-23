// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Networking.V1
{

    /// <summary>
    /// IPAddressSpec describe the attributes in an IP Address.
    /// </summary>
    [OutputType]
    public sealed class IPAddressSpecPatch
    {
        /// <summary>
        /// ParentRef references the resource that an IPAddress is attached to. An IPAddress must reference a parent object.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Networking.V1.ParentReferencePatch ParentRef;

        [OutputConstructor]
        private IPAddressSpecPatch(Pulumi.Kubernetes.Types.Outputs.Networking.V1.ParentReferencePatch parentRef)
        {
            ParentRef = parentRef;
        }
    }
}
