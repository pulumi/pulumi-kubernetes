// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1alpha2.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.String;
import java.util.Objects;

@CustomType
public final class NamedResourcesAllocationResult {
    /**
     * @return Name is the name of the selected resource instance.
     * 
     */
    private String name;

    private NamedResourcesAllocationResult() {}
    /**
     * @return Name is the name of the selected resource instance.
     * 
     */
    public String name() {
        return this.name;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(NamedResourcesAllocationResult defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private String name;
        public Builder() {}
        public Builder(NamedResourcesAllocationResult defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.name = defaults.name;
        }

        @CustomType.Setter
        public Builder name(String name) {
            if (name == null) {
              throw new MissingRequiredPropertyException("NamedResourcesAllocationResult", "name");
            }
            this.name = name;
            return this;
        }
        public NamedResourcesAllocationResult build() {
            final var _resultValue = new NamedResourcesAllocationResult();
            _resultValue.name = name;
            return _resultValue;
        }
    }
}
