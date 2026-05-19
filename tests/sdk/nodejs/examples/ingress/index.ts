// Copyright 2016-2026, Pulumi Corporation.
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

const ns = new k8s.core.v1.Namespace("test");
const namespace = ns.metadata.name;

const helloLabels = { app: "hello", tier: "frontend" };

const helloService = new k8s.core.v1.Service("hello-svc", {
    metadata: {
        namespace: namespace,
        name: "hello",
        labels: helloLabels,
    },
    spec: {
        ports: [{ port: 8080, targetPort: 8080 }],
        selector: helloLabels,
    },
});

new k8s.apps.v1.Deployment("hello-app", {
    metadata: {
        namespace: namespace,
        name: "hello",
    },
    spec: {
        selector: { matchLabels: helloLabels },
        replicas: 1,
        template: {
            metadata: { labels: helloLabels },
            spec: {
                containers: [{
                    name: "hello",
                    image: "gcr.io/google-samples/hello-app:1.0",
                    ports: [{ containerPort: 8080 }],
                }],
            },
        },
    },
});

const feIngress = new k8s.networking.v1.Ingress("feIngress", {
    metadata: {
        namespace: namespace,
        name: "feingress",
    },
    spec: {
        ingressClassName: "traefik",
        rules: [{
            http: {
                paths: [{
                    path: "/hello",
                    pathType: "Prefix",
                    backend: { service: { name: helloService.metadata.name, port: { number: 8080 } } },
                }],
            },
        }],
    },
});

export const ingressIp = feIngress.status.loadBalancer.ingress[0].ip;
