// Code generated by pulumigen DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package v1alpha2

import (
	"context"
	"reflect"

	"errors"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/utilities"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ResourceClaimParametersList is a collection of ResourceClaimParameters.
type ResourceClaimParametersList struct {
	pulumi.CustomResourceState

	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion pulumi.StringOutput `pulumi:"apiVersion"`
	// Items is the list of node resource capacity objects.
	Items ResourceClaimParametersTypeArrayOutput `pulumi:"items"`
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind pulumi.StringOutput `pulumi:"kind"`
	// Standard list metadata
	Metadata metav1.ListMetaOutput `pulumi:"metadata"`
}

// NewResourceClaimParametersList registers a new resource with the given unique name, arguments, and options.
func NewResourceClaimParametersList(ctx *pulumi.Context,
	name string, args *ResourceClaimParametersListArgs, opts ...pulumi.ResourceOption) (*ResourceClaimParametersList, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Items == nil {
		return nil, errors.New("invalid value for required argument 'Items'")
	}
	args.ApiVersion = pulumi.StringPtr("resource.k8s.io/v1alpha2")
	args.Kind = pulumi.StringPtr("ResourceClaimParametersList")
	opts = utilities.PkgResourceDefaultOpts(opts)
	var resource ResourceClaimParametersList
	err := ctx.RegisterResource("kubernetes:resource.k8s.io/v1alpha2:ResourceClaimParametersList", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetResourceClaimParametersList gets an existing ResourceClaimParametersList resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetResourceClaimParametersList(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *ResourceClaimParametersListState, opts ...pulumi.ResourceOption) (*ResourceClaimParametersList, error) {
	var resource ResourceClaimParametersList
	err := ctx.ReadResource("kubernetes:resource.k8s.io/v1alpha2:ResourceClaimParametersList", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering ResourceClaimParametersList resources.
type resourceClaimParametersListState struct {
}

type ResourceClaimParametersListState struct {
}

func (ResourceClaimParametersListState) ElementType() reflect.Type {
	return reflect.TypeOf((*resourceClaimParametersListState)(nil)).Elem()
}

type resourceClaimParametersListArgs struct {
	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion *string `pulumi:"apiVersion"`
	// Items is the list of node resource capacity objects.
	Items []ResourceClaimParametersType `pulumi:"items"`
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind *string `pulumi:"kind"`
	// Standard list metadata
	Metadata *metav1.ListMeta `pulumi:"metadata"`
}

// The set of arguments for constructing a ResourceClaimParametersList resource.
type ResourceClaimParametersListArgs struct {
	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion pulumi.StringPtrInput
	// Items is the list of node resource capacity objects.
	Items ResourceClaimParametersTypeArrayInput
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind pulumi.StringPtrInput
	// Standard list metadata
	Metadata metav1.ListMetaPtrInput
}

func (ResourceClaimParametersListArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*resourceClaimParametersListArgs)(nil)).Elem()
}

type ResourceClaimParametersListInput interface {
	pulumi.Input

	ToResourceClaimParametersListOutput() ResourceClaimParametersListOutput
	ToResourceClaimParametersListOutputWithContext(ctx context.Context) ResourceClaimParametersListOutput
}

func (*ResourceClaimParametersList) ElementType() reflect.Type {
	return reflect.TypeOf((**ResourceClaimParametersList)(nil)).Elem()
}

func (i *ResourceClaimParametersList) ToResourceClaimParametersListOutput() ResourceClaimParametersListOutput {
	return i.ToResourceClaimParametersListOutputWithContext(context.Background())
}

func (i *ResourceClaimParametersList) ToResourceClaimParametersListOutputWithContext(ctx context.Context) ResourceClaimParametersListOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ResourceClaimParametersListOutput)
}

// ResourceClaimParametersListArrayInput is an input type that accepts ResourceClaimParametersListArray and ResourceClaimParametersListArrayOutput values.
// You can construct a concrete instance of `ResourceClaimParametersListArrayInput` via:
//
//	ResourceClaimParametersListArray{ ResourceClaimParametersListArgs{...} }
type ResourceClaimParametersListArrayInput interface {
	pulumi.Input

	ToResourceClaimParametersListArrayOutput() ResourceClaimParametersListArrayOutput
	ToResourceClaimParametersListArrayOutputWithContext(context.Context) ResourceClaimParametersListArrayOutput
}

type ResourceClaimParametersListArray []ResourceClaimParametersListInput

func (ResourceClaimParametersListArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*ResourceClaimParametersList)(nil)).Elem()
}

func (i ResourceClaimParametersListArray) ToResourceClaimParametersListArrayOutput() ResourceClaimParametersListArrayOutput {
	return i.ToResourceClaimParametersListArrayOutputWithContext(context.Background())
}

func (i ResourceClaimParametersListArray) ToResourceClaimParametersListArrayOutputWithContext(ctx context.Context) ResourceClaimParametersListArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ResourceClaimParametersListArrayOutput)
}

// ResourceClaimParametersListMapInput is an input type that accepts ResourceClaimParametersListMap and ResourceClaimParametersListMapOutput values.
// You can construct a concrete instance of `ResourceClaimParametersListMapInput` via:
//
//	ResourceClaimParametersListMap{ "key": ResourceClaimParametersListArgs{...} }
type ResourceClaimParametersListMapInput interface {
	pulumi.Input

	ToResourceClaimParametersListMapOutput() ResourceClaimParametersListMapOutput
	ToResourceClaimParametersListMapOutputWithContext(context.Context) ResourceClaimParametersListMapOutput
}

type ResourceClaimParametersListMap map[string]ResourceClaimParametersListInput

func (ResourceClaimParametersListMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*ResourceClaimParametersList)(nil)).Elem()
}

func (i ResourceClaimParametersListMap) ToResourceClaimParametersListMapOutput() ResourceClaimParametersListMapOutput {
	return i.ToResourceClaimParametersListMapOutputWithContext(context.Background())
}

func (i ResourceClaimParametersListMap) ToResourceClaimParametersListMapOutputWithContext(ctx context.Context) ResourceClaimParametersListMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ResourceClaimParametersListMapOutput)
}

type ResourceClaimParametersListOutput struct{ *pulumi.OutputState }

func (ResourceClaimParametersListOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**ResourceClaimParametersList)(nil)).Elem()
}

func (o ResourceClaimParametersListOutput) ToResourceClaimParametersListOutput() ResourceClaimParametersListOutput {
	return o
}

func (o ResourceClaimParametersListOutput) ToResourceClaimParametersListOutputWithContext(ctx context.Context) ResourceClaimParametersListOutput {
	return o
}

// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (o ResourceClaimParametersListOutput) ApiVersion() pulumi.StringOutput {
	return o.ApplyT(func(v *ResourceClaimParametersList) pulumi.StringOutput { return v.ApiVersion }).(pulumi.StringOutput)
}

// Items is the list of node resource capacity objects.
func (o ResourceClaimParametersListOutput) Items() ResourceClaimParametersTypeArrayOutput {
	return o.ApplyT(func(v *ResourceClaimParametersList) ResourceClaimParametersTypeArrayOutput { return v.Items }).(ResourceClaimParametersTypeArrayOutput)
}

// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (o ResourceClaimParametersListOutput) Kind() pulumi.StringOutput {
	return o.ApplyT(func(v *ResourceClaimParametersList) pulumi.StringOutput { return v.Kind }).(pulumi.StringOutput)
}

// Standard list metadata
func (o ResourceClaimParametersListOutput) Metadata() metav1.ListMetaOutput {
	return o.ApplyT(func(v *ResourceClaimParametersList) metav1.ListMetaOutput { return v.Metadata }).(metav1.ListMetaOutput)
}

type ResourceClaimParametersListArrayOutput struct{ *pulumi.OutputState }

func (ResourceClaimParametersListArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*ResourceClaimParametersList)(nil)).Elem()
}

func (o ResourceClaimParametersListArrayOutput) ToResourceClaimParametersListArrayOutput() ResourceClaimParametersListArrayOutput {
	return o
}

func (o ResourceClaimParametersListArrayOutput) ToResourceClaimParametersListArrayOutputWithContext(ctx context.Context) ResourceClaimParametersListArrayOutput {
	return o
}

func (o ResourceClaimParametersListArrayOutput) Index(i pulumi.IntInput) ResourceClaimParametersListOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *ResourceClaimParametersList {
		return vs[0].([]*ResourceClaimParametersList)[vs[1].(int)]
	}).(ResourceClaimParametersListOutput)
}

type ResourceClaimParametersListMapOutput struct{ *pulumi.OutputState }

func (ResourceClaimParametersListMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*ResourceClaimParametersList)(nil)).Elem()
}

func (o ResourceClaimParametersListMapOutput) ToResourceClaimParametersListMapOutput() ResourceClaimParametersListMapOutput {
	return o
}

func (o ResourceClaimParametersListMapOutput) ToResourceClaimParametersListMapOutputWithContext(ctx context.Context) ResourceClaimParametersListMapOutput {
	return o
}

func (o ResourceClaimParametersListMapOutput) MapIndex(k pulumi.StringInput) ResourceClaimParametersListOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *ResourceClaimParametersList {
		return vs[0].(map[string]*ResourceClaimParametersList)[vs[1].(string)]
	}).(ResourceClaimParametersListOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*ResourceClaimParametersListInput)(nil)).Elem(), &ResourceClaimParametersList{})
	pulumi.RegisterInputType(reflect.TypeOf((*ResourceClaimParametersListArrayInput)(nil)).Elem(), ResourceClaimParametersListArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*ResourceClaimParametersListMapInput)(nil)).Elem(), ResourceClaimParametersListMap{})
	pulumi.RegisterOutputType(ResourceClaimParametersListOutput{})
	pulumi.RegisterOutputType(ResourceClaimParametersListArrayOutput{})
	pulumi.RegisterOutputType(ResourceClaimParametersListMapOutput{})
}
