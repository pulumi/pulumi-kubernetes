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
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/pulumi/pulumi/sdk/v3/go/common/util/deepcopy"

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

REQUIRED:
	// Transfer ownership of any required fields.
	for _, field := range requiredFields {
		var ok bool
		// TODO: wrong Object?
		obj.Object, ok = setRequiredField(liveObj.Object, obj.Object, field)
		if !ok {
			// TODO: better handle failures
			return fmt.Errorf("failed to update field manager")
		}
	}

	yamlObj, err := yaml.Marshal(obj.Object)
	if err != nil {
		return err
	}

	_, err = client.Patch(ctx, input.GetName(), types.ApplyPatchType, yamlObj,
		metav1.PatchOptions{
			FieldManager: fieldManager,
		})
	if k8serrors.IsInvalid(err) {
		var statusErr *k8serrors.StatusError
		if errors.As(err, &statusErr) {
			for _, cause := range statusErr.Status().Details.Causes {
				if cause.Type != metav1.CauseTypeFieldValueRequired {
					continue
				}
				requiredFields = append(requiredFields, cause.Field)
			}
		}
		goto REQUIRED
	}

	return err
}

// fieldRegex matches valid field specifiers.
// Example: a.b.c
// Example: a
// Example: a1.b2
// Example: a[0]
// Example: a.b[0]
// Example: a.b[0].c
var fieldRegex = regexp.MustCompile(`^(?:\w+(?:\[\d+])?\.?)+$`)

// setRequiredField takes a field describing the element in a map, reads that element from the live map, and then sets
// the corresponding value on the obj map. The function returns true if the operation was successful, false otherwise.
func setRequiredField(live, obj map[string]interface{}, field string) (map[string]interface{}, bool) {
	if !fieldRegex.MatchString(field) {
		return nil, false
	}

	type pathToken struct {
		IsSlice bool
		Key     string
		Index   int
	}
	var tokens []pathToken

	dotPath := strings.Split(field, ".")
	for _, p := range dotPath {
		if i := strings.Index(p, "["); i >= 0 {
			tokens = append(tokens, pathToken{IsSlice: false, Key: p[:i]})
			idxStr := p[i+1 : len(p)-1]
			if idx, err := strconv.Atoi(idxStr); err != nil {
				return nil, false
			} else {
				tokens = append(tokens, pathToken{IsSlice: true, Index: idx})
			}
		} else {
			tokens = append(tokens, pathToken{IsSlice: false, Key: p})
		}
	}

	// TODO: example: spec.template.spec.containers[0].image

	resultObj := deepcopy.Copy(obj).(map[string]interface{})

	// Traverse to the specified element in the live map.
	var liveCursor interface{} = live
	var objCursor interface{} = resultObj
	for i, token := range tokens {
		if token.IsSlice {
			liveSlice, ok := liveCursor.([]interface{})
			if !ok || token.Index > len(liveSlice)-1 {
				return nil, false
			}

			objSlice, ok := objCursor.([]interface{})
			if !ok {
				return nil, false
			}
			//if i == len(tokens)-1 {
			if token.Index > len(liveSlice)-1 {
				return nil, false
			}
			objSlice[token.Index] = liveSlice[token.Index]
			break
			//}

			//if v := objSlice[token.Index]; v == nil {
			//	if tokens[i+1].IsSlice {
			//		objSlice[token.Index] = make([]interface{}, len(liveSlice[token.Index].([]interface{})))
			//	} else {
			//		objSlice[token.Index] = map[string]interface{}{}
			//	}
			//}
			//liveCursor = liveSlice[token.Index]
			//objCursor = objSlice[token.Index]
		} else {
			liveMap, ok := liveCursor.(map[string]interface{})
			if !ok {
				return nil, false
			}

			objMap, ok := objCursor.(map[string]interface{})
			if !ok {
				return nil, false
			}
			if i == len(tokens)-1 {
				if _, exists := liveMap[token.Key]; !exists {
					return nil, false
				}
				objMap[token.Key] = liveMap[token.Key]
				break
			}

			if _, exists := objMap[token.Key]; !exists {
				if tokens[i+1].IsSlice {
					objMap[token.Key] = make([]interface{}, len(liveMap[token.Key].([]interface{})))
				} else {
					objMap[token.Key] = map[string]interface{}{}
				}
			}
			liveCursor = liveMap[token.Key]
			objCursor = objMap[token.Key]
		}
	}

	return resultObj, true
}
