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

// This test creates a Provider with `enableServerSideApply` enabled. The following scenarios are tested:
// 1. Patch a Namespace resource with fully-specified configuration.
// 2. Patch a CustomResource.
// 3. Upsert a Deployment resource that already exists.
// 4. Patch the Deployment with a partially-specified configuration.
// 5. Replace a statically-named ConfigMap resource by changing the data on a subsequent update.
// 6. Ignore changes specified in the ignoreChanges resource option.
// 7. Statically-named Namespace can be changed to an auto-named Namespace.

// Create provider with SSA enabled.
const provider = new k8s.Provider("k8s", {enableServerSideApply: true});

new k8s.core.v1.ConfigMap("test", {
    metadata: {
        name: "foo",
    },
    data: {foo: "bar"},
}, {provider});