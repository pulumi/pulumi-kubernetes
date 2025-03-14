// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Storage.V1
{

    /// <summary>
    /// StorageClass describes the parameters for a class of storage for which PersistentVolumes can be dynamically provisioned.
    /// 
    /// StorageClasses are non-namespaced; the name of the storage class according to etcd is in ObjectMeta.Name.
    /// </summary>
    [OutputType]
    public sealed class StorageClass
    {
        /// <summary>
        /// allowVolumeExpansion shows whether the storage class allow volume expand.
        /// </summary>
        public readonly bool AllowVolumeExpansion;
        /// <summary>
        /// allowedTopologies restrict the node topologies where volumes can be dynamically provisioned. Each volume plugin defines its own supported topology specifications. An empty TopologySelectorTerm list means there is no topology restriction. This field is only honored by servers that enable the VolumeScheduling feature.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.TopologySelectorTerm> AllowedTopologies;
        /// <summary>
        /// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        /// </summary>
        public readonly string ApiVersion;
        /// <summary>
        /// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        /// </summary>
        public readonly string Kind;
        /// <summary>
        /// Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Meta.V1.ObjectMeta Metadata;
        /// <summary>
        /// mountOptions controls the mountOptions for dynamically provisioned PersistentVolumes of this storage class. e.g. ["ro", "soft"]. Not validated - mount of the PVs will simply fail if one is invalid.
        /// </summary>
        public readonly ImmutableArray<string> MountOptions;
        /// <summary>
        /// parameters holds the parameters for the provisioner that should create volumes of this storage class.
        /// </summary>
        public readonly ImmutableDictionary<string, string> Parameters;
        /// <summary>
        /// provisioner indicates the type of the provisioner.
        /// </summary>
        public readonly string Provisioner;
        /// <summary>
        /// reclaimPolicy controls the reclaimPolicy for dynamically provisioned PersistentVolumes of this storage class. Defaults to Delete.
        /// </summary>
        public readonly string ReclaimPolicy;
        /// <summary>
        /// volumeBindingMode indicates how PersistentVolumeClaims should be provisioned and bound.  When unset, VolumeBindingImmediate is used. This field is only honored by servers that enable the VolumeScheduling feature.
        /// </summary>
        public readonly string VolumeBindingMode;

        [OutputConstructor]
        private StorageClass(
            bool allowVolumeExpansion,

            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.Core.V1.TopologySelectorTerm> allowedTopologies,

            string apiVersion,

            string kind,

            Pulumi.Kubernetes.Types.Outputs.Meta.V1.ObjectMeta metadata,

            ImmutableArray<string> mountOptions,

            ImmutableDictionary<string, string> parameters,

            string provisioner,

            string reclaimPolicy,

            string volumeBindingMode)
        {
            AllowVolumeExpansion = allowVolumeExpansion;
            AllowedTopologies = allowedTopologies;
            ApiVersion = apiVersion;
            Kind = kind;
            Metadata = metadata;
            MountOptions = mountOptions;
            Parameters = parameters;
            Provisioner = provisioner;
            ReclaimPolicy = reclaimPolicy;
            VolumeBindingMode = volumeBindingMode;
        }
    }
}
