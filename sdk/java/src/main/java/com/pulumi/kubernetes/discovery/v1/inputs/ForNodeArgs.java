// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.discovery.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.String;
import java.util.Objects;


/**
 * ForNode provides information about which nodes should consume this endpoint.
 * 
 */
public final class ForNodeArgs extends com.pulumi.resources.ResourceArgs {

    public static final ForNodeArgs Empty = new ForNodeArgs();

    /**
     * name represents the name of the node.
     * 
     */
    @Import(name="name", required=true)
    private Output<String> name;

    /**
     * @return name represents the name of the node.
     * 
     */
    public Output<String> name() {
        return this.name;
    }

    private ForNodeArgs() {}

    private ForNodeArgs(ForNodeArgs $) {
        this.name = $.name;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ForNodeArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ForNodeArgs $;

        public Builder() {
            $ = new ForNodeArgs();
        }

        public Builder(ForNodeArgs defaults) {
            $ = new ForNodeArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param name name represents the name of the node.
         * 
         * @return builder
         * 
         */
        public Builder name(Output<String> name) {
            $.name = name;
            return this;
        }

        /**
         * @param name name represents the name of the node.
         * 
         * @return builder
         * 
         */
        public Builder name(String name) {
            return name(Output.of(name));
        }

        public ForNodeArgs build() {
            if ($.name == null) {
                throw new MissingRequiredPropertyException("ForNodeArgs", "name");
            }
            return $;
        }
    }

}
