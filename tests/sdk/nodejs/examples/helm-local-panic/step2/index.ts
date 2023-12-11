import * as k8s from "@pulumi/kubernetes";

const namespace = new k8s.core.v1.Namespace("release-ns");

const release = new k8s.helm.v3.Release("local-chart-panic", {
    chart: "./local-test-chart",
    namespace: namespace.metadata.name,
    version: "0.1.1", // <- Updating chart version to trigger an upgrade.
});

