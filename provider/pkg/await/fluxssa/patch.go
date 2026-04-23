/*
Copyright 2022 Stefan Prodan
Copyright 2022 The Flux authors
Copyright 2026 Pulumi Corporation.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package fluxssa is a minimal copy of the subset of github.com/fluxcd/pkg/ssa
// that pulumi-kubernetes uses. Vendored to avoid a transitive dependency on
// k8s.io/api/autoscaling/v2beta2, which was removed from k8s.io/api in v0.36.0
// but is still imported by fluxcd/pkg/ssa/normalize as of v0.71.0.
// Source: https://github.com/fluxcd/pkg/blob/ssa/v0.71.0/ssa/patch.go
package fluxssa

import (
	"bytes"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

const managedFieldsPath = "/metadata/managedFields"

// JSONPatch defines a patch as specified by RFC 6902
// https://www.rfc-editor.org/rfc/rfc6902
type JSONPatch struct {
	Operation string                      `json:"op"`
	Path      string                      `json:"path"`
	Value     []metav1.ManagedFieldsEntry `json:"value,omitempty"`
}

// NewPatchReplace returns a JSONPatch for replacing the specified path with the given value.
func NewPatchReplace(path string, value []metav1.ManagedFieldsEntry) JSONPatch {
	return JSONPatch{
		Operation: "replace",
		Path:      path,
		Value:     value,
	}
}

// FieldManager identifies a workflow that's managing fields.
type FieldManager struct {
	Name          string                            `json:"name"`
	ExactMatch    bool                              `json:"exactMatch"`
	OperationType metav1.ManagedFieldsOperationType `json:"operationType"`
}

func matchFieldManager(entry metav1.ManagedFieldsEntry, manager FieldManager) bool {
	if entry.Operation != manager.OperationType || entry.Subresource != "" {
		return false
	}
	if manager.ExactMatch {
		return entry.Manager == manager.Name
	}
	return strings.HasPrefix(entry.Manager, manager.Name)
}

// PatchReplaceFieldsManagers returns a JSONPatch array for replacing the managers with matching
// name and operation type with the specified manager name and an apply operation.
func PatchReplaceFieldsManagers(object *unstructured.Unstructured, managers []FieldManager, name string) ([]JSONPatch, error) {
	objEntries := object.GetManagedFields()

	var prevManagedFields metav1.ManagedFieldsEntry
	empty := metav1.ManagedFieldsEntry{}

	for _, entry := range objEntries {
		if entry.Manager == name && entry.Operation == metav1.ManagedFieldsOperationApply {
			prevManagedFields = entry
		}
	}

	var patches []JSONPatch
	entries := make([]metav1.ManagedFieldsEntry, 0, len(objEntries))
	edited := false

each_entry:
	for _, entry := range objEntries {
		if entry == prevManagedFields {
			continue
		}

		for _, manager := range managers {
			if matchFieldManager(entry, manager) {
				if prevManagedFields == empty {
					entry.Manager = name
					entry.Operation = metav1.ManagedFieldsOperationApply
					prevManagedFields = entry
					edited = true
					continue each_entry
				}

				mergedField, err := mergeManagedFieldsV1(prevManagedFields.FieldsV1, entry.FieldsV1)
				if err != nil {
					return nil, fmt.Errorf("unable to merge managed fields: '%w'", err)
				}
				prevManagedFields.FieldsV1 = mergedField
				edited = true
				continue each_entry
			}
		}
		entries = append(entries, entry)
	}

	if !edited {
		return nil, nil
	}

	entries = append(entries, prevManagedFields)
	return append(patches, NewPatchReplace(managedFieldsPath, entries)), nil
}

func mergeManagedFieldsV1(prevField *metav1.FieldsV1, newField *metav1.FieldsV1) (*metav1.FieldsV1, error) {
	if prevField == nil && newField == nil {
		return nil, nil
	}
	if prevField == nil {
		return newField, nil
	}
	if newField == nil {
		return prevField, nil
	}

	prevSet, err := FieldsToSet(*prevField)
	if err != nil {
		return nil, err
	}

	newSet, err := FieldsToSet(*newField)
	if err != nil {
		return nil, err
	}

	unionSet := prevSet.Union(&newSet)
	mergedField, err := SetToFields(*unionSet)
	if err != nil {
		return nil, fmt.Errorf("unable to convert managed set to field: %s", err)
	}

	return &mergedField, nil
}

// FieldsToSet and SetToFields are copied from
// https://github.com/kubernetes/apiserver/blob/c4c20f4f7d4ca609906621943c748bc16797a5f3/pkg/endpoints/handlers/fieldmanager/internal/fields.go
// since it is an internal module and can't be imported.

// FieldsToSet creates a set of paths from an input trie of fields.
func FieldsToSet(f metav1.FieldsV1) (s fieldpath.Set, err error) {
	err = s.FromJSON(bytes.NewReader(f.Raw))
	return s, err
}

// SetToFields creates a trie of fields from an input set of paths.
func SetToFields(s fieldpath.Set) (f metav1.FieldsV1, err error) {
	f.Raw, err = s.ToJSON()
	return f, err
}
