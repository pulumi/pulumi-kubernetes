// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Inputs.ApiExtensions.V1
{

    /// <summary>
    /// CustomResourceValidation is a list of validation methods for CustomResources.
    /// </summary>
    public class CustomResourceValidationPatchArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// openAPIV3Schema is the OpenAPI v3 schema to use for validation and pruning.
        /// </summary>
        [Input("openAPIV3Schema")]
        public Input<Pulumi.Kubernetes.Types.Inputs.ApiExtensions.V1.JSONSchemaPropsPatchArgs>? OpenAPIV3Schema { get; set; }

        public CustomResourceValidationPatchArgs()
        {
        }
        public static new CustomResourceValidationPatchArgs Empty => new CustomResourceValidationPatchArgs();
    }
}
