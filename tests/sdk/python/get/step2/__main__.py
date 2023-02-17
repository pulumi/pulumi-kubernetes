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

import pulumi
from pulumi_kubernetes.apiextensions.CustomResource import CustomResource
from pulumi_kubernetes.apiextensions.v1.CustomResourceDefinition import (
    CustomResourceDefinition,
    CustomResourceDefinitionNamesArgs,
    CustomResourceDefinitionSpecArgs,
    CustomResourceDefinitionVersionArgs,
    CustomResourceValidationArgs,
    JSONSchemaPropsArgs,
)
from pulumi_kubernetes.core.v1 import Service
from pulumi_kubernetes.core.v1.Namespace import Namespace
from pulumi_kubernetes.meta.v1 import ObjectMetaArgs

service = Service.get("kube-api", "kubernetes")

crd = CustomResourceDefinition(
    resource_name="foo",
    metadata=ObjectMetaArgs(name="gettests.python.test"),
    spec=CustomResourceDefinitionSpecArgs(
        group="python.test",
        versions=[
            CustomResourceDefinitionVersionArgs(
                name="v1",
                served=True,
                storage=True,
                schema=CustomResourceValidationArgs(
                    open_apiv3_schema=JSONSchemaPropsArgs(
                        type="object",
                        properties={
                            "spec": JSONSchemaPropsArgs(
                                type="object",
                                properties={
                                    "foo": JSONSchemaPropsArgs(
                                        type="string",
                                    )
                                }
                            ),
                        },
                    ),
                ),
            ),
        ],
        scope="Namespaced",
        names=CustomResourceDefinitionNamesArgs(
            plural="gettests",
            singular="gettest",
            kind="GetTest",
        ),
    ))

ns = Namespace("ns")

cr = CustomResource(
    resource_name="foo",
    api_version="python.test/v1",
    kind="GetTest",
    metadata={"namespace": ns.metadata["name"]},
    spec={"foo": "bar"},
    opts=pulumi.ResourceOptions(depends_on=[crd]))

cr_get = CustomResource.get(resource_name="bar", api_version="python.test/v1", kind="GetTest", id=cr.id)
