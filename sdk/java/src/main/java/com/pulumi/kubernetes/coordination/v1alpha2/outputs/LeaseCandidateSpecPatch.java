// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.coordination.v1alpha2.outputs;

import com.pulumi.core.annotations.CustomType;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class LeaseCandidateSpecPatch {
    /**
     * @return BinaryVersion is the binary version. It must be in a semver format without leading `v`. This field is required.
     * 
     */
    private @Nullable String binaryVersion;
    /**
     * @return EmulationVersion is the emulation version. It must be in a semver format without leading `v`. EmulationVersion must be less than or equal to BinaryVersion. This field is required when strategy is &#34;OldestEmulationVersion&#34;
     * 
     */
    private @Nullable String emulationVersion;
    /**
     * @return LeaseName is the name of the lease for which this candidate is contending. This field is immutable.
     * 
     */
    private @Nullable String leaseName;
    /**
     * @return PingTime is the last time that the server has requested the LeaseCandidate to renew. It is only done during leader election to check if any LeaseCandidates have become ineligible. When PingTime is updated, the LeaseCandidate will respond by updating RenewTime.
     * 
     */
    private @Nullable String pingTime;
    /**
     * @return RenewTime is the time that the LeaseCandidate was last updated. Any time a Lease needs to do leader election, the PingTime field is updated to signal to the LeaseCandidate that they should update the RenewTime. Old LeaseCandidate objects are also garbage collected if it has been hours since the last renew. The PingTime field is updated regularly to prevent garbage collection for still active LeaseCandidates.
     * 
     */
    private @Nullable String renewTime;
    /**
     * @return Strategy is the strategy that coordinated leader election will use for picking the leader. If multiple candidates for the same Lease return different strategies, the strategy provided by the candidate with the latest BinaryVersion will be used. If there is still conflict, this is a user error and coordinated leader election will not operate the Lease until resolved.
     * 
     */
    private @Nullable String strategy;

    private LeaseCandidateSpecPatch() {}
    /**
     * @return BinaryVersion is the binary version. It must be in a semver format without leading `v`. This field is required.
     * 
     */
    public Optional<String> binaryVersion() {
        return Optional.ofNullable(this.binaryVersion);
    }
    /**
     * @return EmulationVersion is the emulation version. It must be in a semver format without leading `v`. EmulationVersion must be less than or equal to BinaryVersion. This field is required when strategy is &#34;OldestEmulationVersion&#34;
     * 
     */
    public Optional<String> emulationVersion() {
        return Optional.ofNullable(this.emulationVersion);
    }
    /**
     * @return LeaseName is the name of the lease for which this candidate is contending. This field is immutable.
     * 
     */
    public Optional<String> leaseName() {
        return Optional.ofNullable(this.leaseName);
    }
    /**
     * @return PingTime is the last time that the server has requested the LeaseCandidate to renew. It is only done during leader election to check if any LeaseCandidates have become ineligible. When PingTime is updated, the LeaseCandidate will respond by updating RenewTime.
     * 
     */
    public Optional<String> pingTime() {
        return Optional.ofNullable(this.pingTime);
    }
    /**
     * @return RenewTime is the time that the LeaseCandidate was last updated. Any time a Lease needs to do leader election, the PingTime field is updated to signal to the LeaseCandidate that they should update the RenewTime. Old LeaseCandidate objects are also garbage collected if it has been hours since the last renew. The PingTime field is updated regularly to prevent garbage collection for still active LeaseCandidates.
     * 
     */
    public Optional<String> renewTime() {
        return Optional.ofNullable(this.renewTime);
    }
    /**
     * @return Strategy is the strategy that coordinated leader election will use for picking the leader. If multiple candidates for the same Lease return different strategies, the strategy provided by the candidate with the latest BinaryVersion will be used. If there is still conflict, this is a user error and coordinated leader election will not operate the Lease until resolved.
     * 
     */
    public Optional<String> strategy() {
        return Optional.ofNullable(this.strategy);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(LeaseCandidateSpecPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable String binaryVersion;
        private @Nullable String emulationVersion;
        private @Nullable String leaseName;
        private @Nullable String pingTime;
        private @Nullable String renewTime;
        private @Nullable String strategy;
        public Builder() {}
        public Builder(LeaseCandidateSpecPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.binaryVersion = defaults.binaryVersion;
    	      this.emulationVersion = defaults.emulationVersion;
    	      this.leaseName = defaults.leaseName;
    	      this.pingTime = defaults.pingTime;
    	      this.renewTime = defaults.renewTime;
    	      this.strategy = defaults.strategy;
        }

        @CustomType.Setter
        public Builder binaryVersion(@Nullable String binaryVersion) {

            this.binaryVersion = binaryVersion;
            return this;
        }
        @CustomType.Setter
        public Builder emulationVersion(@Nullable String emulationVersion) {

            this.emulationVersion = emulationVersion;
            return this;
        }
        @CustomType.Setter
        public Builder leaseName(@Nullable String leaseName) {

            this.leaseName = leaseName;
            return this;
        }
        @CustomType.Setter
        public Builder pingTime(@Nullable String pingTime) {

            this.pingTime = pingTime;
            return this;
        }
        @CustomType.Setter
        public Builder renewTime(@Nullable String renewTime) {

            this.renewTime = renewTime;
            return this;
        }
        @CustomType.Setter
        public Builder strategy(@Nullable String strategy) {

            this.strategy = strategy;
            return this;
        }
        public LeaseCandidateSpecPatch build() {
            final var _resultValue = new LeaseCandidateSpecPatch();
            _resultValue.binaryVersion = binaryVersion;
            _resultValue.emulationVersion = emulationVersion;
            _resultValue.leaseName = leaseName;
            _resultValue.pingTime = pingTime;
            _resultValue.renewTime = renewTime;
            _resultValue.strategy = strategy;
            return _resultValue;
        }
    }
}
