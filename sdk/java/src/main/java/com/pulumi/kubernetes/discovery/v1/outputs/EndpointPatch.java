// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.discovery.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.core.v1.outputs.ObjectReferencePatch;
import com.pulumi.kubernetes.discovery.v1.outputs.EndpointConditionsPatch;
import com.pulumi.kubernetes.discovery.v1.outputs.EndpointHintsPatch;
import java.lang.String;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class EndpointPatch {
    /**
     * @return addresses of this endpoint. For EndpointSlices of addressType &#34;IPv4&#34; or &#34;IPv6&#34;, the values are IP addresses in canonical form. The syntax and semantics of other addressType values are not defined. This must contain at least one address but no more than 100. EndpointSlices generated by the EndpointSlice controller will always have exactly 1 address. No semantics are defined for additional addresses beyond the first, and kube-proxy does not look at them.
     * 
     */
    private @Nullable List<String> addresses;
    /**
     * @return conditions contains information about the current status of the endpoint.
     * 
     */
    private @Nullable EndpointConditionsPatch conditions;
    /**
     * @return deprecatedTopology contains topology information part of the v1beta1 API. This field is deprecated, and will be removed when the v1beta1 API is removed (no sooner than kubernetes v1.24).  While this field can hold values, it is not writable through the v1 API, and any attempts to write to it will be silently ignored. Topology information can be found in the zone and nodeName fields instead.
     * 
     */
    private @Nullable Map<String,String> deprecatedTopology;
    /**
     * @return hints contains information associated with how an endpoint should be consumed.
     * 
     */
    private @Nullable EndpointHintsPatch hints;
    /**
     * @return hostname of this endpoint. This field may be used by consumers of endpoints to distinguish endpoints from each other (e.g. in DNS names). Multiple endpoints which use the same hostname should be considered fungible (e.g. multiple A values in DNS). Must be lowercase and pass DNS Label (RFC 1123) validation.
     * 
     */
    private @Nullable String hostname;
    /**
     * @return nodeName represents the name of the Node hosting this endpoint. This can be used to determine endpoints local to a Node.
     * 
     */
    private @Nullable String nodeName;
    /**
     * @return targetRef is a reference to a Kubernetes object that represents this endpoint.
     * 
     */
    private @Nullable ObjectReferencePatch targetRef;
    /**
     * @return zone is the name of the Zone this endpoint exists in.
     * 
     */
    private @Nullable String zone;

    private EndpointPatch() {}
    /**
     * @return addresses of this endpoint. For EndpointSlices of addressType &#34;IPv4&#34; or &#34;IPv6&#34;, the values are IP addresses in canonical form. The syntax and semantics of other addressType values are not defined. This must contain at least one address but no more than 100. EndpointSlices generated by the EndpointSlice controller will always have exactly 1 address. No semantics are defined for additional addresses beyond the first, and kube-proxy does not look at them.
     * 
     */
    public List<String> addresses() {
        return this.addresses == null ? List.of() : this.addresses;
    }
    /**
     * @return conditions contains information about the current status of the endpoint.
     * 
     */
    public Optional<EndpointConditionsPatch> conditions() {
        return Optional.ofNullable(this.conditions);
    }
    /**
     * @return deprecatedTopology contains topology information part of the v1beta1 API. This field is deprecated, and will be removed when the v1beta1 API is removed (no sooner than kubernetes v1.24).  While this field can hold values, it is not writable through the v1 API, and any attempts to write to it will be silently ignored. Topology information can be found in the zone and nodeName fields instead.
     * 
     */
    public Map<String,String> deprecatedTopology() {
        return this.deprecatedTopology == null ? Map.of() : this.deprecatedTopology;
    }
    /**
     * @return hints contains information associated with how an endpoint should be consumed.
     * 
     */
    public Optional<EndpointHintsPatch> hints() {
        return Optional.ofNullable(this.hints);
    }
    /**
     * @return hostname of this endpoint. This field may be used by consumers of endpoints to distinguish endpoints from each other (e.g. in DNS names). Multiple endpoints which use the same hostname should be considered fungible (e.g. multiple A values in DNS). Must be lowercase and pass DNS Label (RFC 1123) validation.
     * 
     */
    public Optional<String> hostname() {
        return Optional.ofNullable(this.hostname);
    }
    /**
     * @return nodeName represents the name of the Node hosting this endpoint. This can be used to determine endpoints local to a Node.
     * 
     */
    public Optional<String> nodeName() {
        return Optional.ofNullable(this.nodeName);
    }
    /**
     * @return targetRef is a reference to a Kubernetes object that represents this endpoint.
     * 
     */
    public Optional<ObjectReferencePatch> targetRef() {
        return Optional.ofNullable(this.targetRef);
    }
    /**
     * @return zone is the name of the Zone this endpoint exists in.
     * 
     */
    public Optional<String> zone() {
        return Optional.ofNullable(this.zone);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(EndpointPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable List<String> addresses;
        private @Nullable EndpointConditionsPatch conditions;
        private @Nullable Map<String,String> deprecatedTopology;
        private @Nullable EndpointHintsPatch hints;
        private @Nullable String hostname;
        private @Nullable String nodeName;
        private @Nullable ObjectReferencePatch targetRef;
        private @Nullable String zone;
        public Builder() {}
        public Builder(EndpointPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.addresses = defaults.addresses;
    	      this.conditions = defaults.conditions;
    	      this.deprecatedTopology = defaults.deprecatedTopology;
    	      this.hints = defaults.hints;
    	      this.hostname = defaults.hostname;
    	      this.nodeName = defaults.nodeName;
    	      this.targetRef = defaults.targetRef;
    	      this.zone = defaults.zone;
        }

        @CustomType.Setter
        public Builder addresses(@Nullable List<String> addresses) {

            this.addresses = addresses;
            return this;
        }
        public Builder addresses(String... addresses) {
            return addresses(List.of(addresses));
        }
        @CustomType.Setter
        public Builder conditions(@Nullable EndpointConditionsPatch conditions) {

            this.conditions = conditions;
            return this;
        }
        @CustomType.Setter
        public Builder deprecatedTopology(@Nullable Map<String,String> deprecatedTopology) {

            this.deprecatedTopology = deprecatedTopology;
            return this;
        }
        @CustomType.Setter
        public Builder hints(@Nullable EndpointHintsPatch hints) {

            this.hints = hints;
            return this;
        }
        @CustomType.Setter
        public Builder hostname(@Nullable String hostname) {

            this.hostname = hostname;
            return this;
        }
        @CustomType.Setter
        public Builder nodeName(@Nullable String nodeName) {

            this.nodeName = nodeName;
            return this;
        }
        @CustomType.Setter
        public Builder targetRef(@Nullable ObjectReferencePatch targetRef) {

            this.targetRef = targetRef;
            return this;
        }
        @CustomType.Setter
        public Builder zone(@Nullable String zone) {

            this.zone = zone;
            return this;
        }
        public EndpointPatch build() {
            final var _resultValue = new EndpointPatch();
            _resultValue.addresses = addresses;
            _resultValue.conditions = conditions;
            _resultValue.deprecatedTopology = deprecatedTopology;
            _resultValue.hints = hints;
            _resultValue.hostname = hostname;
            _resultValue.nodeName = nodeName;
            _resultValue.targetRef = targetRef;
            _resultValue.zone = zone;
            return _resultValue;
        }
    }
}
