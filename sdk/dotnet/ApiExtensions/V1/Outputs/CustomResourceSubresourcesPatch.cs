// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.ApiExtensions.V1
{

    /// <summary>
    /// CustomResourceSubresources defines the status and scale subresources for CustomResources.
    /// </summary>
    [OutputType]
    public sealed class CustomResourceSubresourcesPatch
    {
        /// <summary>
        /// scale indicates the custom resource should serve a `/scale` subresource that returns an `autoscaling/v1` Scale object.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.ApiExtensions.V1.CustomResourceSubresourceScalePatch Scale;
        /// <summary>
        /// status indicates the custom resource should serve a `/status` subresource. When enabled: 1. requests to the custom resource primary endpoint ignore changes to the `status` stanza of the object. 2. requests to the custom resource `/status` subresource ignore changes to anything other than the `status` stanza of the object.
        /// </summary>
        public readonly System.Text.Json.JsonElement Status;

        [OutputConstructor]
        private CustomResourceSubresourcesPatch(
            Pulumi.Kubernetes.Types.Outputs.ApiExtensions.V1.CustomResourceSubresourceScalePatch scale,

            System.Text.Json.JsonElement status)
        {
            Scale = scale;
            Status = status;
        }
    }
}