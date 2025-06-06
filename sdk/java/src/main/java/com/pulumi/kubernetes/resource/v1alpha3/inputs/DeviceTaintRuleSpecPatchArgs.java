// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1alpha3.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.resource.v1alpha3.inputs.DeviceTaintPatchArgs;
import com.pulumi.kubernetes.resource.v1alpha3.inputs.DeviceTaintSelectorPatchArgs;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * DeviceTaintRuleSpec specifies the selector and one taint.
 * 
 */
public final class DeviceTaintRuleSpecPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final DeviceTaintRuleSpecPatchArgs Empty = new DeviceTaintRuleSpecPatchArgs();

    /**
     * DeviceSelector defines which device(s) the taint is applied to. All selector criteria must be satified for a device to match. The empty selector matches all devices. Without a selector, no devices are matches.
     * 
     */
    @Import(name="deviceSelector")
    private @Nullable Output<DeviceTaintSelectorPatchArgs> deviceSelector;

    /**
     * @return DeviceSelector defines which device(s) the taint is applied to. All selector criteria must be satified for a device to match. The empty selector matches all devices. Without a selector, no devices are matches.
     * 
     */
    public Optional<Output<DeviceTaintSelectorPatchArgs>> deviceSelector() {
        return Optional.ofNullable(this.deviceSelector);
    }

    /**
     * The taint that gets applied to matching devices.
     * 
     */
    @Import(name="taint")
    private @Nullable Output<DeviceTaintPatchArgs> taint;

    /**
     * @return The taint that gets applied to matching devices.
     * 
     */
    public Optional<Output<DeviceTaintPatchArgs>> taint() {
        return Optional.ofNullable(this.taint);
    }

    private DeviceTaintRuleSpecPatchArgs() {}

    private DeviceTaintRuleSpecPatchArgs(DeviceTaintRuleSpecPatchArgs $) {
        this.deviceSelector = $.deviceSelector;
        this.taint = $.taint;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(DeviceTaintRuleSpecPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private DeviceTaintRuleSpecPatchArgs $;

        public Builder() {
            $ = new DeviceTaintRuleSpecPatchArgs();
        }

        public Builder(DeviceTaintRuleSpecPatchArgs defaults) {
            $ = new DeviceTaintRuleSpecPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param deviceSelector DeviceSelector defines which device(s) the taint is applied to. All selector criteria must be satified for a device to match. The empty selector matches all devices. Without a selector, no devices are matches.
         * 
         * @return builder
         * 
         */
        public Builder deviceSelector(@Nullable Output<DeviceTaintSelectorPatchArgs> deviceSelector) {
            $.deviceSelector = deviceSelector;
            return this;
        }

        /**
         * @param deviceSelector DeviceSelector defines which device(s) the taint is applied to. All selector criteria must be satified for a device to match. The empty selector matches all devices. Without a selector, no devices are matches.
         * 
         * @return builder
         * 
         */
        public Builder deviceSelector(DeviceTaintSelectorPatchArgs deviceSelector) {
            return deviceSelector(Output.of(deviceSelector));
        }

        /**
         * @param taint The taint that gets applied to matching devices.
         * 
         * @return builder
         * 
         */
        public Builder taint(@Nullable Output<DeviceTaintPatchArgs> taint) {
            $.taint = taint;
            return this;
        }

        /**
         * @param taint The taint that gets applied to matching devices.
         * 
         * @return builder
         * 
         */
        public Builder taint(DeviceTaintPatchArgs taint) {
            return taint(Output.of(taint));
        }

        public DeviceTaintRuleSpecPatchArgs build() {
            return $;
        }
    }

}
