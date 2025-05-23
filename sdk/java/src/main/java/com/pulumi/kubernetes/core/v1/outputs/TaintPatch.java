// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class TaintPatch {
    /**
     * @return Required. The effect of the taint on pods that do not tolerate the taint. Valid effects are NoSchedule, PreferNoSchedule and NoExecute.
     * 
     */
    private @Nullable String effect;
    /**
     * @return Required. The taint key to be applied to a node.
     * 
     */
    private @Nullable String key;
    /**
     * @return TimeAdded represents the time at which the taint was added. It is only written for NoExecute taints.
     * 
     */
    private @Nullable String timeAdded;
    /**
     * @return The taint value corresponding to the taint key.
     * 
     */
    private @Nullable String value;

    private TaintPatch() {}
    /**
     * @return Required. The effect of the taint on pods that do not tolerate the taint. Valid effects are NoSchedule, PreferNoSchedule and NoExecute.
     * 
     */
    public Optional<String> effect() {
        return Optional.ofNullable(this.effect);
    }
    /**
     * @return Required. The taint key to be applied to a node.
     * 
     */
    public Optional<String> key() {
        return Optional.ofNullable(this.key);
    }
    /**
     * @return TimeAdded represents the time at which the taint was added. It is only written for NoExecute taints.
     * 
     */
    public Optional<String> timeAdded() {
        return Optional.ofNullable(this.timeAdded);
    }
    /**
     * @return The taint value corresponding to the taint key.
     * 
     */
    public Optional<String> value() {
        return Optional.ofNullable(this.value);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(TaintPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable String effect;
        private @Nullable String key;
        private @Nullable String timeAdded;
        private @Nullable String value;
        public Builder() {}
        public Builder(TaintPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.effect = defaults.effect;
    	      this.key = defaults.key;
    	      this.timeAdded = defaults.timeAdded;
    	      this.value = defaults.value;
        }

        @CustomType.Setter
        public Builder effect(@Nullable String effect) {

            this.effect = effect;
            return this;
        }
        @CustomType.Setter
        public Builder key(@Nullable String key) {

            this.key = key;
            return this;
        }
        @CustomType.Setter
        public Builder timeAdded(@Nullable String timeAdded) {

            this.timeAdded = timeAdded;
            return this;
        }
        @CustomType.Setter
        public Builder value(@Nullable String value) {

            this.value = value;
            return this;
        }
        public TaintPatch build() {
            final var _resultValue = new TaintPatch();
            _resultValue.effect = effect;
            _resultValue.key = key;
            _resultValue.timeAdded = timeAdded;
            _resultValue.value = value;
            return _resultValue;
        }
    }
}
