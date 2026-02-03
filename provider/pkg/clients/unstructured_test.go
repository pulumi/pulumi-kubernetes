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

package clients

import (
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	validPodObject = &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
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

	validPodUnstructured = &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]any{
				"name": "foo"},
			"spec": map[string]any{
				"containers": []any{
					map[string]any{
						"name":  "foo",
						"image": "nginx"}},
			},
		},
	}

	validDeploymentObject = &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "foo",
							Image: "nginx",
						},
					},
				},
			},
		},
	}

	validDeploymentUnstructured = &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]any{
				"name": "foo"},
			"spec": map[string]any{
				"template": map[string]any{
					"spec": map[string]any{
						"containers": []any{
							map[string]any{
								"name":  "foo",
								"image": "nginx"}},
					},
				},
			},
		},
	}

	// Unregistered GVK
	unregisteredGVK = &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "pulumi/test",
			"kind":       "Foo",
			"metadata": map[string]any{
				"name": "foo"},
			"spec": map[string]any{
				"template": map[string]any{
					"spec": map[string]any{
						"containers": []any{
							map[string]any{
								"name":  "foo",
								"image": "nginx"}},
					},
				},
			},
		},
	}

	crdPreserveUnknownFieldsUnstructured = &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "apiextensions.k8s.io/v1",
			"kind":       "CustomResourceDefinition",
			"metadata": map[string]any{
				"name": "foobars.stable.example.com",
			},
			"spec": map[string]any{
				"group": "stable.example.com",
				"names": map[string]any{
					"kind":   "FooBar",
					"plural": "foobars",
					"shortNames": []string{
						"fb",
					},
					"singular": "foobar",
				},
				"preserveUnknownFields": false,
				"scope":                 "Namespaced",
				"versions": []map[string]any{
					{
						"name": "v1",
						"schema": map[string]any{
							"openAPIV3Schema": map[string]any{
								"properties": map[string]any{
									"spec": map[string]any{
										"properties": map[string]any{
											"foo": map[string]any{
												"type": "string",
											},
										},
										"type": "object",
									},
								},
								"type": "object",
							},
						},
						"served":  true,
						"storage": true,
					},
				},
			},
		},
	}

	crdStatusUnstructured = &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "apiextensions.k8s.io/v1",
			"kind":       "CustomResourceDefinition",
			"metadata": map[string]any{
				"name": "foobars.stable.example.com",
			},
			"spec": map[string]any{
				"group": "stable.example.com",
				"names": map[string]any{
					"kind":   "FooBar",
					"plural": "foobars",
					"shortNames": []string{
						"fb",
					},
					"singular": "foobar",
				},
				"scope": "Namespaced",
				"versions": []map[string]any{
					{
						"name": "v1",
						"schema": map[string]any{
							"openAPIV3Schema": map[string]any{
								"properties": map[string]any{
									"spec": map[string]any{
										"properties": map[string]any{
											"foo": map[string]any{
												"type": "string",
											},
										},
										"type": "object",
									},
								},
								"type": "object",
							},
						},
						"served":  true,
						"storage": true,
					},
				},
			},
			"status": map[string]any{
				"accceptedNames": map[string]any{
					"kind":   "",
					"plural": "",
				},
				"conditions":     []any{},
				"storedVersions": []any{},
			},
		},
	}

	crdUnstructured = &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "apiextensions.k8s.io/v1",
			"kind":       "CustomResourceDefinition",
			"metadata": map[string]any{
				"name": "foobars.stable.example.com",
			},
			"spec": map[string]any{
				"group": "stable.example.com",
				"names": map[string]any{
					"kind":   "FooBar",
					"plural": "foobars",
					"shortNames": []string{
						"fb",
					},
					"singular": "foobar",
				},
				"scope": "Namespaced",
				"versions": []map[string]any{
					{
						"name": "v1",
						"schema": map[string]any{
							"openAPIV3Schema": map[string]any{
								"properties": map[string]any{
									"spec": map[string]any{
										"properties": map[string]any{
											"foo": map[string]any{
												"type": "string",
											},
										},
										"type": "object",
									},
								},
								"type": "object",
							},
						},
						"served":  true,
						"storage": true,
					},
				},
			},
		},
	}

	secretUnstructured = &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]any{
				"name": "foo"},
			"stringData": map[string]any{
				"foo": "bar",
			},
		},
	}

	secretNormalizedUnstructured = &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]any{
				"name": "foo",
			},
			"data": map[string]any{
				"foo": "YmFy",
			},
		},
	}

	secretNewLineUnstructured = &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]any{
				"name": "foo"},
			"data": map[string]any{
				"foo": "dGhpcyBpcyBhIHRlc3Qgc3RyaW5n\n",
			},
		},
	}

	secretNewLineNormalizedUnstructured = &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]any{
				"name": "foo",
			},
			"data": map[string]any{
				"foo": "dGhpcyBpcyBhIHRlc3Qgc3RyaW5n",
			},
		},
	}
)

func TestFromUnstructured(t *testing.T) {
	type args struct {
		obj *unstructured.Unstructured
	}
	tests := []struct {
		name    string
		args    args
		want    metav1.Object
		wantErr bool
	}{
		{"valid Pod", args{obj: validPodUnstructured}, metav1.Object(validPodObject), false},
		{"valid Deployment", args{obj: validDeploymentUnstructured}, metav1.Object(validDeploymentObject), false},
		{"unregistered GVK", args{obj: unregisteredGVK}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromUnstructured(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromUnstructured() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromUnstructured() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	type args struct {
		uns *unstructured.Unstructured
	}
	tests := []struct {
		name    string
		args    args
		want    *unstructured.Unstructured
		wantErr bool
	}{
		{"unregistered GVK", args{uns: unregisteredGVK}, unregisteredGVK, false},
		{"CRD with preserveUnknownFields", args{uns: crdPreserveUnknownFieldsUnstructured}, crdUnstructured, false},
		{"CRD with status", args{uns: crdStatusUnstructured}, crdUnstructured, false},
		{"Secret with stringData input", args{uns: secretUnstructured}, secretNormalizedUnstructured, false},
		{"Secret with data input", args{uns: secretNormalizedUnstructured}, secretNormalizedUnstructured, false},
		{
			"Secret with data containing trailing new line",
			args{uns: secretNewLineUnstructured},
			secretNewLineNormalizedUnstructured,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Normalize(tt.args.uns)
			if (err != nil) != tt.wantErr {
				t.Errorf("Normalize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Normalize() got = %v, want %v", got, tt.want)
			}
		})
	}
}
