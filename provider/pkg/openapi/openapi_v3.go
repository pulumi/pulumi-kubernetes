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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clientopenapi "k8s.io/client-go/openapi"
	"k8s.io/kube-openapi/pkg/util/proto"
	k8sopenapi "k8s.io/kubectl/pkg/util/openapi"

	openapi_v3 "github.com/google/gnostic-models/openapiv3"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
)

// GetResourceSchemasForClientV3 fetches resource schemas using the OpenAPI v3 per-GV
// endpoint. Unlike GetResourceSchemasForClient (v2), schemas are fetched per
// group-version and benefit from the per-GV caching in cachedopenapi.NewClient, so
// repeated calls for the same GV avoid re-downloading the schema.
func GetResourceSchemasForClientV3(client clientopenapi.Client) (k8sopenapi.Resources, error) {
	paths, err := client.Paths()
	if err != nil {
		return nil, err
	}

	gvkToModel := map[schema.GroupVersionKind]string{}
	var allModels []proto.Models

	for path, gv := range paths {
		schemaBytes, err := gv.Schema(runtime.ContentTypeJSON)
		if err != nil {
			// Some paths (e.g. metrics endpoints) don't serve OpenAPI schemas; skip them.
			logger.V(9).Infof("skipping OpenAPI v3 path %q: %v", path, err)
			continue
		}

		doc, err := openapi_v3.ParseDocument(schemaBytes)
		if err != nil {
			logger.V(9).Infof("skipping OpenAPI v3 path %q: failed to parse document: %v", path, err)
			continue
		}

		models, err := proto.NewOpenAPIV3Data(doc)
		if err != nil {
			logger.V(9).Infof("skipping OpenAPI v3 path %q: failed to build models: %v", path, err)
			continue
		}

		for _, name := range models.ListModels() {
			s := models.LookupModel(name)
			if s == nil {
				continue
			}
			for _, gvk := range gvksFromSchema(s) {
				gvkToModel[gvk] = name
			}
		}

		allModels = append(allModels, models)
	}

	return &v3document{
		resources: gvkToModel,
		models:    &mergedModels{models: allModels},
	}, nil
}

// v3document implements k8sopenapi.Resources backed by merged per-GV proto.Models.
type v3document struct {
	resources map[schema.GroupVersionKind]string
	models    proto.Models
}

func (d *v3document) LookupResource(gvk schema.GroupVersionKind) proto.Schema {
	name, ok := d.resources[gvk]
	if !ok {
		return nil
	}
	return d.models.LookupModel(name)
}

// GetConsumes returns the accepted content types for the given GVK and operation.
// The v3 schema does not encode consumes; returning nil matches the behaviour of
// LookupResource for an unknown GVK in the v2 implementation.
func (d *v3document) GetConsumes(_ schema.GroupVersionKind, _ string) []string {
	return nil
}

// mergedModels implements proto.Models by searching across multiple per-GV Models.
type mergedModels struct {
	models []proto.Models
}

func (m *mergedModels) LookupModel(name string) proto.Schema {
	for _, mm := range m.models {
		if s := mm.LookupModel(name); s != nil {
			return s
		}
	}
	return nil
}

func (m *mergedModels) ListModels() []string {
	seen := map[string]bool{}
	var out []string
	for _, mm := range m.models {
		for _, name := range mm.ListModels() {
			if !seen[name] {
				seen[name] = true
				out = append(out, name)
			}
		}
	}
	return out
}

// gvksFromSchema extracts Kubernetes GroupVersionKinds from a proto.Schema by
// reading the x-kubernetes-group-version-kind extension.
func gvksFromSchema(s proto.Schema) []schema.GroupVersionKind {
	exts := s.GetExtensions()
	val, ok := exts["x-kubernetes-group-version-kind"]
	if !ok {
		return nil
	}

	list, ok := val.([]interface{})
	if !ok {
		return nil
	}

	var result []schema.GroupVersionKind
	for _, item := range list {
		group, version, kind := extractGVK(item)
		if version == "" || kind == "" {
			continue
		}
		result = append(result, schema.GroupVersionKind{
			Group:   group,
			Version: version,
			Kind:    kind,
		})
	}
	return result
}

// extractGVK extracts group/version/kind strings from a YAML-unmarshaled map,
// handling both yaml.v2 (map[interface{}]interface{}) and yaml.v3 (map[string]interface{}) formats.
func extractGVK(item interface{}) (group, version, kind string) {
	switch m := item.(type) {
	case map[string]interface{}:
		group, _ = m["group"].(string)
		version, _ = m["version"].(string)
		kind, _ = m["kind"].(string)
	case map[interface{}]interface{}:
		group, _ = m["group"].(string)
		version, _ = m["version"].(string)
		kind, _ = m["kind"].(string)
	}
	return
}
