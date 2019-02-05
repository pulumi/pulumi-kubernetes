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

// Create a test namespace to allow test parallelism.
const namespace = new k8s.core.v1.Namespace("test-namespace");

function addNamespace(o: any) {
  if (o !== undefined) {
    if (o.metadata !== undefined) {
      o.metadata.namespace = namespace.metadata.name;
    } else {
      o.metadata = {namespace: namespace.metadata.name}
    }
  }
}

// Create resources from standard Kubernetes guestbook YAML example in the test namespace.
new k8s.yaml.ConfigFile("guestbook", {
  file: "https://raw.githubusercontent.com/pulumi/pulumi-kubernetes/master/examples/yaml-guestbook/yaml/guestbook.yaml",
  transformations: [addNamespace]
});
