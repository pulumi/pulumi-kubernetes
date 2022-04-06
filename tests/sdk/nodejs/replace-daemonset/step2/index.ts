import * as k8s from "@pulumi/kubernetes";

const appLabels = {
    app: "nginx",
    test: "new-value",
};

const daemonset = new k8s.apps.v1.DaemonSet("test-replacement", {
    spec: {
        selector: { matchLabels: appLabels },
        template: {
            metadata: { labels: appLabels },
            spec: { containers: [{ name: "nginx", image: "nginx" }] },
        },
    },
});

export const name = daemonset.metadata.name;
