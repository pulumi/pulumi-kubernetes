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
from pulumi_kubernetes.core.v1 import Namespace
from pulumi_kubernetes.yaml import ConfigFile

ns = Namespace("ns")


def add_namespace(obj):
    if "metadata" in obj:
        obj["metadata"]["namespace"] = ns.metadata["name"]
    else:
        obj["metadata"] = {"namespace": ns.metadata["name"]}

    return obj


cf_local = ConfigFile(
    "yaml-test",
    "manifest.yaml",
    transformations=[add_namespace],
)

# Create resources from standard Kubernetes guestbook YAML example in the test namespace.
cf_url = ConfigFile(
    name="guestbook",
    file_id="https://raw.githubusercontent.com/pulumi/pulumi-kubernetes/master/tests/examples/yaml-guestbook/yaml"
            "/guestbook.yaml",
    transformations=[add_namespace],
)
