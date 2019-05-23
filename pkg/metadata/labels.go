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
	// during preview, our strategy for if metdata or annotations end up being computed values, is to just not
	// apply an annotiation (since there's no way to insert data into the computed object)
	metadataRaw := obj.Object["metadata"]
	if isComputedValue(metadataRaw) {
		return
	}
	metadata := metadataRaw.(map[string]interface{})
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

// SetManagedByLabel sets the `app.kubernetes.io/managed-by` label to `pulumi`.
func SetManagedByLabel(obj *unstructured.Unstructured) {
	SetLabel(obj, "app.kubernetes.io/managed-by", "pulumi")
}
