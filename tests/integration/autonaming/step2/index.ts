// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";

//
// The image in the Pod's container has changed, triggering a replace. Because `.metadata.name` is
// not specified, Pulumi again will provide a name upon creation of the new Pod resource.
//

const pod = new k8s.core.v1.Pod("autonaming-test", {
  spec: {
    containers: [
      {name: "nginx", image: "nginx:1.15-alpine"},
    ],
  },
});
