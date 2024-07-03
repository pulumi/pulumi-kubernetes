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

const namespace = new k8s.core.v1.Namespace("test-namespace");

//
// Delete the CustomResourceDefinition. On the next refresh, the CustomResource will be
// automatically deleted from the state rather than returning a "kind not found" error.
// This is a contrived example for testing, and it doesn't make any sense to delete a
// CRD while leaving related CRs.
//

new k8s.apiextensions.CustomResource(
    "my-new-foobar-object",
    {
        apiVersion: "stable.example.com/v1",
        kind: "FooBar",
        metadata: {
            namespace: namespace.metadata.name,
            name: "my-new-foobar-object",
        },
        spec: { foo: "such amaze" }
    },
);
