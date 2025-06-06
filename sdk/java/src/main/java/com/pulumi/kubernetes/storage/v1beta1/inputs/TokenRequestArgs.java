// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.storage.v1beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * TokenRequest contains parameters of a service account token.
 * 
 */
public final class TokenRequestArgs extends com.pulumi.resources.ResourceArgs {

    public static final TokenRequestArgs Empty = new TokenRequestArgs();

    /**
     * Audience is the intended audience of the token in &#34;TokenRequestSpec&#34;. It will default to the audiences of kube apiserver.
     * 
     */
    @Import(name="audience", required=true)
    private Output<String> audience;

    /**
     * @return Audience is the intended audience of the token in &#34;TokenRequestSpec&#34;. It will default to the audiences of kube apiserver.
     * 
     */
    public Output<String> audience() {
        return this.audience;
    }

    /**
     * ExpirationSeconds is the duration of validity of the token in &#34;TokenRequestSpec&#34;. It has the same default value of &#34;ExpirationSeconds&#34; in &#34;TokenRequestSpec&#34;
     * 
     */
    @Import(name="expirationSeconds")
    private @Nullable Output<Integer> expirationSeconds;

    /**
     * @return ExpirationSeconds is the duration of validity of the token in &#34;TokenRequestSpec&#34;. It has the same default value of &#34;ExpirationSeconds&#34; in &#34;TokenRequestSpec&#34;
     * 
     */
    public Optional<Output<Integer>> expirationSeconds() {
        return Optional.ofNullable(this.expirationSeconds);
    }

    private TokenRequestArgs() {}

    private TokenRequestArgs(TokenRequestArgs $) {
        this.audience = $.audience;
        this.expirationSeconds = $.expirationSeconds;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(TokenRequestArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private TokenRequestArgs $;

        public Builder() {
            $ = new TokenRequestArgs();
        }

        public Builder(TokenRequestArgs defaults) {
            $ = new TokenRequestArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param audience Audience is the intended audience of the token in &#34;TokenRequestSpec&#34;. It will default to the audiences of kube apiserver.
         * 
         * @return builder
         * 
         */
        public Builder audience(Output<String> audience) {
            $.audience = audience;
            return this;
        }

        /**
         * @param audience Audience is the intended audience of the token in &#34;TokenRequestSpec&#34;. It will default to the audiences of kube apiserver.
         * 
         * @return builder
         * 
         */
        public Builder audience(String audience) {
            return audience(Output.of(audience));
        }

        /**
         * @param expirationSeconds ExpirationSeconds is the duration of validity of the token in &#34;TokenRequestSpec&#34;. It has the same default value of &#34;ExpirationSeconds&#34; in &#34;TokenRequestSpec&#34;
         * 
         * @return builder
         * 
         */
        public Builder expirationSeconds(@Nullable Output<Integer> expirationSeconds) {
            $.expirationSeconds = expirationSeconds;
            return this;
        }

        /**
         * @param expirationSeconds ExpirationSeconds is the duration of validity of the token in &#34;TokenRequestSpec&#34;. It has the same default value of &#34;ExpirationSeconds&#34; in &#34;TokenRequestSpec&#34;
         * 
         * @return builder
         * 
         */
        public Builder expirationSeconds(Integer expirationSeconds) {
            return expirationSeconds(Output.of(expirationSeconds));
        }

        public TokenRequestArgs build() {
            if ($.audience == null) {
                throw new MissingRequiredPropertyException("TokenRequestArgs", "audience");
            }
            return $;
        }
    }

}
