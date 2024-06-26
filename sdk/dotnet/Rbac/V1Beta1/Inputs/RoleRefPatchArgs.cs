// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Inputs.Rbac.V1Beta1
{

    /// <summary>
    /// RoleRef contains information that points to the role being used
    /// </summary>
    public class RoleRefPatchArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// APIGroup is the group for the resource being referenced
        /// </summary>
        [Input("apiGroup")]
        public Input<string>? ApiGroup { get; set; }

        /// <summary>
        /// Kind is the type of resource being referenced
        /// </summary>
        [Input("kind")]
        public Input<string>? Kind { get; set; }

        /// <summary>
        /// Name is the name of resource being referenced
        /// </summary>
        [Input("name")]
        public Input<string>? Name { get; set; }

        public RoleRefPatchArgs()
        {
        }
        public static new RoleRefPatchArgs Empty => new RoleRefPatchArgs();
    }
}
