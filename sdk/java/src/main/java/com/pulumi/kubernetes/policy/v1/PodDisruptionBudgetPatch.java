// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.policy.v1;

import com.pulumi.core.Alias;
import com.pulumi.core.Output;
import com.pulumi.core.annotations.Export;
import com.pulumi.core.annotations.ResourceType;
import com.pulumi.core.internal.Codegen;
import com.pulumi.kubernetes.Utilities;
import com.pulumi.kubernetes.meta.v1.outputs.ObjectMetaPatch;
import com.pulumi.kubernetes.policy.v1.PodDisruptionBudgetPatchArgs;
import com.pulumi.kubernetes.policy.v1.outputs.PodDisruptionBudgetSpecPatch;
import com.pulumi.kubernetes.policy.v1.outputs.PodDisruptionBudgetStatusPatch;
import java.lang.String;
import java.util.List;
import java.util.Optional;
import javax.annotation.Nullable;

/**
 * Patch resources are used to modify existing Kubernetes resources by using
 * Server-Side Apply updates. The name of the resource must be specified, but all other properties are optional. More than
 * one patch may be applied to the same resource, and a random FieldManager name will be used for each Patch resource.
 * Conflicts will result in an error by default, but can be forced using the &#34;pulumi.com/patchForce&#34; annotation. See the
 * [Server-Side Apply Docs](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/managing-resources-with-server-side-apply/) for
 * additional information about using Server-Side Apply to manage Kubernetes resources with Pulumi.
 * PodDisruptionBudget is an object to define the max disruption that can be caused to a collection of pods
 * 
 */
@ResourceType(type="kubernetes:policy/v1:PodDisruptionBudgetPatch")
public class PodDisruptionBudgetPatch extends com.pulumi.resources.CustomResource {
    /**
     * APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
     * 
     */
    @Export(name="apiVersion", refs={String.class}, tree="[0]")
    private Output</* @Nullable */ String> apiVersion;

    /**
     * @return APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
     * 
     */
    public Output<Optional<String>> apiVersion() {
        return Codegen.optional(this.apiVersion);
    }
    /**
     * Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
     * 
     */
    @Export(name="kind", refs={String.class}, tree="[0]")
    private Output</* @Nullable */ String> kind;

    /**
     * @return Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
     * 
     */
    public Output<Optional<String>> kind() {
        return Codegen.optional(this.kind);
    }
    /**
     * Standard object&#39;s metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
     * 
     */
    @Export(name="metadata", refs={ObjectMetaPatch.class}, tree="[0]")
    private Output</* @Nullable */ ObjectMetaPatch> metadata;

    /**
     * @return Standard object&#39;s metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
     * 
     */
    public Output<Optional<ObjectMetaPatch>> metadata() {
        return Codegen.optional(this.metadata);
    }
    /**
     * Specification of the desired behavior of the PodDisruptionBudget.
     * 
     */
    @Export(name="spec", refs={PodDisruptionBudgetSpecPatch.class}, tree="[0]")
    private Output</* @Nullable */ PodDisruptionBudgetSpecPatch> spec;

    /**
     * @return Specification of the desired behavior of the PodDisruptionBudget.
     * 
     */
    public Output<Optional<PodDisruptionBudgetSpecPatch>> spec() {
        return Codegen.optional(this.spec);
    }
    /**
     * Most recently observed status of the PodDisruptionBudget.
     * 
     */
    @Export(name="status", refs={PodDisruptionBudgetStatusPatch.class}, tree="[0]")
    private Output</* @Nullable */ PodDisruptionBudgetStatusPatch> status;

    /**
     * @return Most recently observed status of the PodDisruptionBudget.
     * 
     */
    public Output<Optional<PodDisruptionBudgetStatusPatch>> status() {
        return Codegen.optional(this.status);
    }

    /**
     *
     * @param name The _unique_ name of the resulting resource.
     */
    public PodDisruptionBudgetPatch(java.lang.String name) {
        this(name, PodDisruptionBudgetPatchArgs.Empty);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     */
    public PodDisruptionBudgetPatch(java.lang.String name, @Nullable PodDisruptionBudgetPatchArgs args) {
        this(name, args, null);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param options A bag of options that control this resource's behavior.
     */
    public PodDisruptionBudgetPatch(java.lang.String name, @Nullable PodDisruptionBudgetPatchArgs args, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super("kubernetes:policy/v1:PodDisruptionBudgetPatch", name, makeArgs(args, options), makeResourceOptions(options, Codegen.empty()), false);
    }

    private PodDisruptionBudgetPatch(java.lang.String name, Output<java.lang.String> id, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super("kubernetes:policy/v1:PodDisruptionBudgetPatch", name, null, makeResourceOptions(options, id), false);
    }

    private static PodDisruptionBudgetPatchArgs makeArgs(@Nullable PodDisruptionBudgetPatchArgs args, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        if (options != null && options.getUrn().isPresent()) {
            return null;
        }
        var builder = args == null ? PodDisruptionBudgetPatchArgs.builder() : PodDisruptionBudgetPatchArgs.builder(args);
        return builder
            .apiVersion("policy/v1")
            .kind("PodDisruptionBudget")
            .build();
    }

    private static com.pulumi.resources.CustomResourceOptions makeResourceOptions(@Nullable com.pulumi.resources.CustomResourceOptions options, @Nullable Output<java.lang.String> id) {
        var defaultOptions = com.pulumi.resources.CustomResourceOptions.builder()
            .version(Utilities.getVersion())
            .aliases(List.of(
                Output.of(Alias.builder().type("kubernetes:policy/v1beta1:PodDisruptionBudgetPatch").build())
            ))
            .build();
        return com.pulumi.resources.CustomResourceOptions.merge(defaultOptions, options, id);
    }

    /**
     * Get an existing Host resource's state with the given name, ID, and optional extra
     * properties used to qualify the lookup.
     *
     * @param name The _unique_ name of the resulting resource.
     * @param id The _unique_ provider ID of the resource to lookup.
     * @param options Optional settings to control the behavior of the CustomResource.
     */
    public static PodDisruptionBudgetPatch get(java.lang.String name, Output<java.lang.String> id, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        return new PodDisruptionBudgetPatch(name, id, options);
    }
}
