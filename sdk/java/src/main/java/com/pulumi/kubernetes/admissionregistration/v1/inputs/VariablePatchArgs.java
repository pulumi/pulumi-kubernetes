// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.admissionregistration.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * Variable is the definition of a variable that is used for composition. A variable is defined as a named expression.
 * 
 */
public final class VariablePatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final VariablePatchArgs Empty = new VariablePatchArgs();

    /**
     * Expression is the expression that will be evaluated as the value of the variable. The CEL expression has access to the same identifiers as the CEL expressions in Validation.
     * 
     */
    @Import(name="expression")
    private @Nullable Output<String> expression;

    /**
     * @return Expression is the expression that will be evaluated as the value of the variable. The CEL expression has access to the same identifiers as the CEL expressions in Validation.
     * 
     */
    public Optional<Output<String>> expression() {
        return Optional.ofNullable(this.expression);
    }

    /**
     * Name is the name of the variable. The name must be a valid CEL identifier and unique among all variables. The variable can be accessed in other expressions through `variables` For example, if name is &#34;foo&#34;, the variable will be available as `variables.foo`
     * 
     */
    @Import(name="name")
    private @Nullable Output<String> name;

    /**
     * @return Name is the name of the variable. The name must be a valid CEL identifier and unique among all variables. The variable can be accessed in other expressions through `variables` For example, if name is &#34;foo&#34;, the variable will be available as `variables.foo`
     * 
     */
    public Optional<Output<String>> name() {
        return Optional.ofNullable(this.name);
    }

    private VariablePatchArgs() {}

    private VariablePatchArgs(VariablePatchArgs $) {
        this.expression = $.expression;
        this.name = $.name;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(VariablePatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private VariablePatchArgs $;

        public Builder() {
            $ = new VariablePatchArgs();
        }

        public Builder(VariablePatchArgs defaults) {
            $ = new VariablePatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param expression Expression is the expression that will be evaluated as the value of the variable. The CEL expression has access to the same identifiers as the CEL expressions in Validation.
         * 
         * @return builder
         * 
         */
        public Builder expression(@Nullable Output<String> expression) {
            $.expression = expression;
            return this;
        }

        /**
         * @param expression Expression is the expression that will be evaluated as the value of the variable. The CEL expression has access to the same identifiers as the CEL expressions in Validation.
         * 
         * @return builder
         * 
         */
        public Builder expression(String expression) {
            return expression(Output.of(expression));
        }

        /**
         * @param name Name is the name of the variable. The name must be a valid CEL identifier and unique among all variables. The variable can be accessed in other expressions through `variables` For example, if name is &#34;foo&#34;, the variable will be available as `variables.foo`
         * 
         * @return builder
         * 
         */
        public Builder name(@Nullable Output<String> name) {
            $.name = name;
            return this;
        }

        /**
         * @param name Name is the name of the variable. The name must be a valid CEL identifier and unique among all variables. The variable can be accessed in other expressions through `variables` For example, if name is &#34;foo&#34;, the variable will be available as `variables.foo`
         * 
         * @return builder
         * 
         */
        public Builder name(String name) {
            return name(Output.of(name));
        }

        public VariablePatchArgs build() {
            return $;
        }
    }

}
