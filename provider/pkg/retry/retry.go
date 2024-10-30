// Copyright 2016-2019, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package retry

import (
	"time"
)

type retrier struct {
	try           func(currTry uint) error
	sleep         func(time.Duration)
	waitTime      time.Duration
	tries         uint
	maxRetries    uint
	backOffFactor uint16
}

func SleepingRetry(try func(uint) error) *retrier {
	return &retrier{
		try:           try,
		sleep:         time.Sleep,
		waitTime:      time.Second * 1,
		tries:         0,
		maxRetries:    5,
		backOffFactor: 2,
	}
}

func (r *retrier) WithMaxRetries(n uint) *retrier {
	r.maxRetries = n
	return r
}

// WithSleep uses a custom sleep method for retries, useful for testing.
func (r *retrier) WithSleep(s func(time.Duration)) *retrier {
	if s == nil {
		return r
	}
	r.sleep = s
	return r
}

func (r *retrier) WithBackoffFactor(t uint16) *retrier {
	r.backOffFactor = t
	return r
}

func (r *retrier) Do(allowedErrFuncs ...func(error) bool) error {
	var err error
	for r.tries <= r.maxRetries {
		err = r.try(r.tries)
		r.tries++

		shouldRetry := false
		for _, errFunc := range allowedErrFuncs {
			if errFunc(err) {
				shouldRetry = true
				break
			}
		}
		if !shouldRetry {
			break
		}
		r.sleep(r.waitTime)
		r.waitTime = r.waitTime * time.Duration(r.backOffFactor)
	}
	return err
}
