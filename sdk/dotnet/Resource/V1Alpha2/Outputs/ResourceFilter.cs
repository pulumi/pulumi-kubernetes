// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Resource.V1Alpha2
{

    /// <summary>
    /// ResourceFilter is a filter for resources from one particular driver.
    /// </summary>
    [OutputType]
    public sealed class ResourceFilter
    {
        /// <summary>
        /// DriverName is the name used by the DRA driver kubelet plugin.
        /// </summary>
        public readonly string DriverName;
        /// <summary>
        /// NamedResources describes a resource filter using the named resources model.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Resource.V1Alpha2.NamedResourcesFilter NamedResources;

        [OutputConstructor]
        private ResourceFilter(
            string driverName,

            Pulumi.Kubernetes.Types.Outputs.Resource.V1Alpha2.NamedResourcesFilter namedResources)
        {
            DriverName = driverName;
            NamedResources = namedResources;
        }
    }
}
