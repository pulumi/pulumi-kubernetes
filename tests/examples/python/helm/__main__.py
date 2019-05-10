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
from pulumi_kubernetes.helm.v2 import Chart, ChartOpts

ns = Namespace("ns")


def set_namespace(obj):
    if "metadata" in obj:
        obj["metadata"]["namespace"] = ns.metadata["name"]
    else:
        obj["metadata"] = {"namespace": ns.metadata["name"]}

    return obj


def cluster_ip(obj):
    if obj["kind"] == "Service" and obj["apiVersion"] == "v1":
        if "spec" in obj and "type" in obj["spec"] and obj["spec"]["type"] == "LoadBalancer":
            obj["spec"]["type"] = "ClusterIP"

    return obj


Chart("nginx-lego", ChartOpts(
    "stable/nginx-lego",
    version="0.3.1",
    values={"nginx": None, "default": None, "lego": None},
    transformations=[cluster_ip, set_namespace]))
