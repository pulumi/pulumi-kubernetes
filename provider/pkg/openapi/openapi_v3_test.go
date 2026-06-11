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

package openapi

import (
	"fmt"
	"testing"

	openapi_v2 "github.com/google/gnostic-models/openapiv2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientopenapi "k8s.io/client-go/openapi"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients/fake"
)

// fakeV2Disco wraps an openapi_v2.Document to satisfy discovery.OpenAPISchemaInterface.
type fakeV2Disco struct {
	doc *openapi_v2.Document
}

func (f *fakeV2Disco) OpenAPISchema() (*openapi_v2.Document, error) {
	return f.doc, nil
}

// appsV1Schema is a minimal OpenAPI v3 document for the apps/v1 group version.
const appsV1Schema = `{
  "openapi": "3.0.0",
  "info": {"title": "Kubernetes", "version": "unversioned"},
  "paths": {},
  "components": {
    "schemas": {
      "io.k8s.api.apps.v1.Deployment": {
        "type": "object",
        "x-kubernetes-group-version-kind": [
          {"group": "apps", "version": "v1", "kind": "Deployment"}
        ],
        "properties": {
          "apiVersion": {"type": "string"},
          "kind": {"type": "string"},
          "metadata": {"type": "object"},
          "spec": {"type": "object"}
        }
      }
    }
  }
}`

// coreV1Schema is a minimal OpenAPI v3 document for the core v1 group version.
const coreV1Schema = `{
  "openapi": "3.0.0",
  "info": {"title": "Kubernetes", "version": "unversioned"},
  "paths": {},
  "components": {
    "schemas": {
      "io.k8s.api.core.v1.ConfigMap": {
        "type": "object",
        "x-kubernetes-group-version-kind": [
          {"group": "", "version": "v1", "kind": "ConfigMap"}
        ],
        "properties": {
          "apiVersion": {"type": "string"},
          "kind": {"type": "string"},
          "metadata": {"type": "object"},
          "data": {"type": "object"}
        }
      }
    }
  }
}`

// fakeV3Client is a test-only openapi.Client backed by a map of path → JSON bytes.
type fakeV3Client struct {
	paths map[string][]byte
}

func (f *fakeV3Client) Paths() (map[string]clientopenapi.GroupVersion, error) {
	out := make(map[string]clientopenapi.GroupVersion, len(f.paths))
	for path, data := range f.paths {
		out[path] = &fakeV3GV{data: data}
	}
	return out, nil
}

type fakeV3GV struct{ data []byte }

func (g *fakeV3GV) Schema(_ string) ([]byte, error) { return g.data, nil }
func (g *fakeV3GV) ServerRelativeURL() string        { return "" }

type errorV3GV struct{ err error }

func (g *errorV3GV) Schema(_ string) ([]byte, error) { return nil, g.err }
func (g *errorV3GV) ServerRelativeURL() string        { return "" }

// partialErrorV3Client wraps fakeV3Client and injects an additional error GV.
type partialErrorV3Client struct {
	delegate  *fakeV3Client
	errorPath string
	err       error
}

func (p *partialErrorV3Client) Paths() (map[string]clientopenapi.GroupVersion, error) {
	paths, err := p.delegate.Paths()
	if err != nil {
		return nil, err
	}
	paths[p.errorPath] = &errorV3GV{err: p.err}
	return paths, nil
}

func TestGetResourceSchemasForClientV3_LookupKnownGVK(t *testing.T) {
	client := &fakeV3Client{paths: map[string][]byte{
		"apis/apps/v1": []byte(appsV1Schema),
	}}

	resources, err := GetResourceSchemasForClientV3(client)
	require.NoError(t, err)

	deploymentGVK := schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	s := resources.LookupResource(deploymentGVK)
	assert.NotNil(t, s, "LookupResource should return a non-nil schema for apps/v1/Deployment")
}

func TestGetResourceSchemasForClientV3_NilForUnknownGVK(t *testing.T) {
	client := &fakeV3Client{paths: map[string][]byte{
		"apis/apps/v1": []byte(appsV1Schema),
	}}

	resources, err := GetResourceSchemasForClientV3(client)
	require.NoError(t, err)

	unknownGVK := schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"}
	s := resources.LookupResource(unknownGVK)
	assert.Nil(t, s, "LookupResource should return nil for a GVK not in the schema")
}

func TestGetResourceSchemasForClientV3_MultipleGVs(t *testing.T) {
	client := &fakeV3Client{paths: map[string][]byte{
		"apis/apps/v1": []byte(appsV1Schema),
		"api/v1":       []byte(coreV1Schema),
	}}

	resources, err := GetResourceSchemasForClientV3(client)
	require.NoError(t, err)

	assert.NotNil(t, resources.LookupResource(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}))
	assert.NotNil(t, resources.LookupResource(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ConfigMap"}))
}

func TestGetResourceSchemasForClientV3_PartialFailure(t *testing.T) {
	// An error on one path (e.g. a metrics endpoint) must not prevent schemas from
	// other paths being returned.
	client := &partialErrorV3Client{
		delegate: &fakeV3Client{paths: map[string][]byte{
			"apis/apps/v1": []byte(appsV1Schema),
		}},
		errorPath: "metrics",
		err:       fmt.Errorf("metrics endpoint unavailable"),
	}

	resources, err := GetResourceSchemasForClientV3(client)
	require.NoError(t, err)

	deploymentGVK := schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	assert.NotNil(t, resources.LookupResource(deploymentGVK),
		"Deployment should be found despite error on the metrics path")
}

// TestGetResourceSchemasV2V3Parity verifies that GVKs available via the v3 fixture
// are also present in the v2 embedded swagger.json, confirming the two code paths
// produce consistent results for the same underlying data.
func TestGetResourceSchemasV2V3Parity(t *testing.T) {
	// v3 resources: load from inline fixtures.
	v3Client := &fakeV3Client{paths: map[string][]byte{
		"apis/apps/v1": []byte(appsV1Schema),
		"api/v1":       []byte(coreV1Schema),
	}}
	v3Resources, err := GetResourceSchemasForClientV3(v3Client)
	require.NoError(t, err)

	// v2 resources: load from the embedded swagger.json in the fake package.
	v2Doc, err := fake.LoadOpenAPISchema()
	require.NoError(t, err)
	v2Resources, err := GetResourceSchemasForClient(&fakeV2Disco{doc: v2Doc})
	require.NoError(t, err)

	// Both paths should return non-nil for GVKs covered by the fixtures.
	gvks := []schema.GroupVersionKind{
		{Group: "apps", Version: "v1", Kind: "Deployment"},
		{Group: "", Version: "v1", Kind: "ConfigMap"},
	}
	for _, gvk := range gvks {
		assert.NotNilf(t, v3Resources.LookupResource(gvk), "v3: LookupResource(%v)", gvk)
		assert.NotNilf(t, v2Resources.LookupResource(gvk), "v2: LookupResource(%v)", gvk)
	}
}
