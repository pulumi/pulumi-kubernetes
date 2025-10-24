// Copyright 2016-2022, Pulumi Corporation.
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

const config = new pulumi.Config();

const pw = config.requireSecret("message");
const rawPW = config.require("message");

const provider = new k8s.Provider("k8s");

const cmData = new k8s.core.v1.ConfigMap("cmdata", {
    data: {
        password: pw,
    }
}, {provider});

const cmBinaryData = new k8s.core.v1.ConfigMap("cmbinarydata", {
    binaryData: {
        password: pw.apply(d => Buffer.from(d).toString("base64")),
    }
}, {provider});

const ssStringData = new k8s.core.v1.Secret("ssstringdata", {
    stringData: {
        password: rawPW,
    }
}, {provider});

const ssData = new k8s.core.v1.Secret("ssdata", {
    data: {
        password: Buffer.from(rawPW).toString("base64"),
    }
}, {provider});

const randSuffix = Math.random().toString(36).substring(7);
const name = `test-${randSuffix}`;

// Create a Secret resource directly
const cgSecret = new k8s.core.v1.Secret("cgSecret", {
    metadata: {
        name: name,
    },
    stringData: {
        password: rawPW,
    }
}, {provider});

export const cmDataData = cmData.data;
export const cmBinaryDataData = cmBinaryData.binaryData;
export const ssStringDataStringData = ssStringData.stringData;
export const ssStringDataData = ssStringData.data;
export const ssDataData = ssData.data;
export const cgSecretStringData = cgSecret.stringData;
