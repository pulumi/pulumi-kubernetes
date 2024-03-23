// Copyright 2016-2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clients

import (
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubeversion "k8s.io/apimachinery/pkg/version"
	discoveryfake "k8s.io/client-go/discovery/fake"
	kubetesting "k8s.io/client-go/testing"
)

var fakeCRDs = []unstructured.Unstructured{{Object: map[string]interface{}{
	"apiVersion": "apiextensions.k8s.io/v1",
	"kind":       "CustomResourceDefinition",
	"metadata": map[string]interface{}{
		"name": "crontabs.stable.example.com",
	},
	"spec": map[string]interface{}{
		"group": "stable.example.com",
		"versions": []interface{}{
			map[string]interface{}{
				"name":    "v1",
				"served":  true,
				"storage": true,
				"schema": map[string]interface{}{
					"openAPIV3Schema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"spec": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"cronSpec": map[string]interface{}{
										"type": "string",
									},
									"image": map[string]interface{}{
										"type": "string",
									},
									"replicas": map[string]interface{}{
										"type": "integer",
									},
								},
							},
						},
					},
				},
			},
		},
		"scope": "Namespaced",
		"names": map[string]interface{}{
			"plural":   "crontabs",
			"singular": "crontab",
			"kind":     "CronTab",
		},
	},
}}}

var fakeResources = []*metav1.APIResourceList{
	{
		GroupVersion: "postgresql.example.com/v1alpha1",
		APIResources: []metav1.APIResource{
			{Name: "roles", Namespaced: false, Kind: "Role"},
		},
	},
}

func TestIsNamespacedKind(t *testing.T) {
	// coverage: in-built kinds, discoverable kinds, cached kinds, and kinds based on the supplied CRDs.
	tests := []struct {
		gvk     schema.GroupVersionKind
		cached  []unstructured.Unstructured
		objs    []unstructured.Unstructured
		want    bool
		wantErr bool
	}{
		{schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"}, nil, nil, true, false},
		{schema.GroupVersionKind{Group: "core", Version: "v1", Kind: "Pod"}, nil, nil, true, false},
		{schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "Role"}, nil, nil, true, false},
		{schema.GroupVersionKind{Group: "postgresql.example.com", Version: "v1alpha1", Kind: "Role"}, nil, nil, false, false},
		{schema.GroupVersionKind{Group: "postgresql.example.com", Version: "v1alpha1", Kind: "Missing"}, nil, nil, false, true},
		{schema.GroupVersionKind{Group: "stable.example.com", Version: "v1", Kind: "CronTab"}, fakeCRDs, nil, true, false},
		{schema.GroupVersionKind{Group: "stable.example.com", Version: "v1", Kind: "Missing"}, fakeCRDs, nil, false, true},
		{schema.GroupVersionKind{Group: "stable.example.com", Version: "v1", Kind: "CronTab"}, nil, fakeCRDs, true, false},
		{schema.GroupVersionKind{Group: "stable.example.com", Version: "v1", Kind: "Missing"}, nil, fakeCRDs, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.gvk.String(), func(t *testing.T) {
			version := kubeversion.Info{Major: "1", Minor: "29"}
			clientSet := &DynamicClientSet{
				DiscoveryClientCached: &simpleDiscovery{
					FakeDiscovery: discoveryfake.FakeDiscovery{
						Fake: &kubetesting.Fake{
							Resources: fakeResources,
						},
						FakedServerVersion: &version,
					},
				},
				CRDCache: &CRDCache{},
			}
			for _, crd := range tt.cached {
				_ = clientSet.CRDCache.AddCRD(&crd)
			}
			got, err := IsNamespacedKind(tt.gvk, clientSet, tt.objs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsNamespacedKind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsNamespacedKind() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type simpleDiscovery struct {
	discoveryfake.FakeDiscovery
}

func (d *simpleDiscovery) Fresh() bool {
	return true
}
func (d *simpleDiscovery) Invalidate() {}
