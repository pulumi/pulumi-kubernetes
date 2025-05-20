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
import * as random from "@pulumi/random"

//
// To allow parallel test runs, generate a namespace name with a random suffix.
//

let randomSuffix = new random.RandomString("random-suffix", {
    length: 6,
    special: false,
    upper: false
});
let nsName = pulumi.concat(`test-namespace-`, randomSuffix.result);

//
// Tests that if we force a `Pod` to be created before the `Namespace` it is supposed to exist in,
// it will retry until created.
//

new k8s.core.v1.Pod("nginx", {
    metadata: {
        namespace: nsName,
        name: "nginx"
    },
    spec: {
        containers: [
            {
                name: "nginx",
                image: "nginx:1.27.2",
                ports: [{ containerPort: 80 }]
            }
        ]
    }
});

new k8s.core.v1.Namespace("test", {
    metadata: {
        // Wait 10 seconds before creating the namespace, to make sure we retry the Pod creation.
        annotations: {
            timeout: new Promise(resolve => {
                if (pulumi.runtime.isDryRun()) {
                    return resolve("<output>");
                }
                setTimeout(() => resolve("done"), 10000);
            })
        },
        name: nsName
    }
});
