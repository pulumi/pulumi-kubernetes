// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";

// Create a new provider
const myk8s = new k8s.Provider("myk8s", {});

// Create a new provider with dry run enabled.
const myk8s2 = new k8s.Provider("myk8s2", {
    enableDryRun: true,
});

// Create a Pod using the custom provider
const nginxcontainer = new k8s.core.v1.Pod("nginx", {
    spec: {
        containers: [{
            image: "nginx:1.7.9",
            name: "nginx",
            ports: [{ containerPort: 80 }],
        }],
    },
}, { provider: myk8s });

// TODO(levi): Uncomment once https://github.com/pulumi/pulumi-kubernetes/issues/792 is fixed.
// // Create a Pod using the custom provider
// const nginxcontainer2 = new k8s.core.v1.Pod("nginx2", {
//     spec: {
//         containers: [{
//             image: "nginx:1.7.9",
//             name: "nginx",
//             ports: [{ containerPort: 80 }],
//         }],
//     },
// }, { provider: myk8s2 });
