import pulumi
import pulumi.runtime

class Event(pulumi.CustomResource):
    """
    Event is a report of an event somewhere in the cluster. It generally denotes some state change
    in the system.
    """
    def __init__(self, __name__, __opts__=None, action=None, deprecatedCount=None, deprecatedFirstTimestamp=None, deprecatedLastTimestamp=None, deprecatedSource=None, eventTime=None, metadata=None, note=None, reason=None, regarding=None, related=None, reportingController=None, reportingInstance=None, series=None, type=None):
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
        self.eventTime = eventTime
        """
        Required. Time when this Event was first observed.
        """
        __props__['eventTime'] = eventTime

        if action and not isinstance(action, str):
            raise TypeError('Expected property aliases to be a str')
        self.action = action
        """
        What action was taken/failed regarding to the regarding object.
        """
        __props__['action'] = action

        if deprecatedCount and not isinstance(deprecatedCount, int):
            raise TypeError('Expected property aliases to be a int')
        self.deprecatedCount = deprecatedCount
        """
        Deprecated field assuring backward compatibility with core.v1 Event type
        """
        __props__['deprecatedCount'] = deprecatedCount

        if deprecatedFirstTimestamp and not isinstance(deprecatedFirstTimestamp, str):
            raise TypeError('Expected property aliases to be a str')
        self.deprecatedFirstTimestamp = deprecatedFirstTimestamp
        """
        Deprecated field assuring backward compatibility with core.v1 Event type
        """
        __props__['deprecatedFirstTimestamp'] = deprecatedFirstTimestamp

        if deprecatedLastTimestamp and not isinstance(deprecatedLastTimestamp, str):
            raise TypeError('Expected property aliases to be a str')
        self.deprecatedLastTimestamp = deprecatedLastTimestamp
        """
        Deprecated field assuring backward compatibility with core.v1 Event type
        """
        __props__['deprecatedLastTimestamp'] = deprecatedLastTimestamp

        if deprecatedSource and not isinstance(deprecatedSource, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.deprecatedSource = deprecatedSource
        """
        Deprecated field assuring backward compatibility with core.v1 Event type
        """
        __props__['deprecatedSource'] = deprecatedSource

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

        if reportingController and not isinstance(reportingController, str):
            raise TypeError('Expected property aliases to be a str')
        self.reportingController = reportingController
        """
        Name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`.
        """
        __props__['reportingController'] = reportingController

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
