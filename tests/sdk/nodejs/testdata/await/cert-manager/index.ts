import * as kubernetes from "@pulumi/kubernetes";

const ns = new kubernetes.core.v1.Namespace("ns", {});

const provider = new kubernetes.Provider("provider", {
  namespace: ns.metadata.name,
});

const install = new kubernetes.yaml.v2.ConfigFile(
  "install",
  {
    file: "cert-manager-crds.yaml",
  },
  {
    // Add a "waitFor" annotation which waits for the cainject pod to come
    // online and modify our webhooks with valid CA bundles. Webhooks
    // needed by our Certificate will fail until this happens.
    transforms: [
      (args) => {
        if (
          args.type ===
            "kubernetes:admissionregistration.k8s.io/v1:ValidatingWebhookConfiguration" ||
          args.type ===
            "kubernetes:admissionregistration.k8s.io/v1:MutatingWebhookConfiguration"
        ) {
          args.props.metadata.annotations["pulumi.com/waitFor"] =
            "jsonpath={.webhooks[].clientConfig.caBundle}";
          return {
            props: args.props,
            opts: args.opts,
          };
        }
        return undefined;
      },
    ],
    provider: provider,
  }
);

const issuer = new kubernetes.apiextensions.CustomResource(
  "issuer",
  {
    apiVersion: "cert-manager.io/v1",
    kind: "ClusterIssuer",
    metadata: {
      name: "selfsigned-issuer",
    },
    spec: {
      selfSigned: {},
    },
  },
  { provider: provider, dependsOn: [install] }
);

new kubernetes.apiextensions.CustomResource(
  "certificate",
  {
    apiVersion: "cert-manager.io/v1",
    kind: "Certificate",
    metadata: {
      name: "selfsigned-cert",
    },
    spec: {
      isCA: true,
      commonName: "my-ca",
      secretName: "root-secret",
      privateKey: {
        algorithm: "ECDSA",
        size: 256,
      },
      issuerRef: {
        name: issuer.metadata.name,
        kind: issuer.kind,
        group: "cert-manager.io",
      },
    },
  },
  { provider: provider }
);
