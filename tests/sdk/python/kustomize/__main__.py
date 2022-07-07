# Copyright 2016-2022, Pulumi Corporation.
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

import pulumi_kubernetes as k8s
from pulumi import ResourceOptions

from typing import Any


def set_namespace(namespace):
    def f(obj: Any):
        if "metadata" in obj:
            obj["metadata"]["namespace"] = namespace.metadata["name"]
        else:
            obj["metadata"] = {"namespace": namespace.metadata["name"]}

    return f


provider = k8s.Provider("k8s")

ns = k8s.core.v1.Namespace("ns", opts=ResourceOptions(provider=provider))

k8s.kustomize.Directory(
    "kustomize-local",
    "helloWorld",
    transformations=[set_namespace(ns)],
    opts=ResourceOptions(provider=provider),
)
