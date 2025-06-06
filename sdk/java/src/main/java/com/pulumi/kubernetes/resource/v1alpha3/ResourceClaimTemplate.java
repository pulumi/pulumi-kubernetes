// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1alpha3;

import com.pulumi.core.Alias;
import com.pulumi.core.Output;
import com.pulumi.core.annotations.Export;
import com.pulumi.core.annotations.ResourceType;
import com.pulumi.core.internal.Codegen;
import com.pulumi.kubernetes.Utilities;
import com.pulumi.kubernetes.meta.v1.outputs.ObjectMeta;
import com.pulumi.kubernetes.resource.v1alpha3.ResourceClaimTemplateArgs;
import com.pulumi.kubernetes.resource.v1alpha3.outputs.ResourceClaimTemplateSpec;
import java.lang.String;
import java.util.List;
import javax.annotation.Nullable;

/**
 * ResourceClaimTemplate is used to produce ResourceClaim objects.
 * 
 * This is an alpha type and requires enabling the DynamicResourceAllocation feature gate.
 * 
 */
@ResourceType(type="kubernetes:resource.k8s.io/v1alpha3:ResourceClaimTemplate")
public class ResourceClaimTemplate extends com.pulumi.resources.CustomResource {
    /**
     * APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
     * 
     */
    @Export(name="apiVersion", refs={String.class}, tree="[0]")
    private Output<String> apiVersion;

    /**
     * @return APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
     * 
     */
    public Output<String> apiVersion() {
        return this.apiVersion;
    }
    /**
     * Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
     * 
     */
    @Export(name="kind", refs={String.class}, tree="[0]")
    private Output<String> kind;

    /**
     * @return Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
     * 
     */
    public Output<String> kind() {
        return this.kind;
    }
    /**
     * Standard object metadata
     * 
     */
    @Export(name="metadata", refs={ObjectMeta.class}, tree="[0]")
    private Output<ObjectMeta> metadata;

    /**
     * @return Standard object metadata
     * 
     */
    public Output<ObjectMeta> metadata() {
        return this.metadata;
    }
    /**
     * Describes the ResourceClaim that is to be generated.
     * 
     * This field is immutable. A ResourceClaim will get created by the control plane for a Pod when needed and then not get updated anymore.
     * 
     */
    @Export(name="spec", refs={ResourceClaimTemplateSpec.class}, tree="[0]")
    private Output<ResourceClaimTemplateSpec> spec;

    /**
     * @return Describes the ResourceClaim that is to be generated.
     * 
     * This field is immutable. A ResourceClaim will get created by the control plane for a Pod when needed and then not get updated anymore.
     * 
     */
    public Output<ResourceClaimTemplateSpec> spec() {
        return this.spec;
    }

    /**
     *
     * @param name The _unique_ name of the resulting resource.
     */
    public ResourceClaimTemplate(java.lang.String name) {
        this(name, ResourceClaimTemplateArgs.Empty);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     */
    public ResourceClaimTemplate(java.lang.String name, ResourceClaimTemplateArgs args) {
        this(name, args, null);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param options A bag of options that control this resource's behavior.
     */
    public ResourceClaimTemplate(java.lang.String name, ResourceClaimTemplateArgs args, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super("kubernetes:resource.k8s.io/v1alpha3:ResourceClaimTemplate", name, makeArgs(args, options), makeResourceOptions(options, Codegen.empty()), false);
    }

    private ResourceClaimTemplate(java.lang.String name, Output<java.lang.String> id, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super("kubernetes:resource.k8s.io/v1alpha3:ResourceClaimTemplate", name, null, makeResourceOptions(options, id), false);
    }

    private static ResourceClaimTemplateArgs makeArgs(ResourceClaimTemplateArgs args, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        if (options != null && options.getUrn().isPresent()) {
            return null;
        }
        var builder = args == null ? ResourceClaimTemplateArgs.builder() : ResourceClaimTemplateArgs.builder(args);
        return builder
            .apiVersion("resource.k8s.io/v1alpha3")
            .kind("ResourceClaimTemplate")
            .build();
    }

    private static com.pulumi.resources.CustomResourceOptions makeResourceOptions(@Nullable com.pulumi.resources.CustomResourceOptions options, @Nullable Output<java.lang.String> id) {
        var defaultOptions = com.pulumi.resources.CustomResourceOptions.builder()
            .version(Utilities.getVersion())
            .aliases(List.of(
                Output.of(Alias.builder().type("kubernetes:resource.k8s.io/v1alpha1:ResourceClaimTemplate").build()),
                Output.of(Alias.builder().type("kubernetes:resource.k8s.io/v1alpha2:ResourceClaimTemplate").build()),
                Output.of(Alias.builder().type("kubernetes:resource.k8s.io/v1beta1:ResourceClaimTemplate").build()),
                Output.of(Alias.builder().type("kubernetes:resource.k8s.io/v1beta2:ResourceClaimTemplate").build())
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
    public static ResourceClaimTemplate get(java.lang.String name, Output<java.lang.String> id, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        return new ResourceClaimTemplate(name, id, options);
    }
}
