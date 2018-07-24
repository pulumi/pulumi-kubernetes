// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";

//
// Create a simple Pod.
//

const pod = new k8s.core.v1.Pod("pod-test", {
  metadata: {
    name: "pod-test",
  },
  spec: {
    containers: [
      {name: "nginx", image: "nginx:1.13-alpine"},
    ],
  },
})
