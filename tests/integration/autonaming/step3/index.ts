// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";

//
// Only the labels have changed, so no replace is triggered. Pulumi should update the object
// in-place, and the name should not be changed.
//

const pod = new k8s.core.v1.Pod("autonaming-test", {
  metadata: {
    labels: {app: "autonaming-test"},
  },
  spec: {
    containers: [
      {name: "nginx", image: "nginx:1.15-alpine"},
    ],
  },
});
