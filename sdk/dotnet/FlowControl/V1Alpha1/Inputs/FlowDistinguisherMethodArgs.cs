// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Inputs.FlowControl.V1Alpha1
{

    /// <summary>
    /// FlowDistinguisherMethod specifies the method of a flow distinguisher.
    /// </summary>
    public class FlowDistinguisherMethodArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// `type` is the type of flow distinguisher method The supported types are "ByUser" and "ByNamespace". Required.
        /// </summary>
        [Input("type", required: true)]
        public Input<string> Type { get; set; } = null!;

        public FlowDistinguisherMethodArgs()
        {
        }
        public static new FlowDistinguisherMethodArgs Empty => new FlowDistinguisherMethodArgs();
    }
}
