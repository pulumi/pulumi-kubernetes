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
    /// ResourceHealth represents the health of a resource. It has the latest device health information. This is a part of KEP https://kep.k8s.io/4680.
    /// </summary>
    public class ResourceHealthArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// Health of the resource. can be one of:
        ///  - Healthy: operates as normal
        ///  - Unhealthy: reported unhealthy. We consider this a temporary health issue
        ///               since we do not have a mechanism today to distinguish
        ///               temporary and permanent issues.
        ///  - Unknown: The status cannot be determined.
        ///             For example, Device Plugin got unregistered and hasn't been re-registered since.
        /// 
        /// In future we may want to introduce the PermanentlyUnhealthy Status.
        /// </summary>
        [Input("health")]
        public Input<string>? Health { get; set; }

        /// <summary>
        /// ResourceID is the unique identifier of the resource. See the ResourceID type for more information.
        /// </summary>
        [Input("resourceID", required: true)]
        public Input<string> ResourceID { get; set; } = null!;

        public ResourceHealthArgs()
        {
        }
        public static new ResourceHealthArgs Empty => new ResourceHealthArgs();
    }
}
