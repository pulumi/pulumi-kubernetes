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
	"context"
	"strings"
	"testing"
	"time"

	"go.uber.org/atomic"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func testObjWatcher(ctx context.Context, poll pollFunc) *ObjectWatcher {
	return &ObjectWatcher{
		ctx:      ctx,
		pollFunc: poll,
	}
}

var timeoutErrPrefix = "Timeout occurred polling for"

func Test_WatchUntil_PollFuncTimeout(t *testing.T) {
	type timeoutTest struct {
		name                 string
		targetPollFuncCalls  func(int) bool
		targetWatchFuncCalls func(int) bool
		pollFunc             pollFunc
		predicate            Predicate
		timeout              time.Duration
	}
	timeoutTests := []timeoutTest{
		{
			name:                 "PollFuncTimeout",
			targetPollFuncCalls:  func(i int) bool { return i == 1 },
			targetWatchFuncCalls: func(i int) bool { return i == 0 },
			pollFunc: func() (*unstructured.Unstructured, error) {
				time.Sleep(2 * time.Second)
				return nil, nil
			},
			predicate: func(*unstructured.Unstructured) bool {
				return true
			},
			timeout: 100 * time.Millisecond,
		},
		{
			name:                 "PredicateFuncTimeout",
			targetPollFuncCalls:  func(i int) bool { return i > 1 },
			targetWatchFuncCalls: func(i int) bool { return i > 0 },
			pollFunc: func() (*unstructured.Unstructured, error) {
				return &unstructured.Unstructured{}, nil
			},
			predicate: func(*unstructured.Unstructured) bool {
				return false // Always false.
			},
			timeout: 3 * time.Second,
		},
	}

	testCompleted := make(chan struct{})
	for _, test := range timeoutTests {
		go func(test timeoutTest) {
			pollFuncCalls, watchFuncCalls := atomic.NewInt32(0), atomic.NewInt32(0)
			err := testObjWatcher(
				context.Background(),
				func() (*unstructured.Unstructured, error) {
					pollFuncCalls.Inc()
					return test.pollFunc()
				}).
				WatchUntil(
					func(obj *unstructured.Unstructured) bool {
						watchFuncCalls.Inc()
						return test.predicate(obj)
					},
					test.timeout)
			if err == nil || !strings.HasPrefix(err.Error(), timeoutErrPrefix) {
				t.Errorf("%s: Polling should have timed out", test.name)
			}
			if !test.targetPollFuncCalls(int(pollFuncCalls.Load())) {
				t.Errorf("%s: Got %d poll function calls, which did not satisfy the test predicate", test.name, pollFuncCalls.Load())
			}
			if !test.targetWatchFuncCalls(int(watchFuncCalls.Load())) {
				t.Errorf("%s: Got %d watch function calls, which did not satisfy the test predicate", test.name, watchFuncCalls.Load())
			}
			testCompleted <- struct{}{}
		}(test)
	}

	testsCompleted := 0
	for range testCompleted {
		testsCompleted++
		if testsCompleted == len(timeoutTests) {
			return
		}
	}
}

func Test_WatchUntil_Success(t *testing.T) {
	// Timeout because the `WatchUntil` predicate always returns false.
	err := testObjWatcher(
		context.Background(),
		func() (*unstructured.Unstructured, error) {
			return &unstructured.Unstructured{}, nil
		}).
		WatchUntil(
			func(*unstructured.Unstructured) bool {
				return true // Always true.
			},
			1*time.Second)
	if err != nil {
		t.Error("Expected watch to terminate without error")
	}
}

func Test_RetryUntil_Success(t *testing.T) {
	// Timeout because the `WatchUntil` predicate always returns false.
	err := testObjWatcher(
		context.Background(),
		func() (*unstructured.Unstructured, error) {
			return &unstructured.Unstructured{}, nil
		}).
		RetryUntil(
			func(*unstructured.Unstructured, error) error {
				return nil // Always succeeds.
			},
			1*time.Second)
	if err != nil {
		t.Error("Expected watch to terminate without error")
	}
}

func Test_RetryUntil_Cancel(t *testing.T) {
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	// Timeout because the `WatchUntil` predicate always returns false.
	err := testObjWatcher(
		cancelCtx,
		func() (*unstructured.Unstructured, error) {
			return &unstructured.Unstructured{}, nil
		}).
		RetryUntil(
			func(*unstructured.Unstructured, error) error {
				return nil // Always succeeds.
			},
			1*time.Second)

	if err == nil {
		t.Error("Expected watch to terminate with an initialization error")
	}
}
