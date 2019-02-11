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

package annotations

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const AnnotationTrue = "true"
const AnnotationFalse = "false"
const AnnotationPrefix = "pulumi.com/"
const AnnotationInternalPrefix = AnnotationPrefix + "internal/"
const AnnotationInternalAutonamed = AnnotationInternalPrefix + "autonamed"
const AnnotationSkipAwait = AnnotationPrefix + "skipAwait"

func SetAnnotation(obj *unstructured.Unstructured, key, value string) {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[key] = value
	obj.SetAnnotations(annotations)
}

func SetAnnotationTrue(obj *unstructured.Unstructured, key string) {
	SetAnnotation(obj, key, AnnotationTrue)
}

func IsAnnotationTrue(obj *unstructured.Unstructured, key string) bool {
	annotations := obj.GetAnnotations()
	value := annotations[key]
	return value == AnnotationTrue
}
