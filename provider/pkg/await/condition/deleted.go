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
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"

	checkerlog "github.com/pulumi/cloud-ready-checks/pkg/checker/logging"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/cli-utils/pkg/kstatus/status"
)

var _ Satisfier = (*Deleted)(nil)

// Deleted condition succeeds when GET on a resource 404s or when a Deleted
// event is received for the resource.
type Deleted struct {
	observer *ObjectObserver
	logger   logger
	deleted  atomic.Bool
	getter   objectGetter
}

// NewDeleted constructs a new Deleted condition.
func NewDeleted(
	ctx context.Context,
	source Source,
	getter objectGetter,
	logger logger,
	obj *unstructured.Unstructured,
) (*Deleted, error) {
	dc := &Deleted{
		observer: NewObjectObserver(ctx, source, obj),
		logger:   logger,
		getter:   getter,
	}
	return dc, nil
}

// Range confirms the object exists before establishing an Informer.
// If a Deleted event isn't Observed by the time the underlying Observer is
// exhausted, we attempt a final lookup on the cluster to be absolutely sure it
// still exists.
func (dc *Deleted) Range(yield func(watch.Event) bool) {
	dc.getClusterState()
	if dc.deleted.Load() {
		// Already deleted, nothing more to do.
		return
	}

	// Iterate over the underlying Observer's events.
	dc.observer.Range(yield)

	if dc.deleted.Load() {
		// Nothing more to do.
		return
	}

	// Attempt one last lookup if the object still exists. (This is legacy behavior
	// that might be unnecessary since we're using Informers instead of
	// Watches now.)
	dc.getClusterState()
	if dc.deleted.Load() {
		return
	}

	// Let the user know we might be blocked if the object has finalizers.
	// https://github.com/pulumi/pulumi-kubernetes/issues/1418
	finalizers := dc.Object().GetFinalizers()
	if len(finalizers) == 0 {
		return
	}
	dc.logger.LogMessage(checkerlog.WarningMessage(
		fmt.Sprintf("finalizers might be preventing deletion (%s)", strings.Join(finalizers, ", ")),
	))
}

// Observe watches for Deleted events.
func (dc *Deleted) Observe(e watch.Event) error {
	if e.Type == watch.Deleted {
		dc.deleted.Store(true)
	}
	return dc.observer.Observe(e)
}

// Satisfied returns true if a Deleted event was Observed. Otherwise a status
// message will be logged, if available.
func (dc *Deleted) Satisfied() (bool, error) {
	if dc.deleted.Load() {
		return true, nil
	}

	uns := dc.Object()
	r, _ := status.Compute(uns)
	if r.Message != "" {
		dc.logger.LogMessage(checkerlog.StatusMessage(r.Message))
	}

	return false, nil
}

// Object returns the last-known state of the object we're deleting.
func (dc *Deleted) Object() *unstructured.Unstructured {
	return dc.observer.Object()
}

// getClusterState performs a GET against the cluster and updates state to
// reflect whether the object still exists or not.
func (dc *Deleted) getClusterState() {
	_, err := dc.getter.Get(context.Background(), dc.Object().GetName(), metav1.GetOptions{})
	if err == nil {
		// Still exists.
		dc.deleted.Store(false)
		return
	}
	var statusErr *k8serrors.StatusError
	if errors.As(err, &statusErr) {
		dc.deleted.Store(statusErr.ErrStatus.Code == http.StatusNotFound)
	}
}
