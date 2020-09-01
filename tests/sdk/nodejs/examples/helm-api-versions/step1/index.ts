import * as k8s from "@pulumi/kubernetes";

const namespace = new k8s.core.v1.Namespace("test");

new k8s.helm.v2.Chart("api-versions", {
  apiVersions: ["foo", "bar"],
  namespace: namespace.metadata.name,
  path: "helm-api-versions",
});

new k8s.helm.v2.Chart("single-api-version", {
  apiVersions: ["foo"],
  namespace: namespace.metadata.name,
  path: "helm-single-api-version",
});
