package await

import (
	"time"
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
		if err != nil {
			r.sleep(r.waitTime)
		} else {
			break
		}
		r.waitTime = r.waitTime * time.Duration(r.backOffFactor)
	}
	return err
}
