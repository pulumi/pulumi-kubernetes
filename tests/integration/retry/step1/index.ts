// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

//
// Tests that if we force a `Pod` to be created before the `Namespace` it is supposed to exist in,
// it will retry until created.
//

new k8s.core.v1.Pod("nginx", {
    metadata: { name: "nginx", namespace: "test" },
    spec: {
        containers: [
            {
                name: "nginx",
                image: "nginx:1.7.9",
                ports: [{ containerPort: 80 }]
            }
        ]
    }
});

new k8s.core.v1.Namespace("test", {
    metadata: {
        // Wait 10 seconds before creating the namespace, to make sure we retry the Pod creation.
        annotations: {
            timeout: new Promise(resolve => {
                if (pulumi.runtime.isDryRun()) {
                    return resolve("<output>");
                }
                setTimeout(() => resolve("done"), 10000);
            })
        },
        name: "test"
    }
});
