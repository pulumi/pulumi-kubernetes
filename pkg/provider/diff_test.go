// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package provider

import (
	"reflect"
	"testing"
)

type object map[string]interface{}
type list []interface{}

func TestFieldsChanged(t *testing.T) {
	tests := []struct {
		group    string
		version  string
		kind     string
		old      object
		new      object
		expected []string
	}{
		{
			group: "core", version: "v1", kind: "PersistentVolumeClaim",
			old:      object{"spec": object{}},
			new:      object{"spec": object{"accessModes": object{}}},
			expected: []string{".spec", ".spec.accessModes"},
		},
		{
			group: "core", version: "v1", kind: "Pod",
			old:      object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			new:      object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx:1.15-alpine"}}}},
			expected: []string{".spec.containers[*].image"},
		},
		{
			group: "", version: "v1", kind: "Pod",
			old:      object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			new:      object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx:1.15-alpine"}}}},
			expected: []string{".spec.containers[*].image"},
		},
	}

	for _, test := range tests {
		diff, err := matchingProperties(test.old, test.new,
			forceNew[test.group][test.version][test.kind])
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(diff, test.expected) {
			t.Errorf("Got '%v' expected '%v'", diff, test.expected)
		}
	}
}
