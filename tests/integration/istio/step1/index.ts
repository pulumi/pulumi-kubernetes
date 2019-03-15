import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

import { k8sProvider, k8sConfig, k8sCluster } from "./cluster";
import { istio } from "./istio";

new k8s.core.v1.Namespace(
    "bookinfo",
    { metadata: { name: "bookinfo", labels: { "istio-injection": "enabled" } } },
    { provider: k8sProvider }
);

function addNamespace(o: any) {
    if (o !== undefined) {
        o.metadata.namespace = "bookinfo";
    }
}

const bookinfo = new k8s.yaml.ConfigFile(
    "yaml/bookinfo.yaml",
    { transformations: [addNamespace] },
    { dependsOn: istio, providers: { kubernetes: k8sProvider } }
);

new k8s.yaml.ConfigFile(
    "yaml/bookinfo-gateway.yaml",
    { transformations: [addNamespace] },
    { dependsOn: bookinfo, providers: { kubernetes: k8sProvider } }
);

export const port = istio
    .getResourceProperty("v1/Service","istio-system", "istio-ingressgateway", "spec")
    .apply(p => p.ports.filter(p => p.name == "http2"));

export const frontendIp = pulumi
    .all([
        istio.getResourceProperty("v1/Service", "istio-system", "istio-ingressgateway", "status"),
        istio.getResourceProperty("v1/Service", "istio-system", "istio-ingressgateway", "spec")
    ])
    .apply(([status, spec]) => {
        const port = spec.ports.filter(p => p.name == "http2")[0].port;
        return `${status.loadBalancer.ingress[0].ip}:${port}/productpage`;
    });

export const kubeconfig = k8sConfig;
