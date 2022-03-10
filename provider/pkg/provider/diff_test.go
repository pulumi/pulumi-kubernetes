// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package provider

import (
	"encoding/json"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"

	jsonpatch "github.com/evanphx/json-patch"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
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
		name      string
		group     string
		version   string
		kind      string
		old       object
		new       object
		inputs    object
		oldInputs object
		expected  expected
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
		{
			name:  `State diffs with no corresponding input property are ignored.`,
			group: "core", version: "v1", kind: "Pod",
			old:       object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}, "status": object{"hostIP": "10.0.0.2"}},
			new:       object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}, "status": object{"hostIP": "10.0.0.3"}},
			inputs:    object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			oldInputs: object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			expected:  expected{},
		},
		{
			name:  `State deletes with no corresponding input properties are ignored.`,
			group: "core", version: "v1", kind: "Pod",
			old:       object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}, "status": object{"hostIP": "10.0.0.2"}},
			new:       object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			inputs:    object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			oldInputs: object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			expected:  expected{},
		},
		{
			name:  `PVC resources don't trigger a replace.`,
			group: "core", version: "v1", kind: "PersistentVolumeClaim",
			old: object{"spec": object{"resources": object{"requests": object{"storage": "10Gi"}}}},
			new: object{"spec": object{"resources": object{"requests": object{"storage": "20Gi"}}}},
			expected: expected{
				"spec.resources.requests.storage": U,
			},
		},
		{
			name:  `ConfigMap resources don't trigger a replace.`,
			group: "core", version: "v1", kind: "ConfigMap",
			old: object{"data": object{"property1": "3"}},
			new: object{"data": object{"property1": "4"}},
			expected: expected{
				"data.property1": U,
			},
		},
		{
			name:  `Immutable ConfigMap resources trigger a replace.`,
			group: "core", version: "v1", kind: "ConfigMap",
			old: object{"data": object{"property1": "3"}, "immutable": true},
			new: object{"data": object{"property1": "4"}, "immutable": true},
			expected: expected{
				"data.property1": UR,
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

			inputs := tt.inputs
			if inputs == nil {
				inputs = tt.new
			}
			oldInputs := tt.oldInputs
			if oldInputs == nil {
				oldInputs = tt.old
			}

			gvk := schema.GroupVersionKind{
				Group:   tt.group,
				Version: tt.version,
				Kind:    tt.kind,
			}
			obj := &unstructured.Unstructured{}
			obj.SetUnstructuredContent(tt.old)
			obj.SetGroupVersionKind(gvk)
			diff, err := convertPatchToDiff(patch, tt.old, inputs, oldInputs, forceNewProperties(obj)...)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, diff)
		})
	}
}
