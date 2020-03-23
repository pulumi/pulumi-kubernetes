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

package metadata

import (
	"strings"

	"github.com/pulumi/pulumi/sdk/go/common/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	AnnotationTrue  = "true"
	AnnotationFalse = "false"

	AnnotationPrefix = "pulumi.com/"

	AnnotationAutonamed         = AnnotationPrefix + "autonamed"
	AnnotationSkipAwait         = AnnotationPrefix + "skipAwait"
	AnnotationTimeoutSeconds    = AnnotationPrefix + "timeoutSeconds"
	AnnotationInitialAPIVersion = AnnotationPrefix + "initialApiVersion"
)

// Annotations for internal Pulumi use only.
var internalAnnotationPrefixes = []string{AnnotationAutonamed}

// IsInternalAnnotation returns true if the specified annotation has the `pulumi.com/` prefix, false otherwise.
func IsInternalAnnotation(key string) bool {
	for _, annotationPrefix := range internalAnnotationPrefixes {
		if strings.HasPrefix(key, annotationPrefix) {
			return true
		}
	}

	return false
}

// SetAnnotation sets the specified key, value annotation on the provided Unstructured object.
// TODO(levi): This won't work for Pulumi-computed values. https://github.com/pulumi/pulumi-kubernetes/issues/826
func SetAnnotation(obj *unstructured.Unstructured, key, value string) {
	// Note: Cannot use obj.GetAnnotations() here because it doesn't properly handle computed values from preview.
	// During preview, don't set annotations if the metadata or annotation contains a computed value since there's
	// no way to insert data into the computed object.
	metadataRaw := obj.Object["metadata"]
	if isComputedValue(metadataRaw) {
		return
	}
	metadata := metadataRaw.(map[string]interface{})
	annotationsRaw, ok := metadata["annotations"]
	if isComputedValue(annotationsRaw) {
		return
	}
	var annotations map[string]interface{}
	if !ok {
		annotations = make(map[string]interface{})
	} else {
		annotations = annotationsRaw.(map[string]interface{})
	}
	annotations[key] = value

	metadata["annotations"] = annotations
}

// SetAnnotationTrue sets the specified annotation key to "true" on the provided Unstructured object.
func SetAnnotationTrue(obj *unstructured.Unstructured, key string) {
	SetAnnotation(obj, key, AnnotationTrue)
}

// IsAnnotationTrue returns true if the specified annotation has the value "true", false otherwise.
func IsAnnotationTrue(obj *unstructured.Unstructured, key string) bool {
	annotations := obj.GetAnnotations()
	value := annotations[key]
	return value == AnnotationTrue
}

// GetAnnotationValue returns the value of the specified annotation on the provided Unstructured object.
func GetAnnotationValue(obj *unstructured.Unstructured, key string) string {
	annotations := obj.GetAnnotations()
	return annotations[key]
}

func isComputedValue(v interface{}) bool {
	_, isComputed := v.(resource.Computed)
	return isComputed
}
