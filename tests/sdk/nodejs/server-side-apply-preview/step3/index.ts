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
// 1. Create a namespace and ConfigMap with pulumi.
// 2. Change the provider context to use a ServiceAccount with view-only permission.
// 3. Change the ConfigMap and run a preview to confirm that the provider falls back to a Client-side preview when
//    the user does not have permission to perform the "patch" operation used for SSA diff.
// 4. Change the provider context back to default to allow the update to succeed and clean up.

// Create provider with SSA enabled.
const provider = new k8s.Provider("k8s", {
    enableServerSideApply: true,
    enableConfigMapMutable: true,
    context: "kubeconfig-sa", // Set the context to the "view"-only ServiceAccount user.
});

// Create a randomly-named Namespace.
const ns = new k8s.core.v1.Namespace("test", undefined, {provider});

export const cm = new k8s.core.v1.ConfigMap("test", {
    metadata: {
        name: "foo",
        namespace: ns.metadata.name,
    },
    data: {dataKey: "updated"}, // Update the data.
}, {provider});
