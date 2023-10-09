import * as k8s from "@pulumi/kubernetes";

const namespace = new k8s.core.v1.Namespace("test");

new k8s.helm.v3.Chart("kube-version", {
  kubeVersion: "1.24.1",
  namespace: namespace.metadata.name,
  path: "helm-kube-version",
});