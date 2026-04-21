import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

const release = new k8s.helm.v3.Release("null-test", {
    chart: "./chart",
    namespace: "default",
    values: {},
});

const cm = k8s.core.v1.ConfigMap.get("cm",
    pulumi.interpolate`default/${release.status.name}-cm`);

export const configMapData = cm.data;
