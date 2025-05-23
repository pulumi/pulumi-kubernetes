// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.autoscaling.v2.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.autoscaling.v2.outputs.MetricValueStatusPatch;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class ResourceMetricStatusPatch {
    /**
     * @return current contains the current value for the given metric
     * 
     */
    private @Nullable MetricValueStatusPatch current;
    /**
     * @return name is the name of the resource in question.
     * 
     */
    private @Nullable String name;

    private ResourceMetricStatusPatch() {}
    /**
     * @return current contains the current value for the given metric
     * 
     */
    public Optional<MetricValueStatusPatch> current() {
        return Optional.ofNullable(this.current);
    }
    /**
     * @return name is the name of the resource in question.
     * 
     */
    public Optional<String> name() {
        return Optional.ofNullable(this.name);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ResourceMetricStatusPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable MetricValueStatusPatch current;
        private @Nullable String name;
        public Builder() {}
        public Builder(ResourceMetricStatusPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.current = defaults.current;
    	      this.name = defaults.name;
        }

        @CustomType.Setter
        public Builder current(@Nullable MetricValueStatusPatch current) {

            this.current = current;
            return this;
        }
        @CustomType.Setter
        public Builder name(@Nullable String name) {

            this.name = name;
            return this;
        }
        public ResourceMetricStatusPatch build() {
            final var _resultValue = new ResourceMetricStatusPatch();
            _resultValue.current = current;
            _resultValue.name = name;
            return _resultValue;
        }
    }
}
