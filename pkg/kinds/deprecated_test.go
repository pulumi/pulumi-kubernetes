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
	"testing"

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
	}
	for _, tt := range tests {
		t.Run(tt.gvk.String(), func(t *testing.T) {
			if got := DeprecatedApiVersion(tt.gvk); got != tt.want {
				t.Errorf("DeprecatedApiVersion() = %v, want %v", got, tt.want)
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
			"networking/v1beta1/Ingress",
		},
		{
			GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "PodSecurityPolicy"},
			"policy/v1beta1/PodSecurityPolicy",
		},
		// Current ApiVersions return the same version string.
		{
			GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
			"apps/v1/Deployment",
		},
	}
	for _, tt := range tests {
		t.Run(tt.gvk.String(), func(t *testing.T) {
			if got := SuggestedApiVersion(tt.gvk); got != tt.want {
				t.Errorf("SuggestedApiVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
