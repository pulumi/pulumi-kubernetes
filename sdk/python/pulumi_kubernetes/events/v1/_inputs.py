# coding=utf-8
# *** WARNING: this file was generated by pulumigen. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import builtins as _builtins
import warnings
import sys
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
if sys.version_info >= (3, 11):
    from typing import NotRequired, TypedDict, TypeAlias
else:
    from typing_extensions import NotRequired, TypedDict, TypeAlias
from ... import _utilities
from ... import core as _core
from ... import meta as _meta

__all__ = [
    'EventSeriesPatchArgs',
    'EventSeriesPatchArgsDict',
    'EventSeriesArgs',
    'EventSeriesArgsDict',
    'EventArgs',
    'EventArgsDict',
]

MYPY = False

if not MYPY:
    class EventSeriesPatchArgsDict(TypedDict):
        """
        EventSeries contain information on series of events, i.e. thing that was/is happening continuously for some time. How often to update the EventSeries is up to the event reporters. The default event reporter in "k8s.io/client-go/tools/events/event_broadcaster.go" shows how this struct is updated on heartbeats and can guide customized reporter implementations.
        """
        count: NotRequired[pulumi.Input[_builtins.int]]
        """
        count is the number of occurrences in this series up to the last heartbeat time.
        """
        last_observed_time: NotRequired[pulumi.Input[_builtins.str]]
        """
        lastObservedTime is the time when last Event from the series was seen before last heartbeat.
        """
elif False:
    EventSeriesPatchArgsDict: TypeAlias = Mapping[str, Any]

@pulumi.input_type
class EventSeriesPatchArgs:
    def __init__(__self__, *,
                 count: Optional[pulumi.Input[_builtins.int]] = None,
                 last_observed_time: Optional[pulumi.Input[_builtins.str]] = None):
        """
        EventSeries contain information on series of events, i.e. thing that was/is happening continuously for some time. How often to update the EventSeries is up to the event reporters. The default event reporter in "k8s.io/client-go/tools/events/event_broadcaster.go" shows how this struct is updated on heartbeats and can guide customized reporter implementations.
        :param pulumi.Input[_builtins.int] count: count is the number of occurrences in this series up to the last heartbeat time.
        :param pulumi.Input[_builtins.str] last_observed_time: lastObservedTime is the time when last Event from the series was seen before last heartbeat.
        """
        if count is not None:
            pulumi.set(__self__, "count", count)
        if last_observed_time is not None:
            pulumi.set(__self__, "last_observed_time", last_observed_time)

    @_builtins.property
    @pulumi.getter
    def count(self) -> Optional[pulumi.Input[_builtins.int]]:
        """
        count is the number of occurrences in this series up to the last heartbeat time.
        """
        return pulumi.get(self, "count")

    @count.setter
    def count(self, value: Optional[pulumi.Input[_builtins.int]]):
        pulumi.set(self, "count", value)

    @_builtins.property
    @pulumi.getter(name="lastObservedTime")
    def last_observed_time(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        lastObservedTime is the time when last Event from the series was seen before last heartbeat.
        """
        return pulumi.get(self, "last_observed_time")

    @last_observed_time.setter
    def last_observed_time(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "last_observed_time", value)


if not MYPY:
    class EventSeriesArgsDict(TypedDict):
        """
        EventSeries contain information on series of events, i.e. thing that was/is happening continuously for some time. How often to update the EventSeries is up to the event reporters. The default event reporter in "k8s.io/client-go/tools/events/event_broadcaster.go" shows how this struct is updated on heartbeats and can guide customized reporter implementations.
        """
        count: pulumi.Input[_builtins.int]
        """
        count is the number of occurrences in this series up to the last heartbeat time.
        """
        last_observed_time: pulumi.Input[_builtins.str]
        """
        lastObservedTime is the time when last Event from the series was seen before last heartbeat.
        """
elif False:
    EventSeriesArgsDict: TypeAlias = Mapping[str, Any]

@pulumi.input_type
class EventSeriesArgs:
    def __init__(__self__, *,
                 count: pulumi.Input[_builtins.int],
                 last_observed_time: pulumi.Input[_builtins.str]):
        """
        EventSeries contain information on series of events, i.e. thing that was/is happening continuously for some time. How often to update the EventSeries is up to the event reporters. The default event reporter in "k8s.io/client-go/tools/events/event_broadcaster.go" shows how this struct is updated on heartbeats and can guide customized reporter implementations.
        :param pulumi.Input[_builtins.int] count: count is the number of occurrences in this series up to the last heartbeat time.
        :param pulumi.Input[_builtins.str] last_observed_time: lastObservedTime is the time when last Event from the series was seen before last heartbeat.
        """
        pulumi.set(__self__, "count", count)
        pulumi.set(__self__, "last_observed_time", last_observed_time)

    @_builtins.property
    @pulumi.getter
    def count(self) -> pulumi.Input[_builtins.int]:
        """
        count is the number of occurrences in this series up to the last heartbeat time.
        """
        return pulumi.get(self, "count")

    @count.setter
    def count(self, value: pulumi.Input[_builtins.int]):
        pulumi.set(self, "count", value)

    @_builtins.property
    @pulumi.getter(name="lastObservedTime")
    def last_observed_time(self) -> pulumi.Input[_builtins.str]:
        """
        lastObservedTime is the time when last Event from the series was seen before last heartbeat.
        """
        return pulumi.get(self, "last_observed_time")

    @last_observed_time.setter
    def last_observed_time(self, value: pulumi.Input[_builtins.str]):
        pulumi.set(self, "last_observed_time", value)


if not MYPY:
    class EventArgsDict(TypedDict):
        """
        Event is a report of an event somewhere in the cluster. It generally denotes some state change in the system. Events have a limited retention time and triggers and messages may evolve with time.  Event consumers should not rely on the timing of an event with a given Reason reflecting a consistent underlying trigger, or the continued existence of events with that Reason.  Events should be treated as informative, best-effort, supplemental data.
        """
        event_time: pulumi.Input[_builtins.str]
        """
        eventTime is the time when this Event was first observed. It is required.
        """
        action: NotRequired[pulumi.Input[_builtins.str]]
        """
        action is what action was taken/failed regarding to the regarding object. It is machine-readable. This field cannot be empty for new Events and it can have at most 128 characters.
        """
        api_version: NotRequired[pulumi.Input[_builtins.str]]
        """
        APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        """
        deprecated_count: NotRequired[pulumi.Input[_builtins.int]]
        """
        deprecatedCount is the deprecated field assuring backward compatibility with core.v1 Event type.
        """
        deprecated_first_timestamp: NotRequired[pulumi.Input[_builtins.str]]
        """
        deprecatedFirstTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
        """
        deprecated_last_timestamp: NotRequired[pulumi.Input[_builtins.str]]
        """
        deprecatedLastTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
        """
        deprecated_source: NotRequired[pulumi.Input['_core.v1.EventSourceArgsDict']]
        """
        deprecatedSource is the deprecated field assuring backward compatibility with core.v1 Event type.
        """
        kind: NotRequired[pulumi.Input[_builtins.str]]
        """
        Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        """
        metadata: NotRequired[pulumi.Input['_meta.v1.ObjectMetaArgsDict']]
        """
        Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        """
        note: NotRequired[pulumi.Input[_builtins.str]]
        """
        note is a human-readable description of the status of this operation. Maximal length of the note is 1kB, but libraries should be prepared to handle values up to 64kB.
        """
        reason: NotRequired[pulumi.Input[_builtins.str]]
        """
        reason is why the action was taken. It is human-readable. This field cannot be empty for new Events and it can have at most 128 characters.
        """
        regarding: NotRequired[pulumi.Input['_core.v1.ObjectReferenceArgsDict']]
        """
        regarding contains the object this Event is about. In most cases it's an Object reporting controller implements, e.g. ReplicaSetController implements ReplicaSets and this event is emitted because it acts on some changes in a ReplicaSet object.
        """
        related: NotRequired[pulumi.Input['_core.v1.ObjectReferenceArgsDict']]
        """
        related is the optional secondary object for more complex actions. E.g. when regarding object triggers a creation or deletion of related object.
        """
        reporting_controller: NotRequired[pulumi.Input[_builtins.str]]
        """
        reportingController is the name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`. This field cannot be empty for new Events.
        """
        reporting_instance: NotRequired[pulumi.Input[_builtins.str]]
        """
        reportingInstance is the ID of the controller instance, e.g. `kubelet-xyzf`. This field cannot be empty for new Events and it can have at most 128 characters.
        """
        series: NotRequired[pulumi.Input['EventSeriesArgsDict']]
        """
        series is data about the Event series this event represents or nil if it's a singleton Event.
        """
        type: NotRequired[pulumi.Input[_builtins.str]]
        """
        type is the type of this event (Normal, Warning), new types could be added in the future. It is machine-readable. This field cannot be empty for new Events.
        """
elif False:
    EventArgsDict: TypeAlias = Mapping[str, Any]

@pulumi.input_type
class EventArgs:
    def __init__(__self__, *,
                 event_time: pulumi.Input[_builtins.str],
                 action: Optional[pulumi.Input[_builtins.str]] = None,
                 api_version: Optional[pulumi.Input[_builtins.str]] = None,
                 deprecated_count: Optional[pulumi.Input[_builtins.int]] = None,
                 deprecated_first_timestamp: Optional[pulumi.Input[_builtins.str]] = None,
                 deprecated_last_timestamp: Optional[pulumi.Input[_builtins.str]] = None,
                 deprecated_source: Optional[pulumi.Input['_core.v1.EventSourceArgs']] = None,
                 kind: Optional[pulumi.Input[_builtins.str]] = None,
                 metadata: Optional[pulumi.Input['_meta.v1.ObjectMetaArgs']] = None,
                 note: Optional[pulumi.Input[_builtins.str]] = None,
                 reason: Optional[pulumi.Input[_builtins.str]] = None,
                 regarding: Optional[pulumi.Input['_core.v1.ObjectReferenceArgs']] = None,
                 related: Optional[pulumi.Input['_core.v1.ObjectReferenceArgs']] = None,
                 reporting_controller: Optional[pulumi.Input[_builtins.str]] = None,
                 reporting_instance: Optional[pulumi.Input[_builtins.str]] = None,
                 series: Optional[pulumi.Input['EventSeriesArgs']] = None,
                 type: Optional[pulumi.Input[_builtins.str]] = None):
        """
        Event is a report of an event somewhere in the cluster. It generally denotes some state change in the system. Events have a limited retention time and triggers and messages may evolve with time.  Event consumers should not rely on the timing of an event with a given Reason reflecting a consistent underlying trigger, or the continued existence of events with that Reason.  Events should be treated as informative, best-effort, supplemental data.
        :param pulumi.Input[_builtins.str] event_time: eventTime is the time when this Event was first observed. It is required.
        :param pulumi.Input[_builtins.str] action: action is what action was taken/failed regarding to the regarding object. It is machine-readable. This field cannot be empty for new Events and it can have at most 128 characters.
        :param pulumi.Input[_builtins.str] api_version: APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        :param pulumi.Input[_builtins.int] deprecated_count: deprecatedCount is the deprecated field assuring backward compatibility with core.v1 Event type.
        :param pulumi.Input[_builtins.str] deprecated_first_timestamp: deprecatedFirstTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
        :param pulumi.Input[_builtins.str] deprecated_last_timestamp: deprecatedLastTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
        :param pulumi.Input['_core.v1.EventSourceArgs'] deprecated_source: deprecatedSource is the deprecated field assuring backward compatibility with core.v1 Event type.
        :param pulumi.Input[_builtins.str] kind: Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        :param pulumi.Input['_meta.v1.ObjectMetaArgs'] metadata: Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        :param pulumi.Input[_builtins.str] note: note is a human-readable description of the status of this operation. Maximal length of the note is 1kB, but libraries should be prepared to handle values up to 64kB.
        :param pulumi.Input[_builtins.str] reason: reason is why the action was taken. It is human-readable. This field cannot be empty for new Events and it can have at most 128 characters.
        :param pulumi.Input['_core.v1.ObjectReferenceArgs'] regarding: regarding contains the object this Event is about. In most cases it's an Object reporting controller implements, e.g. ReplicaSetController implements ReplicaSets and this event is emitted because it acts on some changes in a ReplicaSet object.
        :param pulumi.Input['_core.v1.ObjectReferenceArgs'] related: related is the optional secondary object for more complex actions. E.g. when regarding object triggers a creation or deletion of related object.
        :param pulumi.Input[_builtins.str] reporting_controller: reportingController is the name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`. This field cannot be empty for new Events.
        :param pulumi.Input[_builtins.str] reporting_instance: reportingInstance is the ID of the controller instance, e.g. `kubelet-xyzf`. This field cannot be empty for new Events and it can have at most 128 characters.
        :param pulumi.Input['EventSeriesArgs'] series: series is data about the Event series this event represents or nil if it's a singleton Event.
        :param pulumi.Input[_builtins.str] type: type is the type of this event (Normal, Warning), new types could be added in the future. It is machine-readable. This field cannot be empty for new Events.
        """
        pulumi.set(__self__, "event_time", event_time)
        if action is not None:
            pulumi.set(__self__, "action", action)
        if api_version is not None:
            pulumi.set(__self__, "api_version", 'events.k8s.io/v1')
        if deprecated_count is not None:
            pulumi.set(__self__, "deprecated_count", deprecated_count)
        if deprecated_first_timestamp is not None:
            pulumi.set(__self__, "deprecated_first_timestamp", deprecated_first_timestamp)
        if deprecated_last_timestamp is not None:
            pulumi.set(__self__, "deprecated_last_timestamp", deprecated_last_timestamp)
        if deprecated_source is not None:
            pulumi.set(__self__, "deprecated_source", deprecated_source)
        if kind is not None:
            pulumi.set(__self__, "kind", 'Event')
        if metadata is not None:
            pulumi.set(__self__, "metadata", metadata)
        if note is not None:
            pulumi.set(__self__, "note", note)
        if reason is not None:
            pulumi.set(__self__, "reason", reason)
        if regarding is not None:
            pulumi.set(__self__, "regarding", regarding)
        if related is not None:
            pulumi.set(__self__, "related", related)
        if reporting_controller is not None:
            pulumi.set(__self__, "reporting_controller", reporting_controller)
        if reporting_instance is not None:
            pulumi.set(__self__, "reporting_instance", reporting_instance)
        if series is not None:
            pulumi.set(__self__, "series", series)
        if type is not None:
            pulumi.set(__self__, "type", type)

    @_builtins.property
    @pulumi.getter(name="eventTime")
    def event_time(self) -> pulumi.Input[_builtins.str]:
        """
        eventTime is the time when this Event was first observed. It is required.
        """
        return pulumi.get(self, "event_time")

    @event_time.setter
    def event_time(self, value: pulumi.Input[_builtins.str]):
        pulumi.set(self, "event_time", value)

    @_builtins.property
    @pulumi.getter
    def action(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        action is what action was taken/failed regarding to the regarding object. It is machine-readable. This field cannot be empty for new Events and it can have at most 128 characters.
        """
        return pulumi.get(self, "action")

    @action.setter
    def action(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "action", value)

    @_builtins.property
    @pulumi.getter(name="apiVersion")
    def api_version(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        """
        return pulumi.get(self, "api_version")

    @api_version.setter
    def api_version(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "api_version", value)

    @_builtins.property
    @pulumi.getter(name="deprecatedCount")
    def deprecated_count(self) -> Optional[pulumi.Input[_builtins.int]]:
        """
        deprecatedCount is the deprecated field assuring backward compatibility with core.v1 Event type.
        """
        return pulumi.get(self, "deprecated_count")

    @deprecated_count.setter
    def deprecated_count(self, value: Optional[pulumi.Input[_builtins.int]]):
        pulumi.set(self, "deprecated_count", value)

    @_builtins.property
    @pulumi.getter(name="deprecatedFirstTimestamp")
    def deprecated_first_timestamp(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        deprecatedFirstTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
        """
        return pulumi.get(self, "deprecated_first_timestamp")

    @deprecated_first_timestamp.setter
    def deprecated_first_timestamp(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "deprecated_first_timestamp", value)

    @_builtins.property
    @pulumi.getter(name="deprecatedLastTimestamp")
    def deprecated_last_timestamp(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        deprecatedLastTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
        """
        return pulumi.get(self, "deprecated_last_timestamp")

    @deprecated_last_timestamp.setter
    def deprecated_last_timestamp(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "deprecated_last_timestamp", value)

    @_builtins.property
    @pulumi.getter(name="deprecatedSource")
    def deprecated_source(self) -> Optional[pulumi.Input['_core.v1.EventSourceArgs']]:
        """
        deprecatedSource is the deprecated field assuring backward compatibility with core.v1 Event type.
        """
        return pulumi.get(self, "deprecated_source")

    @deprecated_source.setter
    def deprecated_source(self, value: Optional[pulumi.Input['_core.v1.EventSourceArgs']]):
        pulumi.set(self, "deprecated_source", value)

    @_builtins.property
    @pulumi.getter
    def kind(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        """
        return pulumi.get(self, "kind")

    @kind.setter
    def kind(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "kind", value)

    @_builtins.property
    @pulumi.getter
    def metadata(self) -> Optional[pulumi.Input['_meta.v1.ObjectMetaArgs']]:
        """
        Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
        """
        return pulumi.get(self, "metadata")

    @metadata.setter
    def metadata(self, value: Optional[pulumi.Input['_meta.v1.ObjectMetaArgs']]):
        pulumi.set(self, "metadata", value)

    @_builtins.property
    @pulumi.getter
    def note(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        note is a human-readable description of the status of this operation. Maximal length of the note is 1kB, but libraries should be prepared to handle values up to 64kB.
        """
        return pulumi.get(self, "note")

    @note.setter
    def note(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "note", value)

    @_builtins.property
    @pulumi.getter
    def reason(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        reason is why the action was taken. It is human-readable. This field cannot be empty for new Events and it can have at most 128 characters.
        """
        return pulumi.get(self, "reason")

    @reason.setter
    def reason(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "reason", value)

    @_builtins.property
    @pulumi.getter
    def regarding(self) -> Optional[pulumi.Input['_core.v1.ObjectReferenceArgs']]:
        """
        regarding contains the object this Event is about. In most cases it's an Object reporting controller implements, e.g. ReplicaSetController implements ReplicaSets and this event is emitted because it acts on some changes in a ReplicaSet object.
        """
        return pulumi.get(self, "regarding")

    @regarding.setter
    def regarding(self, value: Optional[pulumi.Input['_core.v1.ObjectReferenceArgs']]):
        pulumi.set(self, "regarding", value)

    @_builtins.property
    @pulumi.getter
    def related(self) -> Optional[pulumi.Input['_core.v1.ObjectReferenceArgs']]:
        """
        related is the optional secondary object for more complex actions. E.g. when regarding object triggers a creation or deletion of related object.
        """
        return pulumi.get(self, "related")

    @related.setter
    def related(self, value: Optional[pulumi.Input['_core.v1.ObjectReferenceArgs']]):
        pulumi.set(self, "related", value)

    @_builtins.property
    @pulumi.getter(name="reportingController")
    def reporting_controller(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        reportingController is the name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`. This field cannot be empty for new Events.
        """
        return pulumi.get(self, "reporting_controller")

    @reporting_controller.setter
    def reporting_controller(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "reporting_controller", value)

    @_builtins.property
    @pulumi.getter(name="reportingInstance")
    def reporting_instance(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        reportingInstance is the ID of the controller instance, e.g. `kubelet-xyzf`. This field cannot be empty for new Events and it can have at most 128 characters.
        """
        return pulumi.get(self, "reporting_instance")

    @reporting_instance.setter
    def reporting_instance(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "reporting_instance", value)

    @_builtins.property
    @pulumi.getter
    def series(self) -> Optional[pulumi.Input['EventSeriesArgs']]:
        """
        series is data about the Event series this event represents or nil if it's a singleton Event.
        """
        return pulumi.get(self, "series")

    @series.setter
    def series(self, value: Optional[pulumi.Input['EventSeriesArgs']]):
        pulumi.set(self, "series", value)

    @_builtins.property
    @pulumi.getter
    def type(self) -> Optional[pulumi.Input[_builtins.str]]:
        """
        type is the type of this event (Normal, Warning), new types could be added in the future. It is machine-readable. This field cannot be empty for new Events.
        """
        return pulumi.get(self, "type")

    @type.setter
    def type(self, value: Optional[pulumi.Input[_builtins.str]]):
        pulumi.set(self, "type", value)


