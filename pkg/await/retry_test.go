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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func mockRetrier(try func(uint) error) *retrier {
	return &retrier{
		try:           try,
		sleep:         func(time.Duration) {},
		waitTime:      time.Second * 1,
		tries:         0,
		maxRetries:    5,
		backOffFactor: 3,
	}
}

func notFound(msg string) error {
	return &errors.StatusError{ErrStatus: metav1.Status{
		Reason:  metav1.StatusReasonNotFound,
		Message: msg}}
}

func Test_Retrier(t *testing.T) {
	tests := []struct {
		description   string
		retrier       *retrier
		err           error
		tries         uint
		finalWaitTime time.Duration
	}{
		{
			description:   "Should succeed after 1 retry if maxRetries == 0",
			retrier:       mockRetrier(func(uint) error { return nil }).WithMaxRetries(0),
			tries:         1,
			finalWaitTime: 1 * time.Second,
		},
		{
			description:   "Should succeed after one try if maxRetries > 0",
			retrier:       mockRetrier(func(uint) error { return nil }).WithMaxRetries(10),
			tries:         1,
			finalWaitTime: 1 * time.Second,
		},
		{
			description: "Should back off if first request fails",
			retrier: mockRetrier(
				func(i uint) error {
					if i == 0 {
						return notFound("Operation failed")
					}
					return nil
				}).
				WithMaxRetries(10).
				WithBackoffFactor(5),
			tries:         2,
			finalWaitTime: 5 * time.Second,
		},
		{
			description: "Should fail if retry budget exceeded",
			retrier: mockRetrier(func(uint) error { return notFound("Operation failed") }).
				WithMaxRetries(3).
				WithBackoffFactor(2),
			err:           notFound("Operation failed"),
			tries:         4,
			finalWaitTime: 16 * time.Second,
		},
	}

	for _, test := range tests {
		err := test.retrier.Do()
		assert.Equal(t, test.err, err, test.description)
		assert.Equal(t, test.tries, test.retrier.tries)
		assert.Equal(t, test.finalWaitTime, test.retrier.waitTime)
	}
}
