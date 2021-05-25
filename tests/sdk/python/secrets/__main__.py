# Copyright 2016-2021, Pulumi Corporation.
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
import base64
import random
import string

import pulumi
from pulumi_kubernetes.core.v1 import ConfigMap, ConfigMapArgs, Secret, SecretArgs
from pulumi_kubernetes.yaml import ConfigGroup

conf = pulumi.Config("")
pw = conf.require_secret("message")
rawPW = conf.require("message")

cm_data = ConfigMap(
    "cmdata",
    data={"password": pw}
)

cm_binary_data = ConfigMap(
    "cmbinarydata",
    binary_data={"password": pw.apply(lambda x: base64.b64encode(x.encode('ascii')).decode('utf-8'))}
)

s_string_data = Secret(
    "sstringdata",
    string_data={"password": rawPW}
)

s_data = Secret(
    "sdata",
    data={"password": base64.b64encode(rawPW.encode('ascii')).decode('utf-8')}
)

suffix = ''.join(random.choice(string.ascii_lowercase) for i in range(5))
name = f'test-{suffix}'
secret_yaml = f'''
apiVersion: v1
kind: Secret
metadata:
  name: {name}
stringData:
  password: {rawPW}
'''

cg = ConfigGroup(
    "example",
    yaml=[secret_yaml]
)
cg_secret = cg.get_resource("v1/Secret", name)

pulumi.export("cmData", cm_data.data)
pulumi.export("cmBinaryData", cm_data.binary_data)
pulumi.export("sStringData", s_string_data.string_data)
pulumi.export("sData", s_data.data)
pulumi.export("cgData", cg_secret.data)
