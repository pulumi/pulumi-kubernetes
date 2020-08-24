// Copyright 2016-2019, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import * as k8s from "@pulumi/kubernetes";

// This Chart was hanging and then eventually crashing with an out-of-memory error [1], and fixed in [2].
// Verify that this Chart installs successfully.
//
// [1] https://github.com/pulumi/pulumi-kubernetes/issues/963
// [2] https://github.com/pulumi/pulumi-kubernetes/pull/974

new k8s.yaml.ConfigFile("cert-manager", {file: "./cert-manager.yaml"});
