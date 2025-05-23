// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.core.v1.inputs.DownwardAPIVolumeFileArgs;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * Represents downward API info for projecting into a projected volume. Note that this is identical to a downwardAPI volume source without the default mode.
 * 
 */
public final class DownwardAPIProjectionArgs extends com.pulumi.resources.ResourceArgs {

    public static final DownwardAPIProjectionArgs Empty = new DownwardAPIProjectionArgs();

    /**
     * Items is a list of DownwardAPIVolume file
     * 
     */
    @Import(name="items")
    private @Nullable Output<List<DownwardAPIVolumeFileArgs>> items;

    /**
     * @return Items is a list of DownwardAPIVolume file
     * 
     */
    public Optional<Output<List<DownwardAPIVolumeFileArgs>>> items() {
        return Optional.ofNullable(this.items);
    }

    private DownwardAPIProjectionArgs() {}

    private DownwardAPIProjectionArgs(DownwardAPIProjectionArgs $) {
        this.items = $.items;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(DownwardAPIProjectionArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private DownwardAPIProjectionArgs $;

        public Builder() {
            $ = new DownwardAPIProjectionArgs();
        }

        public Builder(DownwardAPIProjectionArgs defaults) {
            $ = new DownwardAPIProjectionArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param items Items is a list of DownwardAPIVolume file
         * 
         * @return builder
         * 
         */
        public Builder items(@Nullable Output<List<DownwardAPIVolumeFileArgs>> items) {
            $.items = items;
            return this;
        }

        /**
         * @param items Items is a list of DownwardAPIVolume file
         * 
         * @return builder
         * 
         */
        public Builder items(List<DownwardAPIVolumeFileArgs> items) {
            return items(Output.of(items));
        }

        /**
         * @param items Items is a list of DownwardAPIVolume file
         * 
         * @return builder
         * 
         */
        public Builder items(DownwardAPIVolumeFileArgs... items) {
            return items(List.of(items));
        }

        public DownwardAPIProjectionArgs build() {
            return $;
        }
    }

}
