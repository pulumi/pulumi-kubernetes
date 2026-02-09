// Copyright 2021, Pulumi Corporation.  All rights reserved.

package await

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/logging"
)

func mockAwaitConfig(outputs *unstructured.Unstructured) awaitConfig {
	return awaitConfig{
		ctx:            context.Background(),
		currentOutputs: outputs,
		logger:         logging.NewLogger(context.Background(), nil, ""),
	}
}

func mockUpdateConfig(outputs *unstructured.Unstructured, previous *unstructured.Unstructured) awaitConfig {
	c := mockAwaitConfig(outputs)
	c.lastOutputs = previous
	return c
}

func decodeUnstructured(text string) (*unstructured.Unstructured, error) {
	obj, _, err := unstructured.UnstructuredJSONScheme.Decode([]byte(text), nil, nil)
	if err != nil {
		return nil, err
	}
	unst, isUnstructured := obj.(*unstructured.Unstructured)
	if !isUnstructured {
		return nil, fmt.Errorf("could not decode object as *unstructured.Unstructured: %v", unst)
	}
	return unst, nil
}

// mustDecodeUnstructured will panic if the input doesn't deserialize properly.
func mustDecodeUnstructured(s string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(s)
	if err != nil {
		panic(err)
	}
	return obj
}

func TestIsOwnedBy(t *testing.T) {
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
		{
			name:   "nil owner",
			object: pod,
			owner:  nil,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isOwnedBy(tt.object, tt.owner)
			assert.Equal(t, tt.want, got)
		})
	}
}
