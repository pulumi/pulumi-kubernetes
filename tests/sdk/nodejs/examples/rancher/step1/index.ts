import * as k8s from "@pulumi/kubernetes";

const namespace = new k8s.core.v1.Namespace("release-ns");

const release = new k8s.helm.v3.Release("rancher-release", {
    chart: "rancher",
    repositoryOpts: {
        repo: "https://releases.rancher.com/server-charts/latest",
    },
    version: "v2.6.2",
    namespace: namespace.metadata.name,
    values: {
        hostname: "rancher.examples.pulumi.com",
        ingress: {
            tls: {
                source: "secret",
            },
        },
    },
    skipAwait: true,
});

