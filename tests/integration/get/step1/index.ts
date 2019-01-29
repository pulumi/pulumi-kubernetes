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

const namespace = new k8s.core.v1.Namespace("test-namespace");

//
// `get`s the Kubernetes API service.
//

k8s.core.v1.Service.get("kube-api", "default/kubernetes");

//
// Create a CustomResourceDefinition, a CustomResource, and then `.get` it.
//

const ct = new k8s.apiextensions.v1beta1.CustomResourceDefinition("crontab", {
    metadata: { name: "crontabs.stable.example.com" },
    spec: {
        group: "stable.example.com",
        version: "v1",
        scope: "Namespaced",
        names: {
            plural: "crontabs",
            singular: "crontab",
            kind: "CronTab",
            shortNames: ["ct"]
        }
    }
});

new k8s.apiextensions.CustomResource(
    "my-new-cron-object",
    {
        apiVersion: "stable.example.com/v1",
        kind: "CronTab",
      metadata: {
        namespace: namespace.metadata.apply(ns => ns.name),
        name: "my-new-cron-object",
      },
        spec: { cronSpec: "* * * * */5", image: "my-awesome-cron-image" }
    },
    { dependsOn: ct }
);
