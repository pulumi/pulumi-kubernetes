package gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestGVKFromRef(t *testing.T) {
	tests := []struct {
		given string
		want  schema.GroupVersionKind
	}{
		{
			given: "#/definitions/io.k8s.api.apps.v1.Deployment",
			want:  schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
		},
		{
			given: "io.k8s.api.apps.v1.Deployment",
			want:  schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
		},
		{
			given: "#/types/kubernetes:helm.sh/v4:PostRenderer",
			want:  schema.GroupVersionKind{Group: "helm.sh", Version: "v4", Kind: "PostRenderer"},
		},
		{
			given: "argoproj.io/v1alpha1:Rollout",
			want:  schema.GroupVersionKind{Group: "argoproj.io", Version: "v1alpha1", Kind: "Rollout"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.given, func(t *testing.T) {
			actual := GVKFromRef(tt.given)
			assert.Equal(t, tt.want, actual)
		})
	}
}
