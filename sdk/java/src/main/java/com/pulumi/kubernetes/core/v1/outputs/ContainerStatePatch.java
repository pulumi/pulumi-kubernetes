// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.core.v1.outputs.ContainerStateRunningPatch;
import com.pulumi.kubernetes.core.v1.outputs.ContainerStateTerminatedPatch;
import com.pulumi.kubernetes.core.v1.outputs.ContainerStateWaitingPatch;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class ContainerStatePatch {
    /**
     * @return Details about a running container
     * 
     */
    private @Nullable ContainerStateRunningPatch running;
    /**
     * @return Details about a terminated container
     * 
     */
    private @Nullable ContainerStateTerminatedPatch terminated;
    /**
     * @return Details about a waiting container
     * 
     */
    private @Nullable ContainerStateWaitingPatch waiting;

    private ContainerStatePatch() {}
    /**
     * @return Details about a running container
     * 
     */
    public Optional<ContainerStateRunningPatch> running() {
        return Optional.ofNullable(this.running);
    }
    /**
     * @return Details about a terminated container
     * 
     */
    public Optional<ContainerStateTerminatedPatch> terminated() {
        return Optional.ofNullable(this.terminated);
    }
    /**
     * @return Details about a waiting container
     * 
     */
    public Optional<ContainerStateWaitingPatch> waiting() {
        return Optional.ofNullable(this.waiting);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ContainerStatePatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable ContainerStateRunningPatch running;
        private @Nullable ContainerStateTerminatedPatch terminated;
        private @Nullable ContainerStateWaitingPatch waiting;
        public Builder() {}
        public Builder(ContainerStatePatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.running = defaults.running;
    	      this.terminated = defaults.terminated;
    	      this.waiting = defaults.waiting;
        }

        @CustomType.Setter
        public Builder running(@Nullable ContainerStateRunningPatch running) {

            this.running = running;
            return this;
        }
        @CustomType.Setter
        public Builder terminated(@Nullable ContainerStateTerminatedPatch terminated) {

            this.terminated = terminated;
            return this;
        }
        @CustomType.Setter
        public Builder waiting(@Nullable ContainerStateWaitingPatch waiting) {

            this.waiting = waiting;
            return this;
        }
        public ContainerStatePatch build() {
            final var _resultValue = new ContainerStatePatch();
            _resultValue.running = running;
            _resultValue.terminated = terminated;
            _resultValue.waiting = waiting;
            return _resultValue;
        }
    }
}
