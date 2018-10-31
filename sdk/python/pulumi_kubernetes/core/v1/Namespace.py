import pulumi
import pulumi.runtime

from ...tables import _CASING_FORWARD_TABLE, _CASING_BACKWARD_TABLE

class Namespace(pulumi.CustomResource):
    """
    Namespace provides a scope for Names. Use of multiple namespaces is optional.
    """
    def __init__(self, __name__, __opts__=None, metadata=None, spec=None, status=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        self.apiVersion = 'v1'

        __props__['kind'] = 'Namespace'
        self.kind = 'Namespace'

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        if spec and not isinstance(spec, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.spec = spec
        """
        Spec defines the behavior of the Namespace. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
        """
        __props__['spec'] = spec

        if status and not isinstance(status, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.status = status
        """
        Status describes the current status of a Namespace. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
        """
        __props__['status'] = status

        super(Namespace, self).__init__(
            "kubernetes:core/v1:Namespace",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return _CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return _CASING_BACKWARD_TABLE.get(prop) or prop
