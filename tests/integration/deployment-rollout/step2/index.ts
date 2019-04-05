// Copyright 2016-2019, Pulumi Corporation.
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

export const namespace = new k8s.core.v1.Namespace("test-namespace");

//
// Change the image to trigger an update.
//

const appLabels = {app: "nginx"};
new k8s.extensions.v1beta1.Deployment("nginx", {
  metadata: { namespace: namespace.metadata.name },
  spec: {
    selector: {matchLabels: appLabels},
    replicas: 1,
    template: {
      metadata: {labels: appLabels},
      spec: {containers: [{name: "nginx", image: "nginx:stable", ports: [{containerPort: 80}]}]}
    }
  }
});

//
// Change the image to trigger an update.
//

new k8s.apps.v1.Deployment("nginx", {
  metadata: { namespace: namespace.metadata.name },
  spec: {
    selector: {matchLabels: appLabels},
    replicas: 1,
    template: {
      metadata: {labels: appLabels},
      spec: {containers: [{name: "nginx", image: "nginx:stable", ports: [{containerPort: 80}]}]}
    }
  }
});

