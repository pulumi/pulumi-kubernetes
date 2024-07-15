import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const provider = new kubernetes.Provider("provider", {});
const patchRscNamespace = new kubernetes.core.v1.Namespace(
  "patch-rsc-namespace",
  {},
  {
    provider: provider,
  }
);

const deployment = new kubernetes.apps.v1.Deployment(
  "deployment",
  {
    metadata: {
      namespace: patchRscNamespace.metadata.name,
      labels: {
        app: "nginx",
      },
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: {
          app: "nginx",
        },
      },
      template: {
        metadata: {
          labels: {
            app: "nginx",
          },
        },
        spec: {
          containers: [
            {
              image: "nginx:1.14.2",
              name: "nginx",
              ports: [
                {
                  containerPort: 80,
                },
              ],
            },
          ],
        },
      },
    },
  },
  {
    provider: provider,
  }
);

const plainCR = new kubernetes.apiextensions.CustomResource(
  "plain-cr",
  {
    apiVersion: "patchtest.pulumi.com/v1",
    kind: "TestPatchResource",
    metadata: {
      namespace: patchRscNamespace.metadata.name,
    },
    spec: {
      foo: "bar",
    },
  },
  {
    provider: provider,
  }
);

const patchCR = new kubernetes.apiextensions.CustomResource(
  "patch-cr",
  {
    apiVersion: "patchtest.pulumi.com/v1",
    kind: "TestPatchResourcePatch",
    metadata: {
      namespace: patchRscNamespace.metadata.name,
    },
    spec: {
      foo: "bar",
    },
  },
  {
    provider: provider,
  }
);

export const nsName = patchRscNamespace.metadata.name;
export const depName = deployment.metadata.name;
export const plainCRName = plainCR.metadata.name;
export const patchCRName = patchCR.metadata.name;
