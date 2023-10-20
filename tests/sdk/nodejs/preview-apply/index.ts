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
const provider = new k8s.Provider("k8s", { enableServerSideApply: false });

const ns = new k8s.core.v1.Namespace("test-preview-apply", undefined, {
  provider,
});

const dep = new k8s.apps.v1.Deployment(
  "nginx-dep",
  {
    metadata: {
      namespace: ns.metadata.name,
      labels: {
        app: "nginx",
      },
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          app: "nginx",
        },
      },
      template: {
        metadata: {
          labels: {
            app: "nginx",
          },
        },
        spec: {
          containers: [
            {
              name: "nginx",
              image: "nginx:latest",
              ports: [
                {
                  containerPort: 80,
                },
              ],
            },
          ],
        },
      },
    },
  },
  { provider }
);

const svc = new k8s.core.v1.Service(
  "nginx-svc",
  {
    metadata: {
      namespace: ns.metadata.name,
      labels: {
        app: "nginx",
      },
    },
    spec: {
      type: "LoadBalancer",
      ports: [
        {
          port: 80,
          targetPort: 80,
        },
      ],
      selector: {
        app: "nginx",
      },
    },
  },
  { provider }
);

export const ip = svc.status.apply((s) => s.loadBalancer.ingress[0].ip);
export const nsName = ns.metadata.name;
export const svcName = svc.metadata.name;
