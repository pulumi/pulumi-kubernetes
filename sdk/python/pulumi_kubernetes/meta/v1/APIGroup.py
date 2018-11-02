import pulumi
import pulumi.runtime

from ... import tables

class APIGroup(pulumi.CustomResource):
    """
    APIGroup contains the name, the supported versions, and the preferred version of a group.
    """
    def __init__(self, __name__, __opts__=None, name=None, preferred_version=None, server_address_by_client_cid_rs=None, versions=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        __props__['kind'] = 'APIGroup'
        if not name:
            raise TypeError('Missing required property name')
        __props__['name'] = name
        if not serverAddressByClientCIDRs:
            raise TypeError('Missing required property serverAddressByClientCIDRs')
        __props__['serverAddressByClientCIDRs'] = server_address_by_client_cid_rs
        if not versions:
            raise TypeError('Missing required property versions')
        __props__['versions'] = versions
        __props__['preferredVersion'] = preferred_version

        super(APIGroup, self).__init__(
            "kubernetes:core/v1:APIGroup",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return tables._CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return tables._CASING_BACKWARD_TABLE.get(prop) or prop
