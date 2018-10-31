import pulumi
import pulumi.runtime

from ...tables import _CASING_FORWARD_TABLE, _CASING_BACKWARD_TABLE

class Event(pulumi.CustomResource):
    """
    Event is a report of an event somewhere in the cluster.
    """
    def __init__(self, __name__, __opts__=None, action=None, count=None, event_time=None, first_timestamp=None, involved_object=None, last_timestamp=None, message=None, metadata=None, reason=None, related=None, reporting_component=None, reporting_instance=None, series=None, source=None, type=None):
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
        self.involved_object = involved_object
        """
        The object that this event is about.
        """
        __props__['involvedObject'] = involved_object

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

        if event_time and not isinstance(event_time, str):
            raise TypeError('Expected property aliases to be a str')
        self.event_time = event_time
        """
        Time when this Event was first observed.
        """
        __props__['eventTime'] = event_time

        if first_timestamp and not isinstance(first_timestamp, str):
            raise TypeError('Expected property aliases to be a str')
        self.first_timestamp = first_timestamp
        """
        The time at which the event was first recorded. (Time of server receipt is in TypeMeta.)
        """
        __props__['firstTimestamp'] = first_timestamp

        if last_timestamp and not isinstance(last_timestamp, str):
            raise TypeError('Expected property aliases to be a str')
        self.last_timestamp = last_timestamp
        """
        The time at which the most recent occurrence of this event was recorded.
        """
        __props__['lastTimestamp'] = last_timestamp

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

        if reporting_component and not isinstance(reporting_component, str):
            raise TypeError('Expected property aliases to be a str')
        self.reporting_component = reporting_component
        """
        Name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`.
        """
        __props__['reportingComponent'] = reporting_component

        if reporting_instance and not isinstance(reporting_instance, str):
            raise TypeError('Expected property aliases to be a str')
        self.reporting_instance = reporting_instance
        """
        ID of the controller instance, e.g. `kubelet-xyzf`.
        """
        __props__['reportingInstance'] = reporting_instance

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

    def translate_output_property(self, prop: str) -> str:
        return _CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return _CASING_BACKWARD_TABLE.get(prop) or prop
