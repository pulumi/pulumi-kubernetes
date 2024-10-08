// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "../../utilities";

// Export members:
export { IngressArgs } from "./ingress";
export type Ingress = import("./ingress").Ingress;
export const Ingress: typeof import("./ingress").Ingress = null as any;
utilities.lazyLoad(exports, ["Ingress"], () => require("./ingress"));

export { IngressClassArgs } from "./ingressClass";
export type IngressClass = import("./ingressClass").IngressClass;
export const IngressClass: typeof import("./ingressClass").IngressClass = null as any;
utilities.lazyLoad(exports, ["IngressClass"], () => require("./ingressClass"));

export { IngressClassListArgs } from "./ingressClassList";
export type IngressClassList = import("./ingressClassList").IngressClassList;
export const IngressClassList: typeof import("./ingressClassList").IngressClassList = null as any;
utilities.lazyLoad(exports, ["IngressClassList"], () => require("./ingressClassList"));

export { IngressClassPatchArgs } from "./ingressClassPatch";
export type IngressClassPatch = import("./ingressClassPatch").IngressClassPatch;
export const IngressClassPatch: typeof import("./ingressClassPatch").IngressClassPatch = null as any;
utilities.lazyLoad(exports, ["IngressClassPatch"], () => require("./ingressClassPatch"));

export { IngressListArgs } from "./ingressList";
export type IngressList = import("./ingressList").IngressList;
export const IngressList: typeof import("./ingressList").IngressList = null as any;
utilities.lazyLoad(exports, ["IngressList"], () => require("./ingressList"));

export { IngressPatchArgs } from "./ingressPatch";
export type IngressPatch = import("./ingressPatch").IngressPatch;
export const IngressPatch: typeof import("./ingressPatch").IngressPatch = null as any;
utilities.lazyLoad(exports, ["IngressPatch"], () => require("./ingressPatch"));

export { IPAddressArgs } from "./ipaddress";
export type IPAddress = import("./ipaddress").IPAddress;
export const IPAddress: typeof import("./ipaddress").IPAddress = null as any;
utilities.lazyLoad(exports, ["IPAddress"], () => require("./ipaddress"));

export { IPAddressListArgs } from "./ipaddressList";
export type IPAddressList = import("./ipaddressList").IPAddressList;
export const IPAddressList: typeof import("./ipaddressList").IPAddressList = null as any;
utilities.lazyLoad(exports, ["IPAddressList"], () => require("./ipaddressList"));

export { IPAddressPatchArgs } from "./ipaddressPatch";
export type IPAddressPatch = import("./ipaddressPatch").IPAddressPatch;
export const IPAddressPatch: typeof import("./ipaddressPatch").IPAddressPatch = null as any;
utilities.lazyLoad(exports, ["IPAddressPatch"], () => require("./ipaddressPatch"));

export { ServiceCIDRArgs } from "./serviceCIDR";
export type ServiceCIDR = import("./serviceCIDR").ServiceCIDR;
export const ServiceCIDR: typeof import("./serviceCIDR").ServiceCIDR = null as any;
utilities.lazyLoad(exports, ["ServiceCIDR"], () => require("./serviceCIDR"));

export { ServiceCIDRListArgs } from "./serviceCIDRList";
export type ServiceCIDRList = import("./serviceCIDRList").ServiceCIDRList;
export const ServiceCIDRList: typeof import("./serviceCIDRList").ServiceCIDRList = null as any;
utilities.lazyLoad(exports, ["ServiceCIDRList"], () => require("./serviceCIDRList"));

export { ServiceCIDRPatchArgs } from "./serviceCIDRPatch";
export type ServiceCIDRPatch = import("./serviceCIDRPatch").ServiceCIDRPatch;
export const ServiceCIDRPatch: typeof import("./serviceCIDRPatch").ServiceCIDRPatch = null as any;
utilities.lazyLoad(exports, ["ServiceCIDRPatch"], () => require("./serviceCIDRPatch"));


const _module = {
    version: utilities.getVersion(),
    construct: (name: string, type: string, urn: string): pulumi.Resource => {
        switch (type) {
            case "kubernetes:networking.k8s.io/v1beta1:IPAddress":
                return new IPAddress(name, <any>undefined, { urn })
            case "kubernetes:networking.k8s.io/v1beta1:IPAddressList":
                return new IPAddressList(name, <any>undefined, { urn })
            case "kubernetes:networking.k8s.io/v1beta1:IPAddressPatch":
                return new IPAddressPatch(name, <any>undefined, { urn })
            case "kubernetes:networking.k8s.io/v1beta1:Ingress":
                return new Ingress(name, <any>undefined, { urn })
            case "kubernetes:networking.k8s.io/v1beta1:IngressClass":
                return new IngressClass(name, <any>undefined, { urn })
            case "kubernetes:networking.k8s.io/v1beta1:IngressClassList":
                return new IngressClassList(name, <any>undefined, { urn })
            case "kubernetes:networking.k8s.io/v1beta1:IngressClassPatch":
                return new IngressClassPatch(name, <any>undefined, { urn })
            case "kubernetes:networking.k8s.io/v1beta1:IngressList":
                return new IngressList(name, <any>undefined, { urn })
            case "kubernetes:networking.k8s.io/v1beta1:IngressPatch":
                return new IngressPatch(name, <any>undefined, { urn })
            case "kubernetes:networking.k8s.io/v1beta1:ServiceCIDR":
                return new ServiceCIDR(name, <any>undefined, { urn })
            case "kubernetes:networking.k8s.io/v1beta1:ServiceCIDRList":
                return new ServiceCIDRList(name, <any>undefined, { urn })
            case "kubernetes:networking.k8s.io/v1beta1:ServiceCIDRPatch":
                return new ServiceCIDRPatch(name, <any>undefined, { urn })
            default:
                throw new Error(`unknown resource type ${type}`);
        }
    },
};
pulumi.runtime.registerResourceModule("kubernetes", "networking.k8s.io/v1beta1", _module)
