// Code generated by pulumigen DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package v1alpha1

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
// StorageVersionMigration represents a migration of stored data to the latest storage version.
type StorageVersionMigrationPatch struct {
	pulumi.CustomResourceState

	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion pulumi.StringPtrOutput `pulumi:"apiVersion"`
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind pulumi.StringPtrOutput `pulumi:"kind"`
	// Standard object metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata metav1.ObjectMetaPatchPtrOutput `pulumi:"metadata"`
	// Specification of the migration.
	Spec StorageVersionMigrationSpecPatchPtrOutput `pulumi:"spec"`
	// Status of the migration.
	Status StorageVersionMigrationStatusPatchPtrOutput `pulumi:"status"`
}

// NewStorageVersionMigrationPatch registers a new resource with the given unique name, arguments, and options.
func NewStorageVersionMigrationPatch(ctx *pulumi.Context,
	name string, args *StorageVersionMigrationPatchArgs, opts ...pulumi.ResourceOption) (*StorageVersionMigrationPatch, error) {
	if args == nil {
		args = &StorageVersionMigrationPatchArgs{}
	}

	args.ApiVersion = pulumi.StringPtr("storagemigration.k8s.io/v1alpha1")
	args.Kind = pulumi.StringPtr("StorageVersionMigration")
	opts = utilities.PkgResourceDefaultOpts(opts)
	var resource StorageVersionMigrationPatch
	err := ctx.RegisterResource("kubernetes:storagemigration.k8s.io/v1alpha1:StorageVersionMigrationPatch", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetStorageVersionMigrationPatch gets an existing StorageVersionMigrationPatch resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetStorageVersionMigrationPatch(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *StorageVersionMigrationPatchState, opts ...pulumi.ResourceOption) (*StorageVersionMigrationPatch, error) {
	var resource StorageVersionMigrationPatch
	err := ctx.ReadResource("kubernetes:storagemigration.k8s.io/v1alpha1:StorageVersionMigrationPatch", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering StorageVersionMigrationPatch resources.
type storageVersionMigrationPatchState struct {
}

type StorageVersionMigrationPatchState struct {
}

func (StorageVersionMigrationPatchState) ElementType() reflect.Type {
	return reflect.TypeOf((*storageVersionMigrationPatchState)(nil)).Elem()
}

type storageVersionMigrationPatchArgs struct {
	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion *string `pulumi:"apiVersion"`
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind *string `pulumi:"kind"`
	// Standard object metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata *metav1.ObjectMetaPatch `pulumi:"metadata"`
	// Specification of the migration.
	Spec *StorageVersionMigrationSpecPatch `pulumi:"spec"`
}

// The set of arguments for constructing a StorageVersionMigrationPatch resource.
type StorageVersionMigrationPatchArgs struct {
	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion pulumi.StringPtrInput
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind pulumi.StringPtrInput
	// Standard object metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata metav1.ObjectMetaPatchPtrInput
	// Specification of the migration.
	Spec StorageVersionMigrationSpecPatchPtrInput
}

func (StorageVersionMigrationPatchArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*storageVersionMigrationPatchArgs)(nil)).Elem()
}

type StorageVersionMigrationPatchInput interface {
	pulumi.Input

	ToStorageVersionMigrationPatchOutput() StorageVersionMigrationPatchOutput
	ToStorageVersionMigrationPatchOutputWithContext(ctx context.Context) StorageVersionMigrationPatchOutput
}

func (*StorageVersionMigrationPatch) ElementType() reflect.Type {
	return reflect.TypeOf((**StorageVersionMigrationPatch)(nil)).Elem()
}

func (i *StorageVersionMigrationPatch) ToStorageVersionMigrationPatchOutput() StorageVersionMigrationPatchOutput {
	return i.ToStorageVersionMigrationPatchOutputWithContext(context.Background())
}

func (i *StorageVersionMigrationPatch) ToStorageVersionMigrationPatchOutputWithContext(ctx context.Context) StorageVersionMigrationPatchOutput {
	return pulumi.ToOutputWithContext(ctx, i).(StorageVersionMigrationPatchOutput)
}

// StorageVersionMigrationPatchArrayInput is an input type that accepts StorageVersionMigrationPatchArray and StorageVersionMigrationPatchArrayOutput values.
// You can construct a concrete instance of `StorageVersionMigrationPatchArrayInput` via:
//
//	StorageVersionMigrationPatchArray{ StorageVersionMigrationPatchArgs{...} }
type StorageVersionMigrationPatchArrayInput interface {
	pulumi.Input

	ToStorageVersionMigrationPatchArrayOutput() StorageVersionMigrationPatchArrayOutput
	ToStorageVersionMigrationPatchArrayOutputWithContext(context.Context) StorageVersionMigrationPatchArrayOutput
}

type StorageVersionMigrationPatchArray []StorageVersionMigrationPatchInput

func (StorageVersionMigrationPatchArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*StorageVersionMigrationPatch)(nil)).Elem()
}

func (i StorageVersionMigrationPatchArray) ToStorageVersionMigrationPatchArrayOutput() StorageVersionMigrationPatchArrayOutput {
	return i.ToStorageVersionMigrationPatchArrayOutputWithContext(context.Background())
}

func (i StorageVersionMigrationPatchArray) ToStorageVersionMigrationPatchArrayOutputWithContext(ctx context.Context) StorageVersionMigrationPatchArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(StorageVersionMigrationPatchArrayOutput)
}

// StorageVersionMigrationPatchMapInput is an input type that accepts StorageVersionMigrationPatchMap and StorageVersionMigrationPatchMapOutput values.
// You can construct a concrete instance of `StorageVersionMigrationPatchMapInput` via:
//
//	StorageVersionMigrationPatchMap{ "key": StorageVersionMigrationPatchArgs{...} }
type StorageVersionMigrationPatchMapInput interface {
	pulumi.Input

	ToStorageVersionMigrationPatchMapOutput() StorageVersionMigrationPatchMapOutput
	ToStorageVersionMigrationPatchMapOutputWithContext(context.Context) StorageVersionMigrationPatchMapOutput
}

type StorageVersionMigrationPatchMap map[string]StorageVersionMigrationPatchInput

func (StorageVersionMigrationPatchMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*StorageVersionMigrationPatch)(nil)).Elem()
}

func (i StorageVersionMigrationPatchMap) ToStorageVersionMigrationPatchMapOutput() StorageVersionMigrationPatchMapOutput {
	return i.ToStorageVersionMigrationPatchMapOutputWithContext(context.Background())
}

func (i StorageVersionMigrationPatchMap) ToStorageVersionMigrationPatchMapOutputWithContext(ctx context.Context) StorageVersionMigrationPatchMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(StorageVersionMigrationPatchMapOutput)
}

type StorageVersionMigrationPatchOutput struct{ *pulumi.OutputState }

func (StorageVersionMigrationPatchOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**StorageVersionMigrationPatch)(nil)).Elem()
}

func (o StorageVersionMigrationPatchOutput) ToStorageVersionMigrationPatchOutput() StorageVersionMigrationPatchOutput {
	return o
}

func (o StorageVersionMigrationPatchOutput) ToStorageVersionMigrationPatchOutputWithContext(ctx context.Context) StorageVersionMigrationPatchOutput {
	return o
}

// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (o StorageVersionMigrationPatchOutput) ApiVersion() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *StorageVersionMigrationPatch) pulumi.StringPtrOutput { return v.ApiVersion }).(pulumi.StringPtrOutput)
}

// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (o StorageVersionMigrationPatchOutput) Kind() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *StorageVersionMigrationPatch) pulumi.StringPtrOutput { return v.Kind }).(pulumi.StringPtrOutput)
}

// Standard object metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
func (o StorageVersionMigrationPatchOutput) Metadata() metav1.ObjectMetaPatchPtrOutput {
	return o.ApplyT(func(v *StorageVersionMigrationPatch) metav1.ObjectMetaPatchPtrOutput { return v.Metadata }).(metav1.ObjectMetaPatchPtrOutput)
}

// Specification of the migration.
func (o StorageVersionMigrationPatchOutput) Spec() StorageVersionMigrationSpecPatchPtrOutput {
	return o.ApplyT(func(v *StorageVersionMigrationPatch) StorageVersionMigrationSpecPatchPtrOutput { return v.Spec }).(StorageVersionMigrationSpecPatchPtrOutput)
}

// Status of the migration.
func (o StorageVersionMigrationPatchOutput) Status() StorageVersionMigrationStatusPatchPtrOutput {
	return o.ApplyT(func(v *StorageVersionMigrationPatch) StorageVersionMigrationStatusPatchPtrOutput { return v.Status }).(StorageVersionMigrationStatusPatchPtrOutput)
}

type StorageVersionMigrationPatchArrayOutput struct{ *pulumi.OutputState }

func (StorageVersionMigrationPatchArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*StorageVersionMigrationPatch)(nil)).Elem()
}

func (o StorageVersionMigrationPatchArrayOutput) ToStorageVersionMigrationPatchArrayOutput() StorageVersionMigrationPatchArrayOutput {
	return o
}

func (o StorageVersionMigrationPatchArrayOutput) ToStorageVersionMigrationPatchArrayOutputWithContext(ctx context.Context) StorageVersionMigrationPatchArrayOutput {
	return o
}

func (o StorageVersionMigrationPatchArrayOutput) Index(i pulumi.IntInput) StorageVersionMigrationPatchOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *StorageVersionMigrationPatch {
		return vs[0].([]*StorageVersionMigrationPatch)[vs[1].(int)]
	}).(StorageVersionMigrationPatchOutput)
}

type StorageVersionMigrationPatchMapOutput struct{ *pulumi.OutputState }

func (StorageVersionMigrationPatchMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*StorageVersionMigrationPatch)(nil)).Elem()
}

func (o StorageVersionMigrationPatchMapOutput) ToStorageVersionMigrationPatchMapOutput() StorageVersionMigrationPatchMapOutput {
	return o
}

func (o StorageVersionMigrationPatchMapOutput) ToStorageVersionMigrationPatchMapOutputWithContext(ctx context.Context) StorageVersionMigrationPatchMapOutput {
	return o
}

func (o StorageVersionMigrationPatchMapOutput) MapIndex(k pulumi.StringInput) StorageVersionMigrationPatchOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *StorageVersionMigrationPatch {
		return vs[0].(map[string]*StorageVersionMigrationPatch)[vs[1].(string)]
	}).(StorageVersionMigrationPatchOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*StorageVersionMigrationPatchInput)(nil)).Elem(), &StorageVersionMigrationPatch{})
	pulumi.RegisterInputType(reflect.TypeOf((*StorageVersionMigrationPatchArrayInput)(nil)).Elem(), StorageVersionMigrationPatchArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*StorageVersionMigrationPatchMapInput)(nil)).Elem(), StorageVersionMigrationPatchMap{})
	pulumi.RegisterOutputType(StorageVersionMigrationPatchOutput{})
	pulumi.RegisterOutputType(StorageVersionMigrationPatchArrayOutput{})
	pulumi.RegisterOutputType(StorageVersionMigrationPatchMapOutput{})
}
