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

// MutatingAdmissionPolicy describes the definition of an admission mutation policy that mutates the object coming into admission chain.
type MutatingAdmissionPolicy struct {
	pulumi.CustomResourceState

	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion pulumi.StringOutput `pulumi:"apiVersion"`
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind pulumi.StringOutput `pulumi:"kind"`
	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
	Metadata metav1.ObjectMetaOutput `pulumi:"metadata"`
	// Specification of the desired behavior of the MutatingAdmissionPolicy.
	Spec MutatingAdmissionPolicySpecOutput `pulumi:"spec"`
}

// NewMutatingAdmissionPolicy registers a new resource with the given unique name, arguments, and options.
func NewMutatingAdmissionPolicy(ctx *pulumi.Context,
	name string, args *MutatingAdmissionPolicyArgs, opts ...pulumi.ResourceOption) (*MutatingAdmissionPolicy, error) {
	if args == nil {
		args = &MutatingAdmissionPolicyArgs{}
	}

	args.ApiVersion = pulumi.StringPtr("admissionregistration.k8s.io/v1alpha1")
	args.Kind = pulumi.StringPtr("MutatingAdmissionPolicy")
	opts = utilities.PkgResourceDefaultOpts(opts)
	var resource MutatingAdmissionPolicy
	err := ctx.RegisterResource("kubernetes:admissionregistration.k8s.io/v1alpha1:MutatingAdmissionPolicy", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetMutatingAdmissionPolicy gets an existing MutatingAdmissionPolicy resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetMutatingAdmissionPolicy(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *MutatingAdmissionPolicyState, opts ...pulumi.ResourceOption) (*MutatingAdmissionPolicy, error) {
	var resource MutatingAdmissionPolicy
	err := ctx.ReadResource("kubernetes:admissionregistration.k8s.io/v1alpha1:MutatingAdmissionPolicy", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering MutatingAdmissionPolicy resources.
type mutatingAdmissionPolicyState struct {
}

type MutatingAdmissionPolicyState struct {
}

func (MutatingAdmissionPolicyState) ElementType() reflect.Type {
	return reflect.TypeOf((*mutatingAdmissionPolicyState)(nil)).Elem()
}

type mutatingAdmissionPolicyArgs struct {
	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion *string `pulumi:"apiVersion"`
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind *string `pulumi:"kind"`
	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
	Metadata *metav1.ObjectMeta `pulumi:"metadata"`
	// Specification of the desired behavior of the MutatingAdmissionPolicy.
	Spec *MutatingAdmissionPolicySpec `pulumi:"spec"`
}

// The set of arguments for constructing a MutatingAdmissionPolicy resource.
type MutatingAdmissionPolicyArgs struct {
	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion pulumi.StringPtrInput
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind pulumi.StringPtrInput
	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
	Metadata metav1.ObjectMetaPtrInput
	// Specification of the desired behavior of the MutatingAdmissionPolicy.
	Spec MutatingAdmissionPolicySpecPtrInput
}

func (MutatingAdmissionPolicyArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*mutatingAdmissionPolicyArgs)(nil)).Elem()
}

type MutatingAdmissionPolicyInput interface {
	pulumi.Input

	ToMutatingAdmissionPolicyOutput() MutatingAdmissionPolicyOutput
	ToMutatingAdmissionPolicyOutputWithContext(ctx context.Context) MutatingAdmissionPolicyOutput
}

func (*MutatingAdmissionPolicy) ElementType() reflect.Type {
	return reflect.TypeOf((**MutatingAdmissionPolicy)(nil)).Elem()
}

func (i *MutatingAdmissionPolicy) ToMutatingAdmissionPolicyOutput() MutatingAdmissionPolicyOutput {
	return i.ToMutatingAdmissionPolicyOutputWithContext(context.Background())
}

func (i *MutatingAdmissionPolicy) ToMutatingAdmissionPolicyOutputWithContext(ctx context.Context) MutatingAdmissionPolicyOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MutatingAdmissionPolicyOutput)
}

// MutatingAdmissionPolicyArrayInput is an input type that accepts MutatingAdmissionPolicyArray and MutatingAdmissionPolicyArrayOutput values.
// You can construct a concrete instance of `MutatingAdmissionPolicyArrayInput` via:
//
//	MutatingAdmissionPolicyArray{ MutatingAdmissionPolicyArgs{...} }
type MutatingAdmissionPolicyArrayInput interface {
	pulumi.Input

	ToMutatingAdmissionPolicyArrayOutput() MutatingAdmissionPolicyArrayOutput
	ToMutatingAdmissionPolicyArrayOutputWithContext(context.Context) MutatingAdmissionPolicyArrayOutput
}

type MutatingAdmissionPolicyArray []MutatingAdmissionPolicyInput

func (MutatingAdmissionPolicyArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*MutatingAdmissionPolicy)(nil)).Elem()
}

func (i MutatingAdmissionPolicyArray) ToMutatingAdmissionPolicyArrayOutput() MutatingAdmissionPolicyArrayOutput {
	return i.ToMutatingAdmissionPolicyArrayOutputWithContext(context.Background())
}

func (i MutatingAdmissionPolicyArray) ToMutatingAdmissionPolicyArrayOutputWithContext(ctx context.Context) MutatingAdmissionPolicyArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MutatingAdmissionPolicyArrayOutput)
}

// MutatingAdmissionPolicyMapInput is an input type that accepts MutatingAdmissionPolicyMap and MutatingAdmissionPolicyMapOutput values.
// You can construct a concrete instance of `MutatingAdmissionPolicyMapInput` via:
//
//	MutatingAdmissionPolicyMap{ "key": MutatingAdmissionPolicyArgs{...} }
type MutatingAdmissionPolicyMapInput interface {
	pulumi.Input

	ToMutatingAdmissionPolicyMapOutput() MutatingAdmissionPolicyMapOutput
	ToMutatingAdmissionPolicyMapOutputWithContext(context.Context) MutatingAdmissionPolicyMapOutput
}

type MutatingAdmissionPolicyMap map[string]MutatingAdmissionPolicyInput

func (MutatingAdmissionPolicyMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*MutatingAdmissionPolicy)(nil)).Elem()
}

func (i MutatingAdmissionPolicyMap) ToMutatingAdmissionPolicyMapOutput() MutatingAdmissionPolicyMapOutput {
	return i.ToMutatingAdmissionPolicyMapOutputWithContext(context.Background())
}

func (i MutatingAdmissionPolicyMap) ToMutatingAdmissionPolicyMapOutputWithContext(ctx context.Context) MutatingAdmissionPolicyMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MutatingAdmissionPolicyMapOutput)
}

type MutatingAdmissionPolicyOutput struct{ *pulumi.OutputState }

func (MutatingAdmissionPolicyOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**MutatingAdmissionPolicy)(nil)).Elem()
}

func (o MutatingAdmissionPolicyOutput) ToMutatingAdmissionPolicyOutput() MutatingAdmissionPolicyOutput {
	return o
}

func (o MutatingAdmissionPolicyOutput) ToMutatingAdmissionPolicyOutputWithContext(ctx context.Context) MutatingAdmissionPolicyOutput {
	return o
}

// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (o MutatingAdmissionPolicyOutput) ApiVersion() pulumi.StringOutput {
	return o.ApplyT(func(v *MutatingAdmissionPolicy) pulumi.StringOutput { return v.ApiVersion }).(pulumi.StringOutput)
}

// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (o MutatingAdmissionPolicyOutput) Kind() pulumi.StringOutput {
	return o.ApplyT(func(v *MutatingAdmissionPolicy) pulumi.StringOutput { return v.Kind }).(pulumi.StringOutput)
}

// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
func (o MutatingAdmissionPolicyOutput) Metadata() metav1.ObjectMetaOutput {
	return o.ApplyT(func(v *MutatingAdmissionPolicy) metav1.ObjectMetaOutput { return v.Metadata }).(metav1.ObjectMetaOutput)
}

// Specification of the desired behavior of the MutatingAdmissionPolicy.
func (o MutatingAdmissionPolicyOutput) Spec() MutatingAdmissionPolicySpecOutput {
	return o.ApplyT(func(v *MutatingAdmissionPolicy) MutatingAdmissionPolicySpecOutput { return v.Spec }).(MutatingAdmissionPolicySpecOutput)
}

type MutatingAdmissionPolicyArrayOutput struct{ *pulumi.OutputState }

func (MutatingAdmissionPolicyArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*MutatingAdmissionPolicy)(nil)).Elem()
}

func (o MutatingAdmissionPolicyArrayOutput) ToMutatingAdmissionPolicyArrayOutput() MutatingAdmissionPolicyArrayOutput {
	return o
}

func (o MutatingAdmissionPolicyArrayOutput) ToMutatingAdmissionPolicyArrayOutputWithContext(ctx context.Context) MutatingAdmissionPolicyArrayOutput {
	return o
}

func (o MutatingAdmissionPolicyArrayOutput) Index(i pulumi.IntInput) MutatingAdmissionPolicyOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *MutatingAdmissionPolicy {
		return vs[0].([]*MutatingAdmissionPolicy)[vs[1].(int)]
	}).(MutatingAdmissionPolicyOutput)
}

type MutatingAdmissionPolicyMapOutput struct{ *pulumi.OutputState }

func (MutatingAdmissionPolicyMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*MutatingAdmissionPolicy)(nil)).Elem()
}

func (o MutatingAdmissionPolicyMapOutput) ToMutatingAdmissionPolicyMapOutput() MutatingAdmissionPolicyMapOutput {
	return o
}

func (o MutatingAdmissionPolicyMapOutput) ToMutatingAdmissionPolicyMapOutputWithContext(ctx context.Context) MutatingAdmissionPolicyMapOutput {
	return o
}

func (o MutatingAdmissionPolicyMapOutput) MapIndex(k pulumi.StringInput) MutatingAdmissionPolicyOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *MutatingAdmissionPolicy {
		return vs[0].(map[string]*MutatingAdmissionPolicy)[vs[1].(string)]
	}).(MutatingAdmissionPolicyOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*MutatingAdmissionPolicyInput)(nil)).Elem(), &MutatingAdmissionPolicy{})
	pulumi.RegisterInputType(reflect.TypeOf((*MutatingAdmissionPolicyArrayInput)(nil)).Elem(), MutatingAdmissionPolicyArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*MutatingAdmissionPolicyMapInput)(nil)).Elem(), MutatingAdmissionPolicyMap{})
	pulumi.RegisterOutputType(MutatingAdmissionPolicyOutput{})
	pulumi.RegisterOutputType(MutatingAdmissionPolicyArrayOutput{})
	pulumi.RegisterOutputType(MutatingAdmissionPolicyMapOutput{})
}