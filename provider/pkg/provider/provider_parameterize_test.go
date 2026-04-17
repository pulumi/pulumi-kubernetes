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
	"context"
	"encoding/json"
	"strings"
	"testing"

	extensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/kube-openapi/pkg/validation/spec"

	pulumischema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
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
		paths    []string
		expected string
	}{
		{"single yaml file", []string{"/tmp/gateway-api.yaml"}, "gateway-api"},
		{"single yml file", []string{"/tmp/cert-manager.yml"}, "cert-manager"},
		{"multiple files uses first", []string{"/tmp/gateway-api.yaml", "/tmp/cert-manager.yaml"}, "gateway-api"},
		{"no paths falls back to default", nil, defaultPackageName},
		{"empty string path falls back to default", []string{""}, defaultPackageName},
		{"nested path uses base name", []string{"/home/user/crds/istio-crds.yaml"}, "istio-crds"},
		{"no extension", []string{"/tmp/my-crds"}, "my-crds"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := derivePackageName(tt.paths)
			if got != tt.expected {
				t.Errorf("derivePackageName() = %q, want %q", got, tt.expected)
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
	listenerSchema := spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type: spec.StringOrArray{"object"},
			Properties: map[string]spec.Schema{
				"name":     {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"string"}}},
				"port":     {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"integer"}}},
				"protocol": {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"string"}}},
			},
		},
	}
	sw := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Definitions: spec.Definitions{
				"io.k8s.networking.gateway.v1.GatewaySpec": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{"object"},
						Properties: map[string]spec.Schema{
							"listeners": {
								SchemaProps: spec.SchemaProps{
									Type:  spec.StringOrArray{"array"},
									Items: &spec.SchemaOrArray{Schema: &listenerSchema},
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
	backendRefSchema := spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type: spec.StringOrArray{"object"},
			Properties: map[string]spec.Schema{
				"name": {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"string"}}},
				"port": {SchemaProps: spec.SchemaProps{Type: spec.StringOrArray{"integer"}}},
			},
		},
	}
	ruleSchema := spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type: spec.StringOrArray{"object"},
			Properties: map[string]spec.Schema{
				"backendRefs": {
					SchemaProps: spec.SchemaProps{
						Type:  spec.StringOrArray{"array"},
						Items: &spec.SchemaOrArray{Schema: &backendRefSchema},
					},
				},
			},
		},
	}
	sw := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Definitions: spec.Definitions{
				"io.k8s.networking.gateway.v1.HTTPRouteSpec": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{"object"},
						Properties: map[string]spec.Schema{
							"rules": {
								SchemaProps: spec.SchemaProps{
									Type:  spec.StringOrArray{"array"},
									Items: &spec.SchemaOrArray{Schema: &ruleSchema},
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

// TestParameterizeRoundTrip is the back-to-back proof that the parameterization
// pipeline works end-to-end. It simulates what the engine does on two
// successive pulumi runs:
//
//  1. First run: engine invokes Parameterize(Args) with the CLI args that
//     `pulumi package add kubernetes -- -v 1.2.1 -c crds.yaml` produces. The
//     provider generates a schema and persists its Parameter bytes in Pulumi.yaml.
//  2. Subsequent run: engine invokes Parameterize(Value) with those saved
//     bytes, and the provider must reconstruct an identical schema without
//     needing the original CRD file on disk.
//
// Both paths are asserted to produce a schema with:
//   - the derived package name ("gateway-like-crd") as PackageSpec.Name,
//   - resource/type tokens prefixed with "gateway-like-crd:" (not "kubernetes:"),
//   - $refs that don't leak the base "kubernetes:" prefix,
//   - a parameterized Go import base path so the SDK can coexist with the base.
//
// This is the regression guard against the classes of bug that came up during
// development: the Parameterize(Value) stub returning nil, the rewriter
// missing token keys, and identity-leak bugs where some corner of the schema
// kept the base provider's name.
func TestParameterizeRoundTrip(t *testing.T) {
	k := &kubeProvider{
		name:    "kubernetes",
		version: "4.0.0-dev",
	}

	const (
		wantPkg     = "gateway-like-crd"
		wantVersion = "1.2.1"
	)

	// --- Args path: first-ever run ---
	argsResp, err := k.Parameterize(context.Background(), &pulumirpc.ParameterizeRequest{
		Parameters: &pulumirpc.ParameterizeRequest_Args{
			Args: &pulumirpc.ParameterizeRequest_ParametersArgs{
				Args: []string{
					"-v", wantVersion,
					"-c", "testdata/parameterize/gateway-like-crd.yaml",
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Parameterize(Args) failed: %v", err)
	}
	if argsResp.Name != wantPkg {
		t.Errorf("response name: got %q, want %q", argsResp.Name, wantPkg)
	}
	if argsResp.Version != wantVersion {
		t.Errorf("response version: got %q, want %q", argsResp.Version, wantVersion)
	}

	originalSchema := k.crdSchemas.get(wantPkg, wantVersion)
	if originalSchema == nil {
		t.Fatalf("no schema cached for %s@%s after Parameterize(Args)", wantPkg, wantVersion)
	}
	assertParameterizedIdentity(t, originalSchema, wantPkg)

	// Grab the parameter bytes the engine would persist in Pulumi.yaml.
	if originalSchema.Parameterization == nil {
		t.Fatal("schema.Parameterization is nil")
	}
	paramBytes := originalSchema.Parameterization.Parameter
	if len(paramBytes) == 0 {
		t.Fatal("schema.Parameterization.Parameter is empty")
	}

	// --- Clear the cache to simulate a fresh provider process ---
	k.crdSchemas = parameterizedPackageMap{}

	// --- Value path: subsequent run, engine replays saved parameters ---
	valueResp, err := k.Parameterize(context.Background(), &pulumirpc.ParameterizeRequest{
		Parameters: &pulumirpc.ParameterizeRequest_Value{
			Value: &pulumirpc.ParameterizeRequest_ParametersValue{
				Name:    wantPkg,
				Version: wantVersion,
				Value:   paramBytes,
			},
		},
	})
	if err != nil {
		t.Fatalf("Parameterize(Value) failed: %v", err)
	}
	if valueResp.Name != wantPkg || valueResp.Version != wantVersion {
		t.Errorf("Value response: got %s@%s, want %s@%s",
			valueResp.Name, valueResp.Version, wantPkg, wantVersion)
	}

	rebuilt := k.crdSchemas.get(wantPkg, wantVersion)
	if rebuilt == nil {
		t.Fatalf("no schema cached for %s@%s after Parameterize(Value)", wantPkg, wantVersion)
	}
	assertParameterizedIdentity(t, rebuilt, wantPkg)

	// Both runs must produce the same resource and type sets. We don't require
	// byte-identical schemas (some maps may iterate differently on the margins),
	// but the set of tokens must be identical.
	if diff := stringSetDiff(mapKeys(originalSchema.Resources), mapKeys(rebuilt.Resources)); diff != "" {
		t.Errorf("Resources set differs between Args and Value runs: %s", diff)
	}
	if diff := stringSetDiff(mapKeys(originalSchema.Types), mapKeys(rebuilt.Types)); diff != "" {
		t.Errorf("Types set differs between Args and Value runs: %s", diff)
	}
}

// assertParameterizedIdentity verifies the fingerprint of a parameterized
// package: name, token prefixes, ref prefixes, and Go import base path.
func assertParameterizedIdentity(t *testing.T, schema *pulumischema.PackageSpec, pkg string) {
	t.Helper()

	if schema.Name != pkg {
		t.Errorf("PackageSpec.Name: got %q, want %q", schema.Name, pkg)
	}

	prefix := pkg + ":"
	for tok := range schema.Resources {
		if !strings.HasPrefix(tok, prefix) {
			t.Errorf("resource token %q must be prefixed with %q", tok, prefix)
		}
	}
	for tok := range schema.Types {
		if !strings.HasPrefix(tok, prefix) {
			t.Errorf("type token %q must be prefixed with %q", tok, prefix)
		}
	}

	// Walk the JSON representation for lingering base-provider refs. A
	// byte-level scan is the most thorough check — it catches refs nested in
	// OneOf, Items, AdditionalProperties, etc., without needing to walk the
	// struct by hand.
	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("marshalling schema: %v", err)
	}
	for _, leaked := range []string{
		`"#/types/kubernetes:`,
		`"#/resources/kubernetes:`,
	} {
		if strings.Contains(string(data), leaked) {
			t.Errorf("schema still contains base-provider ref prefix %q — identity rewrite is incomplete", leaked)
		}
	}

	// Go import base path should name the parameterized package.
	goLangBytes, ok := schema.Language["go"]
	if !ok {
		t.Fatal("schema.Language[\"go\"] missing")
	}
	var goLang map[string]any
	if err := json.Unmarshal(goLangBytes, &goLang); err != nil {
		t.Fatalf("unmarshalling Go language settings: %v", err)
	}
	importPath, _ := goLang["importBasePath"].(string)
	if !strings.Contains(importPath, pkg) {
		t.Errorf("Go importBasePath %q does not mention parameterized package %q", importPath, pkg)
	}
	if strings.Contains(importPath, "/go/kubernetes") {
		t.Errorf("Go importBasePath %q still points at the base kubernetes SDK path", importPath)
	}
}

func mapKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// stringSetDiff returns "" when a and b represent the same set, or a
// human-readable description of the symmetric difference otherwise.
func stringSetDiff(a, b []string) string {
	seen := make(map[string]int, len(a)+len(b))
	for _, s := range a {
		seen[s] |= 1
	}
	for _, s := range b {
		seen[s] |= 2
	}
	var onlyA, onlyB []string
	for s, bits := range seen {
		switch bits {
		case 1:
			onlyA = append(onlyA, s)
		case 2:
			onlyB = append(onlyB, s)
		}
	}
	if len(onlyA) == 0 && len(onlyB) == 0 {
		return ""
	}
	return "only in first: " + strings.Join(onlyA, ",") + "; only in second: " + strings.Join(onlyB, ",")
}
