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
import * as pulumi from "@pulumi/pulumi";

import { k8sProvider } from "./cluster";
import { crd10, crd11, crd12, istio } from "./istio";

new k8s.core.v1.Namespace(
    "bookinfo",
    { metadata: { name: "bookinfo", labels: { "istio-injection": "enabled" } } },
    { provider: k8sProvider }
);

function addNamespace(o: any) {
    if (o !== undefined) {
        o.metadata.namespace = "bookinfo";
    }
}

const bookinfo = new k8s.yaml.ConfigFile(
    "yaml/bookinfo.yaml",
    { transformations: [addNamespace] },
    { dependsOn: [crd10, crd11, crd12], providers: { kubernetes: k8sProvider } }
);

new k8s.yaml.ConfigFile(
    "yaml/bookinfo-gateway.yaml",
    { transformations: [addNamespace] },
    { dependsOn: [crd10, crd11, crd12], providers: { kubernetes: k8sProvider } }
);

export const port = istio
    .getResourceProperty("v1/Service","istio-system", "istio-ingressgateway", "spec")
    .apply(p => p.ports.filter(p => p.name == "http2"));

export const frontendIp = pulumi
    .all([
        istio.getResourceProperty("v1/Service", "istio-system", "istio-ingressgateway", "status"),
        istio.getResourceProperty("v1/Service", "istio-system", "istio-ingressgateway", "spec")
    ])
    .apply(([status, spec]) => {
        const port = spec.ports.filter(p => p.name == "http2")[0].port;
        return `${status.loadBalancer.ingress[0].ip}:${port}/productpage`;
    });
