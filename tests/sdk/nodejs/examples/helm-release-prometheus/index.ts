import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";


const namespace = new k8s.core.v1.Namespace("release-ns");

const prometheusOperator = new k8s.helm.v3.Release("prometheus", {
    name: "prometheus-operator",
    namespace: namespace.metadata.name,
    chart: "prometheus-operator",
    version: "8.5.2",
    repositoryOpts: {
        repo: "https://charts.helm.sh/stable",
    },
    values: {},
});

const rule = k8s.apiextensions.CustomResource.get("prom-rule",
    {
        apiVersion: "monitoring.coreos.com/v1",
        kind: "PrometheusRule",
        id: pulumi.interpolate `${prometheusOperator.status.namespace}/${prometheusOperator.status.name}-prometheus-operator`
    });

export const metadata = rule.metadata;
