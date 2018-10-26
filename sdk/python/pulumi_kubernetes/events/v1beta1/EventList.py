import pulumi
import pulumi.runtime

class EventList(pulumi.CustomResource):
    """
    EventList is a list of Event objects.
    """
    def __init__(self, __name__, __opts__=None, items=None, metadata=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'events.k8s.io/v1beta1'
        self.apiVersion = 'events.k8s.io/v1beta1'

        __props__['kind'] = 'EventList'
        self.kind = 'EventList'

        if not items:
            raise TypeError('Missing required property items')
        elif not isinstance(items, list):
            raise TypeError('Expected property aliases to be a list')
        self.items = items
        """
        Items is a list of schema objects.
        """
        __props__['items'] = items

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard list metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        super(EventList, self).__init__(
            "kubernetes:events.k8s.io/v1beta1:EventList",
            __name__,
            __props__,
            __opts__)
