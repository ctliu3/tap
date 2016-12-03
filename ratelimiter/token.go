package ratelimiter

import (
	"sync"
	"time"
)

type TokenLimiter struct {
	sync.Mutex
	capacity     int
	fillInterval time.Duration

	// avail and last variables change during process
	// used for local store version
	avail int // available token number
	last  time.Time
}

// rate is QPS
func NewTokenLimiter(opt LimiterOption) (Limiter, error) {
	rate := opt.Rate
	capacity := opt.Capacity

	return &TokenLimiter{
		capacity:     capacity,
		fillInterval: time.Second / time.Duration(rate),
	}, nil
}

func (self *TokenLimiter) Name() string {
	return "TokenLimiter"
}

// maxWait: ms
func (self *TokenLimiter) Acquire(maxWait int) (bool, time.Duration) {
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

func (self *TokenLimiter) fillToken(now time.Time) {
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
	Register("token", NewTokenLimiter)
}
