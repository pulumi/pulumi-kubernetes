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
	"strconv"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// SkipAwaitLogic returns true if the `pulumi.com/skipAwait` annotation is "true", false otherwise.
func SkipAwaitLogic(obj *unstructured.Unstructured) bool {
	return IsAnnotationTrue(obj, AnnotationSkipAwait)
}

// ReplaceUnready returns true if the `pulumi.com/replaceUnready` annotation is "true", false otherwise.
func ReplaceUnready(obj *unstructured.Unstructured) bool {
	return IsAnnotationTrue(obj, AnnotationReplaceUnready)
}

// TimeoutDuration returns the resource timeout duration. There are a number of things it can do here in this order:
// 1. Return the timeout as specified in the customResource options
// 2. Return the timeout as specified in `pulumi.com/timeoutSeconds` annotation,
// 3. Return nil if the annotation is unset/invalid.
func TimeoutDuration(resourceTimeoutSeconds float64, obj *unstructured.Unstructured) *time.Duration {
	if resourceTimeoutSeconds != 0 {
		timeout := time.Duration(resourceTimeoutSeconds) * time.Second
		return &timeout
	} else if s := GetAnnotationValue(obj, AnnotationTimeoutSeconds); s != "" {
		val, err := strconv.Atoi(s)
		if err == nil {
			timeout := time.Duration(val) * time.Second
			return &timeout
		}
	}

	return nil
}

// DeletionPropagation returns the delete propagation policy, Foreground by default.
func DeletionPropagation(obj *unstructured.Unstructured) metav1.DeletionPropagation {
	policy := GetAnnotationValue(obj, AnnotationDeletionPropagation)
	switch strings.ToLower(policy) {
	case "orphan":
		return metav1.DeletePropagationOrphan
	case "background":
		return metav1.DeletePropagationBackground
	default:
		return metav1.DeletePropagationForeground
	}
}
