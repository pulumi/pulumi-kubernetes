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
	"strings"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
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
		initialAPIVersionKey: "",
		fieldManagerKey:      "",
		"oof":                "bar",
		"zab":                float64(4321),
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

	obj := checkpointObject(inputs, live, nil, "", "")
	assert.NotNil(t, obj)

	oldInputs := obj["__inputs"]
	assert.Equal(t, objInputs, oldInputs.Mappable())

	delete(obj, "__inputs")
	assert.Equal(t, objLive, obj.Mappable())
}

// #2300 - Read() on top-level k8s objects of kind "secret" led to unencrypted __input
func TestCheckpointSecretObject(t *testing.T) {
	objInputSecret := map[string]interface{}{
		"kind": "Secret",
		"data": map[string]interface{}{
			"password": "verysecret",
		},
	}
	objSecretLive := map[string]interface{}{
		initialAPIVersionKey: "",
		fieldManagerKey:      "",
		"kind":               "Secret",
		"data": map[string]interface{}{
			"password": "verysecret",
		},
	}

	// Questionable but correct pinning test as of the time of writing
	assert.False(t, resource.NewPropertyMapFromMap(objInputSecret).ContainsSecrets())
	assert.False(t, resource.NewPropertyMapFromMap(objSecretLive).ContainsSecrets())

	inputs := &unstructured.Unstructured{Object: objInputSecret}
	live := &unstructured.Unstructured{Object: objSecretLive}

	obj := checkpointObject(inputs, live, nil, "", "")
	assert.NotNil(t, obj)

	oldInputs := obj["__inputs"]
	assert.True(t, oldInputs.IsObject())
	oldInputsVal := oldInputs.ObjectValue()
	assert.True(t, oldInputsVal["data"].ContainsSecrets())
}

func TestRoundtripCheckpointObject(t *testing.T) {
	old := resource.NewPropertyMapFromMap(objLive)
	old["__inputs"] = resource.NewObjectProperty(resource.NewPropertyMapFromMap(objInputs))

	oldInputs, oldLive := parseCheckpointObject(old)
	assert.Equal(t, objInputs, oldInputs.Object)
	assert.Equal(t, objLive, oldLive.Object)

	obj := checkpointObject(oldInputs, oldLive, nil, "", "")
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

func TestRemoveNonInputtyFields(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		wanted   string
	}{
		{
			name: "Simple Manifest with no status or generated metadata fields",
			resource: `apiVersion: v1
kind: Pod
metadata:
    name: nginx
spec:
    containers:
    - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80`,
			wanted: `apiVersion: v1
kind: Pod
metadata:
    name: nginx
spec:
    containers:
    - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80`,
		},
		{
			name: "Manifest with generated metadata fields are removed",
			resource: `apiVersion: v1
kind: Pod
metadata:
    name: nginx
    generation: 2
    resourceVersion: "1234"
    creationTimestamp: "2023-05-11T02:33:24Z"
    managedFields:
        - manager: kubectl
          operation: Apply
          apiVersion: v1
          time: "2023-05-11T02:33:24Z"
          fieldsType: FieldsV1
    uid: 1234-5678-9012-3456
    annotations:
        kubectl.kubernetes.io/last-applied-configuration: |
            {"apiVersion":"v1","kind":"Pod","metadata":{"annotations":{},"name":"nginx","namespace":"default"},"spec":{"containers":[{"image":"nginx:1.14.2","name":"nginx","ports":[{"containerPort":80}]}]}}
spec:
    containers:
    - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80`,
			wanted: `apiVersion: v1
kind: Pod
metadata:
    name: nginx
    annotations: {}
spec:
    containers:
    - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80`,
		},
		{
			name: "Manifest with .status field is removed",
			resource: `apiVersion: v1
kind: Pod
metadata:
    name: nginx
    generation: 2
spec:
    containers:
    - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80
status:
    conditions:
    - lastProbeTime: null
      lastTransitionTime: "2023-05-11T02:33:24Z"
      status: "True"
      type: Initialized
    hostIP: 172.18.0.2
    phase: Running
    podIP: 10.244.0.7
    podIPs:
    - ip: 10.244.0.7
    qosClass: BestEffort
    startTime: "2023-05-11T02:33:24Z"`,
			wanted: `apiVersion: v1
kind: Pod
metadata:
    name: nginx
spec:
    containers:
    - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Convert yaml manifest to unstructured.Unstructured.
			obj := yamlToUnstructured(t, tt.resource)
			wanted := yamlToUnstructured(t, tt.wanted)

			// Remove non-inputty fields.
			removeNonInputtyFields(obj)

			// Check that the object is equal to the wanted object.
			if !equality.Semantic.DeepEqual(obj, wanted) {
				t.Errorf("removeNonInputtyFields() = %v, want %v", obj, wanted)
			}
		})
	}
}

func yamlToUnstructured(t *testing.T, raw string) *unstructured.Unstructured {
	// Convert YAML to untructured.Unstructured
	decoder := yamlutil.NewYAMLOrJSONDecoder(strings.NewReader(raw), 4096)
	obj := &unstructured.Unstructured{}
	err := decoder.Decode(obj)
	assert.NoError(t, err)
	return obj

}
