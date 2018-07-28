// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import * as helm from "@pulumi/kubernetes/helm";

const nginx = new helm.v2.Chart("simple-nginx", {
  // Represents chart `stable/nginx-lego@v0.3.1`.
  repo: "stable", chart: "nginx-lego", version: "0.3.1",
  values: {
    // Override for the Chart's `values.yml` file. Use `null` to zero out resource requests so it
    // can be scheduled on our (wimpy) CI cluster. (Setting these values to `null` is the "normal"
    // way to delete values.)
    nginx: {resources: null},
    default: {resources: null},
    lego: {resources: null},
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

// Export the ClusterIP of the frontend service.
const frontend = <k8s.core.v1.Service>nginx.resources["v1/Service::default/simple-nginx-nginx-lego"];
exports.frontendClusterIp = frontend.spec.apply(spec => spec.clusterIP);