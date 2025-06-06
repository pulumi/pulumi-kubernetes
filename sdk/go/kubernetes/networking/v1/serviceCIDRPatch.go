// Code generated by pulumigen DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package v1

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
// ServiceCIDR defines a range of IP addresses using CIDR format (e.g. 192.168.0.0/24 or 2001:db2::/64). This range is used to allocate ClusterIPs to Service objects.
type ServiceCIDRPatch struct {
	pulumi.CustomResourceState

	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion pulumi.StringPtrOutput `pulumi:"apiVersion"`
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind pulumi.StringPtrOutput `pulumi:"kind"`
	// Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata metav1.ObjectMetaPatchPtrOutput `pulumi:"metadata"`
	// spec is the desired state of the ServiceCIDR. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec ServiceCIDRSpecPatchPtrOutput `pulumi:"spec"`
	// status represents the current state of the ServiceCIDR. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Status ServiceCIDRStatusPatchPtrOutput `pulumi:"status"`
}

// NewServiceCIDRPatch registers a new resource with the given unique name, arguments, and options.
func NewServiceCIDRPatch(ctx *pulumi.Context,
	name string, args *ServiceCIDRPatchArgs, opts ...pulumi.ResourceOption) (*ServiceCIDRPatch, error) {
	if args == nil {
		args = &ServiceCIDRPatchArgs{}
	}

	args.ApiVersion = pulumi.StringPtr("networking.k8s.io/v1")
	args.Kind = pulumi.StringPtr("ServiceCIDR")
	aliases := pulumi.Aliases([]pulumi.Alias{
		{
			Type: pulumi.String("kubernetes:networking.k8s.io/v1alpha1:ServiceCIDRPatch"),
		},
		{
			Type: pulumi.String("kubernetes:networking.k8s.io/v1beta1:ServiceCIDRPatch"),
		},
	})
	opts = append(opts, aliases)
	opts = utilities.PkgResourceDefaultOpts(opts)
	var resource ServiceCIDRPatch
	err := ctx.RegisterResource("kubernetes:networking.k8s.io/v1:ServiceCIDRPatch", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetServiceCIDRPatch gets an existing ServiceCIDRPatch resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetServiceCIDRPatch(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *ServiceCIDRPatchState, opts ...pulumi.ResourceOption) (*ServiceCIDRPatch, error) {
	var resource ServiceCIDRPatch
	err := ctx.ReadResource("kubernetes:networking.k8s.io/v1:ServiceCIDRPatch", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering ServiceCIDRPatch resources.
type serviceCIDRPatchState struct {
}

type ServiceCIDRPatchState struct {
}

func (ServiceCIDRPatchState) ElementType() reflect.Type {
	return reflect.TypeOf((*serviceCIDRPatchState)(nil)).Elem()
}

type serviceCIDRPatchArgs struct {
	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion *string `pulumi:"apiVersion"`
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind *string `pulumi:"kind"`
	// Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata *metav1.ObjectMetaPatch `pulumi:"metadata"`
	// spec is the desired state of the ServiceCIDR. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec *ServiceCIDRSpecPatch `pulumi:"spec"`
}

// The set of arguments for constructing a ServiceCIDRPatch resource.
type ServiceCIDRPatchArgs struct {
	// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	ApiVersion pulumi.StringPtrInput
	// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind pulumi.StringPtrInput
	// Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	Metadata metav1.ObjectMetaPatchPtrInput
	// spec is the desired state of the ServiceCIDR. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	Spec ServiceCIDRSpecPatchPtrInput
}

func (ServiceCIDRPatchArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*serviceCIDRPatchArgs)(nil)).Elem()
}

type ServiceCIDRPatchInput interface {
	pulumi.Input

	ToServiceCIDRPatchOutput() ServiceCIDRPatchOutput
	ToServiceCIDRPatchOutputWithContext(ctx context.Context) ServiceCIDRPatchOutput
}

func (*ServiceCIDRPatch) ElementType() reflect.Type {
	return reflect.TypeOf((**ServiceCIDRPatch)(nil)).Elem()
}

func (i *ServiceCIDRPatch) ToServiceCIDRPatchOutput() ServiceCIDRPatchOutput {
	return i.ToServiceCIDRPatchOutputWithContext(context.Background())
}

func (i *ServiceCIDRPatch) ToServiceCIDRPatchOutputWithContext(ctx context.Context) ServiceCIDRPatchOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ServiceCIDRPatchOutput)
}

// ServiceCIDRPatchArrayInput is an input type that accepts ServiceCIDRPatchArray and ServiceCIDRPatchArrayOutput values.
// You can construct a concrete instance of `ServiceCIDRPatchArrayInput` via:
//
//	ServiceCIDRPatchArray{ ServiceCIDRPatchArgs{...} }
type ServiceCIDRPatchArrayInput interface {
	pulumi.Input

	ToServiceCIDRPatchArrayOutput() ServiceCIDRPatchArrayOutput
	ToServiceCIDRPatchArrayOutputWithContext(context.Context) ServiceCIDRPatchArrayOutput
}

type ServiceCIDRPatchArray []ServiceCIDRPatchInput

func (ServiceCIDRPatchArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*ServiceCIDRPatch)(nil)).Elem()
}

func (i ServiceCIDRPatchArray) ToServiceCIDRPatchArrayOutput() ServiceCIDRPatchArrayOutput {
	return i.ToServiceCIDRPatchArrayOutputWithContext(context.Background())
}

func (i ServiceCIDRPatchArray) ToServiceCIDRPatchArrayOutputWithContext(ctx context.Context) ServiceCIDRPatchArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ServiceCIDRPatchArrayOutput)
}

// ServiceCIDRPatchMapInput is an input type that accepts ServiceCIDRPatchMap and ServiceCIDRPatchMapOutput values.
// You can construct a concrete instance of `ServiceCIDRPatchMapInput` via:
//
//	ServiceCIDRPatchMap{ "key": ServiceCIDRPatchArgs{...} }
type ServiceCIDRPatchMapInput interface {
	pulumi.Input

	ToServiceCIDRPatchMapOutput() ServiceCIDRPatchMapOutput
	ToServiceCIDRPatchMapOutputWithContext(context.Context) ServiceCIDRPatchMapOutput
}

type ServiceCIDRPatchMap map[string]ServiceCIDRPatchInput

func (ServiceCIDRPatchMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*ServiceCIDRPatch)(nil)).Elem()
}

func (i ServiceCIDRPatchMap) ToServiceCIDRPatchMapOutput() ServiceCIDRPatchMapOutput {
	return i.ToServiceCIDRPatchMapOutputWithContext(context.Background())
}

func (i ServiceCIDRPatchMap) ToServiceCIDRPatchMapOutputWithContext(ctx context.Context) ServiceCIDRPatchMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ServiceCIDRPatchMapOutput)
}

type ServiceCIDRPatchOutput struct{ *pulumi.OutputState }

func (ServiceCIDRPatchOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**ServiceCIDRPatch)(nil)).Elem()
}

func (o ServiceCIDRPatchOutput) ToServiceCIDRPatchOutput() ServiceCIDRPatchOutput {
	return o
}

func (o ServiceCIDRPatchOutput) ToServiceCIDRPatchOutputWithContext(ctx context.Context) ServiceCIDRPatchOutput {
	return o
}

// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
func (o ServiceCIDRPatchOutput) ApiVersion() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *ServiceCIDRPatch) pulumi.StringPtrOutput { return v.ApiVersion }).(pulumi.StringPtrOutput)
}

// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
func (o ServiceCIDRPatchOutput) Kind() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *ServiceCIDRPatch) pulumi.StringPtrOutput { return v.Kind }).(pulumi.StringPtrOutput)
}

// Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
func (o ServiceCIDRPatchOutput) Metadata() metav1.ObjectMetaPatchPtrOutput {
	return o.ApplyT(func(v *ServiceCIDRPatch) metav1.ObjectMetaPatchPtrOutput { return v.Metadata }).(metav1.ObjectMetaPatchPtrOutput)
}

// spec is the desired state of the ServiceCIDR. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
func (o ServiceCIDRPatchOutput) Spec() ServiceCIDRSpecPatchPtrOutput {
	return o.ApplyT(func(v *ServiceCIDRPatch) ServiceCIDRSpecPatchPtrOutput { return v.Spec }).(ServiceCIDRSpecPatchPtrOutput)
}

// status represents the current state of the ServiceCIDR. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
func (o ServiceCIDRPatchOutput) Status() ServiceCIDRStatusPatchPtrOutput {
	return o.ApplyT(func(v *ServiceCIDRPatch) ServiceCIDRStatusPatchPtrOutput { return v.Status }).(ServiceCIDRStatusPatchPtrOutput)
}

type ServiceCIDRPatchArrayOutput struct{ *pulumi.OutputState }

func (ServiceCIDRPatchArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*ServiceCIDRPatch)(nil)).Elem()
}

func (o ServiceCIDRPatchArrayOutput) ToServiceCIDRPatchArrayOutput() ServiceCIDRPatchArrayOutput {
	return o
}

func (o ServiceCIDRPatchArrayOutput) ToServiceCIDRPatchArrayOutputWithContext(ctx context.Context) ServiceCIDRPatchArrayOutput {
	return o
}

func (o ServiceCIDRPatchArrayOutput) Index(i pulumi.IntInput) ServiceCIDRPatchOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *ServiceCIDRPatch {
		return vs[0].([]*ServiceCIDRPatch)[vs[1].(int)]
	}).(ServiceCIDRPatchOutput)
}

type ServiceCIDRPatchMapOutput struct{ *pulumi.OutputState }

func (ServiceCIDRPatchMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*ServiceCIDRPatch)(nil)).Elem()
}

func (o ServiceCIDRPatchMapOutput) ToServiceCIDRPatchMapOutput() ServiceCIDRPatchMapOutput {
	return o
}

func (o ServiceCIDRPatchMapOutput) ToServiceCIDRPatchMapOutputWithContext(ctx context.Context) ServiceCIDRPatchMapOutput {
	return o
}

func (o ServiceCIDRPatchMapOutput) MapIndex(k pulumi.StringInput) ServiceCIDRPatchOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *ServiceCIDRPatch {
		return vs[0].(map[string]*ServiceCIDRPatch)[vs[1].(string)]
	}).(ServiceCIDRPatchOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*ServiceCIDRPatchInput)(nil)).Elem(), &ServiceCIDRPatch{})
	pulumi.RegisterInputType(reflect.TypeOf((*ServiceCIDRPatchArrayInput)(nil)).Elem(), ServiceCIDRPatchArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*ServiceCIDRPatchMapInput)(nil)).Elem(), ServiceCIDRPatchMap{})
	pulumi.RegisterOutputType(ServiceCIDRPatchOutput{})
	pulumi.RegisterOutputType(ServiceCIDRPatchArrayOutput{})
	pulumi.RegisterOutputType(ServiceCIDRPatchMapOutput{})
}
