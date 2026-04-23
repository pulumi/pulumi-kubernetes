// Copyright 2016-2024, Pulumi Corporation.
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

package images

import (
	"sort"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// podSpecPaths maps Kubernetes resource kinds to the nested field path
// of their embedded PodSpec. Only kinds that contain a PodSpec are listed.
var podSpecPaths = map[string][]string{
	"Pod":                   {"spec"},
	"Deployment":            {"spec", "template", "spec"},
	"DaemonSet":             {"spec", "template", "spec"},
	"StatefulSet":           {"spec", "template", "spec"},
	"ReplicaSet":            {"spec", "template", "spec"},
	"ReplicationController": {"spec", "template", "spec"},
	"Job":                   {"spec", "template", "spec"},
	"CronJob":               {"spec", "jobTemplate", "spec", "template", "spec"},
	"PodTemplate":           {"template", "spec"},
}

// FromObjects extracts a deduplicated, sorted list of container image
// references from a slice of Kubernetes unstructured objects. Returns
// an empty slice (never nil) if no images are found.
func FromObjects(objs []unstructured.Unstructured) []string {
	seen := make(map[string]struct{})
	for i := range objs {
		extractFromObject(&objs[i], seen)
	}

	result := make([]string, 0, len(seen))
	for img := range seen {
		result = append(result, img)
	}
	sort.Strings(result)
	return result
}

// extractFromObject extracts images from a single unstructured object.
func extractFromObject(obj *unstructured.Unstructured, seen map[string]struct{}) {
	kind := obj.GetKind()
	path, ok := podSpecPaths[kind]
	if !ok {
		return
	}

	podSpec, found, err := unstructured.NestedMap(obj.Object, path...)
	if err != nil || !found {
		return
	}

	extractFromPodSpec(podSpec, seen)
}

// extractFromPodSpec extracts images from a PodSpec map.
func extractFromPodSpec(podSpec map[string]any, seen map[string]struct{}) {
	for _, field := range []string{"containers", "initContainers", "ephemeralContainers"} {
		containers, found, err := unstructured.NestedSlice(podSpec, field)
		if err != nil || !found {
			continue
		}
		for _, c := range containers {
			container, ok := c.(map[string]any)
			if !ok {
				continue
			}
			image, found, err := unstructured.NestedString(container, "image")
			if err != nil || !found || image == "" {
				continue
			}
			seen[image] = struct{}{}
		}
	}

	// ImageVolumeSource (Kubernetes 1.31+): .volumes[*].image.reference
	volumes, found, err := unstructured.NestedSlice(podSpec, "volumes")
	if err != nil || !found {
		return
	}
	for _, v := range volumes {
		volume, ok := v.(map[string]any)
		if !ok {
			continue
		}
		imageVol, found, err := unstructured.NestedMap(volume, "image")
		if err != nil || !found {
			continue
		}
		ref, found, err := unstructured.NestedString(imageVol, "reference")
		if err != nil || !found || ref == "" {
			continue
		}
		seen[ref] = struct{}{}
	}
}
