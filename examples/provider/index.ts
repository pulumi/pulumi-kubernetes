// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";

// Use the existing ~/.kube/config kubeconfig
const kubeconfig = fs.readFileSync(path.join(os.homedir(), ".kube", "config")).toString();

// Pick a context using config
const config = new pulumi.Config("provider");
const context = config.require("context");

// Create a new provider
const myk8s = new k8s.Provider("myk8s", {
    kubeconfig: kubeconfig,
    context: context,
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
