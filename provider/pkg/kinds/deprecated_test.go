// Copyright 2016-2019, Pulumi Corporation.
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

	"github.com/pulumi/pulumi-kubernetes/provider/v2/pkg/cluster"
	. "k8s.io/apimachinery/pkg/runtime/schema"
)

func TestDeprecatedApiVersion(t *testing.T) {
	tests := []struct {
		gvk  GroupVersionKind
		want bool
	}{
		{toGVK(AdmissionregistrationV1B1, MutatingWebhookConfiguration), true},
		{toGVK(AdmissionregistrationV1B1, ValidatingWebhookConfiguration), true},
		{toGVK(ApiextensionsV1B1, CustomResourceDefinition), true},
		{toGVK(ApiregistrationV1B1, APIService), true},
		{toGVK(AppsV1, Deployment), false},
		{toGVK(AppsV1B1, Deployment), true},
		{toGVK(AppsV1B2, Deployment), true},
		{toGVK(AuthenticationV1B1, TokenReview), true},
		{toGVK(AuthorizationV1B1, LocalSubjectAccessReview), true},
		{toGVK(AuthorizationV1B1, SelfSubjectAccessReview), true},
		{toGVK(AuthorizationV1B1, SelfSubjectRulesReview), true},
		{toGVK(AuthorizationV1B1, SubjectAccessReview), true},
		{toGVK(AutoscalingV2B1, HorizontalPodAutoscaler), true},
		{toGVK(CoordinationV1B1, Lease), true},
		{toGVK(ExtensionsV1B1, DaemonSet), true},
		{toGVK(ExtensionsV1B1, Deployment), true},
		{toGVK(ExtensionsV1B1, Ingress), true},
		{toGVK(ExtensionsV1B1, NetworkPolicy), true},
		{toGVK(ExtensionsV1B1, PodSecurityPolicy), true},
		{toGVK(ExtensionsV1B1, ReplicaSet), true},
		{toGVK(RbacV1A1, ClusterRole), true},
		{toGVK(RbacV1A1, ClusterRoleBinding), true},
		{toGVK(RbacV1A1, Role), true},
		{toGVK(RbacV1A1, RoleBinding), true},
		{toGVK(RbacV1B1, ClusterRole), true},
		{toGVK(RbacV1B1, ClusterRoleBinding), true},
		{toGVK(RbacV1B1, Role), true},
		{toGVK(RbacV1B1, RoleBinding), true},
		{toGVK(SchedulingV1A1, PriorityClass), true},
		{toGVK(SchedulingV1B1, PriorityClass), true},
		{toGVK(StorageV1A1, VolumeAttachment), true},
		{toGVK(StorageV1B1, CSIDriver), true},
		{toGVK(StorageV1B1, CSINode), true},
		{toGVK(StorageV1B1, StorageClass), true},
		{toGVK(StorageV1B1, VolumeAttachment), true},
	}
	for _, tt := range tests {
		t.Run(tt.gvk.String(), func(t *testing.T) {
			if got := DeprecatedAPIVersion(tt.gvk); got != tt.want {
				t.Errorf("DeprecatedAPIVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSuggestedApiVersion(t *testing.T) {
	wantStr := func(gv groupVersion, kind Kind) string {
		return gvkStr(toGVK(gv, kind))
	}

	tests := []struct {
		gvk  GroupVersionKind
		want string
	}{
		// Deprecated ApiVersions return the suggested version string.
		{toGVK(AdmissionregistrationV1B1, MutatingWebhookConfiguration), wantStr(AdmissionregistrationV1, MutatingWebhookConfiguration)},
		{toGVK(AdmissionregistrationV1B1, ValidatingWebhookConfiguration), wantStr(AdmissionregistrationV1, ValidatingWebhookConfiguration)},
		{toGVK(ApiextensionsV1B1, CustomResourceDefinition), wantStr(ApiextensionsV1, CustomResourceDefinition)},
		{toGVK(ApiregistrationV1B1, APIService), wantStr(ApiregistrationV1, APIService)},
		{toGVK(AppsV1B1, Deployment), wantStr(AppsV1, Deployment)},
		{toGVK(AppsV1B2, Deployment), wantStr(AppsV1, Deployment)},
		{toGVK(AuthenticationV1B1, TokenReview), wantStr(AuthenticationV1, TokenReview)},
		{toGVK(AuthorizationV1B1, LocalSubjectAccessReview), wantStr(AuthorizationV1, LocalSubjectAccessReview)},
		{toGVK(AutoscalingV2B1, HorizontalPodAutoscaler), wantStr(AutoscalingV1, HorizontalPodAutoscaler)},
		{toGVK(CoordinationV1B1, Lease), wantStr(CoordinationV1, Lease)},
		{toGVK(ExtensionsV1B1, DaemonSet), wantStr(AppsV1, DaemonSet)},
		{toGVK(ExtensionsV1B1, Deployment), wantStr(AppsV1, Deployment)},
		{toGVK(ExtensionsV1B1, Ingress), wantStr(NetworkingV1B1, Ingress)},
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
		gvk         GroupVersionKind
		wantVersion *cluster.ServerVersion
	}{
		{toGVK(AdmissionregistrationV1B1, MutatingWebhookConfiguration), &cluster.ServerVersion{Major: 1, Minor: 19}},
		{toGVK(ApiextensionsV1B1, CustomResourceDefinition), &cluster.ServerVersion{Major: 1, Minor: 19}},
		{toGVK(AppsV1B1, Deployment), &cluster.ServerVersion{Major: 1, Minor: 16}},
		{toGVK(AppsV1B2, Deployment), &cluster.ServerVersion{Major: 1, Minor: 16}},
		{toGVK(AuthenticationV1B1, TokenReview), &cluster.ServerVersion{Major: 1, Minor: 22}},
		{toGVK(AuthorizationV1B1, LocalSubjectAccessReview), &cluster.ServerVersion{Major: 1, Minor: 22}},
		{toGVK(CoordinationV1B1, Lease), &cluster.ServerVersion{Major: 1, Minor: 22}},
		{toGVK(ExtensionsV1B1, Deployment), &cluster.ServerVersion{Major: 1, Minor: 16}},
		{toGVK(ExtensionsV1B1, Ingress), &cluster.ServerVersion{Major: 1, Minor: 20}},
		{toGVK(RbacV1A1, ClusterRole), &cluster.ServerVersion{Major: 1, Minor: 22}},
		{toGVK(RbacV1B1, ClusterRole), &cluster.ServerVersion{Major: 1, Minor: 22}},
		{toGVK(SchedulingV1A1, PriorityClass), &cluster.ServerVersion{Major: 1, Minor: 17}},
		{toGVK(SchedulingV1B1, PriorityClass), &cluster.ServerVersion{Major: 1, Minor: 17}},
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

func TestRemovedApiVersion(t *testing.T) {
	type args struct {
		gvk     GroupVersionKind
		version cluster.ServerVersion
	}
	tests := []struct {
		name        string
		args        args
		wantRemoved bool
		wantVersion *cluster.ServerVersion
	}{
		{"API exists", args{
			GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
			cluster.ServerVersion{Major: 1, Minor: 16}}, false, nil},
		{"API removed", args{
			GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Deployment"},
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
