import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import * as random from "@pulumi/random";


const redisPassword = pulumi.secret("$053cr3t!");

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
    version: "13.0.0",
    namespace: namespace.metadata.name,
    values: {
        cluster: {
            enabled: true,
            slaveCount: 2,
        },
        metrics: {
            enabled: true,
            service: {
                annotations: {
                    "prometheus.io/port": "9127",
                }
            },
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
export const status = release.status.status;
