import pulumi
import pulumi.runtime

from ...tables import _CASING_FORWARD_TABLE, _CASING_BACKWARD_TABLE

class Eviction(pulumi.CustomResource):
    """
    Eviction evicts a pod from its node subject to certain policies and safety constraints. This is
    a subresource of Pod.  A request to cause such an eviction is created by POSTing to
    .../pods/<pod name>/evictions.
    """
    def __init__(self, __name__, __opts__=None, delete_options=None, metadata=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'policy/v1beta1'
        self.apiVersion = 'policy/v1beta1'

        __props__['kind'] = 'Eviction'
        self.kind = 'Eviction'

        if delete_options and not isinstance(delete_options, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.delete_options = delete_options
        """
        DeleteOptions may be provided
        """
        __props__['deleteOptions'] = delete_options

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        ObjectMeta describes the pod that is being evicted.
        """
        __props__['metadata'] = metadata

        super(Eviction, self).__init__(
            "kubernetes:policy/v1beta1:Eviction",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return _CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return _CASING_BACKWARD_TABLE.get(prop) or prop
