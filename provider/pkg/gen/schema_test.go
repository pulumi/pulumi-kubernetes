// Copyright 2016-2026, Pulumi Corporation.
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

package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pschema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
)

func TestListInputsSpec(t *testing.T) {
	cases := []struct {
		name            string
		kind            string
		expectNamespace bool
	}{
		{"namespaced kind includes namespace", "ConfigMap", true},
		{"namespaced kind (apps/v1) includes namespace", "Deployment", true},
		{"cluster-scoped kind omits namespace", "Namespace", false},
		{"cluster-scoped Node omits namespace", "Node", false},
		{"cluster-scoped ClusterRole omits namespace", "ClusterRole", false},
		{"unknown kind keeps namespace conservatively", "SomeCustomResource", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			spec := listInputsSpec(tc.kind)
			require.NotNil(t, spec)
			assert.Equal(t, "object", spec.Type)
			assert.Contains(t, spec.Properties, "name")
			assert.Contains(t, spec.Properties, "labelSelector")
			assert.Contains(t, spec.Properties, "fieldSelector")
			if tc.expectNamespace {
				assert.Contains(t, spec.Properties, "namespace", "namespaced kinds must advertise the namespace filter")
			} else {
				assert.NotContains(t, spec.Properties, "namespace", "cluster-scoped kinds must omit the namespace filter")
			}
		})
	}
}

func TestExtensionResourcesAliasBaseProviderToken(t *testing.T) {
	swagger := map[string]any{
		"definitions": map[string]any{
			"io.k8s.api.gateway.v1.Gateway": map[string]any{
				"properties": map[string]any{
					"apiVersion": map[string]any{"type": "string"},
					"kind":       map[string]any{"type": "string"},
					"metadata":   map[string]any{"$ref": "#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta"},
					"spec":       map[string]any{"type": "object"},
				},
				"x-kubernetes-group-version-kind": []any{
					map[string]any{"group": "gateway", "version": "v1", "kind": "Gateway"},
				},
			},
		},
	}

	pkg := PulumiSchema(swagger,
		WithExtensionName("crdfoo"),
		WithParameterization(&pschema.ExtensionParameterizationSpec{
			BaseProvider: pschema.BaseProviderRefSpec{Name: "kubernetes", Version: "4.33.0"},
		}),
	)

	require.Equal(t, "crdfoo", pkg.Name)

	res, ok := pkg.Resources["crdfoo:gateway/v1:Gateway"]
	require.True(t, ok, "extension resource must be tokened under the extension package")
	assert.NotContains(t, pkg.Resources, "kubernetes:gateway/v1:Gateway",
		"the base-namespaced token must not be emitted as a resource of its own")

	assert.Contains(t, res.Aliases, pschema.AliasSpec{Type: "kubernetes:gateway/v1:Gateway"},
		"extension resource must alias the base token so crd2pulumi state is adopted rather than replaced")
}

func TestBaseProviderResourcesHaveNoExtensionAlias(t *testing.T) {
	swagger := map[string]any{
		"definitions": map[string]any{
			"io.k8s.api.gateway.v1.Gateway": map[string]any{
				"properties": map[string]any{
					"apiVersion": map[string]any{"type": "string"},
					"kind":       map[string]any{"type": "string"},
					"metadata":   map[string]any{"$ref": "#/definitions/io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta"},
					"spec":       map[string]any{"type": "object"},
				},
				"x-kubernetes-group-version-kind": []any{
					map[string]any{"group": "gateway", "version": "v1", "kind": "Gateway"},
				},
			},
		},
	}

	pkg := PulumiSchema(swagger)

	res, ok := pkg.Resources["kubernetes:gateway/v1:Gateway"]
	require.True(t, ok, "base schema must keep kubernetes-namespaced tokens")

	assert.NotContains(t, res.Aliases, pschema.AliasSpec{Type: "kubernetes:gateway/v1:Gateway"},
		"a resource must not alias its own token when generating the base provider schema")
}
