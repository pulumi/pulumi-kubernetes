// Copyright 2016-2021, Pulumi Corporation.
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

const provider = new k8s.Provider("k8s", {
    enableReplaceCRD: true,
});

//
// Create a CustomResourceDefinition.
//

new k8s.yaml.ConfigFile(
    `operator-manifest`,
    {
        file: 'https://download.elastic.co/downloads/eck/1.8.0/crds.yaml', // New version of the CRD
    },
    {provider},
)
