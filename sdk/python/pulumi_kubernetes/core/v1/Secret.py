import pulumi
import pulumi.runtime

class Secret(pulumi.CustomResource):
    """
    Secret holds secret data of a certain type. The total bytes of the values in the Data field must
    be less than MaxSecretSize bytes.
    """
    def __init__(self, __name__, __opts__=None, data=None, metadata=None, string_data=None, type=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        self.apiVersion = 'v1'

        __props__['kind'] = 'Secret'
        self.kind = 'Secret'

        if data and not isinstance(data, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.data = data
        """
        Data contains the secret data. Each key must consist of alphanumeric characters, '-', '_' or
        '.'. The serialized form of the secret data is a base64 encoded string, representing the
        arbitrary (possibly non-string) data value here. Described in
        https://tools.ietf.org/html/rfc4648#section-4
        """
        __props__['data'] = data

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        if string_data and not isinstance(string_data, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.string_data = string_data
        """
        stringData allows specifying non-binary secret data in string form. It is provided as a
        write-only convenience method. All keys and values are merged into the data field on write,
        overwriting any existing values. It is never output when reading from the API.
        """
        __props__['stringData'] = string_data

        if type and not isinstance(type, str):
            raise TypeError('Expected property aliases to be a str')
        self.type = type
        """
        Used to facilitate programmatic handling of secret data.
        """
        __props__['type'] = type

        super(Secret, self).__init__(
            "kubernetes:core/v1:Secret",
            __name__,
            __props__,
            __opts__)
