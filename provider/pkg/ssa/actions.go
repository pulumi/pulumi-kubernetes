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
	"strconv"
	"strings"

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

func UpdateFieldManager(
	ctx context.Context,
	client dynamic.ResourceInterface,
	input *unstructured.Unstructured,
	requiredFields []string,
	fieldManager string,
) error {
	// Create a minimal resource spec with the same identity as the input resource.
	obj := unstructured.Unstructured{}
	obj.SetAPIVersion(input.GetAPIVersion())
	obj.SetKind(input.GetKind())
	obj.SetNamespace(input.GetNamespace())
	obj.SetName(input.GetName())

	liveObj, err := client.Get(ctx, input.GetName(), metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Transfer ownership of any required fields.
	for _, field := range requiredFields {
		setRequiredField(liveObj.Object, obj.Object, field)
		// TODO: handle failures
	}

	yamlObj, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = client.Patch(ctx, input.GetName(), types.ApplyPatchType, yamlObj,
		metav1.PatchOptions{
			FieldManager: fieldManager,
		})

	return err
}

// Input: live, new, field
// 0. Convert field to slice of path fragments
// 1. Look up value in live object
// 2. Ensure value exists in new object
// Return new, bool

// fieldRegex matches valid field specifiers.
// Example: a.b.c
// Example: a
// Example: a1.b2
// Example: a.b[0].c
// TODO: update regex to allow slices
//var fieldRegex = regexp.MustCompile(`^(?:[a-zA-Z][a-zA-Z0-9]*.?)+[^.]$`)

// setRequiredField takes a field describing the element in a map, reads that element from the live map, and then sets
// the corresponding value on the obj map. The function returns true if the operation was successful, false otherwise.
func setRequiredField(live, obj map[string]interface{}, field string) bool {
	//if !fieldRegex.MatchString(field) {
	//	return false
	//}
	var path []string
	dotPath := strings.Split(field, ".")
	for _, p := range dotPath {
		if i := strings.Index(p, "["); i >= 0 {
			path = append(path, p[:i])
			path = append(path, p[i:])
		} else {
			path = append(path, p)
		}
		//path = append(path, strings.Split(p, "[")...)
	}

	// TODO: example: spec.template.spec.containers[0].image

	// Traverse to the specified element in the live map.
	var err error
	var liveCursor interface{} = live
	for _, component := range path {
		// Make sure we can actually traverse to the current element.
		if strings.Contains(component, "[") {
			idxStr := component[strings.LastIndex(component, "[")+1 : len(component)-1]
			var idx int
			idx, err = strconv.Atoi(idxStr)
			if err != nil {
				return false
			}

			// Make sure we can actually slice into the current element.
			switch currObj := liveCursor.(type) {
			case []string:
				if idx > len(currObj)-1 {
					return false
				}

				liveCursor = currObj[idx]
			}

		} else {
			// Make sure we can actually dot into the current element.
			currObj, isMap := liveCursor.(map[string]interface{})
			if !isMap {
				return false
			}

			// Attempt to dot into the current element.
			var exists bool
			liveCursor, exists = currObj[component]
			if !exists {
				return false
			}
		}
	}

	// 1. Traverse to the specified element in the obj map.
	// 2. If the element does not exist, create it.
	// 3. Set the last element of the path to the liveCursor value.
	var objCursor interface{} = obj
	for i, component := range path {
		// Make sure we can actually dot into the current element.
		currObj, isMap := objCursor.(map[string]interface{})
		if !isMap {
			return false
		}

		// Attempt to dot into the current element.
		var exists bool
		objCursor, exists = currObj[component]
		if !exists {
			currObj[component] = map[string]interface{}{}
			objCursor = currObj[component]
			if i == len(path)-1 {
				currObj[component] = liveCursor
			}
		}

	}

	return true
}

func slice[T any](input any) ([]T, bool) {
	if v, ok := input.([]T); ok {
		return v, true
	}
	return nil, false
}
