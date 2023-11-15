# Copyright 2016-2023, Pulumi Corporation.
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


# This will be unknown during the initial preview.
unknown = k8s.Provider("provider").id.apply(lambda _: True)

# This provider will be unconfigured when the passed-in configuration has an unknown value.
provider = k8s.Provider("k8s", suppress_deprecation_warnings=unknown)

ns = k8s.core.v1.Namespace("unconfiguredtest", opts=ResourceOptions(provider=provider))

# An error shouldn't be raised when called using the unconfigured provider.
k8s.kustomize.Directory(
    "kustomize-local",
    "helloWorld",
    transformations=[set_namespace(ns)],
    opts=ResourceOptions(provider=provider),
)
