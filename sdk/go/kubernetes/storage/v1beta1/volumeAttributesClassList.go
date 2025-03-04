// Code generated by pulumigen DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package v1beta1

import (
	"context"
	"reflect"

	"errors"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/utilities"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// VolumeAttributesClassList is a collection of VolumeAttributesClass objects.
type VolumeAttributesClassList struct {
	pulumi.CustomResourceState

	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion pulumi.StringOutput `pulumi:"apiVersion"`
	// items is the list of VolumeAttributesClass objects.
	Items VolumeAttributesClassTypeArrayOutput `pulumi:"items"`
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind pulumi.StringOutput `pulumi:"kind"`
	// Standard list metadata More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata metav1.ListMetaOutput `pulumi:"metadata"`
}

// NewVolumeAttributesClassList registers a new resource with the given unique name, arguments, and options.
func NewVolumeAttributesClassList(ctx *pulumi.Context,
	name string, args *VolumeAttributesClassListArgs, opts ...pulumi.ResourceOption) (*VolumeAttributesClassList, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Items == nil {
		return nil, errors.New("invalid value for required argument 'Items'")
	}
	args.ApiVersion = pulumi.StringPtr("storage.k8s.io/v1beta1")
	args.Kind = pulumi.StringPtr("VolumeAttributesClassList")
	opts = utilities.PkgResourceDefaultOpts(opts)
	var resource VolumeAttributesClassList
	err := ctx.RegisterResource("kubernetes:storage.k8s.io/v1beta1:VolumeAttributesClassList", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetVolumeAttributesClassList gets an existing VolumeAttributesClassList resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetVolumeAttributesClassList(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *VolumeAttributesClassListState, opts ...pulumi.ResourceOption) (*VolumeAttributesClassList, error) {
	var resource VolumeAttributesClassList
	err := ctx.ReadResource("kubernetes:storage.k8s.io/v1beta1:VolumeAttributesClassList", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering VolumeAttributesClassList resources.
type volumeAttributesClassListState struct {
}

type VolumeAttributesClassListState struct {
}

func (VolumeAttributesClassListState) ElementType() reflect.Type {
	return reflect.TypeOf((*volumeAttributesClassListState)(nil)).Elem()
}

type volumeAttributesClassListArgs struct {
	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion *string `pulumi:"apiVersion"`
	// items is the list of VolumeAttributesClass objects.
	Items []VolumeAttributesClassType `pulumi:"items"`
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind *string `pulumi:"kind"`
	// Standard list metadata More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata *metav1.ListMeta `pulumi:"metadata"`
}

// The set of arguments for constructing a VolumeAttributesClassList resource.
type VolumeAttributesClassListArgs struct {
	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion pulumi.StringPtrInput
	// items is the list of VolumeAttributesClass objects.
	Items VolumeAttributesClassTypeArrayInput
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind pulumi.StringPtrInput
	// Standard list metadata More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata metav1.ListMetaPtrInput
}

func (VolumeAttributesClassListArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*volumeAttributesClassListArgs)(nil)).Elem()
}

type VolumeAttributesClassListInput interface {
	pulumi.Input

	ToVolumeAttributesClassListOutput() VolumeAttributesClassListOutput
	ToVolumeAttributesClassListOutputWithContext(ctx context.Context) VolumeAttributesClassListOutput
}

func (*VolumeAttributesClassList) ElementType() reflect.Type {
	return reflect.TypeOf((**VolumeAttributesClassList)(nil)).Elem()
}

func (i *VolumeAttributesClassList) ToVolumeAttributesClassListOutput() VolumeAttributesClassListOutput {
	return i.ToVolumeAttributesClassListOutputWithContext(context.Background())
}

func (i *VolumeAttributesClassList) ToVolumeAttributesClassListOutputWithContext(ctx context.Context) VolumeAttributesClassListOutput {
	return pulumi.ToOutputWithContext(ctx, i).(VolumeAttributesClassListOutput)
}

// VolumeAttributesClassListArrayInput is an input type that accepts VolumeAttributesClassListArray and VolumeAttributesClassListArrayOutput values.
// You can construct a concrete instance of `VolumeAttributesClassListArrayInput` via:
//
//	VolumeAttributesClassListArray{ VolumeAttributesClassListArgs{...} }
type VolumeAttributesClassListArrayInput interface {
	pulumi.Input

	ToVolumeAttributesClassListArrayOutput() VolumeAttributesClassListArrayOutput
	ToVolumeAttributesClassListArrayOutputWithContext(context.Context) VolumeAttributesClassListArrayOutput
}

type VolumeAttributesClassListArray []VolumeAttributesClassListInput

func (VolumeAttributesClassListArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*VolumeAttributesClassList)(nil)).Elem()
}

func (i VolumeAttributesClassListArray) ToVolumeAttributesClassListArrayOutput() VolumeAttributesClassListArrayOutput {
	return i.ToVolumeAttributesClassListArrayOutputWithContext(context.Background())
}

func (i VolumeAttributesClassListArray) ToVolumeAttributesClassListArrayOutputWithContext(ctx context.Context) VolumeAttributesClassListArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(VolumeAttributesClassListArrayOutput)
}

// VolumeAttributesClassListMapInput is an input type that accepts VolumeAttributesClassListMap and VolumeAttributesClassListMapOutput values.
// You can construct a concrete instance of `VolumeAttributesClassListMapInput` via:
//
//	VolumeAttributesClassListMap{ "key": VolumeAttributesClassListArgs{...} }
type VolumeAttributesClassListMapInput interface {
	pulumi.Input

	ToVolumeAttributesClassListMapOutput() VolumeAttributesClassListMapOutput
	ToVolumeAttributesClassListMapOutputWithContext(context.Context) VolumeAttributesClassListMapOutput
}

type VolumeAttributesClassListMap map[string]VolumeAttributesClassListInput

func (VolumeAttributesClassListMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*VolumeAttributesClassList)(nil)).Elem()
}

func (i VolumeAttributesClassListMap) ToVolumeAttributesClassListMapOutput() VolumeAttributesClassListMapOutput {
	return i.ToVolumeAttributesClassListMapOutputWithContext(context.Background())
}

func (i VolumeAttributesClassListMap) ToVolumeAttributesClassListMapOutputWithContext(ctx context.Context) VolumeAttributesClassListMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(VolumeAttributesClassListMapOutput)
}

type VolumeAttributesClassListOutput struct{ *pulumi.OutputState }

func (VolumeAttributesClassListOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**VolumeAttributesClassList)(nil)).Elem()
}

func (o VolumeAttributesClassListOutput) ToVolumeAttributesClassListOutput() VolumeAttributesClassListOutput {
	return o
}

func (o VolumeAttributesClassListOutput) ToVolumeAttributesClassListOutputWithContext(ctx context.Context) VolumeAttributesClassListOutput {
	return o
}

// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (o VolumeAttributesClassListOutput) ApiVersion() pulumi.StringOutput {
	return o.ApplyT(func(v *VolumeAttributesClassList) pulumi.StringOutput { return v.ApiVersion }).(pulumi.StringOutput)
}

// items is the list of VolumeAttributesClass objects.
func (o VolumeAttributesClassListOutput) Items() VolumeAttributesClassTypeArrayOutput {
	return o.ApplyT(func(v *VolumeAttributesClassList) VolumeAttributesClassTypeArrayOutput { return v.Items }).(VolumeAttributesClassTypeArrayOutput)
}

// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (o VolumeAttributesClassListOutput) Kind() pulumi.StringOutput {
	return o.ApplyT(func(v *VolumeAttributesClassList) pulumi.StringOutput { return v.Kind }).(pulumi.StringOutput)
}

// Standard list metadata More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
func (o VolumeAttributesClassListOutput) Metadata() metav1.ListMetaOutput {
	return o.ApplyT(func(v *VolumeAttributesClassList) metav1.ListMetaOutput { return v.Metadata }).(metav1.ListMetaOutput)
}

type VolumeAttributesClassListArrayOutput struct{ *pulumi.OutputState }

func (VolumeAttributesClassListArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*VolumeAttributesClassList)(nil)).Elem()
}

func (o VolumeAttributesClassListArrayOutput) ToVolumeAttributesClassListArrayOutput() VolumeAttributesClassListArrayOutput {
	return o
}

func (o VolumeAttributesClassListArrayOutput) ToVolumeAttributesClassListArrayOutputWithContext(ctx context.Context) VolumeAttributesClassListArrayOutput {
	return o
}

func (o VolumeAttributesClassListArrayOutput) Index(i pulumi.IntInput) VolumeAttributesClassListOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *VolumeAttributesClassList {
		return vs[0].([]*VolumeAttributesClassList)[vs[1].(int)]
	}).(VolumeAttributesClassListOutput)
}

type VolumeAttributesClassListMapOutput struct{ *pulumi.OutputState }

func (VolumeAttributesClassListMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*VolumeAttributesClassList)(nil)).Elem()
}

func (o VolumeAttributesClassListMapOutput) ToVolumeAttributesClassListMapOutput() VolumeAttributesClassListMapOutput {
	return o
}

func (o VolumeAttributesClassListMapOutput) ToVolumeAttributesClassListMapOutputWithContext(ctx context.Context) VolumeAttributesClassListMapOutput {
	return o
}

func (o VolumeAttributesClassListMapOutput) MapIndex(k pulumi.StringInput) VolumeAttributesClassListOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *VolumeAttributesClassList {
		return vs[0].(map[string]*VolumeAttributesClassList)[vs[1].(string)]
	}).(VolumeAttributesClassListOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*VolumeAttributesClassListInput)(nil)).Elem(), &VolumeAttributesClassList{})
	pulumi.RegisterInputType(reflect.TypeOf((*VolumeAttributesClassListArrayInput)(nil)).Elem(), VolumeAttributesClassListArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*VolumeAttributesClassListMapInput)(nil)).Elem(), VolumeAttributesClassListMap{})
	pulumi.RegisterOutputType(VolumeAttributesClassListOutput{})
	pulumi.RegisterOutputType(VolumeAttributesClassListArrayOutput{})
	pulumi.RegisterOutputType(VolumeAttributesClassListMapOutput{})
}
