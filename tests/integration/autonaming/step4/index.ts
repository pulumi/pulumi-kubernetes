// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";

//
// User has now specified `.metadata.name`, so Pulumi should replace the resource, and NOT allocate
// a name to it.
//

const pod = new k8s.core.v1.Pod("autonaming-test", {
  metadata: {
    name: "autonaming-test",
    labels: {app: "autonaming-test"},
  },
  spec: {
    containers: [
      {name: "nginx", image: "nginx:1.15-alpine"},
    ],
  },
});
