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

// Create provider with SSA enabled.
const provider = new k8s.Provider("k8s", {enableServerSideApply: true});

const ns = new k8s.core.v1.Namespace("test-ignore-changes", undefined, { provider });

const deployment = new k8s.apps.v1.Deployment(
  "test-ignore-changes",
  {
    metadata: {
      namespace: ns.metadata.name,
    },
    spec: {
      selector: { matchLabels: { app: "test-ignore-changes" } },
      replicas: 2,
      template: {
        metadata: {
          labels: { app: "test-ignore-changes" },
        },
        spec: {
          containers: [
            {
              name: "nginx",
              image: "nginx:1.25",
            },
          ],
        },
      },
    },
  },
  { provider: provider, ignoreChanges: ["spec.replicas"] }
);

export const deploymentName = deployment.metadata.name;
export const deploymentNamespace = deployment.metadata.namespace;