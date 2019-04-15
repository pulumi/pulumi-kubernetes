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
import * as fs from "fs";
import * as os from "os";
import * as path from "path";

// Use the existing ~/.kube/config kubeconfig
const kubeconfig = fs.readFileSync(path.join(os.homedir(), ".kube", "config")).toString();

const ns = new k8s.core.v1.Namespace("ns");

// Create a new provider
const myk8s = new k8s.Provider("myk8s", {
    kubeconfig: kubeconfig,
    namespace: ns.metadata.name,
});

// Create a Pod using the custom provider.
// The namespace should be automatically set by the provider override.
new k8s.core.v1.Pod("nginx", {
    spec: {
        containers: [{
            image: "nginx:1.7.9",
            name: "nginx",
            ports: [{ containerPort: 80 }],
        }],
    },
}, { provider: myk8s });

// Create a Pod using the custom provider with a specified namespace.
// The namespace should be overridden by the provider override.
new k8s.core.v1.Pod("namespaced-nginx", {
    metadata: { namespace: ns.metadata.name },
    spec: {
        containers: [{
            image: "nginx:1.7.9",
            name: "nginx",
            ports: [{ containerPort: 80 }],
        }],
    },
}, { provider: myk8s });

// Create a Namespace using the custom provider
// The namespace should not be affected by the provider override since it is a non-namespaceable kind.
new k8s.core.v1.Namespace("other-ns",
    {},
    { provider: myk8s });
