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
from pulumi_kubernetes.yaml import ConfigFile, ConfigGroup

ns = Namespace("ns")
ns2 = Namespace("ns2")


def set_namespace(namespace):
    def f(obj):
        if "metadata" in obj:
            obj["metadata"]["namespace"] = namespace.metadata["name"]
        else:
            obj["metadata"] = {"namespace": namespace.metadata["name"]}

    return f


def secret_status(obj, opts):
    if obj["kind"] == "Pod" and obj["apiVersion"] == "v1":
        opts.additional_secret_outputs = ["apiVersion"]


cf_local = ConfigFile(
    "yaml-test",
    "manifest.yaml",
    transformations=[
        set_namespace(ns),
        # TODO[pulumi/pulumi#2782] Testing of secrets blocked on a bug in Python support for secrets.
        # secret_status,
    ],
)

# Create resources from a simple YAML manifest in a test namespace.
cf_url = ConfigFile(
    "deployment",
    file_id="https://raw.githubusercontent.com/kubernetes/website/master/content/en/examples/controllers/"
            "nginx-deployment.yaml",
    transformations=[set_namespace(ns)],
)
cf_url2 = ConfigFile(
    "deployment",
    file_id="https://raw.githubusercontent.com/kubernetes/website/master/content/en/examples/controllers/"
            "nginx-deployment.yaml",
    transformations=[set_namespace(ns2)],
    resource_prefix="dup"
)
cg = ConfigGroup(
    "deployment",
    files=["ns*.yaml"],
    yaml=["""
apiVersion: v1
kind: Namespace
metadata:
  name: cg3
    """]
)
