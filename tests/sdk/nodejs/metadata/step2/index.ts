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

import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

const appLabels = { app: "nginx" };
const deployment = new k8s.apps.v1.Deployment("nginx", {
    metadata: {
        // step 2 has changed values that are expected to be ignored
        creationTimestamp: "invalid-step2",
        deletionTimestamp: "invalid-step2",
        resourceVersion: "invalid-step2",
        uid: "invalid-step2",
        // and some good properties
        annotations: {
            "foo": "bar"
        }
    },
    spec: {
        selector: { matchLabels: appLabels },
        replicas: 1,
        template: {
            metadata: { labels: appLabels },
            spec: { containers: [{ name: "nginx", image: "nginx" }] }
        }
    }
});
export const name = deployment.metadata.name;
export const resourceVersion = deployment.metadata.resourceVersion;
