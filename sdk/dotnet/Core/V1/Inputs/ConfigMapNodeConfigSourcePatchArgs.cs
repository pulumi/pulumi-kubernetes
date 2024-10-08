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
    /// ConfigMapNodeConfigSource contains the information to reference a ConfigMap as a config source for the Node. This API is deprecated since 1.22: https://git.k8s.io/enhancements/keps/sig-node/281-dynamic-kubelet-configuration
    /// </summary>
    public class ConfigMapNodeConfigSourcePatchArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// KubeletConfigKey declares which key of the referenced ConfigMap corresponds to the KubeletConfiguration structure This field is required in all cases.
        /// </summary>
        [Input("kubeletConfigKey")]
        public Input<string>? KubeletConfigKey { get; set; }

        /// <summary>
        /// Name is the metadata.name of the referenced ConfigMap. This field is required in all cases.
        /// </summary>
        [Input("name")]
        public Input<string>? Name { get; set; }

        /// <summary>
        /// Namespace is the metadata.namespace of the referenced ConfigMap. This field is required in all cases.
        /// </summary>
        [Input("namespace")]
        public Input<string>? Namespace { get; set; }

        /// <summary>
        /// ResourceVersion is the metadata.ResourceVersion of the referenced ConfigMap. This field is forbidden in Node.Spec, and required in Node.Status.
        /// </summary>
        [Input("resourceVersion")]
        public Input<string>? ResourceVersion { get; set; }

        /// <summary>
        /// UID is the metadata.UID of the referenced ConfigMap. This field is forbidden in Node.Spec, and required in Node.Status.
        /// </summary>
        [Input("uid")]
        public Input<string>? Uid { get; set; }

        public ConfigMapNodeConfigSourcePatchArgs()
        {
        }
        public static new ConfigMapNodeConfigSourcePatchArgs Empty => new ConfigMapNodeConfigSourcePatchArgs();
    }
}
