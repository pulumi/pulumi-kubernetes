// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.core.v1.inputs.LinuxContainerUserArgs;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ContainerUser represents user identity information
 * 
 */
public final class ContainerUserArgs extends com.pulumi.resources.ResourceArgs {

    public static final ContainerUserArgs Empty = new ContainerUserArgs();

    /**
     * Linux holds user identity information initially attached to the first process of the containers in Linux. Note that the actual running identity can be changed if the process has enough privilege to do so.
     * 
     */
    @Import(name="linux")
    private @Nullable Output<LinuxContainerUserArgs> linux;

    /**
     * @return Linux holds user identity information initially attached to the first process of the containers in Linux. Note that the actual running identity can be changed if the process has enough privilege to do so.
     * 
     */
    public Optional<Output<LinuxContainerUserArgs>> linux() {
        return Optional.ofNullable(this.linux);
    }

    private ContainerUserArgs() {}

    private ContainerUserArgs(ContainerUserArgs $) {
        this.linux = $.linux;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ContainerUserArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ContainerUserArgs $;

        public Builder() {
            $ = new ContainerUserArgs();
        }

        public Builder(ContainerUserArgs defaults) {
            $ = new ContainerUserArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param linux Linux holds user identity information initially attached to the first process of the containers in Linux. Note that the actual running identity can be changed if the process has enough privilege to do so.
         * 
         * @return builder
         * 
         */
        public Builder linux(@Nullable Output<LinuxContainerUserArgs> linux) {
            $.linux = linux;
            return this;
        }

        /**
         * @param linux Linux holds user identity information initially attached to the first process of the containers in Linux. Note that the actual running identity can be changed if the process has enough privilege to do so.
         * 
         * @return builder
         * 
         */
        public Builder linux(LinuxContainerUserArgs linux) {
            return linux(Output.of(linux));
        }

        public ContainerUserArgs build() {
            return $;
        }
    }

}
