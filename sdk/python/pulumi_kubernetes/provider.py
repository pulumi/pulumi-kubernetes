import pulumi


class Provider(pulumi.ProviderResource):
    """
    The provider type for the kubernetes package.
    """

    def __init__(self,
                 __name__,
                 __opts__=None,
                 cluster=None,
                 context=None,
                 enable_dry_run=None,
                 kubeconfig=None,
                 namespace=None,
                 suppress_deprecation_warnings=None,
                 render_yaml_to_directory=None):
        """
        Create a Provider resource with the given unique name, arguments, and options.

        :param str __name__: The unique name of the resource.
        :param pulumi.ResourceOptions __opts__: An optional bag of options that controls this resource's behavior.
        :param pulumi.Input[str] cluster: If present, the name of the kubeconfig cluster to use.
        :param pulumi.Input[str] context: If present, the name of the kubeconfig context to use.
        :param pulumi.Input[bool] enable_dry_run: BETA FEATURE - If present and set to True, enable server-side diff
                                  calculations. This feature is in developer preview, and is disabled by default.
                                  This config can be specified in the following ways, using this precedence:
                                  1. This `enableDryRun` parameter.
                                  2. The `PULUMI_K8S_ENABLE_DRY_RUN` environment variable.
        :param pulumi.Input[str] kubeconfig: The contents of a kubeconfig file.
                                 If this is set, this config will be used instead of $KUBECONFIG.
        :param pulumi.Input[str] namespace: If present, the default namespace to use.
                                 This flag is ignored for cluster-scoped resources.
                                 A namespace can be specified in multiple places, and the precedence is as follows:
                                 1. `.metadata.namespace` set on the resource.
                                 2. This `namespace` parameter.
                                 3. `namespace` set for the active context in the kubeconfig.
        :param pulumi.Input[bool] suppress_deprecation_warnings: If present and set to True, suppress apiVersion
                                  deprecation warnings from the CLI.
                                  This config can be specified in the following ways, using this precedence:
                                  1. This `suppressDeprecationWarnings` parameter.
                                  2. The `PULUMI_K8S_SUPPRESS_DEPRECATION_WARNINGS` environment variable.
        :param pulumi.Input[str] render_yaml_to_directory: BETA FEATURE - If present, render resource manifests to this
                                 directory. In this mode, resources will not be created on a Kubernetes cluster, but
                                 the rendered manifests will be kept in sync with changes to the Pulumi program.
                                 This feature is in developer preview, and is disabled by default. Note that some
                                 computed Outputs such as status fields will not be populated since the resources are
                                 not created on a Kubernetes cluster. These Output values will remain undefined,
                                 and may result in an error if they are referenced by other resources. Also note that
                                 any secret values used in these resources will be rendered in plaintext to the
                                 resulting YAML.
        """
        __props__ = {
            "cluster": cluster,
            "context": context,
            "enableDryRun": "true" if enable_dry_run else "false",
            "kubeconfig": kubeconfig,
            "namespace": namespace,
            "suppress_deprecation_warnings": "true" if suppress_deprecation_warnings else "false",
            "render_yaml_to_directory": render_yaml_to_directory,
        }
        super(Provider, self).__init__("kubernetes", __name__, __props__, __opts__)
