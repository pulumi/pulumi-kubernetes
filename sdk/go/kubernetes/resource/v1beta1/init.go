// Code generated by pulumigen DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package v1beta1

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/utilities"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type module struct {
	version semver.Version
}

func (m *module) Version() semver.Version {
	return m.version
}

func (m *module) Construct(ctx *pulumi.Context, name, typ, urn string) (r pulumi.Resource, err error) {
	switch typ {
	case "kubernetes:resource.k8s.io/v1beta1:DeviceClass":
		r = &DeviceClass{}
	case "kubernetes:resource.k8s.io/v1beta1:DeviceClassList":
		r = &DeviceClassList{}
	case "kubernetes:resource.k8s.io/v1beta1:DeviceClassPatch":
		r = &DeviceClassPatch{}
	case "kubernetes:resource.k8s.io/v1beta1:ResourceClaim":
		r = &ResourceClaim{}
	case "kubernetes:resource.k8s.io/v1beta1:ResourceClaimList":
		r = &ResourceClaimList{}
	case "kubernetes:resource.k8s.io/v1beta1:ResourceClaimPatch":
		r = &ResourceClaimPatch{}
	case "kubernetes:resource.k8s.io/v1beta1:ResourceClaimTemplate":
		r = &ResourceClaimTemplate{}
	case "kubernetes:resource.k8s.io/v1beta1:ResourceClaimTemplateList":
		r = &ResourceClaimTemplateList{}
	case "kubernetes:resource.k8s.io/v1beta1:ResourceClaimTemplatePatch":
		r = &ResourceClaimTemplatePatch{}
	case "kubernetes:resource.k8s.io/v1beta1:ResourceSlice":
		r = &ResourceSlice{}
	case "kubernetes:resource.k8s.io/v1beta1:ResourceSliceList":
		r = &ResourceSliceList{}
	case "kubernetes:resource.k8s.io/v1beta1:ResourceSlicePatch":
		r = &ResourceSlicePatch{}
	default:
		return nil, fmt.Errorf("unknown resource type: %s", typ)
	}

	err = ctx.RegisterResource(typ, name, nil, r, pulumi.URN_(urn))
	return
}

func init() {
	version, err := utilities.PkgVersion()
	if err != nil {
		version = semver.Version{Major: 1}
	}
	pulumi.RegisterResourceModule(
		"kubernetes",
		"resource.k8s.io/v1beta1",
		&module{version},
	)
}