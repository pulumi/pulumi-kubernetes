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
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]interface{}{
				"name": "foo"},
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{
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
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "foo"},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "foo",
								"image": "nginx"}},
					},
				},
			},
		},
	}

	// Unregistered GVK
	unregisteredGVK = &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "pulumi/test",
			"kind":       "Foo",
			"metadata": map[string]interface{}{
				"name": "foo"},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "foo",
								"image": "nginx"}},
					},
				},
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
