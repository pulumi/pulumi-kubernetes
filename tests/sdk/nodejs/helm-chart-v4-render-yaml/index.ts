import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

let config = new pulumi.Config();
let renderDir: string = config.require("renderDir");

const manifestProvider = new k8s.Provider("k8s_manifest_renderer", {
    renderYamlToDirectory: renderDir,
});

const cnpgOp = new k8s.helm.v4.Chart("cnpg-operator", {
    chart: "cloudnative-pg",
    repositoryOpts: {
        repo: "https://cloudnative-pg.io/charts/",
    },
}, {
    provider: manifestProvider,
});

