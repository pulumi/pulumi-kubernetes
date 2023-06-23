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

// This test creates a Provider with `enableServerSideApply` enabled. The following scenarios are tested:
// 1. Create a Namespace and ConfigMap with pulumi.
// 2. Externally delete labels in the ConfigMap using kubectl.
// 3. Rerun the pulumi program and verify that the labels are restored.

// Create provider with SSA enabled.
const provider = new k8s.Provider("k8s", {
    enableServerSideApply: false,
    enableConfigMapMutable: true, // Reuse the same ConfigMap rather than replacing on change.
});

// Create a randomly-named Namespace.
const ns = new k8s.core.v1.Namespace("test", undefined, {provider});

export const cm = new k8s.core.v1.ConfigMap("test", {
    metadata: {
        namespace: ns.metadata.name,
        annotations: {
            "pulumi.com/patchForce": "true",
        },
    },
    data: {foo: "bar"},
}, {provider});
