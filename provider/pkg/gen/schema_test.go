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
)

func TestListInputsSpec(t *testing.T) {
	cases := []struct {
		name             string
		kind             string
		expectNamespace  bool
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
