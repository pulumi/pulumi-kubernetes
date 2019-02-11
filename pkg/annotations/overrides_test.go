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

package annotations

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
