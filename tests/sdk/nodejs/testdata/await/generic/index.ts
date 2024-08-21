import * as kubernetes from "@pulumi/kubernetes";

const ns = new kubernetes.core.v1.Namespace("ns", {
  metadata: {
    name: "generic-await",
    annotations: {
      "pulumi.com/deletionPropagationPolicy": "background",
    },
  },
});

const provider = new kubernetes.Provider("provider", {
  namespace: ns.metadata.name,
});

const crds = new kubernetes.yaml.v2.ConfigFile(
  "crds",
  { file: "crds.yaml" },
  { provider: provider }
);

export const wantsReady = new kubernetes.apiextensions.CustomResource(
  "wants-ready-condition",
  {
    apiVersion: "test.pulumi.com/v1",
    kind: "GenericAwaiter",
    metadata: {
      name: "wants-ready-condition",
      annotations: {
        "pulumi.com/patchForce": "true", // Don't conflict with kubectl.
        "pulumi.com/timeoutSeconds": "60",
      },
    },
    spec: {
      someField: "untouched",
    },
    status: {
      conditions: [
        {
          type: "Ready",
          status: "False",
        },
      ],
    },
  },
  { provider: provider, dependsOn: [crds] }
);

export const wantsGenerationIncrement =
  new kubernetes.apiextensions.CustomResource(
    "wants-generation-increment",
    {
      apiVersion: "test.pulumi.com/v1",
      kind: "GenericAwaiter",
      metadata: {
        name: "wants-generation-increment",
        generation: 2,
        annotations: {
          "pulumi.com/patchForce": "true", // Don't conflict with kubectl.
          "pulumi.com/timeoutSeconds": "60",
        },
      },
      spec: {
        someField: "untouched",
      },
      status: {
        observedGeneration: 1,
        conditions: [
          {
            type: "Ready",
            status: "True",
          },
        ],
      },
    },

    { provider: provider, dependsOn: [crds] }
  );
