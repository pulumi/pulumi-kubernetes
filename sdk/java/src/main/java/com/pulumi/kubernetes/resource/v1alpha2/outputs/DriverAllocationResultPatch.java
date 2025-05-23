// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1alpha2.outputs;

import com.google.gson.JsonElement;
import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.resource.v1alpha2.outputs.NamedResourcesAllocationResultPatch;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class DriverAllocationResultPatch {
    /**
     * @return NamedResources describes the allocation result when using the named resources model.
     * 
     */
    private @Nullable NamedResourcesAllocationResultPatch namedResources;
    /**
     * @return VendorRequestParameters are the per-request configuration parameters from the time that the claim was allocated.
     * 
     */
    private @Nullable JsonElement vendorRequestParameters;

    private DriverAllocationResultPatch() {}
    /**
     * @return NamedResources describes the allocation result when using the named resources model.
     * 
     */
    public Optional<NamedResourcesAllocationResultPatch> namedResources() {
        return Optional.ofNullable(this.namedResources);
    }
    /**
     * @return VendorRequestParameters are the per-request configuration parameters from the time that the claim was allocated.
     * 
     */
    public Optional<JsonElement> vendorRequestParameters() {
        return Optional.ofNullable(this.vendorRequestParameters);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(DriverAllocationResultPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable NamedResourcesAllocationResultPatch namedResources;
        private @Nullable JsonElement vendorRequestParameters;
        public Builder() {}
        public Builder(DriverAllocationResultPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.namedResources = defaults.namedResources;
    	      this.vendorRequestParameters = defaults.vendorRequestParameters;
        }

        @CustomType.Setter
        public Builder namedResources(@Nullable NamedResourcesAllocationResultPatch namedResources) {

            this.namedResources = namedResources;
            return this;
        }
        @CustomType.Setter
        public Builder vendorRequestParameters(@Nullable JsonElement vendorRequestParameters) {

            this.vendorRequestParameters = vendorRequestParameters;
            return this;
        }
        public DriverAllocationResultPatch build() {
            final var _resultValue = new DriverAllocationResultPatch();
            _resultValue.namedResources = namedResources;
            _resultValue.vendorRequestParameters = vendorRequestParameters;
            return _resultValue;
        }
    }
}
