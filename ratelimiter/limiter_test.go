package ratelimiter

import (
	"testing"
	"time"
)

func TestBasicWorkFlow(t *testing.T) {
	limiter, err := CreateLimiter(LimiterOption{
		Policy:   "leaky",
		Rate:     100,
		Capacity: 10,
	})
	if err != nil {
		t.Errorf("Instantiate limiter error")
	}
	for i := 0; i < 10; i++ {
		ok, dur := limiter.Acquire(10)
		if !ok {
			t.Errorf("i:%v, dur:%v, acuqire should be succ", i, dur)
		}
	}
}

func TestSlackInProperSize(t *testing.T) {
	limiter, err := CreateLimiter(LimiterOption{
		Policy:   "leaky",
		Rate:     10,
		Capacity: 10,
	})
	if err != nil {
		t.Errorf("Instantiate limiter error")
	}

	ok, dur := limiter.Acquire(0)
	if !ok {
		t.Error("dur:%v, first acuqire should be succ", dur)
	}
	// Sleep for 1 sec, make sure the buffer can handle the following request.
	time.Sleep(time.Second * 1)

	for i := 0; i < 10; i++ {
		ok, dur := limiter.Acquire(0)
		if !ok {
			t.Errorf("i:%v, dur:%v, acuqire should be succ", i, dur)
		}
	}
}

func TestConcurrency(t *testing.T) {
	limiter, err := CreateLimiter(LimiterOption{
		Policy:   "leaky",
		Rate:     10,
		Capacity: 10,
	})
	if err != nil {
		t.Errorf("Instantiate limiter error")
	}

	size := 10
	done := make(chan bool, size)

	for i := 0; i < size; i++ {
		go func(id int) {
			ok, _ := limiter.Acquire(0)
			done <- ok
		}(i)
	}

	nOK := 0
	for i := 0; i < size; i++ {
		c := <-done
		if c {
			nOK += 1
		}
	}
	if nOK != 1 {
		t.Errorf("okNum should be 1 (the first request)\n")
	}
}
