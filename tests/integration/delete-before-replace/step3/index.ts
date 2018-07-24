// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";

//
// Cause hard delete-before-replace, when specifying a new namespace ("delete") that is equivalent
// to the old namespace ("", the blanks string). We trigger the hard-replace by again changing the
// Pod's container image tag. If the new pod has succeeded in being created, and the name is
// `pod-test`, then we know that it was deleted before being replaced, because Kubernetes would have
// otherwise complained you can't add two pods with that name.
//

const pod = new k8s.core.v1.Pod("pod-test", {
  metadata: {
    name: "pod-test",
    namespace: "default",
  },
  spec: {
    containers: [
      {name: "nginx", image: "nginx:1.13-alpine"},
    ],
  },
})
