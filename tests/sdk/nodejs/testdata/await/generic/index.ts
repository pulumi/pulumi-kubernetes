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

export const wantsCondition = new kubernetes.apiextensions.CustomResource(
  "wants-foo-condition",
  {
    apiVersion: "test.pulumi.com/v1",
    kind: "GenericAwaiter",
    metadata: {
      name: "wants-foo-condition",
      annotations: {
        "pulumi.com/waitFor": "condition=Foo",
        "pulumi.com/patchForce": "true", // Don't conflict with kubectl.
        "pulumi.com/timeoutSeconds": "60",
      },
    },
    spec: {
      someField: "",
    },
    status: {
      conditions: [
        {
          type: "Foo",
          status: "False",
        },
      ],
    },
  },
  { provider: provider, dependsOn: [crds] }
);

export const wantsField = new kubernetes.apiextensions.CustomResource(
  "wants-field",
  {
    apiVersion: "test.pulumi.com/v1",
    kind: "GenericAwaiter",
    metadata: {
      name: "wants-field",
      annotations: {
        "pulumi.com/waitFor": "jsonpath={.spec.someField}=foo",
        "pulumi.com/patchForce": "true", // Don't conflict with kubectl.
        "pulumi.com/timeoutSeconds": "60",
      },
    },
    spec: {
      someField: "",
    },
    status: {
      conditions: [],
    },
  },
  { provider: provider, dependsOn: [crds] }
);

export const wantsFieldAndCondition =
  new kubernetes.apiextensions.CustomResource(
    "wants-field-and-foo-condition",
    {
      apiVersion: "test.pulumi.com/v1",
      kind: "GenericAwaiter",
      metadata: {
        name: "wants-field-and-foo-condition",
        annotations: {
          "pulumi.com/waitFor":
            '["jsonpath={.spec.someField}=expected", "condition=Foo"]',
          "pulumi.com/patchForce": "true", // Don't conflict with kubectl.
          "pulumi.com/timeoutSeconds": "60",
        },
      },
      spec: {
        someField: "",
      },
      status: {
        conditions: [
          {
            type: "Foo",
            status: "False",
          },
        ],
      },
    },
    { provider: provider, dependsOn: [crds] }
  );
