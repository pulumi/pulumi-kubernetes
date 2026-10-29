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
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/cli-utils/pkg/kstatus/status"

	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
)

var _ Satisfier = (*Deleted)(nil)

// Deleted succeeds when the recorded object no longer exists. When its UID is
// known, a live object with the same name and a different UID also proves deletion.
type Deleted struct {
	ctx      context.Context
	cancel   func()
	observer *ObjectObserver
	logger   logger
	deleted  atomic.Bool
	getter   objectGetter

	// Keep the recorded identity separate from the observer's latest state.
	uid       types.UID
	name      string
	namespace string
}

// NewDeleted constructs a new Deleted condition.
func NewDeleted(
	ctx context.Context,
	source Source,
	getter objectGetter,
	logger logger,
	obj *unstructured.Unstructured,
) (*Deleted, error) {
	ctx, cancel := context.WithCancel(ctx) //nolint:gosec // cancel stored in dc.cancel, called on teardown
	dc := &Deleted{
		ctx:      ctx,
		cancel:   cancel,
		observer: NewObjectObserver(ctx, source, obj),
		logger:   logger,
		getter:   getter,

		uid:       obj.GetUID(),
		name:      obj.GetName(),
		namespace: obj.GetNamespace(),
	}
	return dc, nil
}

// Range establishes an Informer and confirms the object still exists.
// If a Deleted event isn't Observed by the time the underlying Observer is
// exhausted, we attempt a final lookup on the cluster to be absolutely sure it
// still exists.
func (dc *Deleted) Range(yield func(watch.Event) bool) {
	// Start listening to events before we check if the resource has already
	// been deleted. This avoids races where the object is deleted in-between
	// the 404 check and this watch.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for e := range dc.observer.Range {
			if dc.deleted.Load() {
				return
			}
			if !yield(e) {
				return
			}
		}
	}()

	dc.refreshClusterState()
	if dc.deleted.Load() {
		// Already deleted, nothing more to do. Our informer will get cleaned up
		// when its context is canceled.
		dc.cancel()
		wg.Wait()
		return
	}

	wg.Wait()

	if dc.deleted.Load() {
		// Nothing more to do.
		return
	}

	// Attempt one last lookup if the object still exists. (This is legacy
	// behavior that might be unnecessary since we're using Informers instead of
	// Watches now.)
	dc.refreshClusterState()
	if dc.deleted.Load() {
		return
	}

	// Let the user know we might be blocked if the object has finalizers.
	// https://github.com/pulumi/pulumi-kubernetes/issues/1418
	finalizers := dc.Object().GetFinalizers()
	dc.logger.Log(diag.Warning,
		fmt.Sprintf("finalizers might be preventing deletion (%s)", strings.Join(finalizers, ", ")),
	)
}

// Observe tracks deletion without adopting another object with the same name.
func (dc *Deleted) Observe(e watch.Event) error {
	if dc.uid != "" {
		obj, ok := e.Object.(*unstructured.Unstructured)
		if !ok || obj == nil || obj.GetName() != dc.name || obj.GetNamespace() != dc.namespace {
			return nil
		}
		if obj.GetUID() != dc.uid {
			// Informer caches may contain an older object. Confirm another UID
			// with a live read before treating it as a replacement.
			if e.Type != watch.Deleted && dc.replaced(obj) {
				dc.refreshClusterState()
			}
			return nil
		}
	}
	if e.Type == watch.Deleted {
		dc.deleted.Store(true)
	}
	return dc.observer.Observe(e)
}

func (dc *Deleted) replaced(obj *unstructured.Unstructured) bool {
	return dc.uid != "" && obj != nil && obj.GetUID() != "" && obj.GetUID() != dc.uid &&
		obj.GetName() == dc.name && obj.GetNamespace() == dc.namespace
}

// Satisfied reports whether the recorded object is gone, logging its status while pending.
func (dc *Deleted) Satisfied() (bool, error) {
	if dc.deleted.Load() {
		return true, nil
	}

	uns := dc.Object()
	r, _ := status.Compute(uns)
	if r.Message != "" {
		dc.logger.LogStatus(diag.Info, r.Message)
	}

	return false, nil
}

// Object returns the last-known state of the object we're deleting.
func (dc *Deleted) Object() *unstructured.Unstructured {
	return dc.observer.Object()
}

// refreshClusterState performs a GET against the cluster and updates state to
// reflect whether the object still exists or not.
func (dc *Deleted) refreshClusterState() {
	// Our context might be closed, but we still want to issue this request
	// even if we're shutting down.
	ctx := context.WithoutCancel(dc.ctx)
	live, err := dc.getter.Get(ctx, dc.name, metav1.GetOptions{})
	if err == nil {
		if dc.replaced(live) {
			dc.deleted.Store(true)
		}
		return
	}
	if k8serrors.IsNotFound(err) {
		dc.deleted.Store(true)
	} else {
		dc.logger.LogStatus(diag.Warning,
			"unexpected error while checking cluster state: "+err.Error(),
		)
	}
}
