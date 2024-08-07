// Code generated by pulumigen DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package v1alpha1

import (
	"context"
	"reflect"

	"errors"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/utilities"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// StorageVersionMigrationList is a collection of storage version migrations.
type StorageVersionMigrationList struct {
	pulumi.CustomResourceState

	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion pulumi.StringOutput `pulumi:"apiVersion"`
	// Items is the list of StorageVersionMigration
	Items StorageVersionMigrationTypeArrayOutput `pulumi:"items"`
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind pulumi.StringOutput `pulumi:"kind"`
	// Standard list metadata More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata metav1.ListMetaOutput `pulumi:"metadata"`
}

// NewStorageVersionMigrationList registers a new resource with the given unique name, arguments, and options.
func NewStorageVersionMigrationList(ctx *pulumi.Context,
	name string, args *StorageVersionMigrationListArgs, opts ...pulumi.ResourceOption) (*StorageVersionMigrationList, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Items == nil {
		return nil, errors.New("invalid value for required argument 'Items'")
	}
	args.ApiVersion = pulumi.StringPtr("storagemigration.k8s.io/v1alpha1")
	args.Kind = pulumi.StringPtr("StorageVersionMigrationList")
	opts = utilities.PkgResourceDefaultOpts(opts)
	var resource StorageVersionMigrationList
	err := ctx.RegisterResource("kubernetes:storagemigration.k8s.io/v1alpha1:StorageVersionMigrationList", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetStorageVersionMigrationList gets an existing StorageVersionMigrationList resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetStorageVersionMigrationList(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *StorageVersionMigrationListState, opts ...pulumi.ResourceOption) (*StorageVersionMigrationList, error) {
	var resource StorageVersionMigrationList
	err := ctx.ReadResource("kubernetes:storagemigration.k8s.io/v1alpha1:StorageVersionMigrationList", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering StorageVersionMigrationList resources.
type storageVersionMigrationListState struct {
}

type StorageVersionMigrationListState struct {
}

func (StorageVersionMigrationListState) ElementType() reflect.Type {
	return reflect.TypeOf((*storageVersionMigrationListState)(nil)).Elem()
}

type storageVersionMigrationListArgs struct {
	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion *string `pulumi:"apiVersion"`
	// Items is the list of StorageVersionMigration
	Items []StorageVersionMigrationType `pulumi:"items"`
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind *string `pulumi:"kind"`
	// Standard list metadata More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata *metav1.ListMeta `pulumi:"metadata"`
}

// The set of arguments for constructing a StorageVersionMigrationList resource.
type StorageVersionMigrationListArgs struct {
	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion pulumi.StringPtrInput
	// Items is the list of StorageVersionMigration
	Items StorageVersionMigrationTypeArrayInput
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind pulumi.StringPtrInput
	// Standard list metadata More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata metav1.ListMetaPtrInput
}

func (StorageVersionMigrationListArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*storageVersionMigrationListArgs)(nil)).Elem()
}

type StorageVersionMigrationListInput interface {
	pulumi.Input

	ToStorageVersionMigrationListOutput() StorageVersionMigrationListOutput
	ToStorageVersionMigrationListOutputWithContext(ctx context.Context) StorageVersionMigrationListOutput
}

func (*StorageVersionMigrationList) ElementType() reflect.Type {
	return reflect.TypeOf((**StorageVersionMigrationList)(nil)).Elem()
}

func (i *StorageVersionMigrationList) ToStorageVersionMigrationListOutput() StorageVersionMigrationListOutput {
	return i.ToStorageVersionMigrationListOutputWithContext(context.Background())
}

func (i *StorageVersionMigrationList) ToStorageVersionMigrationListOutputWithContext(ctx context.Context) StorageVersionMigrationListOutput {
	return pulumi.ToOutputWithContext(ctx, i).(StorageVersionMigrationListOutput)
}

// StorageVersionMigrationListArrayInput is an input type that accepts StorageVersionMigrationListArray and StorageVersionMigrationListArrayOutput values.
// You can construct a concrete instance of `StorageVersionMigrationListArrayInput` via:
//
//	StorageVersionMigrationListArray{ StorageVersionMigrationListArgs{...} }
type StorageVersionMigrationListArrayInput interface {
	pulumi.Input

	ToStorageVersionMigrationListArrayOutput() StorageVersionMigrationListArrayOutput
	ToStorageVersionMigrationListArrayOutputWithContext(context.Context) StorageVersionMigrationListArrayOutput
}

type StorageVersionMigrationListArray []StorageVersionMigrationListInput

func (StorageVersionMigrationListArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*StorageVersionMigrationList)(nil)).Elem()
}

func (i StorageVersionMigrationListArray) ToStorageVersionMigrationListArrayOutput() StorageVersionMigrationListArrayOutput {
	return i.ToStorageVersionMigrationListArrayOutputWithContext(context.Background())
}

func (i StorageVersionMigrationListArray) ToStorageVersionMigrationListArrayOutputWithContext(ctx context.Context) StorageVersionMigrationListArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(StorageVersionMigrationListArrayOutput)
}

// StorageVersionMigrationListMapInput is an input type that accepts StorageVersionMigrationListMap and StorageVersionMigrationListMapOutput values.
// You can construct a concrete instance of `StorageVersionMigrationListMapInput` via:
//
//	StorageVersionMigrationListMap{ "key": StorageVersionMigrationListArgs{...} }
type StorageVersionMigrationListMapInput interface {
	pulumi.Input

	ToStorageVersionMigrationListMapOutput() StorageVersionMigrationListMapOutput
	ToStorageVersionMigrationListMapOutputWithContext(context.Context) StorageVersionMigrationListMapOutput
}

type StorageVersionMigrationListMap map[string]StorageVersionMigrationListInput

func (StorageVersionMigrationListMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*StorageVersionMigrationList)(nil)).Elem()
}

func (i StorageVersionMigrationListMap) ToStorageVersionMigrationListMapOutput() StorageVersionMigrationListMapOutput {
	return i.ToStorageVersionMigrationListMapOutputWithContext(context.Background())
}

func (i StorageVersionMigrationListMap) ToStorageVersionMigrationListMapOutputWithContext(ctx context.Context) StorageVersionMigrationListMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(StorageVersionMigrationListMapOutput)
}

type StorageVersionMigrationListOutput struct{ *pulumi.OutputState }

func (StorageVersionMigrationListOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**StorageVersionMigrationList)(nil)).Elem()
}

func (o StorageVersionMigrationListOutput) ToStorageVersionMigrationListOutput() StorageVersionMigrationListOutput {
	return o
}

func (o StorageVersionMigrationListOutput) ToStorageVersionMigrationListOutputWithContext(ctx context.Context) StorageVersionMigrationListOutput {
	return o
}

// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (o StorageVersionMigrationListOutput) ApiVersion() pulumi.StringOutput {
	return o.ApplyT(func(v *StorageVersionMigrationList) pulumi.StringOutput { return v.ApiVersion }).(pulumi.StringOutput)
}

// Items is the list of StorageVersionMigration
func (o StorageVersionMigrationListOutput) Items() StorageVersionMigrationTypeArrayOutput {
	return o.ApplyT(func(v *StorageVersionMigrationList) StorageVersionMigrationTypeArrayOutput { return v.Items }).(StorageVersionMigrationTypeArrayOutput)
}

// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (o StorageVersionMigrationListOutput) Kind() pulumi.StringOutput {
	return o.ApplyT(func(v *StorageVersionMigrationList) pulumi.StringOutput { return v.Kind }).(pulumi.StringOutput)
}

// Standard list metadata More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
func (o StorageVersionMigrationListOutput) Metadata() metav1.ListMetaOutput {
	return o.ApplyT(func(v *StorageVersionMigrationList) metav1.ListMetaOutput { return v.Metadata }).(metav1.ListMetaOutput)
}

type StorageVersionMigrationListArrayOutput struct{ *pulumi.OutputState }

func (StorageVersionMigrationListArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*StorageVersionMigrationList)(nil)).Elem()
}

func (o StorageVersionMigrationListArrayOutput) ToStorageVersionMigrationListArrayOutput() StorageVersionMigrationListArrayOutput {
	return o
}

func (o StorageVersionMigrationListArrayOutput) ToStorageVersionMigrationListArrayOutputWithContext(ctx context.Context) StorageVersionMigrationListArrayOutput {
	return o
}

func (o StorageVersionMigrationListArrayOutput) Index(i pulumi.IntInput) StorageVersionMigrationListOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *StorageVersionMigrationList {
		return vs[0].([]*StorageVersionMigrationList)[vs[1].(int)]
	}).(StorageVersionMigrationListOutput)
}

type StorageVersionMigrationListMapOutput struct{ *pulumi.OutputState }

func (StorageVersionMigrationListMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*StorageVersionMigrationList)(nil)).Elem()
}

func (o StorageVersionMigrationListMapOutput) ToStorageVersionMigrationListMapOutput() StorageVersionMigrationListMapOutput {
	return o
}

func (o StorageVersionMigrationListMapOutput) ToStorageVersionMigrationListMapOutputWithContext(ctx context.Context) StorageVersionMigrationListMapOutput {
	return o
}

func (o StorageVersionMigrationListMapOutput) MapIndex(k pulumi.StringInput) StorageVersionMigrationListOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *StorageVersionMigrationList {
		return vs[0].(map[string]*StorageVersionMigrationList)[vs[1].(string)]
	}).(StorageVersionMigrationListOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*StorageVersionMigrationListInput)(nil)).Elem(), &StorageVersionMigrationList{})
	pulumi.RegisterInputType(reflect.TypeOf((*StorageVersionMigrationListArrayInput)(nil)).Elem(), StorageVersionMigrationListArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*StorageVersionMigrationListMapInput)(nil)).Elem(), StorageVersionMigrationListMap{})
	pulumi.RegisterOutputType(StorageVersionMigrationListOutput{})
	pulumi.RegisterOutputType(StorageVersionMigrationListArrayOutput{})
	pulumi.RegisterOutputType(StorageVersionMigrationListMapOutput{})
}
