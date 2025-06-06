// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.events.v1beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * EventSeries contain information on series of events, i.e. thing that was/is happening continuously for some time.
 * 
 */
public final class EventSeriesPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final EventSeriesPatchArgs Empty = new EventSeriesPatchArgs();

    /**
     * Number of occurrences in this series up to the last heartbeat time
     * 
     */
    @Import(name="count")
    private @Nullable Output<Integer> count;

    /**
     * @return Number of occurrences in this series up to the last heartbeat time
     * 
     */
    public Optional<Output<Integer>> count() {
        return Optional.ofNullable(this.count);
    }

    /**
     * Time when last Event from the series was seen before last heartbeat.
     * 
     */
    @Import(name="lastObservedTime")
    private @Nullable Output<String> lastObservedTime;

    /**
     * @return Time when last Event from the series was seen before last heartbeat.
     * 
     */
    public Optional<Output<String>> lastObservedTime() {
        return Optional.ofNullable(this.lastObservedTime);
    }

    /**
     * Information whether this series is ongoing or finished. Deprecated. Planned removal for 1.18
     * 
     */
    @Import(name="state")
    private @Nullable Output<String> state;

    /**
     * @return Information whether this series is ongoing or finished. Deprecated. Planned removal for 1.18
     * 
     */
    public Optional<Output<String>> state() {
        return Optional.ofNullable(this.state);
    }

    private EventSeriesPatchArgs() {}

    private EventSeriesPatchArgs(EventSeriesPatchArgs $) {
        this.count = $.count;
        this.lastObservedTime = $.lastObservedTime;
        this.state = $.state;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(EventSeriesPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private EventSeriesPatchArgs $;

        public Builder() {
            $ = new EventSeriesPatchArgs();
        }

        public Builder(EventSeriesPatchArgs defaults) {
            $ = new EventSeriesPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param count Number of occurrences in this series up to the last heartbeat time
         * 
         * @return builder
         * 
         */
        public Builder count(@Nullable Output<Integer> count) {
            $.count = count;
            return this;
        }

        /**
         * @param count Number of occurrences in this series up to the last heartbeat time
         * 
         * @return builder
         * 
         */
        public Builder count(Integer count) {
            return count(Output.of(count));
        }

        /**
         * @param lastObservedTime Time when last Event from the series was seen before last heartbeat.
         * 
         * @return builder
         * 
         */
        public Builder lastObservedTime(@Nullable Output<String> lastObservedTime) {
            $.lastObservedTime = lastObservedTime;
            return this;
        }

        /**
         * @param lastObservedTime Time when last Event from the series was seen before last heartbeat.
         * 
         * @return builder
         * 
         */
        public Builder lastObservedTime(String lastObservedTime) {
            return lastObservedTime(Output.of(lastObservedTime));
        }

        /**
         * @param state Information whether this series is ongoing or finished. Deprecated. Planned removal for 1.18
         * 
         * @return builder
         * 
         */
        public Builder state(@Nullable Output<String> state) {
            $.state = state;
            return this;
        }

        /**
         * @param state Information whether this series is ongoing or finished. Deprecated. Planned removal for 1.18
         * 
         * @return builder
         * 
         */
        public Builder state(String state) {
            return state(Output.of(state));
        }

        public EventSeriesPatchArgs build() {
            return $;
        }
    }

}
