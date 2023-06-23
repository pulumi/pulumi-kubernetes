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

// Some versions of k8s have started setting managed fields for all resources, even those managed by
// Client-side Apply (CSA). This can cause field manager conflicts when a resource is upgraded to
// Server-side Apply (SSA). This test validates the following scenario does not fail with a field manager conflict:
// 1. Create a Deployment with Client-side Apply.
// 2. Update the Provider to use Server-side Apply.
// 3. Change fields in the Deployment.

// Create provider with SSA disabled.
export const provider = new k8s.Provider("k8s", {enableServerSideApply: true});

// Create a randomly-named Namespace.
const ns = new k8s.core.v1.Namespace("test", undefined, {provider});

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
                    image: "nginx:1.17", // Update a field to ensure that field manager conflicts don't cause an error
                    ports: [
                        // Change the port to ensure that the existing value is updated rather than adding a second port to the array
                        {containerPort: 81},
                    ],
                }],
            }
        }
    }
}, { provider });
