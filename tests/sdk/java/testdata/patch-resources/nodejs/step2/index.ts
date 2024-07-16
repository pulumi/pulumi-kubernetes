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

const patchRscNamespacePatching = new kubernetes.core.v1.NamespacePatch(
  "patch-rsc-namespace-patching",
  {
    metadata: {
      name: patchRscNamespace.metadata.name,
      annotations: {
        "pulumi.com/testPatchAnnotation": "patched",
      },
    },
  }
);

const deploymentPatching = new kubernetes.apps.v1.DeploymentPatch(
  "deployment-patching",
  {
    metadata: {
      name: deployment.metadata.name,
      namespace: patchRscNamespace.metadata.name,
      annotations: {
        "pulumi.com/testPatchAnnotation": "patched",
      },
    },
  }
);

const plainCRPatching = new kubernetes.apiextensions.CustomResourcePatch(
  "plain-cr-patching",
  {
    apiVersion: "patchtest.pulumi.com/v1",
    kind: "TestPatchResource",
    metadata: {
      name: plainCR.metadata.name,
      namespace: patchRscNamespace.metadata.name,
      annotations: {
        "pulumi.com/testPatchAnnotation": "patched",
      },
    },
  }
);

const patchCRPatching = new kubernetes.apiextensions.CustomResourcePatch(
  "patch-cr-patching",
  {
    apiVersion: "patchtest.pulumi.com/v1",
    kind: "TestPatchResourcePatch",
    metadata: {
      name: patchCR.metadata.name,
      namespace: patchRscNamespace.metadata.name,
      annotations: {
        "pulumi.com/testPatchAnnotation": "patched",
      },
    },
  }
);
