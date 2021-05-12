import * as k8s from "@pulumi/kubernetes";

const namespace = new k8s.core.v1.Namespace("test");

new k8s.helm.v3.Chart("skip-crd-rendering", {
  skipCRDRendering: true,
  namespace: namespace.metadata.name,
  path: "helm-skip-crd-rendering",
});

new k8s.helm.v3.Chart("allow-crd-rendering", {
  skipCRDRendering: false,
  namespace: namespace.metadata.name,
  path: "helm-allow-crd-rendering",
});
