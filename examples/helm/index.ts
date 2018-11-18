// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";

const nginx = new k8s.helm.v2.Chart("simple-nginx", {
    // Represents chart `stable/nginx-lego@v0.3.1`.
    repo: "stable",
    chart: "nginx-lego",
    version: "0.3.1",
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
export const frontendClusterIp = nginx.getResourceProperty("v1/Service", "simple-nginx-nginx-lego", "spec")
    .apply(spec => spec.clusterIP);

// Test a variety of other inputs on a chart that creates no resources.
const empty1 = new k8s.helm.v2.Chart("empty1", {
    chart: "https://kubernetes-charts-incubator.storage.googleapis.com/raw-0.1.0.tgz",
});

const empty2 = new k8s.helm.v2.Chart("empty2", {
    chart: "raw",
    version: "0.1.0",
    fetchOpts: {
        repo: "https://kubernetes-charts-incubator.storage.googleapis.com/",
    },
});

const empty3 = new k8s.helm.v2.Chart("empty3", {
    chart: "raw",
    fetchOpts: {
        repo: "https://kubernetes-charts-incubator.storage.googleapis.com/",
    },
});