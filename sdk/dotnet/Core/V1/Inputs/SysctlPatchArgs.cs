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
    /// Sysctl defines a kernel parameter to be set
    /// </summary>
    public class SysctlPatchArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// Name of a property to set
        /// </summary>
        [Input("name")]
        public Input<string>? Name { get; set; }

        /// <summary>
        /// Value of a property to set
        /// </summary>
        [Input("value")]
        public Input<string>? Value { get; set; }

        public SysctlPatchArgs()
        {
        }
        public static new SysctlPatchArgs Empty => new SysctlPatchArgs();
    }
}
