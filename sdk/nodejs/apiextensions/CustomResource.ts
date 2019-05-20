// Copyright 2016-2019, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import * as pulumi from "@pulumi/pulumi"
import * as inputApi from "../types/input";
import * as outputApi from "../types/output";
import { getVersion } from "../version";

/**
 * CustomResourceArgs represents a resource definition we'd use to create an instance of a
 * Kubernetes CustomResourceDefinition (CRD). For example, the CoreOS Prometheus operator
 * exposes a CRD `monitoring.coreos.com/ServiceMonitor`; to create a `ServiceMonitor`, we'd
 * pass a `CustomResourceArgs` containing the `ServiceMonitor` definition to
 * `apiextensions.CustomResource`.
 *
 * NOTE: This type is fairly loose, since other than `apiVersion` and `kind`, there are no
 * fields required across all CRDs.
 */
export interface CustomResourceArgs {
    /**
     * APIVersion defines the versioned schema of this representation of an object. Servers should
     * convert recognized schemas to the latest internal value, and may reject unrecognized
     * values. More info:
     * https://git.k8s.io/community/contributors/devel/api-conventions.md#resources
     */
    apiVersion: pulumi.Input<string>;

    /**
     * Kind is a string value representing the REST resource this object represents. Servers may
     * infer this from the endpoint the client submits requests to. Cannot be updated. In
     * CamelCase. More info:
     * https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
     */
    kind: pulumi.Input<string>;

    /**
     * Standard object metadata; More info:
     * https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
     */
    metadata?: pulumi.Input<inputApi.meta.v1.ObjectMeta>;
    [othersFields: string]: pulumi.Input<any>;
}

/**
 * CustomResourceGetOptions uniquely identifies a Kubernetes CustomResource, primarily for use
 * in supplied to `apiextensions.CustomResource#get`.
 */
export interface CustomResourceGetOptions extends pulumi.CustomResourceOptions {
    /**
     * apiVersion is the API version of the apiExtensions.CustomResource we wish to select,
     * as specified by the CustomResourceDefinition that defines it on the API server.
     */
    apiVersion: pulumi.Input<string>;

    /**
     * kind is the kind of the apiextensions.CustomResource we wish to select, as specified by
     * the CustomResourceDefinition that defines it on the API server.
     */
    kind: pulumi.Input<string>

    /**
     * An ID for the Kubernetes resource to retrive. Takes the form <namespace>/<name> or
     * <name>.
     */
    id: pulumi.Input<pulumi.ID>;
}

/**
 * CustomResource represents an instance of a CustomResourceDefinition (CRD). For example, the
 * CoreOS Prometheus operator exposes a CRD `monitoring.coreos.com/ServiceMonitor`; to
 * instantiate this as a Pulumi resource, one could call `new CustomResource`, passing the
 * `ServiceMonitor` resource definition as an argument.
 */
export class CustomResource extends pulumi.CustomResource {
    /**
     * APIVersion defines the versioned schema of this representation of an object. Servers should
     * convert recognized schemas to the latest internal value, and may reject unrecognized
     * values. More info:
     * https://git.k8s.io/community/contributors/devel/api-conventions.md#resources
     */
    public readonly apiVersion: pulumi.Output<string>;

    /**
     * Kind is a string value representing the REST resource this object represents. Servers may
     * infer this from the endpoint the client submits requests to. Cannot be updated. In
     * CamelCase. More info:
     * https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
     */
    public readonly kind: pulumi.Output<string>;

    /**
     * Standard object metadata; More info:
     * https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
     */
    public readonly metadata: pulumi.Output<outputApi.meta.v1.ObjectMeta>;

    /**
     * Get the state of an existing `CustomResource`, as identified by `id`.
     * Typically this ID  is of the form <namespace>/<name>; if <namespace> is omitted, then (per
     * Kubernetes convention) the ID becomes default/<name>.
     *
     * Pulumi will keep track of this resource using `name` as the Pulumi ID.
     *
     * @param name _Unique_ name used to register this resource with Pulumi.
     * @param opts Uniquely specifies a CustomResource to select.
     */
    public static get(name: string, opts: CustomResourceGetOptions): CustomResource {
        // NOTE: `selectOpts` will be type `pulumi.CustomResource`. If we add a field that does
        // not satisfy that interface, it will cause a compilation error in `...selectOpts` in
        // the constructor call below.
        const {apiVersion, kind, id, ...selectOpts} = opts;
        return new CustomResource(name, {apiVersion: apiVersion, kind: kind}, { ...selectOpts, id: id });
    }

    public getInputs(): CustomResourceArgs { return this.__inputs; }
    private readonly __inputs: CustomResourceArgs;

    /**
     * Create a CustomResource resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args: CustomResourceArgs, opts?: pulumi.CustomResourceOptions) {
        let inputs: pulumi.Inputs = {};
        for (const key of Object.keys(args)) {
            inputs[key] = (args as any)[key];
        }

        if (!opts) {
            opts = {}
        }
        if (!opts.version) {
            opts.version = getVersion();
        }
        super(`kubernetes:${args.apiVersion}:${args.kind}`, name, inputs, opts);
        this.__inputs = args;
    }
}
