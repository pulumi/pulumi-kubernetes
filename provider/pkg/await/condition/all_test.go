package condition

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

func TestAll(t *testing.T) {
	ctx := context.Background()

	source := Static(make(chan watch.Event))

	obj := &unstructured.Unstructured{
		Object: map[string]any{
			"metadata": map[string]any{
				"name": "foo",
			},
		},
	}

	want1 := watch.Event{Type: watch.Added, Object: obj}
	want2 := watch.Event{Type: watch.Deleted, Object: obj}

	cond, err := NewAll(
		NewOn(ctx, source, obj, want1),
		NewOn(ctx, source, obj, want2),
	)
	require.NoError(t, err)

	go func() {
		source <- want1
		source <- watch.Event{Type: watch.Modified, Object: obj}
		source <- want2
		close(source)
	}()

	cond.Range(func(e watch.Event) bool {
		_ = cond.Observe(e)
		return true
	})

	done, err := cond.Satisfied()
	assert.NoError(t, err)
	assert.True(t, done)
}
