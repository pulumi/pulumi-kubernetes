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
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/cluster"
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
// admissionregistration/v1beta1/* / 1.16 / 1.22 (Previously 1.19, see https://github.com/kubernetes/kubernetes/issues/82021#issuecomment-636873001)
// apiextensions/v1beta1/CustomResourceDefinition / 1.16 / 1.22 (Previously 1.19)
// https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.16.md#deprecations-and-removals
//
// rbac/v1alpha1/* / 1.17 / 1.20
// rbac/v1beta1/* / 1.17 / 1.22
// https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.17.md#deprecations-and-removals
//
// apiextensions/v1beta1/* / 1.19 / 1.22
// apiregistration/v1beta1/* / 1.19 / _
// authentication/v1beta1/* / 1.19 / 1.22
// authorization/v1beta1/* / 1.19 / 1.22
// autoscaling/v2beta1/* / 1.19 / _
// coordination/v1beta1/* / 1.19 / 1.22
// storage/v1beta1/* / 1.19 / _
// https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.19.md#deprecation-1
//
// batch/v2alpha1/CronJob / 1.21 / 1.21
// discovery/v1beta1/EndpointSlice / 1.21 / 1.25
// */PodSecurityPolicy / 1.21 / 1.25
// https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.21.md#deprecation-1
//
// storage/v1alpha1/CSIStorageCapacity / 1.24 / 1.24
// https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.24.md#deprecation-1
//
// core/v1/Endpoints / 1.33 / -
// resource/v1beta1/* / 1.33 / 1.36
// https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.33.md#deprecation

var v18 = cluster.ServerVersion{Major: 1, Minor: 8}
var v19 = cluster.ServerVersion{Major: 1, Minor: 9}
var v110 = cluster.ServerVersion{Major: 1, Minor: 10}
var v111 = cluster.ServerVersion{Major: 1, Minor: 11}
var v112 = cluster.ServerVersion{Major: 1, Minor: 12}
var v113 = cluster.ServerVersion{Major: 1, Minor: 13}
var v114 = cluster.ServerVersion{Major: 1, Minor: 14}
var v116 = cluster.ServerVersion{Major: 1, Minor: 16}
var v117 = cluster.ServerVersion{Major: 1, Minor: 17}
var v118 = cluster.ServerVersion{Major: 1, Minor: 18}
var v119 = cluster.ServerVersion{Major: 1, Minor: 19}
var v120 = cluster.ServerVersion{Major: 1, Minor: 20}
var v121 = cluster.ServerVersion{Major: 1, Minor: 21}
var v122 = cluster.ServerVersion{Major: 1, Minor: 22}
var v124 = cluster.ServerVersion{Major: 1, Minor: 24}
var v125 = cluster.ServerVersion{Major: 1, Minor: 25}
var v126 = cluster.ServerVersion{Major: 1, Minor: 26}
var v127 = cluster.ServerVersion{Major: 1, Minor: 27}
var v129 = cluster.ServerVersion{Major: 1, Minor: 29}
var v131 = cluster.ServerVersion{Major: 1, Minor: 31}
var v132 = cluster.ServerVersion{Major: 1, Minor: 32}
var v133 = cluster.ServerVersion{Major: 1, Minor: 33}

func gvkStr(gvk schema.GroupVersionKind) string {
	return gvk.GroupVersion().String() + "/" + gvk.Kind
}

func gvkFromStr(str string) (gvk schema.GroupVersionKind) {
	parts := strings.Split(str, "/")
	if len(parts) == 3 {
		return schema.GroupVersionKind{Group: parts[0], Version: parts[1], Kind: parts[2]}
	}
	return
}

// DeprecatedAPIVersion returns true if the given GVK is deprecated in the given k8s release.
func DeprecatedAPIVersion(gvk schema.GroupVersionKind, version *cluster.ServerVersion) bool {
	suggestedGVKString := SuggestedAPIVersion(gvk)
	if version == nil {
		// If no version is provided, check only against the latest version.
		return suggestedGVKString != gvkStr(gvk)
	}

	suggestedGVK := gvkFromStr(suggestedGVKString)
	suggestedGVKExists := ExistsInVersion(&suggestedGVK, version)

	return suggestedGVKExists && suggestedGVKString != gvkStr(gvk)
}

// AddedInVersion returns the ServerVersion of k8s that a GVK is added in.
func AddedInVersion(gvk *schema.GroupVersionKind) *cluster.ServerVersion {
	gv, k := groupVersion(gvk.GroupVersion().String()), Kind(gvk.Kind)

	switch gv {
	case AdmissionregistrationV1:
		switch k {
		case MutatingWebhookConfiguration, MutatingWebhookConfigurationList, ValidatingWebhookConfiguration, ValidatingWebhookConfigurationList:
			return &v116
		}
	case ApiextensionsV1B1:
		switch k {
		case CustomResourceDefinition, CustomResourceDefinitionList:
			return &v111
		}
	case ApiextensionsV1:
		switch k {
		case CustomResourceDefinition, CustomResourceDefinitionList:
			return &v116
		}
	case ApiregistrationV1, ApiregistrationV1B1:
		switch k {
		case APIService, APIServiceList:
			return &v111
		}
	case AuditregistrationV1A1:
		switch k {
		case AuditSink, AuditSinkList:
			return &v113
		}
	case AutoscalingV2B2:
		switch k {
		case HorizontalPodAutoscaler, HorizontalPodAutoscalerList:
			return &v112
		}
	case BatchV1:
		switch k {
		case CronJob, CronJobList:
			return &v121
		}
	case CoordinationV1B1:
		switch k {
		case Lease, LeaseList:
			return &v112
		}
	case CoordinationV1:
		switch k {
		case Lease, LeaseList:
			return &v114
		}
	case DiscoveryV1B1:
		switch k {
		case EndpointSlice, EndpointSliceList:
			return &v117
		}
	case DiscoveryV1:
		switch k {
		case EndpointSlice, EndpointSliceList:
			return &v121
		}
	case FlowcontrolV1A1:
		switch k {
		case FlowSchema, FlowSchemaList, PriorityLevelConfiguration, PriorityLevelConfigurationList:
			return &v117
		}
	case NetworkingV1B1:
		switch k {
		case Ingress, IngressList:
			return &v114
		case IngressClass, IngressClassList:
			return &v118
		}
	case NodeV1A1, NodeV1B1:
		switch k {
		case RuntimeClass, RuntimeClassList:
			return &v114
		}
	case PolicyV1B1:
		switch k {
		case PodSecurityPolicy, PodSecurityPolicyList:
			return &v110
		}
	case PolicyV1:
		switch k {
		case PodDisruptionBudget, PodDisruptionBudgetList:
			return &v121
		}
	case ResourceV1B1:
		switch k {
		case DeviceClass, DeviceClassList, ResourceClaim, ResourceClaimList, ResourceClaimTemplate, ResourceClaimTemplateList, ResourceSlice, ResourceSliceList:
			return &v132
		}
	case ResourceV1B2:
		switch k {
		case DeviceClass, DeviceClassList, ResourceClaim, ResourceClaimList, ResourceClaimTemplate, ResourceClaimTemplateList, ResourceSlice, ResourceSliceList:
			return &v133
		}
	case SchedulingV1B1:
		switch k {
		case PriorityClass, PriorityClassList:
			return &v111
		}
	case SchedulingV1:
		switch k {
		case PriorityClass, PriorityClassList:
			return &v114
		}
	case StorageV1A1:
		switch k {
		case CSIStorageCapacity, CSIStorageCapacityList:
			return &v121
		case VolumeAttributesClass, VolumeAttributesClassList:
			return &v129
		}
	case StorageV1B1:
		switch k {
		case VolumeAttachment, VolumeAttachmentList:
			return &v110
		case CSIDriver, CSIDriverList, CSINode, CSINodeList:
			return &v114
		case CSIStorageCapacity, CSIStorageCapacityList:
			return &v121
		case VolumeAttributesClass, VolumeAttributesClassList:
			return &v131
		}
	case StorageV1:
		switch k {
		case VolumeAttachment, VolumeAttachmentList:
			return &v113
		case CSINode, CSINodeList:
			return &v117
		case CSIDriver, CSIDriverList:
			return &v118
		case CSIStorageCapacity, CSIStorageCapacityList:
			return &v124
		}
	}

	// We extends this logic back to v1.10, so for all other kinds we return 1.9, meaning that anyone
	// on a 1.9 or earlier cluster will not see deprecation messages.
	return &v19
}

// ExistsInVersion returns true if the given GVK exists in the given k8s version.
func ExistsInVersion(gvk *schema.GroupVersionKind, version *cluster.ServerVersion) bool {
	if gvk == nil || gvk.Empty() {
		return false
	}
	addedIn := AddedInVersion(gvk)

	return version.Compare(*addedIn) >= 0
}

// RemovedInVersion returns the ServerVersion of k8s that a GVK is removed in. The return value is
// nil if the GVK is not scheduled for removal.
func RemovedInVersion(gvk schema.GroupVersionKind) *cluster.ServerVersion {
	gv, k := groupVersion(gvk.GroupVersion().String()), Kind(gvk.Kind)

	switch gv {
	case AdmissionregistrationV1B1:
		return &v122
	case ApiextensionsV1B1:
		return &v122
	case BatchV2A1:
		return &v121
	case BatchV1B1:
		if k == CronJob {
			return &v125
		}
		return nil
	case CoordinationV1B1:
		return &v122
	case DiscoveryV1B1:
		return &v125
	case ExtensionsV1B1, AppsV1B1, AppsV1B2:
		if k == Ingress || k == IngressList {
			return &v120
		}
		return &v116
	case PolicyV1:
		if k == PodSecurityPolicy || k == PodSecurityPolicyList {
			return &v125
		}
		return nil
	case PolicyV1B1:
		return &v125
	case RbacV1A1:
		return &v120
	case RbacV1B1:
		return &v122
	case ResourceV1A2, NetworkingV1A1:
		return &v131
	case SchedulingV1A1, SchedulingV1B1:
		return &v117
	case StorageV1A1:
		if k == CSIStorageCapacity || k == CSIStorageCapacityList {
			return &v127
		}
		return nil
	case AutoscalingV2B2:
		if k == HorizontalPodAutoscaler {
			return &v126
		}
		return nil
	case FlowcontrolV1B1:
		switch k {
		case FlowSchema, PriorityLevelConfiguration:
			return &v126
		}
	case FlowcontrolV1B2:
		switch k {
		case FlowSchema, PriorityLevelConfiguration:
			return &v129
		}
	case FlowcontrolV1B3:
		switch k {
		case FlowSchema, PriorityLevelConfiguration:
			return &v132
		}
	}

	return nil
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
	case AutoscalingV2B1:
		return fmt.Sprintf(gvkFmt, AutoscalingV1, k)
	case BatchV2A1, BatchV1B1:
		return fmt.Sprintf(gvkFmt, BatchV1, k)
	case CoordinationV1B1:
		return fmt.Sprintf(gvkFmt, CoordinationV1, k)
	case CoreV1:
		switch k {
		case Endpoints:
			return fmt.Sprintf(gvkFmt, CoreV1, EndpointSlice)
		case EndpointsList:
			return fmt.Sprintf(gvkFmt, CoreV1, EndpointSliceList)
		default:
			return gvkStr(gvk)
		}
	case DiscoveryV1B1:
		return fmt.Sprintf(gvkFmt, DiscoveryV1, k)
	case ExtensionsV1B1:
		switch k {
		case DaemonSet, DaemonSetList, Deployment, DeploymentList, ReplicaSet, ReplicaSetList:
			return fmt.Sprintf(gvkFmt, AppsV1, k)
		case Ingress, IngressList:
			return fmt.Sprintf(gvkFmt, NetworkingV1B1, k)
		case NetworkPolicy, NetworkPolicyList:
			return fmt.Sprintf(gvkFmt, NetworkingV1, k)
		case PodSecurityPolicy, PodSecurityPolicyList:
			return fmt.Sprintf(gvkFmt, PolicyV1B1, k)
		default:
			return gvkStr(gvk)
		}
	case FlowcontrolV1B3, FlowcontrolV1B2:
		return fmt.Sprintf(gvkFmt, FlowcontrolV1, k)
	case NodeV1B1, NodeV1A1:
		return fmt.Sprintf(gvkFmt, NodeV1, k)
	case PolicyV1B1:
		switch k {
		case PodDisruptionBudget, PodDisruptionBudgetList:
			return fmt.Sprintf(gvkFmt, PolicyV1, k)
		}
	case RbacV1A1, RbacV1B1:
		return fmt.Sprintf(gvkFmt, RbacV1, k)
	case ResourceV1A1, ResourceV1A2, ResourceV1A3, ResourceV1B1:
		return fmt.Sprintf(gvkFmt, ResourceV1B2, k)
	case SchedulingV1A1, SchedulingV1B1:
		return fmt.Sprintf(gvkFmt, SchedulingV1, k)
	case StorageV1A1, StorageV1B1, "storage/v1alpha1", "storage/v1beta1": // The storage group was renamed to storage.k8s.io, so check for both.
		switch k {
		case VolumeAttributesClass, VolumeAttributesClassList:
			return fmt.Sprintf(gvkFmt, StorageV1B1, k)
		}
		return fmt.Sprintf(gvkFmt, StorageV1, k)
	}

	return gvkStr(gvk)
}

// upstreamDocsLink returns a link to information about apiVersion deprecations for the given k8s version.
func upstreamDocsLink(version cluster.ServerVersion) string {
	switch version {
	case v116:
		return "https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.16.md#deprecations-and-removals"
	case v117:
		return "https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.17.md#deprecations-and-removals"
	case v119:
		return "https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.19.md#deprecation-1"
	case v120:
		return "https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.20.md#deprecation"
	case v121:
		return "https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.21.md#deprecation"
	case v127:
		return "https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-1.27.md#deprecation"
	default:
		// If we don't have a specific link for the version, we link to the general changelog deprecation header.
		return fmt.Sprintf("https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-%d.%d.md#deprecation", version.Major, version.Minor)
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
