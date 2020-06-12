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

	"github.com/pulumi/pulumi-kubernetes/provider/v2/pkg/cluster"
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
// https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.14.md#deprecations
//
// scheduling/v1alpha1/PriorityClass / 1.14 / 1.17
// scheduling/v1beta1/PriorityClass / 1.14 / 1.17
// https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.14.md#deprecations
//
// extensions/v1beta1/Ingress / 1.14 / 1.18
// https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.14.md#deprecations
//
// admissionregistration/v1beta1/* / 1.16 / 1.19
// apiextensions/v1beta1/CustomResourceDefinition / 1.16 / 1.19
// https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.16.md#deprecations-and-removals
//
// rbac/v1alpha1/* / 1.17 / 1.22
// rbac/v1beta1/* / 1.17 / 1.22
// https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.17.md#deprecations-and-removals
//
// apiextensions/v1beta1/* / 1.19 / _
// apiregistration/v1beta1/* / 1.19 / _
// authentication/v1beta1/* / 1.19 / 1.22
// authorization/v1beta1/* / 1.19 / 1.22
// autoscaling/v2beta1/* / 1.19 / _
// coordination/v1beta1/* / 1.19 / 1.22
// storage/v1beta1/* / 1.19 / _
// https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.19.md#deprecation-1
//
// TODO: Keep updating this list on every release.

func gvkStr(gvk schema.GroupVersionKind) string {
	return gvk.GroupVersion().String() + "/" + gvk.Kind
}

// DeprecatedAPIVersion returns true if the given GVK is deprecated in the most recent k8s release.
func DeprecatedAPIVersion(gvk schema.GroupVersionKind) bool {
	return SuggestedAPIVersion(gvk) != gvkStr(gvk)
}

// RemovedInVersion returns the ServerVersion of k8s that a GVK is removed in. The return value is
// nil if the GVK is not scheduled for removal.
func RemovedInVersion(gvk schema.GroupVersionKind) *cluster.ServerVersion {
	var removedIn cluster.ServerVersion

	gv, k := groupVersion(gvk.GroupVersion().String()), Kind(gvk.Kind)

	switch gv {
	case AdmissionregistrationV1B1:
		removedIn = cluster.ServerVersion{Major: 1, Minor: 19}
	case ApiextensionsV1B1:
		removedIn = cluster.ServerVersion{Major: 1, Minor: 19}
	case AuthenticationV1B1:
		removedIn = cluster.ServerVersion{Major: 1, Minor: 22}
	case AuthorizationV1B1:
		removedIn = cluster.ServerVersion{Major: 1, Minor: 22}
	case CoordinationV1B1:
		removedIn = cluster.ServerVersion{Major: 1, Minor: 22}
	case ExtensionsV1B1, AppsV1B1, AppsV1B2:
		if k == Ingress {
			removedIn = cluster.ServerVersion{Major: 1, Minor: 20}
		} else {
			removedIn = cluster.ServerVersion{Major: 1, Minor: 16}
		}
	case RbacV1A1, RbacV1B1:
		removedIn = cluster.ServerVersion{Major: 1, Minor: 22}
	case SchedulingV1A1, SchedulingV1B1:
		removedIn = cluster.ServerVersion{Major: 1, Minor: 17}
	default:
		return nil
	}

	return &removedIn
}

// RemovedAPIVersion returns true if the given GVK has been removed in the given k8s version, and the corresponding
// ServerVersion where the GVK was removed.
func RemovedAPIVersion(gvk schema.GroupVersionKind, version cluster.ServerVersion) (bool, *cluster.ServerVersion) {
	removedIn := RemovedInVersion(gvk)

	if removedIn == nil {
		return false, nil
	}
	return version.Compare(*removedIn) >= 0, removedIn
}

// SuggestedAPIVersion returns a string with the suggested apiVersion for a given GVK.
// This is used to provide useful warning messages when a user creates a resource using a deprecated GVK.
func SuggestedAPIVersion(gvk schema.GroupVersionKind) string {
	gv, k := groupVersion(gvk.GroupVersion().String()), Kind(gvk.Kind)

	gvkFmt := `%s/%s`

	switch gv {
	case AdmissionregistrationV1B1:
		return fmt.Sprintf(gvkFmt, AdmissionregistrationV1, k)
	case ApiextensionsV1B1:
		return fmt.Sprintf(gvkFmt, ApiextensionsV1, k)
	case ApiregistrationV1B1:
		return fmt.Sprintf(gvkFmt, ApiregistrationV1, k)
	case AppsV1B1, AppsV1B2:
		return fmt.Sprintf(gvkFmt, AppsV1, k)
	case AuthenticationV1B1:
		return fmt.Sprintf(gvkFmt, AuthenticationV1, k)
	case AuthorizationV1B1:
		return fmt.Sprintf(gvkFmt, AuthorizationV1, k)
	case AutoscalingV2B1:
		return fmt.Sprintf(gvkFmt, AutoscalingV1, k)
	case CoordinationV1B1:
		return fmt.Sprintf(gvkFmt, CoordinationV1, k)
	case ExtensionsV1B1:
		switch k {
		case DaemonSet, Deployment, ReplicaSet:
			return fmt.Sprintf(gvkFmt, AppsV1, k)
		case Ingress:
			return fmt.Sprintf(gvkFmt, NetworkingV1B1, k)
		case NetworkPolicy:
			return fmt.Sprintf(gvkFmt, NetworkingV1, k)
		case PodSecurityPolicy:
			return fmt.Sprintf(gvkFmt, PolicyV1B1, k)
		default:
			return gvkStr(gvk)
		}
	case RbacV1A1, RbacV1B1:
		return fmt.Sprintf(gvkFmt, RbacV1, k)
	case SchedulingV1A1, SchedulingV1B1:
		return fmt.Sprintf(gvkFmt, SchedulingV1, k)
	case StorageV1A1, StorageV1B1, "storage/v1alpha1", "storage/v1beta1": // The storage group was renamed to storage.k8s.io, so check for both.
		return fmt.Sprintf(gvkFmt, StorageV1, k)
	default:
		return gvkStr(gvk)
	}
}

// upstreamDocsLink returns a link to information about apiVersion deprecations for the given k8s version.
func upstreamDocsLink(version cluster.ServerVersion) string {
	switch version {
	case cluster.ServerVersion{Major: 1, Minor: 16}:
		return "https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.16.md#deprecations-and-removals"
	case cluster.ServerVersion{Major: 1, Minor: 17}:
		return "https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.17.md#deprecations-and-removals"
	case cluster.ServerVersion{Major: 1, Minor: 19}:
		return "https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.19.md#deprecation-1"
		// TODO: 1.20
	default:
		return ""
	}
}

// RemovedAPIError is returned if the provided GVK does not exist in the targeted k8s cluster because the apiVersion
// has been deprecated and removed.
type RemovedAPIError struct {
	GVK     schema.GroupVersionKind
	Version *cluster.ServerVersion
}

func (e *RemovedAPIError) Error() string {
	if e.Version == nil {
		return fmt.Sprintf("apiVersion %q was removed in a previous version of Kubernetes", gvkStr(e.GVK))
	}

	link := upstreamDocsLink(*e.Version)
	str := fmt.Sprintf("apiVersion %q was removed in Kubernetes %s. Use %q instead.",
		gvkStr(e.GVK), e.Version, SuggestedAPIVersion(e.GVK))

	if len(link) > 0 {
		str += fmt.Sprintf("\nSee %s for more information.", link)
	}
	return str
}
