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

import * as k8s from "@pulumi/kubernetes";
import * as os from "os";
import * as pulumi from "@pulumi/pulumi";

const namespace = new k8s.core.v1.Namespace("test");
const namespaceName = namespace.metadata.name;

const nginx = new k8s.helm.v3.Chart("simple-nginx", {
    // Represents chart `stable/nginx-lego@v0.3.1`.
    repo: "stable",
    chart: "nginx-lego",
    version: "0.3.1",
    namespace: namespaceName,
    fetchOpts: {
        home: os.homedir(),
    },
    values: {
        // Override for the Chart's `values.yml` file. Use `null` to zero out resource requests so it
        // can be scheduled on our (wimpy) CI cluster. (Setting these values to `null` is the "normal"
        // way to delete values.)
        nginx: { resources: null },
        default: { resources: null },
        lego: { resources: null }
    },
    transformations: [
        // Make every service private to the cluster, i.e., turn all services into ClusterIP instead of
        // LoadBalancer.
        (obj: any) => {
            if (obj.kind == "Service" && obj.apiVersion == "v1") {
                if (obj.spec && obj.spec.type && obj.spec.type == "LoadBalancer") {
                    obj.spec.type = "ClusterIP";
                }
            }
        },
        (obj: any, opts: pulumi.CustomResourceOptions) => {
            if (obj.kind == "Service" && obj.apiVersion == "v1") {
                opts.additionalSecretOutputs = ["status"];
            }
        }
    ]
});

// Export the (cluster-private) IP address of the Guestbook frontend.
const frontendServiceSpec = pulumi.all([namespaceName, nginx]).apply(([nsName, nginx]) =>
    nginx.getResourceProperty("v1/Service", nsName, "simple-nginx-nginx-lego", "spec"));
export const frontendServiceIP = frontendServiceSpec.clusterIP;

// Test a variety of other inputs on a chart that creates no resources.
const empty1 = new k8s.helm.v3.Chart("empty1", {
    chart: "https://kubernetes-charts-incubator.storage.googleapis.com/raw-0.1.0.tgz",
});

const empty2 = new k8s.helm.v3.Chart("empty2", {
    chart: "raw",
    version: "0.1.0",
    fetchOpts: {
        repo: "https://kubernetes-charts-incubator.storage.googleapis.com/",
    },
});

const empty3 = new k8s.helm.v3.Chart("empty3", {
    chart: "raw",
    fetchOpts: {
        repo: "https://kubernetes-charts-incubator.storage.googleapis.com/",
    },
});
