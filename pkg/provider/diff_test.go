// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package provider

import (
	"encoding/json"
	"testing"

	jsonpatch "github.com/evanphx/json-patch"
	pulumirpc "github.com/pulumi/pulumi/sdk/proto/go"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type object = map[string]interface{}
type list = []interface{}
type expected = map[string]*pulumirpc.PropertyDiff

func TestPatchToDiff(t *testing.T) {
	var (
		A  = &pulumirpc.PropertyDiff{Kind: pulumirpc.PropertyDiff_ADD}
		D  = &pulumirpc.PropertyDiff{Kind: pulumirpc.PropertyDiff_DELETE, InputDiff: true}
		U  = &pulumirpc.PropertyDiff{Kind: pulumirpc.PropertyDiff_UPDATE}
		AR = &pulumirpc.PropertyDiff{Kind: pulumirpc.PropertyDiff_ADD_REPLACE}
		DR = &pulumirpc.PropertyDiff{Kind: pulumirpc.PropertyDiff_DELETE_REPLACE, InputDiff: true}
		UR = &pulumirpc.PropertyDiff{Kind: pulumirpc.PropertyDiff_UPDATE_REPLACE}
	)

	tests := []struct {
		name     string
		group    string
		version  string
		kind     string
		old      object
		new      object
		expected expected
	}{
		{
			name:  "Adding spec and nested field results in correct diffs.",
			group: "core", version: "v1", kind: "PersistentVolumeClaim",
			old: object{"spec": object{}},
			new: object{"spec": object{"accessModes": object{}}},
			expected: expected{
				"spec.accessModes": AR,
			},
		},
		{
			name:  "Deleting spec and nested field results in correct diffs.",
			group: "core", version: "v1", kind: "PersistentVolumeClaim",
			old: object{"spec": object{"accessModes": object{}}},
			new: object{"spec": object{}},
			expected: expected{
				"spec.accessModes": DR,
			},
		},
		{
			name:  "Changing image spec results in correct diff.",
			group: "core", version: "v1", kind: "Pod",
			old: object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			new: object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx:1.15-alpine"}}}},
			expected: expected{
				"spec.containers[0].image": UR,
			},
		},
		{
			name:  "Group unspecified and changing image spec results in correct diff.",
			group: "", version: "v1", kind: "Pod",
			old: object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			new: object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx:1.15-alpine"}}}},
			expected: expected{
				"spec.containers[0].image": UR,
			},
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
			expected: expected{
				"spec.containers[0].image": UR,
			},
		},
		{
			name:  `Changing second image spec results in correct diff.`,
			group: "", version: "v1", kind: "Pod",
			old: object{
				"metadata": object{"namespace": "default"},
				"spec":     object{"containers": list{object{"name": "nginx", "image": "nginx"}, object{"name": "nginx", "image": "nginx"}}},
			},
			new: object{
				"metadata": object{"namespace": "default"},
				"spec":     object{"containers": list{object{"name": "nginx", "image": "nginx"}, object{"name": "nginx", "image": "nginx:1.15-alpine"}}},
			},
			expected: expected{
				"spec.containers[1].image": UR,
			},
		},
		{
			name:  `Changing DNS policy and image spec results in correct diff.`,
			group: "core", version: "v1", kind: "Pod",
			old: object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx", "dnsPolicy": "Default"}}}},
			new: object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx:1.15-alpine", "dnsPolicy": "None"}}}},
			expected: expected{
				"spec.containers[0].dnsPolicy": U,
				"spec.containers[0].image":     UR,
			},
		},
		{
			name:  `Adding DNS policy results in correct diff.`,
			group: "core", version: "v1", kind: "Pod",
			old: object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			new: object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx", "dnsPolicy": "Default"}}}},
			expected: expected{
				"spec.containers[0].dnsPolicy": A,
			},
		},
		{
			name:  `Changing DNS policy results in correct diff.`,
			group: "core", version: "v1", kind: "Pod",
			old: object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx", "dnsPolicy": "Default"}}}},
			new: object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx", "dnsPolicy": "None"}}}},
			expected: expected{
				"spec.containers[0].dnsPolicy": U,
			},
		},
		{
			name:  `Deleting DNS policy results in correct diff.`,
			group: "core", version: "v1", kind: "Pod",
			old: object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx", "dnsPolicy": "Default"}}}},
			new: object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			expected: expected{
				"spec.containers[0].dnsPolicy": D,
			},
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
			diff, err := convertPatchToDiff(patch, tt.old, gvk)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, diff)
		})
	}
}
