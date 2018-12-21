// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";

// Create resources from standard Kubernetes guestbook YAML example.
new k8s.yaml.ConfigFile("guestbook", {
    file:
        "https://raw.githubusercontent.com/pulumi/pulumi-kubernetes/master/examples/yaml-guestbook/yaml/guestbook.yaml"
});
