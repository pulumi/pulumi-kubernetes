"""A Kubernetes Python Pulumi program"""

from pulumi_kubernetes.core.v1 import Namespace
from pulumi_kubernetes.helm.v3 import Release, ReleaseArgs, RepositoryOptsArgs

app_labels = {"app": "nginx"}

namespace = Namespace("test")

Release("release1",
        args=ReleaseArgs(
                chart="nginx",
                repository_opts=RepositoryOptsArgs(
                    repo="https://raw.githubusercontent.com/bitnami/charts/eb5f9a9513d987b519f0ecd732e7031241c50328/bitnami"
                ),
                namespace=namespace.metadata["name"],
                values={
                    "service": {"type": "ClusterIP"},
                    "image": {
                        "repository": "bitnamisecure",
                        "tag": "latest",
                    },
                },
                version="6.0.4",
            ))

# Deploy a duplicate chart release to verify that multiple instances of the Chart
# can be managed in the same stack.
Release("release2",
        args=ReleaseArgs(
                chart="nginx",
                repository_opts=RepositoryOptsArgs(
                    repo="https://raw.githubusercontent.com/bitnami/charts/eb5f9a9513d987b519f0ecd732e7031241c50328/bitnami"
                ),
                namespace=namespace.metadata["name"],
                values={
                    "service": {"type": "ClusterIP"},
                    "image": {
                        "repository": "bitnamisecure",
                        "tag": "latest",
                    },
                },
                version="6.0.4",
            ))
