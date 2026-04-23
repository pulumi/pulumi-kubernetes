import * as k8s from "@pulumi/kubernetes";

const appLabels = { app: "null-test" };

const deployment = new k8s.apps.v1.Deployment("null-test", {
    metadata: {
        name: "null-test",
        namespace: "default",
    },
    spec: {
        replicas: null as any,
        selector: { matchLabels: appLabels },
        template: {
            metadata: { labels: appLabels },
            spec: {
                containers: [{
                    name: "pause",
                    image: "registry.k8s.io/pause:3.9",
                }],
            },
        },
    },
});

export const replicas = deployment.spec.replicas;
