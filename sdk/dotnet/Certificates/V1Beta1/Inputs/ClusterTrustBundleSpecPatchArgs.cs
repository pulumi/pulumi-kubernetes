// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Inputs.Certificates.V1Beta1
{

    /// <summary>
    /// ClusterTrustBundleSpec contains the signer and trust anchors.
    /// </summary>
    public class ClusterTrustBundleSpecPatchArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// signerName indicates the associated signer, if any.
        /// 
        /// In order to create or update a ClusterTrustBundle that sets signerName, you must have the following cluster-scoped permission: group=certificates.k8s.io resource=signers resourceName=&lt;the signer name&gt; verb=attest.
        /// 
        /// If signerName is not empty, then the ClusterTrustBundle object must be named with the signer name as a prefix (translating slashes to colons). For example, for the signer name `example.com/foo`, valid ClusterTrustBundle object names include `example.com:foo:abc` and `example.com:foo:v1`.
        /// 
        /// If signerName is empty, then the ClusterTrustBundle object's name must not have such a prefix.
        /// 
        /// List/watch requests for ClusterTrustBundles can filter on this field using a `spec.signerName=NAME` field selector.
        /// </summary>
        [Input("signerName")]
        public Input<string>? SignerName { get; set; }

        /// <summary>
        /// trustBundle contains the individual X.509 trust anchors for this bundle, as PEM bundle of PEM-wrapped, DER-formatted X.509 certificates.
        /// 
        /// The data must consist only of PEM certificate blocks that parse as valid X.509 certificates.  Each certificate must include a basic constraints extension with the CA bit set.  The API server will reject objects that contain duplicate certificates, or that use PEM block headers.
        /// 
        /// Users of ClusterTrustBundles, including Kubelet, are free to reorder and deduplicate certificate blocks in this file according to their own logic, as well as to drop PEM block headers and inter-block data.
        /// </summary>
        [Input("trustBundle")]
        public Input<string>? TrustBundle { get; set; }

        public ClusterTrustBundleSpecPatchArgs()
        {
        }
        public static new ClusterTrustBundleSpecPatchArgs Empty => new ClusterTrustBundleSpecPatchArgs();
    }
}
