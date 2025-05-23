// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.coordination.v1alpha1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * LeaseCandidateSpec is a specification of a Lease.
 * 
 */
public final class LeaseCandidateSpecPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final LeaseCandidateSpecPatchArgs Empty = new LeaseCandidateSpecPatchArgs();

    /**
     * BinaryVersion is the binary version. It must be in a semver format without leading `v`. This field is required when strategy is &#34;OldestEmulationVersion&#34;
     * 
     */
    @Import(name="binaryVersion")
    private @Nullable Output<String> binaryVersion;

    /**
     * @return BinaryVersion is the binary version. It must be in a semver format without leading `v`. This field is required when strategy is &#34;OldestEmulationVersion&#34;
     * 
     */
    public Optional<Output<String>> binaryVersion() {
        return Optional.ofNullable(this.binaryVersion);
    }

    /**
     * EmulationVersion is the emulation version. It must be in a semver format without leading `v`. EmulationVersion must be less than or equal to BinaryVersion. This field is required when strategy is &#34;OldestEmulationVersion&#34;
     * 
     */
    @Import(name="emulationVersion")
    private @Nullable Output<String> emulationVersion;

    /**
     * @return EmulationVersion is the emulation version. It must be in a semver format without leading `v`. EmulationVersion must be less than or equal to BinaryVersion. This field is required when strategy is &#34;OldestEmulationVersion&#34;
     * 
     */
    public Optional<Output<String>> emulationVersion() {
        return Optional.ofNullable(this.emulationVersion);
    }

    /**
     * LeaseName is the name of the lease for which this candidate is contending. This field is immutable.
     * 
     */
    @Import(name="leaseName")
    private @Nullable Output<String> leaseName;

    /**
     * @return LeaseName is the name of the lease for which this candidate is contending. This field is immutable.
     * 
     */
    public Optional<Output<String>> leaseName() {
        return Optional.ofNullable(this.leaseName);
    }

    /**
     * PingTime is the last time that the server has requested the LeaseCandidate to renew. It is only done during leader election to check if any LeaseCandidates have become ineligible. When PingTime is updated, the LeaseCandidate will respond by updating RenewTime.
     * 
     */
    @Import(name="pingTime")
    private @Nullable Output<String> pingTime;

    /**
     * @return PingTime is the last time that the server has requested the LeaseCandidate to renew. It is only done during leader election to check if any LeaseCandidates have become ineligible. When PingTime is updated, the LeaseCandidate will respond by updating RenewTime.
     * 
     */
    public Optional<Output<String>> pingTime() {
        return Optional.ofNullable(this.pingTime);
    }

    /**
     * PreferredStrategies indicates the list of strategies for picking the leader for coordinated leader election. The list is ordered, and the first strategy supersedes all other strategies. The list is used by coordinated leader election to make a decision about the final election strategy. This follows as - If all clients have strategy X as the first element in this list, strategy X will be used. - If a candidate has strategy [X] and another candidate has strategy [Y, X], Y supersedes X and strategy Y
     *   will be used.
     * - If a candidate has strategy [X, Y] and another candidate has strategy [Y, X], this is a user error and leader
     *   election will not operate the Lease until resolved.
     *   (Alpha) Using this field requires the CoordinatedLeaderElection feature gate to be enabled.
     * 
     */
    @Import(name="preferredStrategies")
    private @Nullable Output<List<String>> preferredStrategies;

    /**
     * @return PreferredStrategies indicates the list of strategies for picking the leader for coordinated leader election. The list is ordered, and the first strategy supersedes all other strategies. The list is used by coordinated leader election to make a decision about the final election strategy. This follows as - If all clients have strategy X as the first element in this list, strategy X will be used. - If a candidate has strategy [X] and another candidate has strategy [Y, X], Y supersedes X and strategy Y
     *   will be used.
     * - If a candidate has strategy [X, Y] and another candidate has strategy [Y, X], this is a user error and leader
     *   election will not operate the Lease until resolved.
     *   (Alpha) Using this field requires the CoordinatedLeaderElection feature gate to be enabled.
     * 
     */
    public Optional<Output<List<String>>> preferredStrategies() {
        return Optional.ofNullable(this.preferredStrategies);
    }

    /**
     * RenewTime is the time that the LeaseCandidate was last updated. Any time a Lease needs to do leader election, the PingTime field is updated to signal to the LeaseCandidate that they should update the RenewTime. Old LeaseCandidate objects are also garbage collected if it has been hours since the last renew. The PingTime field is updated regularly to prevent garbage collection for still active LeaseCandidates.
     * 
     */
    @Import(name="renewTime")
    private @Nullable Output<String> renewTime;

    /**
     * @return RenewTime is the time that the LeaseCandidate was last updated. Any time a Lease needs to do leader election, the PingTime field is updated to signal to the LeaseCandidate that they should update the RenewTime. Old LeaseCandidate objects are also garbage collected if it has been hours since the last renew. The PingTime field is updated regularly to prevent garbage collection for still active LeaseCandidates.
     * 
     */
    public Optional<Output<String>> renewTime() {
        return Optional.ofNullable(this.renewTime);
    }

    private LeaseCandidateSpecPatchArgs() {}

    private LeaseCandidateSpecPatchArgs(LeaseCandidateSpecPatchArgs $) {
        this.binaryVersion = $.binaryVersion;
        this.emulationVersion = $.emulationVersion;
        this.leaseName = $.leaseName;
        this.pingTime = $.pingTime;
        this.preferredStrategies = $.preferredStrategies;
        this.renewTime = $.renewTime;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(LeaseCandidateSpecPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private LeaseCandidateSpecPatchArgs $;

        public Builder() {
            $ = new LeaseCandidateSpecPatchArgs();
        }

        public Builder(LeaseCandidateSpecPatchArgs defaults) {
            $ = new LeaseCandidateSpecPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param binaryVersion BinaryVersion is the binary version. It must be in a semver format without leading `v`. This field is required when strategy is &#34;OldestEmulationVersion&#34;
         * 
         * @return builder
         * 
         */
        public Builder binaryVersion(@Nullable Output<String> binaryVersion) {
            $.binaryVersion = binaryVersion;
            return this;
        }

        /**
         * @param binaryVersion BinaryVersion is the binary version. It must be in a semver format without leading `v`. This field is required when strategy is &#34;OldestEmulationVersion&#34;
         * 
         * @return builder
         * 
         */
        public Builder binaryVersion(String binaryVersion) {
            return binaryVersion(Output.of(binaryVersion));
        }

        /**
         * @param emulationVersion EmulationVersion is the emulation version. It must be in a semver format without leading `v`. EmulationVersion must be less than or equal to BinaryVersion. This field is required when strategy is &#34;OldestEmulationVersion&#34;
         * 
         * @return builder
         * 
         */
        public Builder emulationVersion(@Nullable Output<String> emulationVersion) {
            $.emulationVersion = emulationVersion;
            return this;
        }

        /**
         * @param emulationVersion EmulationVersion is the emulation version. It must be in a semver format without leading `v`. EmulationVersion must be less than or equal to BinaryVersion. This field is required when strategy is &#34;OldestEmulationVersion&#34;
         * 
         * @return builder
         * 
         */
        public Builder emulationVersion(String emulationVersion) {
            return emulationVersion(Output.of(emulationVersion));
        }

        /**
         * @param leaseName LeaseName is the name of the lease for which this candidate is contending. This field is immutable.
         * 
         * @return builder
         * 
         */
        public Builder leaseName(@Nullable Output<String> leaseName) {
            $.leaseName = leaseName;
            return this;
        }

        /**
         * @param leaseName LeaseName is the name of the lease for which this candidate is contending. This field is immutable.
         * 
         * @return builder
         * 
         */
        public Builder leaseName(String leaseName) {
            return leaseName(Output.of(leaseName));
        }

        /**
         * @param pingTime PingTime is the last time that the server has requested the LeaseCandidate to renew. It is only done during leader election to check if any LeaseCandidates have become ineligible. When PingTime is updated, the LeaseCandidate will respond by updating RenewTime.
         * 
         * @return builder
         * 
         */
        public Builder pingTime(@Nullable Output<String> pingTime) {
            $.pingTime = pingTime;
            return this;
        }

        /**
         * @param pingTime PingTime is the last time that the server has requested the LeaseCandidate to renew. It is only done during leader election to check if any LeaseCandidates have become ineligible. When PingTime is updated, the LeaseCandidate will respond by updating RenewTime.
         * 
         * @return builder
         * 
         */
        public Builder pingTime(String pingTime) {
            return pingTime(Output.of(pingTime));
        }

        /**
         * @param preferredStrategies PreferredStrategies indicates the list of strategies for picking the leader for coordinated leader election. The list is ordered, and the first strategy supersedes all other strategies. The list is used by coordinated leader election to make a decision about the final election strategy. This follows as - If all clients have strategy X as the first element in this list, strategy X will be used. - If a candidate has strategy [X] and another candidate has strategy [Y, X], Y supersedes X and strategy Y
         *   will be used.
         * - If a candidate has strategy [X, Y] and another candidate has strategy [Y, X], this is a user error and leader
         *   election will not operate the Lease until resolved.
         *   (Alpha) Using this field requires the CoordinatedLeaderElection feature gate to be enabled.
         * 
         * @return builder
         * 
         */
        public Builder preferredStrategies(@Nullable Output<List<String>> preferredStrategies) {
            $.preferredStrategies = preferredStrategies;
            return this;
        }

        /**
         * @param preferredStrategies PreferredStrategies indicates the list of strategies for picking the leader for coordinated leader election. The list is ordered, and the first strategy supersedes all other strategies. The list is used by coordinated leader election to make a decision about the final election strategy. This follows as - If all clients have strategy X as the first element in this list, strategy X will be used. - If a candidate has strategy [X] and another candidate has strategy [Y, X], Y supersedes X and strategy Y
         *   will be used.
         * - If a candidate has strategy [X, Y] and another candidate has strategy [Y, X], this is a user error and leader
         *   election will not operate the Lease until resolved.
         *   (Alpha) Using this field requires the CoordinatedLeaderElection feature gate to be enabled.
         * 
         * @return builder
         * 
         */
        public Builder preferredStrategies(List<String> preferredStrategies) {
            return preferredStrategies(Output.of(preferredStrategies));
        }

        /**
         * @param preferredStrategies PreferredStrategies indicates the list of strategies for picking the leader for coordinated leader election. The list is ordered, and the first strategy supersedes all other strategies. The list is used by coordinated leader election to make a decision about the final election strategy. This follows as - If all clients have strategy X as the first element in this list, strategy X will be used. - If a candidate has strategy [X] and another candidate has strategy [Y, X], Y supersedes X and strategy Y
         *   will be used.
         * - If a candidate has strategy [X, Y] and another candidate has strategy [Y, X], this is a user error and leader
         *   election will not operate the Lease until resolved.
         *   (Alpha) Using this field requires the CoordinatedLeaderElection feature gate to be enabled.
         * 
         * @return builder
         * 
         */
        public Builder preferredStrategies(String... preferredStrategies) {
            return preferredStrategies(List.of(preferredStrategies));
        }

        /**
         * @param renewTime RenewTime is the time that the LeaseCandidate was last updated. Any time a Lease needs to do leader election, the PingTime field is updated to signal to the LeaseCandidate that they should update the RenewTime. Old LeaseCandidate objects are also garbage collected if it has been hours since the last renew. The PingTime field is updated regularly to prevent garbage collection for still active LeaseCandidates.
         * 
         * @return builder
         * 
         */
        public Builder renewTime(@Nullable Output<String> renewTime) {
            $.renewTime = renewTime;
            return this;
        }

        /**
         * @param renewTime RenewTime is the time that the LeaseCandidate was last updated. Any time a Lease needs to do leader election, the PingTime field is updated to signal to the LeaseCandidate that they should update the RenewTime. Old LeaseCandidate objects are also garbage collected if it has been hours since the last renew. The PingTime field is updated regularly to prevent garbage collection for still active LeaseCandidates.
         * 
         * @return builder
         * 
         */
        public Builder renewTime(String renewTime) {
            return renewTime(Output.of(renewTime));
        }

        public LeaseCandidateSpecPatchArgs build() {
            return $;
        }
    }

}
