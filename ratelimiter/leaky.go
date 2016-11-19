package ratelimiter

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	sync.Mutex
	perRequest time.Duration
	sleepFor   time.Duration
	maxSlack   time.Duration // for burstiness
	last       time.Time
}

func NewLeakyBucket(rate int) *LeakyBucket {
	return &LeakyBucket{
		perRequest: time.Second / time.Duration(rate),
		maxSlack:   -10 * time.Second / time.Duration(rate),
	}
}

func (self *LeakyBucket) Acquire(maxWait int) (bool, time.Duration) {
	self.Lock()
	defer self.Unlock()

	now := time.Now()

	if self.last.IsZero() {
		self.last = now
		return true, time.Duration(0)
	}

	waitDur := self.sleepFor + self.perRequest - now.Sub(self.last)
	if waitDur > time.Duration(time.Millisecond*time.Duration(maxWait)) {
		return false, waitDur
	}

	self.sleepFor = waitDur
	if self.sleepFor < self.maxSlack {
		self.sleepFor = self.maxSlack
	}

	if self.sleepFor > 0 {
		time.Sleep(self.sleepFor)
		self.last = now.Add(self.sleepFor)
		self.sleepFor = 0
		return true, self.sleepFor
	}
	self.last = now
	return true, time.Duration(0)
}
