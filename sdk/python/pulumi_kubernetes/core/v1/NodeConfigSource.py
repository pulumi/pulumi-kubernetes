import pulumi
import pulumi.runtime

class NodeConfigSource(pulumi.CustomResource):
    """
    NodeConfigSource specifies a source of node configuration. Exactly one subfield (excluding
    metadata) must be non-nil.
    """
    def __init__(self, __name__, __opts__=None, configMapRef=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        self.apiVersion = 'v1'

        __props__['kind'] = 'NodeConfigSource'
        self.kind = 'NodeConfigSource'

        if configMapRef and not isinstance(configMapRef, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.configMapRef = configMapRef
        
        __props__['configMapRef'] = configMapRef

        super(NodeConfigSource, self).__init__(
            "kubernetes:core/v1:NodeConfigSource",
            __name__,
            __props__,
            __opts__)
