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

import * as k8s from '@pulumi/kubernetes'

const k8sProvider = new k8s.Provider(`k8s-provider`, {
    kubeconfig: '~/.kube/config',
})

const namespace = new k8s.core.v1.Namespace("release-ns");

new k8s.helm.v3.Chart(
    'wordpress',
    {
        fetchOpts: {
            repo: 'https://charts.bitnami.com/bitnami',
        },
        namespace: namespace.metadata.name,
        chart: 'wordpress',
        values: {
            "service": {"type": "ClusterIP"},
        },
    },
    {
        provider: k8sProvider,
    }
)
