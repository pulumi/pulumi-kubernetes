package kinds

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

func TestIsPatchURN(t *testing.T) {
	tests := []struct {
		name string
		urn  resource.URN
		want bool
	}{
		{
			"Simple Patch resource - Happy Path",
			resource.URN("urn:pulumi:dev::kubernetes-ts::kubernetes:apps/v1:DaemonSetPatch::k8s-daemonset-patch"),
			true,
		},
		{
			"Patch resource within a Component Resource",
			resource.URN("urn:pulumi:dev::kubernetes-ts::my:component:Resource$kubernetes:apps/v1:DaemonSetPatch::k8s-daemonset-patch-child"),
			true,
		},
		{
			"Non-Patch DaemonSet resource",
			resource.URN("urn:pulumi:dev::kubernetes-ts::kubernetes:apps/v1:DaemonSet::k8s-daemonset"),
			false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsPatchURN(tc.urn); got != tc.want {
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
			resource.URN("urn:pulumi:dev::kubernetes-ts::my:component:Resource$kubernetes:apps/v1:DaemonSetList::k8s-daemonset-list-child"),
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
