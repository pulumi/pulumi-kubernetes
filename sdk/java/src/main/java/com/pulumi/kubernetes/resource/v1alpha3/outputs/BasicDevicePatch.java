// *** WARNING: this file was generated by pulumi-java-gen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1alpha3.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.resource.v1alpha3.outputs.DeviceAttribute;
import java.lang.String;
import java.util.Map;
import java.util.Objects;
import javax.annotation.Nullable;

@CustomType
public final class BasicDevicePatch {
    /**
     * @return Attributes defines the set of attributes for this device. The name of each attribute must be unique in that set.
     * 
     * The maximum number of attributes and capacities combined is 32.
     * 
     */
    private @Nullable Map<String,DeviceAttribute> attributes;
    /**
     * @return Capacity defines the set of capacities for this device. The name of each capacity must be unique in that set.
     * 
     * The maximum number of attributes and capacities combined is 32.
     * 
     */
    private @Nullable Map<String,String> capacity;

    private BasicDevicePatch() {}
    /**
     * @return Attributes defines the set of attributes for this device. The name of each attribute must be unique in that set.
     * 
     * The maximum number of attributes and capacities combined is 32.
     * 
     */
    public Map<String,DeviceAttribute> attributes() {
        return this.attributes == null ? Map.of() : this.attributes;
    }
    /**
     * @return Capacity defines the set of capacities for this device. The name of each capacity must be unique in that set.
     * 
     * The maximum number of attributes and capacities combined is 32.
     * 
     */
    public Map<String,String> capacity() {
        return this.capacity == null ? Map.of() : this.capacity;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(BasicDevicePatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable Map<String,DeviceAttribute> attributes;
        private @Nullable Map<String,String> capacity;
        public Builder() {}
        public Builder(BasicDevicePatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.attributes = defaults.attributes;
    	      this.capacity = defaults.capacity;
        }

        @CustomType.Setter
        public Builder attributes(@Nullable Map<String,DeviceAttribute> attributes) {

            this.attributes = attributes;
            return this;
        }
        @CustomType.Setter
        public Builder capacity(@Nullable Map<String,String> capacity) {

            this.capacity = capacity;
            return this;
        }
        public BasicDevicePatch build() {
            final var _resultValue = new BasicDevicePatch();
            _resultValue.attributes = attributes;
            _resultValue.capacity = capacity;
            return _resultValue;
        }
    }
}