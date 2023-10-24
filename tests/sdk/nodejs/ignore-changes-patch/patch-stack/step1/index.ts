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
import * as pulumi from "@pulumi/pulumi";

// Create provider with SSA enabled.
let config = new pulumi.Config();
const ns = config.require("DEPLOYMENT_NAMESPACE");
const name = config.require("DEPLOYMENT_NAME");

const provider = new k8s.Provider("k8s", { enableServerSideApply: true });

// Create a patch that changes the number of replicas for the deployment created in
// the deployment-stack example. We are also setting the app labels here to simulate
// shared ownership of the app labels between the deployment stack and this set.
const appLabels = { app: "test-ignore-changes" };
const patch = new k8s.apps.v1.DeploymentPatch(
  "test-ignore-changes-patch",
  {
    metadata: {
      namespace: ns,
      name: name,
      annotations: {
        "pulumi.com/patchForce": "true",
      },
    },
    spec: {
      selector: { matchLabels: appLabels },
      replicas: 2,
      template: {
        metadata: {
          labels: appLabels,
        },
      },
    },
  },
  { provider }
);
