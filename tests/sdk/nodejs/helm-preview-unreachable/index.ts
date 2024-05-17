import * as k8s from "@pulumi/kubernetes";

const kubeconfig = `
apiVersion: v1
clusters:
  - cluster:
      server: 127.0.0.1:1
    name: helm-preview-unreachable
contexts:
  - context:
      cluster: helm-preview-unreachable
    name: default
current-context: default
kind: Config
`;

const namespace = new k8s.core.v1.Namespace("ns", {});

const provider = new k8s.Provider("k8s", { kubeconfig: kubeconfig });

new k8s.helm.v3.Chart(
  "template",
  {
    chart: "redis",
    fetchOpts: {
      repo: "https://raw.githubusercontent.com/bitnami/charts/eb5f9a9513d987b519f0ecd732e7031241c50328/bitnami",
    },
    version: "13.0.0",
    namespace: namespace.metadata.name,
  },
  { provider }
);
