// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.AdmissionRegistration.V1
{

    /// <summary>
    /// ParamKind is a tuple of Group Kind and Version.
    /// </summary>
    [OutputType]
    public sealed class ParamKindPatch
    {
        /// <summary>
        /// APIVersion is the API group version the resources belong to. In format of "group/version". Required.
        /// </summary>
        public readonly string ApiVersion;
        /// <summary>
        /// Kind is the API kind the resources belong to. Required.
        /// </summary>
        public readonly string Kind;

        [OutputConstructor]
        private ParamKindPatch(
            string apiVersion,

            string kind)
        {
            ApiVersion = apiVersion;
            Kind = kind;
        }
    }
}
