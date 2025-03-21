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
    /// Represents storage that is managed by an external CSI volume driver
    /// </summary>
    [OutputType]
    public sealed class CSIPersistentVolumeSourcePatch
    {
        /// <summary>
        /// controllerExpandSecretRef is a reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI ControllerExpandVolume call. This field is optional, and may be empty if no secret is required. If the secret object contains more than one secret, all secrets are passed.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Core.V1.SecretReferencePatch ControllerExpandSecretRef;
        /// <summary>
        /// controllerPublishSecretRef is a reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI ControllerPublishVolume and ControllerUnpublishVolume calls. This field is optional, and may be empty if no secret is required. If the secret object contains more than one secret, all secrets are passed.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Core.V1.SecretReferencePatch ControllerPublishSecretRef;
        /// <summary>
        /// driver is the name of the driver to use for this volume. Required.
        /// </summary>
        public readonly string Driver;
        /// <summary>
        /// fsType to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs".
        /// </summary>
        public readonly string FsType;
        /// <summary>
        /// nodeExpandSecretRef is a reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI NodeExpandVolume call. This field is optional, may be omitted if no secret is required. If the secret object contains more than one secret, all secrets are passed.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Core.V1.SecretReferencePatch NodeExpandSecretRef;
        /// <summary>
        /// nodePublishSecretRef is a reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI NodePublishVolume and NodeUnpublishVolume calls. This field is optional, and may be empty if no secret is required. If the secret object contains more than one secret, all secrets are passed.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Core.V1.SecretReferencePatch NodePublishSecretRef;
        /// <summary>
        /// nodeStageSecretRef is a reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI NodeStageVolume and NodeStageVolume and NodeUnstageVolume calls. This field is optional, and may be empty if no secret is required. If the secret object contains more than one secret, all secrets are passed.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Core.V1.SecretReferencePatch NodeStageSecretRef;
        /// <summary>
        /// readOnly value to pass to ControllerPublishVolumeRequest. Defaults to false (read/write).
        /// </summary>
        public readonly bool ReadOnly;
        /// <summary>
        /// volumeAttributes of the volume to publish.
        /// </summary>
        public readonly ImmutableDictionary<string, string> VolumeAttributes;
        /// <summary>
        /// volumeHandle is the unique volume name returned by the CSI volume plugin’s CreateVolume to refer to the volume on all subsequent calls. Required.
        /// </summary>
        public readonly string VolumeHandle;

        [OutputConstructor]
        private CSIPersistentVolumeSourcePatch(
            Pulumi.Kubernetes.Types.Outputs.Core.V1.SecretReferencePatch controllerExpandSecretRef,

            Pulumi.Kubernetes.Types.Outputs.Core.V1.SecretReferencePatch controllerPublishSecretRef,

            string driver,

            string fsType,

            Pulumi.Kubernetes.Types.Outputs.Core.V1.SecretReferencePatch nodeExpandSecretRef,

            Pulumi.Kubernetes.Types.Outputs.Core.V1.SecretReferencePatch nodePublishSecretRef,

            Pulumi.Kubernetes.Types.Outputs.Core.V1.SecretReferencePatch nodeStageSecretRef,

            bool readOnly,

            ImmutableDictionary<string, string> volumeAttributes,

            string volumeHandle)
        {
            ControllerExpandSecretRef = controllerExpandSecretRef;
            ControllerPublishSecretRef = controllerPublishSecretRef;
            Driver = driver;
            FsType = fsType;
            NodeExpandSecretRef = nodeExpandSecretRef;
            NodePublishSecretRef = nodePublishSecretRef;
            NodeStageSecretRef = nodeStageSecretRef;
            ReadOnly = readOnly;
            VolumeAttributes = volumeAttributes;
            VolumeHandle = volumeHandle;
        }
    }
}
