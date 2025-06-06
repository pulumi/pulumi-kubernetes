// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.discovery.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.discovery.v1.inputs.ForNodeArgs;
import com.pulumi.kubernetes.discovery.v1.inputs.ForZoneArgs;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * EndpointHints provides hints describing how an endpoint should be consumed.
 * 
 */
public final class EndpointHintsArgs extends com.pulumi.resources.ResourceArgs {

    public static final EndpointHintsArgs Empty = new EndpointHintsArgs();

    /**
     * forNodes indicates the node(s) this endpoint should be consumed by when using topology aware routing. May contain a maximum of 8 entries. This is an Alpha feature and is only used when the PreferSameTrafficDistribution feature gate is enabled.
     * 
     */
    @Import(name="forNodes")
    private @Nullable Output<List<ForNodeArgs>> forNodes;

    /**
     * @return forNodes indicates the node(s) this endpoint should be consumed by when using topology aware routing. May contain a maximum of 8 entries. This is an Alpha feature and is only used when the PreferSameTrafficDistribution feature gate is enabled.
     * 
     */
    public Optional<Output<List<ForNodeArgs>>> forNodes() {
        return Optional.ofNullable(this.forNodes);
    }

    /**
     * forZones indicates the zone(s) this endpoint should be consumed by when using topology aware routing. May contain a maximum of 8 entries.
     * 
     */
    @Import(name="forZones")
    private @Nullable Output<List<ForZoneArgs>> forZones;

    /**
     * @return forZones indicates the zone(s) this endpoint should be consumed by when using topology aware routing. May contain a maximum of 8 entries.
     * 
     */
    public Optional<Output<List<ForZoneArgs>>> forZones() {
        return Optional.ofNullable(this.forZones);
    }

    private EndpointHintsArgs() {}

    private EndpointHintsArgs(EndpointHintsArgs $) {
        this.forNodes = $.forNodes;
        this.forZones = $.forZones;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(EndpointHintsArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private EndpointHintsArgs $;

        public Builder() {
            $ = new EndpointHintsArgs();
        }

        public Builder(EndpointHintsArgs defaults) {
            $ = new EndpointHintsArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param forNodes forNodes indicates the node(s) this endpoint should be consumed by when using topology aware routing. May contain a maximum of 8 entries. This is an Alpha feature and is only used when the PreferSameTrafficDistribution feature gate is enabled.
         * 
         * @return builder
         * 
         */
        public Builder forNodes(@Nullable Output<List<ForNodeArgs>> forNodes) {
            $.forNodes = forNodes;
            return this;
        }

        /**
         * @param forNodes forNodes indicates the node(s) this endpoint should be consumed by when using topology aware routing. May contain a maximum of 8 entries. This is an Alpha feature and is only used when the PreferSameTrafficDistribution feature gate is enabled.
         * 
         * @return builder
         * 
         */
        public Builder forNodes(List<ForNodeArgs> forNodes) {
            return forNodes(Output.of(forNodes));
        }

        /**
         * @param forNodes forNodes indicates the node(s) this endpoint should be consumed by when using topology aware routing. May contain a maximum of 8 entries. This is an Alpha feature and is only used when the PreferSameTrafficDistribution feature gate is enabled.
         * 
         * @return builder
         * 
         */
        public Builder forNodes(ForNodeArgs... forNodes) {
            return forNodes(List.of(forNodes));
        }

        /**
         * @param forZones forZones indicates the zone(s) this endpoint should be consumed by when using topology aware routing. May contain a maximum of 8 entries.
         * 
         * @return builder
         * 
         */
        public Builder forZones(@Nullable Output<List<ForZoneArgs>> forZones) {
            $.forZones = forZones;
            return this;
        }

        /**
         * @param forZones forZones indicates the zone(s) this endpoint should be consumed by when using topology aware routing. May contain a maximum of 8 entries.
         * 
         * @return builder
         * 
         */
        public Builder forZones(List<ForZoneArgs> forZones) {
            return forZones(Output.of(forZones));
        }

        /**
         * @param forZones forZones indicates the zone(s) this endpoint should be consumed by when using topology aware routing. May contain a maximum of 8 entries.
         * 
         * @return builder
         * 
         */
        public Builder forZones(ForZoneArgs... forZones) {
            return forZones(List.of(forZones));
        }

        public EndpointHintsArgs build() {
            return $;
        }
    }

}
