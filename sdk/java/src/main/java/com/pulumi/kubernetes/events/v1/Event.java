// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.events.v1;

import com.pulumi.core.Alias;
import com.pulumi.core.Output;
import com.pulumi.core.annotations.Export;
import com.pulumi.core.annotations.ResourceType;
import com.pulumi.core.internal.Codegen;
import com.pulumi.kubernetes.Utilities;
import com.pulumi.kubernetes.core.v1.outputs.EventSource;
import com.pulumi.kubernetes.core.v1.outputs.ObjectReference;
import com.pulumi.kubernetes.events.v1.EventArgs;
import com.pulumi.kubernetes.events.v1.outputs.EventSeries;
import com.pulumi.kubernetes.meta.v1.outputs.ObjectMeta;
import java.lang.Integer;
import java.lang.String;
import java.util.List;
import javax.annotation.Nullable;

/**
 * Event is a report of an event somewhere in the cluster. It generally denotes some state change in the system. Events have a limited retention time and triggers and messages may evolve with time.  Event consumers should not rely on the timing of an event with a given Reason reflecting a consistent underlying trigger, or the continued existence of events with that Reason.  Events should be treated as informative, best-effort, supplemental data.
 * 
 */
@ResourceType(type="kubernetes:events.k8s.io/v1:Event")
public class Event extends com.pulumi.resources.CustomResource {
    /**
     * action is what action was taken/failed regarding to the regarding object. It is machine-readable. This field cannot be empty for new Events and it can have at most 128 characters.
     * 
     */
    @Export(name="action", refs={String.class}, tree="[0]")
    private Output<String> action;

    /**
     * @return action is what action was taken/failed regarding to the regarding object. It is machine-readable. This field cannot be empty for new Events and it can have at most 128 characters.
     * 
     */
    public Output<String> action() {
        return this.action;
    }
    /**
     * APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
     * 
     */
    @Export(name="apiVersion", refs={String.class}, tree="[0]")
    private Output<String> apiVersion;

    /**
     * @return APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
     * 
     */
    public Output<String> apiVersion() {
        return this.apiVersion;
    }
    /**
     * deprecatedCount is the deprecated field assuring backward compatibility with core.v1 Event type.
     * 
     */
    @Export(name="deprecatedCount", refs={Integer.class}, tree="[0]")
    private Output<Integer> deprecatedCount;

    /**
     * @return deprecatedCount is the deprecated field assuring backward compatibility with core.v1 Event type.
     * 
     */
    public Output<Integer> deprecatedCount() {
        return this.deprecatedCount;
    }
    /**
     * deprecatedFirstTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
     * 
     */
    @Export(name="deprecatedFirstTimestamp", refs={String.class}, tree="[0]")
    private Output<String> deprecatedFirstTimestamp;

    /**
     * @return deprecatedFirstTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
     * 
     */
    public Output<String> deprecatedFirstTimestamp() {
        return this.deprecatedFirstTimestamp;
    }
    /**
     * deprecatedLastTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
     * 
     */
    @Export(name="deprecatedLastTimestamp", refs={String.class}, tree="[0]")
    private Output<String> deprecatedLastTimestamp;

    /**
     * @return deprecatedLastTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
     * 
     */
    public Output<String> deprecatedLastTimestamp() {
        return this.deprecatedLastTimestamp;
    }
    /**
     * deprecatedSource is the deprecated field assuring backward compatibility with core.v1 Event type.
     * 
     */
    @Export(name="deprecatedSource", refs={EventSource.class}, tree="[0]")
    private Output<EventSource> deprecatedSource;

    /**
     * @return deprecatedSource is the deprecated field assuring backward compatibility with core.v1 Event type.
     * 
     */
    public Output<EventSource> deprecatedSource() {
        return this.deprecatedSource;
    }
    /**
     * eventTime is the time when this Event was first observed. It is required.
     * 
     */
    @Export(name="eventTime", refs={String.class}, tree="[0]")
    private Output<String> eventTime;

    /**
     * @return eventTime is the time when this Event was first observed. It is required.
     * 
     */
    public Output<String> eventTime() {
        return this.eventTime;
    }
    /**
     * Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
     * 
     */
    @Export(name="kind", refs={String.class}, tree="[0]")
    private Output<String> kind;

    /**
     * @return Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
     * 
     */
    public Output<String> kind() {
        return this.kind;
    }
    /**
     * Standard object&#39;s metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
     * 
     */
    @Export(name="metadata", refs={ObjectMeta.class}, tree="[0]")
    private Output<ObjectMeta> metadata;

    /**
     * @return Standard object&#39;s metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
     * 
     */
    public Output<ObjectMeta> metadata() {
        return this.metadata;
    }
    /**
     * note is a human-readable description of the status of this operation. Maximal length of the note is 1kB, but libraries should be prepared to handle values up to 64kB.
     * 
     */
    @Export(name="note", refs={String.class}, tree="[0]")
    private Output<String> note;

    /**
     * @return note is a human-readable description of the status of this operation. Maximal length of the note is 1kB, but libraries should be prepared to handle values up to 64kB.
     * 
     */
    public Output<String> note() {
        return this.note;
    }
    /**
     * reason is why the action was taken. It is human-readable. This field cannot be empty for new Events and it can have at most 128 characters.
     * 
     */
    @Export(name="reason", refs={String.class}, tree="[0]")
    private Output<String> reason;

    /**
     * @return reason is why the action was taken. It is human-readable. This field cannot be empty for new Events and it can have at most 128 characters.
     * 
     */
    public Output<String> reason() {
        return this.reason;
    }
    /**
     * regarding contains the object this Event is about. In most cases it&#39;s an Object reporting controller implements, e.g. ReplicaSetController implements ReplicaSets and this event is emitted because it acts on some changes in a ReplicaSet object.
     * 
     */
    @Export(name="regarding", refs={ObjectReference.class}, tree="[0]")
    private Output<ObjectReference> regarding;

    /**
     * @return regarding contains the object this Event is about. In most cases it&#39;s an Object reporting controller implements, e.g. ReplicaSetController implements ReplicaSets and this event is emitted because it acts on some changes in a ReplicaSet object.
     * 
     */
    public Output<ObjectReference> regarding() {
        return this.regarding;
    }
    /**
     * related is the optional secondary object for more complex actions. E.g. when regarding object triggers a creation or deletion of related object.
     * 
     */
    @Export(name="related", refs={ObjectReference.class}, tree="[0]")
    private Output<ObjectReference> related;

    /**
     * @return related is the optional secondary object for more complex actions. E.g. when regarding object triggers a creation or deletion of related object.
     * 
     */
    public Output<ObjectReference> related() {
        return this.related;
    }
    /**
     * reportingController is the name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`. This field cannot be empty for new Events.
     * 
     */
    @Export(name="reportingController", refs={String.class}, tree="[0]")
    private Output<String> reportingController;

    /**
     * @return reportingController is the name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`. This field cannot be empty for new Events.
     * 
     */
    public Output<String> reportingController() {
        return this.reportingController;
    }
    /**
     * reportingInstance is the ID of the controller instance, e.g. `kubelet-xyzf`. This field cannot be empty for new Events and it can have at most 128 characters.
     * 
     */
    @Export(name="reportingInstance", refs={String.class}, tree="[0]")
    private Output<String> reportingInstance;

    /**
     * @return reportingInstance is the ID of the controller instance, e.g. `kubelet-xyzf`. This field cannot be empty for new Events and it can have at most 128 characters.
     * 
     */
    public Output<String> reportingInstance() {
        return this.reportingInstance;
    }
    /**
     * series is data about the Event series this event represents or nil if it&#39;s a singleton Event.
     * 
     */
    @Export(name="series", refs={EventSeries.class}, tree="[0]")
    private Output<EventSeries> series;

    /**
     * @return series is data about the Event series this event represents or nil if it&#39;s a singleton Event.
     * 
     */
    public Output<EventSeries> series() {
        return this.series;
    }
    /**
     * type is the type of this event (Normal, Warning), new types could be added in the future. It is machine-readable. This field cannot be empty for new Events.
     * 
     */
    @Export(name="type", refs={String.class}, tree="[0]")
    private Output<String> type;

    /**
     * @return type is the type of this event (Normal, Warning), new types could be added in the future. It is machine-readable. This field cannot be empty for new Events.
     * 
     */
    public Output<String> type() {
        return this.type;
    }

    /**
     *
     * @param name The _unique_ name of the resulting resource.
     */
    public Event(java.lang.String name) {
        this(name, EventArgs.Empty);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     */
    public Event(java.lang.String name, EventArgs args) {
        this(name, args, null);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param options A bag of options that control this resource's behavior.
     */
    public Event(java.lang.String name, EventArgs args, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super("kubernetes:events.k8s.io/v1:Event", name, makeArgs(args, options), makeResourceOptions(options, Codegen.empty()), false);
    }

    private Event(java.lang.String name, Output<java.lang.String> id, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super("kubernetes:events.k8s.io/v1:Event", name, null, makeResourceOptions(options, id), false);
    }

    private static EventArgs makeArgs(EventArgs args, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        if (options != null && options.getUrn().isPresent()) {
            return null;
        }
        var builder = args == null ? EventArgs.builder() : EventArgs.builder(args);
        return builder
            .apiVersion("events.k8s.io/v1")
            .kind("Event")
            .build();
    }

    private static com.pulumi.resources.CustomResourceOptions makeResourceOptions(@Nullable com.pulumi.resources.CustomResourceOptions options, @Nullable Output<java.lang.String> id) {
        var defaultOptions = com.pulumi.resources.CustomResourceOptions.builder()
            .version(Utilities.getVersion())
            .aliases(List.of(
                Output.of(Alias.builder().type("kubernetes:core/v1:Event").build()),
                Output.of(Alias.builder().type("kubernetes:events.k8s.io/v1beta1:Event").build())
            ))
            .build();
        return com.pulumi.resources.CustomResourceOptions.merge(defaultOptions, options, id);
    }

    /**
     * Get an existing Host resource's state with the given name, ID, and optional extra
     * properties used to qualify the lookup.
     *
     * @param name The _unique_ name of the resulting resource.
     * @param id The _unique_ provider ID of the resource to lookup.
     * @param options Optional settings to control the behavior of the CustomResource.
     */
    public static Event get(java.lang.String name, Output<java.lang.String> id, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        return new Event(name, id, options);
    }
}
