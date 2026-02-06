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
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/kinds"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/txtar"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

var (
	objInputs = map[string]any{
		"foo": "bar",
		"baz": float64(1234),
		"qux": map[string]any{
			"xuq": "oof",
		},
		"x_kubernetes_preserve_unknown_fields": true,
	}
	objLive = map[string]any{
		initialAPIVersionKey: "",
		fieldManagerKey:      "",
		"oof":                "bar",
		"zab":                float64(4321),
		"xuq": map[string]any{
			"qux": "foo",
		},
		"x_kubernetes_preserve_unknown_fields": true,
	}

	objInputsWithDash = map[string]any{
		"foo": "bar",
		"baz": float64(1234),
		"qux": map[string]any{
			"xuq": "oof",
		},
		"x-kubernetes-preserve-unknown-fields": true,
	}
	objLiveWithDash = map[string]any{
		initialAPIVersionKey: "",
		fieldManagerKey:      "",
		"oof":                "bar",
		"zab":                float64(4321),
		"xuq": map[string]any{
			"qux": "foo",
		},
		"x-kubernetes-preserve-unknown-fields": true,
	}
)

func TestParseOldCheckpointObject(t *testing.T) {
	old := resource.NewPropertyMapFromMap(map[string]any{
		"inputs": objInputs,
		"live":   objLive,
	})

	oldInputs, live := parseCheckpointObject(old)
	assert.Equal(t, objInputsWithDash, oldInputs.Object)
	assert.Equal(t, objLiveWithDash, live.Object)
}

func TestParseNewCheckpointObject(t *testing.T) {
	old := resource.NewPropertyMapFromMap(objLive)
	old["__inputs"] = resource.NewObjectProperty(resource.NewPropertyMapFromMap(objInputs))

	oldInputs, live := parseCheckpointObject(old)
	assert.Equal(t, objInputsWithDash, oldInputs.Object)
	assert.Equal(t, objLiveWithDash, live.Object)
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
	objInputSecret := map[string]any{
		"kind": "Secret",
		"data": map[string]any{
			"password": "verysecret",
		},
	}
	objSecretLive := map[string]any{
		initialAPIVersionKey: "",
		fieldManagerKey:      "",
		"kind":               "Secret",
		"data": map[string]any{
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

// Ensure that well-known x-kubernetes-* fields are normazlied to x_kubernetes_*
// in the checkpoint object.
func TestCheckpointXKubernetesFields(t *testing.T) {
	objInputWithDash := map[string]any{
		"kind": "fakekind",
		"spec": map[string]any{
			"x-kubernetes-preserve-unknown-fields": "true",
		},
	}
	objLiveWithDash := map[string]any{
		initialAPIVersionKey: "",
		fieldManagerKey:      "",
		"kind":               "fakekind",
		"spec": map[string]any{
			"x-kubernetes-preserve-unknown-fields": "true",
		},
	}

	inputs := &unstructured.Unstructured{Object: objInputWithDash}
	live := &unstructured.Unstructured{Object: objLiveWithDash}

	obj := checkpointObject(inputs, live, nil, "", "")
	assert.NotNil(t, obj)

	// Ensure we do not have the original x-kubernetes-* fields in the checkpoint objects.
	assert.NotContains(t, obj.Mappable()["spec"], "x-kubernetes-preserve-unknown-fields")
	assert.NotContains(t, obj["__inputs"].Mappable(), "x-kubernetes-preserve-unknown-fields")

	// Ensure we have the normalized x_kubernetes_* fields in the checkpoint objects.
	assert.Contains(t, obj.Mappable()["spec"], "x_kubernetes_preserve_unknown_fields")
	assert.Contains(t, obj["__inputs"].Mappable().(map[string]any)["spec"], "x_kubernetes_preserve_unknown_fields")
}

func TestRoundtripCheckpointObject(t *testing.T) {
	old := resource.NewPropertyMapFromMap(objLive)
	old["__inputs"] = resource.NewObjectProperty(resource.NewPropertyMapFromMap(objInputs))

	oldInputs, oldLive := parseCheckpointObject(old)
	assert.Equal(t, objInputsWithDash, oldInputs.Object)
	assert.Equal(t, objLiveWithDash, oldLive.Object)

	obj := checkpointObject(oldInputs, oldLive, nil, "", "")
	assert.Equal(t, old, obj)

	newInputs, newLive := parseCheckpointObject(obj)
	assert.Equal(t, oldInputs, newInputs)
	assert.Equal(t, oldLive, newLive)
}

func Test_equalNumbers(t *testing.T) {
	type args struct {
		a any
		b any
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

func Test_isListURN(t *testing.T) {
	type args struct {
		urn resource.URN
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "list URN",
			args: args{
				urn: resource.NewURN("test", "test", "", "kubernetes:apps/v1:DeploymentList", "test"),
			},
			want: true,
		},
		{
			name: "regular URN",
			args: args{
				urn: resource.NewURN("test", "test", "", "kubernetes:apps/v1:Deployment", "test"),
			},
			want: false,
		},
		{
			name: "CustomResource with List suffix",
			args: args{
				urn: resource.NewURN("test", "test", "", "kubernetes:example/v1alpha1:ExampleList", "test"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, kinds.IsListURN(tt.args.urn), "isListURN(%v)", tt.args.urn)
		})
	}
}

// TestDiffConfig loads txtar files under testdata/diffconfig and uses them to
// craft a DiffConfig request to a provider.
//
// To add a new test case, create a new file with `olds` and `news` files in
// YAML format. `wantDiffs` and/or `wantReplaces` should be a list of key
// names.
func TestDiffConfig(t *testing.T) {
	dir := filepath.Join("testdata", "diffconfig")
	tests, err := os.ReadDir(dir)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.Name(), func(t *testing.T) {
			archive, err := txtar.ParseFile(filepath.Join(dir, tt.Name()))
			require.NoError(t, err)

			var olds, news resource.PropertyMap
			var wantDiffs, wantReplaces []string
			var wantErr string

			for _, f := range archive.Files {
				var parsed map[string]any
				err := yaml.Unmarshal(f.Data, &parsed)
				switch f.Name {
				case "olds":
					require.NoError(t, err, f.Name)
					olds = resource.NewPropertyMapFromMap(parsed)
				case "news":
					require.NoError(t, err, f.Name)
					news = resource.NewPropertyMapFromMap(parsed)
				case "wantDiffs":
					err = yaml.Unmarshal(f.Data, &wantDiffs)
					require.NoError(t, err, f.Name)
				case "wantReplaces":
					err = yaml.Unmarshal(f.Data, &wantReplaces)
					require.NoError(t, err, f.Name)
				case "wantErr":
					wantErr = strings.TrimSpace(string(f.Data))
				default:
					t.Fatal("unrecognized filename", f.Name)
				}
			}

			// Construct a protobuf request from our inputs.
			opts := plugin.MarshalOptions{KeepUnknowns: true, KeepSecrets: true}
			oldspb, err := plugin.MarshalProperties(olds, opts)
			require.NoError(t, err)
			newspb, err := plugin.MarshalProperties(news, opts)
			require.NoError(t, err)
			req := &pulumirpc.DiffRequest{Olds: oldspb, OldInputs: oldspb, News: newspb}

			k := kubeProvider{}
			actual, err := k.DiffConfig(context.Background(), req)
			if wantErr != "" {
				assert.ErrorContains(t, err, wantErr)
				return
			}
			assert.NoError(t, err)

			assert.ElementsMatch(t, wantDiffs, actual.Diffs, "diff mismatch")
			assert.ElementsMatch(t, wantReplaces, actual.Replaces, "replace mismatch")
		})
	}
}

func TestCheckConfig_AlwaysRenderRequiresRenderYamlToDirectory(t *testing.T) {
	tests := []struct {
		name        string
		news        resource.PropertyMap
		wantFailure bool
		wantReason  string
	}{
		{
			name: "alwaysRender without renderYamlToDirectory should fail",
			news: resource.PropertyMap{
				"alwaysRender": resource.NewStringProperty("true"),
			},
			wantFailure: true,
			wantReason:  `"alwaysRender" requires "renderYamlToDirectory" to be set`,
		},
		{
			name: "alwaysRender with renderYamlToDirectory should succeed",
			news: resource.PropertyMap{
				"alwaysRender":          resource.NewStringProperty("true"),
				"renderYamlToDirectory": resource.NewStringProperty("/tmp/yaml"),
			},
			wantFailure: false,
		},
		{
			name: "renderYamlToDirectory without alwaysRender should succeed",
			news: resource.PropertyMap{
				"renderYamlToDirectory": resource.NewStringProperty("/tmp/yaml"),
			},
			wantFailure: false,
		},
		{
			name:        "neither set should succeed",
			news:        resource.PropertyMap{},
			wantFailure: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &kubeProvider{}

			newspb, err := plugin.MarshalProperties(tt.news, plugin.MarshalOptions{})
			require.NoError(t, err)

			req := &pulumirpc.CheckRequest{
				Urn:  "urn:pulumi:test::test::pulumi:providers:kubernetes::k8s",
				News: newspb,
			}

			resp, err := k.CheckConfig(context.Background(), req)
			require.NoError(t, err)

			if tt.wantFailure {
				require.Len(t, resp.Failures, 1)
				assert.Equal(t, "alwaysRender", resp.Failures[0].Property)
				assert.Equal(t, tt.wantReason, resp.Failures[0].Reason)
			} else {
				assert.Empty(t, resp.Failures)
			}
		})
	}
}
