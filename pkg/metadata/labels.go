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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// SetLabel sets the specified key/value pair as a label on the provided Unstructured object.
func SetLabel(obj *unstructured.Unstructured, key, value string) {
	// Note: Cannot use obj.GetLabels() here because it doesn't properly handle computed values from preview.
	// During preview, don't set labels if the metadata or label contains a computed value since there's
	// no way to insert data into the computed object.
	metadataRaw, ok := obj.Object["metadata"]
	if isComputedValue(metadataRaw) {
		return
	}
	var metadata map[string]interface{}
	if !ok {
		metadata = map[string]interface{}{}
		obj.Object["metadata"] = metadata
	} else {
		metadata = metadataRaw.(map[string]interface{})
	}
	labelsRaw, ok := metadata["labels"]
	if isComputedValue(labelsRaw) {
		return
	}
	var labels map[string]interface{}
	if !ok {
		labels = make(map[string]interface{})
	} else {
		labels = labelsRaw.(map[string]interface{})
	}
	labels[key] = value

	metadata["labels"] = labels
}

// GetLabel gets the value of the specified label from the given object.
func GetLabel(obj *unstructured.Unstructured, key string) interface{} {
	metadataRaw := obj.Object["metadata"]
	if isComputedValue(metadataRaw) || metadataRaw == nil {
		return metadataRaw
	}
	metadata := metadataRaw.(map[string]interface{})
	labelsRaw := metadata["labels"]
	if isComputedValue(labelsRaw) || labelsRaw == nil {
		return labelsRaw
	}
	return labelsRaw.(map[string]interface{})[key]
}

// SetManagedByLabel sets the `app.kubernetes.io/managed-by` label to `pulumi`.
func SetManagedByLabel(obj *unstructured.Unstructured) {
	SetLabel(obj, "app.kubernetes.io/managed-by", "pulumi")
}

// HasManagedByLabel returns true if the object has the `app.kubernetes.io/managed-by` label set to `pulumi`.
func HasManagedByLabel(obj *unstructured.Unstructured) bool {
	val := GetLabel(obj, "app.kubernetes.io/managed-by")
	if isComputedValue(val) {
		return true
	}
	str, ok := val.(string)
	return ok && str == "pulumi"
}
