import pulumi
import pulumi.runtime

from ... import tables

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
        __props__['kind'] = 'Eviction'
        __props__['deleteOptions'] = delete_options
        __props__['metadata'] = metadata

        super(Eviction, self).__init__(
            "kubernetes:policy/v1beta1:Eviction",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return tables._CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return tables._CASING_BACKWARD_TABLE.get(prop) or prop
