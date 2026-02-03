package kinds

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

func TestIsPatchURN(t *testing.T) {
	tests := []struct {
		name string
		urn  resource.URN
		kind string
		want bool
	}{
		// Simple resources
		{
			name: "Simple patch URN - Happy Path",
			urn:  resource.NewURN("test", "test", "", "kubernetes:apps/v1:DeploymentPatch", "test"),
			kind: "Deployment",
			want: true,
		},
		{
			name: "CustomResource with Patch suffix that is not a patch resource",
			urn:  resource.NewURN("test", "test", "", "kubernetes:kuma.io/v1alpha1:MeshProxyPatch", "test"),
			kind: "MeshProxyPatch",
			want: false,
		},
		{
			name: "CustomResource with Patch suffix and is a patch resource as well",
			urn:  resource.NewURN("test", "test", "", "kubernetes:kuma.io/v1alpha1:MeshProxyPatchPatch", "test"),
			kind: "MeshProxyPatch",
			want: true,
		},
		// Component resources
		{
			"Patch resource within a Component Resource",
			resource.URN(
				"urn:pulumi:dev::kubernetes-ts::my:component:Resource$kubernetes:apps/v1:DaemonSetPatch::k8s-daemonset-patch-child",
			),
			"DaemonSet",
			true,
		},
		{
			"Non-Patch DaemonSet resource",
			resource.URN("urn:pulumi:dev::kubernetes-ts::kubernetes:apps/v1:DaemonSet::k8s-daemonset"),
			"DaemonSet",
			false,
		},
		{
			"Custom Resource with Patch suffix",
			resource.URN(
				"urn:pulumi:dev::kubernetes-ts::my:component:Resource$kubernetes:kuma.io/v1alpha1:MeshProxyPatch::k8s-meshproxy-patch-child",
			),
			"MeshProxyPatch",
			false,
		},
		{
			"Custom Resource with Patch suffix that is a patch resource",
			resource.URN(
				"urn:pulumi:dev::kubernetes-ts::my:component:Resource$kubernetes:kuma.io/v1alpha1:MeshProxyPatchPatch::k8s-meshproxy-patch-child",
			),
			"MeshProxyPatch",
			true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsPatchResource(tc.urn, tc.kind); got != tc.want {
				t.Errorf("IsPatchURN() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIsListURN(t *testing.T) {
	tests := []struct {
		name string
		urn  resource.URN
		want bool
	}{
		{
			"Simple List resource - Happy Path",
			resource.URN("urn:pulumi:dev::kubernetes-ts::kubernetes:apps/v1:DaemonSetList::k8s-daemonset-list"),
			true,
		},
		{
			"List resource within a Component Resource",
			resource.URN(
				"urn:pulumi:dev::kubernetes-ts::my:component:Resource$kubernetes:apps/v1:DaemonSetList::k8s-daemonset-list-child",
			),
			true,
		},
		{
			"Non-List DaemonSet resource",
			resource.URN("urn:pulumi:dev::kubernetes-ts::kubernetes:apps/v1:DaemonSet::k8s-daemonset"),
			false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsListURN(tc.urn); got != tc.want {
				t.Errorf("IsListURN() = %v, want %v", got, tc.want)
			}
		})
	}
}
