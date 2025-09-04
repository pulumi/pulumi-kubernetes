import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import { FileAsset } from "@pulumi/pulumi/asset";
import * as random from "@pulumi/random";

const password = new random.RandomPassword("password", {
    length: 16,
    special: true,
});
const namespace = new k8s.core.v1.Namespace("release-ns", {});

const releaseName = new random.RandomString("name", {length: 8, special: false, upper: false, numeric: false});

const release = new k8s.helm.v3.Release("release", {
    allowNullValues: true,
    chart: "redis",
    repositoryOpts: {
        repo: "https://raw.githubusercontent.com/bitnami/charts/eb5f9a9513d987b519f0ecd732e7031241c50328/bitnami",
    },
    version: "13.0.0",
    namespace: namespace.metadata.name,
    name: releaseName.result,
    valueYamlFiles: [new FileAsset("./metrics.yml")],
    values: {
        cluster: {
            enabled: true,
            slaveCount: 2,
        },
        image: {
            repository: "bitnamilegacy/redis",
            tag: "8.2.1",
        },
        global: {
            redis: {
                password: password.result,
            }
        },
        rbac: {
            create: true,
        }
    },
});


const srv = k8s.core.v1.Service.get("redis-master-svc", pulumi.interpolate`${release.status.namespace}/${release.status.name}-redis-master`);
export const redisMasterClusterIP = srv.spec.clusterIP;
export const helmStatus = release.status.status;
export const helmDescription = release.description;