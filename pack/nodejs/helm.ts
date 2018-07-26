import { execSync } from "child_process";

import * as k8s from "./index";
import * as pulumi from "@pulumi/pulumi";
import * as yaml from "js-yaml";

export namespace v2 {
    export class Chart extends pulumi.ComponentResource {
        public readonly resources: pulumi.CustomResource[];

        constructor(name: string, path: string, values: object, opts?: pulumi.ResourceOptions) {
            super("kubernetes:v2:Chart", name, { path: path, values: values }, opts);
            // Does not require Tiller. From the `helm template` documentation:
            //
            // >  Render chart templates locally and display the output.
            // >
            // > This does not require Tiller. However, any values that would normally be
            // > looked up or retrieved in-cluster will be faked locally. Additionally, none
            // > of the server-side testing of chart validity (e.g. whether an API is supported)
            // > is done.
            const yamlStream = execSync(
                `helm template ${path} --name-template="${name}"`
            ).toString();
            const resourcesObjects = yaml.safeLoadAll(yamlStream);
            this.resources = k8s.fromJson(resourcesObjects, { ...opts, parent: this });
        }
    }
}