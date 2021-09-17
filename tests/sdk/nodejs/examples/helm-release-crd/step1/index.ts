import * as k8s from "@pulumi/kubernetes";

const prometheusRelease = new k8s.helm.v3.Release("prometheus", {
    name: "kube-prometheus-stack",
    chart: "kube-prometheus-stack",
    repositoryOpts: {
        repo: "https://prometheus-community.github.io/helm-charts",
    },
    values: {},
});

