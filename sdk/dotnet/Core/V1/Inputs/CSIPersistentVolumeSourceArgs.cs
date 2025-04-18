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
    /// Represents storage that is managed by an external CSI volume driver
    /// </summary>
    public class CSIPersistentVolumeSourceArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// controllerExpandSecretRef is a reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI ControllerExpandVolume call. This field is optional, and may be empty if no secret is required. If the secret object contains more than one secret, all secrets are passed.
        /// </summary>
        [Input("controllerExpandSecretRef")]
        public Input<Pulumi.Kubernetes.Types.Inputs.Core.V1.SecretReferenceArgs>? ControllerExpandSecretRef { get; set; }

        /// <summary>
        /// controllerPublishSecretRef is a reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI ControllerPublishVolume and ControllerUnpublishVolume calls. This field is optional, and may be empty if no secret is required. If the secret object contains more than one secret, all secrets are passed.
        /// </summary>
        [Input("controllerPublishSecretRef")]
        public Input<Pulumi.Kubernetes.Types.Inputs.Core.V1.SecretReferenceArgs>? ControllerPublishSecretRef { get; set; }

        /// <summary>
        /// driver is the name of the driver to use for this volume. Required.
        /// </summary>
        [Input("driver", required: true)]
        public Input<string> Driver { get; set; } = null!;

        /// <summary>
        /// fsType to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs".
        /// </summary>
        [Input("fsType")]
        public Input<string>? FsType { get; set; }

        /// <summary>
        /// nodeExpandSecretRef is a reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI NodeExpandVolume call. This field is optional, may be omitted if no secret is required. If the secret object contains more than one secret, all secrets are passed.
        /// </summary>
        [Input("nodeExpandSecretRef")]
        public Input<Pulumi.Kubernetes.Types.Inputs.Core.V1.SecretReferenceArgs>? NodeExpandSecretRef { get; set; }

        /// <summary>
        /// nodePublishSecretRef is a reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI NodePublishVolume and NodeUnpublishVolume calls. This field is optional, and may be empty if no secret is required. If the secret object contains more than one secret, all secrets are passed.
        /// </summary>
        [Input("nodePublishSecretRef")]
        public Input<Pulumi.Kubernetes.Types.Inputs.Core.V1.SecretReferenceArgs>? NodePublishSecretRef { get; set; }

        /// <summary>
        /// nodeStageSecretRef is a reference to the secret object containing sensitive information to pass to the CSI driver to complete the CSI NodeStageVolume and NodeStageVolume and NodeUnstageVolume calls. This field is optional, and may be empty if no secret is required. If the secret object contains more than one secret, all secrets are passed.
        /// </summary>
        [Input("nodeStageSecretRef")]
        public Input<Pulumi.Kubernetes.Types.Inputs.Core.V1.SecretReferenceArgs>? NodeStageSecretRef { get; set; }

        /// <summary>
        /// readOnly value to pass to ControllerPublishVolumeRequest. Defaults to false (read/write).
        /// </summary>
        [Input("readOnly")]
        public Input<bool>? ReadOnly { get; set; }

        [Input("volumeAttributes")]
        private InputMap<string>? _volumeAttributes;

        /// <summary>
        /// volumeAttributes of the volume to publish.
        /// </summary>
        public InputMap<string> VolumeAttributes
        {
            get => _volumeAttributes ?? (_volumeAttributes = new InputMap<string>());
            set => _volumeAttributes = value;
        }

        /// <summary>
        /// volumeHandle is the unique volume name returned by the CSI volume plugin’s CreateVolume to refer to the volume on all subsequent calls. Required.
        /// </summary>
        [Input("volumeHandle", required: true)]
        public Input<string> VolumeHandle { get; set; } = null!;

        public CSIPersistentVolumeSourceArgs()
        {
        }
        public static new CSIPersistentVolumeSourceArgs Empty => new CSIPersistentVolumeSourceArgs();
    }
}
