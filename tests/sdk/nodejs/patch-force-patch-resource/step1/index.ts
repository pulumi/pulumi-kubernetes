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
const provider = new k8s.Provider("k8s", {
    enableServerSideApply: true,
    kubeconfig: new pulumi.Config().require("kubeconfig"),
});

// Do not apply the patchForce annotation, and this step will fail.
const patch = new k8s.apps.v1.DaemonSetPatch("kube-proxy-image", {
    metadata: {
        name: "kube-proxy",
        namespace: "kube-system",
    },
    spec: {
        template: {
            spec: {
                containers: [
                    {
                        name: "kube-proxy",
                        image: "registry.k8s.io/kube-proxy:v1.27.1",
                        command: [
                            "/usr/local/bin/kube-proxy",
                            "--config=/var/lib/kube-proxy/config.conf",
                            "--hostname-override=$(NODE_NAME)",
                            "--v=2",
                        ],
                    },
                ],
            },
        },
    },
}, { provider, retainOnDelete: true });