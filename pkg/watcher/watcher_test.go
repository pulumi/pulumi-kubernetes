package watcher

import (
	"strings"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func testObjWatcher(poll pollFunc) *ObjectWatcher {
	return &ObjectWatcher{
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
			pollFuncCalls, watchFuncCalls := 0, 0
			err := testObjWatcher(
				func() (*unstructured.Unstructured, error) {
					pollFuncCalls++
					return test.pollFunc()
				}).
				WatchUntil(
					func(obj *unstructured.Unstructured) bool {
						watchFuncCalls++
						return test.predicate(obj)
					},
					test.timeout)
			if err == nil || !strings.HasPrefix(err.Error(), timeoutErrPrefix) {
				t.Errorf("%s: Polling should have timed out", test.name)
			}
			if !test.targetPollFuncCalls(pollFuncCalls) {
				t.Errorf("%s: Got %d poll function calls, which did not satisfy the test predicate", test.name, pollFuncCalls)
			}
			if !test.targetWatchFuncCalls(watchFuncCalls) {
				t.Errorf("%s: Got %d watch function calls, which did not satisfy the test predicate", test.name, watchFuncCalls)
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
