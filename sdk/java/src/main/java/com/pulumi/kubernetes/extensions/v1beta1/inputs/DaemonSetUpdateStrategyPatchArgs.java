// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.extensions.v1beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.extensions.v1beta1.inputs.RollingUpdateDaemonSetPatchArgs;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


public final class DaemonSetUpdateStrategyPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final DaemonSetUpdateStrategyPatchArgs Empty = new DaemonSetUpdateStrategyPatchArgs();

    /**
     * Rolling update config params. Present only if type = &#34;RollingUpdate&#34;.
     * 
     */
    @Import(name="rollingUpdate")
    private @Nullable Output<RollingUpdateDaemonSetPatchArgs> rollingUpdate;

    /**
     * @return Rolling update config params. Present only if type = &#34;RollingUpdate&#34;.
     * 
     */
    public Optional<Output<RollingUpdateDaemonSetPatchArgs>> rollingUpdate() {
        return Optional.ofNullable(this.rollingUpdate);
    }

    /**
     * Type of daemon set update. Can be &#34;RollingUpdate&#34; or &#34;OnDelete&#34;. Default is OnDelete.
     * 
     */
    @Import(name="type")
    private @Nullable Output<String> type;

    /**
     * @return Type of daemon set update. Can be &#34;RollingUpdate&#34; or &#34;OnDelete&#34;. Default is OnDelete.
     * 
     */
    public Optional<Output<String>> type() {
        return Optional.ofNullable(this.type);
    }

    private DaemonSetUpdateStrategyPatchArgs() {}

    private DaemonSetUpdateStrategyPatchArgs(DaemonSetUpdateStrategyPatchArgs $) {
        this.rollingUpdate = $.rollingUpdate;
        this.type = $.type;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(DaemonSetUpdateStrategyPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private DaemonSetUpdateStrategyPatchArgs $;

        public Builder() {
            $ = new DaemonSetUpdateStrategyPatchArgs();
        }

        public Builder(DaemonSetUpdateStrategyPatchArgs defaults) {
            $ = new DaemonSetUpdateStrategyPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param rollingUpdate Rolling update config params. Present only if type = &#34;RollingUpdate&#34;.
         * 
         * @return builder
         * 
         */
        public Builder rollingUpdate(@Nullable Output<RollingUpdateDaemonSetPatchArgs> rollingUpdate) {
            $.rollingUpdate = rollingUpdate;
            return this;
        }

        /**
         * @param rollingUpdate Rolling update config params. Present only if type = &#34;RollingUpdate&#34;.
         * 
         * @return builder
         * 
         */
        public Builder rollingUpdate(RollingUpdateDaemonSetPatchArgs rollingUpdate) {
            return rollingUpdate(Output.of(rollingUpdate));
        }

        /**
         * @param type Type of daemon set update. Can be &#34;RollingUpdate&#34; or &#34;OnDelete&#34;. Default is OnDelete.
         * 
         * @return builder
         * 
         */
        public Builder type(@Nullable Output<String> type) {
            $.type = type;
            return this;
        }

        /**
         * @param type Type of daemon set update. Can be &#34;RollingUpdate&#34; or &#34;OnDelete&#34;. Default is OnDelete.
         * 
         * @return builder
         * 
         */
        public Builder type(String type) {
            return type(Output.of(type));
        }

        public DaemonSetUpdateStrategyPatchArgs build() {
            return $;
        }
    }

}
