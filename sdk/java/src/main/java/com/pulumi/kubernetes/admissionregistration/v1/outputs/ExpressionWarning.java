// *** WARNING: this file was generated by pulumi-java-gen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.admissionregistration.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import java.lang.String;
import java.util.Objects;

@CustomType
public final class ExpressionWarning {
    /**
     * @return The path to the field that refers the expression. For example, the reference to the expression of the first item of validations is &#34;spec.validations[0].expression&#34;
     * 
     */
    private String fieldRef;
    /**
     * @return The content of type checking information in a human-readable form. Each line of the warning contains the type that the expression is checked against, followed by the type check error from the compiler.
     * 
     */
    private String warning;

    private ExpressionWarning() {}
    /**
     * @return The path to the field that refers the expression. For example, the reference to the expression of the first item of validations is &#34;spec.validations[0].expression&#34;
     * 
     */
    public String fieldRef() {
        return this.fieldRef;
    }
    /**
     * @return The content of type checking information in a human-readable form. Each line of the warning contains the type that the expression is checked against, followed by the type check error from the compiler.
     * 
     */
    public String warning() {
        return this.warning;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ExpressionWarning defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private String fieldRef;
        private String warning;
        public Builder() {}
        public Builder(ExpressionWarning defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.fieldRef = defaults.fieldRef;
    	      this.warning = defaults.warning;
        }

        @CustomType.Setter
        public Builder fieldRef(String fieldRef) {
            this.fieldRef = Objects.requireNonNull(fieldRef);
            return this;
        }
        @CustomType.Setter
        public Builder warning(String warning) {
            this.warning = Objects.requireNonNull(warning);
            return this;
        }
        public ExpressionWarning build() {
            final var o = new ExpressionWarning();
            o.fieldRef = fieldRef;
            o.warning = warning;
            return o;
        }
    }
}