// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";

//
// Cause hard delete-before-replace. Changing the Pod's container image tag, causes it to be
// replaced, and since we've manually specified a name, the engine has no choice but to delete it
// first, since names are unique, and it can't generate the name by itself.
//

const pod = new k8s.core.v1.Pod("pod-test", {
  metadata: {
    name: "pod-test",
  },
  spec: {
    containers: [
      {name: "nginx", image: "nginx:1.15-alpine"},
    ],
  },
})
