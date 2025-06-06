// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1alpha2.outputs;

import com.google.gson.JsonElement;
import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.resource.v1alpha2.outputs.DriverAllocationResult;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class StructuredResourceHandle {
    /**
     * @return NodeName is the name of the node providing the necessary resources if the resources are local to a node.
     * 
     */
    private @Nullable String nodeName;
    /**
     * @return Results lists all allocated driver resources.
     * 
     */
    private List<DriverAllocationResult> results;
    /**
     * @return VendorClaimParameters are the per-claim configuration parameters from the resource claim parameters at the time that the claim was allocated.
     * 
     */
    private @Nullable JsonElement vendorClaimParameters;
    /**
     * @return VendorClassParameters are the per-claim configuration parameters from the resource class at the time that the claim was allocated.
     * 
     */
    private @Nullable JsonElement vendorClassParameters;

    private StructuredResourceHandle() {}
    /**
     * @return NodeName is the name of the node providing the necessary resources if the resources are local to a node.
     * 
     */
    public Optional<String> nodeName() {
        return Optional.ofNullable(this.nodeName);
    }
    /**
     * @return Results lists all allocated driver resources.
     * 
     */
    public List<DriverAllocationResult> results() {
        return this.results;
    }
    /**
     * @return VendorClaimParameters are the per-claim configuration parameters from the resource claim parameters at the time that the claim was allocated.
     * 
     */
    public Optional<JsonElement> vendorClaimParameters() {
        return Optional.ofNullable(this.vendorClaimParameters);
    }
    /**
     * @return VendorClassParameters are the per-claim configuration parameters from the resource class at the time that the claim was allocated.
     * 
     */
    public Optional<JsonElement> vendorClassParameters() {
        return Optional.ofNullable(this.vendorClassParameters);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(StructuredResourceHandle defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable String nodeName;
        private List<DriverAllocationResult> results;
        private @Nullable JsonElement vendorClaimParameters;
        private @Nullable JsonElement vendorClassParameters;
        public Builder() {}
        public Builder(StructuredResourceHandle defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.nodeName = defaults.nodeName;
    	      this.results = defaults.results;
    	      this.vendorClaimParameters = defaults.vendorClaimParameters;
    	      this.vendorClassParameters = defaults.vendorClassParameters;
        }

        @CustomType.Setter
        public Builder nodeName(@Nullable String nodeName) {

            this.nodeName = nodeName;
            return this;
        }
        @CustomType.Setter
        public Builder results(List<DriverAllocationResult> results) {
            if (results == null) {
              throw new MissingRequiredPropertyException("StructuredResourceHandle", "results");
            }
            this.results = results;
            return this;
        }
        public Builder results(DriverAllocationResult... results) {
            return results(List.of(results));
        }
        @CustomType.Setter
        public Builder vendorClaimParameters(@Nullable JsonElement vendorClaimParameters) {

            this.vendorClaimParameters = vendorClaimParameters;
            return this;
        }
        @CustomType.Setter
        public Builder vendorClassParameters(@Nullable JsonElement vendorClassParameters) {

            this.vendorClassParameters = vendorClassParameters;
            return this;
        }
        public StructuredResourceHandle build() {
            final var _resultValue = new StructuredResourceHandle();
            _resultValue.nodeName = nodeName;
            _resultValue.results = results;
            _resultValue.vendorClaimParameters = vendorClaimParameters;
            _resultValue.vendorClassParameters = vendorClassParameters;
            return _resultValue;
        }
    }
}
