// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class NodeAddressPatch {
    /**
     * @return The node address.
     * 
     */
    private @Nullable String address;
    /**
     * @return Node address type, one of Hostname, ExternalIP or InternalIP.
     * 
     */
    private @Nullable String type;

    private NodeAddressPatch() {}
    /**
     * @return The node address.
     * 
     */
    public Optional<String> address() {
        return Optional.ofNullable(this.address);
    }
    /**
     * @return Node address type, one of Hostname, ExternalIP or InternalIP.
     * 
     */
    public Optional<String> type() {
        return Optional.ofNullable(this.type);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(NodeAddressPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable String address;
        private @Nullable String type;
        public Builder() {}
        public Builder(NodeAddressPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.address = defaults.address;
    	      this.type = defaults.type;
        }

        @CustomType.Setter
        public Builder address(@Nullable String address) {

            this.address = address;
            return this;
        }
        @CustomType.Setter
        public Builder type(@Nullable String type) {

            this.type = type;
            return this;
        }
        public NodeAddressPatch build() {
            final var _resultValue = new NodeAddressPatch();
            _resultValue.address = address;
            _resultValue.type = type;
            return _resultValue;
        }
    }
}
