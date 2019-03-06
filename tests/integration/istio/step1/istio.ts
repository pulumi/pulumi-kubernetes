import { execSync } from "child_process";

import * as k8s from "@pulumi/kubernetes";

import * as config from "./config";
import { k8sProvider } from "./cluster";

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

export const istio = new k8s.helm.v2.Chart(
    appName,
    {
        chart: "istio",
        namespace: namespace.metadata.name,
        version: "1.0.1",
        fetchOpts: { repo: "https://istio.io/charts/" },
        // for all options check https://github.com/istio/istio/tree/master/install/kubernetes/helm/istio
        values: { kiali: { enabled: true } }
    },
    { dependsOn: [namespace, adminBinding], providers: { kubernetes: k8sProvider } }
);
