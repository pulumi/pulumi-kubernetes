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

package provider

import (
	"testing"

	extensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

func TestParseCrdArgsWithName(t *testing.T) {
	args, err := parseCrdArgs([]string{"-n", "gateway-api", "-v", "1.2.1", "-c", "crds.yaml"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if args.PackageName != "gateway-api" {
		t.Errorf("expected package name %q, got %q", "gateway-api", args.PackageName)
	}
	if args.PackageVersion != "1.2.1" {
		t.Errorf("expected version %q, got %q", "1.2.1", args.PackageVersion)
	}
}

func TestParseCrdArgsWithoutName(t *testing.T) {
	args, err := parseCrdArgs([]string{"-v", "1.0.0", "-c", "crds.yaml"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if args.PackageName != "" {
		t.Errorf("expected empty package name, got %q", args.PackageName)
	}
}

func TestDerivePackageName(t *testing.T) {
	tests := []struct {
		name     string
		crds     []*extensionv1.CustomResourceDefinition
		expected string
	}{
		{
			"gateway API group",
			[]*extensionv1.CustomResourceDefinition{
				{Spec: extensionv1.CustomResourceDefinitionSpec{Group: "gateway.networking.k8s.io"}},
			},
			"gateway-networking",
		},
		{
			"cert-manager group",
			[]*extensionv1.CustomResourceDefinition{
				{Spec: extensionv1.CustomResourceDefinitionSpec{Group: "cert-manager.io"}},
			},
			"cert-manager",
		},
		{
			"multiple groups",
			[]*extensionv1.CustomResourceDefinition{
				{Spec: extensionv1.CustomResourceDefinitionSpec{Group: "gateway.networking.k8s.io"}},
				{Spec: extensionv1.CustomResourceDefinitionSpec{Group: "cert-manager.io"}},
			},
			"gateway-networking-cert-manager",
		},
		{
			"deduplicates same group",
			[]*extensionv1.CustomResourceDefinition{
				{Spec: extensionv1.CustomResourceDefinitionSpec{Group: "gateway.networking.k8s.io"}},
				{Spec: extensionv1.CustomResourceDefinitionSpec{Group: "gateway.networking.k8s.io"}},
			},
			"gateway-networking",
		},
		{
			"no groups falls back to default",
			[]*extensionv1.CustomResourceDefinition{
				{Spec: extensionv1.CustomResourceDefinitionSpec{Group: ""}},
			},
			defaultPackageName,
		},
		{
			"x-k8s.io suffix",
			[]*extensionv1.CustomResourceDefinition{
				{Spec: extensionv1.CustomResourceDefinitionSpec{Group: "multicluster.x-k8s.io"}},
			},
			"multicluster",
		},
		{
			"plain group without known suffix",
			[]*extensionv1.CustomResourceDefinition{
				{Spec: extensionv1.CustomResourceDefinitionSpec{Group: "example.com"}},
			},
			"example-com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := derivePackageName(tt.crds)
			if got != tt.expected {
				t.Errorf("derivePackageName() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestStripK8sSuffix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"gateway.networking.k8s.io", "gateway.networking"},
		{"cert-manager.io", "cert-manager"},
		{"multicluster.x-k8s.io", "multicluster"},
		{"example.com", "example.com"},
		{"storage", "storage"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := stripK8sSuffix(tt.input)
			if got != tt.expected {
				t.Errorf("stripK8sSuffix(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestSetCRDDefaults(t *testing.T) {
	tests := []struct {
		name     string
		crd      extensionv1.CustomResourceDefinition
		expected extensionv1.CustomResourceDefinition
	}{
		{
			"No defaults need to be set",
			extensionv1.CustomResourceDefinition{
				Spec: extensionv1.CustomResourceDefinitionSpec{
					Names: extensionv1.CustomResourceDefinitionNames{
						Singular: "foo",
						ListKind: "fooList",
						Kind:     "foo",
					},
				},
			},
			extensionv1.CustomResourceDefinition{
				Spec: extensionv1.CustomResourceDefinitionSpec{
					Names: extensionv1.CustomResourceDefinitionNames{
						Singular: "foo",
						ListKind: "fooList",
						Kind:     "foo",
					},
				},
			},
		},
		{
			"Need to set singular name",
			extensionv1.CustomResourceDefinition{
				Spec: extensionv1.CustomResourceDefinitionSpec{
					Names: extensionv1.CustomResourceDefinitionNames{
						ListKind: "fooCustomList",
						Kind:     "foo",
					},
				},
			},
			extensionv1.CustomResourceDefinition{
				Spec: extensionv1.CustomResourceDefinitionSpec{
					Names: extensionv1.CustomResourceDefinitionNames{
						Singular: "foo",
						ListKind: "fooCustomList",
						Kind:     "foo",
					},
				},
			},
		},
		{
			"Need to set list name",
			extensionv1.CustomResourceDefinition{
				Spec: extensionv1.CustomResourceDefinitionSpec{
					Names: extensionv1.CustomResourceDefinitionNames{
						Singular: "foocustomsingular",
						Kind:     "foo",
					},
				},
			},
			extensionv1.CustomResourceDefinition{
				Spec: extensionv1.CustomResourceDefinitionSpec{
					Names: extensionv1.CustomResourceDefinitionNames{
						Singular: "foocustomsingular",
						ListKind: "fooList",
						Kind:     "foo",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setCRDDefaults(
				&tt.crd,
			) //nolint:gosec // This is a false positive on older versions of golangci-lint. We are already using Go v1.22+

			if !equality.Semantic.DeepEqual(tt.crd, tt.expected) {
				t.Errorf("setCRDDefaults() got = %v, want %v", tt.crd, tt.expected)
			}
		})
	}
}

func TestFlattenOpenAPIArrayOfObjects(t *testing.T) {
	// Simulate a CRD with an array-of-objects property (like spec.listeners).
	sw := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Definitions: spec.Definitions{
				"io.k8s.networking.gateway.v1.GatewaySpec": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{"object"},
						Properties: map[string]spec.Schema{
							"listeners": {
								SchemaProps: spec.SchemaProps{
									Type: spec.StringOrArray{"array"},
									Items: &spec.SchemaOrArray{
										Schema: &spec.Schema{
											SchemaProps: spec.SchemaProps{
												Type: spec.StringOrArray{"object"},
												Properties: map[string]spec.Schema{
													"name":     {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"string"}}},
													"port":     {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"integer"}}},
													"protocol": {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"string"}}},
												},
											},
										},
									},
								},
							},
							"gatewayClassName": {
								SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"string"}},
							},
						},
					},
				},
			},
		},
	}

	err := flattenOpenAPI(sw)
	if err != nil {
		t.Fatalf("flattenOpenAPI returned error: %v", err)
	}

	// The inner object should be hoisted to a top-level definition.
	hoistedName := "io.k8s.networking.gateway.v1.GatewaySpecListeners"
	hoisted, ok := sw.Definitions[hoistedName]
	if !ok {
		t.Fatalf("expected hoisted definition %q not found; have: %v", hoistedName, definitionNames(sw))
	}

	if _, ok := hoisted.Properties["name"]; !ok {
		t.Error("hoisted definition missing 'name' property")
	}
	if _, ok := hoisted.Properties["port"]; !ok {
		t.Error("hoisted definition missing 'port' property")
	}
	if _, ok := hoisted.Properties["protocol"]; !ok {
		t.Error("hoisted definition missing 'protocol' property")
	}

	// The original array property's items should now be a $ref.
	parentDef := sw.Definitions["io.k8s.networking.gateway.v1.GatewaySpec"]
	listeners := parentDef.Properties["listeners"]
	if listeners.Items == nil || listeners.Items.Schema == nil {
		t.Fatal("listeners.Items.Schema is nil after flattening")
	}
	refURL := listeners.Items.Schema.Ref.GetURL()
	if refURL == nil {
		t.Fatal("listeners items should be a $ref after flattening")
	}
	expectedRef := "#/definitions/" + hoistedName
	if refURL.String() != expectedRef {
		t.Errorf("listeners items ref = %q, want %q", refURL.String(), expectedRef)
	}

	// Simple string property should be unchanged.
	className := parentDef.Properties["gatewayClassName"]
	if !className.Type.Contains("string") {
		t.Error("gatewayClassName should still be a string")
	}
}

func TestFlattenOpenAPINestedArrayOfObjects(t *testing.T) {
	// Simulate rules[].backendRefs[] — nested array-of-objects inside array-of-objects.
	sw := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Definitions: spec.Definitions{
				"io.k8s.networking.gateway.v1.HTTPRouteSpec": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{"object"},
						Properties: map[string]spec.Schema{
							"rules": {
								SchemaProps: spec.SchemaProps{
									Type: spec.StringOrArray{"array"},
									Items: &spec.SchemaOrArray{
										Schema: &spec.Schema{
											SchemaProps: spec.SchemaProps{
												Type: spec.StringOrArray{"object"},
												Properties: map[string]spec.Schema{
													"backendRefs": {
														SchemaProps: spec.SchemaProps{
															Type: spec.StringOrArray{"array"},
															Items: &spec.SchemaOrArray{
																Schema: &spec.Schema{
																	SchemaProps: spec.SchemaProps{
																		Type: spec.StringOrArray{"object"},
																		Properties: map[string]spec.Schema{
																			"name": {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"string"}}},
																			"port": {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"integer"}}},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	err := flattenOpenAPI(sw)
	if err != nil {
		t.Fatalf("flattenOpenAPI returned error: %v", err)
	}

	if _, ok := sw.Definitions["io.k8s.networking.gateway.v1.HTTPRouteSpecRules"]; !ok {
		t.Fatalf("expected Rules definition; have: %v", definitionNames(sw))
	}

	if _, ok := sw.Definitions["io.k8s.networking.gateway.v1.HTTPRouteSpecRulesBackendRefs"]; !ok {
		t.Fatalf("expected BackendRefs definition; have: %v", definitionNames(sw))
	}

	backendRefs := sw.Definitions["io.k8s.networking.gateway.v1.HTTPRouteSpecRulesBackendRefs"]
	if _, ok := backendRefs.Properties["name"]; !ok {
		t.Error("backendRefs missing 'name' property")
	}
	if _, ok := backendRefs.Properties["port"]; !ok {
		t.Error("backendRefs missing 'port' property")
	}
}

func definitionNames(sw *spec.Swagger) []string {
	var names []string
	for k := range sw.Definitions {
		names = append(names, k)
	}
	return names
}
