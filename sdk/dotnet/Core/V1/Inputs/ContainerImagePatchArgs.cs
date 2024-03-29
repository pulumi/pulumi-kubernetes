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
    /// Describe a container image
    /// </summary>
    public class ContainerImagePatchArgs : Pulumi.ResourceArgs
    {
        [Input("names")]
        private InputList<string>? _names;

        /// <summary>
        /// Names by which this image is known. e.g. ["k8s.gcr.io/hyperkube:v1.0.7", "dockerhub.io/google_containers/hyperkube:v1.0.7"]
        /// </summary>
        public InputList<string> Names
        {
            get => _names ?? (_names = new InputList<string>());
            set => _names = value;
        }

        /// <summary>
        /// The size of the image in bytes.
        /// </summary>
        [Input("sizeBytes")]
        public Input<int>? SizeBytes { get; set; }

        public ContainerImagePatchArgs()
        {
        }
    }
}
