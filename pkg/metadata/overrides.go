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

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// SkipAwaitLogic returns true if the `pulumi.com/skipAwait` annotation is "true", false otherwise.
func SkipAwaitLogic(obj *unstructured.Unstructured) bool {
	return IsAnnotationTrue(obj, AnnotationSkipAwait)
}

// TimeoutSeconds returns the int value of the `pulumi.com/timeoutSeconds` annotation, or the defaultSeconds value
// if the annotation is unset/invalid.
func TimeoutSeconds(obj *unstructured.Unstructured, defaultSeconds int) int {
	if s := GetAnnotationValue(obj, AnnotationTimeoutSeconds); s != "" {
		val, err := strconv.Atoi(s)
		if err == nil {
			return val
		}
	}

	return defaultSeconds
}
