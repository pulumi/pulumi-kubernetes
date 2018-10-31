import pulumi
import pulumi.runtime

from ...tables import _CASING_FORWARD_TABLE, _CASING_BACKWARD_TABLE

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
        self.apiVersion = 'v1'

        __props__['kind'] = 'APIVersions'
        self.kind = 'APIVersions'

        if not serverAddressByClientCIDRs:
            raise TypeError('Missing required property serverAddressByClientCIDRs')
        elif not isinstance(serverAddressByClientCIDRs, list):
            raise TypeError('Expected property aliases to be a list')
        self.server_address_by_client_cid_rs = server_address_by_client_cid_rs
        """
        a map of client CIDR to server address that is serving this group. This is to help clients
        reach servers in the most network-efficient way possible. Clients can use the appropriate
        server address as per the CIDR that they match. In case of multiple matches, clients should
        use the longest matching CIDR. The server returns only those CIDRs that it thinks that the
        client can match. For example: the master will return an internal IP CIDR only, if the
        client reaches the server using an internal IP. Server looks at X-Forwarded-For header or
        X-Real-Ip header or request.RemoteAddr (in that order) to get the client IP.
        """
        __props__['serverAddressByClientCIDRs'] = server_address_by_client_cid_rs

        if not versions:
            raise TypeError('Missing required property versions')
        elif not isinstance(versions, list):
            raise TypeError('Expected property aliases to be a list')
        self.versions = versions
        """
        versions are the api versions that are available.
        """
        __props__['versions'] = versions

        super(APIVersions, self).__init__(
            "kubernetes:core/v1:APIVersions",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return _CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return _CASING_BACKWARD_TABLE.get(prop) or prop
