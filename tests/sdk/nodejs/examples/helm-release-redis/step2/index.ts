import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import * as random from "@pulumi/random";

const redisPassword = new random.RandomPassword("password", {length: 16}); // Change password

const nsName = new random.RandomPet("test");
const namespace = new k8s.core.v1.Namespace("release-ns", {
    metadata: {
        name: nsName.id
    }
});

// Validates fix for https://github.com/pulumi/pulumi-kubernetes/issues/1933
function values(password: pulumi.Output<string>): pulumi.Input<{ [key: string]: any }> {
    return pulumi.output({
        cluster: {
            enabled: true,
            slaveCount: 1,
        },
        image: {
            repository: "bitnamilegacy/redis",
            tag: "8.2.1",
        },
        global: {
            redis: {
                password: password,
            }
        },
        rbac: {
            create: true,
        },
        service: {
            type: "ClusterIP"
        }
    });
};

const release = new k8s.helm.v3.Release("redis", {
    chart: "redis",
    repositoryOpts: {
        repo: "https://raw.githubusercontent.com/bitnami/charts/eb5f9a9513d987b519f0ecd732e7031241c50328/bitnami",
    },
    namespace: namespace.metadata.name,
    values: values(redisPassword.result),
    verify: false, // Turn off verification explicitly.
});


const srv = k8s.core.v1.Service.get("redis-master-svc", pulumi.interpolate`${release.status.namespace}/${release.status.name}-master`);
export const redisMasterClusterIP = srv.spec.clusterIP;
export const status = release.status.status;
