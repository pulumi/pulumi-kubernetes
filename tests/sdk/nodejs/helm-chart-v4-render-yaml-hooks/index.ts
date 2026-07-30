import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

let config = new pulumi.Config();
let renderDir: string = config.require("renderDir");

const manifestProvider = new k8s.Provider("k8s_manifest_renderer", {
    renderYamlToDirectory: renderDir,
});

// The local chart contains a normal ConfigMap, a pre-install hook ConfigMap
// (helm.sh/hook: pre-install) and a test hook Pod (helm.sh/hook: test). With
// includeHooks set, the pre-install hook should be rendered while the test hook
// is excluded.
const chart = new k8s.helm.v4.Chart("hooks", {
    chart: "./testchart",
    includeHooks: true,
}, {
    provider: manifestProvider,
});
