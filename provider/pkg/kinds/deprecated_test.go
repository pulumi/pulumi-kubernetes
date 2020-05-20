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
		{GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}, false},
		{GroupVersionKind{Group: "apps", Version: "v1beta1", Kind: "Deployment"}, true},
		{GroupVersionKind{Group: "apps", Version: "v1beta2", Kind: "Deployment"}, true},
		{GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "DaemonSet"}, true},
		{GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Deployment"}, true},
		{GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Ingress"}, true},
		{GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "NetworkPolicy"}, true},
		{GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "PodSecurityPolicy"}, true},
		{GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "ReplicaSet"}, true},
		{GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1alpha1", Kind: "ClusterRole"}, true},
		{GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1beta1", Kind: "ClusterRole"}, true},
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
	tests := []struct {
		gvk  GroupVersionKind
		want string
	}{
		// Deprecated ApiVersions return the suggested version string.
		{
			GroupVersionKind{Group: "apps", Version: "v1beta1", Kind: "Deployment"},
			"apps/v1/Deployment",
		},
		{
			GroupVersionKind{Group: "apps", Version: "v1beta2", Kind: "Deployment"},
			"apps/v1/Deployment",
		},
		{
			GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Deployment"},
			"apps/v1/Deployment",
		},
		{
			GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Ingress"},
			"networking.k8s.io/v1beta1/Ingress",
		},
		{
			GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "NetworkPolicy"},
			"networking.k8s.io/v1/NetworkPolicy",
		},
		{
			GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "PodSecurityPolicy"},
			"policy/v1beta1/PodSecurityPolicy",
		},
		{
			GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1alpha1", Kind: "ClusterRole"},
			"rbac.authorization.k8s.io/v1/ClusterRole",
		},
		{
			GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1beta1", Kind: "ClusterRole"},
			"rbac.authorization.k8s.io/v1/ClusterRole",
		},
		{
			GroupVersionKind{Group: "scheduling.k8s.io", Version: "v1beta1", Kind: "PriorityClass"},
			"scheduling.k8s.io/v1/PriorityClass",
		},
		{
			GroupVersionKind{Group: "scheduling.k8s.io", Version: "v1alpha1", Kind: "PriorityClass"},
			"scheduling.k8s.io/v1/PriorityClass",
		},
		// Current ApiVersions return the same version string.
		{
			GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
			"apps/v1/Deployment",
		},
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
	type args struct {
		gvk GroupVersionKind
	}
	tests := []struct {
		name        string
		args        args
		wantVersion *cluster.ServerVersion
	}{
		{"extensions/v1beta1:Deployment", args{
			GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Deployment"},
		}, &cluster.ServerVersion{Major: 1, Minor: 16}},
		{"extensions/v1beta1:Ingress", args{
			GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Ingress"},
		}, &cluster.ServerVersion{Major: 1, Minor: 20}},
		{"rbac.authorization.k8s.io/v1beta1:ClusterRole", args{
			GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1beta1", Kind: "ClusterRole"},
		}, &cluster.ServerVersion{Major: 1, Minor: 20}},
		{"rbac.authorization.k8s.io/v1alpha1:ClusterRole", args{
			GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1alpha1", Kind: "ClusterRole"},
		}, &cluster.ServerVersion{Major: 1, Minor: 20}},
		{"scheduling.k8s.io/v1beta1:PriorityClass", args{
			GroupVersionKind{Group: "scheduling.k8s.io", Version: "v1beta1", Kind: "PriorityClass"},
		}, &cluster.ServerVersion{Major: 1, Minor: 17}},
		{"scheduling.k8s.io/v1alpha1:PriorityClass", args{
			GroupVersionKind{Group: "scheduling.k8s.io", Version: "v1alpha1", Kind: "PriorityClass"},
		}, &cluster.ServerVersion{Major: 1, Minor: 17}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemovedInVersion(tt.args.gvk)
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
