// Copyright 2016-2021, Pulumi Corporation.
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

// This test creates a Provider with `enableDryRun` (server side apply) enabled. The namespace option is removed, which
// causes the Deployment to be recreated in the default namespace.

const ns = new k8s.core.v1.Namespace("test");
const provider = new k8s.Provider("k8s", {enableDryRun: true}); // Use the default namespace.

const appLabels = { app: "nginx" };
const deployment = new k8s.apps.v1.Deployment("nginx", {
    spec: {
        selector: { matchLabels: appLabels },
        template: {
            metadata: { labels: appLabels },
            spec: {
                containers: [{
                    name: "nginx",
                    image: "nginx",
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
}, { provider });
export const name = deployment.metadata.name;
