package ratelimiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	sync.Mutex
	capacity     int
	avail        int // available token number
	fillInterval time.Duration
	last         time.Time
}

// rate is QPS
func NewTokenBucket(opt BucketOption) (Bucket, error) {
	rate := opt.Rate
	capacity := opt.Capacity

	return &TokenBucket{
		capacity:     capacity,
		fillInterval: time.Second / time.Duration(rate),
	}, nil
}

func (self *TokenBucket) Name() string {
	return "TokenBucket"
}

// maxWait: ms
func (self *TokenBucket) Acquire(maxWait int) (bool, time.Duration) {
	self.Lock()
	defer self.Unlock()

	now := time.Now()

	if self.last.IsZero() {
		self.last = now
		return true, 0
	}

	// Fill the bucket with the tokens in the time range [last, now]
	self.fillToken(now)

	if self.avail > 0 {
		self.avail -= 1
		return true, 0
	}

	waitDur := time.Duration(-time.Duration(self.avail) + self.fillInterval)
	if waitDur > time.Millisecond*time.Duration(maxWait) {
		return false, waitDur
	}
	time.Sleep(waitDur)
	self.avail -= 1
	return true, waitDur
}

func (self *TokenBucket) fillToken(now time.Time) {
	// Default, each interval add one token.
	fillTickNum := int(now.Sub(self.last) / self.fillInterval)
	if fillTickNum <= 0 {
		return
	}

	self.avail += fillTickNum
	self.last = now
	if self.avail > self.capacity {
		self.avail = self.capacity
	}
}

func init() {
	Register("token", NewTokenBucket)
}
