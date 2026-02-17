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

const ns = new k8s.core.v1.Namespace("test");
const namespace = ns.metadata.name;

// REDIS MASTER

let redisMasterLabels = { app: "redis", tier: "backend", role: "master"};
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
let frontendService = new k8s.core.v1.Service("frontend", {
    metadata: {
        namespace: namespace,
        name: "frontend",
        labels: frontendLabels,
    },
    spec: {
        // If your cluster supports it, uncomment the following to automatically create
        // an external load-balanced IP for the frontend service.
        // type: "LoadBalancer",
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
