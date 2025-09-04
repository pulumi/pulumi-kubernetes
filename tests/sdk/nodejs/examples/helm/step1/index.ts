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

const nginx = new k8s.helm.v3.Chart("test", {
    chart: "nginx",
    version: "6.0.4",
    namespace: namespaceName,
    fetchOpts: {
        home: os.homedir(),
        repo: "https://raw.githubusercontent.com/bitnami/charts/eb5f9a9513d987b519f0ecd732e7031241c50328/bitnami",
    },
    values: {
        service: { type: "ClusterIP" },
        image: {
            repository: "bitnamilegacy/nginx",
            tag: "1.29.1",
        }
    },
    transformations: [
        (obj: any, opts: pulumi.CustomResourceOptions) => {
            if (obj.kind == "Service" && obj.apiVersion == "v1") {
                opts.additionalSecretOutputs = ["status"];
            }
        }
    ]
});

// Export the (cluster-private) IP address of the Guestbook frontend.
const frontendServiceSpec = pulumi.all([namespaceName, nginx]).apply(([nsName, nginx]) =>
    nginx.getResourceProperty("v1/Service", nsName, "test-nginx", "spec"));
export const frontendServiceIP = frontendServiceSpec.clusterIP;

// Test a variety of other inputs on a chart that creates no resources.
const empty1 = new k8s.helm.v3.Chart("empty1", {
    chart: "https://charts.helm.sh/incubator/packages/raw-0.1.0.tgz",
});

const empty2 = new k8s.helm.v3.Chart("empty2", {
    chart: "raw",
    version: "0.1.0",
    fetchOpts: {
        repo: "https://charts.helm.sh/incubator",
    },
});

const empty3 = new k8s.helm.v3.Chart("empty3", {
    chart: "raw",
    fetchOpts: {
        repo: "https://charts.helm.sh/incubator",
    },
});
