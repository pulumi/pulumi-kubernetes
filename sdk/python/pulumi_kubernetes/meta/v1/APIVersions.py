import pulumi
import pulumi.runtime

from ... import tables

class APIVersions(pulumi.CustomResource):
    """
    APIVersions lists the versions that are available, to allow clients to discover the API at /api,
    which is the root path of the legacy v1 API.
    """
    def __init__(self, __name__, __opts__=None, server_address_by_client_cid_rs=None, versions=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        __props__['kind'] = 'APIVersions'
        if not serverAddressByClientCIDRs:
            raise TypeError('Missing required property serverAddressByClientCIDRs')
        __props__['serverAddressByClientCIDRs'] = server_address_by_client_cid_rs
        if not versions:
            raise TypeError('Missing required property versions')
        __props__['versions'] = versions

        super(APIVersions, self).__init__(
            "kubernetes:core/v1:APIVersions",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return tables._CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return tables._CASING_BACKWARD_TABLE.get(prop) or prop
