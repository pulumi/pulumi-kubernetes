// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package provider

import (
	"encoding/json"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type (
	object   = map[string]any
	list     = []any
	expected = map[string]*pulumirpc.PropertyDiff
)

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
		name              string
		group             string
		version           string
		kind              string
		old               object
		new               object
		inputs            object
		oldInputs         object
		expected          expected
		customizeProvider func(provider *kubeProvider)
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
			name:  `ConfigMap resources don't trigger a replace when mutable.`,
			group: "core", version: "v1", kind: "ConfigMap",
			old: object{"data": object{"property1": "3"}},
			new: object{"data": object{"property1": "4"}},
			customizeProvider: func(p *kubeProvider) {
				p.enableConfigMapMutable = true
			},
			expected: expected{
				"data.property1": U,
			},
		},
		{
			name:  `ConfigMap resources trigger a replace when enableConfigMapMutable is not set.`,
			group: "core", version: "v1", kind: "ConfigMap",
			old: object{"data": object{"property1": "3"}},
			new: object{"data": object{"property1": "4"}},
			expected: expected{
				"data.property1": UR,
			},
		},
		{
			name:  `ConfigMap resources trigger a replace when marked as immutable even when enableConfigMapMutable is set.`,
			group: "core", version: "v1", kind: "ConfigMap",
			old: object{"data": object{"property1": "3"}, "immutable": true},
			new: object{"data": object{"property1": "4"}, "immutable": true},
			customizeProvider: func(p *kubeProvider) {
				p.enableConfigMapMutable = true
			},
			expected: expected{
				"data.property1": UR,
			},
		},
		{
			name:  `Secret resources don't trigger a replace when mutable.`,
			group: "core", version: "v1", kind: "Secret",
			old: object{"data": object{"property1": "3"}},
			new: object{"data": object{"property1": "4"}},
			customizeProvider: func(p *kubeProvider) {
				p.enableSecretMutable = true
			},
			expected: expected{
				"data.property1": U,
			},
		},
		{
			name:  `Secret resources trigger a replace when enableSecretMutable is not set.`,
			group: "core", version: "v1", kind: "Secret",
			old: object{"data": object{"property1": "3"}},
			new: object{"data": object{"property1": "4"}},
			expected: expected{
				"data.property1": UR,
			},
		},
		{
			name:  `Secret resources trigger a replace when type changes even if enableSecretMutable is set.`,
			group: "core", version: "v1", kind: "Secret",
			old: object{"type": "kubernetes.io/dockerconfigjson", "data": object{"property1": "3"}},
			new: object{"type": "Opaque", "data": object{"property1": "3"}},
			customizeProvider: func(p *kubeProvider) {
				p.enableSecretMutable = true
			},
			expected: expected{
				"type": UR,
			},
		},
		{
			name:  `Secret resources trigger a replace when marked as immutable even if enableSecretMutable is set.`,
			group: "core", version: "v1", kind: "Secret",
			old: object{"data": object{"property1": "3"}, "immutable": true},
			new: object{"data": object{"property1": "4"}, "immutable": true},
			customizeProvider: func(p *kubeProvider) {
				p.enableSecretMutable = true
			},
			expected: expected{
				"data.property1": UR,
			},
		},
		{
			name:  `Changing computed object values results in correct diff`,
			group: "core", version: "v1", kind: "Pod",
			old:    object{"spec": object{"containers": list{object{"name": "nginx", "image": "nginx"}}}},
			new:    object{"spec": object{"containers": list{object{"name": "nginx", "image": nil}}}},
			inputs: object{"spec": object{"containers": list{object{"name": "nginx", "image": resource.Computed{}}}}},
			expected: expected{
				"spec.containers[0].image": UR,
			},
		},
		{
			name:  `Adding computed object values results in correct diff`,
			group: "core", version: "v1", kind: "Pod",
			old:    object{"spec": object{"containers": list{object{"name": "nginx"}}}},
			new:    object{"spec": object{"containers": list{object{"name": "nginx", "image": nil}}}},
			inputs: object{"spec": object{"containers": list{object{"name": "nginx", "image": resource.Computed{}}}}},
			expected: expected{
				"spec.containers[0].image": AR,
			},
		},
		{
			name:  `Adding computed array values results in correct diff`,
			group: "core", version: "v1", kind: "Pod",
			old:    object{"spec": object{"containers": list{}}},
			new:    object{"spec": object{"containers": list{nil}}},
			inputs: object{"spec": object{"containers": list{resource.Computed{}}}},
			expected: expected{
				"spec.containers[0]": A,
			},
		},
		{
			name:  `Changing computed array values results in correct diff`,
			group: "core", version: "v1", kind: "Pod",
			old:    object{"spec": object{"containers": list{object{"name": "nginx"}}}},
			new:    object{"spec": object{"containers": list{nil}}},
			inputs: object{"spec": object{"containers": list{resource.Computed{}}}},
			expected: expected{
				"spec.containers[0]": U,
			},
		},
		{
			name:  `Removing array values results in correct diff`,
			group: "core", version: "v1", kind: "Pod",
			old:    object{"spec": object{"containers": list{object{"name": "nginx"}}}},
			new:    object{"spec": object{"containers": list{}}},
			inputs: object{"spec": object{"containers": list{}}},
			expected: expected{
				"spec.containers[0]": D,
			},
		},
		{
			name:  "Removing a field that was nil should not panic (#1970).",
			group: "tekton.dev", version: "v1beta1", kind: "Pipeline",
			old: object{"taskSpec": object{"spec": nil, "steps": list{object{"name": "something"}}}},
			new: object{"taskSpec": object{"steps": list{object{"name": "something-else"}}}},
			expected: expected{
				"taskSpec.spec":          D,
				"taskSpec.steps[0].name": U,
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

			patch := map[string]any{}
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
			k := &kubeProvider{}
			if tt.customizeProvider != nil {
				tt.customizeProvider(k)
			}
			diff, err := convertPatchToDiff(patch, tt.old, inputs, oldInputs, k.forceNewProperties(obj)...)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, diff)
		})
	}
}
