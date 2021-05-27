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
import { k8sProvider } from "./cluster";

import * as config from "./config";

const appName = "istio";
const namespace = new k8s.core.v1.Namespace(
    `${appName}-system`,
    { metadata: { name: `${appName}-system` } },
    { provider: k8sProvider }
);

const adminBinding = new k8s.rbac.v1.ClusterRoleBinding(
    "cluster-admin-binding",
    {
        metadata: { name: "cluster-admin-binding" },
        roleRef: {
            apiGroup: "rbac.authorization.k8s.io",
            kind: "ClusterRole",
            name: "cluster-admin"
        },
        subjects: [
            { apiGroup: "rbac.authorization.k8s.io", kind: "User", name: config.gcpUsername }
        ]
    },
    { provider: k8sProvider }
);

export const istio_init = new k8s.helm.v2.Chart(
    `${appName}-init`,
    {
        path: "charts/istio-init",
        namespace: namespace.metadata.name,
        values: { kiali: { enabled: true } }
    },
    { dependsOn: [namespace, adminBinding], providers: { kubernetes: k8sProvider } }
);

export const crd10 = namespace.metadata.name.apply(ns => istio_init.getResource("batch/v1/Job", ns, "istio-init-crd-10"));
export const crd11 = namespace.metadata.name.apply(ns => istio_init.getResource("batch/v1/Job", ns, "istio-init-crd-11"));
export const crd12 = namespace.metadata.name.apply(ns => istio_init.getResource("batch/v1/Job", ns, "istio-init-crd-12"));

export const istio = new k8s.helm.v2.Chart(
    appName,
    {
        path: "charts/istio",
        namespace: namespace.metadata.name,
        values: { kiali: { enabled: true } }
    },
    { dependsOn: [adminBinding, crd10, crd11, crd12], providers: { kubernetes: k8sProvider } }
);
