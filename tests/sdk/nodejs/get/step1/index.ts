// Copyright 2016-2023, Pulumi Corporation.
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

const namespace = new k8s.core.v1.Namespace("test-namespace");

//
// `get`s the Kubernetes API service.
//

export const svc = k8s.core.v1.Service.get("kube-api", "default/kubernetes");

// This will fail with a TypeError if the status was not populated (i.e. the .get isn't working)
export const loadBalancer = svc.status.loadBalancer;

//
// Create a Service resource with skipAwait that would fail to initialize due to await logic.
// get should return the resource state without gating on await logic.
//

export const awaitSvc = new k8s.core.v1.Service("svc", {
    metadata: {
        name: "test",
        namespace: namespace.metadata.name,
        annotations: {
            "pulumi.com/skipAwait": "true",
        },
    },
    spec: {
        type: k8s.types.enums.core.v1.ServiceSpecType.ClusterIP,
        ports: [{
            name: "http",
            port: 8080,
            targetPort: 80,
        }],
        selector: { app: "nginx" }, // selector doesn't match Pods, so await logic would fail
    }
});

//
// Create a CustomResourceDefinition, a CustomResource, and then `.get` it.
//

export const ct = new k8s.apiextensions.v1.CustomResourceDefinition("crontab", {
    metadata: { name: "crontabs.stable.example.com" },
    spec: {
        group: "stable.example.com",
        versions: [
            {
                name: "v1",
                served: true,
                storage: true,
                schema: {
                    openAPIV3Schema: {
                        type: "object",
                        properties: {
                            spec: {
                                type: "object",
                                properties: {
                                    cronSpec: {
                                        type: "string"
                                    },
                                    image: {
                                        type: "string"
                                    },
                                    replicas: {
                                        type: "integer"
                                    }
                                }
                            }
                        }
                    }
                }
            },
        ],
        scope: "Namespaced",
        names: {
            plural: "crontabs",
            singular: "crontab",
            kind: "CronTab",
            shortNames: ["ct"]
        },
    }
});

export const cr = new k8s.apiextensions.CustomResource(
    "my-new-cron-object",
    {
        apiVersion: "stable.example.com/v1",
        kind: "CronTab",
      metadata: {
        namespace: namespace.metadata.name,
        name: "my-new-cron-object",
      },
        spec: { cronSpec: "* * * * */5", image: "my-awesome-cron-image" }
    },
    { dependsOn: ct }
);
