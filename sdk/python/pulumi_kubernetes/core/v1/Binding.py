import pulumi
import pulumi.runtime

class Binding(pulumi.CustomResource):
    """
    Binding ties one object to another; for example, a pod is bound to a node by a scheduler.
    Deprecated in 1.7, please use the bindings subresource of pods instead.
    """
    def __init__(self, __name__, __opts__=None, metadata=None, target=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        self.apiVersion = 'v1'

        __props__['kind'] = 'Binding'
        self.kind = 'Binding'

        if not target:
            raise TypeError('Missing required property target')
        elif not isinstance(target, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.target = target
        """
        The target object that you want to bind to the standard object.
        """
        __props__['target'] = target

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        super(Binding, self).__init__(
            "kubernetes:core/v1:Binding",
            __name__,
            __props__,
            __opts__)
