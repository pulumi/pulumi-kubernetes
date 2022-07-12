#  Copyright 2016-2022, Pulumi Corporation.
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
from os import path

import pulumi
import pulumi_kubernetes as k8s

kubeconfig_file = path.join(path.expanduser("~"), ".kube", "config")
with open(kubeconfig_file) as f:
    kubeconfig = f.read()

provider = k8s.Provider("myk8s", kubeconfig=kubeconfig, enable_server_side_apply=True)

ns = k8s.core.v1.Namespace("test", opts=pulumi.ResourceOptions(provider=provider))

crd = k8s.apiextensions.v1.CustomResourceDefinition(
    "crd",
    args=k8s.apiextensions.v1.CustomResourceDefinitionInitArgs(
        metadata=k8s.meta.v1.ObjectMetaArgs(
            name="tests.pyssa.example.com",
            namespace=ns.metadata.name,
        ),
        spec=k8s.apiextensions.v1.CustomResourceDefinitionSpecArgs(
            group="pyssa.example.com",
            versions=[
                k8s.apiextensions.v1.CustomResourceDefinitionVersionArgs(
                    name="v1",
                    served=True,
                    storage=True,
                    schema=k8s.apiextensions.v1.CustomResourceValidationArgs(
                        open_apiv3_schema=k8s.apiextensions.v1.JSONSchemaPropsArgs(
                            type="object",
                            properties={
                                "spec": k8s.apiextensions.v1.JSONSchemaPropsArgs(
                                    type="object",
                                    properties={
                                        "foo": k8s.apiextensions.v1.JSONSchemaPropsArgs(
                                            type="string",
                                        ),
                                    },
                                ),
                            },
                        ),
                    ),
                ),
            ],
            scope="Namespaced",
            names=k8s.apiextensions.v1.CustomResourceDefinitionNamesArgs(
                plural="tests",
                singular="test",
                kind="Test",
            ),
        ),
    ),
    opts=pulumi.ResourceOptions(provider=provider)
)

cr = k8s.apiextensions.CustomResource(
    "cr",
    api_version="pyssa.example.com/v1",
    kind="Test",
    metadata=k8s.meta.v1.ObjectMetaArgs(
        name="foo",
        namespace=ns.metadata.name,
    ),
    spec={
        "foo": "bar",
    },
    opts=pulumi.ResourceOptions(provider=provider, depends_on=[crd])
)

cr_patch = k8s.apiextensions.CustomResourcePatch(
    "cr_patch",
    api_version="pyssa.example.com/v1",
    kind="Test",
    metadata=k8s.meta.v1.ObjectMetaArgs(
        labels={
            "foo": "foo"
        },
        # name=cr.metadata.name, # TODO: CustomResource.metadata is not working. Switch to auto-name once fixed.
        name="foo",
        namespace=ns.metadata.name,
    ),
    opts=pulumi.ResourceOptions(provider=provider, depends_on=[cr])
)

pulumi.export("crPatched", cr_patch)
