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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	extensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kube-openapi/pkg/validation/spec"

	pulumischema "github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

func TestParseCrdArgs(t *testing.T) {
	parseTests := []struct {
		name    string
		args    []string
		wantErr string // substring expected in the error message; empty means no error
		check   func(t *testing.T, args *ParameterizedArgs)
	}{
		{
			name: "with name",
			args: []string{"name=gateway-api", "version=1.2.1", "crd-manifest=crds.yaml"},
			check: func(t *testing.T, args *ParameterizedArgs) {
				if args.ExtensionName != "gateway-api" {
					t.Errorf("expected extension name %q, got %q", "gateway-api", args.ExtensionName)
				}
				if args.ExtensionVersion != "1.2.1" {
					t.Errorf("expected version %q, got %q", "1.2.1", args.ExtensionVersion)
				}
			},
		},
		{
			name: "without name",
			args: []string{"version=1.0.0", "crd-manifest=crds.yaml"},
			check: func(t *testing.T, args *ParameterizedArgs) {
				if args.ExtensionName != "" {
					t.Errorf("expected empty extension name, got %q", args.ExtensionName)
				}
			},
		},
		{
			name: "defaults version when omitted",
			args: []string{"name=gateway-api", "crd-manifest=crds.yaml"},
			check: func(t *testing.T, args *ParameterizedArgs) {
				if args.ExtensionVersion != defaultExtensionVersion {
					t.Errorf("expected default version %q, got %q", defaultExtensionVersion, args.ExtensionVersion)
				}
			},
		},
		{
			name: "collects multiple manifests",
			args: []string{"crd-manifest=a.yaml", "crd-manifest=b.yaml"},
			check: func(t *testing.T, args *ParameterizedArgs) {
				want := []string{"a.yaml", "b.yaml"}
				if len(args.CRDManifestPaths) != len(want) {
					t.Fatalf("expected %d manifest paths, got %v", len(want), args.CRDManifestPaths)
				}
				for i, path := range want {
					if args.CRDManifestPaths[i] != path {
						t.Errorf("manifest path %d = %q, want %q", i, args.CRDManifestPaths[i], path)
					}
				}
			},
		},
		{
			name:    "rejects unknown key",
			args:    []string{"crd-manifest=crds.yaml", "flavor=spicy"},
			wantErr: "flavor",
		},
		{
			name:    "rejects malformed arg",
			args:    []string{"crd-manifest=crds.yaml", "noequalssign"},
			wantErr: "key=value",
		},
		{
			name:    "requires crd-manifest",
			args:    []string{"name=gateway-api", "version=1.0.0"},
			wantErr: "crd-manifest",
		},
	}

	for _, tt := range parseTests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := parseCrdArgs(tt.args)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("error = %q, want it to contain %q", err.Error(), tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.check != nil {
				tt.check(t, args)
			}
		})
	}
}

func TestDeriveExtensionName(t *testing.T) {
	tests := []struct {
		name     string
		paths    []string
		expected string
	}{
		{"single yaml file", []string{"gateway-api-crds.yaml"}, "gateway-api-crds"},
		{"yml extension", []string{"cert-manager.yml"}, "cert-manager"},
		{"path with directory", []string{"manifests/gateway.yaml"}, "gateway"},
		{"first path wins", []string{"gateway-api.yaml", "cert-manager.yaml"}, "gateway-api"},
		{"empty list falls back to default", nil, defaultExtensionName},
		{"empty path falls back to default", []string{""}, defaultExtensionName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deriveExtensionName(tt.paths)
			if got != tt.expected {
				t.Errorf("deriveExtensionName() = %q, want %q", got, tt.expected)
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

func TestCRDToOpenAPIPreservesNullableFieldTypes(t *testing.T) {
	stringItem := extensionv1.JSONSchemaProps{Type: "string"}
	crd := &extensionv1.CustomResourceDefinition{
		Spec: extensionv1.CustomResourceDefinitionSpec{
			Group: "example.com",
			Names: extensionv1.CustomResourceDefinitionNames{Kind: "Widget", Plural: "widgets"},
			Scope: extensionv1.NamespaceScoped,
			Versions: []extensionv1.CustomResourceDefinitionVersion{{
				Name:    "v1",
				Served:  true,
				Storage: true,
				Schema: &extensionv1.CustomResourceValidation{
					OpenAPIV3Schema: &extensionv1.JSONSchemaProps{
						Type: "object",
						Properties: map[string]extensionv1.JSONSchemaProps{
							"timeout": {Type: "integer", Nullable: true},
							"hosts": {
								Type:     "array",
								Nullable: true,
								Items:    &extensionv1.JSONSchemaPropsOrArray{Schema: &stringItem},
							},
						},
					},
				},
			}},
		},
	}

	specs, err := crdToOpenAPI(crd)
	if err != nil {
		t.Fatalf("crdToOpenAPI returned error: %v", err)
	}
	if len(specs) != 1 {
		t.Fatalf("expected 1 spec, got %d", len(specs))
	}

	widget, ok := specs[0].Definitions["com.example.v1.Widget"]
	if !ok {
		t.Fatalf("expected Widget definition; have: %v", definitionNames(specs[0]))
	}

	timeout := widget.Properties["timeout"]
	if !timeout.Type.Contains("integer") {
		t.Errorf("nullable scalar lost its type: %v", timeout.Type)
	}

	hosts := widget.Properties["hosts"]
	switch {
	case !hosts.Type.Contains("array"):
		t.Errorf("nullable array lost its type: %v", hosts.Type)
	case hosts.Items == nil || hosts.Items.Schema == nil:
		t.Error("nullable array lost its items")
	case !hosts.Items.Schema.Type.Contains("string"):
		t.Errorf("nullable array lost its element type: %v", hosts.Items.Schema.Type)
	}
}

// mustParameterizeExtension builds a CRD, converts it with the provider's own crdToOpenAPI
// and mergeSpecs, then calls Parameterize with the result.
func mustParameterizeExtension(t *testing.T, k *kubeProvider, extension, group, kind string) {
	t.Helper()

	crd := &extensionv1.CustomResourceDefinition{
		Spec: extensionv1.CustomResourceDefinitionSpec{
			Group: group,
			Names: extensionv1.CustomResourceDefinitionNames{Kind: kind, Plural: strings.ToLower(kind) + "s"},
			Scope: extensionv1.NamespaceScoped,
			Versions: []extensionv1.CustomResourceDefinitionVersion{{
				Name:    "v1",
				Served:  true,
				Storage: true,
				Schema: &extensionv1.CustomResourceValidation{
					OpenAPIV3Schema: &extensionv1.JSONSchemaProps{
						Type:       "object",
						Properties: map[string]extensionv1.JSONSchemaProps{"size": {Type: "string"}},
					},
				},
			}},
		},
	}

	specs, err := crdToOpenAPI(crd)
	if err != nil {
		t.Fatalf("crdToOpenAPI(%s) returned error: %v", kind, err)
	}
	merged, err := mergeSpecs(specs)
	if err != nil {
		t.Fatalf("mergeSpecs(%s) returned error: %v", kind, err)
	}
	paramBytes, err := json.Marshal(merged)
	if err != nil {
		t.Fatalf("marshalling %s OpenAPI spec: %v", kind, err)
	}

	resp, err := k.Parameterize(context.Background(), &pulumirpc.ParameterizeRequest{
		Parameters: &pulumirpc.ParameterizeRequest_Value{
			Value: &pulumirpc.ParameterizeRequest_ParametersValue{
				Name:    extension,
				Version: "1.0.0",
				Value:   paramBytes,
			},
		},
	})
	if err != nil {
		t.Fatalf("Parameterize(%s) returned error: %v", extension, err)
	}
	if resp.GetName() != extension {
		t.Fatalf("Parameterize(%s) returned name %q", extension, resp.GetName())
	}
}

func TestParameterizeTwoExtensionsKeepsBothServed(t *testing.T) {
	k, err := makeKubeProvider(nil, "kubernetes", "4.34.1", []byte("{}"), []byte("{}"), []byte("{}"))
	if err != nil {
		t.Fatalf("makeKubeProvider returned error: %v", err)
	}

	mustParameterizeExtension(t, k, "dragon-ext", "dragon.multiext.pulumi.com", "Dragon")
	mustParameterizeExtension(t, k, "unicorn-ext", "unicorn.multiext.pulumi.com", "Unicorn")

	tokens := map[string]schema.GroupVersionKind{
		"dragon-ext:dragon.multiext.pulumi.com/v1:Dragon": {
			Group: "dragon.multiext.pulumi.com", Version: "v1", Kind: "Dragon",
		},
		"unicorn-ext:unicorn.multiext.pulumi.com/v1:Unicorn": {
			Group: "unicorn.multiext.pulumi.com", Version: "v1", Kind: "Unicorn",
		},
		"kubernetes:core/v1:ConfigMap": {
			Group: "", Version: "v1", Kind: "ConfigMap",
		},
	}
	for token, want := range tokens {
		got, err := k.gvkFromTypeToken(token)
		if err != nil {
			t.Errorf("gvkFromTypeToken(%q) returned error: %v", token, err)
			continue
		}
		if got != want {
			t.Errorf("gvkFromTypeToken(%q) = %v, want %v", token, got, want)
		}
	}

	if _, err := k.gvkFromTypeToken("griffin-ext:griffin.multiext.pulumi.com/v1:Griffin"); err == nil {
		t.Error("gvkFromTypeToken accepted a package the provider was never parameterized with")
	}
}

const hyphenatedPropertyCRD = `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ciliumidentities.cilium.io
spec:
  group: cilium.io
  names:
    kind: CiliumIdentity
    listKind: CiliumIdentityList
    plural: ciliumidentities
    singular: ciliumidentity
  scope: Cluster
  versions:
    - name: v2
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                replicas:
                  type: integer
                security-labels:
                  type: object
                  additionalProperties:
                    type: string
`

func TestParameterizePreservesHyphenatedPropertyNames(t *testing.T) {
	manifest := filepath.Join(t.TempDir(), "crd.yaml")
	require.NoError(t, os.WriteFile(manifest, []byte(hyphenatedPropertyCRD), 0o600))

	fromArgs := &kubeProvider{name: "kubernetes", version: "4.0.0"}
	_, err := fromArgs.parameterizeRequestArgs(&pulumirpc.ParameterizeRequest_Args{
		Args: &pulumirpc.ParameterizeRequest_ParametersArgs{
			Args: []string{"name=cilium", "version=1.0.0", "crd-manifest=" + manifest},
		},
	})
	require.NoError(t, err)

	argsSchema := fromArgs.crdSchemas.get("cilium", "1.0.0")
	require.NotNil(t, argsSchema)
	assertHyphenatedProperty(t, argsSchema)

	require.NotNil(t, argsSchema.ExtensionParameterization)
	savedParameter := argsSchema.ExtensionParameterization.Parameter
	require.NotEmpty(t, savedParameter)

	fromValue := &kubeProvider{name: "kubernetes", version: "4.0.0"}
	_, err = fromValue.parameterizeRequestValue(&pulumirpc.ParameterizeRequest_Value{
		Value: &pulumirpc.ParameterizeRequest_ParametersValue{
			Name:    "cilium",
			Version: "1.0.0",
			Value:   savedParameter,
		},
	})
	require.NoError(t, err)

	valueSchema := fromValue.crdSchemas.get("cilium", "1.0.0")
	require.NotNil(t, valueSchema)
	assertHyphenatedProperty(t, valueSchema)
}

func assertHyphenatedProperty(t *testing.T, pkg *pulumischema.PackageSpec) {
	t.Helper()

	for token, typ := range pkg.Types {
		if _, ok := typ.Properties["security-labels"]; ok {
			assert.NotContains(t, typ.Properties, "security_labels",
				"%s should carry only the wire name", token)
			return
		}
	}

	t.Fatalf("no type declared a %q property; the schema has %d types", "security-labels", len(pkg.Types))
}
