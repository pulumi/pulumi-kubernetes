import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import * as random from "@pulumi/random";
import { FileAsset } from "@pulumi/pulumi/asset";

export const redisPassword = pulumi.secret("$053cr3t!");

const nsName = new random.RandomPet("test");
const namespace = new k8s.core.v1.Namespace("release-ns", {
    metadata: {
        name: nsName.id
    }
});

const release = new k8s.helm.v3.Release("release", {
    chart: "redis",
    repositoryOpts: {
        repo: "https://charts.bitnami.com/bitnami",
    },
    version: "13.0.1", // <--- change
    namespace: namespace.metadata.name,
    valueYamlFiles: [new FileAsset("./metrics.yml")],
    values: {
        cluster: {
            enabled: true,
            slaveCount: 2,
        },
        global: {
            redis: {
                password: redisPassword,
            }
        },
        rbac: {
            create: true,
        }
    },
});


const srv = k8s.core.v1.Service.get("redis-master-svc", pulumi.interpolate`${release.status.namespace}/${release.status.name}-redis-master`);
export const redisMasterClusterIP = srv.spec.clusterIP;
export const status = release.status;
