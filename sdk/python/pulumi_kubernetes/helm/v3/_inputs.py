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

__all__ = [
    'RepositoryOptsArgs',
    'RepositoryOptsArgsDict',
]

MYPY = False

if not MYPY:
    class RepositoryOptsArgsDict(TypedDict):
        """
        Specification defining the Helm chart repository to use.
        """
        ca_file: NotRequired[pulumi.Input[_builtins.str]]
        """
        The Repository's CA File
        """
        cert_file: NotRequired[pulumi.Input[_builtins.str]]
        """
        The repository's cert file
        """
        key_file: NotRequired[pulumi.Input[_builtins.str]]
        """
        The repository's cert key file
        """
        password: NotRequired[pulumi.Input[_builtins.str]]
        """
        Password for HTTP basic authentication
        """
        repo: NotRequired[pulumi.Input[_builtins.str]]
        """
        Repository where to locate the requested chart. If it's a URL the chart is installed without installing the repository.
        """
        username: NotRequired[pulumi.Input[_builtins.str]]
        """
        Username for HTTP basic authentication
        """
elif False:
    RepositoryOptsArgsDict: TypeAlias = Mapping[str, Any]

@pulumi.input_type
class RepositoryOptsArgs:
    def __init__(__self__, *,
                 ca_file: Optional[pulumi.Input[_builtins.str]] = None,
                 cert_file: Optional[pulumi.Input[_builtins.str]] = None,
                 key_file: Optional[pulumi.Input[_builtins.str]] = None,
                 password: Optional[pulumi.Input[_builtins.str]] = None,
                 repo: Optional[pulumi.Input[_builtins.str]] = None,
                 username: Optional[pulumi.Input[_builtins.str]] = None):
        """
        Specification defining the Helm chart repository to use.
        :param pulumi.Input[_builtins.str] ca_file: The Repository's CA File
        :param pulumi.Input[_builtins.str] cert_file: The repository's cert file
        :param pulumi.Input[_builtins.str] key_file: The repository's cert key file
        :param pulumi.Input[_builtins.str] password: Password for HTTP basic authentication
        :param pulumi.Input[_builtins.str] repo: Repository where to locate the requested chart. If it's a URL the chart is installed without installing the repository.
        :param pulumi.Input[_builtins.str] username: Username for HTTP basic authentication
        """
        if ca_file is not None:
            pulumi.set(__self__, "ca_file", ca_file)
        if cert_file is not None:
            pulumi.set(__self__, "cert_file", cert_file)
        if key_file is not None:
            pulumi.set(__self__, "key_file", key_file)
        if password is not None:
            pulumi.set(__self__, "password", password)
        if repo is not None:
            pulumi.set(__self__, "repo", repo)
        if username is not None:
            pulumi.set(__self__, "username", username)

    @_builtins.property
    @pulumi.getter(name="caFile")
    def ca_file(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        The Repository's CA File
        """
        return pulumi.get(self, "ca_file")

    @ca_file.setter
    def ca_file(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "ca_file", value)

    @_builtins.property
    @pulumi.getter(name="certFile")
    def cert_file(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        The repository's cert file
        """
        return pulumi.get(self, "cert_file")

    @cert_file.setter
    def cert_file(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "cert_file", value)

    @_builtins.property
    @pulumi.getter(name="keyFile")
    def key_file(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        The repository's cert key file
        """
        return pulumi.get(self, "key_file")

    @key_file.setter
    def key_file(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "key_file", value)

    @_builtins.property
    @pulumi.getter
    def password(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        Password for HTTP basic authentication
        """
        return pulumi.get(self, "password")

    @password.setter
    def password(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "password", value)

    @_builtins.property
    @pulumi.getter
    def repo(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        Repository where to locate the requested chart. If it's a URL the chart is installed without installing the repository.
        """
        return pulumi.get(self, "repo")

    @repo.setter
    def repo(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "repo", value)

    @_builtins.property
    @pulumi.getter
    def username(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        Username for HTTP basic authentication
        """
        return pulumi.get(self, "username")

    @username.setter
    def username(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "username", value)


