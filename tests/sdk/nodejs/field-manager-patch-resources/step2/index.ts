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
const provider = new k8s.Provider("k8s", { enableServerSideApply: true });

const config = new pulumi.Config();
const namespace = config.require("namespace");

const depPatch = new k8s.apps.v1.DeploymentPatch(
  "nginx-patch",
  {
    metadata: {
      namespace: namespace,
      name: "test-mgr-nginx",
      annotations: {
        "pulumi.com/patchForce": "true",
      },
    },
    spec: {
        template: {
            spec: {
                containers: [
                    {
                        name: "nginx",
                        image: "nginx:1.14.0",
                    },
                ],
            },
    }
  },
}, { provider: provider, retainOnDelete: true, ignoreChanges: ["spec.template.metadata", "spec.selector"]});