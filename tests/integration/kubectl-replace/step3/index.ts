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

import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

//
// No changes to Kubernetes resources.
//

const namespace = new k8s.core.v1.Namespace("test-namespace");

const cm = new k8s.core.v1.ConfigMap("game-config", {
    metadata: {
        name: "game-config",
        namespace: namespace.metadata.name,
    },
    data: {
        "game.properties":
            "enemies=aliens\nlives=3\nenemies.cheat=true\nenemies.cheat.level=noGoodRotten\nsecret.code.passphrase=UUDDLRLRBABAS\nsecret.code.allowed=true\nsecret.code.lives=30\n",
        "ui.properties":
            "color.good=purple\ncolor.bad=yellow\nallow.textmode=true\nhow.nice.to.look=fairlyNice",
    },
});

//
// Now run `kubectl replace` again. This time, it should work.
//

export const replaced = pulumi.all([namespace.metadata.name, cm.metadata.name]).apply(([ns]) =>
    pulumi.runtime.invoke("kubernetes:kubernetes:kubectlReplace", {
        apiVersion: "v1",
        kind: "ConfigMap",
        metadata: {
            name: "game-config",
            namespace: ns,
        },
        data: { foo: "bar" },
    }),
);
