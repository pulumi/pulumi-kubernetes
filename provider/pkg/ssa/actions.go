// Copyright 2016-2022, Pulumi Corporation.
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

package ssa

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/yaml"
)

// Relinquish is used to "undo" a patch operation by relinquishing management of all fields in a resource for the
// specified field manager. Any fields that no longer have a manager will return to a default value, so this isn't
// guaranteed to return the resource back to the same state that it was in before the patch.
//
// This operation will fail unless another manager is responsible for a valid resource spec. In other words, if the
// specified field manager is the sole manager of a resource, then this will fail with a validation error for most
// resources.
func Relinquish(
	ctx context.Context,
	client dynamic.ResourceInterface,
	input *unstructured.Unstructured,
	fieldManager string,
) error {
	// Create a minimal resource spec with the same identity as the input resource.
	obj := unstructured.Unstructured{}
	obj.SetAPIVersion(input.GetAPIVersion())
	obj.SetKind(input.GetKind())
	obj.SetNamespace(input.GetNamespace())
	obj.SetName(input.GetName())

	yamlObj, err := yaml.Marshal(obj.Object)
	if err != nil {
		return err
	}

	// Patching with a minimal spec tells the cluster that this field manager will no longer be managing any fields.
	_, err = client.Patch(ctx, input.GetName(), types.ApplyPatchType, yamlObj,
		metav1.PatchOptions{
			FieldManager: fieldManager,
		})

	return err
}
