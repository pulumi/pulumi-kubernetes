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

import * as pulumi from "@pulumi/pulumi";
import * as k8s from '@pulumi/kubernetes'

const namespace = new k8s.core.v1.Namespace("release-ns");

const k8sProvider = new k8s.Provider(`k8s-provider`, {namespace: namespace.metadata.name})

const alertManager = new k8s.helm.v3.Release("alertmanager", {
    name: "alertmanager",
    chart: "alertmanager",
    version: "0.12.2",
    repositoryOpts: {
        repo: "https://prometheus-community.github.io/helm-charts",
    },
}, {provider: k8sProvider});

// Ensure we get the expected namespace for the stateful set.
export const alertManagerNamespace = k8s.apps.v1.StatefulSet.get(
    "alertmanager-statefulset",
    pulumi.interpolate`${alertManager.status.namespace}/alertmanager`).metadata.namespace;

export const providerNamespace = k8sProvider.namespace;
