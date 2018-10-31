import pulumi
import pulumi.runtime

class Event(pulumi.CustomResource):
    """
    Event is a report of an event somewhere in the cluster. It generally denotes some state change
    in the system.
    """
    def __init__(self, __name__, __opts__=None, action=None, deprecated_count=None, deprecated_first_timestamp=None, deprecated_last_timestamp=None, deprecated_source=None, event_time=None, metadata=None, note=None, reason=None, regarding=None, related=None, reporting_controller=None, reporting_instance=None, series=None, type=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'events.k8s.io/v1beta1'
        self.apiVersion = 'events.k8s.io/v1beta1'

        __props__['kind'] = 'Event'
        self.kind = 'Event'

        if not eventTime:
            raise TypeError('Missing required property eventTime')
        elif not isinstance(eventTime, str):
            raise TypeError('Expected property aliases to be a str')
        self.event_time = event_time
        """
        Required. Time when this Event was first observed.
        """
        __props__['eventTime'] = event_time

        if action and not isinstance(action, str):
            raise TypeError('Expected property aliases to be a str')
        self.action = action
        """
        What action was taken/failed regarding to the regarding object.
        """
        __props__['action'] = action

        if deprecated_count and not isinstance(deprecated_count, int):
            raise TypeError('Expected property aliases to be a int')
        self.deprecated_count = deprecated_count
        """
        Deprecated field assuring backward compatibility with core.v1 Event type
        """
        __props__['deprecatedCount'] = deprecated_count

        if deprecated_first_timestamp and not isinstance(deprecated_first_timestamp, str):
            raise TypeError('Expected property aliases to be a str')
        self.deprecated_first_timestamp = deprecated_first_timestamp
        """
        Deprecated field assuring backward compatibility with core.v1 Event type
        """
        __props__['deprecatedFirstTimestamp'] = deprecated_first_timestamp

        if deprecated_last_timestamp and not isinstance(deprecated_last_timestamp, str):
            raise TypeError('Expected property aliases to be a str')
        self.deprecated_last_timestamp = deprecated_last_timestamp
        """
        Deprecated field assuring backward compatibility with core.v1 Event type
        """
        __props__['deprecatedLastTimestamp'] = deprecated_last_timestamp

        if deprecated_source and not isinstance(deprecated_source, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.deprecated_source = deprecated_source
        """
        Deprecated field assuring backward compatibility with core.v1 Event type
        """
        __props__['deprecatedSource'] = deprecated_source

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        
        __props__['metadata'] = metadata

        if note and not isinstance(note, str):
            raise TypeError('Expected property aliases to be a str')
        self.note = note
        """
        Optional. A human-readable description of the status of this operation. Maximal length of
        the note is 1kB, but libraries should be prepared to handle values up to 64kB.
        """
        __props__['note'] = note

        if reason and not isinstance(reason, str):
            raise TypeError('Expected property aliases to be a str')
        self.reason = reason
        """
        Why the action was taken.
        """
        __props__['reason'] = reason

        if regarding and not isinstance(regarding, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.regarding = regarding
        """
        The object this Event is about. In most cases it's an Object reporting controller
        implements. E.g. ReplicaSetController implements ReplicaSets and this event is emitted
        because it acts on some changes in a ReplicaSet object.
        """
        __props__['regarding'] = regarding

        if related and not isinstance(related, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.related = related
        """
        Optional secondary object for more complex actions. E.g. when regarding object triggers a
        creation or deletion of related object.
        """
        __props__['related'] = related

        if reporting_controller and not isinstance(reporting_controller, str):
            raise TypeError('Expected property aliases to be a str')
        self.reporting_controller = reporting_controller
        """
        Name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`.
        """
        __props__['reportingController'] = reporting_controller

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

        if type and not isinstance(type, str):
            raise TypeError('Expected property aliases to be a str')
        self.type = type
        """
        Type of this event (Normal, Warning), new types could be added in the future.
        """
        __props__['type'] = type

        super(Event, self).__init__(
            "kubernetes:events.k8s.io/v1beta1:Event",
            __name__,
            __props__,
            __opts__)
