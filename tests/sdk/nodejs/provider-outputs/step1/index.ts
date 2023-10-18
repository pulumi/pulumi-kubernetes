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
import * as fs from "fs";
import * as os from "os";
import * as path from "path";

const ns1 = new k8s.core.v1.Namespace("ns1");

// Create a new provider using the current context.
const k8s1 = new k8s.Provider("k8s1", {
    kubeconfig: "./kubeconfig",
});

// Create a new provider using an overridden namespace.
const k8s2 = new k8s.Provider("k8s2", {
    kubeconfig: "./kubeconfig",
    namespace: ns1.metadata.name,
});

// Create a new provider using an overridden context.
const k8s3 = new k8s.Provider("k8s3", {
    kubeconfig: "./kubeconfig",
    context: "context2"
});

// Create a new provider using an overridden cluster.
const k8s4 = new k8s.Provider("k8s4", {
    kubeconfig: "./kubeconfig",
    cluster: "cluster2"
});

export const ns1Name = ns1.metadata.name

export const k8s1Namespace = k8s1.namespace
export const k8s1Context = k8s1.context
export const k8s1Config = k8s1.kubeconfig
export const k8s1Cluster = k8s1.cluster

export const k8s2Namespace = k8s2.namespace
export const k8s2Context = k8s2.context
export const k8s2Config = k8s2.kubeconfig
export const k8s2Cluster = k8s2.cluster

export const k8s3Namespace = k8s3.namespace
export const k8s3Context = k8s3.context
export const k8s3Config = k8s3.kubeconfig
export const k8s3Cluster = k8s3.cluster

export const k8s4Namespace = k8s4.namespace
export const k8s4Context = k8s4.context
export const k8s4Config = k8s4.kubeconfig
export const k8s4Cluster = k8s4.cluster