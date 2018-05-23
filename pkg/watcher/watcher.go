// Copyright 2016-2018, Pulumi Corporation.
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
	"fmt"
	"math/rand"
	"time"

	"github.com/golang/glog"
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
	objName  string
	pollFunc pollFunc
}

// ForObject creates an `ObjectWatcher` to watch some object.
func ForObject(
	clientForResource dynamic.ResourceInterface, name string,
) *ObjectWatcher {
	return &ObjectWatcher{
		objName: name,
		pollFunc: func() (*unstructured.Unstructured, error) {
			obj, err := clientForResource.Get(name, metav1.GetOptions{})
			if err != nil {
				// Log the error.
				glog.V(3).Infof("Received error polling for '%s': %#v", name, err)
				return nil, err
			}
			return obj, nil
		},
	}
}

// WatchUntil will block and watch an object until some predicate is true, or a timeout occurs.
func (ow *ObjectWatcher) WatchUntil(pred Predicate, timeout time.Duration) error {
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
func (ow *ObjectWatcher) RetryUntil(r Retry, timeout time.Duration) error {
	untilNoError := func(obj *unstructured.Unstructured, err error) (bool, error) {
		if retryErr := r(obj, err); retryErr != nil {
			rerr, isRetryable := retryErr.(*RetryError)
			// Non-retryable error; return error and true to stop watching.
			if !isRetryable {
				return true, retryErr
			}
			glog.V(3).Infof("Retrying operation with message: %s", rerr.Error())

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
) error {
	const maxTimeout = 30000 * time.Millisecond

	timeoutCh := make(chan struct{})
	go func() {
		time.Sleep(timeout)
		timeoutCh <- struct{}{}
	}()

	results := make(chan result)
	poll := func() {
		obj, err := ow.pollFunc()
		results <- result{Obj: obj, Err: err}
	}

	// Poll with exponential backoff.
	wait := 500 * time.Millisecond
	for {
		// Race between timeout and getting one Kubernetes object from the polling function. If the
		// object does not satisfy `until`, we'll back off.
		go poll()
		select {
		case <-timeoutCh:
			return fmt.Errorf("Timeout occurred polling for '%s'", ow.objName)
		case res := <-results:
			if stop, err := until(res.Obj, res.Err); err != nil {
				return err
			} else if stop {
				return nil
			}
			time.Sleep(wait + time.Duration(rand.Intn(int(float64(wait)*0.2))))
			wait *= 2
			if maxTimeout < wait {
				wait = maxTimeout
			}
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

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
