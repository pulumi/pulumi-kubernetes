import pulumi
import pulumi.runtime

class TokenReview(pulumi.CustomResource):
    """
    TokenReview attempts to authenticate a token to a known user. Note: TokenReview requests may be
    cached by the webhook token authenticator plugin in the kube-apiserver.
    """
    def __init__(self, __name__, __opts__=None, metadata=None, spec=None, status=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'authentication.k8s.io/v1beta1'
        self.apiVersion = 'authentication.k8s.io/v1beta1'

        __props__['kind'] = 'TokenReview'
        self.kind = 'TokenReview'

        if not spec:
            raise TypeError('Missing required property spec')
        elif not isinstance(spec, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.spec = spec
        """
        Spec holds information about the request being evaluated
        """
        __props__['spec'] = spec

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        
        __props__['metadata'] = metadata

        if status and not isinstance(status, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.status = status
        """
        Status is filled in by the server and indicates whether the request can be authenticated.
        """
        __props__['status'] = status

        super(TokenReview, self).__init__(
            "kubernetes:authentication.k8s.io/v1beta1:TokenReview",
            __name__,
            __props__,
            __opts__)
