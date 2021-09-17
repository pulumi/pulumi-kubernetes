import * as k8s from "@pulumi/kubernetes";

const dapr = new k8s.helm.v3.Release("dapr", {
    name: "dapr",
    chart: "dapr",
    version: "1.4.0",
    repositoryOpts: {
        repo: "https://dapr.github.io/helm-charts/",
    },
    values: {
        global: {
            prometheus: {
                "enabled": false,
            }
        }
    },
});

