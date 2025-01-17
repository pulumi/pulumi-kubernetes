// Code generated by pulumigen DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package v4

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/utilities"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var _ = utilities.GetEnvOrDefault

// Specification defining the post-renderer to use.
type PostRenderer struct {
	// Arguments to pass to the post-renderer command.
	Args []string `pulumi:"args"`
	// Path to an executable to be used for post rendering.
	Command string `pulumi:"command"`
}

// PostRendererInput is an input type that accepts PostRendererArgs and PostRendererOutput values.
// You can construct a concrete instance of `PostRendererInput` via:
//
//	PostRendererArgs{...}
type PostRendererInput interface {
	pulumi.Input

	ToPostRendererOutput() PostRendererOutput
	ToPostRendererOutputWithContext(context.Context) PostRendererOutput
}

// Specification defining the post-renderer to use.
type PostRendererArgs struct {
	// Arguments to pass to the post-renderer command.
	Args pulumi.StringArrayInput `pulumi:"args"`
	// Path to an executable to be used for post rendering.
	Command pulumi.StringInput `pulumi:"command"`
}

func (PostRendererArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*PostRenderer)(nil)).Elem()
}

func (i PostRendererArgs) ToPostRendererOutput() PostRendererOutput {
	return i.ToPostRendererOutputWithContext(context.Background())
}

func (i PostRendererArgs) ToPostRendererOutputWithContext(ctx context.Context) PostRendererOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PostRendererOutput)
}

func (i PostRendererArgs) ToPostRendererPtrOutput() PostRendererPtrOutput {
	return i.ToPostRendererPtrOutputWithContext(context.Background())
}

func (i PostRendererArgs) ToPostRendererPtrOutputWithContext(ctx context.Context) PostRendererPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PostRendererOutput).ToPostRendererPtrOutputWithContext(ctx)
}

// PostRendererPtrInput is an input type that accepts PostRendererArgs, PostRendererPtr and PostRendererPtrOutput values.
// You can construct a concrete instance of `PostRendererPtrInput` via:
//
//	        PostRendererArgs{...}
//
//	or:
//
//	        nil
type PostRendererPtrInput interface {
	pulumi.Input

	ToPostRendererPtrOutput() PostRendererPtrOutput
	ToPostRendererPtrOutputWithContext(context.Context) PostRendererPtrOutput
}

type postRendererPtrType PostRendererArgs

func PostRendererPtr(v *PostRendererArgs) PostRendererPtrInput {
	return (*postRendererPtrType)(v)
}

func (*postRendererPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**PostRenderer)(nil)).Elem()
}

func (i *postRendererPtrType) ToPostRendererPtrOutput() PostRendererPtrOutput {
	return i.ToPostRendererPtrOutputWithContext(context.Background())
}

func (i *postRendererPtrType) ToPostRendererPtrOutputWithContext(ctx context.Context) PostRendererPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PostRendererPtrOutput)
}

// Specification defining the post-renderer to use.
type PostRendererOutput struct{ *pulumi.OutputState }

func (PostRendererOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*PostRenderer)(nil)).Elem()
}

func (o PostRendererOutput) ToPostRendererOutput() PostRendererOutput {
	return o
}

func (o PostRendererOutput) ToPostRendererOutputWithContext(ctx context.Context) PostRendererOutput {
	return o
}

func (o PostRendererOutput) ToPostRendererPtrOutput() PostRendererPtrOutput {
	return o.ToPostRendererPtrOutputWithContext(context.Background())
}

func (o PostRendererOutput) ToPostRendererPtrOutputWithContext(ctx context.Context) PostRendererPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v PostRenderer) *PostRenderer {
		return &v
	}).(PostRendererPtrOutput)
}

// Arguments to pass to the post-renderer command.
func (o PostRendererOutput) Args() pulumi.StringArrayOutput {
	return o.ApplyT(func(v PostRenderer) []string { return v.Args }).(pulumi.StringArrayOutput)
}

// Path to an executable to be used for post rendering.
func (o PostRendererOutput) Command() pulumi.StringOutput {
	return o.ApplyT(func(v PostRenderer) string { return v.Command }).(pulumi.StringOutput)
}

type PostRendererPtrOutput struct{ *pulumi.OutputState }

func (PostRendererPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**PostRenderer)(nil)).Elem()
}

func (o PostRendererPtrOutput) ToPostRendererPtrOutput() PostRendererPtrOutput {
	return o
}

func (o PostRendererPtrOutput) ToPostRendererPtrOutputWithContext(ctx context.Context) PostRendererPtrOutput {
	return o
}

func (o PostRendererPtrOutput) Elem() PostRendererOutput {
	return o.ApplyT(func(v *PostRenderer) PostRenderer {
		if v != nil {
			return *v
		}
		var ret PostRenderer
		return ret
	}).(PostRendererOutput)
}

// Arguments to pass to the post-renderer command.
func (o PostRendererPtrOutput) Args() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *PostRenderer) []string {
		if v == nil {
			return nil
		}
		return v.Args
	}).(pulumi.StringArrayOutput)
}

// Path to an executable to be used for post rendering.
func (o PostRendererPtrOutput) Command() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *PostRenderer) *string {
		if v == nil {
			return nil
		}
		return &v.Command
	}).(pulumi.StringPtrOutput)
}

// Specification defining the Helm chart repository to use.
type RepositoryOpts struct {
	// The Repository's CA File
	CaFile pulumi.AssetOrArchive `pulumi:"caFile"`
	// The repository's cert file
	CertFile pulumi.AssetOrArchive `pulumi:"certFile"`
	// The repository's cert key file
	KeyFile pulumi.AssetOrArchive `pulumi:"keyFile"`
	// Password for HTTP basic authentication
	Password *string `pulumi:"password"`
	// Repository where to locate the requested chart. If it's a URL the chart is installed without installing the repository.
	Repo *string `pulumi:"repo"`
	// Username for HTTP basic authentication
	Username *string `pulumi:"username"`
}

// RepositoryOptsInput is an input type that accepts RepositoryOptsArgs and RepositoryOptsOutput values.
// You can construct a concrete instance of `RepositoryOptsInput` via:
//
//	RepositoryOptsArgs{...}
type RepositoryOptsInput interface {
	pulumi.Input

	ToRepositoryOptsOutput() RepositoryOptsOutput
	ToRepositoryOptsOutputWithContext(context.Context) RepositoryOptsOutput
}

// Specification defining the Helm chart repository to use.
type RepositoryOptsArgs struct {
	// The Repository's CA File
	CaFile pulumi.AssetOrArchiveInput `pulumi:"caFile"`
	// The repository's cert file
	CertFile pulumi.AssetOrArchiveInput `pulumi:"certFile"`
	// The repository's cert key file
	KeyFile pulumi.AssetOrArchiveInput `pulumi:"keyFile"`
	// Password for HTTP basic authentication
	Password pulumi.StringPtrInput `pulumi:"password"`
	// Repository where to locate the requested chart. If it's a URL the chart is installed without installing the repository.
	Repo pulumi.StringPtrInput `pulumi:"repo"`
	// Username for HTTP basic authentication
	Username pulumi.StringPtrInput `pulumi:"username"`
}

func (RepositoryOptsArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*RepositoryOpts)(nil)).Elem()
}

func (i RepositoryOptsArgs) ToRepositoryOptsOutput() RepositoryOptsOutput {
	return i.ToRepositoryOptsOutputWithContext(context.Background())
}

func (i RepositoryOptsArgs) ToRepositoryOptsOutputWithContext(ctx context.Context) RepositoryOptsOutput {
	return pulumi.ToOutputWithContext(ctx, i).(RepositoryOptsOutput)
}

func (i RepositoryOptsArgs) ToRepositoryOptsPtrOutput() RepositoryOptsPtrOutput {
	return i.ToRepositoryOptsPtrOutputWithContext(context.Background())
}

func (i RepositoryOptsArgs) ToRepositoryOptsPtrOutputWithContext(ctx context.Context) RepositoryOptsPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(RepositoryOptsOutput).ToRepositoryOptsPtrOutputWithContext(ctx)
}

// RepositoryOptsPtrInput is an input type that accepts RepositoryOptsArgs, RepositoryOptsPtr and RepositoryOptsPtrOutput values.
// You can construct a concrete instance of `RepositoryOptsPtrInput` via:
//
//	        RepositoryOptsArgs{...}
//
//	or:
//
//	        nil
type RepositoryOptsPtrInput interface {
	pulumi.Input

	ToRepositoryOptsPtrOutput() RepositoryOptsPtrOutput
	ToRepositoryOptsPtrOutputWithContext(context.Context) RepositoryOptsPtrOutput
}

type repositoryOptsPtrType RepositoryOptsArgs

func RepositoryOptsPtr(v *RepositoryOptsArgs) RepositoryOptsPtrInput {
	return (*repositoryOptsPtrType)(v)
}

func (*repositoryOptsPtrType) ElementType() reflect.Type {
	return reflect.TypeOf((**RepositoryOpts)(nil)).Elem()
}

func (i *repositoryOptsPtrType) ToRepositoryOptsPtrOutput() RepositoryOptsPtrOutput {
	return i.ToRepositoryOptsPtrOutputWithContext(context.Background())
}

func (i *repositoryOptsPtrType) ToRepositoryOptsPtrOutputWithContext(ctx context.Context) RepositoryOptsPtrOutput {
	return pulumi.ToOutputWithContext(ctx, i).(RepositoryOptsPtrOutput)
}

// Specification defining the Helm chart repository to use.
type RepositoryOptsOutput struct{ *pulumi.OutputState }

func (RepositoryOptsOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*RepositoryOpts)(nil)).Elem()
}

func (o RepositoryOptsOutput) ToRepositoryOptsOutput() RepositoryOptsOutput {
	return o
}

func (o RepositoryOptsOutput) ToRepositoryOptsOutputWithContext(ctx context.Context) RepositoryOptsOutput {
	return o
}

func (o RepositoryOptsOutput) ToRepositoryOptsPtrOutput() RepositoryOptsPtrOutput {
	return o.ToRepositoryOptsPtrOutputWithContext(context.Background())
}

func (o RepositoryOptsOutput) ToRepositoryOptsPtrOutputWithContext(ctx context.Context) RepositoryOptsPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v RepositoryOpts) *RepositoryOpts {
		return &v
	}).(RepositoryOptsPtrOutput)
}

// The Repository's CA File
func (o RepositoryOptsOutput) CaFile() pulumi.AssetOrArchiveOutput {
	return o.ApplyT(func(v RepositoryOpts) pulumi.AssetOrArchive { return v.CaFile }).(pulumi.AssetOrArchiveOutput)
}

// The repository's cert file
func (o RepositoryOptsOutput) CertFile() pulumi.AssetOrArchiveOutput {
	return o.ApplyT(func(v RepositoryOpts) pulumi.AssetOrArchive { return v.CertFile }).(pulumi.AssetOrArchiveOutput)
}

// The repository's cert key file
func (o RepositoryOptsOutput) KeyFile() pulumi.AssetOrArchiveOutput {
	return o.ApplyT(func(v RepositoryOpts) pulumi.AssetOrArchive { return v.KeyFile }).(pulumi.AssetOrArchiveOutput)
}

// Password for HTTP basic authentication
func (o RepositoryOptsOutput) Password() pulumi.StringPtrOutput {
	return o.ApplyT(func(v RepositoryOpts) *string { return v.Password }).(pulumi.StringPtrOutput)
}

// Repository where to locate the requested chart. If it's a URL the chart is installed without installing the repository.
func (o RepositoryOptsOutput) Repo() pulumi.StringPtrOutput {
	return o.ApplyT(func(v RepositoryOpts) *string { return v.Repo }).(pulumi.StringPtrOutput)
}

// Username for HTTP basic authentication
func (o RepositoryOptsOutput) Username() pulumi.StringPtrOutput {
	return o.ApplyT(func(v RepositoryOpts) *string { return v.Username }).(pulumi.StringPtrOutput)
}

type RepositoryOptsPtrOutput struct{ *pulumi.OutputState }

func (RepositoryOptsPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**RepositoryOpts)(nil)).Elem()
}

func (o RepositoryOptsPtrOutput) ToRepositoryOptsPtrOutput() RepositoryOptsPtrOutput {
	return o
}

func (o RepositoryOptsPtrOutput) ToRepositoryOptsPtrOutputWithContext(ctx context.Context) RepositoryOptsPtrOutput {
	return o
}

func (o RepositoryOptsPtrOutput) Elem() RepositoryOptsOutput {
	return o.ApplyT(func(v *RepositoryOpts) RepositoryOpts {
		if v != nil {
			return *v
		}
		var ret RepositoryOpts
		return ret
	}).(RepositoryOptsOutput)
}

// The Repository's CA File
func (o RepositoryOptsPtrOutput) CaFile() pulumi.AssetOrArchiveOutput {
	return o.ApplyT(func(v *RepositoryOpts) pulumi.AssetOrArchive {
		if v == nil {
			return nil
		}
		return v.CaFile
	}).(pulumi.AssetOrArchiveOutput)
}

// The repository's cert file
func (o RepositoryOptsPtrOutput) CertFile() pulumi.AssetOrArchiveOutput {
	return o.ApplyT(func(v *RepositoryOpts) pulumi.AssetOrArchive {
		if v == nil {
			return nil
		}
		return v.CertFile
	}).(pulumi.AssetOrArchiveOutput)
}

// The repository's cert key file
func (o RepositoryOptsPtrOutput) KeyFile() pulumi.AssetOrArchiveOutput {
	return o.ApplyT(func(v *RepositoryOpts) pulumi.AssetOrArchive {
		if v == nil {
			return nil
		}
		return v.KeyFile
	}).(pulumi.AssetOrArchiveOutput)
}

// Password for HTTP basic authentication
func (o RepositoryOptsPtrOutput) Password() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *RepositoryOpts) *string {
		if v == nil {
			return nil
		}
		return v.Password
	}).(pulumi.StringPtrOutput)
}

// Repository where to locate the requested chart. If it's a URL the chart is installed without installing the repository.
func (o RepositoryOptsPtrOutput) Repo() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *RepositoryOpts) *string {
		if v == nil {
			return nil
		}
		return v.Repo
	}).(pulumi.StringPtrOutput)
}

// Username for HTTP basic authentication
func (o RepositoryOptsPtrOutput) Username() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *RepositoryOpts) *string {
		if v == nil {
			return nil
		}
		return v.Username
	}).(pulumi.StringPtrOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*PostRendererInput)(nil)).Elem(), PostRendererArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*PostRendererPtrInput)(nil)).Elem(), PostRendererArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*RepositoryOptsInput)(nil)).Elem(), RepositoryOptsArgs{})
	pulumi.RegisterInputType(reflect.TypeOf((*RepositoryOptsPtrInput)(nil)).Elem(), RepositoryOptsArgs{})
	pulumi.RegisterOutputType(PostRendererOutput{})
	pulumi.RegisterOutputType(PostRendererPtrOutput{})
	pulumi.RegisterOutputType(RepositoryOptsOutput{})
	pulumi.RegisterOutputType(RepositoryOptsPtrOutput{})
}
