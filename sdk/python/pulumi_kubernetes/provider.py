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
                 suppress_deprecation_warnings=None):
        """
        Create a Provider resource with the given unique name, arguments, and options.

        :param str __name__: The unique name of the resource.
        :param pulumi.ResourceOptions __opts__: An optional bag of options that controls this resource's behavior.
        :param pulumi.Input[str] cluster: If present, the name of the kubeconfig cluster to use.
        :param pulumi.Input[str] context: If present, the name of the kubeconfig context to use.
        :param pulumi.Input[bool] enable_dry_run: BETA FEATURE - If present and set to True, enable server-side diff
                                                  calculations. This feature is in developer preview, and is disabled by
                                                  default.
        :param pulumi.Input[str] kubeconfig: The contents of a kubeconfig file.
                                             If this is set, this config will be used instead of $KUBECONFIG.
        :param pulumi.Input[str] namespace: If present, the default namespace to use.
                                            This flag is ignored for cluster-scoped resources.
                                            Note: if .metadata.namespace is set on a resource, that value takes
                                            precedence over the provider default.
        :param pulumi.Input[bool] suppress_deprecation_warnings: If present and set to True, suppress apiVersion
                                                                 deprecation warnings from the CLI.
        """
        __props__ = {
            "cluster": cluster,
            "context": context,
            "enableDryRun": "true" if enable_dry_run else "false",
            "kubeconfig": kubeconfig,
            "namespace": namespace,
            "suppress_deprecation_warnings": "true" if suppress_deprecation_warnings else "false",
        }
        super(Provider, self).__init__("kubernetes", __name__, __props__, __opts__)
