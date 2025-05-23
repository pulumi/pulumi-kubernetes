// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.core.v1.inputs.ExecActionArgs;
import com.pulumi.kubernetes.core.v1.inputs.GRPCActionArgs;
import com.pulumi.kubernetes.core.v1.inputs.HTTPGetActionArgs;
import com.pulumi.kubernetes.core.v1.inputs.TCPSocketActionArgs;
import java.lang.Integer;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.
 * 
 */
public final class ProbeArgs extends com.pulumi.resources.ResourceArgs {

    public static final ProbeArgs Empty = new ProbeArgs();

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
     * Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1.
     * 
     */
    @Import(name="failureThreshold")
    private @Nullable Output<Integer> failureThreshold;

    /**
     * @return Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1.
     * 
     */
    public Optional<Output<Integer>> failureThreshold() {
        return Optional.ofNullable(this.failureThreshold);
    }

    /**
     * GRPC specifies a GRPC HealthCheckRequest.
     * 
     */
    @Import(name="grpc")
    private @Nullable Output<GRPCActionArgs> grpc;

    /**
     * @return GRPC specifies a GRPC HealthCheckRequest.
     * 
     */
    public Optional<Output<GRPCActionArgs>> grpc() {
        return Optional.ofNullable(this.grpc);
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
     * Number of seconds after the container has started before liveness probes are initiated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
     * 
     */
    @Import(name="initialDelaySeconds")
    private @Nullable Output<Integer> initialDelaySeconds;

    /**
     * @return Number of seconds after the container has started before liveness probes are initiated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
     * 
     */
    public Optional<Output<Integer>> initialDelaySeconds() {
        return Optional.ofNullable(this.initialDelaySeconds);
    }

    /**
     * How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
     * 
     */
    @Import(name="periodSeconds")
    private @Nullable Output<Integer> periodSeconds;

    /**
     * @return How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
     * 
     */
    public Optional<Output<Integer>> periodSeconds() {
        return Optional.ofNullable(this.periodSeconds);
    }

    /**
     * Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
     * 
     */
    @Import(name="successThreshold")
    private @Nullable Output<Integer> successThreshold;

    /**
     * @return Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
     * 
     */
    public Optional<Output<Integer>> successThreshold() {
        return Optional.ofNullable(this.successThreshold);
    }

    /**
     * TCPSocket specifies a connection to a TCP port.
     * 
     */
    @Import(name="tcpSocket")
    private @Nullable Output<TCPSocketActionArgs> tcpSocket;

    /**
     * @return TCPSocket specifies a connection to a TCP port.
     * 
     */
    public Optional<Output<TCPSocketActionArgs>> tcpSocket() {
        return Optional.ofNullable(this.tcpSocket);
    }

    /**
     * Optional duration in seconds the pod needs to terminate gracefully upon probe failure. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. If this value is nil, the pod&#39;s terminationGracePeriodSeconds will be used. Otherwise, this value overrides the value provided by the pod spec. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate. Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
     * 
     */
    @Import(name="terminationGracePeriodSeconds")
    private @Nullable Output<Integer> terminationGracePeriodSeconds;

    /**
     * @return Optional duration in seconds the pod needs to terminate gracefully upon probe failure. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. If this value is nil, the pod&#39;s terminationGracePeriodSeconds will be used. Otherwise, this value overrides the value provided by the pod spec. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate. Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
     * 
     */
    public Optional<Output<Integer>> terminationGracePeriodSeconds() {
        return Optional.ofNullable(this.terminationGracePeriodSeconds);
    }

    /**
     * Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
     * 
     */
    @Import(name="timeoutSeconds")
    private @Nullable Output<Integer> timeoutSeconds;

    /**
     * @return Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
     * 
     */
    public Optional<Output<Integer>> timeoutSeconds() {
        return Optional.ofNullable(this.timeoutSeconds);
    }

    private ProbeArgs() {}

    private ProbeArgs(ProbeArgs $) {
        this.exec = $.exec;
        this.failureThreshold = $.failureThreshold;
        this.grpc = $.grpc;
        this.httpGet = $.httpGet;
        this.initialDelaySeconds = $.initialDelaySeconds;
        this.periodSeconds = $.periodSeconds;
        this.successThreshold = $.successThreshold;
        this.tcpSocket = $.tcpSocket;
        this.terminationGracePeriodSeconds = $.terminationGracePeriodSeconds;
        this.timeoutSeconds = $.timeoutSeconds;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ProbeArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ProbeArgs $;

        public Builder() {
            $ = new ProbeArgs();
        }

        public Builder(ProbeArgs defaults) {
            $ = new ProbeArgs(Objects.requireNonNull(defaults));
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
         * @param failureThreshold Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1.
         * 
         * @return builder
         * 
         */
        public Builder failureThreshold(@Nullable Output<Integer> failureThreshold) {
            $.failureThreshold = failureThreshold;
            return this;
        }

        /**
         * @param failureThreshold Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1.
         * 
         * @return builder
         * 
         */
        public Builder failureThreshold(Integer failureThreshold) {
            return failureThreshold(Output.of(failureThreshold));
        }

        /**
         * @param grpc GRPC specifies a GRPC HealthCheckRequest.
         * 
         * @return builder
         * 
         */
        public Builder grpc(@Nullable Output<GRPCActionArgs> grpc) {
            $.grpc = grpc;
            return this;
        }

        /**
         * @param grpc GRPC specifies a GRPC HealthCheckRequest.
         * 
         * @return builder
         * 
         */
        public Builder grpc(GRPCActionArgs grpc) {
            return grpc(Output.of(grpc));
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
         * @param initialDelaySeconds Number of seconds after the container has started before liveness probes are initiated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
         * 
         * @return builder
         * 
         */
        public Builder initialDelaySeconds(@Nullable Output<Integer> initialDelaySeconds) {
            $.initialDelaySeconds = initialDelaySeconds;
            return this;
        }

        /**
         * @param initialDelaySeconds Number of seconds after the container has started before liveness probes are initiated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
         * 
         * @return builder
         * 
         */
        public Builder initialDelaySeconds(Integer initialDelaySeconds) {
            return initialDelaySeconds(Output.of(initialDelaySeconds));
        }

        /**
         * @param periodSeconds How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
         * 
         * @return builder
         * 
         */
        public Builder periodSeconds(@Nullable Output<Integer> periodSeconds) {
            $.periodSeconds = periodSeconds;
            return this;
        }

        /**
         * @param periodSeconds How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
         * 
         * @return builder
         * 
         */
        public Builder periodSeconds(Integer periodSeconds) {
            return periodSeconds(Output.of(periodSeconds));
        }

        /**
         * @param successThreshold Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
         * 
         * @return builder
         * 
         */
        public Builder successThreshold(@Nullable Output<Integer> successThreshold) {
            $.successThreshold = successThreshold;
            return this;
        }

        /**
         * @param successThreshold Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
         * 
         * @return builder
         * 
         */
        public Builder successThreshold(Integer successThreshold) {
            return successThreshold(Output.of(successThreshold));
        }

        /**
         * @param tcpSocket TCPSocket specifies a connection to a TCP port.
         * 
         * @return builder
         * 
         */
        public Builder tcpSocket(@Nullable Output<TCPSocketActionArgs> tcpSocket) {
            $.tcpSocket = tcpSocket;
            return this;
        }

        /**
         * @param tcpSocket TCPSocket specifies a connection to a TCP port.
         * 
         * @return builder
         * 
         */
        public Builder tcpSocket(TCPSocketActionArgs tcpSocket) {
            return tcpSocket(Output.of(tcpSocket));
        }

        /**
         * @param terminationGracePeriodSeconds Optional duration in seconds the pod needs to terminate gracefully upon probe failure. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. If this value is nil, the pod&#39;s terminationGracePeriodSeconds will be used. Otherwise, this value overrides the value provided by the pod spec. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate. Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
         * 
         * @return builder
         * 
         */
        public Builder terminationGracePeriodSeconds(@Nullable Output<Integer> terminationGracePeriodSeconds) {
            $.terminationGracePeriodSeconds = terminationGracePeriodSeconds;
            return this;
        }

        /**
         * @param terminationGracePeriodSeconds Optional duration in seconds the pod needs to terminate gracefully upon probe failure. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. If this value is nil, the pod&#39;s terminationGracePeriodSeconds will be used. Otherwise, this value overrides the value provided by the pod spec. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate. Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
         * 
         * @return builder
         * 
         */
        public Builder terminationGracePeriodSeconds(Integer terminationGracePeriodSeconds) {
            return terminationGracePeriodSeconds(Output.of(terminationGracePeriodSeconds));
        }

        /**
         * @param timeoutSeconds Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
         * 
         * @return builder
         * 
         */
        public Builder timeoutSeconds(@Nullable Output<Integer> timeoutSeconds) {
            $.timeoutSeconds = timeoutSeconds;
            return this;
        }

        /**
         * @param timeoutSeconds Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
         * 
         * @return builder
         * 
         */
        public Builder timeoutSeconds(Integer timeoutSeconds) {
            return timeoutSeconds(Output.of(timeoutSeconds));
        }

        public ProbeArgs build() {
            return $;
        }
    }

}
