// Copyright 2016-2022, Pulumi Corporation.
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

package watcher

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

// Predicate is a function that returns true when we can stop watching in `WatchUntil`.
type Predicate func(*unstructured.Unstructured) bool

// Retry is a function that returns nil when an operation that can fail some number of times has
// succeeded, and `RetryUntil` can stop retrying.
type Retry func(*unstructured.Unstructured, error) error

type pollFunc func() (*unstructured.Unstructured, error)

// ObjectWatcher will block and watch or retry some operation on an object until a timeout, or some
// condition is true.
type ObjectWatcher struct {
	ctx      context.Context
	objName  string
	pollFunc pollFunc
}

// ForObject creates an `ObjectWatcher` to watch some object.
func ForObject(
	ctx context.Context, clientForResource dynamic.ResourceInterface, name string,
) *ObjectWatcher {
	return &ObjectWatcher{
		ctx:     ctx,
		objName: name,
		pollFunc: func() (*unstructured.Unstructured, error) {
			obj, err := clientForResource.Get(ctx, name, metav1.GetOptions{})
			if err != nil {
				// Log the error.
				logger.V(3).Infof("Received error polling for %q: %#v", name, err)
				return nil, err
			}
			return obj, nil
		},
	}
}

// WatchUntil will block and watch an object until some predicate is true, or a timeout occurs.
func (ow *ObjectWatcher) WatchUntil(pred Predicate, timeout time.Duration) (*unstructured.Unstructured, error) {
	untilTrue := func(obj *unstructured.Unstructured, err error) (bool, error) {
		if err != nil {
			// Error. Return error and true to stop watching.
			return true, err
		} else if pred(obj) {
			// Success. Return true to stop watching.
			return true, nil
		}
		// No error and predicate does not say to stop.
		return false, nil
	}

	return ow.watch(untilTrue, timeout)
}

// RetryUntil will block and retry getting an object until the `Retry` operation succeeds (i.e.,
// does not return an error).
func (ow *ObjectWatcher) RetryUntil(r Retry, timeout time.Duration) (*unstructured.Unstructured, error) {
	untilNoError := func(obj *unstructured.Unstructured, err error) (bool, error) {
		if retryErr := r(obj, err); retryErr != nil {
			rerr, isRetryable := retryErr.(*RetryError)
			// Non-retryable error; return error and true to stop watching.
			if !isRetryable {
				return true, retryErr
			}
			logger.V(3).Infof("Retrying operation with message: %s", rerr.Error())

			// Retryable error. Return false to continue watching.
			return false, nil
		}

		// No error occurred. Return true to stop watching.
		return true, nil
	}

	return ow.watch(untilNoError, timeout)
}

func (ow *ObjectWatcher) watch(
	until func(*unstructured.Unstructured, error) (bool, error), timeout time.Duration,
) (*unstructured.Unstructured, error) {
	timeoutCh := make(chan struct{})
	go func() {
		time.Sleep(timeout)
		timeoutCh <- struct{}{}
	}()

	var obj *unstructured.Unstructured
	results := make(chan result)
	poll := func() {
		o, err := ow.pollFunc()
		results <- result{Obj: o, Err: err}
	}

	wait := 500 * time.Millisecond
	for {
		// Race between timeout and getting one Kubernetes object from the polling function.
		go poll()
		select {
		case <-timeoutCh:
			return nil, timeoutErr(ow.objName, obj)
		case <-ow.ctx.Done():
			return nil, cancellationErr(ow.objName, obj)
		case res := <-results:
			obj = res.Obj
			if stop, err := until(res.Obj, res.Err); err != nil {
				return res.Obj, err
			} else if stop {
				return res.Obj, nil
			}
			// nolint:gosec
			time.Sleep(wait + time.Duration(rand.Intn(int(float64(wait)*0.2))))
		}
	}
}

// --------------------------------------------------------------------------

// Helper utilities.

// --------------------------------------------------------------------------

type result struct {
	Err error
	Obj *unstructured.Unstructured
}

// --------------------------------------------------------------------------

// Errors.
//
// A collection of errors used to implement retry and watch logic.

// --------------------------------------------------------------------------

func timeoutErr(name string, obj *unstructured.Unstructured) error {
	return &watchError{
		object:  obj,
		message: fmt.Sprintf("Timeout occurred polling for '%s'", name),
	}
}

func cancellationErr(name string, obj *unstructured.Unstructured) error {
	return &watchError{
		object:  obj,
		message: fmt.Sprintf("Resource operation was cancelled for '%s'", name),
	}
}

type watchError struct {
	object  *unstructured.Unstructured
	message string
}

var _ error = (*watchError)(nil)

func (we *watchError) Error() string {
	return we.message
}

func (we *watchError) Object() *unstructured.Unstructured {
	return we.object
}

// RetryError is the required return type of RetryFunc. It forces client code
// to choose whether or not a given error is retryable.
type RetryError struct {
	Err       error
	Retryable bool
}

func (re RetryError) Error() string {
	return re.Err.Error()
}

// RetryableError is a helper to create a RetryError that's retryable from a
// given error.
func RetryableError(err error) *RetryError {
	if err == nil {
		return nil
	}
	return &RetryError{Err: err, Retryable: true}
}
