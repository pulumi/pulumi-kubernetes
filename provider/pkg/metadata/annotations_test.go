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
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/logging"
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

func TestSkipAwait(t *testing.T) {
	tests := []struct {
		name   string
		getter func(context.Context, condition.Source, clientGetter, *logging.DedupLogger, *unstructured.Unstructured) (condition.Satisfier, error)
		obj    *unstructured.Unstructured
		want   condition.Satisfier
	}{
		{
			name: "skipAwait=true takes priority over delete condition",
			obj: &unstructured.Unstructured{Object: map[string]any{
				"metadata": map[string]any{
					"annotations": map[string]any{
						AnnotationSkipAwait: "true",
					},
				},
			}},
			getter: GetDeletedCondition,
			want:   condition.Immediate{},
		},
		{
			name: "skipAwait=false doesn't take priority over delete condition",
			obj: &unstructured.Unstructured{Object: map[string]any{
				"metadata": map[string]any{
					"annotations": map[string]any{
						AnnotationSkipAwait: "false",
					},
				},
			}},
			getter: GetDeletedCondition,
			want:   &condition.Deleted{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition, err := tt.getter(context.Background(), nil, noopClientGetter{}, nil, tt.obj)
			require.NoError(t, err)

			assert.IsType(t, tt.want, condition)
		})
	}
}

type noopClientGetter struct{}

func (noopClientGetter) ResourceClientForObject(*unstructured.Unstructured) (dynamic.ResourceInterface, error) {
	return nil, nil
}
