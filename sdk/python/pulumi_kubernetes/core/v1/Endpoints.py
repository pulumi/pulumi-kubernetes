import pulumi
import pulumi.runtime

class Endpoints(pulumi.CustomResource):
    """
    Endpoints is a collection of endpoints that implement the actual service. Example:
      Name: "mysvc",
      Subsets: [
        {
          Addresses: [{"ip": "10.10.1.1"}, {"ip": "10.10.2.2"}],
          Ports: [{"name": "a", "port": 8675}, {"name": "b", "port": 309}]
        },
        {
          Addresses: [{"ip": "10.10.3.3"}],
          Ports: [{"name": "a", "port": 93}, {"name": "b", "port": 76}]
        },
     ]
    """
    def __init__(self, __name__, __opts__=None, metadata=None, subsets=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        self.apiVersion = 'v1'

        __props__['kind'] = 'Endpoints'
        self.kind = 'Endpoints'

        if not subsets:
            raise TypeError('Missing required property subsets')
        elif not isinstance(subsets, list):
            raise TypeError('Expected property aliases to be a list')
        self.subsets = subsets
        """
        The set of all endpoints is the union of all subsets. Addresses are placed into subsets
        according to the IPs they share. A single address with multiple ports, some of which are
        ready and some of which are not (because they come from different containers) will result in
        the address being displayed in different subsets for the different ports. No address will
        appear in both Addresses and NotReadyAddresses in the same subset. Sets of addresses and
        ports that comprise a service.
        """
        __props__['subsets'] = subsets

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        super(Endpoints, self).__init__(
            "kubernetes:core/v1:Endpoints",
            __name__,
            __props__,
            __opts__)
