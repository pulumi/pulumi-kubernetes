// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.core.v1.outputs.NodeAffinity;
import com.pulumi.kubernetes.core.v1.outputs.PodAffinity;
import com.pulumi.kubernetes.core.v1.outputs.PodAntiAffinity;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class Affinity {
    /**
     * @return Describes node affinity scheduling rules for the pod.
     * 
     */
    private @Nullable NodeAffinity nodeAffinity;
    /**
     * @return Describes pod affinity scheduling rules (e.g. co-locate this pod in the same node, zone, etc. as some other pod(s)).
     * 
     */
    private @Nullable PodAffinity podAffinity;
    /**
     * @return Describes pod anti-affinity scheduling rules (e.g. avoid putting this pod in the same node, zone, etc. as some other pod(s)).
     * 
     */
    private @Nullable PodAntiAffinity podAntiAffinity;

    private Affinity() {}
    /**
     * @return Describes node affinity scheduling rules for the pod.
     * 
     */
    public Optional<NodeAffinity> nodeAffinity() {
        return Optional.ofNullable(this.nodeAffinity);
    }
    /**
     * @return Describes pod affinity scheduling rules (e.g. co-locate this pod in the same node, zone, etc. as some other pod(s)).
     * 
     */
    public Optional<PodAffinity> podAffinity() {
        return Optional.ofNullable(this.podAffinity);
    }
    /**
     * @return Describes pod anti-affinity scheduling rules (e.g. avoid putting this pod in the same node, zone, etc. as some other pod(s)).
     * 
     */
    public Optional<PodAntiAffinity> podAntiAffinity() {
        return Optional.ofNullable(this.podAntiAffinity);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(Affinity defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable NodeAffinity nodeAffinity;
        private @Nullable PodAffinity podAffinity;
        private @Nullable PodAntiAffinity podAntiAffinity;
        public Builder() {}
        public Builder(Affinity defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.nodeAffinity = defaults.nodeAffinity;
    	      this.podAffinity = defaults.podAffinity;
    	      this.podAntiAffinity = defaults.podAntiAffinity;
        }

        @CustomType.Setter
        public Builder nodeAffinity(@Nullable NodeAffinity nodeAffinity) {

            this.nodeAffinity = nodeAffinity;
            return this;
        }
        @CustomType.Setter
        public Builder podAffinity(@Nullable PodAffinity podAffinity) {

            this.podAffinity = podAffinity;
            return this;
        }
        @CustomType.Setter
        public Builder podAntiAffinity(@Nullable PodAntiAffinity podAntiAffinity) {

            this.podAntiAffinity = podAntiAffinity;
            return this;
        }
        public Affinity build() {
            final var _resultValue = new Affinity();
            _resultValue.nodeAffinity = nodeAffinity;
            _resultValue.podAffinity = podAffinity;
            _resultValue.podAntiAffinity = podAntiAffinity;
            return _resultValue;
        }
    }
}
