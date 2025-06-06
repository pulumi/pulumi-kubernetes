// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.flowcontrol.v1beta1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.flowcontrol.v1beta1.outputs.LimitResponse;
import java.lang.Integer;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class LimitedPriorityLevelConfiguration {
    /**
     * @return `assuredConcurrencyShares` (ACS) configures the execution limit, which is a limit on the number of requests of this priority level that may be exeucting at a given time.  ACS must be a positive number. The server&#39;s concurrency limit (SCL) is divided among the concurrency-controlled priority levels in proportion to their assured concurrency shares. This produces the assured concurrency value (ACV) --- the number of requests that may be executing at a time --- for each such priority level:
     * 
     *             ACV(l) = ceil( SCL * ACS(l) / ( sum[priority levels k] ACS(k) ) )
     * 
     * bigger numbers of ACS mean more reserved concurrent requests (at the expense of every other PL). This field has a default value of 30.
     * 
     */
    private @Nullable Integer assuredConcurrencyShares;
    /**
     * @return `limitResponse` indicates what to do with requests that can not be executed right now
     * 
     */
    private @Nullable LimitResponse limitResponse;

    private LimitedPriorityLevelConfiguration() {}
    /**
     * @return `assuredConcurrencyShares` (ACS) configures the execution limit, which is a limit on the number of requests of this priority level that may be exeucting at a given time.  ACS must be a positive number. The server&#39;s concurrency limit (SCL) is divided among the concurrency-controlled priority levels in proportion to their assured concurrency shares. This produces the assured concurrency value (ACV) --- the number of requests that may be executing at a time --- for each such priority level:
     * 
     *             ACV(l) = ceil( SCL * ACS(l) / ( sum[priority levels k] ACS(k) ) )
     * 
     * bigger numbers of ACS mean more reserved concurrent requests (at the expense of every other PL). This field has a default value of 30.
     * 
     */
    public Optional<Integer> assuredConcurrencyShares() {
        return Optional.ofNullable(this.assuredConcurrencyShares);
    }
    /**
     * @return `limitResponse` indicates what to do with requests that can not be executed right now
     * 
     */
    public Optional<LimitResponse> limitResponse() {
        return Optional.ofNullable(this.limitResponse);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(LimitedPriorityLevelConfiguration defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable Integer assuredConcurrencyShares;
        private @Nullable LimitResponse limitResponse;
        public Builder() {}
        public Builder(LimitedPriorityLevelConfiguration defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.assuredConcurrencyShares = defaults.assuredConcurrencyShares;
    	      this.limitResponse = defaults.limitResponse;
        }

        @CustomType.Setter
        public Builder assuredConcurrencyShares(@Nullable Integer assuredConcurrencyShares) {

            this.assuredConcurrencyShares = assuredConcurrencyShares;
            return this;
        }
        @CustomType.Setter
        public Builder limitResponse(@Nullable LimitResponse limitResponse) {

            this.limitResponse = limitResponse;
            return this;
        }
        public LimitedPriorityLevelConfiguration build() {
            final var _resultValue = new LimitedPriorityLevelConfiguration();
            _resultValue.assuredConcurrencyShares = assuredConcurrencyShares;
            _resultValue.limitResponse = limitResponse;
            return _resultValue;
        }
    }
}
