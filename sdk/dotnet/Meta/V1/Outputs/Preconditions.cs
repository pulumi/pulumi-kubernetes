// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Meta.V1
{

    [OutputType]
    public sealed class Preconditions
    {
        /// <summary>
        /// Specifies the target ResourceVersion
        /// </summary>
        public readonly string ResourceVersion;
        /// <summary>
        /// Specifies the target UID.
        /// </summary>
        public readonly string Uid;

        [OutputConstructor]
        private Preconditions(
            string resourceVersion,

            string uid)
        {
            ResourceVersion = resourceVersion;
            Uid = uid;
        }
    }
}
