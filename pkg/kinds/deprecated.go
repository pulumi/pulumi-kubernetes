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
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func gvkStr(gvk schema.GroupVersionKind) string {
	return gvk.GroupVersion().String() + "/" + gvk.Kind
}

// DeprecatedApiVersion returns true if the given GVK is deprecated in the most recent k8s release.
func DeprecatedApiVersion(gvk schema.GroupVersionKind) bool {
	return SuggestedApiVersion(gvk) != gvkStr(gvk)
}

// SuggestedApiVersion returns a string with the suggested apiVersion for a given GVK.
// This is used to provide useful warning messages when a user creates a resource using a deprecated GVK.
func SuggestedApiVersion(gvk schema.GroupVersionKind) string {
	switch gvk.GroupVersion() {
	case schema.GroupVersion{Group: "apps", Version: "v1beta1"},
		schema.GroupVersion{Group: "apps", Version: "v1beta2"}:
		return "apps/v1/" + gvk.Kind
	case schema.GroupVersion{Group: "extensions", Version: "v1beta1"}:
		switch Kind(gvk.Kind) {
		case DaemonSet, Deployment, NetworkPolicy, ReplicaSet:
			return "apps/v1/" + gvk.Kind
		case Ingress:
			return "networking/v1beta1/" + gvk.Kind
		case PodSecurityPolicy:
			return "policy/v1beta1/" + gvk.Kind
		default:
			return gvkStr(gvk)
		}
	default:
		return gvkStr(gvk)
	}

}
