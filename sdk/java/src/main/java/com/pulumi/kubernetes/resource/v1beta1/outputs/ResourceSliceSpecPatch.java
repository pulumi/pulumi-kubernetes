// *** WARNING: this file was generated by pulumi-java-gen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1beta1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.core.v1.outputs.NodeSelectorPatch;
import com.pulumi.kubernetes.resource.v1beta1.outputs.DevicePatch;
import com.pulumi.kubernetes.resource.v1beta1.outputs.ResourcePoolPatch;
import java.lang.Boolean;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class ResourceSliceSpecPatch {
    /**
     * @return AllNodes indicates that all nodes have access to the resources in the pool.
     * 
     * Exactly one of NodeName, NodeSelector and AllNodes must be set.
     * 
     */
    private @Nullable Boolean allNodes;
    /**
     * @return Devices lists some or all of the devices in this pool.
     * 
     * Must not have more than 128 entries.
     * 
     */
    private @Nullable List<DevicePatch> devices;
    /**
     * @return Driver identifies the DRA driver providing the capacity information. A field selector can be used to list only ResourceSlice objects with a certain driver name.
     * 
     * Must be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver. This field is immutable.
     * 
     */
    private @Nullable String driver;
    /**
     * @return NodeName identifies the node which provides the resources in this pool. A field selector can be used to list only ResourceSlice objects belonging to a certain node.
     * 
     * This field can be used to limit access from nodes to ResourceSlices with the same node name. It also indicates to autoscalers that adding new nodes of the same type as some old node might also make new resources available.
     * 
     * Exactly one of NodeName, NodeSelector and AllNodes must be set. This field is immutable.
     * 
     */
    private @Nullable String nodeName;
    /**
     * @return NodeSelector defines which nodes have access to the resources in the pool, when that pool is not limited to a single node.
     * 
     * Must use exactly one term.
     * 
     * Exactly one of NodeName, NodeSelector and AllNodes must be set.
     * 
     */
    private @Nullable NodeSelectorPatch nodeSelector;
    /**
     * @return Pool describes the pool that this ResourceSlice belongs to.
     * 
     */
    private @Nullable ResourcePoolPatch pool;

    private ResourceSliceSpecPatch() {}
    /**
     * @return AllNodes indicates that all nodes have access to the resources in the pool.
     * 
     * Exactly one of NodeName, NodeSelector and AllNodes must be set.
     * 
     */
    public Optional<Boolean> allNodes() {
        return Optional.ofNullable(this.allNodes);
    }
    /**
     * @return Devices lists some or all of the devices in this pool.
     * 
     * Must not have more than 128 entries.
     * 
     */
    public List<DevicePatch> devices() {
        return this.devices == null ? List.of() : this.devices;
    }
    /**
     * @return Driver identifies the DRA driver providing the capacity information. A field selector can be used to list only ResourceSlice objects with a certain driver name.
     * 
     * Must be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver. This field is immutable.
     * 
     */
    public Optional<String> driver() {
        return Optional.ofNullable(this.driver);
    }
    /**
     * @return NodeName identifies the node which provides the resources in this pool. A field selector can be used to list only ResourceSlice objects belonging to a certain node.
     * 
     * This field can be used to limit access from nodes to ResourceSlices with the same node name. It also indicates to autoscalers that adding new nodes of the same type as some old node might also make new resources available.
     * 
     * Exactly one of NodeName, NodeSelector and AllNodes must be set. This field is immutable.
     * 
     */
    public Optional<String> nodeName() {
        return Optional.ofNullable(this.nodeName);
    }
    /**
     * @return NodeSelector defines which nodes have access to the resources in the pool, when that pool is not limited to a single node.
     * 
     * Must use exactly one term.
     * 
     * Exactly one of NodeName, NodeSelector and AllNodes must be set.
     * 
     */
    public Optional<NodeSelectorPatch> nodeSelector() {
        return Optional.ofNullable(this.nodeSelector);
    }
    /**
     * @return Pool describes the pool that this ResourceSlice belongs to.
     * 
     */
    public Optional<ResourcePoolPatch> pool() {
        return Optional.ofNullable(this.pool);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ResourceSliceSpecPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable Boolean allNodes;
        private @Nullable List<DevicePatch> devices;
        private @Nullable String driver;
        private @Nullable String nodeName;
        private @Nullable NodeSelectorPatch nodeSelector;
        private @Nullable ResourcePoolPatch pool;
        public Builder() {}
        public Builder(ResourceSliceSpecPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.allNodes = defaults.allNodes;
    	      this.devices = defaults.devices;
    	      this.driver = defaults.driver;
    	      this.nodeName = defaults.nodeName;
    	      this.nodeSelector = defaults.nodeSelector;
    	      this.pool = defaults.pool;
        }

        @CustomType.Setter
        public Builder allNodes(@Nullable Boolean allNodes) {

            this.allNodes = allNodes;
            return this;
        }
        @CustomType.Setter
        public Builder devices(@Nullable List<DevicePatch> devices) {

            this.devices = devices;
            return this;
        }
        public Builder devices(DevicePatch... devices) {
            return devices(List.of(devices));
        }
        @CustomType.Setter
        public Builder driver(@Nullable String driver) {

            this.driver = driver;
            return this;
        }
        @CustomType.Setter
        public Builder nodeName(@Nullable String nodeName) {

            this.nodeName = nodeName;
            return this;
        }
        @CustomType.Setter
        public Builder nodeSelector(@Nullable NodeSelectorPatch nodeSelector) {

            this.nodeSelector = nodeSelector;
            return this;
        }
        @CustomType.Setter
        public Builder pool(@Nullable ResourcePoolPatch pool) {

            this.pool = pool;
            return this;
        }
        public ResourceSliceSpecPatch build() {
            final var _resultValue = new ResourceSliceSpecPatch();
            _resultValue.allNodes = allNodes;
            _resultValue.devices = devices;
            _resultValue.driver = driver;
            _resultValue.nodeName = nodeName;
            _resultValue.nodeSelector = nodeSelector;
            _resultValue.pool = pool;
            return _resultValue;
        }
    }
}
