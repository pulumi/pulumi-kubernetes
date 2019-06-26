// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package provider

import (
	"encoding/json"
	"reflect"
	"testing"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type object map[string]interface{}
type list []interface{}

func TestForceNewProperties(t *testing.T) {
	tests := []struct {
		name     string
		group    string
		version  string
		kind     string
		old      object
		new      object
		expected []string
	}{
		{
			name:  "Adding spec and nested field results in correct diffs.",
			group: "core", version: "v1", kind: "PersistentVolumeClaim",
			old:      object{"spec": object{}},
			new:      object{"spec": object{"accessModes": object{}}},
			expected: []string{"spec", "spec.accessModes"},
		},
		{
			name:  "Changing image spec results in correct diff.",
			group: "core", version: "v1", kind: "Pod",
			old:      object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			new:      object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx:1.15-alpine"}}}},
			expected: []string{"spec.containers[0].image"},
		},
		{
			name:  "Group unspecified and changing image spec results in correct diff.",
			group: "", version: "v1", kind: "Pod",
			old:      object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			new:      object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx:1.15-alpine"}}}},
			expected: []string{"spec.containers[0].image"},
		},
		{
			name:  `Changing image spec results in correct diff.`,
			group: "", version: "v1", kind: "Pod",
			old: object{
				"metadata": object{"namespace": "default"},
				"spec":     object{"containers": list{object{"name": "nginx", "image": "nginx"}}},
			},
			new: object{
				"metadata": object{"namespace": "default"},
				"spec":     object{"containers": list{object{"name": "nginx", "image": "nginx:1.15-alpine"}}},
			},
			expected: []string{"spec.containers[0].image"},
		},
		{
			name:  `Changing one image spec results in correct diff.`,
			group: "", version: "v1", kind: "Pod",
			old: object{
				"metadata": object{"namespace": "default"},
				"spec":     object{"containers": list{object{"name": "nginx", "image": "nginx"}, object{"name": "nginx", "image": "nginx"}}},
			},
			new: object{
				"metadata": object{"namespace": "default"},
				"spec":     object{"containers": list{object{"name": "nginx", "image": "nginx"}, object{"name": "nginx", "image": "nginx:1.15-alpine"}}},
			},
			expected: []string{"spec.containers[1].image"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldJSON, err := json.Marshal(tt.old)
			assert.NoError(t, err)

			newJSON, err := json.Marshal(tt.new)
			assert.NoError(t, err)

			patchBytes, err := jsonpatch.CreateMergePatch(oldJSON, newJSON)
			assert.NoError(t, err)

			patch := map[string]interface{}{}
			err = json.Unmarshal(patchBytes, &patch)
			assert.NoError(t, err)

			gvk := schema.GroupVersionKind{
				Group:   tt.group,
				Version: tt.version,
				Kind:    tt.kind,
			}
			diff, err := forceNewProperties(patch, tt.old, gvk)

			if err != nil {
				t.Errorf("forceNewProperties() error = %v, wantErr %v", err, nil)
			}

			if !reflect.DeepEqual(diff, tt.expected) {
				t.Errorf("forceNewProperties() = %v, want %v", diff, tt.expected)
			}
		})
	}
}
