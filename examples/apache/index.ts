// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as k8s from "@pulumi/kubernetes";

const appLabels = {app: "apache"};
const defaults = {
  name: "apache",
  namespace: "default",
  labels: appLabels,
  serviceSelector: appLabels,
  // ref: https://hub.docker.com/r/bitnami/apache/tags/
  imageTag: "2.4.23-r12",
  // ref: http://kubernetes.io/docs/user-guide/images/#pre-pulling-images
  imagePullPolicy: "IfNotPresent",
};

const apacheDeployment = new k8s.apps.v1.Deployment(
  defaults.name,
  {
    metadata: {
      namespace: defaults.namespace,
      name: defaults.name,
      labels: defaults.labels
    },
    spec: {
      replicas: 1,
      selector: {
        matchLabels: defaults.labels,
      },
      template: {
        metadata: {
          labels: defaults.labels
        },
        spec: {
          containers: [
            {
              name: defaults.name,
              image: `bitnami/apache:${defaults.imageTag}`,
              imagePullPolicy: defaults.imagePullPolicy,
              ports: [
                {
                  name: "http",
                  containerPort: 80,
                },
                {
                  name: "https",
                  containerPort: 443,
                }
              ],
              livenessProbe: {
                httpGet: {
                  path: "/",
                  port: "http",
                },
                initialDelaySeconds: 30,
                timeoutSeconds: 5,
                failureThreshold: 6,
              },
              readinessProbe: {
                httpGet: {
                  path: "/",
                  port: "http",
                },
                initialDelaySeconds: 5,
                timeoutSeconds: 3,
                periodSeconds: 5,
              },
              volumeMounts: [
                {
                  name: "apache-data",
                  mountPath: "/bitnami/apache",
                }
              ],
            }
          ],
          volumes: [
            {
              name: "apache-data",
              emptyDir: {},
            }
          ]
        },
      },
    }
  });

const apacheService = new k8s.core.v1.Service(
  defaults.name,
  {
    metadata: {
      name: defaults.name,
      namespace: defaults.namespace,
      labels: {
        app: defaults.name,
        run: "apache",
      },
    },
    spec: {
      // NOTE: Uncomment this if your cloud provider supports services of type LoadBalancer.
      // type: "LoadBalancer",
      ports: [
        {
          name: "http",
          port: 80,
          targetPort: "http",
        },
        {
          name: "https",
          port: 443,
          targetPort: "https",
        },
      ],
      selector: defaults.serviceSelector
    },
  }
);
