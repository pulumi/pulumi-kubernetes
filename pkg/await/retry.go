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

package await

import (
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
)

type retrier struct {
	try           func(currTry uint) error
	sleep         func(time.Duration)
	waitTime      time.Duration
	tries         uint
	maxRetries    uint
	backOffFactor uint
}

func sleepingRetry(try func(uint) error) *retrier {
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

func (r *retrier) WithBackoffFactor(t uint) *retrier {
	r.backOffFactor = t
	return r
}

func (r *retrier) Do() error {
	var err error
	for r.tries <= r.maxRetries {
		err = r.try(r.tries)
		r.tries++
		if errors.IsNotFound(err) || meta.IsNoMatchError(err) {
			r.sleep(r.waitTime)
		} else {
			break
		}
		r.waitTime = r.waitTime * time.Duration(r.backOffFactor)
	}
	return err
}
