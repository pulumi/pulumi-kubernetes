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

package await

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestRelatedResource(t *testing.T) {
	p := &corev1.Pod{
		TypeMeta: v1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:       "foo",
			Namespace:  "bar",
			Generation: 0,
			OwnerReferences: []v1.OwnerReference{
				{
					APIVersion: "batch/v1",
					Kind:       "Job",
					Name:       "baz",
					UID:        "14ba58cc-cf83-11e9-8c3a-025000000001",
				},
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "foo",
					Image: "nginx",
				},
			},
		},
	}
	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(p)
	require.NoError(t, err)
	pod := &unstructured.Unstructured{Object: obj}

	tests := []struct {
		name   string
		owner  *unstructured.Unstructured
		object *unstructured.Unstructured
		want   bool
	}{
		{
			name: "Matching pod",
			owner: &unstructured.Unstructured{Object: map[string]any{
				"apiVersion": "batch/v1",
				"kind":       "Job",
				"metadata": map[string]any{
					"name":       "baz",
					"namespace":  "bar",
					"uid":        "14ba58cc-cf83-11e9-8c3a-025000000001",
					"generation": 0,
				},
			}},
			object: pod,
			want:   true,
		},
		{
			name: "Different namespace",
			owner: &unstructured.Unstructured{Object: map[string]any{
				"apiVersion": "batch/v1",
				"kind":       "Job",
				"metadata": map[string]any{
					"name":       "baz",
					"namespace":  "default",
					"generation": 0,
				},
			}},
			object: pod,
			want:   false,
		},
		{
			name: "Different name",
			owner: &unstructured.Unstructured{Object: map[string]any{
				"apiVersion": "batch/v1",
				"kind":       "Job",
				"metadata": map[string]any{
					"name":       "different",
					"namespace":  "bar",
					"generation": 0,
				},
			}},
			object: pod,
			want:   false,
		},
		{
			name: "Different GVK",
			owner: &unstructured.Unstructured{Object: map[string]any{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]any{
					"name":       "baz",
					"namespace":  "bar",
					"generation": 0,
				},
			}},
			object: pod,
			want:   false,
		},
		{
			name: "pod owned by daemonset, UID match",
			object: &unstructured.Unstructured{
				Object: map[string]any{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]any{
						"name":      "uid-match",
						"namespace": "default",
						"labels": map[string]any{
							"pod-template-generation": "12",
						},
						"ownerReferences": []any{map[string]any{
							"apiVersion": "apps/v1",
							"kind":       "DaemonSet",
							"name":       "ds-owner",
							"uid":        "7a83550a-3ccf-4ee2-b5b2-dd3b2bd6061b",
						}},
						"uid": "de21966d-f6a2-4214-a33b-ecf94a361bee",
					},
				},
			},
			owner: &unstructured.Unstructured{
				Object: map[string]any{
					"apiVersion": "apps/v1",
					"kind":       "DaemonSet",
					"metadata": map[string]any{
						"name":      "ds-owner",
						"namespace": "default",
						"ownerReferences": []any{map[string]any{
							"apiVersion": "pulumi.com/v1",
							"kind":       "Foo",
							"name":       "foo",
							"uid":        "a0a92f2e-c3cb-4948-9c81-136c3ea72490",
						}},
						"generation": int64(10),
						"uid":        "7a83550a-3ccf-4ee2-b5b2-dd3b2bd6061b",
					},
				},
			},

			want: true,
		},
		{
			name: "indirect pod/deployment relationship",
			object: &unstructured.Unstructured{
				Object: map[string]any{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]any{
						"name":      "replica-set-owner",
						"namespace": "default",
						"ownerReferences": []any{map[string]any{
							"apiVersion": "apps/v1",
							"kind":       "ReplicaSet",
							"name":       "deployment-owner",
							"uid":        "1ff8ab6a-eb8a-47e0-b198-ac505289714d",
						}},
						"uid": "29045f16-ecc8-4d1a-af3b-06089f005ad3",
					},
				},
			},
			owner: &unstructured.Unstructured{
				Object: map[string]any{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]any{
						"name":      "deployment-owner",
						"namespace": "default",
						"uid":       "81beaa42-2aae-47dc-95b2-23036c408adc",
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := relatedResource(tt.owner, tt.object)
			assert.Equal(t, tt.want, got)
		})
	}
}
