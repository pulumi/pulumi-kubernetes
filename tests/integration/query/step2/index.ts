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

import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

//
// Test a query that retrieves only Pods.
//

pulumi.runtime
    .listResourceOutputs(k8s.core.v1.Pod.isInstance, `moolumi/${pulumi.runtime.getStack()}`)
    .toArray()
    .then(pods => {
        //
        // Test that there is exactly 1 Pod. This query should filter out the stack, provider, and
        // namespace resources.
        //
        if (pods.length !== 1) {
            throw Error("Expected 1 pod, got " + pods.length);
        }

        //
        // Test that query is well-typed. `pod[0]` should be a Pod, and the compiler should let us
        // access its fields.
        //
        if (pods[0].kind !== "Pod") {
            throw Error("Expected Pods to have `.kind === 'Pod'`, got " + pods[0].kind);
        }
    });
