# coding=utf-8
# *** WARNING: this file was generated by pulumigen. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import builtins as _builtins
import warnings
import sys
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
if sys.version_info >= (3, 11):
    from typing import NotRequired, TypedDict, TypeAlias
else:
    from typing_extensions import NotRequired, TypedDict, TypeAlias
from ... import _utilities
from . import outputs
from ... import meta as _meta
from ._inputs import *

__all__ = ['ClusterRoleBindingListArgs', 'ClusterRoleBindingList']

@pulumi.input_type
class ClusterRoleBindingListArgs:
    def __init__(__self__, *,
                 items: pulumi.Input[Sequence[pulumi.Input['ClusterRoleBindingArgs']]],
                 api_version: Optional[pulumi.Input[_builtins.str]] = None,
                 kind: Optional[pulumi.Input[_builtins.str]] = None,
                 metadata: Optional[pulumi.Input['_meta.v1.ListMetaArgs']] = None):
        """
        The set of arguments for constructing a ClusterRoleBindingList resource.
        :param pulumi.Input[Sequence[pulumi.Input['ClusterRoleBindingArgs']]] items: Items is a list of ClusterRoleBindings
        :param pulumi.Input[_builtins.str] api_version: APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        :param pulumi.Input[_builtins.str] kind: Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        :param pulumi.Input['_meta.v1.ListMetaArgs'] metadata: Standard object's metadata.
        """
        pulumi.set(__self__, "items", items)
        if api_version is not None:
            pulumi.set(__self__, "api_version", 'rbac.authorization.k8s.io/v1alpha1')
        if kind is not None:
            pulumi.set(__self__, "kind", 'ClusterRoleBindingList')
        if metadata is not None:
            pulumi.set(__self__, "metadata", metadata)

    @_builtins.property
    @pulumi.getter
    def items(self) -> pulumi.Input[Sequence[pulumi.Input['ClusterRoleBindingArgs']]]:
        """
        Items is a list of ClusterRoleBindings
        """
        return pulumi.get(self, "items")

    @items.setter
    def items(self, value: pulumi.Input[Sequence[pulumi.Input['ClusterRoleBindingArgs']]]):
        pulumi.set(self, "items", value)

    @_builtins.property
    @pulumi.getter(name="apiVersion")
    def api_version(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        """
        return pulumi.get(self, "api_version")

    @api_version.setter
    def api_version(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "api_version", value)

    @_builtins.property
    @pulumi.getter
    def kind(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        """
        return pulumi.get(self, "kind")

    @kind.setter
    def kind(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "kind", value)

    @_builtins.property
    @pulumi.getter
    def metadata(self) -> Optional[pulumi.Input['_meta.v1.ListMetaArgs']]:
        """
        Standard object's metadata.
        """
        return pulumi.get(self, "metadata")

    @metadata.setter
    def metadata(self, value: Optional[pulumi.Input['_meta.v1.ListMetaArgs']]):
        pulumi.set(self, "metadata", value)


@pulumi.type_token("kubernetes:rbac.authorization.k8s.io/v1alpha1:ClusterRoleBindingList")
class ClusterRoleBindingList(pulumi.CustomResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 api_version: Optional[pulumi.Input[_builtins.str]] = None,
                 items: Optional[pulumi.Input[Sequence[pulumi.Input[Union['ClusterRoleBindingArgs', 'ClusterRoleBindingArgsDict']]]]] = None,
                 kind: Optional[pulumi.Input[_builtins.str]] = None,
                 metadata: Optional[pulumi.Input[Union['_meta.v1.ListMetaArgs', '_meta.v1.ListMetaArgsDict']]] = None,
                 __props__=None):
        """
        ClusterRoleBindingList is a collection of ClusterRoleBindings. Deprecated in v1.17 in favor of rbac.authorization.k8s.io/v1 ClusterRoleBindings, and will no longer be served in v1.20.

        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        :param pulumi.Input[_builtins.str] api_version: APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        :param pulumi.Input[Sequence[pulumi.Input[Union['ClusterRoleBindingArgs', 'ClusterRoleBindingArgsDict']]]] items: Items is a list of ClusterRoleBindings
        :param pulumi.Input[_builtins.str] kind: Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        :param pulumi.Input[Union['_meta.v1.ListMetaArgs', '_meta.v1.ListMetaArgsDict']] metadata: Standard object's metadata.
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: ClusterRoleBindingListArgs,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        ClusterRoleBindingList is a collection of ClusterRoleBindings. Deprecated in v1.17 in favor of rbac.authorization.k8s.io/v1 ClusterRoleBindings, and will no longer be served in v1.20.

        :param str resource_name: The name of the resource.
        :param ClusterRoleBindingListArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(ClusterRoleBindingListArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 api_version: Optional[pulumi.Input[_builtins.str]] = None,
                 items: Optional[pulumi.Input[Sequence[pulumi.Input[Union['ClusterRoleBindingArgs', 'ClusterRoleBindingArgsDict']]]]] = None,
                 kind: Optional[pulumi.Input[_builtins.str]] = None,
                 metadata: Optional[pulumi.Input[Union['_meta.v1.ListMetaArgs', '_meta.v1.ListMetaArgsDict']]] = None,
                 __props__=None):
        opts = pulumi.ResourceOptions.merge(_utilities.get_resource_opts_defaults(), opts)
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.id is None:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = ClusterRoleBindingListArgs.__new__(ClusterRoleBindingListArgs)

            __props__.__dict__["api_version"] = 'rbac.authorization.k8s.io/v1alpha1'
            if items is None and not opts.urn:
                raise TypeError("Missing required property 'items'")
            __props__.__dict__["items"] = items
            __props__.__dict__["kind"] = 'ClusterRoleBindingList'
            __props__.__dict__["metadata"] = metadata
        super(ClusterRoleBindingList, __self__).__init__(
            'kubernetes:rbac.authorization.k8s.io/v1alpha1:ClusterRoleBindingList',
            resource_name,
            __props__,
            opts)

    @staticmethod
    def get(resource_name: str,
            id: pulumi.Input[str],
            opts: Optional[pulumi.ResourceOptions] = None) -> 'ClusterRoleBindingList':
        """
        Get an existing ClusterRoleBindingList resource's state with the given name, id, and optional extra
        properties used to qualify the lookup.

        :param str resource_name: The unique name of the resulting resource.
        :param pulumi.Input[str] id: The unique provider ID of the resource to lookup.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        opts = pulumi.ResourceOptions.merge(opts, pulumi.ResourceOptions(id=id))

        __props__ = ClusterRoleBindingListArgs.__new__(ClusterRoleBindingListArgs)

        __props__.__dict__["api_version"] = None
        __props__.__dict__["items"] = None
        __props__.__dict__["kind"] = None
        __props__.__dict__["metadata"] = None
        return ClusterRoleBindingList(resource_name, opts=opts, __props__=__props__)

    @_builtins.property
    @pulumi.getter(name="apiVersion")
    def api_version(self) -> pulumi.Output[_builtins.str]:
        """
        APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        """
        return pulumi.get(self, "api_version")

    @_builtins.property
    @pulumi.getter
    def items(self) -> pulumi.Output[Sequence['outputs.ClusterRoleBinding']]:
        """
        Items is a list of ClusterRoleBindings
        """
        return pulumi.get(self, "items")

    @_builtins.property
    @pulumi.getter
    def kind(self) -> pulumi.Output[_builtins.str]:
        """
        Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        """
        return pulumi.get(self, "kind")

    @_builtins.property
    @pulumi.getter
    def metadata(self) -> pulumi.Output['_meta.v1.outputs.ListMeta']:
        """
        Standard object's metadata.
        """
        return pulumi.get(self, "metadata")

