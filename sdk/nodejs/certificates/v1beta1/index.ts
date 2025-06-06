// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "../../utilities";

// Export members:
export { CertificateSigningRequestArgs } from "./certificateSigningRequest";
export type CertificateSigningRequest = import("./certificateSigningRequest").CertificateSigningRequest;
export const CertificateSigningRequest: typeof import("./certificateSigningRequest").CertificateSigningRequest = null as any;
utilities.lazyLoad(exports, ["CertificateSigningRequest"], () => require("./certificateSigningRequest"));

export { CertificateSigningRequestListArgs } from "./certificateSigningRequestList";
export type CertificateSigningRequestList = import("./certificateSigningRequestList").CertificateSigningRequestList;
export const CertificateSigningRequestList: typeof import("./certificateSigningRequestList").CertificateSigningRequestList = null as any;
utilities.lazyLoad(exports, ["CertificateSigningRequestList"], () => require("./certificateSigningRequestList"));

export { CertificateSigningRequestPatchArgs } from "./certificateSigningRequestPatch";
export type CertificateSigningRequestPatch = import("./certificateSigningRequestPatch").CertificateSigningRequestPatch;
export const CertificateSigningRequestPatch: typeof import("./certificateSigningRequestPatch").CertificateSigningRequestPatch = null as any;
utilities.lazyLoad(exports, ["CertificateSigningRequestPatch"], () => require("./certificateSigningRequestPatch"));

export { ClusterTrustBundleArgs } from "./clusterTrustBundle";
export type ClusterTrustBundle = import("./clusterTrustBundle").ClusterTrustBundle;
export const ClusterTrustBundle: typeof import("./clusterTrustBundle").ClusterTrustBundle = null as any;
utilities.lazyLoad(exports, ["ClusterTrustBundle"], () => require("./clusterTrustBundle"));

export { ClusterTrustBundleListArgs } from "./clusterTrustBundleList";
export type ClusterTrustBundleList = import("./clusterTrustBundleList").ClusterTrustBundleList;
export const ClusterTrustBundleList: typeof import("./clusterTrustBundleList").ClusterTrustBundleList = null as any;
utilities.lazyLoad(exports, ["ClusterTrustBundleList"], () => require("./clusterTrustBundleList"));

export { ClusterTrustBundlePatchArgs } from "./clusterTrustBundlePatch";
export type ClusterTrustBundlePatch = import("./clusterTrustBundlePatch").ClusterTrustBundlePatch;
export const ClusterTrustBundlePatch: typeof import("./clusterTrustBundlePatch").ClusterTrustBundlePatch = null as any;
utilities.lazyLoad(exports, ["ClusterTrustBundlePatch"], () => require("./clusterTrustBundlePatch"));


const _module = {
    version: utilities.getVersion(),
    construct: (name: string, type: string, urn: string): pulumi.Resource => {
        switch (type) {
            case "kubernetes:certificates.k8s.io/v1beta1:CertificateSigningRequest":
                return new CertificateSigningRequest(name, <any>undefined, { urn })
            case "kubernetes:certificates.k8s.io/v1beta1:CertificateSigningRequestList":
                return new CertificateSigningRequestList(name, <any>undefined, { urn })
            case "kubernetes:certificates.k8s.io/v1beta1:CertificateSigningRequestPatch":
                return new CertificateSigningRequestPatch(name, <any>undefined, { urn })
            case "kubernetes:certificates.k8s.io/v1beta1:ClusterTrustBundle":
                return new ClusterTrustBundle(name, <any>undefined, { urn })
            case "kubernetes:certificates.k8s.io/v1beta1:ClusterTrustBundleList":
                return new ClusterTrustBundleList(name, <any>undefined, { urn })
            case "kubernetes:certificates.k8s.io/v1beta1:ClusterTrustBundlePatch":
                return new ClusterTrustBundlePatch(name, <any>undefined, { urn })
            default:
                throw new Error(`unknown resource type ${type}`);
        }
    },
};
pulumi.runtime.registerResourceModule("kubernetes", "certificates.k8s.io/v1beta1", _module)
