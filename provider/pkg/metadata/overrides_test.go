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

// nolint: lll
package metadata

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/condition"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

func TestSkipAwaitLogic(t *testing.T) {
	resource := &unstructured.Unstructured{}

	annotatedResourceTrue := &unstructured.Unstructured{}
	annotatedResourceTrue.SetAnnotations(map[string]string{AnnotationSkipAwait: AnnotationTrue})

	annotatedResourceFalse := &unstructured.Unstructured{}
	annotatedResourceFalse.SetAnnotations(map[string]string{AnnotationSkipAwait: AnnotationFalse})

	type args struct {
		obj *unstructured.Unstructured
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Skip annotation unset", args: args{resource}, want: false},
		{name: "Skip annotation set true", args: args{annotatedResourceTrue}, want: true},
		{name: "Skip annotation set false", args: args{annotatedResourceFalse}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SkipAwaitLogic(tt.args.obj); got != tt.want {
				t.Errorf("SkipAwaitLogic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeoutSeconds(t *testing.T) {
	resource := &unstructured.Unstructured{}

	annotatedResource15 := &unstructured.Unstructured{}
	annotatedResource15.SetAnnotations(map[string]string{AnnotationTimeoutSeconds: "15"})

	annotatedResourceZero := &unstructured.Unstructured{}
	annotatedResourceZero.SetAnnotations(map[string]string{AnnotationTimeoutSeconds: "0"})

	annotatedResourceInvalid := &unstructured.Unstructured{}
	annotatedResourceInvalid.SetAnnotations(map[string]string{AnnotationTimeoutSeconds: "foo"})

	ptr := func(t time.Duration) *time.Duration {
		return &t
	}

	type args struct {
		customTimeout float64
		obj           *unstructured.Unstructured
	}
	tests := []struct {
		name string
		args args
		want *time.Duration
	}{
		{"Timeout annotation unset", args{customTimeout: 0, obj: resource}, nil},
		{"Timeout annotation set", args{customTimeout: 0, obj: annotatedResource15}, ptr(15 * time.Second)},
		{"Timeout annotation zero", args{customTimeout: 0, obj: annotatedResourceZero}, ptr(0 * time.Second)},
		{"Timeout annotation invalid", args{customTimeout: 0, obj: annotatedResourceInvalid}, nil},
		{"Timeout from customResource", args{customTimeout: 600, obj: annotatedResource15}, ptr(10 * time.Minute)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeoutDuration(tt.args.customTimeout, tt.args.obj); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeoutDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeletionPropagation(t *testing.T) {
	resource := &unstructured.Unstructured{}

	annotatedResourceInvalid := &unstructured.Unstructured{}
	annotatedResourceInvalid.SetAnnotations(map[string]string{AnnotationTimeoutSeconds: "foo"})

	annotatedResourceOrphan := &unstructured.Unstructured{}
	annotatedResourceOrphan.SetAnnotations(map[string]string{AnnotationDeletionPropagation: "orphan"})

	annotatedResourceUpper := &unstructured.Unstructured{}
	annotatedResourceUpper.SetAnnotations(map[string]string{AnnotationDeletionPropagation: "Orphan"})

	type args struct {
		obj *unstructured.Unstructured
	}
	tests := []struct {
		name string
		args args
		want metav1.DeletionPropagation
	}{
		{"undefined", args{obj: resource}, metav1.DeletePropagationForeground},
		{"invalid", args{obj: annotatedResourceInvalid}, metav1.DeletePropagationForeground},
		{"orphan", args{obj: annotatedResourceOrphan}, metav1.DeletePropagationOrphan},
		{"upper", args{obj: annotatedResourceUpper}, metav1.DeletePropagationOrphan},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeletionPropagation(tt.args.obj); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeoutDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeletedCondition(t *testing.T) {
	tests := []struct {
		name   string
		inputs *unstructured.Unstructured
		obj    *unstructured.Unstructured
		want   condition.Satisfier
	}{
		{
			name: "skipAwait=true doesn't affect generic resources",
			inputs: &unstructured.Unstructured{
				Object: map[string]any{
					"metadata": map[string]any{
						"annotations": map[string]any{
							AnnotationSkipAwait: "true",
						},
					},
				},
			},
			want: &condition.Deleted{},
		},
		{
			name: "skipAwait=true does affect legacy resources",
			inputs: &unstructured.Unstructured{
				Object: map[string]any{
					"apiVersion": "v1",
					"kind":       "Namespace",
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
			condition, err := DeletedCondition(context.Background(), nil, noopClientGetter{}, nil, tt.inputs, obj)
			require.NoError(t, err)

			assert.IsType(t, tt.want, condition)
		})
	}
}

type noopClientGetter struct{}

func (noopClientGetter) ResourceClientForObject(*unstructured.Unstructured) (dynamic.ResourceInterface, error) {
	return nil, nil
}
