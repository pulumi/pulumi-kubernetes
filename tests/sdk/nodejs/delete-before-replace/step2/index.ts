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
// Cause hard delete-before-replace. Changing the Pod's container image tag, causes it to be
// replaced, and since we've manually specified a name, the engine has no choice but to delete it
// first, since names are unique, and it can't generate the name by itself.
//

const pod = new k8s.core.v1.Pod("pod-test", {
  metadata: {
    namespace: namespace.metadata.name,
    name: "pod-test",
  },
  spec: {
    containers: [
      {name: "nginx", image: "nginx:1.15-alpine"},
    ],
  },
});
