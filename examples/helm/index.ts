// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as helm from "@pulumi/kubernetes/helm";

// const mysql = new helm.v2.Chart("simple-mysql", "stable/mysql", {});
const nginx = new helm.v2.Chart("simple-nginx", "stable/nginx-lego", {});
