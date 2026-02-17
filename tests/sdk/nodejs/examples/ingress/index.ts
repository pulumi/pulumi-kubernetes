// Copyright 2016-2021, Pulumi Corporation.
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

import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import { output as outputs } from "@pulumi/kubernetes/types";

const ns = new k8s.core.v1.Namespace("test");
const namespace = ns.metadata.name;

const config = new pulumi.Config();
const useV1Beta1Ingress = config.get("use-v1beta1-ingress") != undefined;

// Install nginx ingress controller first
const ingressNs = new k8s.core.v1.Namespace("ingress-nginx");
const ingressController = new k8s.helm.v3.Release(
    "nginx-ingress", {
    name: "nginx-ingress",
    namespace: ingressNs.metadata.name,
    chart: "ingress-nginx",
    version: "4.1.4",
    repositoryOpts: {
        repo: "https://kubernetes.github.io/ingress-nginx",
    },
});

// nginx hello app
let helloLabels = { app: "hello", tier: "frontend" };
const helloService = new k8s.core.v1.Service("hello-svc", {
    metadata: {
        namespace: namespace,
        name: "hello",
        labels: helloLabels,
    },
    spec: {
        // GKE ingress will only work with service type of NodePort or LoadBalancer.
        // We only want one endpoint through the ingress so choosing NodePort.
        ports: [{ port: 8080, targetPort: 8080 }],
        selector: helloLabels,
    },
});

let helloDeployment = new k8s.apps.v1.Deployment("hello-app", {
    metadata: {
        namespace: namespace,
        name: "hello",
    },
    spec: {
        selector: {
            matchLabels: helloLabels,
        },
        replicas: 3,
        template: {
            metadata: {
                labels: helloLabels,
            },
            spec: {
                containers: [{
                    name: "hello",
                    image: "gcr.io/google-samples/hello-app:1.0",
                    resources: {
                        requests: {
                            cpu: "100m",
                            memory: "100Mi",
                        },
                    },
                    ports: [{
                        containerPort: 8080,
                    }],
                }],
            },
        },
    },
});

let feIngressNginx = undefined;

if (!useV1Beta1Ingress) {
    // Note - this uses the nginx ingress controller which should work across k8s providers.
    feIngressNginx = new k8s.networking.v1.Ingress("feIngressNginx", {
        metadata: {
            namespace: namespace,
            name: "feingress-nginx",
            annotations: {
                "kubernetes.io/ingress.class": "nginx",
                "nginx.ingress.kubernetes.io/ssl-redirect": "false"
            },
        },
        spec: {
            rules: [{
                host: "ingresshello.io",
                http: {
                    paths: [{
                        path: "/hello",
                        pathType: "Prefix",
                        backend: { service: { name: helloService.metadata.name, port: { number: 8080 } } }
                    }]
                }
            }],
        }
    }, { dependsOn: [ingressController] });

} else {
    // Note - this uses the nginx ingress controller which should work across k8s providers.
    feIngressNginx = new k8s.networking.v1beta1.Ingress("feIngressNginx", {
        metadata: {
            namespace: namespace,
            name: "feingress-nginx",
            annotations: {
                "kubernetes.io/ingress.class": "nginx",
                "nginx.ingress.kubernetes.io/ssl-redirect": "false"
            },
        },
        spec: {
            rules: [{
                host: "ingresshello.io",
                http: {
                    paths: [{
                        path: "/hello",
                        pathType: "Prefix",
                        backend: { serviceName: helloService.metadata.name, servicePort: 8080 }
                    }]
                }
            }],
        }
    }, { dependsOn: [ingressController] });
}
export const ingressNginxIp = feIngressNginx.status.loadBalancer.ingress[0].ip;

// REDIS MASTER

let redisMasterLabels = { app: "redis", tier: "backend", role: "master" };
let redisMasterService = new k8s.core.v1.Service("redis-master", {
    metadata: {
        namespace: namespace,
        name: "redis-master",
        labels: redisMasterLabels,
    },
    spec: {
        ports: [{ port: 6379, targetPort: 6379 }],
        selector: redisMasterLabels,
    },
});

let redisMasterDeployment = new k8s.apps.v1.Deployment("redis-master", {
    metadata: {
        namespace: namespace,
        name: "redis-master",
    },
    spec: {
        selector: {
            matchLabels: redisMasterLabels,
        },
        replicas: 1,
        template: {
            metadata: {
                labels: redisMasterLabels,
            },
            spec: {
                containers: [{
                    name: "master",
                    image: "docker.io/redis:6.0.5",
                    resources: {
                        requests: {
                            cpu: "100m",
                            memory: "100Mi",
                        },
                    },
                    ports: [{
                        containerPort: 6379,
                    }],
                }],
            },
        },
    },
});

// REDIS SLAVE
let redisSlaveLabels = { app: "redis", tier: "backend", role: "slave" };
let redisSlaveService = new k8s.core.v1.Service("redis-slave", {
    metadata: {
        namespace: namespace,
        name: "redis-slave",
        labels: redisSlaveLabels,
    },
    spec: {
        ports: [{ port: 6379, targetPort: 6379 }],
        selector: redisSlaveLabels,
    },
});

let redisSlaveDeployment = new k8s.apps.v1.Deployment("redis-slave", {
    metadata: {
        namespace: namespace,
        name: "redis-slave",
    },
    spec: {
        selector: {
            matchLabels: redisSlaveLabels
        },
        replicas: 1,
        template: {
            metadata: {
                labels: redisSlaveLabels,
            },
            spec: {
                containers: [{
                    name: "slave",
                    image: "us-docker.pkg.dev/google-samples/containers/gke/gb-redis-follower:v2",
                    resources: {
                        requests: {
                            cpu: "100m",
                            memory: "100Mi",
                        },
                    },
                    env: [{
                        name: "GET_HOSTS_FROM",
                        value: "dns",
                        // If your cluster config does not include a dns service, then to instead access an environment
                        // variable to find the master service's host, comment out the 'value: dns' line above, and
                        // uncomment the line below:
                        // value: "env"
                    }],
                    ports: [{
                        containerPort: 6379,
                    }],
                }],
            },
        },
    },
});

// FRONTEND
let frontendLabels = { app: "guestbook", tier: "frontend" };
const frontendService = new k8s.core.v1.Service("frontend", {
    metadata: {
        namespace: namespace,
        name: "frontend",
        labels: frontendLabels,
    },
    spec: {
        // GKE ingress will only work with service type of NodePort or LoadBalancer.
        // We only want one endpoint through the ingress so choosing NodePort.
        type: "NodePort",
        ports: [{ port: 80 }],
        selector: frontendLabels,
    },
});



let frontendDeployment = new k8s.apps.v1.Deployment("frontend", {
    metadata: {
        namespace: namespace,
        name: "frontend",
    },
    spec: {
        selector: {
            matchLabels: frontendLabels,
        },
        replicas: 3,
        template: {
            metadata: {
                labels: frontendLabels,
            },
            spec: {
                containers: [{
                    name: "php-redis",
                    image: "us-docker.pkg.dev/google-samples/containers/gke/gb-frontend:v5",
                    resources: {
                        requests: {
                            cpu: "100m",
                            memory: "100Mi",
                        },
                    },
                    env: [{
                        name: "GET_HOSTS_FROM",
                        value: "dns",
                        // If your cluster config does not include a dns service, then to instead access an environment
                        // variable to find the master service's host, comment out the 'value: dns' line above, and
                        // uncomment the line below:
                        // value: "env"
                    }],
                    ports: [{
                        containerPort: 80,
                    }],
                }],
            },
        },
    },
});

let feIngress = undefined;

if (!useV1Beta1Ingress) {
    // Note - this uses the built-in GKE ingress controller. This will not work across other Kubernetes providers.
    feIngress = new k8s.networking.v1.Ingress("feIngress", {
        metadata: {
            namespace: namespace,
            name: "feingress",
            annotations: {
                "kubernetes.io/ingress.class": "gce"
            },
        },
        spec: {
            rules: [{
                http: {
                    paths: [{
                        path: "/*",
                        pathType: "ImplementationSpecific",
                        backend: { service: { name: frontendService.metadata.name, port: { number: 80 } } }
                    }]
                }
            }],
        }
        // Don't want to race with the ingress controller admission webhook being fully available.
    }, { dependsOn: [ingressController] });
} else {
    // Note - this uses the built-in GKE ingress controller. This will not work across other Kubernetes providers.
    feIngress = new k8s.networking.v1beta1.Ingress("feIngress", {
        metadata: {
            namespace: namespace,
            name: "feingress",
            annotations: {
                "kubernetes.io/ingress.class": "gce"
            },
        },
        spec: {
            rules: [{
                http: {
                    paths: [{
                        path: "/*",
                        pathType: "ImplementationSpecific",
                        backend: { serviceName: frontendService.metadata.name, servicePort: 80 }
                    }]
                }
            }],
        }
        // Don't want to race with the ingress controller admission controller being fully available.
    }, { dependsOn: [ingressController] });
}
export const ingressIp = feIngress.status.loadBalancer.ingress[0].ip;
