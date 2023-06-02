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

// This test validates the following restrictions enforced by "strict mode":
// 1. Default providers are not allowed.
// 2. Each Provider requires a "context" argument.
// 3. Each Provider requires a "kubeconfig" argument.

// Create a new provider with no context specified.
const provider = new k8s.Provider("missingContext", {
    kubeconfig: "~/.kube/config",
});

// The ConfigMap should not be created since the Provider is invalid under "strict mode".
const cm = new k8s.core.v1.ConfigMap("default", {
    data: {foo: "bar"},
}, {provider});
