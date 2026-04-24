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
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
	kscheme "k8s.io/client-go/kubernetes/scheme"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/cluster"
)

// deprecationInfo stores lifecycle data for a single GVK, populated from
// upstream k8s.io/api type annotations at init time.
type deprecationInfo struct {
	IntroducedMajor, IntroducedMinor int // 0,0 = unknown
	RemovedMajor, RemovedMinor       int // 0,0 = not removed
	Replacement                      *schema.GroupVersionKind
}

// deprecations maps GVK → lifecycle info.
var deprecations map[schema.GroupVersionKind]deprecationInfo

// Lifecycle interfaces implemented by k8s.io/api types via generated
// zz_generated.prerelease-lifecycle.go files.
type apiLifecycleIntroduced interface {
	APILifecycleIntroduced() (major, minor int)
}
type apiLifecycleRemoved interface {
	APILifecycleRemoved() (major, minor int)
}
type apiLifecycleReplacement interface {
	APILifecycleReplacement() schema.GroupVersionKind
}

func init() {
	deprecations = make(map[schema.GroupVersionKind]deprecationInfo, 256)

	// Populate from the standard Kubernetes client-go scheme, which registers
	// all types from k8s.io/api. Each type that implements the lifecycle
	// interfaces provides its introduced/removed/replacement metadata.
	for gvk, typ := range kscheme.Scheme.AllKnownTypes() {
		obj := reflect.New(typ).Interface()

		var info deprecationInfo
		var hasData bool

		if lc, ok := obj.(apiLifecycleIntroduced); ok {
			info.IntroducedMajor, info.IntroducedMinor = lc.APILifecycleIntroduced()
			hasData = true
		}
		if lc, ok := obj.(apiLifecycleRemoved); ok {
			info.RemovedMajor, info.RemovedMinor = lc.APILifecycleRemoved()
			hasData = true
		}
		if lc, ok := obj.(apiLifecycleReplacement); ok {
			r := lc.APILifecycleReplacement()
			info.Replacement = &r
			hasData = true
		}

		if !hasData {
			continue
		}

		deprecations[gvk] = info

		// The Kubernetes scheme uses Group="" for core types, but this
		// codebase uses Group="core". Index under both for compatibility.
		if gvk.Group == "" {
			coreGVK := schema.GroupVersionKind{Group: "core", Version: gvk.Version, Kind: gvk.Kind}
			deprecations[coreGVK] = info
		}
	}

	// Manual entries for types not in k8s.io/client-go/kubernetes/scheme.
	// These are either in separate modules (apiextensions, apiregistration)
	// or were removed from k8s.io/api in older Kubernetes versions.

	// apiextensions.k8s.io (from k8s.io/apiextensions-apiserver, not in client-go scheme)
	addManualEntry("apiextensions.k8s.io", "v1beta1", "CustomResourceDefinition",
		1, 7, 1, 22, &schema.GroupVersionKind{Group: "apiextensions.k8s.io", Version: "v1", Kind: "CustomResourceDefinition"})
	crdListReplacement := &schema.GroupVersionKind{
		Group: "apiextensions.k8s.io", Version: "v1", Kind: "CustomResourceDefinitionList",
	}
	addManualEntry("apiextensions.k8s.io", "v1beta1", "CustomResourceDefinitionList",
		1, 7, 1, 22, crdListReplacement)
	addManualEntry("apiextensions.k8s.io", "v1", "CustomResourceDefinition", 1, 16, 0, 0, nil)
	addManualEntry("apiextensions.k8s.io", "v1", "CustomResourceDefinitionList", 1, 16, 0, 0, nil)

	// apiregistration.k8s.io (from k8s.io/kube-aggregator, not in client-go scheme)
	addManualEntry("apiregistration.k8s.io", "v1beta1", "APIService",
		1, 10, 0, 0, &schema.GroupVersionKind{Group: "apiregistration.k8s.io", Version: "v1", Kind: "APIService"})
	addManualEntry("apiregistration.k8s.io", "v1beta1", "APIServiceList",
		1, 10, 0, 0, &schema.GroupVersionKind{Group: "apiregistration.k8s.io", Version: "v1", Kind: "APIServiceList"})
	addManualEntry("apiregistration.k8s.io", "v1", "APIService", 1, 10, 0, 0, nil)
	addManualEntry("apiregistration.k8s.io", "v1", "APIServiceList", 1, 10, 0, 0, nil)

	// batch/v2alpha1 (removed from k8s.io/api, CronJob moved to batch/v1)
	addManualEntry("batch", "v2alpha1", "CronJob",
		1, 8, 1, 21, &schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "CronJob"})
	addManualEntry("batch", "v2alpha1", "CronJobList",
		1, 8, 1, 21, &schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "CronJobList"})

	// autoscaling/v2beta1 and v2beta2 (removed from k8s.io/api in v0.36.0,
	// HPA moved to autoscaling/v2).
	addManualEntry("autoscaling", "v2beta1", "HorizontalPodAutoscaler",
		1, 8, 1, 25, &schema.GroupVersionKind{Group: "autoscaling", Version: "v2", Kind: "HorizontalPodAutoscaler"})
	addManualEntry("autoscaling", "v2beta1", "HorizontalPodAutoscalerList",
		1, 8, 1, 25, &schema.GroupVersionKind{Group: "autoscaling", Version: "v2", Kind: "HorizontalPodAutoscalerList"})
	addManualEntry("autoscaling", "v2beta2", "HorizontalPodAutoscaler",
		1, 12, 1, 26, &schema.GroupVersionKind{Group: "autoscaling", Version: "v2", Kind: "HorizontalPodAutoscaler"})
	addManualEntry("autoscaling", "v2beta2", "HorizontalPodAutoscalerList",
		1, 12, 1, 26, &schema.GroupVersionKind{Group: "autoscaling", Version: "v2", Kind: "HorizontalPodAutoscalerList"})

	// auditregistration.k8s.io/v1alpha1 (removed from k8s.io/api)
	addManualEntry("auditregistration.k8s.io", "v1alpha1", "AuditSink", 1, 13, 0, 0, nil)
	addManualEntry("auditregistration.k8s.io", "v1alpha1", "AuditSinkList", 1, 13, 0, 0, nil)

	// extensions/v1beta1/PodSecurityPolicy — this type was served by the API
	// server under extensions/ but the Go type only exists in policy/. It needs
	// a manual entry since the scheme doesn't register it under extensions/.
	addManualEntry("extensions", "v1beta1", "PodSecurityPolicy",
		1, 3, 1, 16, &schema.GroupVersionKind{Group: "policy", Version: "v1beta1", Kind: "PodSecurityPolicy"})
	addManualEntry("extensions", "v1beta1", "PodSecurityPolicyList",
		1, 3, 1, 16, &schema.GroupVersionKind{Group: "policy", Version: "v1beta1", Kind: "PodSecurityPolicyList"})

	// policy/v1beta1/PodSecurityPolicy — the Go type was removed from k8s.io/api;
	// only PodDisruptionBudget remains in policy/v1beta1.
	addManualEntry("policy", "v1beta1", "PodSecurityPolicy",
		1, 3, 1, 25, nil)
	addManualEntry("policy", "v1beta1", "PodSecurityPolicyList",
		1, 3, 1, 25, nil)

	// policy/v1/PodSecurityPolicy — never existed as a Go type but was served
	// by the API server; removed in 1.25.
	addManualEntry("policy", "v1", "PodSecurityPolicy", 0, 0, 1, 25, nil)
	addManualEntry("policy", "v1", "PodSecurityPolicyList", 0, 0, 1, 25, nil)

	// rbac.authorization.k8s.io/v1alpha1 (no lifecycle files in k8s.io/api)
	for _, kind := range []string{"ClusterRole", "ClusterRoleBinding", "Role", "RoleBinding"} {
		addManualEntry("rbac.authorization.k8s.io", "v1alpha1", kind,
			1, 6, 1, 20, &schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: kind})
		addManualEntry("rbac.authorization.k8s.io", "v1alpha1", kind+"List",
			1, 6, 1, 20, &schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: kind + "List"})
	}

	// scheduling.k8s.io/v1alpha1 (no lifecycle files in k8s.io/api)
	addManualEntry("scheduling.k8s.io", "v1alpha1", "PriorityClass",
		1, 10, 1, 17, &schema.GroupVersionKind{Group: "scheduling.k8s.io", Version: "v1", Kind: "PriorityClass"})
	addManualEntry("scheduling.k8s.io", "v1alpha1", "PriorityClassList",
		1, 10, 1, 17, &schema.GroupVersionKind{Group: "scheduling.k8s.io", Version: "v1", Kind: "PriorityClassList"})

	// resource.k8s.io/v1beta1 — has lifecycle data but no APILifecycleReplacement.
	// The replacement is v1beta2.
	setReplacement("resource.k8s.io", "v1beta1", "DeviceClass",
		schema.GroupVersionKind{Group: "resource.k8s.io", Version: "v1beta2", Kind: "DeviceClass"})
	setReplacement("resource.k8s.io", "v1beta1", "DeviceClassList",
		schema.GroupVersionKind{Group: "resource.k8s.io", Version: "v1beta2", Kind: "DeviceClassList"})
	// NOTE: DeviceTaintRule only exists in v1alpha3 in k8s.io/api v0.35.2.
	// It has not been promoted to v1beta1 or v1beta2 yet, so there is no
	// replacement to point to. The v1alpha3 introduced/removed data is
	// picked up automatically by the scheme iteration above.
	setReplacement("resource.k8s.io", "v1beta1", "ResourceClaim",
		schema.GroupVersionKind{Group: "resource.k8s.io", Version: "v1beta2", Kind: "ResourceClaim"})
	setReplacement("resource.k8s.io", "v1beta1", "ResourceClaimList",
		schema.GroupVersionKind{Group: "resource.k8s.io", Version: "v1beta2", Kind: "ResourceClaimList"})
	setReplacement("resource.k8s.io", "v1beta1", "ResourceClaimTemplate",
		schema.GroupVersionKind{Group: "resource.k8s.io", Version: "v1beta2", Kind: "ResourceClaimTemplate"})
	setReplacement("resource.k8s.io", "v1beta1", "ResourceClaimTemplateList",
		schema.GroupVersionKind{Group: "resource.k8s.io", Version: "v1beta2", Kind: "ResourceClaimTemplateList"})
	setReplacement("resource.k8s.io", "v1beta1", "ResourceSlice",
		schema.GroupVersionKind{Group: "resource.k8s.io", Version: "v1beta2", Kind: "ResourceSlice"})
	setReplacement("resource.k8s.io", "v1beta1", "ResourceSliceList",
		schema.GroupVersionKind{Group: "resource.k8s.io", Version: "v1beta2", Kind: "ResourceSliceList"})

	// Also index under truncated group names for compatibility with callers
	// that receive groups from OpenAPI definition names (e.g., "io.k8s.api.storage"
	// gets truncated to "storage" by APIVersionComment in additionalComments.go).
	// The scheme uses "storage.k8s.io" but the caller looks up "storage".
	// This runs after setReplacement so that aliases inherit patched Replacement data.
	extras := make(map[schema.GroupVersionKind]deprecationInfo)
	for gvk, info := range deprecations {
		if idx := strings.Index(gvk.Group, "."); idx > 0 {
			shortGVK := schema.GroupVersionKind{Group: gvk.Group[:idx], Version: gvk.Version, Kind: gvk.Kind}
			if _, exists := deprecations[shortGVK]; !exists {
				extras[shortGVK] = info
			}
		}
	}
	for gvk, info := range extras {
		deprecations[gvk] = info
	}
}

// setReplacement adds a replacement GVK to an existing deprecation entry
// (populated by the scheme iteration) without overwriting its other fields.
// Panics if the GVK is not already in the deprecations map, since a missing
// entry indicates a bug in the caller (e.g., referencing a GVK that doesn't
// exist in the scheme).
func setReplacement(group, version, kind string, replacement schema.GroupVersionKind) {
	gvk := schema.GroupVersionKind{Group: group, Version: version, Kind: kind}
	info, ok := deprecations[gvk]
	if !ok {
		panic(fmt.Sprintf("setReplacement: %s not found in deprecations map; use addManualEntry instead", gvk))
	}
	info.Replacement = &replacement
	deprecations[gvk] = info
}

func addManualEntry(
	group, version, kind string,
	introMajor, introMinor, removedMajor, removedMinor int,
	replacement *schema.GroupVersionKind,
) {
	gvk := schema.GroupVersionKind{Group: group, Version: version, Kind: kind}
	deprecations[gvk] = deprecationInfo{
		IntroducedMajor: introMajor,
		IntroducedMinor: introMinor,
		RemovedMajor:    removedMajor,
		RemovedMinor:    removedMinor,
		Replacement:     replacement,
	}
}

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
	if info, ok := deprecations[*gvk]; ok && info.IntroducedMajor > 0 {
		return &cluster.ServerVersion{Major: info.IntroducedMajor, Minor: info.IntroducedMinor}
	}

	// We extend this logic back to v1.10, so for all other kinds we return 1.9, meaning that anyone
	// on a 1.9 or earlier cluster will not see deprecation messages.
	return &cluster.ServerVersion{Major: 1, Minor: 9}
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
	if info, ok := deprecations[gvk]; ok && info.RemovedMajor > 0 {
		return &cluster.ServerVersion{Major: info.RemovedMajor, Minor: info.RemovedMinor}
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
	if info, ok := deprecations[gvk]; ok && info.Replacement != nil {
		r := *info.Replacement
		// Normalize empty group to "core" for consistent string representation.
		if r.Group == "" {
			r.Group = "core"
		}
		return gvkStr(r)
	}
	return gvkStr(gvk)
}

// upstreamDocsLink returns a link to information about apiVersion deprecations for the given k8s version.
func upstreamDocsLink(version cluster.ServerVersion) string {
	return fmt.Sprintf(
		"https://git.k8s.io/kubernetes/CHANGELOG/CHANGELOG-%d.%d.md#deprecation",
		version.Major,
		version.Minor,
	)
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
