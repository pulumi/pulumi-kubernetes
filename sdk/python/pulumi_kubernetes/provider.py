import pulumi
from . import tables


class Provider(pulumi.ProviderResource):
    """
    The provider type for the kubernetes package.
    """
    def __init__(self,
                 __name__,
                 __opts__=None,
                 cluster=None,
                 context=None,
                 kubeconfig=None,
                 namespace=None):
        """
        Create a Provider resource with the given unique name, arguments, and options.

        :param str __name__: The unique name of the resource.
        :param pulumi.ResourceOptions __opts__: An optional bag of options that controls this resource's behavior.
        :param pulumi.Input[str] cluster: If present, the name of the kubeconfig cluster to use.
        :param pulumi.Input[str] context: If present, the name of the kubeconfig context to use.
        :param pulumi.Input[str] kubeconfig: The contents of a kubeconfig file. If this is set, this config will be used instead
                               of $KUBECONFIG.
        :param pulumi.Input[str] namespace: If present, the namespace scope to use.
        """
        __props__ = {
            "cluster": cluster,
            "context": context,
            "kubeconfig": kubeconfig,
            "namespace": namespace,
        }
        super(Provider, self).__init__("kubernetes", __name__, __props__, __opts__)
