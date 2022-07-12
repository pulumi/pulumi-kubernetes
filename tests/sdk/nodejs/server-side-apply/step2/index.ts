// Copyright 2016-2022, Pulumi Corporation.
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

// This test creates a Provider with `enableServerSideApply` enabled. The following scenarios are tested:
// 1. Patch a Namespace resource with fully-specified configuration.
// 2. Patch a CustomResource.
// 3. Upsert a Deployment resource that already exists.
// 4. Patch the Deployment with a partially-specified configuration.
// 5. Replace a statically-named ConfigMap resource by changing the data on a subsequent update.

// Create provider with SSA enabled.
const provider = new k8s.Provider("k8s", {enableServerSideApply: true});

// Create a randomly-named Namespace.
const ns = new k8s.core.v1.Namespace("test", undefined, {provider});

// Patch the Namespace with additional labels.
export const nsPatched = new k8s.core.v1.NamespacePatch("test", {
    metadata: {
        labels: {
            foo: "foo",
        },
        name: ns.metadata.name,
    }
}, {provider});

const crd = new k8s.apiextensions.v1.CustomResourceDefinition("crontab", {
    metadata: { name: "crontabs.nodessa.example.com" },
    spec: {
        group: "nodessa.example.com",
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
}, {provider});

// Create a k8s CustomResource of type "CronTab".
const cr = new k8s.apiextensions.CustomResource("my-new-cron-object",
    {
        apiVersion: "nodessa.example.com/v1",
        kind: "CronTab",
        metadata: {
            name: "my-new-cron-object",
            namespace: ns.metadata.name,
        },
        spec: { cronSpec: "* * * * */6", image: "my-awesome-cron-image" }
    }, {provider, dependsOn: [crd]});

// Patch the CronTab CustomResource to add a label.
export const crPatched = new k8s.apiextensions.CustomResourcePatch("label-cr",
    {
        apiVersion: "nodessa.example.com/v1",
        kind: "CronTab",
        metadata: {
            labels: {
                foo: "foo",
            },
            name: "my-new-cron-object",
            namespace: ns.metadata.name,
        },
    }, {provider, dependsOn: [cr]});

const appLabels = { app: "nginx" };
const deployment = new k8s.apps.v1.Deployment("nginx", {
    metadata: {
        namespace: ns.metadata.name,
    },
    spec: {
        selector: { matchLabels: appLabels },
        template: {
            metadata: { labels: appLabels },
            spec: {
                containers: [{
                    name: "nginx",
                    image: "nginx:1.16",
                }],
            }
        }
    }
}, { provider });

// Upsert an identical Deployment to the one we just created.
const upsert = new k8s.apps.v1.Deployment("nginx-upsert", {
    metadata: {
        name: deployment.metadata.name,
        namespace: ns.metadata.name,
    },
    spec: {
        selector: { matchLabels: appLabels },
        template: {
            metadata: { labels: appLabels },
            spec: {
                containers: [{
                    name: "nginx",
                    image: "nginx:1.16",
                }],
            }
        }
    }
}, { provider, dependsOn: [deployment], retainOnDelete: true });

// Patch the Deployment to add memory limits. This is a partial spec that is
// missing required values, but it works as a patch. All fields are optional
// in the SDK, and validation is disabled. The patch spec must be
// unambiguous, but doesn't need to specify fields that we don't care about.
export const deploymentPatched = new k8s.apps.v1.DeploymentPatch("nginx", {
    metadata: {
        name: deployment.metadata.name,
        namespace: ns.metadata.name,
    },
    spec: {
        template: {
            spec: {
                containers: [{
                    name: "nginx",
                    resources: {
                        limits: {
                            // This value will be normalized by the API server to "1Gi" in the returned state. Without
                            // dry-run support, Pulumi will detect a spurious diff due to the normalization.
                            memory: "1024Mi",
                        },
                    },
                }],
            }
        }
    }
}, { provider, dependsOn: [upsert]});

new k8s.core.v1.ConfigMap("test", {
    metadata: {
        name: "foo", // Specify the name to force resource replacement on change.
        namespace: ns.metadata.name,
    },
    data: {foo: "baz"}, // <-- Updated value
}, {provider});
