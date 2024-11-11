// Copyright 2016-2020, Pulumi Corporation.
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
import * as fs from "fs";
import * as os from "os";
import * as path from "path";

const ns1 = new k8s.core.v1.Namespace("ns1");
const ns2 = new k8s.core.v1.Namespace("ns2");

const kcfg = process.env.KUBECONFIG|| path.join(os.homedir(), ".kube", "config");

// Create a new provider using the contents of a k8s config.
const kubeconfigContentsProvider = new k8s.Provider("kubeconfigContentsProvider", {
    kubeconfig: fs.readFileSync(kcfg).toString(),
    namespace: ns1.metadata.name,
});

// Create a new provider using the path to a k8s config.
const kubeconfigPathProvider = new k8s.Provider("kubeconfigPathProvider", {
    kubeconfig: kcfg,
    namespace: ns1.metadata.name,
});

// Create a Pod using the contents provider.
// The namespace should be automatically set by the provider default.
new k8s.core.v1.Pod("nginx1", {
    spec: {
        containers: [{
            image: "nginx:1.27.2",
            name: "nginx",
            ports: [{ containerPort: 80 }],
        }],
    },
}, { provider: kubeconfigContentsProvider });

// Create a Pod using the path provider.
// The namespace should be automatically set by the provider default.
new k8s.core.v1.Pod("nginx2", {
    spec: {
        containers: [{
            image: "nginx:1.27.2",
            name: "nginx",
            ports: [{ containerPort: 80 }],
        }],
    },
}, { provider: kubeconfigPathProvider });

// Create a Pod using the contents provider with a specified namespace.
// The namespace should not be overridden by the provider default.
new k8s.core.v1.Pod("namespaced-nginx", {
    metadata: { namespace: ns2.metadata.name },
    spec: {
        containers: [{
            image: "nginx:1.27.2",
            name: "nginx",
            ports: [{ containerPort: 80 }],
        }],
    },
}, { provider: kubeconfigContentsProvider });

// Create a Namespace using the contents provider
// The namespace should not be affected by the provider override since it is a cluster-scoped kind.
new k8s.core.v1.Namespace("other-ns",
    {},
    { provider: kubeconfigContentsProvider });
