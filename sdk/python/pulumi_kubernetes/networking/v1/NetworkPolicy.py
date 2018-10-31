import pulumi
import pulumi.runtime

from ...tables import _CASING_FORWARD_TABLE, _CASING_BACKWARD_TABLE

class NetworkPolicy(pulumi.CustomResource):
    """
    NetworkPolicy describes what network traffic is allowed for a set of Pods
    """
    def __init__(self, __name__, __opts__=None, metadata=None, spec=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'networking.k8s.io/v1'
        self.apiVersion = 'networking.k8s.io/v1'

        __props__['kind'] = 'NetworkPolicy'
        self.kind = 'NetworkPolicy'

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
        Specification of the desired behavior for this NetworkPolicy.
        """
        __props__['spec'] = spec

        super(NetworkPolicy, self).__init__(
            "kubernetes:networking.k8s.io/v1:NetworkPolicy",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return _CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return _CASING_BACKWARD_TABLE.get(prop) or prop
