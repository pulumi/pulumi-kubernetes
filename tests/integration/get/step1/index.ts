// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";

//
// `get`s the Kubernetes Dashboard, which is deployed by default in minikube.
//

const dashboard = k8s.core.v1.Service.get("kube-dashboard", "kubernetes");
