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

import { Config } from "@pulumi/pulumi";
import * as random from "@pulumi/random";

const config = new Config();

export const appName = "test-ci-istio";

export const gcpProject = "pulumi-development";

export const gcpRegion = "us-west1";
export const gcpZone = "a";

// nodeCount is the number of cluster nodes to provision. Defaults to 3 if unspecified.
export const nodeCount = config.getNumber("nodeCount") || 5;

// nodeMachineType is the machine type to use for cluster nodes. Defaults to n1-standard-1 if unspecified.
// See https://cloud.google.com/compute/docs/machine-types for more details on available machine types.
export const nodeMachineType = config.get("nodeMachineType") || "n1-standard-2";

// masterUsername is the admin username for the cluster.
export const masterUsername = config.get("masterUsername") || "admin";

// masterPassword is the password for the admin user in the cluster.
export const masterPassword =
    config.get("password") ||
    new random.RandomString("password", {
        length: 16,
        special: true
    }).result;

// username is the admin username for the cluster.
export const gcpUsername =
    config.get("gcpUsername") || "test-ci@pulumi-development.iam.gserviceaccount.com";
