import pulumi
import pulumi.runtime

class CertificateSigningRequest(pulumi.CustomResource):
    """
    Describes a certificate signing request
    """
    def __init__(self, __name__, __opts__=None, metadata=None, spec=None, status=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'certificates.k8s.io/v1beta1'
        self.apiVersion = 'certificates.k8s.io/v1beta1'

        __props__['kind'] = 'CertificateSigningRequest'
        self.kind = 'CertificateSigningRequest'

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        
        __props__['metadata'] = metadata

        if spec and not isinstance(spec, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.spec = spec
        """
        The certificate request itself and any additional information.
        """
        __props__['spec'] = spec

        if status and not isinstance(status, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.status = status
        """
        Derived information about the request.
        """
        __props__['status'] = status

        super(CertificateSigningRequest, self).__init__(
            "kubernetes:certificates.k8s.io/v1beta1:CertificateSigningRequest",
            __name__,
            __props__,
            __opts__)
