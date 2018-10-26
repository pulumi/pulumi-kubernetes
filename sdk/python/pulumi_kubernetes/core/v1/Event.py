import pulumi
import pulumi.runtime

class Event(pulumi.CustomResource):
    """
    Event is a report of an event somewhere in the cluster.
    """
    def __init__(self, __name__, __opts__=None, action=None, count=None, eventTime=None, firstTimestamp=None, involvedObject=None, lastTimestamp=None, message=None, metadata=None, reason=None, related=None, reportingComponent=None, reportingInstance=None, series=None, source=None, type=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        self.apiVersion = 'v1'

        __props__['kind'] = 'Event'
        self.kind = 'Event'

        if not involvedObject:
            raise TypeError('Missing required property involvedObject')
        elif not isinstance(involvedObject, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.involvedObject = involvedObject
        """
        The object that this event is about.
        """
        __props__['involvedObject'] = involvedObject

        if not metadata:
            raise TypeError('Missing required property metadata')
        elif not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        if action and not isinstance(action, str):
            raise TypeError('Expected property aliases to be a str')
        self.action = action
        """
        What action was taken/failed regarding to the Regarding object.
        """
        __props__['action'] = action

        if count and not isinstance(count, int):
            raise TypeError('Expected property aliases to be a int')
        self.count = count
        """
        The number of times this event has occurred.
        """
        __props__['count'] = count

        if eventTime and not isinstance(eventTime, str):
            raise TypeError('Expected property aliases to be a str')
        self.eventTime = eventTime
        """
        Time when this Event was first observed.
        """
        __props__['eventTime'] = eventTime

        if firstTimestamp and not isinstance(firstTimestamp, str):
            raise TypeError('Expected property aliases to be a str')
        self.firstTimestamp = firstTimestamp
        """
        The time at which the event was first recorded. (Time of server receipt is in TypeMeta.)
        """
        __props__['firstTimestamp'] = firstTimestamp

        if lastTimestamp and not isinstance(lastTimestamp, str):
            raise TypeError('Expected property aliases to be a str')
        self.lastTimestamp = lastTimestamp
        """
        The time at which the most recent occurrence of this event was recorded.
        """
        __props__['lastTimestamp'] = lastTimestamp

        if message and not isinstance(message, str):
            raise TypeError('Expected property aliases to be a str')
        self.message = message
        """
        A human-readable description of the status of this operation.
        """
        __props__['message'] = message

        if reason and not isinstance(reason, str):
            raise TypeError('Expected property aliases to be a str')
        self.reason = reason
        """
        This should be a short, machine understandable string that gives the reason for the
        transition into the object's current status.
        """
        __props__['reason'] = reason

        if related and not isinstance(related, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.related = related
        """
        Optional secondary object for more complex actions.
        """
        __props__['related'] = related

        if reportingComponent and not isinstance(reportingComponent, str):
            raise TypeError('Expected property aliases to be a str')
        self.reportingComponent = reportingComponent
        """
        Name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`.
        """
        __props__['reportingComponent'] = reportingComponent

        if reportingInstance and not isinstance(reportingInstance, str):
            raise TypeError('Expected property aliases to be a str')
        self.reportingInstance = reportingInstance
        """
        ID of the controller instance, e.g. `kubelet-xyzf`.
        """
        __props__['reportingInstance'] = reportingInstance

        if series and not isinstance(series, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.series = series
        """
        Data about the Event series this event represents or nil if it's a singleton Event.
        """
        __props__['series'] = series

        if source and not isinstance(source, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.source = source
        """
        The component reporting this event. Should be a short machine understandable string.
        """
        __props__['source'] = source

        if type and not isinstance(type, str):
            raise TypeError('Expected property aliases to be a str')
        self.type = type
        """
        Type of this event (Normal, Warning), new types could be added in the future
        """
        __props__['type'] = type

        super(Event, self).__init__(
            "kubernetes:core/v1:Event",
            __name__,
            __props__,
            __opts__)
