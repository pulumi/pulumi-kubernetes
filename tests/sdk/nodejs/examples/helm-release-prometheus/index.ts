// Copyright 2016-2021, Pulumi Corporation.
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


const namespace = new k8s.core.v1.Namespace("release-ns");

const prometheusOperator = new k8s.helm.v3.Release("prometheus", {
    name: "prometheus-operator",
    namespace: namespace.metadata.name,
    chart: "prometheus-operator",
    version: "8.5.2",
    repositoryOpts: {
        repo: "https://charts.helm.sh/stable",
    },
    values: {},
});

export const namespaceName = namespace.metadata.name
// Explicitly look for the rule in the release's namespace.
export const ruleUrn = k8s.apiextensions.CustomResource.get("prom-rule",
    {
        apiVersion: "monitoring.coreos.com/v1",
        kind: "PrometheusRule",
        id: pulumi.interpolate`${prometheusOperator.status.namespace}/${prometheusOperator.status.name}-prometheus-operator`
    }).urn;
