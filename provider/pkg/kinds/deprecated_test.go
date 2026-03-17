// Copyright 2016-2021, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kinds

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/cluster"
)

// Version constants used by tests.
var (
	v18  = cluster.ServerVersion{Major: 1, Minor: 8}
	v19  = cluster.ServerVersion{Major: 1, Minor: 9}
	v114 = cluster.ServerVersion{Major: 1, Minor: 14}
	v116 = cluster.ServerVersion{Major: 1, Minor: 16}
	v117 = cluster.ServerVersion{Major: 1, Minor: 17}
	v118 = cluster.ServerVersion{Major: 1, Minor: 18}
	v120 = cluster.ServerVersion{Major: 1, Minor: 20}
	v121 = cluster.ServerVersion{Major: 1, Minor: 21}
	v122 = cluster.ServerVersion{Major: 1, Minor: 22}
	v124 = cluster.ServerVersion{Major: 1, Minor: 24}
	v125 = cluster.ServerVersion{Major: 1, Minor: 25}
	v126 = cluster.ServerVersion{Major: 1, Minor: 26}
	v129 = cluster.ServerVersion{Major: 1, Minor: 29}
	v132 = cluster.ServerVersion{Major: 1, Minor: 32}
	v133 = cluster.ServerVersion{Major: 1, Minor: 33}
)

func TestDeprecatedApiVersion(t *testing.T) {
	tests := []struct {
		gvk     schema.GroupVersionKind
		version *cluster.ServerVersion
		want    bool
	}{
		{toGVK(AdmissionregistrationV1B1, MutatingWebhookConfiguration), nil, true},
		{toGVK(AdmissionregistrationV1B1, MutatingWebhookConfiguration), &v114, false},
		{toGVK(AdmissionregistrationV1B1, MutatingWebhookConfiguration), &v116, true},
		{toGVK(AdmissionregistrationV1B1, ValidatingWebhookConfiguration), &v118, true},
		{toGVK(ApiextensionsV1B1, CustomResourceDefinition), nil, true},
		{toGVK(ApiregistrationV1B1, APIService), nil, true},
		{toGVK(AppsV1, Deployment), nil, false},
		{toGVK(AppsV1B1, Deployment), nil, true},
		{toGVK(AppsV1B2, Deployment), nil, true},
		{toGVK(AutoscalingV2B1, HorizontalPodAutoscaler), nil, true},
		{toGVK(BatchV2A1, CronJob), &v121, true},
		{toGVK(CoordinationV1B1, Lease), nil, true},
		{toGVK(CoreV1, Endpoints), nil, true},
		{toGVK(DiscoveryV1B1, EndpointSlice), &v121, true},
		{toGVK(ExtensionsV1B1, DaemonSet), nil, true},
		{toGVK(ExtensionsV1B1, Deployment), nil, true},
		{toGVK(ExtensionsV1B1, Ingress), nil, true},
		{toGVK(ExtensionsV1B1, NetworkPolicy), nil, true},
		{toGVK(ExtensionsV1B1, PodSecurityPolicy), nil, true},
		{toGVK(ExtensionsV1B1, ReplicaSet), nil, true},
		{toGVK(RbacV1A1, ClusterRole), nil, true},
		{toGVK(RbacV1A1, ClusterRoleBinding), nil, true},
		{toGVK(RbacV1A1, Role), nil, true},
		{toGVK(RbacV1A1, RoleBinding), nil, true},
		{toGVK(RbacV1B1, ClusterRole), nil, true},
		{toGVK(RbacV1B1, ClusterRoleBinding), nil, true},
		{toGVK(RbacV1B1, Role), nil, true},
		{toGVK(RbacV1B1, RoleBinding), nil, true},
		{toGVK(ResourceV1B1, DeviceClass), nil, true},
		{toGVK(SchedulingV1A1, PriorityClass), nil, true},
		{toGVK(SchedulingV1B1, PriorityClass), nil, true},
		{toGVK(StorageV1A1, CSIStorageCapacity), nil, true},
		{toGVK(StorageV1A1, VolumeAttachment), nil, true},
		{toGVK(StorageV1B1, CSIDriver), nil, true},
		{toGVK(StorageV1B1, CSIDriver), &v118, true},
		{toGVK(StorageV1B1, CSIDriver), &v117, false},
		{toGVK(StorageV1B1, CSIDriver), &v116, false},
		{toGVK(StorageV1B1, CSINode), &v118, true},
		{toGVK(StorageV1B1, CSINode), &v117, true},
		{toGVK(StorageV1B1, CSINode), &v116, false},
		{toGVK(StorageV1B1, StorageClass), nil, true},
		{toGVK(StorageV1B1, StorageClass), &v114, true},
		{toGVK(StorageV1B1, VolumeAttachment), nil, true},
		{toGVK(StorageV1, CSINode), &v118, false},
		{toGVK(StorageV1, CSINode), &v120, false},
		{toGVK(AppsV1, Deployment), &v18, false},
		{toGVK(AppsV1, Deployment), &v19, false},
		{toGVK(StorageV1A1, VolumeAttributesClass), &v133, false},
		{toGVK(StorageV1B1, VolumeAttributesClass), &v133, false},
	}
	for _, tt := range tests {
		t.Run(tt.gvk.String(), func(t *testing.T) {
			if got := DeprecatedAPIVersion(tt.gvk, tt.version); got != tt.want {
				t.Errorf("DeprecatedAPIVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExistsInVersion(t *testing.T) {
	tests := []struct {
		gvk     schema.GroupVersionKind
		version *cluster.ServerVersion
		want    bool
	}{
		{toGVK(StorageV1, CSINode), &v118, true},
		{toGVK(StorageV1, CSINode), &v117, true},
		{toGVK(StorageV1, CSINode), &v116, false},
		{schema.GroupVersionKind{}, nil, false},
		{toGVK(AppsV1, Deployment), &v18, false},
		{toGVK(AppsV1, Deployment), &v19, true},
		{toGVK(ResourceV1B1, DeviceClass), &v132, true},
		{toGVK(ResourceV1B2, DeviceClass), &v132, false},
		{toGVK(ResourceV1B2, DeviceClass), &v133, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.gvk.String(), func(t *testing.T) {
			if got := ExistsInVersion(&tt.gvk, tt.version); got != tt.want {
				t.Errorf("ExistsInVersion() = %v, want %v", got, tt.want)
			}
		})
	}
	t.Run("nil GVK and version", func(t *testing.T) {
		if got := ExistsInVersion(nil, nil); got != false {
			t.Errorf("ExistsInVersion() = %v, want %v", got, false)
		}
	})
	t.Run("nil GVK only", func(t *testing.T) {
		if got := ExistsInVersion(nil, &v118); got != false {
			t.Errorf("ExistsInVersion() = %v, want %v", got, false)
		}
	})
}

func TestGvkFromStr(t *testing.T) {
	tests := []struct {
		gvkString string
		want      schema.GroupVersionKind
	}{
		{"storage.k8s.io/v1/CSINode", schema.GroupVersionKind{Group: "storage.k8s.io", Version: "v1", Kind: "CSINode"}},
		{
			"networking.k8s.io/v1beta1/IngressList",
			schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1beta1", Kind: "IngressList"},
		},
		{"something/else", schema.GroupVersionKind{}},
	}
	for _, tt := range tests {
		t.Run(tt.gvkString, func(t *testing.T) {
			if got := gvkFromStr(tt.gvkString); got != tt.want {
				t.Errorf("TestGvkFromStr() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestSuggestedApiVersion(t *testing.T) {
	wantStr := func(gv groupVersion, kind Kind) string {
		return gvkStr(toGVK(gv, kind))
	}

	tests := []struct {
		gvk  schema.GroupVersionKind
		want string
	}{
		// Deprecated ApiVersions return the suggested version string.
		{
			toGVK(AdmissionregistrationV1B1, MutatingWebhookConfiguration),
			wantStr(AdmissionregistrationV1, MutatingWebhookConfiguration),
		},
		{
			toGVK(AdmissionregistrationV1B1, ValidatingWebhookConfiguration),
			wantStr(AdmissionregistrationV1, ValidatingWebhookConfiguration),
		},
		{toGVK(ApiextensionsV1B1, CustomResourceDefinition), wantStr(ApiextensionsV1, CustomResourceDefinition)},
		{toGVK(ApiregistrationV1B1, APIService), wantStr(ApiregistrationV1, APIService)},
		{toGVK(ApiregistrationV1B1, APIServiceList), wantStr(ApiregistrationV1, APIServiceList)},
		{toGVK(AppsV1B1, Deployment), wantStr(AppsV1, Deployment)},
		{toGVK(AppsV1B2, Deployment), wantStr(AppsV1, Deployment)},
		{toGVK(AutoscalingV2B1, HorizontalPodAutoscaler), wantStr(AutoscalingV2, HorizontalPodAutoscaler)},
		{toGVK(BatchV2A1, CronJob), wantStr(BatchV1, CronJob)},
		{toGVK(BatchV1B1, CronJob), wantStr(BatchV1, CronJob)},
		{toGVK(CoordinationV1B1, Lease), wantStr(CoordinationV1, Lease)},
		{toGVK(DiscoveryV1B1, EndpointSlice), wantStr(DiscoveryV1, EndpointSlice)},
		{toGVK(ExtensionsV1B1, DaemonSet), wantStr(AppsV1, DaemonSet)},
		{toGVK(ExtensionsV1B1, Deployment), wantStr(AppsV1, Deployment)},
		{toGVK(ExtensionsV1B1, DeploymentList), wantStr(AppsV1, DeploymentList)},
		{toGVK(ExtensionsV1B1, Ingress), wantStr(NetworkingV1, Ingress)},
		{toGVK(ExtensionsV1B1, IngressList), wantStr(NetworkingV1, IngressList)},
		{toGVK(ExtensionsV1B1, NetworkPolicy), wantStr(NetworkingV1, NetworkPolicy)},
		{toGVK(ExtensionsV1B1, PodSecurityPolicy), wantStr(PolicyV1B1, PodSecurityPolicy)},
		{toGVK(ExtensionsV1B1, ReplicaSet), wantStr(AppsV1, ReplicaSet)},
		{toGVK(RbacV1A1, ClusterRole), wantStr(RbacV1, ClusterRole)},
		{toGVK(RbacV1A1, ClusterRoleBinding), wantStr(RbacV1, ClusterRoleBinding)},
		{toGVK(RbacV1B1, ClusterRole), wantStr(RbacV1, ClusterRole)},
		{toGVK(RbacV1B1, ClusterRoleBinding), wantStr(RbacV1, ClusterRoleBinding)},
		{toGVK(SchedulingV1A1, PriorityClass), wantStr(SchedulingV1, PriorityClass)},
		{toGVK(SchedulingV1B1, PriorityClass), wantStr(SchedulingV1, PriorityClass)},
		{toGVK(StorageV1A1, VolumeAttachment), wantStr(StorageV1, VolumeAttachment)},
		{toGVK(StorageV1B1, CSIDriver), wantStr(StorageV1, CSIDriver)},
		{toGVK(StorageV1B1, CSINode), wantStr(StorageV1, CSINode)},
		{toGVK(StorageV1B1, StorageClass), wantStr(StorageV1, StorageClass)},
		{toGVK(StorageV1B1, VolumeAttachment), wantStr(StorageV1, VolumeAttachment)},
		// Current ApiVersions return the same version string.
		{toGVK(AppsV1, Deployment), wantStr(AppsV1, Deployment)},
	}
	for _, tt := range tests {
		t.Run(tt.gvk.String(), func(t *testing.T) {
			if got := SuggestedAPIVersion(tt.gvk); got != tt.want {
				t.Errorf("SuggestedAPIVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemovedInVersion(t *testing.T) {
	tests := []struct {
		gvk         schema.GroupVersionKind
		wantVersion *cluster.ServerVersion
	}{
		{toGVK(AdmissionregistrationV1B1, MutatingWebhookConfiguration), &v122},
		{toGVK(ApiextensionsV1B1, CustomResourceDefinition), &v122},
		{toGVK(AppsV1B1, Deployment), &v116},
		{toGVK(AppsV1B2, Deployment), &v116},
		{toGVK(BatchV2A1, CronJob), &v121},
		{toGVK(BatchV1B1, CronJob), &v125},
		{toGVK(CoordinationV1B1, Lease), &v122},
		{toGVK(DiscoveryV1B1, EndpointSlice), &v125},
		{toGVK(ExtensionsV1B1, Deployment), &v116},
		{toGVK(ExtensionsV1B1, DeploymentList), &v116},
		{toGVK(ExtensionsV1B1, Ingress), &v122},
		{toGVK(ExtensionsV1B1, IngressList), &v122},
		{toGVK(ExtensionsV1B1, PodSecurityPolicy), &v116},
		{toGVK(FlowcontrolV1B1, FlowSchema), &v126},
		{toGVK(FlowcontrolV1B1, PriorityLevelConfiguration), &v126},
		{toGVK(FlowcontrolV1B2, FlowSchema), &v129},
		{toGVK(FlowcontrolV1B2, PriorityLevelConfiguration), &v129},
		{toGVK(FlowcontrolV1B3, FlowSchema), &v132},
		{toGVK(FlowcontrolV1B3, PriorityLevelConfiguration), &v132},
		{toGVK(PolicyV1, PodSecurityPolicy), &v125},
		{toGVK(PolicyV1B1, PodDisruptionBudget), &v125},
		{toGVK(PolicyV1B1, PodSecurityPolicy), &v125},
		{toGVK(RbacV1A1, ClusterRole), &v120},
		{toGVK(RbacV1B1, ClusterRole), &v122},
		{toGVK(SchedulingV1A1, PriorityClass), &v117},
		{toGVK(SchedulingV1B1, PriorityClass), &v122},
		{toGVK(StorageV1A1, CSIStorageCapacity), &v124},
	}
	for _, tt := range tests {
		t.Run(tt.gvk.String(), func(t *testing.T) {
			got := RemovedInVersion(tt.gvk)
			if !reflect.DeepEqual(got, tt.wantVersion) {
				t.Errorf("RemovedInVersion() got = %v, want %v", got, tt.wantVersion)
			}
		})
	}
}

// TestSetReplacementPanicsOnMissingEntry verifies that setReplacement panics
// when called with a GVK not in the deprecations map, rather than silently
// no-oping (which previously caused DeviceTaintRule to have no deprecation data).
func TestSetReplacementPanicsOnMissingEntry(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("setReplacement did not panic for a missing GVK")
		}
	}()
	setReplacement("nonexistent.k8s.io", "v1", "FakeKind",
		schema.GroupVersionKind{Group: "nonexistent.k8s.io", Version: "v2", Kind: "FakeKind"})
}

// TestDeviceTaintRuleV1Alpha3HasDeprecationData verifies that DeviceTaintRule
// in v1alpha3 is picked up by the scheme iteration with introduced/removed data.
func TestDeviceTaintRuleV1Alpha3HasDeprecationData(t *testing.T) {
	gvk := schema.GroupVersionKind{Group: "resource.k8s.io", Version: "v1alpha3", Kind: "DeviceTaintRule"}
	info, ok := deprecations[gvk]
	if !ok {
		t.Fatalf("DeviceTaintRule v1alpha3 not found in deprecations map")
	}
	if info.IntroducedMajor == 0 && info.IntroducedMinor == 0 {
		t.Error("DeviceTaintRule v1alpha3 has no IntroducedIn data")
	}
	if info.RemovedMajor == 0 && info.RemovedMinor == 0 {
		t.Error("DeviceTaintRule v1alpha3 has no RemovedIn data")
	}
}

// TestAllSetReplacementTargetsExist verifies that every GVK with a replacement
// pointer actually has the replacement target's group/version/kind set (not zero-value).
func TestAllReplacementsAreValid(t *testing.T) {
	for gvk, info := range deprecations {
		if info.Replacement != nil {
			if info.Replacement.Kind == "" {
				t.Errorf("%s has a replacement with empty Kind", gvk)
			}
			if info.Replacement.Version == "" {
				t.Errorf("%s has a replacement with empty Version", gvk)
			}
		}
	}
}

// TestShortGroupAliasesInheritReplacement verifies that truncated-group aliases
// (e.g., "resource" for "resource.k8s.io") carry the same Replacement pointer
// as the canonical full-group entry.
func TestShortGroupAliasesInheritReplacement(t *testing.T) {
	// DeviceClass resource.k8s.io/v1beta1 has a setReplacement pointing to v1beta2.
	// The short alias "resource/v1beta1/DeviceClass" should also have it.
	full := schema.GroupVersionKind{Group: "resource.k8s.io", Version: "v1beta1", Kind: "DeviceClass"}
	short := schema.GroupVersionKind{Group: "resource", Version: "v1beta1", Kind: "DeviceClass"}

	fullInfo, ok := deprecations[full]
	if !ok {
		t.Fatalf("%s not found in deprecations map", full)
	}
	if fullInfo.Replacement == nil {
		t.Fatalf("%s has no replacement", full)
	}

	shortInfo, ok := deprecations[short]
	if !ok {
		t.Fatalf("%s (short-group alias) not found in deprecations map", short)
	}
	if shortInfo.Replacement == nil {
		t.Errorf("%s (short-group alias) has Replacement == nil; expected it to inherit from %s", short, full)
	} else if *shortInfo.Replacement != *fullInfo.Replacement {
		t.Errorf("%s replacement = %s, want %s", short, *shortInfo.Replacement, *fullInfo.Replacement)
	}
}

func TestRemovedApiVersion(t *testing.T) {
	type args struct {
		gvk     schema.GroupVersionKind
		version cluster.ServerVersion
	}
	tests := []struct {
		name        string
		args        args
		wantRemoved bool
		wantVersion *cluster.ServerVersion
	}{
		{"API exists", args{
			schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
			cluster.ServerVersion{Major: 1, Minor: 16}}, false, nil},
		{"API removed", args{
			schema.GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Deployment"},
			cluster.ServerVersion{Major: 1, Minor: 16}},
			true, &cluster.ServerVersion{Major: 1, Minor: 16}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := RemovedAPIVersion(tt.args.gvk, tt.args.version)
			if got != tt.wantRemoved {
				t.Errorf("RemovedAPIVersion() got = %v, want %v", got, tt.wantRemoved)
			}
			if !reflect.DeepEqual(got1, tt.wantVersion) {
				t.Errorf("RemovedAPIVersion() got1 = %v, want %v", got1, tt.wantVersion)
			}
		})
	}
}
