// Copyright 2016-2019, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import * as k8s from "@pulumi/kubernetes";

const namespace = new k8s.core.v1.Namespace("test");

const appLabels = {app: "apache"};
const defaults = {
  name: "apache",
  namespace: namespace.metadata.apply(metadata => metadata.namespace),
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
