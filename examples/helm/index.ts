// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as azure from "@pulumi/azure";
import * as helm from "@pulumi/kubernetes/helm";
import * as fs from "fs";
import * as azurecontainerservice from "./acs";
import * as k8s from "./kubernetes";

const skymanconfig = JSON.parse(fs.readFileSync("./skyman.json").toString());

const resourceGroup = new azure.core.ResourceGroup("acs", {
    location: "West US 2",
})

const acs = new azurecontainerservice.ACS("acs", {
    resourceGroupName: resourceGroup.name,
    location: "westus2",
    parameters: skymanconfig.acs,
});

let kubeProvider = new k8s.KubernetesProvider("acsKube", {
    kubeconfig: acs.kubeconfig.apply(JSON.stringify),
});

for (let chart of skymanconfig.charts) {
    new helm.v2.Chart(chart.name, chart.path, chart.args || {}, { provider: kubeProvider });
}
