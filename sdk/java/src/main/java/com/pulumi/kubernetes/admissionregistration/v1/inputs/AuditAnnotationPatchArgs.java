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
 * AuditAnnotation describes how to produce an audit annotation for an API request.
 * 
 */
public final class AuditAnnotationPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final AuditAnnotationPatchArgs Empty = new AuditAnnotationPatchArgs();

    /**
     * key specifies the audit annotation key. The audit annotation keys of a ValidatingAdmissionPolicy must be unique. The key must be a qualified name ([A-Za-z0-9][-A-Za-z0-9_.]*) no more than 63 bytes in length.
     * 
     * The key is combined with the resource name of the ValidatingAdmissionPolicy to construct an audit annotation key: &#34;{ValidatingAdmissionPolicy name}/{key}&#34;.
     * 
     * If an admission webhook uses the same resource name as this ValidatingAdmissionPolicy and the same audit annotation key, the annotation key will be identical. In this case, the first annotation written with the key will be included in the audit event and all subsequent annotations with the same key will be discarded.
     * 
     * Required.
     * 
     */
    @Import(name="key")
    private @Nullable Output<String> key;

    /**
     * @return key specifies the audit annotation key. The audit annotation keys of a ValidatingAdmissionPolicy must be unique. The key must be a qualified name ([A-Za-z0-9][-A-Za-z0-9_.]*) no more than 63 bytes in length.
     * 
     * The key is combined with the resource name of the ValidatingAdmissionPolicy to construct an audit annotation key: &#34;{ValidatingAdmissionPolicy name}/{key}&#34;.
     * 
     * If an admission webhook uses the same resource name as this ValidatingAdmissionPolicy and the same audit annotation key, the annotation key will be identical. In this case, the first annotation written with the key will be included in the audit event and all subsequent annotations with the same key will be discarded.
     * 
     * Required.
     * 
     */
    public Optional<Output<String>> key() {
        return Optional.ofNullable(this.key);
    }

    /**
     * valueExpression represents the expression which is evaluated by CEL to produce an audit annotation value. The expression must evaluate to either a string or null value. If the expression evaluates to a string, the audit annotation is included with the string value. If the expression evaluates to null or empty string the audit annotation will be omitted. The valueExpression may be no longer than 5kb in length. If the result of the valueExpression is more than 10kb in length, it will be truncated to 10kb.
     * 
     * If multiple ValidatingAdmissionPolicyBinding resources match an API request, then the valueExpression will be evaluated for each binding. All unique values produced by the valueExpressions will be joined together in a comma-separated list.
     * 
     * Required.
     * 
     */
    @Import(name="valueExpression")
    private @Nullable Output<String> valueExpression;

    /**
     * @return valueExpression represents the expression which is evaluated by CEL to produce an audit annotation value. The expression must evaluate to either a string or null value. If the expression evaluates to a string, the audit annotation is included with the string value. If the expression evaluates to null or empty string the audit annotation will be omitted. The valueExpression may be no longer than 5kb in length. If the result of the valueExpression is more than 10kb in length, it will be truncated to 10kb.
     * 
     * If multiple ValidatingAdmissionPolicyBinding resources match an API request, then the valueExpression will be evaluated for each binding. All unique values produced by the valueExpressions will be joined together in a comma-separated list.
     * 
     * Required.
     * 
     */
    public Optional<Output<String>> valueExpression() {
        return Optional.ofNullable(this.valueExpression);
    }

    private AuditAnnotationPatchArgs() {}

    private AuditAnnotationPatchArgs(AuditAnnotationPatchArgs $) {
        this.key = $.key;
        this.valueExpression = $.valueExpression;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(AuditAnnotationPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private AuditAnnotationPatchArgs $;

        public Builder() {
            $ = new AuditAnnotationPatchArgs();
        }

        public Builder(AuditAnnotationPatchArgs defaults) {
            $ = new AuditAnnotationPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param key key specifies the audit annotation key. The audit annotation keys of a ValidatingAdmissionPolicy must be unique. The key must be a qualified name ([A-Za-z0-9][-A-Za-z0-9_.]*) no more than 63 bytes in length.
         * 
         * The key is combined with the resource name of the ValidatingAdmissionPolicy to construct an audit annotation key: &#34;{ValidatingAdmissionPolicy name}/{key}&#34;.
         * 
         * If an admission webhook uses the same resource name as this ValidatingAdmissionPolicy and the same audit annotation key, the annotation key will be identical. In this case, the first annotation written with the key will be included in the audit event and all subsequent annotations with the same key will be discarded.
         * 
         * Required.
         * 
         * @return builder
         * 
         */
        public Builder key(@Nullable Output<String> key) {
            $.key = key;
            return this;
        }

        /**
         * @param key key specifies the audit annotation key. The audit annotation keys of a ValidatingAdmissionPolicy must be unique. The key must be a qualified name ([A-Za-z0-9][-A-Za-z0-9_.]*) no more than 63 bytes in length.
         * 
         * The key is combined with the resource name of the ValidatingAdmissionPolicy to construct an audit annotation key: &#34;{ValidatingAdmissionPolicy name}/{key}&#34;.
         * 
         * If an admission webhook uses the same resource name as this ValidatingAdmissionPolicy and the same audit annotation key, the annotation key will be identical. In this case, the first annotation written with the key will be included in the audit event and all subsequent annotations with the same key will be discarded.
         * 
         * Required.
         * 
         * @return builder
         * 
         */
        public Builder key(String key) {
            return key(Output.of(key));
        }

        /**
         * @param valueExpression valueExpression represents the expression which is evaluated by CEL to produce an audit annotation value. The expression must evaluate to either a string or null value. If the expression evaluates to a string, the audit annotation is included with the string value. If the expression evaluates to null or empty string the audit annotation will be omitted. The valueExpression may be no longer than 5kb in length. If the result of the valueExpression is more than 10kb in length, it will be truncated to 10kb.
         * 
         * If multiple ValidatingAdmissionPolicyBinding resources match an API request, then the valueExpression will be evaluated for each binding. All unique values produced by the valueExpressions will be joined together in a comma-separated list.
         * 
         * Required.
         * 
         * @return builder
         * 
         */
        public Builder valueExpression(@Nullable Output<String> valueExpression) {
            $.valueExpression = valueExpression;
            return this;
        }

        /**
         * @param valueExpression valueExpression represents the expression which is evaluated by CEL to produce an audit annotation value. The expression must evaluate to either a string or null value. If the expression evaluates to a string, the audit annotation is included with the string value. If the expression evaluates to null or empty string the audit annotation will be omitted. The valueExpression may be no longer than 5kb in length. If the result of the valueExpression is more than 10kb in length, it will be truncated to 10kb.
         * 
         * If multiple ValidatingAdmissionPolicyBinding resources match an API request, then the valueExpression will be evaluated for each binding. All unique values produced by the valueExpressions will be joined together in a comma-separated list.
         * 
         * Required.
         * 
         * @return builder
         * 
         */
        public Builder valueExpression(String valueExpression) {
            return valueExpression(Output.of(valueExpression));
        }

        public AuditAnnotationPatchArgs build() {
            return $;
        }
    }

}
