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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestRelatedResource(t *testing.T) {
	pod := &corev1.Pod{
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
	type args struct {
		owner  ResourceId
		object v1.Object
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Matching pod", args{
			owner: ResourceId{
				Name:       "baz",
				Namespace:  "bar",
				GVK:        schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Job"},
				Generation: 0,
			},
			object: pod,
		}, true},
		{"Different namespace", args{
			owner: ResourceId{
				Name:       "baz",
				Namespace:  "default",
				GVK:        schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Job"},
				Generation: 0,
			},
			object: pod,
		}, false},
		{"Different name", args{
			owner: ResourceId{
				Name:       "different",
				Namespace:  "bar",
				GVK:        schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Job"},
				Generation: 0,
			},
			object: pod,
		}, false},
		{"Different GVK", args{
			owner: ResourceId{
				Name:       "baz",
				Namespace:  "bar",
				GVK:        schema.GroupVersionKind{Group: "core", Version: "v1", Kind: "Pod"},
				Generation: 0,
			},
			object: pod,
		}, false},
		{"Different generation", args{
			owner: ResourceId{
				Name:       "baz",
				Namespace:  "bar",
				GVK:        schema.GroupVersionKind{Group: "batch", Version: "v1", Kind: "Job"},
				Generation: 1,
			},
			object: pod,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := relatedResource(tt.args.owner, tt.args.object); got != tt.want {
				t.Errorf("relatedResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
