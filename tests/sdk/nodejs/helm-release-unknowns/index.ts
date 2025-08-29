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
    chart: "ingress-nginx",
    repositoryOpts: {
        repo: "https://kubernetes.github.io/ingress-nginx",
    },
    version: "4.13.2",
    namespace: namespace.metadata.name,
    name: releaseName.result,
    valueYamlFiles: [new FileAsset("./metrics.yml")],
    values: {
        cluster: {
            enabled: true,
            slaveCount: 2,
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
