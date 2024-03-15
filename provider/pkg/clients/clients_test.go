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
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubeversion "k8s.io/apimachinery/pkg/version"
	discoveryfake "k8s.io/client-go/discovery/fake"
	kubetesting "k8s.io/client-go/testing"
)

var fakeResources = []*metav1.APIResourceList{
	{
		GroupVersion: "stable.example.com/v1",
		APIResources: []metav1.APIResource{
			{Name: "crontabs", Namespaced: true, Kind: "CronTab"},
			{Name: "roles", Namespaced: false, Kind: "Role"},
		},
	},
}

func TestIsNamespacedKind(t *testing.T) {
	tests := []struct {
		gvk     schema.GroupVersionKind
		want    bool
		wantErr bool
	}{
		{schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"}, true, false},
		{schema.GroupVersionKind{Group: "core", Version: "v1", Kind: "Pod"}, true, false},
		{schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "Role"}, true, false},
		{schema.GroupVersionKind{Group: "stable.example.com", Version: "v1", Kind: "CronTab"}, true, false},
		{schema.GroupVersionKind{Group: "stable.example.com", Version: "v1", Kind: "Role"}, false, false},
		{schema.GroupVersionKind{Group: "stable.example.com", Version: "v1", Kind: "Missing"}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.gvk.String(), func(t *testing.T) {
			version := kubeversion.Info{Major: "1", Minor: "29"}
			disco := &discoveryfake.FakeDiscovery{
				Fake: &kubetesting.Fake{
					Resources: fakeResources,
				},
				FakedServerVersion: &version,
			}

			got, err := IsNamespacedKind(tt.gvk, disco)
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
