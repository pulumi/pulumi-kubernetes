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

const provider = new k8s.Provider("k8s");

// Create test namespace to allow test parallelism.
const namespace = new k8s.core.v1.Namespace("test-namespace", {}, {provider});

new k8s.kustomize.Directory("moria", {
    directory: "moria/base",
    transformations: [
        (obj: any) => {
            if (obj !== undefined) {
                if (obj.metadata !== undefined) {
                    obj.metadata.namespace = namespace;
                } else {
                    obj.metadata = {namespace: namespace};
                }
            }
        },
    ]
}, {provider});
