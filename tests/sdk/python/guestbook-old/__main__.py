# Copyright 2016-2019, Pulumi Corporation.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
from pulumi_kubernetes.apps.v1 import Deployment
from pulumi_kubernetes.core.v1 import Service, Namespace

namespace = Namespace("test")

redis_leader_labels = {
    "app": "redis",
    "tier": "backend",
    "role": "master"
}

redis_leader_deployment = Deployment(
    "redis-leader",
    metadata={
        "namespace": namespace.metadata.apply(lambda x: x.name),
    },
    spec={
        "selector": {
            "match_labels": redis_leader_labels,
        },
        "replicas": 1,
        "template": {
            "metadata": {
                "labels": redis_leader_labels,
            },
            "spec": {
                "containers": [{
                    "name": "master",
                    "image": "docker.io/redis:6.0.5",
                    "resources": {
                        "requests": {
                            "cpu": "100m",
                            "memory": "100Mi",
                        },
                    },
                    "ports": [{
                        "container_port": 6379,
                    }],
                }],
            },
        },
    })

redis_leader_service = Service(
    "redis-leader",
    metadata={
        "namespace": namespace.metadata.apply(lambda x: x.name),
        "labels": redis_leader_labels
    },
    spec={
        "ports": [{
            "port": 6379,
            "target_port": 6379,
        }],
        "selector": redis_leader_labels
    })

redis_follower_labels = {
    "app": "redis",
    "tier": "backend",
    "role": "slave"
}

redis_follower_deployment = Deployment(
    "redis-follower",
    metadata={
        "namespace": namespace.metadata.apply(lambda x: x.name),
    },
    spec={
        "selector": {
            "match_labels": redis_follower_labels
        },
        "replicas": 1,
        "template": {
            "metadata": {
                "labels": redis_follower_labels,
            },
            "spec": {
                "containers": [{
                    "name": "slave",
                    "image": "gcr.io/google_samples/gb-redisslave:v1",
                    "resources": {
                        "requests": {
                            "cpu": "100m",
                            "memory": "100Mi",
                        },
                    },
                    "env": [{
                        "name": "GET_HOSTS_FROM",
                        "value": "dns",
                        # If your cluster config does not include a dns service, then to instead access an environment
                        # variable to find the master service's host, comment out the 'value: dns' line above, and
                        # uncomment the line below:
                        # value: "env"
                    }],
                    "ports": [{
                        "container_port": 6379,
                    }],
                }],
            },
        },
    })

redis_follower_service = Service(
    "redis-follower",
    metadata={
        "namespace": namespace.metadata.apply(lambda x: x.name),
        "labels": redis_follower_labels
    },
    spec={
        "ports": [{
            "port": 6379,
            "target_port": 6379,
        }],
        "selector": redis_follower_labels
    })

# Frontend
frontend_labels = {
    "app": "guestbook",
    "tier": "frontend"
}
frontend_service = Service(
    "frontend",
    metadata={
        "namespace": namespace.metadata.apply(lambda x: x.name),
        "labels": frontend_labels
    },
    spec={
        # If your cluster supports it, uncomment the following to automatically create
        # an external load-balanced IP for the frontend service.
        # "type": "LoadBalancer",
        "ports": [{
            "port": 80
        }],
        "selector": frontend_labels,
    })

frontend_deployment = Deployment(
    "frontend",
    metadata={
        "namespace": namespace.metadata.apply(lambda x: x.name),
    },
    spec={
        "selector": {
            "match_labels": frontend_labels,
        },
        "replicas": 3,
        "template": {
            "metadata": {
                "labels": frontend_labels,
            },
            "spec": {
                "containers": [{
                    "name": "php-redis",
                    "image": "us-docker.pkg.dev/google-samples/containers/gke/gb-frontend:v5",
                    "resources": {
                        "requests": {
                            "cpu": "100m",
                            "memory": "100Mi",
                        },
                    },
                    "env": [{
                        "name": "GET_HOSTS_FROM",
                        "value": "dns",
                        # If your cluster config does not include a dns service, then to instead access an environment
                        # variable to find the master service's host, comment out the 'value: dns' line above, and
                        # uncomment the line below:
                        # "value": "env"
                    }],
                    "ports": [{
                        "container_port": 80,
                    }],
                }],
            },
        },
    })
