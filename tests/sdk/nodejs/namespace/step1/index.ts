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

///
/// Create a test namespace with no specified metadata.
///

new k8s.core.v1.Namespace("test");

///
/// Create a test Pod with no metadata specified. This will be created in the "default" namespace,
/// but the object registered with Pulumi to create will not have the .metadata.namespace field set.
///

new k8s.core.v1.Pod("no-metadata-pod", {
    spec: {
        containers: [
            {
                name: "nginx",
                image: "nginx:1.27.2",
                ports: [{containerPort: 80}]
            }
        ]
    }
});

///
/// Create a test Pod with the "default" namespace explicitly set.
///

new k8s.core.v1.Pod("default-ns-pod", {
    metadata: {
        namespace: "default"
    },
    spec: {
        containers: [
            {
                name: "nginx",
                image: "nginx:1.27.2",
                ports: [{containerPort: 80}]
            }
        ]
    }
});
