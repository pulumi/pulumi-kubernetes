// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.apiextensions.v1beta1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.apiextensions.v1beta1.outputs.CustomResourceDefinitionConditionPatch;
import com.pulumi.kubernetes.apiextensions.v1beta1.outputs.CustomResourceDefinitionNamesPatch;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class CustomResourceDefinitionStatusPatch {
    /**
     * @return acceptedNames are the names that are actually being used to serve discovery. They may be different than the names in spec.
     * 
     */
    private @Nullable CustomResourceDefinitionNamesPatch acceptedNames;
    /**
     * @return conditions indicate state for particular aspects of a CustomResourceDefinition
     * 
     */
    private @Nullable List<CustomResourceDefinitionConditionPatch> conditions;
    /**
     * @return storedVersions lists all versions of CustomResources that were ever persisted. Tracking these versions allows a migration path for stored versions in etcd. The field is mutable so a migration controller can finish a migration to another version (ensuring no old objects are left in storage), and then remove the rest of the versions from this list. Versions may not be removed from `spec.versions` while they exist in this list.
     * 
     */
    private @Nullable List<String> storedVersions;

    private CustomResourceDefinitionStatusPatch() {}
    /**
     * @return acceptedNames are the names that are actually being used to serve discovery. They may be different than the names in spec.
     * 
     */
    public Optional<CustomResourceDefinitionNamesPatch> acceptedNames() {
        return Optional.ofNullable(this.acceptedNames);
    }
    /**
     * @return conditions indicate state for particular aspects of a CustomResourceDefinition
     * 
     */
    public List<CustomResourceDefinitionConditionPatch> conditions() {
        return this.conditions == null ? List.of() : this.conditions;
    }
    /**
     * @return storedVersions lists all versions of CustomResources that were ever persisted. Tracking these versions allows a migration path for stored versions in etcd. The field is mutable so a migration controller can finish a migration to another version (ensuring no old objects are left in storage), and then remove the rest of the versions from this list. Versions may not be removed from `spec.versions` while they exist in this list.
     * 
     */
    public List<String> storedVersions() {
        return this.storedVersions == null ? List.of() : this.storedVersions;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(CustomResourceDefinitionStatusPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable CustomResourceDefinitionNamesPatch acceptedNames;
        private @Nullable List<CustomResourceDefinitionConditionPatch> conditions;
        private @Nullable List<String> storedVersions;
        public Builder() {}
        public Builder(CustomResourceDefinitionStatusPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.acceptedNames = defaults.acceptedNames;
    	      this.conditions = defaults.conditions;
    	      this.storedVersions = defaults.storedVersions;
        }

        @CustomType.Setter
        public Builder acceptedNames(@Nullable CustomResourceDefinitionNamesPatch acceptedNames) {

            this.acceptedNames = acceptedNames;
            return this;
        }
        @CustomType.Setter
        public Builder conditions(@Nullable List<CustomResourceDefinitionConditionPatch> conditions) {

            this.conditions = conditions;
            return this;
        }
        public Builder conditions(CustomResourceDefinitionConditionPatch... conditions) {
            return conditions(List.of(conditions));
        }
        @CustomType.Setter
        public Builder storedVersions(@Nullable List<String> storedVersions) {

            this.storedVersions = storedVersions;
            return this;
        }
        public Builder storedVersions(String... storedVersions) {
            return storedVersions(List.of(storedVersions));
        }
        public CustomResourceDefinitionStatusPatch build() {
            final var _resultValue = new CustomResourceDefinitionStatusPatch();
            _resultValue.acceptedNames = acceptedNames;
            _resultValue.conditions = conditions;
            _resultValue.storedVersions = storedVersions;
            return _resultValue;
        }
    }
}
