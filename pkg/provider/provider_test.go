// Copyright 2016-2018, Pulumi Corporation.
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

package provider

import (
	"testing"

	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var (
	objInputs = map[string]interface{}{
		"foo": "bar",
		"baz": float64(1234),
		"qux": map[string]interface{}{
			"xuq": "oof",
		},
	}
	objLive = map[string]interface{}{
		"oof": "bar",
		"zab": float64(4321),
		"xuq": map[string]interface{}{
			"qux": "foo",
		},
	}
)

func TestParseOldCheckpointObject(t *testing.T) {
	old := resource.NewPropertyMapFromMap(map[string]interface{}{
		"inputs": objInputs,
		"live":   objLive,
	})

	oldInputs, live := parseCheckpointObject(old)
	assert.Equal(t, objInputs, oldInputs.Object)
	assert.Equal(t, objLive, live.Object)
}

func TestParseNewCheckpointObject(t *testing.T) {
	old := resource.NewPropertyMapFromMap(objLive)
	old["__inputs"] = resource.NewObjectProperty(resource.NewPropertyMapFromMap(objInputs))

	oldInputs, live := parseCheckpointObject(old)
	assert.Equal(t, objInputs, oldInputs.Object)
	assert.Equal(t, objLive, live.Object)
}

func TestCheckpointObject(t *testing.T) {
	inputs := &unstructured.Unstructured{Object: objInputs}
	live := &unstructured.Unstructured{Object: objLive}

	obj := checkpointObject(inputs, live)
	assert.NotNil(t, obj)

	oldInputs := obj["__inputs"]
	assert.Equal(t, objInputs, oldInputs.Mappable())

	delete(obj, "__inputs")
	assert.Equal(t, objLive, obj.Mappable())
}

func TestRoundtripCheckpointObject(t *testing.T) {
	old := resource.NewPropertyMapFromMap(objLive)
	old["__inputs"] = resource.NewObjectProperty(resource.NewPropertyMapFromMap(objInputs))

	oldInputs, oldLive := parseCheckpointObject(old)
	assert.Equal(t, objInputs, oldInputs.Object)
	assert.Equal(t, objLive, oldLive.Object)

	obj := checkpointObject(oldInputs, oldLive)
	assert.Equal(t, old, obj)

	newInputs, newLive := parseCheckpointObject(obj)
	assert.Equal(t, oldInputs, newInputs)
	assert.Equal(t, oldLive, newLive)
}

func Test_equalNumbers(t *testing.T) {
	type args struct {
		a interface{}
		b interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"a = b, int64", args{a: int64(1), b: int64(1)}, true},
		{"a = b, float64", args{a: float64(1), b: float64(1)}, true},
		{"a = b, int64, float64", args{a: int64(1), b: float64(1)}, true},
		{"a = b, float64, int64", args{a: float64(1), b: int64(1)}, true},
		{"a != b, int64", args{a: int64(1), b: int64(2)}, false},
		{"a != b, float64", args{a: float64(1), b: float64(2)}, false},
		{"a != b, int64, float64", args{a: int64(1), b: float64(2)}, false},
		{"a != b, float64, int64", args{a: float64(1), b: int64(2)}, false},
		{"unsupported a", args{a: "", b: int64(1)}, false},
		{"unsupported b", args{a: int64(1), b: ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := equalNumbers(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("equalNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}
