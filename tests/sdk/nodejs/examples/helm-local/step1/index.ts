// Copyright 2016-2020, Pulumi Corporation.
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
import * as pulumi from "@pulumi/pulumi";

const namespace = new k8s.core.v1.Namespace("test");
const namespaceName = namespace.metadata.name;

function chart(resourcePrefix?: string): k8s.helm.v2.Chart {
    return new k8s.helm.v2.Chart("nginx", {
        path: "nginx",
        namespace: namespaceName,
        resourcePrefix: resourcePrefix,
        values: {
            service: { type: "ClusterIP" }
        },
        transformations: [
            (obj: any, opts: pulumi.CustomResourceOptions) => {
                if (obj.kind == "Service" && obj.apiVersion == "v1") {
                    opts.additionalSecretOutputs = ["status"];
                }
            }
        ]
    });
}

// Create the first instance of the Chart.
const nginx = chart();

// Create a ConfigMap depending on the Chart. The ConfigMap should not be created until after all of the Chart
// resources are ready.
new k8s.core.v1.ConfigMap("foo", {
    metadata: { namespace: namespaceName },
    data: {foo: "bar"}
}, {dependsOn: nginx})

// Export the (cluster-private) IP address of the Guestbook frontend.
const frontendServiceSpec = pulumi.all([namespaceName, nginx]).apply(([nsName, nginx]) =>
    nginx.getResourceProperty("v1/Service", nsName, "nginx", "spec"));
export const frontendServiceIP = frontendServiceSpec.clusterIP;

// Deploy a duplicate chart with a different resource prefix to verify that multiple instances of the Chart
// can be managed in the same stack.
chart("dup");
