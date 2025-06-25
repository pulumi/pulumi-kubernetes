// Copyright 2024, Pulumi Corporation.
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

package condition

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

var pod = &unstructured.Unstructured{
	Object: map[string]any{
		"apiVersion": "v1",
		"kind":       "Pod",
		"metadata": map[string]any{
			"name":      "foo",
			"namespace": "default",
		},
		"spec": map[string]any{
			"containers": []any{
				map[string]any{
					"name":  "foo",
					"image": "nginx",
				},
			},
		},
	},
}

type get404 struct{}

func (get404) Get(ctx context.Context, name string, opts metav1.GetOptions, sub ...string) (*unstructured.Unstructured, error) {
	return nil, k8serrors.NewNotFound(schema.GroupResource{}, name)
}

type get503 struct{}

func (get503) Get(ctx context.Context, name string, opts metav1.GetOptions, sub ...string) (*unstructured.Unstructured, error) {
	return nil, k8serrors.NewServiceUnavailable("boom")
}

type get200 struct{ obj *unstructured.Unstructured }

func (g *get200) Get(context.Context, string, metav1.GetOptions, ...string) (*unstructured.Unstructured, error) {
	return g.obj, nil
}

type getsequence struct {
	getters []objectGetter
	idx     int
}

func (g *getsequence) Get(ctx context.Context, name string, opts metav1.GetOptions, sub ...string) (*unstructured.Unstructured, error) {
	defer func() { g.idx++ }()
	return g.getters[g.idx].Get(ctx, name, opts, sub...)
}

type getdeferred struct {
	wg sync.WaitGroup
	objectGetter
}

func (g *getdeferred) Get(ctx context.Context, name string, opts metav1.GetOptions, sub ...string) (*unstructured.Unstructured, error) {
	g.wg.Wait()
	return g.objectGetter.Get(ctx, name, opts, sub...)
}

func TestDeleted(t *testing.T) {
	stdout := logbuf{os.Stdout}

	t.Run("already deleted", func(t *testing.T) {
		ctx := context.Background()
		getter := get404{}

		source := &DeletionSource{obj: pod, getter: getter, source: Static(nil)}

		cond, err := NewDeleted(ctx, source, getter, stdout, pod)
		assert.NoError(t, err)

		cond.Range(func(e watch.Event) bool {
			err := cond.Observe(e)
			assert.NoError(t, err)
			return true
		})

		done, err := cond.Satisfied()
		assert.NoError(t, err)
		assert.True(t, done)
	})

	t.Run("deleted during watch", func(t *testing.T) {
		ctx := context.Background()

		getter := &get200{pod}
		source := Static(make(chan watch.Event, 1))

		cond, err := NewDeleted(ctx, source, getter, stdout, pod)
		assert.NoError(t, err)

		seen := make(chan struct{})
		go cond.Range(func(e watch.Event) bool {
			err := cond.Observe(e)
			assert.NoError(t, err)
			seen <- struct{}{}
			return true
		})

		source <- watch.Event{Type: watch.Modified, Object: pod}
		<-seen
		done, err := cond.Satisfied()
		assert.NoError(t, err)
		assert.False(t, done)

		source <- watch.Event{Type: watch.Deleted, Object: pod}
		<-seen
		done, err = cond.Satisfied()
		assert.NoError(t, err)
		assert.True(t, done)
	})

	t.Run("issue-3317", func(t *testing.T) {
		ctx := context.Background()

		getter := &getdeferred{objectGetter: &get200{pod}}
		getter.wg.Add(1)
		source := Static(make(chan watch.Event, 1))

		cond, err := NewDeleted(ctx, source, getter, stdout, pod)
		assert.NoError(t, err)

		wg := sync.WaitGroup{}
		wg.Add(1)
		seen := make(chan struct{})
		go func() {
			cond.Range(func(e watch.Event) bool {
				err := cond.Observe(e)
				assert.NoError(t, err)
				seen <- struct{}{}
				return true
			})
			wg.Done()
		}()

		// Fire a 'deleted' event to be observed before the getter returns.
		source <- watch.Event{Type: watch.Deleted, Object: pod}
		<-seen
		assert.True(t, cond.deleted.Load())

		// Unblock the getter to return a (stale) object, which should be ignored.
		getter.wg.Done()

		// Wait for the Range operation to complete.
		wg.Wait()

		done, err := cond.Satisfied()
		assert.NoError(t, err)
		assert.True(t, done)
	})

	t.Run("times out", func(t *testing.T) {
		getter := &get200{pod}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		cond, err := NewDeleted(ctx, Static(nil), getter, stdout, pod)
		assert.NoError(t, err)

		cond.Range(nil)

		done, err := cond.Satisfied()
		assert.NoError(t, err)
		assert.False(t, done)
	})

	t.Run("times out with finalizers", func(t *testing.T) {
		getter := &get200{pod}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		podWithFinalizers := pod.DeepCopy()
		podWithFinalizers.SetFinalizers([]string{"stuck"})

		buf := &strings.Builder{}
		cond, err := NewDeleted(ctx, Static(nil), getter, logbuf{buf}, podWithFinalizers)
		assert.NoError(t, err)

		cond.Range(nil)

		assert.Contains(t, buf.String(), "finalizers might be preventing deletion")
	})

	// TODO: It's questionable whether we still need to test this behavior. I
	// suspect this stems from earlier error handling code around our watch
	// logic, which is largely obviated by our use of informers now. In other
	// words, we needed this when we weren't handling the sort of watch errors
	// Informers handle automatically.
	t.Run("times out with recovery", func(t *testing.T) {
		getter := &getsequence{[]objectGetter{&get200{pod}, get404{}}, 0}

		source := &DeletionSource{obj: pod, getter: getter, source: Static(nil)}

		ctx, cancel := context.WithCancel(context.Background())
		cond, err := NewDeleted(ctx, source, getter, stdout, pod)
		assert.NoError(t, err)

		cancel()
		cond.Range(nil)

		done, err := cond.Satisfied()
		assert.NoError(t, err)
		assert.True(t, done)
	})

	t.Run("unexpected error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		buf := &strings.Builder{}
		cond, err := NewDeleted(ctx, Static(nil), get503{}, logbuf{buf}, pod)
		assert.NoError(t, err)

		cancel()
		cond.Range(nil)

		done, err := cond.Satisfied()
		assert.NoError(t, err)
		assert.False(t, done)
		assert.Contains(t, buf.String(), "boom")
	})
}
