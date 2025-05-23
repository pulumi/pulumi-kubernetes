// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.admissionregistration.v1alpha1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.admissionregistration.v1alpha1.inputs.MatchResourcesPatchArgs;
import com.pulumi.kubernetes.admissionregistration.v1alpha1.inputs.ParamRefPatchArgs;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * MutatingAdmissionPolicyBindingSpec is the specification of the MutatingAdmissionPolicyBinding.
 * 
 */
public final class MutatingAdmissionPolicyBindingSpecPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final MutatingAdmissionPolicyBindingSpecPatchArgs Empty = new MutatingAdmissionPolicyBindingSpecPatchArgs();

    /**
     * matchResources limits what resources match this binding and may be mutated by it. Note that if matchResources matches a resource, the resource must also match a policy&#39;s matchConstraints and matchConditions before the resource may be mutated. When matchResources is unset, it does not constrain resource matching, and only the policy&#39;s matchConstraints and matchConditions must match for the resource to be mutated. Additionally, matchResources.resourceRules are optional and do not constraint matching when unset. Note that this is differs from MutatingAdmissionPolicy matchConstraints, where resourceRules are required. The CREATE, UPDATE and CONNECT operations are allowed.  The DELETE operation may not be matched. &#39;*&#39; matches CREATE, UPDATE and CONNECT.
     * 
     */
    @Import(name="matchResources")
    private @Nullable Output<MatchResourcesPatchArgs> matchResources;

    /**
     * @return matchResources limits what resources match this binding and may be mutated by it. Note that if matchResources matches a resource, the resource must also match a policy&#39;s matchConstraints and matchConditions before the resource may be mutated. When matchResources is unset, it does not constrain resource matching, and only the policy&#39;s matchConstraints and matchConditions must match for the resource to be mutated. Additionally, matchResources.resourceRules are optional and do not constraint matching when unset. Note that this is differs from MutatingAdmissionPolicy matchConstraints, where resourceRules are required. The CREATE, UPDATE and CONNECT operations are allowed.  The DELETE operation may not be matched. &#39;*&#39; matches CREATE, UPDATE and CONNECT.
     * 
     */
    public Optional<Output<MatchResourcesPatchArgs>> matchResources() {
        return Optional.ofNullable(this.matchResources);
    }

    /**
     * paramRef specifies the parameter resource used to configure the admission control policy. It should point to a resource of the type specified in spec.ParamKind of the bound MutatingAdmissionPolicy. If the policy specifies a ParamKind and the resource referred to by ParamRef does not exist, this binding is considered mis-configured and the FailurePolicy of the MutatingAdmissionPolicy applied. If the policy does not specify a ParamKind then this field is ignored, and the rules are evaluated without a param.
     * 
     */
    @Import(name="paramRef")
    private @Nullable Output<ParamRefPatchArgs> paramRef;

    /**
     * @return paramRef specifies the parameter resource used to configure the admission control policy. It should point to a resource of the type specified in spec.ParamKind of the bound MutatingAdmissionPolicy. If the policy specifies a ParamKind and the resource referred to by ParamRef does not exist, this binding is considered mis-configured and the FailurePolicy of the MutatingAdmissionPolicy applied. If the policy does not specify a ParamKind then this field is ignored, and the rules are evaluated without a param.
     * 
     */
    public Optional<Output<ParamRefPatchArgs>> paramRef() {
        return Optional.ofNullable(this.paramRef);
    }

    /**
     * policyName references a MutatingAdmissionPolicy name which the MutatingAdmissionPolicyBinding binds to. If the referenced resource does not exist, this binding is considered invalid and will be ignored Required.
     * 
     */
    @Import(name="policyName")
    private @Nullable Output<String> policyName;

    /**
     * @return policyName references a MutatingAdmissionPolicy name which the MutatingAdmissionPolicyBinding binds to. If the referenced resource does not exist, this binding is considered invalid and will be ignored Required.
     * 
     */
    public Optional<Output<String>> policyName() {
        return Optional.ofNullable(this.policyName);
    }

    private MutatingAdmissionPolicyBindingSpecPatchArgs() {}

    private MutatingAdmissionPolicyBindingSpecPatchArgs(MutatingAdmissionPolicyBindingSpecPatchArgs $) {
        this.matchResources = $.matchResources;
        this.paramRef = $.paramRef;
        this.policyName = $.policyName;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(MutatingAdmissionPolicyBindingSpecPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private MutatingAdmissionPolicyBindingSpecPatchArgs $;

        public Builder() {
            $ = new MutatingAdmissionPolicyBindingSpecPatchArgs();
        }

        public Builder(MutatingAdmissionPolicyBindingSpecPatchArgs defaults) {
            $ = new MutatingAdmissionPolicyBindingSpecPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param matchResources matchResources limits what resources match this binding and may be mutated by it. Note that if matchResources matches a resource, the resource must also match a policy&#39;s matchConstraints and matchConditions before the resource may be mutated. When matchResources is unset, it does not constrain resource matching, and only the policy&#39;s matchConstraints and matchConditions must match for the resource to be mutated. Additionally, matchResources.resourceRules are optional and do not constraint matching when unset. Note that this is differs from MutatingAdmissionPolicy matchConstraints, where resourceRules are required. The CREATE, UPDATE and CONNECT operations are allowed.  The DELETE operation may not be matched. &#39;*&#39; matches CREATE, UPDATE and CONNECT.
         * 
         * @return builder
         * 
         */
        public Builder matchResources(@Nullable Output<MatchResourcesPatchArgs> matchResources) {
            $.matchResources = matchResources;
            return this;
        }

        /**
         * @param matchResources matchResources limits what resources match this binding and may be mutated by it. Note that if matchResources matches a resource, the resource must also match a policy&#39;s matchConstraints and matchConditions before the resource may be mutated. When matchResources is unset, it does not constrain resource matching, and only the policy&#39;s matchConstraints and matchConditions must match for the resource to be mutated. Additionally, matchResources.resourceRules are optional and do not constraint matching when unset. Note that this is differs from MutatingAdmissionPolicy matchConstraints, where resourceRules are required. The CREATE, UPDATE and CONNECT operations are allowed.  The DELETE operation may not be matched. &#39;*&#39; matches CREATE, UPDATE and CONNECT.
         * 
         * @return builder
         * 
         */
        public Builder matchResources(MatchResourcesPatchArgs matchResources) {
            return matchResources(Output.of(matchResources));
        }

        /**
         * @param paramRef paramRef specifies the parameter resource used to configure the admission control policy. It should point to a resource of the type specified in spec.ParamKind of the bound MutatingAdmissionPolicy. If the policy specifies a ParamKind and the resource referred to by ParamRef does not exist, this binding is considered mis-configured and the FailurePolicy of the MutatingAdmissionPolicy applied. If the policy does not specify a ParamKind then this field is ignored, and the rules are evaluated without a param.
         * 
         * @return builder
         * 
         */
        public Builder paramRef(@Nullable Output<ParamRefPatchArgs> paramRef) {
            $.paramRef = paramRef;
            return this;
        }

        /**
         * @param paramRef paramRef specifies the parameter resource used to configure the admission control policy. It should point to a resource of the type specified in spec.ParamKind of the bound MutatingAdmissionPolicy. If the policy specifies a ParamKind and the resource referred to by ParamRef does not exist, this binding is considered mis-configured and the FailurePolicy of the MutatingAdmissionPolicy applied. If the policy does not specify a ParamKind then this field is ignored, and the rules are evaluated without a param.
         * 
         * @return builder
         * 
         */
        public Builder paramRef(ParamRefPatchArgs paramRef) {
            return paramRef(Output.of(paramRef));
        }

        /**
         * @param policyName policyName references a MutatingAdmissionPolicy name which the MutatingAdmissionPolicyBinding binds to. If the referenced resource does not exist, this binding is considered invalid and will be ignored Required.
         * 
         * @return builder
         * 
         */
        public Builder policyName(@Nullable Output<String> policyName) {
            $.policyName = policyName;
            return this;
        }

        /**
         * @param policyName policyName references a MutatingAdmissionPolicy name which the MutatingAdmissionPolicyBinding binds to. If the referenced resource does not exist, this binding is considered invalid and will be ignored Required.
         * 
         * @return builder
         * 
         */
        public Builder policyName(String policyName) {
            return policyName(Output.of(policyName));
        }

        public MutatingAdmissionPolicyBindingSpecPatchArgs build() {
            return $;
        }
    }

}
