import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

// Test for https://github.com/pulumi/pulumi-kubernetes/issues/2997
//
// The local chart defaults config.alpha="default-alpha" and config.beta="default-beta".
// When nullAlpha is true, config.alpha is set to null, which should delete it
// from the merged chart values.

const config = new pulumi.Config();
const nullAlpha = config.getBoolean("nullAlpha") ?? false;

const values: any = {};
if (nullAlpha) {
    values.config = { alpha: null };
}

const release = new k8s.helm.v3.Release("null-test", {
    chart: "./chart",
    namespace: "default",
    values: values,
});

const cm = k8s.core.v1.ConfigMap.get("cm",
    pulumi.interpolate`default/${release.status.name}-cm`);

export const configMapData = cm.data;
