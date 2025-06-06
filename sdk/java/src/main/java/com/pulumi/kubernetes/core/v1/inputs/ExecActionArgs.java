// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ExecAction describes a &#34;run in container&#34; action.
 * 
 */
public final class ExecActionArgs extends com.pulumi.resources.ResourceArgs {

    public static final ExecActionArgs Empty = new ExecActionArgs();

    /**
     * Command is the command line to execute inside the container, the working directory for the command  is root (&#39;/&#39;) in the container&#39;s filesystem. The command is simply exec&#39;d, it is not run inside a shell, so traditional shell instructions (&#39;|&#39;, etc) won&#39;t work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
     * 
     */
    @Import(name="command")
    private @Nullable Output<List<String>> command;

    /**
     * @return Command is the command line to execute inside the container, the working directory for the command  is root (&#39;/&#39;) in the container&#39;s filesystem. The command is simply exec&#39;d, it is not run inside a shell, so traditional shell instructions (&#39;|&#39;, etc) won&#39;t work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
     * 
     */
    public Optional<Output<List<String>>> command() {
        return Optional.ofNullable(this.command);
    }

    private ExecActionArgs() {}

    private ExecActionArgs(ExecActionArgs $) {
        this.command = $.command;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ExecActionArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ExecActionArgs $;

        public Builder() {
            $ = new ExecActionArgs();
        }

        public Builder(ExecActionArgs defaults) {
            $ = new ExecActionArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param command Command is the command line to execute inside the container, the working directory for the command  is root (&#39;/&#39;) in the container&#39;s filesystem. The command is simply exec&#39;d, it is not run inside a shell, so traditional shell instructions (&#39;|&#39;, etc) won&#39;t work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
         * 
         * @return builder
         * 
         */
        public Builder command(@Nullable Output<List<String>> command) {
            $.command = command;
            return this;
        }

        /**
         * @param command Command is the command line to execute inside the container, the working directory for the command  is root (&#39;/&#39;) in the container&#39;s filesystem. The command is simply exec&#39;d, it is not run inside a shell, so traditional shell instructions (&#39;|&#39;, etc) won&#39;t work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
         * 
         * @return builder
         * 
         */
        public Builder command(List<String> command) {
            return command(Output.of(command));
        }

        /**
         * @param command Command is the command line to execute inside the container, the working directory for the command  is root (&#39;/&#39;) in the container&#39;s filesystem. The command is simply exec&#39;d, it is not run inside a shell, so traditional shell instructions (&#39;|&#39;, etc) won&#39;t work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
         * 
         * @return builder
         * 
         */
        public Builder command(String... command) {
            return command(List.of(command));
        }

        public ExecActionArgs build() {
            return $;
        }
    }

}
