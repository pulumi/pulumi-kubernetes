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

import pulumi as p
import pulumi_kubernetes as k8s

k8s_provider = k8s.Provider("gke")
default_opts = p.ResourceOptions(
  providers=[
    k8s_provider,
  ],
)

# use cert-manager as an example
k8s.yaml.ConfigFile(
  "cert-manager",
  file="https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.yaml",
  opts=default_opts,
)