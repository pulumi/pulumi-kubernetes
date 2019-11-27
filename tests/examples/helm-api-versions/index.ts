import * as k8s from "@pulumi/kubernetes";

const namespace = new k8s.core.v1.Namespace("test");
const namespaceName = namespace.metadata.name;

new k8s.helm.v2.Chart("api-versions", {
  apiVersions: ["foo", "bar"],
  path: "helm-api-versions",
  transformations: [
    // Put every resource in the created namespace.
    (obj: any) => {
      if (obj.metadata !== undefined) {
        obj.metadata.namespace = namespaceName;
      } else {
        obj.metadata = { namespace: namespaceName };
      }
    }
  ]
});
