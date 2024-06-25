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

import * as k8s from "@pulumi/kubernetes";
import * as pulumi from "@pulumi/pulumi";

let config = new pulumi.Config();
let renderDir: string = config.require("renderDir");

const provider = new k8s.Provider("render-yaml", {
    renderYamlToDirectory: renderDir,
});

const appLabels = {app: "nginx"};
const deployment = new k8s.apps.v1.Deployment("nginx", {
    spec: {
        selector: {matchLabels: appLabels},
        replicas: 1,
        template: {
            metadata: {labels: appLabels},
            spec: {containers: [{name: "nginx-fake", image: "nginx-fake", ports: [{containerPort: 80}]}]}
        }
    }
}, {provider});
const service = new k8s.core.v1.Service("nginx", {
    spec: {
        ports: [{port: 8080, protocol: "TCP"}],
        selector: deployment.metadata.labels,
    }
}, {provider});
