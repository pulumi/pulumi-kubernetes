// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.core.v1.inputs.ExecActionArgs;
import com.pulumi.kubernetes.core.v1.inputs.HTTPGetActionArgs;
import com.pulumi.kubernetes.core.v1.inputs.SleepActionArgs;
import com.pulumi.kubernetes.core.v1.inputs.TCPSocketActionArgs;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * LifecycleHandler defines a specific action that should be taken in a lifecycle hook. One and only one of the fields, except TCPSocket must be specified.
 * 
 */
public final class LifecycleHandlerArgs extends com.pulumi.resources.ResourceArgs {

    public static final LifecycleHandlerArgs Empty = new LifecycleHandlerArgs();

    /**
     * Exec specifies a command to execute in the container.
     * 
     */
    @Import(name="exec")
    private @Nullable Output<ExecActionArgs> exec;

    /**
     * @return Exec specifies a command to execute in the container.
     * 
     */
    public Optional<Output<ExecActionArgs>> exec() {
        return Optional.ofNullable(this.exec);
    }

    /**
     * HTTPGet specifies an HTTP GET request to perform.
     * 
     */
    @Import(name="httpGet")
    private @Nullable Output<HTTPGetActionArgs> httpGet;

    /**
     * @return HTTPGet specifies an HTTP GET request to perform.
     * 
     */
    public Optional<Output<HTTPGetActionArgs>> httpGet() {
        return Optional.ofNullable(this.httpGet);
    }

    /**
     * Sleep represents a duration that the container should sleep.
     * 
     */
    @Import(name="sleep")
    private @Nullable Output<SleepActionArgs> sleep;

    /**
     * @return Sleep represents a duration that the container should sleep.
     * 
     */
    public Optional<Output<SleepActionArgs>> sleep() {
        return Optional.ofNullable(this.sleep);
    }

    /**
     * Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept for backward compatibility. There is no validation of this field and lifecycle hooks will fail at runtime when it is specified.
     * 
     */
    @Import(name="tcpSocket")
    private @Nullable Output<TCPSocketActionArgs> tcpSocket;

    /**
     * @return Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept for backward compatibility. There is no validation of this field and lifecycle hooks will fail at runtime when it is specified.
     * 
     */
    public Optional<Output<TCPSocketActionArgs>> tcpSocket() {
        return Optional.ofNullable(this.tcpSocket);
    }

    private LifecycleHandlerArgs() {}

    private LifecycleHandlerArgs(LifecycleHandlerArgs $) {
        this.exec = $.exec;
        this.httpGet = $.httpGet;
        this.sleep = $.sleep;
        this.tcpSocket = $.tcpSocket;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(LifecycleHandlerArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private LifecycleHandlerArgs $;

        public Builder() {
            $ = new LifecycleHandlerArgs();
        }

        public Builder(LifecycleHandlerArgs defaults) {
            $ = new LifecycleHandlerArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param exec Exec specifies a command to execute in the container.
         * 
         * @return builder
         * 
         */
        public Builder exec(@Nullable Output<ExecActionArgs> exec) {
            $.exec = exec;
            return this;
        }

        /**
         * @param exec Exec specifies a command to execute in the container.
         * 
         * @return builder
         * 
         */
        public Builder exec(ExecActionArgs exec) {
            return exec(Output.of(exec));
        }

        /**
         * @param httpGet HTTPGet specifies an HTTP GET request to perform.
         * 
         * @return builder
         * 
         */
        public Builder httpGet(@Nullable Output<HTTPGetActionArgs> httpGet) {
            $.httpGet = httpGet;
            return this;
        }

        /**
         * @param httpGet HTTPGet specifies an HTTP GET request to perform.
         * 
         * @return builder
         * 
         */
        public Builder httpGet(HTTPGetActionArgs httpGet) {
            return httpGet(Output.of(httpGet));
        }

        /**
         * @param sleep Sleep represents a duration that the container should sleep.
         * 
         * @return builder
         * 
         */
        public Builder sleep(@Nullable Output<SleepActionArgs> sleep) {
            $.sleep = sleep;
            return this;
        }

        /**
         * @param sleep Sleep represents a duration that the container should sleep.
         * 
         * @return builder
         * 
         */
        public Builder sleep(SleepActionArgs sleep) {
            return sleep(Output.of(sleep));
        }

        /**
         * @param tcpSocket Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept for backward compatibility. There is no validation of this field and lifecycle hooks will fail at runtime when it is specified.
         * 
         * @return builder
         * 
         */
        public Builder tcpSocket(@Nullable Output<TCPSocketActionArgs> tcpSocket) {
            $.tcpSocket = tcpSocket;
            return this;
        }

        /**
         * @param tcpSocket Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept for backward compatibility. There is no validation of this field and lifecycle hooks will fail at runtime when it is specified.
         * 
         * @return builder
         * 
         */
        public Builder tcpSocket(TCPSocketActionArgs tcpSocket) {
            return tcpSocket(Output.of(tcpSocket));
        }

        public LifecycleHandlerArgs build() {
            return $;
        }
    }

}
