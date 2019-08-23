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

const pw = (new pulumi.Config()).requireSecret("message");

const cmData = new k8s.core.v1.ConfigMap("cmdata", {
    data: {
        password: pw,
    }
})

const cmBinaryData = new k8s.core.v1.ConfigMap("cmbinarydata", {
    binaryData: {
        password: pw.apply(d => new Buffer(d).toString("base64")),
    }
})

const ssStringData = new k8s.core.v1.Secret("ssstringdata", {
    stringData: {
        password: pw,
    }
})

const ssData = new k8s.core.v1.Secret("ssdata", {
    data: {
        password: pw.apply(d => new Buffer(d).toString("base64")),
    }
})

export const cmDataData = cmData.data;
export const cmBinaryDataData = cmBinaryData.binaryData;
export const ssStringDataData = ssStringData.data;
export const ssDataData = ssData.data;
