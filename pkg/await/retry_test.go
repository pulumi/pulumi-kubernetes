package await

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
						return fmt.Errorf("Operation failed")
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
			retrier: mockRetrier(func(uint) error { return fmt.Errorf("Operation failed") }).
				WithMaxRetries(3).
				WithBackoffFactor(2),
			err:           fmt.Errorf("Operation failed"),
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
