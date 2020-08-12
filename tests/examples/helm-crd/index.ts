// Copyright 2016-2020, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import * as k8s from "@pulumi/kubernetes";
import * as os from "os";
import * as pulumi from "@pulumi/pulumi";

const namespace = new k8s.core.v1.Namespace("argocd");
const namespaceName = namespace.metadata.name;

const argoCrd = new k8s.helm.v3.Chart("argocd-crd", {
    chart: "argo-cd",
    namespace: namespaceName,
    installCRDs: true,
    fetchOpts: {
        home: os.homedir(),
        repo: "https://argoproj.github.io/argo-helm"
    },
});

const argoNoCrd = new k8s.helm.v3.Chart("argocd-nocrd", {
    chart: "argo-cd",
    namespace: namespaceName,
    installCRDs: false,
    fetchOpts: {
        home: os.homedir(),
        repo: "https://argoproj.github.io/argo-helm"
    },
});
