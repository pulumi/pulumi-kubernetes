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
from os.path import expanduser

from pulumi_kubernetes.core.v1 import Namespace
from pulumi_kubernetes.helm.v3 import Chart, ChartOpts, FetchOpts
from pulumi_random import RandomString

namespace = Namespace("test")

rs = RandomString("random-string", length=8).result

values = {"service": {"type": "ClusterIP"}, "random-string": rs}

Chart(
    "remote-chart",
    ChartOpts(
        chart="ingress-nginx",
        fetch_opts=FetchOpts(
            home=expanduser("~"),
            repo="https://kubernetes.github.io/ingress-nginx"
        ),
        namespace=namespace.metadata["name"],
        values={"service": {"type": "ClusterIP"}},
        version="4.13.2",
    ))

# Deploy a duplicate chart with a different resource prefix to verify that multiple instances of the Chart
# can be managed in the same stack.
Chart(
    "remote-chart",
    ChartOpts(
        chart="ingress-nginx",
        resource_prefix="dup",
        fetch_opts=FetchOpts(
            home=expanduser("~"),
            repo="https://kubernetes.github.io/ingress-nginx"
        ),
        namespace=namespace.metadata["name"],
        values={"service": {"type": "ClusterIP"}},
        version="4.13.2",
    ))
