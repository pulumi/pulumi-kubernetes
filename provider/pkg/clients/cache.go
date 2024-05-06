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

package clients

import (
	"fmt"
	"sync"

	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// CRDCache facilitates program-wide resolution of Kubernetes kinds.
// For example, it allows the yaml/v2 package to lookup a CRD during preview that would be installed by
// a "CustomResourceDefinition" resource. The user is expected to use DependsOn to ensure that
// the CRD is created first.
type CRDCache struct {
	mu   sync.Mutex
	crds map[schema.GroupKind]*unstructured.Unstructured
}

// GetCRD returns the CRD for the given kind, if it exists in the cache.
func (c *CRDCache) GetCRD(kind schema.GroupKind) *unstructured.Unstructured {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.crds[kind]
}

// AddCRD adds a CRD to the cache.
func (c *CRDCache) AddCRD(crd *unstructured.Unstructured) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	kind, err := groupKindFromCRD(crd)
	if err != nil {
		return err
	}
	if c.crds == nil {
		c.crds = make(map[schema.GroupKind]*unstructured.Unstructured)
	}
	c.crds[kind] = crd
	return nil
}

// RemoveCRD removes a CRD from the cache.
func (c *CRDCache) RemoveCRD(crd *unstructured.Unstructured) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	kind, err := groupKindFromCRD(crd)
	if err != nil {
		return err
	}
	delete(c.crds, kind)
	return nil
}

func groupKindFromCRD(crd *unstructured.Unstructured) (schema.GroupKind, error) {
	contract.Requiref(IsCRD(crd), "crd", "expected a CRD")
	crdGroup, _, err := unstructured.NestedString(crd.Object, "spec", "group")
	if err != nil {
		return schema.GroupKind{}, err
	}
	if crdGroup == "" {
		return schema.GroupKind{}, fmt.Errorf("expected .spec.group")
	}
	crdKind, _, err := unstructured.NestedString(crd.Object, "spec", "names", "kind")
	if err != nil {
		return schema.GroupKind{}, err
	}
	if crdKind == "" {
		return schema.GroupKind{}, fmt.Errorf("expected .spec.names.kind")
	}
	kind := schema.GroupKind{Group: crdGroup, Kind: crdKind}
	return kind, nil
}
