// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";

//
// A simple Pod definition. `.metadata.name` is not provided, so Pulumi will allocate a unique name
// to the resource upon creation.
//

const pod = new k8s.core.v1.Pod("autonaming-test", {
  spec: {
    containers: [
      {name: "nginx", image: "nginx"},
    ],
  },
});
