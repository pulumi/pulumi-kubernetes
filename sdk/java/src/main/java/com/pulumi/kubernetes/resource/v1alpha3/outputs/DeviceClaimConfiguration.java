// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1alpha3.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.resource.v1alpha3.outputs.OpaqueDeviceConfiguration;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class DeviceClaimConfiguration {
    /**
     * @return Opaque provides driver-specific configuration parameters.
     * 
     */
    private @Nullable OpaqueDeviceConfiguration opaque;
    /**
     * @return Requests lists the names of requests where the configuration applies. If empty, it applies to all requests.
     * 
     * References to subrequests must include the name of the main request and may include the subrequest using the format &lt;main request&gt;[/&lt;subrequest&gt;]. If just the main request is given, the configuration applies to all subrequests.
     * 
     */
    private @Nullable List<String> requests;

    private DeviceClaimConfiguration() {}
    /**
     * @return Opaque provides driver-specific configuration parameters.
     * 
     */
    public Optional<OpaqueDeviceConfiguration> opaque() {
        return Optional.ofNullable(this.opaque);
    }
    /**
     * @return Requests lists the names of requests where the configuration applies. If empty, it applies to all requests.
     * 
     * References to subrequests must include the name of the main request and may include the subrequest using the format &lt;main request&gt;[/&lt;subrequest&gt;]. If just the main request is given, the configuration applies to all subrequests.
     * 
     */
    public List<String> requests() {
        return this.requests == null ? List.of() : this.requests;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(DeviceClaimConfiguration defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable OpaqueDeviceConfiguration opaque;
        private @Nullable List<String> requests;
        public Builder() {}
        public Builder(DeviceClaimConfiguration defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.opaque = defaults.opaque;
    	      this.requests = defaults.requests;
        }

        @CustomType.Setter
        public Builder opaque(@Nullable OpaqueDeviceConfiguration opaque) {

            this.opaque = opaque;
            return this;
        }
        @CustomType.Setter
        public Builder requests(@Nullable List<String> requests) {

            this.requests = requests;
            return this;
        }
        public Builder requests(String... requests) {
            return requests(List.of(requests));
        }
        public DeviceClaimConfiguration build() {
            final var _resultValue = new DeviceClaimConfiguration();
            _resultValue.opaque = opaque;
            _resultValue.requests = requests;
            return _resultValue;
        }
    }
}
