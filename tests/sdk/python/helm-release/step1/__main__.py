"""A Kubernetes Python Pulumi program"""

from pulumi_kubernetes.core.v1 import Namespace
from pulumi_kubernetes.helm.v3 import Release, ReleaseArgs, RepositoryOptsArgs

app_labels = {"app": "nginx"}

namespace = Namespace("test")

Release("release1",
        args=ReleaseArgs(
                chart="nginx",
                repository_opts=RepositoryOptsArgs(
                    repo="https://charts.bitnami.com/bitnami"
                ),
                namespace=namespace.metadata["name"],
                values={"service": {"type": "ClusterIP"}},
                version="6.0.4",
            ))

# Deploy a duplicate chart release to verify that multiple instances of the Chart
# can be managed in the same stack.
Release("release2",
        args=ReleaseArgs(
                chart="nginx",
                repository_opts=RepositoryOptsArgs(
                    repo="https://charts.bitnami.com/bitnami"
                ),
                namespace=namespace.metadata["name"],
                values={"service": {"type": "ClusterIP"}},
                version="6.0.4",
            ))
