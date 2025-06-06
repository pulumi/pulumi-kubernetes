// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.autoscaling.v2.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.autoscaling.v2.outputs.MetricValueStatus;
import java.lang.String;
import java.util.Objects;

@CustomType
public final class ResourceMetricStatus {
    /**
     * @return current contains the current value for the given metric
     * 
     */
    private MetricValueStatus current;
    /**
     * @return name is the name of the resource in question.
     * 
     */
    private String name;

    private ResourceMetricStatus() {}
    /**
     * @return current contains the current value for the given metric
     * 
     */
    public MetricValueStatus current() {
        return this.current;
    }
    /**
     * @return name is the name of the resource in question.
     * 
     */
    public String name() {
        return this.name;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ResourceMetricStatus defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private MetricValueStatus current;
        private String name;
        public Builder() {}
        public Builder(ResourceMetricStatus defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.current = defaults.current;
    	      this.name = defaults.name;
        }

        @CustomType.Setter
        public Builder current(MetricValueStatus current) {
            if (current == null) {
              throw new MissingRequiredPropertyException("ResourceMetricStatus", "current");
            }
            this.current = current;
            return this;
        }
        @CustomType.Setter
        public Builder name(String name) {
            if (name == null) {
              throw new MissingRequiredPropertyException("ResourceMetricStatus", "name");
            }
            this.name = name;
            return this;
        }
        public ResourceMetricStatus build() {
            final var _resultValue = new ResourceMetricStatus();
            _resultValue.current = current;
            _resultValue.name = name;
            return _resultValue;
        }
    }
}
