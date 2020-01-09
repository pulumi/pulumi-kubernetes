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
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/pkg/cluster"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

//
// Reference links:
//
// GVK	/ Deprecated Version / Removed Version
// Upstream Docs Link
// -----------------------------------------------------------
// extensions/v1beta1/DaemonSet / 1.14 / 1.16
// apps/v1beta1/DaemonSet / 1.14 / 1.16
// apps/v1beta2/DaemonSet / 1.14 / 1.16
// extensions/v1beta1/Deployment / 1.14 / 1.16
// apps/v1beta1/Deployment / 1.14 / 1.16
// apps/v1beta2/Deployment / 1.14 / 1.16
// extensions/v1beta1/NetworkPolicy / 1.14 / 1.16
// extensions/v1beta1/PodSecurityPolicy / 1.14 / 1.16
// extensions/v1beta1/ReplicaSet / 1.14 / 1.16
// apps/v1beta1/ReplicaSet / 1.14 / 1.16
// apps/v1beta2/ReplicaSet / 1.14 / 1.16
// https://git.k8s.io/kubernetes/CHANGELOG-1.14.md#deprecations
//
// scheduling/v1alpha1/PriorityClass / 1.14 / 1.17
// scheduling/v1beta1/PriorityClass / 1.14 / 1.17
// https://git.k8s.io/kubernetes/CHANGELOG-1.14.md#deprecations
//
// extensions/v1beta1/Ingress / 1.14 / 1.18
// https://git.k8s.io/kubernetes/CHANGELOG-1.14.md#deprecations
//
// rbac/v1alpha1/* / 1.17 / 1.20
// rbac/v1beta1/* / 1.17 / 1.20
// https://git.k8s.io/kubernetes/CHANGELOG-1.17.md#deprecations-and-removals

func gvkStr(gvk schema.GroupVersionKind) string {
	return gvk.GroupVersion().String() + "/" + gvk.Kind
}

// DeprecatedApiVersion returns true if the given GVK is deprecated in the most recent k8s release.
func DeprecatedApiVersion(gvk schema.GroupVersionKind) bool {
	return SuggestedApiVersion(gvk) != gvkStr(gvk)
}

// RemovedInVersion returns the ServerVersion of k8s that a GVK is removed in. The return value is
// nil if the GVK is not scheduled for removal.
func RemovedInVersion(gvk schema.GroupVersionKind) *cluster.ServerVersion {
	var removedIn cluster.ServerVersion

	switch gvk.GroupVersion() {
	case schema.GroupVersion{Group: "extensions", Version: "v1beta1"},
		schema.GroupVersion{Group: "apps", Version: "v1beta1"},
		schema.GroupVersion{Group: "apps", Version: "v1beta2"}:

		if gvk.Kind == "Ingress" {
			removedIn = cluster.ServerVersion{Major: 1, Minor: 20}
		} else {
			removedIn = cluster.ServerVersion{Major: 1, Minor: 16}
		}
	case schema.GroupVersion{Group: "rbac", Version: "v1beta1"},
		schema.GroupVersion{Group: "rbac", Version: "v1alpha1"}:
		removedIn = cluster.ServerVersion{Major: 1, Minor: 20}
	case schema.GroupVersion{Group: "scheduling", Version: "v1beta1"},
		schema.GroupVersion{Group: "scheduling", Version: "v1alpha1"}:
		removedIn = cluster.ServerVersion{Major: 1, Minor: 17}
	default:
		return nil
	}

	return &removedIn
}

// RemovedApiVersion returns true if the given GVK has been removed in the given k8s version, and the corresponding
// ServerVersion where the GVK was removed.
func RemovedApiVersion(gvk schema.GroupVersionKind, version cluster.ServerVersion) (bool, *cluster.ServerVersion) {
	removedIn := RemovedInVersion(gvk)

	if removedIn == nil {
		return false, nil
	}
	return version.Compare(*removedIn) >= 0, removedIn
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
		case DaemonSet, Deployment, ReplicaSet:
			return "apps/v1/" + gvk.Kind
		case Ingress:
			return "networking/v1beta1/" + gvk.Kind
		case NetworkPolicy:
			return "networking/v1/" + gvk.Kind
		case PodSecurityPolicy:
			return "policy/v1beta1/" + gvk.Kind
		default:
			return gvkStr(gvk)
		}
	case schema.GroupVersion{Group: "rbac", Version: "v1beta1"},
		schema.GroupVersion{Group: "rbac", Version: "v1alpha1"}:
		return "rbac/v1/" + gvk.Kind
	case schema.GroupVersion{Group: "scheduling", Version: "v1beta1"},
		schema.GroupVersion{Group: "scheduling", Version: "v1alpha1"}:
		return "scheduling/v1/" + gvk.Kind
	case schema.GroupVersion{Group: "storage", Version: "v1beta1"}:
		switch Kind(gvk.Kind) {
		case CSINode:
			return "storage/v1/" + gvk.Kind
		default:
			return gvkStr(gvk)
		}
	default:
		return gvkStr(gvk)
	}
}

// upstreamDocsLink returns a link to information about apiVersion deprecations for the given k8s version.
func upstreamDocsLink(version cluster.ServerVersion) string {
	switch version {
	case cluster.ServerVersion{Major: 1, Minor: 16}:
		return "https://git.k8s.io/kubernetes/CHANGELOG-1.16.md#deprecations-and-removals"
	case cluster.ServerVersion{Major: 1, Minor: 17}:
		return "https://git.k8s.io/kubernetes/CHANGELOG-1.17.md#deprecations-and-removals"
	default:
		return ""
	}
}

// RemovedApiError is returned if the provided GVK does not exist in the targeted k8s cluster because the apiVersion
// has been deprecated and removed.
type RemovedApiError struct {
	GVK     schema.GroupVersionKind
	Version *cluster.ServerVersion
}

func (e *RemovedApiError) Error() string {
	if e.Version == nil {
		return fmt.Sprintf("apiVersion %q was removed in a previous version of Kubernetes", gvkStr(e.GVK))
	}

	link := upstreamDocsLink(*e.Version)
	str := fmt.Sprintf("apiVersion %q was removed in Kubernetes %s. Use %q instead.",
		gvkStr(e.GVK), e.Version, SuggestedApiVersion(e.GVK))

	if len(link) > 0 {
		str += fmt.Sprintf("\nSee %s for more information.", link)
	}
	return str
}
