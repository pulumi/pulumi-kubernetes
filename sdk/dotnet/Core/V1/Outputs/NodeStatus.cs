// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Core.V1
{

    /// <summary>
    /// NodeStatus is information about the current status of a node.
    /// </summary>
    [OutputType]
    public sealed class NodeStatus
    {
        /// <summary>
        /// List of addresses reachable to the node. Queried from cloud provider, if available. More info: https://kubernetes.io/docs/reference/node/node-status/#addresses Note: This field is declared as mergeable, but the merge key is not sufficiently unique, which can cause data corruption when it is merged. Callers should instead use a full-replacement patch. See https://pr.k8s.io/79391 for an example. Consumers should assume that addresses can change during the lifetime of a Node. However, there are some exceptions where this may not be possible, such as Pods that inherit a Node's address in its own status or consumers of the downward API (status.hostIP).
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeAddress> Addresses;
        /// <summary>
        /// Allocatable represents the resources of a node that are available for scheduling. Defaults to Capacity.
        /// </summary>
        public readonly ImmutableDictionary<string, string> Allocatable;
        /// <summary>
        /// Capacity represents the total resources of a node. More info: https://kubernetes.io/docs/reference/node/node-status/#capacity
        /// </summary>
        public readonly ImmutableDictionary<string, string> Capacity;
        /// <summary>
        /// Conditions is an array of current observed node conditions. More info: https://kubernetes.io/docs/reference/node/node-status/#condition
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeCondition> Conditions;
        /// <summary>
        /// Status of the config assigned to the node via the dynamic Kubelet config feature.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeConfigStatus Config;
        /// <summary>
        /// Endpoints of daemons running on the Node.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeDaemonEndpoints DaemonEndpoints;
        /// <summary>
        /// Features describes the set of features implemented by the CRI implementation.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeFeatures Features;
        /// <summary>
        /// List of container images on this node
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.ContainerImage> Images;
        /// <summary>
        /// Set of ids/uuids to uniquely identify the node. More info: https://kubernetes.io/docs/reference/node/node-status/#info
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeSystemInfo NodeInfo;
        /// <summary>
        /// NodePhase is the recently observed lifecycle phase of the node. More info: https://kubernetes.io/docs/concepts/nodes/node/#phase The field is never populated, and now is deprecated.
        /// </summary>
        public readonly string Phase;
        /// <summary>
        /// The available runtime handlers.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeRuntimeHandler> RuntimeHandlers;
        /// <summary>
        /// List of volumes that are attached to the node.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.AttachedVolume> VolumesAttached;
        /// <summary>
        /// List of attachable volumes in use (mounted) by the node.
        /// </summary>
        public readonly ImmutableArray<string> VolumesInUse;

        [OutputConstructor]
        private NodeStatus(
            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeAddress> addresses,

            ImmutableDictionary<string, string> allocatable,

            ImmutableDictionary<string, string> capacity,

            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeCondition> conditions,

            Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeConfigStatus config,

            Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeDaemonEndpoints daemonEndpoints,

            Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeFeatures features,

            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.ContainerImage> images,

            Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeSystemInfo nodeInfo,

            string phase,

            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.NodeRuntimeHandler> runtimeHandlers,

            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.AttachedVolume> volumesAttached,

            ImmutableArray<string> volumesInUse)
        {
            Addresses = addresses;
            Allocatable = allocatable;
            Capacity = capacity;
            Conditions = conditions;
            Config = config;
            DaemonEndpoints = daemonEndpoints;
            Features = features;
            Images = images;
            NodeInfo = nodeInfo;
            Phase = phase;
            RuntimeHandlers = runtimeHandlers;
            VolumesAttached = volumesAttached;
            VolumesInUse = volumesInUse;
        }
    }
}
