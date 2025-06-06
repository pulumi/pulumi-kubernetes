// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.apps.v1beta2.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.apps.v1beta2.outputs.RollingUpdateDaemonSetPatch;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class DaemonSetUpdateStrategyPatch {
    /**
     * @return Rolling update config params. Present only if type = &#34;RollingUpdate&#34;.
     * 
     */
    private @Nullable RollingUpdateDaemonSetPatch rollingUpdate;
    /**
     * @return Type of daemon set update. Can be &#34;RollingUpdate&#34; or &#34;OnDelete&#34;. Default is RollingUpdate.
     * 
     */
    private @Nullable String type;

    private DaemonSetUpdateStrategyPatch() {}
    /**
     * @return Rolling update config params. Present only if type = &#34;RollingUpdate&#34;.
     * 
     */
    public Optional<RollingUpdateDaemonSetPatch> rollingUpdate() {
        return Optional.ofNullable(this.rollingUpdate);
    }
    /**
     * @return Type of daemon set update. Can be &#34;RollingUpdate&#34; or &#34;OnDelete&#34;. Default is RollingUpdate.
     * 
     */
    public Optional<String> type() {
        return Optional.ofNullable(this.type);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(DaemonSetUpdateStrategyPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable RollingUpdateDaemonSetPatch rollingUpdate;
        private @Nullable String type;
        public Builder() {}
        public Builder(DaemonSetUpdateStrategyPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.rollingUpdate = defaults.rollingUpdate;
    	      this.type = defaults.type;
        }

        @CustomType.Setter
        public Builder rollingUpdate(@Nullable RollingUpdateDaemonSetPatch rollingUpdate) {

            this.rollingUpdate = rollingUpdate;
            return this;
        }
        @CustomType.Setter
        public Builder type(@Nullable String type) {

            this.type = type;
            return this;
        }
        public DaemonSetUpdateStrategyPatch build() {
            final var _resultValue = new DaemonSetUpdateStrategyPatch();
            _resultValue.rollingUpdate = rollingUpdate;
            _resultValue.type = type;
            return _resultValue;
        }
    }
}
