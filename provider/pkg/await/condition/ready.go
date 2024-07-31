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

	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/cli-utils/pkg/kstatus/status"
)

var _ Satisfier = (*Ready)(nil)

// NewReady creates a new Ready condition.
func NewReady(
	ctx context.Context,
	source Source,
	logger logger,
	obj *unstructured.Unstructured,
) *Ready {
	return &Ready{observer: NewObjectObserver(ctx, source, obj), logger: logger}
}

type Ready struct {
	observer *ObjectObserver
	logger   logger
}

// Satisfied returns `true` if the object doesn't have any recognizable status
// conditions.
func (r *Ready) Satisfied() (bool, error) {
	message := "Waiting for readiness"
	s, err := status.Compute(r.Object())
	if err != nil {
		return false, err
	}
	if s.Message != "" {
		message = s.Message
	}
	if r.logger != nil {
		r.logger.LogStatus(diag.Info, message)
	}
	return s.Status == status.CurrentStatus, nil
}

// Range watches events on the object. It is the caller's responsibility to
// ensure the object exists on the cluster beforehand.
func (r *Ready) Range(yield func(watch.Event) bool) {
	r.observer.Range(yield)
}

// Object returns the last-known state of the object we're watching.
func (r *Ready) Object() *unstructured.Unstructured {
	return r.observer.Object()
}

// Observe is a passthrough to the underlying Observer.
func (r *Ready) Observe(e watch.Event) error {
	return r.observer.Observe(e)
}
