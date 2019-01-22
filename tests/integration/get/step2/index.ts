// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";


k8s.apiextensions.CustomResource.get("my-new-cron-object-get", {
    apiVersion: "stable.example.com/v1",
    kind: "CronTab",
    id: "my-new-cron-object",
});
