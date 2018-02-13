// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

import * as pulumi from "pulumi";
import * as kubernetes from "@pulumi/kubernetes";

// REDIS MASTER

let redisMasterLabels = { app: "redis", tier: "backend", role: "master"};
let redisMasterService = new kubernetes.Service("redis-master", {
    metadata: [{
        name: "redis-master",
        labels: [redisMasterLabels],
    }],
    spec: [{
        port: [{ port: 6379, targetPort: 6379 }],
        selector: [redisMasterLabels],
    }],
});
let redisMasterDeployment = new kubernetes.Deployment("redis-master", {
    metadata: [{
        name: "redis-master",
    }],
    spec: [{
        selector: [redisMasterLabels],
        replicas: 1,
        strategy: [],
        template: [{
            metadata: [{
                labels: [redisMasterLabels],
            }],
            spec: [{
                container: [{
                    name: "master",
                    image: "k8s.gcr.io/redis:e2e",
                    resources: [{
                        requests: [{
                            cpu: "100m",
                            memory: "100Mi",
                        }]
                    }],
                    port: [{
                        containerPort: 6379,
                    }],
                }],
            }],
        }],
    }],
});

// REDIS SLAVE
let redisSlaveLabels = { app: "redis", tier: "backend", role: "slave" };
let redisSlaveService = new kubernetes.Service("redis-slave", {
    metadata: [{
        name: "redis-slave",
        labels: [redisSlaveLabels],
    }],
    spec: [{
        port: [{ port: 6379, targetPort: 6379 }],
        selector: [redisSlaveLabels],
    }],
});
let redisSlaveDeployment = new kubernetes.Deployment("redis-slave", {
    metadata: [{
        name: "redis-slave",
    }],
    spec: [{
        selector: [redisSlaveLabels],
        replicas: 1,
        strategy: [],
        template: [{
            metadata: [{
                labels: [redisSlaveLabels],
            }],
            spec: [{
                container: [{
                    name: "slave",
                    image: "gcr.io/google_samples/gb-redisslave:v1",
                    resources: [{
                        requests: [{
                            cpu: "100m",
                            memory: "100Mi",
                        }]
                    }],
                    env: [{
                        name: "GET_HOSTS_FROM",
                        value: "dns",
                        // If your cluster config does not include a dns service, then to instead access an environment
                        // variable to find the master service's host, comment out the 'value: dns' line above, and
                        // uncomment the line below: 
                        // value: env
                    }],
                    port: [{
                        containerPort: 6379,
                    }],
                }],
            }],
        }],
    }],
});

// FRONTEND
let frontendLabels = { app: "guestbook", tier: "frontend" };
let frontendService = new kubernetes.Service("frontend", {
    metadata: [{
        name: "frontend",
        labels: [frontendLabels],
    }],
    spec: [{
        type: "LoadBalancer",
        port: [{ port: 80 }],
        selector: [frontendLabels],
    }],
});
let frontendDeployment = new kubernetes.Deployment("frontend", {
    metadata: [{
        name: "frontend",
    }],
    spec: [{
        selector: [frontendLabels],
        replicas: 3,
        strategy: [],
        template: [{
            metadata: [{
                labels: [frontendLabels],
            }],
            spec: [{
                container: [{
                    name: "php-redis",
                    image: "gcr.io/google-samples/gb-frontend:v4",
                    resources: [{
                        requests: [{
                            cpu: "100m",
                            memory: "100Mi",
                        }]
                    }],
                    env: [{
                        name: "GET_HOSTS_FROM",
                        value: "dns",
                        // If your cluster config does not include a dns service, then to instead access an environment
                        // variable to find the master service's host, comment out the 'value: dns' line above, and
                        // uncomment the line below: 
                        // value: env
                    }],
                    port: [{
                        containerPort: 80,
                    }],
                }],
            }],
        }],
    }],
});

export let frontendPort: pulumi.Output<number> = frontendService.spec.apply(spec => spec[0].port![0].nodePort);
