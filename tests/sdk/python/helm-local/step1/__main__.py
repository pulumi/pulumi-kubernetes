# Copyright 2016-2020, Pulumi Corporation.
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

from pulumi import ResourceOptions
from pulumi_kubernetes.core.v1 import ConfigMap
from pulumi_kubernetes.helm.v2 import Chart, LocalChartOpts

values = {"service": {"type": "ClusterIP"}}

chart = Chart("nginx", LocalChartOpts(path="nginx", values=values))

ConfigMap("foo", data={"foo": "bar"}, opts=ResourceOptions(depends_on=[chart.ready]))

# Deploy a duplicate chart with a different resource prefix to verify that multiple instances of the Chart
# can be managed in the same stack.
Chart("nginx", LocalChartOpts(path="nginx", resource_prefix="dup", values=values))
