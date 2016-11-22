package ratelimiter

import (
	"testing"
	"time"
)

func TestBasicWorkFlow(t *testing.T) {
	limiter, err := CreateBucket(BucketOption{
		Name:     "leaky",
		Rate:     100,
		Capacity: 10,
	})
	if err != nil {
		t.Errorf("Instantiate bucket error")
	}
	// limiter := NewLeakyBucket(100, 10)
	for i := 0; i < 10; i++ {
		ok, dur := limiter.Acquire(10)
		if !ok {
			t.Errorf("i:%v, dur:%v, acuqire should be succ", i, dur)
		}
	}
}

func TestSlackInProperSize(t *testing.T) {
	limiter, err := CreateBucket(BucketOption{
		Name:     "leaky",
		Rate:     10,
		Capacity: 10,
	})
	if err != nil {
		t.Errorf("Instantiate bucket error")
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
	limiter, err := CreateBucket(BucketOption{
		Name:     "leaky",
		Rate:     10,
		Capacity: 10,
	})
	if err != nil {
		t.Errorf("Instantiate bucket error")
	}

	requestNum := 10
	done := make(chan bool, requestNum)

	for i := 0; i < requestNum; i++ {
		go func(id int) {
			ok, _ := limiter.Acquire(0)
			done <- ok
		}(i)
	}

	okNum := 0
	for i := 0; i < requestNum; i++ {
		c := <-done
		if c {
			okNum += 1
		}
	}
	if okNum != 1 {
		t.Errorf("okNum should be 1 (the first request)\n")
	}
}
