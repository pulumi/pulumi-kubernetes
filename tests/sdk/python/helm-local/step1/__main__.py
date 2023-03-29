#  Copyright 2016-2021, Pulumi Corporation.
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

import pulumi
from pulumi_kubernetes.core.v1 import ConfigMap, ConfigMapInitArgs
from pulumi_kubernetes.helm.v2 import Chart, LocalChartOpts
from pulumi_kubernetes.helm.v3 import Release, ReleaseArgs

values = {"service": {"type": "ClusterIP"}}

chart = Chart("nginx", LocalChartOpts(path="nginx-chart", values=values))
foo = ConfigMap("foo", ConfigMapInitArgs(data={"foo": "bar"}), opts=pulumi.ResourceOptions(depends_on=chart.ready))

# Deploy a duplicate chart with a different resource prefix to verify that multiple instances of the Chart
# can be managed in the same stack.
Chart("nginx", LocalChartOpts(path="nginx-chart", resource_prefix="dup", values=values))

# Deploy the same chart but using v3.Release with local path argument
Release('nginx-release', ReleaseArgs(chart="nginx", path="nginx-chart", values=values))
