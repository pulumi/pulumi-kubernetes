// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.meta.v1.inputs;

import com.google.gson.JsonElement;
import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ManagedFieldsEntry is a workflow-id, a FieldSet and the group version of the resource that the fieldset applies to.
 * 
 */
public final class ManagedFieldsEntryPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final ManagedFieldsEntryPatchArgs Empty = new ManagedFieldsEntryPatchArgs();

    /**
     * APIVersion defines the version of this resource that this field set applies to. The format is &#34;group/version&#34; just like the top-level APIVersion field. It is necessary to track the version of a field set because it cannot be automatically converted.
     * 
     */
    @Import(name="apiVersion")
    private @Nullable Output<String> apiVersion;

    /**
     * @return APIVersion defines the version of this resource that this field set applies to. The format is &#34;group/version&#34; just like the top-level APIVersion field. It is necessary to track the version of a field set because it cannot be automatically converted.
     * 
     */
    public Optional<Output<String>> apiVersion() {
        return Optional.ofNullable(this.apiVersion);
    }

    /**
     * FieldsType is the discriminator for the different fields format and version. There is currently only one possible value: &#34;FieldsV1&#34;
     * 
     */
    @Import(name="fieldsType")
    private @Nullable Output<String> fieldsType;

    /**
     * @return FieldsType is the discriminator for the different fields format and version. There is currently only one possible value: &#34;FieldsV1&#34;
     * 
     */
    public Optional<Output<String>> fieldsType() {
        return Optional.ofNullable(this.fieldsType);
    }

    /**
     * FieldsV1 holds the first JSON version format as described in the &#34;FieldsV1&#34; type.
     * 
     */
    @Import(name="fieldsV1")
    private @Nullable Output<JsonElement> fieldsV1;

    /**
     * @return FieldsV1 holds the first JSON version format as described in the &#34;FieldsV1&#34; type.
     * 
     */
    public Optional<Output<JsonElement>> fieldsV1() {
        return Optional.ofNullable(this.fieldsV1);
    }

    /**
     * Manager is an identifier of the workflow managing these fields.
     * 
     */
    @Import(name="manager")
    private @Nullable Output<String> manager;

    /**
     * @return Manager is an identifier of the workflow managing these fields.
     * 
     */
    public Optional<Output<String>> manager() {
        return Optional.ofNullable(this.manager);
    }

    /**
     * Operation is the type of operation which lead to this ManagedFieldsEntry being created. The only valid values for this field are &#39;Apply&#39; and &#39;Update&#39;.
     * 
     */
    @Import(name="operation")
    private @Nullable Output<String> operation;

    /**
     * @return Operation is the type of operation which lead to this ManagedFieldsEntry being created. The only valid values for this field are &#39;Apply&#39; and &#39;Update&#39;.
     * 
     */
    public Optional<Output<String>> operation() {
        return Optional.ofNullable(this.operation);
    }

    /**
     * Subresource is the name of the subresource used to update that object, or empty string if the object was updated through the main resource. The value of this field is used to distinguish between managers, even if they share the same name. For example, a status update will be distinct from a regular update using the same manager name. Note that the APIVersion field is not related to the Subresource field and it always corresponds to the version of the main resource.
     * 
     */
    @Import(name="subresource")
    private @Nullable Output<String> subresource;

    /**
     * @return Subresource is the name of the subresource used to update that object, or empty string if the object was updated through the main resource. The value of this field is used to distinguish between managers, even if they share the same name. For example, a status update will be distinct from a regular update using the same manager name. Note that the APIVersion field is not related to the Subresource field and it always corresponds to the version of the main resource.
     * 
     */
    public Optional<Output<String>> subresource() {
        return Optional.ofNullable(this.subresource);
    }

    /**
     * Time is the timestamp of when the ManagedFields entry was added. The timestamp will also be updated if a field is added, the manager changes any of the owned fields value or removes a field. The timestamp does not update when a field is removed from the entry because another manager took it over.
     * 
     */
    @Import(name="time")
    private @Nullable Output<String> time;

    /**
     * @return Time is the timestamp of when the ManagedFields entry was added. The timestamp will also be updated if a field is added, the manager changes any of the owned fields value or removes a field. The timestamp does not update when a field is removed from the entry because another manager took it over.
     * 
     */
    public Optional<Output<String>> time() {
        return Optional.ofNullable(this.time);
    }

    private ManagedFieldsEntryPatchArgs() {}

    private ManagedFieldsEntryPatchArgs(ManagedFieldsEntryPatchArgs $) {
        this.apiVersion = $.apiVersion;
        this.fieldsType = $.fieldsType;
        this.fieldsV1 = $.fieldsV1;
        this.manager = $.manager;
        this.operation = $.operation;
        this.subresource = $.subresource;
        this.time = $.time;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ManagedFieldsEntryPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ManagedFieldsEntryPatchArgs $;

        public Builder() {
            $ = new ManagedFieldsEntryPatchArgs();
        }

        public Builder(ManagedFieldsEntryPatchArgs defaults) {
            $ = new ManagedFieldsEntryPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param apiVersion APIVersion defines the version of this resource that this field set applies to. The format is &#34;group/version&#34; just like the top-level APIVersion field. It is necessary to track the version of a field set because it cannot be automatically converted.
         * 
         * @return builder
         * 
         */
        public Builder apiVersion(@Nullable Output<String> apiVersion) {
            $.apiVersion = apiVersion;
            return this;
        }

        /**
         * @param apiVersion APIVersion defines the version of this resource that this field set applies to. The format is &#34;group/version&#34; just like the top-level APIVersion field. It is necessary to track the version of a field set because it cannot be automatically converted.
         * 
         * @return builder
         * 
         */
        public Builder apiVersion(String apiVersion) {
            return apiVersion(Output.of(apiVersion));
        }

        /**
         * @param fieldsType FieldsType is the discriminator for the different fields format and version. There is currently only one possible value: &#34;FieldsV1&#34;
         * 
         * @return builder
         * 
         */
        public Builder fieldsType(@Nullable Output<String> fieldsType) {
            $.fieldsType = fieldsType;
            return this;
        }

        /**
         * @param fieldsType FieldsType is the discriminator for the different fields format and version. There is currently only one possible value: &#34;FieldsV1&#34;
         * 
         * @return builder
         * 
         */
        public Builder fieldsType(String fieldsType) {
            return fieldsType(Output.of(fieldsType));
        }

        /**
         * @param fieldsV1 FieldsV1 holds the first JSON version format as described in the &#34;FieldsV1&#34; type.
         * 
         * @return builder
         * 
         */
        public Builder fieldsV1(@Nullable Output<JsonElement> fieldsV1) {
            $.fieldsV1 = fieldsV1;
            return this;
        }

        /**
         * @param fieldsV1 FieldsV1 holds the first JSON version format as described in the &#34;FieldsV1&#34; type.
         * 
         * @return builder
         * 
         */
        public Builder fieldsV1(JsonElement fieldsV1) {
            return fieldsV1(Output.of(fieldsV1));
        }

        /**
         * @param manager Manager is an identifier of the workflow managing these fields.
         * 
         * @return builder
         * 
         */
        public Builder manager(@Nullable Output<String> manager) {
            $.manager = manager;
            return this;
        }

        /**
         * @param manager Manager is an identifier of the workflow managing these fields.
         * 
         * @return builder
         * 
         */
        public Builder manager(String manager) {
            return manager(Output.of(manager));
        }

        /**
         * @param operation Operation is the type of operation which lead to this ManagedFieldsEntry being created. The only valid values for this field are &#39;Apply&#39; and &#39;Update&#39;.
         * 
         * @return builder
         * 
         */
        public Builder operation(@Nullable Output<String> operation) {
            $.operation = operation;
            return this;
        }

        /**
         * @param operation Operation is the type of operation which lead to this ManagedFieldsEntry being created. The only valid values for this field are &#39;Apply&#39; and &#39;Update&#39;.
         * 
         * @return builder
         * 
         */
        public Builder operation(String operation) {
            return operation(Output.of(operation));
        }

        /**
         * @param subresource Subresource is the name of the subresource used to update that object, or empty string if the object was updated through the main resource. The value of this field is used to distinguish between managers, even if they share the same name. For example, a status update will be distinct from a regular update using the same manager name. Note that the APIVersion field is not related to the Subresource field and it always corresponds to the version of the main resource.
         * 
         * @return builder
         * 
         */
        public Builder subresource(@Nullable Output<String> subresource) {
            $.subresource = subresource;
            return this;
        }

        /**
         * @param subresource Subresource is the name of the subresource used to update that object, or empty string if the object was updated through the main resource. The value of this field is used to distinguish between managers, even if they share the same name. For example, a status update will be distinct from a regular update using the same manager name. Note that the APIVersion field is not related to the Subresource field and it always corresponds to the version of the main resource.
         * 
         * @return builder
         * 
         */
        public Builder subresource(String subresource) {
            return subresource(Output.of(subresource));
        }

        /**
         * @param time Time is the timestamp of when the ManagedFields entry was added. The timestamp will also be updated if a field is added, the manager changes any of the owned fields value or removes a field. The timestamp does not update when a field is removed from the entry because another manager took it over.
         * 
         * @return builder
         * 
         */
        public Builder time(@Nullable Output<String> time) {
            $.time = time;
            return this;
        }

        /**
         * @param time Time is the timestamp of when the ManagedFields entry was added. The timestamp will also be updated if a field is added, the manager changes any of the owned fields value or removes a field. The timestamp does not update when a field is removed from the entry because another manager took it over.
         * 
         * @return builder
         * 
         */
        public Builder time(String time) {
            return time(Output.of(time));
        }

        public ManagedFieldsEntryPatchArgs build() {
            return $;
        }
    }

}
