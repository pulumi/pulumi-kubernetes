// Copyright 2016-2020, Pulumi Corporation.
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

// Create two test namespaces to allow test parallelism.
const namespace = new k8s.core.v1.Namespace("test-namespace");
const namespace2 = new k8s.core.v1.Namespace("test-namespace2");

function kustomizeDirectory(name: string, directory: string, namespace: string, resourcePrefix?: string): k8s.kustomize.Directory {
    return new k8s.kustomize.Directory(name, {
        directory,
        resourcePrefix,
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
            (obj: any) => {
                if (obj.kind === "Service") {
                    obj.spec.type = "ClusterIP"
                }
            }
        ]
    });
}

// Create resources from local kustomize directory in the first namespace.
namespace.metadata.name.apply(ns => kustomizeDirectory("helloWorld", "helloWorld", ns));

// Create resources from remote kustomize directory in the second namespace.
// Disambiguate resource names with a specified prefix.
namespace2.metadata.name.apply(ns => kustomizeDirectory(
    "helloWorld",
    "https://github.com/kubernetes-sigs/kustomize/tree/master/examples/helloWorld",
    ns,
    "remote")
);
