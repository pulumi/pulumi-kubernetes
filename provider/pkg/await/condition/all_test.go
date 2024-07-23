package condition

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

// TODO: unit test the event filter stuff so the aggregator doesn't need to

func TestAll(t *testing.T) {
	ctx := context.Background()

	source1 := Static(make(chan watch.Event))
	source2 := Static(make(chan watch.Event))

	obj := &unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{
			"name": "foo",
		},
	}}

	want1 := watch.Event{Type: watch.Added, Object: obj}
	want2 := watch.Event{Type: watch.Deleted, Object: obj}

	cond, err := NewAll(
		NewOn(ctx, source1, obj, want1),
		NewOn(ctx, source2, obj, want2),
	)
	require.NoError(t, err)

	go func() {
		source1 <- want1
		source2 <- watch.Event{Type: watch.Modified, Object: obj}
		source2 <- want2
		close(source1)
		close(source2)
	}()

	cond.Range(func(e watch.Event) bool {
		_ = cond.Observe(e)
		return true
	})

	done, err := cond.Satisfied()
	assert.NoError(t, err)
	assert.True(t, done)
}
