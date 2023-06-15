#  Copyright 2016-2022, Pulumi Corporation.
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

from pulumi import ResourceOptions
from pulumi_kubernetes import Provider
from pulumi_kubernetes.core.v1 import Pod, Namespace, PodSpecArgs, ContainerArgs, ContainerPortArgs
from pulumi_kubernetes.meta.v1 import ObjectMetaArgs

my_k8s = Provider("myk8s")

namespace = Namespace("test")

nginx = Pod(
    "nginx",
    metadata=ObjectMetaArgs(
        namespace=namespace.metadata.apply(lambda x: x.name),
    ),
    spec=PodSpecArgs(
        containers=[ContainerArgs(
            image="nginx:1.7.9",
            name="nginx",
            ports=[ContainerPortArgs(
                container_port=80,
            )],
        )],
    ), opts=ResourceOptions(provider=my_k8s))
