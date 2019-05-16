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

	"github.com/pulumi/pulumi-kubernetes/pkg/await/fixtures"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestFromUnstructured(t *testing.T) {
	pod := fixtures.PodBasic()
	deployment := fixtures.DeploymentBasic()

	type args struct {
		obj *unstructured.Unstructured
	}
	tests := []struct {
		name    string
		args    args
		want    v1.Object
		wantErr bool
	}{
		{"valid-pod", args{obj: pod.Unstructured}, v1.Object(pod.Object), false},
		{"valid-deployment", args{obj: deployment.Unstructured},
			v1.Object(deployment.Object), false},
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

func TestPodFromUnstructured(t *testing.T) {
	type args struct {
		uns *unstructured.Unstructured
	}
	tests := []struct {
		name    string
		args    args
		want    *corev1.Pod
		wantErr bool
	}{
		{"valid", args{uns: fixtures.PodBasic_Uns()}, fixtures.PodBasic(), false},
		{"wrong-type", args{uns: fixtures.DeploymentBasic_Uns()}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PodFromUnstructured(tt.args.uns)
			if (err != nil) != tt.wantErr {
				t.Errorf("PodFromUnstructured() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PodFromUnstructured() = %v, want %v", got, tt.want)
			}
		})
	}
}
}
