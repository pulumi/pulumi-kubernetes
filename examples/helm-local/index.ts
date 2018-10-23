// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";

const nginx = new k8s.helm.v2.Chart("simple-nginx-local", {
    // Represents chart `stable/nginx-lego@v0.3.1`.
    path: "nginx-lego",
    values: {
        // Override for the Chart's `values.yml` file. Use `null` to zero out resource requests so it
        // can be scheduled on our (wimpy) CI cluster. (Setting these values to `null` is the "normal"
        // way to delete values.)
        nginx: { resources: null },
        default: { resources: null },
        lego: { resources: null }
    },
    transformations: [
        // Make every service private to the cluster, i.e., turn all services into ClusterIP instead of
        // LoadBalancer.
        (obj: any) => {
            if (obj.kind == "Service" && obj.apiVersion == "v1") {
                if (obj.spec && obj.spec.type && obj.spec.type == "LoadBalancer") {
                    obj.spec.type = "ClusterIP";
                }
            }
        }
    ]
});

// Export the (cluster-private) IP address of the Guestbook frontend.
const frontend = nginx.getResource("v1/Service", "simple-nginx-local-nginx");
export const frontendClusterIp = frontend.spec.apply(spec => spec.clusterIP);
