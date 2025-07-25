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
from ... import meta as _meta

__all__ = ['SecretPatchArgs', 'SecretPatch']

@pulumi.input_type
class SecretPatchArgs:
    def __init__(__self__, *,
                 api_version: Optional[pulumi.Input[_builtins.str]] = None,
                 data: Optional[pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]]] = None,
                 immutable: Optional[pulumi.Input[_builtins.bool]] = None,
                 kind: Optional[pulumi.Input[_builtins.str]] = None,
                 metadata: Optional[pulumi.Input['_meta.v1.ObjectMetaPatchArgs']] = None,
                 string_data: Optional[pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]]] = None,
                 type: Optional[pulumi.Input[_builtins.str]] = None):
        """
        The set of arguments for constructing a SecretPatch resource.
        :param pulumi.Input[_builtins.str] api_version: APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        :param pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]] data: Data contains the secret data. Each key must consist of alphanumeric characters, '-', '_' or '.'. The serialized form of the secret data is a base64 encoded string, representing the arbitrary (possibly non-string) data value here. Described in https://tools.ietf.org/html/rfc4648#section-4
        :param pulumi.Input[_builtins.bool] immutable: Immutable, if set to true, ensures that data stored in the Secret cannot be updated (only object metadata can be modified). If not set to true, the field can be modified at any time. Defaulted to nil.
        :param pulumi.Input[_builtins.str] kind: Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        :param pulumi.Input['_meta.v1.ObjectMetaPatchArgs'] metadata: Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        :param pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]] string_data: stringData allows specifying non-binary secret data in string form. It is provided as a write-only input field for convenience. All keys and values are merged into the data field on write, overwriting any existing values. The stringData field is never output when reading from the API.
        :param pulumi.Input[_builtins.str] type: Used to facilitate programmatic handling of secret data. More info: https://kubernetes.io/docs/concepts/configuration/secret/#secret-types
        """
        if api_version is not None:
            pulumi.set(__self__, "api_version", 'v1')
        if data is not None:
            pulumi.set(__self__, "data", data)
        if immutable is not None:
            pulumi.set(__self__, "immutable", immutable)
        if kind is not None:
            pulumi.set(__self__, "kind", 'Secret')
        if metadata is not None:
            pulumi.set(__self__, "metadata", metadata)
        if string_data is not None:
            pulumi.set(__self__, "string_data", string_data)
        if type is not None:
            pulumi.set(__self__, "type", type)

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
    def data(self) -> Optional[pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]]]:
        """
        Data contains the secret data. Each key must consist of alphanumeric characters, '-', '_' or '.'. The serialized form of the secret data is a base64 encoded string, representing the arbitrary (possibly non-string) data value here. Described in https://tools.ietf.org/html/rfc4648#section-4
        """
        return pulumi.get(self, "data")

    @data.setter
    def data(self, value: Optional[pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]]]):
        pulumi.set(self, "data", value)

    @_builtins.property
    @pulumi.getter
    def immutable(self) -> Optional[pulumi.Input[_builtins.bool]]:
        """
        Immutable, if set to true, ensures that data stored in the Secret cannot be updated (only object metadata can be modified). If not set to true, the field can be modified at any time. Defaulted to nil.
        """
        return pulumi.get(self, "immutable")

    @immutable.setter
    def immutable(self, value: Optional[pulumi.Input[_builtins.bool]]):
        pulumi.set(self, "immutable", value)

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
    def metadata(self) -> Optional[pulumi.Input['_meta.v1.ObjectMetaPatchArgs']]:
        """
        Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        """
        return pulumi.get(self, "metadata")

    @metadata.setter
    def metadata(self, value: Optional[pulumi.Input['_meta.v1.ObjectMetaPatchArgs']]):
        pulumi.set(self, "metadata", value)

    @_builtins.property
    @pulumi.getter(name="stringData")
    def string_data(self) -> Optional[pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]]]:
        """
        stringData allows specifying non-binary secret data in string form. It is provided as a write-only input field for convenience. All keys and values are merged into the data field on write, overwriting any existing values. The stringData field is never output when reading from the API.
        """
        return pulumi.get(self, "string_data")

    @string_data.setter
    def string_data(self, value: Optional[pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]]]):
        pulumi.set(self, "string_data", value)

    @_builtins.property
    @pulumi.getter
    def type(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        Used to facilitate programmatic handling of secret data. More info: https://kubernetes.io/docs/concepts/configuration/secret/#secret-types
        """
        return pulumi.get(self, "type")

    @type.setter
    def type(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "type", value)


@pulumi.type_token("kubernetes:core/v1:SecretPatch")
class SecretPatch(pulumi.CustomResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 api_version: Optional[pulumi.Input[_builtins.str]] = None,
                 data: Optional[pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]]] = None,
                 immutable: Optional[pulumi.Input[_builtins.bool]] = None,
                 kind: Optional[pulumi.Input[_builtins.str]] = None,
                 metadata: Optional[pulumi.Input[Union['_meta.v1.ObjectMetaPatchArgs', '_meta.v1.ObjectMetaPatchArgsDict']]] = None,
                 string_data: Optional[pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]]] = None,
                 type: Optional[pulumi.Input[_builtins.str]] = None,
                 __props__=None):
        """
        Patch resources are used to modify existing Kubernetes resources by using
        Server-Side Apply updates. The name of the resource must be specified, but all other properties are optional. More than
        one patch may be applied to the same resource, and a random FieldManager name will be used for each Patch resource.
        Conflicts will result in an error by default, but can be forced using the "pulumi.com/patchForce" annotation. See the
        [Server-Side Apply Docs](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/managing-resources-with-server-side-apply/) for
        additional information about using Server-Side Apply to manage Kubernetes resources with Pulumi.
        Secret holds secret data of a certain type. The total bytes of the values in the Data field must be less than MaxSecretSize bytes.

        Note: While Pulumi automatically encrypts the 'data' and 'stringData'
        fields, this encryption only applies to Pulumi's context, including the state file,
        the Service, the CLI, etc. Kubernetes does not encrypt Secret resources by default,
        and the contents are visible to users with access to the Secret in Kubernetes using
        tools like 'kubectl'.

        For more information on securing Kubernetes Secrets, see the following links:
        https://kubernetes.io/docs/concepts/configuration/secret/#security-properties
        https://kubernetes.io/docs/concepts/configuration/secret/#risks

        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        :param pulumi.Input[_builtins.str] api_version: APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        :param pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]] data: Data contains the secret data. Each key must consist of alphanumeric characters, '-', '_' or '.'. The serialized form of the secret data is a base64 encoded string, representing the arbitrary (possibly non-string) data value here. Described in https://tools.ietf.org/html/rfc4648#section-4
        :param pulumi.Input[_builtins.bool] immutable: Immutable, if set to true, ensures that data stored in the Secret cannot be updated (only object metadata can be modified). If not set to true, the field can be modified at any time. Defaulted to nil.
        :param pulumi.Input[_builtins.str] kind: Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        :param pulumi.Input[Union['_meta.v1.ObjectMetaPatchArgs', '_meta.v1.ObjectMetaPatchArgsDict']] metadata: Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        :param pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]] string_data: stringData allows specifying non-binary secret data in string form. It is provided as a write-only input field for convenience. All keys and values are merged into the data field on write, overwriting any existing values. The stringData field is never output when reading from the API.
        :param pulumi.Input[_builtins.str] type: Used to facilitate programmatic handling of secret data. More info: https://kubernetes.io/docs/concepts/configuration/secret/#secret-types
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: Optional[SecretPatchArgs] = None,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        Patch resources are used to modify existing Kubernetes resources by using
        Server-Side Apply updates. The name of the resource must be specified, but all other properties are optional. More than
        one patch may be applied to the same resource, and a random FieldManager name will be used for each Patch resource.
        Conflicts will result in an error by default, but can be forced using the "pulumi.com/patchForce" annotation. See the
        [Server-Side Apply Docs](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/managing-resources-with-server-side-apply/) for
        additional information about using Server-Side Apply to manage Kubernetes resources with Pulumi.
        Secret holds secret data of a certain type. The total bytes of the values in the Data field must be less than MaxSecretSize bytes.

        Note: While Pulumi automatically encrypts the 'data' and 'stringData'
        fields, this encryption only applies to Pulumi's context, including the state file,
        the Service, the CLI, etc. Kubernetes does not encrypt Secret resources by default,
        and the contents are visible to users with access to the Secret in Kubernetes using
        tools like 'kubectl'.

        For more information on securing Kubernetes Secrets, see the following links:
        https://kubernetes.io/docs/concepts/configuration/secret/#security-properties
        https://kubernetes.io/docs/concepts/configuration/secret/#risks

        :param str resource_name: The name of the resource.
        :param SecretPatchArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(SecretPatchArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 api_version: Optional[pulumi.Input[_builtins.str]] = None,
                 data: Optional[pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]]] = None,
                 immutable: Optional[pulumi.Input[_builtins.bool]] = None,
                 kind: Optional[pulumi.Input[_builtins.str]] = None,
                 metadata: Optional[pulumi.Input[Union['_meta.v1.ObjectMetaPatchArgs', '_meta.v1.ObjectMetaPatchArgsDict']]] = None,
                 string_data: Optional[pulumi.Input[Mapping[str, pulumi.Input[_builtins.str]]]] = None,
                 type: Optional[pulumi.Input[_builtins.str]] = None,
                 __props__=None):
        opts = pulumi.ResourceOptions.merge(_utilities.get_resource_opts_defaults(), opts)
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.id is None:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = SecretPatchArgs.__new__(SecretPatchArgs)

            __props__.__dict__["api_version"] = 'v1'
            __props__.__dict__["data"] = data
            __props__.__dict__["immutable"] = immutable
            __props__.__dict__["kind"] = 'Secret'
            __props__.__dict__["metadata"] = metadata
            __props__.__dict__["string_data"] = string_data
            __props__.__dict__["type"] = type
        super(SecretPatch, __self__).__init__(
            'kubernetes:core/v1:SecretPatch',
            resource_name,
            __props__,
            opts)

    @staticmethod
    def get(resource_name: str,
            id: pulumi.Input[str],
            opts: Optional[pulumi.ResourceOptions] = None) -> 'SecretPatch':
        """
        Get an existing SecretPatch resource's state with the given name, id, and optional extra
        properties used to qualify the lookup.

        :param str resource_name: The unique name of the resulting resource.
        :param pulumi.Input[str] id: The unique provider ID of the resource to lookup.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        opts = pulumi.ResourceOptions.merge(opts, pulumi.ResourceOptions(id=id))

        __props__ = SecretPatchArgs.__new__(SecretPatchArgs)

        __props__.__dict__["api_version"] = None
        __props__.__dict__["data"] = None
        __props__.__dict__["immutable"] = None
        __props__.__dict__["kind"] = None
        __props__.__dict__["metadata"] = None
        __props__.__dict__["string_data"] = None
        __props__.__dict__["type"] = None
        return SecretPatch(resource_name, opts=opts, __props__=__props__)

    @_builtins.property
    @pulumi.getter(name="apiVersion")
    def api_version(self) -> pulumi.Output[Optional[_builtins.str]]:
        """
        APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        """
        return pulumi.get(self, "api_version")

    @_builtins.property
    @pulumi.getter
    def data(self) -> pulumi.Output[Optional[Mapping[str, _builtins.str]]]:
        """
        Data contains the secret data. Each key must consist of alphanumeric characters, '-', '_' or '.'. The serialized form of the secret data is a base64 encoded string, representing the arbitrary (possibly non-string) data value here. Described in https://tools.ietf.org/html/rfc4648#section-4
        """
        return pulumi.get(self, "data")

    @_builtins.property
    @pulumi.getter
    def immutable(self) -> pulumi.Output[Optional[_builtins.bool]]:
        """
        Immutable, if set to true, ensures that data stored in the Secret cannot be updated (only object metadata can be modified). If not set to true, the field can be modified at any time. Defaulted to nil.
        """
        return pulumi.get(self, "immutable")

    @_builtins.property
    @pulumi.getter
    def kind(self) -> pulumi.Output[Optional[_builtins.str]]:
        """
        Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        """
        return pulumi.get(self, "kind")

    @_builtins.property
    @pulumi.getter
    def metadata(self) -> pulumi.Output[Optional['_meta.v1.outputs.ObjectMetaPatch']]:
        """
        Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        """
        return pulumi.get(self, "metadata")

    @_builtins.property
    @pulumi.getter(name="stringData")
    def string_data(self) -> pulumi.Output[Optional[Mapping[str, _builtins.str]]]:
        """
        stringData allows specifying non-binary secret data in string form. It is provided as a write-only input field for convenience. All keys and values are merged into the data field on write, overwriting any existing values. The stringData field is never output when reading from the API.
        """
        return pulumi.get(self, "string_data")

    @_builtins.property
    @pulumi.getter
    def type(self) -> pulumi.Output[Optional[_builtins.str]]:
        """
        Used to facilitate programmatic handling of secret data. More info: https://kubernetes.io/docs/concepts/configuration/secret/#secret-types
        """
        return pulumi.get(self, "type")

