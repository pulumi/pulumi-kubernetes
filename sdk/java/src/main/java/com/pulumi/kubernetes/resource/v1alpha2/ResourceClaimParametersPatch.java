// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1alpha2;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Export;
import com.pulumi.core.annotations.ResourceType;
import com.pulumi.core.internal.Codegen;
import com.pulumi.kubernetes.Utilities;
import com.pulumi.kubernetes.meta.v1.outputs.ObjectMetaPatch;
import com.pulumi.kubernetes.resource.v1alpha2.ResourceClaimParametersPatchArgs;
import com.pulumi.kubernetes.resource.v1alpha2.outputs.DriverRequestsPatch;
import com.pulumi.kubernetes.resource.v1alpha2.outputs.ResourceClaimParametersReferencePatch;
import java.lang.Boolean;
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
 * ResourceClaimParameters defines resource requests for a ResourceClaim in an in-tree format understood by Kubernetes.
 * 
 */
@ResourceType(type="kubernetes:resource.k8s.io/v1alpha2:ResourceClaimParametersPatch")
public class ResourceClaimParametersPatch extends com.pulumi.resources.CustomResource {
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
     * DriverRequests describes all resources that are needed for the allocated claim. A single claim may use resources coming from different drivers. For each driver, this array has at most one entry which then may have one or more per-driver requests.
     * 
     * May be empty, in which case the claim can always be allocated.
     * 
     */
    @Export(name="driverRequests", refs={List.class,DriverRequestsPatch.class}, tree="[0,1]")
    private Output</* @Nullable */ List<DriverRequestsPatch>> driverRequests;

    /**
     * @return DriverRequests describes all resources that are needed for the allocated claim. A single claim may use resources coming from different drivers. For each driver, this array has at most one entry which then may have one or more per-driver requests.
     * 
     * May be empty, in which case the claim can always be allocated.
     * 
     */
    public Output<Optional<List<DriverRequestsPatch>>> driverRequests() {
        return Codegen.optional(this.driverRequests);
    }
    /**
     * If this object was created from some other resource, then this links back to that resource. This field is used to find the in-tree representation of the claim parameters when the parameter reference of the claim refers to some unknown type.
     * 
     */
    @Export(name="generatedFrom", refs={ResourceClaimParametersReferencePatch.class}, tree="[0]")
    private Output</* @Nullable */ ResourceClaimParametersReferencePatch> generatedFrom;

    /**
     * @return If this object was created from some other resource, then this links back to that resource. This field is used to find the in-tree representation of the claim parameters when the parameter reference of the claim refers to some unknown type.
     * 
     */
    public Output<Optional<ResourceClaimParametersReferencePatch>> generatedFrom() {
        return Codegen.optional(this.generatedFrom);
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
     * Standard object metadata
     * 
     */
    @Export(name="metadata", refs={ObjectMetaPatch.class}, tree="[0]")
    private Output</* @Nullable */ ObjectMetaPatch> metadata;

    /**
     * @return Standard object metadata
     * 
     */
    public Output<Optional<ObjectMetaPatch>> metadata() {
        return Codegen.optional(this.metadata);
    }
    /**
     * Shareable indicates whether the allocated claim is meant to be shareable by multiple consumers at the same time.
     * 
     */
    @Export(name="shareable", refs={Boolean.class}, tree="[0]")
    private Output</* @Nullable */ Boolean> shareable;

    /**
     * @return Shareable indicates whether the allocated claim is meant to be shareable by multiple consumers at the same time.
     * 
     */
    public Output<Optional<Boolean>> shareable() {
        return Codegen.optional(this.shareable);
    }

    /**
     *
     * @param name The _unique_ name of the resulting resource.
     */
    public ResourceClaimParametersPatch(java.lang.String name) {
        this(name, ResourceClaimParametersPatchArgs.Empty);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     */
    public ResourceClaimParametersPatch(java.lang.String name, @Nullable ResourceClaimParametersPatchArgs args) {
        this(name, args, null);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param options A bag of options that control this resource's behavior.
     */
    public ResourceClaimParametersPatch(java.lang.String name, @Nullable ResourceClaimParametersPatchArgs args, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super("kubernetes:resource.k8s.io/v1alpha2:ResourceClaimParametersPatch", name, makeArgs(args, options), makeResourceOptions(options, Codegen.empty()), false);
    }

    private ResourceClaimParametersPatch(java.lang.String name, Output<java.lang.String> id, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super("kubernetes:resource.k8s.io/v1alpha2:ResourceClaimParametersPatch", name, null, makeResourceOptions(options, id), false);
    }

    private static ResourceClaimParametersPatchArgs makeArgs(@Nullable ResourceClaimParametersPatchArgs args, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        if (options != null && options.getUrn().isPresent()) {
            return null;
        }
        var builder = args == null ? ResourceClaimParametersPatchArgs.builder() : ResourceClaimParametersPatchArgs.builder(args);
        return builder
            .apiVersion("resource.k8s.io/v1alpha2")
            .kind("ResourceClaimParameters")
            .build();
    }

    private static com.pulumi.resources.CustomResourceOptions makeResourceOptions(@Nullable com.pulumi.resources.CustomResourceOptions options, @Nullable Output<java.lang.String> id) {
        var defaultOptions = com.pulumi.resources.CustomResourceOptions.builder()
            .version(Utilities.getVersion())
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
    public static ResourceClaimParametersPatch get(java.lang.String name, Output<java.lang.String> id, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        return new ResourceClaimParametersPatch(name, id, options);
    }
}
