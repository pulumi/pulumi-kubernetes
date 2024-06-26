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
    /// A node selector represents the union of the results of one or more label queries over a set of nodes; that is, it represents the OR of the selectors represented by the node selector terms.
    /// </summary>
    public class NodeSelectorArgs : global::Pulumi.ResourceArgs
    {
        [Input("nodeSelectorTerms", required: true)]
        private InputList<Pulumi.Kubernetes.Types.Inputs.Core.V1.NodeSelectorTermArgs>? _nodeSelectorTerms;

        /// <summary>
        /// Required. A list of node selector terms. The terms are ORed.
        /// </summary>
        public InputList<Pulumi.Kubernetes.Types.Inputs.Core.V1.NodeSelectorTermArgs> NodeSelectorTerms
        {
            get => _nodeSelectorTerms ?? (_nodeSelectorTerms = new InputList<Pulumi.Kubernetes.Types.Inputs.Core.V1.NodeSelectorTermArgs>());
            set => _nodeSelectorTerms = value;
        }

        public NodeSelectorArgs()
        {
        }
        public static new NodeSelectorArgs Empty => new NodeSelectorArgs();
    }
}
