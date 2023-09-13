// Copyright 2016-2023, Pulumi Corporation.
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

// Read the cluster kubeconfig.

// Create 2 Kubernetes providers with to simulate 2 clusters. The KUBECONFIG is obtained from the environment
// variables `KUBECONFIG_CLUSTER_0` and `KUBECONFIG_CLUSTER_1` respectively.
const provider0 = new k8s.Provider("k8s", {
  kubeconfig: process.env.KUBECONFIG_CLUSTER_0,
  enableConfigMapMutable: true,
});
const provider1 = new k8s.Provider("k8s2", {
  kubeconfig: process.env.KUBECONFIG_CLUSTER_1,
  enableConfigMapMutable: true,
});

// Create a namespace and a configmap in the provided cluster.
function createCM(clusterNumber: number, provider: k8s.Provider): [k8s.core.v1.Namespace, k8s.core.v1.ConfigMap] {
  const ns = new k8s.core.v1.Namespace(
    `test-cluster-${clusterNumber}`,
    undefined,
    { provider }
  );

  const cm = new k8s.core.v1.ConfigMap(
    `test-cluster-${clusterNumber}`,
    {
      metadata: {
        namespace: ns.metadata.name,
      },
      data: { foo: "step1" },
    },
    { provider }
  );

    return [ns, cm];  
}

const [ns0, cm0] = createCM(0, provider0);
const [ns1, cm1] = createCM(1, provider1);

// Export the resource names.
export const nsName0 = ns0.metadata.name;
export const cmName0 = cm0.metadata.name;
export const nsName1 = ns1.metadata.name;
export const cmName1 = cm1.metadata.name;