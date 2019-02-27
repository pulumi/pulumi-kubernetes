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

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	AnnotationTrue  = "true"
	AnnotationFalse = "false"

	AnnotationPrefix = "pulumi.com/"

	AnnotationAutonamed      = AnnotationPrefix + "autonamed"
	AnnotationSkipAwait      = AnnotationPrefix + "skipAwait"
	AnnotationTimeoutSeconds = AnnotationPrefix + "timeoutSeconds"
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
func SetAnnotation(obj *unstructured.Unstructured, key, value string) {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[key] = value
	obj.SetAnnotations(annotations)
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
