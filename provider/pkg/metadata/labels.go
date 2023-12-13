// Copyright 2016-2023, Pulumi Corporation.
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
	"fmt"
	"os"
	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	LabelManagedBy = "app.kubernetes.io/managed-by"
)

// TrySetLabel attempts to set the specified key/value pair as a label on the provided Unstructured
// object, reporting whether the write was successful, and an error if (e.g.) the underlying object
// is mistyped. In particular, TrySetLabel will fail if the underlying object has a Pulumi computed
// value.
func TrySetLabel(obj *unstructured.Unstructured, key, value string) (succeeded bool, err error) {
	// Note: Cannot use obj.GetLabels() here because it doesn't properly handle computed values from preview.
	// During preview, don't set labels if the metadata or label contains a computed value since there's
	// no way to insert data into the computed object.
	metadataRaw, ok := obj.Object["metadata"]
	if isComputedValue(metadataRaw) {
		return false, nil
	}

	var isMap bool
	var metadata map[string]any
	if !ok {
		metadata = map[string]any{}
		obj.Object["metadata"] = metadata
	} else {
		metadata, isMap = metadataRaw.(map[string]any)
		if !isMap {
			return false, fmt.Errorf("expected .metadata to be a map[string]interface{}, got %q",
				reflect.TypeOf(metadata))
		}
	}

	labelsRaw, ok := metadata["labels"]
	if isComputedValue(labelsRaw) {
		return false, nil
	}
	var labels map[string]any
	if !ok {
		labels = make(map[string]any)
	} else {
		labels, isMap = labelsRaw.(map[string]any)
		if !isMap {
			return false, fmt.Errorf("expected .metadata.labels to be a map[string]interface{}, got %q",
				reflect.TypeOf(labels))
		}
	}

	labels[key] = value
	metadata["labels"] = labels
	return true, nil
}

// GetLabel gets the value of the specified label from the given object.
func GetLabel(obj *unstructured.Unstructured, key string) any {
	metadataRaw := obj.Object["metadata"]
	if isComputedValue(metadataRaw) || metadataRaw == nil {
		return metadataRaw
	}
	metadata := metadataRaw.(map[string]any)
	labelsRaw := metadata["labels"]
	if isComputedValue(labelsRaw) || labelsRaw == nil {
		return labelsRaw
	}
	return labelsRaw.(map[string]any)[key]
}

// TrySetManagedByLabel attempts to set the `app.kubernetes.io/managed-by` label to the `pulumi` key
// on the provided Unstructured object, reporting whether the write was successful, and an error if
// (e.g.) the underlying object is mistyped. In particular, TrySetLabel will fail if the underlying
// object has a Pulumi computed value.
func TrySetManagedByLabel(obj *unstructured.Unstructured) (bool, error) {
	managedBy := "pulumi"
	labelVal, exists := os.LookupEnv("PULUMI_KUBERNETES_MANAGED_BY_LABEL")
	if exists {
		managedBy = labelVal
	}
	return TrySetLabel(obj, LabelManagedBy, managedBy)
}

// HasManagedByLabel returns true if the object has the `app.kubernetes.io/managed-by` label set to the value
// of `PULUMI_KUBERNETES_MANAGED_BY_LABEL` EnvVar, `pulumi`, or is a computed value.
func HasManagedByLabel(obj *unstructured.Unstructured) bool {
	val := GetLabel(obj, LabelManagedBy)
	if isComputedValue(val) {
		return true
	}

	// Check if the value is a string/exists.
	valStr, ok := val.(string)
	if !ok {
		return false
	}

	// Compare the obtained value with the value of `PULUMI_KUBERNETES_MANAGED_BY_LABEL` EnvVar if it exists.
	labelVal, exists := os.LookupEnv("PULUMI_KUBERNETES_MANAGED_BY_LABEL")
	if exists {
		return labelVal == valStr
	}

	return valStr == "pulumi"
}
