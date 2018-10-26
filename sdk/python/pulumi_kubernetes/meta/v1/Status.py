import pulumi
import pulumi.runtime

class Status(pulumi.CustomResource):
    """
    Status is a return value for calls that don't return other objects.
    """
    def __init__(self, __name__, __opts__=None, code=None, details=None, message=None, metadata=None, reason=None, status=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        self.apiVersion = 'v1'

        __props__['kind'] = 'Status'
        self.kind = 'Status'

        if code and not isinstance(code, int):
            raise TypeError('Expected property aliases to be a int')
        self.code = code
        """
        Suggested HTTP return code for this status, 0 if not set.
        """
        __props__['code'] = code

        if details and not isinstance(details, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.details = details
        """
        Extended data associated with the reason.  Each reason may define its own extended details.
        This field is optional and the data returned is not guaranteed to conform to any schema
        except that defined by the reason type.
        """
        __props__['details'] = details

        if message and not isinstance(message, str):
            raise TypeError('Expected property aliases to be a str')
        self.message = message
        """
        A human-readable description of the status of this operation.
        """
        __props__['message'] = message

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard list metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
        """
        __props__['metadata'] = metadata

        if reason and not isinstance(reason, str):
            raise TypeError('Expected property aliases to be a str')
        self.reason = reason
        """
        A machine-readable description of why this operation is in the "Failure" status. If this
        value is empty there is no information available. A Reason clarifies an HTTP status code but
        does not override it.
        """
        __props__['reason'] = reason

        if status and not isinstance(status, str):
            raise TypeError('Expected property aliases to be a str')
        self.status = status
        """
        Status of the operation. One of: "Success" or "Failure". More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
        """
        __props__['status'] = status

        super(Status, self).__init__(
            "kubernetes:core/v1:Status",
            __name__,
            __props__,
            __opts__)
