// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package provider

import (
	"reflect"
	"testing"

	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
)

type object map[string]interface{}
type list []interface{}

func TestPropertiesChanged(t *testing.T) {
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
			expected: []string{".spec", ".spec.accessModes"},
		},
		{
			name:  "Changing image spec results in correct diff.",
			group: "core", version: "v1", kind: "Pod",
			old:      object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			new:      object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx:1.15-alpine"}}}},
			expected: []string{".spec.containers[*].image"},
		},
		{
			name:  "Group unspecified and changing image spec results in correct diff.",
			group: "", version: "v1", kind: "Pod",
			old:      object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			new:      object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx:1.15-alpine"}}}},
			expected: []string{".spec.containers[*].image"},
		},
		{
			name:  `Changing namespace from "" to "default" produces no diff.`,
			group: "core", version: "v1", kind: "Pod",
			old:      object{"metadata": object{"namespace": ""}},
			new:      object{"metadata": object{"namespace": "default"}},
			expected: []string{},
		},
		{
			name:  `Changing image spec results in correct diff and changing namespace from "" to "default" produces no diff.`,
			group: "", version: "v1", kind: "Pod",
			old: object{
				"metadata": object{"namespace": ""},
				"spec":     object{"containers": list{object{"name": "nginx", "image": "nginx"}}},
			},
			new: object{
				"metadata": object{"namespace": "default"},
				"spec":     object{"containers": list{object{"name": "nginx", "image": "nginx:1.15-alpine"}}},
			},
			expected: []string{".spec.containers[*].image"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff, err := openapi.PropertiesChanged(tt.old, tt.new,
				forceNew[tt.group][tt.version][tt.kind])
			if err != nil {
				t.Errorf("PropertiesChanged() error = %v, wantErr %v", err, nil)
			}

			if !reflect.DeepEqual(diff, tt.expected) {
				t.Errorf("PropertiesChanged() = %v, want %v", diff, tt.expected)
			}
		})
	}
}
