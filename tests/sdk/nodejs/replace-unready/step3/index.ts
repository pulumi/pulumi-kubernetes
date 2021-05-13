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

// Create a Job that will fail.
// The replaceUnready annotation is set to true, so the replace behavior is enabled.
new k8s.batch.v1.Job("test", {
  metadata: {
    annotations: {
      "pulumi.com/replaceUnready": "true"
    }
  },
  spec: {
    template: {
      spec: {
        containers: [
          {
            name: "boom",
            image: "alpine",
            command: ["fail"],
          }
        ],
        restartPolicy: "Never"
      }
    },
    backoffLimit: 0
  }
});
