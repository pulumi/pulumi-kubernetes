// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1alpha2.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.resource.v1alpha2.inputs.NamedResourcesAttributeArgs;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * NamedResourcesInstance represents one individual hardware instance that can be selected based on its attributes.
 * 
 */
public final class NamedResourcesInstanceArgs extends com.pulumi.resources.ResourceArgs {

    public static final NamedResourcesInstanceArgs Empty = new NamedResourcesInstanceArgs();

    /**
     * Attributes defines the attributes of this resource instance. The name of each attribute must be unique.
     * 
     */
    @Import(name="attributes")
    private @Nullable Output<List<NamedResourcesAttributeArgs>> attributes;

    /**
     * @return Attributes defines the attributes of this resource instance. The name of each attribute must be unique.
     * 
     */
    public Optional<Output<List<NamedResourcesAttributeArgs>>> attributes() {
        return Optional.ofNullable(this.attributes);
    }

    /**
     * Name is unique identifier among all resource instances managed by the driver on the node. It must be a DNS subdomain.
     * 
     */
    @Import(name="name", required=true)
    private Output<String> name;

    /**
     * @return Name is unique identifier among all resource instances managed by the driver on the node. It must be a DNS subdomain.
     * 
     */
    public Output<String> name() {
        return this.name;
    }

    private NamedResourcesInstanceArgs() {}

    private NamedResourcesInstanceArgs(NamedResourcesInstanceArgs $) {
        this.attributes = $.attributes;
        this.name = $.name;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(NamedResourcesInstanceArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private NamedResourcesInstanceArgs $;

        public Builder() {
            $ = new NamedResourcesInstanceArgs();
        }

        public Builder(NamedResourcesInstanceArgs defaults) {
            $ = new NamedResourcesInstanceArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param attributes Attributes defines the attributes of this resource instance. The name of each attribute must be unique.
         * 
         * @return builder
         * 
         */
        public Builder attributes(@Nullable Output<List<NamedResourcesAttributeArgs>> attributes) {
            $.attributes = attributes;
            return this;
        }

        /**
         * @param attributes Attributes defines the attributes of this resource instance. The name of each attribute must be unique.
         * 
         * @return builder
         * 
         */
        public Builder attributes(List<NamedResourcesAttributeArgs> attributes) {
            return attributes(Output.of(attributes));
        }

        /**
         * @param attributes Attributes defines the attributes of this resource instance. The name of each attribute must be unique.
         * 
         * @return builder
         * 
         */
        public Builder attributes(NamedResourcesAttributeArgs... attributes) {
            return attributes(List.of(attributes));
        }

        /**
         * @param name Name is unique identifier among all resource instances managed by the driver on the node. It must be a DNS subdomain.
         * 
         * @return builder
         * 
         */
        public Builder name(Output<String> name) {
            $.name = name;
            return this;
        }

        /**
         * @param name Name is unique identifier among all resource instances managed by the driver on the node. It must be a DNS subdomain.
         * 
         * @return builder
         * 
         */
        public Builder name(String name) {
            return name(Output.of(name));
        }

        public NamedResourcesInstanceArgs build() {
            if ($.name == null) {
                throw new MissingRequiredPropertyException("NamedResourcesInstanceArgs", "name");
            }
            return $;
        }
    }

}
