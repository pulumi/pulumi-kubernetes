# coding=utf-8
# *** WARNING: this file was generated by pulumigen. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import copy
import warnings
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
from ... import _utilities

__all__ = ['ConfigFileArgs', 'ConfigFile']

@pulumi.input_type
class ConfigFileArgs:
    def __init__(__self__, *,
                 file: pulumi.Input[str],
                 resource_prefix: Optional[pulumi.Input[str]] = None,
                 skip_await: Optional[pulumi.Input[bool]] = None):
        """
        The set of arguments for constructing a ConfigFile resource.
        :param pulumi.Input[str] file: Path or a URL that uniquely identifies a file.
        :param pulumi.Input[str] resource_prefix: An optional prefix for the auto-generated resource names. Example: A resource created with resourcePrefix="foo" would produce a resource named "foo-resourceName".
        :param pulumi.Input[bool] skip_await: Indicates that child resources should skip the await logic.
        """
        pulumi.set(__self__, "file", file)
        if resource_prefix is not None:
            pulumi.set(__self__, "resource_prefix", resource_prefix)
        if skip_await is not None:
            pulumi.set(__self__, "skip_await", skip_await)

    @property
    @pulumi.getter
    def file(self) -> pulumi.Input[str]:
        """
        Path or a URL that uniquely identifies a file.
        """
        return pulumi.get(self, "file")

    @file.setter
    def file(self, value: pulumi.Input[str]):
        pulumi.set(self, "file", value)

    @property
    @pulumi.getter(name="resourcePrefix")
    def resource_prefix(self) -> Optional[pulumi.Input[str]]:
        """
        An optional prefix for the auto-generated resource names. Example: A resource created with resourcePrefix="foo" would produce a resource named "foo-resourceName".
        """
        return pulumi.get(self, "resource_prefix")

    @resource_prefix.setter
    def resource_prefix(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "resource_prefix", value)

    @property
    @pulumi.getter(name="skipAwait")
    def skip_await(self) -> Optional[pulumi.Input[bool]]:
        """
        Indicates that child resources should skip the await logic.
        """
        return pulumi.get(self, "skip_await")

    @skip_await.setter
    def skip_await(self, value: Optional[pulumi.Input[bool]]):
        pulumi.set(self, "skip_await", value)


class ConfigFile(pulumi.ComponentResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 file: Optional[pulumi.Input[str]] = None,
                 resource_prefix: Optional[pulumi.Input[str]] = None,
                 skip_await: Optional[pulumi.Input[bool]] = None,
                 __props__=None):
        """
        ConfigFile creates a set of Kubernetes resources from a Kubernetes YAML file.

        ## Example Usage
        ### Local File
        ```python
        from pulumi_kubernetes.yaml.v2 import ConfigFile

        example = ConfigFile(
            "example",
            file="foo.yaml",
        )
        ```
        {% /examples %}}

        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        :param pulumi.Input[str] file: Path or a URL that uniquely identifies a file.
        :param pulumi.Input[str] resource_prefix: An optional prefix for the auto-generated resource names. Example: A resource created with resourcePrefix="foo" would produce a resource named "foo-resourceName".
        :param pulumi.Input[bool] skip_await: Indicates that child resources should skip the await logic.
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: ConfigFileArgs,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        ConfigFile creates a set of Kubernetes resources from a Kubernetes YAML file.

        ## Example Usage
        ### Local File
        ```python
        from pulumi_kubernetes.yaml.v2 import ConfigFile

        example = ConfigFile(
            "example",
            file="foo.yaml",
        )
        ```
        {% /examples %}}

        :param str resource_name: The name of the resource.
        :param ConfigFileArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(ConfigFileArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 file: Optional[pulumi.Input[str]] = None,
                 resource_prefix: Optional[pulumi.Input[str]] = None,
                 skip_await: Optional[pulumi.Input[bool]] = None,
                 __props__=None):
        opts = pulumi.ResourceOptions.merge(_utilities.get_resource_opts_defaults(), opts)
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.id is not None:
            raise ValueError('ComponentResource classes do not support opts.id')
        else:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = ConfigFileArgs.__new__(ConfigFileArgs)

            if file is None and not opts.urn:
                raise TypeError("Missing required property 'file'")
            __props__.__dict__["file"] = file
            __props__.__dict__["resource_prefix"] = resource_prefix
            __props__.__dict__["skip_await"] = skip_await
            __props__.__dict__["resources"] = None
        super(ConfigFile, __self__).__init__(
            'kubernetes:yaml/v2:ConfigFile',
            resource_name,
            __props__,
            opts,
            remote=True)

    @property
    @pulumi.getter
    def resources(self) -> pulumi.Output[Optional[Sequence[Any]]]:
        """
        Resources created by the ConfigFile.
        """
        return pulumi.get(self, "resources")
