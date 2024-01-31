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

package fake

import (
	"embed"
	_ "embed"
	"encoding/json"
	"io"
	"io/fs"

	openapi_v2 "github.com/google/gnostic-models/openapiv2"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubeversion "k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	discoveryfake "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/restmapper"
	kubetesting "k8s.io/client-go/testing"
)

//go:embed swagger.json
var swagger []byte

func LoadOpenAPISchema() (*openapi_v2.Document, error) {
	return openapi_v2.ParseDocument(swagger)
}

//go:embed serverresources
var serverresources embed.FS

func loadServerResources() ([]*metav1.APIResourceList, error) {
	all := []*metav1.APIResourceList{}
	err := fs.WalkDir(serverresources, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && d.Name() == "serverresources.json" {
			f, err := serverresources.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			data, err := io.ReadAll(f)
			if err != nil {
				return err
			}
			serverResources := &metav1.APIResourceList{}
			if err := json.Unmarshal(data, &serverResources); err != nil {
				return err
			}
			all = append(all, serverResources)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return all, nil
}

// SimpleDiscovery provides a fake discovery client with core Kubernetes types.
type SimpleDiscovery struct {
	*discoveryfake.FakeDiscovery
}

var _ discovery.CachedDiscoveryInterface = &SimpleDiscovery{}

func (*SimpleDiscovery) Fresh() bool {
	return true
}

func (*SimpleDiscovery) Invalidate() {}

func (*SimpleDiscovery) OpenAPISchema() (*openapi_v2.Document, error) {
	return LoadOpenAPISchema()
}

func NewSimpleDiscovery(serverVersion kubeversion.Info) *SimpleDiscovery {
	// make a fake discovery client based on an embedded openapi document and corresponding server metadata.
	serverResources, err := loadServerResources()
	if err != nil {
		panic(err)
	}
	return &SimpleDiscovery{
		FakeDiscovery: &discoveryfake.FakeDiscovery{
			Fake: &kubetesting.Fake{
				Resources: serverResources,
			},
			FakedServerVersion: &serverVersion,
		},
	}
}

type SimpleRESTMapper struct {
	meta.RESTMapper
}

var _ meta.ResettableRESTMapper = &SimpleRESTMapper{}

func (m *SimpleRESTMapper) Reset() {}

// NewSimpleRESTMapper creates a simple REST mapper for testing purposes, backed by the provided resource mappings.
//
// Note: Generate resource mappings by calling discovery.GetAPIGroupResources with a working discovery client.
func NewSimpleRESTMapper(resources []*restmapper.APIGroupResources) *SimpleRESTMapper {
	return &SimpleRESTMapper{RESTMapper: restmapper.NewDiscoveryRESTMapper(resources)}
}

type StubResettableRESTMapper struct {
	meta.ResettableRESTMapper
	ResetF                func()
	KindForF              func(resource schema.GroupVersionResource) (schema.GroupVersionKind, error)
	KindsForF             func(resource schema.GroupVersionResource) ([]schema.GroupVersionKind, error)
	ResourceForF          func(input schema.GroupVersionResource) (schema.GroupVersionResource, error)
	ResourcesForF         func(input schema.GroupVersionResource) ([]schema.GroupVersionResource, error)
	RESTMappingF          func(gk schema.GroupKind, versions ...string) (*meta.RESTMapping, error)
	RESTMappingsF         func(gk schema.GroupKind, versions ...string) ([]*meta.RESTMapping, error)
	ResourceSingularizerF func(resource string) (singular string, err error)
}

var _ meta.ResettableRESTMapper = &StubResettableRESTMapper{}

func (m *StubResettableRESTMapper) Reset() {
	if m.ResetF != nil {
		m.ResetF()
		return
	}
	m.ResettableRESTMapper.Reset()
}

var _ meta.RESTMapper = &StubResettableRESTMapper{}

func (m *StubResettableRESTMapper) KindFor(resource schema.GroupVersionResource) (schema.GroupVersionKind, error) {
	if m.KindsForF != nil {
		return m.KindForF(resource)
	}
	return m.ResettableRESTMapper.KindFor(resource)
}

func (m *StubResettableRESTMapper) KindsFor(resource schema.GroupVersionResource) ([]schema.GroupVersionKind, error) {
	if m.KindsForF != nil {
		return m.KindsForF(resource)
	}
	return m.ResettableRESTMapper.KindsFor(resource)
}

func (m *StubResettableRESTMapper) ResourceFor(input schema.GroupVersionResource) (schema.GroupVersionResource, error) {
	if m.ResourceForF != nil {
		return m.ResourceForF(input)
	}
	return m.ResettableRESTMapper.ResourceFor(input)
}

func (m *StubResettableRESTMapper) ResourcesFor(input schema.GroupVersionResource) ([]schema.GroupVersionResource, error) {
	if m.ResourcesForF != nil {
		return m.ResourcesForF(input)
	}
	return m.ResettableRESTMapper.ResourcesFor(input)
}

func (m *StubResettableRESTMapper) RESTMapping(gk schema.GroupKind, versions ...string) (*meta.RESTMapping, error) {
	if m.RESTMappingF != nil {
		return m.RESTMappingF(gk, versions...)
	}
	return m.ResettableRESTMapper.RESTMapping(gk, versions...)
}

func (m *StubResettableRESTMapper) RESTMappings(gk schema.GroupKind, versions ...string) ([]*meta.RESTMapping, error) {
	if m.RESTMappingsF != nil {
		return m.RESTMappingsF(gk, versions...)
	}
	return m.ResettableRESTMapper.RESTMappings(gk, versions...)
}

func (m *StubResettableRESTMapper) ResourceSingularizer(resource string) (singular string, err error) {
	if m.ResourceSingularizerF != nil {
		return m.ResourceSingularizerF(resource)
	}
	return m.ResettableRESTMapper.ResourceSingularizer(resource)
}
