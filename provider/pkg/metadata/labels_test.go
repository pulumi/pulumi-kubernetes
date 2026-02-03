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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

func TestSetLabel(t *testing.T) {
	noLabelObj := &unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{},
	}}
	existingLabelObj := &unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{
			"labels": map[string]any{
				"pulumi": "rocks",
			},
		},
	}}
	incorrectMetadataType := &unstructured.Unstructured{Object: map[string]any{
		"metadata": "badtyping",
	}}
	incorrectLabelsType := &unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{"labels": "badtyping"},
	}}
	computedMetadataObj := &unstructured.Unstructured{Object: map[string]any{
		"metadata": resource.Computed{Element: resource.NewObjectProperty(nil)},
	}}
	computedLabelObj := &unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{
			"labels": resource.Computed{Element: resource.NewObjectProperty(nil)},
		},
	}}

	type args struct {
		obj         *unstructured.Unstructured
		shouldError bool
		key         string
		value       string
		expectSet   bool // True if SetLabel is expected to set the label.
		expectKey   string
		expectValue string
	}
	tests := []struct {
		name string
		args args
	}{
		{"set-with-no-label", args{
			obj: noLabelObj, key: "foo", value: "bar", expectSet: true, expectKey: "foo", expectValue: "bar"}},
		{"set-with-existing-labels", args{
			obj: existingLabelObj, key: "foo", value: "bar", expectSet: true, expectKey: "foo", expectValue: "bar"}},
		{"fail-if-metadata-type-incorrect", args{obj: incorrectMetadataType, shouldError: true}},
		{"fail-if-label-type-incorrect", args{obj: incorrectLabelsType, shouldError: true}},

		// Computed fields cannot be set, so SetLabel is a no-op.
		{"set-with-computed-metadata", args{
			obj: computedMetadataObj, key: "foo", value: "bar", expectSet: false}},
		{"set-with-computed-label", args{
			obj: computedLabelObj, key: "foo", value: "bar", expectSet: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := TrySetLabel(tt.args.obj, tt.args.key, tt.args.value)
			assert.Equal(t, tt.args.shouldError, err != nil,
				fmt.Sprintf("Expected error: %t, got error: %t", tt.args.shouldError, err != nil))
			if tt.args.shouldError {
				return
			}
			labels := tt.args.obj.GetLabels()
			value, ok := labels[tt.args.expectKey]
			assert.Equal(t, tt.args.expectSet, ok)
			if ok {
				assert.Equal(t, tt.args.expectValue, value)
			}
		})
	}
}

func TestHasManagedByLabel(t *testing.T) {
	tests := []struct {
		name           string
		obj            *unstructured.Unstructured
		managedByLabel string
		want           bool
	}{
		{
			"Has pulumi manager - happy path",
			&unstructured.Unstructured{Object: map[string]any{
				"metadata": map[string]any{
					"labels": map[string]any{
						"app.kubernetes.io/managed-by": "pulumi",
					},
				},
			}},
			"",
			true,
		},
		{
			"Has pulumi manager - happy path with env matching",
			&unstructured.Unstructured{Object: map[string]any{
				"metadata": map[string]any{
					"labels": map[string]any{
						"app.kubernetes.io/managed-by": "pulumi",
					},
				},
			}},
			"pulumi",
			true,
		},
		{
			"Has non-pulumi manager - happy path with env matching",
			&unstructured.Unstructured{Object: map[string]any{
				"metadata": map[string]any{
					"labels": map[string]any{
						"app.kubernetes.io/managed-by": "fake-mgr",
					},
				},
			}},
			"fake-mgr",
			true,
		},
		{
			"Has non-pulumi manager - with mismatch env",
			&unstructured.Unstructured{Object: map[string]any{
				"metadata": map[string]any{
					"labels": map[string]any{
						"app.kubernetes.io/managed-by": "pulumi-mgr",
					},
				},
			}},
			"fake-mgr",
			false,
		},
		{
			"Has no manager label - with set env",
			&unstructured.Unstructured{Object: map[string]any{
				"metadata": map[string]any{},
			}},
			"fake-mgr",
			false,
		},
		{
			"Has no manager label - without env",
			&unstructured.Unstructured{Object: map[string]any{
				"metadata": map[string]any{},
			}},
			"",
			false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.managedByLabel != "" {
				t.Setenv("PULUMI_KUBERNETES_MANAGED_BY_LABEL", tc.managedByLabel)
			}

			if got := HasManagedByLabel(tc.obj); got != tc.want {
				t.Errorf("HasManagedByLabel() = %v, want %v", got, tc.want)
			}
		})
	}
}
