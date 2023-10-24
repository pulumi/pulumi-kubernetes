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

const provider = new k8s.Provider("k8s", {enableServerSideApply: true});

// Update the patch to reduce replicas back to 1, and also drop specifying the 
// app labels. We drop labels here to simulate a strategic merge patch where
// we just want to update the replicas field, and not update the app labels. However,
// this is not possible as Pulumi Kubernetes Patch resources use SSA so this step
// WILL fail since SSA will attempt to unset the app label fields (.spec.template.metadata.labels
// and .spec.selector) since they are not specified in the SSA patch.
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
      replicas: 1,
    },
  },
  { provider }
);