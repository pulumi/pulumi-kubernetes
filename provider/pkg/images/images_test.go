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

package images

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestFromObjects_Deployment(t *testing.T) {
	objs := []unstructured.Unstructured{
		{Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata":   map[string]any{"name": "test", "namespace": "default"},
			"spec": map[string]any{
				"template": map[string]any{
					"spec": map[string]any{
						"initContainers": []any{
							map[string]any{"name": "init", "image": "busybox:1.36"},
						},
						"containers": []any{
							map[string]any{"name": "app", "image": "nginx:1.25"},
							map[string]any{"name": "sidecar", "image": "envoyproxy/envoy:v1.28"},
						},
					},
				},
			},
		}},
	}
	got := FromObjects(objs)
	assert.Equal(t, []string{"busybox:1.36", "envoyproxy/envoy:v1.28", "nginx:1.25"}, got)
}

func TestFromObjects_DaemonSet(t *testing.T) {
	objs := []unstructured.Unstructured{
		{Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "DaemonSet",
			"metadata":   map[string]any{"name": "cni"},
			"spec": map[string]any{
				"template": map[string]any{
					"spec": map[string]any{
						"containers": []any{
							map[string]any{"name": "cni", "image": "quay.io/cilium/cilium:v1.18.5"},
						},
						"initContainers": []any{
							map[string]any{"name": "install-cni", "image": "quay.io/cilium/cilium:v1.18.5"},
						},
					},
				},
			},
		}},
	}
	got := FromObjects(objs)
	// Dedup: same image in containers and initContainers.
	assert.Equal(t, []string{"quay.io/cilium/cilium:v1.18.5"}, got)
}

func TestFromObjects_CronJob(t *testing.T) {
	objs := []unstructured.Unstructured{
		{Object: map[string]any{
			"apiVersion": "batch/v1",
			"kind":       "CronJob",
			"metadata":   map[string]any{"name": "backup"},
			"spec": map[string]any{
				"jobTemplate": map[string]any{
					"spec": map[string]any{
						"template": map[string]any{
							"spec": map[string]any{
								"containers": []any{
									map[string]any{"name": "backup", "image": "postgres:16"},
								},
							},
						},
					},
				},
			},
		}},
	}
	got := FromObjects(objs)
	assert.Equal(t, []string{"postgres:16"}, got)
}

func TestFromObjects_Pod(t *testing.T) {
	objs := []unstructured.Unstructured{
		{Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata":   map[string]any{"name": "test"},
			"spec": map[string]any{
				"containers": []any{
					map[string]any{"name": "app", "image": "alpine:3.19"},
				},
				"ephemeralContainers": []any{
					map[string]any{"name": "debug", "image": "busybox:latest"},
				},
			},
		}},
	}
	got := FromObjects(objs)
	assert.Equal(t, []string{"alpine:3.19", "busybox:latest"}, got)
}

func TestFromObjects_ImageVolumeSource(t *testing.T) {
	objs := []unstructured.Unstructured{
		{Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata":   map[string]any{"name": "test"},
			"spec": map[string]any{
				"containers": []any{
					map[string]any{"name": "app", "image": "app:1.0"},
				},
				"volumes": []any{
					map[string]any{
						"name": "data",
						"image": map[string]any{
							"reference": "registry.example.com/data:v1",
						},
					},
				},
			},
		}},
	}
	got := FromObjects(objs)
	assert.Equal(t, []string{"app:1.0", "registry.example.com/data:v1"}, got)
}

func TestFromObjects_Dedup(t *testing.T) {
	obj := unstructured.Unstructured{Object: map[string]any{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata":   map[string]any{"name": "a"},
		"spec": map[string]any{
			"template": map[string]any{
				"spec": map[string]any{
					"containers": []any{
						map[string]any{"name": "a", "image": "nginx:1.25"},
					},
				},
			},
		},
	}}
	got := FromObjects([]unstructured.Unstructured{obj, obj})
	assert.Equal(t, []string{"nginx:1.25"}, got)
}

func TestFromObjects_NonWorkload(t *testing.T) {
	objs := []unstructured.Unstructured{
		{Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata":   map[string]any{"name": "test"},
			"data":       map[string]any{"key": "value"},
		}},
		{Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata":   map[string]any{"name": "test"},
		}},
		{Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Namespace",
			"metadata":   map[string]any{"name": "test"},
		}},
	}
	got := FromObjects(objs)
	assert.Empty(t, got)
}

func TestFromObjects_Empty(t *testing.T) {
	got := FromObjects(nil)
	assert.Empty(t, got)
	assert.NotNil(t, got) // Must be [] not nil.
}

func TestFromObjects_StatefulSet(t *testing.T) {
	objs := []unstructured.Unstructured{
		{Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "StatefulSet",
			"metadata":   map[string]any{"name": "db"},
			"spec": map[string]any{
				"template": map[string]any{
					"spec": map[string]any{
						"containers": []any{
							map[string]any{"name": "db", "image": "postgres:16-alpine"},
						},
					},
				},
			},
		}},
	}
	got := FromObjects(objs)
	assert.Equal(t, []string{"postgres:16-alpine"}, got)
}

func TestFromObjects_Job(t *testing.T) {
	objs := []unstructured.Unstructured{
		{Object: map[string]any{
			"apiVersion": "batch/v1",
			"kind":       "Job",
			"metadata":   map[string]any{"name": "migrate"},
			"spec": map[string]any{
				"template": map[string]any{
					"spec": map[string]any{
						"containers": []any{
							map[string]any{"name": "migrate", "image": "flyway/flyway:10"},
						},
					},
				},
			},
		}},
	}
	got := FromObjects(objs)
	assert.Equal(t, []string{"flyway/flyway:10"}, got)
}

func TestFromObjects_MixedWorkloadsAndNonWorkloads(t *testing.T) {
	objs := []unstructured.Unstructured{
		{Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata":   map[string]any{"name": "config"},
		}},
		{Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata":   map[string]any{"name": "app"},
			"spec": map[string]any{
				"template": map[string]any{
					"spec": map[string]any{
						"containers": []any{
							map[string]any{"name": "app", "image": "myapp:v2"},
						},
					},
				},
			},
		}},
		{Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Service",
			"metadata":   map[string]any{"name": "svc"},
		}},
		{Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "DaemonSet",
			"metadata":   map[string]any{"name": "agent"},
			"spec": map[string]any{
				"template": map[string]any{
					"spec": map[string]any{
						"containers": []any{
							map[string]any{"name": "agent", "image": "datadog/agent:7"},
						},
					},
				},
			},
		}},
	}
	got := FromObjects(objs)
	assert.Equal(t, []string{"datadog/agent:7", "myapp:v2"}, got)
}

func TestFromObjects_ContainerMissingImageField(t *testing.T) {
	objs := []unstructured.Unstructured{
		{Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata":   map[string]any{"name": "test"},
			"spec": map[string]any{
				"template": map[string]any{
					"spec": map[string]any{
						"containers": []any{
							map[string]any{"name": "no-image"},
							map[string]any{"name": "has-image", "image": "nginx:1.25"},
							map[string]any{"name": "empty-image", "image": ""},
						},
					},
				},
			},
		}},
	}
	got := FromObjects(objs)
	assert.Equal(t, []string{"nginx:1.25"}, got)
}

func TestFromObjects_SpecWrongType(t *testing.T) {
	objs := []unstructured.Unstructured{
		// spec is a string instead of a map
		{Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata":   map[string]any{"name": "test"},
			"spec":       "not-a-map",
		}},
		// template.spec is a string instead of a map
		{Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata":   map[string]any{"name": "test2"},
			"spec": map[string]any{
				"template": map[string]any{
					"spec": "also-not-a-map",
				},
			},
		}},
		// containers is a string instead of a slice
		{Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata":   map[string]any{"name": "test3"},
			"spec": map[string]any{
				"template": map[string]any{
					"spec": map[string]any{
						"containers": "not-a-slice",
					},
				},
			},
		}},
		// container entry is a string instead of a map
		{Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata":   map[string]any{"name": "test4"},
			"spec": map[string]any{
				"template": map[string]any{
					"spec": map[string]any{
						"containers": []any{"not-a-map"},
					},
				},
			},
		}},
	}
	got := FromObjects(objs)
	assert.Empty(t, got)
}

func TestFromObjects_MissingSpec(t *testing.T) {
	objs := []unstructured.Unstructured{
		{Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata":   map[string]any{"name": "no-spec"},
		}},
	}
	got := FromObjects(objs)
	assert.Empty(t, got)
}

func TestFromObjects_MissingKind(t *testing.T) {
	objs := []unstructured.Unstructured{
		{Object: map[string]any{
			"apiVersion": "apps/v1",
			"metadata":   map[string]any{"name": "no-kind"},
			"spec": map[string]any{
				"template": map[string]any{
					"spec": map[string]any{
						"containers": []any{
							map[string]any{"name": "app", "image": "nginx:1.25"},
						},
					},
				},
			},
		}},
	}
	got := FromObjects(objs)
	assert.Empty(t, got)
}

func TestFromObjects_PodTemplate(t *testing.T) {
	objs := []unstructured.Unstructured{
		{Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "PodTemplate",
			"metadata":   map[string]any{"name": "worker-template"},
			"template": map[string]any{
				"spec": map[string]any{
					"containers": []any{
						map[string]any{"name": "worker", "image": "worker:v3"},
					},
					"initContainers": []any{
						map[string]any{"name": "setup", "image": "setup:latest"},
					},
				},
			},
		}},
	}
	got := FromObjects(objs)
	assert.Equal(t, []string{"setup:latest", "worker:v3"}, got)
}
