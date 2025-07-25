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
from . import _utilities
from ._inputs import *

__all__ = ['ProviderArgs', 'Provider']

@pulumi.input_type
class ProviderArgs:
    def __init__(__self__, *,
                 cluster: Optional[pulumi.Input[_builtins.str]] = None,
                 cluster_identifier: Optional[pulumi.Input[_builtins.str]] = None,
                 context: Optional[pulumi.Input[_builtins.str]] = None,
                 delete_unreachable: Optional[pulumi.Input[_builtins.bool]] = None,
                 enable_config_map_mutable: Optional[pulumi.Input[_builtins.bool]] = None,
                 enable_secret_mutable: Optional[pulumi.Input[_builtins.bool]] = None,
                 enable_server_side_apply: Optional[pulumi.Input[_builtins.bool]] = None,
                 helm_release_settings: Optional[pulumi.Input['HelmReleaseSettingsArgs']] = None,
                 kube_client_settings: Optional[pulumi.Input['KubeClientSettingsArgs']] = None,
                 kubeconfig: Optional[pulumi.Input[_builtins.str]] = None,
                 namespace: Optional[pulumi.Input[_builtins.str]] = None,
                 render_yaml_to_directory: Optional[pulumi.Input[_builtins.str]] = None,
                 skip_update_unreachable: Optional[pulumi.Input[_builtins.bool]] = None,
                 suppress_deprecation_warnings: Optional[pulumi.Input[_builtins.bool]] = None,
                 suppress_helm_hook_warnings: Optional[pulumi.Input[_builtins.bool]] = None):
        """
        The set of arguments for constructing a Provider resource.
        :param pulumi.Input[_builtins.str] cluster: If present, the name of the kubeconfig cluster to use.
        :param pulumi.Input[_builtins.str] cluster_identifier: If present, this value will control the provider's replacement behavior. In particular, the provider will _only_ be replaced when `clusterIdentifier` changes; all other changes to provider configuration will be treated as updates.
               
               Kubernetes does not yet offer an API for cluster identification, so Pulumi uses heuristics to decide when a provider resource should be replaced or updated. These heuristics can sometimes lead to destructive replace operations when an update would be more appropriate, or vice versa.
               
               Use `clusterIdentifier` for more fine-grained control of the provider resource's lifecycle.
        :param pulumi.Input[_builtins.str] context: If present, the name of the kubeconfig context to use.
        :param pulumi.Input[_builtins.bool] delete_unreachable: If present and set to true, the provider will delete resources associated with an unreachable Kubernetes cluster from Pulumi state
        :param pulumi.Input[_builtins.bool] enable_config_map_mutable: BETA FEATURE - If present and set to true, allow ConfigMaps to be mutated.
               This feature is in developer preview, and is disabled by default.
               
               This config can be specified in the following ways using this precedence:
               1. This `enableConfigMapMutable` parameter.
               2. The `PULUMI_K8S_ENABLE_CONFIGMAP_MUTABLE` environment variable.
        :param pulumi.Input[_builtins.bool] enable_secret_mutable: BETA FEATURE - If present and set to true, allow Secrets to be mutated.
               This feature is in developer preview, and is disabled by default.
               
               This config can be specified in the following ways using this precedence:
               1. This `enableSecretMutable` parameter.
               2. The `PULUMI_K8S_ENABLE_SECRET_MUTABLE` environment variable.
        :param pulumi.Input[_builtins.bool] enable_server_side_apply: If present and set to false, disable Server-Side Apply mode.
               See https://github.com/pulumi/pulumi-kubernetes/issues/2011 for additional details.
        :param pulumi.Input['HelmReleaseSettingsArgs'] helm_release_settings: Options to configure the Helm Release resource.
        :param pulumi.Input['KubeClientSettingsArgs'] kube_client_settings: Options for tuning the Kubernetes client used by a Provider.
        :param pulumi.Input[_builtins.str] kubeconfig: The contents of a kubeconfig file or the path to a kubeconfig file.
        :param pulumi.Input[_builtins.str] namespace: If present, the default namespace to use. This flag is ignored for cluster-scoped resources.
               
               A namespace can be specified in multiple places, and the precedence is as follows:
               1. `.metadata.namespace` set on the resource.
               2. This `namespace` parameter.
               3. `namespace` set for the active context in the kubeconfig.
        :param pulumi.Input[_builtins.str] render_yaml_to_directory: BETA FEATURE - If present, render resource manifests to this directory. In this mode, resources will not
               be created on a Kubernetes cluster, but the rendered manifests will be kept in sync with changes
               to the Pulumi program. This feature is in developer preview, and is disabled by default.
               
               Note that some computed Outputs such as status fields will not be populated
               since the resources are not created on a Kubernetes cluster. These Output values will remain undefined,
               and may result in an error if they are referenced by other resources. Also note that any secret values
               used in these resources will be rendered in plaintext to the resulting YAML.
        :param pulumi.Input[_builtins.bool] skip_update_unreachable: If present and set to true, the provider will skip resources update associated with an unreachable Kubernetes cluster from Pulumi state
        :param pulumi.Input[_builtins.bool] suppress_deprecation_warnings: If present and set to true, suppress apiVersion deprecation warnings from the CLI.
        :param pulumi.Input[_builtins.bool] suppress_helm_hook_warnings: If present and set to true, suppress unsupported Helm hook warnings from the CLI.
        """
        if cluster is not None:
            pulumi.set(__self__, "cluster", cluster)
        if cluster_identifier is not None:
            pulumi.set(__self__, "cluster_identifier", cluster_identifier)
        if context is not None:
            pulumi.set(__self__, "context", context)
        if delete_unreachable is None:
            delete_unreachable = _utilities.get_env_bool('PULUMI_K8S_DELETE_UNREACHABLE')
        if delete_unreachable is not None:
            pulumi.set(__self__, "delete_unreachable", delete_unreachable)
        if enable_config_map_mutable is None:
            enable_config_map_mutable = _utilities.get_env_bool('PULUMI_K8S_ENABLE_CONFIGMAP_MUTABLE')
        if enable_config_map_mutable is not None:
            pulumi.set(__self__, "enable_config_map_mutable", enable_config_map_mutable)
        if enable_secret_mutable is None:
            enable_secret_mutable = _utilities.get_env_bool('PULUMI_K8S_ENABLE_SECRET_MUTABLE')
        if enable_secret_mutable is not None:
            pulumi.set(__self__, "enable_secret_mutable", enable_secret_mutable)
        if enable_server_side_apply is None:
            enable_server_side_apply = _utilities.get_env_bool('PULUMI_K8S_ENABLE_SERVER_SIDE_APPLY')
        if enable_server_side_apply is not None:
            pulumi.set(__self__, "enable_server_side_apply", enable_server_side_apply)
        if helm_release_settings is not None:
            pulumi.set(__self__, "helm_release_settings", helm_release_settings)
        if kube_client_settings is not None:
            pulumi.set(__self__, "kube_client_settings", kube_client_settings)
        if kubeconfig is None:
            kubeconfig = _utilities.get_env('KUBECONFIG')
        if kubeconfig is not None:
            pulumi.set(__self__, "kubeconfig", kubeconfig)
        if namespace is not None:
            pulumi.set(__self__, "namespace", namespace)
        if render_yaml_to_directory is not None:
            pulumi.set(__self__, "render_yaml_to_directory", render_yaml_to_directory)
        if skip_update_unreachable is None:
            skip_update_unreachable = _utilities.get_env_bool('PULUMI_K8S_SKIP_UPDATE_UNREACHABLE')
        if skip_update_unreachable is not None:
            pulumi.set(__self__, "skip_update_unreachable", skip_update_unreachable)
        if suppress_deprecation_warnings is None:
            suppress_deprecation_warnings = _utilities.get_env_bool('PULUMI_K8S_SUPPRESS_DEPRECATION_WARNINGS')
        if suppress_deprecation_warnings is not None:
            pulumi.set(__self__, "suppress_deprecation_warnings", suppress_deprecation_warnings)
        if suppress_helm_hook_warnings is None:
            suppress_helm_hook_warnings = _utilities.get_env_bool('PULUMI_K8S_SUPPRESS_HELM_HOOK_WARNINGS')
        if suppress_helm_hook_warnings is not None:
            pulumi.set(__self__, "suppress_helm_hook_warnings", suppress_helm_hook_warnings)

    @_builtins.property
    @pulumi.getter
    def cluster(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        If present, the name of the kubeconfig cluster to use.
        """
        return pulumi.get(self, "cluster")

    @cluster.setter
    def cluster(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "cluster", value)

    @_builtins.property
    @pulumi.getter(name="clusterIdentifier")
    def cluster_identifier(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        If present, this value will control the provider's replacement behavior. In particular, the provider will _only_ be replaced when `clusterIdentifier` changes; all other changes to provider configuration will be treated as updates.

        Kubernetes does not yet offer an API for cluster identification, so Pulumi uses heuristics to decide when a provider resource should be replaced or updated. These heuristics can sometimes lead to destructive replace operations when an update would be more appropriate, or vice versa.

        Use `clusterIdentifier` for more fine-grained control of the provider resource's lifecycle.
        """
        return pulumi.get(self, "cluster_identifier")

    @cluster_identifier.setter
    def cluster_identifier(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "cluster_identifier", value)

    @_builtins.property
    @pulumi.getter
    def context(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        If present, the name of the kubeconfig context to use.
        """
        return pulumi.get(self, "context")

    @context.setter
    def context(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "context", value)

    @_builtins.property
    @pulumi.getter(name="deleteUnreachable")
    def delete_unreachable(self) -> Optional[pulumi.Input[_builtins.bool]]:
        """
        If present and set to true, the provider will delete resources associated with an unreachable Kubernetes cluster from Pulumi state
        """
        return pulumi.get(self, "delete_unreachable")

    @delete_unreachable.setter
    def delete_unreachable(self, value: Optional[pulumi.Input[_builtins.bool]]):
        pulumi.set(self, "delete_unreachable", value)

    @_builtins.property
    @pulumi.getter(name="enableConfigMapMutable")
    def enable_config_map_mutable(self) -> Optional[pulumi.Input[_builtins.bool]]:
        """
        BETA FEATURE - If present and set to true, allow ConfigMaps to be mutated.
        This feature is in developer preview, and is disabled by default.

        This config can be specified in the following ways using this precedence:
        1. This `enableConfigMapMutable` parameter.
        2. The `PULUMI_K8S_ENABLE_CONFIGMAP_MUTABLE` environment variable.
        """
        return pulumi.get(self, "enable_config_map_mutable")

    @enable_config_map_mutable.setter
    def enable_config_map_mutable(self, value: Optional[pulumi.Input[_builtins.bool]]):
        pulumi.set(self, "enable_config_map_mutable", value)

    @_builtins.property
    @pulumi.getter(name="enableSecretMutable")
    def enable_secret_mutable(self) -> Optional[pulumi.Input[_builtins.bool]]:
        """
        BETA FEATURE - If present and set to true, allow Secrets to be mutated.
        This feature is in developer preview, and is disabled by default.

        This config can be specified in the following ways using this precedence:
        1. This `enableSecretMutable` parameter.
        2. The `PULUMI_K8S_ENABLE_SECRET_MUTABLE` environment variable.
        """
        return pulumi.get(self, "enable_secret_mutable")

    @enable_secret_mutable.setter
    def enable_secret_mutable(self, value: Optional[pulumi.Input[_builtins.bool]]):
        pulumi.set(self, "enable_secret_mutable", value)

    @_builtins.property
    @pulumi.getter(name="enableServerSideApply")
    def enable_server_side_apply(self) -> Optional[pulumi.Input[_builtins.bool]]:
        """
        If present and set to false, disable Server-Side Apply mode.
        See https://github.com/pulumi/pulumi-kubernetes/issues/2011 for additional details.
        """
        return pulumi.get(self, "enable_server_side_apply")

    @enable_server_side_apply.setter
    def enable_server_side_apply(self, value: Optional[pulumi.Input[_builtins.bool]]):
        pulumi.set(self, "enable_server_side_apply", value)

    @_builtins.property
    @pulumi.getter(name="helmReleaseSettings")
    def helm_release_settings(self) -> Optional[pulumi.Input['HelmReleaseSettingsArgs']]:
        """
        Options to configure the Helm Release resource.
        """
        return pulumi.get(self, "helm_release_settings")

    @helm_release_settings.setter
    def helm_release_settings(self, value: Optional[pulumi.Input['HelmReleaseSettingsArgs']]):
        pulumi.set(self, "helm_release_settings", value)

    @_builtins.property
    @pulumi.getter(name="kubeClientSettings")
    def kube_client_settings(self) -> Optional[pulumi.Input['KubeClientSettingsArgs']]:
        """
        Options for tuning the Kubernetes client used by a Provider.
        """
        return pulumi.get(self, "kube_client_settings")

    @kube_client_settings.setter
    def kube_client_settings(self, value: Optional[pulumi.Input['KubeClientSettingsArgs']]):
        pulumi.set(self, "kube_client_settings", value)

    @_builtins.property
    @pulumi.getter
    def kubeconfig(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        The contents of a kubeconfig file or the path to a kubeconfig file.
        """
        return pulumi.get(self, "kubeconfig")

    @kubeconfig.setter
    def kubeconfig(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "kubeconfig", value)

    @_builtins.property
    @pulumi.getter
    def namespace(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        If present, the default namespace to use. This flag is ignored for cluster-scoped resources.

        A namespace can be specified in multiple places, and the precedence is as follows:
        1. `.metadata.namespace` set on the resource.
        2. This `namespace` parameter.
        3. `namespace` set for the active context in the kubeconfig.
        """
        return pulumi.get(self, "namespace")

    @namespace.setter
    def namespace(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "namespace", value)

    @_builtins.property
    @pulumi.getter(name="renderYamlToDirectory")
    def render_yaml_to_directory(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        BETA FEATURE - If present, render resource manifests to this directory. In this mode, resources will not
        be created on a Kubernetes cluster, but the rendered manifests will be kept in sync with changes
        to the Pulumi program. This feature is in developer preview, and is disabled by default.

        Note that some computed Outputs such as status fields will not be populated
        since the resources are not created on a Kubernetes cluster. These Output values will remain undefined,
        and may result in an error if they are referenced by other resources. Also note that any secret values
        used in these resources will be rendered in plaintext to the resulting YAML.
        """
        return pulumi.get(self, "render_yaml_to_directory")

    @render_yaml_to_directory.setter
    def render_yaml_to_directory(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "render_yaml_to_directory", value)

    @_builtins.property
    @pulumi.getter(name="skipUpdateUnreachable")
    def skip_update_unreachable(self) -> Optional[pulumi.Input[_builtins.bool]]:
        """
        If present and set to true, the provider will skip resources update associated with an unreachable Kubernetes cluster from Pulumi state
        """
        return pulumi.get(self, "skip_update_unreachable")

    @skip_update_unreachable.setter
    def skip_update_unreachable(self, value: Optional[pulumi.Input[_builtins.bool]]):
        pulumi.set(self, "skip_update_unreachable", value)

    @_builtins.property
    @pulumi.getter(name="suppressDeprecationWarnings")
    def suppress_deprecation_warnings(self) -> Optional[pulumi.Input[_builtins.bool]]:
        """
        If present and set to true, suppress apiVersion deprecation warnings from the CLI.
        """
        return pulumi.get(self, "suppress_deprecation_warnings")

    @suppress_deprecation_warnings.setter
    def suppress_deprecation_warnings(self, value: Optional[pulumi.Input[_builtins.bool]]):
        pulumi.set(self, "suppress_deprecation_warnings", value)

    @_builtins.property
    @pulumi.getter(name="suppressHelmHookWarnings")
    def suppress_helm_hook_warnings(self) -> Optional[pulumi.Input[_builtins.bool]]:
        """
        If present and set to true, suppress unsupported Helm hook warnings from the CLI.
        """
        return pulumi.get(self, "suppress_helm_hook_warnings")

    @suppress_helm_hook_warnings.setter
    def suppress_helm_hook_warnings(self, value: Optional[pulumi.Input[_builtins.bool]]):
        pulumi.set(self, "suppress_helm_hook_warnings", value)


@pulumi.type_token("pulumi:providers:kubernetes")
class Provider(pulumi.ProviderResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 cluster: Optional[pulumi.Input[_builtins.str]] = None,
                 cluster_identifier: Optional[pulumi.Input[_builtins.str]] = None,
                 context: Optional[pulumi.Input[_builtins.str]] = None,
                 delete_unreachable: Optional[pulumi.Input[_builtins.bool]] = None,
                 enable_config_map_mutable: Optional[pulumi.Input[_builtins.bool]] = None,
                 enable_secret_mutable: Optional[pulumi.Input[_builtins.bool]] = None,
                 enable_server_side_apply: Optional[pulumi.Input[_builtins.bool]] = None,
                 helm_release_settings: Optional[pulumi.Input[Union['HelmReleaseSettingsArgs', 'HelmReleaseSettingsArgsDict']]] = None,
                 kube_client_settings: Optional[pulumi.Input[Union['KubeClientSettingsArgs', 'KubeClientSettingsArgsDict']]] = None,
                 kubeconfig: Optional[pulumi.Input[_builtins.str]] = None,
                 namespace: Optional[pulumi.Input[_builtins.str]] = None,
                 render_yaml_to_directory: Optional[pulumi.Input[_builtins.str]] = None,
                 skip_update_unreachable: Optional[pulumi.Input[_builtins.bool]] = None,
                 suppress_deprecation_warnings: Optional[pulumi.Input[_builtins.bool]] = None,
                 suppress_helm_hook_warnings: Optional[pulumi.Input[_builtins.bool]] = None,
                 __props__=None):
        """
        The provider type for the kubernetes package.

        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        :param pulumi.Input[_builtins.str] cluster: If present, the name of the kubeconfig cluster to use.
        :param pulumi.Input[_builtins.str] cluster_identifier: If present, this value will control the provider's replacement behavior. In particular, the provider will _only_ be replaced when `clusterIdentifier` changes; all other changes to provider configuration will be treated as updates.
               
               Kubernetes does not yet offer an API for cluster identification, so Pulumi uses heuristics to decide when a provider resource should be replaced or updated. These heuristics can sometimes lead to destructive replace operations when an update would be more appropriate, or vice versa.
               
               Use `clusterIdentifier` for more fine-grained control of the provider resource's lifecycle.
        :param pulumi.Input[_builtins.str] context: If present, the name of the kubeconfig context to use.
        :param pulumi.Input[_builtins.bool] delete_unreachable: If present and set to true, the provider will delete resources associated with an unreachable Kubernetes cluster from Pulumi state
        :param pulumi.Input[_builtins.bool] enable_config_map_mutable: BETA FEATURE - If present and set to true, allow ConfigMaps to be mutated.
               This feature is in developer preview, and is disabled by default.
               
               This config can be specified in the following ways using this precedence:
               1. This `enableConfigMapMutable` parameter.
               2. The `PULUMI_K8S_ENABLE_CONFIGMAP_MUTABLE` environment variable.
        :param pulumi.Input[_builtins.bool] enable_secret_mutable: BETA FEATURE - If present and set to true, allow Secrets to be mutated.
               This feature is in developer preview, and is disabled by default.
               
               This config can be specified in the following ways using this precedence:
               1. This `enableSecretMutable` parameter.
               2. The `PULUMI_K8S_ENABLE_SECRET_MUTABLE` environment variable.
        :param pulumi.Input[_builtins.bool] enable_server_side_apply: If present and set to false, disable Server-Side Apply mode.
               See https://github.com/pulumi/pulumi-kubernetes/issues/2011 for additional details.
        :param pulumi.Input[Union['HelmReleaseSettingsArgs', 'HelmReleaseSettingsArgsDict']] helm_release_settings: Options to configure the Helm Release resource.
        :param pulumi.Input[Union['KubeClientSettingsArgs', 'KubeClientSettingsArgsDict']] kube_client_settings: Options for tuning the Kubernetes client used by a Provider.
        :param pulumi.Input[_builtins.str] kubeconfig: The contents of a kubeconfig file or the path to a kubeconfig file.
        :param pulumi.Input[_builtins.str] namespace: If present, the default namespace to use. This flag is ignored for cluster-scoped resources.
               
               A namespace can be specified in multiple places, and the precedence is as follows:
               1. `.metadata.namespace` set on the resource.
               2. This `namespace` parameter.
               3. `namespace` set for the active context in the kubeconfig.
        :param pulumi.Input[_builtins.str] render_yaml_to_directory: BETA FEATURE - If present, render resource manifests to this directory. In this mode, resources will not
               be created on a Kubernetes cluster, but the rendered manifests will be kept in sync with changes
               to the Pulumi program. This feature is in developer preview, and is disabled by default.
               
               Note that some computed Outputs such as status fields will not be populated
               since the resources are not created on a Kubernetes cluster. These Output values will remain undefined,
               and may result in an error if they are referenced by other resources. Also note that any secret values
               used in these resources will be rendered in plaintext to the resulting YAML.
        :param pulumi.Input[_builtins.bool] skip_update_unreachable: If present and set to true, the provider will skip resources update associated with an unreachable Kubernetes cluster from Pulumi state
        :param pulumi.Input[_builtins.bool] suppress_deprecation_warnings: If present and set to true, suppress apiVersion deprecation warnings from the CLI.
        :param pulumi.Input[_builtins.bool] suppress_helm_hook_warnings: If present and set to true, suppress unsupported Helm hook warnings from the CLI.
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: Optional[ProviderArgs] = None,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        The provider type for the kubernetes package.

        :param str resource_name: The name of the resource.
        :param ProviderArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(ProviderArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 cluster: Optional[pulumi.Input[_builtins.str]] = None,
                 cluster_identifier: Optional[pulumi.Input[_builtins.str]] = None,
                 context: Optional[pulumi.Input[_builtins.str]] = None,
                 delete_unreachable: Optional[pulumi.Input[_builtins.bool]] = None,
                 enable_config_map_mutable: Optional[pulumi.Input[_builtins.bool]] = None,
                 enable_secret_mutable: Optional[pulumi.Input[_builtins.bool]] = None,
                 enable_server_side_apply: Optional[pulumi.Input[_builtins.bool]] = None,
                 helm_release_settings: Optional[pulumi.Input[Union['HelmReleaseSettingsArgs', 'HelmReleaseSettingsArgsDict']]] = None,
                 kube_client_settings: Optional[pulumi.Input[Union['KubeClientSettingsArgs', 'KubeClientSettingsArgsDict']]] = None,
                 kubeconfig: Optional[pulumi.Input[_builtins.str]] = None,
                 namespace: Optional[pulumi.Input[_builtins.str]] = None,
                 render_yaml_to_directory: Optional[pulumi.Input[_builtins.str]] = None,
                 skip_update_unreachable: Optional[pulumi.Input[_builtins.bool]] = None,
                 suppress_deprecation_warnings: Optional[pulumi.Input[_builtins.bool]] = None,
                 suppress_helm_hook_warnings: Optional[pulumi.Input[_builtins.bool]] = None,
                 __props__=None):
        opts = pulumi.ResourceOptions.merge(_utilities.get_resource_opts_defaults(), opts)
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.id is None:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = ProviderArgs.__new__(ProviderArgs)

            __props__.__dict__["cluster"] = cluster
            __props__.__dict__["cluster_identifier"] = cluster_identifier
            __props__.__dict__["context"] = context
            if delete_unreachable is None:
                delete_unreachable = _utilities.get_env_bool('PULUMI_K8S_DELETE_UNREACHABLE')
            __props__.__dict__["delete_unreachable"] = pulumi.Output.from_input(delete_unreachable).apply(pulumi.runtime.to_json) if delete_unreachable is not None else None
            if enable_config_map_mutable is None:
                enable_config_map_mutable = _utilities.get_env_bool('PULUMI_K8S_ENABLE_CONFIGMAP_MUTABLE')
            __props__.__dict__["enable_config_map_mutable"] = pulumi.Output.from_input(enable_config_map_mutable).apply(pulumi.runtime.to_json) if enable_config_map_mutable is not None else None
            if enable_secret_mutable is None:
                enable_secret_mutable = _utilities.get_env_bool('PULUMI_K8S_ENABLE_SECRET_MUTABLE')
            __props__.__dict__["enable_secret_mutable"] = pulumi.Output.from_input(enable_secret_mutable).apply(pulumi.runtime.to_json) if enable_secret_mutable is not None else None
            if enable_server_side_apply is None:
                enable_server_side_apply = _utilities.get_env_bool('PULUMI_K8S_ENABLE_SERVER_SIDE_APPLY')
            __props__.__dict__["enable_server_side_apply"] = pulumi.Output.from_input(enable_server_side_apply).apply(pulumi.runtime.to_json) if enable_server_side_apply is not None else None
            __props__.__dict__["helm_release_settings"] = pulumi.Output.from_input(helm_release_settings).apply(pulumi.runtime.to_json) if helm_release_settings is not None else None
            __props__.__dict__["kube_client_settings"] = pulumi.Output.from_input(kube_client_settings).apply(pulumi.runtime.to_json) if kube_client_settings is not None else None
            if kubeconfig is None:
                kubeconfig = _utilities.get_env('KUBECONFIG')
            __props__.__dict__["kubeconfig"] = kubeconfig
            __props__.__dict__["namespace"] = namespace
            __props__.__dict__["render_yaml_to_directory"] = render_yaml_to_directory
            if skip_update_unreachable is None:
                skip_update_unreachable = _utilities.get_env_bool('PULUMI_K8S_SKIP_UPDATE_UNREACHABLE')
            __props__.__dict__["skip_update_unreachable"] = pulumi.Output.from_input(skip_update_unreachable).apply(pulumi.runtime.to_json) if skip_update_unreachable is not None else None
            if suppress_deprecation_warnings is None:
                suppress_deprecation_warnings = _utilities.get_env_bool('PULUMI_K8S_SUPPRESS_DEPRECATION_WARNINGS')
            __props__.__dict__["suppress_deprecation_warnings"] = pulumi.Output.from_input(suppress_deprecation_warnings).apply(pulumi.runtime.to_json) if suppress_deprecation_warnings is not None else None
            if suppress_helm_hook_warnings is None:
                suppress_helm_hook_warnings = _utilities.get_env_bool('PULUMI_K8S_SUPPRESS_HELM_HOOK_WARNINGS')
            __props__.__dict__["suppress_helm_hook_warnings"] = pulumi.Output.from_input(suppress_helm_hook_warnings).apply(pulumi.runtime.to_json) if suppress_helm_hook_warnings is not None else None
        super(Provider, __self__).__init__(
            'kubernetes',
            resource_name,
            __props__,
            opts)

