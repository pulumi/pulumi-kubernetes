// Code generated by pulumigen DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package v1alpha3

import (
	"context"
	"reflect"

	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/utilities"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Patch resources are used to modify existing Kubernetes resources by using
// Server-Side Apply updates. The name of the resource must be specified, but all other properties are optional. More than
// one patch may be applied to the same resource, and a random FieldManager name will be used for each Patch resource.
// Conflicts will result in an error by default, but can be forced using the "pulumi.com/patchForce" annotation. See the
// [Server-Side Apply Docs](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/managing-resources-with-server-side-apply/) for
// additional information about using Server-Side Apply to manage Kubernetes resources with Pulumi.
// ResourceSlice represents one or more resources in a pool of similar resources, managed by a common driver. A pool may span more than one ResourceSlice, and exactly how many ResourceSlices comprise a pool is determined by the driver.
//
// At the moment, the only supported resources are devices with attributes and capacities. Each device in a given pool, regardless of how many ResourceSlices, must have a unique name. The ResourceSlice in which a device gets published may change over time. The unique identifier for a device is the tuple <driver name>, <pool name>, <device name>.
//
// Whenever a driver needs to update a pool, it increments the pool.Spec.Pool.Generation number and updates all ResourceSlices with that new number and new resource definitions. A consumer must only use ResourceSlices with the highest generation number and ignore all others.
//
// When allocating all resources in a pool matching certain criteria or when looking for the best solution among several different alternatives, a consumer should check the number of ResourceSlices in a pool (included in each ResourceSlice) to determine whether its view of a pool is complete and if not, should wait until the driver has completed updating the pool.
//
// For resources that are not local to a node, the node name is not set. Instead, the driver may use a node selector to specify where the devices are available.
//
// This is an alpha type and requires enabling the DynamicResourceAllocation feature gate.
type ResourceSlicePatch struct {
	pulumi.CustomResourceState

	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion pulumi.StringPtrOutput `pulumi:"apiVersion"`
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind pulumi.StringPtrOutput `pulumi:"kind"`
	// Standard object metadata
	Metadata metav1.ObjectMetaPatchPtrOutput `pulumi:"metadata"`
	// Contains the information published by the driver.
	//
	// Changing the spec automatically increments the metadata.generation number.
	Spec ResourceSliceSpecPatchPtrOutput `pulumi:"spec"`
}

// NewResourceSlicePatch registers a new resource with the given unique name, arguments, and options.
func NewResourceSlicePatch(ctx *pulumi.Context,
	name string, args *ResourceSlicePatchArgs, opts ...pulumi.ResourceOption) (*ResourceSlicePatch, error) {
	if args == nil {
		args = &ResourceSlicePatchArgs{}
	}

	args.ApiVersion = pulumi.StringPtr("resource.k8s.io/v1alpha3")
	args.Kind = pulumi.StringPtr("ResourceSlice")
	aliases := pulumi.Aliases([]pulumi.Alias{
		{
			Type: pulumi.String("kubernetes:resource.k8s.io/v1alpha2:ResourceSlicePatch"),
		},
	})
	opts = append(opts, aliases)
	opts = utilities.PkgResourceDefaultOpts(opts)
	var resource ResourceSlicePatch
	err := ctx.RegisterResource("kubernetes:resource.k8s.io/v1alpha3:ResourceSlicePatch", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetResourceSlicePatch gets an existing ResourceSlicePatch resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetResourceSlicePatch(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *ResourceSlicePatchState, opts ...pulumi.ResourceOption) (*ResourceSlicePatch, error) {
	var resource ResourceSlicePatch
	err := ctx.ReadResource("kubernetes:resource.k8s.io/v1alpha3:ResourceSlicePatch", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering ResourceSlicePatch resources.
type resourceSlicePatchState struct {
}

type ResourceSlicePatchState struct {
}

func (ResourceSlicePatchState) ElementType() reflect.Type {
	return reflect.TypeOf((*resourceSlicePatchState)(nil)).Elem()
}

type resourceSlicePatchArgs struct {
	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion *string `pulumi:"apiVersion"`
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind *string `pulumi:"kind"`
	// Standard object metadata
	Metadata *metav1.ObjectMetaPatch `pulumi:"metadata"`
	// Contains the information published by the driver.
	//
	// Changing the spec automatically increments the metadata.generation number.
	Spec *ResourceSliceSpecPatch `pulumi:"spec"`
}

// The set of arguments for constructing a ResourceSlicePatch resource.
type ResourceSlicePatchArgs struct {
	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion pulumi.StringPtrInput
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind pulumi.StringPtrInput
	// Standard object metadata
	Metadata metav1.ObjectMetaPatchPtrInput
	// Contains the information published by the driver.
	//
	// Changing the spec automatically increments the metadata.generation number.
	Spec ResourceSliceSpecPatchPtrInput
}

func (ResourceSlicePatchArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*resourceSlicePatchArgs)(nil)).Elem()
}

type ResourceSlicePatchInput interface {
	pulumi.Input

	ToResourceSlicePatchOutput() ResourceSlicePatchOutput
	ToResourceSlicePatchOutputWithContext(ctx context.Context) ResourceSlicePatchOutput
}

func (*ResourceSlicePatch) ElementType() reflect.Type {
	return reflect.TypeOf((**ResourceSlicePatch)(nil)).Elem()
}

func (i *ResourceSlicePatch) ToResourceSlicePatchOutput() ResourceSlicePatchOutput {
	return i.ToResourceSlicePatchOutputWithContext(context.Background())
}

func (i *ResourceSlicePatch) ToResourceSlicePatchOutputWithContext(ctx context.Context) ResourceSlicePatchOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ResourceSlicePatchOutput)
}

// ResourceSlicePatchArrayInput is an input type that accepts ResourceSlicePatchArray and ResourceSlicePatchArrayOutput values.
// You can construct a concrete instance of `ResourceSlicePatchArrayInput` via:
//
//	ResourceSlicePatchArray{ ResourceSlicePatchArgs{...} }
type ResourceSlicePatchArrayInput interface {
	pulumi.Input

	ToResourceSlicePatchArrayOutput() ResourceSlicePatchArrayOutput
	ToResourceSlicePatchArrayOutputWithContext(context.Context) ResourceSlicePatchArrayOutput
}

type ResourceSlicePatchArray []ResourceSlicePatchInput

func (ResourceSlicePatchArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*ResourceSlicePatch)(nil)).Elem()
}

func (i ResourceSlicePatchArray) ToResourceSlicePatchArrayOutput() ResourceSlicePatchArrayOutput {
	return i.ToResourceSlicePatchArrayOutputWithContext(context.Background())
}

func (i ResourceSlicePatchArray) ToResourceSlicePatchArrayOutputWithContext(ctx context.Context) ResourceSlicePatchArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ResourceSlicePatchArrayOutput)
}

// ResourceSlicePatchMapInput is an input type that accepts ResourceSlicePatchMap and ResourceSlicePatchMapOutput values.
// You can construct a concrete instance of `ResourceSlicePatchMapInput` via:
//
//	ResourceSlicePatchMap{ "key": ResourceSlicePatchArgs{...} }
type ResourceSlicePatchMapInput interface {
	pulumi.Input

	ToResourceSlicePatchMapOutput() ResourceSlicePatchMapOutput
	ToResourceSlicePatchMapOutputWithContext(context.Context) ResourceSlicePatchMapOutput
}

type ResourceSlicePatchMap map[string]ResourceSlicePatchInput

func (ResourceSlicePatchMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*ResourceSlicePatch)(nil)).Elem()
}

func (i ResourceSlicePatchMap) ToResourceSlicePatchMapOutput() ResourceSlicePatchMapOutput {
	return i.ToResourceSlicePatchMapOutputWithContext(context.Background())
}

func (i ResourceSlicePatchMap) ToResourceSlicePatchMapOutputWithContext(ctx context.Context) ResourceSlicePatchMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ResourceSlicePatchMapOutput)
}

type ResourceSlicePatchOutput struct{ *pulumi.OutputState }

func (ResourceSlicePatchOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**ResourceSlicePatch)(nil)).Elem()
}

func (o ResourceSlicePatchOutput) ToResourceSlicePatchOutput() ResourceSlicePatchOutput {
	return o
}

func (o ResourceSlicePatchOutput) ToResourceSlicePatchOutputWithContext(ctx context.Context) ResourceSlicePatchOutput {
	return o
}

// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (o ResourceSlicePatchOutput) ApiVersion() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *ResourceSlicePatch) pulumi.StringPtrOutput { return v.ApiVersion }).(pulumi.StringPtrOutput)
}

// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (o ResourceSlicePatchOutput) Kind() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *ResourceSlicePatch) pulumi.StringPtrOutput { return v.Kind }).(pulumi.StringPtrOutput)
}

// Standard object metadata
func (o ResourceSlicePatchOutput) Metadata() metav1.ObjectMetaPatchPtrOutput {
	return o.ApplyT(func(v *ResourceSlicePatch) metav1.ObjectMetaPatchPtrOutput { return v.Metadata }).(metav1.ObjectMetaPatchPtrOutput)
}

// Contains the information published by the driver.
//
// Changing the spec automatically increments the metadata.generation number.
func (o ResourceSlicePatchOutput) Spec() ResourceSliceSpecPatchPtrOutput {
	return o.ApplyT(func(v *ResourceSlicePatch) ResourceSliceSpecPatchPtrOutput { return v.Spec }).(ResourceSliceSpecPatchPtrOutput)
}

type ResourceSlicePatchArrayOutput struct{ *pulumi.OutputState }

func (ResourceSlicePatchArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*ResourceSlicePatch)(nil)).Elem()
}

func (o ResourceSlicePatchArrayOutput) ToResourceSlicePatchArrayOutput() ResourceSlicePatchArrayOutput {
	return o
}

func (o ResourceSlicePatchArrayOutput) ToResourceSlicePatchArrayOutputWithContext(ctx context.Context) ResourceSlicePatchArrayOutput {
	return o
}

func (o ResourceSlicePatchArrayOutput) Index(i pulumi.IntInput) ResourceSlicePatchOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *ResourceSlicePatch {
		return vs[0].([]*ResourceSlicePatch)[vs[1].(int)]
	}).(ResourceSlicePatchOutput)
}

type ResourceSlicePatchMapOutput struct{ *pulumi.OutputState }

func (ResourceSlicePatchMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*ResourceSlicePatch)(nil)).Elem()
}

func (o ResourceSlicePatchMapOutput) ToResourceSlicePatchMapOutput() ResourceSlicePatchMapOutput {
	return o
}

func (o ResourceSlicePatchMapOutput) ToResourceSlicePatchMapOutputWithContext(ctx context.Context) ResourceSlicePatchMapOutput {
	return o
}

func (o ResourceSlicePatchMapOutput) MapIndex(k pulumi.StringInput) ResourceSlicePatchOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *ResourceSlicePatch {
		return vs[0].(map[string]*ResourceSlicePatch)[vs[1].(string)]
	}).(ResourceSlicePatchOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*ResourceSlicePatchInput)(nil)).Elem(), &ResourceSlicePatch{})
	pulumi.RegisterInputType(reflect.TypeOf((*ResourceSlicePatchArrayInput)(nil)).Elem(), ResourceSlicePatchArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*ResourceSlicePatchMapInput)(nil)).Elem(), ResourceSlicePatchMap{})
	pulumi.RegisterOutputType(ResourceSlicePatchOutput{})
	pulumi.RegisterOutputType(ResourceSlicePatchArrayOutput{})
	pulumi.RegisterOutputType(ResourceSlicePatchMapOutput{})
}