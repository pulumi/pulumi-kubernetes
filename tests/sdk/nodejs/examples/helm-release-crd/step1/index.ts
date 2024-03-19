import * as k8s from "@pulumi/kubernetes";

const postgresOperator = new k8s.helm.v3.Release("postgres-operator", {
    name: "postgres-operator",
    chart: "postgres-operator",
    version: "1.11.0",
    repositoryOpts: {
        repo: "https://opensource.zalando.com/postgres-operator/charts/postgres-operator/",
    },
    values: {},
});
