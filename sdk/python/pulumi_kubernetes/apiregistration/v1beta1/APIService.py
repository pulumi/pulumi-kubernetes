import pulumi
import pulumi.runtime

class APIService(pulumi.CustomResource):
    """
    APIService represents a server for a particular GroupVersion. Name must be "version.group".
    """
    def __init__(self, __name__, __opts__=None, metadata=None, spec=None, status=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'apiregistration/v1beta1'
        self.apiVersion = 'apiregistration/v1beta1'

        __props__['kind'] = 'APIService'
        self.kind = 'APIService'

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        
        __props__['metadata'] = metadata

        if spec and not isinstance(spec, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.spec = spec
        """
        Spec contains information for locating and communicating with a server
        """
        __props__['spec'] = spec

        if status and not isinstance(status, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.status = status
        """
        Status contains derived information about an API server
        """
        __props__['status'] = status

        super(APIService, self).__init__(
            "kubernetes:apiregistration/v1beta1:APIService",
            __name__,
            __props__,
            __opts__)
