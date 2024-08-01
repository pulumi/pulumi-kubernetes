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
	"context"
	"testing"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/condition"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

func TestSetAnnotation(t *testing.T) {
	noAnnotationObj := &unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{},
	}}
	existingAnnotationObj := &unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{
			"annotations": map[string]any{
				"pulumi": "rocks",
			},
		},
	}}
	computedMetadataObj := &unstructured.Unstructured{Object: map[string]any{
		"metadata": resource.Computed{Element: resource.NewObjectProperty(nil)},
	}}
	computedAnnotationObj := &unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{
			"annotations": resource.Computed{Element: resource.NewObjectProperty(nil)},
		},
	}}

	type args struct {
		obj         *unstructured.Unstructured
		key         string
		value       string
		expectSet   bool // True if SetAnnotation is expected to set the annotation.
		expectKey   string
		expectValue string
	}
	tests := []struct {
		name string
		args args
	}{
		{"set-with-no-annotation", args{
			obj: noAnnotationObj, key: "foo", value: "bar", expectSet: true, expectKey: "foo", expectValue: "bar",
		}},
		{"set-with-existing-annotations", args{
			obj: existingAnnotationObj, key: "foo", value: "bar", expectSet: true, expectKey: "foo", expectValue: "bar",
		}},

		// Computed fields cannot be set, so SetAnnotation is a no-op.
		{"set-with-computed-metadata", args{
			obj: computedMetadataObj, key: "foo", value: "bar", expectSet: false,
		}},
		{"set-with-computed-annotation", args{
			obj: computedAnnotationObj, key: "foo", value: "bar", expectSet: false,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetAnnotation(tt.args.obj, tt.args.key, tt.args.value)
			annotations := tt.args.obj.GetAnnotations()
			value, ok := annotations[tt.args.expectKey]
			assert.Equal(t, tt.args.expectSet, ok)
			if ok {
				assert.Equal(t, tt.args.expectValue, value)
			}
		})
	}
}

func TestGetDeletedCondition(t *testing.T) {
	tests := []struct {
		name   string
		inputs *unstructured.Unstructured
		obj    *unstructured.Unstructured
		want   condition.Satisfier
	}{
		{
			name: "skipAwait=true",
			inputs: &unstructured.Unstructured{
				Object: map[string]any{
					"metadata": map[string]any{
						"annotations": map[string]any{
							AnnotationSkipAwait: "true",
						},
					},
				},
			},
			want: condition.Immediate{},
		},
		{
			name: "skipAwait=false",
			inputs: &unstructured.Unstructured{
				Object: map[string]any{
					"metadata": map[string]any{
						"annotations": map[string]any{
							AnnotationSkipAwait: "false",
						},
					},
				},
			},
			want: &condition.Deleted{},
		},
		{
			name: "skipAwait unset",
			inputs: &unstructured.Unstructured{
				Object: map[string]any{
					"metadata": map[string]any{},
				},
			},
			want: &condition.Deleted{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := tt.obj
			if obj == nil {
				obj = tt.inputs
			}
			condition, err := GetDeletedCondition(context.Background(), nil, noopClientGetter{}, nil, tt.inputs, obj)
			require.NoError(t, err)

			assert.IsType(t, tt.want, condition)
		})
	}
}

type noopClientGetter struct{}

func (noopClientGetter) ResourceClientForObject(*unstructured.Unstructured) (dynamic.ResourceInterface, error) {
	return nil, nil
}
