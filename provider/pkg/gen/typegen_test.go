// Copyright 2016-2024, Pulumi Corporation.
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

// nolint: lll
package gen

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	pschema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/txtar"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// gvkToString concatenates a group, version, and kind into a string in the format group.version.kind.
func gvkToString(group, version, kind string) string {
	return group + "." + version + "." + kind
}

// sliceToSet converts a slice of strings into a set of strings for easy membership checking.
func sliceToSet(slice []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, item := range slice {
		set[item] = struct{}{}
	}
	return set
}

// TestCreateGroups_IdentifyListKinds loads txtar files under testdata/identify-list-kinds and uses them to
// craft an unstructured map[string]any definitions file for createGroups.
// The goal of this test is to ensure we can accurately distinuguish between singletons and lists of kinds.
//
// The test files should contain the following files:
// - definitions: a JSON file containing the definitions
// - kinds: a YAML file containing a list of kinds that are singletons
// - listKinds: a YAML file containing a list of kinds that are lists/collections of singleton kinds
func TestCreateGroups_IdentifyListKinds(t *testing.T) {
	dir := filepath.Join("testdata/identify-list-kinds")
	tests, err := os.ReadDir(dir)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.Name(), func(t *testing.T) {
			archive, err := txtar.ParseFile(filepath.Join(dir, tt.Name()))
			require.NoError(t, err)

			var definitions map[string]any
			var kinds, listKinds map[string]struct{}

			for _, f := range archive.Files {
				var parsed []string
				switch f.Name {
				case "definitions":
					err := json.Unmarshal(f.Data, &definitions)
					require.NoError(t, err, f.Name)
				case "kinds":
					err = yaml.Unmarshal(f.Data, &parsed)
					require.NoError(t, err, f.Name)
					kinds = sliceToSet(parsed)
				case "listKinds":
					err = yaml.Unmarshal(f.Data, &parsed)
					require.NoError(t, err, f.Name)
					listKinds = sliceToSet(parsed)
				default:
					t.Fatal("unrecognized filename", f.Name)
				}
			}

			configGroups := createGroups(definitions, true)

			// Loop through all parsed kinds and ensure they are accounted for.
			for _, g := range configGroups {
				for _, v := range g.versions {
					for _, kind := range v.kinds {
						gvk := gvkToString(kind.gvk.Group, kind.gvk.Version, kind.gvk.Kind)
						if kind.isList {
							delete(listKinds, gvk)
						} else {
							delete(kinds, gvk)
						}
					}
				}
			}

			assert.Equal(t, 0, len(kinds), "kinds not found while parsing: %v", kinds)
			assert.Equal(t, 0, len(listKinds), "listKinds not found while parsing: %v", listKinds)
		})
	}
}

func TestAliasesForKind(t *testing.T) {
	tests := []struct {
		kind       string
		apiVersion string
		aliases    map[string][]any
		expected   []string
	}{
		// No aliases mapping, should return no types.
		{
			kind:       "Deployment",
			apiVersion: "apps/v1",
			aliases:    map[string][]any{},
			expected:   []string{},
		},
		{
			kind:       "Deployment",
			apiVersion: "apps/v1",
			aliases: map[string][]any{
				"Deployment": {
					"kubernetes:apps/v1beta1:Deployment",
					"kubernetes:extensions/v1beta1:Deployment",
				},
			},
			expected: []string{
				"kubernetes:apps/v1beta1:Deployment",
				"kubernetes:extensions/v1beta1:Deployment",
			},
		},
		{
			kind:       "CSIStorageCapacity",
			apiVersion: "storage.k8s.io/v1beta1",
			aliases: map[string][]any{
				"CSIStorageCapacity": {
					"kubernetes:storage.k8s.io/v1alpha1:CSIStorageCapacity",
				},
			},
			expected: []string{
				"kubernetes:storage.k8s.io/v1alpha1:CSIStorageCapacity",
				"kubernetes:storage.k8s.io/v1alpha1:CSIStorageCapacity",
			},
		},
		{
			kind:       "APIService",
			apiVersion: "apiregistration.k8s.io/v1",
			aliases: map[string][]any{
				"APIService": {
					"kubernetes:apiregistration.k8s.io/v1beta1:APIService",
				},
			},
			expected: []string{
				"kubernetes:apiregistration.k8s.io/v1beta1:APIService",
				"kubernetes:apiregistration/v1beta1:APIService",
				"kubernetes:apiregistration/v1:APIService",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			actual := aliasesForKind(tt.kind, tt.apiVersion, tt.aliases)
			assert.ElementsMatch(t, tt.expected, actual)
		})
	}
}
func TestCreateGroupsFromVersions(t *testing.T) {
	tests := []struct {
		name     string
		versions []VersionConfig
		expected []GroupConfig
	}{
		{
			name: "single group with single version",
			versions: []VersionConfig{
				{
					version: "v1",
					gv:      schema.GroupVersion{Group: "apps", Version: "v1"},
					kinds: []KindConfig{
						{kind: "Deployment"},
					},
				},
			},
			expected: []GroupConfig{
				{
					group: "apps",
					versions: []VersionConfig{
						{
							version: "v1",
							gv:      schema.GroupVersion{Group: "apps", Version: "v1"},
							kinds: []KindConfig{
								{kind: "Deployment"},
							},
						},
					},
				},
			},
		},
		{
			name: "single group with multiple versions",
			versions: []VersionConfig{
				{
					version: "v1",
					gv:      schema.GroupVersion{Group: "apps", Version: "v1"},
					kinds: []KindConfig{
						{kind: "Deployment"},
					},
				},
				{
					version: "v1beta1",
					gv:      schema.GroupVersion{Group: "apps", Version: "v1beta1"},
					kinds: []KindConfig{
						{kind: "Deployment"},
					},
				},
			},
			expected: []GroupConfig{
				{
					group: "apps",
					versions: []VersionConfig{
						{
							version: "v1",
							gv:      schema.GroupVersion{Group: "apps", Version: "v1"},
							kinds: []KindConfig{
								{kind: "Deployment"},
							},
						},
						{
							version: "v1beta1",
							gv:      schema.GroupVersion{Group: "apps", Version: "v1beta1"},
							kinds: []KindConfig{
								{kind: "Deployment"},
							},
						},
					},
				},
			},
		},
		{
			name: "multiple groups with multiple versions",
			versions: []VersionConfig{
				{
					version: "v1",
					gv:      schema.GroupVersion{Group: "apps", Version: "v1"},
					kinds: []KindConfig{
						{kind: "Deployment"},
					},
				},
				{
					version: "v1beta1",
					gv:      schema.GroupVersion{Group: "apps", Version: "v1beta1"},
					kinds: []KindConfig{
						{kind: "Deployment"},
					},
				},
				{
					version: "v1",
					gv:      schema.GroupVersion{Group: "core", Version: "v1"},
					kinds: []KindConfig{
						{kind: "Pod"},
					},
				},
			},
			expected: []GroupConfig{
				{
					group: "apps",
					versions: []VersionConfig{
						{
							version: "v1",
							gv:      schema.GroupVersion{Group: "apps", Version: "v1"},
							kinds: []KindConfig{
								{kind: "Deployment"},
							},
						},
						{
							version: "v1beta1",
							gv:      schema.GroupVersion{Group: "apps", Version: "v1beta1"},
							kinds: []KindConfig{
								{kind: "Deployment"},
							},
						},
					},
				},
				{
					group: "core",
					versions: []VersionConfig{
						{
							version: "v1",
							gv:      schema.GroupVersion{Group: "core", Version: "v1"},
							kinds: []KindConfig{
								{kind: "Pod"},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := createGroupsFromVersions(tt.versions)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
func TestCreateVersions(t *testing.T) {
	tests := []struct {
		name     string
		kinds    []KindConfig
		expected []VersionConfig
	}{
		{
			name: "single version with single kind",
			kinds: []KindConfig{
				{
					kind:       "Deployment",
					gvk:        schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
					apiVersion: "apps/v1",
				},
			},
			expected: []VersionConfig{
				{
					version:    "v1",
					gv:         schema.GroupVersion{Group: "apps", Version: "v1"},
					kinds:      []KindConfig{{kind: "Deployment", gvk: schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}, apiVersion: "apps/v1"}},
					apiVersion: "apps/v1",
				},
			},
		},
		{
			name: "single version with multiple kinds",
			kinds: []KindConfig{
				{
					kind:       "Deployment",
					gvk:        schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
					apiVersion: "apps/v1",
				},
				{
					kind:       "StatefulSet",
					gvk:        schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"},
					apiVersion: "apps/v1",
				},
			},
			expected: []VersionConfig{
				{
					version:    "v1",
					gv:         schema.GroupVersion{Group: "apps", Version: "v1"},
					kinds:      []KindConfig{{kind: "Deployment", gvk: schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}, apiVersion: "apps/v1"}, {kind: "StatefulSet", gvk: schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"}, apiVersion: "apps/v1"}},
					apiVersion: "apps/v1",
				},
			},
		},
		{
			name: "multiple versions with multiple kinds",
			kinds: []KindConfig{
				{
					kind:       "Deployment",
					gvk:        schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
					apiVersion: "apps/v1",
				},
				{
					kind:       "Deployment",
					gvk:        schema.GroupVersionKind{Group: "apps", Version: "v1beta1", Kind: "Deployment"},
					apiVersion: "apps/v1beta1",
				},
				{
					kind:       "Pod",
					gvk:        schema.GroupVersionKind{Group: "core", Version: "v1", Kind: "Pod"},
					apiVersion: "core/v1",
				},
			},
			expected: []VersionConfig{
				{
					version:    "v1",
					gv:         schema.GroupVersion{Group: "apps", Version: "v1"},
					kinds:      []KindConfig{{kind: "Deployment", gvk: schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}, apiVersion: "apps/v1"}},
					apiVersion: "apps/v1",
				},
				{
					version:    "v1beta1",
					gv:         schema.GroupVersion{Group: "apps", Version: "v1beta1"},
					kinds:      []KindConfig{{kind: "Deployment", gvk: schema.GroupVersionKind{Group: "apps", Version: "v1beta1", Kind: "Deployment"}, apiVersion: "apps/v1beta1"}},
					apiVersion: "apps/v1beta1",
				},
				{
					version:    "v1",
					gv:         schema.GroupVersion{Group: "core", Version: "v1"},
					kinds:      []KindConfig{{kind: "Pod", gvk: schema.GroupVersionKind{Group: "core", Version: "v1", Kind: "Pod"}, apiVersion: "core/v1"}},
					apiVersion: "core/v1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := createVersions(tt.kinds)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
func TestCreateDefinitions(t *testing.T) {
	definitionsJSON := map[string]any{
		"io.k8s.api.apps.v1.Deployment": map[string]any{
			"properties": map[string]any{
				"apiVersion": map[string]any{"type": "string"},
				"kind":       map[string]any{"type": "string"},
				"metadata":   map[string]any{"$ref": "#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta"},
			},
			"x-kubernetes-group-version-kind": []any{
				map[string]any{"group": "apps", "version": "v1", "kind": "Deployment"},
			},
		},
		"io.k8s.api.core.v1.Pod": map[string]any{
			"properties": map[string]any{
				"apiVersion": map[string]any{"type": "string"},
				"kind":       map[string]any{"type": "string"},
				"metadata":   map[string]any{"$ref": "#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta"},
			},
			"x-kubernetes-group-version-kind": []any{
				map[string]any{"group": "", "version": "v1", "kind": "Pod"},
			},
		},
	}

	canonicalGroups := map[string]string{
		"io.k8s.api.apps": "apps",
		"io.k8s.api.core": "core",
	}

	expected := []definition{
		{
			gvk:  schema.GroupVersionKind{Group: "io.k8s.api.apps", Version: "v1", Kind: "Deployment"},
			name: "io.k8s.api.apps.v1.Deployment",
			data: map[string]any{
				"properties": map[string]any{
					"apiVersion": map[string]any{"type": "string"},
					"kind":       map[string]any{"type": "string"},
					"metadata":   map[string]any{"$ref": "#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta"},
				},
				"x-kubernetes-group-version-kind": []any{
					map[string]any{"group": "apps", "version": "v1", "kind": "Deployment"},
				},
			},
			canonicalGroup: "apps",
		},
		{
			gvk:  schema.GroupVersionKind{Group: "io.k8s.api.core", Version: "v1", Kind: "Pod"},
			name: "io.k8s.api.core.v1.Pod",
			data: map[string]any{
				"properties": map[string]any{
					"apiVersion": map[string]any{"type": "string"},
					"kind":       map[string]any{"type": "string"},
					"metadata":   map[string]any{"$ref": "#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta"},
				},
				"x-kubernetes-group-version-kind": []any{
					map[string]any{"group": "", "version": "v1", "kind": "Pod"},
				},
			},
			canonicalGroup: "core",
		},
	}

	actual := createDefinitions(definitionsJSON, canonicalGroups)
	assert.ElementsMatch(t, expected, actual)
}

func TestCreateCanonicalGroups(t *testing.T) {
	tests := []struct {
		name            string
		definitionsJSON map[string]any
		expectedGroups  map[string]string
	}{
		{
			name: "basic canonical groups",
			definitionsJSON: map[string]any{
				"io.k8s.api.core.v1.Pod": map[string]any{
					"x-kubernetes-group-version-kind": []any{
						map[string]any{"group": "", "version": "v1", "kind": "Pod"},
					},
				},
				"io.k8s.api.apps.v1.Deployment": map[string]any{
					"x-kubernetes-group-version-kind": []any{
						map[string]any{"group": "apps", "version": "v1", "kind": "Deployment"},
					},
				},
			},
			expectedGroups: map[string]string{
				"io.k8s.api.core":                   "core",
				"io.k8s.api.apps":                   "apps",
				"io.k8s.apimachinery.pkg.apis.meta": "meta",
				"io.k8s.apimachinery.pkg":           "pkg",
			},
		},
		{
			name: "meta group",
			definitionsJSON: map[string]any{
				"io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta": map[string]any{
					"x-kubernetes-group-version-kind": []any{
						map[string]any{"group": "meta", "version": "v1", "kind": "ObjectMeta"},
					},
				},
			},
			expectedGroups: map[string]string{
				"io.k8s.apimachinery.pkg.apis.meta": "meta",
				"io.k8s.apimachinery.pkg":           "pkg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualGroups := createCanonicalGroups(tt.definitionsJSON)
			assert.Equal(t, tt.expectedGroups, actualGroups)
		})
	}
}
func TestMakeSchemaTypeSpec(t *testing.T) {
	tests := []struct {
		name             string
		prop             map[string]any
		canonicalGroups  map[string]string
		expectedTypeSpec pschema.TypeSpec
	}{
		{
			name: "simple string type",
			prop: map[string]any{
				"type": "string",
			},
			canonicalGroups:  map[string]string{},
			expectedTypeSpec: pschema.TypeSpec{Type: "string"},
		},
		{
			name: "array type",
			prop: map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
			},
			canonicalGroups: map[string]string{},
			expectedTypeSpec: pschema.TypeSpec{
				Type: "array",
				Items: &pschema.TypeSpec{
					Type: "string",
				},
			},
		},
		{
			name: "object type with additional properties",
			prop: map[string]any{
				"type": "object",
				"additionalProperties": map[string]any{
					"type": "string",
				},
			},
			canonicalGroups: map[string]string{},
			expectedTypeSpec: pschema.TypeSpec{
				Type: "object",
				AdditionalProperties: &pschema.TypeSpec{
					Type: "string",
				},
			},
		},
		{
			name: "preserve unknown fields",
			prop: map[string]any{
				"x-kubernetes-preserve-unknown-fields": true,
			},
			canonicalGroups: map[string]string{},
			expectedTypeSpec: pschema.TypeSpec{
				Type:                 "object",
				AdditionalProperties: &pschema.TypeSpec{Ref: "pulumi.json#/Any"},
			},
		},
		{
			name: "int or string",
			prop: map[string]any{
				"x-kubernetes-int-or-string": true,
			},
			canonicalGroups: map[string]string{},
			expectedTypeSpec: pschema.TypeSpec{
				OneOf: []pschema.TypeSpec{
					{Type: "integer"},
					{Type: "string"},
				},
			},
		},
		{
			name: "quantity type",
			prop: map[string]any{
				"$ref": "io.k8s.apimachinery.pkg.api.resource.Quantity",
			},
			canonicalGroups:  map[string]string{},
			expectedTypeSpec: pschema.TypeSpec{Type: "string"},
		},
		{
			name: "GVK reference",
			prop: map[string]any{
				"$ref": "#/definitions/io.k8s.api.apps.v1.Deployment",
			},
			canonicalGroups: map[string]string{
				"io.k8s.api.apps": "apps",
			},
			expectedTypeSpec: pschema.TypeSpec{
				Ref: "#/types/kubernetes:apps/v1:Deployment",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := makeSchemaTypeSpec(tt.prop, tt.canonicalGroups)
			assert.Equal(t, tt.expectedTypeSpec, actual)
		})
	}
}

func TestIsTopLevel(t *testing.T) {
	tests := []struct {
		name       string
		definition definition
		expected   bool
	}{
		{
			name: "top-level resource with ObjectMeta",
			definition: definition{
				gvk: schema.GroupVersionKind{Group: "io.k8s.api.apps", Version: "v1", Kind: "Deployment"},
				data: map[string]any{
					"properties": map[string]any{
						"metadata": map[string]any{
							"$ref": "#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta",
						},
					},
					"x-kubernetes-group-version-kind": []any{
						map[string]any{"group": "apps", "version": "v1", "kind": "Deployment"},
					},
				},
			},
			expected: true,
		},
		{
			name: "top-level resource with ListMeta",
			definition: definition{
				gvk: schema.GroupVersionKind{Group: "io.k8s.api.core", Version: "v1", Kind: "PodList"},
				data: map[string]any{
					"properties": map[string]any{
						"metadata": map[string]any{
							"$ref": "#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.ListMeta",
						},
					},
					"x-kubernetes-group-version-kind": []any{
						map[string]any{"group": "core", "version": "v1", "kind": "PodList"},
					},
				},
			},
			expected: true,
		},
		{
			name: "non-top-level resource",
			definition: definition{
				gvk: schema.GroupVersionKind{Group: "io.k8s.api.apps", Version: "v1", Kind: "ReplicaSet"},
				data: map[string]any{
					"properties": map[string]any{
						"spec": map[string]any{
							"type": "object",
						},
					},
					"x-kubernetes-group-version-kind": []any{
						map[string]any{"group": "apps", "version": "v1", "kind": "ReplicaSet"},
					},
				},
			},
			expected: false,
		},
		{
			name: "imperative resource type",
			definition: definition{
				gvk: schema.GroupVersionKind{Group: "io.k8s.api.authentication", Version: "v1", Kind: "TokenRequest"},
				data: map[string]any{
					"x-kubernetes-group-version-kind": []any{
						map[string]any{"group": "authentication", "version": "v1", "kind": "TokenRequest"},
					},
				},
			},
			expected: false,
		},
		{
			name: "missing properties",
			definition: definition{
				gvk: schema.GroupVersionKind{Group: "io.k8s.api.apps", Version: "v1", Kind: "Deployment"},
				data: map[string]any{
					"x-kubernetes-group-version-kind": []any{
						map[string]any{"group": "apps", "version": "v1", "kind": "Deployment"},
					},
				},
			},
			expected: false,
		},
		{
			name: "missing metadata",
			definition: definition{
				gvk: schema.GroupVersionKind{Group: "io.k8s.api.apps", Version: "v1", Kind: "Deployment"},
				data: map[string]any{
					"properties": map[string]any{
						"spec": map[string]any{
							"type": "object",
						},
					},
					"x-kubernetes-group-version-kind": []any{
						map[string]any{"group": "apps", "version": "v1", "kind": "Deployment"},
					},
				},
			},
			expected: false,
		},
		{
			name: "missing $ref in metadata",
			definition: definition{
				gvk: schema.GroupVersionKind{Group: "io.k8s.api.apps", Version: "v1", Kind: "Deployment"},
				data: map[string]any{
					"properties": map[string]any{
						"metadata": map[string]any{
							"type": "object",
						},
					},
					"x-kubernetes-group-version-kind": []any{
						map[string]any{"group": "apps", "version": "v1", "kind": "Deployment"},
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.definition.isTopLevel()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestGVKFromRef(t *testing.T) {
	tests := []struct {
		ref      string
		expected schema.GroupVersionKind
	}{
		{
			ref: "io.k8s.api.apps.v1.Deployment",
			expected: schema.GroupVersionKind{
				Group:   "io.k8s.api.apps",
				Version: "v1",
				Kind:    "Deployment",
			},
		},
		{
			ref: "io.k8s.api.core.v1.Pod",
			expected: schema.GroupVersionKind{
				Group:   "io.k8s.api.core",
				Version: "v1",
				Kind:    "Pod",
			},
		},
		{
			ref: "io.k8s.api.extensions.v1beta1.Ingress",
			expected: schema.GroupVersionKind{
				Group:   "io.k8s.api.extensions",
				Version: "v1beta1",
				Kind:    "Ingress",
			},
		},
		{
			ref: "io.k8s.api.networking.v1.NetworkPolicy",
			expected: schema.GroupVersionKind{
				Group:   "io.k8s.api.networking",
				Version: "v1",
				Kind:    "NetworkPolicy",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.ref, func(t *testing.T) {
			actual := GVKFromRef(tt.ref)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
